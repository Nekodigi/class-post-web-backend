package firestore

import (
	"context"
	"fmt"
	"log"

	"github.com/Nekodigi/class-post-web-backend/models"
	"github.com/thoas/go-funk"
	"google.golang.org/api/calendar/v3"
)

// * List up all events doc id as delete.
// * find corresponding data in events foreach sheets value. Remove from delete hit doc id if hit. Add doc to incoming if don't.
// * current A F X / incoming G F I
// * delete A X / add G I

func (fs *Firestore) GetEventsDelta(current []*calendar.Event, incoming []*calendar.Event) *models.EventsDelta {
	delta := &models.EventsDelta{Delete: current, Add: []*calendar.Event{}}
	for _, event := range incoming {
		currentEvent, err := fs.GetEventByAttr(event)
		if err != nil {
			delta.Delete = funk.Filter(delta.Delete, func(event *calendar.Event) bool {
				return !(event.Summary == currentEvent.Summary && event.Start.Date == currentEvent.Start.Date)
			}).([]*calendar.Event)
		} else {
			delta.Add = append(delta.Add, event)
		}
	}
	return delta
}

func (fs *Firestore) GetEvents() []*calendar.Event {

	ctx := context.Background()

	docs, err := fs.client.Collection("class_posts").Doc("events").Collection(fs.cal.Id).Documents(ctx).GetAll()
	if err != nil {
		log.Fatalf("Failed getting events: %v", err)
	}
	var events []*calendar.Event
	for _, doc := range docs {

		event := &calendar.Event{}
		doc.DataTo(event)
		events = append(events, event)
	}
	return events
}

func (fs *Firestore) DeleteEvents() {

	ctx := context.Background()

	refs, err := fs.client.Collection("class_posts").Doc("events").Collection(fs.cal.Id).DocumentRefs(ctx).GetAll()
	if err != nil {
		log.Fatalf("Failed getting events: %v", err)
	}
	for _, ref := range refs {
		ref.Delete(ctx)
	}
}

// EVENT TYPE is updated to store data
func (fs *Firestore) GetEventsByDateInCalendar(date string) ([]*calendar.Event, error) {
	ctx := context.Background()

	docs, err := fs.client.Collection("class_posts").Doc("events").Collection(fs.cal.Id).Where("Start.Date", "==", date).Documents(ctx).GetAll()
	if err != nil {
		log.Fatalf("Failed getting events: %v\n", err)
		return nil, err
	}

	var events []*calendar.Event
	for _, doc := range docs {

		event := &calendar.Event{}
		doc.DataTo(event)
		event.EventType = fs.cal.Summary
		events = append(events, event)
	}
	if len(events) == 0 {
		return nil, fmt.Errorf("No match events.")
	}
	return events, nil
}

func (fs *Firestore) GetEventByAttr(event *calendar.Event) (*calendar.Event, error) {
	ctx := context.Background()

	docs, err := fs.client.Collection("class_posts").Doc("events").Collection(fs.cal.Id).Where("Summary", "==", event.Summary).Where("Start.Date", "==", event.Start.Date).Documents(ctx).GetAll()
	if err != nil {
		log.Fatalf("Failed getting events: %v\n", err)
		return nil, err
	}
	event_ := &calendar.Event{}
	if len(docs) != 1 {
		return nil, fmt.Errorf("no matching event:\n")
	}
	docs[0].DataTo(event_)
	return event_, nil
}

func (fs *Firestore) DeleteEvent(id string) {
	ctx := context.Background()

	_, err := fs.client.Collection("class_posts").Doc("events").Collection(fs.cal.Id).Doc(id).Delete(ctx)

	if err != nil {
		log.Fatalf("Failed deleting event: %v", err)
	}
}

func (fs *Firestore) AddEvent(event *calendar.Event) {
	ctx := context.Background()

	_, err := fs.client.Collection("class_posts").Doc("events").Collection(fs.cal.Id).Doc(event.Id).Set(ctx, event)

	if err != nil {
		log.Fatalf("Failed adding event: %v", err)
	}

}
