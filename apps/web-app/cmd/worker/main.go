package main

import (
	"fmt"
	"time"

	"github.com/yuya-takeyama/monotonix-playground/apps/web-app/pkg/common"
)

func main() {
	fmt.Println(common.FormatMessage("info", "Starting background worker"))

	// Simulate background processing
	for {
		fmt.Println(common.FormatMessage("info", "Processing background tasks..."))

		// Example: Process some tasks
		taskIDs := []string{"task1", "task2", "task3"}
		for _, taskID := range taskIDs {
			message := common.GetTimestampedMessage("WORKER", fmt.Sprintf("Processing %s", taskID))
			fmt.Println(common.FormatMessage("success", message))
		}

		fmt.Println(common.FormatMessage("info", "Background task completed. Sleeping for 30 seconds..."))
		time.Sleep(30 * time.Second)
	}
}
