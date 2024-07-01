package app

import (
	"slices"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/glebarez/sqlite"
	slogGorm "github.com/orandin/slog-gorm"
)

const thresholdSlowQueries = 100 * time.Millisecond

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

func (a *App) GetEvent(id uint) (*Event, error) {
	e := Event{}

	if err := a.DB().First(&e, id).Error; err != nil {
		return nil, err
	}

	return &e, nil
}

func (a *App) FindEvent(e *Event) error {
	return a.DB().Where(e).First(e).Error
}

func (a *App) AllEvents() (Events, error) {
	e := Events{}

	if err := a.DB().Find(&e).Order("start asc").Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (a *App) FilterEvents(cats Categories) (Events, error) {
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
