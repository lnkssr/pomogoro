package kernel

import "time"

type Timer struct {
	Duration time.Duration
	Phase    string
}

func NewTimer(duration time.Duration, phase string) *Timer {
	return &Timer{
		Duration: duration,
		Phase:    phase,
	}
}
