package detector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnsupportedSource(t *testing.T) {
	_, check, isDisabled := DetectComic("http://example.com")

	assert.False(t, check)
	assert.False(t, isDisabled)
}
