package handler

import (
	"net/http"
	"net/http/httptest"
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
			h := NewURLHandler(tt.svc)
			h.GetByIDGet(tt.res, tt.req)
			rr := tt.res.(*httptest.ResponseRecorder)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedHeader, rr.Header().Get("Location"))
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}
