package app

import (
	"fmt"
	"time"

	"github.com/bounoable/ical"
	"github.com/bounoable/ical/parse"
	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Summary    string      `gorm:"not null;uniqueIndex:idx_event"`
	Start      time.Time   `gorm:"not null;uniqueIndex:idx_event"`
	End        time.Time   `gorm:"not null;uniqueIndex:idx_event"`
	Event      parse.Event `gorm:"serializer:json"`
	Categories Categories  `gorm:"serializer:json"`
}

type Events []*Event

func (e *Event) String() string {
	return fmt.Sprintf("[%03d] [%s] %s", e.ID, e.Start.Format("2006-01-02"), e.Summary)
}

type (
	Category   string
	Categories []Category
)

func (e *Events) ToCalendar() ([]byte, error) {
	c := ical.Calendar{}

	for _, event := range *e {
		c.Events = append(c.Events, event.Event)
	}

	return ical.Marshal(c)
}
