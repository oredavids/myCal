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
	demoMode := flag.Bool("demo", false, "Run with demo data (for screenshots)")
	themeName := flag.String("theme", "default", "Color theme (default, catppuccin, dracula, nord, tokyonight, gruvbox)")
	flag.BoolVar(new(bool), "themes", false, "List available themes")
	flag.Parse()

	// Handle --themes flag to list themes
	for _, arg := range flag.Args() {
		if arg == "themes" {
			fmt.Println("Available themes:")
			for _, name := range tui.GetThemeNames() {
				fmt.Printf("  - %s\n", name)
			}
			return
		}
	}

	// Set theme
	if !tui.SetTheme(*themeName) {
		fmt.Printf("Unknown theme: %s\nAvailable: default, catppuccin, dracula, nord, tokyonight, gruvbox\n", *themeName)
		return
	}

	// Set up logging
	log.SetPrefix("myCalApp: ")
	log.SetFlags(0)

	// Demo mode doesn't need auth
	if *demoMode {
		runDemoMode()
		return
	}

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
	todayEvents, _ := calendar.FetchTodayEvents(srv)
	nextEvent, _ := calendar.FetchNextEvent(srv)

	var upcomingEvents []*calendar.Event
	if len(todayEvents) < 3 {
		upcomingEvents, _ = calendar.FetchUpcomingEvents(srv, 5, true)
	}

	fmt.Print(tui.RenderStatic(tui.RenderData{
		UserName:       tui.GetUserName(),
		TodayEvents:    todayEvents,
		UpcomingEvents: upcomingEvents,
		NextEvent:      nextEvent,
	}))
}

func runDemoMode() {
	todayEvents, upcomingEvents, nextEvent := calendar.GetDemoEvents()

	fmt.Print(tui.RenderStatic(tui.RenderData{
		UserName:       "acme-user", // Empty for demo
		TodayEvents:    todayEvents,
		UpcomingEvents: upcomingEvents,
		NextEvent:      nextEvent,
	}))
}
