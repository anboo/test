package timer

import (
	"time"
)

type Timer struct {
}

func NewTimer() *Timer {
	return &Timer{}
}

func (Timer) Now() time.Time {
	return time.Now()
}
