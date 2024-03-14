package menu

import (
	"fmt"
	"giogo/server"
	"giogo/ui"
	"giogo/ui/component"
	"giogo/ui/pages/minesweeper/engine"
	routerModule "giogo/ui/router"
	"giogo/ui/styles"
	"image"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type MinesweeperCreateLobbyMenu struct {
	w         *app.Window
	router    *routerModule.Router[ui.ApplicationCycles, string]
	container *component.CentralizedContainer
	footer    *component.CentralizedContainer

	backClickable     widget.Clickable
	startClickable    widget.Clickable
	isPrivateCheckbox widget.Bool
	portEditor        widget.Editor
	playerLimitEditor widget.Editor
}

var _ ui.ApplicationCycles = (*MinesweeperCreateLobbyMenu)(nil)

const (
	private_host = "localhost"
	public_host  = "0.0.0.0"
)

func NewMinesweeperCreateLobbyMenu(w *app.Window, router *routerModule.Router[ui.ApplicationCycles, string]) *MinesweeperCreateLobbyMenu {
	m := &MinesweeperCreateLobbyMenu{
		w:                 w,
		router:            router,
		container:         component.NewCentralizedContainer(false, true),
		footer:            component.NewCentralizedContainer(true, false),
		isPrivateCheckbox: widget.Bool{Value: true},
	}

	m.footer.Container.Axis = layout.Horizontal
	m.footer.Container.Alignment = layout.End

	m.portEditor.SetText("4222")
	m.portEditor.InputHint = key.HintNumeric
	m.portEditor.Filter = "0123456789"
	m.portEditor.MaxLen = 7
	m.portEditor.Alignment = text.Middle
	m.portEditor.SingleLine = true

	m.playerLimitEditor.SetText("2")
	m.playerLimitEditor.InputHint = key.HintNumeric
	m.playerLimitEditor.Filter = "0123456789"
	m.playerLimitEditor.MaxLen = 3
	m.playerLimitEditor.Alignment = text.Middle
	m.playerLimitEditor.SingleLine = true

	return m
}

func (m *MinesweeperCreateLobbyMenu) Initialize() {
	m.w.Option(func(_ unit.Metric, c *app.Config) {
		c.MaxSize = image.Point{}
		c.MinSize = image.Point{280, 200}
		c.Size = styles.MenuWindowSizes
		c.Title = "Aknakereső Lobby"
		c.Decorated = true
	})
}

func (m *MinesweeperCreateLobbyMenu) Close() {
	m.w.Option(func(_ unit.Metric, c *app.Config) {
		styles.MenuWindowSizes = c.Size
	})
}

func (m *MinesweeperCreateLobbyMenu) Restart() {}

func (m *MinesweeperCreateLobbyMenu) Layout(gtx layout.Context) layout.Dimensions {
	res := m.handleEvents(&gtx)
	if res, isLD := res.(layout.Dimensions); isLD {
		return res
	}

	return m.container.Layout(gtx,
		layout.Rigid(layout.Spacer{Height: 16}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = 150

			macro := op.Record(gtx.Ops)
			editorDim := material.Editor(styles.MaterialTheme, &m.portEditor, "Port").Layout(gtx)
			recordCallOp := macro.Stop()

			editorDim.Size.X += 16
			editorDim.Size.Y += 16

			gtx.Constraints.Max = editorDim.Size
			paint.FillShape(gtx.Ops, styles.HEADER_BACKGROUND, clip.Rect{Max: gtx.Constraints.Max}.Op())
			defer op.Offset(image.Point{8, 8}).Push(gtx.Ops).Pop()

			recordCallOp.Add(gtx.Ops)

			return editorDim
		}),
		layout.Rigid(layout.Spacer{Height: 6}.Layout),
		layout.Rigid(material.CheckBox(styles.MaterialTheme, &m.isPrivateCheckbox, "Privát lobby").Layout),
		layout.Rigid(layout.Spacer{Height: 6}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = 150

			macro := op.Record(gtx.Ops)
			editorDim := material.Editor(styles.MaterialTheme, &m.playerLimitEditor, "Max játékosszám").Layout(gtx)
			recordCallOp := macro.Stop()

			editorDim.Size.X += 16
			editorDim.Size.Y += 16

			gtx.Constraints.Max = editorDim.Size
			paint.FillShape(gtx.Ops, styles.HEADER_BACKGROUND, clip.Rect{Max: gtx.Constraints.Max}.Op())
			defer op.Offset(image.Point{8, 8}).Push(gtx.Ops).Pop()

			recordCallOp.Add(gtx.Ops)

			return editorDim
		}),
		layout.Flexed(1, layout.Spacer{}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = image.Point{}
			gtx.Constraints.Max = image.Point{gtx.Constraints.Max.X, 75}

			return m.footer.Layout(gtx,
				layout.Flexed(1, layout.Spacer{}.Layout),
				layout.Rigid(material.Button(styles.MaterialTheme, &m.startClickable, "Indítás").Layout),
				layout.Flexed(2, layout.Spacer{}.Layout),
				layout.Rigid(material.Button(styles.CancelTheme, &m.backClickable, "Vissza").Layout),
				layout.Flexed(1, layout.Spacer{}.Layout),
			)
		}),
		layout.Rigid(layout.Spacer{Height: 16}.Layout),
	)
}

func (m *MinesweeperCreateLobbyMenu) onClickStart() (res interface{}) {
	var port uint
	fmt.Sscanf(m.portEditor.Text(), "%d", &port)
	if port < 1000 {
		fmt.Println("Lobby | invalid port number, must be above 1000")

		return
	}

	var host string = private_host
	if !m.isPrivateCheckbox.Value {
		host = public_host
	}

	var playerLimit uint
	fmt.Sscanf(m.playerLimitEditor.Text(), "%d", &playerLimit)
	if playerLimit == 0 || playerLimit > 255 {
		fmt.Println("Lobby | invalid player limit")

		return
	}

	mserver := server.NewMinesweeperServer(host, port, uint8(playerLimit))
	clientEngine := engine.NewMinesweeperClientEngine(host, port)

	m.router.Add(routerModule.MinesweeperLobbyPage, NewMinesweeperLobby(m.w, m.router, clientEngine, mserver))
	m.router.GoTo(routerModule.MinesweeperLobbyPage)

	res = layout.Dimensions{}

	return
}

func (m *MinesweeperCreateLobbyMenu) handleEvents(gtx *layout.Context) interface{} {
	var isSubmited bool = false

	if events := m.portEditor.Events(); len(events) > 0 {
		_, isSubmited = events[0].(widget.SubmitEvent)
	}

	for _, event := range m.playerLimitEditor.Events() {
		if _, isChange := event.(widget.ChangeEvent); isChange {
			var limit int
			fmt.Sscanf(m.playerLimitEditor.Text(), "%d", &limit)

			if limit > 255 {
				m.playerLimitEditor.SetText("255")
				m.playerLimitEditor.SetCaret(3, 3)

				continue
			}
		}
	}

	if m.startClickable.Clicked(*gtx) || isSubmited {
		return m.onClickStart()
	}

	if m.backClickable.Clicked(*gtx) {
		m.router.GoBack()

		return layout.Dimensions{}
	}

	return nil
}
