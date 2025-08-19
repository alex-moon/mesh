package column

import (
	"log/slog"
	"mesh/src/components/card"
	"mesh/src/services"
	"net/http"

	"mesh/src/components/base"

	"github.com/a-h/templ"
)

// Handler represents the column component's HTTP handler
type Handler struct {
	*base.BaseHandler
	CardHandler *card.Handler
	*services.CardService
}

// New creates a new column handler with dependencies
func New(log *slog.Logger, cardService *services.CardService, cardHandler *card.Handler) *Handler {
	return &Handler{
		BaseHandler: base.NewBaseHandler(log),
		CardHandler: cardHandler,
		CardService: cardService,
	}
}

// ServeHTTP handles HTTP requests for the column component
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.BaseHandler.ServeHTTP(w, r, h.Get, nil)
}

// Get renders the full column component
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Convert to column components
	var columnWithCards = h.CardService.GetColumn(ctx.Value("columnID").(int))
	var cardComponents []templ.Component
	for _, card := range columnWithCards.Cards {
		cardComponent := h.CardHandler.RenderComponent(card)
		cardComponents = append(cardComponents, cardComponent)
	}
	props := ColumnProps{
		Title:  columnWithCards.Column.Title,
		Cards:  cardComponents,
		IsHTMX: h.IsHTMXRequest(r),
	}

	// Render using the base handler
	c := Column(props)
	h.RenderTemplate(ctx, w, c, "column")
}

func (h *Handler) RenderComponent(column *services.ColumnWithCards) templ.Component {
	// Convert to card components
	var cardComponents []templ.Component
	for _, card := range column.Cards {
		columnComponent := h.CardHandler.RenderComponent(card)
		cardComponents = append(cardComponents, columnComponent)
	}
	props := ColumnProps{
		Title:  column.Column.Title,
		Cards:  cardComponents,
		IsHTMX: false,
	}

	return Column(props)
}
