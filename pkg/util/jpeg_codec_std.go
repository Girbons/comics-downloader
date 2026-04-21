//go:build !libjpeg

package util

import (
	"image"
	"image/jpeg"
	"io"
)

func decodeJPEG(content io.Reader) (image.Image, error) {
	return jpeg.Decode(content)
}

func encodeJPEG(w io.Writer, img image.Image) error {
	return jpeg.Encode(w, img, &jpeg.Options{Quality: 100})
}
