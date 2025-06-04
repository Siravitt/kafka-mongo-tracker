package logger

import (
	"log/slog"
	"os"
)

var logLevel = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

var defaultLogLevel = "error"
var LogLevel slog.Level

func init() {
	logLevelConf := os.Getenv("LOG_LEVEL")
	if logLevelConf == "" {
		logLevelConf = defaultLogLevel
	}

	level, ok := logLevel[logLevelConf]
	if !ok {
		level = logLevel[defaultLogLevel]
	}

	LogLevel = level
}
