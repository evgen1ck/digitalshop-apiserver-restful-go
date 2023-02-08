package mailer

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/smtp"
)

const confirmationURL = "http://localhost:8080/confirm?token="

// ConfirmationData is the data structure to store the confirmation token and email
type ConfirmationData struct {
	Token string
	Email string
}

// GenerateConfirmationToken generates a random token for email confirmation
func GenerateConfirmationToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	hash := sha256.Sum256(bytes)
	return hex.EncodeToString(hash[:]), nil
}

// SendConfirmationEmail sends an email to the user with the confirmation link
func SendConfirmationEmail(confirmationData *ConfirmationData) error {
	to := confirmationData.Email
	token := confirmationData.Token
	body := "To confirm your email, please click on the following link: " + confirmationURL + token
	subject := "Email Confirmation"

	err := smtp.SendMail("smtp.example.com:587",
		smtp.PlainAuth("", "from@example.com", "password", "smtp.example.com"),
		"from@example.com", []string{to}, []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)))

	if err != nil {
		return err
	}
	return nil
}
