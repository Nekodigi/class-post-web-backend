package firestore

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/Nekodigi/class-post-web-backend/config"
	"google.golang.org/api/calendar/v3"
)

var (
	fs *Firestore
)

type (
	Firestore struct {
		app    *firebase.App
		client *firestore.Client
		cal    *calendar.Calendar
	}
)

func NewFirestore(conf *config.Config) *Firestore {
	ctx := context.Background()
	//sa := option.WithCredentialsFile("credentials/serviceAccount.json")
	config := &firebase.Config{ProjectID: conf.ProjectId}

	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	fs = &Firestore{
		app:    app,
		client: client,
		cal:    &calendar.Calendar{Id: ""},
	}

	return fs
}

func (fs *Firestore) SetCalendarId(id string) {
	cal, err := fs.GetCalendar(id)
	if err != nil {
		fmt.Errorf("Failed to set calendar")
	}
	fs.cal = cal
}
