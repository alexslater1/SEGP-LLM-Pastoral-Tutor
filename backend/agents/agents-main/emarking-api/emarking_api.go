package emarking_api

type EmarkingApi interface {
}

type MockEmarkingApiClient struct {
}

func NewMockAbcApiClient() *MockEmarkingApiClient {
	return &MockEmarkingApiClient{}
}