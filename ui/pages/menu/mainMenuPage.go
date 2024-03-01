package menu

import (
	"fmt"
	"giogo/ui"
	"giogo/ui/component"
	routerModule "giogo/ui/router"
	"giogo/ui/styles"
	"image"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Menu struct {
	w                   *app.Window
	router              *routerModule.Router[ui.ApplicationCycles, string]
	size                image.Point
	container           *component.CentralizedContainer
	minesweeperClicable widget.Clickable
	scrollerClicable    widget.Clickable
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
	fmt.Println("Menu initialized")

	m.w.Option(func(_ unit.Metric, c *app.Config) {
		c.MaxSize = image.Point{}
		c.MinSize = image.Point{280, 200}
		c.Size = image.Point{600, 400}
		c.Title = "Menü"
		c.Decorated = true

		if m.size != (image.Point{}) {
			c.Size = m.size
		}
	})
}

func (m *Menu) Close() {
	m.w.Option(func(_ unit.Metric, c *app.Config) {
		m.size = c.Size
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

	return m.container.Layout(gtx,
		layout.Rigid(layout.Spacer{Height: 16}.Layout),
		layout.Rigid(material.Button(styles.MaterialTheme, &m.minesweeperClicable, "Minesweeper").Layout),
		layout.Rigid(layout.Spacer{Height: 6}.Layout),
		layout.Rigid(material.Button(styles.MaterialTheme, &m.scrollerClicable, "Scroller").Layout),
	)
}
