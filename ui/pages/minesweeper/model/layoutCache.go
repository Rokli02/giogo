package model

import (
	"gioui.org/layout"
	"gioui.org/op"
)

type LayoutCache struct {
	Dimensions layout.Dimensions
	Macro      *op.CallOp
	Ops        *op.Ops
}
