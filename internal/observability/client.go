package observability

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/jgfranco17/observability-platform/internal/logging"
	"github.com/sirupsen/logrus"
)

type Client struct {
	baseURL url.URL
	doer    httpDoer
}

// NewClient creates a new observability client with the given base URL.
func NewClient(baseURL string) (*Client, error) {
	parsedURL, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	return &Client{
		baseURL: *parsedURL,
		doer: &http.Client{
			Timeout: 5 * time.Second,
			Jar:     &jar{},
		},
	}, nil
}

// Send sends an observability report to the configured endpoint.
func (c *Client) Send(ctx context.Context, report Report) error {
	resp, err := c.makeJSONRequest(ctx, "/reports", report)
	if err != nil {
		return fmt.Errorf("failed to send report: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

// makeJSONRequest is a helper method to send a JSON-encoded POST request to the
// specified platform endpoint.
func (c *Client) makeJSONRequest(ctx context.Context, endpoint string, data any) (*http.Response, error) {
	logger := logging.FromContext(ctx)

	url := c.baseURL.JoinPath(endpoint).String()
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	switch {
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		logger.WithFields(logrus.Fields{
			"endpoint": endpoint,
			"status":   resp.StatusCode,
		}).Trace("Request successful")
	case resp.StatusCode == 401 || resp.StatusCode == 403:
		return resp, fmt.Errorf("access denied, service returned HTTP %d", resp.StatusCode)
	default:
		return resp, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return resp, nil
}
