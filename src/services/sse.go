package services

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/r3labs/sse/v2"
)

type BatchedUpdate struct {
	ID   string `json:"id"`
	HTML string `json:"html"`
}

type UpdateBatch struct {
	BatchID string          `json:"batchId"`
	Updates []BatchedUpdate `json:"updates"`
}

type SSEService struct {
	log            *slog.Logger
	server         *sse.Server
	pendingUpdates []BatchedUpdate
	batchMutex     sync.Mutex
	batchTimer     *time.Timer
	batchDuration  time.Duration
}

func NewSSEService(log *slog.Logger) *SSEService {
	server := sse.New()

	server.AutoReplay = false
	server.AutoStream = true

	server.OnSubscribe = func(streamID string, sub *sse.Subscriber) {
		log.Info("SSE client connected", "streamID", streamID)
	}

	server.OnUnsubscribe = func(streamID string, sub *sse.Subscriber) {
		log.Info("SSE client disconnected", "streamID", streamID)
	}

	return &SSEService{
		log:           log,
		server:        server,
		batchDuration: 50 * time.Millisecond,
	}
}

func (s *SSEService) BroadcastOOBUpdate(component templ.Component) {
	var buf strings.Builder
	err := component.Render(context.Background(), &buf)
	if err != nil {
		s.log.Error("Failed to render component for SSE broadcast", "error", err)
		return
	}

	html := buf.String()

	var componentID = "?"
	if idStart := strings.Index(html, "id=\""); idStart >= 0 {
		idStart += 4 // Skip `id="`
		if idEnd := strings.Index(html[idStart:], "\""); idEnd >= 0 {
			componentID = html[idStart : idStart+idEnd]
		}
	}
	if componentID == "?" {
		s.log.Error("Failed to find component ID for OOB update")
		return
	}

	s.log.Info("Queueing OOB update for batch", "componentID", componentID)

	update := BatchedUpdate{
		ID:   componentID,
		HTML: html,
	}

	s.addToBatch(update)
}

func (s *SSEService) addToBatch(update BatchedUpdate) {
	s.batchMutex.Lock()
	defer s.batchMutex.Unlock()

	found := false
	for i, existing := range s.pendingUpdates {
		if existing.ID == update.ID {
			s.pendingUpdates[i] = update
			found = true
			break
		}
	}

	if !found {
		s.pendingUpdates = append(s.pendingUpdates, update)
	}

	if s.batchTimer != nil {
		s.batchTimer.Stop()
	}

	s.batchTimer = time.AfterFunc(s.batchDuration, func() {
		s.flushBatch()
	})
}

func (s *SSEService) flushBatch() {
	s.batchMutex.Lock()
	if len(s.pendingUpdates) == 0 {
		s.batchMutex.Unlock()
		return
	}

	updates := make([]BatchedUpdate, len(s.pendingUpdates))
	copy(updates, s.pendingUpdates)
	s.pendingUpdates = s.pendingUpdates[:0] // Clear the slice
	s.batchMutex.Unlock()

	batch := UpdateBatch{
		BatchID: uuid.New().String(),
		Updates: updates,
	}

	batchData, err := json.Marshal(batch)
	if err != nil {
		s.log.Error("Failed to serialize batch", "error", err)
		return
	}

	s.log.Info("Broadcasting batch", "batchID", batch.BatchID, "updateCount", len(updates))

	s.server.Publish("oob-updates", &sse.Event{
		Event: []byte("oob-batch"),
		Data:  batchData,
	})

	s.log.Debug("Broadcasted batch to all clients", "batchID", batch.BatchID, "count", len(updates))
}

func (s *SSEService) ServeSSE(w http.ResponseWriter, r *http.Request) {
	s.server.ServeHTTP(w, r)
}
