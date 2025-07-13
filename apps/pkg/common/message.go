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
	return "v1.0.0"
}