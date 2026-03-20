package observability

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCounter(t *testing.T) {
	name := "requests.total"
	value := 42.0
	labels := MetricLabels{"method": "GET", "status": "200"}

	before := time.Now()
	metric := NewCounter(name, value, labels)
	after := time.Now()

	assert.Equal(t, name, metric.Name)
	assert.Equal(t, value, metric.Value)
	assert.Equal(t, MetricTypeCounter, metric.Type)
	assert.Equal(t, labels, metric.Labels)
	assert.True(t, metric.Timestamp.After(before) || metric.Timestamp.Equal(before))
	assert.True(t, metric.Timestamp.Before(after) || metric.Timestamp.Equal(after))
}

func TestNewGauge(t *testing.T) {
	name := "memory.usage"
	value := 1024.5
	labels := MetricLabels{"host": "server1"}

	metric := NewGauge(name, value, labels)

	assert.Equal(t, name, metric.Name)
	assert.Equal(t, value, metric.Value)
	assert.Equal(t, MetricTypeGauge, metric.Type)
	assert.Equal(t, labels, metric.Labels)
	assert.NotZero(t, metric.Timestamp)
}

func TestNewHistogram(t *testing.T) {
	name := "request.duration"
	value := 123.45
	labels := MetricLabels{"endpoint": "/api/v1/users"}

	metric := NewHistogram(name, value, labels)

	assert.Equal(t, name, metric.Name)
	assert.Equal(t, value, metric.Value)
	assert.Equal(t, MetricTypeHistogram, metric.Type)
	assert.Equal(t, labels, metric.Labels)
	assert.NotZero(t, metric.Timestamp)
}

func TestNewCounter_NilLabels(t *testing.T) {
	metric := NewCounter("test.counter", 1.0, nil)

	assert.Equal(t, "test.counter", metric.Name)
	assert.Equal(t, 1.0, metric.Value)
	assert.Equal(t, MetricTypeCounter, metric.Type)
	assert.Nil(t, metric.Labels)
}

func TestNewGauge_EmptyLabels(t *testing.T) {
	labels := MetricLabels{}
	metric := NewGauge("test.gauge", 100.0, labels)

	assert.Equal(t, "test.gauge", metric.Name)
	assert.Equal(t, 100.0, metric.Value)
	assert.Equal(t, MetricTypeGauge, metric.Type)
	assert.NotNil(t, metric.Labels)
	assert.Empty(t, metric.Labels)
}

func TestMetricTypes(t *testing.T) {
	testCases := []struct {
		name       string
		metricType MetricType
		expected   string
	}{
		{name: "counter type", metricType: MetricTypeCounter, expected: "counter"},
		{name: "gauge type", metricType: MetricTypeGauge, expected: "gauge"},
		{name: "histogram type", metricType: MetricTypeHistogram, expected: "histogram"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, string(tc.metricType))
		})
	}
}

func TestMetric_FloatPrecision(t *testing.T) {
	testCases := []struct {
		name  string
		value float64
	}{
		{name: "integer", value: 42.0},
		{name: "decimal", value: 3.14159},
		{name: "large number", value: 1234567.89},
		{name: "small number", value: 0.00001},
		{name: "negative", value: -123.45},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			metric := NewCounter("test", tc.value, nil)
			assert.Equal(t, tc.value, metric.Value)
		})
	}
}
