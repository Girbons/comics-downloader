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
