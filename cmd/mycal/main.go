package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	gcal "google.golang.org/api/calendar/v3"

	"oredavids.com/myCal/internal/auth"
	"oredavids.com/myCal/internal/calendar"
	"oredavids.com/myCal/internal/config"
	"oredavids.com/myCal/internal/tui"
)

func main() {
	// Parse flags
	watchMode := flag.Bool("watch", false, "Run in interactive watch mode")
	flag.BoolVar(watchMode, "w", false, "Run in interactive watch mode (shorthand)")
	flag.Parse()

	// Set up logging
	log.SetPrefix("myCalApp: ")
	log.SetFlags(0)

	// Check credentials directory
	if config.GetCredsDirectory() == "" {
		fmt.Printf("Credentials directory not configured. Current working directory will be used.\n Set '%s' env variable to configure\n", config.CredsDirectoryEnv)
	}

	// Get calendar service
	ctx := context.Background()
	srv, err := auth.GetCalendarService(ctx)
	if err != nil {
		log.Fatalf("Failed to get calendar service: %v", err)
	}

	if *watchMode {
		// Interactive TUI mode
		if err := tui.Run(srv); err != nil {
			log.Fatalf("Error running TUI: %v", err)
		}
	} else {
		// Static output mode
		runStaticMode(srv)
	}
}

func runStaticMode(srv *gcal.Service) {
	// Header
	fmt.Println(tui.RenderHeader())

	// Next meeting countdown
	nextEvent, _ := calendar.FetchNextEvent(srv)
	if nextEvent != nil {
		countdown := tui.RenderCountdown(nextEvent)
		if countdown != "" {
			fmt.Println(countdown)
		}
	}

	// Today's events
	fmt.Println()
	fmt.Println(tui.RenderSectionTitle("Today", "📅"))
	todayEvents, err := calendar.FetchTodayEvents(srv)
	if err != nil {
		log.Printf("Error fetching today's events: %v", err)
	} else if len(todayEvents) == 0 {
		fmt.Println(tui.NoEventsStyle.Render("No events remaining today"))
	} else {
		fmt.Println(tui.RenderEventList(todayEvents, true, -1))
	}

	// Upcoming events (if today has few events)
	if len(todayEvents) < 3 {
		upcomingEvents, err := calendar.FetchUpcomingEvents(srv, 5, true)
		if err != nil {
			log.Printf("Error fetching upcoming events: %v", err)
		} else if len(upcomingEvents) > 0 {
			fmt.Println()
			fmt.Println(tui.RenderSectionTitle("Upcoming", "🗓"))
			fmt.Println(tui.RenderEventList(upcomingEvents, false, -1))
		}
	}
}
