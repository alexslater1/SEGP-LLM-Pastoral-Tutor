package jobs

import (
	"fmt"
	"time"

	"github.com/segp/agents-main/history"
	"github.com/segp/agents-main/storage"
	"github.com/segp/agents-main/utils"
)

type Feedback struct {
	IsUpvote         bool
	Reason           string
	CompletionResult string
	UserID           string
	RequestID        string
	SessionID        string
}

type AdjustPromptsJob struct {
	store   storage.Storage
	history history.History
}

func NewAdjustPromptsJob(store storage.Storage, history history.History) *AdjustPromptsJob {
	return &AdjustPromptsJob{store: store, history: history}
}

func (a *AdjustPromptsJob) Name() string {
	return "adjust_prompts"
}

func (a *AdjustPromptsJob) Interval() time.Duration {
	return 1 * time.Hour
}

func (a *AdjustPromptsJob) Run() error {
	feedback, err := a.getFeedbck()
	if err != nil {
		return fmt.Errorf("failed to get feedback: %w", err)
	}

	allSessionIds := make([]string, 0, len(feedback))
	for _, feedback := range feedback {
		allSessionIds = append(allSessionIds, feedback.SessionID)
	}

	chatsForSessionIds, err := a.chatsForSessionIds(allSessionIds)
	if err != nil {
		return fmt.Errorf("failed to get surrounding chats: %w", err)
	}

	return nil
}

func (a *AdjustPromptsJob) chatsForSessionIds(sessionIds []string) (map[string][]string, error) {
	chatHistoryTasks := utils.DoAsyncList(sessionIds, func(sessionId string) ([]string, error) {
		return a.history.GetMessageHistory(sessionId)
	})

	chatHistories, err := utils.GetAsyncList(chatHistoryTasks)
	if err != nil {
		return nil, err
	}

	chatHistoriesMap := make(map[string][]string)
	for i, sessionId := range sessionIds {
		chatHistoriesMap[sessionId] = chatHistories[i]
	}

	return chatHistoriesMap, nil
}

func (a *AdjustPromptsJob) getFeedbck() ([]Feedback, error) {
	upvotesTask := utils.DoAsync(func() ([]storage.UpvotedResponses, error) {
		return storage.GetAll[storage.UpvotedResponses](a.store, nil)
	})

	downvotes, err := storage.GetAll[storage.DownvotedResponses](a.store, nil)
	if err != nil {
		return nil, err
	}

	upvotes, err := upvotesTask.Get()
	if err != nil {
		return nil, err
	}

	upvoteSessionIDsTasks := utils.DoAsyncList(upvotes, func(upvote storage.UpvotedResponses) (string, error) {
		requestSessions, err := storage.GetAll[storage.RequestSession](a.store, nil)
		if err != nil {
			return "", err
		}

		if len(requestSessions) != 1 {
			return "", fmt.Errorf("expected 1 request session, got %d", len(requestSessions))
		}

		return requestSessions[0].SessionID, nil
	})

	downvoteSessionIDsTasks := utils.DoAsyncList(downvotes, func(downvote storage.DownvotedResponses) (string, error) {
		requestSessions, err := storage.GetAll[storage.RequestSession](a.store, nil)
		if err != nil {
			return "", err
		}

		if len(requestSessions) != 1 {
			return "", fmt.Errorf("expected 1 request session, got %d", len(requestSessions))
		}

		return requestSessions[0].SessionID, nil
	})

	upvoteCompletionResultTasks := utils.DoAsyncList(upvotes, func(upvote storage.UpvotedResponses) (string, error) {
		completionResults, err := storage.GetAll[storage.CompletionResult](a.store, storage.NewQueryBuilder().Eq("request_id", upvote.RequestID))
		if err != nil {
			return "", err
		}

		if len(completionResults) != 1 {
			return "", fmt.Errorf("expected 1 completion result, got %d", len(completionResults))
		}

		return *completionResults[0].Result, nil
	})

	downvoteCompletionResultTasks := utils.DoAsyncList(downvotes, func(downvote storage.DownvotedResponses) (string, error) {
		completionResults, err := storage.GetAll[storage.CompletionResult](a.store, storage.NewQueryBuilder().Eq("request_id", downvote.RequestID))
		if err != nil {
			return "", err
		}

		if len(completionResults) != 1 {
			return "", fmt.Errorf("expected 1 completion result, got %d", len(completionResults))
		}

		return *completionResults[0].Result, nil
	})

	upvoteCompletionResults, err := utils.GetAsyncList(upvoteCompletionResultTasks)
	if err != nil {
		return nil, err
	}

	downvoteCompletionResults, err := utils.GetAsyncList(downvoteCompletionResultTasks)
	if err != nil {
		return nil, err
	}

	upvoteSessionsIds, err := utils.GetAsyncList(upvoteSessionIDsTasks)
	if err != nil {
		return nil, err
	}

	downvoteSessionsIds, err := utils.GetAsyncList(downvoteSessionIDsTasks)
	if err != nil {
		return nil, err
	}

	feedback := make([]Feedback, len(upvotes)+len(downvotes))
	for i, upvote := range upvotes {
		feedback[i] = Feedback{
			IsUpvote:         true,
			Reason:           upvote.Reason,
			UserID:           upvote.UserID,
			RequestID:        upvote.RequestID,
			SessionID:        upvoteSessionsIds[i],
			CompletionResult: upvoteCompletionResults[i],
		}
	}

	for i, downvote := range downvotes {
		feedback[i+len(upvotes)] = Feedback{
			IsUpvote:         false,
			Reason:           downvote.Reason,
			UserID:           downvote.UserID,
			RequestID:        downvote.RequestID,
			SessionID:        downvoteSessionsIds[i],
			CompletionResult: downvoteCompletionResults[i],
		}
	}

	return feedback, nil
}

func surroundingChatForFeedbackAndHistory(feedback Feedback, history []string) (string, error) {

}
