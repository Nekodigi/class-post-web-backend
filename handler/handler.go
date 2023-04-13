package handler

import (
	"fmt"
	"net/http"

	"github.com/Nekodigi/class-post-web-backend/config"
	ctrCalendar "github.com/Nekodigi/class-post-web-backend/handler/calendar"
	"github.com/Nekodigi/class-post-web-backend/infrastructure"
	infraCalendar "github.com/Nekodigi/class-post-web-backend/infrastructure/calendar"
	infraFirestore "github.com/Nekodigi/class-post-web-backend/infrastructure/firestore"
	"github.com/Nekodigi/class-post-web-backend/infrastructure/spreadsheet"
	infraSheets "github.com/Nekodigi/class-post-web-backend/infrastructure/spreadsheet"

	"github.com/gin-gonic/gin"
)

var (
	client      *http.Client
	sheetsSrv   *infraSheets.Sheets
	calendarSrv *infraCalendar.Calendar
	fs          *infraFirestore.Firestore
)

func init() {
	conf := config.Load()

	client := infrastructure.GetClient(conf)

	spreadsheetId := "1BeddTk5XSxR6Kg1iPUVnp_yB_QdBLxnfmwnf3XSH8do"
	//TODO
	//calendarId := "1nm497i6vlgj0fc0o16hms2vts@group.calendar.google.com"

	sheetsSrv = spreadsheet.NewService(client, spreadsheetId)
	calendarSrv = infraCalendar.NewService(client)
	fs = infraFirestore.NewFirestore(conf)
}

func Firestore() {
	//cal, _ := fs.GetCalendarBySummary("掃除当番")
	//fs.SetCalendarId(cal.Id)
	events, _ := fs.GetEventsByDate("2023-04-13")
	if events != nil {
		for _, event := range events {

			fmt.Printf("Events:%v\n", event.EventType)
		}
	}
	//fs.AddEvent()
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func Router(e *gin.Engine) {
	e.Use(CORSMiddleware())
	e.GET("/ping", func(ctx *gin.Context) { ctx.String(http.StatusOK, "pong") })
	(&ctrCalendar.Calendar{SheetsSrv: sheetsSrv, CalendarSrv: calendarSrv, Fs: fs}).Handle(e.Group("/calendar"))
}
