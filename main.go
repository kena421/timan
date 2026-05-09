package main

import (
	"github.com/kena421/timan/internal/domain"
	"github.com/kena421/timan/internal/engine"
	"github.com/kena421/timan/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

// main is the entry point for the Timan HUD utility.
// It initializes the persistence layer, restores application state,
// bootstraps the timing engine, and launches the graphical HUD.
func main() {
	// 1. Initialize Application & Window
	myApp := app.NewWithID("com.timan.timer")
	window := myApp.NewWindow("Timan")
	
	// 2. Initialize Persistence Layer
	store := domain.NewEventStore()
	state, _ := store.LoadState()
	
	// 3. Prepare Default/Last-Used Profile
	// Default to 60 minutes if no library exists
	defaultPhases := []domain.Phase{{Name: "Session", Duration: 3600}}
	eventName := "Quick Timer"
	warningSec := 300 // 5 minute default alert
	
	// Try to restore the last used profile from the library
	if state.LastEventID != "" {
		events, _ := store.LoadAll()
		for _, e := range events {
			if e.ID == state.LastEventID {
				defaultPhases = e.Phases
				eventName = e.Name
				// Calculate alert threshold
				wMins := e.WarningValue
				if e.WarningIsPercent {
					wMins = (e.WarningValue * e.Total) / 100
				}
				warningSec = wMins * 60
				break
			}
		}
	}

	// 4. Initialize Core Engine & UI
	tm := engine.NewTimerEngine(eventName, defaultPhases, warningSec)
	tm.SetCurrentEventID(state.LastEventID)
	
	// TimerUI acts as an observer to the engine
	timerUI := ui.NewTimerUI(window, tm)
	
	// 5. Start Background Ticker
	tm.Start()
	
	// 6. Launch HUD and enter main event loop
	timerUI.Show()
	myApp.Run()
}
