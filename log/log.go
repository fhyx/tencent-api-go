package log

import (
	syslog "log"
	"log/slog"
)

type logger struct{}

// Default 默认实例
var Default Logger

func init() {
	syslog.SetFlags(syslog.Ltime | syslog.Lshortfile)
	Default = &logger{}
}

func SetLogger(logger Logger) {
	if logger != nil {
		Default = logger
	}
}

func GetLogger() Logger {
	return Default
}

func (z *logger) Debugw(msg string, keysAndValues ...interface{}) {
	slog.Debug(msg, keysAndValues...)
}

func (z *logger) Infow(msg string, keysAndValues ...interface{}) {
	slog.Info(msg, keysAndValues...)
}

func (z *logger) Warnw(msg string, keysAndValues ...interface{}) {
	slog.Warn(msg, keysAndValues...)
}

func (z *logger) Errorw(msg string, keysAndValues ...interface{}) {
	slog.Error(msg, keysAndValues...)
}

func (z *logger) Fatalw(msg string, keysAndValues ...interface{}) {
	syslog.Fatal(msg, keysAndValues)
}
