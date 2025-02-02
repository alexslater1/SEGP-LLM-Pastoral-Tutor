package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/segp/agents-main/knowledge"
	"github.com/segp/agents-main/llm"
	"github.com/segp/agents-main/tools"
)

type FastAgent struct {
	Background  string
	ToolHandler *tools.ToolHandler
	LLM         llm.LLM
	Knowledge   knowledge.Knowledge
}

func NewFastAgent(background string, toolHandler *tools.ToolHandler, llm llm.LLM, knowledge knowledge.Knowledge) *FastAgent {
	return &FastAgent{
		Background:  background,
		ToolHandler: toolHandler,
		LLM:         llm,
		Knowledge:   knowledge,
	}
}

func (a *FastAgent) Run(query string, requestId string) (*string, *string, error) {
	return a.logicLoop(query, requestId)
}

func (a *FastAgent) Subscribe() <-chan AgentEvent {
	panic("not implemented")
}

func (a *FastAgent) Unsubscribe(ch <-chan AgentEvent) {
	panic("not implemented")
}

func (a *FastAgent) logicLoop(query string, requestId string) (*string, *string, error) {
	slog.Info("Starting FAST AGENT logic loop for query", "query", query)

	var prevThoughts *string
	var prevToolCallResult *string
	var prevToolCall *tools.ToolCall

	for i := 0; i < 10; i++ {
		knowledgeContext, err := a.Knowledge.Get(query)
		if err != nil {
			return nil, nil, err
		}

		thoughts, toolChoice, err := a.thinkAndChooseTool(i == 0, query, knowledgeContext, prevThoughts, prevToolCall, prevToolCallResult, requestId)
		if err != nil {
			return nil, nil, err
		}

		if toolChoice.Name == "no_tool" {
			return a.extractAnswerAndReason(requestId, toolChoice)
		}

		result, err := a.ToolHandler.Call(*toolChoice)
		if err != nil {
			return nil, nil, err
		}

		prevThoughts = thoughts
		prevToolCallResult = result
		prevToolCall = toolChoice
	}

	return nil, nil, fmt.Errorf("failed to find answer after 10 iterations")
}

func (a *FastAgent) thinkAndChooseTool(isFirstIteration bool, query string, knowledgeContext *string, prevThoughts *string, prevToolCall *tools.ToolCall, prevToolCallResult *string, requestId string) (*string, *tools.ToolCall, error) {
	kc := ""
	if knowledgeContext != nil {
		kc = *knowledgeContext
	}

	prompt, err := a.thinkingAndActPrompt(isFirstIteration, query, &kc, prevThoughts, prevToolCall, prevToolCallResult, requestId)
	if err != nil {
		return nil, nil, err
	}

	toolCalls, err := a.LLM.ChatCompletionWithTools(context.TODO(), *prompt, a.ToolHandler.ToolDefinitionsWithThoughts(a.ToolHandler.ToolDefinitionsWithDescribingAction()), tools.ToolChoice{Type: tools.ToolChoiceTypeRequired})
	if err != nil {
		return nil, nil, err
	}

	if len(toolCalls) != 1 {
		return nil, nil, fmt.Errorf("expected 1 tool call, got %d", len(toolCalls))
	}

	var parsedArgs map[string]interface{}
	err = json.Unmarshal([]byte(toolCalls[0].Arguments), &parsedArgs)
	if err != nil {
		return nil, nil, err
	}

	if parsedArgs["_thoughts"] == nil {
		return nil, nil, fmt.Errorf("thoughts field is required")
	}

	thoughts, ok := parsedArgs["_thoughts"].(string)
	if !ok {
		return nil, nil, fmt.Errorf("conversion of thoughts field to string failed")
	}

	return &thoughts, &toolCalls[0], nil
}

func (a *FastAgent) thinkingAndActPrompt(isFirstIteration bool, query string, knowledgeContext *string, prevThoughts *string, prevToolCall *tools.ToolCall, prevToolCallResult *string, requestId string) (*string, error) {
	prompt := ""

	if isFirstIteration {
		prompt = fmt.Sprintf("You are a reAct agent. Your goal is to solve the following query: `%s`. Here is some (potentially relevant) knowledge from a rag source: `%s`.  ", query, *knowledgeContext)
	} else {
		prompt = fmt.Sprintf("You are a reAct agent, currently in the process of solving the query: `%s`. In the previous iteration, you thought `%s` and then called the tool `%s`. The results of this tool where `%s`.", query, *prevThoughts, *prevToolCall, *prevToolCallResult)
	}

	prompt += "Now, give some thoughts about what you already know, and then generate a plan (based on what you need to find out), of how to solve the problem."

	toolChoiceString, err := a.toolChoicesString()
	if err != nil {
		return nil, err
	}

	prompt += ` You have these tools at your disposal: ` + toolChoiceString + ` It is also essential that you give your thoughts in the _thoughts field. If you believe you already know the answer to the query, or that you will be unable to get the answer, pick the no_tool tool.`

	return &prompt, nil
}

func (a *FastAgent) toolChoicesString() (string, error) {
	availableTools := a.ToolHandler.ToolDefinitions()

	toolChoiceString, err := json.Marshal(availableTools)
	if err != nil {
		return "", err
	}

	return string(toolChoiceString), nil
}

// TODO: remove duplication of this
func (a *FastAgent) extractAnswerAndReason(requestId string, toolCall *tools.ToolCall) (*string, *string, error) {
	var arguments map[string]string
	if err := json.Unmarshal([]byte(toolCall.Arguments), &arguments); err != nil {
		return nil, nil, err
	}

	reason := arguments["reason"]
	answer := arguments["answer"]

	// a.publish(NewAnswerSuccessEvent(requestId, answer, reason))

	return &answer, &reason, nil
}
