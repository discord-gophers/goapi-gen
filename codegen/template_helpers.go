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
	"sort"
	"strings"
	"text/template"
	"unicode"

	"github.com/kenshaw/snaker"
)

var (
	contentTypesJSON = []string{"application/json", "text/x-json"}
	contentTypesYAML = []string{"application/yaml", "application/x-yaml", "text/yaml", "text/x-yaml"}
	contentTypesXML  = []string{"application/xml", "text/xml"}
)

// This function takes an array of Parameter definition, and generates a valid
// Go parameter declaration from them, eg:
// ", foo int, bar string, baz float32". The preceding comma is there to save
// a lot of work in the template engine.
func genParamArgs(params []ParameterDefinition) string {
	if len(params) == 0 {
		return ""
	}
	parts := make([]string, len(params))
	for i, p := range params {
		paramName := p.GoVariableName()
		parts[i] = fmt.Sprintf("%s %s", paramName, p.TypeDef())
	}
	return ", " + strings.Join(parts, ", ")
}

// This is another variation of the function above which generates only the
// parameter names:
// ", foo, bar, baz"
func genParamNames(params []ParameterDefinition) string {
	if len(params) == 0 {
		return ""
	}
	parts := make([]string, len(params))
	for i, p := range params {
		parts[i] = p.GoVariableName()
	}
	return ", " + strings.Join(parts, ", ")
}

func getResponseTypeDefinitions(op *OperationDefinition) []ResponseTypeDefinition {
	td, err := op.GetResponseTypeDefinitions()
	if err != nil {
		panic(err)
	}
	return td
}

func getTaggedMiddlewares(ops []OperationDefinition) []string {
	middlewares := make(map[string]struct{})
	for _, op := range ops {
		for _, m := range op.Middlewares {
			middlewares[m] = struct{}{}
		}
	}

	keys := make([]string, 0, len(middlewares))
	for k := range middlewares {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}

// This outputs a string array
func toStringArray(sarr []string) string {
	return `[]string{"` + strings.Join(sarr, `","`) + `"}`
}

func responseNameToStatusCode(responseName string) string {
	switch strings.ToUpper(responseName) {
	case "DEFAULT":
		return "200"
	case "1XX", "2XX", "3XX", "4XX", "5XX":
		return fmt.Sprintf("%s00", responseName[:1])
	default:
		return responseName
	}
}

func responseToStatusRangeString(responseName string) string {
	switch strings.ToUpper(responseName) {
	case "DEFAULT":
		return "resp.StatusCode != 0"
	case "1XX", "2XX", "3XX", "4XX", "5XX":
		return fmt.Sprintf("resp.StatusCode < %s00 && resp.StatusCode > %s99", responseName[:1], responseName[:1])
	}
	return fmt.Sprintf("resp.StatusCode != %s", responseName)
}

// TitleWord converts a single worded string to title case.
// This is a replacement to `strings.Title` which we used previously.
// We didn't need strings.Title word boundary rules, and just want to Title the words directly,
// to export them in most cases.
// Thus, this function is simpler and more efficient.
func TitleWord(s string) string {
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// ToComment converts a string to a comment.
func ToComment(s string) string {
	return "// " + strings.ReplaceAll(s, "\n", "\n// ")
}

// TemplateFunctions generates the list of utlity and helpfer functions used by
// the templates.
var TemplateFunctions = template.FuncMap{
	"genParamArgs":               genParamArgs,
	"genParamNames":              genParamNames,
	"getResponseTypeDefinitions": getResponseTypeDefinitions,
	"genTaggedMiddleware":        getTaggedMiddlewares,
	"toStringArray":              toStringArray,
	"toComment":                  ToComment,

	"swaggerURIToChiURI": SwaggerURIToChiURI,

	"statusCode":      responseNameToStatusCode,
	"statusCodeRange": responseToStatusRangeString,

	"ucFirst": snaker.ForceCamelIdentifier,
	"lower":   strings.ToLower,
	"title":   TitleWord,
}
