package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type HTTPObserver struct {
	url        string
	client     *http.Client
	timeout    time.Duration
	maxRetries int
}

func NewHTTPObserver(url string) *HTTPObserver {
	return &HTTPObserver{
		url: url,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		timeout:    5 * time.Second,
		maxRetries: 3,
	}
}

func (h *HTTPObserver) Notify(event AuditEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	var lastErr error
	for i := 0; i < h.maxRetries; i++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.url, bytes.NewReader(data))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := h.client.Do(req)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return nil
			}
			lastErr = err
		} else {
			lastErr = err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(i+1) * 100 * time.Millisecond):
		}
	}
	return lastErr
}
