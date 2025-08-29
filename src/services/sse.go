package services

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/a-h/templ"
	"github.com/google/uuid"
)

// SSEClient represents a connected client
type SSEClient struct {
	ID      string
	Writer  http.ResponseWriter
	Flusher http.Flusher
	Done    chan bool
	Context context.Context
}

// SSEService manages Server-Sent Events connections and broadcasting
type SSEService struct {
	log     *slog.Logger
	clients map[string]*SSEClient
	mutex   sync.RWMutex
}

// NewSSEService creates a new SSE service
func NewSSEService(log *slog.Logger) *SSEService {
	return &SSEService{
		log:     log,
		clients: make(map[string]*SSEClient),
	}
}

// AddClient adds a new SSE client connection
func (s *SSEService) AddClient(clientID string, w http.ResponseWriter, ctx context.Context) *SSEClient {
	flusher, ok := w.(http.Flusher)
	if !ok {
		s.log.Error("Response writer does not support flushing")
		return nil
	}

	client := &SSEClient{
		ID:      clientID,
		Writer:  w,
		Flusher: flusher,
		Done:    make(chan bool),
		Context: ctx,
	}

	s.mutex.Lock()
	s.clients[clientID] = client
	s.mutex.Unlock()

	s.log.Info("SSE client connected", "clientID", clientID, "total", len(s.clients))
	return client
}

// RemoveClient removes an SSE client connection
func (s *SSEService) RemoveClient(clientID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if client, exists := s.clients[clientID]; exists {
		close(client.Done)
		delete(s.clients, clientID)
		s.log.Info("SSE client disconnected", "clientID", clientID, "total", len(s.clients))
	}
}

// BroadcastOOBUpdate sends an OOB update to all connected clients
func (s *SSEService) BroadcastOOBUpdate(component templ.Component) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if len(s.clients) == 0 {
		return
	}

	// Render component to string
	var buf strings.Builder
	err := component.Render(context.Background(), &buf)
	if err != nil {
		s.log.Error("Failed to render component for SSE broadcast", "error", err)
		return
	}

	html := buf.String()
	data := fmt.Sprintf("event: oob-update\ndata: %s\n\n", html)

	// Broadcast to all clients
	for clientID, client := range s.clients {
		select {
		case <-client.Done:
			// Client is done, skip
			continue
		default:
			_, err := client.Writer.Write([]byte(data))
			if err != nil {
				s.log.Error("Failed to write to SSE client", "clientID", clientID, "error", err)
				// Remove client asynchronously to avoid blocking
				go s.RemoveClient(clientID)
				continue
			}
			client.Flusher.Flush()
		}
	}

	s.log.Debug("Broadcasted OOB update to clients", "count", len(s.clients))
}

// ServeSSE handles SSE endpoint requests
func (s *SSEService) ServeSSE(w http.ResponseWriter, r *http.Request) {
	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Generate client ID (in real app, you might use session ID or user ID)
	clientID := uuid.New().String()

	// Add client
	client := s.AddClient(clientID, w, r.Context())
	if client == nil {
		http.Error(w, "Failed to create SSE connection", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "data: {\"type\":\"connected\",\"clientID\":\"%s\"}\n\n", clientID)
	client.Flusher.Flush()

	// Wait for client disconnect or context cancellation
	select {
	case <-r.Context().Done():
		s.RemoveClient(clientID)
	case <-client.Done:
		// Client was removed elsewhere
	}
}
