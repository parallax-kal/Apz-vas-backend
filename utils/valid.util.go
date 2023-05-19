package utils

import (
	"net/mail"
	"unicode"
)

// ValidateEmail validates the email
func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

// ValidatePassword validates the password
func ValidatePassword(password string) string {
	if password == "" {
		return "Password is required"
	}
	if len(password) < 8 {
		return "Password must be at least 8 characters"
	}
	if !checkPassword(password) {
		return "Password must have at least one uppercase letter, one lowercase letter, one digit and one special character"
	}
	return ""
}

func checkPassword(password string) bool {
	var (
		hasUpper = false
		hasLower = false
		hasDigit = false
		hasSpecial = false
	)
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasDigit && hasSpecial
}

