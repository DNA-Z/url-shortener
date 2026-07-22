// Package compress реализована библиотека для сжатия ответов в формате gzip.
package compress

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

func NewCompressWriter(w http.ResponseWriter) *compressWriter {
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

type compressReader struct {
	reader io.ReadCloser
	gzip   *gzip.Reader
}

func NewCompressReader(rc io.ReadCloser) (*compressReader, error) {
	reader, err := gzip.NewReader(rc)
	if err != nil {
		return nil, err
	}
	return &compressReader{
		gzip:   reader,
		reader: rc,
	}, nil
}

func (cr *compressReader) Read(p []byte) (n int, err error) {
	return cr.gzip.Read(p)
}

func (cr *compressReader) Close() error {
	if err := cr.reader.Close(); err != nil {
		return err
	}
	return cr.gzip.Close()
}
