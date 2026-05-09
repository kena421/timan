package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type InteractionWrapper struct {
	widget.BaseWidget
	OnTap  func()
	Menu   *fyne.Menu
	Window fyne.Window
}

func (i *InteractionWrapper) Tapped(_ *fyne.PointEvent) {
	if i.OnTap != nil {
		i.OnTap()
	}
}

func (i *InteractionWrapper) TappedSecondary(e *fyne.PointEvent) {
	widget.ShowPopUpMenuAtPosition(i.menuFix(), i.Window.Canvas(), e.AbsolutePosition)
}

func (i *InteractionWrapper) menuFix() *fyne.Menu {
    return i.Menu
}

func (i *InteractionWrapper) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(canvas.NewRectangle(color.Transparent))
}
