package minesweeper

import (
	"fmt"
	"giogo/ui/pages/minesweeper/model"
	"giogo/ui/styles"
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

var (
	shineColor  = color.NRGBA{A: 0xFF, R: 0xD6, G: 0xD6, B: 0xD6}
	baseColor   = color.NRGBA{A: 0xFF, R: 0xC3, G: 0xC3, B: 0xC3}
	borderColor = color.NRGBA{A: 0xFF, R: 0xBC, G: 0xBC, B: 0xBC}
	shadowColor = color.NRGBA{A: 0xFF, R: 0x80, G: 0x80, B: 0x80}
)

const (
	shadowThicknes = float32(4)
)

type MineButton struct {
	Size image.Point

	Value  int8
	Hidden bool
	Marked bool
}

func (mb *MineButton) Layout(gtx layout.Context, state model.MinesweeperState) layout.Dimensions {
	gtx.Constraints.Max = mb.Size

	mb.layout(&gtx, state)

	return layout.Dimensions{
		Size: gtx.Constraints.Max,
	}
}

func (mb *MineButton) layout(gtx *layout.Context, state model.MinesweeperState) {
	defer op.TransformOp{}.Push(gtx.Ops).Pop()
	defer clip.Rect{Max: mb.Size}.Push(gtx.Ops).Pop()

	paint.ColorOp{Color: baseColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	if state == model.WIN && mb.Value == -1 {
		drawMarkedCell(gtx, mb.Size)

		return
	}

	if mb.Hidden {
		if mb.Marked {
			drawMarkedCell(gtx, mb.Size)

			return
		}

		drawHiddenCell(gtx, mb.Size)

		return
	}

	paint.FillShape(gtx.Ops, borderColor, clip.Stroke{Path: clip.Rect{Max: mb.Size}.Path(), Width: 2}.Op())

	switch mb.Value {
	case 0:
		if !mb.Marked {
			return
		}

		drawMissMarkedCell(gtx, mb.Size)
	case -1:
		if mb.Marked {
			drawMarkedCell(gtx, mb.Size)

			return
		}

		drawMineCell(gtx, mb.Size)
	default:
		if mb.Marked {
			drawMissMarkedCell(gtx, mb.Size)

			return
		}

		fontSize := unit.Sp(mb.Size.Y / 2)
		offsetPoint := image.Pt(mb.Size.X/3, mb.Size.Y/5-1)

		textOffsetStack := op.Offset(offsetPoint).Push(gtx.Ops)
		textShadow := material.Label(styles.MaterialTheme, fontSize+2, fmt.Sprint(mb.Value))
		textShadow.Color = styles.TEXT_SHADOW
		textShadow.Font.Weight = font.Bold
		textShadow.Layout(*gtx)
		textOffsetStack.Pop()

		offsetPoint.X += 1
		offsetPoint.Y += 1
		textOffsetStack = op.Offset(offsetPoint).Push(gtx.Ops)
		textContent := material.Label(styles.MaterialTheme, fontSize, fmt.Sprint(mb.Value))
		textContent.Color = mb.getColorOfValue()
		textContent.Layout(*gtx)
		textOffsetStack.Pop()
	}
}

func (mb *MineButton) getColorOfValue() color.NRGBA {
	switch mb.Value {
	case 1:
		return styles.BLUE
	case 2:
		return styles.GREEN
	case 3:
		return styles.YELLOW
	case 4:
		return styles.ORANGE
	case 5:
		return styles.BLOOD_ORANGE
	case 6:
		return styles.RED
	default:
		return styles.NIGHT_BLACK
	}
}
