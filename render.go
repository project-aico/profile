package main

import (
	"image"
	"image/color"
	"math"
)

// ProcessRetro8Bit transforms the image into a pixelated 8-bit style.
// It maintains the original image dimensions while creating a "fat pixel" effect.
func ProcessRetro8Bit(src image.Image) image.Image {
	bounds := src.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	// Adjust this value to change blockiness.
	// Lower number = larger blocks.
	lowResW := 100
	lowResH := (h * lowResW) / w

	// 1. Create a small "lo-fi" version of the image
	lowResImg := image.NewRGBA(image.Rect(0, 0, lowResW, lowResH))
	for y := 0; y < lowResH; y++ {
		for x := 0; x < lowResW; x++ {
			// Sample from original
			srcX := x * w / lowResW
			srcY := y * h / lowResH
			origColor := colorToRGBA(src.At(srcX+bounds.Min.X, srcY+bounds.Min.Y))

			// Quantize to 8-bit Xterm palette
			idx := rgbToXterm256(origColor.R, origColor.G, origColor.B)

			// Set the lo-fi pixel using the palette color
			lowResImg.Set(x, y, xtermIndexToRGBA(idx))
		}
	}

	// 2. Scale it back up to original size using Nearest Neighbor (pixelated)
	out := image.NewRGBA(bounds)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			lx := x * lowResW / w
			ly := y * lowResH / h
			out.Set(x+bounds.Min.X, y+bounds.Min.Y, lowResImg.At(lx, ly))
		}
	}

	return out
}

func colorToRGBA(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}

func rgbToXterm256(r, g, b uint8) int {
	r6 := int(math.Round(float64(r) * 5.0 / 255.0))
	g6 := int(math.Round(float64(g) * 5.0 / 255.0))
	b6 := int(math.Round(float64(b) * 5.0 / 255.0))
	return 16 + 36*r6 + 6*g6 + b6
}

func xtermIndexToRGBA(idx int) color.RGBA {
	if idx < 16 {
		return standardColors[idx]
	}
	if idx < 232 {
		idx -= 16
		vals := []uint8{0, 95, 135, 175, 215, 255}
		return color.RGBA{vals[(idx/36)%6], vals[(idx/6)%6], vals[idx%6], 255}
	}
	gray := uint8((idx-232)*10 + 8)
	return color.RGBA{gray, gray, gray, 255}
}

var standardColors = []color.RGBA{
	{0, 0, 0, 255}, {128, 0, 0, 255}, {0, 128, 0, 255}, {128, 128, 0, 255},
	{0, 0, 128, 255}, {128, 0, 128, 255}, {0, 128, 128, 255}, {192, 192, 192, 255},
	{128, 128, 128, 255}, {255, 0, 0, 255}, {0, 255, 0, 255}, {255, 255, 0, 255},
	{0, 0, 255, 255}, {255, 0, 255, 255}, {0, 255, 255, 255}, {255, 255, 255, 255},
}
