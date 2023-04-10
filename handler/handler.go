package handler

import (
	"context"
	"log"
	

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"

)

func init() {
	ctx := context.Background()

	client :=
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}
}
