package menu

import (
	"giogo/ui"
	"giogo/ui/component"
	routerModule "giogo/ui/router"
	"giogo/ui/styles"
	"image"

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

	joinEditor    widget.Editor
	joinClickable widget.Clickable
	hostClickable widget.Clickable
	backClickable widget.Clickable
}

var _ ui.ApplicationCycles = (*MinesweeperMultiplayerMenu)(nil)

func NewMinesweeperMultiplayerMenu(w *app.Window, router *routerModule.Router[ui.ApplicationCycles, string]) *MinesweeperMultiplayerMenu {
	m := &MinesweeperMultiplayerMenu{
		w:          w,
		router:     router,
		container:  component.NewCentralizedContainer(false, true),
		joinEditor: widget.Editor{Alignment: text.Start, SingleLine: true, MaxLen: 128, Submit: true},
	}

	return m
}

func (m *MinesweeperMultiplayerMenu) Initialize() {
	m.w.Option(func(_ unit.Metric, c *app.Config) {
		c.MaxSize = image.Point{}
		c.MinSize = image.Point{280, 200}
		c.Size = styles.MenuWindowSizes
		c.Title = "Multiplayer Aknakereső"
		c.Decorated = true
	})
}

func (m *MinesweeperMultiplayerMenu) Close() {
	m.w.Option(func(_ unit.Metric, c *app.Config) {
		styles.MenuWindowSizes = c.Size
	})
}

func (m *MinesweeperMultiplayerMenu) Layout(gtx layout.Context) layout.Dimensions {

	// for _, event := range m.joinEditor.Events() {
	// 	switch evt := event.(type) {
	// 	case widget.ChangeEvent:
	// 		fmt.Println("Change event:", evt)
	// 	case widget.SubmitEvent:
	// 		fmt.Println("Submit event:", evt)
	// 	default:
	// 		fmt.Println("Editor event:", evt)
	// 	}
	// }
	var isSubmited bool = false

	if events := m.joinEditor.Events(); len(events) > 0 {
		_, isSubmited = events[0].(widget.SubmitEvent)
	}

	if m.joinClickable.Clicked(gtx) || isSubmited {
		// TODO: Szöveg mező legyen felette, amibe bele lehet írni az IP-t és arra megpróbál csatlakozni
		m.joinEditor.SetText("")

		// return layout.Dimensions{}
	}

	if m.hostClickable.Clicked(gtx) {
		// TODO: Egy új oldalra átvisz, ahol meglehet adni milyen beállításokkal induljon le a játék
		//			 és felhasználók várhatnak itt az indulásra.
		//			 Számukra nem szerkeszthető módon látják a beállításokat, csat a tényleges host tudja azokat írni és broadcastolja a belépett felhasználók számára.

		// return layout.Dimensions{}
	}

	if m.backClickable.Clicked(gtx) {
		m.router.GoTo(routerModule.MinesweeperMenuPage)

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
		layout.Flexed(1, layout.Spacer{}.Layout),
		layout.Rigid(material.Button(styles.MaterialTheme, &m.backClickable, "Vissza").Layout),
		layout.Rigid(layout.Spacer{Height: 16}.Layout),
	)
}

func (*MinesweeperMultiplayerMenu) Restart() {}
