package styles

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/widget/material"
)

var (
	BACKGROUND_COLOR = color.NRGBA{R: 0x20, G: 0x20, B: 0x20, A: 0xFF}
	GREEN            = color.NRGBA{A: 0xFF, R: 0x04, G: 0x80, B: 0x25} // 04 80 25
	YELLOWISH_GREEN  = color.NRGBA{A: 0xFF, R: 0x6c, G: 0x96, B: 0x1e} // 6c 96 1e
	YELLOW           = color.NRGBA{A: 0xFF, R: 0xbd, G: 0xbd, B: 0x0d} // bd bd 0d
	GOLD             = color.NRGBA{A: 0xFF, R: 0xb3, G: 0x87, B: 0x19} // b3 87 19
	ORANGE           = color.NRGBA{A: 0xFF, R: 0xc2, G: 0x63, B: 0x0a} // c2 63 0a
	BLOOD_ORANGE     = color.NRGBA{A: 0xFF, R: 0xb8, G: 0x30, B: 0x0b} // b8 30 0b
	RED              = color.NRGBA{A: 0xFF, R: 0xad, G: 0x02, B: 0x02} // ad 02 02
	BLUE             = color.NRGBA{A: 0xFF, R: 0x05, G: 0x08, B: 0xb0} // 05 08 b0
	NIGHT_BLACK      = color.NRGBA{A: 0xFF}
)

var (
	HEADER_BACKGROUND = color.NRGBA{A: 0xFF, R: 0xa3, G: 0xa3, B: 0xa3}
	TEXT_SHADOW       = color.NRGBA{A: 0x9A}
)

var (
	MenuWindowSizes = image.Point{360, 500}
)

var (
	MaterialTheme *material.Theme
	CancelTheme   *material.Theme
)

func InitializeStyles() {
	fmt.Println("Initializing Styles and Themes")

	MaterialTheme = material.NewTheme()
	CancelTheme = material.NewTheme()

	CancelTheme.ContrastFg = color.NRGBA{A: 0xFF, R: 0xe8, G: 0xd9, B: 0xda} //#E8D9DA
	// CancelTheme.ContrastFg = color.NRGBA{A: 0xFF, R: 0x36, G: 0x09, B: 0x05} //#360905	✅
	CancelTheme.ContrastBg = color.NRGBA{A: 0xFF, R: 0xc2, G: 0x28, B: 0x1d} //#C2281D
}
