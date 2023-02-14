package mailer

import (
	"github.com/go-mail/mail/v2"
	"test-server-go/internal/config"
	"time"
)

type Mailer struct {
	dialer *mail.Dialer
	from   string
}

func NewSmtp(cfg config.Config) *Mailer {
	dialer := mail.NewDialer(cfg.Smtp1.Host, cfg.Smtp1.Port, cfg.Smtp1.Username, cfg.Smtp1.Password)
	dialer.Timeout = 5 * time.Second

	return &Mailer{
		dialer: dialer,
		from:   cfg.Smtp1.From,
	}
}

// SendEmail function sends an email with provided parameters
func (m *Mailer) SendEmail(to []string, subject string, body string) error {
	msg := mail.NewMessage()
	msg.SetHeader("From", m.from)
	msg.SetHeader("To", to...)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)

	if err := m.dialer.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}
