package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"

	"oredavids.com/myCal/internal/config"
)

// GetCalendarService creates and returns an authenticated calendar service
func GetCalendarService(ctx context.Context) (*calendar.Service, error) {
	b, err := os.ReadFile(config.GetCredentialsPath())
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %v", err)
	}

	oauthConfig, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	client := getClient(oauthConfig)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Calendar client: %v", err)
	}

	return srv, nil
}

// getClient retrieves a token, saves the token, then returns the generated client
func getClient(oauthConfig *oauth2.Config) *http.Client {
	tokFile := config.GetTokenPath()
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		fmt.Println("Token required...")
		tok = getTokenFromWeb(oauthConfig)
		saveToken(tokFile, tok)
	}
	return oauthConfig.Client(context.Background(), tok)
}

// getTokenFromWeb requests a token from the web, then returns the retrieved token
func getTokenFromWeb(oauthConfig *oauth2.Config) *oauth2.Token {
	// Channel to receive the auth code
	codeChan := make(chan string)

	// Find an available port
	listener, err := findAvailablePort()
	if err != nil {
		log.Fatalf("Unable to find available port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port

	// Update the redirect URL to use the actual port
	oauthConfig.RedirectURL = fmt.Sprintf("http://localhost:%d", port)

	// Set up handler
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code != "" {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, "<html><body><h1>Authorization successful!</h1><p>You can close this window.</p></body></html>")
			codeChan <- code
		} else {
			fmt.Fprintf(w, "No code received")
		}
	})

	server := &http.Server{Handler: mux}
	go func() {
		if err := server.Serve(listener); err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	// Open browser for authorization
	authURL := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Opening browser for authorization (callback on port %d)...\n", port)
	browser.OpenURL(authURL)

	// Wait for the auth code
	authCode := <-codeChan

	// Shutdown the server
	server.Shutdown(context.Background())

	tok, err := oauthConfig.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// findAvailablePort tries preferred ports first, then falls back to any available port
func findAvailablePort() (net.Listener, error) {
	preferredPorts := []int{3000, 3001, 8080, 8000, 9000}

	for _, port := range preferredPorts {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			return listener, nil
		}
	}

	// Fall back to any available port
	return net.Listen("tcp", ":0")
}

// tokenFromFile retrieves a token from a local file
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

// saveToken saves a token to a file path
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
