package utils

import (
	"errors"
	"net/mail"
	"unicode"
)

// ValidateEmail validates the email
func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

// ValidatePassword validates the password
func ValidatePassword(password string) error {
	if password == "" {
		// return error
		return errors.New("Password is required")
	}
	if len(password) < 8 {
		return errors.New("Password must be at least 8 characters")
	}
	if len(password) > 72 {
		return errors.New("Password must be at most 72 characters")
	}
	// if !checkPassword(password) {
	// 	return errors.New("Password must contain at least 1 uppercase, 1 lowercase, 1 digit and 1 special character")
	// }
	return nil
}

func checkPassword(password string) bool {
	var (
		hasUpper   = false
		hasLower   = false
		hasDigit   = false
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
