package component

import (
	"giogo/ui/styles"
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	textModule "gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Input struct {
	MaxWidth        int
	Editor          material.EditorStyle
	label           material.LabelStyle
	Text            string
	BackgroundColor color.NRGBA
}

func NewInput(text string, backgroundColor color.NRGBA, maxWidth int) *Input {
	we := &widget.Editor{Alignment: textModule.Middle, SingleLine: true, ReadOnly: false}
	i := &Input{
		Editor:          material.Editor(styles.MaterialTheme, we, ""),
		label:           material.Label(styles.MaterialTheme, unit.Sp(16), text),
		Text:            text,
		BackgroundColor: backgroundColor,
		MaxWidth:        maxWidth,
	}

	i.Editor.Editor.Alignment = textModule.Middle
	i.Editor.Editor.SingleLine = true

	return i
}

func (i *Input) Layout(gtx layout.Context) layout.Dimensions {
	if i.MaxWidth != 0 {
		gtx.Constraints.Max.X = i.MaxWidth
	}

	macroLabel := op.Record(gtx.Ops)
	labelDim := i.label.Layout(gtx)
	recordLabelCallOp := macroLabel.Stop()

	tempGtxConstr := gtx.Constraints
	gtx.Constraints.Max.X -= labelDim.Size.X + 4
	gtx.Constraints.Min.X = gtx.Constraints.Max.X

	macro := op.Record(gtx.Ops)
	editorDim := i.Editor.Layout(gtx)
	recordCallOp := macro.Stop()

	gtx.Constraints = tempGtxConstr
	editorDim.Size.X += 12
	editorDim.Size.Y += 8

	gtx.Constraints.Max.Y = max(editorDim.Size.Y, labelDim.Size.Y)
	labelOffset := op.Offset(image.Point{0, (gtx.Constraints.Max.Y - labelDim.Size.Y) >> 1}).Push(gtx.Ops)
	recordLabelCallOp.Add(gtx.Ops)
	labelOffset.Pop()

	op.Offset(image.Point{labelDim.Size.X + 4, 0}).Add(gtx.Ops)
	paint.FillShape(gtx.Ops, i.BackgroundColor, clip.Rect{Max: editorDim.Size}.Op())
	op.Offset(image.Point{6, 4}).Add(gtx.Ops)

	recordCallOp.Add(gtx.Ops)

	return layout.Dimensions{Size: gtx.Constraints.Max}
}
