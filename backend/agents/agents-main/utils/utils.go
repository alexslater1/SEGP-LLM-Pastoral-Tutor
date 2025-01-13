package utils

import (
	"encoding/json"
	"fmt"
)

func SchemaMapFrom(schema interface{}) (map[string]interface{}, error) {
	var jsonStr string

	// Convert schema to JSON string
	jsonBytes, err := json.Marshal(schema)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal schema: %v", err)
	}
	jsonStr = string(jsonBytes)

	// Convert JSON string to map
	var schemaMap map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &schemaMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schema: %v", err)
	}

	return schemaMap, nil
}
