package auth

import (
	"fmt"
	"tasklist/internal/models"
	"tasklist/pkg/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}
type TokenData struct {
	UserID int `json:"user_id"`
}

func GenerateToken(userID int, cfg *config.Config) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.JwtTtlMin) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JwtSecret))
}

func ParseToken(tokenStr string) (*TokenData, error) {
	var claims TokenData
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		claims.UserID = int(tokenClaims[models.UserCtxKey].(float64))
		return &claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
