package tools

type NoTool struct {
}

func NewNoTool() *NoTool {
	return &NoTool{}
}

func (n *NoTool) Definition() ToolDefinition {
	return ToolDefinition{
		Name:        "no_tool",
		Description: "Do not use any tools",
		Parameters: []Parameter{
			{Name: "reason", Description: "The reason why no tool was used", Type: ParameterTypeString},
			{Name: "answer", Description: "The answer to the question", Type: ParameterTypeString},
		},
	}
}
