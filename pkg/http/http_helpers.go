package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	urlpkg "net/url"
)

type RequestConfig struct {
	request *http.Request
}

func (rc *RequestConfig) apply(opts []RequestOption) {
	for _, opt := range opts {
		opt(rc)
	}
}

type RequestOption func(*RequestConfig)

// withHttpHeader adds a header to the request, but only if it doesn't already exist. This allows users to set headers at the client level and override them on a per-request basis without worrying about conflicts.
func WithHttpHeader(key, value string) RequestOption {
	return func(rc *RequestConfig) {
		if rc.request == nil {
			return
		}
		if rc.request.Header == nil {
			rc.request.Header = make(http.Header)
		}
		// avoid overwriting existing headers, user may have set it in the client already
		_, ok := rc.request.Header[key]
		if ok {
			return
		}
		rc.request.Header.Set(key, value)
	}
}

// hostFromURL extracts the hostname from a given URL string. If the URL is invalid, it returns an empty string.
func hostFromURL(link string) string {
	parsed, err := urlpkg.Parse(link)
	if err != nil {
		return ""
	}
	return parsed.Host
}

// buildRequest constructs an HTTP request based on the provided link and options. It applies any request options to the request configuration and attaches the context to the request if provided.
func (c *ComicClient) buildRequest(ctx context.Context, link string, opts ...RequestOption) (*http.Request, error) {
	reqConfig := &RequestConfig{}
	req, err := c.PrepareRequest(link, hostFromURL(link))
	if err != nil {
		return nil, err
	}
	reqConfig.request = req

	// apply config to request
	reqConfig.apply(opts)

	if ctx != nil {
		req = req.WithContext(ctx)
	}
	return req, nil
}

// FetchHTML retrieves the HTML content from the specified link with retry, timeout, rate limiting, and caching policies applied.
func (c *ComicClient) FetchHTML(ctx context.Context, link string, opts ...RequestOption) (string, error) {
	response, err := c.FetchBytes(ctx, link, opts...)
	if err != nil {
		return "", err
	}
	return string(response), nil
}

// FetchJSON retrieves the JSON content from the specified link, applies retry, timeout, rate limiting, and caching policies, and unmarshals the response into the provided target structure.
func (c *ComicClient) FetchJSON(ctx context.Context, link string, target interface{}, opts ...RequestOption) error {
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}

	data, err := c.FetchBytes(ctx, link, opts...)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return fmt.Errorf("empty response for %s", link)
	}
	return json.Unmarshal(data, target)
}

// FetchBytes retrieves the raw byte content from the specified link, applying retry, timeout, rate limiting, and caching policies. It returns the response body as a byte slice.
func (c *ComicClient) FetchBytes(ctx context.Context, link string, opts ...RequestOption) ([]byte, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	req, err := c.buildRequest(ctx, link, opts...)
	if err != nil {
		return nil, err
	}

	cache := c.ResponseCache()
	cacheKey := ""
	if cache != nil {
		cacheKey = MetadataCacheKey(req)
		if cacheKey != "" {
			if cached, ok := cache.Get(cacheKey); ok {
				// log.Printf("Cache hit for request %s (%s)\n", req.URL.String(), cacheKey)
				return cached, nil
			}
		}
	}

	resp, err := c.DoRaw(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("sites: failed to close response body for %s: %v", link, closeErr)
		}
	}()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("unexpected status code %d for %s", resp.StatusCode, link)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if cache != nil && cacheKey != "" {
		cache.Set(cacheKey, body)
	}
	return body, nil
}
