package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComicOutputFormatHandlesInvalid(t *testing.T) {
	_, err := ToComicOutputFormat("invalid_format")

	if assert.Error(t, err) {
		assert.Equal(t, InvalidOutputFormatError, err)
	}
}

func TestComicOutputFormatHandlesValid(t *testing.T) {
	format, err := ToComicOutputFormat("cbz")

	if assert.NoError(t, err) {
		assert.Equal(t, CBZ, format)
	}
}
