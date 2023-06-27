package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

const (
	HTTPReadTimeout  = 30 * time.Second
	HTTPWriteTimeout = 30 * time.Second
)

func main() {
	targetDir := os.Getenv("TARGET_DIR")
	if targetDir == "" {
		log.Fatal("TARGET_DIR is required")
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(targetDir)))

	server := http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  HTTPReadTimeout,
		WriteTimeout: HTTPWriteTimeout,
	}

	log.Fatal(server.ListenAndServe())
}
