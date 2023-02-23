package mailer

import (
	"github.com/go-mail/mail/v2"
	"net"
	"strings"
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

// CheckEmailDomainExistence checks if an email domain exists using SPF (Sender Policy Framework) record.
// It extracts the domain from the email address, performs a DNS lookup to get the TXT record of the domain,
// and checks if the TXT record contains the "v=spf1" flag. It returns true if the flag is found, and false otherwise.
func CheckEmailDomainExistence(addr string) (bool, error) {
	// Extract the domain from the email address
	domain := strings.Split(addr, "@")[1]

	// Get the TXT record of the domain
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		return false, err
	}

	// Search for the "v=spf1" flag in the TXT record
	for _, txt := range txtRecords {
		if strings.Contains(txt, "v=spf1") {
			return true, nil
		}
	}

	return false, nil
}
