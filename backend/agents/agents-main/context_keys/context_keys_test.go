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
	got := GetSessionID(ctxWithSession)
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
	got := GetRequestID(ctxWithRequest)
	if got != requestID {
		t.Errorf("Expected request ID %q, got %q", requestID, got)
	}
}

// Test that GetSessionID panics when the key is not set.
func TestGetSessionIDWithoutSetting(t *testing.T) {
	ctx := context.Background()
	defer func() {
		if r := recover(); r == nil {
			t.Error("GetSessionID did not panic when no session ID was set")
		}
	}()
	// This should panic.
	_ = GetSessionID(ctx)
}

// Test that GetRequestID panics when the key is not set.
func TestGetRequestIDWithoutSetting(t *testing.T) {
	ctx := context.Background()
	defer func() {
		if r := recover(); r == nil {
			t.Error("GetRequestID did not panic when no request ID was set")
		}
	}()
	// This should panic.
	_ = GetRequestID(ctx)
}
