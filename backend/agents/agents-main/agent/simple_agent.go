package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segp/agents-main/llm"
)

const (
	agentPrompt = "You must reply with an answer, and a reason which is your reason for giving this answer."
)

type SimpleAgent struct {
	systemPrompt string
	llm          llm.LLM
}

func NewSimpleAgent(systemPrompt string, llm llm.LLM) *SimpleAgent {
	return &SimpleAgent{systemPrompt: systemPrompt, llm: llm}
}

func (s *SimpleAgent) Run(query string, requestID string) (*string, *string, error) {

	type response struct {
		Answer string `json:"answer"`
		Reason string `json:"reason"`
	}

	resp, err := s.llm.StructuredOutputCompletion(context.Background(), promptFrom(s.systemPrompt, query), response{})
	if err != nil {
		return nil, nil, err
	}

	var parsedResponse response
	if err := json.Unmarshal([]byte(*resp), &parsedResponse); err != nil {
		return nil, nil, err
	}

	return &parsedResponse.Answer, &parsedResponse.Reason, nil
}

func (s *SimpleAgent) Subscribe() <-chan AgentEvent {
	panic("not implemented")
}

func (s *SimpleAgent) Unsubscribe(ch <-chan AgentEvent) {
	panic("not implemented")
}

func promptFrom(systemPrompt string, userPrompt string) string {
	return fmt.Sprintf("System prompt: %s\n\nUser query: %s\n\n%s", systemPrompt, userPrompt, agentPrompt)
}
