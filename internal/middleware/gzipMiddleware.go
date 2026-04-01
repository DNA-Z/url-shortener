package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type compressWriter struct {
	writer http.ResponseWriter
	gzip   *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		writer: w,
		gzip:   gzip.NewWriter(w),
	}
}

func (cw *compressWriter) Header() http.Header {
	return cw.writer.Header()
}

func (cw *compressWriter) Write(data []byte) (int, error) {
	if isCompressibleType(cw.writer.Header().Get("Content-Type")) {
		cw.writer.Header().Set("Content-Encoding", "gzip")
		return cw.gzip.Write(data)
	}
	cw.writer.Header().Del("Content-Encoding")
	return cw.writer.Write(data)
}

func (cw *compressWriter) WriteHeader(statusCode int) {
	if isCompressibleType(cw.writer.Header().Get("Content-Type")) {
		cw.writer.Header().Set("Content-Encoding", "gzip")
	} else {
		cw.writer.Header().Del("Content-Encoding")
	}

	cw.writer.WriteHeader(statusCode)
}

func (cw *compressWriter) Close() error {
	if isCompressibleType(cw.writer.Header().Get("Content-Type")) {
		return cw.gzip.Close()
	}
	return nil
}

func isCompressibleType(contentType string) bool {
	if strings.HasPrefix(contentType, "application/json") ||
		strings.HasPrefix(contentType, "text/html") {
		return true
	}
	return false
}

type CompressReader struct {
	reader io.ReadCloser
	gzip   *gzip.Reader
}

func newCompressReader(rc io.ReadCloser) (*CompressReader, error) {
	reader, err := gzip.NewReader(rc)
	if err != nil {
		return nil, err
	}
	return &CompressReader{
		gzip:   reader,
		reader: rc,
	}, nil
}

func (cr *CompressReader) Read(p []byte) (n int, err error) {
	return cr.gzip.Read(p)
}

func (cr *CompressReader) Close() error {
	if err := cr.reader.Close(); err != nil {
		return err
	}
	return cr.gzip.Close()
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := w

		contentEncoding := r.Header.Get("Content-Encoding")
		if strings.Contains(contentEncoding, "gzip") {
			compressReader, err := newCompressReader(r.Body)
			if err != nil {
				http.Error(w, "Invalid gzip body", http.StatusBadRequest)
				return
			}
			r.Body = compressReader
			defer compressReader.Close()
		}

		acceptsGzip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
		if acceptsGzip {
			compressWriter := newCompressWriter(w)
			writer = compressWriter
			defer compressWriter.Close()
		}

		next.ServeHTTP(writer, r)
	})
}
