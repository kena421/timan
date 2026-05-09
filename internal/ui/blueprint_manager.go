package ui

import (
	"fmt"
	"timan/internal/domain"
	"timan/internal/engine"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
)

type BlueprintManager struct {
	window fyne.Window
	engine *engine.TimerEngine
	store  *domain.BlueprintStore
	ui     *TimerUI
}

func NewBlueprintManager(ui *TimerUI, e *engine.TimerEngine) *BlueprintManager {
	return &BlueprintManager{
		ui:     ui,
		engine: e,
		store:  domain.NewBlueprintStore(),
	}
}

func (m *BlueprintManager) Show() {
	m.window = fyne.CurrentApp().NewWindow("Blueprint Library")
	m.refresh()
	m.window.Resize(fyne.NewSize(400, 500))
	m.window.Show()
}

func (m *BlueprintManager) refresh() {
	blueprints, _ := m.store.LoadAll()

	list := container.NewVBox()
	for _, b := range blueprints {
		blueprint := b // capture for closure
		loadBtn := widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
			m.engine.UpdatePhases(blueprint.Phases)
			m.window.Close()
		})
		editBtn := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
			m.showEditor(&blueprint)
		})
		deleteBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
			m.deleteBlueprint(blueprint.ID)
		})

		info := fmt.Sprintf("%s (%d mins)", blueprint.Name, blueprint.Total)
		row := container.NewBorder(nil, nil, nil, container.NewHBox(loadBtn, editBtn, deleteBtn), widget.NewLabel(info))
		list.Add(row)
	}

	addBtn := widget.NewButtonWithIcon("Create New Blueprint", theme.ContentAddIcon(), func() {
		m.showEditor(nil)
	})

	m.window.SetContent(container.NewBorder(
		widget.NewLabelWithStyle("Your Blueprints", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		addBtn,
		nil, nil,
		container.NewVScroll(list),
	))
}

func (m *BlueprintManager) showEditor(b *domain.Blueprint) {
	// Re-using the logic from TimerUI.showConfig but adapted for CRUD
	// For simplicity, I'll just open the existing showConfig logic but wrapped with "Save to Library"
    // Actually, let's build a dedicated editor here that saves to the store.
    m.ui.showBlueprintEditor(b, func(saved *domain.Blueprint) {
        blueprints, _ := m.store.LoadAll()
        if b == nil {
            saved.ID = uuid.New().String()
            blueprints = append(blueprints, *saved)
        } else {
            for i, bp := range blueprints {
                if bp.ID == b.ID {
                    blueprints[i] = *saved
                    blueprints[i].ID = b.ID
                    break
                }
            }
        }
        m.store.SaveAll(blueprints)
        m.refresh()
    })
}

func (m *BlueprintManager) deleteBlueprint(id string) {
	blueprints, _ := m.store.LoadAll()
	var updated []domain.Blueprint
	for _, b := range blueprints {
		if b.ID != id {
			updated = append(updated, b)
		}
	}
	m.store.SaveAll(updated)
	m.refresh()
}
