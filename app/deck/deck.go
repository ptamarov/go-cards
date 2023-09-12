package deck

import (
	"github.com/google/uuid"
	"github.com/ptamarov/go-cards/app/card"
)

type MemoryCardDeck struct {
	UserID   uuid.UUID
	DeckID   uuid.UUID
	Cards    map[uuid.UUID]card.MemoryCard
	DeckName string
}

func (mcd *MemoryCardDeck) PutBack(card card.MemoryCard, cardID uuid.UUID) {
	mcd.Cards[cardID] = card
}

func (mcd *MemoryCardDeck) Remove(cardID uuid.UUID) card.MemoryCard {
	out := mcd.Cards[cardID]
	delete(mcd.Cards, cardID)
	return out
}
