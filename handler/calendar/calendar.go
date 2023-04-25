package calendar

import (
	"fmt"
	"net/http"
	"time"

	infraCalendar "github.com/Nekodigi/class-post-web-backend/infrastructure/calendar"
	infraFirestore "github.com/Nekodigi/class-post-web-backend/infrastructure/firestore"
	infraSheets "github.com/Nekodigi/class-post-web-backend/infrastructure/spreadsheet"
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
	"google.golang.org/api/calendar/v3"
)

var (
	c *Calendar
)

type (
	Calendar struct {
		SheetsSrv   *infraSheets.Sheets
		CalendarSrv *infraCalendar.Calendar
		Fs          *infraFirestore.Firestore
	}
)

func (c *Calendar) Handle(rg *gin.RouterGroup) {
	rg.GET("/update", func(gc *gin.Context) { //urlupdate for client
		fmt.Println("sheet update:")
		c.CalendarUpdate()
		gc.String(http.StatusOK, "ok")
	})
	rg.GET("/calendars", func(gc *gin.Context) {
		cals, err := c.Fs.GetCalendars()
		if err != nil {
			var empty []*calendar.Calendar
			gc.JSON(http.StatusOK, empty)
		} else {
			gc.JSON(http.StatusOK, cals)
		}
	})
	rg.GET("/day", func(gc *gin.Context) {
		date := gc.Query("date")
		if date == "" {
			date = time.Now().Format("2006-01-02")
		}
		fmt.Println("Query:" + date)
		events, err := c.Fs.GetEventsByDate(date)
		if err != nil {
			gc.JSON(http.StatusOK, make([]int, 0))
		} else {
			gc.JSON(http.StatusOK, events)
		}
	})
}

func (c *Calendar) CalendarUpdate() {
	// Prints the names and majors of students in a sample spreadsheet:
	// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit

	sheets := c.SheetsSrv.GetSheets()
	//TODO find unused data in fs. Exist in fs but not in spreadsheet.

	//* DELETE ALL
	// for _, event := range events {
	// 	fs.DeleteEvent(event.Id)
	// 	calendarSrv.DeleteEvent(event.Id)
	// }
	//TODO Delete calendar it self!
	//* List up unused calendar
	//* Delete calendar, contained events

	outdatedCalendars, err := c.Fs.GetCalendars()
	if err != nil {
		fmt.Errorf("Failed to get calendars")
	}

	for _, sheet := range sheets {
		title := sheet.Properties.Title
		fmt.Println(title)

		//* Init calendar
		cal, _ := c.Fs.GetCalendarBySummary(title)
		if cal != nil {
			c.CalendarSrv.SetCalendarId(cal.Id)
			outdatedCalendars = funk.Filter(outdatedCalendars, func(calendar *calendar.Calendar) bool {
				return cal.Id != calendar.Id
			}).([]*calendar.Calendar)
		} else {
			cal := c.CalendarSrv.NewCalendar(title)
			c.Fs.AddCalendar(cal)
		}
		c.Fs.SetCalendarId(c.CalendarSrv.Id)

		//* Add events
		outdatedEvents := c.Fs.GetEvents()

		resp := c.SheetsSrv.GetValues(title + "!A2:D")

		fmt.Println("title, date, assign")
		for _, row := range resp.Values {
			//* required!
			if len(row) < 2 || row[0] == nil || row[1] == nil {
				continue
			}

			// Print columns A and E, which correspond to indices 0 and 4.
			date, _ := time.Parse("1/2/2006", fmt.Sprintf("%s", row[1]))

			desc := ""
			if len(row) >= 4 {
				desc = fmt.Sprintf("担当:%s", row[3])
			}
			event := c.CalendarSrv.CreateEventAllday(fmt.Sprintf("%s", row[0]), date.Format("2006-01-02"), desc)
			oldEvent, _ := c.Fs.GetEventByAttr(event)
			//remove event from fs list as in use.
			if oldEvent != nil {
				outdatedEvents = funk.Filter(outdatedEvents, func(event *calendar.Event) bool {
					return oldEvent.Id != event.Id
				}).([]*calendar.Event)
			}

			if oldEvent != nil {
				// fs.DeleteEvent(oldEvent.Id)
				// calendarSrv.DeleteEvent(oldEvent.Id)
				//fmt.Println("omit create")
			} else {
				event = c.CalendarSrv.AddEvent(event)
				c.Fs.AddEvent(event)
			}
			//fmt.Printf("%s, %s, %s\n", row[0], date.String(), row[3])
		}

		//* Remove unused event
		for _, event := range outdatedEvents {
			c.CalendarSrv.DeleteEvent(event.Id)
			c.Fs.DeleteEvent(event.Id)
			fmt.Println("delete outdated event:", event.Id)
		}
	}
	//* Remove unused calendar
	for _, cal := range outdatedCalendars {
		//TODO implement delete calendar
		c.CalendarSrv.DeleteCalendar(cal.Id)
		c.Fs.DeleteCalendar(cal.Id)
		fmt.Println("delete outdated calendar:", cal.Id)
	}
}
