package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
)

// decode the image
func decodeImage(r io.Reader) (image.Image, error) {
	img, _, err := image.Decode(r)
	if err == nil {
		return img, nil
	}
	if seeker, ok := r.(io.ReadSeeker); ok {
		seeker.Seek(0, io.SeekStart)
		if img, err := png.Decode(seeker); err == nil {
			return img, nil
		}
		seeker.Seek(0, io.SeekStart)
		if img, err := jpeg.Decode(seeker); err == nil {
			return img, nil
		}
	}
	return nil, fmt.Errorf("unsupported image format or decoding failure: %v", err)
}
