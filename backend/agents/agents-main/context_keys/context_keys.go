package context_keys

import "context"

type contextKey string

const (
	ContextKeySessionID contextKey = "session_id"
	ContextKeyRequestID contextKey = "request_id"
)

func SetSessionID(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, ContextKeySessionID, sessionID)
}

func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, ContextKeyRequestID, requestID)
}

func GetSessionID(ctx context.Context) (string, bool) {
	sessionID, ok := ctx.Value(ContextKeySessionID).(string)
	return sessionID, ok
}

func GetRequestID(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(ContextKeyRequestID).(string)
	return requestID, ok
}
