package minesweeper

import (
	"giogo/assets"
	"giogo/ui/styles"
	"giogo/utils"
	"image"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

var markedCellOps *op.Ops
var markedCellOpsCall op.CallOp

func drawMarkedCell(gtx *layout.Context, size image.Point) {
	if markedCellOps != nil {
		markedCellOpsCall.Add(gtx.Ops)

		return
	}

	markedCellOps = new(op.Ops)
	localGtx := *gtx
	localGtx.Ops = markedCellOps

	macro := op.Record(markedCellOps)
	drawHiddenCell(&localGtx, size)

	tmpGtx := localGtx.Constraints
	localGtx.Constraints.Max.X = localGtx.Constraints.Max.X - (int(shadowThicknes) << 1)
	localGtx.Constraints.Max.Y = localGtx.Constraints.Max.Y - (int(shadowThicknes) << 1)

	offset := op.Offset(image.Point{X: int(shadowThicknes), Y: int(shadowThicknes)}).Push(markedCellOps)

	markedField := assets.GetWidgetImage("marked", size.X)
	markedField.Position = layout.Center
	markedField.Layout(localGtx)

	localGtx.Constraints = tmpGtx

	offset.Pop()

	markedCellOpsCall = macro.Stop()
	markedCellOpsCall.Add(gtx.Ops)
}

func drawMineCell(gtx *layout.Context, size image.Point) {
	mineField := assets.GetWidgetImage("bomb", size.X*2)
	mineField.Position = layout.Center
	mineField.Layout(*gtx)
}

// ??? mi ez a név ???
func drawMissMarkedCell(gtx *layout.Context, size image.Point) {
	drawMarkedCell(gtx, size)
	drawBigRedX(gtx, 2)
}

var hiddenCellOps *op.Ops
var hiddenCellOpsCall op.CallOp

func drawHiddenCell(gtx *layout.Context, size image.Point) {
	if hiddenCellOps != nil {
		hiddenCellOpsCall.Add(gtx.Ops)

		return
	}

	hiddenCellOps = new(op.Ops)

	macro := op.Record(hiddenCellOps)

	shine := clip.Path{}
	shine.Begin(hiddenCellOps)
	shine.MoveTo(f32.Pt(float32(0), float32(size.Y)))
	shine.LineTo(f32.Pt(0, 0))
	shine.LineTo(f32.Pt(float32(size.X), 0))
	shine.LineTo(f32.Pt(float32(size.X)-shadowThicknes, shadowThicknes))
	shine.LineTo(f32.Pt(shadowThicknes, shadowThicknes))
	shine.LineTo(f32.Pt(shadowThicknes, float32(size.Y)-shadowThicknes))
	shine.Close()

	paint.FillShape(hiddenCellOps, shineColor, clip.Outline{Path: shine.End()}.Op())

	shadow := clip.Path{}
	shadow.Begin(hiddenCellOps)
	shadow.MoveTo(f32.Pt(float32(0), float32(size.Y)))
	shadow.LineTo(f32.Pt(float32(size.X), float32(size.Y)))
	shadow.LineTo(f32.Pt(float32(size.X), 0))
	shadow.LineTo(f32.Pt(float32(size.X)-shadowThicknes, shadowThicknes))
	shadow.LineTo(f32.Pt(float32(size.X)-shadowThicknes, float32(size.Y)-shadowThicknes))
	shadow.LineTo(f32.Pt(shadowThicknes, float32(size.Y)-shadowThicknes))
	shadow.Close()

	paint.FillShape(hiddenCellOps, shadowColor, clip.Outline{Path: shadow.End()}.Op())

	hiddenCellOpsCall = macro.Stop()
	hiddenCellOpsCall.Add(gtx.Ops)
}

func drawBigRedX(gtx *layout.Context, weigth int) layout.Dimensions {
	bigRedX := clip.Path{}
	middlePoint := image.Point{gtx.Constraints.Max.X >> 1, gtx.Constraints.Max.Y >> 1}

	bigRedX.Begin(gtx.Ops)
	bigRedX.MoveTo(f32.Pt(0, 0))
	bigRedX.LineTo(f32.Pt(0, float32(weigth)))
	bigRedX.LineTo(f32.Pt(float32(middlePoint.X-weigth), float32(middlePoint.Y)))
	bigRedX.LineTo(f32.Pt(0, float32(gtx.Constraints.Max.Y-weigth)))
	bigRedX.LineTo(f32.Pt(0, float32(gtx.Constraints.Max.Y)))
	bigRedX.LineTo(f32.Pt(float32(weigth), float32(gtx.Constraints.Max.Y)))
	bigRedX.LineTo(f32.Pt(float32(middlePoint.X), float32(middlePoint.Y+weigth)))
	bigRedX.LineTo(f32.Pt(float32(gtx.Constraints.Max.X-weigth), float32(gtx.Constraints.Max.Y)))
	bigRedX.LineTo(f32.Pt(float32(gtx.Constraints.Max.X), float32(gtx.Constraints.Max.Y)))
	bigRedX.LineTo(f32.Pt(float32(gtx.Constraints.Max.X), float32(gtx.Constraints.Max.Y-weigth)))
	bigRedX.LineTo(f32.Pt(float32(middlePoint.X+weigth), float32(middlePoint.Y)))
	bigRedX.LineTo(f32.Pt(float32(gtx.Constraints.Max.X), float32(weigth)))
	bigRedX.LineTo(f32.Pt(float32(gtx.Constraints.Max.X), float32(0)))
	bigRedX.LineTo(f32.Pt(float32(gtx.Constraints.Max.X-weigth), float32(0)))
	bigRedX.LineTo(f32.Pt(float32(middlePoint.X), float32(middlePoint.Y-weigth)))
	bigRedX.LineTo(f32.Pt(float32(weigth), 0))
	bigRedX.Close()

	paint.FillShape(gtx.Ops, styles.RED, clip.Outline{Path: bigRedX.End()}.Op())

	return layout.Dimensions{Size: utils.GetSize(gtx)}
}

var neutralFaceCache *op.Ops
var neutralFaceCallOp op.CallOp

func drawNeutralFace(gtx *layout.Context) {
	if neutralFaceCache != nil {
		neutralFaceCallOp.Add(gtx.Ops)

		return
	}

	neutralFaceCache = new(op.Ops)
	macro := op.Record(neutralFaceCache)

	middlePoint := image.Point{X: gtx.Constraints.Max.X >> 1, Y: gtx.Constraints.Max.Y >> 1}
	eyeOffset := image.Point{X: middlePoint.X >> 1, Y: middlePoint.Y / 3}
	eyeSize := 2
	mouthOffset := image.Point{X: middlePoint.X / 3, Y: middlePoint.Y / 3}
	mouthSize := 1

	// Left eye
	paint.FillShape(
		neutralFaceCache,
		styles.NIGHT_BLACK,
		clip.Ellipse{
			Min: image.Point{
				X: middlePoint.X - eyeOffset.X - eyeSize,
				Y: middlePoint.Y - eyeOffset.Y - eyeSize,
			},
			Max: image.Point{
				X: middlePoint.X - eyeOffset.X + eyeSize,
				Y: middlePoint.Y - eyeOffset.Y + eyeSize,
			},
		}.Op(neutralFaceCache),
	)

	// Right eye
	paint.FillShape(
		neutralFaceCache,
		styles.NIGHT_BLACK,
		clip.Ellipse{
			Min: image.Point{
				X: middlePoint.X + eyeOffset.X - eyeSize,
				Y: middlePoint.Y - eyeOffset.Y - eyeSize,
			},
			Max: image.Point{
				X: middlePoint.X + eyeOffset.X + eyeSize,
				Y: middlePoint.Y - eyeOffset.Y + eyeSize,
			},
		}.Op(neutralFaceCache),
	)

	// Mouth
	paint.FillShape(
		neutralFaceCache,
		styles.NIGHT_BLACK,
		clip.Rect{
			Min: image.Point{
				X: middlePoint.X - mouthOffset.X,
				Y: middlePoint.Y + mouthOffset.Y - mouthSize,
			},
			Max: image.Point{
				X: middlePoint.X + mouthOffset.X,
				Y: middlePoint.Y + mouthOffset.Y + mouthSize,
			},
		}.Op(),
	)

	neutralFaceCallOp = macro.Stop()
	neutralFaceCallOp.Add(gtx.Ops)
}

var madFaceCache *op.Ops
var madFaceCallOp op.CallOp

func drawMadFace(gtx *layout.Context) {
	if madFaceCache != nil {
		madFaceCallOp.Add(gtx.Ops)

		return
	}

	madFaceCache = new(op.Ops)
	macro := op.Record(madFaceCache)

	middlePoint := image.Point{X: gtx.Constraints.Max.X >> 1, Y: gtx.Constraints.Max.Y >> 1}
	eyeOffset := image.Point{X: middlePoint.X >> 1, Y: middlePoint.Y >> 2}
	mouthOffset := image.Point{X: middlePoint.X >> 1, Y: middlePoint.Y >> 1}
	const (
		eyeSize       = 2
		eyebrowSize   = 3
		eyebrowWidth  = 10
		eyebrowOffset = 4
		mouthSize     = 2
		mouthCurve    = 4
	)

	// Left eye
	paint.FillShape(
		madFaceCache,
		styles.NIGHT_BLACK,
		clip.Ellipse{
			Min: image.Point{
				X: middlePoint.X - eyeOffset.X - eyeSize,
				Y: middlePoint.Y - eyeOffset.Y - eyeSize,
			},
			Max: image.Point{
				X: middlePoint.X - eyeOffset.X + eyeSize,
				Y: middlePoint.Y - eyeOffset.Y + eyeSize,
			},
		}.Op(madFaceCache),
	)

	// Left Eyebrow
	leftEyebrowPath := clip.Path{}
	leftEyebrowPath.Begin(madFaceCache)
	leftEyebrowPath.MoveTo(f32.Pt(float32(middlePoint.X-eyeOffset.X-(eyebrowWidth>>1)), float32(middlePoint.Y-eyeOffset.Y-eyebrowOffset)))
	leftEyebrowPath.Line(f32.Pt(eyebrowWidth, eyebrowOffset))
	leftEyebrowPath.Line(f32.Pt(0, -eyebrowSize))
	leftEyebrowPath.Line(f32.Pt(-eyebrowWidth, -eyebrowOffset))
	leftEyebrowPath.Close()

	paint.FillShape(
		madFaceCache,
		styles.NIGHT_BLACK,
		clip.Outline{Path: leftEyebrowPath.End()}.Op(),
	)

	// Right eye
	paint.FillShape(
		madFaceCache,
		styles.NIGHT_BLACK,
		clip.Ellipse{
			Min: image.Point{
				X: middlePoint.X + eyeOffset.X - eyeSize,
				Y: middlePoint.Y - eyeOffset.Y - eyeSize,
			},
			Max: image.Point{
				X: middlePoint.X + eyeOffset.X + eyeSize,
				Y: middlePoint.Y - eyeOffset.Y + eyeSize,
			},
		}.Op(madFaceCache),
	)

	// Right eyebrow
	rightEyebrowPath := clip.Path{}
	rightEyebrowPath.Begin(madFaceCache)
	rightEyebrowPath.MoveTo(f32.Pt(float32(middlePoint.X+eyeOffset.X+(eyebrowWidth>>1)), float32(middlePoint.Y-eyeOffset.Y-eyebrowOffset)))
	rightEyebrowPath.Line(f32.Pt(-eyebrowWidth, eyebrowOffset))
	rightEyebrowPath.Line(f32.Pt(0, -eyebrowSize))
	rightEyebrowPath.Line(f32.Pt(eyebrowWidth, -eyebrowOffset))
	rightEyebrowPath.Close()

	paint.FillShape(
		madFaceCache,
		styles.NIGHT_BLACK,
		clip.Outline{Path: rightEyebrowPath.End()}.Op(),
	)

	// Mouth
	mouthPath := clip.Path{}
	mouthPath.Begin(madFaceCache)
	mouthPath.MoveTo(f32.Pt(float32(middlePoint.X-mouthOffset.X), float32(middlePoint.Y+mouthOffset.Y-mouthSize)))
	mouthPath.QuadTo(f32.Pt(float32(middlePoint.X), float32(middlePoint.Y+mouthOffset.Y-mouthCurve-mouthSize)), f32.Pt(float32(middlePoint.X+mouthOffset.X), float32(middlePoint.Y+mouthOffset.Y-mouthSize)))
	mouthPath.Line(f32.Pt(0, float32(mouthSize)))
	mouthPath.QuadTo(f32.Pt(float32(middlePoint.X), float32(middlePoint.Y+mouthOffset.Y-mouthCurve)), f32.Pt(float32(middlePoint.X-mouthOffset.X), float32(middlePoint.Y+mouthOffset.Y)))
	mouthPath.Close()
	paint.FillShape(
		madFaceCache,
		styles.NIGHT_BLACK,
		clip.Outline{Path: mouthPath.End()}.Op(),
	)

	madFaceCallOp = macro.Stop()
	madFaceCallOp.Add(gtx.Ops)
}

var happyFaceCache *op.Ops
var happyFaceCallOp op.CallOp

func drawHappyFace(gtx *layout.Context) {
	if happyFaceCache != nil {
		happyFaceCallOp.Add(gtx.Ops)

		return
	}

	happyFaceCache = new(op.Ops)
	macro := op.Record(happyFaceCache)

	middlePoint := image.Point{X: gtx.Constraints.Max.X >> 1, Y: gtx.Constraints.Max.Y >> 1}
	eyeOffset := image.Point{X: middlePoint.X >> 1, Y: middlePoint.Y >> 2}
	mouthOffset := image.Point{X: middlePoint.X >> 1, Y: middlePoint.Y >> 1}
	const (
		eyeSize    = 2
		mouthSize  = 2
		mouthCurve = 4
	)

	// Left eye
	paint.FillShape(
		happyFaceCache,
		styles.NIGHT_BLACK,
		clip.Ellipse{
			Min: image.Point{
				X: middlePoint.X - eyeOffset.X - eyeSize,
				Y: middlePoint.Y - eyeOffset.Y - eyeSize,
			},
			Max: image.Point{
				X: middlePoint.X - eyeOffset.X + eyeSize,
				Y: middlePoint.Y - eyeOffset.Y + eyeSize,
			},
		}.Op(happyFaceCache),
	)

	// Right eye
	paint.FillShape(
		happyFaceCache,
		styles.NIGHT_BLACK,
		clip.Ellipse{
			Min: image.Point{
				X: middlePoint.X + eyeOffset.X - eyeSize,
				Y: middlePoint.Y - eyeOffset.Y - eyeSize,
			},
			Max: image.Point{
				X: middlePoint.X + eyeOffset.X + eyeSize,
				Y: middlePoint.Y - eyeOffset.Y + eyeSize,
			},
		}.Op(happyFaceCache),
	)

	// Mouth
	mouthPath := clip.Path{}
	mouthPath.Begin(happyFaceCache)
	mouthPath.MoveTo(f32.Pt(float32(middlePoint.X-mouthOffset.X), float32(middlePoint.Y+mouthOffset.Y-mouthSize)))
	mouthPath.QuadTo(f32.Pt(float32(middlePoint.X), float32(middlePoint.Y+mouthOffset.Y+mouthCurve-mouthSize)), f32.Pt(float32(middlePoint.X+mouthOffset.X), float32(middlePoint.Y+mouthOffset.Y-mouthSize)))
	mouthPath.Line(f32.Pt(0, float32(mouthSize)))
	mouthPath.QuadTo(f32.Pt(float32(middlePoint.X), float32(middlePoint.Y+mouthOffset.Y+mouthCurve)), f32.Pt(float32(middlePoint.X-mouthOffset.X), float32(middlePoint.Y+mouthOffset.Y)))
	mouthPath.Close()
	paint.FillShape(
		happyFaceCache,
		styles.NIGHT_BLACK,
		clip.Outline{Path: mouthPath.End()}.Op(),
	)

	happyFaceCallOp = macro.Stop()
	happyFaceCallOp.Add(gtx.Ops)
}
