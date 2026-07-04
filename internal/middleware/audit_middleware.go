package middleware

import (
	"context"
	"net/http"

	"github.com/DNA-Z/url-shortener/internal/audit"
)

type auditContextKey string

const AuditPublisherKey auditContextKey = "audit_publisher"

func AuditMiddleware(publisher audit.IPublisher) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), AuditPublisherKey, publisher)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetAuditPublisher(ctx context.Context) *audit.Publisher {
	if publisher, ok := ctx.Value(AuditPublisherKey).(*audit.Publisher); ok {
		return publisher
	}
	return nil
}
