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

var calwCmd = &cobra.Command{
	Use:   "calw",
	Short: "get a week calendar",
	Run: func(cmd *cobra.Command, args []string) {
		now := time.Now()
		weekday := int(now.Weekday())
		sunday := now.AddDate(0, 0, -weekday)
		var weekHeaderSlc [7]string

		for i := 0; i < 7; i++ {
			day := sunday.AddDate(0, 0, i)
			weekHeaderSlc[i] = day.Format("Mon (02)")
		}
		sun := time.Date(sunday.Year(), sunday.Month(), sunday.Day(), 0, 0, 0, 0, now.Location())

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(weekHeaderSlc[:])

		srv := api.GetCalendar()

		loading := make(chan bool)
		go utils.ShowLoading(loading)
		events, err := srv.Events.List("primary").ShowDeleted(false).SingleEvents(true).TimeMin(sun.Format(time.RFC3339)).TimeMax(timeutils.EndOfDay(sunday.AddDate(0, 0, 6)).Format(time.RFC3339)).OrderBy("startTime").Do()
		loading <- true

		if err != nil {
			log.Fatalf("unable to retreive calendar event: %v", err)
		}

		var weekSlc [7][]string
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				allDay, err := time.Parse(timeutils.AllDayDefaultLayout, item.Start.Date)
				if err != nil {
					log.Fatalf("unable to parse date: %v\n", err)
				}
				date = allDay.Format(time.RFC3339)
			}

			day, err := time.Parse(time.RFC3339, date)
			if err != nil {
				log.Fatalf("eror: %v", err)
			}
			dayInt := int(day.Weekday())
			data := fmt.Sprintf("%v (%v)", item.Summary, day.Format("15:04"))
			weekSlc[dayInt] = append(weekSlc[dayInt], data)
		}

		maxLen := 0

		for _, row := range weekSlc {
			if len(row) > maxLen {
				maxLen = len(row)
			}
		}

		transposed := make([][7]string, maxLen)

		for i := range weekSlc {
			for j := range weekSlc[i] {
				transposed[j][i] = weekSlc[i][j]
			}
		}

		for _, v := range transposed {
			table.Append(v[:])
		}

		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(calwCmd)
}
