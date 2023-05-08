package config

import (
	"github.com/sirupsen/logrus"
)

type CustomLogger struct {
	*logrus.Entry
}

func NewLogger(instanceName string) *CustomLogger {
	logger := logrus.New()

	// Set logger formatting
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Set logger level
	logger.SetLevel(logrus.DebugLevel)

	// Create log entry with instance field
	entry := logger.WithFields(logrus.Fields{
		"instance": instanceName,
	})

	return &CustomLogger{entry}
}

func (l *CustomLogger) Infof(format string, args ...interface{}) {
	l.Entry.Infof(format, args...)
}

func (l *CustomLogger) Debugf(format string, args ...interface{}) {
	l.Entry.Debugf(format, args...)
}

func (l *CustomLogger) Warnf(format string, args ...interface{}) {
	l.Entry.Warnf(format, args...)
}

func (l *CustomLogger) Errorf(format string, args ...interface{}) {
	l.Entry.Errorf(format, args...)
}
