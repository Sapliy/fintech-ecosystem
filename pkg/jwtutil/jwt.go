package jwtutil

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// SecretKey is the key used to sign the tokens.
// It is loaded from the JWT_SECRET environment variable.
var SecretKey []byte

func init() {
	SecretKey = []byte(os.Getenv("JWT_SECRET"))
	if len(SecretKey) == 0 {
		if os.Getenv("GO_ENV") == "production" {
			panic("FATAL: JWT_SECRET environment variable must be set in production")
		}
		SecretKey = []byte("super-secret-key-change-me-dev")
	}
}

// Claims defines the custom claims structure for our JWTs.
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT token for a given user ID and email.
// The token expires in 24 hours.
func GenerateToken(userID, email string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "microservices-auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

// ValidateToken parses and validates a JWT token string.
// It returns the claims if the token is valid, or an error otherwise.
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
