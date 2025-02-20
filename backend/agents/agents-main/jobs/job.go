package jobs

import "time"

type Job interface {
	Name() string
	Interval() time.Duration
	Run() error
}
