package emarking_api

type EmarkingApiEndpoint string

const (
	GetExercises       EmarkingApiEndpoint = "get-exercises"
	GetFeedback        EmarkingApiEndpoint = "get-feedback"
	GetExerciseSummary EmarkingApiEndpoint = "get-exercise-summary"
	GetSubmissionGroup EmarkingApiEndpoint = "get-submission-group"
)

type EmarkingApi interface {
}

type MockEmarkingApiClient struct {
}

func NewMockEmarkingApiClient() *MockEmarkingApiClient {
	return &MockEmarkingApiClient{}
}

func (e *MockEmarkingApiClient) AllApiEndpoints() []EmarkingApiEndpoint {
	return []EmarkingApiEndpoint{
		GetExercises,
		GetFeedback,
		GetExerciseSummary,
		GetSubmissionGroup,
	}
}

func (e *MockEmarkingApiClient) EndpointMappings() map[EmarkingApiEndpoint]func() any {
	return map[EmarkingApiEndpoint]func() any{
		GetExercises:       func() any { return e.GetExercises() },
		GetFeedback:        func() any { return e.GetFeedback() },
		GetExerciseSummary: func() any { return e.GetExerciseSummary() },
		GetSubmissionGroup: func() any { return e.GetSubmissionGroup() },
	}
}
