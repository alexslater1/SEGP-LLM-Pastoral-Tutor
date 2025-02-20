package tools

type GiveAnswerTool struct {
}

func NewGiveAnswerTool() *GiveAnswerTool {
	return &GiveAnswerTool{}
}

func (n *GiveAnswerTool) Definition() ToolDefinition {
	return ToolDefinition{
		Name:        "give_answer",
		Description: "This tool is used to give an answer to the user's query. It is used when you have an answer, or you deem that after sufficient attempts, it will not be possible to feasibly find an accurate answer.",
		Parameters: []Parameter{
			{Name: "reason", Description: "The reason why give_answer tool was used", Type: ParameterTypeString},
			{Name: "answer", Description: "The answer to the question / reason why not possible to answer the question", Type: ParameterTypeString},
		},
	}
}
