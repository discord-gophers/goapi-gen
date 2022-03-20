// Package components provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/discord-gophers/goapi-gen version (devel) DO NOT EDIT.
package components

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/discord-gophers/goapi-gen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// Has additional properties of type int
type AdditionalPropertiesObject1 struct {
	ID                   int            `json:"id"`
	Name                 string         `json:"name"`
	Optional             *string        `json:"optional,omitempty"`
	AdditionalProperties map[string]int `json:"-"`
}

// Does not allow additional properties
type AdditionalPropertiesObject2 struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Allows any additional property
type AdditionalPropertiesObject3 struct {
	Name                 string                 `json:"name"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// Has anonymous field which has additional properties
type AdditionalPropertiesObject4 struct {
	Inner                AdditionalPropertiesObject4_Inner `json:"inner"`
	Name                 string                            `json:"name"`
	AdditionalProperties map[string]interface{}            `json:"-"`
}

// AdditionalPropertiesObject4_Inner defines model for AdditionalPropertiesObject4.Inner.
type AdditionalPropertiesObject4_Inner struct {
	Name                 string                 `json:"name"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// Has additional properties with schema for dictionaries
type AdditionalPropertiesObject5 struct {
	AdditionalProperties map[string]SchemaObject `json:"-"`
}

// Tests x-go-string
type GoStringTag struct {
	ID   int    `json:"id,string"`
	Name string `json:"name"`
}

// ObjectWithJSONField defines model for ObjectWithJsonField.
type ObjectWithJSONField struct {
	Name   string          `json:"name"`
	Value1 json.RawMessage `json:"value1"`
	Value2 json.RawMessage `json:"value2,omitempty"`
}

// Tests x-go-optional-value on the whole object
type OptionalObject struct {
	Str1    string  `json:"str1,omitempty"`
	Str2    string  `json:"str2,omitempty"`
	StrPtr1 *string `json:"str_ptr1,omitempty"`
	StrPtr2 *string `json:"str_ptr2,omitempty"`
}

// Tests x-go-optional-value
type OptionalValue struct {
	Str1 string `json:"str1,omitempty"`
	Str2 string `json:"str2,omitempty"`
}

// SchemaObject defines model for SchemaObject.
type SchemaObject struct {
	FirstName string `json:"firstName"`
	Role      string `json:"role"`
}

// ResponseObject defines model for ResponseObject.
type ResponseObject struct {
	Field SchemaObject `json:"Field"`
}

// RequestBody defines model for RequestBody.
type RequestBody struct {
	Field SchemaObject `json:"Field"`
}

// ParamsWithAddPropsParams_P1 defines parameters for ParamsWithAddProps.
type ParamsWithAddPropsParams_P1 struct {
	AdditionalProperties map[string]interface{} `json:"-"`
}

// ParamsWithAddPropsParams defines parameters for ParamsWithAddProps.
type ParamsWithAddPropsParams struct {
	// This parameter has additional properties
	P1 ParamsWithAddPropsParams_P1 `json:"p1"`

	// This parameter has an anonymous inner property which needs to be
	// turned into a proper type for additionalProperties to work
	P2 struct {
		Inner ParamsWithAddPropsParams_P2_Inner `json:"inner"`
	} `json:"p2"`
}

// ParamsWithAddPropsParams_P2_Inner defines parameters for ParamsWithAddProps.
type ParamsWithAddPropsParams_P2_Inner struct {
	AdditionalProperties map[string]string `json:"-"`
}

// BodyWithAddPropsJSONBody defines parameters for BodyWithAddProps.
type BodyWithAddPropsJSONBody struct {
	Inner                BodyWithAddPropsJSONBody_Inner `json:"inner"`
	Name                 string                         `json:"name"`
	AdditionalProperties map[string]interface{}         `json:"-"`
}

// BodyWithAddPropsJSONBody_Inner defines parameters for BodyWithAddProps.
type BodyWithAddPropsJSONBody_Inner struct {
	AdditionalProperties map[string]int `json:"-"`
}

// EnsureEverythingIsReferencedJSONRequestBody defines body for EnsureEverythingIsReferenced for application/json ContentType.
type EnsureEverythingIsReferencedJSONRequestBody RequestBody

// Bind implements render.Binder.
func (EnsureEverythingIsReferencedJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// BodyWithAddPropsJSONRequestBody defines body for BodyWithAddProps for application/json ContentType.
type BodyWithAddPropsJSONRequestBody BodyWithAddPropsJSONBody

// Bind implements render.Binder.
func (BodyWithAddPropsJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// Response is a common response struct for all the API calls.
// A Response object may be instantiated via functions for specific operation responses.
// It may also be instantiated directly, for the purpose of responding with a single status code.
type Response struct {
	body        interface{}
	StatusCode  int
	contentType string
}

// Render implements the render.Renderer interface. It sets the Content-Type header
// and status code based on the response definition.
func (resp *Response) Render(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", resp.contentType)
	render.Status(r, resp.StatusCode)
	return nil
}

// Status is a builder method to override the default status code for a response.
func (resp *Response) Status(code int) *Response {
	resp.StatusCode = code
	return resp
}

// ContentType is a builder method to override the default content type for a response.
func (resp *Response) ContentType(contentType string) *Response {
	resp.contentType = contentType
	return resp
}

// MarshalJSON implements the json.Marshaler interface.
// This is used to only marshal the body of the response.
func (resp *Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.body)
}

// MarshalXML implements the xml.Marshaler interface.
// This is used to only marshal the body of the response.
func (resp *Response) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.Encode(resp.body)
}

// EnsureEverythingIsReferencedJSON200Response is a constructor method for a EnsureEverythingIsReferenced response.
// A *Response is returned with the configured status code and content type from the spec.
func EnsureEverythingIsReferencedJSON200Response(body struct {
	// Has additional properties with schema for dictionaries
	Five *AdditionalPropertiesObject5 `json:"five,omitempty"`

	// Has anonymous field which has additional properties
	Four *AdditionalPropertiesObject4 `json:"four,omitempty"`

	// Tests x-go-string
	GoStringTag *GoStringTag         `json:"goStringTag,omitempty"`
	JSONField   *ObjectWithJSONField `json:"jsonField,omitempty"`

	// Has additional properties of type int
	One *AdditionalPropertiesObject1 `json:"one,omitempty"`

	// Tests x-go-optional-value on the whole object
	OptionalObject OptionalObject `json:"optionalObject,omitempty"`

	// Tests x-go-optional-value
	OptionalValue *OptionalValue `json:"optionalValue,omitempty"`

	// Allows any additional property
	Three *AdditionalPropertiesObject3 `json:"three,omitempty"`

	// Does not allow additional properties
	Two *AdditionalPropertiesObject2 `json:"two,omitempty"`
}) *Response {
	return &Response{
		body:        body,
		StatusCode:  200,
		contentType: "application/json",
	}
}

// EnsureEverythingIsReferencedJSONDefaultResponse is a constructor method for a EnsureEverythingIsReferenced response.
// A *Response is returned with the configured status code and content type from the spec.
func EnsureEverythingIsReferencedJSONDefaultResponse(body struct {
	Field SchemaObject `json:"Field"`
}) *Response {
	return &Response{
		body:        body,
		StatusCode:  200,
		contentType: "application/json",
	}
}

// Getter for additional properties for ParamsWithAddPropsParams_P1. Returns the specified
// element and whether it was found
func (a ParamsWithAddPropsParams_P1) Get(fieldName string) (value interface{}, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for ParamsWithAddPropsParams_P1
func (a *ParamsWithAddPropsParams_P1) Set(fieldName string, value interface{}) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]interface{})
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for ParamsWithAddPropsParams_P1 to handle AdditionalProperties
func (a *ParamsWithAddPropsParams_P1) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]interface{})
		for fieldName, fieldBuf := range object {
			var fieldVal interface{}
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for ParamsWithAddPropsParams_P1 to handle AdditionalProperties
func (a ParamsWithAddPropsParams_P1) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}

// Getter for additional properties for ParamsWithAddPropsParams_P2_Inner. Returns the specified
// element and whether it was found
func (a ParamsWithAddPropsParams_P2_Inner) Get(fieldName string) (value string, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for ParamsWithAddPropsParams_P2_Inner
func (a *ParamsWithAddPropsParams_P2_Inner) Set(fieldName string, value string) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]string)
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for ParamsWithAddPropsParams_P2_Inner to handle AdditionalProperties
func (a *ParamsWithAddPropsParams_P2_Inner) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]string)
		for fieldName, fieldBuf := range object {
			var fieldVal string
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for ParamsWithAddPropsParams_P2_Inner to handle AdditionalProperties
func (a ParamsWithAddPropsParams_P2_Inner) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}

// Getter for additional properties for BodyWithAddPropsJSONBody. Returns the specified
// element and whether it was found
func (a BodyWithAddPropsJSONBody) Get(fieldName string) (value interface{}, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for BodyWithAddPropsJSONBody
func (a *BodyWithAddPropsJSONBody) Set(fieldName string, value interface{}) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]interface{})
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for BodyWithAddPropsJSONBody to handle AdditionalProperties
func (a *BodyWithAddPropsJSONBody) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if raw, found := object["inner"]; found {
		err = json.Unmarshal(raw, &a.Inner)
		if err != nil {
			return fmt.Errorf("error reading 'inner': %w", err)
		}
		delete(object, "inner")
	}

	if raw, found := object["name"]; found {
		err = json.Unmarshal(raw, &a.Name)
		if err != nil {
			return fmt.Errorf("error reading 'name': %w", err)
		}
		delete(object, "name")
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]interface{})
		for fieldName, fieldBuf := range object {
			var fieldVal interface{}
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for BodyWithAddPropsJSONBody to handle AdditionalProperties
func (a BodyWithAddPropsJSONBody) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	object["inner"], err = json.Marshal(a.Inner)
	if err != nil {
		return nil, fmt.Errorf("error marshaling 'inner': %w", err)
	}

	object["name"], err = json.Marshal(a.Name)
	if err != nil {
		return nil, fmt.Errorf("error marshaling 'name': %w", err)
	}

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}

// Getter for additional properties for BodyWithAddPropsJSONBody_Inner. Returns the specified
// element and whether it was found
func (a BodyWithAddPropsJSONBody_Inner) Get(fieldName string) (value int, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for BodyWithAddPropsJSONBody_Inner
func (a *BodyWithAddPropsJSONBody_Inner) Set(fieldName string, value int) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]int)
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for BodyWithAddPropsJSONBody_Inner to handle AdditionalProperties
func (a *BodyWithAddPropsJSONBody_Inner) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]int)
		for fieldName, fieldBuf := range object {
			var fieldVal int
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for BodyWithAddPropsJSONBody_Inner to handle AdditionalProperties
func (a BodyWithAddPropsJSONBody_Inner) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}

// Getter for additional properties for AdditionalPropertiesObject1. Returns the specified
// element and whether it was found
func (a AdditionalPropertiesObject1) Get(fieldName string) (value int, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for AdditionalPropertiesObject1
func (a *AdditionalPropertiesObject1) Set(fieldName string, value int) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]int)
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for AdditionalPropertiesObject1 to handle AdditionalProperties
func (a *AdditionalPropertiesObject1) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if raw, found := object["id"]; found {
		err = json.Unmarshal(raw, &a.ID)
		if err != nil {
			return fmt.Errorf("error reading 'id': %w", err)
		}
		delete(object, "id")
	}

	if raw, found := object["name"]; found {
		err = json.Unmarshal(raw, &a.Name)
		if err != nil {
			return fmt.Errorf("error reading 'name': %w", err)
		}
		delete(object, "name")
	}

	if raw, found := object["optional"]; found {
		err = json.Unmarshal(raw, &a.Optional)
		if err != nil {
			return fmt.Errorf("error reading 'optional': %w", err)
		}
		delete(object, "optional")
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]int)
		for fieldName, fieldBuf := range object {
			var fieldVal int
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for AdditionalPropertiesObject1 to handle AdditionalProperties
func (a AdditionalPropertiesObject1) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	object["id"], err = json.Marshal(a.ID)
	if err != nil {
		return nil, fmt.Errorf("error marshaling 'id': %w", err)
	}

	object["name"], err = json.Marshal(a.Name)
	if err != nil {
		return nil, fmt.Errorf("error marshaling 'name': %w", err)
	}

	if a.Optional != nil {
		object["optional"], err = json.Marshal(a.Optional)
		if err != nil {
			return nil, fmt.Errorf("error marshaling 'optional': %w", err)
		}
	}

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}

// Getter for additional properties for AdditionalPropertiesObject3. Returns the specified
// element and whether it was found
func (a AdditionalPropertiesObject3) Get(fieldName string) (value interface{}, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for AdditionalPropertiesObject3
func (a *AdditionalPropertiesObject3) Set(fieldName string, value interface{}) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]interface{})
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for AdditionalPropertiesObject3 to handle AdditionalProperties
func (a *AdditionalPropertiesObject3) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if raw, found := object["name"]; found {
		err = json.Unmarshal(raw, &a.Name)
		if err != nil {
			return fmt.Errorf("error reading 'name': %w", err)
		}
		delete(object, "name")
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]interface{})
		for fieldName, fieldBuf := range object {
			var fieldVal interface{}
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for AdditionalPropertiesObject3 to handle AdditionalProperties
func (a AdditionalPropertiesObject3) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	object["name"], err = json.Marshal(a.Name)
	if err != nil {
		return nil, fmt.Errorf("error marshaling 'name': %w", err)
	}

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}

// Getter for additional properties for AdditionalPropertiesObject4. Returns the specified
// element and whether it was found
func (a AdditionalPropertiesObject4) Get(fieldName string) (value interface{}, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for AdditionalPropertiesObject4
func (a *AdditionalPropertiesObject4) Set(fieldName string, value interface{}) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]interface{})
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for AdditionalPropertiesObject4 to handle AdditionalProperties
func (a *AdditionalPropertiesObject4) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if raw, found := object["inner"]; found {
		err = json.Unmarshal(raw, &a.Inner)
		if err != nil {
			return fmt.Errorf("error reading 'inner': %w", err)
		}
		delete(object, "inner")
	}

	if raw, found := object["name"]; found {
		err = json.Unmarshal(raw, &a.Name)
		if err != nil {
			return fmt.Errorf("error reading 'name': %w", err)
		}
		delete(object, "name")
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]interface{})
		for fieldName, fieldBuf := range object {
			var fieldVal interface{}
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for AdditionalPropertiesObject4 to handle AdditionalProperties
func (a AdditionalPropertiesObject4) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	object["inner"], err = json.Marshal(a.Inner)
	if err != nil {
		return nil, fmt.Errorf("error marshaling 'inner': %w", err)
	}

	object["name"], err = json.Marshal(a.Name)
	if err != nil {
		return nil, fmt.Errorf("error marshaling 'name': %w", err)
	}

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}

// Getter for additional properties for AdditionalPropertiesObject4_Inner. Returns the specified
// element and whether it was found
func (a AdditionalPropertiesObject4_Inner) Get(fieldName string) (value interface{}, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for AdditionalPropertiesObject4_Inner
func (a *AdditionalPropertiesObject4_Inner) Set(fieldName string, value interface{}) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]interface{})
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for AdditionalPropertiesObject4_Inner to handle AdditionalProperties
func (a *AdditionalPropertiesObject4_Inner) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if raw, found := object["name"]; found {
		err = json.Unmarshal(raw, &a.Name)
		if err != nil {
			return fmt.Errorf("error reading 'name': %w", err)
		}
		delete(object, "name")
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]interface{})
		for fieldName, fieldBuf := range object {
			var fieldVal interface{}
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for AdditionalPropertiesObject4_Inner to handle AdditionalProperties
func (a AdditionalPropertiesObject4_Inner) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	object["name"], err = json.Marshal(a.Name)
	if err != nil {
		return nil, fmt.Errorf("error marshaling 'name': %w", err)
	}

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}

// Getter for additional properties for AdditionalPropertiesObject5. Returns the specified
// element and whether it was found
func (a AdditionalPropertiesObject5) Get(fieldName string) (value SchemaObject, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for AdditionalPropertiesObject5
func (a *AdditionalPropertiesObject5) Set(fieldName string, value SchemaObject) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]SchemaObject)
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for AdditionalPropertiesObject5 to handle AdditionalProperties
func (a *AdditionalPropertiesObject5) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]SchemaObject)
		for fieldName, fieldBuf := range object {
			var fieldVal SchemaObject
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for AdditionalPropertiesObject5 to handle AdditionalProperties
func (a AdditionalPropertiesObject5) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /ensure-everything-is-referenced)
	EnsureEverythingIsReferenced(w http.ResponseWriter, r *http.Request) *Response

	// (GET /params_with_add_props)
	ParamsWithAddProps(w http.ResponseWriter, r *http.Request, params ParamsWithAddPropsParams) *Response

	// (POST /params_with_add_props)
	BodyWithAddProps(w http.ResponseWriter, r *http.Request) *Response
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler          ServerInterface
	Middlewares      map[string]func(http.Handler) http.Handler
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// EnsureEverythingIsReferenced operation middleware
func (siw *ServerInterfaceWrapper) EnsureEverythingIsReferenced(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.EnsureEverythingIsReferenced(w, r)
		if resp != nil {
			render.Render(w, r, resp)
		}
	})

	handler(w, r.WithContext(ctx))
}

// ParamsWithAddProps operation middleware
func (siw *ServerInterfaceWrapper) ParamsWithAddProps(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parameter object where we will unmarshal all parameters from the context
	var params ParamsWithAddPropsParams

	// ------------- Required query parameter "p1" -------------

	if err := runtime.BindQueryParameter("simple", true, true, "p1", r.URL.Query(), &params.P1); err != nil {
		err = fmt.Errorf("invalid format for parameter p1: %w", err)
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "p1"})
		return
	}

	// ------------- Required query parameter "p2" -------------

	if err := runtime.BindQueryParameter("form", true, true, "p2", r.URL.Query(), &params.P2); err != nil {
		err = fmt.Errorf("invalid format for parameter p2: %w", err)
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "p2"})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.ParamsWithAddProps(w, r, params)
		if resp != nil {
			render.Render(w, r, resp)
		}
	})

	handler(w, r.WithContext(ctx))
}

// BodyWithAddProps operation middleware
func (siw *ServerInterfaceWrapper) BodyWithAddProps(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.BodyWithAddProps(w, r)
		if resp != nil {
			render.Render(w, r, resp)
		}
	})

	handler(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	error
}

type UnmarshalingParamError struct {
	error
	paramName string
}

type RequiredParamError struct {
	error
	paramName string
}

type RequiredHeaderError struct {
	error
	paramName string
}

type InvalidParamFormatError struct {
	error
	paramName string
}

type TooManyValuesForParamError struct {
	error
	paramName string
}

// ParameterName is an interface that is implemented by error types that are
// relevant to a specific parameter.
type ParameterError interface {
	error
	// ParamName is the name of the parameter that the error is referring to.
	ParamName() string
}

func (err UnmarshalingParamError) ParamName() string     { return err.paramName }
func (err RequiredParamError) ParamName() string         { return err.paramName }
func (err RequiredHeaderError) ParamName() string        { return err.paramName }
func (err InvalidParamFormatError) ParamName() string    { return err.paramName }
func (err TooManyValuesForParamError) ParamName() string { return err.paramName }

type ServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      map[string]func(http.Handler) http.Handler
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

type ServerOption func(*ServerOptions)

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface, opts ...ServerOption) http.Handler {
	options := &ServerOptions{
		BaseURL:     "/",
		BaseRouter:  chi.NewRouter(),
		Middlewares: make(map[string]func(http.Handler) http.Handler),
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
	}

	for _, f := range opts {
		f(options)
	}

	r := options.BaseRouter
	wrapper := ServerInterfaceWrapper{
		Handler:          si,
		Middlewares:      options.Middlewares,
		ErrorHandlerFunc: options.ErrorHandlerFunc,
	}

	r.Route(options.BaseURL, func(r chi.Router) {
		r.Get("/ensure-everything-is-referenced", wrapper.EnsureEverythingIsReferenced)
		r.Get("/params_with_add_props", wrapper.ParamsWithAddProps)
		r.Post("/params_with_add_props", wrapper.BodyWithAddProps)

	})
	return r
}

func WithRouter(r chi.Router) ServerOption {
	return func(s *ServerOptions) {
		s.BaseRouter = r
	}
}

func WithServerBaseURL(url string) ServerOption {
	return func(s *ServerOptions) {
		s.BaseURL = url
	}
}

func WithMiddleware(key string, middleware func(http.Handler) http.Handler) ServerOption {
	return func(s *ServerOptions) {
		s.Middlewares[key] = middleware
	}
}

func WithMiddlewares(middlewares map[string]func(http.Handler) http.Handler) ServerOption {
	return func(s *ServerOptions) {
		s.Middlewares = middlewares
	}
}

func WithErrorHandler(handler func(w http.ResponseWriter, r *http.Request, err error)) ServerOption {
	return func(s *ServerOptions) {
		s.ErrorHandlerFunc = handler
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/9RYT2/jthP9KgR/v6Mcx872oluKbtsUaBtsgvawCQJGHFncykMtScURFvruxVCyJduU",
	"VnZy6V5Wljj/HmceH/ONJ3pdaAR0lsffuIGvJVj3o5YK/ItPuxcV/Uw0OkBHj6IocpUIpzTOv1iN9M4m",
	"GawFPRVGF2Bc6+VnBbmkh/8bSHnM/zfvws4bIzu/8///+fwFEsfrOvLJKAOSx59bD4/02sGrmxe5UAch",
	"XVUAj7l1RuGK1/SPfNhCo90W0/xoY/zX6om4BJsYVVCOPObXzKp1kQPbFsl0F6zNghxdS6nIROS3uyqa",
	"tBa+8MDnXnyFDlZg+FH4X4VlnS3rEGI6ZWTMFDoeHUCnZNg3ijUEqo64LpoAIUj2MfUuIorwGG2XbhGJ",
	"RlBYDqOQitzCYeE/abAMtWMiz/UmjMFb636n0q6GS3OmPKrsmgqyTGAVqKo6qumE3E9L+8NpaftORI3V",
	"WpeWpTRabJOpJGPZUI8e7w8imO+Ffdfyp5k3eUXnoPjD2HRPZ67pc79RLmONE5Zqw6RK/CLTAH6U+i/6",
	"zhd9L1aU0n6Ye7DOstfZSs9aaEZHaremb9Fum3/Vjd57jV0D0N/KZb9ZjbszYVKTRPxF5CV4Ak61WQvH",
	"Y+6PnWhg6XLC0nD6baRgCS25difi4CZseXjm3TGNzGXANpnOd+fO4QZZZxbB2q0zy6EPT0XQqt3F/Sxa",
	"gu4Mlyca1gNfqW968PzVvJ2MzkQkxoPbU8rxNnVgh/dm+ag7U2Ws+2OoRY3OJwyJXxX1XD16raIw1WSc",
	"qwTQQjcL/Pebe/LulCP3HkV2B+bF89wLGNugu7i4vLhsFACgKBSP+dXF5cWCwBUu8/nPAW1pYAYvYCqX",
	"KVzNlJ0ZSMEAJuDncQWhxs6UZYCy0Aodg1dFO2k1c5lwrKNElghkz8ASA8KBZIq6XtkHtAUkTKD0OuAZ",
	"WGFKBPlAM0n4ehl5I3nMP/oEP+7yu7GfuuyinuCuhlh5T5PP+4L8UN8uLy/fIGpT9QLfOxnGDps64qku",
	"zfkuPpCL1f6ZMOapf3zUkafEScI8xNvUZviG8hd9qdpN22ga+6t79ju6mWLeLKZ5ygy8oYIr72Ojz/ew",
	"5AEGau8tqShzN9zhbRPPD25ojfW8EEas7RPJiych5RP1rR0c7WtG9NCIEW8JDoz1wyrYs5ZVqw1bDgtr",
	"mcAk3/osqG2upbz1KRATbQPw+HOQZHYrhsWoD0bXP/61BFNt1V7MiwXvc22jZrr5HZOqRweBdZWn2+bO",
	"yOtoSrbY09Veie4uAy2ICCAtc5o9wwO60qAnSaeZaFc2N0FSg6FsyXKjzT/DCCxHEThJwwdOuMNmDWnv",
	"x7p+3CNaLPO8jnihbaD7vDpmLWf3241YWSikr1IZSFwQkIj69AFHgafGDtkGepaOiYOONWf+RWf6vWjq",
	"NvRuwWdejrb6fLtP9WGv1McbV9f1vwEAAP//YqjbhPYSAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
