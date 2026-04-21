package util

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"strings"

	"github.com/chai2010/webp"
)

// IMAGEREGEX to extract the image html tag
const IMAGEREGEX = `<img[^>]+src="([^">]+)"`

type ImageFormat string

func (f ImageFormat) String() string {
	return string(f)
}

const (
	ImgFormatPNG     ImageFormat = "png"
	ImgFormatJPG     ImageFormat = "jpg"
	ImgFormatGIF     ImageFormat = "gif"
	ImgFormatWEBP    ImageFormat = "webp"
	ImgFormatIMG     ImageFormat = "img"
	ImgFormatUnknown ImageFormat = "unknown"
)

// ImageType return the image type
func ImageType(mimeStr string) (format ImageFormat) {
	mimeStr = strings.ToLower(mimeStr)
	switch mimeStr {
	case "image/png", "png":
		format = ImgFormatPNG
	case "image/jpg", "jpg", "image/jpeg", "jpeg":
		format = ImgFormatJPG
	case "image/gif", "gif":
		format = ImgFormatGIF
	case "image/webp", "webp":
		format = ImgFormatWEBP
	case "img":
		format = ImgFormatIMG
	default:
		format = ImgFormatUnknown
	}
	return
}

// SaveImage saves an image from a given format
func SaveImage(w io.Writer, content io.Reader, outputFormat ImageFormat, providedImageFormat ImageFormat) error {
	var (
		img image.Image
		err error
	)

	// TODO: we can optimize this by only decoding the image if the output format is different from the input format, otherwise we can just copy the content to the writer without decoding and encoding again
	// TODO: add avif support

	if providedImageFormat == ImgFormatWEBP {
		img, err = webp.Decode(content)
	} else {
		img, _, err = image.Decode(content)
	}

	if err != nil {
		return err
	}

	switch strings.ToLower(outputFormat.String()) {
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
	case "webp":
		return webp.Encode(w, img, &webp.Options{Lossless: true})
	default:
		return errors.New("format not found")
	}
}
