package log

import (
	"go.uber.org/zap"
)

const (
	LevelDebug = "DEBUG"
	LevelInfo  = "INFO"
	LevelWarn  = "WARN"
	LevelError = "ERROR"
	LevelFatal = "FATAL"
)

type Logger interface {
	Error(string, ...zap.Field)
	Debug(string, ...zap.Field)
	Info(string, ...zap.Field)
	Warn(string, ...zap.Field)
	Fatal(string, ...zap.Field)
}

func New(refID string, appName string, level string) Logger {
	logger := Init(refID, appName, ZapLevel(level))
	return logger
}
