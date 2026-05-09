package ui

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/timan-org/timan/internal/engine"
	"github.com/timan-org/timan/internal/platform"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// TimerUI handles the main floating window of the application.
type TimerUI struct {
	window             fyne.Window
	timerLabel         *canvas.Text
	totalLabel         *canvas.Text
	upcomingLabel      *canvas.Text
	background         *canvas.Rectangle

	playIcon *TappableIcon

	engine    *engine.TimerEngine
	dashboard *Dashboard
}

// TappableIcon is a widget that wraps an icon and provides a tap handler.
type TappableIcon struct {
	widget.Icon
	OnTap   func()
	minSize fyne.Size
}

// Tapped implements the fyne.Tappable interface.
func (t *TappableIcon) Tapped(_ *fyne.PointEvent) {
	if t.OnTap != nil {
		t.OnTap()
	}
}

// MinSize overrides the default minimum size of the icon.
func (t *TappableIcon) MinSize() fyne.Size {
	return t.minSize
}

// NewTappableIcon creates a new TappableIcon instance.
func NewTappableIcon(res fyne.Resource, size fyne.Size, onTap func()) *TappableIcon {
	t := &TappableIcon{OnTap: onTap, minSize: size}
	t.SetResource(res)
	t.ExtendBaseWidget(t)
	return t
}

// NewTimerUI initializes the TimerUI and connects it to the TimerEngine.
func NewTimerUI(w fyne.Window, e *engine.TimerEngine) *TimerUI {
	ui := &TimerUI{
		window: w,
		engine: e,
	}
	ui.dashboard = NewDashboard(ui, e)
	ui.setup()
	e.AddObserver(ui)
	ui.OnTick(e.GetState())
	return ui
}

func (ui *TimerUI) setup() {
	ui.timerLabel = canvas.NewText("00:00", color.NRGBA{R: 200, G: 200, B: 200, A: 255})
	ui.timerLabel.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
	ui.timerLabel.TextSize = 28
	ui.timerLabel.Alignment = fyne.TextAlignCenter

	ui.totalLabel = canvas.NewText("PHASE: 00:00", color.NRGBA{R: 180, G: 220, B: 255, A: 255})
	ui.totalLabel.TextSize = 10
	ui.totalLabel.TextStyle = fyne.TextStyle{Bold: true}
	ui.totalLabel.Alignment = fyne.TextAlignCenter

	ui.upcomingLabel = canvas.NewText("", color.NRGBA{R: 120, G: 120, B: 120, A: 255})
	ui.upcomingLabel.TextSize = 7
	ui.upcomingLabel.Alignment = fyne.TextAlignCenter

	ui.background = canvas.NewRectangle(color.NRGBA{R: 30, G: 30, B: 30, A: 200})

	// Micro Icons (with enlarged hit area)
	iconSize := fyne.NewSize(16, 16)
	
	ui.playIcon = NewTappableIcon(theme.MediaPlayIcon(), iconSize, func() { ui.engine.Toggle() })
	reset := NewTappableIcon(theme.ViewRefreshIcon(), iconSize, func() { ui.engine.Reset() })
	dash := NewTappableIcon(theme.SettingsIcon(), iconSize, func() { ui.dashboard.Show() })

	playWrap := container.NewStack(canvas.NewRectangle(color.Transparent), ui.playIcon)
	resetWrap := container.NewStack(canvas.NewRectangle(color.Transparent), reset)
	dashWrap := container.NewStack(canvas.NewRectangle(color.Transparent), dash)

	controls := container.NewHBox(playWrap, resetWrap, dashWrap)
	topBar := container.NewBorder(nil, nil, nil, controls, container.NewCenter(ui.upcomingLabel))

	content := container.NewStack(
		ui.background,
		container.NewVBox(
			topBar,
			container.NewCenter(ui.timerLabel),
			container.NewCenter(ui.totalLabel),
		),
	)

	ui.window.SetContent(content)
	ui.window.Resize(fyne.NewSize(180, 80))
}

func (ui *TimerUI) formatTime(s int) string {
	return fmt.Sprintf("%02d:%02d", s/60, s%60)
}

// OnTick updates the UI with the latest timer state.
func (ui *TimerUI) OnTick(state engine.TimerState) {
	// Main Focus: ONLY Remaining Session Time
	ui.timerLabel.Text = ui.formatTime(state.TotalRemainingSeconds)
	
	// Phase Focus (Emphasized)
	ui.totalLabel.Text = fmt.Sprintf("%s: %s", 
		strings.ToUpper(state.CurrentPhase.Name),
		ui.formatTime(state.RemainingSeconds))
	ui.totalLabel.Color = color.NRGBA{R: 180, G: 220, B: 255, A: 255} // Light Cyan/Blue for distinction
	
	// Upcoming Phase (Subtle)
	if state.UpcomingPhaseName != "" {
		ui.upcomingLabel.Text = "NEXT: " + strings.ToUpper(state.UpcomingPhaseName)
	} else {
		ui.upcomingLabel.Text = ""
	}

	// Update Play/Pause Icon
	if state.IsRunning {
		ui.playIcon.SetResource(theme.MediaPauseIcon())
	} else {
		ui.playIcon.SetResource(theme.MediaPlayIcon())
	}

	// Background Alert & Timer Color
	if state.TotalRemainingSeconds <= state.WarningSeconds && state.IsRunning {
		ui.timerLabel.Color = color.NRGBA{R: 255, G: 100, B: 0, A: 255}
		// Flash/Change background to alert
		ui.background.FillColor = color.NRGBA{R: 100, G: 0, B: 0, A: 180}
	} else if !state.IsRunning {
		ui.timerLabel.Color = color.NRGBA{R: 200, G: 200, B: 200, A: 255}
		ui.background.FillColor = color.NRGBA{R: 30, G: 30, B: 30, A: 200}
	} else {
		ui.timerLabel.Color = color.NRGBA{R: 50, G: 255, B: 50, A: 255}
		ui.background.FillColor = color.NRGBA{R: 30, G: 30, B: 30, A: 200}
	}

	ui.timerLabel.Refresh()
	ui.totalLabel.Refresh()
	ui.upcomingLabel.Refresh()
	ui.background.Refresh()
	ui.playIcon.Refresh()
}

// Show displays the main timer window and applies platform tweaks.
func (ui *TimerUI) Show() {
	ui.window.Show()
	platform.TweakWindow("Timan") // Updated window title for tweaking
}
