package ui

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/kena421/timan/internal/engine"
	"github.com/kena421/timan/internal/platform"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// UI Color Palette (Industrial Neon)
var (
	colorNeonGreen   = color.NRGBA{R: 34, G: 197, B: 94, A: 255}
	colorBrightGreen = color.NRGBA{R: 167, G: 243, B: 208, A: 255}
	colorGhostGreen  = color.NRGBA{R: 20, G: 60, B: 30, A: 255}
	colorNeonRed     = color.NRGBA{R: 239, G: 68, B: 68, A: 255}
	colorMutedGray   = color.NRGBA{R: 100, G: 116, B: 139, A: 255}
	colorSkyBlue     = color.NRGBA{R: 186, G: 230, B: 253, A: 255}
	colorBgDeep      = color.NRGBA{R: 9, G: 9, B: 11, A: 245}
	colorBgWell      = color.NRGBA{R: 39, G: 39, B: 42, A: 255}
)

// TimerUI manages the primary HUD window and its interactive elements.
// It implements the engine.TimerObserver interface to react to ticks.
type TimerUI struct {
	window             fyne.Window
	timerLabel         *canvas.Text      // Large center countdown (total time)
	totalLabel         *canvas.Text      // Sub-countdown (phase name + phase time)
	upcomingLabel      *canvas.Text      // Small top-left preview of next phase
	background         *canvas.Rectangle // The main HUD background
	border             *canvas.Rectangle // Optional window border (currently transparent)
	eventNameLabel     *canvas.Text      // Sub-info top-center showing active profile name

	playIcon *TappableIcon // Toggle button for start/pause

	engine    *engine.TimerEngine
	dashboard *Dashboard
	privacyOn bool // Tracks if Screen Sharing Privacy is active
}

// TappableIcon is a lightweight custom widget that makes an icon respond to mouse clicks.
type TappableIcon struct {
	widget.Icon
	OnTap   func()
	minSize fyne.Size
}

// Tapped captures the click event from the Fyne driver.
func (t *TappableIcon) Tapped(_ *fyne.PointEvent) {
	if t.OnTap != nil {
		t.OnTap()
	}
}

// MinSize allows us to control the "hit area" and visual size of the icon.
func (t *TappableIcon) MinSize() fyne.Size {
	return t.minSize
}

// NewTappableIcon factory for interactive UI elements.
func NewTappableIcon(res fyne.Resource, size fyne.Size, onTap func()) *TappableIcon {
	t := &TappableIcon{OnTap: onTap, minSize: size}
	t.SetResource(res)
	t.ExtendBaseWidget(t)
	return t
}

// NewTimerUI constructs the HUD and links it to the provided engine.
func NewTimerUI(w fyne.Window, e *engine.TimerEngine) *TimerUI {
	ui := &TimerUI{
		window: w,
		engine: e,
	}
	ui.dashboard = NewDashboard(ui, e)
	ui.setup()
	e.AddObserver(ui)
	ui.OnTick(e.GetState()) // Set initial state
	return ui
}

// setup builds the visual tree of the HUD using a balanced 4-corner layout.
func (ui *TimerUI) setup() {
	// 1. Scaled Main Timer (Total Remaining)
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

	// 4. Background Layers
	ui.background = canvas.NewRectangle(colorBgDeep)
	ui.border = canvas.NewRectangle(color.Transparent)
	ui.border.StrokeWidth = 0

	// 5. Flushed Controls (Top Right)
	iconSize := fyne.NewSize(14, 14)
	ui.playIcon = NewTappableIcon(theme.MediaPlayIcon(), iconSize, func() { ui.engine.Toggle() })
	reset := NewTappableIcon(theme.ViewRefreshIcon(), iconSize, func() { ui.engine.Reset() })
	dash := NewTappableIcon(theme.SettingsIcon(), iconSize, func() { ui.dashboard.Show() })
	
	// Privacy (Screen Sharing Hide) Button
	ui.privacyOn = false
	privacyBtn := NewTappableIcon(theme.VisibilityIcon(), iconSize, nil)
	privacyBtn.OnTap = func() {
		ui.privacyOn = !ui.privacyOn
		platform.SetPrivacyMode("Timan", ui.privacyOn)
		if ui.privacyOn {
			privacyBtn.SetResource(theme.VisibilityOffIcon())
		} else {
			privacyBtn.SetResource(theme.VisibilityIcon())
		}
		privacyBtn.Refresh()
	}

	controlBg := canvas.NewRectangle(colorBgWell)
	controlBg.CornerRadius = 4
	
	iconContainer := container.NewHBox(ui.playIcon, reset, privacyBtn, dash)
	controlWell := container.NewStack(controlBg, container.NewPadded(iconContainer))

	// ASSEMBLY (Four-Corner Balanced Distribution)
	topRow := container.NewHBox(
		container.NewPadded(ui.upcomingLabel),
		layout.NewSpacer(),
		controlWell, // No padding to allow corner flushing
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
		mainLayout,
	)

	ui.window.SetContent(content)
	ui.window.Resize(fyne.NewSize(180, 85))
}

func (ui *TimerUI) formatTime(s int) string {
	return fmt.Sprintf("%02d:%02d", s/60, s%60)
}

// OnTick is the reactive callback that updates the HUD every second.
func (ui *TimerUI) OnTick(state engine.TimerState) {
	// Update text content
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

	// Update play/pause icon state
	if state.IsRunning {
		ui.playIcon.SetResource(theme.MediaPauseIcon())
	} else {
		ui.playIcon.SetResource(theme.MediaPlayIcon())
	}

	// Update Visual Modes (Alert, Paused, Active)
	if state.TotalRemainingSeconds <= state.WarningSeconds && state.IsRunning {
		// Alert Mode: Neon Red Digit + Slight Red Background Glow
		ui.timerLabel.Color = colorNeonRed
		ui.background.FillColor = color.NRGBA{R: 45, G: 10, B: 10, A: 245}
	} else if !state.IsRunning {
		// Paused Mode: Dimmer Ghostly Green
		ui.timerLabel.Color = colorGhostGreen
		ui.background.FillColor = colorBgDeep
	} else {
		// Active Mode: Vibrant Bright Green
		ui.timerLabel.Color = colorBrightGreen
		ui.background.FillColor = colorBgDeep
	}

	// Force canvas refreshes
	ui.timerLabel.Refresh()
	ui.totalLabel.Refresh()
	ui.upcomingLabel.Refresh()
	ui.eventNameLabel.Refresh()
	ui.background.Refresh()
	ui.border.Refresh()
	ui.playIcon.Refresh()
}

// Show renders the window and applies platform-specific HUD tweaks (Always on Top, etc.)
func (ui *TimerUI) Show() {
	ui.window.Show()
	platform.TweakWindow("Timan")
}
