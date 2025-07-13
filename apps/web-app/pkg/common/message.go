package common

import (
	"fmt"
	"time"
)

// GetTimestampedMessage returns a message with the current timestamp
func GetTimestampedMessage(prefix string, message string) string {
	return fmt.Sprintf("[%s] %s: %s", time.Now().Format(time.RFC3339), prefix, message)
}

// GetVersion returns the common library version
func GetVersion() string {
	return "v1.1.0"
}

// FormatMessage creates a nicely formatted message with emoji
func FormatMessage(level, message string) string {
	emoji := "ℹ️"
	switch level {
	case "info":
		emoji = "ℹ️"
	case "warn":
		emoji = "⚠️"
	case "error":
		emoji = "❌"
	case "success":
		emoji = "✅"
	}
	return fmt.Sprintf("%s %s", emoji, message)
}
