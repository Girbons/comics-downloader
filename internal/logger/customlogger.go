package logger

import (
	"fmt"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

// Logger is the custom app logger that can optionally emit messages to a GUI channel.
type Logger struct {
	inner         *logrus.Logger
	bindToChannel bool
	messages      chan string
	channelMu     sync.RWMutex
}

// NewLogger returns a logger instance.
func NewLogger(bindToChannel bool, messages chan string) *Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
	})
	return &Logger{
		inner:         log,
		bindToChannel: bindToChannel && messages != nil,
		messages:      messages,
	}
}

// SetLevel sets the logger level.
func (logger *Logger) SetLevel(level logrus.Level) {
	logger.inner.SetLevel(level)
}

// Writer exposes the underlying logrus logger for advanced usage.
func (logger *Logger) Writer() *logrus.Logger {
	return logger.inner
}

func (logger *Logger) sendToChannel(level, msg string) {
	logger.channelMu.RLock()
	defer logger.channelMu.RUnlock()

	if !logger.bindToChannel || logger.messages == nil {
		return
	}

	formatted := fmt.Sprintf("%s: %s", strings.ToUpper(level), msg)
	select {
	case logger.messages <- formatted:
	default:
	}
}

// Debug logs at Debug level.
func (logger *Logger) Debug(msg string) {
	logger.inner.Debug(msg)
	logger.sendToChannel("DEBUG", msg)
}

// Debugf logs formatted entries at Debug level.
func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.Debug(fmt.Sprintf(format, args...))
}

// Info logs at Info level.
func (logger *Logger) Info(msg string) {
	logger.inner.Info(msg)
	logger.sendToChannel("INFO", msg)
}

// Infof logs formatted entries at Info level.
func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.Info(fmt.Sprintf(format, args...))
}

// Warning logs at Warning level.
func (logger *Logger) Warning(msg string) {
	logger.inner.Warning(msg)
	logger.sendToChannel("WARNING", msg)
}

// Warningf logs formatted entries at Warning level.
func (logger *Logger) Warningf(format string, args ...interface{}) {
	logger.Warning(fmt.Sprintf(format, args...))
}

// Error logs at Error level.
func (logger *Logger) Error(msg string) {
	logger.inner.Error(msg)
	logger.sendToChannel("ERROR", msg)
}

// Errorf logs formatted entries at Error level.
func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.Error(fmt.Sprintf(format, args...))
}
