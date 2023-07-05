package httprangelimit

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
)

type Config struct {
	MaxLengthBytes int64
}

func handleRangeLimit(r *http.Request, cfg Config, fileSystem http.FileSystem) error {
	if r.Method != "GET" && r.Method != "HEAD" {
		return nil
	}

	rangeHeader := r.Header.Get("Range")
	if rangeHeader == "" ||
		strings.Contains(rangeHeader, ",") ||
		!strings.HasPrefix(rangeHeader, "bytes=") ||
		!strings.HasSuffix(rangeHeader, "-") {
		return nil
	}

	size, err := fileSize(fileSystem, requestFilepath(r))
	if err != nil {
		return fmt.Errorf("fileSize: %w", err)
	}

	var startBytes int64

	_, err = fmt.Fscanf(strings.NewReader(rangeHeader), "bytes=%d-", &startBytes)
	if err != nil {
		return fmt.Errorf("Fscanf: %w", err)
	}

	endBytes := startBytes + cfg.MaxLengthBytes
	if endBytes > size {
		endBytes = size
	}

	r.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", startBytes, endBytes))

	return nil
}

func requestFilepath(r *http.Request) string {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}

	return path.Clean(upath)
}

func fileSize(fs http.FileSystem, name string) (int64, error) {
	file, err := fs.Open(name)
	if err != nil {
		return 0, fmt.Errorf("Open: %w", err)
	}
	defer file.Close()

	d, err := file.Stat()
	if err != nil {
		return 0, fmt.Errorf("Stat: %w", err)
	}

	return d.Size(), nil
}

func HTTPRangeLimit(cfg Config, fs http.FileSystem, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := handleRangeLimit(r, cfg, fs); err != nil {
			log.Println("handleRangeLimit", err)
		}

		h.ServeHTTP(w, r)
	})
}
