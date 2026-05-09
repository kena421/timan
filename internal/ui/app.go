package ui

import (
	"fmt"
	"image/color"
	"strings"

	"timan/internal/engine"
	"timan/internal/platform"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type TimerUI struct {
	window     fyne.Window
	phaseLabel *canvas.Text
	timerLabel *canvas.Text
	totalLabel *canvas.Text
	progress   *widget.ProgressBar
	background *canvas.Rectangle

	engine    *engine.TimerEngine
	dashboard *Dashboard
}

func NewTimerUI(w fyne.Window, e *engine.TimerEngine) *TimerUI {
	ui := &TimerUI{
		window: w,
		engine: e,
	}
	ui.dashboard = NewDashboard(ui, e)
	ui.setup()
	e.AddObserver(ui)
	return ui
}

func (ui *TimerUI) setup() {
	ui.phaseLabel = canvas.NewText("INITIALIZING", color.NRGBA{R: 200, G: 200, B: 200, A: 255})
	ui.phaseLabel.TextSize = 10
	ui.phaseLabel.Alignment = fyne.TextAlignCenter

	ui.timerLabel = canvas.NewText("00:00", color.NRGBA{R: 200, G: 200, B: 200, A: 255})
	ui.timerLabel.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
	ui.timerLabel.TextSize = 32
	ui.timerLabel.Alignment = fyne.TextAlignCenter

	ui.totalLabel = canvas.NewText("PHASE: 00:00 / 00:00", color.NRGBA{R: 150, G: 150, B: 150, A: 255})
	ui.totalLabel.TextSize = 10
	ui.totalLabel.Alignment = fyne.TextAlignCenter

	ui.progress = widget.NewProgressBar()
	ui.progress.TextFormatter = func() string { return "" }

	ui.background = canvas.NewRectangle(color.NRGBA{R: 30, G: 30, B: 30, A: 200})

	menu := fyne.NewMenu("",
		fyne.NewMenuItem("Open Dashboard", ui.dashboard.Show),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Reset current phase", ui.engine.ResetPhase),
		fyne.NewMenuItem("Reset All", ui.engine.Reset),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Quit", func() { fyne.CurrentApp().Quit() }),
	)

	content := container.NewStack(
		ui.background,
		container.NewVBox(
			ui.phaseLabel,
			container.NewCenter(ui.timerLabel),
			ui.totalLabel,
			ui.progress,
		),
		&InteractionWrapper{
			OnTap:  ui.engine.Toggle,
			Menu:   menu,
			Window: ui.window,
		},
	)

	ui.window.SetContent(content)
	ui.window.Resize(fyne.NewSize(220, 110))
}

func (ui *TimerUI) formatTime(s int) string {
	return fmt.Sprintf("%02d:%02d", s/60, s%60)
}

func (ui *TimerUI) OnTick(state engine.TimerState) {
	ui.phaseLabel.Text = strings.ToUpper(state.CurrentPhase.Name)
	
	// Main Focus: Total Session Time (Remaining / Total)
	ui.timerLabel.Text = fmt.Sprintf("%s / %s", ui.formatTime(state.TotalRemainingSeconds), ui.formatTime(state.TotalDurationSeconds))
	
	// Secondary: Current Phase Time (Remaining / Total)
	ui.totalLabel.Text = fmt.Sprintf("PHASE: %s / %s", ui.formatTime(state.RemainingSeconds), ui.formatTime(state.CurrentPhase.Duration))
	
	// Progress bar reflects Total Session Progress
	ui.progress.Max = float64(state.TotalDurationSeconds)
	ui.progress.Value = float64(state.TotalDurationSeconds - state.TotalRemainingSeconds)
	
	if state.TotalRemainingSeconds < 300 && state.IsRunning { // Red in last 5 mins
		ui.timerLabel.Color = color.NRGBA{R: 255, G: 100, B: 0, A: 255}
	} else if !state.IsRunning {
		ui.timerLabel.Color = color.NRGBA{R: 200, G: 200, B: 200, A: 255}
	} else {
		ui.timerLabel.Color = color.NRGBA{R: 50, G: 255, B: 50, A: 255}
	}

	ui.phaseLabel.Refresh()
	ui.timerLabel.Refresh()
	ui.totalLabel.Refresh()
	ui.progress.Refresh()
}

func (ui *TimerUI) Show() {
	ui.window.Show()
	platform.TweakWindow("InterviewTimer")
}
