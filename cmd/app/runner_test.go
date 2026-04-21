package app

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/Girbons/comics-downloader/pkg/config"
	httpclient "github.com/Girbons/comics-downloader/pkg/http"
)

func TestRunnerPrepareOptionsProvidesDependencies(t *testing.T) {
	runner := NewRunner(config.Options{})

	opts := runner.prepareOptions()

	if opts.Logger == nil {
		t.Fatalf("expected logger to be initialized")
	}

	if opts.Client == nil {
		t.Fatalf("expected http client to be initialized")
	}

	if runner.base.Logger != nil {
		t.Fatalf("expected runner base logger to remain nil")
	}
	if runner.base.Client != nil {
		t.Fatalf("expected runner base client to remain nil")
	}
}

func TestRunnerRunRequiresURL(t *testing.T) {
	msgs := make(chan string, 1)

	runner := NewRunner(config.Options{})
	runner.WithChannelBinding(msgs)

	runner.Run()

	select {
	case msg := <-msgs:
		if !strings.Contains(msg, "url parameter is required") {
			t.Fatalf("expected error message about missing url, got %q", msg)
		}
	default:
		t.Fatalf("expected an error message to be sent")
	}
}

func TestBuildClientOptionsCacheToggle(t *testing.T) {
	clientWithCache := httpclient.NewComicClient(buildClientOptions(config.Options{})...)
	if clientWithCache.ResponseCache() == nil {
		t.Fatalf("expected response cache to be enabled by default")
	}

	clientWithoutCache := httpclient.NewComicClient(buildClientOptions(config.Options{NoCache: true})...)
	if clientWithoutCache.ResponseCache() != nil {
		t.Fatalf("expected response cache to be disabled when no-cache is set")
	}
}

func TestBuildClientOptionsProxy(t *testing.T) {
	client := httpclient.NewComicClient(buildClientOptions(config.Options{HTTPProxy: "http://127.0.0.1:8080"})...)

	transport, ok := client.HTTPClient().Transport.(*http.Transport)
	if !ok {
		t.Fatalf("expected *http.Transport, got %T", client.HTTPClient().Transport)
	}
	if transport.Proxy == nil {
		t.Fatalf("expected proxy function to be configured")
	}

	reqURL, err := url.Parse("https://example.com")
	if err != nil {
		t.Fatalf("failed to parse test URL: %v", err)
	}

	proxyURL, err := transport.Proxy(&http.Request{URL: reqURL})
	if err != nil {
		t.Fatalf("unexpected proxy resolution error: %v", err)
	}
	if proxyURL == nil {
		t.Fatalf("expected resolved proxy URL")
	}
	if proxyURL.String() != "http://127.0.0.1:8080" {
		t.Fatalf("expected proxy URL %q, got %q", "http://127.0.0.1:8080", proxyURL.String())
	}
}

func TestProxyHTTPClientRejectsInvalidValue(t *testing.T) {
	client, ok := proxyHTTPClient("not-a-url")
	if ok {
		t.Fatalf("expected invalid proxy to be rejected")
	}
	if client != nil {
		t.Fatalf("expected nil client for invalid proxy")
	}
}
