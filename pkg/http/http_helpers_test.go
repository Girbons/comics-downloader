package http

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestFetchHTMLUsesMetadataCache(t *testing.T) {
	var hits int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		fmt.Fprint(w, "<html>ok</html>")
	}))
	defer server.Close()

	client := NewComicClient(
		WithHTTPClient(server.Client()),
		WithRetry(0, 0),
		WithResponseCache(NewInMemoryResponseCache(64)),
	)

	ctx := context.Background()
	first, err := client.FetchHTML(ctx, server.URL+"/metadata")
	if err != nil {
		t.Fatalf("first fetch failed: %v", err)
	}
	second, err := client.FetchHTML(ctx, server.URL+"/metadata")
	if err != nil {
		t.Fatalf("second fetch failed: %v", err)
	}
	if first != second {
		t.Fatalf("expected identical response bodies")
	}
	if got := atomic.LoadInt32(&hits); got != 1 {
		t.Fatalf("expected one upstream hit, got %d", got)
	}
}

func TestFetchHTMLCacheKeyIncludesRequestHeaders(t *testing.T) {
	var hits int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		fmt.Fprint(w, "<html>ok</html>")
	}))
	defer server.Close()

	client := NewComicClient(
		WithHTTPClient(server.Client()),
		WithRetry(0, 0),
		WithResponseCache(NewInMemoryResponseCache(64)),
	)

	ctx := context.Background()
	_, err := client.FetchHTML(ctx, server.URL+"/metadata", WithHttpHeader("Accept", "text/html"))
	if err != nil {
		t.Fatalf("first fetch failed: %v", err)
	}
	_, err = client.FetchHTML(ctx, server.URL+"/metadata", WithHttpHeader("Accept", "application/json"))
	if err != nil {
		t.Fatalf("second fetch failed: %v", err)
	}

	if got := atomic.LoadInt32(&hits); got != 2 {
		t.Fatalf("expected two upstream hits for distinct headers, got %d", got)
	}
}

func TestFetchHTMLDoesNotCacheErrorResponses(t *testing.T) {
	var hits int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewComicClient(
		WithHTTPClient(server.Client()),
		WithRetry(0, 0),
		WithResponseCache(NewInMemoryResponseCache(64)),
	)

	ctx := context.Background()
	_, err := client.FetchHTML(ctx, server.URL+"/metadata")
	if err == nil {
		t.Fatalf("expected first fetch to fail")
	}
	_, err = client.FetchHTML(ctx, server.URL+"/metadata")
	if err == nil {
		t.Fatalf("expected second fetch to fail")
	}

	if got := atomic.LoadInt32(&hits); got != 2 {
		t.Fatalf("expected two upstream hits because errors are not cached, got %d", got)
	}
}
