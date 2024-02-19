package utils

import (
	"image"

	"gioui.org/layout"
)

func GetSize(gtx *layout.Context) image.Point {
	return image.Point{gtx.Constraints.Max.X - gtx.Constraints.Min.X, gtx.Constraints.Max.Y - gtx.Constraints.Min.Y}
}
