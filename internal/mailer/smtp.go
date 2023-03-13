package mailer

import (
	"bytes"
	"github.com/go-mail/mail/v2"
	"html/template"
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
func (m *Mailer) sendEmail(to []string, title, body string) error {
	msg := mail.NewMessage()
	msg.SetHeader("From", m.from)
	msg.SetHeader("To", to...)
	msg.SetHeader("Subject", title)
	msg.SetBody("text/plain", body)

	if err := m.dialer.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}

func (m *Mailer) SendEmailConfirmation(nickname, email, confirmationUrl string) error {
	tmpl, err := template.ParseFiles("tmplEmailConfirmation.tmpl")
	if err != nil {
		return err
	}

	resources := map[string]interface{}{
		"Nickname":         nickname,
		"Email":            email,
		"ConfirmationLink": confirmationUrl,
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "tmplEmailConfirmation.tmpl", resources)
	if err != nil {
		return err
	}

	err = m.sendEmail([]string{email}, "Evgenick's Digitals: подтверждение учётной записи", buf.String())
	if err != nil {
		return err
	}

	return nil
}
