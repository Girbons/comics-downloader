package http

import (
	"net/http"
	"testing"
)

func TestMetadataCacheKeyDeterministic(t *testing.T) {
	reqA, err := http.NewRequest(http.MethodGet, "https://example.com/a?x=1", nil)
	if err != nil {
		t.Fatalf("failed creating request A: %v", err)
	}
	reqA.Header.Set("Accept", "text/html")
	reqA.Header.Set("User-Agent", "UA-1")
	reqA.Header.Set("Cookie", "cf=abc")

	reqB, err := http.NewRequest(http.MethodGet, "https://example.com/a?x=1", nil)
	if err != nil {
		t.Fatalf("failed creating request B: %v", err)
	}
	reqB.Header.Set("Cookie", "cf=abc")
	reqB.Header.Set("User-Agent", "UA-1")
	reqB.Header.Set("Accept", "text/html")

	keyA := MetadataCacheKey(reqA)
	keyB := MetadataCacheKey(reqB)
	if keyA == "" || keyB == "" {
		t.Fatalf("expected non-empty keys")
	}
	if keyA != keyB {
		t.Fatalf("expected deterministic keys, got %q and %q", keyA, keyB)
	}
}

func TestMetadataCacheKeyVariesByHeaders(t *testing.T) {
	reqA, _ := http.NewRequest(http.MethodGet, "https://example.com/a", nil)
	reqB, _ := http.NewRequest(http.MethodGet, "https://example.com/a", nil)
	reqA.Header.Set("Accept", "text/html")
	reqB.Header.Set("Accept", "application/json")

	if MetadataCacheKey(reqA) == MetadataCacheKey(reqB) {
		t.Fatalf("expected cache key to differ when request headers differ")
	}
}

func TestInMemoryResponseCacheStoresCopiesAndEvicts(t *testing.T) {
	cache := NewInMemoryResponseCache(1)
	original := []byte("value-1")
	cache.Set("a", original)
	original[0] = 'X'

	cached, ok := cache.Get("a")
	if !ok {
		t.Fatalf("expected cached value")
	}
	if string(cached) != "value-1" {
		t.Fatalf("expected cached copy, got %q", string(cached))
	}

	cached[0] = 'Y'
	reloaded, ok := cache.Get("a")
	if !ok {
		t.Fatalf("expected cached value on second get")
	}
	if string(reloaded) != "value-1" {
		t.Fatalf("expected cache to return immutable copy, got %q", string(reloaded))
	}

	cache.Set("b", []byte("value-2"))
	if _, ok := cache.Get("a"); ok {
		t.Fatalf("expected oldest key to be evicted")
	}
	if _, ok := cache.Get("b"); !ok {
		t.Fatalf("expected newest key to remain")
	}
}
