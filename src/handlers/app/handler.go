// handlers/app/handler.go
package app

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	appComponent "mesh/src/components/app" // Import the generated package
)

// BoardHandler interface for dependency injection
type BoardHandler interface {
	RenderComponent(ctx context.Context) templ.Component
}

// Handler represents the app component's HTTP handler
type Handler struct {
	Log          *slog.Logger
	BoardHandler BoardHandler
}

// New creates a new app handler with dependencies
func New(log *slog.Logger, boardHandler BoardHandler) *Handler {
	return &Handler{
		Log:          log,
		BoardHandler: boardHandler,
	}
}

// ServeHTTP handles HTTP requests for the app component
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Get renders the full app component
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get the board component
	boardComponent := h.BoardHandler.RenderComponent(ctx)

	// Create props for the app template
	props := appComponent.AppProps{
		BoardComponent: boardComponent,
		IsHTMX:         r.Header.Get("HX-Request") == "true",
	}

	// Render using the generated templ function
	component := appComponent.App(props)
	if err := component.Render(ctx, w); err != nil {
		h.Log.Error("failed to render app component", slog.Any("error", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// RenderComponent provides a way for other components to render this one
func (h *Handler) RenderComponent(ctx context.Context) templ.Component {
	boardComponent := h.BoardHandler.RenderComponent(ctx)

	props := appComponent.AppProps{
		BoardComponent: boardComponent,
		IsHTMX:         false,
	}

	return appComponent.App(props)
}

// Context helpers for dependency injection
type contextKey string

const handlerKey contextKey = "app.handler"

func WithHandler(ctx context.Context, h *Handler) context.Context {
	return context.WithValue(ctx, handlerKey, h)
}

func FromContext(ctx context.Context) *Handler {
	if h, ok := ctx.Value(handlerKey).(*Handler); ok {
		return h
	}
	return nil
}
