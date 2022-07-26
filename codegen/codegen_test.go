package codegen

import (
	"go/format"
	"testing"
	"text/template"

	examplePetstore "github.com/discord-gophers/goapi-gen/examples/petstore-expanded/api"

	"github.com/discord-gophers/goapi-gen/templates"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/golangci/lint-1"
	"github.com/stretchr/testify/assert"
)

func TestExamplePetStoreCodeGeneration(t *testing.T) {
	// Input vars for code generation:
	packageName := "api"
	opts := Options{
		GenerateServer: true,
		GenerateTypes:  true,
		EmbedSpec:      true,
	}

	// Get a spec from the example PetStore definition:
	swagger, err := examplePetstore.GetSwagger()
	assert.NoError(t, err)

	// Run our code generation:
	code, err := Generate(swagger, packageName, opts)
	assert.NoError(t, err)
	assert.NotEmpty(t, code)

	// Check that we have valid (formattable) code:
	_, err = format.Source([]byte(code))
	assert.NoError(t, err)

	// Check that we have a package:
	assert.Contains(t, code, "package api")

	// Check that the property comments were generated
	assert.Contains(t, code, "// Unique id of the pet")

	// Check that the summary comment contains newlines
	assert.Contains(t, code, `// Deletes a pet by ID
	// (DELETE /pets/{id})
`)

	// Make sure the generated code is valid:
	linter := new(lint.Linter)
	problems, err := linter.Lint("test.gen.go", []byte(code))
	assert.NoError(t, err)
	assert.Len(t, problems, 0)
}

func TestGenerateRequestBindMethods(t *testing.T) {
	packageName := "api"
	opts := Options{
		GenerateTypes: true,
	}

	// Get the sample spec:
	swagger, err := openapi3.NewLoader().LoadFromData([]byte(testOpenAPIDefinition))
	assert.NoError(t, err)

	// Run our code generation:
	code, err := Generate(swagger, packageName, opts)
	assert.NoError(t, err)
	assert.NotEmpty(t, code)

	// Check that we have valid (formattable) code:
	_, err = format.Source([]byte(code))
	assert.NoError(t, err)

	// Check that we have a package:
	assert.Contains(t, code, "package api")

	// func (AddPetJSONRequestBody) Bind(*http.Request) error {

	// Check for expected request binders:
	assert.Contains(t, code, "func (CreateLiveCatJSONRequestBody) Bind(*http.Request) error {")

	// Check for unbindable request types:
	assert.NotContains(t, code, "func (CreateCatJSONRequestBody) Bind(*http.Request) error {")
}

func TestExamplePetStoreCodeGenerationWithUserTemplates(t *testing.T) {
	userTemplates := map[string]string{"typedef.tmpl": "//blah"}

	// Input vars for code generation:
	packageName := "api"
	opts := Options{
		GenerateTypes: true,
		UserTemplates: userTemplates,
	}

	// Get a spec from the example PetStore definition:
	swagger, err := examplePetstore.GetSwagger()
	assert.NoError(t, err)

	// Run our code generation:
	code, err := Generate(swagger, packageName, opts)
	assert.NoError(t, err)
	assert.NotEmpty(t, code)

	// Check that we have valid (formattable) code:
	_, err = format.Source([]byte(code))
	assert.NoError(t, err)

	// Check that we have a package:
	assert.Contains(t, code, "package api")

	// Check that the built-in template has been overridden
	assert.Contains(t, code, "//blah")
}

func TestExampleOpenAPICodeGeneration(t *testing.T) {
	// Input vars for code generation:
	packageName := "testswagger"
	opts := Options{
		GenerateTypes: true,
		EmbedSpec:     true,
	}

	// Get a spec from the test definition in this file:
	swagger, err := openapi3.NewLoader().LoadFromData([]byte(testOpenAPIDefinition))
	assert.NoError(t, err)

	// Run our code generation:
	code, err := Generate(swagger, packageName, opts)
	assert.NoError(t, err)
	assert.NotEmpty(t, code)

	// Check that we have valid (formattable) code:
	_, err = format.Source([]byte(code))
	assert.NoError(t, err)
}

const testOpenAPIDefinition = `
openapi: 3.0.1

info:
  title: OpenAPI-CodeGen Test
  description: 'This is a test OpenAPI Spec'
  version: 1.0.0

servers:
- url: https://test.goapi-gen.com/v2
- url: http://test.goapi-gen.com/v2

paths:
  /test/{name}:
    get:
      tags:
      - test
      summary: Get test
      operationId: getTestByName
      parameters:
      - name: name
        in: path
        required: true
        schema:
          type: string
      - name: $top
        in: query
        required: false
        schema:
          type: integer
      responses:
        200:
          description: Success
          content:
            application/xml:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Test'
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Test'
        422:
          description: InvalidArray
          content:
            application/xml:
              schema:
                type: array
            application/json:
              schema:
                type: array
        default:
          description: Error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
  /cat:
    get:
      tags:
      - cat
      summary: Get cat status
      operationId: getCatStatus
      responses:
        200:
          description: Success
          content:
            application/json:
              schema:
                oneOf:
                - $ref: '#/components/schemas/CatAlive'
                - $ref: '#/components/schemas/CatDead'
            application/xml:
              schema:
                anyOf:
                - $ref: '#/components/schemas/CatAlive'
                - $ref: '#/components/schemas/CatDead'
            application/yaml:
              schema:
                allOf:
                - $ref: '#/components/schemas/CatAlive'
                - $ref: '#/components/schemas/CatDead'
        default:
          description: Error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
    post:
      tags:
      - cat
      summary: Create a new cat
      operationId: createCat
      requestBody:
        description: Cat to add to the store
        required: true
        content:
          application/json:
            schema:
              oneOf:
              - $ref: '#/components/schemas/CatAlive'
              - $ref: '#/components/schemas/CatDead'
      responses:
        201:
          description: Success
          content:
            application/json:
              schema:
                oneOf:
                - $ref: '#/components/schemas/CatAlive'
                - $ref: '#/components/schemas/CatDead'
        default:
          description: Error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /live-cat:
    post:
      summary: Create a new live cat
      operationId: createLiveCat
      requestBody:
        description: Cat to add to the store
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CatAlive'
      responses:
        201:
          description: Success
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/CatAlive'
        default:
          description: Error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

components:
  schemas:

    Test:
      properties:
        name:
          type: string
        cases:
          type: array
          items:
            $ref: '#/components/schemas/TestCase'

    TestCase:
      properties:
        name:
          type: string
        command:
          type: string

    Error:
      properties:
        code:
          type: integer
          format: int32
        message:
          type: string

    CatAlive:
      properties:
        name:
          type: string
        alive_since:
          type: string
          format: date-time

    CatDead:
      properties:
        name:
          type: string
        dead_since:
          type: string
          format: date-time
          x-go-extra-tags:
            tag1: value1
            tag2: value2
        cause:
          type: string
          enum: [car, dog, oldage]
`

func TestGenerateEnumTypes(t *testing.T) {
	tests := []struct {
		name    string
		types   []TypeDefinition
		want    string
		wantErr bool
	}{
		{
			name: "string",

			types: []TypeDefinition{{
				JSONName: "my_type",
				TypeName: "MyType",
				Schema: Schema{
					GoType:     "string",
					RefType:    "string",
					EnumValues: map[string]string{"some": "value"},
				},
			}},
			want: `
// MyType defines model for my_type.
type MyType struct {
    value string
}
func (t *MyType) ToValue() string {
    return t.value
}
func (t *MyType) MarshalJSON() ([]byte, error) {
    return json.Marshal(t.value)
}
func (t *MyType) UnmarshalJSON(data []byte) error {
    var value string
    if err := json.Unmarshal(data, &value); err != nil {
        return err
    }
    return t.FromValue(value)
}
func (t *MyType) FromValue(value string) error {
    switch value {
    
    case some.value:
        t.value = value
        return nil
    
    }
    return fmt.Errorf("unknown enum value: %v", value)
}
`,
		},
		{
			name: "int64",

			types: []TypeDefinition{{
				JSONName: "my_type",
				TypeName: "MyType",
				Schema: Schema{
					GoType:     "int64",
					RefType:    "int64",
					EnumValues: map[string]string{"some": "value"},
				},
			}},
			want: `
// MyType defines model for my_type.
type MyType struct {
    value int64
}
func (t *MyType) ToValue() int64 {
    return t.value
}
func (t *MyType) MarshalJSON() ([]byte, error) {
    return json.Marshal(t.value)
}
func (t *MyType) UnmarshalJSON(data []byte) error {
    var value int64
    if err := json.Unmarshal(data, &value); err != nil {
        return err
    }
    return t.FromValue(value)
}
func (t *MyType) FromValue(value int64) error {
    switch value {
    
    case some.value:
        t.value = value
        return nil
    
    }
    return fmt.Errorf("unknown enum value: %v", value)
}
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TemplateFunctions["opts"] = func() Options { return Options{} }
			tmpl := template.New("goapi-gen").Funcs(TemplateFunctions)
			// This parses all of our own template files into the template object
			// above
			tmpl, _ = templates.Parse(tmpl)

			got, err := GenerateEnumTypes(tmpl, tt.types)
			if (err != nil) != tt.wantErr {
				assert.NotNil(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
