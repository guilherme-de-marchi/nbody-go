package util

import (
	"image/color"
)

const MAX_RGB_INT = 16777215 // 16777215 = r256^2 + g256 + b = int(rgb(255, 255, 255))

func IntToRgb(c int) color.RGBA {
	return color.RGBA{
		R: uint8((c >> 24) & 0xFF),
		G: uint8((c >> 16) & 0xFF),
		B: uint8((c >> 8) & 0xFF),
		A: uint8(c & 0xFF),
	}
}

func IntToRgbRange(c, r int) color.RGBA {
	c /= r / MAX_RGB_INT
	return color.RGBA{
		R: uint8((c >> 24) & 0xFF),
		G: uint8((c >> 16) & 0xFF),
		B: uint8((c >> 8) & 0xFF),
		A: uint8(c & 0xFF),
	}
}
