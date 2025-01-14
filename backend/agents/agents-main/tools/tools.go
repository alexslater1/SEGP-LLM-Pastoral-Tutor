package tools

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

type Tool struct {
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

func CheckWeatherTool() Tool {
	return Tool{
		Name:        "check_weather",
		Description: "Check the weather in a given location",
		Parameters: []Parameter{
			{Name: "location", Description: "The location to check the weather for", Type: ParameterTypeString},
			{Name: "date", Description: "The date to check the weather for", Type: ParameterTypeString},
		},
	}
}
