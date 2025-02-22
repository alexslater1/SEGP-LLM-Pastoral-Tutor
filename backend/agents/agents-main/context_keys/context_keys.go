package context_keys

import "context"

type contextKey string

const (
	ContextKeySessionID contextKey = "session_id"
	ContextKeyRequestID contextKey = "request_id"
	ContextKeyUserID    contextKey = "user_id"
	ContextKeyAgentID   contextKey = "agent_id"
)

func SetSessionID(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, ContextKeySessionID, sessionID)
}

func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, ContextKeyRequestID, requestID)
}

func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, ContextKeyUserID, userID)
}

func SetAgentID(ctx context.Context, agentID string) context.Context {
	return context.WithValue(ctx, ContextKeyAgentID, agentID)
}

func GetSessionID(ctx context.Context) (string, bool) {
	sessionID, ok := ctx.Value(ContextKeySessionID).(string)
	return sessionID, ok
}

func GetRequestID(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(ContextKeyRequestID).(string)
	return requestID, ok
}

func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(ContextKeyUserID).(string)
	return userID, ok
}

func GetAgentID(ctx context.Context) (string, bool) {
	agentID, ok := ctx.Value(ContextKeyAgentID).(string)
	return agentID, ok
}
