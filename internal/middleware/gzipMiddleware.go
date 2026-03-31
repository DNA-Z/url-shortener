package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"
)

var (
	gzipWriterPool = sync.Pool{
		New: func() interface{} {
			w, _ := gzip.NewWriterLevel(nil, gzip.BestSpeed)
			return w
		},
	}

	compressibleContentTypes = []string{
		"application/javascript",
		"application/json",
		"text/css",
		"text/html",
		"text/plain",
		"text/xml",
	}
)

type CompressWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (cw *CompressWriter) Write(data []byte) (int, error) {
	if !cw.shouldCompress() {
		return cw.ResponseWriter.Write(data)
	}

	if cw.Header().Get("Content-Encoding") == "" {
		cw.Header().Set("Content-Encoding", "gzip")
		cw.Header().Del("Content-Length")
	}

	if cw.Writer == nil {
		cw.Writer = gzipWriterPool.Get().(*gzip.Writer)
		cw.Writer.Reset(cw.ResponseWriter)
	}

	return cw.Writer.Write(data)
}

func (cw *CompressWriter) Close() error {
	if cw.Writer != nil {
		cw.Writer.Close()
		gzipWriterPool.Put(cw.Writer)
		cw.Writer = nil
	}
	return nil
}

func (cw *CompressWriter) shouldCompress() bool {
	if cw.Header().Get("Content-Encoding") != "" {
		return false
	}

	contentType := cw.Header().Get("Content-Type")
	for _, ct := range compressibleContentTypes {
		if strings.HasPrefix(contentType, ct) {
			return true
		}
	}
	return false
}

type CompressReader struct {
	*gzip.Reader
	rc io.ReadCloser
}

func NewCompressReader(r io.ReadCloser) (*CompressReader, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &CompressReader{Reader: gz, rc: r}, nil
}

func (cr *CompressReader) Close() error {
	cr.Reader.Close()
	return cr.rc.Close()
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := w

		sendsGzip := r.Header.Get("Content-Encoding") == "gzip"
		if sendsGzip {
			cr, err := NewCompressReader(r.Body)
			if err != nil {
				http.Error(w, "Invalid gzip body", http.StatusBadRequest)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		acceptsGzip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
		supportsGzip := acceptsGzip && (r.ProtoMajor == 1 && r.ProtoMinor >= 1 || r.ProtoMajor > 1)
		if supportsGzip {
			cw := &CompressWriter{ResponseWriter: w}
			writer = cw
			defer cw.Close()
		}

		next.ServeHTTP(writer, r)
	})
}
