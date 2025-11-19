package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT payload we use across the app.
type Claims struct {
	UserID uint   `json:"user_id"`
	WorkNo string `json:"work_no"`
	jwt.RegisteredClaims
}

// GenerateToken creates a signed JWT with the provided claims and TTL.
func GenerateToken(secret string, issuer string, ttl time.Duration, userID uint, workNo string) (string, error) {
	claims := Claims{
		UserID: userID,
		WorkNo: workNo,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken validates a token string and returns the claims.
func ParseToken(secret, tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
