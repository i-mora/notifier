package google

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

const (
	viperGmailCredentialsJSON = "gmail.credentials.json"
	viperGmailToken           = "gmail.token"
)

func init() {
	viper.BindEnv(viperGmailCredentialsJSON, "GMAIL_CREDENTIALS_JSON")
	viper.BindEnv(viperGmailToken, "GMAIL_TOKEN")
}

func NewService() *gmail.Service {
	ctx := context.Background()

	credentials := viper.GetString(viperGmailCredentialsJSON)

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON([]byte(credentials), gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	return srv
}

// getClient retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	rawToken := viper.GetString(viperGmailToken)

	tok := &oauth2.Token{}
	err := json.NewDecoder(strings.NewReader(rawToken)).Decode(tok)
	if err != nil {
		log.Fatalf("Unable to decode Gmail token: %v", err)
	}
	return config.Client(context.Background(), tok)
}
