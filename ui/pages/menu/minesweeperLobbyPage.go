package menu

import (
	"fmt"
	"image"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"giogo/server"
	"giogo/ui"
	"giogo/ui/component"
	"giogo/ui/pages/minesweeper"
	"giogo/ui/pages/minesweeper/engine"
	"giogo/ui/pages/minesweeper/model"
	routerModule "giogo/ui/router"
	"giogo/ui/styles"
)

type MinesweeperLobby struct {
	w         *app.Window
	router    *routerModule.Router[ui.ApplicationCycles, string]
	container *component.CentralizedContainer
	footer    *component.CentralizedContainer
	gameSize  image.Point
	mines     uint16

	server       *server.MinesweeperServer
	clientEngine *engine.MinesweeperClientEngine

	startClickable widget.Clickable
	exitClickable  widget.Clickable

	minePage *minesweeper.MineField
}

var _ ui.ApplicationCycles = (*MinesweeperLobby)(nil)

func NewMinesweeperLobby(w *app.Window, router *routerModule.Router[ui.ApplicationCycles, string], clientEngine *engine.MinesweeperClientEngine, server *server.MinesweeperServer, gameSize image.Point, mines uint16) *MinesweeperLobby {
	m := &MinesweeperLobby{
		w:            w,
		router:       router,
		gameSize:     gameSize,
		mines:        mines,
		server:       server,
		clientEngine: clientEngine,
		container:    component.NewCentralizedContainer(false, true),
		footer:       component.NewCentralizedContainer(true, false),
		minePage:     minesweeper.NewMinefield(w, router, 20).SetEngine(clientEngine),
	}

	m.footer.Container.Alignment = layout.Middle
	m.footer.Container.Axis = layout.Horizontal

	return m
}

func (m *MinesweeperLobby) Initialize() {
	if m.server != nil {
		m.server.Open()

		m.clientEngine.Client.Host = m.server.GetHost()
		m.clientEngine.Client.Port = m.server.GetPort()
	}

	m.minePage.Initialize()

	m.w.Option(func(_ unit.Metric, c *app.Config) {
		c.MaxSize = image.Point{}
		c.MinSize = image.Point{280, 200}
		c.Size = styles.MenuWindowSizes
		c.Title = "Aknakereső Lobby"
		c.Decorated = true
	})
}

func (m *MinesweeperLobby) Restart() {
	m.minePage.Restart()
}

func (m *MinesweeperLobby) Close() {
	m.w.Option(func(_ unit.Metric, c *app.Config) {
		styles.MenuWindowSizes = c.Size
	})

	m.router.Remove(routerModule.MinesweeperLobbyPage)
	m.clientEngine.Close()

	if m.server != nil {
		m.server.Close()
	}

	m.minePage.Close()
}

func (m *MinesweeperLobby) Layout(gtx layout.Context) layout.Dimensions {
	if m.clientEngine.GetState() != model.WAITING {
		return m.minePage.Layout(gtx)
	}

	m.handleEvents(&gtx)

	// TODO: Felhasználók név hozzáfűzése és azok megjelenítése
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = image.Point{}

			macro := op.Record(gtx.Ops)
			labelDim := material.Label(styles.MaterialTheme, unit.Sp(12), fmt.Sprintf("%s:%d", m.clientEngine.Client.Host, m.clientEngine.Client.Port)).Layout(gtx)
			labelCallOp := macro.Stop()

			op.Offset(image.Point{4, gtx.Constraints.Max.Y - labelDim.Size.Y - 4}).Add(gtx.Ops)

			labelCallOp.Add(gtx.Ops)

			return labelDim
		}),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = image.Point{}

			macro := op.Record(gtx.Ops)
			labelDim := material.Label(styles.MaterialTheme, unit.Sp(12), fmt.Sprintf("%d/%d", m.clientEngine.ServerStatus.Joined, m.clientEngine.ServerStatus.Limit)).Layout(gtx)
			labelCallOp := macro.Stop()

			op.Offset(image.Point{4, 4}).Add(gtx.Ops)

			labelCallOp.Add(gtx.Ops)

			return labelDim
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return m.container.Layout(gtx,
				layout.Rigid(layout.Spacer{Height: 16}.Layout),
				layout.Rigid(material.Label(styles.MaterialTheme, styles.MaterialTheme.TextSize, "Lobby").Layout),
				layout.Rigid(layout.Spacer{Height: 16}.Layout),
				layout.Flexed(1, layout.Spacer{}.Layout),
				layout.Rigid(layout.Spacer{Height: 16}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Max.Y = 75

					return m.footer.Layout(gtx, showStartButton(m.server != nil, &m.startClickable, &m.exitClickable)...)
				}),
				layout.Rigid(layout.Spacer{Height: 16}.Layout),
			)
		}),
	)
}

func (m *MinesweeperLobby) handleEvents(gtx *layout.Context) {
	if m.exitClickable.Clicked(*gtx) {
		m.router.GoBackTo(routerModule.MinesweeperMenuPage)
	}

	if m.startClickable.Clicked(*gtx) {
		m.server.DisableJoin()
		m.clientEngine.Resize(uint16(m.gameSize.X), uint16(m.gameSize.Y), m.mines)
	}
}

func showStartButton(canShow bool, startClickable *widget.Clickable, exitClickable *widget.Clickable) []layout.FlexChild {
	if canShow {
		return []layout.FlexChild{
			layout.Flexed(1, layout.Spacer{}.Layout),
			layout.Rigid(material.Button(styles.MaterialTheme, startClickable, "Start").Layout),
			layout.Flexed(2, layout.Spacer{}.Layout),
			layout.Rigid(material.Button(styles.CancelTheme, exitClickable, "Kilépés").Layout),
			layout.Flexed(1, layout.Spacer{}.Layout),
		}
	}

	return []layout.FlexChild{
		layout.Flexed(1, layout.Spacer{}.Layout),
		layout.Rigid(material.Button(styles.CancelTheme, exitClickable, "Kilépés").Layout),
		layout.Flexed(1, layout.Spacer{}.Layout),
	}
}
