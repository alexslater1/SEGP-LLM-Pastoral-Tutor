package history

import (
	"testing"
	"time"

	"github.com/segp/agents-main/storage"
	"github.com/stretchr/testify/assert"
)

var (
	now = time.Now()

	requestId1 = "00aaa2a7-d4f1-4dda-b3ed-131f1963d416"
	requestId2 = "0243d463-f0a0-42dd-a305-6414702d4dda"

	sessionId = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"

	metadataQuery1 = map[string]interface{}{
		"query": "What is imperials policy on late coursework submissions?",
	}

	metadataQuery2 = map[string]interface{}{
		"query": "What is the current exchange rate of USD to GBP?",
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

	eventQuery1 = storage.AgentEvent{
		Type:      "query",
		Metadata:  metadataQuery1,
		RequestID: requestId1,
		CreatedAt: timeAdd(now, 1*time.Second),
	}

	eventQuery2 = storage.AgentEvent{
		Type:      "query",
		Metadata:  metadataQuery2,
		RequestID: requestId2,
		CreatedAt: timeAdd(now, 2*time.Second),
	}

	event1 = storage.AgentEvent{
		Type:      "tool_call_choice",
		Metadata:  metadata1,
		RequestID: requestId1,
		CreatedAt: timeAdd(now, 3*time.Second),
	}

	event2 = storage.AgentEvent{
		Type:      "tool_call_result",
		Metadata:  metadata2,
		RequestID: requestId1,
		CreatedAt: timeAdd(now, 4*time.Second),
	}

	event3 = storage.AgentEvent{
		Type:      "answer_success",
		Metadata:  metadata3,
		RequestID: requestId1,
		CreatedAt: timeAdd(now, 5*time.Second),
	}

	event4 = storage.AgentEvent{
		Type:      "tool_call_choice",
		Metadata:  metadata4,
		RequestID: requestId2,
		CreatedAt: timeAdd(now, 6*time.Second),
	}

	event5 = storage.AgentEvent{
		Type:      "tool_call_result",
		Metadata:  metadata5,
		RequestID: requestId2,
		CreatedAt: timeAdd(now, 7*time.Second),
	}

	requestSession1 = storage.RequestSession{
		SessionID: sessionId,
		RequestID: requestId1,
		CreatedAt: timeAdd(now, 8*time.Second),
	}

	requestSession2 = storage.RequestSession{
		SessionID: sessionId,
		RequestID: requestId2,
		CreatedAt: timeAdd(now, 9*time.Second),
	}

	store = storage.NewMemoryStorage()
)

func TestAgentEventGetMessageHistory(t *testing.T) {
	storage.StoreAll(store, requestSession1, requestSession2)
	storage.StoreAll(store, eventQuery1, eventQuery2, event1, event2, event3, event4, event5)

	history := NewAgentEventHistory(store)

	h, err := history.GetMessageHistory(sessionId)
	if err != nil {
		t.Fatalf("failed to get message history: %v", err)
	}

	assert.Equal(t, h, []string{
		"Query: What is imperials policy on late coursework submissions?\nResponse: The current exchange rate is 1 USD = 0.814 GBP.",
		"Query: What is the current exchange rate of USD to GBP?\nResponse: [PENDING]	",
	})
}

func TestAgentEventGetMessagesAndActions(t *testing.T) {
	storage.StoreAll(store, requestSession1, requestSession2)
	storage.StoreAll(store, eventQuery1, eventQuery2, event1, event2, event3, event4, event5)

	history := NewAgentEventHistory(store)

	messagesAndActions, err := history.GetMessagesAndActions(sessionId)
	if err != nil {
		t.Fatalf("failed to get message history: %v", err)
	}

	// Expect exactly two returned message/action items.
	assert.Equal(t, 2, len(messagesAndActions))

	assert.Equal(t, StatusResponseTypeCompleted, messagesAndActions[0].Type)
	assert.Equal(t, "What is imperials policy on late coursework submissions?", messagesAndActions[0].Query)
	assert.Equal(t, requestId1, messagesAndActions[0].RequestID)
	assert.Equal(t, []string{"I am searching for the current USD to GBP exchange rate.", "Thinking"}, messagesAndActions[0].Actions)
	assert.Equal(t, "The current exchange rate is 1 USD = 0.814 GBP.", messagesAndActions[0].Answer)
	assert.Equal(t, "", messagesAndActions[0].Error)
	assert.Equal(t, "", messagesAndActions[0].CurrentAction)

	assert.Equal(t, StatusResponseTypePending, messagesAndActions[1].Type)
	assert.Equal(t, "What is the current exchange rate of USD to GBP?", messagesAndActions[1].Query)
	assert.Equal(t, requestId2, messagesAndActions[1].RequestID)
	assert.Equal(t, []string{"this is a description.", "Thinking"}, messagesAndActions[1].Actions)
	assert.Equal(t, "", messagesAndActions[1].Answer)
	assert.Equal(t, "", messagesAndActions[1].Error)
	assert.Equal(t, "this is a description.", messagesAndActions[1].CurrentAction)
}

func timeAdd(t time.Time, d time.Duration) *time.Time {
	t = t.Add(d)
	return &t
}
