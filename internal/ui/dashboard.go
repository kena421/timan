package ui

import (
	"fmt"
	"strconv"
	"strings"
	"github.com/timan-org/timan/internal/domain"
	"github.com/timan-org/timan/internal/engine"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
)

// Dashboard provides a user interface for managing Event templates and library.
type Dashboard struct {
	window fyne.Window
	engine *engine.TimerEngine
	store  *domain.EventStore
	
	// State management
	mainContent *fyne.Container
}

// NewDashboard creates a new Dashboard instance.
func NewDashboard(ui *TimerUI, e *engine.TimerEngine) *Dashboard {
	return &Dashboard{
		engine: e,
		store:  domain.NewEventStore(),
	}
}

// Show displays the Dashboard window.
func (d *Dashboard) Show() {
	if d.window != nil {
		d.window.Close()
	}

	d.window = fyne.CurrentApp().NewWindow("Timan Dashboard")
	d.mainContent = container.NewStack()
	d.window.SetContent(d.mainContent)
	d.window.Resize(fyne.NewSize(500, 600))
	
	d.window.SetOnClosed(func() {
		d.window = nil
	})

	d.ShowLibrary()
	d.window.Show()
	d.window.RequestFocus()
}

// ShowLibrary displays the list of saved Event templates.
func (d *Dashboard) ShowLibrary() {
	events, _ := d.store.LoadAll()

	if len(events) == 0 {
		d.mainContent.Objects = []fyne.CanvasObject{container.NewCenter(
			container.NewVBox(
				widget.NewLabel("Your library is empty."),
				widget.NewButton("Create Your First Event", func() { d.ShowEditor(nil) }),
			),
		)}
		d.mainContent.Refresh()
		return
	}

	// Auto-select first if none active
	if d.engine.GetCurrentEventID() == "" {
		e := events[0]
		wMins := e.WarningValue
		if e.WarningIsPercent {
			wMins = (e.WarningValue * e.Total) / 100
		}
		d.engine.UpdatePhases(e.Phases, wMins)
		d.engine.SetCurrentEventID(e.ID)
	}

	list := container.NewVBox()
	for _, e := range events {
		event := e
		isActive := d.engine.GetCurrentEventID() == event.ID

		loadBtn := widget.NewButton("Select", func() {
			wMins := event.WarningValue
			if event.WarningIsPercent {
				wMins = (event.WarningValue * event.Total) / 100
			}
			d.engine.UpdatePhases(event.Phases, wMins)
			d.engine.SetCurrentEventID(event.ID)
			d.ShowLibrary() // Refresh to show selection
		})
		loadBtn.Importance = widget.HighImportance

		editBtn := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
			d.ShowEditor(&event)
		})
		deleteBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
			d.deleteEvent(event.ID)
		})

		info := fmt.Sprintf("%s (%d mins)", event.Name, event.Total)
		label := widget.NewLabel(info)
		if isActive {
			label.Importance = widget.SuccessImportance
			row := container.NewBorder(nil, nil, widget.NewIcon(theme.ConfirmIcon()), container.NewHBox(loadBtn, editBtn, deleteBtn), label)
			list.Add(row)
		} else {
			row := container.NewBorder(nil, nil, nil, container.NewHBox(loadBtn, editBtn, deleteBtn), label)
			list.Add(row)
		}
	}

	addBtn := widget.NewButtonWithIcon("Create New Event", theme.ContentAddIcon(), func() {
		d.ShowEditor(nil)
	})

	simpleBtn := widget.NewButtonWithIcon("Simple Quick Timer", theme.HistoryIcon(), func() {
		minsEntry := widget.NewEntry()
		minsEntry.SetText("60")
		
		warnEntry := widget.NewEntry()
		warnEntry.SetText("5")

		dialog.ShowForm("Simple Timer", "Start", "Cancel", []*widget.FormItem{
			{Text: "Total Minutes", Widget: minsEntry},
			{Text: "Alert at (mins remaining)", Widget: warnEntry},
		}, func(ok bool) {
			if ok {
				mins, _ := strconv.Atoi(minsEntry.Text)
				warn, _ := strconv.Atoi(warnEntry.Text)
				if mins > 0 {
					d.engine.UpdatePhases([]domain.Phase{{Name: "Timer", Duration: mins * 60}}, warn)
					d.engine.SetCurrentEventID("quick-timer")
					d.ShowLibrary()
				}
			}
		}, d.window)
	})

	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("Event Library", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewSeparator(),
		),
		container.NewVBox(simpleBtn, addBtn),
		nil, nil,
		container.NewVScroll(list),
	)

	d.mainContent.Objects = []fyne.CanvasObject{content}
	d.mainContent.Refresh()
}

// ShowEditor provides a form to create or edit an Event template.
func (d *Dashboard) ShowEditor(existing *domain.Event) {
	title := "Event Designer"
	if existing != nil {
		title = "Edit: " + existing.Name
	}

	phases := d.engine.GetPhases()
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
		nameEntry.SetText("New Event")
	}

	totalEntry := widget.NewEntry()
	totalEntry.SetText(strconv.Itoa(totalMins))
	totalEntry.PlaceHolder = "Total Duration (mins)"

	warningEntry := widget.NewEntry()
	if existing != nil {
		val := strconv.Itoa(existing.WarningValue)
		if existing.WarningIsPercent {
			val += "%"
		}
		warningEntry.SetText(val)
	} else {
		warningEntry.SetText("5")
	}
	warningEntry.PlaceHolder = "Alert at (mins or %)"

	rows := container.NewVBox()
	summaryLabel := widget.NewLabel("")

	phaseHeader := container.NewGridWithColumns(2,
		widget.NewLabelWithStyle("Phase Name", fyne.TextAlignLeading, fyne.TextStyle{Italic: true}),
		widget.NewLabelWithStyle("Duration (mins or %)", fyne.TextAlignLeading, fyne.TextStyle{Italic: true}),
	)

	updateSummary := func() {
		currentSumSec := 0
		targetMins, _ := strconv.Atoi(totalEntry.Text)

		for _, row := range rows.Objects {
			if box, ok := row.(*fyne.Container); ok {
				if grid, ok := box.Objects[0].(*fyne.Container); ok {
					if minsEntry, ok := grid.Objects[1].(*widget.Entry); ok {
						txt := strings.TrimSpace(minsEntry.Text)
						if strings.HasSuffix(txt, "%") {
							p, _ := strconv.Atoi(strings.TrimSuffix(txt, "%"))
							currentSumSec += (p * targetMins * 60) / 100
						} else {
							m, _ := strconv.Atoi(txt)
							currentSumSec += m * 60
						}
					}
				}
			}
		}

		allocatedMins := currentSumSec / 60
		remainingMins := targetMins - allocatedMins
		
		summaryText := fmt.Sprintf("Allocated: %d / %d mins", allocatedMins, targetMins)
		if remainingMins > 0 {
			summaryText += fmt.Sprintf(" (%d mins remaining)", remainingMins)
			summaryLabel.Importance = widget.WarningImportance
		} else if remainingMins < 0 {
			summaryText += fmt.Sprintf(" (%d mins OVER)", -remainingMins)
			summaryLabel.Importance = widget.DangerImportance
		} else {
			summaryText += " (Perfect!)"
			summaryLabel.Importance = widget.SuccessImportance
		}
		
		summaryLabel.SetText(summaryText)
		summaryLabel.Refresh()
	}

	addPhaseRow := func(name string, val string) {
		pNameEntry := widget.NewEntry()
		pNameEntry.SetText(name)
		pNameEntry.PlaceHolder = "Phase Name"
		pNameEntry.OnChanged = func(string) { updateSummary() }

		minsEntry := widget.NewEntry()
		minsEntry.SetText(val)
		minsEntry.PlaceHolder = "Mins or %"
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
		val := strconv.Itoa(p.Duration / 60)
		if p.IsPercent {
			val = fmt.Sprintf("%d%%", p.Percent)
		}
		addPhaseRow(p.Name, val)
	}

	totalEntry.OnChanged = func(string) { updateSummary() }

	saveBtn := widget.NewButtonWithIcon("Save Event", theme.ConfirmIcon(), func() {
		newPhases := []domain.Phase{}
		targetMins, _ := strconv.Atoi(totalEntry.Text)
		
		for _, row := range rows.Objects {
			if box, ok := row.(*fyne.Container); ok {
				if grid, ok := box.Objects[0].(*fyne.Container); ok {
					pNameEntry := grid.Objects[0].(*widget.Entry)
					minsEntry := grid.Objects[1].(*widget.Entry)
					
					txt := strings.TrimSpace(minsEntry.Text)
					phase := domain.Phase{Name: pNameEntry.Text}
					if strings.HasSuffix(txt, "%") {
						p, _ := strconv.Atoi(strings.TrimSuffix(txt, "%"))
						phase.IsPercent = true
						phase.Percent = p
						phase.Duration = (p * targetMins * 60) / 100
					} else {
						m, _ := strconv.Atoi(txt)
						phase.Duration = m * 60
					}
					
					if phase.Name != "" && phase.Duration > 0 {
						newPhases = append(newPhases, phase)
					}
				}
			}
		}

		if err := domain.ValidatePhases(newPhases, targetMins); err == nil {
			txt := strings.TrimSpace(warningEntry.Text)
			isPercent := strings.HasSuffix(txt, "%")
			val, _ := strconv.Atoi(strings.TrimSuffix(txt, "%"))
			
			d.saveEvent(&domain.Event{
				ID:               func() string { if existing != nil { return existing.ID }; return uuid.New().String() }(),
				Name:             nameEntry.Text,
				Total:            targetMins,
				WarningValue:     val,
				WarningIsPercent: isPercent,
				Phases:           newPhases,
			})
			d.ShowLibrary()
		} else {
			dialog.ShowError(err, d.window)
		}
	})

	cancelBtn := widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
		d.ShowLibrary()
	})

	addBtn := widget.NewButtonWithIcon("Add Phase", theme.ContentAddIcon(), func() {
		addPhaseRow("New Phase", "0")
	})

	header := container.NewGridWithColumns(3,
		container.NewVBox(widget.NewLabel("Event Name"), nameEntry),
		container.NewVBox(widget.NewLabel("Total Duration (mins)"), totalEntry),
		container.NewVBox(widget.NewLabel("Alert Threshold (mins or %)"), warningEntry),
	)

	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle(title, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			header,
			widget.NewSeparator(),
			phaseHeader,
		),
		container.NewVBox(summaryLabel, container.NewGridWithColumns(3, addBtn, cancelBtn, saveBtn)),
		nil, nil,
		container.NewVScroll(rows),
	)

	d.mainContent.Objects = []fyne.CanvasObject{content}
	d.mainContent.Refresh()
	updateSummary()
}

func (d *Dashboard) saveEvent(saved *domain.Event) {
	events, _ := d.store.LoadAll()
	found := false
	for i, ev := range events {
		if ev.ID == saved.ID {
			events[i] = *saved
			found = true
			break
		}
	}
	if !found {
		events = append(events, *saved)
	}
	d.store.SaveAll(events)
}

func (d *Dashboard) deleteEvent(id string) {
	events, _ := d.store.LoadAll()
	var updated []domain.Event
	for _, e := range events {
		if e.ID != id {
			updated = append(updated, e)
		}
	}
	d.store.SaveAll(updated)
	d.ShowLibrary()
}
