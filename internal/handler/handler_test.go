package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DNA-Z/url-shortener/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestURLHandler_GetByIDGet(t *testing.T) {
	tests := []struct {
		name           string
		svc            *service.URL
		req            *http.Request
		res            http.ResponseWriter
		expectedStatus int
		expectedHeader string
		expectedBody   string
	}{
		{
			name: "successful redirect with valid ID",
			svc: func() *service.URL {
				svc := service.NewURL()
				svc.Urls = map[string]string{
					"123": "https://example.com",
				}
				return svc
			}(),
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/url/123", nil)
				req.SetPathValue("id", "123")
				return req
			}(),
			res:            httptest.NewRecorder(),
			expectedStatus: http.StatusTemporaryRedirect,
			expectedHeader: "https://example.com",
			expectedBody:   "",
		},
		{
			name: "non-existent ID returns error",
			svc: func() *service.URL {
				svc := service.NewURL()
				svc.Urls = map[string]string{
					"existing-id": "https://example.com",
				}
				return svc
			}(),
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/url/non-existent", nil)
				req.SetPathValue("id", "non-existent")
				return req
			}(),
			res:            httptest.NewRecorder(),
			expectedStatus: http.StatusBadRequest,
			expectedHeader: "",
			expectedBody:   "URL not found\n",
		},
		{
			name: "empty URLs map returns error",
			svc: func() *service.URL {
				return service.NewURL()
			}(),
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/url/any-id", nil)
				req.SetPathValue("id", "any-id")
				return req
			}(),
			res:            httptest.NewRecorder(),
			expectedStatus: http.StatusBadRequest,
			expectedHeader: "",
			expectedBody:   "URL not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewURLHandler(tt.svc, "localhost:8080", "http://localhost:8080/")
			h.GetByIDGet(tt.res, tt.req)
			rr := tt.res.(*httptest.ResponseRecorder)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedHeader, rr.Header().Get("Location"))
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestURLHandler_ShortenerPost(t *testing.T) {
	tests := []struct {
		name        string
		url1        string
		url2        string
		shouldMatch bool
	}{
		{
			name:        "different URLs generate different IDs",
			url1:        "https://example.com",
			url2:        "https://google.com",
			shouldMatch: false,
		},
		{
			name:        "same URL returns same ID",
			url1:        "https://example.com",
			url2:        "https://example.com",
			shouldMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := service.NewURL()
			h := NewURLHandler(svc, "localhost:8080", "http://localhost:8080/")

			req1 := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBufferString(tt.url1))
			req1.Host = "localhost:8080"
			rr1 := httptest.NewRecorder()
			h.ShortenerPost(rr1, req1)

			req2 := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBufferString(tt.url2))
			req2.Host = "localhost:8080"
			rr2 := httptest.NewRecorder()
			h.ShortenerPost(rr2, req2)

			id1 := strings.TrimPrefix(rr1.Body.String(), "http://localhost:8080/")
			id2 := strings.TrimPrefix(rr2.Body.String(), "http://localhost:8080/")

			if tt.shouldMatch {
				assert.Equal(t, id1, id2, "IDs should match for duplicate URLs")
			} else {
				assert.NotEqual(t, id1, id2, "IDs should be different for different URLs")
			}
		})
	}
}
