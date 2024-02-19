package utils

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

func SetBackgroundColor(gtx *layout.Context, color color.NRGBA) clip.Stack {
	var stack = clip.Rect{Min: gtx.Constraints.Min, Max: gtx.Constraints.Max}.Push(gtx.Ops)
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return stack
}
