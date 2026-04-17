package backend

import (
	"fmt"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type TimerState string

const (
	StateStopped TimerState = "stopped"
	StateRunning TimerState = "running"
	StatePaused  TimerState = "paused"
)

type SessionType string

const (
	SessionFocus      SessionType = "focus"
	SessionShortBreak SessionType = "short_break"
	SessionLongBreak  SessionType = "long_break"
)

type TimerTick struct {
	RemainingSeconds int         `json:"remaining_seconds"`
	TotalSeconds     int         `json:"total_seconds"`
	State            TimerState  `json:"state"`
	SessionType      SessionType `json:"session_type"`
	FormattedTime    string      `json:"formatted_time"`
}

type TimerService struct {
	mu               sync.Mutex
	state            TimerState
	sessionType      SessionType
	remainingSeconds int
	totalSeconds     int
	ticker           *time.Ticker
	stopChan         chan struct{}
}

func NewTimerService() *TimerService {
	ts := &TimerService{
		state:       StateStopped,
		sessionType: SessionFocus,
	}
	ts.setDurationForSession(SessionFocus)
	return ts
}

func (ts *TimerService) setDurationForSession(sType SessionType) {
	ts.sessionType = sType
	switch sType {
	case SessionFocus:
		ts.totalSeconds = 25 * 60
	case SessionShortBreak:
		ts.totalSeconds = 5 * 60
	case SessionLongBreak:
		ts.totalSeconds = 15 * 60
	}
	ts.remainingSeconds = ts.totalSeconds
}

func (ts *TimerService) Play() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if ts.state == StateRunning {
		return
	}

	if ts.state == StateStopped {
		ts.setDurationForSession(ts.sessionType)
	}

	ts.state = StateRunning
	ts.startTicker()
	ts.emitTick()
}

func (ts *TimerService) Pause() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if ts.state == StateRunning {
		ts.state = StatePaused
		ts.stopTicker()
		ts.emitTick()
	}
}

func (ts *TimerService) Reset() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.stopTicker()
	ts.state = StateStopped
	ts.setDurationForSession(ts.sessionType)
	ts.emitTick()
}

func (ts *TimerService) Skip() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.stopTicker()
	ts.state = StateStopped

	if ts.sessionType == SessionFocus {
		ts.setDurationForSession(SessionShortBreak)
	} else {
		ts.setDurationForSession(SessionFocus)
	}
	ts.emitTick()
}

func (ts *TimerService) startTicker() {
	ts.ticker = time.NewTicker(1 * time.Second)
	ts.stopChan = make(chan struct{})

	go func() {
		for {
			select {
			case <-ts.ticker.C:
				ts.tick()
			case <-ts.stopChan:
				return
			}
		}
	}()
}

func (ts *TimerService) stopTicker() {
	if ts.ticker != nil {
		ts.ticker.Stop()
		close(ts.stopChan)
		ts.ticker = nil
	}
}

func (ts *TimerService) tick() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if ts.state != StateRunning {
		return
	}

	ts.remainingSeconds--
	if ts.remainingSeconds <= 0 {
		ts.remainingSeconds = 0
		ts.state = StateStopped
		ts.stopTicker()

		app := application.Get()
		if app != nil {
			app.Event.Emit("timer_complete", string(ts.sessionType))
		}

		if ts.sessionType == SessionFocus {
			ts.setDurationForSession(SessionShortBreak)
		} else {
			ts.setDurationForSession(SessionFocus)
		}
	}

	ts.emitTick()
}

func (ts *TimerService) emitTick() {
	app := application.Get()
	if app == nil {
		return
	}

	m := ts.remainingSeconds / 60
	s := ts.remainingSeconds % 60
	formatted := fmt.Sprintf("%02d:%02d", m, s)

	tickInfo := TimerTick{
		RemainingSeconds: ts.remainingSeconds,
		TotalSeconds:     ts.totalSeconds,
		State:            ts.state,
		SessionType:      ts.sessionType,
		FormattedTime:    formatted,
	}

	app.Event.Emit("timer_tick", tickInfo)
}

func (ts *TimerService) SetSessionType(sType string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.stopTicker()
	ts.state = StateStopped
	ts.setDurationForSession(SessionType(sType))
	ts.emitTick()
}
