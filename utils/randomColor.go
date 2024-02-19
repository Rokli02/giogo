package utils

import (
	"image/color"
)

const spectrumBitShift = 5
const colorSpectrumDivision uint8 = 1 << spectrumBitShift
const colorUnit uint8 = uint8(256 / int(colorSpectrumDivision))
const (
	redToPurpleGroup uint8 = iota
	purpleToBlueGroup
	blueToBluishGreenGroup
	bluishGreenToGreenGroup
	greenToGreenishRedGroup
	greenishRedToRedGroup
	groupCount
)
const groupModulo = int(colorSpectrumDivision * groupCount)

func GetColor(i int) color.NRGBA {
	var result color.NRGBA
	colorGroupModulo := int16(i % groupModulo)

	if colorGroupModulo < 0 {
		colorGroupModulo += int16(groupModulo)
	}

	colorGroup := uint8(colorGroupModulo >> spectrumBitShift)
	colorIntensity := uint8(colorGroupModulo) - colorGroup*colorSpectrumDivision

	switch colorGroup {
	// Red-Blue
	case redToPurpleGroup:
		result = color.NRGBA{A: 0xFF, R: 0xFF, B: colorUnit * colorIntensity}
	// Blue
	case purpleToBlueGroup:
		result = color.NRGBA{A: 0xFF, B: 0xFF, R: colorUnit*(colorSpectrumDivision-colorIntensity) - 1}
	// Blue-Green
	case blueToBluishGreenGroup:
		result = color.NRGBA{A: 0xFF, B: 0xFF, G: colorUnit * colorIntensity}
	// Green
	case bluishGreenToGreenGroup:
		result = color.NRGBA{A: 0xFF, G: 0xFF, B: colorUnit*(colorSpectrumDivision-colorIntensity) - 1}
	// Green-Red
	case greenToGreenishRedGroup:
		result = color.NRGBA{A: 0xFF, G: 0xFF, R: colorUnit * colorIntensity}
	// Red
	case greenishRedToRedGroup:
		result = color.NRGBA{A: 0xFF, R: 0xFF, G: colorUnit*(colorSpectrumDivision-colorIntensity) - 1}
	}

	return result
}
