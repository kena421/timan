package engine

import (
	"time"
	"timan/internal/domain"
)

type TimerState struct {
	CurrentPhaseIndex int
	RemainingSeconds  int
	IsRunning         bool
	CurrentPhase      domain.Phase
	TotalMinutes      int
}

type TimerObserver interface {
	OnTick(state TimerState)
}

type TimerEngine struct {
	phases      []domain.Phase
	state       TimerState
	observers   []TimerObserver
	stopChan    chan bool
}

func NewTimerEngine(phases []domain.Phase) *TimerEngine {
	total := 0
	for _, p := range phases {
		total += p.Duration
	}
	return &TimerEngine{
		phases:   phases,
		stopChan: make(chan bool),
		state: TimerState{
			CurrentPhaseIndex: 0,
			RemainingSeconds:  phases[0].Duration,
			IsRunning:         false,
			CurrentPhase:      phases[0],
			TotalMinutes:      total / 60,
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
	} else if e.state.CurrentPhaseIndex < len(e.phases)-1 {
		e.state.CurrentPhaseIndex++
		e.state.CurrentPhase = e.phases[e.state.CurrentPhaseIndex]
		e.state.RemainingSeconds = e.state.CurrentPhase.Duration
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
	e.state.CurrentPhaseIndex = 0
	e.state.CurrentPhase = e.phases[0]
	e.state.RemainingSeconds = e.phases[0].Duration
	e.state.IsRunning = false
	e.notify()
}

func (e *TimerEngine) ResetPhase() {
	e.state.RemainingSeconds = e.phases[e.state.CurrentPhaseIndex].Duration
	e.notify()
}

func (e *TimerEngine) UpdatePhases(phases []domain.Phase) {
	e.phases = phases
	total := 0
	for _, p := range phases {
		total += p.Duration
	}
	e.state.TotalMinutes = total / 60
	e.Reset()
}

func (e *TimerEngine) GetPhases() []domain.Phase {
	return e.phases
}

func (e *TimerEngine) notify() {
	for _, o := range e.observers {
		o.OnTick(e.state)
	}
}
