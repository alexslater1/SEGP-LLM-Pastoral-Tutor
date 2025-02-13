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

func GetSessionID(ctx context.Context) string {
	return ctx.Value(ContextKeySessionID).(string)
}

func GetRequestID(ctx context.Context) string {
	return ctx.Value(ContextKeyRequestID).(string)
}
