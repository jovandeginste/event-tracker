package app

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/bounoable/ical"
	"github.com/bounoable/ical/parse"
	"github.com/dustin/go-humanize"
	"gorm.io/gorm"
)

const categoryProperty = "CATEGORIES"

type Event struct {
	gorm.Model
	Summary      string      `gorm:"not null;uniqueIndex:idx_event"`
	Start        time.Time   `gorm:"not null;uniqueIndex:idx_event"`
	End          time.Time   `gorm:"not null;uniqueIndex:idx_event"`
	Categories   []string    `gorm:"serializer:json"`
	AICategories []string    `gorm:"serializer:json"`
	Event        parse.Event `gorm:"serializer:json"`

	HumanStart     string `gorm:"-"`
	HumanEnd       string `gorm:"-"`
	HumanTimeRange string `gorm:"-"`
	TimeRange      string `gorm:"-"`
}

type Events []*Event

func (e *Event) Description() string {
	if e == nil {
		return ""
	}

	return e.Event.Description
}

func (e *Event) Location() string {
	if e == nil {
		return ""
	}

	if r, ok := e.Event.Property("LOCATION"); ok {
		return r.Value
	}

	return ""
}

func (e *Events) CalculateAttributes() {
	for i := range *e {
		(*e)[i].CalculateAttributes()
	}
}

func (e *Event) CalculateAttributes() {
	e.CalculateDates()

	e.AICategories = slices.DeleteFunc(e.AICategories, func(c string) bool {
		return slices.Contains(e.Categories, c)
	})
}

func (e *Event) CalculateDates() {
	e.HumanStart = humanize.Time(e.Start)
	e.HumanEnd = humanize.Time(e.End)

	if e.HumanEnd == e.HumanStart ||
		e.End.Sub(e.Start).Hours() == 24 {
		e.HumanTimeRange = e.HumanStart
	} else {
		e.HumanTimeRange = fmt.Sprintf("%s - %s", e.HumanStart, e.HumanEnd)
	}

	if e.End.Sub(e.Start).Hours() == 24 {
		e.TimeRange = niceDate(e.Start)
	} else {
		e.TimeRange = fmt.Sprintf("%s - %s", niceDate(e.Start), niceDate(e.End))
	}
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
