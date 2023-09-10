package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/ptamarov/go-cards/app/card"
	"github.com/ptamarov/go-cards/app/history"
)

type DatabaseRepository interface {
	// GET ROUTINES
	GetRedactedPromptFromCardID(id uuid.UUID) (string, error)
	GetCardByID(cardID uuid.UUID) (card.MemoryCard, error)
	GetTopCardFromDeck(userID, deckID uuid.UUID) (card.MemoryCard, error)
	GetNumberOfCardsAnsweredCorrectlyForEpoch(start, end time.Time) (int, error)
	GetTimesCardAnsweredInDeck(userid, deckid, cardid uuid.UUID) (int, error)
	GetUserHistoryForDeck(userid, deckid uuid.UUID) history.UserHistoryForDeck
	GetRandomCardInDatabase() (card.MemoryCard, error)

	// PUT ROUTINES
	RecordActionForUserAndDeck(history.UserAction) error
	UpdateCardSeeNextDate(userID, deckID, cardID uuid.UUID, newDate time.Time) error
}
