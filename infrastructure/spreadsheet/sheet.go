package spreadsheet

import (
	"log"

	"google.golang.org/api/sheets/v4"
)

func (spr *Sheets) GetSheets() []*sheets.Sheet {
	sheets, err := spr.srv.Spreadsheets.Get(spr.Id).Do()
	if err != nil {
		log.Fatalf("Cannot get spreadsheet: %v", err)
	}

	return sheets.Sheets
}
