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

	playIcon *TappableIcon

	engine    *engine.TimerEngine
	dashboard *Dashboard
}

type TappableIcon struct {
	widget.Icon
	OnTap   func()
	minSize fyne.Size
}

func (t *TappableIcon) Tapped(_ *fyne.PointEvent) {
	if t.OnTap != nil {
		t.OnTap()
	}
}

func (t *TappableIcon) MinSize() fyne.Size {
	return t.minSize
}

func NewTappableIcon(res fyne.Resource, size fyne.Size, onTap func()) *TappableIcon {
	t := &TappableIcon{OnTap: onTap, minSize: size}
	t.SetResource(res)
	t.ExtendBaseWidget(t)
	return t
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
	ui.phaseLabel.TextSize = 8
	ui.phaseLabel.Alignment = fyne.TextAlignCenter

	ui.timerLabel = canvas.NewText("00:00", color.NRGBA{R: 200, G: 200, B: 200, A: 255})
	ui.timerLabel.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
	ui.timerLabel.TextSize = 28
	ui.timerLabel.Alignment = fyne.TextAlignCenter

	ui.totalLabel = canvas.NewText("00:00/00:00", color.NRGBA{R: 150, G: 150, B: 150, A: 255})
	ui.totalLabel.TextSize = 8
	ui.totalLabel.Alignment = fyne.TextAlignCenter

	ui.progressBackground = canvas.NewRectangle(color.NRGBA{R: 60, G: 60, B: 60, A: 255})
	ui.progressBackground.SetMinSize(fyne.NewSize(160, 2))

	ui.progress = canvas.NewRectangle(color.NRGBA{R: 50, G: 255, B: 50, A: 255})
	ui.progress.SetMinSize(fyne.NewSize(0, 2))

	ui.background = canvas.NewRectangle(color.NRGBA{R: 30, G: 30, B: 30, A: 200})

	// Micro Icons
	size := fyne.NewSize(12, 12)
	ui.playIcon = NewTappableIcon(theme.MediaPlayIcon(), size, ui.engine.Toggle)
	reset := NewTappableIcon(theme.ViewRefreshIcon(), size, ui.engine.Reset)
	dash := NewTappableIcon(theme.SettingsIcon(), size, ui.dashboard.Show)

	controls := container.NewHBox(ui.playIcon, reset, dash)
	topBar := container.NewBorder(nil, nil, nil, controls, container.NewCenter(ui.phaseLabel))

	content := container.NewStack(
		ui.background,
		container.NewVBox(
			topBar,
			container.NewCenter(ui.timerLabel),
			container.NewCenter(ui.totalLabel),
			container.NewStack(ui.progressBackground, ui.progress),
		),
	)

	ui.window.SetContent(content)
	ui.window.Resize(fyne.NewSize(180, 75))
}

func (ui *TimerUI) formatTime(s int) string {
	return fmt.Sprintf("%02d:%02d", s/60, s%60)
}

func (ui *TimerUI) OnTick(state engine.TimerState) {
	ui.phaseLabel.Text = strings.ToUpper(state.CurrentPhase.Name)
	
	// Main Focus: ONLY Remaining Session Time
	ui.timerLabel.Text = ui.formatTime(state.TotalRemainingSeconds)
	
	// Ultra-compact: P for Phase, G for Goal
	ui.totalLabel.Text = fmt.Sprintf("P: %s/%s | G: %s", 
		ui.formatTime(state.RemainingSeconds), 
		ui.formatTime(state.CurrentPhase.Duration),
		ui.formatTime(state.TotalDurationSeconds))
	
	// Update Play/Pause Icon
	if state.IsRunning {
		ui.playIcon.SetResource(theme.MediaPauseIcon())
	} else {
		ui.playIcon.SetResource(theme.MediaPlayIcon())
	}

	// Progress bar reflects Total Session Progress
	ratio := 0.0
	if state.TotalDurationSeconds > 0 {
		ratio = float64(state.TotalDurationSeconds-state.TotalRemainingSeconds) / float64(state.TotalDurationSeconds)
	}
	
	fullWidth := float64(160) // Compact width
	ui.progress.SetMinSize(fyne.NewSize(float32(fullWidth*ratio), 2))
	
	// Background Alert & Timer Color
	if state.TotalRemainingSeconds <= state.WarningSeconds && state.IsRunning {
		ui.timerLabel.Color = color.NRGBA{R: 255, G: 100, B: 0, A: 255}
		ui.progress.FillColor = color.NRGBA{R: 255, G: 50, B: 50, A: 255}
		// Flash/Change background to alert
		ui.background.FillColor = color.NRGBA{R: 100, G: 0, B: 0, A: 180}
	} else if !state.IsRunning {
		ui.timerLabel.Color = color.NRGBA{R: 200, G: 200, B: 200, A: 255}
		ui.progress.FillColor = color.NRGBA{R: 150, G: 150, B: 150, A: 255}
		ui.background.FillColor = color.NRGBA{R: 30, G: 30, B: 30, A: 200}
	} else {
		ui.timerLabel.Color = color.NRGBA{R: 50, G: 255, B: 50, A: 255}
		ui.progress.FillColor = color.NRGBA{R: 50, G: 255, B: 50, A: 255}
		ui.background.FillColor = color.NRGBA{R: 30, G: 30, B: 30, A: 200}
	}

	ui.phaseLabel.Refresh()
	ui.timerLabel.Refresh()
	ui.totalLabel.Refresh()
	ui.progress.Refresh()
	ui.progressBackground.Refresh()
	ui.background.Refresh()
	ui.playIcon.Refresh()
}

func (ui *TimerUI) Show() {
	ui.window.Show()
	platform.TweakWindow("InterviewTimer")
}
