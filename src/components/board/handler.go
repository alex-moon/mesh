package board

import (
	"log/slog"
	"mesh/src/components/base"
	"mesh/src/components/column"
	"mesh/src/services"
	"net/http"

	"github.com/a-h/templ"
)

type Handler struct {
	*base.BaseHandler
	CardService   *services.CardService
	ColumnHandler *column.Handler
}

func New(log *slog.Logger, eventService *services.EventService, cardService *services.CardService, columnHandler *column.Handler) *Handler {
	return &Handler{
		BaseHandler:   base.NewBaseHandler(log, "board", eventService),
		CardService:   cardService,
		ColumnHandler: columnHandler,
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
	columnsWithCards := h.CardService.GetColumns()
	var columnComponents []templ.Component
	for _, columnWithCards := range columnsWithCards {
		columnComponent := h.ColumnHandler.RenderComponent(&columnWithCards, false)
		columnComponents = append(columnComponents, columnComponent)
	}
	props := BoardProps{
		Columns: columnComponents,
	}
	return Board(props)
}
