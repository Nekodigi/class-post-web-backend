package calendar

import (
	"log"

	"google.golang.org/api/calendar/v3"
)

func (cld *Calendar) CreateEventAllday(title string, date string, desc string) *calendar.Event {

	event := &calendar.Event{
		Summary:     title,
		Description: desc,
		Start: &calendar.EventDateTime{
			Date:     date,
			TimeZone: "Asia/Tokyo",
		},
		End: &calendar.EventDateTime{
			Date:     date,
			TimeZone: "Asia/Tokyo",
		},
	}
	// events := &calendar.Events{
	// 	Summary: title,
	// 	Items:   []*calendar.Event{event},
	// }

	return event
}

func (cld *Calendar) DeleteEvent(id string) {
	err := cld.srv.Events.Delete(cld.Id, id).Do()
	if err != nil {
		log.Fatalf("Unable to delete event. %v\n", err)
	}
}

func (cld *Calendar) AddEvent(event *calendar.Event) *calendar.Event {
	event, err := cld.srv.Events.Insert(cld.Id, event).Do()
	if err != nil {
		log.Fatalf("Unable to create event. %v\n", err)
	}
	return event
}
