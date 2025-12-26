package logger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

var Log *slog.Logger

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

func LogDebug(msg string, args ...any) {
	Log.Debug(fmt.Sprintf(msg, args...))
}

func LogInfo(msg string, args ...any) {
	Log.Info(fmt.Sprintf(msg, args...))
}

func LogWarn(msg string, args ...any) {
	Log.Warn(fmt.Sprintf(msg, args...))
}

func LogError(msg string, args ...any) {
	Log.Error(fmt.Sprintf(msg, args...))
}
