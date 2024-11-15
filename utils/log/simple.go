package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Simple struct {
	refID   string
	appName string
	*zap.Logger
}

func Init(refID string, appName string, level zapcore.Level) *Simple {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Level.SetLevel(level)
	logger, _ := loggerConfig.Build()
	logger = logger.With(zap.String("ref_id", refID), zap.String("app_name", appName))
	return &Simple{Logger: logger, refID: refID, appName: appName}
}

func ZapLevel(level string) zapcore.Level {
	var lvl zapcore.Level
	switch level {
	case LevelDebug:
		lvl = zap.DebugLevel
	case LevelInfo:
		lvl = zap.InfoLevel
	case LevelWarn:
		lvl = zap.WarnLevel
	case LevelError:
		lvl = zap.ErrorLevel
	case LevelFatal:
		lvl = zap.FatalLevel
	default:
		lvl = zap.InfoLevel
	}

	return lvl
}
