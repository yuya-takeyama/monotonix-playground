package common

import (
	"strings"
	"testing"
)

func TestGetTimestampedMessage(t *testing.T) {
	tests := []struct {
		name     string
		prefix   string
		message  string
		contains []string
	}{
		{
			name:     "basic message",
			prefix:   "INFO",
			message:  "Hello World",
			contains: []string{"INFO", "Hello World", "["},
		},
		{
			name:     "empty prefix",
			prefix:   "",
			message:  "Test message",
			contains: []string{": Test message"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetTimestampedMessage(tt.prefix, tt.message)
			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("GetTimestampedMessage() = %v, want substring %v", result, substr)
				}
			}
		})
	}
}

func TestGetVersion(t *testing.T) {
	expected := "v1.0.0"
	if got := GetVersion(); got != expected {
		t.Errorf("GetVersion() = %v, want %v", got, expected)
	}
}