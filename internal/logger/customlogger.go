package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Logger is the custom app logger.
type Logger struct {
	l             *logrus.Logger
	bindToChannel bool
	Messages      chan string
}

// NewLogger returns a logger instance
func NewLogger(bindToChannel bool, messages chan string) *Logger {
	return &Logger{
		l:             logrus.New(),
		Messages:      messages,
		bindToChannel: bindToChannel,
	}
}

// SetLevel set logger level.
func (logger *Logger) SetLevel(level logrus.Level) {
	logger.l.SetLevel(level)
}

func (logger *Logger) sendToChannel(msg string) {
	if logger.bindToChannel {
		logger.Messages <- msg
	}
}

// Info logs info level log.
func (logger *Logger) Info(msg string) {
	logger.l.Info(msg)
	logger.sendToChannel(fmt.Sprintf("INFO: %s", msg))
}

// Debug logs debug level log.
func (logger *Logger) Debug(msg string) {
	logger.l.Debug(msg)
	logger.sendToChannel(fmt.Sprintf("DEBUG: %s", msg))
}

// Warning logs Warning level log.
func (logger *Logger) Warning(msg string) {
	logger.l.Warning(msg)
	logger.sendToChannel(fmt.Sprintf("WARNING: %s", msg))
}

// Error logs error level log.
func (logger *Logger) Error(msg string) {
	logger.l.Error(msg)
	logger.sendToChannel(fmt.Sprintf("ERROR: %s", msg))
}
