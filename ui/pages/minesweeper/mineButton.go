package minesweeper

import (
	"fmt"
	"giogo/ui/styles"
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

const (
	EXPLODE = iota
	FLIP
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

type MinesweeperState uint8

const (
	UNDEFINED MinesweeperState = iota
	START
	RUNNING
	LOSE
	WIN
	END
	LOADING
)

type MineAction struct {
	Value int8
	Pos   image.Point
	Type  pointer.Buttons
}

type MineButton struct {
	Size   image.Point
	Value  int8
	Hidden bool
	State  MinesweeperState
	Pos    image.Point
	Marked bool
	Parent *MineField

	pid       pointer.ID
	clickType pointer.Buttons
}

func (mb *MineButton) Layout(gtx layout.Context) layout.Dimensions {
	for _, event := range gtx.Events(mb) {
		switch event := event.(type) {
		case pointer.Event:
			switch event.Kind {
			case pointer.Press:
				mb.pressEvent(&event)
			case pointer.Release:
				mb.releaseEvent(&event, gtx.Ops)
			}
		}
	}

	gtx.Constraints.Max = mb.Size

	mb.layout(&gtx)

	return layout.Dimensions{
		Size: gtx.Constraints.Max,
	}
}

func (mb *MineButton) layout(gtx *layout.Context) {
	defer op.TransformOp{}.Push(gtx.Ops).Pop()
	defer clip.Rect{Max: mb.Size}.Push(gtx.Ops).Pop()

	paint.ColorOp{Color: baseColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	if mb.Parent.State == WIN && mb.Value == -1 {
		drawMarkedCell(gtx, unit.Dp(mb.Size.Y/8))

		return
	}

	if mb.Hidden {
		if mb.State == START || mb.State == RUNNING {
			mb.registerEvents(gtx.Ops)
		}

		if mb.Marked {
			drawMarkedCell(gtx, unit.Dp(mb.Size.Y/8))

			return
		}

		drawHiddenCell(gtx, mb.Size)
	} else {
		paint.FillShape(gtx.Ops, borderColor, clip.Stroke{Path: clip.Rect{Max: mb.Size}.Path(), Width: 2}.Op())

		switch mb.Value {
		// Üres mező -> semmi
		case 0:
			if mb.Marked {
				inset := unit.Dp(mb.Size.Y / 8)

				drawMissMarkedCell(gtx, inset)

				return
			}
		// Bomba
		case -1:
			inset := unit.Dp(mb.Size.Y / 8)

			if mb.Marked {
				drawMarkedCell(gtx, inset)

				return
			}

			drawMineCell(gtx, inset)
		// Számos mező
		default:
			if mb.Marked {
				drawMissMarkedCell(gtx, unit.Dp(mb.Size.Y/8))

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
}

func (mb *MineButton) registerEvents(ops *op.Ops) {
	pointer.InputOp{
		Tag:   mb,
		Kinds: pointer.Press | pointer.Release,
	}.Add(ops)
}

func (mb *MineButton) pressEvent(event *pointer.Event) {
	mb.pid = event.PointerID
	mb.clickType = event.Buttons
}

func (mb *MineButton) releaseEvent(event *pointer.Event, ops *op.Ops) {
	if mb.pid != event.PointerID {
		mb.pid = 0
		mb.clickType = 0

		return
	}

	if mb.Marked && mb.clickType != pointer.ButtonSecondary {
		return
	}

	switch mb.Parent.State {
	case START:
		fmt.Println("--- Game Start ---")
		mb.Parent.State = RUNNING

		// Legenerálni a bombákat
		mb.Parent.generateMines(mb.Pos)
		mb.Hidden = false
		mb.Parent.RevealedCells++

		if mb.Value == 0 {
			go mb.Parent.floodNeutralCells(mb.Pos)
		}
	case RUNNING:
		if mb.Hidden && mb.clickType == pointer.ButtonSecondary {
			mb.Marked = !mb.Marked

			if mb.Marked {
				mb.Parent.MarkedCells++
			} else {
				mb.Parent.MarkedCells--
			}

			break
		}

		mb.Parent.RevealedCells++

		switch mb.Value {
		case -1:
			mb.Parent.State = LOSE

			mb.Parent.revealFields()
		case 0:
			mb.Hidden = false

			go mb.Parent.floodNeutralCells(mb.Pos)
		default:
			mb.Hidden = false
		}
	}

	if mb.Parent.State == RUNNING && mb.Parent.RevealedCells == mb.Parent.countOfHorizontalCells*mb.Parent.countOfVerticalCells-mb.Parent.plantedMines {
		mb.Parent.State = WIN
		mb.Parent.MarkedCells = mb.Parent.plantedMines
		fmt.Println("--- GG, WIN ---")
	}

	op.InvalidateOp{}.Add(ops)
	mb.pid = 0
	mb.clickType = 0
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
