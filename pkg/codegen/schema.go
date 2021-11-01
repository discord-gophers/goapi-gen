package codegen

import (
	"errors"
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Schema represents an OpenAPI type definition.
type Schema struct {
	GoType  string // The Go type needed to represent the schema
	RefType string // If the type has a type name, this is set

	ArrayType *Schema // The schema of array element

	EnumValues map[string]string // Enum values

	Properties               []Property       // For an object, the fields with names
	HasAdditionalProperties  bool             // Whether we support additional properties
	AdditionalPropertiesType *Schema          // And if we do, their type
	AdditionalTypes          []TypeDefinition // We may need to generate auxiliary helper types, stored here

	SkipOptionalPointer bool // Some types don't need a * in front when they're optional
	Bindable            bool // Indicates whether this type can implement render.Binder

	Description string // The description of the element

	// The original OpenAPIv3 Schema.
	OAPISchema *openapi3.Schema
}

// IsRef returns if s references another type.
func (s Schema) IsRef() bool {
	return s.RefType != ""
}

// TypeDecl returns the Go type needed to represent the schema. If s has
// a reference, it returns the type of the reference.
func (s Schema) TypeDecl() string {
	if s.IsRef() {
		return s.RefType
	}
	return s.GoType
}

// MergeProperty adds p to the properties of s. p must not already be a property
// in the schema.
func (s *Schema) MergeProperty(p Property) error {
	// Scan all existing properties for a conflict
	for _, e := range s.Properties {
		if e.JSONFieldName == p.JSONFieldName && !PropertiesEqual(e, p) {
			return fmt.Errorf("property '%s' already exists with a different type", e.JSONFieldName)
		}
	}
	s.Properties = append(s.Properties, p)
	return nil
}

// AdditionalTypeDefs returns all type definitions of the properties of s,
// and then any additional types in s.
func (s Schema) AdditionalTypeDefs() []TypeDefinition {
	var result []TypeDefinition
	for _, p := range s.Properties {
		result = append(result, p.Schema.AdditionalTypeDefs()...)
	}
	result = append(result, s.AdditionalTypes...)
	return result
}

// Property represents an OpenAPI property.
type Property struct {
	Description    string
	JSONFieldName  string
	Schema         Schema
	Required       bool
	Nullable       bool
	ExtensionProps *openapi3.ExtensionProps
}

// GoFieldName returns the Go name of p.
func (p Property) GoFieldName() string {
	return SchemaNameToTypeName(p.JSONFieldName)
}

// GoTypeDef returns the go type of p.
func (p Property) GoTypeDef() string {
	typeDef := p.Schema.TypeDecl()
	if !p.Schema.SkipOptionalPointer && (!p.Required || p.Nullable) {
		typeDef = "*" + typeDef
	}
	return typeDef
}

// EnumDefinition holds type information for enum
type EnumDefinition struct {
	Schema       Schema
	TypeName     string
	ValueWrapper string
}

// Constants holds the list of definitions that will be defined as constants.
type Constants struct {
	// SecuritySchemeProviderNames holds all provider names for security schemes.
	SecuritySchemeProviderNames []string
	// EnumDefinitions holds type and value information for all enums
	EnumDefinitions []EnumDefinition
}

// TypeDefinition describes a Go type definition in generated code.
//
// Let's use this example schema:
// components:
//  schemas:
//    Person:
//      type: object
//      properties:
//      name:
//        type: string
type TypeDefinition struct {
	// The name of the type, eg, type <...> Person
	TypeName string

	// The name of the corresponding JSON description, as it will sometimes
	// differ due to invalid characters.
	JSONName string

	// This is the Schema wrapper is used to populate the type description
	Schema Schema
}

// ResponseTypeDefinition is an extension of TypeDefinition, specifically for
// response unmarshaling.
type ResponseTypeDefinition struct {
	TypeDefinition
	// The content type name where this is used, eg, application/json
	ContentTypeName string

	// The type name of a response model.
	ResponseName string
}

// CanAlias returns whether the name of the type can be aliased.
// It can be aliased if t's Schema is a reference or an array to a reference.
func (t *TypeDefinition) CanAlias() bool {
	return t.Schema.IsRef() || (t.Schema.ArrayType != nil && t.Schema.ArrayType.IsRef())
}

// PropertiesEqual returns if a and b can be considered to be the same.
// a and b are the same if they have the same field name, same type, and are
// both required (or not).
func PropertiesEqual(a, b Property) bool {
	return a.JSONFieldName == b.JSONFieldName && a.Schema.TypeDecl() == b.Schema.TypeDecl() && a.Required == b.Required
}

// GenerateGoSchema generates the schema for sref.
// If it cannot properly resolve the type of sref, it returns
// map[string]interface{} or interface{}.
func GenerateGoSchema(sref *openapi3.SchemaRef, path []string) (Schema, error) {
	// Add a fallback value in case the sref is nil.
	// i.e. the parent schema defines a type:array, but the array has
	// no items defined. Therefore we have at least valid Go-Code.
	if sref == nil {
		return Schema{GoType: "interface{}"}, nil
	}

	schema := sref.Value

	// If Ref is set on the SchemaRef, it means that this type is actually a reference to
	// another type. We're not de-referencing, so simply use the referenced type.
	if IsGoTypeReference(sref.Ref) {
		// Convert the reference path to Go type
		refType, err := RefPathToGoType(sref.Ref)
		if err != nil {
			return Schema{}, fmt.Errorf("error turning reference (%s) into a Go type: %s",
				sref.Ref, err)
		}
		return Schema{
			GoType:      refType,
			Description: StringToGoComment(schema.Description),
			Bindable:    true,
		}, nil
	}

	outSchema := Schema{
		Description: StringToGoComment(schema.Description),
		OAPISchema:  schema,
		Bindable:    true,
	}

	// FIXME(hhhapz): We can probably support this in a meaningful way.
	// We can't support this in any meaningful way
	if schema.AnyOf != nil || schema.OneOf != nil {
		outSchema.GoType = "interface{}"
		outSchema.Bindable = false
		return outSchema, nil
	}

	if schema.AllOf != nil {
		mergedSchema, err := MergeSchemas(schema.AllOf, path)
		if err != nil {
			return Schema{}, fmt.Errorf("error merging schemas: %w", err)
		}
		mergedSchema.OAPISchema = schema
		return mergedSchema, nil
	}

	// Check for custom Go type extension
	if extension, ok := schema.Extensions[extPropGoType]; ok {
		typeName, err := extTypeName(extension)
		if err != nil {
			return outSchema, fmt.Errorf("invalid value for %q: %w", extPropGoType, err)
		}
		outSchema.GoType = typeName
		return outSchema, nil
	}

	// Schema type and format, eg. string / binary
	t := schema.Type
	// Handle objects and empty schemas first as a special case
	if t == "" || t == "object" {
		var outType string

		if len(schema.Properties) == 0 && !SchemaHasAdditionalProperties(schema) {
			// If the object has no properties or additional properties, we
			// have some special cases for its type.
			if t == "object" {
				// We have an object with no properties. This is a generic object
				// expressed as a map.
				outType = "map[string]interface{}"
			} else { // t == ""
				// If we don't even have the object designator, we're a completely
				// generic type.
				outType = "interface{}"
				outSchema.Bindable = false
			}
			outSchema.GoType = outType
		} else {
			// We've got an object with some properties.
			for _, pName := range SortedSchemaKeys(schema.Properties) {
				p := schema.Properties[pName]
				propertyPath := append(path, pName)
				pSchema, err := GenerateGoSchema(p, propertyPath)
				if err != nil {
					return Schema{}, fmt.Errorf("error generating Go schema for property '%s': %w", pName, err)
				}

				required := StringInArray(pName, schema.Required)

				if pSchema.HasAdditionalProperties && pSchema.RefType == "" {
					// If we have fields present which have additional properties,
					// but are not a pre-defined type, we need to define a type
					// for them, which will be based on the field names we followed
					// to get to the type.
					typeName := PathToTypeName(propertyPath)

					typeDef := TypeDefinition{
						TypeName: typeName,
						JSONName: strings.Join(propertyPath, "."),
						Schema:   pSchema,
					}
					pSchema.AdditionalTypes = append(pSchema.AdditionalTypes, typeDef)

					pSchema.RefType = typeName
				}
				description := ""
				if p.Value != nil {
					description = p.Value.Description
				}
				prop := Property{
					JSONFieldName:  pName,
					Schema:         pSchema,
					Required:       required,
					Description:    description,
					Nullable:       p.Value.Nullable,
					ExtensionProps: &p.Value.ExtensionProps,
				}
				outSchema.Properties = append(outSchema.Properties, prop)
			}

			outSchema.HasAdditionalProperties = SchemaHasAdditionalProperties(schema)
			outSchema.AdditionalPropertiesType = &Schema{
				GoType: "interface{}",
			}
			if schema.AdditionalProperties != nil {
				additionalSchema, err := GenerateGoSchema(schema.AdditionalProperties, path)
				if err != nil {
					return Schema{}, fmt.Errorf("error generating type for additional properties: %w", err)
				}
				outSchema.AdditionalPropertiesType = &additionalSchema
			}

			outSchema.GoType = GenStructFromSchema(outSchema)
		}
		return outSchema, nil
	} else if len(schema.Enum) > 0 {
		err := resolveType(schema, path, &outSchema)
		if err != nil {
			return Schema{}, fmt.Errorf("error resolving primitive type: %w", err)
		}
		enumValues := make([]string, len(schema.Enum))
		for i, enumValue := range schema.Enum {
			enumValues[i] = fmt.Sprintf("%v", enumValue)
		}

		sanitizedValues := SanitizeEnumNames(enumValues)
		outSchema.EnumValues = make(map[string]string, len(sanitizedValues))
		var constNamePath []string
		for k, v := range sanitizedValues {
			if v == "" {
				constNamePath = append(path, "Empty")
			} else {
				constNamePath = append(path, k)
			}
			outSchema.EnumValues[SchemaNameToTypeName(PathToTypeName(constNamePath))] = v
		}
		if len(path) > 1 { // handle additional type only on non-toplevel types
			typeName := SchemaNameToTypeName(PathToTypeName(path))
			typeDef := TypeDefinition{
				TypeName: typeName,
				JSONName: strings.Join(path, "."),
				Schema:   outSchema,
			}
			outSchema.AdditionalTypes = append(outSchema.AdditionalTypes, typeDef)
			outSchema.RefType = typeName
		}
		// outSchema.RefType = typeName
	} else {
		err := resolveType(schema, path, &outSchema)
		if err != nil {
			return Schema{}, fmt.Errorf("error resolving primitive type")
		}
	}
	return outSchema, nil
}

// resolveType resolves primitive  type or array for schema
func resolveType(schema *openapi3.Schema, path []string, outSchema *Schema) error {
	f := schema.Format
	t := schema.Type

	switch t {
	case "array":
		// For arrays, we'll get the type of the Items and throw a
		// [] in front of it.
		arrayType, err := GenerateGoSchema(schema.Items, path)
		if err != nil {
			return fmt.Errorf("error generating type for array: %w", err)
		}
		outSchema.ArrayType = &arrayType
		outSchema.GoType = "[]" + arrayType.TypeDecl()
		outSchema.AdditionalTypes = arrayType.AdditionalTypes
		outSchema.Properties = arrayType.Properties
		outSchema.SkipOptionalPointer = true
	case "integer":
		// We default to int if format doesn't ask for something else.
		switch f {
		case "int64":
			outSchema.GoType = "int64"
		case "int32":
			outSchema.GoType = "int32"
		case "int16":
			outSchema.GoType = "int16"
		case "int8":
			outSchema.GoType = "int8"
		case "int":
			outSchema.GoType = "int"
		case "uint64":
			outSchema.GoType = "uint64"
		case "uint32":
			outSchema.GoType = "uint32"
		case "uint16":
			outSchema.GoType = "uint16"
		case "uint8":
			outSchema.GoType = "uint8"
		case "uint":
			outSchema.GoType = "uint"
		case "":
			outSchema.GoType = "int"
		default:
			return fmt.Errorf("invalid integer format: %s", f)
		}
	case "number":
		// We default to float for "number"
		switch f {
		case "double":
			outSchema.GoType = "float64"
		case "float", "":
			outSchema.GoType = "float32"
		default:
			return fmt.Errorf("invalid number format: %s", f)
		}
	case "boolean":
		if f != "" {
			return fmt.Errorf("invalid format (%s) for boolean", f)
		}
		outSchema.GoType = "bool"
	case "string":
		// Special case string formats here.
		switch f {
		case "byte":
			outSchema.GoType = "[]byte"
			outSchema.SkipOptionalPointer = true
		case "email":
			outSchema.GoType = "openapi_types.Email"
		case "date":
			outSchema.GoType = "openapi_types.Date"
		case "date-time":
			outSchema.GoType = "time.Time"
		case "json":
			outSchema.GoType = "json.RawMessage"
			outSchema.SkipOptionalPointer = true
		default:
			// All unrecognized formats are simply a regular string.
			outSchema.GoType = "string"
		}
	default:
		return fmt.Errorf("unhandled Schema type: %s", t)
	}
	return nil
}

// SchemaDescriptor describes a Schema, a type definition.
type SchemaDescriptor struct {
	Fields                   []FieldDescriptor
	HasAdditionalProperties  bool
	AdditionalPropertiesType string
}

// FieldDescriptor describes a field.
type FieldDescriptor struct {
	Required bool   // Is the schema required? If not, we'll pass by pointer
	GoType   string // The Go type needed to represent the json type.
	GoName   string // The Go compatible type name for the type
	JSONName string // The json type name for the type
	IsRef    bool   // Is this schema a reference to predefined object?
}

// GenFieldsFromProperties produces corresponding field names with JSON
// annotations
func GenFieldsFromProperties(props []Property) []string {
	var fields []string
	for i, p := range props {
		field := ""
		// Add a comment to a field in case we have one, otherwise skip.
		if p.Description != "" {
			// Separate the comment from a previous-defined, unrelated field.
			// Make sure the actual field is separated by a newline.
			if i != 0 {
				field += "\n"
			}
			field += fmt.Sprintf("%s\n", StringToGoComment(p.Description))
		}
		field += fmt.Sprintf("    %s %s", p.GoFieldName(), p.GoTypeDef())

		// Support x-omitempty
		omitEmpty := true
		if _, ok := p.ExtensionProps.Extensions[extPropOmitEmpty]; ok {
			if extOmitEmpty, err := extParseOmitEmpty(p.ExtensionProps.Extensions[extPropOmitEmpty]); err == nil {
				omitEmpty = extOmitEmpty
			}
		}

		fieldTags := make(map[string]string)

		if p.Required || p.Nullable || !omitEmpty {
			fieldTags["json"] = p.JSONFieldName
		} else {
			fieldTags["json"] = p.JSONFieldName + ",omitempty"
		}
		if extension, ok := p.ExtensionProps.Extensions[extPropExtraTags]; ok {
			if tags, err := extExtraTags(extension); err == nil {
				keys := SortedStringKeys(tags)
				for _, k := range keys {
					fieldTags[k] = tags[k]
				}
			}
		}
		// Convert the fieldTags map into Go field annotations.
		keys := SortedStringKeys(fieldTags)
		tags := make([]string, len(keys))
		for i, k := range keys {
			tags[i] = fmt.Sprintf(`%s:"%s"`, k, fieldTags[k])
		}
		field += "`" + strings.Join(tags, " ") + "`"
		fields = append(fields, field)
	}
	return fields
}

// GenStructFromSchema creates a struct definition from the given Schema.
// If the schema has additional properties, it is defined as a map[string]Type.
func GenStructFromSchema(schema Schema) string {
	// Start out with struct {
	objectParts := []string{"struct {"}
	// Append all the field definitions
	objectParts = append(objectParts, GenFieldsFromProperties(schema.Properties)...)
	// Close the struct
	if schema.HasAdditionalProperties {
		addPropsType := schema.AdditionalPropertiesType.GoType
		if schema.AdditionalPropertiesType.RefType != "" {
			addPropsType = schema.AdditionalPropertiesType.RefType
		}

		objectParts = append(objectParts,
			fmt.Sprintf("AdditionalProperties map[string]%s `json:\"-\"`", addPropsType))
	}
	objectParts = append(objectParts, "}")
	return strings.Join(objectParts, "\n")
}

// MergeSchemas merges all the fields in the schemas supplied together.
func MergeSchemas(allOf []*openapi3.SchemaRef, path []string) (Schema, error) {
	var outSchema Schema
	for _, schemaOrRef := range allOf {
		ref := schemaOrRef.Ref

		var refType string
		var err error
		if IsGoTypeReference(ref) {
			refType, err = RefPathToGoType(ref)
			if err != nil {
				return Schema{}, fmt.Errorf("error converting reference path to a go type: %w", err)
			}
		}

		schema, err := GenerateGoSchema(schemaOrRef, path)
		if err != nil {
			return Schema{}, fmt.Errorf("error generating Go schema in allOf: %w", err)
		}
		schema.RefType = refType

		for _, p := range schema.Properties {
			err = outSchema.MergeProperty(p)
			if err != nil {
				return Schema{}, fmt.Errorf("error merging properties: %w", err)
			}
		}

		if schema.HasAdditionalProperties {
			if outSchema.HasAdditionalProperties {
				// Both this schema, and the aggregate schema have additional
				// properties, they must match.
				if schema.AdditionalPropertiesType.TypeDecl() != outSchema.AdditionalPropertiesType.TypeDecl() {
					return Schema{}, errors.New("additional properties in allOf have incompatible types")
				}
			} else {
				// We're switching from having no additional properties to having
				// them
				outSchema.HasAdditionalProperties = true
				outSchema.AdditionalPropertiesType = schema.AdditionalPropertiesType
			}
		}
	}

	// Now, we generate the struct which merges together all the fields.
	var err error
	outSchema.GoType, err = GenStructFromAllOf(allOf, path)
	if err != nil {
		return Schema{}, fmt.Errorf("unable to generate aggregate type for AllOf: %w", err)
	}
	return outSchema, nil
}

// GenStructFromAllOf function generates an object that is the union of the objects in the
// input array. In the case of Ref objects, we use an embedded struct, otherwise,
// we inline the fields.
func GenStructFromAllOf(allOf []*openapi3.SchemaRef, path []string) (string, error) {
	// Start out with struct {
	objectParts := []string{"struct {"}
	for _, schemaOrRef := range allOf {
		ref := schemaOrRef.Ref
		if IsGoTypeReference(ref) {
			// We have a referenced type, we will generate an inlined struct
			// member.
			// struct {
			//   InlinedMember
			//   ...
			// }
			goType, err := RefPathToGoType(ref)
			if err != nil {
				return "", err
			}
			objectParts = append(objectParts,
				fmt.Sprintf("   // Embedded struct due to allOf(%s)", ref))
			objectParts = append(objectParts,
				fmt.Sprintf("   %s `yaml:\",inline\"`", goType))
		} else {
			// Inline all the fields from the schema into the output struct,
			// just like in the simple case of generating an object.
			goSchema, err := GenerateGoSchema(schemaOrRef, path)
			if err != nil {
				return "", err
			}
			objectParts = append(objectParts, "   // Embedded fields due to inline allOf schema")
			objectParts = append(objectParts, GenFieldsFromProperties(goSchema.Properties)...)

			if goSchema.HasAdditionalProperties {
				addPropsType := goSchema.AdditionalPropertiesType.GoType
				if goSchema.AdditionalPropertiesType.RefType != "" {
					addPropsType = goSchema.AdditionalPropertiesType.RefType
				}

				additionalPropertiesPart := fmt.Sprintf("AdditionalProperties map[string]%s `json:\"-\"`", addPropsType)
				if !StringInArray(additionalPropertiesPart, objectParts) {
					objectParts = append(objectParts, additionalPropertiesPart)
				}
			}
		}
	}
	objectParts = append(objectParts, "}")
	return strings.Join(objectParts, "\n"), nil
}

// This constructs a Go type for a parameter, looking at either the schema or
// the content, whichever is available
func paramToGoType(param *openapi3.Parameter, path []string) (Schema, error) {
	if param.Content == nil && param.Schema == nil {
		return Schema{}, fmt.Errorf("parameter '%s' has no schema or content", param.Name)
	}

	// We can process the schema through the generic schema processor
	if param.Schema != nil {
		return GenerateGoSchema(param.Schema, path)
	}

	// At this point, we have a content type. We know how to deal with
	// application/json, but if multiple formats are present, we can't do anything,
	// so we'll return the parameter as a string, not bothering to decode it.
	if len(param.Content) > 1 {
		return Schema{
			GoType:      "string",
			Description: StringToGoComment(param.Description),
		}, nil
	}

	// Otherwise, look for application/json in there
	mt, found := param.Content["application/json"]
	if !found {
		// If we don't have json, it's a string
		return Schema{
			GoType:      "string",
			Description: StringToGoComment(param.Description),
		}, nil
	}

	// For json, we go through the standard schema mechanism
	return GenerateGoSchema(mt.Schema, path)
}
