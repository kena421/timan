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
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	colorNeonGreen   = color.NRGBA{R: 34, G: 197, B: 94, A: 255}
	colorBrightGreen = color.NRGBA{R: 167, G: 243, B: 208, A: 255}
	colorGhostGreen  = color.NRGBA{R: 20, G: 60, B: 30, A: 255}
	colorNeonRed     = color.NRGBA{R: 239, G: 68, B: 68, A: 255}
	colorMutedGray   = color.NRGBA{R: 100, G: 116, B: 139, A: 255}
	colorSkyBlue     = color.NRGBA{R: 186, G: 230, B: 253, A: 255}
	colorBgDeep      = color.NRGBA{R: 9, G: 9, B: 11, A: 245}
	colorBgWell      = color.NRGBA{R: 39, G: 39, B: 42, A: 255}
	colorBorderLight = color.NRGBA{R: 63, G: 63, B: 70, A: 255}
)

// TimerUI handles the main floating window of the application.
type TimerUI struct {
	window             fyne.Window
	timerLabel         *canvas.Text
	totalLabel         *canvas.Text
	upcomingLabel      *canvas.Text
	background         *canvas.Rectangle
	border             *canvas.Rectangle
	eventNameLabel     *canvas.Text

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
	// 1. Scaled Main Timer
	ui.timerLabel = canvas.NewText("00:00", colorBrightGreen)
	ui.timerLabel.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
	ui.timerLabel.TextSize = 24
	ui.timerLabel.Alignment = fyne.TextAlignCenter

	// 2. Current Phase Info (Bottom Left)
	ui.totalLabel = canvas.NewText("PHASE: 00:00", colorSkyBlue)
	ui.totalLabel.TextSize = 9
	ui.totalLabel.TextStyle = fyne.TextStyle{Bold: true}
	ui.totalLabel.Alignment = fyne.TextAlignLeading

	// 3. Compact Upcoming Hint (Top Left)
	ui.upcomingLabel = canvas.NewText("NEXT: --", colorMutedGray)
	ui.upcomingLabel.TextSize = 7
	ui.upcomingLabel.TextStyle = fyne.TextStyle{Bold: true}
	ui.upcomingLabel.Alignment = fyne.TextAlignLeading

	// 3b. Event Name Context (Bottom Right)
	ui.eventNameLabel = canvas.NewText("SESSION", colorMutedGray)
	ui.eventNameLabel.TextSize = 7
	ui.eventNameLabel.TextStyle = fyne.TextStyle{Bold: true}
	ui.eventNameLabel.Alignment = fyne.TextAlignTrailing

	// 4. Backgrounds
	ui.background = canvas.NewRectangle(colorBgDeep)
	ui.border = canvas.NewRectangle(color.Transparent)
	ui.border.StrokeColor = color.Transparent
	ui.border.StrokeWidth = 0

	// 5. Accessible Controls (Top Right)
	iconSize := fyne.NewSize(14, 14)
	ui.playIcon = NewTappableIcon(theme.MediaPlayIcon(), iconSize, func() { ui.engine.Toggle() })
	reset := NewTappableIcon(theme.ViewRefreshIcon(), iconSize, func() { ui.engine.Reset() })
	dash := NewTappableIcon(theme.SettingsIcon(), iconSize, func() { ui.dashboard.Show() })

	controlBg := canvas.NewRectangle(colorBgWell)
	// Only bottom-left corner rounded to "flush" against top and right edges
	controlBg.CornerRadius = 4
	
	iconContainer := container.NewHBox(ui.playIcon, reset, dash)
	controlWell := container.NewStack(controlBg, container.NewPadded(iconContainer))

	// ASSEMBLY (Zero-Padding on Controls)
	topRow := container.NewHBox(
		container.NewPadded(ui.upcomingLabel),
		layout.NewSpacer(),
		controlWell, // Flushed to the corner
	)

	bottomRow := container.NewHBox(
		container.NewPadded(ui.totalLabel),
		layout.NewSpacer(),
		container.NewPadded(ui.eventNameLabel),
	)

	mainLayout := container.NewBorder(
		topRow,
		bottomRow,
		nil, nil,
		container.NewCenter(ui.timerLabel),
	)

	content := container.NewStack(
		ui.background,
		ui.border,
		mainLayout, // Removed outer padding to allow flushing
	)

	ui.window.SetContent(content)
	ui.window.Resize(fyne.NewSize(180, 85))
}

func (ui *TimerUI) formatTime(s int) string {
	return fmt.Sprintf("%02d:%02d", s/60, s%60)
}

// OnTick updates the UI with the latest timer state.
func (ui *TimerUI) OnTick(state engine.TimerState) {
	ui.timerLabel.Text = ui.formatTime(state.TotalRemainingSeconds)
	ui.totalLabel.Text = fmt.Sprintf("%s: %s", 
		strings.ToUpper(state.CurrentPhase.Name),
		ui.formatTime(state.RemainingSeconds))
	
	if state.UpcomingPhaseName != "" {
		ui.upcomingLabel.Text = "NEXT: " + strings.ToUpper(state.UpcomingPhaseName)
	} else {
		ui.upcomingLabel.Text = ""
	}

	ui.eventNameLabel.Text = strings.ToUpper(state.EventName)

	if state.IsRunning {
		ui.playIcon.SetResource(theme.MediaPauseIcon())
	} else {
		ui.playIcon.SetResource(theme.MediaPlayIcon())
	}

	if state.TotalRemainingSeconds <= state.WarningSeconds && state.IsRunning {
		ui.timerLabel.Color = colorNeonRed
		ui.background.FillColor = color.NRGBA{R: 45, G: 10, B: 10, A: 245}
	} else if !state.IsRunning {
		ui.timerLabel.Color = colorGhostGreen
		ui.background.FillColor = colorBgDeep
	} else {
		ui.timerLabel.Color = colorBrightGreen
		ui.background.FillColor = colorBgDeep
	}

	ui.timerLabel.Refresh()
	ui.totalLabel.Refresh()
	ui.upcomingLabel.Refresh()
	ui.eventNameLabel.Refresh()
	ui.background.Refresh()
	ui.border.Refresh()
	ui.playIcon.Refresh()
}

// Show displays the main timer window and applies platform tweaks.
func (ui *TimerUI) Show() {
	ui.window.Show()
	platform.TweakWindow("Timan")
}
