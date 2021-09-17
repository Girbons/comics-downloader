package util

import (
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
	validURL := IsURLValid("http://example.com")
	gifURL := IsURLValid("http://foo.gif")
	logoURL := IsURLValid("http://foo.logo")

	assert.True(t, validURL)
	assert.False(t, gifURL)
	assert.False(t, logoURL)
}

func TestParse(t *testing.T) {
	result := Parse("aaa/bbb/ccc")

	assert.Equal(t, "aaa_bbb_ccc", result)
}
