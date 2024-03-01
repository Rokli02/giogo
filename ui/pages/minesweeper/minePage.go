package minesweeper

import (
	"fmt"
	"giogo/ui"
	routerModule "giogo/ui/router"
	"giogo/ui/styles"
	"giogo/utils"
	"image"
	"time"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
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

const (
	heightOfHeader = 42
)

type MineField struct {
	BtnSize   image.Point
	BtnMatrix [][]MineButton

	engine              *MinesweeperEngine
	w                   *app.Window
	router              *routerModule.Router[ui.ApplicationCycles, string]
	mineChannel         chan *MineElement
	acks                chan uint8
	returnHomeClickable widget.Clickable
}

func NewMinefield(w *app.Window, router *routerModule.Router[ui.ApplicationCycles, string], horizontalCells, verticalCells, numberOfMines uint16) *MineField {
	mineField := &MineField{
		BtnSize: image.Pt(32, 32),

		engine:              NewMinesweeperEngine(),
		w:                   w,
		router:              router,
		returnHomeClickable: widget.Clickable{},
	}

	mineField.engine.Resize(horizontalCells, verticalCells, numberOfMines)

	mineField.Initialize()

	return mineField
}

/*
 * Public methods
 */

func (mf *MineField) Initialize() {
	sizeOfMinesweeper := image.Pt(mf.BtnSize.X*int(mf.engine.Width), mf.BtnSize.Y*int(mf.engine.Height)+heightOfHeader)
	mf.w.Option(func(_ unit.Metric, c *app.Config) {
		c.Title = "Minesweeper COPY"
		c.MinSize = sizeOfMinesweeper
		c.MaxSize = sizeOfMinesweeper
		c.Size = sizeOfMinesweeper
		c.Decorated = true
	})

	mf.BtnMatrix = make([][]MineButton, mf.engine.Height)
	mf.mineChannel = make(chan *MineElement, 4)
	mf.acks = make(chan uint8)

	mf.engine.mineChannel = mf.mineChannel
	mf.engine.acks = mf.acks

	for rowIndex := range mf.BtnMatrix {
		mf.BtnMatrix[rowIndex] = make([]MineButton, mf.engine.Width)

		for colIndex := range mf.BtnMatrix[rowIndex] {
			btn := &mf.BtnMatrix[rowIndex][colIndex]

			btn.onClick = mf.onButtonClick
			btn.Pos = image.Point{colIndex, rowIndex}
			btn.Size = mf.BtnSize
		}
	}

	go func() {
		for message := range mf.mineChannel {
			btn := &mf.BtnMatrix[message.Pos.Y][message.Pos.X]

			btn.Value = message.Value
			btn.Hidden = message.IsHidden()
			btn.Marked = message.IsMarked()

			mf.acks <- 1
		}
	}()

	mf.Restart()
}

func (mf *MineField) Restart() {
	mf.engine.Resize(mf.engine.Width, mf.engine.Height, mf.engine.MaxMines)

	for rowIndex := range mf.engine.matrix {
		for colIndex := range mf.engine.matrix[rowIndex] {
			element := mf.engine.matrix[rowIndex][colIndex]
			btn := &mf.BtnMatrix[rowIndex][colIndex]

			btn.Value = element.Value
			btn.Hidden = element.IsHidden()
			btn.Marked = element.IsMarked()
		}
	}
}

func (mf *MineField) Close() {
	close(mf.mineChannel)
	close(mf.acks)
}

func (mf *MineField) Layout(gtx layout.Context) layout.Dimensions {
	for _, event := range gtx.Events(mf) {
		switch evt := event.(type) {
		case key.Event:
			switch evt.Name {
			case "R":
				if evt.Modifiers == key.ModCtrl && evt.State == key.Release {
					mf.Restart()
				}
			}
		case pointer.Event:
			switch evt.Kind {
			case pointer.Press:
				mf.Restart()
			}
		}
	}

	defer utils.SetBackgroundColor(&gtx, styles.BACKGROUND_COLOR).Pop()

	if mf.engine.State == LOADING {
		op.InvalidateOp{At: time.Now().Add(mf.engine.animationDuration >> 1)}.Add(gtx.Ops)
	} else {
		key.InputOp{
			Tag:  mf,
			Keys: key.Set("Ctrl-R"),
		}.Add(gtx.Ops)
	}

	return layout.N.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		flexD := layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(mf.headerComponent),
			layout.Rigid(mf.bodyComponent),
		)

		flexD.Size.X = mf.BtnSize.X * int(mf.engine.Width)

		return flexD
	})
}

/*
 * Private methods
 */

func (mf *MineField) headerComponent(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min = image.Pt(0, 0)
	gtx.Constraints.Max = image.Pt(int(mf.engine.Width)*mf.BtnSize.X, heightOfHeader)

	defer utils.SetBackgroundColor(&gtx, styles.HEADER_BACKGROUND).Pop()

	headerD := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
		layout.Flexed(3, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = image.Point{gtx.Constraints.Max.Y - 4, gtx.Constraints.Max.Y - 4}
			icon, _ := widget.NewIcon(icons.ActionHome)
			iconColor := styles.BLOOD_ORANGE

			if mf.returnHomeClickable.Clicked(gtx) {
				mf.router.GoTo(routerModule.MenuPage)
			}

			if mf.returnHomeClickable.Pressed() {
				iconColor.R -= 10
				iconColor.G -= 10
				iconColor.B -= 10
			} else if mf.returnHomeClickable.Hovered() {
				iconColor.R += 20
				iconColor.G += 20
				iconColor.B += 20
			}

			mf.returnHomeClickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				pointer.Cursor.Add(pointer.CursorPointer, gtx.Ops)

				return icon.Layout(gtx, iconColor)
			})

			return layout.Dimensions{Size: gtx.Constraints.Max}
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(4).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min = image.Pt(0, 0)
				gtx.Constraints.Max = image.Pt(heightOfHeader-8, heightOfHeader-8)

				paint.FillShape(gtx.Ops, styles.NIGHT_BLACK, clip.Stroke{Width: 2, Path: clip.Ellipse{Min: gtx.Constraints.Min, Max: gtx.Constraints.Max}.Path(gtx.Ops)}.Op())
				defer clip.Ellipse{Min: gtx.Constraints.Min, Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()

				switch mf.engine.State {
				case START:
					fallthrough
				case END:
					fallthrough
				case LOADING:
					fallthrough
				case RUNNING:
					// Sárga neutrális
					paint.ColorOp{Color: styles.YELLOW}.Add(gtx.Ops)
					paint.PaintOp{}.Add(gtx.Ops)

					drawNeutralFace(&gtx)
				case WIN:
					// Zöld mosolygó
					paint.ColorOp{Color: styles.GREEN}.Add(gtx.Ops)
					paint.PaintOp{}.Add(gtx.Ops)

					drawHappyFace(&gtx)
				case LOSE:
					// Piros mérges/szomorú
					paint.ColorOp{Color: styles.RED}.Add(gtx.Ops)
					paint.PaintOp{}.Add(gtx.Ops)

					drawMadFace(&gtx)
				}

				pointer.InputOp{
					Tag:   mf,
					Kinds: pointer.Press,
				}.Add(gtx.Ops)

				return layout.Dimensions{Size: gtx.Constraints.Max}
			})
		}),
		layout.Flexed(3, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = image.Point{0, 0}

			width := min(74, gtx.Constraints.Max.X-gtx.Constraints.Max.X>>2)
			defer op.Offset(image.Pt(gtx.Constraints.Max.X-width, 0)).Push(gtx.Ops).Pop()
			gtx.Constraints.Max.X = width

			return layout.UniformInset(4).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				clr := styles.YELLOWISH_GREEN
				clr.A = 0x6F

				paint.FillShape(gtx.Ops, styles.NIGHT_BLACK, clip.Stroke{Width: 2, Path: clip.Rect{Min: gtx.Constraints.Min, Max: gtx.Constraints.Max}.Path()}.Op())
				defer utils.SetBackgroundColor(&gtx, clr).Pop()

				macro := op.Record(gtx.Ops)
				label := material.Label(styles.MaterialTheme, unit.Sp(16), fmt.Sprintf("%v/%v", mf.engine.Marked, mf.engine.mines))
				label.MaxLines = 1
				label.Font.Weight = font.Medium
				label.Alignment = text.Middle
				labelDim := label.Layout(gtx)
				labelMacro := macro.Stop()

				labelOffset := image.Point{X: (gtx.Constraints.Max.X - labelDim.Size.X) >> 1, Y: (gtx.Constraints.Max.Y - labelDim.Size.Y) >> 1}
				defer op.Offset(labelOffset).Push(gtx.Ops).Pop()

				labelMacro.Add(gtx.Ops)

				return layout.Dimensions{Size: gtx.Constraints.Max}
			})
		}),
	)

	headerD.Size = gtx.Constraints.Max

	return headerD
}

func (mf *MineField) bodyComponent(gtx layout.Context) layout.Dimensions {
	rowList := &layout.List{Axis: layout.Vertical, Alignment: layout.Start}

	return rowList.Layout(gtx, len(mf.BtnMatrix), func(gtx layout.Context, rowIndex int) layout.Dimensions {
		list := &layout.List{Axis: layout.Horizontal, Alignment: layout.Start}

		return list.Layout(gtx, len(mf.BtnMatrix[rowIndex]), func(gtx layout.Context, colIndex int) layout.Dimensions {
			btnDimension := mf.BtnMatrix[rowIndex][colIndex].Layout(gtx, mf.engine.State)

			return btnDimension
		})
	})
}

func (mf *MineField) onButtonClick(pos image.Point, clickType pointer.Buttons) {
	state, element := mf.engine.OnButtonClick(pos, clickType)

	switch state {
	case RUNNING:
		mf.mineChannel <- element
		<-mf.acks
	case LOSE:
		fallthrough
	case WIN:
		for rowIndex := range mf.engine.matrix {
			for colIndex := range mf.engine.matrix[rowIndex] {
				btn := &mf.BtnMatrix[rowIndex][colIndex]
				element := mf.engine.matrix[rowIndex][colIndex]
				element.PropOff(hiddenBits)

				btn.Value = element.Value
				btn.Hidden = element.IsHidden()
				btn.Marked = element.IsMarked()
			}
		}
	}
}
