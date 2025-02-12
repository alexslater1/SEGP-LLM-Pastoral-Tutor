package handlers

import (
	"fmt"
	"testing"
	"time"

	"github.com/segp/agents-main/storage"
	"github.com/stretchr/testify/assert"
)

func TestQueryAndResponseFrom(t *testing.T) {
	var (
		requestId1 = "00aaa2a7-d4f1-4dda-b3ed-131f1963d416"
		requestId2 = "0243d463-f0a0-42dd-a305-6414702d4dda"
		requestId3 = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"

		timeRequest1 = time.Now().Add(-1 * time.Hour)
		timeRequest2 = time.Now().Add(-2 * time.Hour)
		timeRequest3 = time.Now().Add(-3 * time.Hour)

		chatId1 = "a chat id"
		chatId2 = "another chat id"
		chatId3 = "yet another chat id"

		query1 = "What is imperials policy on late coursework submissions?"
		query2 = "What is the sum of the temperature of in japan and the temperature in london, both in farenheight?"
		query3 = "What is the current USD to GBP exchange rate?"

		agentRequest1 = storage.AgentRequest{
			ID:        requestId1,
			CreatedAt: &timeRequest1,
			Endpoint:  "endpoint1",
			Metadata:  map[string]interface{}{"query": query1},
			ChatID:    chatId1,
		}

		agentRequest2 = storage.AgentRequest{
			ID:        requestId2,
			CreatedAt: &timeRequest2,
			Endpoint:  "endpoint2",
			Metadata:  map[string]interface{}{"query": query2},
			ChatID:    chatId2,
		}

		agentRequest3 = storage.AgentRequest{
			ID:        requestId3,
			CreatedAt: &timeRequest3,
			Endpoint:  "endpoint3",
			Metadata:  map[string]interface{}{"query": query3},
			ChatID:    chatId3,
		}

		metadata1 = map[string]interface{}{
			"toolCallChoice": map[string]interface{}{
				"name":      "google_search_first_results_page_contents",
				"arguments": "{\"_thoughts\":\"The user is asking for the current USD price. This is a very general request and needs more context. I assume the user is asking for the current exchange rate of USD to some other currency, likely GBP since the context seems to be about UK institutions. However, without knowing the target currency, I need to use a search query to find the current USD to GBP exchange rate. I can specify 'current usd to gbp exchange rate' in my query to get the needed information.\",\"description_of_action\":\"I am searching for the current USD to GBP exchange rate.\",\"query\":\"current usd to gbp exchange rate\"}",
			},
		}

		metadata2 = map[string]interface{}{
			"toolCallResult": "tool call result",
		}

		metadata3 = map[string]interface{}{
			"answer": "The current exchange rate is 1 USD = 0.814 GBP.",
			"reason": "I have the current USD to GBP exchange rate from the previous search.",
		}

		metadata4 = map[string]interface{}{
			"toolCallChoice": map[string]interface{}{
				"name":      "google_search_first_results_page_contents",
				"arguments": "{\"_thoughts\":\"this  is a thought.\",\"description_of_action\":\"this is a description.\",\"query\":\"this is a query.\"}",
			},
		}

		metadata5 = map[string]interface{}{
			"toolCallResult": "another tool call result",
		}

		metadata6 = map[string]interface{}{
			"error": "an error",
		}

		metadata7 = map[string]interface{}{
			"toolCallChoice": map[string]interface{}{
				"name":      "google_search_first_results_page_contents",
				"arguments": "{\"_thoughts\":\"this  is such a good thought.\",\"description_of_action\":\"this is a really good yeah description.\",\"query\":\"this is a good good query.\"}",
			},
		}

		metadata8 = map[string]interface{}{
			"toolCallResult": "another tool call result",
		}

		event1 = storage.AgentEvent{
			Type:     "tool_call_choice",
			Metadata: metadata1,
		}

		event2 = storage.AgentEvent{
			Type:     "tool_call_result",
			Metadata: metadata2,
		}

		event3 = storage.AgentEvent{
			Type:     "answer_success",
			Metadata: metadata3,
		}

		event4 = storage.AgentEvent{
			Type:     "tool_call_choice",
			Metadata: metadata4,
		}

		event5 = storage.AgentEvent{
			Type:     "tool_call_result",
			Metadata: metadata5,
		}

		event6 = storage.AgentEvent{
			Type:     "error",
			Metadata: metadata6,
		}

		event7 = storage.AgentEvent{
			Type:     "tool_call_choice",
			Metadata: metadata7,
		}

		event8 = storage.AgentEvent{
			Type:     "tool_call_result",
			Metadata: metadata8,
		}

		requests = []storage.AgentRequest{agentRequest1, agentRequest2, agentRequest3}
		events   = [][]storage.AgentEvent{{event3, event2, event1}, {event6, event5, event4}, {event8, event7}}
	)

	queryAndResponses := queryAndResponsesFrom(requests, events)

	fmt.Printf("%+v\n", queryAndResponses)

	assert.Equal(t, len(queryAndResponses), 3)
	assert.Equal(t, queryAndResponses[0].Query, query1)
	assert.Equal(t, queryAndResponses[1].Query, query2)
	assert.Equal(t, queryAndResponses[2].Query, query3)

	assert.Equal(t, queryAndResponses[0].RequestID, requestId1)
	assert.Equal(t, queryAndResponses[1].RequestID, requestId2)
	assert.Equal(t, queryAndResponses[2].RequestID, requestId3)

	assert.Equal(t, queryAndResponses[0].Type, ChatCompletionV2StatusResponseTypeCompleted)
	assert.Equal(t, queryAndResponses[0].Actions, []string{"Thinking", "I am searching for the current USD to GBP exchange rate."})
	assert.Equal(t, queryAndResponses[0].Answer, "The current exchange rate is 1 USD = 0.814 GBP.")
	assert.Equal(t, queryAndResponses[0].CurrentAction, "")
	assert.Equal(t, queryAndResponses[0].Error, "")

	assert.Equal(t, queryAndResponses[1].Type, ChatCompletionV2StatusResponseTypeError)
	assert.Equal(t, queryAndResponses[1].Actions, []string{"Thinking", "this is a description."})
	assert.Equal(t, queryAndResponses[1].Error, "an error")
	assert.Equal(t, queryAndResponses[1].Answer, "")
	assert.Equal(t, queryAndResponses[1].CurrentAction, "")

	assert.Equal(t, queryAndResponses[2].Type, ChatCompletionV2StatusResponseTypePending)
	assert.Equal(t, queryAndResponses[2].Actions, []string{"Thinking", "this is a really good yeah description."})
	assert.Equal(t, queryAndResponses[2].Answer, "")
	assert.Equal(t, queryAndResponses[2].CurrentAction, "Thinking")
	assert.Equal(t, queryAndResponses[2].Error, "")
}
