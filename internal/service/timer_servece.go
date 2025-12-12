package service

import (
	"time"

	"pomogoro/internal/kernel"
)

type TimerUseCase struct {
	WorkDuration  time.Duration
	BreakDuration time.Duration
	Cycles        int
}

func NewTimerUseCase(work, br time.Duration, cycles int) *TimerUseCase {
	return &TimerUseCase{
		WorkDuration:  work,
		BreakDuration: br,
		Cycles:        cycles,
	}
}

func (t *TimerUseCase) Start(onTick func(time.Duration, string), onPhaseEnd func(string)) {
	for i := 0; i < t.Cycles; i++ {
		workTimer := kernel.NewTimer(t.WorkDuration, "Work")
		t.runTimer(workTimer, onTick, onPhaseEnd)

		breakTimer := kernel.NewTimer(t.BreakDuration, "Break")
		t.runTimer(breakTimer, onTick, onPhaseEnd)
	}
}

func (t *TimerUseCase) runTimer(timer *kernel.Timer, onTick func(time.Duration, string), onPhaseEnd func(string)) {
	end := time.Now().Add(timer.Duration)
	for time.Now().Before(end) {
		remain := time.Until(time.Now())
		onTick(remain, timer.Phase)
		time.Sleep(1 * time.Second)
	}
	onPhaseEnd(timer.Phase)
}
