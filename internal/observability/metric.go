package observability

import "time"

// MetricType represents the type of metric.
type MetricType string

const (
	MetricTypeCounter   MetricType = "counter"
	MetricTypeGauge     MetricType = "gauge"
	MetricTypeHistogram MetricType = "histogram"
)

type MetricLabels map[string]string

// Metric represents a single metric measurement with labels.
type Metric struct {
	Name      string       `json:"name"`
	Value     float64      `json:"value"`
	Type      MetricType   `json:"type"`
	Labels    MetricLabels `json:"labels,omitempty"`
	Timestamp time.Time    `json:"timestamp"`
}

// NewCounter creates a new counter metric.
func NewCounter(name string, value float64, labels MetricLabels) Metric {
	return Metric{
		Name:      name,
		Value:     value,
		Type:      MetricTypeCounter,
		Labels:    labels,
		Timestamp: time.Now(),
	}
}

// NewGauge creates a new gauge metric.
func NewGauge(name string, value float64, labels MetricLabels) Metric {
	return Metric{
		Name:      name,
		Value:     value,
		Type:      MetricTypeGauge,
		Labels:    labels,
		Timestamp: time.Now(),
	}
}

// NewHistogram creates a new histogram metric.
func NewHistogram(name string, value float64, labels MetricLabels) Metric {
	return Metric{
		Name:      name,
		Value:     value,
		Type:      MetricTypeHistogram,
		Labels:    labels,
		Timestamp: time.Now(),
	}
}
