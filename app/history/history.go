package history

import (
	"time"

	"github.com/google/uuid"
	"github.com/ptamarov/go-cards/app/deck"
)

type UserAction struct {
	UserID   uuid.UUID
	DeckID   uuid.UUID
	CardID   uuid.UUID
	Guess    string
	Duration float64
	Date     time.Time // include
}

type DeckWithHistory struct {
	History map[uuid.UUID][]UserAction // History maps a card ID to the user actions for that card
	deck.MemoryCardDeck
}
