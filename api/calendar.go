package api

import (
	"context"
	"log"

	"github.com/agungfir98/gcal-cli/auth"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func GetCalendar() *calendar.Service {
	client := auth.NewClient().GetClient()

	srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("unable to retreive calendar client: %v\n", err)
	}
	return srv
}
