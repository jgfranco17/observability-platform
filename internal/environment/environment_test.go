package environment

import (
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestEnvironmentRetrieval(t *testing.T) {
	environments := []string{
		"local",
		"dev",
		"stage",
		"prod",
		"ci",
	}
	for _, sampleEnv := range environments {
		t.Run(fmt.Sprintf("Check environment [%s]", sampleEnv), func(t *testing.T) {
			t.Setenv(ENV_KEY_ENVIRONMENT, sampleEnv)
			valueFromEnv := GetEnvWithDefault(ENV_KEY_ENVIRONMENT, "")
			assert.Equal(t, sampleEnv, valueFromEnv)
		})
	}
}

func TestCheckLocalEnvironment(t *testing.T) {
	testCases := []struct {
		envValue        string
		expectedIsLocal bool
	}{
		{envValue: APPLICATION_ENV_LOCAL, expectedIsLocal: true},
		{envValue: APPLICATION_ENV_DEV, expectedIsLocal: false},
		{envValue: APPLICATION_ENV_STAGE, expectedIsLocal: false},
		{envValue: APPLICATION_ENV_PROD, expectedIsLocal: false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Environment [%s] is local: %t", tc.envValue, tc.expectedIsLocal), func(t *testing.T) {
			t.Setenv(ENV_KEY_ENVIRONMENT, tc.envValue)
			isLocal := IsRunningLocally()
			assert.Equal(t, tc.expectedIsLocal, isLocal)
		})
	}
}

func TestGetLogLevel(t *testing.T) {
	testCases := []struct {
		envValue      string
		expectedLevel logrus.Level
	}{
		{envValue: "DEBUG", expectedLevel: logrus.DebugLevel},
		{envValue: "INFO", expectedLevel: logrus.InfoLevel},
		{envValue: "WARN", expectedLevel: logrus.WarnLevel},
		{envValue: "ERROR", expectedLevel: logrus.ErrorLevel},
		{envValue: "PANIC", expectedLevel: logrus.PanicLevel},
		{envValue: "FATAL", expectedLevel: logrus.FatalLevel},
		{envValue: "TRACE", expectedLevel: logrus.TraceLevel},
		{envValue: "UNKNOWN", expectedLevel: logrus.InfoLevel}, // Default case
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Using log level [%s]", tc.envValue), func(t *testing.T) {
			t.Setenv(ENV_KEY_LOG_LEVEL, tc.envValue)
			logLevel := GetLogLevel()
			assert.Equal(t, tc.expectedLevel, logLevel)
		})
	}
}

func TestGetEnvWithDefaultReturnsDefault(t *testing.T) {
	mockKey := "MOCK_KEY"
	defaultValue := "default"

	result := GetEnvWithDefault(mockKey, defaultValue)

	assert.Equal(t, defaultValue, result)
}

func TestGetLogFormatter(t *testing.T) {
	testCases := []struct {
		name         string
		envValue     string
		expectedType string
	}{
		{
			name:         "JSON formatter",
			envValue:     "JSON",
			expectedType: "*logrus.JSONFormatter",
		},
		{
			name:         "TEXT formatter",
			envValue:     "TEXT",
			expectedType: "*logrus.TextFormatter",
		},
		{
			name:         "lowercase json",
			envValue:     "json",
			expectedType: "*logrus.JSONFormatter",
		},
		{
			name:         "Default formatter",
			envValue:     "DEFAULT",
			expectedType: "*logrus.TextFormatter",
		},
		{
			name:         "Unknown formatter defaults to TextFormatter",
			envValue:     "UNKNOWN",
			expectedType: "*logrus.TextFormatter",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv(ENV_LOG_FORMAT, tc.envValue)
			formatter := GetLogFormatter()

			assert.NotNil(t, formatter)
			assert.Equal(t, tc.expectedType, fmt.Sprintf("%T", formatter))
		})
	}
}
