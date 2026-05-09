package engine

import (
	"time"
	"github.com/timan-org/timan/internal/domain"
)

// TimerState holds the current snapshot of the timer's progress and status.
type TimerState struct {
	EventName             string
	CurrentPhaseIndex     int
	RemainingSeconds      int
	IsRunning             bool
	CurrentPhase          domain.Phase
	TotalMinutes          int
	TotalRemainingSeconds int
	TotalDurationSeconds  int
	WarningSeconds        int
	UpcomingPhaseName     string
}

// TimerObserver defines the interface for components that wish to be notified of timer ticks.
type TimerObserver interface {
	OnTick(state TimerState)
}

// TimerEngine manages the core logic of the timer, including phase transitions and state updates.
type TimerEngine struct {
	phases         []domain.Phase
	state          TimerState
	observers      []TimerObserver
	stopChan       chan bool
	currentEventID string
}

// GetCurrentEventID returns the ID of the Event currently loaded in the engine.
func (e *TimerEngine) GetCurrentEventID() string {
	return e.currentEventID
}

// SetCurrentEventID sets the ID of the current Event.
func (e *TimerEngine) SetCurrentEventID(id string) {
	e.currentEventID = id
}

// NewTimerEngine initializes a new TimerEngine with a set of phases and warning threshold.
func NewTimerEngine(name string, phases []domain.Phase, warningSec int) *TimerEngine {
	totalSec := 0
	for _, p := range phases {
		totalSec += p.Duration
	}
	return &TimerEngine{
		phases:   phases,
		stopChan: make(chan bool),
		state: TimerState{
			EventName:             name,
			CurrentPhaseIndex:     0,
			RemainingSeconds:      phases[0].Duration,
			IsRunning:             false,
			CurrentPhase:          phases[0],
			TotalMinutes:          totalSec / 60,
			TotalDurationSeconds:  totalSec,
			TotalRemainingSeconds: totalSec,
			WarningSeconds:        warningSec,
		},
	}
}

// AddObserver registers a new observer to receive timer updates.
func (e *TimerEngine) AddObserver(o TimerObserver) {
	e.observers = append(e.observers, o)
}

// Start begins the timer loop in a separate goroutine.
func (e *TimerEngine) Start() {
	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				if e.state.IsRunning {
					e.tick()
				}
			case <-e.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}

func (e *TimerEngine) tick() {
	if e.state.RemainingSeconds > 0 {
		e.state.RemainingSeconds--
		e.state.TotalRemainingSeconds--
	} else if e.state.CurrentPhaseIndex < len(e.phases)-1 {
		e.state.CurrentPhaseIndex++
		e.state.CurrentPhase = e.phases[e.state.CurrentPhaseIndex]
		e.state.RemainingSeconds = e.state.CurrentPhase.Duration
		e.state.TotalRemainingSeconds--
	} else {
		e.state.IsRunning = false
	}
	e.notify()
}

// Toggle flips the running state of the timer (Start/Pause).
func (e *TimerEngine) Toggle() {
	e.state.IsRunning = !e.state.IsRunning
	e.notify()
}

// Reset returns the timer to its initial state (first phase, paused).
func (e *TimerEngine) Reset() {
	totalSec := 0
	for _, p := range e.phases {
		totalSec += p.Duration
	}
	e.state.TotalDurationSeconds = totalSec
	e.state.TotalRemainingSeconds = totalSec
	e.state.CurrentPhaseIndex = 0
	e.state.CurrentPhase = e.phases[0]
	e.state.RemainingSeconds = e.phases[0].Duration
	e.state.IsRunning = false
	e.notify()
}

// ResetPhase restarts the current phase from the beginning.
func (e *TimerEngine) ResetPhase() {
	oldRemaining := e.state.RemainingSeconds
	e.state.RemainingSeconds = e.phases[e.state.CurrentPhaseIndex].Duration
	e.state.TotalRemainingSeconds += (e.state.RemainingSeconds - oldRemaining)
	e.notify()
}

// UpdatePhases reconfigures the engine with a new set of phases and alert threshold.
func (e *TimerEngine) UpdatePhases(name string, phases []domain.Phase, warningMins int) {
	e.phases = phases
	totalSec := 0
	for _, p := range phases {
		totalSec += p.Duration
	}
	e.state.EventName = name
	e.state.TotalMinutes = totalSec / 60
	e.state.TotalDurationSeconds = totalSec
	e.state.TotalRemainingSeconds = totalSec
	e.state.WarningSeconds = warningMins * 60
	e.Reset()
}

// GetPhases returns the current set of phases.
func (e *TimerEngine) GetPhases() []domain.Phase {
	return e.phases
}

// GetState returns the current state of the timer.
func (e *TimerEngine) GetState() TimerState {
	state := e.state
	state.EventName = e.state.EventName // Explicitly ensure it's copied
	if state.CurrentPhaseIndex < len(e.phases)-1 {
		state.UpcomingPhaseName = e.phases[state.CurrentPhaseIndex+1].Name
	} else {
		state.UpcomingPhaseName = ""
	}
	return state
}

func (e *TimerEngine) notify() {
	state := e.state
	if state.CurrentPhaseIndex < len(e.phases)-1 {
		state.UpcomingPhaseName = e.phases[state.CurrentPhaseIndex+1].Name
	} else {
		state.UpcomingPhaseName = ""
	}

	for _, o := range e.observers {
		o.OnTick(state)
	}
}
