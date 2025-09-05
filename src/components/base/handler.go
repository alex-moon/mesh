package base

import (
	"context"
	"log/slog"
	"mesh/src/services"
	"net/http"

	"github.com/a-h/templ"
)

type BaseHandler struct {
	Log          *slog.Logger
	name         string
	EventService *services.EventService
}

func NewBaseHandler(log *slog.Logger, name string, eventService *services.EventService) *BaseHandler {
	return &BaseHandler{
		Log:          log,
		name:         name,
		EventService: eventService,
	}
}

func (h *BaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, f map[string]http.HandlerFunc) {
	if handler, exists := f[r.Method]; exists && handler != nil {
		handler(w, r)
	} else {
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
