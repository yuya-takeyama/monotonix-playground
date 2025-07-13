package main

import (
	"fmt"
	"net/http"

	"github.com/yuya-takeyama/monotonix-playground/apps/foo/internal/service"
	"github.com/yuya-takeyama/monotonix-playground/apps/foo/pkg/config"
	"github.com/yuya-takeyama/monotonix-playground/apps/foo/pkg/database"
	"github.com/yuya-takeyama/monotonix-playground/apps/foo/pkg/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New("API-SERVER")
	db := database.New(cfg.DBHost, log)
	userService := service.NewUserService(db, log)

	log.Info(fmt.Sprintf("Starting API server on port %s", cfg.Port))

	if err := db.Connect(); err != nil {
		log.Error(fmt.Sprintf("Failed to connect to database: %v", err))
		return
	}

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Path[len("/users/"):]
		if userID == "" {
			http.Error(w, "User ID required", http.StatusBadRequest)
			return
		}

		user, err := userService.GetUser(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"user": "%s"}`, user)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	log.Info(fmt.Sprintf("Server listening on :%s", cfg.Port))
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Error(fmt.Sprintf("Server failed: %v", err))
	}
}