package mailer

import (
	"bytes"
	"github.com/go-mail/mail/v2"
	"html/template"
	"log"
	"path/filepath"
	"test-server-go/internal/config"
	tl "test-server-go/internal/tools"
	"time"
)

type Mailer struct {
	*mail.Dialer
	from string
}

func getPath(file string) (string, error) {
	path, err := tl.GetExecutablePath()
	if err != nil {
		log.Fatal()
	}

	dir := filepath.Join(path, "resources", "templates")

	tmplFilePath := filepath.Join(dir, file)

	return tmplFilePath, nil
}

func NewSmtp(cfg config.Config) *Mailer {
	dialer := mail.NewDialer(cfg.MailNoreply.Host, cfg.MailNoreply.Port, cfg.MailNoreply.Username, cfg.MailNoreply.Password)
	dialer.Timeout = 5 * time.Second

	return &Mailer{
		Dialer: dialer,
		from:   cfg.MailNoreply.From,
	}
}

// SendEmail function sends an email with provided parameters
func (m *Mailer) sendEmail(to []string, title, body string) error {
	msg := mail.NewMessage()
	msg.SetHeader("From", m.from)
	msg.SetHeader("To", to...)
	msg.SetHeader("Subject", title)
	msg.SetBody("text/plain", body)

	if err := m.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}

func (m *Mailer) SendEmailConfirmation(email, nickname, confirmationUrl, shopName, clientAppUrl string) error {
	templateFile, err := getPath("mailConfirmationEmail.tmpl")
	if err != nil {
		return err
	}

	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}

	resources := map[string]interface{}{
		"Nickname":         nickname,
		"Email":            email,
		"ConfirmationLink": confirmationUrl,
		"ShopName":         shopName,
		"ClientAppUrl":     clientAppUrl,
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, resources); err != nil {
		return err
	}

	if err = m.sendEmail([]string{email}, shopName+": подтверждение учётной записи", buf.String()); err != nil {
		return err
	}

	return nil
}

func (m *Mailer) SendOrderContent(email, variantName, orderContent, shopName, clientAppUrl string) error {
	templateFile, err := getPath("mailOrder.tmpl")
	if err != nil {
		return err
	}

	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}

	resources := map[string]interface{}{
		"VariantName":  variantName,
		"OrderContent": orderContent,
		"ShopName":     shopName,
		"ClientAppUrl": clientAppUrl,
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, resources); err != nil {
		return err
	}

	if err = m.sendEmail([]string{email}, shopName+": содержимое заказа", buf.String()); err != nil {
		return err
	}

	return nil
}
