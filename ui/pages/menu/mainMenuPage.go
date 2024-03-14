package menu

import (
	"giogo/ui"
	"giogo/ui/component"
	routerModule "giogo/ui/router"
	"giogo/ui/styles"
	"image"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Menu struct {
	w                   *app.Window
	router              *routerModule.Router[ui.ApplicationCycles, string]
	container           *component.CentralizedContainer
	minesweeperClicable widget.Clickable
	scrollerClicable    widget.Clickable
	exitClicable        widget.Clickable
}

var _ ui.ApplicationCycles = (*Menu)(nil)

func NewMenu(w *app.Window, router *routerModule.Router[ui.ApplicationCycles, string]) *Menu {
	m := &Menu{
		w:         w,
		router:    router,
		container: component.NewCentralizedContainer(false, true),
	}

	return m
}

func (m *Menu) Initialize() {
	m.w.Option(func(_ unit.Metric, c *app.Config) {
		c.MaxSize = image.Point{}
		c.MinSize = image.Point{280, 200}
		c.Size = styles.MenuWindowSizes
		c.Title = "Menü"
		c.Decorated = true
	})
}

func (m *Menu) Close() {
	m.w.Option(func(_ unit.Metric, c *app.Config) {
		styles.MenuWindowSizes = c.Size
	})
}

func (*Menu) Restart() {}

func (m *Menu) Layout(gtx layout.Context) layout.Dimensions {
	if m.minesweeperClicable.Clicked(gtx) {
		m.router.GoTo(routerModule.MinesweeperMenuPage)

		return layout.Dimensions{}
	}

	if m.scrollerClicable.Clicked(gtx) {
		m.router.GoTo(routerModule.ScrollerPage)

		return layout.Dimensions{}
	}

	if m.exitClicable.Clicked(gtx) {
		m.w.Perform(system.ActionClose)

		return layout.Dimensions{}
	}

	return m.container.Layout(gtx,
		layout.Rigid(layout.Spacer{Height: 16}.Layout),
		layout.Rigid(material.Button(styles.MaterialTheme, &m.minesweeperClicable, "Minesweeper").Layout),
		layout.Rigid(layout.Spacer{Height: 6}.Layout),
		layout.Rigid(material.Button(styles.MaterialTheme, &m.scrollerClicable, "Scroller").Layout),
		layout.Flexed(1, layout.Spacer{}.Layout),
		layout.Rigid(material.Button(styles.CancelTheme, &m.exitClicable, "Kilépés").Layout),
		layout.Rigid(layout.Spacer{Height: 16}.Layout),
	)
}
