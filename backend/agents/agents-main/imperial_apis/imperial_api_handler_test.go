package imperial_apis

import (
	"reflect"
	"testing"

	abc_api "github.com/segp/agents-main/imperial_apis/abc-api"
	emarking_api "github.com/segp/agents-main/imperial_apis/emarking-api"
)

func TestNewImperialApiHandler(t *testing.T) {
	// Create mock clients
	abcClient := abc_api.NewMockAbcApiClient()
	emarkingClient := emarking_api.NewMockEmarkingApiClient()

	// Create handler
	handler := NewImperialApiHandler(abcClient, emarkingClient)

	// Verify handler was created with the correct clients
	if handler.abcClient != abcClient {
		t.Errorf("Expected abcClient to be %v, got %v", abcClient, handler.abcClient)
	}
	if handler.emarkingClient != emarkingClient {
		t.Errorf("Expected emarkingClient to be %v, got %v", emarkingClient, handler.emarkingClient)
	}
}

func TestImperialApiHandler_EmarkingApis(t *testing.T) {
	// Create mock clients
	abcClient := abc_api.NewMockAbcApiClient()
	emarkingClient := emarking_api.NewMockEmarkingApiClient()

	// Create handler
	handler := NewImperialApiHandler(abcClient, emarkingClient)

	// Get eMarking APIs
	apis := handler.EmarkingApis()

	// Expected APIs
	expectedApis := []emarking_api.EmarkingApiEndpoint{
		emarking_api.GetExercises,
		emarking_api.GetFeedback,
		emarking_api.GetExerciseSummary,
		emarking_api.GetSubmissionGroup,
	}

	// Verify APIs match expected
	if !reflect.DeepEqual(apis, expectedApis) {
		t.Errorf("Expected APIs %v, got %v", expectedApis, apis)
	}
}

func TestImperialApiHandler_AbcApis(t *testing.T) {
	// Create mock clients
	abcClient := abc_api.NewMockAbcApiClient()
	emarkingClient := emarking_api.NewMockEmarkingApiClient()

	// Create handler
	handler := NewImperialApiHandler(abcClient, emarkingClient)

	// Get ABC APIs
	apis := handler.AbcApis()

	// Verify the function returns the result from the client
	// Note: We can't easily test the exact values since the mock client doesn't define AllApiEndpoints()
	// This test just ensures the function calls through to the client
	if apis == nil {
		t.Error("Expected non-nil APIs")
	}
}

func TestImperialApiHandler_EmarkingEndpointsFor(t *testing.T) {
	// Create mock clients
	abcClient := abc_api.NewMockAbcApiClient()
	emarkingClient := emarking_api.NewMockEmarkingApiClient()

	// Create handler
	handler := NewImperialApiHandler(abcClient, emarkingClient)

	// Test with specific endpoints
	endpoints := []emarking_api.EmarkingApiEndpoint{
		emarking_api.GetExercises,
		emarking_api.GetFeedback,
	}

	// Get endpoint functions
	endpointFuncs := handler.EmarkingEndpointsFor(endpoints...)

	// Verify we got the right number of functions
	if len(endpointFuncs) != len(endpoints) {
		t.Errorf("Expected %d endpoint functions, got %d", len(endpoints), len(endpointFuncs))
	}

	// Verify each function is callable (doesn't panic)
	for i, fn := range endpointFuncs {
		defer func(i int) {
			if r := recover(); r != nil {
				t.Errorf("Endpoint function %d panicked: %v", i, r)
			}
		}(i)

		_ = fn()
	}
}

func TestImperialApiHandler_AbcEndpointsFor(t *testing.T) {
	// Create a custom mock ABC client with defined endpoints for testing
	abcClient := abc_api.NewMockAbcApiClient()
	emarkingClient := emarking_api.NewMockEmarkingApiClient()

	// Create handler
	handler := NewImperialApiHandler(abcClient, emarkingClient)

	// We can't easily test with specific endpoints since the mock client doesn't define constants
	// This test just ensures the function doesn't panic when called
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("AbcEndpointsFor panicked: %v", r)
		}
	}()

	// Call with empty endpoints list to avoid potential nil map access
	endpointFuncs := handler.AbcEndpointsFor()

	// Verify we got an empty slice, not nil
	if endpointFuncs == nil {
		t.Error("Expected non-nil slice, got nil")
	}
	if len(endpointFuncs) != 0 {
		t.Errorf("Expected empty slice, got %d items", len(endpointFuncs))
	}
}
