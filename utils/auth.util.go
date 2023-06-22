package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"os"
	"strings"
	"time"
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

type UserEmailedDataClaims struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

type UserEmailedData struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

func GenerateTokenFromUserData(data UserEmailedData) (string, error) {
	// expires in 5 minutes
	expriresAt := time.Now().Add(time.Minute * 5).Unix()

	claims := &UserEmailedDataClaims{
		Name:     data.Name,
		Email:    data.Email,
		Password: data.Password,
		Role:     data.Role,
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

func ExtractDataFromUserEmailedDataToken(tokenString string) (*UserEmailedData, error) {
	claims := &UserEmailedDataClaims{}
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

	return &UserEmailedData{
		Name:     claims.Name,
		Email:    claims.Email,
		Password: claims.Password,
		Role:     claims.Role,
	}, nil
}

// GenerateToken generates a jwt token for the user
func GenerateToken(data UserData) (string, error) {

	// expires in 7 days
	expriresAt := time.Now().Add(time.Hour * 24 * 7).Unix()

	claims := &Claims{
		ID:   data.ID,
		Role: data.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expriresAt,
		},
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
