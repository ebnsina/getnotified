package notify

import (
	"context"
	"fmt"
	"net/smtp"
	"os"

	"github.com/fajrlabs/getnotified/internal/store"
)

type email struct{}

func (email) Send(_ context.Context, e Event, c store.Channel) error {
	host, port := os.Getenv("SMTP_HOST"), os.Getenv("SMTP_PORT")
	from := os.Getenv("SMTP_FROM")
	if host == "" || port == "" || from == "" {
		return fmt.Errorf("email channel needs SMTP_HOST, SMTP_PORT and SMTP_FROM")
	}

	var auth smtp.Auth
	if user := os.Getenv("SMTP_USER"); user != "" {
		auth = smtp.PlainAuth("", user, os.Getenv("SMTP_PASS"), host)
	}
	to := cfg(c, "to")
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		from, to, e.Subject(), e.Body())
	return smtp.SendMail(host+":"+port, auth, from, []string{to}, []byte(msg))
}
