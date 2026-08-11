// Package auth служит для аутентификации и авторизации пользователей.
package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const TokenExpiration = time.Hour * 3

var jwtSecretKey string

func SetJWTSecretKey(key string) {
	jwtSecretKey = key
}

func GetJWTSecretKey() string {
	return jwtSecretKey
}

func IsAuthEnabled() bool {
	return jwtSecretKey != ""
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func BuildJWTString(userID string) (string, error) {
	if !IsAuthEnabled() {
		return "", nil
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpiration)),
		},

		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserID(tokenStrring string) (string, error) {
	if !IsAuthEnabled() {
		return "", nil
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStrring, claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecretKey), nil
		})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("token is not valid")
	}

	return claims.UserID, nil
}

func GenerateUserID() string {
	return uuid.New().String()
}
