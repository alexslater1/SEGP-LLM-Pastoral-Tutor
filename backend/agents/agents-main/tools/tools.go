package tools

type ParameterType string

const (
	ParameterTypeString  ParameterType = "string"
	ParameterTypeNumber  ParameterType = "number"
	ParameterTypeBoolean ParameterType = "boolean"
	ParameterTypeArray   ParameterType = "array"
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

func CheckWeatherTool(tool Tool) Tool {
	return Tool{
		Name:        "Check Weather",
		Description: "Check the weather in a given location",
		Parameters: []Parameter{
			{Name: "location", Description: "The location to check the weather for", Type: ParameterTypeString},
			{Name: "date", Description: "The date to check the weather for", Type: ParameterTypeString},
		},
	}
}
