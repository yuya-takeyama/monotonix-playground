package main

import (
	"fmt"
	"net/http"

	"github.com/yuya-takeyama/monotonix-playground/apps/web-app/pkg/common"
)

func main() {
	fmt.Println(common.FormatMessage("info", "Starting API server"))

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Path[len("/users/"):]
		if userID == "" {
			http.Error(w, "User ID required", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"user_id": "%s", "message": "%s"}`, userID, common.GetTimestampedMessage("API", "User data"))
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status": "OK", "version": "%s"}`, common.GetVersion())
	})

	fmt.Println(common.FormatMessage("info", "Server listening on :8080"))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(common.FormatMessage("error", fmt.Sprintf("Server failed: %v", err)))
	}
}
