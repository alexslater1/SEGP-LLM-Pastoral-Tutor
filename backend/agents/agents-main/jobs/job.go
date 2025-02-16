package jobs

import "time"

type Job interface {
	Interval() time.Duration
	Run()
}
