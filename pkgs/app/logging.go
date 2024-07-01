package app

import (
	"log/slog"
	"os"

	"github.com/fsouza/slognil"
)

func newLogger(enabled bool) *slog.Logger {
	if !enabled {
		return slognil.NewLogger()
	}

	return slog.New(newLogHandler())
}

func newLogHandler() slog.Handler {
	return slog.NewJSONHandler(os.Stdout, nil)
}

func (a *App) ConfigureLogger() {
	logger := newLogger(a.Config.Logging).
		With("app", "event-tracker")

	a.rawLogger = logger
	a.logger = logger.With("module", "app")
}

func (a *App) Logger() *slog.Logger {
	return a.logger
}
