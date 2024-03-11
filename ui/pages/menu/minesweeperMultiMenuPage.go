package menu

import (
	"giogo/server"
	"giogo/ui"
	"giogo/ui/component"
	routerModule "giogo/ui/router"
	"giogo/ui/styles"
	"image"
	"strconv"
	"strings"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type MinesweeperMultiplayerMenu struct {
	w         *app.Window
	router    *routerModule.Router[ui.ApplicationCycles, string]
	container *component.CentralizedContainer
	// socketReader chan []byte
	client *server.MinesweeperServerClient

	testSocketText string
	joinEditor     widget.Editor
	joinClickable  widget.Clickable
	hostClickable  widget.Clickable
	backClickable  widget.Clickable
}

var _ ui.ApplicationCycles = (*MinesweeperMultiplayerMenu)(nil)

func NewMinesweeperMultiplayerMenu(w *app.Window, router *routerModule.Router[ui.ApplicationCycles, string]) *MinesweeperMultiplayerMenu {
	m := &MinesweeperMultiplayerMenu{
		w:          w,
		router:     router,
		container:  component.NewCentralizedContainer(false, true),
		client:     server.NewMinesweeperServerClient("localhost", 4222),
		joinEditor: widget.Editor{Alignment: text.Start, SingleLine: true, MaxLen: 128, Submit: true},
	}

	m.client.OnClosedConnection = func() {
		m.router.GoBack()
	}

	return m
}

func (m *MinesweeperMultiplayerMenu) Initialize() {
	m.w.Option(func(_ unit.Metric, c *app.Config) {
		c.MaxSize = image.Point{}
		c.MinSize = image.Point{280, 200}
		c.Size = styles.MenuWindowSizes
		c.Title = "Aknakereső Lobby"
		c.Decorated = true
	})
}

func (m *MinesweeperMultiplayerMenu) Close() {
	m.w.Option(func(_ unit.Metric, c *app.Config) {
		styles.MenuWindowSizes = c.Size
	})

	m.testSocketText = ""
	m.client.Disconnect()
}

func (m *MinesweeperMultiplayerMenu) Layout(gtx layout.Context) layout.Dimensions {
	var isSubmited bool = false

	if events := m.joinEditor.Events(); len(events) > 0 {
		_, isSubmited = events[0].(widget.SubmitEvent)
	}

	if m.joinClickable.Clicked(gtx) || isSubmited {
		txts := strings.Split(m.joinEditor.Text(), ":")
		m.client.Host = txts[0]
		m.client.Port = 80

		if len(txts) > 1 {
			port, err := strconv.ParseUint(txts[1], 10, 0)

			if err == nil {
				m.client.Port = uint(port)
			}
		}

		m.joinEditor.SetText("")
		m.client.Join()
		// return layout.Dimensions{}
	}

	if m.hostClickable.Clicked(gtx) {
		// TODO: Egy új oldalra átvisz, ahol meglehet adni milyen beállításokkal induljon le a játék
		//			 és felhasználók várhatnak itt az indulásra.
		//			 Számukra nem szerkeszthető módon látják a beállításokat, csat a tényleges host tudja azokat írni és broadcastolja a belépett felhasználók számára.

		socketData := server.SocketData{}
		socketData.DataType = server.POSITION
		socketData.Data = []byte{0, 1, 2, 3, 4, 5, 6, 100, 101}
		m.client.WriteData(socketData.ToBytes())

		// return layout.Dimensions{}
	}

	if m.backClickable.Clicked(gtx) {
		m.router.GoBack()

		return layout.Dimensions{}
	}

	return m.container.Layout(gtx,
		layout.Rigid(layout.Spacer{Height: 16}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = 200

			macro := op.Record(gtx.Ops)
			editorDim := material.Editor(styles.MaterialTheme, &m.joinEditor, "192.168.1.1").Layout(gtx)
			recordCallOp := macro.Stop()

			editorDim.Size.X = editorDim.Size.X + 16
			editorDim.Size.Y += 16

			gtx.Constraints.Max = editorDim.Size
			paint.FillShape(gtx.Ops, styles.HEADER_BACKGROUND, clip.Rect{Max: gtx.Constraints.Max}.Op())
			defer op.Offset(image.Point{8, 8}).Push(gtx.Ops).Pop()

			recordCallOp.Add(gtx.Ops)

			return editorDim
		}),
		layout.Rigid(layout.Spacer{Height: 4}.Layout),
		layout.Rigid(material.Button(styles.MaterialTheme, &m.joinClickable, "Csatlakozás játékhoz").Layout),
		layout.Rigid(layout.Spacer{Height: 6}.Layout),
		layout.Rigid(material.Button(styles.MaterialTheme, &m.hostClickable, "Játék létrehozás").Layout),
		layout.Rigid(layout.Spacer{Height: 6}.Layout),
		layout.Rigid(material.Label(styles.MaterialTheme, unit.Sp(16), m.testSocketText).Layout),
		layout.Flexed(1, layout.Spacer{}.Layout),
		layout.Rigid(material.Button(styles.MaterialTheme, &m.backClickable, "Vissza").Layout),
		layout.Rigid(layout.Spacer{Height: 16}.Layout),
	)
}

func (*MinesweeperMultiplayerMenu) Restart() {}
