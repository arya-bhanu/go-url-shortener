package helpers

import (
	"os"
	"testing"
)

func TestRemoveDomainError(t *testing.T) {
	appDomain := "localhost:9000"
	os.Setenv("APP_DOMAIN", appDomain)

	tests := []struct {
		name     string
		inputURL string
		expected bool
	}{
		{
			name:     "Exact match with domain",
			inputURL: "localhost:9000",
			expected: false,
		},
		{
			name:     "Match with http prefix",
			inputURL: "http://localhost:9000",
			expected: false,
		},
		{
			name:     "Match with https prefix",
			inputURL: "https://localhost:9000",
			expected: false,
		},
		{
			name:     "Match with www and http",
			inputURL: "http://www.localhost:9000",
			expected: false,
		},
		{
			name:     "Different domain",
			inputURL: "https://example.com",
			expected: true,
		},
		{
			name:     "Different domain",
			inputURL: "example.com",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveDomainError(tt.inputURL)
			if result != tt.expected {
				t.Errorf("RemoveDomainError(%q) = %v; want %v", tt.inputURL, result, tt.expected)
			}
		})
	}
}
