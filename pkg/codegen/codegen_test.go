package codegen

import (
	"bytes"
	"go/format"
	"io"
	"net/http"
	"testing"

	examplePetstore "github.com/discord-gophers/goapi-gen/examples/petstore-expanded/api"
	examplePetstoreClient "github.com/discord-gophers/goapi-gen/examples/petstore-expanded/client"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/golangci/lint-1"
	"github.com/stretchr/testify/assert"
)

func TestExamplePetStoreCodeGeneration(t *testing.T) {
	// Input vars for code generation:
	packageName := "api"
	opts := Options{
		GenerateChiServer: true,
		GenerateClient:    true,
		GenerateTypes:     true,
		EmbedSpec:         true,
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

	// Check that the client method signatures return response structs:
	assert.Contains(t, code, "func (c *Client) FindPetByID(ctx context.Context, id int64, reqEditors ...RequestEditorFn) (*http.Response, error) {")

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

func TestGenerateResponseBodyTypes(t *testing.T) {
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

	// Check for response type
	assert.Contains(t, code, "type Response struct {")

	// Check for expected response constructors:
	assert.Contains(t, code, "func GetTestByNameJSON200Response(body []Test) *Response {")
	assert.Contains(t, code, "func GetTestByNameXML200Response(body []Test) *Response {")
	assert.Contains(t, code, "func GetTestByNameJSON422Response(body []interface{}) *Response {")
	assert.Contains(t, code, "func GetTestByNameXML422Response(body []interface{}) *Response {")
	assert.Contains(t, code, "func GetTestByNameJSONDefaultResponse(body Error) *Response {")
	assert.Contains(t, code, "func GetCatStatusJSON200Response(body interface{}) *Response {")
	assert.Contains(t, code, "func GetCatStatusXML200Response(body interface{}) *Response {")
	assert.Contains(t, code, "func GetCatStatusJSONDefaultResponse(body Error) *Response {")
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

	// Check that the built-in template has been overriden
	assert.Contains(t, code, "//blah")
}

func TestExamplePetStoreParseFunction(t *testing.T) {
	bodyBytes := []byte(`{"id": 5, "name": "testpet", "tag": "cat"}`)

	cannedResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader(bodyBytes)),
		Header:     http.Header{},
	}
	cannedResponse.Header.Add("Content-type", "application/json")

	findPetByIDResponse, err := examplePetstoreClient.ParseFindPetByIDResponse(cannedResponse)
	assert.NoError(t, err)
	assert.NotNil(t, findPetByIDResponse.JSON200)
	assert.Equal(t, int64(5), findPetByIDResponse.JSON200.ID)
	assert.Equal(t, "testpet", findPetByIDResponse.JSON200.Name)
	assert.NotNil(t, findPetByIDResponse.JSON200.Tag)
	assert.Equal(t, "cat", *findPetByIDResponse.JSON200.Tag)
}

func TestExampleOpenAPICodeGeneration(t *testing.T) {
	// Input vars for code generation:
	packageName := "testswagger"
	opts := Options{
		GenerateClient: true,
		GenerateTypes:  true,
		EmbedSpec:      true,
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

	// Check that we have a package:
	assert.Contains(t, code, "package testswagger")

	// Check that response structs are generated correctly:
	assert.Contains(t, code, "type GetTestByNameResponse struct {")

	// Check that response structs contains fallbacks to interface for invalid types:
	// Here an invalid array with no items.
	assert.Contains(t, code, `
type GetTestByNameResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]Test
	XML200       *[]Test
	JSON422      *[]interface{}
	XML422       *[]interface{}
	JSONDefault  *Error
}`)

	// Check that the helper methods are generated correctly:
	assert.Contains(t, code, "func (r GetTestByNameResponse) Status() string {")
	assert.Contains(t, code, "func (r GetTestByNameResponse) StatusCode() int {")
	assert.Contains(t, code, "func ParseGetTestByNameResponse(rsp *http.Response) (*GetTestByNameResponse, error) {")

	// Check the client method signatures:
	assert.Contains(t, code, "type GetTestByNameParams struct {")
	assert.Contains(t, code, "Top *int `json:\"$top,omitempty\"`")
	assert.Contains(t, code, "func (c *Client) GetTestByName(ctx context.Context, name string, params *GetTestByNameParams, reqEditors ...RequestEditorFn) (*http.Response, error) {")
	assert.Contains(t, code, "func (c *ClientWithResponses) GetTestByNameWithResponse(ctx context.Context, name string, params *GetTestByNameParams, reqEditors ...RequestEditorFn) (*GetTestByNameResponse, error) {")
	assert.Contains(t, code, "DeadSince *time.Time    `json:\"dead_since,omitempty\" tag1:\"value1\" tag2:\"value2\"`")

	// Make sure the generated code is valid:
	linter := new(lint.Linter)
	problems, err := linter.Lint("test.gen.go", []byte(code))
	assert.NoError(t, err)
	assert.Len(t, problems, 0)
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
