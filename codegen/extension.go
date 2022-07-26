package codegen

import (
	"encoding/json"
	"fmt"
)

const (
	extPropGoType        = "x-go-type"
	extPropOmitEmpty     = "x-omitempty"
	extPropExtraTags     = "x-go-extra-tags"
	extPropOptionalValue = "x-go-optional-value"
	extPropString        = "x-go-string"
	extMiddlewares       = "x-go-middlewares"
)

func extTypeName(extPropValue interface{}) (string, error) {
	var name string
	err := extParseAny(extPropValue, &name)
	return name, err
}

func extExtraTags(extPropValue interface{}) (map[string]string, error) {
	var tags map[string]string
	err := extParseAny(extPropValue, &tags)
	return tags, err
}

func extParseMiddlewares(extPropValue interface{}) ([]string, error) {
	var middlewares []string
	err := extParseAny(extPropValue, &middlewares)
	return middlewares, err
}

func extParseBool(extPropValue interface{}) (bool, error) {
	var b bool
	err := extParseAny(extPropValue, &b)
	return b, err
}

func extParseAny(extPropValue, target interface{}) error {
	raw, ok := extPropValue.(json.RawMessage)
	if !ok {
		return fmt.Errorf("failed to convert type: %T", extPropValue)
	}

	if err := json.Unmarshal(raw, target); err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return nil
}
