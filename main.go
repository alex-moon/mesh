package main

import (
	"log"
	"log/slog"
	"mesh/src"
	"mesh/src/components"
	"net/http"
	"os"
)

func main() {
	// Create logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create registry with all handlers
	registry := components.NewRegistry(logger)

	// Index page handler
	http.HandleFunc("/", src.IndexHandler(registry))

	// Route handlers with registry context middleware
	http.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		registry.AppHandler.ServeHTTP(w, r)
	})

	http.HandleFunc("/board", func(w http.ResponseWriter, r *http.Request) {
		registry.BoardHandler.ServeHTTP(w, r)
	})

	http.HandleFunc("/column", func(w http.ResponseWriter, r *http.Request) {
		registry.ColumnHandler.ServeHTTP(w, r)
	})

	http.HandleFunc("/card", func(w http.ResponseWriter, r *http.Request) {
		registry.CardHandler.ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(":8000", nil))
}
