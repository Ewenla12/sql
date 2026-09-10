package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTSecret is set once at startup from config. Handlers read it from here
// instead of every function needing the secret passed in as an argument.
var JWTSecret string

// Claims holds the data embed inside every JWT that is issued.
type Claims struct {
	Subject string `json:"sub"`
	Role    string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken creates a signed JWT for the given subject and role,
// valid for 24 hours.
func GenerateToken(subject, role string) (string, error) {
	claims := Claims{
		Subject: subject,
		Role:    role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTSecret))
}

// ParseToken validates a JWT string and returns its claims if valid.
func ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
