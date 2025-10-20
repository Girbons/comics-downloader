package app

import (
	"strings"
	"testing"

	"github.com/Girbons/comics-downloader/pkg/config"
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
