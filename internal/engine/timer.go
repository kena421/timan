package engine

import (
	"time"
	"timan/internal/domain"
)

type TimerState struct {
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

type TimerObserver interface {
	OnTick(state TimerState)
}

type TimerEngine struct {
	phases             []domain.Phase
	state              TimerState
	observers          []TimerObserver
	stopChan           chan bool
	currentBlueprintID string
}

func (e *TimerEngine) GetCurrentBlueprintID() string {
	return e.currentBlueprintID
}

func (e *TimerEngine) SetCurrentBlueprintID(id string) {
	e.currentBlueprintID = id
}

func NewTimerEngine(phases []domain.Phase) *TimerEngine {
	totalSec := 0
	for _, p := range phases {
		totalSec += p.Duration
	}
	return &TimerEngine{
		phases:   phases,
		stopChan: make(chan bool),
		state: TimerState{
			CurrentPhaseIndex:     0,
			RemainingSeconds:      phases[0].Duration,
			IsRunning:             false,
			CurrentPhase:          phases[0],
			TotalMinutes:          totalSec / 60,
			TotalDurationSeconds:  totalSec,
			TotalRemainingSeconds: totalSec,
		},
	}
}

func (e *TimerEngine) AddObserver(o TimerObserver) {
	e.observers = append(e.observers, o)
}

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

func (e *TimerEngine) Toggle() {
	e.state.IsRunning = !e.state.IsRunning
	e.notify()
}

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

func (e *TimerEngine) ResetPhase() {
	// Need to adjust TotalRemainingSeconds when resetting a phase
	oldRemaining := e.state.RemainingSeconds
	e.state.RemainingSeconds = e.phases[e.state.CurrentPhaseIndex].Duration
	e.state.TotalRemainingSeconds += (e.state.RemainingSeconds - oldRemaining)
	e.notify()
}

func (e *TimerEngine) UpdatePhases(phases []domain.Phase, warningMins int) {
	e.phases = phases
	totalSec := 0
	for _, p := range phases {
		totalSec += p.Duration
	}
	e.state.TotalMinutes = totalSec / 60
	e.state.TotalDurationSeconds = totalSec
	e.state.TotalRemainingSeconds = totalSec
	e.state.WarningSeconds = warningMins * 60
	e.Reset()
}

func (e *TimerEngine) GetPhases() []domain.Phase {
	return e.phases
}

func (e *TimerEngine) GetState() TimerState {
	state := e.state
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
