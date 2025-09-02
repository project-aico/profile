package main

import (
	"image"
	"image/color"
	"math"
	"os"

	"golang.org/x/term"
)

// get the size of terminal
func getTermSize() (cols, rows int, ok bool) {
	fd := int(os.Stdout.Fd())
	w, h, err := term.GetSize(fd)
	if err != nil {
		return 0, 0, false
	}
	return w, h, true
}

// RGBA conversion
func colorToRGBA(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}

// nearest neighbor scaling
func resizeNearest(src image.Image, w, h int) *image.RGBA {
	r := src.Bounds()
	sw, sh := r.Dx(), r.Dy()
	out := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		srcY := int(float64(y) * float64(sh) / float64(h))
		if srcY >= sh {
			srcY = sh - 1
		}
		for x := 0; x < w; x++ {
			srcX := int(float64(x) * float64(sw) / float64(w))
			if srcX >= sw {
				srcX = sw - 1
			}
			out.Set(x, y, src.At(srcX+r.Min.X, srcY+r.Min.Y))
		}
	}
	return out
}

// RGB to Xterm-256
func rgbToXterm256(r, g, b uint8) int {
	r6 := int(math.Round(float64(r) * 5.0 / 255.0))
	g6 := int(math.Round(float64(g) * 5.0 / 255.0))
	b6 := int(math.Round(float64(b) * 5.0 / 255.0))
	cubeIdx := 16 + 36*r6 + 6*g6 + b6

	avg := float64(r+g+b) / 3.0
	grayIdx := int(math.Round((avg - 8.0) / 247.0 * 24.0))
	if grayIdx < 0 {
		grayIdx = 0
	}
	if grayIdx > 23 {
		grayIdx = 23
	}
	grayCode := 232 + grayIdx

	cr := int(float64(r6) * 255.0 / 5.0)
	cg := int(float64(g6) * 255.0 / 5.0)
	cb := int(float64(b6) * 255.0 / 5.0)
	cubeDist := colorDistSquared(int(r), int(g), int(b), cr, cg, cb)

	gval := 8 + grayIdx*10
	grayDist := colorDistSquared(int(r), int(g), int(b), gval, gval, gval)

	if grayDist < cubeDist {
		return grayCode
	}
	return cubeIdx
}

func colorDistSquared(r1, g1, b1, r2, g2, b2 int) int {
	dr, dg, db := r1-r2, g1-g2, b1-b2
	return dr*dr + dg*dg + db*db
}
