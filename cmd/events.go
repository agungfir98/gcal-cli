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
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
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

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"summary", "start", "end", "attendee"})

		for _, item := range events.Items {
			itemStart := item.Start.DateTime
			if itemStart == "" {
				itemStart = item.Start.Date
			}

			start, err := time.Parse(time.RFC3339, itemStart)
			if err != nil {
				log.Fatalf("unable to parse time start")
			}

			itemEnd := item.End.DateTime
			if itemEnd == "" {
				itemEnd = item.Start.Date
			}

			end, err := time.Parse(time.RFC3339, itemEnd)
			if err != nil {
				log.Fatalf("unable to parse time end")
			}

			var attendeeInfo string

			for i, attendee := range item.Attendees {
				if i > 0 {
					attendeeInfo += "\n"
				}
				if attendee.DisplayName != "" {
					attendeeInfo += attendee.DisplayName + "(" + attendee.Email + ")"
				}
				attendeeInfo += attendee.Email
			}

			event := Event{
				// Id:        item.Id,
				Summary:   item.Summary,
				Start:     start.Format(timeutils.DefaultLayout),
				End:       end.Format(timeutils.DefaultLayout),
				Attendees: attendeeInfo,
			}
			row := []string{event.Summary, event.Start, event.End, event.Attendees}
			table.Append(row)
		}

		table.Render()
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

type Event struct {
	// Id        string `json:"id"`
	Start     string `json:"start"`
	Summary   string `json:"summary"`
	End       string `json:"end"`
	Attendees string `json:"attendee"`
}
