package app

import (
	"context"
	"log/slog"
	"mesh/src/components/board"
	"net/http"

	"mesh/src/components/base"

	"github.com/a-h/templ"
)

// Handler represents the app component's HTTP handler
type Handler struct {
	*base.BaseHandler
	BoardHandler *board.Handler
}

// New creates a new app handler with dependencies
func New(log *slog.Logger, boardHandler *board.Handler) *Handler {
	return &Handler{
		BaseHandler:  base.NewBaseHandler(log),
		BoardHandler: boardHandler,
	}
}

// ServeHTTP handles HTTP requests for the app component
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.BaseHandler.ServeHTTP(w, r, h.Get, nil)
}

// Get renders the full app component
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get the board component
	boardComponent := h.BoardHandler.RenderComponent(ctx)

	// Create props for the app template
	props := AppProps{
		BoardComponent: boardComponent,
		IsHTMX:         h.IsHTMXRequest(r),
	}

	// Render using the base handler
	c := App(props)
	h.RenderTemplate(ctx, w, c, "app")
}

// RenderComponent provides a way for other components to render this one
func (h *Handler) RenderComponent(ctx context.Context) templ.Component {
	boardComponent := h.BoardHandler.RenderComponent(ctx)

	props := AppProps{
		BoardComponent: boardComponent,
		IsHTMX:         false,
	}

	return App(props)
}
