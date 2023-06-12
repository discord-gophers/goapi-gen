// Copyright 2019 DeepMap, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package codegen provides functions to generate Go code from OpenAPI
// specifications.
package codegen

import (
	"bufio"
	"bytes"
	"fmt"
	"runtime/debug"
	"sort"
	"strings"
	"text/template"

	"github.com/getkin/kin-openapi/openapi3"
	"golang.org/x/tools/imports"

	"github.com/discord-gophers/goapi-gen/templates"
)

// Options defines the optional code to generate.
//
// Most callers to this package will use Generate.
type Options struct {
	GenerateServer bool              // GenerateChiServer specifies whether to generate chi server boilerplate
	GenerateTypes  bool              // GenerateTypes specifies whether to generate type definitions
	EmbedSpec      bool              // Whether to embed the swagger spec in the generated code
	SkipFmt        bool              // Whether to skip go imports on the generated code
	SkipPrune      bool              // Whether to skip pruning unused components on the generated code
	AliasTypes     bool              // Whether to alias types if possible
	IncludeTags    []string          // Only include operations that have one of these tags. Ignored when empty.
	ExcludeTags    []string          // Exclude operations that have one of these tags. Ignored when empty.
	UserTemplates  map[string]string // Override built-in templates from user-provided files
	ImportMapping  map[string]string // ImportMapping specifies the golang package path for each external reference
	ExcludeSchemas []string          // Exclude from generation schemas with given names. Ignored when empty.
}

// goImport represents a go package to be imported in the generated code
type goImport struct {
	Name string // package name
	Path string // package path
}

// String returns a go import statement
func (gi goImport) String() string {
	if gi.Name != "" {
		return fmt.Sprintf("%s %q", gi.Name, gi.Path)
	}
	return fmt.Sprintf("%q", gi.Path)
}

// importMap maps external OpenAPI specifications files/urls to external go packages
type importMap map[string]goImport

// GoImports returns a slice of go import statements
func (im importMap) GoImports() []string {
	goImports := make([]string, 0, len(im))
	for _, v := range im {
		goImports = append(goImports, v.String())
	}
	return goImports
}

var importMapping importMap

func constructImportMapping(input map[string]string) importMap {
	var (
		pathToName = map[string]string{}
		result     = importMap{}
	)

	{
		var packagePaths []string
		for _, packageName := range input {
			packagePaths = append(packagePaths, packageName)
		}
		sort.Strings(packagePaths)

		for _, packagePath := range packagePaths {
			if _, ok := pathToName[packagePath]; !ok {
				pathToName[packagePath] = fmt.Sprintf("externalRef%d", len(pathToName))
			}
		}
	}
	for specPath, packagePath := range input {
		result[specPath] = goImport{Name: pathToName[packagePath], Path: packagePath}
	}
	return result
}

// Generate uses the Go templating engine to generate all of our server wrappers from
// the descriptions we've built up above from the schema objects.
func Generate(swagger *openapi3.T, packageName string, opts Options) (string, error) {
	importMapping = constructImportMapping(opts.ImportMapping)

	filterOperationsByTag(swagger, opts)
	if !opts.SkipPrune {
		pruneUnusedComponents(swagger)
	}

	// This creates the golang templates text package
	TemplateFunctions["opts"] = func() Options { return opts }
	t := template.New("goapi-gen").Funcs(TemplateFunctions)
	// This parses all of our own template files into the template object
	// above
	t, err := templates.Parse(t)
	if err != nil {
		return "", fmt.Errorf("error parsing goapi-gen templates: %w", err)
	}

	// Override built-in templates with user-provided versions
	for _, tpl := range t.Templates() {
		if _, ok := opts.UserTemplates[tpl.Name()]; ok {
			utpl := t.New(tpl.Name())
			if _, err := utpl.Parse(opts.UserTemplates[tpl.Name()]); err != nil {
				return "", fmt.Errorf("error parsing user-provided template %q: %w", tpl.Name(), err)
			}
		}
	}

	var finalCustomImports []string
	ops, err := OperationDefinitions(swagger)
	if err != nil {
		return "", fmt.Errorf("error creating operation definitions: %w", err)
	}

	for _, op := range ops {
		finalCustomImports = append(finalCustomImports, op.CustomImports...)
	}

	var typeDefinitions, constantDefinitions string
	var customImports []string
	if opts.GenerateTypes {
		typeDefinitions, customImports, err = GenerateTypeDefinitions(t, swagger, ops, opts.ExcludeSchemas)
		if err != nil {
			return "", fmt.Errorf("error generating type definitions: %w", err)
		}

		constantDefinitions, err = GenerateConstants(t, ops)
		if err != nil {
			return "", fmt.Errorf("error generating constants: %w", err)
		}

	}

	// TODO: check for exact double imports and merge them together with 1 alias, otherwise we might run into double imports under different names
	finalCustomImports = append(finalCustomImports, customImports...)

	var serverOut string
	if opts.GenerateServer {
		serverOut, err = GenerateChiServer(t, ops)
		if err != nil {
			return "", fmt.Errorf("error generating Go handlers for Paths: %w", err)
		}
	}

	var inlinedSpec string
	if opts.EmbedSpec {
		inlinedSpec, err = GenerateInlinedSpec(t, importMapping, swagger)
		if err != nil {
			return "", fmt.Errorf("error generating Go handlers for Paths: %w", err)
		}
	}

	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)

	externalImports := importMapping.GoImports()
	externalImports = append(externalImports, finalCustomImports...)
	importsOut, err := GenerateImports(t, externalImports, packageName)
	if err != nil {
		return "", fmt.Errorf("error generating imports: %w", err)
	}

	_, err = w.WriteString(importsOut)
	if err != nil {
		return "", fmt.Errorf("error writing imports: %w", err)
	}

	_, err = w.WriteString(constantDefinitions)
	if err != nil {
		return "", fmt.Errorf("error writing constants: %w", err)
	}

	_, err = w.WriteString(typeDefinitions)
	if err != nil {
		return "", fmt.Errorf("error writing type definitions: %w", err)
	}

	if opts.GenerateServer {
		_, err = w.WriteString(serverOut)
		if err != nil {
			return "", fmt.Errorf("error writing server path handlers: %w", err)
		}
	}

	if opts.EmbedSpec {
		_, err = w.WriteString(inlinedSpec)
		if err != nil {
			return "", fmt.Errorf("error writing inlined spec: %w", err)
		}
	}

	err = w.Flush()
	if err != nil {
		return "", fmt.Errorf("error flushing output buffer: %w", err)
	}

	// remove any byte-order-marks which break Go-Code
	goCode := SanitizeCode(buf.String())

	// The generation code produces unindented horrors. Use the Go Imports
	// to make it all pretty.
	if opts.SkipFmt {
		return goCode, nil
	}

	outBytes, err := imports.Process(packageName+".go", []byte(goCode), nil)
	if err != nil {
		return "", fmt.Errorf("error formatting Go code: %w", err)
	}
	return string(outBytes), nil
}

// GenerateTypeDefinitions produces the type definitions in ops and executes
// the template.
func GenerateTypeDefinitions(t *template.Template, swagger *openapi3.T, ops []OperationDefinition, excludeSchemas []string) (string, []string, error) {
	schemaTypes, err := GenerateTypesForSchemas(t, swagger.Components.Schemas, excludeSchemas)
	if err != nil {
		return "", nil, fmt.Errorf("error generating Go types for component schemas: %w", err)
	}

	paramTypes, err := GenerateTypesForParameters(t, swagger.Components.Parameters)
	if err != nil {
		return "", nil, fmt.Errorf("error generating Go types for component parameters: %w", err)
	}
	allTypes := append(schemaTypes, paramTypes...)

	responseTypes, err := GenerateTypesForResponses(t, swagger.Components.Responses)
	if err != nil {
		return "", nil, fmt.Errorf("error generating Go types for component responses: %w", err)
	}
	allTypes = append(allTypes, responseTypes...)

	bodyTypes, err := GenerateTypesForRequestBodies(t, swagger.Components.RequestBodies)
	if err != nil {
		return "", nil, fmt.Errorf("error generating Go types for component request bodies: %w", err)
	}
	allTypes = append(allTypes, bodyTypes...)

	paramTypesOut, err := GenerateTypesForOperations(t, ops)
	if err != nil {
		return "", nil, fmt.Errorf("error generating Go types for component request bodies: %w", err)
	}

	enumsOut, err := GenerateEnums(t, allTypes)
	if err != nil {
		return "", nil, fmt.Errorf("error generating code for type enums: %w", err)
	}

	enumTypesOut, err := GenerateEnumTypes(t, allTypes)
	if err != nil {
		return "", nil, fmt.Errorf("error generating code for enum type definitions: %w", err)
	}

	typesOut, err := GenerateTypes(t, allTypes)
	if err != nil {
		return "", nil, fmt.Errorf("error generating code for type definitions: %w", err)
	}

	allOfBoilerplate, err := GenerateAdditionalPropertyBoilerplate(t, allTypes)
	if err != nil {
		return "", nil, fmt.Errorf("error generating allOf boilerplate: %w", err)
	}

	var customImports []string
	for _, allType := range allTypes {
		customImports = append(customImports, allType.Schema.CustomImports...)
		for _, prop := range allType.Schema.Properties {
			customImports = append(customImports, prop.Schema.CustomImports...)
		}
	}

	typeDefinitions := enumsOut + typesOut + enumTypesOut + paramTypesOut + allOfBoilerplate
	return typeDefinitions, customImports, nil
}

// GenerateConstants creates operation ids, context keys, paths, etc. to be
// exported as constants
func GenerateConstants(t *template.Template, ops []OperationDefinition) (string, error) {
	constants := Constants{
		SecuritySchemeProviderNames: []string{},
		Middlewares:                 []string{},
	}

	providerNameMap := map[string]struct{}{}
	middlewaresMap := map[string]struct{}{}
	for _, op := range ops {
		for _, def := range op.SecurityDefinitions {
			providerName := SanitizeGoIdentity(def.ProviderName)
			providerNameMap[providerName] = struct{}{}
		}

		for _, middleware := range op.Middlewares {
			middlewaresMap[middleware] = struct{}{}
		}
	}

	for providerName := range providerNameMap {
		constants.SecuritySchemeProviderNames = append(constants.SecuritySchemeProviderNames, providerName)
	}

	for middleware := range middlewaresMap {
		constants.Middlewares = append(constants.Middlewares, middleware)
	}

	sort.Strings(constants.SecuritySchemeProviderNames)
	sort.Strings(constants.Middlewares)

	return GenerateTemplates([]string{"constants.tmpl"}, t, constants)
}

// GenerateTypesForSchemas creates type definitions for any custom types defined in the
// components/schemas section of the Swagger spec.
func GenerateTypesForSchemas(t *template.Template, schemas map[string]*openapi3.SchemaRef, excludeSchemas []string) ([]TypeDefinition, error) {
	excludeSchemasMap := make(map[string]bool)
	for _, schema := range excludeSchemas {
		excludeSchemasMap[schema] = true
	}
	types := make([]TypeDefinition, 0)
	// We're going to define Go types for every object under components/schemas
	for _, schemaName := range SortedSchemaKeys(schemas) {
		if _, ok := excludeSchemasMap[schemaName]; ok {
			continue
		}
		schemaRef := schemas[schemaName]

		goSchema, err := GenerateGoSchema(schemaRef, []string{schemaName})
		if err != nil {
			return nil, fmt.Errorf("error converting Schema %s to Go type: %w", schemaName, err)
		}

		types = append(types, TypeDefinition{
			JSONName: schemaName,
			TypeName: SchemaNameToTypeName(schemaName),
			Schema:   goSchema,
		})

		types = append(types, goSchema.AdditionalTypeDefs()...)
	}
	return types, nil
}

// GenerateTypesForParameters creates type definitions for any custom types defined in
// the components/parameters section of the Swagger spec.
func GenerateTypesForParameters(t *template.Template, params map[string]*openapi3.ParameterRef) ([]TypeDefinition, error) {
	var types []TypeDefinition
	for _, paramName := range SortedParameterKeys(params) {
		paramOrRef := params[paramName]

		goType, err := paramToGoType(paramOrRef.Value, nil)
		if err != nil {
			return nil, fmt.Errorf("error generating Go type for schema in parameter %s: %w", paramName, err)
		}

		typeDef := TypeDefinition{
			JSONName: paramName,
			Schema:   goType,
			TypeName: SchemaNameToTypeName(paramName),
		}

		if paramOrRef.Ref != "" {
			// Generate a reference type for referenced parameters
			refType, err := RefPathToGoType(paramOrRef.Ref)
			if err != nil {
				return nil, fmt.Errorf("error generating Go type for (%s) in parameter %s: %w", paramOrRef.Ref, paramName, err)
			}
			typeDef.TypeName = SchemaNameToTypeName(refType)
		}

		types = append(types, typeDef)
	}
	return types, nil
}

// GenerateTypesForResponses makes definitions for any custom types defined in
// the components/responses section of the Swagger spec.
func GenerateTypesForResponses(t *template.Template, responses openapi3.Responses) ([]TypeDefinition, error) {
	var types []TypeDefinition

	for _, responseName := range SortedResponsesKeys(responses) {
		responseOrRef := responses[responseName]

		// We have to generate the response object. We're only going to
		// handle application/json media types here. Other responses should
		// simply be specified as strings or byte arrays.
		response := responseOrRef.Value
		jsonResponse, found := response.Content["application/json"]
		if found {
			goType, err := GenerateGoSchema(jsonResponse.Schema, []string{responseName})
			if err != nil {
				return nil, fmt.Errorf("error generating Go type for schema in response %s: %w", responseName, err)
			}

			typeDef := TypeDefinition{
				JSONName: responseName,
				Schema:   goType,
				TypeName: SchemaNameToTypeName(responseName),
			}

			if responseOrRef.Ref != "" {
				// Generate a reference type for referenced parameters
				refType, err := RefPathToGoType(responseOrRef.Ref)
				if err != nil {
					return nil, fmt.Errorf("error generating Go type for (%s) in parameter %s: %w", responseOrRef.Ref, responseName, err)
				}
				typeDef.TypeName = SchemaNameToTypeName(refType)
			}
			types = append(types, typeDef)
		}
	}
	return types, nil
}

// GenerateTypesForRequestBodies creates definitions for any custom types
// defined in the components/requestBodies section of the Swagger spec.
func GenerateTypesForRequestBodies(t *template.Template, bodies map[string]*openapi3.RequestBodyRef) ([]TypeDefinition, error) {
	var types []TypeDefinition

	for _, bodyName := range SortedRequestBodyKeys(bodies) {
		bodyOrRef := bodies[bodyName]

		// As for responses, we will only generate Go code for JSON bodies,
		// the other body formats are up to the user.
		response := bodyOrRef.Value
		jsonBody, found := response.Content["application/json"]
		if found {
			goType, err := GenerateGoSchema(jsonBody.Schema, []string{bodyName})
			if err != nil {
				return nil, fmt.Errorf("error generating Go type for schema in body %s: %w", bodyName, err)
			}

			typeDef := TypeDefinition{
				JSONName: bodyName,
				Schema:   goType,
				TypeName: SchemaNameToTypeName(bodyName),
			}

			if bodyOrRef.Ref != "" {
				// Generate a reference type for referenced bodies
				refType, err := RefPathToGoType(bodyOrRef.Ref)
				if err != nil {
					return nil, fmt.Errorf("error generating Go type for (%s) in body %s: %w", bodyOrRef.Ref, bodyName, err)
				}
				typeDef.TypeName = SchemaNameToTypeName(refType)
			}
			types = append(types, typeDef)
		}
	}
	return types, nil
}

// GenerateTypes is a helper function to execute the templates for type defs.
func GenerateTypes(t *template.Template, types []TypeDefinition) (string, error) {
	m := map[string]bool{}
	ts := []TypeDefinition{}

	for _, t := range types {
		if found := m[t.TypeName]; found {
			continue
		}

		m[t.TypeName] = true

		if len(t.Schema.EnumValues) > 0 {
			continue
		}

		ts = append(ts, t)
	}

	context := struct {
		Types []TypeDefinition
	}{
		Types: ts,
	}

	return GenerateTemplates([]string{"typedef.tmpl"}, t, context)
}

// GenerateEnumTypes makes types for any enums in types.
func GenerateEnumTypes(t *template.Template, types []TypeDefinition) (string, error) {
	m := map[string]bool{}
	ts := []TypeDefinition{}

	for _, t := range types {
		if found := m[t.TypeName]; found {
			continue
		}

		m[t.TypeName] = true

		if len(t.Schema.EnumValues) <= 0 {
			continue
		}

		ts = append(ts, t)
	}

	context := struct {
		Types []TypeDefinition
	}{
		Types: ts,
	}

	return GenerateTemplates([]string{"enum-typedef.tmpl"}, t, context)
}

// GenerateEnums makes definitions for any enums defined in types.
func GenerateEnums(t *template.Template, types []TypeDefinition) (string, error) {
	c := Constants{
		EnumDefinitions: []EnumDefinition{},
	}

	m := map[string]bool{}

	for _, tp := range types {
		if found := m[tp.TypeName]; found {
			continue
		}

		m[tp.TypeName] = true

		if len(tp.Schema.EnumValues) > 0 {
			wrapper := ""
			if tp.Schema.GoType == "string" {
				wrapper = `"`
			}
			c.EnumDefinitions = append(c.EnumDefinitions, EnumDefinition{
				Schema:       tp.Schema,
				TypeName:     tp.TypeName,
				ValueWrapper: wrapper,
			})
		}
	}

	return GenerateTemplates([]string{"enum-values.tmpl"}, t, c)
}

// GenerateImports creates import statements and the package definition.
func GenerateImports(t *template.Template, externalImports []string, packageName string) (string, error) {
	// Read build version for incorporating into generated files
	var modulePath string
	var moduleVersion string
	if bi, ok := debug.ReadBuildInfo(); ok {
		modulePath = bi.Main.Path
		moduleVersion = bi.Main.Version
	} else {
		// Unit tests have ok=false, so we'll just use "unknown" for the
		// version if we can't read this.
		modulePath = "unknown module path"
		moduleVersion = "unknown version"
	}

	context := struct {
		ExternalImports []string
		PackageName     string
		ModuleName      string
		Version         string
	}{
		ExternalImports: externalImports,
		PackageName:     packageName,
		ModuleName:      modulePath,
		Version:         moduleVersion,
	}

	return GenerateTemplates([]string{"imports.tmpl"}, t, context)
}

// GenerateAdditionalPropertyBoilerplate creates any glue code for interfacing
// with additional properties and JSON marshaling.
func GenerateAdditionalPropertyBoilerplate(t *template.Template, typeDefs []TypeDefinition) (string, error) {
	var filteredTypes []TypeDefinition

	m := map[string]bool{}

	for _, t := range typeDefs {
		if found := m[t.TypeName]; found {
			continue
		}

		m[t.TypeName] = true

		if t.Schema.HasAdditionalProperties {
			filteredTypes = append(filteredTypes, t)
		}
	}

	context := struct {
		Types []TypeDefinition
	}{
		Types: filteredTypes,
	}

	return GenerateTemplates([]string{"additional-properties.tmpl"}, t, context)
}

// SanitizeCode runs sanitizers across the generated Go code to ensure the
// generated code will be able to compile.
func SanitizeCode(goCode string) string {
	// remove any byte-order-marks which break Go-Code
	// See: https://groups.google.com/forum/#!topic/golang-nuts/OToNIPdfkks
	return strings.Replace(goCode, "\uFEFF", "", -1)
}
