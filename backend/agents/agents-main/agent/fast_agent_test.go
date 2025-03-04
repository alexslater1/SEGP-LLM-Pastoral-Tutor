package agent

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/segp/agents-main/clock"
	"github.com/segp/agents-main/email"
	googleSearch "github.com/segp/agents-main/google_search"
	"github.com/segp/agents-main/history"
	"github.com/segp/agents-main/knowledge"
	"github.com/segp/agents-main/llm"
	"github.com/segp/agents-main/tools"
	"github.com/segp/agents-main/utils"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	os.Exit(m.Run())
}

func TestToolCallChoiceString(t *testing.T) {
	var (
		knowledge    = knowledge.NewLocalKnowledge()
		googleSearch = googleSearch.NewMockGoogleSearchClient()
		toolHandler  = tools.NewToolHandler([]tools.Tool{
			tools.NewGoogleSearchFirstResultsPageContentsTool(googleSearch, 3),
			tools.NewRagTool(knowledge),
			tools.NewGoogleMapsResultsTool(googleSearch),
			tools.NewGoogleMapsPlaceTool(googleSearch),
		})
		llm   = llm.NewMockLLM()
		clock = clock.NewMockClock()
		h     = history.NewLocalHistory()
		agent = newFastAgent("test", "test", "test", toolHandler, llm, knowledge, clock, h, newLoggingCallback())
	)

	toolChoicesString, err := agent.toolChoicesString()
	assert.NoError(t, err)

	expectedStr := "[{\"Name\":\"google_maps_place\",\"Description\":\"Returns lots of information about a specific place on Google Maps, given a GoogleMapsPlaceURL.\",\"Parameters\":[{\"Name\":\"google_maps_place_url\",\"Description\":\"The URL of the place on Google Maps. Comes from a previous google_maps_results tool call.\",\"Type\":\"string\"}]},{\"Name\":\"google_maps_results\",\"Description\":\"Returns the results of a google maps search. For each place listed, will get the title, rating, PLACES URL and then some info + tags about the place.\",\"Parameters\":[{\"Name\":\"query\",\"Description\":\"The query to search for results for\",\"Type\":\"string\"}]},{\"Name\":\"google_search_first_results_page_contents\",\"Description\":\"Returns the contents of the first 3 search results for a query.\",\"Parameters\":[{\"Name\":\"query\",\"Description\":\"The query to search for.\",\"Type\":\"string\"}]},{\"Name\":\"no_tool\",\"Description\":\"Do not use any tools. This could be because you have an answer, or you deem that after sufficient attempts, it will not be possible to feasibly find an accurate answer.\",\"Parameters\":[{\"Name\":\"reason\",\"Description\":\"The reason why no tool was used\",\"Type\":\"string\"},{\"Name\":\"answer\",\"Description\":\"The answer to the question / reason why not possible to answer the question\",\"Type\":\"string\"}]},{\"Name\":\"rag_tool\",\"Description\":\"Get more relevant context and/or important links/contact information about the given query in respect to Imperial College London. Uses RAG\",\"Parameters\":[{\"Name\":\"query\",\"Description\":\"The query to send to the RAG model.\",\"Type\":\"string\"}]}]"
	assert.Equal(t, expectedStr, toolChoicesString)
}

func TestThinkingAndActPromptFirstIteration(t *testing.T) {
	var (
		knowledge   = knowledge.NewLocalKnowledge()
		toolHandler = tools.NewToolHandler([]tools.Tool{
			tools.NewRagTool(knowledge),
		})
		llm   = llm.NewMockLLM()
		clock = clock.NewMockClock()
		h     = history.NewLocalHistory()
		agent = newFastAgent("test", "test", "test", toolHandler, llm, knowledge, clock, h, newLoggingCallback())

		query            = "test query"
		knowledgeContext = "test knowledge context"
		prevThoughts     = "test prev thoughts"
		prevToolCall     = tools.ToolCall{
			Name:      "test_tool_call",
			Arguments: `{"x": 1, "y": 2}`,
		}
		prevToolCallResult = "test prev tool call result"
	)

	prompt, err := agent.thinkingAndActPrompt(context.Background(), 0, query, &knowledgeContext, &prevThoughts, &prevToolCall, &prevToolCallResult, []tools.ToolCall{})
	assert.NoError(t, err)

	// Only check the static parts of the prompt (i.e. not the date)
	assert.Equal(t, *prompt, "You are a reAct agent. Your goal is to solve the following query: `test query`. Here is some (potentially relevant) knowledge from a rag source: `test knowledge context`.  Now, give some thoughts about what you already know, and then generate a plan (based on what you need to find out), of how to solve the problem. You have these tools at your disposal: [{\"Name\":\"no_tool\",\"Description\":\"Do not use any tools. This could be because you have an answer, or you deem that after sufficient attempts, it will not be possible to feasibly find an accurate answer.\",\"Parameters\":[{\"Name\":\"reason\",\"Description\":\"The reason why no tool was used\",\"Type\":\"string\"},{\"Name\":\"answer\",\"Description\":\"The answer to the question / reason why not possible to answer the question\",\"Type\":\"string\"}]},{\"Name\":\"rag_tool\",\"Description\":\"Get more relevant context and/or important links/contact information about the given query in respect to Imperial College London. Uses RAG\",\"Parameters\":[{\"Name\":\"query\",\"Description\":\"The query to send to the RAG model.\",\"Type\":\"string\"}]}] It is also essential that you give your thoughts in the _thoughts field. If you believe you already know the answer to the query, or that you will be unable to get the answer, pick the no_tool tool. Information: The date and time is 2025-02-04T12:00:00Z.  Ensure to also provide a \"description_of_action\" which is a short description of what you will be doing when calling this tool, in present progressive tense. This will be shown to the user progressively as an interactive loading indicator.")
}

func TestThinkingAndActPromptSubsequentIteration(t *testing.T) {
	var (
		knowledge   = knowledge.NewLocalKnowledge()
		toolHandler = tools.NewToolHandler([]tools.Tool{
			tools.NewRagTool(knowledge),
		})
		llm   = llm.NewMockLLM()
		clock = clock.NewMockClock()
		h     = history.NewLocalHistory()
		agent = newFastAgent("test", "test", "test", toolHandler, llm, knowledge, clock, h, newLoggingCallback())

		query            = "test query"
		knowledgeContext = "test knowledge context"
		prevThoughts     = "test prev thoughts"
		prevToolCall     = tools.ToolCall{
			Name:      "test_tool_call",
			Arguments: `{"x": 1, "y": 2}`,
		}
		prevToolCallResult = "test prev tool call result"
	)

	prompt, err := agent.thinkingAndActPrompt(context.Background(), 1, query, &knowledgeContext, &prevThoughts, &prevToolCall, &prevToolCallResult, []tools.ToolCall{prevToolCall})
	assert.NoError(t, err)

	// Only check the static parts of the prompt (i.e. not the date)
	assert.Equal(t, *prompt, "You are a reAct agent, currently in the process of solving the query: `test query`. In the previous iteration, you thought `test prev thoughts` and then called the tool `{test_tool_call {\"x\": 1, \"y\": 2}}`. The results of this tool where `test prev tool call result`. Now, give some thoughts about what you already know, and then generate a plan (based on what you need to find out), of how to solve the problem. You have these tools at your disposal: [{\"Name\":\"no_tool\",\"Description\":\"Do not use any tools. This could be because you have an answer, or you deem that after sufficient attempts, it will not be possible to feasibly find an accurate answer.\",\"Parameters\":[{\"Name\":\"reason\",\"Description\":\"The reason why no tool was used\",\"Type\":\"string\"},{\"Name\":\"answer\",\"Description\":\"The answer to the question / reason why not possible to answer the question\",\"Type\":\"string\"}]},{\"Name\":\"rag_tool\",\"Description\":\"Get more relevant context and/or important links/contact information about the given query in respect to Imperial College London. Uses RAG\",\"Parameters\":[{\"Name\":\"query\",\"Description\":\"The query to send to the RAG model.\",\"Type\":\"string\"}]}] It is also essential that you give your thoughts in the _thoughts field. If you believe you already know the answer to the query, or that you will be unable to get the answer, pick the no_tool tool. Information: The date and time is 2025-02-04T12:00:00Z. This is now your second iteration in attempting to solve the query. Ensure to also provide a \"description_of_action\" which is a short description of what you will be doing when calling this tool, in present progressive tense. This will be shown to the user progressively as an interactive loading indicator. The tools you have alreaady called, in order of oldest to newest are: Tool: test_tool_call, Arguments: {\"x\": 1, \"y\": 2}\n. Refrain from doing things you have already done.")
}

func TestThinkAndChooseTool(t *testing.T) {
	var (
		knowledge   = knowledge.NewLocalKnowledge()
		toolHandler = tools.NewToolHandler([]tools.Tool{
			tools.NewRagTool(knowledge),
		})
		llm   = llm.NewMockLLM()
		clock = clock.NewMockClock()
		h     = history.NewLocalHistory()
		agent = newFastAgent("test", "test", "test", toolHandler, llm, knowledge, clock, h, newLoggingCallback())

		query            = "test query"
		knowledgeContext = "test knowledge context"
		prevThoughts     = "test prev thoughts"
		prevToolCall     = tools.ToolCall{
			Name:      "test_tool_call",
			Arguments: `{"x": 1, "y": 2}`,
		}
		prevToolCallResult = "test prev tool call result"

		// prompt = "You are a reAct agent. Your goal is to solve the following query: `test query`. Here is some (potentially relevant) knowledge from a rag source: `test knowledge context`.  Now, give some thoughts about what you already know, and then generate a plan (based on what you need to find out), of how to solve the problem. You have these tools at your disposal: [{\"Name\":\"rag_tool\",\"Description\":\"Get more relevant context and/or important links/contact information about the given query in respect to Imperial College London. Uses RAG\",\"Parameters\":[{\"Name\":\"query\",\"Description\":\"The query to send to the RAG model.\",\"Type\":\"string\"}]},{\"Name\":\"no_tool\",\"Description\":\"Do not use any tools. This could be because you have an answer, or you deem that after sufficient attempts, it will not be possible to feasibly find an accurate answer.\",\"Parameters\":[{\"Name\":\"reason\",\"Description\":\"The reason why no tool was used\",\"Type\":\"string\"},{\"Name\":\"answer\",\"Description\":\"The answer to the question / reason why not possible to answer the question\",\"Type\":\"string\"}]}] It is also essential that you give your thoughts in the _thoughts field. If you believe you already know the answer to the query, or that you will be unable to get the answer, pick the no_tool tool."
	)

	llm.NewCallChain().
		ThenStructured(`{"_thoughts": "test thoughts", "tool_call_name": "rag_tool", "tool_call_args": [{"tool_call_arg_name": "query", "tool_call_arg_value": "test query"}], "description_of_action": "test description of action"}`).
		Set()

	thoughts, toolCall, err := agent.thinkAndChooseTool(context.Background(), 0, query, &knowledgeContext, &prevThoughts, &prevToolCall, &prevToolCallResult, []tools.ToolCall{})
	assert.NoError(t, err)

	assert.Equal(t, "rag_tool", toolCall.Name)
	assert.Equal(t, "test thoughts", *thoughts)
}

// func TestFastAgentRun(t *testing.T) {
// 	if os.Getenv("CICD") == "true" {
// 		t.Skip("skipping test in CI")
// 	}

// 	agent := NewDefaultLoggingUserQueryAgent()
// 	ctx := context_keys.SetRequestID(context.Background(), "2fea8a5f-b82c-4261-9889-3e42136d9ef0")

// 	response, err := agent.Run(ctx, "what about in 3 days?")
// 	if err != nil {
// 		t.Fatalf("error running agent: %v", err)
// 	}

// 	t.Logf("answer: %s", *response.Answer)
// 	t.Logf("reason: %s", *response.Reason)
// }

func TestFastAgentAskQuestion(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("skipping test in CI")
	}

	agent := newFastAgent("test", "test", "test", tools.NewToolHandler([]tools.Tool{
		// tools.NewGoogleSearchFirstResultsPageContentsTool(googleSearch.NewMockGoogleSearchClient(), 3),
	}), llm.NewGeminiLLM(context.Background(), utils.Required(os.Getenv("GEMINI_API_KEY"), "GEMINI_API_KEY is not set")), knowledge.NewLocalKnowledge(), clock.NewMockClock(), history.NewLocalHistory(), newLoggingCallback())

	query := "What is the temperature?"
	response, err := agent.Run(context.Background(), query)
	if err != nil {
		t.Fatalf("error running agent: %v", err)
	}

	t.Logf("answer: %s", *response.Answer)
	t.Logf("reason: %s", *response.Reason)
}

func TestFastAgentDescription(t *testing.T) {
	agent := newFastAgent("test", "a description", "a prompt", tools.NewToolHandler([]tools.Tool{}), llm.NewMockLLM(), knowledge.NewLocalKnowledge(), clock.NewMockClock(), history.NewLocalHistory(), newLoggingCallback())
	assert.Equal(t, "a description", agent.Description())
}

func TestFastAgentId(t *testing.T) {
	agent := newFastAgent("test", "a description", "a prompt", tools.NewToolHandler([]tools.Tool{}), llm.NewMockLLM(), knowledge.NewLocalKnowledge(), clock.NewMockClock(), history.NewLocalHistory(), newLoggingCallback())
	assert.Equal(t, "test", agent.Id())
}

func TestPersonalTutorAgent(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("skipping test in CI")
	}

	var (
		knowledge = knowledge.NewLocalKnowledge()
		llm       = llm.NewGeminiLLM(context.Background(), utils.Required(os.Getenv("GEMINI_API_KEY"), "GEMINI_API_KEY is not set"))
		clock     = clock.NewMockClock()
		h         = history.NewLocalHistory()

		id     = "personal_tutor_agent"
		desc   = "A personal tutor agent"
		prompt = "You are a personal tutor agent. You are meant to provide support for a student at imperial college london. You are a layer between the students and their personal tutor. Students interact with you via a chatbot. In the case where you have flagged something concerning, you must use the email tool to send an email to the personal tutor, raising this concern and your reasons. You must also always reply to the user in a way which is supportive."

		agent = newFastAgent(id, desc, prompt, tools.NewToolHandler([]tools.Tool{
			tools.NewEmailTool("personal.tutor@imperial.ac.uk", "Personal Tutor", email.NewMockEmailClient(), "To be used to send an email to a personal tutor, in case of a concern."),
		}), llm, knowledge, clock, h, newLoggingCallback())
	)

	resp, err := agent.Run(context.Background(), "What's 1 + 1?")
	if err != nil {
		t.Fatalf("error running agent: %v", err)
	}

	t.Logf("answer: %s", *resp.Answer)
	t.Logf("reason: %s", *resp.Reason)
}
