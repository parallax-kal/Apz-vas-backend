package utils

import (
	"errors"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type Claims struct {
	ID   uuid.UUID `json:"id"`
	Role string    `json:"role"`
	jwt.StandardClaims
}

type UserData struct {
	ID   uuid.UUID `json:"id"`
	Role string    `json:"role"`
}

// GenerateToken generates a jwt token for the user
func GenerateToken(data UserData) (string, error) {

	claims := &Claims{
		ID:   data.ID,
		Role: data.Role,
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ExtractDataFromToken(tokenString string) (*UserData, error) {
	claims := &Claims{}
	// it is bearer token
	// split at bearer
	tokenSplit := strings.Split(tokenString, "Bearer ")[1]
	token, err := jwt.ParseWithClaims(tokenSplit, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, errors.New("Invalid token")
	}
	if !token.Valid {
		return nil, errors.New("Invalid token")
	}

	return &UserData{
		ID:   claims.ID,
		Role: claims.Role,
	}, nil
}
