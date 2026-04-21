//go:build libjpeg

package util

import (
	"image"
	"io"

	libjpeg "github.com/pixiv/go-libjpeg/jpeg"
)

func decodeJPEG(content io.Reader) (image.Image, error) {
	return libjpeg.Decode(content, &libjpeg.DecoderOptions{})
}

func encodeJPEG(w io.Writer, img image.Image) error {
	return libjpeg.Encode(w, img, &libjpeg.EncoderOptions{Quality: 100})
}
