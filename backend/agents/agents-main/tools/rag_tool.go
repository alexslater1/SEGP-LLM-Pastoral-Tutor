package tools

import (
	"github.com/segp/agents-main/knowledge"
)

type RagTool struct{
	ragKnowledge knowledge.Knowledge
}

func NewRagTool(ragKnowledge knowledge.Knowledge) *RagTool {
	return &RagTool{
		ragKnowledge: ragKnowledge,
	}
}

func (r *RagTool) Definition() ToolDefinition {
	return ToolDefinition{
		Name:        "rag_tool",
		Description: "Get more relevant context and/or important links/contact information about the given query in respect to Imperial College London. Uses RAG",
		Parameters: []Parameter{
			{Name: "query", Type: "string", Description: "The query to send to the RAG model."},
		},
	}
}

func (r *RagTool) SearchRagFor(query string) (*string, error) {
	result, err := r.ragKnowledge.Get(query)

	if err != nil {
		return nil, err
	}

	return result, nil
}