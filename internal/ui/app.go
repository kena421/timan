package ui

import (
	"fmt"
	"image/color"
	"strings"

	"timan/internal/domain"
	"timan/internal/engine"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

type TimerUI struct {
	window     fyne.Window
	phaseLabel *canvas.Text
	timerLabel *canvas.Text
	progress   *widget.ProgressBar
	background *canvas.Rectangle

	engine *engine.TimerEngine
	manager *BlueprintManager
}

func NewTimerUI(w fyne.Window, e *engine.TimerEngine) *TimerUI {
	ui := &TimerUI{
		window: w,
		engine: e,
	}
	ui.manager = NewBlueprintManager(ui, e)
	ui.setup()
	e.AddObserver(ui)
	return ui
}

func (ui *TimerUI) setup() {
	ui.phaseLabel = canvas.NewText("INITIALIZING", color.NRGBA{R: 200, G: 200, B: 200, A: 255})
	ui.phaseLabel.TextSize = 12
	ui.phaseLabel.Alignment = fyne.TextAlignCenter

	ui.timerLabel = canvas.NewText("00:00", color.NRGBA{R: 200, G: 200, B: 200, A: 255})
	ui.timerLabel.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
	ui.timerLabel.TextSize = 36
	ui.timerLabel.Alignment = fyne.TextAlignCenter

	ui.progress = widget.NewProgressBar()
	ui.progress.TextFormatter = func() string { return "" }

	ui.background = canvas.NewRectangle(color.NRGBA{R: 30, G: 30, B: 30, A: 200})

	menu := fyne.NewMenu("",
		fyne.NewMenuItem("Design Quick Blueprint", ui.showConfig),
		fyne.NewMenuItem("Blueprint Library", ui.manager.Show),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Reset current phase", ui.engine.ResetPhase),
		fyne.NewMenuItem("Reset All", ui.engine.Reset),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Quit", func() { fyne.CurrentApp().Quit() }),
	)

	content := container.NewStack(
		ui.background,
		container.NewVBox(
			container.NewPadded(ui.phaseLabel),
			container.NewCenter(ui.timerLabel),
			ui.progress,
		),
		&InteractionWrapper{
			OnTap:  ui.engine.Toggle,
			Menu:   menu,
			Window: ui.window,
		},
	)

	ui.window.SetContent(content)
}

func (ui *TimerUI) OnTick(state engine.TimerState) {
	ui.phaseLabel.Text = strings.ToUpper(state.CurrentPhase.Name)
	ui.timerLabel.Text = state.CurrentPhase.FormatDuration() // This is wrong, should be RemainingSeconds
    ui.timerLabel.Text = fmt.Sprintf("%02d:%02d", state.RemainingSeconds/60, state.RemainingSeconds%60)
	ui.progress.Max = float64(state.CurrentPhase.Duration)
	ui.progress.Value = float64(state.RemainingSeconds)

	if state.RemainingSeconds < 60 && state.IsRunning {
		ui.timerLabel.Color = color.NRGBA{R: 255, G: 100, B: 0, A: 255}
	} else if !state.IsRunning {
		ui.timerLabel.Color = color.NRGBA{R: 200, G: 200, B: 200, A: 255}
	} else {
		ui.timerLabel.Color = color.NRGBA{R: 50, G: 255, B: 50, A: 255}
	}

	ui.phaseLabel.Refresh()
	ui.timerLabel.Refresh()
	ui.progress.Refresh()
}

func (ui *TimerUI) showConfig() {
	ui.showBlueprintEditor(nil, func(b *domain.Blueprint) {
		ui.engine.UpdatePhases(b.Phases)
	})
}

func (ui *TimerUI) showBlueprintEditor(existing *domain.Blueprint, onSave func(*domain.Blueprint)) {
	title := "Blueprint Designer"
	if existing != nil {
		title = "Edit: " + existing.Name
	}
	configWindow := fyne.CurrentApp().NewWindow(title)
	
	phases := ui.engine.GetPhases()
	if existing != nil {
		phases = existing.Phases
	}

	totalMins := 0
	for _, p := range phases {
		totalMins += p.Duration / 60
	}

	nameEntry := widget.NewEntry()
	if existing != nil {
		nameEntry.SetText(existing.Name)
	} else {
		nameEntry.SetText("New Blueprint")
	}

	totalEntry := widget.NewEntry()
	totalEntry.SetText(strconv.Itoa(totalMins))
	totalEntry.PlaceHolder = "Total Duration (mins)"

	rows := container.NewVBox()
	summaryLabel := widget.NewLabel("")

	updateSummary := func() {
		currentSum := 0
		for _, row := range rows.Objects {
			if box, ok := row.(*fyne.Container); ok {
				if grid, ok := box.Objects[0].(*fyne.Container); ok {
					if minsEntry, ok := grid.Objects[1].(*widget.Entry); ok {
						m, _ := strconv.Atoi(minsEntry.Text)
						currentSum += m
					}
				}
			}
		}
		target, _ := strconv.Atoi(totalEntry.Text)
		summaryLabel.SetText(fmt.Sprintf("Allocated: %d / %d mins", currentSum, target))
		summaryLabel.Importance = widget.DangerImportance
		if currentSum == target && target > 0 {
			summaryLabel.Importance = widget.SuccessImportance
		}
		summaryLabel.Refresh()
	}

	addPhaseRow := func(name string, mins int) {
		pNameEntry := widget.NewEntry()
		pNameEntry.SetText(name)
		pNameEntry.PlaceHolder = "Phase Name"
		pNameEntry.OnChanged = func(string) { updateSummary() }

		minsEntry := widget.NewEntry()
		minsEntry.SetText(strconv.Itoa(mins))
		minsEntry.PlaceHolder = "Mins"
		minsEntry.OnChanged = func(string) { updateSummary() }

		removeBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)
		grid := container.NewGridWithColumns(2, pNameEntry, minsEntry)
		row := container.NewBorder(nil, nil, nil, removeBtn, grid)
		
		removeBtn.OnTapped = func() {
			rows.Remove(row)
			updateSummary()
		}

		rows.Add(row)
		updateSummary()
	}

	for _, p := range phases {
		addPhaseRow(p.Name, p.Duration/60)
	}

	totalEntry.OnChanged = func(string) { updateSummary() }

	scroll := container.NewVScroll(rows)
	scroll.SetMinSize(fyne.NewSize(400, 300))

	saveBtn := widget.NewButtonWithIcon("Save Blueprint", theme.ConfirmIcon(), func() {
		newPhases := []domain.Phase{}
		target, _ := strconv.Atoi(totalEntry.Text)
		
		for _, row := range rows.Objects {
			if box, ok := row.(*fyne.Container); ok {
				if grid, ok := box.Objects[0].(*fyne.Container); ok {
					pNameEntry := grid.Objects[0].(*widget.Entry)
					minsEntry := grid.Objects[1].(*widget.Entry)
					
					m, _ := strconv.Atoi(minsEntry.Text)
					if pNameEntry.Text != "" && m > 0 {
						newPhases = append(newPhases, domain.Phase{Name: pNameEntry.Text, Duration: m * 60})
					}
				}
			}
		}

		if err := domain.ValidatePhases(newPhases, target); err == nil {
			onSave(&domain.Blueprint{
				Name:   nameEntry.Text,
				Total:  target,
				Phases: newPhases,
			})
			configWindow.Close()
		} else {
			dialog.ShowError(err, configWindow)
		}
	})

	addBtn := widget.NewButtonWithIcon("Add Phase", theme.ContentAddIcon(), func() {
		addPhaseRow("New Phase", 0)
	})

	footer := container.NewVBox(
		summaryLabel,
		container.NewGridWithColumns(2, addBtn, saveBtn),
	)

	configWindow.SetContent(container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("Blueprint Architect", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			container.NewGridWithColumns(2, nameEntry, totalEntry),
		),
		footer,
		nil, nil,
		scroll,
	))
	
	configWindow.Resize(fyne.NewSize(450, 550))
	configWindow.Show()
	updateSummary()
}
