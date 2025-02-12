package tools

type NoTool struct {
}

func NewNoTool() *NoTool {
	return &NoTool{}
}

func (n *NoTool) Definition() ToolDefinition {
	return ToolDefinition{
		Name:        "no_tool",
		Description: "Do not use any tools. This could be because you have an answer, you have a question for the user which you require an answer from the user in order to answer the user's query, or you deem that after sufficient attempts, it will not be possible to feasibly find an accurate answer.",
		Parameters: []Parameter{
			{Name: "reason", Description: "The reason why no tool was used", Type: ParameterTypeString},
			{Name: "answer", Description: "The answer to the question / reason why not possible to answer the question / your question to the user", Type: ParameterTypeString},
		},
	}
}
