package app

import (
	"log/slog"
	"mesh/src/components/board"
	"mesh/src/services"
	"net/http"

	"mesh/src/components/base"

	"github.com/a-h/templ"
)

type Handler struct {
	*base.BaseHandler
	BoardHandler *board.Handler
}

func New(
	log *slog.Logger,
	eventService *services.EventService,
	boardHandler *board.Handler,
) *Handler {
	return &Handler{
		BaseHandler:  base.NewBaseHandler(log, "app", eventService),
		BoardHandler: boardHandler,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.BaseHandler.ServeHTTP(w, r, map[string]http.HandlerFunc{
		http.MethodGet: h.Get,
	})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	h.RenderTemplate(r.Context(), w, h.RenderComponent())
}

func (h *Handler) RenderComponent() templ.Component {
	boardComponent := h.BoardHandler.RenderComponent()
	props := AppProps{
		BoardComponent: boardComponent,
	}
	return App(props)
}
