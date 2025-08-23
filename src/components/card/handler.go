package card

import (
	"fmt"
	"log/slog"
	"mesh/src/services"
	"net/http"
	"strconv"
	"strings"

	"mesh/src/components/base"

	"github.com/a-h/templ"
)

type Handler struct {
	*base.BaseHandler
	*services.CardService
}

func New(log *slog.Logger, cardService *services.CardService) *Handler {
	return &Handler{
		BaseHandler: base.NewBaseHandler(log, "card"),
		CardService: cardService,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.BaseHandler.ServeHTTP(w, r, h.Get, h.Post, h.Patch)
}

func (h *Handler) getCardFromRequest(r *http.Request) (*services.Card, error) {
	cardIDString := r.FormValue("cardID")
	if cardIDString == "" {
		return nil, fmt.Errorf("missing card ID")
	}

	cardID, err := strconv.Atoi(cardIDString)
	if err != nil {
		return nil, fmt.Errorf("invalid card ID %s", cardIDString)
	}

	card, err := h.CardService.GetCard(cardID)
	if err != nil {
		return nil, fmt.Errorf("card not found %d", cardID)
	}

	return card, nil
}

func (h *Handler) getColumnFromRequest(r *http.Request) (*services.Column, error) {
	columnIDString := r.FormValue("columnID")
	if columnIDString == "" {
		return nil, fmt.Errorf("missing column ID")
	}

	columnID, err := strconv.Atoi(columnIDString)
	if err != nil {
		return nil, fmt.Errorf("invalid column ID %s", columnIDString)
	}

	column, err := h.CardService.GetColumn(columnID)
	if err != nil {
		return nil, fmt.Errorf("column not found %d", columnID)
	}

	return &column.Column, nil
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	card, err := h.getCardFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	props := h.getProps(card)
	props.IsHTMX = h.IsHTMXRequest(r)
	c := Card(props)
	h.RenderTemplate(r.Context(), w, c)
}

func (h *Handler) validate(r *http.Request) (Data, Errors) {
	errors := Errors{}

	var data = Data{
		Title:   strings.TrimSpace(r.FormValue("title")),
		Content: strings.TrimSpace(r.FormValue("content")),
	}

	if data.Title == "" {
		errors.Title = "Title is required"
	}

	if len(data.Title) > 100 {
		errors.Title = "Title must be less than 100 characters"
	}

	if len(data.Content) > 1000 {
		errors.Content = "Content must be less than 1000 characters"
	}

	if r.FormValue("columnID") != "" {
		var column, err = h.getColumnFromRequest(r)
		if err != nil {
			errors.ColumnID = err.Error()
		} else {
			data.ColumnID = column.ID
		}
	}

	return data, errors
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	var data, errors = h.validate(r)
	if errors.Any() {
		var props = h.getPropsWithData(&services.Card{}, data, errors)
		props.IsHTMX = h.IsHTMXRequest(r)
		h.RenderTemplate(r.Context(), w, Card(props))
		return
	}

	card, err := h.CardService.AddCard(
		data.Title,
		data.Content,
		data.ColumnID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var props = h.getProps(card)
	props.IsHTMX = h.IsHTMXRequest(r)
	h.RenderTemplate(r.Context(), w, Card(props))
}

func (h *Handler) Patch(w http.ResponseWriter, r *http.Request) {
	card, err := h.getCardFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.Log.Info("Validating card")
	var data, errors = h.validate(r)

	if errors.Any() {
		var props = h.getPropsWithData(card, data, errors)
		props.IsHTMX = h.IsHTMXRequest(r)
		h.RenderTemplate(r.Context(), w, Card(props))
		return
	}

	h.Log.Info("Updating card")
	err = h.CardService.UpdateCard(
		card.ID,
		data.Title,
		data.Content,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Log.Info("Rendering card")
	var props = h.getProps(card)
	props.IsHTMX = h.IsHTMXRequest(r)
	h.RenderTemplate(r.Context(), w, Card(props))
}

func (h *Handler) RenderComponent(card *services.Card) templ.Component {
	props := h.getProps(card)
	return Card(props)
}

func (h *Handler) getProps(card *services.Card) CardProps {
	return h.getPropsWithData(card, Data{
		ID:      card.ID,
		Title:   card.Title,
		Content: card.Content,
	}, Errors{})
}

func (h *Handler) getPropsWithData(card *services.Card, data Data, errors Errors) CardProps {
	return CardProps{
		Card:       card,
		Data:       data,
		Errors:     errors,
		HasErrors:  errors.Any(),
		CanDemote:  h.CardService.CanDemote(card.ID),
		CanPromote: h.CardService.CanPromote(card.ID),
	}
}
