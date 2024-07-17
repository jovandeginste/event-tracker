package app

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func (a *App) VersionHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, a.Version)
}

func (a *App) SearchEventsHandler(c echo.Context) error {
	var resp APIResponse

	term := c.QueryParam("term")

	events, err := a.SearchEvents(term)
	if err != nil {
		resp.AddError(err)
	}

	resp.Results = events

	resp.ParseErrors()

	return c.JSON(http.StatusOK, resp)
}

func (a *App) EventsHandler(c echo.Context) error {
	type CalendarEvent struct {
		Title string    `json:"title"`
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	}

	events, err := a.AllEvents()
	if err != nil {
		return err
	}

	resp := []CalendarEvent{}

	for _, e := range events {
		resp = append(resp, CalendarEvent{
			Title: e.Summary,
			Start: e.Start,
			End:   e.End,
		})
	}

	return c.JSON(http.StatusOK, resp)
}

func (a *App) JSONHandler(c echo.Context) error {
	var (
		categories []string
		resp       APIResponse
	)

	if err := echo.QueryParamsBinder(c).
		Strings("categories", &categories).
		BindError(); err != nil {
		resp.AddError(err)
	}

	events, err := a.FilterEvents(categories)
	if err != nil {
		resp.AddError(err)
	}

	resp.Results = events

	resp.ParseErrors()

	return c.JSON(http.StatusOK, resp)
}

func (a *App) CalendarHandler(c echo.Context) error {
	var categories []string

	if err := echo.QueryParamsBinder(c).
		Strings("categories", &categories).
		BindError(); err != nil {
		return err
	}

	events, err := a.FilterEvents(categories)
	if err != nil {
		return err
	}

	resp, err := events.ToCalendar()
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, string(resp))
}

func (a *App) ShowIcalEventHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	event, err := a.GetEvent(uint(id))
	if err != nil {
		return err
	}

	events := Events{event}

	resp, err := events.ToCalendar()
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, string(resp))
}

func (a *App) AddEventsHandler(c echo.Context) error {
	return a.addEventFromFile(c)
}

func (a *App) DeleteEventHandler(c echo.Context) error {
	var resp APIResponse

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	event, err := a.GetEvent(uint(id))
	if err != nil {
		resp.AddError(err)
	} else {
		if err := a.DeleteEvent(event); err != nil {
			resp.AddError(err)
		} else {
			resp.AddNotification("Successfully deleted event '" + event.Summary + "'")
		}
	}

	c.Response().Header().Set("HX-Trigger", "events-updated")

	resp.Results = event

	return c.JSON(http.StatusOK, resp)
}

func (a *App) ShowEventHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	event, err := a.GetEvent(uint(id))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, event)
}

func (a *App) AddCategoryHandler(c echo.Context) error {
	var (
		resp  APIResponse
		event *Event
	)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		resp.AddError(err)
	}

	cat := c.FormValue("category")

	if resp.NoErrors() {
		event, err = a.GetEvent(uint(id))
		if err != nil {
			resp.AddError(err)
		}

		event.AddCategory(cat)
	}

	if resp.NoErrors() {
		if err := a.UpdateEvent(event); err != nil {
			resp.AddError(err)
		}
	}

	if resp.NoErrors() {
		resp.Results = event
		resp.AddNotification(fmt.Sprintf("Added category '%s' to event '%s'", c.FormValue("category"), event.Summary))
	}

	resp.ParseErrors()

	c.Response().Header().Set("HX-Trigger", "events-updated")

	return c.JSON(http.StatusOK, resp)
}

func (a *App) RemoveCategoryHandler(c echo.Context) error {
	var (
		resp  APIResponse
		event *Event
	)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		resp.AddError(err)
	}

	if resp.NoErrors() {
		event, err = a.GetEvent(uint(id))
		if err != nil {
			resp.AddError(err)
		}

		event.RemoveCategory(c.Param("category"))
	}

	if resp.NoErrors() {
		if err := a.UpdateEvent(event); err != nil {
			resp.AddError(err)
		}
	}

	if resp.NoErrors() {
		resp.Results = event
		resp.AddNotification(fmt.Sprintf("Removed category '%s' from event '%s'", c.Param("category"), event.Summary))
	}

	resp.ParseErrors()

	c.Response().Header().Set("HX-Trigger", "events-updated")

	return c.JSON(http.StatusOK, resp)
}
