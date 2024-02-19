package component

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"
)

const (
	checkboxDotOffset = 3
)

type Checkbox struct {
	Click     *widget.Clickable
	Size      int
	IsChecked bool
}

func (cb Checkbox) Layout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min = image.Point{}
	gtx.Constraints.Max = image.Point{cb.Size, cb.Size}

	return cb.Click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		rrectPath := clip.RRect{SE: 4, SW: 4, NW: 4, NE: 4, Rect: image.Rectangle{Max: gtx.Constraints.Max}}

		paint.FillShape(gtx.Ops, color.NRGBA{A: 0x66}, clip.Outline{
			Path: rrectPath.Path(gtx.Ops),
		}.Op())

		paint.FillShape(gtx.Ops, color.NRGBA{A: 0xFF}, clip.Stroke{
			Width: 2,
			Path:  rrectPath.Path(gtx.Ops),
		}.Op())

		if cb.IsChecked {
			tempGtxConstraints := gtx.Constraints
			circleOffset := op.Offset(image.Pt(checkboxDotOffset, checkboxDotOffset)).Push(gtx.Ops)

			gtx.Constraints.Max.X -= (checkboxDotOffset << 1)
			gtx.Constraints.Max.Y -= (checkboxDotOffset << 1)
			rrectPath.Rect.Max = gtx.Constraints.Max
			rrectPath.NE = 2
			rrectPath.NW = 2
			rrectPath.SE = 2
			rrectPath.SW = 2

			paint.FillShape(gtx.Ops, color.NRGBA{A: 0x90}, clip.Outline{
				Path: rrectPath.Path(gtx.Ops),
			}.Op())

			circleOffset.Pop()

			gtx.Constraints = tempGtxConstraints
		}

		return layout.Dimensions{Size: gtx.Constraints.Max}
	})
}
