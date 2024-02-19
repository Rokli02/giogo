package scroller

import (
	"fmt"
	"giogo/ui/component"
	"giogo/ui/styles"
	"giogo/utils"
	"image"
	"image/color"

	"gioui.org/app"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Indicator struct {
	offset     int
	normOffset int
	height     int
}

type ScrollPage struct {
	NumOfElements int
	Width         int

	heightOfChildren     int
	size                 image.Point
	w                    *app.Window
	currentStartingIndex int
	list                 layout.List
	clicks               []widget.Clickable

	indicator            Indicator
	prevScrollingDims    image.Point
	trackHeight          int
	visibleChildrenCount int
	dragStart            float32
}

func NewScrollPage(w *app.Window, numOfElements int) *ScrollPage {
	scrollPage := &ScrollPage{
		w:             w,
		NumOfElements: numOfElements,
		Width:         400,
	}

	return scrollPage
}

func (sp *ScrollPage) Initialize() {
	sp.w.Option(func(_ unit.Metric, c *app.Config) {
		c.Title = "Görgető"
		c.MinSize = image.Pt(int(sp.Width)+16, 200)
		c.MaxSize = image.Point{}

		if sp.size != (image.Point{}) {
			c.Size = sp.size
		}
	})

	sp.dragStart = -1

	if len(sp.clicks) == 0 || len(sp.clicks) != int(sp.NumOfElements) {
		sp.list = layout.List{Axis: layout.Vertical, Alignment: layout.Start}
		sp.clicks = make([]widget.Clickable, sp.NumOfElements)
		sp.heightOfChildren = 62
	}
}

func (sp *ScrollPage) Close() {
	sp.w.Option(func(_ unit.Metric, c *app.Config) {
		sp.size = c.Size
	})
}

func (sp *ScrollPage) Restart() {
	sp.list.Position.First = 0
	sp.list.Position.Offset = 0
}

func (sp *ScrollPage) Layout(gtx layout.Context) layout.Dimensions {
	defer utils.SetBackgroundColor(&gtx, styles.BACKGROUND_COLOR).Pop()

	inlineMargin := (gtx.Constraints.Max.X - sp.Width) >> 1

	return layout.Inset{
		Left:   unit.Dp(inlineMargin),
		Right:  unit.Dp(inlineMargin),
		Top:    unit.Dp(8),
		Bottom: unit.Dp(8),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		defer utils.SetBackgroundColor(&gtx, color.NRGBA{A: 0xFF, R: 0xA8, G: 0xA8, B: 0xA8}).Pop()

		return layout.Stack{Alignment: layout.E}.Layout(gtx,
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return sp.list.Layout(gtx, int(sp.NumOfElements), func(gtx layout.Context, index int) layout.Dimensions {
					defer utils.SetBackgroundColor(&gtx, utils.GetColor(index-sp.currentStartingIndex)).Pop()

					if sp.clicks[index].Clicked(gtx) {
						sp.currentStartingIndex = index
						op.InvalidateOp{}.Add(gtx.Ops)
					}

					return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						// Main Item Macro
						height := sp.heightOfChildren - 32

						centerDim := layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							material.Label(styles.MaterialTheme, unit.Sp(14), fmt.Sprintf("%d.: Test String", index)).Layout(gtx)

							return layout.Dimensions{Size: image.Point{X: gtx.Constraints.Max.X, Y: height}}
						})

						// Checkbox
						checkboxSize := 16
						checkboxOffset := op.Offset(image.Pt(gtx.Constraints.Max.X-checkboxSize-6, (height-checkboxSize)>>1)).Push(gtx.Ops)
						component.Checkbox{
							Click:     &sp.clicks[index],
							Size:      checkboxSize,
							IsChecked: sp.currentStartingIndex == index,
						}.Layout(gtx)
						checkboxOffset.Pop()

						return centerDim
					})
				})
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Left: unit.Dp(sp.Width - 10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Max.X = 10

					trackStack := utils.SetBackgroundColor(&gtx, color.NRGBA{A: 0x77})

					sp.setChangedScrollingStates(&gtx)
					sp.handleScrolling(&gtx)

					pointer.InputOp{
						Tag:   &sp.dragStart,
						Kinds: pointer.Drag | pointer.Press,
					}.Add(gtx.Ops)

					sp.indicator.normOffset = min(sp.indicator.offset*sp.trackHeight/(sp.list.Position.Length-sp.visibleChildrenCount*sp.heightOfChildren), sp.trackHeight)

					indicatorOffset := op.Offset(image.Point{0, sp.indicator.normOffset}).Push(gtx.Ops)
					gtx.Constraints.Max.Y = sp.indicator.height

					utils.SetBackgroundColor(&gtx, color.NRGBA{A: 0x66, G: 0x30}).Pop()

					gtx.Constraints.Max.Y = sp.prevScrollingDims.Y
					indicatorOffset.Pop()

					trackStack.Pop()

					return layout.Dimensions{Size: gtx.Constraints.Max}
				})
			}),
		)
	})
}

func (sp *ScrollPage) handleScrolling(gtx *layout.Context) {
	deepestScrollHeight := sp.list.Position.Length - sp.visibleChildrenCount*sp.heightOfChildren

	for _, evt := range gtx.Events(&sp.dragStart) {
		switch e := evt.(type) {
		case pointer.Event:
			switch e.Kind {
			case pointer.Press:
				clickedInIndicator := int(e.Position.Y) > sp.indicator.normOffset && int(e.Position.Y) < sp.indicator.normOffset+sp.indicator.height

				sp.dragStart = e.Position.Y

				if !clickedInIndicator {
					jumpTemp := ((int(e.Position.Y) - (sp.indicator.height >> 1)) * deepestScrollHeight / sp.trackHeight)

					if jumpTemp < 0 {
						sp.list.Position.First = 0
						sp.list.Position.Offset = 0

						break
					} else if jumpTemp > deepestScrollHeight {
						sp.list.Position.First = int(sp.NumOfElements) - sp.list.Position.Count
						sp.list.Position.Offset = (deepestScrollHeight - sp.prevScrollingDims.Y) % (sp.heightOfChildren) // Elosztani úgy, hogy a first elemek + még egy kics

						break
					}

					sp.list.Position.First = jumpTemp / sp.heightOfChildren
					sp.list.Position.Offset = jumpTemp % sp.heightOfChildren
				}
			case pointer.Drag:
				if sp.dragStart == -1 {
					break
				}

				dragLength := int(e.Position.Y-sp.dragStart) * deepestScrollHeight / sp.trackHeight // 'sp.trackHeight' helyett lehet még 'sp.prevScrollingDims.Y' is talán

				listScrolledOffset := dragLength % sp.heightOfChildren
				listScrolledFirst := dragLength / sp.heightOfChildren

				sp.list.Position.First += listScrolledFirst
				sp.list.Position.Offset += listScrolledOffset

				// FIXME: Ha túl nagy a dragLength, akkor még itt nem fogja meg, picit elcsúszik
				if sp.indicator.normOffset > 0 && sp.indicator.normOffset < sp.trackHeight {
					sp.dragStart = e.Position.Y
				}
			}

			op.InvalidateOp{}.Add(gtx.Ops)
		}
	}

	sp.indicator.offset = max(0, min(deepestScrollHeight, sp.list.Position.First*sp.heightOfChildren+sp.list.Position.Offset))
}

func (sp *ScrollPage) setChangedScrollingStates(gtx *layout.Context) {
	if gtx.Constraints.Max.Y != sp.prevScrollingDims.Y {
		sp.prevScrollingDims.Y = gtx.Constraints.Max.Y

		sp.visibleChildrenCount = sp.list.Position.Count
		sp.indicator.height = max(40, int(float32(sp.prevScrollingDims.Y*sp.list.Position.Count)/float32(sp.NumOfElements)))
		sp.trackHeight = sp.prevScrollingDims.Y - sp.indicator.height
	}

	if gtx.Constraints.Max.X != sp.prevScrollingDims.X {
		sp.prevScrollingDims.X = gtx.Constraints.Max.X
	}
}
