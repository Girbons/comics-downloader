package util

import (
	"bytes"
	"image/png"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrlSource(t *testing.T) {
	result, _ := UrlSource("http://example.com")
	assert.Equal(t, "example.com", result)
}

func TestIsValueInSlice(t *testing.T) {
	s := []string{"foo"}
	assert.Equal(t, true, IsValueInSlice("foo", s))
	assert.Equal(t, false, IsValueInSlice("bar", s))
}

func TestIsUrlValid(t *testing.T) {
	validUrl := IsUrlValid("http://example.com")
	gifUrl := IsUrlValid("http://foo.gif")
	logoUrl := IsUrlValid("http://foo.logo")

	assert.Equal(t, true, validUrl)
	assert.Equal(t, false, gifUrl)
	assert.Equal(t, false, logoUrl)
}

func TestImageType(t *testing.T) {
	assert.Equal(t, ImageType("image/jpg"), "jpg")
	assert.Equal(t, ImageType("image/jpeg"), "jpg")
	assert.Equal(t, ImageType("image/png"), "png")
	assert.Equal(t, ImageType("image/gif"), "gif")
	assert.Equal(t, ImageType("foo"), "unknown")
}

func TestConvertImage(t *testing.T) {
	imgURL := "http://via.placeholder.com/150"

	resp, _ := http.Get(imgURL)
	defer resp.Body.Close()

	img, _ := png.Decode(resp.Body)
	imgData := new(bytes.Buffer)
	ConvertTo8BitPNG(img, imgData)
}
