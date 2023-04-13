package spreadsheet

import (
	"context"
	"log"
	"net/http"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var (
	spr *Sheets
)

type (
	Sheets struct {
		srv *sheets.Service
		Id  string
	}
)

func NewService(client *http.Client, spreadsheetId string) *Sheets {
	ctx := context.Background()

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	spr = &Sheets{
		srv: srv,
		Id:  spreadsheetId,
	}

	return spr
}
