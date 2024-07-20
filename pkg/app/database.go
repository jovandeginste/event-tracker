package app

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/glebarez/sqlite"
	slogGorm "github.com/orandin/slog-gorm"
)

const thresholdSlowQueries = 100 * time.Millisecond

var filteredProperties = []string{
	"LOCATION", "SUMMARY", "DESCRIPTION", "ORGANIZER", "ATTENDEE",
}

type databaseConfig struct {
	File string
}

func (a *App) DB() *gorm.DB {
	return a.db.Preload(clause.Associations)
}

func (a *App) Migrate() error {
	return a.db.AutoMigrate(
		Event{},
	)
}

func (a *App) initializeDatabase() error {
	loggerOptions := []slogGorm.Option{
		slogGorm.WithHandler(a.Logger().With("module", "database").Handler()),
		slogGorm.WithSlowThreshold(thresholdSlowQueries),
	}

	if a.Config.Trace {
		loggerOptions = append(loggerOptions, slogGorm.WithTraceAll())
	}

	gormLogger := slogGorm.New(
		loggerOptions...,
	)

	db, err := gorm.Open(sqlite.Open(a.Config.Database.File), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return err
	}

	a.db = db

	return a.Migrate()
}

func (a *App) UpdateEvent(e *Event) error {
	return a.DB().Save(e).Error
}

func (a *App) CreateEvent(e *Event) error {
	return a.DB().Create(e).Error
}

func (a *App) DeleteEvent(e *Event) error {
	return a.DB().Unscoped().Delete(e).Error
}

func (a *App) GetEvent(id uint) (*Event, error) {
	e := Event{}

	if err := a.DB().First(&e, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("event %d not found", id)
		}

		return nil, err
	}

	return &e, nil
}

func (a *App) FindEvent(e *Event) error {
	qe := &Event{
		Summary: e.Summary,
		Start:   e.Start,
		End:     e.End,
	}

	return a.DB().Where(qe).First(e).Error
}

func (a *App) AllEvents() (Events, error) {
	e := Events{}

	if err := a.DB().Order("start asc").Find(&e).Error; err != nil {
		return nil, err
	}

	e.CalculateAttributes()

	return e, nil
}

func (a *App) FilterEvents(cats []string) (Events, error) {
	allEvents, err := a.AllEvents()
	if err != nil {
		return nil, err
	}

	slices.Sort(cats)
	cats = slices.Compact(cats)

	if len(cats) == 0 {
		return allEvents, nil
	}

	filteredEvents := Events{}

	for _, e := range allEvents {
		for _, c := range e.Categories {
			if slices.Contains(cats, c) {
				filteredEvents = append(filteredEvents, e)
			}
		}
	}

	return filteredEvents, nil
}

func (a *App) SearchEvents(t string) (Events, error) {
	allEvents, err := a.AllEvents()
	if err != nil {
		return nil, err
	}

	if t == "" {
		return allEvents, nil
	}

	t = strings.ToLower(t)

	filteredEvents := Events{}

	for _, e := range allEvents {
		if e.Matches(t) {
			filteredEvents = append(filteredEvents, e)
		}
	}

	return filteredEvents, nil
}

func (a *App) AllCategories() ([]string, error) {
	evs := Events{}
	if err := a.DB().Order("start asc").Find(&evs).Error; err != nil {
		return nil, err
	}

	var categories []string

	for _, e := range evs {
		categories = append(categories, e.Categories...)
	}

	slices.Sort(categories)
	categories = slices.Compact(categories)

	return categories, nil
}

func (e *Event) Matches(term string) bool {
	term = strings.ToLower(term)

	for _, p := range filteredProperties {
		if l, ok := e.Property(p); ok {
			if strings.Contains(strings.ToLower(l.Value), term) {
				return true
			}
		}
	}

	for _, c := range e.Categories {
		if strings.Contains(strings.ToLower(c), term) {
			return true
		}
	}

	return false
}
