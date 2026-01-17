package tui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"oredavids.com/myCal/internal/calendar"
)

// RenderHeader returns the styled header with date and greeting
func RenderHeader() string {
	now := time.Now()
	dateStr := now.Format("Monday, January 2, 2006")
	greeting := fmt.Sprintf("%s, %s!", getGreeting(), getUserName())

	header := lipgloss.JoinVertical(
		lipgloss.Left,
		DateStyle.Render("  "+dateStr),
		GreetingStyle.Render("  "+greeting),
	)

	headerBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor).
		Padding(0, 1).
		Render(header)

	return headerBox
}

// RenderSectionTitle returns a styled section title
func RenderSectionTitle(title string, icon string) string {
	return SectionTitleStyle.Render(icon + " " + title)
}

// RenderEventList renders a list of events in a styled box
func RenderEventList(events []*calendar.Event, isToday bool, selectedIndex int) string {
	if len(events) == 0 {
		return NoEventsStyle.Render("No events")
	}

	var eventRows []string

	for i, event := range events {
		eventRow := RenderEvent(event, isToday, i == selectedIndex)
		eventRows = append(eventRows, eventRow)

		// Add divider between events (but not after last one)
		if i < len(events)-1 {
			divider := DividerStyle.Render("─────────────────────────────────")
			eventRows = append(eventRows, divider)
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left, eventRows...)
	return EventBoxStyle.Render(content)
}

// RenderEvent renders a single event with styling
func RenderEvent(event *calendar.Event, isToday bool, selected bool) string {
	var rows []string

	// Title
	title := event.Summary
	if selected {
		title = SelectedStyle.Render("▸ " + title)
	} else {
		title = EventTitleStyle.Render(title)
	}
	rows = append(rows, title)

	// Time
	var timeStr string
	if event.IsAllDay {
		timeStr = EventAllDayStyle.Render(event.StartTime.Format("Monday") + " · All day")
	} else {
		var formattedTime string
		if isToday {
			formattedTime = event.StartTime.Format("3:04 PM")
		} else {
			formattedTime = event.StartTime.Format("Mon · 3:04 PM")
		}
		timeStr = EventTimeStyle.Render(formattedTime)
	}

	if HyperlinkSupport {
		// Terminals with hyperlink support: compact clickable links
		var links []string
		if event.MeetingURL != "" {
			links = append(links, RenderLink("[Join]", event.MeetingURL, "green"))
		}
		if event.HtmlLink != "" {
			links = append(links, RenderLink("[Cal]", event.HtmlLink, "gray"))
		}
		infoRow := timeStr
		if len(links) > 0 {
			infoRow += "  " + strings.Join(links, " ")
		}
		rows = append(rows, infoRow)
	} else {
		// Terminals without hyperlink support: show URL on separate line
		rows = append(rows, timeStr)
		if event.MeetingURL != "" {
			rows = append(rows, RenderFallbackURL(event.MeetingURL))
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

// RenderCountdown renders the next meeting countdown
func RenderCountdown(event *calendar.Event) string {
	if event == nil {
		return ""
	}

	duration := event.TimeUntilStart()
	if duration < 0 {
		return ""
	}

	return fmt.Sprintf("%s %s %s",
		LabelStyle.Render("Next:"),
		EventTitleStyle.Render(event.Summary),
		CountdownStyle.Render(FormatDuration(duration)),
	)
}

// FormatDuration formats a duration into a human-readable string
func FormatDuration(d time.Duration) string {
	if d < 0 {
		return "now"
	}

	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		if days == 1 {
			return "in 1 day"
		}
		return fmt.Sprintf("in %d days", days)
	}
	if hours > 0 {
		if hours == 1 {
			if minutes > 0 {
				return fmt.Sprintf("in 1 hr %d min", minutes)
			}
			return "in 1 hour"
		}
		return fmt.Sprintf("in %d hours", hours)
	}
	if minutes > 0 {
		if minutes == 1 {
			return "in 1 minute"
		}
		return fmt.Sprintf("in %d minutes", minutes)
	}
	return "starting now"
}

// RenderHelp renders the help text
func RenderHelp() string {
	return HelpStyle.Render("↑/↓ navigate • enter join • r refresh • q quit")
}

func getGreeting() string {
	hour := time.Now().Hour()
	switch {
	case hour < 12:
		return "Good morning"
	case hour < 17:
		return "Good afternoon"
	default:
		return "Good evening"
	}
}

func getUserName() string {
	if name := os.Getenv("USER"); name != "" {
		return strings.ToUpper(string(name[0])) + name[1:]
	}
	return "there"
}
