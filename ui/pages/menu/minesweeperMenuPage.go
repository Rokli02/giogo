package menu

import (
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

type MinesweeperMenu struct {
	w         *app.Window
	router    *routerModule.Router[ui.ApplicationCycles, string]
	container *component.CentralizedContainer

	singlePlayerClickable widget.Clickable
	multiPlayerClickable  widget.Clickable
	backClickable         widget.Clickable
}

var _ ui.ApplicationCycles = (*MinesweeperMenu)(nil)

func NewMinesweeperMenu(w *app.Window, router *routerModule.Router[ui.ApplicationCycles, string]) *MinesweeperMenu {
	m := &MinesweeperMenu{
		w:         w,
		router:    router,
		container: component.NewCentralizedContainer(false, true),
	}

	return m
}

func (m *MinesweeperMenu) Initialize() {
	m.w.Option(func(_ unit.Metric, c *app.Config) {
		c.MaxSize = image.Point{}
		c.MinSize = image.Point{280, 200}
		c.Size = styles.MenuWindowSizes
		c.Title = "Aknakereső Menü"
		c.Decorated = true
	})
}

func (m *MinesweeperMenu) Close() {
	m.w.Option(func(_ unit.Metric, c *app.Config) {
		styles.MenuWindowSizes = c.Size
	})
}

// Layout implements ui.ApplicationCycles.
func (m *MinesweeperMenu) Layout(gtx layout.Context) layout.Dimensions {
	if m.singlePlayerClickable.Clicked(gtx) {
		m.router.GoTo(routerModule.MinesweeperPage)

		return layout.Dimensions{}
	}

	if m.multiPlayerClickable.Clicked(gtx) {
		m.router.GoTo(routerModule.MinesweeperMultiplayerMenuPage)

		return layout.Dimensions{}
	}

	if m.backClickable.Clicked(gtx) {
		m.router.GoBack()

		return layout.Dimensions{}
	}

	return m.container.Layout(gtx,
		layout.Rigid(layout.Spacer{Height: 16}.Layout),
		layout.Rigid(material.Button(styles.MaterialTheme, &m.singlePlayerClickable, "Egyjátékos mód").Layout),
		layout.Rigid(layout.Spacer{Height: 6}.Layout),
		layout.Rigid(material.Button(styles.MaterialTheme, &m.multiPlayerClickable, "Többjátékos mód").Layout),
		layout.Flexed(1, layout.Spacer{}.Layout),
		layout.Rigid(material.Button(styles.MaterialTheme, &m.backClickable, "Vissza").Layout),
		layout.Rigid(layout.Spacer{Height: 16}.Layout),
	)
}

// Restart implements ui.ApplicationCycles.
func (*MinesweeperMenu) Restart() {}
