package menu

import (
	"fmt"
	"giogo/assets"
	"giogo/server"
	"giogo/ui"
	"giogo/ui/component"
	"giogo/ui/pages/minesweeper/engine"
	routerModule "giogo/ui/router"
	"giogo/ui/styles"
	"image"
	"strconv"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type MinesweeperMultiplayerMenu struct {
	w         *app.Window
	router    *routerModule.Router[ui.ApplicationCycles, string]
	container *component.CentralizedContainer
	// socketReader chan []byte
	client *server.MinesweeperServerClient

	loadingStart                time.Time
	isJoinableServerFound       bool
	failedToEstablishConnection bool

	joinEditor    widget.Editor
	joinClickable widget.Clickable
	hostClickable widget.Clickable
	backClickable widget.Clickable
}

var _ ui.ApplicationCycles = (*MinesweeperMultiplayerMenu)(nil)

const stateMarkMargin = 4

func NewMinesweeperMultiplayerMenu(w *app.Window, router *routerModule.Router[ui.ApplicationCycles, string]) *MinesweeperMultiplayerMenu {
	m := &MinesweeperMultiplayerMenu{
		w:          w,
		router:     router,
		container:  component.NewCentralizedContainer(false, true),
		client:     server.NewMinesweeperServerClient("", 0), // localhost:4222
		joinEditor: widget.Editor{Alignment: text.Start, SingleLine: true, MaxLen: 128, Submit: true},
	}

	m.joinEditor.Alignment = text.Middle

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

	m.loadingStart = time.Time{}
	m.isJoinableServerFound = false
	m.failedToEstablishConnection = false
	m.joinEditor.SetText("")
	m.joinEditor.ReadOnly = false
}

func (m *MinesweeperMultiplayerMenu) Close() {
	m.w.Option(func(_ unit.Metric, c *app.Config) {
		styles.MenuWindowSizes = c.Size
	})
}

func (m *MinesweeperMultiplayerMenu) Layout(gtx layout.Context) layout.Dimensions {
	var isSubmited bool = false

	if events := m.joinEditor.Events(); len(events) > 0 {
		_, isSubmited = events[0].(widget.SubmitEvent)
	}

	if m.loadingStart.IsZero() && (m.joinClickable.Clicked(gtx) || isSubmited) {
		m.joinClickHandler()
	}

	if m.hostClickable.Clicked(gtx) {
		m.router.GoTo(routerModule.MinesweeperCreateLobbyMenuPage)

		return layout.Dimensions{}
	}

	if m.backClickable.Clicked(gtx) {
		m.router.GoBack()

		return layout.Dimensions{}
	}

	// TODO: Részletezze miért nem tudott csatlakozni
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return m.container.Layout(gtx,
				layout.Rigid(layout.Spacer{Height: 16}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Max.X = 200

					macro := op.Record(gtx.Ops)
					editorDim := material.Editor(styles.MaterialTheme, &m.joinEditor, "192.168.1.1:8080").Layout(gtx)
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
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					var txt string = "Lobby keresés"

					if m.isJoinableServerFound {
						txt = "Csatlakozás lobbyhoz"
					}

					return material.Button(styles.MaterialTheme, &m.joinClickable, txt).Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: 6}.Layout),
				layout.Rigid(material.Button(styles.MaterialTheme, &m.hostClickable, "Játék létrehozás").Layout),
				layout.Flexed(1, layout.Spacer{}.Layout),
				layout.Rigid(material.Button(styles.CancelTheme, &m.backClickable, "Vissza").Layout),
				layout.Rigid(layout.Spacer{Height: 16}.Layout),
			)
		}),
		layout.Expanded(m.drawStackedStatus),
	)
}

func (*MinesweeperMultiplayerMenu) Restart() {}

func (m *MinesweeperMultiplayerMenu) joinClickHandler() {
	if !m.isJoinableServerFound {
		m.failedToEstablishConnection = false
		m.loadingStart = time.Now()

		go func() {
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
			joined, limit, canJoin := m.client.GetStatus()

			if limit != 0 {
				fmt.Printf("szerver stat: %d/%d | in Lobby (%t))\n", joined, limit, canJoin)

				if joined < limit && canJoin {
					m.isJoinableServerFound = true
					m.joinEditor.ReadOnly = true
					m.joinEditor.SetText(fmt.Sprintf("%d/%d", joined, limit))
				} else {
					m.failedToEstablishConnection = true
					fmt.Println("Szerver teli van!")
				}
			} else {
				m.failedToEstablishConnection = true
			}

			m.loadingStart = time.Time{}
		}()
	} else if !m.failedToEstablishConnection {
		clientEngine := engine.NewMinesweeperClientEngine(m.w, m.client.Host, m.client.Port)

		m.router.Add(routerModule.MinesweeperLobbyPage, NewMinesweeperLobby(m.w, m.router, clientEngine, nil, image.Point{}, 0))
		m.router.GoTo(routerModule.MinesweeperLobbyPage)
	}

}

func (m *MinesweeperMultiplayerMenu) drawStackedStatus(gtx layout.Context) layout.Dimensions {
	size := image.Point{24, 24}

	op.Offset(image.Point{
		gtx.Constraints.Max.X - size.X - stateMarkMargin,
		stateMarkMargin,
	}).Add(gtx.Ops)

	gtx.Constraints.Max = size

	if !m.loadingStart.IsZero() {
		radian := float32(time.Since(m.loadingStart).Seconds()) * 2
		affine2d := f32.Affine2D{}
		affine2d = affine2d.Rotate(f32.Pt(float32(size.X>>1), float32(size.Y>>1)), -radian)
		op.Affine(affine2d).Add(gtx.Ops)

		icon, _ := widget.NewIcon(icons.NotificationSync)
		icon.Layout(gtx, styles.BLUE)

		op.InvalidateOp{At: time.Now().Add(time.Millisecond * 10)}.Add(gtx.Ops)
	} else {
		if m.failedToEstablishConnection {
			assets.GetWidgetImage("joined_X", size.X).Layout(gtx)
		} else if m.isJoinableServerFound {
			assets.GetWidgetImage("joined_L", size.X).Layout(gtx)
		}
	}

	return layout.Dimensions{Size: size}
}
