package app

import (
	"io/fs"
	"log/slog"
	"time"

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
	working   bool
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

	go a.backgroundTasks()

	return a.echo.Start(a.Config.Bind)
}

func (a *App) backgroundTasks() {
	if a.working {
		return
	}

	a.working = true
	defer func() {
		a.working = false
	}()

	for {
		a.logger.Info("Running background tasks")

		e, err := a.AllEvents()
		if err != nil {
			a.logger.Error("Failed to fetch events: " + err.Error())
			continue
		}

		for _, ev := range e {
			if err := a.AddAITags(ev); err != nil {
				a.logger.Error("Failed to add tags to event: " + err.Error())
				continue
			}
		}

		time.Sleep(5 * time.Minute)
	}
}
