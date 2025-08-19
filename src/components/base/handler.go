package base

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
)

// BaseHandler provides common functionality for all component handlers
type BaseHandler struct {
	Log *slog.Logger
}

// NewBaseHandler creates a new base handler
func NewBaseHandler(log *slog.Logger) *BaseHandler {
	return &BaseHandler{
		Log: log,
	}
}

// ServeHTTP provides common HTTP method routing
func (h *BaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, getHandler func(w http.ResponseWriter, r *http.Request), postHandler func(w http.ResponseWriter, r *http.Request)) {
	switch r.Method {
	case http.MethodGet:
		getHandler(w, r)
	case http.MethodPost:
		if postHandler != nil {
			postHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// RenderTemplate renders a templ component with error handling
func (h *BaseHandler) RenderTemplate(ctx context.Context, w http.ResponseWriter, component templ.Component, componentName string) {
	if err := component.Render(ctx, w); err != nil {
		h.Log.Error("failed to render "+componentName+" component", slog.Any("error", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// IsHTMXRequest checks if the request is an HTMX request
func (h *BaseHandler) IsHTMXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}
