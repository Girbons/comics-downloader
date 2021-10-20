package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareRequestMangakakalot(t *testing.T) {
	cc := NewComicClient()
	link := "http://mangakakalot.com"
	req, err := cc.PrepareRequest(link)

	assert.Equal(t, req.Header["Referer"], []string{link})
	assert.Nil(t, err)
}

func TestPrepareRequest(t *testing.T) {
	cc := NewComicClient()
	link := "http://foo.com"
	req, err := cc.PrepareRequest(link)

	assert.Equal(t, len(req.Header["Referer"]), 0)
	assert.Nil(t, err)
}
