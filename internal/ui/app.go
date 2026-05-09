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
	ui.timerLabel.TextSize = 24 // Reduced from 28
	ui.timerLabel.Alignment = fyne.TextAlignCenter

	// 2. Phase Info
	ui.totalLabel = canvas.NewText("PHASE: 00:00", colorSkyBlue)
	ui.totalLabel.TextSize = 10
	ui.totalLabel.TextStyle = fyne.TextStyle{Bold: true}
	ui.totalLabel.Alignment = fyne.TextAlignCenter

	// 3. Upcoming Hint
	ui.upcomingLabel = canvas.NewText("NEXT: --", colorMutedGray)
	ui.upcomingLabel.TextSize = 7
	ui.upcomingLabel.TextStyle = fyne.TextStyle{Bold: true}
	ui.upcomingLabel.Alignment = fyne.TextAlignLeading

	// 4. Backgrounds
	ui.background = canvas.NewRectangle(colorBgDeep)
	ui.border = canvas.NewRectangle(color.Transparent)
	ui.border.StrokeColor = color.Transparent // Removed border color
	ui.border.StrokeWidth = 0                 // Set to 0

	// 5. Accessible Controls (Larger)
	iconSize := fyne.NewSize(14, 14) // Increased from 11
	ui.playIcon = NewTappableIcon(theme.MediaPlayIcon(), iconSize, func() { ui.engine.Toggle() })
	reset := NewTappableIcon(theme.ViewRefreshIcon(), iconSize, func() { ui.engine.Reset() })
	dash := NewTappableIcon(theme.SettingsIcon(), iconSize, func() { ui.dashboard.Show() })

	controlBg := canvas.NewRectangle(colorBgWell)
	controlBg.CornerRadius = 4
	
	// Spaced container for icons
	iconContainer := container.NewHBox(
		ui.playIcon,
		reset,
		dash,
	)
	
	controlWell := container.NewStack(
		controlBg,
		container.NewPadded(iconContainer), // Re-added padding for better "well" feel
	)

	// ASSEMBLY
	topRow := container.NewBorder(nil, nil, 
		ui.upcomingLabel, 
		controlWell, 
		nil)

	mainLayout := container.NewBorder(
		topRow,
		ui.totalLabel,
		nil, nil,
		ui.timerLabel,
	)

	content := container.NewStack(
		ui.background,
		ui.border,
		container.NewPadded(mainLayout),
	)

	ui.window.SetContent(content)
	ui.window.Resize(fyne.NewSize(160, 80))
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

	if state.IsRunning {
		ui.playIcon.SetResource(theme.MediaPauseIcon())
	} else {
		ui.playIcon.SetResource(theme.MediaPlayIcon())
	}

	if state.TotalRemainingSeconds <= state.WarningSeconds && state.IsRunning {
		// Alert Mode: Neon Red Flash
		ui.timerLabel.Color = colorNeonRed
		ui.background.FillColor = color.NRGBA{R: 45, G: 10, B: 10, A: 245}
	} else if !state.IsRunning {
		// Paused Mode: Ghostly Green (Visible but dim)
		ui.timerLabel.Color = colorGhostGreen
		ui.background.FillColor = colorBgDeep
	} else {
		// Active Mode: Vibrant Bright Green
		ui.timerLabel.Color = colorBrightGreen
		ui.background.FillColor = colorBgDeep
	}

	ui.timerLabel.Refresh()
	ui.totalLabel.Refresh()
	ui.upcomingLabel.Refresh()
	ui.background.Refresh()
	ui.border.Refresh()
	ui.playIcon.Refresh()
}

// Show displays the main timer window and applies platform tweaks.
func (ui *TimerUI) Show() {
	ui.window.Show()
	platform.TweakWindow("Timan")
}
