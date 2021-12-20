package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type Logger struct {
	logger *zap.Logger
}

func (l Logger) Error(v string) {
	l.logger.Error(v)
}

func (l Logger) Info(v string) {
	l.logger.Info(v)
}

func (l Logger) Infof(tmp string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(tmp, args...))
}

func NewLogger() Logger {
	logger, _ := zap.NewDevelopment()

	return Logger{
		logger: logger,
	}
}
