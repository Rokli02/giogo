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

func drawMarkedCell(gtx *layout.Context, size image.Point) {
	drawHiddenCell(gtx, size)

	tmpGtx := gtx.Constraints
	gtx.Constraints.Max.X = gtx.Constraints.Max.X - (int(shadowThicknes) << 1)
	gtx.Constraints.Max.Y = gtx.Constraints.Max.Y - (int(shadowThicknes) << 1)

	defer op.Offset(image.Point{X: int(shadowThicknes), Y: int(shadowThicknes)}).Push(gtx.Ops).Pop()

	markedField := assets.GetWidgetImage("marked", size.X)
	markedField.Position = layout.Center
	markedField.Layout(*gtx)

	gtx.Constraints = tmpGtx
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

var hiddenCellOps = new(op.Ops)

func drawHiddenCell(gtx *layout.Context, size image.Point) {
	macro := op.Record(hiddenCellOps)

	shine := clip.Path{}
	shine.Begin(gtx.Ops)
	shine.MoveTo(f32.Pt(float32(0), float32(size.Y)))
	shine.LineTo(f32.Pt(0, 0))
	shine.LineTo(f32.Pt(float32(size.X), 0))
	shine.LineTo(f32.Pt(float32(size.X)-shadowThicknes, shadowThicknes))
	shine.LineTo(f32.Pt(shadowThicknes, shadowThicknes))
	shine.LineTo(f32.Pt(shadowThicknes, float32(size.Y)-shadowThicknes))
	shine.Close()

	paint.FillShape(gtx.Ops, shineColor, clip.Outline{Path: shine.End()}.Op())

	shadow := clip.Path{}
	shadow.Begin(gtx.Ops)
	shadow.MoveTo(f32.Pt(float32(0), float32(size.Y)))
	shadow.LineTo(f32.Pt(float32(size.X), float32(size.Y)))
	shadow.LineTo(f32.Pt(float32(size.X), 0))
	shadow.LineTo(f32.Pt(float32(size.X)-shadowThicknes, shadowThicknes))
	shadow.LineTo(f32.Pt(float32(size.X)-shadowThicknes, float32(size.Y)-shadowThicknes))
	shadow.LineTo(f32.Pt(shadowThicknes, float32(size.Y)-shadowThicknes))
	shadow.Close()

	paint.FillShape(gtx.Ops, shadowColor, clip.Outline{Path: shadow.End()}.Op())

	macro.Stop().Add(gtx.Ops)
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

var neutralFaceCache = new(op.Ops)

func drawNeutralFace(gtx *layout.Context) {
	macro := op.Record(neutralFaceCache)

	middlePoint := image.Point{X: gtx.Constraints.Max.X >> 1, Y: gtx.Constraints.Max.Y >> 1}
	eyeOffset := image.Point{X: middlePoint.X >> 1, Y: middlePoint.Y / 3}
	eyeSize := 2
	mouthOffset := image.Point{X: middlePoint.X / 3, Y: middlePoint.Y / 3}
	mouthSize := 1

	// Left eye
	paint.FillShape(
		gtx.Ops,
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
		}.Op(gtx.Ops),
	)

	// Right eye
	paint.FillShape(
		gtx.Ops,
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
		}.Op(gtx.Ops),
	)

	// Mouth
	paint.FillShape(
		gtx.Ops,
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

	macro.Stop().Add(gtx.Ops)
}

var madFaceCache = new(op.Ops)

func drawMadFace(gtx *layout.Context) {
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
		gtx.Ops,
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
		}.Op(gtx.Ops),
	)

	// Left Eyebrow
	leftEyebrowPath := clip.Path{}
	leftEyebrowPath.Begin(gtx.Ops)
	leftEyebrowPath.MoveTo(f32.Pt(float32(middlePoint.X-eyeOffset.X-(eyebrowWidth>>1)), float32(middlePoint.Y-eyeOffset.Y-eyebrowOffset)))
	leftEyebrowPath.Line(f32.Pt(eyebrowWidth, eyebrowOffset))
	leftEyebrowPath.Line(f32.Pt(0, -eyebrowSize))
	leftEyebrowPath.Line(f32.Pt(-eyebrowWidth, -eyebrowOffset))
	leftEyebrowPath.Close()

	paint.FillShape(
		gtx.Ops,
		styles.NIGHT_BLACK,
		clip.Outline{Path: leftEyebrowPath.End()}.Op(),
	)

	// Right eye
	paint.FillShape(
		gtx.Ops,
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
		}.Op(gtx.Ops),
	)

	// Right eyebrow
	rightEyebrowPath := clip.Path{}
	rightEyebrowPath.Begin(gtx.Ops)
	rightEyebrowPath.MoveTo(f32.Pt(float32(middlePoint.X+eyeOffset.X+(eyebrowWidth>>1)), float32(middlePoint.Y-eyeOffset.Y-eyebrowOffset)))
	rightEyebrowPath.Line(f32.Pt(-eyebrowWidth, eyebrowOffset))
	rightEyebrowPath.Line(f32.Pt(0, -eyebrowSize))
	rightEyebrowPath.Line(f32.Pt(eyebrowWidth, -eyebrowOffset))
	rightEyebrowPath.Close()

	paint.FillShape(
		gtx.Ops,
		styles.NIGHT_BLACK,
		clip.Outline{Path: rightEyebrowPath.End()}.Op(),
	)

	// Mouth
	mouthPath := clip.Path{}
	mouthPath.Begin(gtx.Ops)
	mouthPath.MoveTo(f32.Pt(float32(middlePoint.X-mouthOffset.X), float32(middlePoint.Y+mouthOffset.Y-mouthSize)))
	mouthPath.QuadTo(f32.Pt(float32(middlePoint.X), float32(middlePoint.Y+mouthOffset.Y-mouthCurve-mouthSize)), f32.Pt(float32(middlePoint.X+mouthOffset.X), float32(middlePoint.Y+mouthOffset.Y-mouthSize)))
	mouthPath.Line(f32.Pt(0, float32(mouthSize)))
	mouthPath.QuadTo(f32.Pt(float32(middlePoint.X), float32(middlePoint.Y+mouthOffset.Y-mouthCurve)), f32.Pt(float32(middlePoint.X-mouthOffset.X), float32(middlePoint.Y+mouthOffset.Y)))
	mouthPath.Close()
	paint.FillShape(
		gtx.Ops,
		styles.NIGHT_BLACK,
		clip.Outline{Path: mouthPath.End()}.Op(),
	)

	macro.Stop().Add(madFaceCache)
}

var happyFaceCache = new(op.Ops)

func drawHappyFace(gtx *layout.Context) {
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
		gtx.Ops,
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
		}.Op(gtx.Ops),
	)

	// Right eye
	paint.FillShape(
		gtx.Ops,
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
		}.Op(gtx.Ops),
	)

	// Mouth
	mouthPath := clip.Path{}
	mouthPath.Begin(gtx.Ops)
	mouthPath.MoveTo(f32.Pt(float32(middlePoint.X-mouthOffset.X), float32(middlePoint.Y+mouthOffset.Y-mouthSize)))
	mouthPath.QuadTo(f32.Pt(float32(middlePoint.X), float32(middlePoint.Y+mouthOffset.Y+mouthCurve-mouthSize)), f32.Pt(float32(middlePoint.X+mouthOffset.X), float32(middlePoint.Y+mouthOffset.Y-mouthSize)))
	mouthPath.Line(f32.Pt(0, float32(mouthSize)))
	mouthPath.QuadTo(f32.Pt(float32(middlePoint.X), float32(middlePoint.Y+mouthOffset.Y+mouthCurve)), f32.Pt(float32(middlePoint.X-mouthOffset.X), float32(middlePoint.Y+mouthOffset.Y)))
	mouthPath.Close()
	paint.FillShape(
		gtx.Ops,
		styles.NIGHT_BLACK,
		clip.Outline{Path: mouthPath.End()}.Op(),
	)

	macro.Stop().Add(happyFaceCache)
}
