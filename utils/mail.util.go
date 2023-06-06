package utils

import (
	"net/smtp"
)

// SendMail sends mail
func SendMail(email string, subject string, body string) error {
	from := "gthecoderkalisaineza@gmail.com"
	password := "KALISA123."
	to := []string{
		email,
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message.
	message := []byte("This is a test email message.")
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return err
	}

	return nil
}
