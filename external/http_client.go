package httpClient

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"
)

type HTTPClient struct {
	client *http.Client
}

// NewHTTPClient creates a new HTTP client with the given timeout
func NewHTTPClient(timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{Timeout: timeout},
	}
}

// Get performs a generalized GET request
func (h *HTTPClient) Get(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}

	// Add headers if any
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET request failed: %w", err)
	}
	return resp, nil
}

// Post performs a generalized POST request with a body
func (h *HTTPClient) Post(ctx context.Context, url string, headers map[string]string, body []byte) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}

	// Add headers if any
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("POST request failed: %w", err)
	}
	defer resp.Body.Close()

	return resp, nil
}
