package card

import (
	"log/slog"
	"mesh/src/services"
	"net/http"

	"mesh/src/components/base"

	"github.com/a-h/templ"
)

// Handler represents the card component's HTTP handler
type Handler struct {
	*base.BaseHandler
	*services.CardService
}

// New creates a new card handler with dependencies
func New(log *slog.Logger, cardService *services.CardService) *Handler {
	return &Handler{
		BaseHandler: base.NewBaseHandler(log),
		CardService: cardService,
	}
}

// ServeHTTP handles HTTP requests for the card component
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.BaseHandler.ServeHTTP(w, r, h.Get, nil)
}

// Get renders the full card component
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	card, err := h.CardService.GetCard(ctx.Value("cardID").(int))
	if err != nil {
		http.Error(w, "Card not found", http.StatusNotFound)
		return
	}

	// Create props for the card template
	props := CardProps{
		Title:   card.Title,
		Content: card.Content,
		IsHTMX:  h.IsHTMXRequest(r),
	}

	// Render using the base handler
	c := Card(props)
	h.RenderTemplate(ctx, w, c, "card")
}

func (h *Handler) RenderComponent(card services.Card) templ.Component {
	props := CardProps{
		Title:   card.Title,
		Content: card.Content,
		IsHTMX:  false,
	}

	return Card(props)
}
