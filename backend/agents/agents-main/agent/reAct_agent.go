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
	thinkPrompt = `You are a ReAct (Reasoning and Acting) agent tasked with answering the following query:

Query: "%s"

Your goal is to reason about the query and decide on the best course of action to answer it accurately. You must answer using chain of thought reasoning, so explain your though process before giving the answer. You should start with ## Thoughts 
and give your thoughts in a detailed manner. Then you should give your answer in 
## Answer where you give the answer to the query.

Previous reasoning steps and observations: "%s"

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
)

type ReActAgent struct {
	Description string
	Tools       map[string]tools.Tool
	LLM         llm.LLM
	Memory      memory.Memory[memory.ReActMemorySteps]
	Storage     storage.Storage
	Knowledge   knowledge.Knowledge
}

func NewReActAgent(description string, ts []tools.Tool, llm llm.LLM, memory memory.Memory[memory.ReActMemorySteps], storage storage.Storage, knowledge knowledge.Knowledge) *ReActAgent {
	// TODO: add check for no_tool tool

	toolsMap := make(map[string]tools.Tool)
	for _, tool := range ts {
		toolsMap[tool.Definition().Name] = tool
	}

	return &ReActAgent{
		Description: description,
		Tools:       toolsMap,
		LLM:         llm,
		Memory:      memory,
		Storage:     storage,
		Knowledge:   knowledge,
	}
}

func (a *ReActAgent) Run(input string) (string, error) {
	return "", nil
}

func (a *ReActAgent) think() (*string, error) {
	prompt := a.getThinkPrompt()
	if prompt == nil {
		return nil, errors.New("failed to get think prompt")
	}
	return a.LLM.ChatCompletion(context.Background(), *prompt)
}

func (a *ReActAgent) decide(thought string) (*tools.ToolCall, error) {
	prompt := fmt.Sprintf(decidePrompt, thought)

	toolDefinitions := []tools.ToolDefinition{}
	for _, tool := range a.Tools {
		toolDefinitions = append(toolDefinitions, tool.Definition())
	}

	toolChoice := tools.ToolChoice{
		Type: tools.ToolChoiceTypeAuto,
	}

	chosenTool, err := a.LLM.ChatCompletionWithTools(context.Background(), prompt, toolDefinitions, toolChoice)
	if err != nil {
		return nil, err
	}

	if len(chosenTool) == 0 {
		return nil, errors.New("no tool chosen")
	}

	return &chosenTool[0], nil
}

func (a *ReActAgent) act(toolCall tools.ToolCall) (*string, error) {
	tool := a.Tools[toolCall.Name]

	switch tool.(type) {
	case *tools.GoogleSearchResultsTool:
		//todo:
	}

	return nil, nil
}

func (a *ReActAgent) getThinkPrompt() *string {
	mem, err := a.Memory.Get()
	if err != nil {
		return nil
	}

	memoryStr, err := json.Marshal(mem)
	if err != nil {
		return nil
	}

	toolsStr := ""
	for _, tool := range a.Tools {
		jsonTool, err := json.Marshal(tool.Definition())
		if err != nil {
			return nil
		}
		toolsStr += string(jsonTool) + ", "
	}

	prompt := fmt.Sprintf(thinkPrompt, a.Description, string(memoryStr), toolsStr)
	return &prompt
}
