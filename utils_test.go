package main

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

func TestValueInSlice(t *testing.T) {
	s := []string{"foo"}
	assert.Equal(t, true, CheckValueInSlice("foo", s))
	assert.Equal(t, false, CheckValueInSlice("bar", s))
}
