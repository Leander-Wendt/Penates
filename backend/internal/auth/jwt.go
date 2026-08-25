package auth

import (
	"errors"
	"time"

	"github.com/Leander-Wendt/Penates/backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

// ErrInvalidToken is returned when a JWT cannot be parsed or fails validation.
var ErrInvalidToken = errors.New("invalid token")

// Claims are the custom JWT claims carried by Penates access tokens.
type Claims struct {
	UserID uint        `json:"userId"`
	Role   models.Role `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken creates a signed JWT for the given user, embedding their id and role.
// It returns the signed token string and its expiry time.
func GenerateToken(secret string, userID uint, role models.Role, expiry time.Duration) (string, time.Time, error) {
	expiresAt := time.Now().Add(expiry)
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, expiresAt, nil
}

// ParseToken validates the given JWT string and returns its claims.
// It returns ErrInvalidToken if the token is malformed, unsigned with the wrong secret, or expired.
func ParseToken(secret string, tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
