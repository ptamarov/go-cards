package deck

import (
	"github.com/google/uuid"
	"github.com/ptamarov/go-cards/app/card"
)

type MemoryCardDeck struct {
	Cards    []card.MemoryCard
	DeckName string
	DeckID   uuid.UUID // TODO: make this a uuid
}

func (mcd *MemoryCardDeck) PutBack(card card.MemoryCard) {
	mcd.Cards = append(mcd.Cards, card)
}

func (mcd *MemoryCardDeck) PopLast() card.MemoryCard {
	var out card.MemoryCard
	mcd.Cards, out = mcd.Cards[:len(mcd.Cards)-1], mcd.Cards[len(mcd.Cards)-1]
	return out
}
