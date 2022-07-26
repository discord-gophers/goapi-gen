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
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/kenshaw/snaker"
)

var pathParamRE *regexp.Regexp

func init() {
	pathParamRE = regexp.MustCompile(`{[.;?]?([^{}*]+)\*?}`)
}

// UppercaseFirstCharacter uppercases the first character of str.
func UppercaseFirstCharacter(str string) string {
	if str == "" {
		return ""
	}
	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// LowercaseFirstCharacter lowercases te first character of str.
func LowercaseFirstCharacter(str string) string {
	if str == "" {
		return ""
	}
	runes := []rune(str)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// ToCamelCase converts a string to camel case with proper Go initialisms.
func ToCamelCase(str string) string {
	if str != "" && unicode.IsDigit([]rune(str)[0]) {
		// FIXME this is so hacky please help
		str = "F" + str
		str = snaker.ForceCamelIdentifier(str)
		return str[1:]
	}
	return snaker.ForceCamelIdentifier(str)
}

// ToSnakeCase converts a string to snake case.
func ToSnakeCase(str string) string {
	str = strings.Replace(str, "-", "_", -1)
	return snaker.CamelToSnake(str)
}

// SortedSchemaKeys returns the keys of dict in alphabetically.
func SortedSchemaKeys(dict map[string]*openapi3.SchemaRef) []string {
	keys := make([]string, len(dict))
	i := 0
	for key := range dict {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}

// SortedPathsKeys returns the keys of dict alphabetically.
func SortedPathsKeys(dict openapi3.Paths) []string {
	keys := make([]string, len(dict))
	i := 0
	for key := range dict {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}

// SortedOperationsKeys returns the keys of dict alphabetically.
func SortedOperationsKeys(dict map[string]*openapi3.Operation) []string {
	keys := make([]string, len(dict))
	i := 0
	for key := range dict {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}

// SortedResponsesKeys returns the keys of dict alphabetically.
func SortedResponsesKeys(dict openapi3.Responses) []string {
	keys := make([]string, len(dict))
	i := 0
	for key := range dict {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}

// SortedContentKeys returns the keys of dict alphabetically.
func SortedContentKeys(dict openapi3.Content) []string {
	keys := make([]string, len(dict))
	i := 0
	for key := range dict {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}

// SortedStringKeys returns the keys of dict alphabetically.
func SortedStringKeys(dict map[string]string) []string {
	keys := make([]string, len(dict))
	i := 0
	for key := range dict {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}

// SortedParameterKeys returns the keys of dict alphabetically.
func SortedParameterKeys(dict map[string]*openapi3.ParameterRef) []string {
	keys := make([]string, len(dict))
	i := 0
	for key := range dict {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}

// SortedRequestBodyKeys returns the keys of dict alphabetically.
func SortedRequestBodyKeys(dict map[string]*openapi3.RequestBodyRef) []string {
	keys := make([]string, len(dict))
	i := 0
	for key := range dict {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}

// SortedSecurityRequirementKeys eturns the keys of dict alphabetically.
func SortedSecurityRequirementKeys(dict openapi3.SecurityRequirement) []string {
	keys := make([]string, len(dict))
	i := 0
	for key := range dict {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}

// StringInArray returns if strs contains str.
func StringInArray(str string, strs []string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}
	return false
}

// RefPathToGoType converts refPath to a Go type name.
// #/components/schemas/Foo -> Foo
// #/components/parameters/Bar -> Bar
// #/components/responses/Baz -> Baz
// Remote components (document.json#/Foo) are supported if they present in --import-mapping
// URL components (http://deepmap.com/schemas/document.json#/Foo) are supported if they present in --import-mapping
// Remote and URL also support standard local paths even though the spec doesn't mention them.
func RefPathToGoType(refPath string) (string, error) {
	return refPathToGoType(refPath, true)
}

// refPathToGoType returns the Go typename for refPath given its
func refPathToGoType(refPath string, local bool) (string, error) {
	if refPath[0] == '#' {
		pathParts := strings.Split(refPath, "/")
		depth := len(pathParts)
		if local {
			if depth != 4 {
				return "", fmt.Errorf("unexpected reference depth: %d for ref: %s local: %t", depth, refPath, local)
			}
		} else if depth != 4 && depth != 2 {
			return "", fmt.Errorf("unexpected reference depth: %d for ref: %s local: %t", depth, refPath, local)
		}
		return SchemaNameToTypeName(pathParts[len(pathParts)-1]), nil
	}
	pathParts := strings.Split(refPath, "#")
	if len(pathParts) != 2 {
		return "", fmt.Errorf("unsupported reference: %s", refPath)
	}

	remoteComponent, flatComponent := pathParts[0], pathParts[1]
	goImport, ok := importMapping[remoteComponent]
	if !ok {
		return "", fmt.Errorf("unrecognized external reference '%s'; please provide the known import for this reference using option --import-mapping", remoteComponent)
	}

	goType, err := refPathToGoType("#"+flatComponent, false)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s.%s", goImport.Name, goType), nil
}

// IsGoTypeReference checks if ref links to a valid go type.
// #/components/schemas/Foo                     -> true
// ./local/file.yml#/components/parameters/Bar  -> true
// ./local/file.yml                             -> false
func IsGoTypeReference(ref string) bool {
	return ref != "" && !IsWholeDocumentReference(ref)
}

// IsWholeDocumentReference checks if ref is a whole document reference.
// #/components/schemas/Foo                             -> false
// ./local/file.yml#/components/parameters/Bar          -> false
// ./local/file.yml                                     -> true
// http://example.com/schemas/document.json             -> true
// http://example.com/schemas/document.json#/Foo        -> false
func IsWholeDocumentReference(ref string) bool {
	return ref != "" && !strings.ContainsAny(ref, "#")
}

// SwaggerURIToChiURI converts uri to a Chi-style URI.
// It replaces all swagger parameters with {param}.
//
// Valid input parameters are:
//   {param}
//   {param*}
//   {.param}
//   {.param*}
//   {;param}
//   {;param*}
//   {?param}
//   {?param*}
func SwaggerURIToChiURI(uri string) string {
	return pathParamRE.ReplaceAllString(uri, "{$1}")
}

// OrderedParamsFromURI returns argument names in uri.
// Given /path/{param1}/{.param2*}/{?param3},
// returns [param1, param2, param3]
func OrderedParamsFromURI(uri string) []string {
	matches := pathParamRE.FindAllStringSubmatch(uri, -1)
	result := make([]string, len(matches))
	for i, m := range matches {
		result[i] = m[1]
	}
	return result
}

// ReplacePathParamsWithStr replaces uri parameters with %s
func ReplacePathParamsWithStr(uri string) string {
	return pathParamRE.ReplaceAllString(uri, "%s")
}

// SortParamsByPath sorts in to match those in the path URI.
func SortParamsByPath(path string, in []ParameterDefinition) ([]ParameterDefinition, error) {
	pathParams := OrderedParamsFromURI(path)
	if len(pathParams) != len(in) {
		return nil, fmt.Errorf("path '%s' has %d positional parameters, but spec has %d declared",
			path, len(pathParams), len(in))
	}

	out := make([]ParameterDefinition, len(in))
	for i, name := range pathParams {
		p := ParameterDefinitions(in).FindByName(name)
		if p == nil {
			return nil, fmt.Errorf("path '%s' refers to parameter '%s', which doesn't exist in specification",
				path, name)
		}
		out[i] = *p
	}
	return out, nil
}

// IsGoKeyword checks if str is a reserved keyword.
// Returns whether the given string is a go keyword
func IsGoKeyword(str string) bool {
	keywords := []string{
		"break",
		"case",
		"chan",
		"const",
		"continue",
		"default",
		"defer",
		"else",
		"fallthrough",
		"for",
		"func",
		"go",
		"goto",
		"if",
		"import",
		"interface",
		"map",
		"package",
		"range",
		"return",
		"select",
		"struct",
		"switch",
		"type",
		"var",
	}

	for _, k := range keywords {
		if k == str {
			return true
		}
	}
	return false
}

// IsPredeclaredGoIdentifier returns whether str is a go indentifier.
func IsPredeclaredGoIdentifier(str string) bool {
	predeclaredIdentifiers := []string{
		// Types
		"bool",
		"byte",
		"complex64",
		"complex128",
		"error",
		"float32",
		"float64",
		"int",
		"int8",
		"int16",
		"int32",
		"int64",
		"rune",
		"string",
		"uint",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"uintptr",
		// Constants
		"true",
		"false",
		"iota",
		// Zero value
		"nil",
		// Functions
		"append",
		"cap",
		"close",
		"complex",
		"copy",
		"delete",
		"imag",
		"len",
		"make",
		"new",
		"panic",
		"print",
		"println",
		"real",
		"recover",
	}

	for _, k := range predeclaredIdentifiers {
		if k == str {
			return true
		}
	}

	return false
}

// IsGoIdentity checks if str is a valid go identifier.
func IsGoIdentity(str string) bool {
	for i, c := range str {
		if !isValidRuneForGoID(i, c) {
			return false
		}
	}

	return IsGoKeyword(str)
}

func isValidRuneForGoID(index int, char rune) bool {
	if index == 0 && unicode.IsNumber(char) {
		return false
	}

	return unicode.IsLetter(char) || char == '_' || unicode.IsNumber(char)
}

// IsValidGoIdentity checks if str can be used as a name of variable, constant,
// or type.
func IsValidGoIdentity(str string) bool {
	return !IsPredeclaredGoIdentifier(str) && !IsGoIdentity(str)
}

// SanitizeGoIdentity replaces illegal characters in str.
func SanitizeGoIdentity(str string) string {
	sanitized := []rune(str)

	for i, c := range sanitized {
		if !isValidRuneForGoID(i, c) {
			sanitized[i] = '_'
		} else {
			sanitized[i] = c
		}
	}

	str = string(sanitized)

	if IsGoKeyword(str) || IsPredeclaredGoIdentifier(str) {
		str = "_" + str
	}

	if !IsValidGoIdentity(str) {
		panic("here is a bug")
	}

	return str
}

// SanitizeEnumNames removes illegal and duplicates chars in enum names.
func SanitizeEnumNames(enumNames []string) map[string]string {
	dupCheck := make(map[string]int, len(enumNames))
	deDup := make([]string, 0, len(enumNames))

	for _, n := range enumNames {
		if _, dup := dupCheck[n]; !dup {
			deDup = append(deDup, n)
		}
		dupCheck[n] = 0
	}

	dupCheck = make(map[string]int, len(deDup))
	sanitizedDeDup := make(map[string]string, len(deDup))

	for _, n := range deDup {
		sanitized := SanitizeGoIdentity(SchemaNameToTypeName(n))

		if _, dup := dupCheck[sanitized]; !dup {
			sanitizedDeDup[sanitized] = n
		} else {
			sanitizedDeDup[sanitized+strconv.Itoa(dupCheck[sanitized])] = n
		}
		dupCheck[sanitized]++
	}

	return sanitizedDeDup
}

// SchemaNameToTypeName converts name to a valid Go type name.
// It converts name to camel case and is valid in Go.
func SchemaNameToTypeName(name string) string {
	if name == "$" {
		name = "DollarSign"
	} else {
		name = ToCamelCase(name)
		// Prepend "N" to schemas starting with a number
		if name != "" && unicode.IsDigit([]rune(name)[0]) {
			name = "N" + name
		}
	}
	return name
}

// SchemaHasAdditionalProperties checks if schema has additional properties.
//
// According to the spec, additionalProperties may be true, false, or a
// schema. If not present, true is implied. If it's a schema, true is implied.
// If it's false, no additional properties are allowed. We're going to act a little
// differently, in that if you want additionalProperties code to be generated,
// you must specify an additionalProperties type
// If additionalProperties it true/false, this field will be non-nil.
func SchemaHasAdditionalProperties(schema *openapi3.Schema) bool {
	if schema.AdditionalPropertiesAllowed != nil && *schema.AdditionalPropertiesAllowed {
		return true
	}

	if schema.AdditionalProperties != nil {
		return true
	}
	return false
}

// PathToTypeName converts path to a go type name.
// It converts each entry in path to camel case and joins them with _.
func PathToTypeName(path []string) string {
	for i, p := range path {
		path[i] = ToCamelCase(p)
	}
	return strings.Join(path, "_")
}

// StringToGoComment renders a possible multi-line string to a valid Go-Comment.
// Each line is prefixed as a comment.
func StringToGoComment(in string) string {
	if len(in) == 0 || len(strings.TrimSpace(in)) == 0 { // ignore empty comment
		return ""
	}

	// Normalize newlines from Windows/Mac to Linux
	in = strings.ReplaceAll(in, "\r\n", "\n")
	in = strings.ReplaceAll(in, "\r", "\n")

	// Add comment to each line
	var lines []string
	for _, line := range strings.Split(in, "\n") {
		lines = append(lines, fmt.Sprintf("// %s", line))
	}
	in = strings.Join(lines, "\n")

	// in case we have a multiline string which ends with \n, we would generate
	// empty-line-comments, like `// `. Therefore remove this line comment.
	in = strings.TrimSuffix(in, "\n// ")
	return in
}

// EscapePathElements escapes non path parameters in path and url encodes them.
func EscapePathElements(path string) string {
	elems := strings.Split(path, "/")
	for i, e := range elems {
		if strings.HasPrefix(e, "{") && strings.HasSuffix(e, "}") {
			// This is a path parameter, we don't want to mess with its value
			continue
		}
		elems[i] = url.QueryEscape(e)
	}
	return strings.Join(elems, "/")
}
