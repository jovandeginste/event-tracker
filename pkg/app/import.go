package app

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"strings"

	"github.com/DusanKasan/parsemail"
	"github.com/bounoable/ical"
	_ "github.com/paulrosania/go-charset/data"
)

func (a *App) ImportICS(content string) []error {
	c, err := ical.ParseText(content)
	if err != nil {
		return []error{err}
	}

	var errs []error

	for _, e := range c.Events {
		cats := []string{}
		if c, ok := e.Property(categoryProperty); ok {
			cats = strings.Split(c.Value, ",")
		}

		n := &Event{
			Summary:    e.Summary,
			Start:      e.Start,
			End:        e.End,
			Event:      e,
			Categories: cats,
		}

		if err := a.FindEvent(n); err == nil {
			errs = append(errs, fmt.Errorf("skipping: %s (already exists)", n.String()))
			continue
		}

		if err := a.CreateEvent(n); err != nil {
			errs = append(errs, fmt.Errorf("could not create: %w", err))
		}
	}

	return errs
}

func (a *App) ImportEML(content []byte) []error {
	e, err := parsemail.Parse(bytes.NewReader(content))
	if err != nil {
		return []error{err}
	}

	var errs []error

	for _, p := range e.EmbeddedFiles {
		if err := a.importEMLFile(p.ContentType, p.Data); err != nil {
			errs = append(errs, err...)
			continue
		}
	}

	for _, p := range e.Attachments {
		if err := a.importEMLFile(p.ContentType, p.Data); err != nil {
			errs = append(errs, err...)
			continue
		}
	}

	return errs
}

func (a *App) importEMLFile(contentType string, data io.Reader) []error {
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return []error{err}
	}

	if mediatype != "text/calendar" {
		return nil
	}

	buf := new(strings.Builder)

	if _, err := io.Copy(buf, data); err != nil {
		return []error{err}
	}

	return a.ImportICS(buf.String())
}
