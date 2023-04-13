package infrastructure

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/Nekodigi/class-post-web-backend/config"
	"golang.org/x/oauth2"
)

func GetClient(conf *config.Config) *http.Client {
	//ctx := context.Background()

	// b, err := os.ReadFile("credentials/credentials.json")
	// if err != nil {
	// 	log.Fatalf("Unable to read client secret file: %v", err)
	// }

	// //If modifying these scopes, delete your previously saved token.json.
	// config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)

	// if err != nil {
	// 	log.Fatalf("Unable to parse client secret file to config: %v", err)
	// }

	endpoint := &oauth2.Endpoint{AuthURL: "https://accounts.google.com/o/oauth2/auth", TokenURL: "https://oauth2.googleapis.com/token"}
	config := &oauth2.Config{ClientID: conf.ClientId, ClientSecret: conf.ClientSecret, Endpoint: *endpoint}
	client := getClient_(conf, config)

	return client
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient_(conf *config.Config, config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tok := &oauth2.Token{TokenType: "authorized_user", RefreshToken: conf.RefreshToken}

	// tokFile := "credentials/token.json"
	// tok, err := tokenFromFile(tokFile)

	// if err != nil {
	// 	fmt.Errorf("No token found. check https://developers.google.com/calendar/api/quickstart/go")
	// 	// tok = getTokenFromWeb(config)
	// 	// saveToken(tokFile, tok)
	// }
	return config.Client(context.Background(), tok)
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
