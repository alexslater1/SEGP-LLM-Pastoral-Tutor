package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/segp/agents-main/clock"
	"github.com/segp/agents-main/context_keys"
	"github.com/segp/agents-main/history"
	"github.com/segp/agents-main/knowledge"
	"github.com/segp/agents-main/llm"
	"github.com/segp/agents-main/tools"
)

const (
	maxIterations = 10
)

type FastAgent struct {
	ID     string
	Desc   string
	Prompt string

	ToolHandler *tools.ToolHandler
	LLM         llm.LLM
	Knowledge   knowledge.Knowledge
	Clock       clock.Clock
	History     history.History

	subscribers []chan AgentEvent
}

func NewFastAgent(id string, description string, prompt string, toolHandler *tools.ToolHandler, llm llm.LLM, knowledge knowledge.Knowledge, clock clock.Clock, history history.History) *FastAgent {
	return &FastAgent{
		ID:     id,
		Desc:   description,
		Prompt: prompt,

		ToolHandler: toolHandler,
		LLM:         llm,
		Knowledge:   knowledge,
		Clock:       clock,
		History:     history,

		subscribers: []chan AgentEvent{},
	}
}

func (a *FastAgent) Id() string {
	return a.ID
}

func (a *FastAgent) Description() string {
	return a.Desc
}

func (a *FastAgent) Run(ctx context.Context, query string) (*AgentResponse, error) {
	a.publish(NewQueryEvent(ctx, query))
	response, err := a.logicLoop(ctx, query)
	if err != nil {
		a.publish(NewAnswerErrorEvent(ctx, err.Error()))
	}
	return response, err
}

func (a *FastAgent) Subscribe() <-chan AgentEvent {
	ch := make(chan AgentEvent, defaultSubscriberBufferSize)
	a.subscribers = append(a.subscribers, ch)
	return ch
}

func (a *FastAgent) Unsubscribe(ch <-chan AgentEvent) {
	for i, subscriber := range a.subscribers {
		if subscriber == ch {
			close(subscriber)
			a.subscribers = append(a.subscribers[:i], a.subscribers[i+1:]...)
			break
		}
	}
}

func (a *FastAgent) publish(event AgentEvent) {
	for _, subscriber := range a.subscribers {
		select {
		case subscriber <- event:
		default:
			slog.Warn("Subscriber buffer is full, skipping event", "event", event)
		}
	}
}

func (a *FastAgent) handleGiveAnswer(ctx context.Context, toolChoice *tools.ToolCall) (*AgentResponse, error) {
	answer, reason, err := a.extractAnswerAndReason(toolChoice)
	if err != nil {
		return nil, err
	}
	a.publish(NewAnswerSuccessEvent(ctx, *answer, *reason))
	return &AgentResponse{
		Answer: answer,
		Reason: reason,
	}, nil
}

func (a *FastAgent) handleOffloadTask(ctx context.Context, toolChoice *tools.ToolCall) (*AgentResponse, error) {
	entityId, task, err := a.extractEntityIdAndTask(toolChoice)
	if err != nil {
		return nil, err
	}
	a.publish(NewOffloadTaskEvent(ctx, *entityId, *task))
	return &AgentResponse{
		OffloadTask: &OffloadTask{
			EntityID: *entityId,
			Task:     *task,
		},
	}, nil
}

func (a *FastAgent) logicLoop(ctx context.Context, query string) (*AgentResponse, error) {
	slog.Info("Starting FAST AGENT logic loop for query", "query", query)

	var prevThoughts *string
	var prevToolCallResult *string
	var prevToolCall *tools.ToolCall

	var prevToolCalls []tools.ToolCall

	for i := 0; i < maxIterations; i++ {
		knowledgeContext, err := a.Knowledge.Get(query)
		if err != nil {
			return nil, err
		}

		thoughts, toolChoice, err := a.thinkAndChooseTool(ctx, i, query, knowledgeContext, prevThoughts, prevToolCall, prevToolCallResult, prevToolCalls)
		if err != nil {
			return nil, err
		}

		if toolChoice.Name == "give_answer" {
			return a.handleGiveAnswer(ctx, toolChoice)
		}

		if toolChoice.Name == "offload_task" {
			return a.handleOffloadTask(ctx, toolChoice)
		}

		result, err := a.ToolHandler.Call(*toolChoice)
		if err != nil {
			return nil, err
		}
		a.publish(NewToolCallResultEvent(ctx, result))

		prevThoughts = thoughts
		prevToolCallResult = result
		prevToolCall = toolChoice
		prevToolCalls = append(prevToolCalls, *toolChoice)
	}

	return nil, fmt.Errorf("failed to find answer after 10 iterations")
}

type StructuredOutput struct {
	Thoughts     string `json:"_thoughts"`
	ToolCallName string `json:"tool_call_name"`
	ToolCallArgs []struct {
		ToolCallArgName  string `json:"tool_call_arg_name"`
		ToolCallArgValue string `json:"tool_call_arg_value"`
	} `json:"tool_call_args"`
	DescriptionOfAction string `json:"description_of_action"`
}

func structuredOutputToToolCall(structuredOutputCompletion *string) *tools.ToolCall {
	var structuredOutput StructuredOutput
	err := json.Unmarshal([]byte(*structuredOutputCompletion), &structuredOutput)
	if err != nil {
		return nil
	}

	argsMap := make(map[string]string)
	for _, arg := range structuredOutput.ToolCallArgs {
		argsMap[arg.ToolCallArgName] = arg.ToolCallArgValue
	}

	argsMap["_thoughts"] = structuredOutput.Thoughts
	argsMap["description_of_action"] = structuredOutput.DescriptionOfAction

	argsJson, err := json.Marshal(argsMap)
	if err != nil {
		return nil
	}

	return &tools.ToolCall{
		Name:      structuredOutput.ToolCallName,
		Arguments: string(argsJson),
	}
}

func (a *FastAgent) thinkAndChooseTool(ctx context.Context, iteration int, query string, knowledgeContext *string, prevThoughts *string, prevToolCall *tools.ToolCall, prevToolCallResult *string, prevToolCalls []tools.ToolCall) (*string, *tools.ToolCall, error) {
	kc := ""
	if knowledgeContext != nil {
		kc = *knowledgeContext
	}

	prompt, err := a.thinkingAndActPrompt(ctx, iteration, query, &kc, prevThoughts, prevToolCall, prevToolCallResult, prevToolCalls)
	if err != nil {
		return nil, nil, err
	}

	structuredCompletion, err := a.LLM.StructuredOutputCompletion(context.TODO(), *prompt, StructuredOutput{})
	if err != nil {
		return nil, nil, err
	}

	toolCall := structuredOutputToToolCall(structuredCompletion)
	if toolCall == nil {
		return nil, nil, fmt.Errorf("failed to parse tool call from structured output")
	}

	if toolCall.Name != "give_answer" {
		a.publish(NewToolCallChoiceEvent(ctx, *toolCall))
	}

	var parsedArgs map[string]interface{}
	err = json.Unmarshal([]byte(toolCall.Arguments), &parsedArgs)
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

	return &thoughts, toolCall, nil
}

func (a *FastAgent) thinkingAndActPrompt(ctx context.Context, iteration int, query string, knowledgeContext *string, prevThoughts *string, prevToolCall *tools.ToolCall, prevToolCallResult *string, prevToolCalls []tools.ToolCall) (*string, error) {
	prompt := ""

	chatHistory, err := a.chatHistory(ctx)
	if err != nil {
		return nil, err
	}

	if iteration == 0 {
		prompt = fmt.Sprintf("You are a reAct agent. Here are previous messages: %+v. Your goal is to solve the following: `%s`. Here is some (potentially relevant) knowledge from a rag source: `%s`.  ", chatHistory, query, *knowledgeContext)
	} else {
		prompt = fmt.Sprintf("You are a reAct agent, currently in the process of solving: `%s`. In the previous iteration, you thought `%s` and then called the tool `%s`. The results of this tool where `%s`. ", query, *prevThoughts, *prevToolCall, *prevToolCallResult)
	}

	prompt += "Now, give some thoughts about what you already know, and then generate a plan (based on what you need to find out), of how to solve the problem."

	toolChoiceString, err := a.toolChoicesString()
	if err != nil {
		return nil, err
	}

	prompt += ` You have these tools at your disposal: ` + toolChoiceString + ` It is also essential that you give your thoughts in the _thoughts field. If you believe you already know the answer to the query, or that you will be unable to get the answer, pick the give_answer tool. Information: The date and time is ` + a.Clock.CurrentDateTime().Format(time.RFC3339) + `. ` + a.iterationBasedPrompt(iteration)

	prompt += ` Ensure to also provide a "description_of_action" which is a short description of what you will be doing when calling this tool, in present progressive tense. This will be shown to the user progressively as an interactive loading indicator.`

	if iteration > 0 {
		prompt += fmt.Sprintf(" The tools you have alreaady called, in order of oldest to newest are: %s. Refrain from doing things you have already done.", a.formattedToolsStringFrom(prevToolCalls))
	}

	return &prompt, nil
}

func (a *FastAgent) chatHistory(ctx context.Context) ([]string, error) {
	sessionId, ok := context_keys.GetSessionID(ctx)
	if !ok {
		return nil, nil
	}
	return a.History.GetMessageHistory(sessionId)
}

func (a *FastAgent) iterationBasedPrompt(iteration int) string {
	switch iteration {
	case 0:
		return ""
	case 1:
		return "This is now your second iteration in attempting to solve the query."
	default:
		return fmt.Sprintf("This is now iteration number %v in attempting to solve the query. Unless you cannot answer the question correctly (even with additional user clarification in the case of a vague question), you should really think about giving your answer.", iteration)
	}
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
func (a *FastAgent) extractAnswerAndReason(toolCall *tools.ToolCall) (*string, *string, error) {
	var arguments map[string]string
	if err := json.Unmarshal([]byte(toolCall.Arguments), &arguments); err != nil {
		return nil, nil, err
	}

	reason := arguments["reason"]
	answer := arguments["answer"]

	// a.publish(NewAnswerSuccessEvent(requestId, answer, reason))

	return &answer, &reason, nil
}

func (a *FastAgent) extractEntityIdAndTask(toolCall *tools.ToolCall) (*string, *string, error) {
	var arguments map[string]string
	if err := json.Unmarshal([]byte(toolCall.Arguments), &arguments); err != nil {
		return nil, nil, err
	}

	entityId := arguments["entity_id"]
	task := arguments["task"]

	return &entityId, &task, nil
}

func (a *FastAgent) formattedToolsStringFrom(prevToolCalls []tools.ToolCall) string {
	formattedToolsString := ""
	for _, toolCall := range prevToolCalls {
		formattedToolsString += fmt.Sprintf("Tool: %s, Arguments: %s\n", toolCall.Name, toolCall.Arguments)
	}
	return formattedToolsString
}
