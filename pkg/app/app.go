package app

import (
	"io/fs"
	"log/slog"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Version struct {
	BuildTime string
	Ref       string
	RefName   string
	RefType   string
	Sha       string
}

type App struct {
	Version Version
	Config  Config
	Assets  fs.FS

	db        *gorm.DB
	echo      *echo.Echo
	logger    *slog.Logger
	rawLogger *slog.Logger
}

func NewApp(version Version) *App {
	return &App{
		Version:   version,
		logger:    newLogger(false),
		rawLogger: newLogger(false),
	}
}

func (a *App) Initialize() error {
	if err := a.Config.Initialize(); err != nil {
		return err
	}

	a.ConfigureLogger()
	a.configureRoutes()

	return a.initializeDatabase()
}

func (a *App) Serve() error {
	a.logger.Info("Starting web server on " + a.Config.Bind)

	return a.echo.Start(a.Config.Bind)
}
