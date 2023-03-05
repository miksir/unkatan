package schedule

import (
	"time"
)

type Timer struct {
	done  chan bool
	timer *time.Ticker
	call  func()
}

func RunTimer(d time.Duration, call func()) *Timer {
	timer := Timer{
		done:  make(chan bool),
		timer: time.NewTicker(d),
		call:  call,
	}
	go timer.timerTick()
	return &timer
}

func (t *Timer) StopTimer() {
	t.done <- true
	close(t.done)
	t.timer = nil
	t.call = nil
}

func (t *Timer) timerTick() {
	defer t.timer.Stop()
	for {
		select {
		case <-t.timer.C:
			t.call()
		case <-t.done:
			return
		}
	}
}
