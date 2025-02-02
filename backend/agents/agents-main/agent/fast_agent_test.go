package agent

import (
	"testing"

	googleSearch "github.com/segp/agents-main/google_search"
	"github.com/segp/agents-main/knowledge"
	"github.com/segp/agents-main/llm"
	"github.com/segp/agents-main/tools"
	"github.com/stretchr/testify/assert"
)

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
		agent = NewFastAgent("test", toolHandler, llm, knowledge)
	)

	expectedStr := `[{"Name":"google_search_first_results_page_contents","Description":"Returns the contents of the first 3 search results for a query.","Parameters":[{"Name":"query","Description":"The query to search for.","Type":"string"}]},{"Name":"rag_tool","Description":"Get more relevant context and/or important links/contact information about the given query in respect to Imperial College London. Uses RAG","Parameters":[{"Name":"query","Description":"The query to send to the RAG model.","Type":"string"}]},{"Name":"google_maps_results","Description":"Returns the results of a google maps search. For each place listed, will get the title, rating, PLACES URL and then some info + tags about the place.","Parameters":[{"Name":"query","Description":"The query to search for results for","Type":"string"}]},{"Name":"google_maps_place","Description":"Returns lots of information about a specific place on Google Maps, given a GoogleMapsPlaceURL.","Parameters":[{"Name":"google_maps_place_url","Description":"The URL of the place on Google Maps. Comes from a previous google_maps_results tool call.","Type":"string"}]},{"Name":"no_tool","Description":"Do not use any tools. This could be because you have an answer, or you deem that after sufficient attempts, it will not be possible to feasibly find an accurate answer.","Parameters":[{"Name":"reason","Description":"The reason why no tool was used","Type":"string"},{"Name":"answer","Description":"The answer to the question / reason why not possible to answer the question","Type":"string"}]}]`

	toolChoicesString, err := agent.toolChoicesString()
	assert.NoError(t, err)
	assert.Equal(t, expectedStr, toolChoicesString)
}

func TestThinkingAndActPromptFirstIteration(t *testing.T) {
	var (
		knowledge   = knowledge.NewLocalKnowledge()
		toolHandler = tools.NewToolHandler([]tools.Tool{
			tools.NewRagTool(knowledge),
		})
		llm   = llm.NewMockLLM()
		agent = NewFastAgent("test", toolHandler, llm, knowledge)

		query            = "test query"
		knowledgeContext = "test knowledge context"
		prevThoughts     = "test prev thoughts"
		prevToolCall     = tools.ToolCall{
			Name:      "test_tool_call",
			Arguments: `{"x": 1, "y": 2}`,
		}
		prevToolCallResult = "test prev tool call result"
		requestId          = "test_request_id"
	)

	prompt, err := agent.thinkingAndActPrompt(true, query, &knowledgeContext, &prevThoughts, &prevToolCall, &prevToolCallResult, requestId)
	assert.NoError(t, err)

	expectedPrompt := `You are a reAct agent. Your goal is to solve the following query: ` + "`test query`" + `. Here is some (potentially relevant) knowledge from a rag source: ` + "`test knowledge context`" + `.  Now, give some thoughts about what you already know, and then generate a plan (based on what you need to find out), of how to solve the problem. You have these tools at your disposal: [{"Name":"rag_tool","Description":"Get more relevant context and/or important links/contact information about the given query in respect to Imperial College London. Uses RAG","Parameters":[{"Name":"query","Description":"The query to send to the RAG model.","Type":"string"}]},{"Name":"no_tool","Description":"Do not use any tools. This could be because you have an answer, or you deem that after sufficient attempts, it will not be possible to feasibly find an accurate answer.","Parameters":[{"Name":"reason","Description":"The reason why no tool was used","Type":"string"},{"Name":"answer","Description":"The answer to the question / reason why not possible to answer the question","Type":"string"}]}] It is also essential that you give your thoughts in the _thoughts field. If you believe you already know the answer to the query, or that you will be unable to get the answer, pick the no_tool tool.`

	assert.Equal(t, expectedPrompt, *prompt)
}

func TestThinkingAndActPromptSubsequentIteration(t *testing.T) {

}
