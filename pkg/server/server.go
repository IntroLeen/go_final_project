package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"go_final_project/pkg/api"
)

const DefaultPort = 7540
const DefaultWebDir = "./web"

func Run() error {
	port := resolvePort()
	webDir := DefaultWebDir

	if _, err := os.Stat(webDir); err != nil {
		abs, _ := filepath.Abs(webDir)
		return fmt.Errorf("web directory not found: %s (abs: %s): %w", webDir, abs, err)
	}

	mux := http.NewServeMux()

	api.Init(mux)

	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("todo-server: serving %s on http://%s/\n", webDir, addr)
	return srv.ListenAndServe()
}

func resolvePort() int {
	if v := os.Getenv("TODO_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 && p < 65536 {
			return p
		}
		log.Printf("todo-server: invalid TODO_PORT=%q, fallback to %d\n", v, DefaultPort)
	}
	return DefaultPort
}
