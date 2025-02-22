package emarking_api

type EmarkingApi interface {
}

type MockEmarkingApiClient struct {
}

func NewMockEmarkingApiClient() *MockEmarkingApiClient {
	return &MockEmarkingApiClient{}
}