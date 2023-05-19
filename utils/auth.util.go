package utils

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type Claims struct {
	ID uuid.UUID  `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

type Data struct {
	ID uuid.UUID `json:"id"`
	Email string `json:"email"`
}




// GenerateToken generates a jwt token for the user
func GenerateToken(data *Data) (string, error) {
	// Declare the expiration time of the token
	// Here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)

	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		ID: data.ID,
		Email: data.Email,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
			
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

func ExtractDataFromToken(tokenString string) (*Data, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, err
	}

	return &Data{
		ID: claims.ID,
	}, nil
}
