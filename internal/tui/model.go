package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/browser"
	gcal "google.golang.org/api/calendar/v3"

	"oredavids.com/myCal/internal/calendar"
)

// Model is the bubbletea model for the TUI
type Model struct {
	calendarService *gcal.Service
	todayEvents     []*calendar.Event
	upcomingEvents  []*calendar.Event
	nextEvent       *calendar.Event
	selectedIndex   int
	allEvents       []*calendar.Event // combined list for selection
	status          string
	lastRefresh     time.Time
	err             error
}

// tickMsg is sent every second to update the countdown
type tickMsg time.Time

// refreshMsg is sent when data needs to be refreshed
type refreshMsg struct{}

// NewModel creates a new TUI model
func NewModel(srv *gcal.Service) Model {
	return Model{
		calendarService: srv,
		selectedIndex:   0,
		lastRefresh:     time.Now(),
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.fetchEvents(),
		tickEvery(),
	)
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit

		case "up", "k":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}

		case "down", "j":
			if m.selectedIndex < len(m.allEvents)-1 {
				m.selectedIndex++
			}

		case "enter":
			if m.selectedIndex < len(m.allEvents) {
				event := m.allEvents[m.selectedIndex]
				if event.MeetingURL != "" {
					browser.OpenURL(event.MeetingURL)
					m.status = fmt.Sprintf("Opening %s...", event.Summary)
				} else {
					m.status = "No meeting link for this event"
				}
			}

		case "r":
			m.status = "Refreshing..."
			return m, m.fetchEvents()
		}

	case tickMsg:
		// Check if we should auto-refresh (every 5 minutes)
		if time.Since(m.lastRefresh) > 5*time.Minute {
			return m, tea.Batch(tickEvery(), m.fetchEvents())
		}
		return m, tickEvery()

	case eventsMsg:
		m.todayEvents = msg.today
		m.upcomingEvents = msg.upcoming
		m.nextEvent = msg.next
		m.allEvents = append(m.todayEvents, m.upcomingEvents...)
		m.lastRefresh = time.Now()
		m.status = ""
		if m.selectedIndex >= len(m.allEvents) && len(m.allEvents) > 0 {
			m.selectedIndex = len(m.allEvents) - 1
		}

	case errMsg:
		m.err = msg.err
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	var b strings.Builder

	// Header
	b.WriteString(renderHeader(getUserName()))
	b.WriteString("\n")

	// Next meeting countdown
	if m.nextEvent != nil {
		countdown := RenderCountdown(m.nextEvent)
		if countdown != "" {
			b.WriteString(countdown)
			b.WriteString("\n")
		}
	}

	// Status message
	if m.status != "" {
		b.WriteString(StatusStyle.Render(m.status))
		b.WriteString("\n")
	}

	// Today's events
	b.WriteString("\n")
	b.WriteString(RenderSectionTitle("Today", "🗓"))
	b.WriteString("\n")
	if len(m.todayEvents) == 0 {
		b.WriteString(NoEventsStyle.Render("No events remaining today"))
	} else {
		b.WriteString(RenderEventList(m.todayEvents, true, m.selectedIndex))
	}
	b.WriteString("\n")

	// Upcoming events (only show if today has < 3 events)
	if len(m.todayEvents) < 3 && len(m.upcomingEvents) > 0 {
		b.WriteString("\n")
		b.WriteString(RenderSectionTitle("Upcoming", "🗓"))
		b.WriteString("\n")
		// Adjust selected index for upcoming section
		upcomingSelectedIndex := m.selectedIndex - len(m.todayEvents)
		b.WriteString(RenderEventList(m.upcomingEvents, false, upcomingSelectedIndex))
		b.WriteString("\n")
	}

	// Help
	b.WriteString(RenderHelp())
	b.WriteString("\n")

	// Error display
	if m.err != nil {
		errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444"))
		b.WriteString(errStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		b.WriteString("\n")
	}

	return b.String()
}

// eventsMsg carries fetched events
type eventsMsg struct {
	today    []*calendar.Event
	upcoming []*calendar.Event
	next     *calendar.Event
}

// errMsg carries an error
type errMsg struct {
	err error
}

// fetchEvents returns a command that fetches calendar events
func (m Model) fetchEvents() tea.Cmd {
	return func() tea.Msg {
		today, err := calendar.FetchTodayEvents(m.calendarService)
		if err != nil {
			return errMsg{err}
		}

		var upcoming []*calendar.Event
		if len(today) < 3 {
			upcoming, err = calendar.FetchUpcomingEvents(m.calendarService, 5, true)
			if err != nil {
				return errMsg{err}
			}
		}

		next, _ := calendar.FetchNextEvent(m.calendarService)

		return eventsMsg{
			today:    today,
			upcoming: upcoming,
			next:     next,
		}
	}
}

// tickEvery returns a command that sends a tick every second
func tickEvery() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Run starts the TUI
func Run(srv *gcal.Service) error {
	p := tea.NewProgram(NewModel(srv), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
