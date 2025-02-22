package context_keys

import (
	"context"
	"testing"
)

func TestSetAndGetSessionID(t *testing.T) {
	const sessionID = "session123"
	ctx := context.Background()

	// Add the session ID to the context.
	ctxWithSession := SetSessionID(ctx, sessionID)

	// Retrieve the session ID from the context.
	got, ok := GetSessionID(ctxWithSession)
	if !ok {
		t.Errorf("Expected session ID %q, got %q", sessionID, got)
	}

	if got != sessionID {
		t.Errorf("Expected session ID %q, got %q", sessionID, got)
	}
}

func TestSetAndGetRequestID(t *testing.T) {
	const requestID = "req456"
	ctx := context.Background()

	// Add the request ID to the context.
	ctxWithRequest := SetRequestID(ctx, requestID)

	// Retrieve the request ID from the context.
	got, ok := GetRequestID(ctxWithRequest)
	if !ok {
		t.Errorf("Expected request ID %q, got %q", requestID, got)
	}

	if got != requestID {
		t.Errorf("Expected request ID %q, got %q", requestID, got)
	}
}

// Test that GetSessionID panics when the key is not set.
func TestGetSessionIDWithoutSetting(t *testing.T) {
	ctx := context.Background()

	// This should panic.
	_, ok := GetSessionID(ctx)
	if ok {
		t.Error("GetSessionID was not false when no session ID was set")
	}
}

// Test that GetRequestID panics when the key is not set.
func TestGetRequestIDWithoutSetting(t *testing.T) {
	ctx := context.Background()

	_, ok := GetRequestID(ctx)
	if ok {
		t.Error("GetRequestID was not false when no request ID was set")
	}
}

func TestSetAndGetUserID(t *testing.T) {
	const userID = "user123"
	ctx := context.Background()

	// Add the user ID to the context.
	ctxWithUser := SetUserID(ctx, userID)

	// Retrieve the user ID from the context.
	got, ok := GetUserID(ctxWithUser)
	if !ok {
		t.Errorf("Expected user ID %q, got %q", userID, got)
	}
}

func TestSetAndGetAgentID(t *testing.T) {
	const agentID = "agent123"
	ctx := context.Background()

	// Add the agent ID to the context.
	ctxWithAgent := SetAgentID(ctx, agentID)

	// Retrieve the agent ID from the context.
	got, ok := GetAgentID(ctxWithAgent)
	if !ok {
		t.Errorf("Expected agent ID %q, got %q", agentID, got)
	}
}
