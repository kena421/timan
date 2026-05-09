package main

import (
	"github.com/timan-org/timan/internal/domain"
	"github.com/timan-org/timan/internal/engine"
	"github.com/timan-org/timan/internal/ui"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

func main() {
	// 1. Initial State: Default Event Phases
	initialPhases := []domain.Phase{
		{Name: "Intro", Duration: 5 * 60},
		{Name: "Main Content", Duration: 45 * 60},
		{Name: "Wrap-up", Duration: 10 * 60},
	}

	// 2. Initialize Engine
	timerEngine := engine.NewTimerEngine(initialPhases)
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
