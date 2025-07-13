package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yuya-takeyama/monotonix-playground/apps/foo/pkg/common"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

func main() {
	http.HandleFunc("/", helloHandler)
	startupMsg := common.GetTimestampedMessage("HELLO-WORLD", "Starting server on :8080")
	versionMsg := common.GetTimestampedMessage("HELLO-WORLD", fmt.Sprintf("Using common library version: %s", common.GetVersion()))
	infoMsg := common.FormatMessage("info", "Hello World service ready to serve requests!")
	log.Println(startupMsg)
	log.Println(versionMsg)
	log.Println(infoMsg)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
