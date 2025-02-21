package tools

type NoToolTool struct {
}

func NewNoToolTool() *NoToolTool {
	return &NoToolTool{}
}

func (n *NoToolTool) Definition() ToolDefinition {
	return ToolDefinition{
		Name:        "no_tool",
		Description: "It is used when you have an answer, or you deem that after sufficient attempts, it will not be possible to feasibly find an accurate answer, OR you want to ask the user something.",
		Parameters: []Parameter{
			{Name: "reason", Description: "The reason why no_tool tool was used", Type: ParameterTypeString},
			{Name: "response", Description: "The question you want to ask the user / the reason why you are not able to answer the question / your answer to the question", Type: ParameterTypeString},
		},
	}
}
