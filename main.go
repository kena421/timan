package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

void makeWindowTopmostAndFrameless(const char* title) {
    @autoreleasepool {
        NSString* nsTitle = [NSString stringWithUTF8String:title];
        NSArray* windows = [NSApp windows];
        for (NSWindow* window in windows) {
            if ([[window title] isEqualToString:nsTitle]) {
                [window setStyleMask:NSWindowStyleMaskBorderless];
                [window setLevel:NSStatusWindowLevel]; // Always on top
                [window setBackgroundColor:[NSColor clearColor]];
                [window setOpaque:NO];
                [window setHasShadow:YES];
                [window setMovableByWindowBackground:YES];
                // Make it visible on all spaces
                [window setCollectionBehavior:NSWindowCollectionBehaviorCanJoinAllSpaces | NSWindowCollectionBehaviorFullScreenAuxiliary];
                break;
            }
        }
    }
}
*/
import "C"

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	windowTitle    = "InterviewTimer"
	defaultLimit   = 45 * 60 // 45 minutes
	refreshRate    = time.Second
)

type TimerMode int

const (
	ModeUp TimerMode = iota
	ModeDown
)

type timerApp struct {
	window      fyne.Window
	label       *canvas.Text
	background  *canvas.Rectangle
	
	seconds     int
	limit       int
	isRunning   bool
	mode        TimerMode
	
	ticker      *time.Ticker
	stopChan    chan bool
}

func newTimerApp(w fyne.Window) *timerApp {
	t := &timerApp{
		window:    w,
		seconds:   0,
		limit:     defaultLimit,
		isRunning: true,
		mode:      ModeUp,
		stopChan:  make(chan bool),
	}
	
	t.setupUI()
	return t
}

func (t *timerApp) setupUI() {
	// Monospaced high-contrast font
	t.label = canvas.NewText("00:00", color.NRGBA{R: 0, G: 255, B: 0, A: 255})
	t.label.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
	t.label.TextSize = 42
	t.label.Alignment = fyne.TextAlignCenter

	// Semi-transparent dark background (70% opacity)
	t.background = canvas.NewRectangle(color.NRGBA{R: 30, G: 30, B: 30, A: 180})
	// Rounded corners hack: Fyne doesn't have a CornerRadius on Rectangle, 
	// but we can use a Card or just wait for the Mac native window to handle it if we used a path.
	// For simplicity and "lightweight" feel, we'll use a stack.

	// Context Menu
	menu := fyne.NewMenu("",
		fyne.NewMenuItem("Toggle Mode (Up/Down)", t.toggleMode),
		fyne.NewMenuItem("Reset", t.reset),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Quit", func() { fyne.CurrentApp().Quit() }),
	)

	content := container.NewStack(
		t.background,
		container.NewCenter(t.label),
		&contextMenuWrapper{
			menu:   menu,
			window: t.window,
		},
	)

	t.window.SetContent(content)
	t.window.Resize(fyne.NewSize(220, 80))
	t.window.SetFixedSize(true)
	t.window.SetPadded(false)
}

func (t *timerApp) start() {
	t.ticker = time.NewTicker(refreshRate)
	go func() {
		for {
			select {
			case <-t.ticker.C:
				if t.isRunning {
					if t.mode == ModeUp {
						t.seconds++
					} else {
						if t.seconds > 0 {
							t.seconds--
						}
					}
					t.updateUI()
				}
			case <-t.stopChan:
				return
			}
		}
	}()
}

func (t *timerApp) updateUI() {
	mins := t.seconds / 60
	secs := t.seconds % 60
	t.label.Text = fmt.Sprintf("%02d:%02d", mins, secs)
	
	// Visual Cues: Red if over limit or count down reached zero
	if (t.mode == ModeUp && t.seconds >= t.limit) || (t.mode == ModeDown && t.seconds == 0) {
		t.label.Color = color.NRGBA{R: 255, G: 50, B: 50, A: 255}
	} else {
		t.label.Color = color.NRGBA{R: 50, G: 255, B: 50, A: 255}
	}
	t.label.Refresh()
}

func (t *timerApp) toggleMode() {
	if t.mode == ModeUp {
		t.mode = ModeDown
		t.seconds = t.limit
	} else {
		t.mode = ModeUp
		t.seconds = 0
	}
	t.updateUI()
}

func (t *timerApp) reset() {
	if t.mode == ModeUp {
		t.seconds = 0
	} else {
		t.seconds = t.limit
	}
	t.updateUI()
}

// --- Custom Widgets ---

type contextMenuWrapper struct {
	widget.BaseWidget
	menu   *fyne.Menu
	window fyne.Window
}

func (c *contextMenuWrapper) TappedSecondary(e *fyne.PointEvent) {
	widget.ShowPopUpMenuAtPosition(c.menu, c.window.Canvas(), e.AbsolutePosition)
}

func (c *contextMenuWrapper) CreateRenderer() fyne.WidgetRenderer {
	// Invisible layer to capture right clicks
	return widget.NewSimpleRenderer(canvas.NewRectangle(color.Transparent))
}

// --- Main ---

func main() {
	a := app.NewWithID("com.interview.timer")
	a.SetIcon(theme.SettingsIcon())
	
	w := a.NewWindow(windowTitle)
	
	timer := newTimerApp(w)
	timer.start()

	// Show window first so it's created in the OS
	w.Show()

	// Apply native tweaks for Mac (Always on Top + Frameless)
	// We run this in a goroutine to ensure the window has time to initialize
	go func() {
		time.Sleep(100 * time.Millisecond)
		C.makeWindowTopmostAndFrameless(C.CString(windowTitle))
	}()

	a.Run()
}
