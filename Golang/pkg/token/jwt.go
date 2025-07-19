package token

import (
	"time"

	"my-go-project/internal/common"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID uint, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(common.SecretKey))
}
