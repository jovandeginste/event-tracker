package app

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
)

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

func (a *App) configureRoutes() {
	e := newEcho(a.logger)

	publicGroup := e.Group("")
	publicGroup.StaticFS("/", a.Assets)

	publicGroup.GET("/version", a.VersionHandler).Name = "version"

	feedsGroup := e.Group("/feed")
	feedsGroup.GET("/calendar", a.CalendarHandler).Name = "feed-calendar"

	eventsGroup := e.Group("/events")
	eventsGroup.GET("", a.JSONHandler).Name = "events-json"
	eventsGroup.POST("", a.AddEventsHandler).Name = "events-create"
	eventsGroup.GET("/search", a.SearchEventsHandler).Name = "events-search"
	eventsGroup.GET("/:id/json", a.ShowEventHandler).Name = "event-show"
	eventsGroup.GET("/:id/ical", a.ShowIcalEventHandler).Name = "event-ical"
	eventsGroup.DELETE("/:id", a.DeleteEventHandler).Name = "event-delete"
	eventsGroup.POST("/:id/category-add", a.AddCategoryHandler).Name = "event-category-add"
	eventsGroup.POST("/:id/category-rm/:category", a.RemoveCategoryHandler).Name = "event-category-rm"

	a.echo = e
}
