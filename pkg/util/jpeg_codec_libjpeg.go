//go:build libjpeg

package util

import (
	"image"
	stdjpeg "image/jpeg"
	"io"

	libjpeg "github.com/pixiv/go-libjpeg/jpeg"
)

// // https://github.com/pixiv/go-libjpeg/blob/3da21a74767d9ffe29fcad7484ddd745f99e9f4c/jpeg/compress.go#L243
// var libjpegUnsupportedFormat = errors.New("unsupported image type")

func decodeJPEG(content io.Reader) (image.Image, error) {
	return libjpeg.Decode(content, &libjpeg.DecoderOptions{})
}

func encodeJPEG(w io.Writer, img image.Image) error {
	// err := libjpeg.Encode(w, img, &libjpeg.EncoderOptions{Quality: 100})
	// if err == nil {
	// 	return nil
	// }

	// if errors.Is(err, libjpegUnsupportedFormat) {
	// 	// try to use stdjpeg to encode the image because libjpeg doesn't support some image type
	// 	return stdjpeg.Encode(w, img, &stdjpeg.Options{Quality: 100})
	// }
	// return err

	return stdjpeg.Encode(w, img, &stdjpeg.Options{Quality: 100})
}
