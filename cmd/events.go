/*
Copyright © 2024 Agung Firmansyah agungfir98@gmail.com
*/
package cmd

import (
	"fmt"
	"gcal-cli/api"
	"gcal-cli/utils"
	timeutils "gcal-cli/utils/time_utils"
	"log"
	"time"

	"github.com/spf13/cobra"
)

var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "shows event in your google calendar",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		//TODO: make a table view of this.

		t := time.Now()
		tMin, err := time.ParseInLocation(timeutils.DefaultLayout, timeMin, t.Location())
		if err != nil {
			log.Fatalf("min-time format is not valid: %v", err)
		}
		tMax, err := time.ParseInLocation(timeutils.DefaultLayout, timeMax, t.Location())
		if err != nil {
			log.Fatalf("max-time format is not valid: %v", err)
		}

		srv := api.GetCalendar()

		loading := make(chan bool)
		go utils.ShowLoading(loading)
		events, err := srv.Events.List("primary").ShowDeleted(false).SingleEvents(true).TimeMin(tMin.Local().Format(time.RFC3339)).TimeMax(tMax.Format(time.RFC3339)).MaxResults(max).OrderBy("startTime").Do()
		loading <- true

		if err != nil {
			log.Fatalf("unable to retreive next ten events: %v\n", err)
		}
		fmt.Println("upcoming events")
		if len(events.Items) == 0 {
			fmt.Println("no upcomeing event")
			return
		}

		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}

			fmt.Printf("%s (%v)\n", item.Summary, date)
		}

	},
}

func init() {
	eventsCmd.Flags().BoolVar(&showAttendee, "show-attendee", false, "show attendee")
	eventsCmd.Flags().StringVar(&timeMin, "min-time", time.Now().Format(timeutils.DefaultLayout), "min-time of the events to be fetched\n")
	eventsCmd.Flags().StringVar(&timeMax, "max-time", timeutils.EndOfDay(time.Now()).Format(timeutils.DefaultLayout), "max-time of the events to be fetched\n")
	eventsCmd.Flags().Int64VarP(&max, "max", "m", 10, "max events to be fetched")
	rootCmd.AddCommand(eventsCmd)
}

var showAttendee bool
var max int64
var timeMin string
var timeMax string
