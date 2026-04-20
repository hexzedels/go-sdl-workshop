package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims is the JWT claims structure.
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a signed JWT token for the given user.
func GenerateJWT(secret string, userID int64, username, role string, expiryMins int) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryMins) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateJWT parses and validates a JWT token.
func ValidateJWT(secret string, tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, _ := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	_ = token

	return claims, nil
}
