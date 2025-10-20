package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type stubLimiter struct {
	count int32
}

func (s *stubLimiter) Wait(ctx context.Context) error {
	atomic.AddInt32(&s.count, 1)
	return nil
}

func TestPrepareRequestMangakakalot(t *testing.T) {
	cc := NewComicClient()
	link := "http://mangakakalot.com"
	source := "mangakakalot.com"
	req, err := cc.PrepareRequest(link, source)

	require.NoError(t, err)
	require.Equal(t, link, req.Header.Get("Referer"))
	require.Equal(t, defaultUserAgent, req.Header.Get("User-Agent"))
}

func TestPrepareRequestGenericHost(t *testing.T) {
	cc := NewComicClient()
	link := "http://foo.com"
	req, err := cc.PrepareRequest(link, "foo")

	require.NoError(t, err)
	require.Empty(t, req.Header.Values("Referer"))
	require.Equal(t, defaultUserAgent, req.Header.Get("User-Agent"))
}

func TestGetRetriesOnServerError(t *testing.T) {
	var hits int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		call := atomic.AddInt32(&hits, 1)
		if call == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewComicClient(
		WithHTTPClient(server.Client()),
		WithRetry(1, 0),
	)

	resp, err := client.Get(server.URL, "example.com")
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, int32(2), atomic.LoadInt32(&hits))
}

func TestRateLimiterInvoked(t *testing.T) {
	limiter := &stubLimiter{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewComicClient(
		WithHTTPClient(server.Client()),
		WithRetry(0, 0),
		WithRateLimiter(limiter),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req, err := client.PrepareRequest(server.URL, "example.com")
	require.NoError(t, err)

	req = req.WithContext(ctx)
	_, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, int32(1), atomic.LoadInt32(&limiter.count))
}

func TestUserAgentRotation(t *testing.T) {
	agents := []string{"UA-1", "UA-2"}
	// we need to ensure rotation occurs across calls
	client := NewComicClient(WithUserAgents(agents))

	req1, err := client.PrepareRequest("http://example.com", "example.com")
	require.NoError(t, err)
	req2, err := client.PrepareRequest("http://example.com", "example.com")
	require.NoError(t, err)

	if req1.Header.Get("User-Agent") == req2.Header.Get("User-Agent") {
		t.Fatalf("expected rotating user agents, got identical headers %q", req1.Header.Get("User-Agent"))
	}

	req3, err := client.PrepareRequest("http://example.com", "example.com")
	require.NoError(t, err)

	require.Equal(t, req1.Header.Get("User-Agent"), req3.Header.Get("User-Agent"), "rotation should loop back to first entry")
}

func TestAdditionalHeadersApplied(t *testing.T) {
	client := NewComicClient(WithHeaders(map[string]string{
		"Cookie":      "cf_clearance=abc",
		"X-Custom-Id": "123",
	}))

	req, err := client.PrepareRequest("http://example.com", "example.com")
	require.NoError(t, err)

	require.Equal(t, "cf_clearance=abc", req.Header.Get("Cookie"))
	require.Equal(t, "123", req.Header.Get("X-Custom-Id"))
}
