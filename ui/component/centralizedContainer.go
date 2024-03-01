package component

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
)

type CentralizedContainer struct {
	Container            layout.Flex
	VerticallyCentered   bool
	HorizontallyCentered bool

	marginInset unit.Dp
	marginBlock unit.Dp
}

func NewCentralizedContainer(verticallyCentered, horizontallyCentered bool) *CentralizedContainer {
	c := &CentralizedContainer{
		Container:            layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle},
		VerticallyCentered:   verticallyCentered,
		HorizontallyCentered: horizontallyCentered,
	}

	return c
}

func (cc *CentralizedContainer) Layout(gtx layout.Context, children ...layout.FlexChild) layout.Dimensions {
	rec := op.Record(gtx.Ops)
	cDims := cc.Container.Layout(gtx, children...)
	macro := rec.Stop()

	if cc.HorizontallyCentered {
		cc.marginInset = unit.Dp((gtx.Constraints.Max.X - cDims.Size.X) >> 1)
	}

	if cc.VerticallyCentered {
		cc.marginBlock = unit.Dp((gtx.Constraints.Max.Y - cDims.Size.Y) >> 1)
	}

	return layout.Inset{Left: cc.marginInset, Right: cc.marginInset, Top: cc.marginBlock, Bottom: cc.marginBlock}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		macro.Add(gtx.Ops)

		return cDims
	})
}
