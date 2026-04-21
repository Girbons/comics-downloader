package http

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"sort"
	"strings"
	"sync"
)

const defaultResponseCacheCapacity = 2048

// ResponseCache stores response bodies for metadata requests.
type ResponseCache interface {
	Get(key string) ([]byte, bool)
	Set(key string, value []byte)
	Clear()
}

// InMemoryResponseCache is a bounded in-memory cache implementation.
type InMemoryResponseCache struct {
	mu       sync.RWMutex
	maxItems int
	items    map[string][]byte
	order    []string
}

// NewInMemoryResponseCache creates a process-local in-memory response cache.
func NewInMemoryResponseCache(maxItems int) *InMemoryResponseCache {
	if maxItems <= 0 {
		maxItems = defaultResponseCacheCapacity
	}

	return &InMemoryResponseCache{
		maxItems: maxItems,
		items:    make(map[string][]byte),
		order:    make([]string, 0, maxItems),
	}
}

// Get returns a copy of the cached value when present.
func (c *InMemoryResponseCache) Get(key string) ([]byte, bool) {
	if c == nil || key == "" {
		return nil, false
	}

	c.mu.RLock()
	value, ok := c.items[key]
	c.mu.RUnlock()
	if !ok {
		return nil, false
	}

	out := make([]byte, len(value))
	copy(out, value)
	return out, true
}

// Set stores a copy of the value in the cache.
func (c *InMemoryResponseCache) Set(key string, value []byte) {
	if c == nil || key == "" || len(value) == 0 {
		return
	}

	cached := make([]byte, len(value))
	copy(cached, value)

	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.items[key]; !exists {
		c.order = append(c.order, key)
	}
	c.items[key] = cached

	for len(c.order) > c.maxItems {
		oldest := c.order[0]
		c.order = c.order[1:]
		delete(c.items, oldest)
	}
}

// Clear removes all cache entries.
func (c *InMemoryResponseCache) Clear() {
	if c == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string][]byte)
	c.order = c.order[:0]
}

// MetadataCacheKey builds a deterministic cache key for metadata GET requests.
func MetadataCacheKey(req *http.Request) string {
	if req == nil || req.URL == nil {
		return ""
	}

	builder := strings.Builder{}
	builder.WriteString(req.Method)
	builder.WriteString("\n")
	builder.WriteString(req.URL.String())
	builder.WriteString("\n")

	headers := []string{"Accept", "Cookie", "Referer", "User-Agent"}
	sort.Strings(headers)
	for _, key := range headers {
		values := req.Header.Values(key)
		if len(values) == 0 {
			continue
		}
		copied := append([]string(nil), values...)
		sort.Strings(copied)
		builder.WriteString(strings.ToLower(key))
		builder.WriteString(":")
		builder.WriteString(strings.Join(copied, ","))
		builder.WriteString("\n")
	}

	sum := sha256.Sum256([]byte(builder.String()))
	return hex.EncodeToString(sum[:])
}
