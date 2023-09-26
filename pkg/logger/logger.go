package logger

import (
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var (
	level   zerolog.Level
	zLog    *zerolog.Logger
	logOnce sync.Once
)

func init() {
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		level = zerolog.DebugLevel
	case "info":
		level = zerolog.InfoLevel
	case "warn":
		fallthrough
	case "warning":
		level = zerolog.WarnLevel
	case "err":
		fallthrough
	case "error":
		level = zerolog.ErrorLevel
	case "fatal":
		level = zerolog.FatalLevel
	default:
		level = zerolog.Disabled
	}
}

func Log() *zerolog.Logger {
	logOnce.Do(func() {
		writter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
		l := zerolog.New(writter).With().Timestamp().Logger().Level(level)
		zLog = &l
	})
	return zLog
}
