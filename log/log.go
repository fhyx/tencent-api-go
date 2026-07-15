package log

import (
	syslog "log"
	"log/slog"
)

type logger struct{}

// Default 默认实例
var Default Logger

func init() {
	// syslog.SetFlags(syslog.Ltime | syslog.Lshortfile)
	Default = &logger{}
}

// Deprecated: the package's default logger is now backed by slog.Default()
// via the log/slog standard library. To configure the process-wide logger,
// use slog.SetDefault instead. SetLogger only existed to inject custom
// loggers (e.g. zap.SugaredLogger) into the package-level default and is
// no longer the recommended way to control logging.
func SetLogger(logger Logger) {
	if logger != nil {
		Default = logger
	}
}

// Deprecated: callers should use slog.Default() directly to honor
// slog.SetDefault updates made elsewhere in the process. GetLogger only
// returns the init-time wrapper snapshot.
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
