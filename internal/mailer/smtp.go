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

func (m *Mailer) Send(recipient string, data any, patterns ...string) error {
	//for i := range patterns {
	//	patterns[i] = "emails/" + patterns[i]
	//}
	//
	//msg := mail.NewMessage()
	//msg.SetHeader("To", recipient)
	//msg.SetHeader("From", m.from)
	//
	//ts, err := template.New("").Funcs(funcs.TemplateFuncs).ParseFS(assets.EmbeddedFiles, patterns...)
	//if err != nil {
	//	return err
	//}
	//
	//subject := new(bytes.Buffer)
	//err = ts.ExecuteTemplate(subject, "subject", data)
	//if err != nil {
	//	return err
	//}
	//
	//msg.SetHeader("Subject", subject.String())
	//
	//plainBody := new(bytes.Buffer)
	//err = ts.ExecuteTemplate(plainBody, "plainBody", data)
	//if err != nil {
	//	return err
	//}
	//
	//msg.SetBody("text/plain", plainBody.String())
	//
	//if ts.Lookup("htmlBody") != nil {
	//	htmlBody := new(bytes.Buffer)
	//	err = ts.ExecuteTemplate(htmlBody, "htmlBody", data)
	//	if err != nil {
	//		return err
	//	}
	//
	//	msg.AddAlternative("text/html", htmlBody.String())
	//}
	//
	//for i := 1; i <= 3; i++ {
	//	err = m.dialer.DialAndSend(msg)
	//
	//	if nil == err {
	//		return nil
	//	}
	//
	//	time.Sleep(2 * time.Second)
	//}
	//
	//return err
	return nil
}
