package util

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"net/url"
	"strings"
)

// UrlSource will retrieve the url hostname.
func UrlSource(u string) (string, error) {
	parsedUrl, err := url.Parse(u)

	if err != nil {
		return "", err
	}

	return parsedUrl.Hostname(), nil
}

// IsUrlValid will exclude those url containing `.gif` and `logo`.
func IsUrlValid(url string) bool {
	return !strings.Contains(url, ".gif") && !strings.Contains(url, "logo") && !strings.Contains(url, "mobilebanner")
}

// ValueInSlice will check if a value is already inside the slice.
func IsValueInSlice(valueToCheck string, values []string) bool {
	for _, v := range values {
		if v == valueToCheck {
			return true
		}
	}
	return false
}

// Converts an image of any type to a PNG with 8-bit color depth
func ConvertTo8BitPNG(img image.Image, imgData *bytes.Buffer) error {
	b := img.Bounds()
	imgSet := image.NewRGBA(b)
	// Converts each pixel to a 32-bit RGBA pixel
	for y := 0; y < b.Max.Y; y++ {
		for x := 0; x < b.Max.X; x++ {
			newPixel := color.RGBAModel.Convert(img.At(x, y))
			imgSet.Set(x, y, newPixel)
		}
	}

	err := png.Encode(imgData, imgSet)
	if err != nil {
		return err
	}

	return nil
}

// ImageType return the image type
func ImageType(mimeStr string) (tp string) {
	switch mimeStr {
	case "image/png":
		tp = "png"
	case "image/jpg":
		tp = "jpg"
	case "image/jpeg":
		tp = "jpg"
	case "image/gif":
		tp = "gif"
	default:
		tp = "unknown"
	}
	return
}
