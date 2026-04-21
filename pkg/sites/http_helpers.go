package sites

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	urlpkg "net/url"

	httpclient "github.com/Girbons/comics-downloader/pkg/http"
)

func defaultClient(client *httpclient.ComicClient) *httpclient.ComicClient {
	if client != nil {
		return client
	}
	return httpclient.NewComicClient()
}

type requestConfig struct {
	request *http.Request
}

func (rc *requestConfig) apply(opts []requestOption) {
	for _, opt := range opts {
		opt(rc)
	}
}

type requestOption func(*requestConfig)

func withHttpHeader(key, value string) requestOption {
	return func(rc *requestConfig) {
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

func hostFromURL(link string) string {
	parsed, err := urlpkg.Parse(link)
	if err != nil {
		return ""
	}
	return parsed.Host
}

func buildRequest(ctx context.Context, client *httpclient.ComicClient, link string, opts ...requestOption) (*http.Request, error) {
	reqConfig := &requestConfig{}
	req, err := client.PrepareRequest(link, hostFromURL(link))
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

func fetchHTML(ctx context.Context, client *httpclient.ComicClient, link string, opts ...requestOption) (string, error) {
	response, err := fetchBytes(ctx, client, link, opts...)
	if err != nil {
		return "", err
	}
	return string(response), nil
}

func fetchJSON(ctx context.Context, client *httpclient.ComicClient, link string, target interface{}, opts ...requestOption) error {
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}

	data, err := fetchBytes(ctx, client, link, opts...)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return fmt.Errorf("empty response for %s", link)
	}
	return json.Unmarshal(data, target)
}

func fetchBytes(ctx context.Context, client *httpclient.ComicClient, link string, opts ...requestOption) ([]byte, error) {
	cc := defaultClient(client)
	if ctx == nil {
		ctx = context.Background()
	}

	req, err := buildRequest(ctx, cc, link, opts...)
	if err != nil {
		return nil, err
	}

	cache := cc.ResponseCache()
	cacheKey := ""
	if cache != nil {
		cacheKey = httpclient.MetadataCacheKey(req)
		if cacheKey != "" {
			if cached, ok := cache.Get(cacheKey); ok {
				// log.Printf("Cache hit for request %s (%s)\n", req.URL.String(), cacheKey)
				return cached, nil
			}
		}
	}

	resp, err := cc.Do(req)
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
