package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yuya-takeyama/monotonix-playground/apps/pkg/common"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

func main() {
	http.HandleFunc("/", helloHandler)
	startupMsg := common.GetTimestampedMessage("HELLO-WORLD", "Starting server on :8080")
	versionMsg := common.GetTimestampedMessage("HELLO-WORLD", fmt.Sprintf("Using common library version: %s", common.GetVersion()))
	log.Println(startupMsg)
	log.Println(versionMsg)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
