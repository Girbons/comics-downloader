package util

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"strings"
)

// IMAGEREGEX to extract the image html tag
const IMAGEREGEX = `<img[^>]+src="([^">]+)"`

// ImageType return the image type
func ImageType(mimeStr string) (tp string) {
	switch mimeStr {
	case "image/png", "png":
		tp = "png"
	case "image/jpg", "jpg":
		tp = "jpg"
	case "image/jpeg", "jpeg":
		tp = "jpg"
	case "image/gif", "gif":
		tp = "gif"
	case "img":
		tp = "img"
	default:
		tp = "unknown"
	}
	return
}

// SaveImage saves an image from a given format
func SaveImage(w io.Writer, content io.Reader, format string) error {
	img, _, err := image.Decode(content)

	if err != nil {
		return err
	}

	switch strings.ToLower(format) {
	case "img":
		_, err = io.Copy(w, content)
		return err
	case "gif":
		return gif.Encode(w, img, nil)
	case "jpg", "jpeg":
		return jpeg.Encode(w, img, &jpeg.Options{Quality: 100})
	case "png":
		pngEncoder := png.Encoder{CompressionLevel: png.BestCompression}
		return pngEncoder.Encode(w, img)
	default:
		return errors.New("format not found")
	}
}
