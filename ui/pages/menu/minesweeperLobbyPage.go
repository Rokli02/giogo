package menu

import (
	"image"

	"gioui.org/app"
	"gioui.org/layout"
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

	server       *server.MinesweeperServer
	clientEngine *engine.MinesweeperClientEngine

	startClickable widget.Clickable
	exitClickable  widget.Clickable

	minePage *minesweeper.MineField
}

var _ ui.ApplicationCycles = (*MinesweeperLobby)(nil)

func NewMinesweeperLobby(w *app.Window, router *routerModule.Router[ui.ApplicationCycles, string], clientEngine *engine.MinesweeperClientEngine, server *server.MinesweeperServer) *MinesweeperLobby {
	m := &MinesweeperLobby{
		w:            w,
		router:       router,
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

		<-m.server.HealthCheckChan

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
	m.handleEvents(&gtx)

	if m.clientEngine.GetState() == model.WAITING {
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
	}

	return m.minePage.Layout(gtx)
}

func (m *MinesweeperLobby) handleEvents(gtx *layout.Context) {
	if m.exitClickable.Clicked(*gtx) {
		m.router.GoBackTo(routerModule.MinesweeperMenuPage)
	}

	if m.startClickable.Clicked(*gtx) {
		m.server.DisableJoin()
		m.clientEngine.Resize(8, 10, 10)
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
