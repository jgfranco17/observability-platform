package observability

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewReport_Success(t *testing.T) {
	identifier := "test-app"
	message := "test message"
	level := "info"

	report := NewReport(identifier, message, level)

	assert.Contains(t, report.ID, identifier)
	assert.Contains(t, report.ID, "-")
	assert.Equal(t, message, report.Message)
	assert.Equal(t, level, report.Level)
	assert.NotNil(t, report.Timestamp)
}

func TestNewReport_UniqueIDs(t *testing.T) {
	identifier := "test-app"
	message := "test message"
	level := "info"

	report1 := NewReport(identifier, message, level)
	report2 := NewReport(identifier, message, level)

	assert.NotEqual(t, report1.ID, report2.ID, "Report IDs should be unique")
}

func TestNewReport_DifferentLevels(t *testing.T) {
	testCases := []struct {
		name  string
		level string
	}{
		{name: "info level", level: "info"},
		{name: "warn level", level: "warn"},
		{name: "error level", level: "error"},
		{name: "debug level", level: "debug"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			report := NewReport("test-app", "test message", tc.level)
			assert.Equal(t, tc.level, report.Level)
		})
	}
}
