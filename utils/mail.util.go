package utils

import (
	"encoding/base64"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

var layout = "2006-01-02 15:04:05.000000000 -0700 MST"

var (
	oauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			gmail.GmailComposeScope,
			gmail.GmailSendScope,
			gmail.GmailModifyScope,
			gmail.MailGoogleComScope,
		},
		Endpoint: google.Endpoint,
	}
	fromEmail = os.Getenv("GOOGLE_FROM_EMAIL")
)

type EmailClaims struct {
	Email string    `json:"email"`
	ID    uuid.UUID `json:"id"`
	jwt.StandardClaims
}

type EmailData struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

func GenerateEmailToken(email string, id uuid.UUID) (string, error) {
	expriresAt := time.Now().Add(time.Minute * 5).Unix()
	claims := &EmailClaims{
		Email: email,
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expriresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ExtractEmailData(token string) (*EmailData, error) {
	claims := &EmailClaims{}
	// it is bearer token
	// split at bearer
	tokenSplit := strings.Split(token, "Bearer ")[1]
	tokenData, err := jwt.ParseWithClaims(tokenSplit, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, errors.New("Invalid token")
	}
	if !tokenData.Valid {
		return nil, errors.New("Invalid token")
	}

	return &EmailData{
		Email: claims.Email,
		ID: claims.ID,
	}, nil
}

// SendMail sends mail
func SendMail(toEmail string, subject string, message string) error {
	t, err := time.Parse(layout, os.Getenv("GOOGLE_TOKEN_EXPIRY"))
	if err != nil {
		return err
	}

	// Get the Unix timestamp in int64
	var token = &oauth2.Token{
		RefreshToken: os.Getenv("GOOGLE_REFRESH_TOKEN"),
		AccessToken:  os.Getenv("GOOGLE_ACCESS_TOKEN"),
		Expiry:       t,
	}
	client := oauthConfig.Client(oauth2.NoContext, token)

	srv, err := gmail.NewService(oauth2.NoContext, option.WithHTTPClient(client))

	if err != nil {
		return err
	}

	var gmailmessage gmail.Message

	emailStr := []byte(
		"From: " + fromEmail + "\r\n" +
			"To: " + toEmail + "\r\n" +
			"Subject: " + subject + "\r\n\r\n" +
			message + "\r\n")

	gmailmessage.Raw = base64.URLEncoding.EncodeToString(emailStr)

	_, err = srv.Users.Messages.Send("me", &gmailmessage).Do()

	if err != nil {
		return err
	}
	return nil

}
