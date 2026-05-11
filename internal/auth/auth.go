package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const TOKEN_EXP = time.Hour * 3
const SECRET_KEY = "supersecretkey"

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func BuildJWTString(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},

		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserID(tokenStrring string) (string, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStrring, claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(SECRET_KEY), nil
		})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		fmt.Println("Token is not valid")
		return "", err
	}

	fmt.Println("Token is valid")
	return claims.UserID, nil
}

func GenerateUserID() string {
	return uuid.New().String()
}
