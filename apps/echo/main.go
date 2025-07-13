package main

import (
	"fmt"
	"log"
	"net/http"
)

func echoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Method: %s\n", r.Method)
	fmt.Fprintf(w, "URL: %s\n", r.URL.String())
	fmt.Fprintf(w, "Headers:\n")

	for name, values := range r.Header {
		for _, value := range values {
			fmt.Fprintf(w, "%s: %s\n", name, value)
		}
	}

	fmt.Fprintf(w, "\nBody:\n")
	if _, err := r.Body.Read([]byte{}); err == nil {
		fmt.Fprintf(w, "Body is empty or error reading body.")
	}
}

func main() {
	http.HandleFunc("/", echoHandler)
	log.Println("Starting echo server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
