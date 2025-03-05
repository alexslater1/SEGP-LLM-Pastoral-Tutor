package jobs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segp/agents-main/agent"
	"github.com/segp/agents-main/history"
	"github.com/segp/agents-main/llm"
	"github.com/segp/agents-main/storage"
	"github.com/segp/agents-main/utils"
)

const (
	name = "adjust_prompts"

	numChatsToConsider = 5
)

type AdjustPromptsJob struct {
	store         storage.Storage
	history       history.History
	agentProvider *agent.AgentProvider
	llm           llm.LLM
	interval      time.Duration
}

func NewAdjustPromptsJob(store storage.Storage, history history.History, llm llm.LLM, agentProvider *agent.AgentProvider, interval time.Duration) *AdjustPromptsJob {
	return &AdjustPromptsJob{store: store, history: history, llm: llm, agentProvider: agentProvider, interval: interval}
}

func (j *AdjustPromptsJob) Name() string {
	return name
}

func (j *AdjustPromptsJob) Interval() time.Duration {
	return j.interval
}

func (j *AdjustPromptsJob) Run() error {
	feedbackChecksEnabled, err := storage.GetAll[storage.FeedbackChecksEnabled](j.store, storage.NewQueryBuilder().OrderBy("created_at", storage.OrderByDesc).Limit(1))
	if err != nil {
		return fmt.Errorf("error getting feedback checks enabled: %v", err)
	}

	if len(feedbackChecksEnabled) == 0 || !feedbackChecksEnabled[0].Enabled {
		log.Println("Feedback checks are not enabled: skipping")
		return nil
	}

	log.Println("Executing chat checker job")

	lastJob, err := storage.GetAll[storage.FeedbackCheck](j.store, storage.NewQueryBuilder().OrderBy("created_at", storage.OrderByDesc).Limit(1))
	if err != nil {
		return err
	}

	since := time.Now().Add(-999999 * time.Hour)
	if len(lastJob) > 0 {
		since = *lastJob[0].CreatedAt
	}

	feedback, err := getFeedback(j.store, since)
	if err != nil {
		return err
	}

	// fmt.Printf("Feedback: %+v\n", feedback)

	enrichedFeedback, err := j.getEnrichedFeedback(feedback)
	if err != nil {
		return err
	}

	groupedByAgentFeedback := groupByAgentID(enrichedFeedback)

	agents := j.agentProvider.GetAgents()
	newAgentPromptsTasks := utils.DoAsyncList(agents, func(agent agent.Agent) (*string, error) {
		agentID := agent.Id()
		feedback, ok := groupedByAgentFeedback[agentID]
		if !ok {
			return nil, nil
		}

		log.Printf("Adjusting prompt for agent %s", agentID)
		prompt := promptFrom(feedback, agent.Prompt())
		return j.llm.ChatCompletion(context.TODO(), prompt)
	})

	ps, err := utils.GetAsyncList(newAgentPromptsTasks)
	if err != nil {
		return err
	}

	updatingAgentPrompts := make(map[string]string)
	for i, prompt := range ps {
		if prompt == nil {
			continue
		}
		updatingAgentPrompts[agents[i].Id()] = *prompt
	}

	fmt.Printf("%+v\n", updatingAgentPrompts)

	storage.StoreAll(j.store, storage.NewFeedbackCheck())
	return nil
	// return j.agentProvider.UpdatePrompts(updatingAgentPrompts)
}

type Feedback struct {
	id        string
	createdAt time.Time
	requestID string
	reason    string
	userID    string
	isUpvote  bool
}

type EnrichedFeedback struct {
	Feedback
	AgentID         string
	Context         []string
	MessageResponse string
}

func (j *AdjustPromptsJob) getEnrichedFeedback(feedback []Feedback) ([]EnrichedFeedback, error) {
	enrichedFeedbackTasks := utils.DoAsyncList(feedback, func(feedback Feedback) (EnrichedFeedback, error) {
		return j.getSingleEnrichedFeedbackFrom(feedback)
	})
	return utils.GetAsyncList(enrichedFeedbackTasks)
}

func (j *AdjustPromptsJob) getSingleEnrichedFeedbackFrom(feedback Feedback) (EnrichedFeedback, error) {
	requestSessions, err := storage.GetAll[storage.RequestSession](j.store, storage.NewQueryBuilder().Eq("request_id", feedback.requestID))
	if err != nil {
		return EnrichedFeedback{}, fmt.Errorf("failed to get request session for request ID %s: %w", feedback.requestID, err)
	}

	if len(requestSessions) != 1 {
		return EnrichedFeedback{}, fmt.Errorf("expected 1 request session for request ID %s, got %d", feedback.requestID, len(requestSessions))
	}

	requestSession := requestSessions[0]

	agentMessagesAndActions, err := j.history.GetMessagesAndActions(requestSession.SessionID)
	if err != nil {
		return EnrichedFeedback{}, fmt.Errorf("failed to get agent messages and actions for request session session ID %s: %w", requestSession.SessionID, err)
	}

	index, err := indexOfMatchingRequestID(agentMessagesAndActions, feedback.requestID)
	if err != nil {
		return EnrichedFeedback{}, fmt.Errorf("failed to get index of matching request ID %s: %w", feedback.requestID, err)
	}

	context := history.MessagesAndActionsToMessageHistory(agentMessagesAndActions[max(0, index-numChatsToConsider) : index+1])

	relevantAgentId := agentMessagesAndActions[index].AgentIDs[len(agentMessagesAndActions[index].AgentIDs)-1]

	return EnrichedFeedback{
		Feedback:        feedback,
		AgentID:         relevantAgentId,
		Context:         context,
		MessageResponse: context[len(context)-1],
	}, nil
}

func indexOfMatchingRequestID(agentEvents []history.MessagesAndActions, requestID string) (int, error) {
	for i, event := range agentEvents {
		if event.RequestID == requestID {
			return i, nil
		}
	}
	return -1, fmt.Errorf("request ID not found in agent events: %s", requestID)
}

func getFeedback(store storage.Storage, since time.Time) ([]Feedback, error) {
	upvotesTask := utils.DoAsync(func() ([]storage.UpvotedResponses, error) {
		return storage.GetAll[storage.UpvotedResponses](store, storage.NewQueryBuilder().Gt("created_at", since))
	})

	downvotes, err := storage.GetAll[storage.DownvotedResponses](store, storage.NewQueryBuilder().Gt("created_at", since))
	if err != nil {
		return nil, err
	}

	upvotes, err := upvotesTask.Get()
	if err != nil {
		return nil, err
	}

	feedback := make([]Feedback, 0, len(upvotes)+len(downvotes))

	for _, upvote := range upvotes {
		feedback = append(feedback, Feedback{
			id:        upvote.ID,
			createdAt: *upvote.CreatedAt,
			requestID: upvote.RequestID,
			reason:    upvote.Reason,
			userID:    upvote.UserID,
			isUpvote:  true,
		})
	}

	for _, downvote := range downvotes {
		feedback = append(feedback, Feedback{
			id:        downvote.ID,
			createdAt: *downvote.CreatedAt,
			requestID: downvote.RequestID,
			reason:    downvote.Reason,
			userID:    downvote.UserID,
			isUpvote:  false,
		})
	}

	return feedback, nil
}

func promptFrom(enrichedFeedback []EnrichedFeedback, currentPrompt string) string {
	prompt := fmt.Sprintf(`You are an expert agent LLM prompt finetuner. You are tasked with adjusting the prompt of the following Imperial College London agent: %s. The user interacts with this agent via a chatbot interface. The user sends a message, and the chatbot responds with an answer.

	`, enrichedFeedback[0].AgentID)

	for _, enrichedFeedback := range enrichedFeedback {
		prompt += specificEnrichedFeedbackPrompt(enrichedFeedback)
	}

	prompt += fmt.Sprintf(`
	Here is the current prompt: %s. You are free to adjust this prompt to improve the agent's performance, however you must only make minor changes IF NECESSARY. There is no obligation to make any changes - it is perfectly fine to leave the prompt as is. Respond with the just the new prompt (or the old prompt if you don't want to make any changes), no other text.
	`, currentPrompt)

	return prompt
}
func specificEnrichedFeedbackPrompt(enrichedFeedback EnrichedFeedback) string {
	return fmt.Sprintf(
		`The user and agent had this interaction: """%s""". %s. Here are some surrounding messages for context: %+v
	`, enrichedFeedback.MessageResponse, promptFromFeedback(enrichedFeedback.Feedback), enrichedFeedback.Context)
}

func promptFromFeedback(feedback Feedback) string {
	if feedback.reason == "" {
		if feedback.isUpvote {
			return "The user liked the response."
		} else {
			return "The user disliked the response."
		}
	}

	if feedback.isUpvote {
		return "The user liked the response and left this feedback: " + feedback.reason
	}

	return "The user disliked the response and left this feedback: " + feedback.reason
}

func groupByAgentID(enrichedFeedback []EnrichedFeedback) map[string][]EnrichedFeedback {
	agentFeedback := make(map[string][]EnrichedFeedback)
	for _, feedback := range enrichedFeedback {
		agentFeedback[feedback.AgentID] = append(agentFeedback[feedback.AgentID], feedback)
	}
	return agentFeedback
}
