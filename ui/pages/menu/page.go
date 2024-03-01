package menu

import (
	"fmt"
	"giogo/ui"
	routerModule "giogo/ui/router"
	"giogo/ui/styles"
	"image"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Menu struct {
	w                   *app.Window
	router              *routerModule.Router[ui.ApplicationCycles, string]
	size                image.Point
	container           layout.Flex
	minesweeperClicable widget.Clickable
	scrollerClicable    widget.Clickable
}

var _ ui.ApplicationCycles = (*Menu)(nil)

func NewMenu(w *app.Window, router *routerModule.Router[ui.ApplicationCycles, string]) *Menu {
	m := &Menu{
		w:                   w,
		router:              router,
		container:           layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle},
		minesweeperClicable: widget.Clickable{},
		scrollerClicable:    widget.Clickable{},
	}

	return m
}

func (m *Menu) Initialize() {
	fmt.Println("Menu initialized")

	m.w.Option(func(_ unit.Metric, c *app.Config) {
		c.MaxSize = image.Point{}
		c.MinSize = image.Point{80, 140}
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
	var marginInset unit.Dp = 0
	var marginBlock unit.Dp = 0

	if m.minesweeperClicable.Clicked(gtx) {
		m.router.GoTo(routerModule.MinesweeperPage)

		return layout.Dimensions{}
	}

	if m.scrollerClicable.Clicked(gtx) {
		m.router.GoTo(routerModule.ScrollerPage)

		return layout.Dimensions{}
	}

	rec := op.Record(gtx.Ops)
	cDims := m.container.Layout(gtx,
		layout.Rigid(material.Button(styles.MaterialTheme, &m.minesweeperClicable, "Minesweeper").Layout),
		layout.Rigid(layout.Spacer{Height: 6}.Layout),
		layout.Rigid(material.Button(styles.MaterialTheme, &m.scrollerClicable, "Scroller").Layout),
	)

	marginInset = unit.Dp((gtx.Constraints.Max.X - cDims.Size.X) >> 1)
	marginBlock = unit.Dp((gtx.Constraints.Max.Y - cDims.Size.Y) >> 1)

	macro := rec.Stop()
	return layout.Inset{Left: marginInset, Right: marginInset, Top: marginBlock, Bottom: marginBlock}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		macro.Add(gtx.Ops)

		return cDims
	})
}
