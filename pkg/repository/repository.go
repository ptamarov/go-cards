package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/ptamarov/go-cards/app/card"
	"github.com/ptamarov/go-cards/app/history"
)

type DatabaseRepository interface {
	GetCardPromptFromCardID(id uuid.UUID) string
	GetCardsWithGuessForDate(t time.Time) []string
	GetNumberOfCardsAnsweredCorrectlyForDate(t time.Time) int
	GetNumberOfCardsAnsweredForDate(t time.Time) int
	GetTimesCardAnsweredInDeck(userid, deckid, cardid uuid.UUID) int
	GetUserHistoryForDeck(userid, deckid uuid.UUID) history.UserHistoryForDeck
	RecordActionForUserAndDeck(history.UserAction) error
	GetRandomCardInDatabase() (card.MemoryCard, error)
}
