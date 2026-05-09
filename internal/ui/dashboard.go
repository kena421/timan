package ui

import (
	"fmt"
	"strconv"
	"strings"
	"timan/internal/domain"
	"timan/internal/engine"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
)

type Dashboard struct {
	window fyne.Window
	engine *engine.TimerEngine
	store  *domain.BlueprintStore
	
	// State management
	mainContent *fyne.Container
}

func NewDashboard(ui *TimerUI, e *engine.TimerEngine) *Dashboard {
	return &Dashboard{
		engine: e,
		store:  domain.NewBlueprintStore(),
	}
}

func (d *Dashboard) Show() {
	if d.window == nil {
		d.window = fyne.CurrentApp().NewWindow("Timer Dashboard")
		d.mainContent = container.NewStack()
		d.window.SetContent(d.mainContent)
		d.window.Resize(fyne.NewSize(500, 600))
	}
	d.ShowLibrary()
	d.window.Show()
}

func (d *Dashboard) ShowLibrary() {
	blueprints, _ := d.store.LoadAll()

	if len(blueprints) == 0 {
		d.mainContent.Objects = []fyne.CanvasObject{container.NewCenter(
			container.NewVBox(
				widget.NewLabel("Your library is empty."),
				widget.NewButton("Create Your First Blueprint", func() { d.ShowEditor(nil) }),
			),
		)}
		d.mainContent.Refresh()
		return
	}

	// Auto-select first if none active
	if d.engine.GetCurrentBlueprintID() == "" {
		b := blueprints[0]
		wMins := b.WarningValue
		if b.WarningIsPercent {
			wMins = (b.WarningValue * b.Total) / 100
		}
		d.engine.UpdatePhases(b.Phases, wMins)
		d.engine.SetCurrentBlueprintID(b.ID)
	}

	list := container.NewVBox()
	for _, b := range blueprints {
		blueprint := b
		isActive := d.engine.GetCurrentBlueprintID() == blueprint.ID

		loadBtn := widget.NewButton("Select", func() {
			wMins := blueprint.WarningValue
			if blueprint.WarningIsPercent {
				wMins = (blueprint.WarningValue * blueprint.Total) / 100
			}
			d.engine.UpdatePhases(blueprint.Phases, wMins)
			d.engine.SetCurrentBlueprintID(blueprint.ID)
			d.ShowLibrary() // Refresh to show selection
		})
		loadBtn.Importance = widget.HighImportance

		editBtn := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
			d.ShowEditor(&blueprint)
		})
		deleteBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
			d.deleteBlueprint(blueprint.ID)
		})

		info := fmt.Sprintf("%s (%d mins)", blueprint.Name, blueprint.Total)
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

	addBtn := widget.NewButtonWithIcon("Create New Blueprint", theme.ContentAddIcon(), func() {
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
					d.engine.SetCurrentBlueprintID("quick-timer")
					d.ShowLibrary()
				}
			}
		}, d.window)
	})

	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("Blueprint Library", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewSeparator(),
		),
		container.NewVBox(simpleBtn, addBtn),
		nil, nil,
		container.NewVScroll(list),
	)

	d.mainContent.Objects = []fyne.CanvasObject{content}
	d.mainContent.Refresh()
}

func (d *Dashboard) ShowEditor(existing *domain.Blueprint) {
	title := "Blueprint Designer"
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
		nameEntry.SetText("New Blueprint")
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

	saveBtn := widget.NewButtonWithIcon("Save Blueprint", theme.ConfirmIcon(), func() {
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
			
			d.saveBlueprint(&domain.Blueprint{
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
		container.NewVBox(widget.NewLabel("Blueprint Name"), nameEntry),
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

func (d *Dashboard) saveBlueprint(saved *domain.Blueprint) {
	blueprints, _ := d.store.LoadAll()
	found := false
	for i, bp := range blueprints {
		if bp.ID == saved.ID {
			blueprints[i] = *saved
			found = true
			break
		}
	}
	if !found {
		blueprints = append(blueprints, *saved)
	}
	d.store.SaveAll(blueprints)
}

func (d *Dashboard) deleteBlueprint(id string) {
	blueprints, _ := d.store.LoadAll()
	var updated []domain.Blueprint
	for _, b := range blueprints {
		if b.ID != id {
			updated = append(updated, b)
		}
	}
	d.store.SaveAll(updated)
	d.ShowLibrary()
}
