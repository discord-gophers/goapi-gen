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

package codegen

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"unicode"

	"github.com/kenshaw/snaker"

	"github.com/getkin/kin-openapi/openapi3"
)

// ParameterDefinition is a definition for a parameter.
type ParameterDefinition struct {
	ParamName string // The original json parameter name, eg param_name
	In        string // Where the parameter is defined - path, header, cookie, query
	Required  bool   // Is this a required parameter?
	Spec      *openapi3.Parameter
	Schema    Schema
}

// TypeDef returns the type definition for a parameter without the leading '*'
// for optionl types.
func (pd ParameterDefinition) TypeDef() string {
	typeDecl := pd.Schema.TypeDecl()
	return typeDecl
}

// JSONTag returns the tag for pd. It includes omitempty if it is
// required.
func (pd *ParameterDefinition) JSONTag() string {
	if pd.Required {
		return fmt.Sprintf("`json:\"%s\"`", pd.ParamName)
	}
	return fmt.Sprintf("`json:\"%s,omitempty\"`", pd.ParamName)
}

// IsJSON returns if pd.Spec.Content contains application/json.
func (pd *ParameterDefinition) IsJSON() bool {
	p := pd.Spec
	if len(p.Content) == 1 {
		_, ok := p.Content["application/json"]
		return ok
	}
	return false
}

// IsPassThrough returns if pd is passthrough.
func (pd *ParameterDefinition) IsPassThrough() bool {
	p := pd.Spec
	if len(p.Content) > 1 {
		return true
	}
	if len(p.Content) == 1 {
		return !pd.IsJSON()
	}
	return false
}

// IsStyled returns if pdf.Spec.Schema is not nil.
func (pd *ParameterDefinition) IsStyled() bool {
	p := pd.Spec
	return p.Schema != nil
}

// Style returns pd.Spec.Style. If not provided, it determines the Style based
// on pd.Spec.In.
func (pd *ParameterDefinition) Style() string {
	style := pd.Spec.Style
	if style == "" {
		in := pd.Spec.In
		switch in {
		case "path", "header":
			return "simple"
		case "query", "cookie":
			return "form"
		default:
			panic("unknown parameter format")
		}
	}
	return style
}

// Explode returns if pd.Spec is explode.
// If Explode is nil, it determines Explode based on pd.Spec.In.
func (pd *ParameterDefinition) Explode() bool {
	if pd.Spec.Explode == nil {
		in := pd.Spec.In
		switch in {
		case "path", "header":
			return false
		case "query", "cookie":
			return true
		default:
			panic("unknown parameter format")
		}
	}
	return *pd.Spec.Explode
}

// GoVariableName returns a safe version of the name of pd's GoName.
func (pd ParameterDefinition) GoVariableName() string {
	name := snaker.ForceLowerCamelIdentifier(pd.GoName())
	if IsGoKeyword(name) {
		name = "p" + UppercaseFirstCharacter(name)
	}
	if unicode.IsNumber([]rune(name)[0]) {
		name = "n" + name
	}
	return name
}

// GoName returns the Go name of pd.ParamName.
func (pd ParameterDefinition) GoName() string {
	return SchemaNameToTypeName(pd.ParamName)
}

// IndirectOptional returns if pd is optiona, directlry or indirectly.
func (pd ParameterDefinition) IndirectOptional() bool {
	return !pd.Required && !pd.Schema.SkipOptionalPointer
}

// ParameterDefinitions is a slice of ParameterDefinition.
type ParameterDefinitions []ParameterDefinition

// FindByName finds a definition by name.
func (p ParameterDefinitions) FindByName(name string) *ParameterDefinition {
	for _, param := range p {
		if param.ParamName == name {
			return &param
		}
	}
	return nil
}

// DescribeParameters generates descriptors based on params and path.
// This makes it a lot easier to traverse the data in the template engine.
func DescribeParameters(params openapi3.Parameters, path []string) ([]ParameterDefinition, error) {
	outParams := make([]ParameterDefinition, 0)
	for _, paramOrRef := range params {
		param := paramOrRef.Value

		goType, err := paramToGoType(param, append(path, param.Name))
		if err != nil {
			return nil, fmt.Errorf("error generating type for param (%s): %s",
				param.Name, err)
		}

		pd := ParameterDefinition{
			ParamName: param.Name,
			In:        param.In,
			Required:  param.Required,
			Spec:      param,
			Schema:    goType,
		}

		// If this is a reference to a predefined type, simply use the reference
		// name as the type. $ref: "#/components/schemas/custom_type" becomes
		// "CustomType".
		if IsGoTypeReference(paramOrRef.Ref) {
			goType, err := RefPathToGoType(paramOrRef.Ref)
			if err != nil {
				return nil, fmt.Errorf("error dereferencing (%s) for param (%s): %s",
					paramOrRef.Ref, param.Name, err)
			}
			pd.Schema.GoType = goType
		}
		outParams = append(outParams, pd)
	}
	return outParams, nil
}

// SecurityDefinition represents an OpenAPI security definition.
type SecurityDefinition struct {
	ProviderName string
	Scopes       []string
}

// DescribeSecurityDefinition returns all security definitions in srs.
func DescribeSecurityDefinition(srs openapi3.SecurityRequirements) []SecurityDefinition {
	outDefs := make([]SecurityDefinition, 0)

	for _, sr := range srs {
		for _, k := range SortedSecurityRequirementKeys(sr) {
			v := sr[k]
			outDefs = append(outDefs, SecurityDefinition{ProviderName: k, Scopes: v})
		}
	}

	return outDefs
}

// OperationDefinition represents an Operation.
type OperationDefinition struct {
	OperationID string // The operation_id description from Swagger, used to generate function names

	PathParams          []ParameterDefinition // Parameters in the path, eg, /path/:param
	HeaderParams        []ParameterDefinition // Parameters in HTTP headers
	QueryParams         []ParameterDefinition // Parameters in the query, /path?param
	CookieParams        []ParameterDefinition // Parameters in cookies
	TypeDefinitions     []TypeDefinition      // These are all the types we need to define for this operation
	SecurityDefinitions []SecurityDefinition  // These are the security providers
	BodyRequired        bool
	Bodies              []RequestBodyDefinition // The list of bodies for which to generate handlers.
	Summary             string                  // Summary string from Swagger, used to generate a comment
	Method              string                  // GET, POST, DELETE, etc.
	Path                string                  // The Swagger path for the operation, like /resource/{id}
	Middlewares         []string                // Sent as part of x-go-middlewares.
	Spec                *openapi3.Operation
}

// Params returns the list of all parameters except Path parameters.
// Path parameters are handled differently from the rest, since they're
// mandatory.
func (o *OperationDefinition) Params() []ParameterDefinition {
	result := append(o.QueryParams, o.HeaderParams...)
	result = append(result, o.CookieParams...)
	return result
}

// AllParams returns all params.
func (o *OperationDefinition) AllParams() []ParameterDefinition {
	result := append(o.QueryParams, o.HeaderParams...)
	result = append(result, o.CookieParams...)
	result = append(result, o.PathParams...)
	return result
}

// RequiresParamObject returns if any path parameters exist.
// If we have parameters other than path parameters, they're bundled into an
// object.
func (o *OperationDefinition) RequiresParamObject() bool {
	return len(o.Params()) > 0
}

// SummaryAsComment returns the summary as a multiline comment for o.
func (o *OperationDefinition) SummaryAsComment() string {
	if o.Summary == "" {
		return ""
	}
	trimmed := strings.TrimSuffix(o.Summary, "\n")
	parts := strings.Split(trimmed, "\n")
	for i, p := range parts {
		parts[i] = "// " + p
	}
	return strings.Join(parts, "\n")
}

// GetResponseTypeDefinitions produces a list of type definitions for a given
// Operation for the response types which we know how to parse. These will be
// turned into fields on a response object for automatic deserialization of
// responses.
func (o *OperationDefinition) GetResponseTypeDefinitions() ([]ResponseTypeDefinition, error) {
	var tds []ResponseTypeDefinition

	responses := o.Spec.Responses
	sortedResponsesKeys := SortedResponsesKeys(responses)
	for _, responseName := range sortedResponsesKeys {
		responseRef := responses[responseName]

		// We can only generate a type if we have a value:
		if responseRef.Value != nil {
			sortedContentKeys := SortedContentKeys(responseRef.Value.Content)
			for _, contentTypeName := range sortedContentKeys {
				contentType := responseRef.Value.Content[contentTypeName]
				// We can only generate a type if we have a schema:
				if contentType.Schema != nil {
					responseSchema, err := GenerateGoSchema(contentType.Schema, []string{responseName})
					if err != nil {
						return nil, fmt.Errorf("unable to determine Go type for %s.%s: %w", o.OperationID, contentTypeName, err)
					}

					var typeName string
					switch {
					case StringInArray(contentTypeName, contentTypesJSON):
						typeName = fmt.Sprintf("JSON%s", ToCamelCase(responseName))
					// YAML:
					case StringInArray(contentTypeName, contentTypesYAML):
						typeName = fmt.Sprintf("YAML%s", ToCamelCase(responseName))
					// XML:
					case StringInArray(contentTypeName, contentTypesXML):
						typeName = fmt.Sprintf("XML%s", ToCamelCase(responseName))
					default:
						continue
					}

					td := ResponseTypeDefinition{
						TypeDefinition: TypeDefinition{
							TypeName: typeName,
							Schema:   responseSchema,
						},
						ResponseName:    responseName,
						ContentTypeName: contentTypeName,
					}
					if IsGoTypeReference(contentType.Schema.Ref) {
						refType, err := RefPathToGoType(contentType.Schema.Ref)
						if err != nil {
							return nil, fmt.Errorf("error dereferencing response Ref: %w", err)
						}
						td.Schema.RefType = refType
					}
					tds = append(tds, td)
				}
			}
		}
	}
	return tds, nil
}

// RequestBodyDefinition describes a request body
type RequestBodyDefinition struct {
	Required bool
	Schema   Schema

	// When we generate type names, we need a Tag for it, such as JSON, in
	// which case we will produce "JSONBody".
	NameTag string

	// This is the content type corresponding to the body, eg, application/json
	ContentType string

	// Whether this is the default body type. For an operation named OpFoo, we
	// will not add suffixes like OpFooJSONBody for this one.
	Default bool
}

// TypeDef returns the Go type definition for a request body
func (r RequestBodyDefinition) TypeDef(opID string) *TypeDefinition {
	return &TypeDefinition{
		TypeName: fmt.Sprintf("%s%sRequestBody", opID, r.NameTag),
		Schema:   r.Schema,
	}
}

// CustomType returns if body is a custom inline type, or pre-defined.
// TODO: clean up the templates code, it can be simpler.
func (r RequestBodyDefinition) CustomType() bool {
	return r.Schema.RefType == ""
}

// Suffix returns "With{r.nameTag}Body"
// Operation DoFoo would be suffixed with DoFooWithXMLBody.
func (r RequestBodyDefinition) Suffix() string {
	// The default response is never suffixed.
	if r.Default {
		return ""
	}
	return "With" + r.NameTag + "Body"
}

// FilterParameterDefinitionByType returns params which match the the type with
// in.
func FilterParameterDefinitionByType(params []ParameterDefinition, in string) []ParameterDefinition {
	var out []ParameterDefinition
	for _, p := range params {
		if p.In == in {
			out = append(out, p)
		}
	}
	return out
}

// OperationDefinitions returns all operations for a swagger definition.
func OperationDefinitions(swagger *openapi3.T) ([]OperationDefinition, error) {
	var operations []OperationDefinition

	for _, requestPath := range SortedPathsKeys(swagger.Paths) {
		pathItem := swagger.Paths[requestPath]
		// These are parameters defined for all methods on a given path. They
		// are shared by all methods.
		globalParams, err := DescribeParameters(pathItem.Parameters, nil)
		if err != nil {
			return nil, fmt.Errorf("error describing global parameters for %s: %s",
				requestPath, err)
		}

		var pathMiddlewares []string
		if extension, ok := pathItem.Extensions[extMiddlewares]; ok {
			var err error
			pathMiddlewares, err = extParseMiddlewares(extension)
			if err != nil {
				return nil, fmt.Errorf("invalid value for %q: %w", extMiddlewares, err)
			}
		}

		// Each path can have a number of operations, POST, GET, OPTIONS, etc.
		pathOps := pathItem.Operations()
		for _, opName := range SortedOperationsKeys(pathOps) {
			op := pathOps[opName]
			if pathItem.Servers != nil {
				op.Servers = &pathItem.Servers
			}

			// We rely on OperationID to generate function names, it's required
			op.OperationID = ToCamelCase(op.OperationID)
			if op.OperationID == "" {
				op.OperationID, err = generateDefaultOperationID(opName, requestPath)
				if err != nil {
					return nil, fmt.Errorf("error generating default OperationID for %s/%s: %s",
						opName, requestPath, err)
				}
			}

			// These are parameters defined for the specific path method that
			// we're iterating over.
			localParams, err := DescribeParameters(op.Parameters, []string{op.OperationID + "Params"})
			if err != nil {
				return nil, fmt.Errorf("error describing global parameters for %s/%s: %s",
					opName, requestPath, err)
			}
			// All the parameters required by a handler are the union of the
			// global parameters and the local parameters.
			allParams := append(globalParams, localParams...)

			// Order the path parameters to match the order as specified in
			// the path, not in the swagger spec, and validate that the parameter
			// names match, as downstream code depends on that.
			pathParams := FilterParameterDefinitionByType(allParams, "path")
			pathParams, err = SortParamsByPath(requestPath, pathParams)
			if err != nil {
				return nil, err
			}

			middlewares := pathMiddlewares
			if extension, ok := op.Extensions[extMiddlewares]; ok {
				opMiddlewares, err := extParseMiddlewares(extension)
				if err != nil {
					return nil, fmt.Errorf("invalid value for %q: %w", extMiddlewares, err)
				}
				middlewares = append(middlewares, opMiddlewares...)
			}

			bodyDefinitions, typeDefinitions, err := GenerateBodyDefinitions(op.OperationID, op.RequestBody)
			if err != nil {
				return nil, fmt.Errorf("error generating body definitions: %w", err)
			}

			opDef := OperationDefinition{
				PathParams:   pathParams,
				HeaderParams: FilterParameterDefinitionByType(allParams, "header"),
				QueryParams:  FilterParameterDefinitionByType(allParams, "query"),
				CookieParams: FilterParameterDefinitionByType(allParams, "cookie"),
				OperationID:  ToCamelCase(op.OperationID),
				// Replace newlines in summary.
				Summary:         op.Summary,
				Method:          opName,
				Path:            requestPath,
				Spec:            op,
				Bodies:          bodyDefinitions,
				TypeDefinitions: typeDefinitions,
				Middlewares:     middlewares,
			}

			// check for overrides of SecurityDefinitions.
			// See: "Step 2. Applying security:" from the spec:
			// https://swagger.io/docs/specification/authentication/
			if op.Security != nil {
				opDef.SecurityDefinitions = DescribeSecurityDefinition(*op.Security)
			} else {
				// use global securityDefinitions
				// globalSecurityDefinitions contains the top-level securityDefinitions.
				// They are the default securityPermissions which are injected into each
				// path, except for the case where a path explicitly overrides them.
				opDef.SecurityDefinitions = DescribeSecurityDefinition(swagger.Security)
			}

			if op.RequestBody != nil {
				opDef.BodyRequired = op.RequestBody.Value.Required
			}

			// Generate all the type definitions needed for this operation
			opDef.TypeDefinitions = append(opDef.TypeDefinitions, GenerateTypeDefsForOperation(opDef)...)

			operations = append(operations, opDef)
		}
	}
	return operations, nil
}

func generateDefaultOperationID(opName string, requestPath string) (string, error) {
	operationID := strings.ToLower(opName)

	if opName == "" {
		return "", fmt.Errorf("operation name cannot be an empty string")
	}

	if requestPath == "" {
		return "", fmt.Errorf("request path cannot be an empty string")
	}

	for _, part := range strings.Split(requestPath, "/") {
		if part != "" {
			operationID = operationID + "-" + part
		}
	}

	return ToCamelCase(operationID), nil
}

// GenerateBodyDefinitions returns  the Swagger body definitions into a list of
// our body definitions which will be used for code generation.
func GenerateBodyDefinitions(operationID string, bodyOrRef *openapi3.RequestBodyRef) ([]RequestBodyDefinition, []TypeDefinition, error) {
	if bodyOrRef == nil {
		return nil, nil, nil
	}
	body := bodyOrRef.Value

	var bodyDefinitions []RequestBodyDefinition
	var typeDefinitions []TypeDefinition

	for contentType, content := range body.Content {
		var tag string
		var defaultBody bool

		switch contentType {
		case "application/json":
			tag = "JSON"
			defaultBody = true
		default:
			continue
		}

		bodyTypeName := operationID + tag + "Body"
		bodySchema, err := GenerateGoSchema(content.Schema, []string{bodyTypeName})
		if err != nil {
			return nil, nil, fmt.Errorf("error generating request body definition: %w", err)
		}

		// If the body is a pre-defined type
		if IsGoTypeReference(bodyOrRef.Ref) {
			// Convert the reference path to Go type
			refType, err := RefPathToGoType(bodyOrRef.Ref)
			if err != nil {
				return nil, nil, fmt.Errorf("error turning reference (%s) into a Go type: %w", bodyOrRef.Ref, err)
			}
			bodySchema.RefType = refType
		}

		// If the request has a body, but it's not a user defined
		// type under #/components, we'll define a type for it, so
		// that we have an easy to use type for marshaling.
		if bodySchema.RefType == "" {
			td := TypeDefinition{
				TypeName: bodyTypeName,
				Schema:   bodySchema,
			}
			typeDefinitions = append(typeDefinitions, td)
			// The body schema now is a reference to a type
			bodySchema.RefType = bodyTypeName
		}

		bd := RequestBodyDefinition{
			Required:    body.Required,
			Schema:      bodySchema,
			NameTag:     tag,
			ContentType: contentType,
			Default:     defaultBody,
		}
		bodyDefinitions = append(bodyDefinitions, bd)
	}
	return bodyDefinitions, typeDefinitions, nil
}

// GenerateTypeDefsForOperation returns the type definitions for op.
func GenerateTypeDefsForOperation(op OperationDefinition) []TypeDefinition {
	var typeDefs []TypeDefinition
	// Start with the params object itself
	if len(op.Params()) != 0 {
		typeDefs = append(typeDefs, GenerateParamsTypes(op)...)
	}

	// Now, go through all the additional types we need to declare.
	for _, param := range op.AllParams() {
		typeDefs = append(typeDefs, param.Schema.AdditionalTypeDefs()...)
	}

	for _, body := range op.Bodies {
		typeDefs = append(typeDefs, body.Schema.AdditionalTypeDefs()...)
	}
	return typeDefs
}

// GenerateParamsTypes defines the schema for a parameters definition object
// which encapsulates all the query, header and cookie parameters for an
// operation.
func GenerateParamsTypes(op OperationDefinition) []TypeDefinition {
	var typeDefs []TypeDefinition

	objectParams := op.QueryParams
	objectParams = append(objectParams, op.HeaderParams...)
	objectParams = append(objectParams, op.CookieParams...)

	typeName := op.OperationID + "Params"

	s := Schema{}
	for _, param := range objectParams {
		pSchema := param.Schema
		if pSchema.HasAdditionalProperties {
			propRefName := strings.Join([]string{typeName, param.GoName()}, "_")
			pSchema.RefType = propRefName
			typeDefs = append(typeDefs, TypeDefinition{
				TypeName: propRefName,
				Schema:   param.Schema,
			})
		}
		prop := Property{
			Description:    param.Spec.Description,
			JSONFieldName:  param.ParamName,
			Required:       param.Required,
			Schema:         pSchema,
			ExtensionProps: &param.Spec.ExtensionProps,
		}
		s.Properties = append(s.Properties, prop)
	}

	s.Description = op.Spec.Description
	s.GoType = GenStructFromSchema(s)

	td := TypeDefinition{
		TypeName: typeName,
		Schema:   s,
	}
	return append(typeDefs, td)
}

// GenerateTypesForOperations prooduces code all types used by ops.
func GenerateTypesForOperations(t *template.Template, ops []OperationDefinition) (string, error) {
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)

	addTypes, err := GenerateTemplates([]string{"param-types.tmpl", "request-bodies.tmpl", "response-bodies.tmpl"}, t, ops)
	if err != nil {
		return "", fmt.Errorf("error generating type boilerplate for operations: %w", err)
	}
	if _, err := w.WriteString(addTypes); err != nil {
		return "", fmt.Errorf("error writing boilerplate to buffer: %w", err)
	}

	// Generate boiler plate for all additional types.
	var td []TypeDefinition
	for _, op := range ops {
		td = append(td, op.TypeDefinitions...)
	}

	addProps, err := GenerateAdditionalPropertyBoilerplate(t, td)
	if err != nil {
		return "", fmt.Errorf("error generating additional properties boilerplate for operations: %w", err)
	}

	if _, err := w.WriteString("\n"); err != nil {
		return "", fmt.Errorf("error generating additional properties boilerplate for operations: %w", err)
	}

	if _, err := w.WriteString(addProps); err != nil {
		return "", fmt.Errorf("error generating additional properties boilerplate for operations: %w", err)
	}

	if err = w.Flush(); err != nil {
		return "", fmt.Errorf("error flushing output buffer for server interface: %w", err)
	}

	return buf.String(), nil
}

// GenerateChiServer generates codee for the chi server for ops.
func GenerateChiServer(t *template.Template, operations []OperationDefinition) (string, error) {
	return GenerateTemplates([]string{"interface.tmpl", "middleware.tmpl", "handler.tmpl"}, t, operations)
}

// GenerateTemplates generates templates
func GenerateTemplates(templates []string, t *template.Template, ops interface{}) (string, error) {
	var generatedTemplates []string
	for _, tmpl := range templates {
		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)

		if err := t.ExecuteTemplate(w, tmpl, ops); err != nil {
			return "", fmt.Errorf("error generating %s: %s", tmpl, err)
		}
		if err := w.Flush(); err != nil {
			return "", fmt.Errorf("error flushing output buffer for %s: %s", tmpl, err)
		}
		generatedTemplates = append(generatedTemplates, buf.String())
	}

	return strings.Join(generatedTemplates, "\n"), nil
}
