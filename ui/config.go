package ui

import (
	"gioui.org/layout"
)

type ApplicationCycles interface {
	Initialize()
	Restart()
	Close()
	Layout(gtx layout.Context) layout.Dimensions
}
