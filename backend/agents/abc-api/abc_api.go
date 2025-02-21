package abc_api

type AbcApi interface {
}

type MockAbcApiClient struct {
}

func NewMockAbcApiClient() *MockAbcApiClient {
	return &MockAbcApiClient{}
}
