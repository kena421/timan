package main

import (
	"timan/internal/domain"
	"timan/internal/engine"
	"timan/internal/ui"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

func main() {
	// 1. Domain / Initial State
	initialPhases := []domain.Phase{
		{Name: "Intro", Duration: 5 * 60},
		{Name: "Understand", Duration: 10 * 60},
		{Name: "Design", Duration: 20 * 60},
		{Name: "Discussion", Duration: 25 * 60},
	}

	// 2. Initialize Engine
	timerEngine := engine.NewTimerEngine(initialPhases)
	timerEngine.Start()

	// 3. Initialize App & UI
	a := app.NewWithID("com.interview.timer.v3")
	a.SetIcon(theme.SettingsIcon())
	w := a.NewWindow("InterviewTimer")
	w.SetFixedSize(true)

	// Dependency Injection: UI depends on the engine
	timerUI := ui.NewTimerUI(w, timerEngine)

	// 4. Show (handles internal platform tweaks)
	timerUI.Show()

	// 5. Run
	a.Run()
}
