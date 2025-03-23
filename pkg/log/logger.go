package log

import (
	"io"
	"log/slog"
	"os"

	"github.com/nobbmaestro/lazyhis/pkg/config"
)

var logLevelMapping = map[config.LogLevel]slog.Level{
	config.LevelDebug: slog.LevelDebug,
	config.LevelInfo:  slog.LevelInfo,
	config.LevelWarn:  slog.LevelWarn,
	config.LevelError: slog.LevelError,
}

type Logger struct {
	Logger *slog.Logger
	file   *os.File
}

func NewLogger(cfg config.LogConfig) (*Logger, error) {
	var file *os.File
	var handler *slog.TextHandler

	opts := &slog.HandlerOptions{Level: logLevelMapping[cfg.LogLevel]}

	if cfg.LogEnabled {
		file, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}
		handler = slog.NewTextHandler(file, opts)
	} else {
		handler = slog.NewTextHandler(io.Discard, opts)
	}

	return &Logger{
		slog.New(handler),
		file,
	}, nil
}

func (l *Logger) Close() {
	if l.file != nil {
		l.file.Close()
	}
}
