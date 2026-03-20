package observability

import (
	"testing"

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
			client, err := NewClient("mock-system", tc.baseURL)

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
