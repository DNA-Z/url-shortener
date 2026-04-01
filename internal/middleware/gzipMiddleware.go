package middleware

import (
	"net/http"
	"strings"

	"github.com/DNA-Z/url-shortener/internal/compress"
)

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := w

		sendsGzip := r.Header.Get("Content-Encoding") == "gzip"
		if sendsGzip {
			compressReader, err := compress.NewCompressReader(r.Body)
			if err != nil {
				http.Error(w, "Invalid gzip body", http.StatusBadRequest)
				return
			}
			r.Body = compressReader
			defer compressReader.Close()
		}

		acceptsGzip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
		if acceptsGzip {
			compressWriter := &compress.CompressWriter{ResponseWriter: w}
			writer = compressWriter
			defer compressWriter.Close()
		}

		next.ServeHTTP(writer, r)
	})
}
