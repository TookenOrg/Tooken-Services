package logger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
)

// Log is usable before Init runs. It used to be a nil pointer until main called
// Init, so any code that logged earlier - a spec parsed at startup, anything
// exercised by a test - brought the process down with a nil dereference inside
// slog rather than printing its message. A logging package is the last thing
// that should be able to crash a program, so it starts on the standard library
// default and Init swaps in the configured handler.
var Log = slog.Default()

func Init(isDevelopment bool) {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Raccourcir le chemin source
			if a.Key == slog.SourceKey {
				source := a.Value.Any().(*slog.Source)
				source.File = filepath.Base(source.File)
			}
			return a
		},
	}

	if isDevelopment {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	Log = slog.New(handler)
	slog.SetDefault(Log)
}

func LogTrace(msg string, args ...any) {
	logTraceEnabled := os.Getenv("LOG_TRACE")

	logTrace, _ := strconv.ParseBool(logTraceEnabled)
	if logTrace {
		Log.Debug(fmt.Sprintf(msg, args...))
	}
}

func LogDebug(msg string, args ...any) {
	Log.Debug(fmt.Sprintf(msg, args...))
}

func LogInfo(msg string, args ...any) {
	Log.Info(fmt.Sprintf(msg, args...))
}

func LogWarn(msg string, args ...any) {
	Log.Warn(fmt.Sprintf(msg, args...))
}

func LogError(msg string, args ...any) error {
	Log.Error(fmt.Sprintf(msg, args...))
	return fmt.Errorf(msg, args...)
}
