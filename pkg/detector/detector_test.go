package detector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectComicExtra(t *testing.T) {
	source, check, isDisabled := DetectComic("http://www.comicextra.com/daredevil-2015/chapter-600/full")

	assert.True(t, check)
	assert.True(t, isDisabled)
	assert.Equal(t, "www.comicextra.com", source)
}

func TestUnsupportedSource(t *testing.T) {
	_, check, isDisabled := DetectComic("http://example.com")

	assert.False(t, check)
	assert.False(t, isDisabled)
}
