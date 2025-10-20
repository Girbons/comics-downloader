package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	defaultTimeout    = 15 * time.Second
	defaultRetryCount = 2
	defaultRetryWait  = 500 * time.Millisecond
	defaultUserAgent  = "comics-downloader-client"
)

// RateLimiter exposes a minimal interface for throttling outgoing requests.
type RateLimiter interface {
	Wait(context.Context) error
}

// Option modifies ComicClient behaviour.
type Option func(*ComicClient)

// WithHTTPClient allows supplying a custom http.Client.
func WithHTTPClient(client *http.Client) Option {
	return func(cc *ComicClient) {
		if client != nil {
			cc.client = client
		}
	}
}

// WithRetry configures retry attempts and wait duration between retries.
func WithRetry(count int, wait time.Duration) Option {
	return func(cc *ComicClient) {
		if count >= 0 {
			cc.retryCount = count
		}
		if wait >= 0 {
			cc.retryWait = wait
		}
	}
}

// WithRateLimiter sets a rate limiter to gate outbound requests.
func WithRateLimiter(limiter RateLimiter) Option {
	return func(cc *ComicClient) {
		cc.rateLimiter = limiter
	}
}

// WithUserAgent overrides the default user-agent header.
func WithUserAgent(agent string) Option {
	return func(cc *ComicClient) {
		if strings.TrimSpace(agent) != "" {
			cc.userAgent = agent
		}
	}
}

// ComicClient is the custom HTTP helper used across the downloader.
type ComicClient struct {
	client      *http.Client
	retryCount  int
	retryWait   time.Duration
	rateLimiter RateLimiter
	userAgent   string
}

// NewComicClient returns a ComicClient instance with sane defaults.
func NewComicClient(options ...Option) *ComicClient {
	cc := &ComicClient{
		client: &http.Client{
			Timeout: defaultTimeout,
		},
		retryCount: defaultRetryCount,
		retryWait:  defaultRetryWait,
		userAgent:  defaultUserAgent,
	}

	for _, opt := range options {
		opt(cc)
	}

	if cc.client.Timeout == 0 {
		cc.client.Timeout = defaultTimeout
	}

	return cc
}

// HTTPClient exposes the underlying http.Client instance.
func (c *ComicClient) HTTPClient() *http.Client {
	if c == nil {
		return nil
	}
	return c.client
}

// PrepareRequest setup a `GET` request with custom headers.
func (c *ComicClient) PrepareRequest(link, hostname string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return nil, err
	}

	if strings.Contains(hostname, "manganato") || strings.Contains(hostname, "mangakakalot") {
		req.Header.Set("Referer", link)
	}

	if c != nil && c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	return req, nil
}

// Do executes an HTTP request applying retry, timeout, and rate limiting policies.
func (c *ComicClient) Do(req *http.Request) (*http.Response, error) {
	if c == nil {
		return nil, errors.New("comic client is nil")
	}
	if req == nil {
		return nil, errors.New("request is nil")
	}

	attempts := c.retryCount + 1
	if attempts < 1 {
		attempts = 1
	}

	var lastErr error
	for attempt := 0; attempt < attempts; attempt++ {
		if attempt > 0 && c.retryWait > 0 {
			select {
			case <-time.After(c.retryWait):
			case <-req.Context().Done():
				return nil, req.Context().Err()
			}
		}

		if err := c.wait(req.Context()); err != nil {
			return nil, err
		}

		currentReq := req
		if attempt > 0 {
			currentReq = req.Clone(req.Context())
		}

		resp, err := c.client.Do(currentReq)
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("server error: %d", resp.StatusCode)
			resp.Body.Close()
			continue
		}

		return resp, nil
	}

	return nil, lastErr
}

// GET performs a GET request applying the configured policies.
func (c *ComicClient) Get(link, hostname string) (*http.Response, error) {
	request, err := c.PrepareRequest(link, hostname)
	if err != nil {
		return nil, err
	}

	return c.Do(request)
}

func (c *ComicClient) wait(ctx context.Context) error {
	if c.rateLimiter == nil {
		return nil
	}
	return c.rateLimiter.Wait(ctx)
}
