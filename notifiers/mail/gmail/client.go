package gmail

import (
	"context"
	"encoding/base64"
	"log"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func SendMail(ctx context.Context, template string) (*gmail.Message, error) {
	client, err := google.DefaultClient(ctx,
		"https://mail.google.com/",
		"https://www.googleapis.com/auth/gmail.modify",
		"https://www.googleapis.com/auth/gmail.compose",
		"https://www.googleapis.com/auth/gmail.send",
	)
	if err != nil {
		log.Fatal(err)
	}

	svc, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create Gmail client: %v", err)
	}
	message := gmail.Message{
		Raw: base64.RawURLEncoding.EncodeToString([]byte(template)),
	}

	return svc.Users.Messages.Send("quesomora@gmail.com", &message).Do()
}
