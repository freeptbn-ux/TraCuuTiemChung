package logger

import (
	"log/slog"
	"os"
	"sync"
)

var (
	Log  *slog.Logger
	once sync.Once
)

// InitLogger initializes the global slog logger with JSON handler for Production.
func InitLogger() {
	once.Do(func() {
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		Log = slog.New(handler)
		slog.SetDefault(Log)
	})
}
