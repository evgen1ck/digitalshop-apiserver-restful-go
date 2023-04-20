package mailer

import (
	"bytes"
	"fmt"
	"github.com/go-mail/mail/v2"
	"html/template"
	"os"
	"path/filepath"
	"runtime"
	"test-server-go/internal/config"
	"time"
)

type Mailer struct {
	*mail.Dialer
	from string
}

const folderName = "templates"

func getPath(file string) (string, error) {
	var tmplFilePath string

	if runtime.GOOS == "windows" {
		if err := os.MkdirAll(folderName, os.ModePerm); err != nil {
			return tmplFilePath, fmt.Errorf("failed to create log directory: %v", err)
		}

		tmplFilePath = filepath.Join(folderName, file)
	} else {
		executablePath, err := os.Executable()
		if err != nil {
			return tmplFilePath, fmt.Errorf("failed to get executable path: %v", err)
		}
		executableDir := filepath.Dir(executablePath)

		logsPath := filepath.Join(executableDir, folderName)
		if err := os.MkdirAll(logsPath, os.ModePerm); err != nil {
			return tmplFilePath, fmt.Errorf("failed to create log directory: %v", err)
		}

		tmplFilePath = filepath.Join(logsPath, file)
	}

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

func (m *Mailer) SendEmailConfirmation(nickname, email, confirmationUrl string) error {
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
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, resources); err != nil {
		return err
	}

	if err := m.sendEmail([]string{email}, "Evgenick's Digitals: подтверждение учётной записи", buf.String()); err != nil {
		return err
	}

	return nil
}
