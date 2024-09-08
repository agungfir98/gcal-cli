/*
Copyright © 2024 Agung Firmansyah agungfir98@gmail.com
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/agungfir98/gcal-cli/api"
	timeutils "github.com/agungfir98/gcal-cli/utils/time_utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"google.golang.org/api/calendar/v3"
)

var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "shows event in your google calendar",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		srv := api.GetCalendar()
		calendar := NewCalendar()
		calendar.GetEvents(srv)
		calendar.Filter()
		calendar.Render()
	},
}

var (
	startDate     string
	endDate       string
	limit         int64
	dateFormat    string
	excludeAllDay bool
)

type Show string

const (
	Table Show = "table"
	List  Show = "list"
)

func (s *Show) String() string {
	return string(*s)
}

func (s *Show) Type() string {
	return "Show"
}

func (s *Show) Set(v string) error {
	switch v {
	case "table", "list":
		*s = Show(v)
		return nil
	default:
		return errors.New(`must be one of "table" or "list" `)
	}
}

var show Show = Table

func init() {
	eventsCmd.Flags().StringVar(&startDate, "start-date", "", "start date to query")
	eventsCmd.Flags().StringVar(&endDate, "end-date", "", "end date to query")
	eventsCmd.Flags().StringVar(&dateFormat, "date-format", "%d %m %Y %H:%M", "date format for input and output date")
	eventsCmd.Flags().Int64VarP(&limit, "limit", "L", 10, "limit event result fetched")
	eventsCmd.Flags().Var(&show, "show", "show enum. allowed('table', 'list')")
	eventsCmd.Flags().BoolVar(&excludeAllDay, "exclude-all-day", false, "exclude all day event")

	rootCmd.AddCommand(eventsCmd)
}

type EventList struct {
	Summary   string
	Start     string
	End       string
	Attendees []string
	Link      string
}

type Flags struct {
	limit         int64
	dateFormat    string
	startDate     string
	excludeAllDay bool
	endDate       string
	show          Show
}

func NewFlags() *Flags {
	f := &Flags{
		limit:         limit,
		dateFormat:    dateFormat,
		startDate:     startDate,
		endDate:       endDate,
		show:          show,
		excludeAllDay: excludeAllDay,
	}
	return f
}

type Calendar struct {
	flags     *Flags
	events    *calendar.Events
	eventlist []EventList
}

func NewCalendar() *Calendar {
	c := &Calendar{}
	c.flags = NewFlags()
	return c
}

func (c *Calendar) GetEvents(service *calendar.Service) {

	if c.flags.startDate == "" {
		c.flags.startDate = time.Now().Format(timeutils.ConvertToGoLayout(c.flags.dateFormat))
	}
	if c.flags.endDate == "" {
		c.flags.endDate = time.Now().Add(time.Hour * 12).Format(timeutils.ConvertToGoLayout(c.flags.dateFormat))
	}

	tMin, err := timeutils.ParseWithCustomFormat(c.flags.dateFormat, c.flags.startDate)
	if err != nil {
		log.Fatalf("unable to parse start date parameter")
	}
	tMax, err := timeutils.ParseWithCustomFormat(c.flags.dateFormat, c.flags.endDate)
	if err != nil {
		log.Fatalf("unable to parse end date parameter")
	}

	if tMin.After(tMax) {
		log.Fatalf("start date can't be greater than end date")
	}

	events, err := service.Events.List("primary").MaxResults(c.flags.limit).TimeMin(tMin.Format(time.RFC3339)).TimeMax(tMax.Format(time.RFC3339)).Do()
	if err != nil {
		log.Fatalf("unable to retreive events")
	}
	c.events = events
}

func (c *Calendar) Filter() {
	var Items []*calendar.Event

	for _, item := range c.events.Items {
		if c.flags.excludeAllDay && item.Start.DateTime == "" {
			continue
		}
		Items = append(Items, item)
	}
	c.events.Items = Items
}

func (c *Calendar) Render() {

	for _, item := range c.events.Items {
		var start string
		if item.Start.DateTime == "" {
			start = item.Start.Date
		} else {
			parsed, err := time.Parse(time.RFC3339, item.Start.DateTime)
			if err != nil {
				log.Fatalf("unable to parse start date")
			}
			start = parsed.Format(timeutils.ConvertToGoLayout(dateFormat))
		}

		var end string
		if item.End.DateTime == "" {
			end = item.End.Date
		} else {
			parsed, err := time.Parse(time.RFC3339, item.End.DateTime)
			if err != nil {
				log.Fatalf("unable to parse end date 2")
			}
			end = parsed.Format(timeutils.ConvertToGoLayout(c.flags.dateFormat))
		}

		var attendees []string
		for _, attendee := range item.Attendees {
			attendees = append(attendees, attendee.Email)
		}

		c.eventlist = append(c.eventlist, EventList{
			Summary:   item.Summary,
			Start:     start,
			End:       end,
			Attendees: attendees,
			Link:      item.HangoutLink,
		})
	}

	switch c.flags.show {
	case Table:
		table := tablewriter.NewWriter(os.Stdout)
		for _, event := range c.eventlist {
			var attendees string
			for i, attendee := range event.Attendees {
				if i > 0 {
					attendees += "\n"
				}
				attendees += attendee

			}
			table.Append([]string{event.Summary, event.Start, event.End, attendees, event.Link})
		}

		table.SetHeader([]string{"summary", "start", "end", "attendees", "link"})
		table.Render()
	case List:
		// TODO: continue tomorrow 2024-09-08 implement the list render
		fmt.Println("a listttt")
		fmt.Println(show)
	default:
		log.Fatalf("enum is not correct and this message should not show. please contact the project owner")
	}
}
