package gmail

import (
	"github.com/i-mora/notifier/internal/google"
	"google.golang.org/api/gmail/v1"
)

func SendMail(template string) (*gmail.Message, error) {
	svc := google.NewService()

	message := gmail.Message{
		Raw: template,
	}

	return svc.Users.Messages.Send("me", &message).Do()
}
