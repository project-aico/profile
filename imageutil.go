package main

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// decodeImage handles decoding PNG and JPEG
func decodeImage(r io.Reader) (image.Image, error) {
	img, _, err := image.Decode(r)
	return img, err
}

// saveImage saves the resulting image as PNG or JPEG
func saveImage(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg":
		return jpeg.Encode(f, img, &jpeg.Options{Quality: 95})
	case ".png":
		return png.Encode(f, img)
	default:
		// Default to PNG if unknown
		return png.Encode(f, img)
	}
}
