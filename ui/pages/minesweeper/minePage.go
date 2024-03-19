package minesweeper

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
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

	"giogo/ui"
	"giogo/ui/pages/minesweeper/engine"
	"giogo/ui/pages/minesweeper/model"
	routerModule "giogo/ui/router"
	"giogo/ui/styles"
	"giogo/utils"
)

const (
	heightOfHeader       = 32
	game_end_txt_padding = 4
)

var (
	game_end_cover_shadow = color.NRGBA{A: 0xCC}
	game_end_txt_highligh = color.NRGBA{A: 0x69, R: 0xbf, G: 0xbf, B: 0xbf}
)

type MineField struct {
	BtnSize   image.Point
	BtnMatrix [][]MineButton

	engine              engine.MinesweeperEngine
	w                   *app.Window
	router              *routerModule.Router[ui.ApplicationCycles, string]
	mineChannel         chan model.MineElement
	acks                chan uint8
	engineCommand       chan engine.EngineCommand
	returnHomeClickable widget.Clickable
	refreshRate         time.Duration
}

func NewMinefield(w *app.Window, router *routerModule.Router[ui.ApplicationCycles, string], refreshRate time.Duration) *MineField {
	mineField := &MineField{
		BtnSize: image.Pt(24, 24),

		engine:              nil,
		w:                   w,
		router:              router,
		refreshRate:         refreshRate,
		returnHomeClickable: widget.Clickable{},
	}

	return mineField
}

/*
 * Public methods
 */

func (mf *MineField) SetEngine(minesweeperEngine engine.MinesweeperEngine) *MineField {
	mf.engine = minesweeperEngine

	return mf
}

func (mf *MineField) Initialize() {
	if mf.engine == nil {
		panic("Engine is not set for Minesweeper")
	}

	mf.engine.Initialize()

	mf.w.Option(func(_ unit.Metric, c *app.Config) {
		c.Title = "Minesweeper COPY"
		c.Decorated = true
	})

	mf.mineChannel = make(chan model.MineElement, 4)
	mf.acks = make(chan uint8)

	mf.engineCommand = make(chan engine.EngineCommand)

	go func() {
		for command := range mf.engineCommand {
			fmt.Printf("MinePage received command from engine (%s)\n", command.ToString())

			switch command {
			case engine.RESIZE:
				mf.ResizeGui()
				mf.w.Invalidate()
			case engine.RESTART:
				mf.RestartGui()
				mf.w.Invalidate()
			case engine.GO_BACK:
				mf.router.GoBackTo(routerModule.MinesweeperMenuPage)
			case engine.AFTER_CLICK_LOSE:
				fallthrough
			case engine.AFTER_CLICK_WIN:
				remainingMines := mf.engine.GetRemainingMines()

				for _, mine := range remainingMines {
					btn := &mf.BtnMatrix[mine.Pos.Y][mine.Pos.X]

					btn.Value = mine.Value
					btn.Hidden = mine.IsHidden()
					btn.Marked = mine.IsMarked()
				}

				mf.w.Invalidate()
			}

			fmt.Printf("(%s) State after command: (width=%d) | (height=%d) | (marked=%d) | (mines=%d)\n",
				mf.engine.GetState().ToString(), mf.engine.GetWidth(), mf.engine.GetHeight(),
				mf.engine.GetMarked(), mf.engine.GetMines(),
			)
		}
	}()

	go func() {
		for {
			message, isOpen := <-mf.mineChannel

			if !isOpen {
				close(mf.acks)
				mf.engine.Close()
				mf.acks = nil

				return
			}

			btn := &mf.BtnMatrix[message.Pos.Y][message.Pos.X]

			btn.Value = message.Value
			btn.Hidden = message.IsHidden()
			btn.Marked = message.IsMarked()

			mf.acks <- 1
		}
	}()

	mf.engine.SetChannels(mf.mineChannel, mf.acks, mf.engineCommand)
	mf.ResizeGui()
}

func (mf *MineField) Restart() {
	mf.engine.Restart()

	if _, isLocalEngine := mf.engine.(*engine.MinesweeperLocalEngine); isLocalEngine {
		mf.RestartGui()
	}
}

func (mf *MineField) RestartGui() {
	for rowIndex := range mf.BtnMatrix {
		for colIndex := range mf.BtnMatrix[rowIndex] {
			btn := &mf.BtnMatrix[rowIndex][colIndex]

			btn.Value = 0
			btn.Hidden = true
			btn.Marked = false
		}
	}
}

func (mf *MineField) ResizeGui() {
	mf.BtnMatrix = make([][]MineButton, mf.engine.GetHeight())

	for rowIndex := range mf.BtnMatrix {
		mf.BtnMatrix[rowIndex] = make([]MineButton, mf.engine.GetWidth())

		for colIndex := range mf.BtnMatrix[rowIndex] {
			btn := &mf.BtnMatrix[rowIndex][colIndex]

			btn.Value = 0
			btn.Hidden = true
			btn.Marked = false
			btn.onClick = mf.engine.OnButtonClick
			btn.Pos = image.Point{colIndex, rowIndex}
			btn.Size = mf.BtnSize
		}
	}

	sizeOfMinesweeper := image.Pt(mf.BtnSize.X*mf.engine.GetWidth(), mf.BtnSize.Y*mf.engine.GetHeight()+heightOfHeader)
	mf.w.Option(func(_ unit.Metric, c *app.Config) {
		c.MinSize = sizeOfMinesweeper
		c.MaxSize = sizeOfMinesweeper
		c.Size = sizeOfMinesweeper
	})
}

func (mf *MineField) Close() {
	close(mf.mineChannel)
	mf.engine.Close()
	mf.mineChannel = nil
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

	if mf.engine.GetState() == model.LOADING {
		op.InvalidateOp{At: time.Now().Add(mf.refreshRate)}.Add(gtx.Ops)
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

		flexD.Size.X = mf.BtnSize.X * mf.engine.GetWidth()

		return flexD
	})
}

/*
 * Private methods
 */

func (mf *MineField) headerComponent(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min = image.Pt(0, 0)
	gtx.Constraints.Max = image.Pt(mf.engine.GetWidth()*mf.BtnSize.X, heightOfHeader)

	defer utils.SetBackgroundColor(&gtx, styles.HEADER_BACKGROUND).Pop()

	headerD := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
		layout.Flexed(3, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = image.Point{gtx.Constraints.Max.Y - 4, gtx.Constraints.Max.Y - 4}
			icon, _ := widget.NewIcon(icons.ActionHome)
			iconColor := styles.BLOOD_ORANGE

			if mf.returnHomeClickable.Clicked(gtx) {
				mf.router.GoBackTo(routerModule.MinesweeperMenuPage)
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

				switch mf.engine.GetState() {
				case model.START:
					fallthrough
				case model.END:
					fallthrough
				case model.LOADING:
					fallthrough
				case model.RUNNING:
					// Sárga neutrális
					paint.ColorOp{Color: styles.YELLOW}.Add(gtx.Ops)
					paint.PaintOp{}.Add(gtx.Ops)

					drawNeutralFace(&gtx)
				case model.WIN:
					// Zöld mosolygó
					paint.ColorOp{Color: styles.GREEN}.Add(gtx.Ops)
					paint.PaintOp{}.Add(gtx.Ops)

					drawHappyFace(&gtx)
				case model.LOSE:
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
				label := material.Label(styles.MaterialTheme, unit.Sp(16), fmt.Sprintf("%v/%v", mf.engine.GetMarked(), mf.engine.GetMines()))
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
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			rowList := &layout.List{Axis: layout.Vertical, Alignment: layout.Start}

			return rowList.Layout(gtx, len(mf.BtnMatrix), func(gtx layout.Context, rowIndex int) layout.Dimensions {
				list := &layout.List{Axis: layout.Horizontal, Alignment: layout.Start}

				return list.Layout(gtx, len(mf.BtnMatrix[rowIndex]), func(gtx layout.Context, colIndex int) layout.Dimensions {
					btnDimension := mf.BtnMatrix[rowIndex][colIndex].Layout(gtx, mf.engine.GetState())

					return btnDimension
				})
			})
		}),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			isGameEnded := false
			var txt string

			switch mf.engine.GetState() {
			case model.WIN:
				isGameEnded = true
				txt = "GG EZ"
			case model.LOSE:
				isGameEnded = true
				txt = "Na majd holnap"
			}

			if isGameEnded {
				defer clip.Rect{Max: gtx.Constraints.Min}.Push(gtx.Ops).Pop()
				paint.Fill(gtx.Ops, game_end_cover_shadow)
				tempGtxConstraints := gtx.Constraints
				gtx.Constraints.Min.Y = 0

				// Szöveg
				rec := op.Record(gtx.Ops)
				lbl := material.Label(styles.MaterialTheme, unit.Sp(20), txt)
				lbl.Alignment = text.Middle
				lbl.Font.Weight = font.Medium
				lblDim := lbl.Layout(gtx)
				macro := rec.Stop()

				op.Offset(image.Point{
					0,
					((gtx.Constraints.Max.Y - lblDim.Size.Y) >> 1) - game_end_txt_padding - heightOfHeader,
				}).Add(gtx.Ops)

				txtCoverHighlighter(&gtx, &lblDim)

				op.Offset(image.Point{
					0,
					game_end_txt_padding,
				}).Add(gtx.Ops)

				macro.Add(gtx.Ops)

				gtx.Constraints = tempGtxConstraints
				return layout.Dimensions{Size: gtx.Constraints.Min}
			}

			return layout.Dimensions{}
		}),
	)
}

func txtCoverHighlighter(gtx *layout.Context, lblDim *layout.Dimensions) {
	// LinearGradient the first part of text cover
	highlightCover := clip.Rect{Max: image.Point{gtx.Constraints.Max.X >> 1, lblDim.Size.Y + game_end_txt_padding<<1}}

	hcPop := highlightCover.Push(gtx.Ops)
	middleOfHighlightY := float32(game_end_txt_padding + lblDim.Size.Y>>1)
	paint.LinearGradientOp{
		Stop1:  f32.Pt(0, middleOfHighlightY),
		Stop2:  f32.Pt(float32(gtx.Constraints.Max.X>>2), middleOfHighlightY),
		Color1: color.NRGBA{},
		Color2: game_end_txt_highligh,
	}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	hcPop.Pop()

	// LinearGradient the second part of text cover
	highlightOffset := op.Offset(image.Point{gtx.Constraints.Max.X >> 1, 0}).Push(gtx.Ops)
	hcPop = highlightCover.Push(gtx.Ops)

	paint.LinearGradientOp{
		Stop1:  f32.Pt(0, middleOfHighlightY),
		Stop2:  f32.Pt(float32(gtx.Constraints.Max.X-gtx.Constraints.Max.X>>2), middleOfHighlightY),
		Color1: game_end_txt_highligh,
		Color2: color.NRGBA{},
	}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	hcPop.Pop()
	highlightOffset.Pop()
}
