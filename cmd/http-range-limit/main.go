package main

import (
	"log"
	"net/http"
	"os"
	"time"

	httprangelimit "github.com/nakatanakatana/http-range-limit"
)

const (
	HTTPReadTimeout  = 5 * time.Second
	HTTPWriteTimeout = 5 * time.Second

	RangeMaxLengthBytes = 1 * 1024 * 1024
)

func main() {
	targetDir := os.Getenv("TARGET_DIR")
	if targetDir == "" {
		log.Fatal("TARGET_DIR is required")
	}

	cfg := httprangelimit.Config{
		MaxLengthBytes: RangeMaxLengthBytes,
	}

	mux := http.NewServeMux()
	fs := http.Dir(targetDir)
	mux.Handle("/",
		httprangelimit.HTTPRangeLimit(
			cfg,
			fs,
			http.FileServer(fs),
		),
	)

	server := http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  HTTPReadTimeout,
		WriteTimeout: HTTPWriteTimeout,
	}

	log.Fatal(server.ListenAndServe())
}
