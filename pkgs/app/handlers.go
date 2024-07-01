package app

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type APIResponse struct {
	Errors  []string    `json:"errors"`
	Results interface{} `json:"results"`
}

func (a *App) EventsHandler(c echo.Context) error {
	var (
		resp          APIResponse
		strCategories []string
		categories    Categories
	)

	if err := echo.QueryParamsBinder(c).
		Strings("categories", &strCategories).
		BindError(); err != nil {
		resp.Errors = append(resp.Errors, err.Error())
	}

	for _, c := range strCategories {
		categories = append(categories, Category(c))
	}

	events, err := a.FilterEvents(categories)
	if err != nil {
		resp.Errors = append(resp.Errors, err.Error())
	}

	resp.Results = events

	return c.JSON(http.StatusOK, resp)
}

func (a *App) CalendarHandler(c echo.Context) error {
	var (
		strCategories []string
		categories    Categories
	)

	if err := echo.QueryParamsBinder(c).
		Strings("categories", &strCategories).
		BindError(); err != nil {
		return err
	}

	for _, c := range strCategories {
		categories = append(categories, Category(c))
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
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	event, err := a.GetEvent(uint(id))
	if err != nil {
		return err
	}

	category := Category(c.Param("category"))

	event.Categories = append(event.Categories, category)
	if err := a.UpdateEvent(event); err != nil {
		return err
	}

	return c.String(http.StatusOK, "OK")
}

func (a *App) RemoveCategoryHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	event, err := a.GetEvent(uint(id))
	if err != nil {
		return err
	}

	category := Category(c.Param("category"))
	newCats := Categories{}

	for _, c := range event.Categories {
		if c != category {
			newCats = append(newCats, c)
		}
	}

	event.Categories = newCats
	if err := a.UpdateEvent(event); err != nil {
		return err
	}

	return c.String(http.StatusOK, "OK")
}
