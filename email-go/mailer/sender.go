package mailer

import (
	"mail-go/mailer/notifier"
	"mail-go/config"
)

type Sender interface {
	Send(payload []byte) (string, error)
}

const (
	APIKEY = "MAILGUN_API_KEY"
	DOMAIN = "MAILGUN_DOMAIN_NAME"
	FROM   = "MAILGUN_FROM"
)

// Sender ..
func NewSender() Sender {
	apiKey := config.GetConfig(APIKEY)
	domain := config.GetConfig(DOMAIN)
	from := config.GetConfig(FROM)
	return notifier.MailgunNotifier{
		ApiKey: apiKey,
		Domain: domain,
		From:   from,
	}
}
