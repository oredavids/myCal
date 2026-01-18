package calendar

import (
	"strings"
	"time"

	"google.golang.org/api/calendar/v3"
)

// Event wraps a calendar event with additional computed fields
type Event struct {
	*calendar.Event
	StartTime  time.Time
	IsAllDay   bool
	MeetingURL string
}

// FetchTodayEvents retrieves events for the rest of today
func FetchTodayEvents(srv *calendar.Service) ([]*Event, error) {
	now := time.Now()
	timeToStart := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location()).Format(time.RFC3339)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location()).Format(time.RFC3339)

	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(timeToStart).TimeMax(endOfDay).OrderBy("startTime").Do()
	if err != nil {
		return nil, err
	}

	return wrapEvents(events.Items), nil
}

// FetchUpcomingEvents retrieves the next N events
func FetchUpcomingEvents(srv *calendar.Service, maxResults int64, excludeToday bool) ([]*Event, error) {
	now := time.Now()

	var timeToStart string
	if excludeToday {
		timeToStart = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location()).Format(time.RFC3339)
	} else {
		timeToStart = now.Format(time.RFC3339)
	}

	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(timeToStart).MaxResults(maxResults).OrderBy("startTime").Do()
	if err != nil {
		return nil, err
	}

	return wrapEvents(events.Items), nil
}

// FetchNextEvent retrieves the next upcoming event (timed, not all-day)
func FetchNextEvent(srv *calendar.Service) (*Event, error) {
	now := time.Now()
	timeMin := now.Format(time.RFC3339)

	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(timeMin).MaxResults(5).OrderBy("startTime").Do()
	if err != nil {
		return nil, err
	}

	// Find first timed event (not all-day)
	for _, e := range events.Items {
		if e.Start.DateTime != "" {
			wrapped := wrapEvent(e)
			if wrapped.StartTime.After(now) {
				return wrapped, nil
			}
		}
	}

	return nil, nil
}

// wrapEvents converts calendar events to our Event type
func wrapEvents(items []*calendar.Event) []*Event {
	events := make([]*Event, 0, len(items))
	for _, item := range items {
		events = append(events, wrapEvent(item))
	}
	return events
}

// wrapEvent converts a single calendar event to our Event type
func wrapEvent(e *calendar.Event) *Event {
	event := &Event{Event: e}

	// Parse start time
	if e.Start.DateTime != "" {
		event.StartTime, _ = time.Parse(time.RFC3339, e.Start.DateTime)
		event.IsAllDay = false
	} else {
		event.StartTime, _ = time.Parse("2006-01-02", e.Start.Date)
		event.IsAllDay = true
	}

	// Find meeting URL
	if e.ConferenceData != nil && len(e.ConferenceData.EntryPoints) > 0 {
		event.MeetingURL = e.ConferenceData.EntryPoints[0].Uri
	} else if e.HangoutLink != "" {
		event.MeetingURL = e.HangoutLink
	} else if e.Location != "" && strings.HasPrefix(e.Location, "http") {
		event.MeetingURL = e.Location
	}

	return event
}

// TimeUntilStart returns the duration until the event starts
func (e *Event) TimeUntilStart() time.Duration {
	return e.StartTime.Sub(time.Now())
}

// GetDemoEvents returns mock events for demo/screenshot purposes
func GetDemoEvents() ([]*Event, []*Event, *Event) {
	now := time.Now()

	// Create demo upcoming events
	upcoming := []*Event{
		{
			Event: &calendar.Event{
				Summary:  "Team Standup",
				HtmlLink: "https://calendar.google.com/event/123",
			},
			StartTime:  time.Date(now.Year(), now.Month(), now.Day()+1, 10, 30, 0, 0, now.Location()),
			IsAllDay:   false,
			MeetingURL: "https://meet.google.com/abc-defg-hij",
		},
		{
			Event: &calendar.Event{
				Summary:  "Focus Time",
				HtmlLink: "https://calendar.google.com/event/456",
			},
			StartTime: time.Date(now.Year(), now.Month(), now.Day()+2, 0, 0, 0, 0, now.Location()),
			IsAllDay:  true,
		},
		{
			Event: &calendar.Event{
				Summary:  "All Hands",
				HtmlLink: "https://calendar.google.com/event/789",
			},
			StartTime:  time.Date(now.Year(), now.Month(), now.Day()+2, 10, 0, 0, 0, now.Location()),
			IsAllDay:   false,
			MeetingURL: "https://meet.google.com/xyz-uvwx-yz",
		},
		{
			Event: &calendar.Event{
				Summary:  "1:1 Meeting",
				HtmlLink: "https://calendar.google.com/event/101",
			},
			StartTime:  time.Date(now.Year(), now.Month(), now.Day()+3, 14, 30, 0, 0, now.Location()),
			IsAllDay:   false,
			MeetingURL: "https://meet.google.com/one-two-three",
		},
	}

	// No events for today in demo
	today := []*Event{}

	// Next event is the first upcoming
	next := upcoming[0]

	return today, upcoming, next
}
