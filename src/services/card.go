package services

import (
	"fmt"
	"log/slog"
	"sync"
)

type Card struct {
	ID       int
	Title    string
	Content  string
	ColumnID int
}

type Column struct {
	ID    int
	Title string
	Order int // For maintaining column display order
}

type CardService struct {
	mu      sync.RWMutex
	cards   map[int]*Card   // cardID -> Card
	columns map[int]*Column // columnID -> Column

	// For maintaining order within columns
	columnCards map[int][]int // columnID -> []cardID (ordered)

	nextCardID   int
	nextColumnID int

	log *slog.Logger
}

func NewCardService(log *slog.Logger) *CardService {
	service := &CardService{
		mu:           sync.RWMutex{},
		cards:        make(map[int]*Card),
		columns:      make(map[int]*Column),
		columnCards:  make(map[int][]int),
		nextCardID:   4, // Starting after our seed data
		nextColumnID: 4,
		log:          log,
	}

	// Seed data
	service.seedData()
	return service
}

func (c *CardService) seedData() {
	// Create columns
	c.columns[1] = &Column{ID: 1, Title: "To Do", Order: 0}
	c.columns[2] = &Column{ID: 2, Title: "In Progress", Order: 1}
	c.columns[3] = &Column{ID: 3, Title: "Done", Order: 2}

	// Create cards
	c.cards[1] = &Card{ID: 1, Title: "Blog post", Content: "Once the app is working and looking good, write it up", ColumnID: 1}
	c.cards[2] = &Card{ID: 2, Title: "Post to HN", Content: "", ColumnID: 1}
	c.cards[3] = &Card{ID: 3, Title: "Build app", Content: "Implement minimal Kanban Board with columns and draggable/editable cards", ColumnID: 2}

	// Set up column ordering
	c.columnCards[1] = []int{1, 2}
	c.columnCards[2] = []int{3}
	c.columnCards[3] = []int{}
}

func (c *CardService) CanPromote(cardID int) bool {
	for _, column := range c.columns {
		for _, id := range c.columnCards[column.ID] {
			if id == cardID {
				return column.Order < len(c.columns)
			}
		}
	}
	return false
}

func (c *CardService) CanDemote(cardID int) bool {
	for _, column := range c.columns {
		for _, id := range c.columnCards[column.ID] {
			if id == cardID {
				return column.Order > 0
			}
		}
	}
	return false
}

func (c *CardService) GetColumn(id int) (*ColumnWithCards, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	column, exists := c.columns[id]
	if !exists {
		return nil, fmt.Errorf("column with id %d not found", id)
	}

	cards := c.getCardsForColumn(id)
	return &ColumnWithCards{
		Column: *column,
		Cards:  cards,
	}, nil
}

// GetColumns returns columns with their cards in order
func (c *CardService) GetColumns() []ColumnWithCards {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var result []ColumnWithCards

	// Sort columns by order
	for i := 0; i < len(c.columns); i++ {
		for _, column := range c.columns {
			if column.Order == i {
				cards := c.getCardsForColumn(column.ID)
				result = append(result, ColumnWithCards{
					Column: *column,
					Cards:  cards,
				})
				break
			}
		}
	}

	return result
}

type ColumnWithCards struct {
	Column Column
	Cards  []Card
}

func (c *CardService) getCardsForColumn(columnID int) []Card {
	cardIDs := c.columnCards[columnID]
	cards := make([]Card, 0, len(cardIDs))

	for _, cardID := range cardIDs {
		if card, exists := c.cards[cardID]; exists {
			cards = append(cards, *card)
		}
	}

	return cards
}

// GetCard gets a single card by ID
func (c *CardService) GetCard(cardID int) (*Card, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if card, exists := c.cards[cardID]; exists {
		return card, nil
	}
	return nil, fmt.Errorf("card with ID %d not found", cardID)
}

// AddCard adds a new card to a column
func (c *CardService) AddCard(title, content string, columnID int) (*Card, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if column exists
	if _, exists := c.columns[columnID]; !exists {
		return nil, fmt.Errorf("column with ID %d not found", columnID)
	}

	card := &Card{
		ID:       c.nextCardID,
		Title:    title,
		Content:  content,
		ColumnID: columnID,
	}

	c.cards[c.nextCardID] = card
	c.columnCards[columnID] = append(c.columnCards[columnID], c.nextCardID)
	c.nextCardID++

	return card, nil
}

// UpdateCard updates an existing card's content
func (c *CardService) UpdateCard(cardID int, title, content string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	card, exists := c.cards[cardID]
	if !exists {
		return fmt.Errorf("card with ID %d not found", cardID)
	}

	card.Title = title
	card.Content = content
	return nil
}

// MoveCard moves a card to a different column and/or position
func (c *CardService) MoveCard(cardID, newColumnID, newPosition int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	card, exists := c.cards[cardID]
	if !exists {
		return fmt.Errorf("card with ID %d not found", cardID)
	}

	if _, exists := c.columns[newColumnID]; !exists {
		return fmt.Errorf("column with ID %d not found", newColumnID)
	}

	oldColumnID := card.ColumnID

	// Remove from old column
	c.removeCardFromColumn(cardID, oldColumnID)

	// Add to new column at specified position
	c.insertCardInColumn(cardID, newColumnID, newPosition)

	// Update card's column reference
	card.ColumnID = newColumnID

	return nil
}

// ReorderCardInColumn moves a card to a different position within the same column
func (c *CardService) ReorderCardInColumn(cardID, newPosition int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	card, exists := c.cards[cardID]
	if !exists {
		return fmt.Errorf("card with ID %d not found", cardID)
	}

	columnID := card.ColumnID
	cardList := c.columnCards[columnID]

	// Find current position
	currentPos := -1
	for i, id := range cardList {
		if id == cardID {
			currentPos = i
			break
		}
	}

	if currentPos == -1 {
		return fmt.Errorf("card %d not found in column %d", cardID, columnID)
	}

	// If position hasn't changed, do nothing
	if currentPos == newPosition {
		return nil
	}

	// Remove from current position
	cardList = append(cardList[:currentPos], cardList[currentPos+1:]...)

	// Insert at new position
	if newPosition >= len(cardList) {
		cardList = append(cardList, cardID)
	} else {
		cardList = append(cardList[:newPosition], append([]int{cardID}, cardList[newPosition:]...)...)
	}

	c.columnCards[columnID] = cardList
	return nil
}

// DeleteCard removes a card completely
func (c *CardService) DeleteCard(cardID int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	card, exists := c.cards[cardID]
	if !exists {
		return fmt.Errorf("card with ID %d not found", cardID)
	}

	// Remove from column
	c.removeCardFromColumn(cardID, card.ColumnID)

	// Remove from cards map
	delete(c.cards, cardID)

	return nil
}

// Helper: remove card from column's card list
func (c *CardService) removeCardFromColumn(cardID, columnID int) {
	cardList := c.columnCards[columnID]
	for i, id := range cardList {
		if id == cardID {
			c.columnCards[columnID] = append(cardList[:i], cardList[i+1:]...)
			break
		}
	}
}

// Helper: insert card into column at specific position
func (c *CardService) insertCardInColumn(cardID, columnID, position int) {
	cardList := c.columnCards[columnID]

	if position >= len(cardList) {
		// Append to end
		c.columnCards[columnID] = append(cardList, cardID)
	} else {
		// Insert at position
		cardList = append(cardList[:position], append([]int{cardID}, cardList[position:]...)...)
		c.columnCards[columnID] = cardList
	}
}
