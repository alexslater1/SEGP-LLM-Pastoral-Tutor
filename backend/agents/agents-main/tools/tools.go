package tools

type Tool interface {
	Definition() ToolDefinition
}

type ParameterType string

const (
	ParameterTypeString  ParameterType = "string"
	ParameterTypeNumber  ParameterType = "number"
	ParameterTypeBoolean ParameterType = "boolean"
	ParameterTypeArray   ParameterType = "array"
)

type ToolChoiceType string

const (
	ToolChoiceTypeAuto      ToolChoiceType = "auto"
	ToolChoiceTypeRequired  ToolChoiceType = "required"
	ToolChoiceTypeForcedOne ToolChoiceType = "forced_one"
)

type Parameter struct {
	Name        string
	Description string
	Type        ParameterType
}

type ToolDefinition struct {
	Name        string
	Description string
	Parameters  []Parameter
}

type ToolChoice struct {
	Type         ToolChoiceType
	FunctionName string
}

type ToolCall struct {
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments,omitempty"`
}

func CheckWeatherTool() ToolDefinition {
	return ToolDefinition{
		Name:        "check_weather",
		Description: "Check the weather in a given location",
		Parameters: []Parameter{
			{Name: "location", Description: "The location to check the weather for", Type: ParameterTypeString},
			{Name: "date", Description: "The date to check the weather for", Type: ParameterTypeString},
		},
	}
}

func CheckBestAnimalNameTool() ToolDefinition {
	return ToolDefinition{
		Name:        "check_best_animal_name",
		Description: "Check the best animal name for a given animal",
		Parameters: []Parameter{
			{Name: "animal", Description: "The animal to check the best name for", Type: ParameterTypeString},
		},
	}
}

// Always need to pass this
func NoTool() ToolDefinition {
	return ToolDefinition{
		Name:        "no_tool",
		Description: "Do not use any tools",
		Parameters: []Parameter{
			{Name: "reason", Description: "The reason why no tool was used", Type: ParameterTypeString},
			{Name: "answer", Description: "The answer to the question", Type: ParameterTypeString},
		},
	}
}
