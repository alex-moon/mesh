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

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Index page handler
	http.HandleFunc("/", src.IndexHandler(registry))

	// Route handlers with registry context middleware
	http.Handle("/app", registry.AppHandler)
	http.Handle("/board", registry.BoardHandler)
	http.Handle("/column", registry.ColumnHandler)
	http.Handle("/card", registry.CardHandler)

	// SSE endpoint for real-time updates
	http.HandleFunc("/sse", registry.SSEService.ServeSSE)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
