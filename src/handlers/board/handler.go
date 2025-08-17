// handlers/board/handler.go
package board

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	boardComponent "mesh/src/components/board" // Import the generated package

	"github.com/a-h/templ"
)

// CounterService interface for dependency injection
type CounterService interface {
	GetCount(ctx context.Context) int
	Increment(ctx context.Context) int
}

// Handler represents the board component's HTTP handler
type Handler struct {
	Log            *slog.Logger
	CounterService CounterService
}

// New creates a new board handler with dependencies
func New(log *slog.Logger, counterService CounterService) *Handler {
	return &Handler{
		Log:            log,
		CounterService: counterService,
	}
}

// ServeHTTP handles HTTP requests for the board component
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	case http.MethodPost:
		h.Post(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Get renders the board component
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get current count from service
	count := h.CounterService.GetCount(ctx)

	// Create props for the board template
	props := boardComponent.BoardProps{
		Count:  count,
		IsHTMX: r.Header.Get("HX-Request") == "true",
	}

	// Render using the generated templ function
	component := boardComponent.Board(props)
	if err := component.Render(ctx, w); err != nil {
		h.Log.Error("failed to render board component", slog.Any("error", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Post handles form submissions and actions
func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	action := r.FormValue("action")

	switch action {
	case "increment":
		h.handleIncrement(ctx, w, r)
	default:
		http.Error(w, "Unknown action", http.StatusBadRequest)
	}
}

// handleIncrement increments the counter
func (h *Handler) handleIncrement(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	newCount := h.CounterService.Increment(ctx)
	h.Log.Info("counter incremented", slog.Int("newCount", newCount))

	// Re-render the component with updated count
	h.Get(w, r.WithContext(ctx))
}

// RenderComponent provides a way for other components to render this one
func (h *Handler) RenderComponent(ctx context.Context) templ.Component {
	count := h.CounterService.GetCount(ctx)

	fmt.Println("count", count)

	props := boardComponent.BoardProps{
		Count:  count,
		IsHTMX: false,
	}

	return boardComponent.Board(props)
}

// Context helpers for dependency injection
type contextKey string

const handlerKey contextKey = "board.handler"

func WithHandler(ctx context.Context, h *Handler) context.Context {
	return context.WithValue(ctx, handlerKey, h)
}

func FromContext(ctx context.Context) *Handler {
	if h, ok := ctx.Value(handlerKey).(*Handler); ok {
		return h
	}
	return nil
}
