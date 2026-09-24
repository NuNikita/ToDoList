package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

var ErrInvalidToken = errors.New("invalid token")

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func CheckSecretKey() error {
	if len(secretKey) == 0 {
		return errors.New("secret key is required")
	}

	return nil
}

func GenerateToken(userID int) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("sign JWT: %w", err)
	}

	return tokenStr, nil
}

func ParseToken(tokenStr string) (int, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok || t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method %q", t.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return 0, fmt.Errorf("%w: parse JWT: %w", ErrInvalidToken, err)
	}

	if !token.Valid {
		return 0, ErrInvalidToken
	}

	return claims.UserID, nil
}
