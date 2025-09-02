package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"os"
)

func main() {
	// no parameters passed: print usage
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <the path to the image>\n", os.Args[0])
		os.Exit(1)
	}

	// reading input files
	var img image.Image
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to open file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	img, err = decodeImage(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode the image: %v\n", err)
		os.Exit(1)
	}

	// terminal size
	cols, rows, ok := getTermSize()
	if !ok {
		cols, rows = 80, 24
	}

	// resize the image
	targetW := cols
	targetH := rows * 2
	resized := resizeNearest(img, targetW, targetH)

	// output to terminal
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	for y := 0; y < targetH; y += 2 {
		for x := 0; x < targetW; x++ {
			top := colorToRGBA(resized.At(x, y))
			var bottom color.RGBA
			if y+1 < targetH {
				bottom = colorToRGBA(resized.At(x, y+1))
			} else {
				bottom = color.RGBA{0, 0, 0, 255}
			}

			fg := rgbToXterm256(top.R, top.G, top.B)
			bg := rgbToXterm256(bottom.R, bottom.G, bottom.B)
			fmt.Fprintf(writer, "\x1b[38;5;%dm\x1b[48;5;%dm▀", fg, bg)
		}
		fmt.Fprint(writer, "\x1b[0m\n")
	}
}
