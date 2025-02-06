package clock

import "time"

type MockClock struct {}

func NewMockClock() *MockClock {
	return &MockClock{}
}

func (c *MockClock) CurrentDateTime() time.Time {
	return time.Date(2025, 2, 4, 12, 0, 0, 0, time.UTC)
}
