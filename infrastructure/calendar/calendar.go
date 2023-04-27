package calendar

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

var (
	cld *Calendar
)

type (
	Calendar struct {
		srv *calendar.Service
		Id  string
	}
)

func NewService(client *http.Client) *Calendar {
	ctx := context.Background()

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	cld = &Calendar{
		srv: srv,
		Id:  "",
	}

	// if calendarId == "" {
	// 	cld.NewCalendar()
	// }

	return cld
}

func (cld *Calendar) SetCalendarId(calendarId string) {
	cld.Id = calendarId
}

func (cld *Calendar) NewCalendar(summary string) *calendar.Calendar {
	calendar_ := &calendar.Calendar{
		Summary:  summary,
		TimeZone: "Asia/Tokyo",
	}

	newCalendar, _ := cld.srv.Calendars.Insert(calendar_).Do()
	scope := &calendar.AclRuleScope{Type: "default", Value: ""}
	cld.srv.Acl.Insert(newCalendar.Id, &calendar.AclRule{Scope: scope, Role: "reader"}).Do()

	fmt.Printf("New calendar created: %s\n", newCalendar.Summary)
	cld.Id = newCalendar.Id
	return newCalendar
}

func (cld *Calendar) DeleteCalendar(id string) {
	cld.srv.Calendars.Delete(id).Do()
	fmt.Printf("Calendar deleted: %s\n", id)
}
