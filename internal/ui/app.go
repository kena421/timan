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
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type TimerUI struct {
	window     fyne.Window
	phaseLabel *canvas.Text
	timerLabel *canvas.Text
	totalLabel *canvas.Text
	progress           *canvas.Rectangle
	progressBackground *canvas.Rectangle
	background         *canvas.Rectangle

	playBtn      *widget.Button
	resetBtn     *widget.Button
	dashboardBtn *widget.Button

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

	ui.progressBackground = canvas.NewRectangle(color.NRGBA{R: 60, G: 60, B: 60, A: 255})
	ui.progressBackground.SetMinSize(fyne.NewSize(200, 2))

	ui.progress = canvas.NewRectangle(color.NRGBA{R: 50, G: 255, B: 50, A: 255})
	ui.progress.SetMinSize(fyne.NewSize(0, 2))

	ui.background = canvas.NewRectangle(color.NRGBA{R: 30, G: 30, B: 30, A: 200})

	// Small Control Buttons
	ui.playBtn = widget.NewButtonWithIcon("", theme.MediaPlayIcon(), ui.engine.Toggle)
	ui.playBtn.Importance = widget.LowImportance

	ui.resetBtn = widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), ui.engine.Reset)
	ui.resetBtn.Importance = widget.LowImportance

	ui.dashboardBtn = widget.NewButtonWithIcon("", theme.SettingsIcon(), ui.dashboard.Show)
	ui.dashboardBtn.Importance = widget.LowImportance

	controls := container.NewHBox(ui.playBtn, ui.resetBtn, ui.dashboardBtn)

	topBar := container.NewBorder(nil, nil, nil, controls, container.NewCenter(ui.phaseLabel))

	content := container.NewStack(
		ui.background,
		container.NewVBox(
			topBar,
			container.NewCenter(ui.timerLabel),
			ui.totalLabel,
			container.NewPadded(container.NewStack(ui.progressBackground, ui.progress)),
		),
	)

	ui.window.SetContent(content)
	ui.window.Resize(fyne.NewSize(240, 110))
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
	
	// Update Play/Pause Icon
	if state.IsRunning {
		ui.playBtn.SetIcon(theme.MediaPauseIcon())
	} else {
		ui.playBtn.SetIcon(theme.MediaPlayIcon())
	}

	// Progress bar reflects Total Session Progress
	ratio := 0.0
	if state.TotalDurationSeconds > 0 {
		ratio = float64(state.TotalDurationSeconds-state.TotalRemainingSeconds) / float64(state.TotalDurationSeconds)
	}
	
	// Total width of the container is approx 220
	fullWidth := ui.window.Content().Size().Width - 20
	ui.progress.SetMinSize(fyne.NewSize(float32(float64(fullWidth)*ratio), 2))
	
	if state.TotalRemainingSeconds < 300 && state.IsRunning { // Red in last 5 mins
		ui.timerLabel.Color = color.NRGBA{R: 255, G: 100, B: 0, A: 255}
		ui.progress.FillColor = color.NRGBA{R: 255, G: 100, B: 0, A: 255}
	} else if !state.IsRunning {
		ui.timerLabel.Color = color.NRGBA{R: 200, G: 200, B: 200, A: 255}
		ui.progress.FillColor = color.NRGBA{R: 150, G: 150, B: 150, A: 255}
	} else {
		ui.timerLabel.Color = color.NRGBA{R: 50, G: 255, B: 50, A: 255}
		ui.progress.FillColor = color.NRGBA{R: 50, G: 255, B: 50, A: 255}
	}

	ui.phaseLabel.Refresh()
	ui.timerLabel.Refresh()
	ui.totalLabel.Refresh()
	ui.progress.Refresh()
	ui.progressBackground.Refresh()
	ui.playBtn.Refresh()
}

func (ui *TimerUI) Show() {
	ui.window.Show()
	platform.TweakWindow("InterviewTimer")
}
