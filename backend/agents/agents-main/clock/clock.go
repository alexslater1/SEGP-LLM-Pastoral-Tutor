package clock

import "time"

type Clock interface {
	CurrentDateTime() time.Time
}