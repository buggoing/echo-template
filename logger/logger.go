package logger

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// New 返回一个新的日志实例。
func New(prefix string) echo.Logger {
	logger := log.New(prefix)
	InitLogger(logger)
	return logger
}

// InitLogger 用于初始化日志实例。
func InitLogger(logger echo.Logger) {
	var lvl log.Lvl
	logLevel := os.Getenv("LOG_SEVERITY_LEVEL")
	switch logLevel {
	case "ERROR", "error":
		lvl = log.ERROR
	case "WARNING", "warning":
		lvl = log.WARN
	// If env is unset, set level to INFO.
	case "INFO", "info":
		lvl = log.INFO
	default:
		lvl = log.DEBUG
	}
	logger.SetLevel(lvl)
	if l, ok := logger.(*log.Logger); ok {
		l.SetHeader("[${time_rfc3339}][${level}][${prefix}][${short_file}:${line}]")
	}
}
