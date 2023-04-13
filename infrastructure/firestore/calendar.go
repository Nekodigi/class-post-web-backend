package firestore

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/calendar/v3"
)

func (fs *Firestore) GetCalendars() ([]*calendar.Calendar, error) {
	ctx := context.Background()

	docs, err := fs.client.Collection("class_posts").Doc("calendars").Collection("list").Documents(ctx).GetAll()
	if err != nil {
		log.Fatalf("Failed getting calendars: %v", err)
		return nil, err
	}
	var calendars []*calendar.Calendar
	for _, doc := range docs {

		calendar := &calendar.Calendar{}
		doc.DataTo(calendar)
		calendars = append(calendars, calendar)
	}
	return calendars, nil
}

// Error without cal id
func (fs *Firestore) GetCalendar(id string) (*calendar.Calendar, error) {
	ctx := context.Background()

	doc, err := fs.client.Collection("class_posts").Doc("calendars").Collection("list").Doc(id).Get(ctx)
	if err != nil {
		log.Fatalf("Failed getting calendar: %v", err)
		return nil, err
	}
	calendar := &calendar.Calendar{}
	doc.DataTo(calendar)
	return calendar, nil
}

// EVENT TYPE is updated to store data
func (fs *Firestore) GetEventsByDate(date string) ([]*calendar.Event, error) {
	ctx := context.Background()

	calendarId := fs.cal.Id
	fmt.Printf("calId:%s\n", calendarId)
	docs, err := fs.client.Collection("class_posts").Doc("calendars").Collection("list").Documents(ctx).GetAll()
	if err != nil {
		log.Fatalf("Failed getting calendars: %v\n", err)
		return nil, err
	}

	var events []*calendar.Event
	for _, doc := range docs {
		cal := &calendar.Calendar{}
		doc.DataTo(cal)
		fs.SetCalendarId(cal.Id)
		events_, _ := fs.GetEventsByDateInCalendar(date)
		if events_ != nil {
			events = append(events, events_...)
		}
	}
	//FAIL if cal id is invalid!
	if calendarId != "" {
		fs.SetCalendarId(calendarId)
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("No match events.")
	}
	return events, nil
}

func (fs *Firestore) GetCalendarBySummary(summary string) (*calendar.Calendar, error) {
	ctx := context.Background()

	docs, err := fs.client.Collection("class_posts").Doc("calendars").Collection("list").Where("Summary", "==", summary).Documents(ctx).GetAll()
	if err != nil {
		log.Fatalf("Failed getting calendar: %v", err)
		return nil, err
	}
	if len(docs) != 1 {
		return nil, fmt.Errorf("No matching calendar:\n")
	}
	cal := &calendar.Calendar{}

	docs[0].DataTo(cal)
	return cal, nil
}

func (fs *Firestore) AddCalendar(cal *calendar.Calendar) {
	ctx := context.Background()

	fs.client.Collection("class_posts").Doc("calendars").Collection("list").Doc(cal.Id).Set(ctx, cal)
}

func (fs *Firestore) DeleteCalendar(id string) {
	ctx := context.Background()

	fs.SetCalendarId(id)
	fs.client.Collection("class_posts").Doc("calendars").Collection("list").Doc(id).Delete(ctx)
	fs.DeleteEvents()
}
