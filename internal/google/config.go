package google

import (
	"context"
	"log"
	"os"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

const (
	viperGmailAccessToken     = "gmail.access.token"
	viperGmailCredentialsJSON = "gmail.credentials.json"
)

func init() {
	viper.BindEnv(viperGmailAccessToken, "GMAIL_ACCESS_TOKEN")
	viper.BindEnv(viperGmailCredentialsJSON, "GMAIL_CREDENTIALS_JSON")
}

func NewService() *gmail.Service {
	ctx := context.Background()

	b, err := os.ReadFile(viper.GetString(viperGmailCredentialsJSON))
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	token := oauth2.Token{
		AccessToken: viper.GetString(viperGmailAccessToken),
	}

	srv, err := gmail.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, &token)))
	if err != nil {
		log.Fatalf("Unable to create Gmail client: %v", err)
	}

	return srv
}
