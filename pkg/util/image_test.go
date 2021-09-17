package util

import (
	"bytes"
	"image/png"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertImage(t *testing.T) {
	imgURL := "http://via.placeholder.com/150"

	resp, err := http.Get(imgURL)
	assert.Nil(t, err)

	defer resp.Body.Close()

	img, _ := png.Decode(resp.Body)
	imgData := new(bytes.Buffer)
	err = ConvertToJPG(img, imgData)

	assert.Nil(t, err)
}

func TestImageType(t *testing.T) {
	assert.Equal(t, ImageType("image/jpg"), "jpg")
	assert.Equal(t, ImageType("image/jpeg"), "jpg")
	assert.Equal(t, ImageType("image/png"), "png")
	assert.Equal(t, ImageType("image/gif"), "gif")
	assert.Equal(t, ImageType("foo"), "unknown")
}
