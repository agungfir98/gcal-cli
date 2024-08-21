package api

import (
	"context"
	"gcal-cli/auth"
	"gcal-cli/utils"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func GetCalendar() *calendar.Service {
	path := utils.GetCredentialsFile()
	b, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("unable to read credentials file: %v\n", err)
	}

	cfg, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("unable to parse client secret file to config: %v\n", err)
	}

	client := auth.GetClient(cfg)

	srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("unable to retreive calendar client: %v\n", err)
	}
	return srv
}
