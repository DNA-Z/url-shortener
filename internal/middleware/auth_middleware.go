package middleware

import (
	"context"
	"net/http"

	"github.com/DNA-Z/url-shortener/internal/auth"
)

type contextKey string

const userIDKey contextKey = "userID"

var authEnabled bool

func InitAuthMiddleware(enabled bool) {
	authEnabled = enabled
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !authEnabled {
			ctx := context.WithValue(r.Context(), userIDKey, "")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		userID, err := GetUserIDFromCookie(r)

		if err != nil {
			newUserID := auth.GenerateUserID()
			if err := SetUserCookie(w, newUserID); err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			userID = newUserID
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}

func GetUserIDFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("user_jwt")
	if err != nil {
		return "", err
	}

	userID, err := auth.GetUserID(cookie.Value)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func SetUserCookie(w http.ResponseWriter, userID string) error {
	tokenString, err := auth.BuildJWTString(userID)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "user_jwt",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(auth.TokenExpiration.Seconds()),
	})

	return nil
}
