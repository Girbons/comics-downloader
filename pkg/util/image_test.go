package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageType(t *testing.T) {
	assert.Equal(t, ImageType("image/jpg"), "jpg")
	assert.Equal(t, ImageType("image/jpeg"), "jpg")
	assert.Equal(t, ImageType("image/png"), "png")
	assert.Equal(t, ImageType("image/gif"), "gif")
	assert.Equal(t, ImageType("foo"), "unknown")
}
