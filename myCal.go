package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"

	"oredavids.com/greetings"
)

const credsFilePathEnv = "MYCAL_GOOGLE_CALENDAR_CREDENTIALS_FILE_PATH"

// init sets initial values for variables used in the function.
func init() {
	rand.Seed(time.Now().UnixNano())
}

func greet() {
	// A slice of names.
	names := []string{"Ore", "Oreoluwa", "OreDavids"}

	name := names[rand.Intn(len(names))]

	// Request greeting messages for the name
	messages, err := greetings.Hello(name)

	if err != nil {
		log.Fatal(err)
	}

	// If no error was returned, print the returned message to the console.
	log := fmt.Sprintf("\n%s\n", messages)

	fmt.Println(log)
}

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {

	// load the .env file in the current directory, if one exists
	godotenv.Load()

	return os.Getenv(key)
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "myCalAppToken.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser, authorize the app, and then paste JUST the "+
		"authorization code(from the browser url into the CLI): \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

var color = map[string]string{
	"Reset":  "\033[0m",
	"Red":    "\033[31m",
	"Green":  "\033[32m",
	"Yellow": "\033[33m",
	"Blue":   "\033[34m",
	"Purple": "\033[35m",
	"Cyan":   "\033[36m",
	"Gray":   "\033[37m",
	"White":  "\033[97m",
}

func logEventsForToday(calendarService *calendar.Service) {
	now := time.Now()

	// Start at the top of the hour so any meetings started in the same hour show up
	timeToStart := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, now.Nanosecond(), now.Location()).Format(time.RFC3339)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, now.Nanosecond(), now.Location()).Format(time.RFC3339)

	events, err := calendarService.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(timeToStart).TimeMax(endOfDay).OrderBy("startTime").Do()

	if err != nil {
		log.Fatalf("Unable to retrieve events for today: %v", err)
	}

	fmt.Println(color["Green"] + "Upcoming Events for Today:" + color["Reset"])
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found for today.")
	} else {
		for _, eventItem := range events.Items {
			getEventTime(eventItem, true)
		}
	}

	if len(events.Items) < 3 {
		maxNumberOfEvents := 5
		excludeToday := true
		logUpcomingEvents(calendarService, int64(maxNumberOfEvents), excludeToday)
	}
}

func logUpcomingEvents(calendarService *calendar.Service, maxNumberOfEvents int64, excludeToday bool) {
	now := time.Now()

	var timeToStart string
	if excludeToday {
		timeToStart = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, now.Nanosecond(), now.Location()).Format(time.RFC3339)
	} else {
		timeToStart = time.Now().Format(time.RFC3339)
	}

	events, err := calendarService.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(timeToStart).MaxResults(maxNumberOfEvents).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next couple of your events: %v", err)
	}

	fmt.Println(color["Blue"] + "Upcoming Events:" + color["Reset"])
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found!")
	} else {
		for _, eventItem := range events.Items {
			getEventTime(eventItem, false)
		}
	}
}

func getEventTime(event *calendar.Event, isToday bool) time.Time {
	date := event.Start.DateTime
	if date == "" {
		date = event.Start.Date
	}
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		t, err := time.Parse("2006-01-02", date)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%v -- %v\n", event.Summary, t.Format("Monday")+" All day")
		return t
	}

	var formattedTime string
	if isToday {
		formattedTime = t.Format(time.Kitchen)
	} else {
		formattedTime = t.Format("Monday, 3:00PM")
	}

	fmt.Printf("%v -- %v\n", event.Summary, formattedTime)
	return t
}

func main() {
	// Set properties of the predefined Logger, including
	// the log entry prefix and a flag to disable printing
	// the time, source file, and line number.
	log.SetPrefix("myCalApp: ")
	log.SetFlags(0)

	greet()

	credsFilePath := goDotEnvVariable(credsFilePathEnv)

	if credsFilePath == "" {
		log.Fatalf("Env variable %s is required", credsFilePathEnv)
	}

	ctx := context.Background()
	b, err := os.ReadFile(credsFilePath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	logEventsForToday(srv)
}
