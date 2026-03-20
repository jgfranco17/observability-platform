package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jgfranco17/observability-platform/internal/config"
	"github.com/jgfranco17/observability-platform/internal/db"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// RouterTestCase represents a single test case for router testing.
type RouterTestCase struct {
	Name           string
	Method         string
	Endpoint       string
	Body           any
	ExpectedStatus int
	ValidateBody   func(t *testing.T, body []byte)
}

// RouterTestsBuilder provides setup dependency injection for router tests.
type RouterTestsBuilder struct {
	router http.Handler
}

// NewTestBuilder creates a new builder instance.
func NewTestBuilder(t *testing.T) *RouterTestsBuilder {
	t.Helper()
	return &RouterTestsBuilder{}
}

// WithRouter sets the HTTP handler to test (dependency injection).
func (b *RouterTestsBuilder) WithRouter(router http.Handler) *RouterTestsBuilder {
	b.router = router
	return b
}

// Run executes all provided test cases against the configured router.
func (b *RouterTestsBuilder) Run(t *testing.T, testCases []RouterTestCase) {
	t.Helper()

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var body io.Reader
			if tc.Body != nil {
				jsonData, err := json.Marshal(tc.Body)
				require.NoError(t, err, "Failed to marshal test body")
				body = bytes.NewBuffer(jsonData)
			}

			req, err := http.NewRequest(tc.Method, tc.Endpoint, body)
			require.NoError(t, err)

			if body != nil {
				req.Header.Set("Content-Type", "application/json")
			}

			rr := httptest.NewRecorder()
			b.router.ServeHTTP(rr, req)

			assert.Equal(t, tc.ExpectedStatus, rr.Code)

			if tc.ValidateBody != nil {
				tc.ValidateBody(t, rr.Body.Bytes())
			}
		})
	}
}

func TestHealthEndpoint(t *testing.T) {
	service := createTestService(t)

	testCases := []RouterTestCase{
		{
			Name:           "should return OK for health check",
			Method:         http.MethodGet,
			Endpoint:       "/api/health",
			ExpectedStatus: http.StatusOK,
		},
	}

	tests := NewTestBuilder(t).WithRouter(service.router)
	tests.Run(t, testCases)
}

func TestTraceEndpoints(t *testing.T) {
	service := createTestService(t)

	validTrace := map[string]interface{}{
		"trace_id":   "trace-123",
		"span_id":    "span-456",
		"parent_id":  "",
		"start_time": "2024-01-01T00:00:00Z",
		"end_time":   "2024-01-01T00:00:01Z",
		"status":     "ok",
		"attributes": map[string]string{"service": "api"},
	}

	testCases := []RouterTestCase{
		{
			Name:           "should accept valid trace",
			Method:         http.MethodPost,
			Endpoint:       "/api/v1/traces",
			Body:           []any{validTrace},
			ExpectedStatus: http.StatusCreated,
		},
		{
			Name:           "should accept empty trace array",
			Method:         http.MethodPost,
			Endpoint:       "/api/v1/traces",
			Body:           []any{},
			ExpectedStatus: http.StatusCreated,
		},
		{
			Name:           "should retrieve all traces",
			Method:         http.MethodGet,
			Endpoint:       "/api/v1/traces",
			ExpectedStatus: http.StatusOK,
		},
		{
			Name:           "should filter traces by trace_id",
			Method:         http.MethodGet,
			Endpoint:       "/api/v1/traces?trace_id=trace-123",
			ExpectedStatus: http.StatusOK,
		},
	}

	tests := NewTestBuilder(t).WithRouter(service.router)
	tests.Run(t, testCases)
}

func TestMetricEndpoints(t *testing.T) {
	service := createTestService(t)

	validMetric := map[string]interface{}{
		"name":      "http_requests_total",
		"type":      "counter",
		"value":     42.0,
		"labels":    map[string]string{"method": "GET", "status": "200"},
		"timestamp": "2024-01-01T00:00:00Z",
	}

	testCases := []RouterTestCase{
		// POST tests
		{
			Name:           "should accept valid metric",
			Method:         http.MethodPost,
			Endpoint:       "/api/v1/metrics",
			Body:           []any{validMetric},
			ExpectedStatus: http.StatusCreated,
		},
		{
			Name:           "should accept multiple metrics",
			Method:         http.MethodPost,
			Endpoint:       "/api/v1/metrics",
			Body:           []any{validMetric, validMetric},
			ExpectedStatus: http.StatusCreated,
		},
		// GET tests
		{
			Name:           "should retrieve all metrics",
			Method:         http.MethodGet,
			Endpoint:       "/api/v1/metrics",
			ExpectedStatus: http.StatusOK,
		},
		{
			Name:           "should filter metrics by name",
			Method:         http.MethodGet,
			Endpoint:       "/api/v1/metrics?name=http_requests_total",
			ExpectedStatus: http.StatusOK,
		},
	}

	tests := NewTestBuilder(t).WithRouter(service.router)
	tests.Run(t, testCases)
}

func TestLogEndpoints(t *testing.T) {
	service := createTestService(t)

	validLog := map[string]interface{}{
		"trace_id":   "trace-123",
		"span_id":    "span-456",
		"level":      "info",
		"message":    "Test log entry",
		"attributes": map[string]string{"key": "value"},
		"timestamp":  "2024-01-01T00:00:00Z",
	}

	errorLog := map[string]interface{}{
		"level":     "error",
		"message":   "Error occurred",
		"timestamp": "2024-01-01T00:00:00Z",
	}

	testCases := []RouterTestCase{
		// POST tests
		{
			Name:           "should accept valid log",
			Method:         http.MethodPost,
			Endpoint:       "/api/v1/logs",
			Body:           []any{validLog},
			ExpectedStatus: http.StatusCreated,
		},
		{
			Name:           "should accept log with different levels",
			Method:         http.MethodPost,
			Endpoint:       "/api/v1/logs",
			Body:           []any{errorLog},
			ExpectedStatus: http.StatusCreated,
		},
		// GET tests
		{
			Name:           "should retrieve all logs",
			Method:         http.MethodGet,
			Endpoint:       "/api/v1/logs",
			ExpectedStatus: http.StatusOK,
		},
		{
			Name:           "should filter logs by trace_id",
			Method:         http.MethodGet,
			Endpoint:       "/api/v1/logs?trace_id=trace-123",
			ExpectedStatus: http.StatusOK,
		},
	}

	tests := NewTestBuilder(t).WithRouter(service.router)
	tests.Run(t, testCases)
}

func TestRouterNotFound(t *testing.T) {
	service := createTestService(t)

	testCases := []RouterTestCase{
		{
			Name:           "should return 404 for unknown routes",
			Method:         http.MethodGet,
			Endpoint:       "/api/unknown",
			ExpectedStatus: http.StatusNotFound,
		},
		{
			Name:           "should return 404 for missing v1 prefix",
			Method:         http.MethodGet,
			Endpoint:       "/api/traces",
			ExpectedStatus: http.StatusNotFound,
		},
		{
			Name:           "should return 404 for unknown POST route",
			Method:         http.MethodPost,
			Endpoint:       "/api/v1/unknown",
			ExpectedStatus: http.StatusNotFound,
		},
	}

	tests := NewTestBuilder(t).WithRouter(service.router)
	tests.Run(t, testCases)
}

func TestMethodNotAllowed(t *testing.T) {
	service := createTestService(t)

	testCases := []RouterTestCase{
		{
			Name:           "should reject PUT on traces endpoint",
			Method:         http.MethodPut,
			Endpoint:       "/api/v1/traces",
			ExpectedStatus: http.StatusMethodNotAllowed,
		},
		{
			Name:           "should reject DELETE on metrics endpoint",
			Method:         http.MethodDelete,
			Endpoint:       "/api/v1/metrics",
			ExpectedStatus: http.StatusMethodNotAllowed,
		},
	}

	tests := NewTestBuilder(t).WithRouter(service.router)
	tests.Run(t, testCases)
}

func createTestService(t *testing.T) *Service {
	t.Helper()
	ctx := context.Background()
	logger := logrus.New()
	logger.SetOutput(io.Discard) // Silence logs during tests

	cfg := config.ServiceSettings{
		Host: "localhost",
		Port: 20000,
	}

	service, err := New(ctx, cfg, logger, db.NewClient)
	require.NoError(t, err, "Failed to create test service")

	return service
}
