package utils

import (
	"fmt"
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
	fmt.Println("Email Sent Successfully!")

	return nil
}
