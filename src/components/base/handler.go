package base

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
)

type BaseHandler struct {
	Log  *slog.Logger
	name string
}

func NewBaseHandler(log *slog.Logger, name string) *BaseHandler {
	return &BaseHandler{
		Log:  log,
		name: name,
	}
}

func (h *BaseHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
	getHandler func(
		w http.ResponseWriter,
		r *http.Request,
	),
	postHandler func(
		w http.ResponseWriter,
		r *http.Request,
	),
	patchHandler func(
		w http.ResponseWriter,
		r *http.Request,
	),
) {
	switch r.Method {
	case http.MethodGet:
		if getHandler == nil {
			getHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case http.MethodPost:
		if postHandler != nil {
			postHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case http.MethodPatch:
		if patchHandler != nil {
			patchHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *BaseHandler) RenderTemplate(
	ctx context.Context,
	w http.ResponseWriter,
	component templ.Component,
) {
	if err := component.Render(ctx, w); err != nil {
		h.Log.Error("failed to render "+h.name+" component", slog.Any("error", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// IsHTMXRequest checks if the request is an HTMX request
func (h *BaseHandler) IsHTMXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}
