/*
Copyright © 2024 Agung Firmansyah agungfir98@gmail.com
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/agungfir98/gcal-cli/api"
	"github.com/agungfir98/gcal-cli/utils"
	timeutils "github.com/agungfir98/gcal-cli/utils/time_utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "shows event in your google calendar",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		tMin, tMax := ParseDates(startDate, endDate)

		srv := api.GetCalendar()

		loading := make(chan bool)
		go utils.ShowLoading(loading)
		events, err := srv.Events.List("primary").ShowDeleted(false).SingleEvents(true).TimeMin(tMin.Format(time.RFC3339)).TimeMax(tMax.Format(time.RFC3339)).MaxResults(max).OrderBy("startTime").Do()
		loading <- true
		close(loading)

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
				dateStart, err := time.Parse(timeutils.AllDayDefaultLayout, item.Start.Date)
				if err != nil {
					fmt.Printf("unable to parse all day date: %v\n", err)
				}
				itemStart = dateStart.Format(time.RFC3339)
			}

			start, err := time.Parse(time.RFC3339, itemStart)
			if err != nil {
				log.Fatalf("unable to parse time start")
			}

			itemEnd := item.End.DateTime
			if itemEnd == "" {
				dateEnd, err := time.Parse(timeutils.AllDayDefaultLayout, item.Start.Date)
				if err != nil {
					fmt.Printf("unable to parse all day date: %v\n", err)
				}
				itemEnd = timeutils.EndOfDay(dateEnd).Format(time.RFC3339)
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
				Start:     start.Format(timeutils.DefaultLayoutWithTime),
				End:       end.Format(timeutils.DefaultLayoutWithTime),
				Attendees: attendeeInfo,
			}
			row := []string{event.Summary, event.Start, event.End, event.Attendees}
			table.Append(row)
		}

		table.Render()
	},
}

func init() {
	eventsCmd.Flags().StringVar(&startDate, "start-date", "", "start date to query\n format('31 08 2024' or '31 08 2024 15:00')")
	eventsCmd.Flags().StringVar(&endDate, "end-date", "", "end date to query\n format('31 08 2024' or '31 08 2024 15:00')")
	eventsCmd.Flags().Int64VarP(&max, "max", "m", 10, "max events to be fetched")
	rootCmd.AddCommand(eventsCmd)
}

func ParseDates(start, end string) (time.Time, time.Time) {
	t := time.Now()

	if start == "" {
		start = t.Format(timeutils.DefaultLayoutWithTime)
	}
	s, err := time.ParseInLocation(timeutils.DefaultLayout, start, t.Location())
	if err != nil {
		s, err = time.ParseInLocation(timeutils.DefaultLayoutWithTime, start, t.Location())
		if err != nil {
			log.Fatalf("unable to parse start date: %v\n", err)
		}
	}

	if end == "" {
		end = timeutils.EndOfDay(s).Format(timeutils.DefaultLayout)
	}

	e, err := time.ParseInLocation(timeutils.DefaultLayout, end, t.Location())
	if err != nil {
		log.Fatalf("unable to parse end date: %v\n", err)
	}
	e = timeutils.EndOfDay(e)

	if s.After(e) {
		log.Fatalf("end date must not be greater than start date")
	}

	return s, e
}

var max int64
var startDate string
var endDate string

type Event struct {
	// Id        string `json:"id"`
	Start     string `json:"start"`
	Summary   string `json:"summary"`
	End       string `json:"end"`
	Attendees string `json:"attendee"`
}
