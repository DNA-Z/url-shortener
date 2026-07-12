package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

var sugar *zap.SugaredLogger

func InitLogger(logger *zap.Logger) {
	sugar = logger.Sugar()
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI
		method := r.Method

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			bodySize:       0,
		}

		next.ServeHTTP(w, r)

		duration := time.Since(start)

		sugar.Infoln(
			"uri", uri,
			"method", method,
			"status", rw.statusCode,
			"size", rw.bodySize,
			"duration", duration,
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	bodySize   int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(data)
	rw.bodySize += size
	return size, err
}
