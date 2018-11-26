package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitUrl(t *testing.T) {
	result := SplitUrl("http://example.com/")
	assert.Equal(t, 4, len(result))
}

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
