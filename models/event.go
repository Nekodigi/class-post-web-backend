package models

import "google.golang.org/api/calendar/v3"

type EventsDelta struct {
	Delete []*calendar.Event
	Add    []*calendar.Event
}

type Person struct {
	name string
	age  int
}
