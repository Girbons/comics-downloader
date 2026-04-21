package version

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func testHTTPClient(target *url.URL) *http.Client {
	return &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			cloned := req.Clone(req.Context())
			cloned.URL.Scheme = target.Scheme
			cloned.URL.Host = target.Host
			cloned.Host = target.Host
			return http.DefaultTransport.RoundTrip(cloned)
		}),
		Timeout: time.Second * 5,
	}
}

func TestIsNewAvailableCachesResult(t *testing.T) {
	ResetCache()
	originalTag := Tag
	Tag = "v0.1.0"
	defer func() {
		Tag = originalTag
		ResetCache()
	}()

	var hits int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{"tag_name":"v0.2.0","html_url":"http://example.com/latest"}]`)
	}))
	defer server.Close()

	serverURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("failed to parse test server url: %v", err)
	}

	client := testHTTPClient(serverURL)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	available, link, err := IsNewAvailable(ctx, client)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !available {
		t.Fatalf("expected new version to be available")
	}
	if link != "http://example.com/latest" {
		t.Fatalf("unexpected link: %s", link)
	}

	if got := atomic.LoadInt32(&hits); got != 1 {
		t.Fatalf("expected single HTTP call, got %d", got)
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second)
	defer cancel2()
	available, _, err = IsNewAvailable(ctx2, client)
	if err != nil {
		t.Fatalf("expected cached call to succeed, got %v", err)
	}
	if !available {
		t.Fatalf("expected cached call to preserve availability")
	}
	if got := atomic.LoadInt32(&hits); got != 1 {
		t.Fatalf("expected cached call to avoid HTTP, got %d requests", got)
	}
}

func TestIsNewAvailableInvalidSemver(t *testing.T) {
	ResetCache()
	originalTag := Tag
	Tag = "not-a-semver"
	defer func() {
		Tag = originalTag
		ResetCache()
	}()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{"tag_name":"v0.2.0","html_url":"http://example.com/latest"}]`)
	}))
	defer server.Close()

	serverURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("failed to parse test server url: %v", err)
	}

	client := testHTTPClient(serverURL)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	available, _, err := IsNewAvailable(ctx, client)
	if err == nil {
		t.Fatalf("expected error for invalid semver")
	}
	if available {
		t.Fatalf("expected availability to be false on error")
	}
}
