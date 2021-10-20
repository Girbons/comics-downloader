package http

import (
	"net/http"
	"strings"
)

// ComicClient is the custom client.
type ComicClient struct {
	Client *http.Client
}

// NewComicClient returns a ComicClient instance.
func NewComicClient() *ComicClient {
	return &ComicClient{
		Client: &http.Client{},
	}
}

// PrepareRequest setup a `GET` request with customs headers.
func (c *ComicClient) PrepareRequest(link string) (*http.Request, error) {
	req, err := http.NewRequest("GET", link, nil)

	if strings.Contains(link, "mangakakalot") {
		// avoid that MangaKakalot forbids the request.
		req.Header.Add("Referer", link)
	}

	return req, err
}

// GET Performs a Get request..
func (c *ComicClient) Get(link string) (*http.Response, error) {
	request, err := c.PrepareRequest(link)

	if err != nil {
		return nil, err
	}

	return c.Client.Do(request)
}
