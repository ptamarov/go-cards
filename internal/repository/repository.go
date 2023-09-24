package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/ptamarov/go-cards/app/algorithm"
	"github.com/ptamarov/go-cards/app/card"
	"github.com/ptamarov/go-cards/app/history"
)

type DatabaseRepository interface {
	// RETRIEVE ROUTINES
	GetCardByID(cardID uuid.UUID) (card.MemoryCard, error)
	GetCardStatus(userID, deckID, cardID uuid.UUID) (algorithm.CardStatus, error)
	GetCardToLearn(userID, deckID uuid.UUID) (card.MemoryCard, error)
	GetAllActionsForDeck(userID, deckID uuid.UUID) ([]history.UserAction, error)
	GetAllActionsForCard(userID, deckID, cardID uuid.UUID) ([]history.UserAction, error)
	GetAllActionsFromTo(userID, deckID uuid.UUID, start, end time.Time) ([]history.UserAction, error)
	GetAllActionsForToday(userID, deckID uuid.UUID) ([]history.UserAction, error)
	GetRedactedPromptFromCardID(cardID uuid.UUID) (string, error)
	GetTimeSeen(userID, deckID, cardID uuid.UUID) (int, error)

	// RETRIEVE STATISTICS
	GetCountAllActionsForToday(userID, deckID uuid.UUID) (int, error)
	GetCountAllActionsFromTo(userID, deckID uuid.UUID, start, end time.Time) (int, error)
	GetCountCardsReady(userID, deckID uuid.UUID) (int, error)
	GetCountCardsInProgress(userID, deckID uuid.UUID) (int, error)
	GetCountCardsNotSeen(userID, deckID uuid.UUID) (int, error)
	GetCountCardsLearned(userID, deckID uuid.UUID) (int, error)

	// PUT ROUTINES
	RecordAction(history.UserAction) error
	UpdateCardStatus(userID, deckID, cardID uuid.UUID, status algorithm.CardStatus) error
}
