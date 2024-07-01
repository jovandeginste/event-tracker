package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/bounoable/ical"
)

func (a *App) ImportDir(dir string) error {
	files, err := os.ReadDir("testfiles")
	if err != nil {
		return err
	}

	for _, f := range files {
		if !f.IsDir() && !strings.HasSuffix(f.Name(), ".ics") {
			continue
		}

		if err := a.Import("./testfiles/" + f.Name()); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) Import(filename string) error {
	a.Logger().Info(fmt.Sprintf("Parsing: %s", filename))

	c, err := ical.ParseFile(filename)
	if err != nil {
		return err
	}

	for _, e := range c.Events {
		n := &Event{
			Summary: e.Summary,
			Start:   e.Start,
			End:     e.End,
			Event:   e,
		}

		if err := a.FindEvent(n); err == nil {
			a.Logger().Debug(fmt.Sprintf("Skipping: %s (already exists)", n.String()))
			continue
		}

		if err := a.CreateEvent(n); err != nil {
			return err
		}
	}

	return nil
}
