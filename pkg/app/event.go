package app

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/bounoable/ical"
	"github.com/bounoable/ical/parse"
	"gorm.io/gorm"
)

const categoryProperty = "CATEGORIES"

type Event struct {
	gorm.Model
	Summary    string      `gorm:"not null;uniqueIndex:idx_event"`
	Start      time.Time   `gorm:"not null;uniqueIndex:idx_event"`
	End        time.Time   `gorm:"not null;uniqueIndex:idx_event"`
	Categories []string    `gorm:"serializer:json"`
	Event      parse.Event `gorm:"serializer:json"`

	HumanStart string `gorm:"-"`
	HumanEnd   string `gorm:"-"`
}

type Events []*Event

func (e *Events) CalculateDates() {
	for i := range *e {
		(*e)[i].CalculateDates()
	}
}

func (e *Event) CalculateDates() {
	e.HumanStart = niceDate(e.Start)
	e.HumanEnd = niceDate(e.End)
}

func niceDate(t time.Time) string {
	if t.Hour() == 0 && t.Minute() == 0 {
		return t.Format("2006-01-02")
	}

	return t.Format("2006-01-02 15:04")
}

func (e *Event) String() string {
	return fmt.Sprintf("[%03d] [%s] %s", e.ID, e.Start.Format("2006-01-02"), e.Summary)
}

func (e Events) ToCalendar() ([]byte, error) {
	c := ical.Calendar{}

	for i := range e {
		c.Events = append(c.Events, e[i].Event)
	}

	return ical.Marshal(c)
}

func (e *Event) AddCategory(cat string) {
	if cat == "" {
		return
	}

	if slices.Contains(e.Categories, cat) {
		return
	}

	e.UpdateCategories(append(e.Categories, cat))
}

func (e *Event) RemoveCategory(cat string) {
	newCats := []string{}

	for _, c := range e.Categories {
		if c != cat {
			newCats = append(newCats, c)
		}
	}

	e.UpdateCategories(newCats)
}

func (e *Event) UpdateCategories(cats []string) {
	e.Categories = cats
	propCats := strings.Join(cats, ",")

	if c, ok := e.Property(categoryProperty); ok {
		c.Value = propCats
		return
	}

	e.Event.Properties = append(e.Event.Properties, parse.Property{
		Name:  categoryProperty,
		Value: propCats,
	})
}

func (e *Event) Property(name string) (*parse.Property, bool) {
	for i, prop := range e.Event.Properties {
		if prop.Name == name {
			return &e.Event.Properties[i], true
		}
	}

	return nil, false
}
