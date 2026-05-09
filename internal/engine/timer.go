package engine

import (
	"time"
	"github.com/kena421/timan/internal/domain"
)

// TimerState holds a snapshot of the current timer status at any point in time.
type TimerState struct {
	EventName             string       // Name of the active event profile
	CurrentPhaseIndex     int          // Index of the active phase in the slice
	RemainingSeconds      int          // Seconds left in the current active phase
	IsRunning             bool         // True if the timer is actively counting down
	CurrentPhase          domain.Phase // Full details of the current phase
	TotalMinutes          int          // Total duration of the event in minutes
	TotalRemainingSeconds int          // Sum of all remaining seconds across all phases
	TotalDurationSeconds  int          // Sum of all phase durations in seconds
	WarningSeconds        int          // Threshold in seconds for triggering alert (red) mode
	UpcomingPhaseName     string       // Name of the next phase in the sequence (if any)
}

// TimerObserver defines the interface for components (like UI) that need to
// react to every second tick of the timer engine.
type TimerObserver interface {
	OnTick(state TimerState)
}

// TimerEngine is the core state machine that manages phase transitions,
// timing logic, and observer notifications.
type TimerEngine struct {
	phases         []domain.Phase  // The ordered sequence of phases to execute
	state          TimerState     // The current internal state
	observers      []TimerObserver // Registered UI components
	stopChan       chan bool      // Channel to gracefully shutdown the ticker
	currentEventID string         // ID of the event profile currently loaded
}

// GetCurrentEventID returns the persistent ID of the loaded event profile.
func (e *TimerEngine) GetCurrentEventID() string {
	return e.currentEventID
}

// SetCurrentEventID updates the engine with the ID of the selected profile.
func (e *TimerEngine) SetCurrentEventID(id string) {
	e.currentEventID = id
}

// NewTimerEngine initializes the engine with a set of phases and an alert threshold.
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

// AddObserver registers a new listener to receive 1Hz state updates.
func (e *TimerEngine) AddObserver(o TimerObserver) {
	e.observers = append(e.observers, o)
}

// Start initiates the internal 1-second ticker loop in a background goroutine.
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

// tick handles the decrement of counters and phase transition logic.
func (e *TimerEngine) tick() {
	if e.state.RemainingSeconds > 0 {
		e.state.RemainingSeconds--
		e.state.TotalRemainingSeconds--
	} else if e.state.CurrentPhaseIndex < len(e.phases)-1 {
		// Transition to next phase
		e.state.CurrentPhaseIndex++
		e.state.CurrentPhase = e.phases[e.state.CurrentPhaseIndex]
		e.state.RemainingSeconds = e.state.CurrentPhase.Duration
		e.state.TotalRemainingSeconds--
	} else {
		// End of all phases
		e.state.IsRunning = false
	}
	e.notify()
}

// Toggle pauses or resumes the countdown.
func (e *TimerEngine) Toggle() {
	e.state.IsRunning = !e.state.IsRunning
	e.notify()
}

// Reset restores the timer to the beginning of the first phase.
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

// ResetPhase restarts only the current phase from its maximum duration.
func (e *TimerEngine) ResetPhase() {
	oldRemaining := e.state.RemainingSeconds
	e.state.RemainingSeconds = e.phases[e.state.CurrentPhaseIndex].Duration
	e.state.TotalRemainingSeconds += (e.state.RemainingSeconds - oldRemaining)
	e.notify()
}

// UpdatePhases re-initializes the engine with a completely new event profile.
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

// GetPhases returns the current ordered list of phase templates.
func (e *TimerEngine) GetPhases() []domain.Phase {
	return e.phases
}

// GetState returns a calculated snapshot of the current state for UI consumption.
func (e *TimerEngine) GetState() TimerState {
	state := e.state
	state.EventName = e.state.EventName 
	if state.CurrentPhaseIndex < len(e.phases)-1 {
		state.UpcomingPhaseName = e.phases[state.CurrentPhaseIndex+1].Name
	} else {
		state.UpcomingPhaseName = ""
	}
	return state
}

// notify broadcasts the current state to all registered observers.
func (e *TimerEngine) notify() {
	state := e.GetState()
	for _, o := range e.observers {
		o.OnTick(state)
	}
}
