package pkg

import "time"

type Timer struct{}

func NewTimer() *Timer {
	return &Timer{}
}

func (t *Timer) NowUTC() time.Time {
	return time.Now().UTC()
}
