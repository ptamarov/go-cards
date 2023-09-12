package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/ptamarov/go-cards/app/card"
	"github.com/ptamarov/go-cards/app/history"
	"github.com/ptamarov/go-cards/app/judges"
)

type DatabaseRepository interface {
	// GET ROUTINES
	GetCardByID(cardID uuid.UUID) (card.MemoryCard, error)
	GetCardStatus(userID, deckID, cardID uuid.UUID) (card.CardStatus, error)
	GetCardToLearn(userID, deckID uuid.UUID) (card.MemoryCard, error)
	GetAllActionsForDeck(userid, deckid uuid.UUID) ([]history.UserAction, error)
	GetAllActionsForCard(userID, deckID, cardID uuid.UUID) ([]history.UserAction, error)
	GetAnsweredCorrectlyFromTo(userID, deckID uuid.UUID, start, end time.Time, judge judges.Judge) (int, error)
	GetRedactedPromptFromCardID(id uuid.UUID) (string, error)
	GetTimeSeen(userID, deckID, cardID uuid.UUID) (int, error)

	GetRandomCardInDatabase() (card.MemoryCard, error)
	// PUT ROUTINES
	RecordAction(history.UserAction) error
	UpdateCardStatus(userID, deckID, cardID uuid.UUID, status card.CardStatus) error
}
