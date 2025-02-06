package clock

import "time"

type RealClock struct{}

func NewRealClock() *RealClock {
	return &RealClock{}
}

func (c *RealClock) CurrentDateTime() time.Time {
	return time.Now()
}
