package logger

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestInfofSendsToChannel(t *testing.T) {
	ch := make(chan string, 1)
	log := NewLogger(true, ch)
	log.SetLevel(logrus.InfoLevel)

	log.Infof("downloaded %d files", 3)

	select {
	case msg := <-ch:
		if msg != "INFO: downloaded 3 files" {
			t.Fatalf("unexpected message %q", msg)
		}
	default:
		t.Fatalf("expected message to be published")
	}
}

func TestChannelSendDoesNotBlock(t *testing.T) {
	ch := make(chan string, 1)
	log := NewLogger(true, ch)

	log.Info("first")
	log.Info("second") // should not block even if buffer is full

	select {
	case <-ch:
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("expected to dequeue at least one message")
	}
}

func TestLoggerHandlesNilChannel(t *testing.T) {
	log := NewLogger(false, nil)
	log.Debug("noop")
	log.Infof("hello %s", "world")
	log.Errorf("error %s", "msg")
}
