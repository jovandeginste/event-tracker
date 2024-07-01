package app

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
	"gorm.io/gorm"
)

type App struct {
	Config Config

	db        *gorm.DB
	echo      *echo.Echo
	logger    *slog.Logger
	rawLogger *slog.Logger
}

func New() *App {
	return &App{}
}

func (a *App) Initialize() error {
	if err := a.Config.Initialize(); err != nil {
		return err
	}

	a.ConfigureLogger()

	if err := a.configureRoutes(); err != nil {
		return err
	}

	return a.initializeDatabase()
}

func (a *App) configureRoutes() error {
	a.echo = newEcho(a.logger)

	a.echo.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})
	a.echo.GET("/calendar", a.CalendarHandler)
	a.echo.GET("/events", a.EventsHandler)
	a.echo.GET("/events/:id", a.ShowEventHandler)
	a.echo.GET("/events/:id/category-add/:category", a.AddCategoryHandler)
	a.echo.GET("/events/:id/category-rm/:category", a.RemoveCategoryHandler)

	return nil
}

func (a *App) Serve() error {
	a.logger.Info("Starting web server on " + a.Config.Bind)

	return a.echo.Start(a.Config.Bind)
}

func newEcho(logger *slog.Logger) *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	e.Use(slogecho.New(logger.With("module", "webserver")))
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())
	e.Pre(middleware.RemoveTrailingSlash())

	return e
}
