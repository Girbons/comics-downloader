package detector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnsupportedSource(t *testing.T) {
	_, check, isDisabled := DetectSource("http://example.com")

	assert.False(t, check)
	assert.False(t, isDisabled)
}

func TestKnownSupportedSource(t *testing.T) {
	source, isSupported, isDisabled := DetectSource("https://comicextra.com/comic/some-comic/issue-1")

	assert.Contains(t, source, "comicextra")
	assert.True(t, isSupported)
	assert.False(t, isDisabled)
}

func TestKnownSupportedSourceMangadex(t *testing.T) {
	source, isSupported, isDisabled := DetectSource("https://mangadex.org/chapter/abc123")

	assert.Contains(t, source, "mangadex")
	assert.True(t, isSupported)
	assert.False(t, isDisabled)
}
