package sites

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

func hostFromURL(link string) string {
	parsed, err := urlpkg.Parse(link)
	if err != nil {
		return ""
	}
	return parsed.Host
}

func buildRequest(ctx context.Context, client *httpclient.ComicClient, link string) (*http.Request, error) {
	req, err := client.PrepareRequest(link, hostFromURL(link))
	if err != nil {
		return nil, err
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	return req, nil
}

func fetchJSON(ctx context.Context, client *httpclient.ComicClient, link string, target interface{}) error {
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}

	data, err := fetchBytes(ctx, client, link)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return fmt.Errorf("empty response for %s", link)
	}
	return json.Unmarshal(data, target)
}

func fetchBytes(ctx context.Context, client *httpclient.ComicClient, link string) ([]byte, error) {
	cc := defaultClient(client)
	if ctx == nil {
		ctx = context.Background()
	}

	req, err := buildRequest(ctx, cc, link)
	if err != nil {
		return nil, err
	}

	resp, err := cc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("unexpected status code %d for %s", resp.StatusCode, link)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
