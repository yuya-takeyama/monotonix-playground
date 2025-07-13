package main

import (
	"fmt"
	"time"

	"github.com/yuya-takeyama/monotonix-playground/apps/foo/internal/service"
	"github.com/yuya-takeyama/monotonix-playground/apps/foo/pkg/config"
	"github.com/yuya-takeyama/monotonix-playground/apps/foo/pkg/database"
	"github.com/yuya-takeyama/monotonix-playground/apps/foo/pkg/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New("WORKER")
	db := database.New(cfg.DBHost, log)
	userService := service.NewUserService(db, log)

	log.Info("Starting background worker")

	if err := db.Connect(); err != nil {
		log.Error(fmt.Sprintf("Failed to connect to database: %v", err))
		return
	}

	// Simulate background processing
	for {
		log.Info("Processing background tasks...")

		// Example: Process some users
		userIDs := []string{"user1", "user2", "user3"}
		for _, userID := range userIDs {
			user, err := userService.GetUser(userID)
			if err != nil {
				log.Error(fmt.Sprintf("Failed to process user %s: %v", userID, err))
				continue
			}
			log.Info(fmt.Sprintf("Processed: %s", user))
		}

		log.Info("Background task completed. Sleeping for 30 seconds...")
		time.Sleep(30 * time.Second)
	}
}
