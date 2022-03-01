package util

import (
	"image/color"
)

func IntToRgb(c int) color.RGBA {
	return color.RGBA{
		R: uint8((c >> 24) & 0xFF),
		G: uint8((c >> 16) & 0xFF),
		B: uint8((c >> 8) & 0xFF),
		A: uint8(c & 0xFF),
	}
}

func IntToRgbRange(c, r int) color.RGBA {
	// 16777215 = r256^2 + g256 + b = int(rgb(255, 255, 255))
	c /= r / 16777215
	return color.RGBA{
		R: uint8((c >> 24) & 0xFF),
		G: uint8((c >> 16) & 0xFF),
		B: uint8((c >> 8) & 0xFF),
		A: uint8(c & 0xFF),
	}
}
