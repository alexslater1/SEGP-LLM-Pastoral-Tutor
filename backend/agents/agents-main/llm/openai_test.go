package llm

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestOpenAiLLM(t *testing.T) {
	jsonStr := `
	{
		"name": "John Doe",
		"age": 30,
		"address": {
			"city": "New York",
			"zipcode": "10001"
		},
		"hobbies": ["reading", "travelling"],
		"attributes": {
			"height": 5.9,
			"weight": 70,
			"skills": {
				"coding": true,
				"languages": ["Go", "Python"]
			}
		}
	}`

	// Unmarshal the JSON into a generic map
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	// Iterate over the JSON
	iterateJSON(data, "")
}
