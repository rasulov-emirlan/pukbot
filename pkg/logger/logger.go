package logger

import "github.com/sirupsen/logrus"

type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type logger struct {
	ChatID string
}

func NewLogger() Logger {
	return logrus.New()
}
