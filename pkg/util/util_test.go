package util

import (
	"bytes"
	"image/png"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrimAndSplitURL(t *testing.T) {
	result := TrimAndSplitURL("http://example.com/path/")
	assert.Equal(t, []string{"http:", "", "example.com", "path"}, result)
}

func TestUrlSource(t *testing.T) {
	result, _ := URLSource("http://example.com")
	assert.Equal(t, "example.com", result)
}

func TestIsValueInSlice(t *testing.T) {
	s := []string{"foo"}
	assert.True(t, IsValueInSlice("foo", s))
	assert.False(t, IsValueInSlice("bar", s))
}

func TestIsUrlValid(t *testing.T) {
	validUrl := IsURLValid("http://example.com")
	gifUrl := IsURLValid("http://foo.gif")
	logoUrl := IsURLValid("http://foo.logo")

	assert.True(t, validUrl)
	assert.False(t, gifUrl)
	assert.False(t, logoUrl)
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

	resp, err := http.Get(imgURL)
	assert.Nil(t, err)

	defer resp.Body.Close()

	img, _ := png.Decode(resp.Body)
	imgData := new(bytes.Buffer)
	err = ConvertToJPG(img, imgData)

	assert.Nil(t, err)
}

func TestPathSetup(t *testing.T) {
	result, err := PathSetup("example-source", "comic-name")

	assert.Nil(t, err)
	assert.Contains(t, result, "example-source")
	assert.Contains(t, result, "comic-name")
}

func TestParse(t *testing.T) {
	result := Parse("aaa/bbb/ccc")

	assert.Equal(t, "aaa_bbb_ccc", result)
}

func TestGenerateFileName(t *testing.T) {
	result := GenerateFileName("path/to/something", "invalid_character", "pdf")

	assert.Equal(t, "path/to/something/invalid_character.pdf", result)
}
