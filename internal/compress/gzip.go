package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

var compressibleContentTypes = []string{
	"application/javascript",
	"application/json",
	"text/css",
	"text/html",
	"text/plain",
	"text/xml",
}

type CompressWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (cw *CompressWriter) Write(data []byte) (int, error) {
	if !cw.shouldCompress() {
		return cw.ResponseWriter.Write(data)
	}
	return cw.Writer.Write(data)
}

func (cw *CompressWriter) Close() error {
	return cw.Writer.Close()
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
	Closer io.ReadCloser
	*gzip.Reader
}

func NewCompressReader(rc io.ReadCloser) (*CompressReader, error) {
	reader, err := gzip.NewReader(rc)
	if err != nil {
		return nil, err
	}
	return &CompressReader{
		Reader: reader,
		Closer: rc,
	}, nil
}

func (cr *CompressReader) Read(p []byte) (n int, err error) {
	return cr.Closer.Read(p)
}

func (cr *CompressReader) Close() error {
	cr.Reader.Close()
	return cr.Closer.Close()
}
