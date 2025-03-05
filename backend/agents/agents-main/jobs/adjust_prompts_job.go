package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/segp/agents-main/history"
	"github.com/segp/agents-main/storage"
	"github.com/segp/agents-main/utils"
)

const (
	name     = "adjust_prompts"
	interval = 1 * time.Hour

	numChatsToConsider = 5
)

type AdjustPromptsJob struct {
	store   storage.Storage
	history history.History
}

func NewAdjustPromptsJob(store storage.Storage, history history.History) *AdjustPromptsJob {
	return &AdjustPromptsJob{store: store, history: history}
}

func (j *AdjustPromptsJob) Name() string {
	return name
}

func (j *AdjustPromptsJob) Interval() time.Duration {
	return interval
}

func (j *AdjustPromptsJob) Run(ctx context.Context) error {
	feedback, err := getFeedback(j.store)
	if err != nil {
		return err
	}

	enrichedFeedback, err := j.getEnrichedFeedback(feedback)
	if err != nil {
		return err
	}

	fmt.Println(enrichedFeedback)

	return nil
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

func getFeedback(store storage.Storage) ([]Feedback, error) {
	upvotesTask := utils.DoAsync(func() ([]storage.UpvotedResponses, error) {
		return storage.GetAll[storage.UpvotedResponses](store, nil)
	})

	downvotes, err := storage.GetAll[storage.DownvotedResponses](store, nil)
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
