package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/segp/agents-main/knowledge"
	"github.com/segp/agents-main/llm"
	"github.com/segp/agents-main/memory"
	"github.com/segp/agents-main/storage"
	"github.com/segp/agents-main/tools"
)

const (
	thinkPrompt = `%s
You are tasked with answering the following query:
Query: "%s"

Your goal is to reason about the query and decide on the best course of action to answer it accurately. You must answer using chain of thought reasoning, so explain your though process before giving the answer. You should start with ## Thoughts 
and give your thoughts in a detailed manner. Then you should give your answer in 
## Answer where you give the answer to the query.

Previous reasoning steps and observations: "%+v"

Available tools: "%s"

Instructions:
1. Analyze the query, previous reasoning steps, and observations.
2. Decide on the next action: use a tool or provide a final answer.


If you need to use a tool: 
1. Give a detailed reasoning about what to do next
2. Choose the tool that is most likely to answer the query
3. Give a reason for choosing the tool
4. Give the input for the tool
5. You must only give the next ONE tool call, and nothing else.

If instead you have enough information to answer the query:
1. Give a detailed reasoning process
2. Give the answer to the query


Remember:
- Use chain of thought reasoning
- Be thorough in your reasoning.
- Use tools when you need more information.
- If you use a tool, you must only give the next ONE tool call, and nothing else.
- Always base your reasoning on the actual observations from tool use.
- Provide a final answer only when you're confident you have sufficient information.
- If you cannot find the necessary information after using available tools, admit that you don't have enough information to answer the query confidently.`

	decidePrompt = `Given the following thought, choose the correct tool to use: %s`

	observerPrompt = `Previously, you thought %s
	
	You then decided to call tool %+v
	For which this was the result: %+v
	
	Now,describe thoroughly what can be observed from the result of the tool call, relevant to the given task.`
)

type ReActAgent struct {
	Background  string
	ToolHandler *tools.ToolHandler
	LLM         llm.LLM
	Memory      memory.Memory[memory.ReActMemorySteps]
	Storage     storage.Storage
	Knowledge   knowledge.Knowledge
}

func NewReActAgent(background string, toolHandler *tools.ToolHandler, llm llm.LLM, memory memory.Memory[memory.ReActMemorySteps], storage storage.Storage, knowledge knowledge.Knowledge) *ReActAgent {
	// TODO: add check for no_tool tool

	return &ReActAgent{
		Background:  background,
		ToolHandler: toolHandler,
		LLM:         llm,
		Memory:      memory,
		Storage:     storage,
		Knowledge:   knowledge,
	}
}

func (a *ReActAgent) Run(query string) (*string, *string, error) {
	return a.logicLoop(query)
}

func (a *ReActAgent) logicLoop(query string) (*string, *string, error) {

	for {
		mem, err := a.Memory.Get()
		if err != nil {
			return nil, nil, err
		}

		thoughts, err := a.think(mem, a.ToolHandler.ToolDefinitions(), query)
		if err != nil {
			return nil, nil, err
		}

		toolCall, err := a.decide(*thoughts)
		if err != nil {
			return nil, nil, err
		}

		if toolCall.Name == "no_tool" {
			return extractAnswerAndReason(toolCall)
		}

		toolCallResult, err := a.act(*toolCall)
		if err != nil {
			return nil, nil, err
		}

		observation, err := a.observe(toolCallResult, thoughts, *toolCall)
		if err != nil {
			return nil, nil, err
		}

		action := fmt.Sprintf("%+v", toolCall)
		a.Memory.Add(memory.ReActMemorySteps{
			Thought:     *thoughts,
			Action:      action,
			Observation: *observation,
		})
	}

}

func (a *ReActAgent) think(mem []memory.ReActMemorySteps, toolDefinitions []tools.ToolDefinition, query string) (*string, error) {
	prompt := a.getThinkPrompt(mem, toolDefinitions, query)
	if prompt == nil {
		return nil, errors.New("failed to get think prompt")
	}
	return a.LLM.ChatCompletion(context.Background(), *prompt)
}

func (a *ReActAgent) decide(thought string) (*tools.ToolCall, error) {
	prompt := fmt.Sprintf(decidePrompt, thought)

	toolChoice := tools.ToolChoice{
		Type: tools.ToolChoiceTypeAuto,
	}

	chosenTool, err := a.LLM.ChatCompletionWithTools(context.Background(), prompt, a.ToolHandler.ToolDefinitions(), toolChoice)
	if err != nil {
		return nil, err
	}

	if len(chosenTool) == 0 {
		return nil, errors.New("no tool chosen")
	}

	return &chosenTool[0], nil
}

func (a *ReActAgent) act(toolCall tools.ToolCall) (*string, error) {
	return a.ToolHandler.Call(toolCall)
}

func (a *ReActAgent) observe(toolCallResult *string, thoughts *string, chosenToolCall tools.ToolCall) (*string, error) {
	prompt := fmt.Sprintf(observerPrompt, *thoughts, chosenToolCall, *toolCallResult)

	return a.LLM.ChatCompletion(context.Background(), prompt)
}

func (a *ReActAgent) getThinkPrompt(mem []memory.ReActMemorySteps, toolDefinitions []tools.ToolDefinition, query string) *string {
	memoryStr, err := json.Marshal(mem)
	if err != nil {
		return nil
	}

	toolsStr := ""
	for _, tool := range toolDefinitions {
		jsonTool, err := json.Marshal(tool)
		if err != nil {
			return nil
		}
		toolsStr += string(jsonTool) + ", "
	}

	prompt := fmt.Sprintf(thinkPrompt, a.Background, query, memoryStr, toolsStr)
	return &prompt
}

func extractAnswerAndReason(toolCall *tools.ToolCall) (*string, *string, error) {
	var arguments map[string]string
	if err := json.Unmarshal([]byte(toolCall.Arguments), &arguments); err != nil {
		return nil, nil, err
	}

	reason := arguments["reason"]
	answer := arguments["answer"]

	return &answer, &reason, nil
}
