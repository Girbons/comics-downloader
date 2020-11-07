package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	l             *logrus.Logger
	bindToChannel bool
	Messages      chan string
}

func NewLogger(bindToChannel bool, messages chan string) *Logger {
	return &Logger{
		l:             logrus.New(),
		Messages:      messages,
		bindToChannel: bindToChannel,
	}
}

func (logger *Logger) SetLevel(level logrus.Level) {
	logger.l.SetLevel(level)
}

func (logger *Logger) sendToChannel(msg string) {
	if logger.bindToChannel {
		logger.Messages <- msg
	}
}

func (logger *Logger) Info(msg string) {
	logger.l.Info(msg)
	logger.sendToChannel(fmt.Sprintf("INFO: %s", msg))
}

func (logger *Logger) Debug(msg string) {
	logger.l.Debug(msg)
	logger.sendToChannel(fmt.Sprintf("DEBUG: %s", msg))
}

func (logger *Logger) Warning(msg string) {
	logger.l.Warning(msg)
	logger.sendToChannel(fmt.Sprintf("WARNING: %s", msg))
}

func (logger *Logger) Error(msg string) {
	logger.l.Error(msg)
	logger.sendToChannel(fmt.Sprintf("ERROR: %s", msg))
}
