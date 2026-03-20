package observability

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jgfranco17/observability-platform/internal/logging"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	testCases := []struct {
		name      string
		baseURL   string
		expectErr bool
	}{
		{
			name:      "valid HTTP URL",
			baseURL:   "http://localhost:8080",
			expectErr: false,
		},
		{
			name:      "valid HTTPS URL",
			baseURL:   "https://api.example.com",
			expectErr: false,
		},
		{
			name:      "invalid URL",
			baseURL:   "not a valid url",
			expectErr: true,
		},
		{
			name:      "empty URL",
			baseURL:   "",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client, err := NewClient(tc.baseURL)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, client)
				assert.Equal(t, tc.baseURL, client.baseURL.String())
			}
		})
	}
}

func TestClientSend_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/reports", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var report Report
		err := json.NewDecoder(r.Body).Decode(&report)
		require.NoError(t, err)
		assert.Equal(t, "test-report", report.ID)
		assert.Equal(t, "test message", report.Message)

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"status":"success"}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	require.NoError(t, err)

	var buf bytes.Buffer
	logger := logging.New(&buf, logrus.TraceLevel)
	ctx := logging.AddToContext(context.Background(), logger)

	report := Report{
		ID:      "test-report",
		Message: "test message",
		Level:   "info",
	}

	err = client.Send(ctx, report)
	assert.NoError(t, err)
}

func TestClientSend_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"unauthorized"}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	require.NoError(t, err)

	var buf bytes.Buffer
	logger := logging.New(&buf, logrus.TraceLevel)
	ctx := logging.AddToContext(context.Background(), logger)

	report := Report{
		ID:      "test-report",
		Message: "test message",
		Level:   "info",
	}

	err = client.Send(ctx, report)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "access denied")
}

func TestClientSend_Forbidden(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error":"forbidden"}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	require.NoError(t, err)

	var buf bytes.Buffer
	logger := logging.New(&buf, logrus.TraceLevel)
	ctx := logging.AddToContext(context.Background(), logger)

	report := Report{
		ID:      "test-report",
		Message: "test message",
		Level:   "info",
	}

	err = client.Send(ctx, report)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "access denied")
}

func TestClientSend_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	require.NoError(t, err)

	var buf bytes.Buffer
	logger := logging.New(&buf, logrus.TraceLevel)
	ctx := logging.AddToContext(context.Background(), logger)

	report := Report{
		ID:      "test-report",
		Message: "test message",
		Level:   "info",
	}

	err = client.Send(ctx, report)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status code")
}

func TestClientSend_BadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"bad request"}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	require.NoError(t, err)

	var buf bytes.Buffer
	logger := logging.New(&buf, logrus.TraceLevel)
	ctx := logging.AddToContext(context.Background(), logger)

	report := Report{
		ID:      "test-report",
		Message: "test message",
		Level:   "info",
	}

	err = client.Send(ctx, report)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status code")
}
