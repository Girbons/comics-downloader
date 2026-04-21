package util

import (
	"bufio"
	"errors"
	"image"
	"image/gif"
	"image/png"
	"io"
	"strings"

	"github.com/Girbons/comics-downloader/internal/logger"
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
	switch strings.ToLower(strings.TrimSpace(mimeStr)) {
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
func SaveImage(logger *logger.Logger, w io.Writer, content io.Reader, outputFormat ImageFormat, providedImageFormat ImageFormat) error {
	var (
		img image.Image
		err error
	)

	if strings.EqualFold(outputFormat.String(), ImgFormatIMG.String()) {
		_, err = io.Copy(w, content)
		return err
	}

	// TODO: we can optimize this by only decoding the image if the output format is different from the input format, otherwise we can just copy the content to the writer without decoding and encoding again
	// TODO: add avif support

	img, err = decodeInputImage(content, providedImageFormat)
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
		return encodeJPEG(w, img)
	case "png":
		pngEncoder := png.Encoder{CompressionLevel: png.BestCompression}
		return pngEncoder.Encode(w, img)
	case "webp":
		return webp.Encode(w, img, &webp.Options{Lossless: true})
	default:
		return errors.New("format not found")
	}
}

func decodeInputImage(content io.Reader, providedImageFormat ImageFormat) (image.Image, error) {
	if providedImageFormat == ImgFormatWEBP {
		return webp.Decode(content)
	}

	bufferedContent := bufio.NewReader(content)
	if providedImageFormat == ImgFormatJPG {
		return decodeJPEG(bufferedContent)
	}

	img, _, err := image.Decode(bufferedContent)
	return img, err
}
