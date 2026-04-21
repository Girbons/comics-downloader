package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageType(t *testing.T) {
	assert.Equal(t, ImageType("image/jpg"), ImageFormat("jpg"))
	assert.Equal(t, ImageType("image/jpeg"), ImageFormat("jpg"))
	assert.Equal(t, ImageType("image/png"), ImageFormat("png"))
	assert.Equal(t, ImageType("image/gif"), ImageFormat("gif"))
	assert.Equal(t, ImageType("image/webp"), ImageFormat("webp"))
	assert.Equal(t, ImageType("foo"), ImageFormat("unknown"))
}
