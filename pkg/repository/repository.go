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
	GetAllActionsForDeck(userID, deckID uuid.UUID) ([]history.UserAction, error)
	GetAllActionsForCard(userID, deckID, cardID uuid.UUID) ([]history.UserAction, error)
	GetAnsweredCorrectlyFromTo(userID, deckID uuid.UUID, start, end time.Time, judge judges.Judge) (int, error)
	GetAnsweredCorrectlyToday(userID, deckID uuid.UUID, j judges.Judge) (int, error)
	GetRedactedPromptFromCardID(cardID uuid.UUID) (string, error)
	GetTimeSeen(userID, deckID, cardID uuid.UUID) (int, error)

	// STATISTICS
	GetCountCardsReady(userID, deckID uuid.UUID) (int, error)
	GetCountCardsInProgress(userID, deckID uuid.UUID) (int, error)
	GetCountCardsNotSeen(userID, deckID uuid.UUID) (int, error)
	GetCountCardsLearned(userID, deckID uuid.UUID) (int, error)

	GetRandomCardInDatabase() (card.MemoryCard, error)
	// PUT ROUTINES
	RecordAction(history.UserAction) error
	UpdateCardStatus(userID, deckID, cardID uuid.UUID, status card.CardStatus) error
}
