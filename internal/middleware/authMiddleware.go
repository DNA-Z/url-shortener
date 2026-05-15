package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/DNA-Z/url-shortener/internal/auth"
)

var authEnabled bool

func InitAuthMiddleware(enabled bool) {
	authEnabled = enabled
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !authEnabled {
			ctx := context.WithValue(r.Context(), "userID", "")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		userID, err := GetUserIDFromCookie(r)

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func PublicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !authEnabled {
			ctx := context.WithValue(r.Context(), "userID", "")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		userId, err := GetUserIDFromCookie(r)

		if err != nil {
			newUserID := auth.GenerateUserID()
			if err := SetUserCookie(w, newUserID); err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			userId = newUserID
		}

		ctx := context.WithValue(r.Context(), "userID", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
		MaxAge:   int(auth.TOKEN_EXP.Seconds()),
	})

	return nil
}
