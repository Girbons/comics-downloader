package sites

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	httpclient "github.com/Girbons/comics-downloader/pkg/http"
)

func TestFetchHTMLUsesMetadataCache(t *testing.T) {
	var hits int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		fmt.Fprint(w, "<html>ok</html>")
	}))
	defer server.Close()

	client := httpclient.NewComicClient(
		httpclient.WithHTTPClient(server.Client()),
		httpclient.WithRetry(0, 0),
		httpclient.WithResponseCache(httpclient.NewInMemoryResponseCache(64)),
	)

	ctx := context.Background()
	first, err := fetchHTML(ctx, client, server.URL+"/metadata")
	if err != nil {
		t.Fatalf("first fetch failed: %v", err)
	}
	second, err := fetchHTML(ctx, client, server.URL+"/metadata")
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

	client := httpclient.NewComicClient(
		httpclient.WithHTTPClient(server.Client()),
		httpclient.WithRetry(0, 0),
		httpclient.WithResponseCache(httpclient.NewInMemoryResponseCache(64)),
	)

	ctx := context.Background()
	_, err := fetchHTML(ctx, client, server.URL+"/metadata", withHttpHeader("Accept", "text/html"))
	if err != nil {
		t.Fatalf("first fetch failed: %v", err)
	}
	_, err = fetchHTML(ctx, client, server.URL+"/metadata", withHttpHeader("Accept", "application/json"))
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

	client := httpclient.NewComicClient(
		httpclient.WithHTTPClient(server.Client()),
		httpclient.WithRetry(0, 0),
		httpclient.WithResponseCache(httpclient.NewInMemoryResponseCache(64)),
	)

	ctx := context.Background()
	_, err := fetchHTML(ctx, client, server.URL+"/metadata")
	if err == nil {
		t.Fatalf("expected first fetch to fail")
	}
	_, err = fetchHTML(ctx, client, server.URL+"/metadata")
	if err == nil {
		t.Fatalf("expected second fetch to fail")
	}

	if got := atomic.LoadInt32(&hits); got != 2 {
		t.Fatalf("expected two upstream hits because errors are not cached, got %d", got)
	}
}
