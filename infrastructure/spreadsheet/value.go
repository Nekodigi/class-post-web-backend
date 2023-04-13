package spreadsheet

import (
	"fmt"
	"log"

	"google.golang.org/api/sheets/v4"
)

func (spr *Sheets) GetValues(readRange string) *sheets.ValueRange {
	resp, err := spr.srv.Spreadsheets.Values.Get(spr.Id, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	}
	return resp
}
