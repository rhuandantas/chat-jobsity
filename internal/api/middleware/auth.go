package middleware

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type jwtCustomClaims struct {
	Email string `json:"email"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

// TODO env var ACCESS_SECRET
var signingKey = []byte("secret")

func GenerateAccessToken(email string) (*string, error) {
	claims := &jwtCustomClaims{
		email,
		true,
		jwt.StandardClaims{
			//TODO read from config file
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString(signingKey)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Couldn't generate access token %s", err.Error()))
	}

	return &t, nil
}
