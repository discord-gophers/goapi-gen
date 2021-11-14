package codegen

import (
	"encoding/json"
	"fmt"
)

const (
	extPropGoType         = "x-go-type"
	extPropGoTypeExternal = "x-go-type-external"
	extPropOmitEmpty      = "x-omitempty"
	extPropExtraTags      = "x-go-extra-tags"
	extMiddlewares        = "x-go-middlewares"
)

type extImportPathDetails struct {
	Import string `json:"import"`
	Alias  string `json:"alias"`
	Type   string `json:"type"`
}

func extTypeName(extPropValue interface{}) (extImportPathDetails, error) {
	var details extImportPathDetails
	raw, ok := extPropValue.(json.RawMessage)
	if !ok {
		return details, fmt.Errorf("failed to convert type: %T", extPropValue)
	}
	var name string
	if err := json.Unmarshal(raw, &name); err == nil {
		details.Type = name
		return details, nil
	}

	return extImportPath(raw)
}

func extImportPath(raw json.RawMessage) (extImportPathDetails, error) {
	var details extImportPathDetails
	if err := json.Unmarshal(raw, &details); err != nil {
		return extImportPathDetails{}, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return details, nil
}

func extParseOmitEmpty(extPropValue interface{}) (bool, error) {
	raw, ok := extPropValue.(json.RawMessage)
	if !ok {
		return false, fmt.Errorf("failed to convert type: %T", extPropValue)
	}

	var omitEmpty bool
	if err := json.Unmarshal(raw, &omitEmpty); err != nil {
		return false, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return omitEmpty, nil
}

func extExtraTags(extPropValue interface{}) (map[string]string, error) {
	raw, ok := extPropValue.(json.RawMessage)
	if !ok {
		return nil, fmt.Errorf("failed to convert type: %T", extPropValue)
	}
	var tags map[string]string
	if err := json.Unmarshal(raw, &tags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}
	return tags, nil
}

func extParseMiddlewares(extPropValue interface{}) ([]string, error) {
	raw, ok := extPropValue.(json.RawMessage)
	if !ok {
		return nil, fmt.Errorf("failed to convert type: %T", extPropValue)
	}
	var middlewares []string
	if err := json.Unmarshal(raw, &middlewares); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}
	return middlewares, nil
}
