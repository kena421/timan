package main

import (
	"github.com/timan-org/timan/internal/domain"
	"github.com/timan-org/timan/internal/engine"
	"github.com/timan-org/timan/internal/ui"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

func main() {
	// 1. Load Persistence
	store := domain.NewEventStore()
	events, _ := store.LoadAll()
	state := store.LoadState()

	initialPhases := []domain.Phase{
		{Name: "Intro", Duration: 5 * 60},
		{Name: "Main Content", Duration: 45 * 60},
		{Name: "Wrap-up", Duration: 10 * 60},
	}
	lastID := ""

	// Find the last used event OR the first available event
	targetID := state.LastEventID
	if targetID == "" && len(events) > 0 {
		targetID = events[0].ID
	}

	if targetID != "" {
		for _, e := range events {
			if e.ID == targetID {
				initialPhases = e.Phases
				lastID = e.ID
				break
			}
		}
	}

	// 2. Initialize Engine
	timerEngine := engine.NewTimerEngine(initialPhases)
	if lastID != "" {
		timerEngine.SetCurrentEventID(lastID)
	}
	timerEngine.Start()

	// 3. Initialize App & UI
	// Using a generic app ID
	a := app.NewWithID("com.timan.timer")
	a.SetIcon(theme.SettingsIcon())
	
	// Main Window
	w := a.NewWindow("Timan")
	w.SetFixedSize(true)

	// Dependency Injection: UI depends on the engine
	timerUI := ui.NewTimerUI(w, timerEngine)

	// 4. Show (handles internal platform tweaks)
	timerUI.Show()

	// 5. Run
	a.Run()
}
