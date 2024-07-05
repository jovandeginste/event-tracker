package app

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path"

	"github.com/labstack/echo/v4"
)

func uploadedFile(file *multipart.FileHeader) ([]byte, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Read all from r into a bytes slice
	content, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (a *App) addEventFromFile(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	files := form.File["file"]

	var resp APIResponse

	for _, file := range files {
		content, parseErr := uploadedFile(file)
		if parseErr != nil {
			resp.AddError(parseErr)
			continue
		}

		ext := path.Ext(file.Filename)

		switch ext {
		case ".ics":
			if err := a.ImportICS(string(content)); err != nil {
				resp.AddErrors(err)
			} else {
				resp.AddNotification("Successfully added file '" + file.Filename + "'")
			}
		case ".eml":
			if err := a.ImportEML(content); err != nil {
				resp.AddErrors(err)
			} else {
				resp.AddNotification("Successfully added file '" + file.Filename + "'")
			}
		default:
			resp.AddError(fmt.Errorf("unknown file extension: %v", ext))
		}
	}

	resp.ParseErrors()

	c.Response().Header().Set("HX-Trigger", "events-updated")

	return c.JSON(http.StatusOK, resp)
}
