package history

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/segp/agents-main/storage"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	os.Exit(m.Run())
}

func TestStoreHistory(t *testing.T) {
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
			Type:      "tool_call_choice",
			Metadata:  metadata1,
			RequestID: requestId1,
		}

		event2 = storage.AgentEvent{
			Type:      "tool_call_result",
			Metadata:  metadata2,
			RequestID: requestId1,
		}

		event3 = storage.AgentEvent{
			Type:      "answer_success",
			Metadata:  metadata3,
			RequestID: requestId1,
		}

		event4 = storage.AgentEvent{
			Type:      "tool_call_choice",
			Metadata:  metadata4,
			RequestID: requestId2,
		}

		event5 = storage.AgentEvent{
			Type:      "tool_call_result",
			Metadata:  metadata5,
			RequestID: requestId2,
		}

		event6 = storage.AgentEvent{
			Type:      "error",
			Metadata:  metadata6,
			RequestID: requestId2,
		}

		event7 = storage.AgentEvent{
			Type:      "tool_call_choice",
			Metadata:  metadata7,
			RequestID: requestId3,
		}

		event8 = storage.AgentEvent{
			Type:      "tool_call_result",
			Metadata:  metadata8,
			RequestID: requestId3,
		}

		requests = []storage.AgentRequest{agentRequest1, agentRequest2, agentRequest3}
		events   = [][]storage.AgentEvent{{event3, event2, event1}, {event6, event5, event4}, {event8, event7}}
	)

	store := storage.NewMemoryStorage()
	history := NewStoreHistory(store)

	storage.StoreAll(store, requests...)
	storage.StoreAll(store, events[0]...)
	storage.StoreAll(store, events[1]...)
	storage.StoreAll(store, events[2]...)

	h1, err := history.GetChatHistory(chatId1)
	if err != nil {
		t.Errorf("Error getting chat history: %v", err)
	}

	h2, err := history.GetChatHistory(chatId2)
	if err != nil {
		t.Errorf("Error getting chat history: %v", err)
	}

	h3, err := history.GetChatHistory(chatId3)
	if err != nil {
		t.Errorf("Error getting chat history: %v", err)
	}

	assert.Equal(t, h1, []string{"Query: What is imperials policy on late coursework submissions?\nResponse: The current exchange rate is 1 USD = 0.814 GBP."})
	assert.Equal(t, h2, []string{"Query: What is the sum of the temperature of in japan and the temperature in london, both in farenheight?\nResponse: [ERROR]"})
	assert.Equal(t, h3, []string{"Query: What is the current USD to GBP exchange rate?\nResponse: [PENDING]"})
}

func TestStoreHistoryWithSupabase(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CI/CD")
	}

	store := storage.NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))
	history := NewStoreHistory(store)

	h, err := history.GetChatHistory("2fea8a5f-b82c-4261-9889-3e42136d9ef0")
	if err != nil {
		t.Errorf("Error getting chat history: %v", err)
	}

	fmt.Printf("History: %v\n", h)
}
