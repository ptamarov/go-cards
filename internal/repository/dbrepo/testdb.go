package dbrepo

import (
	"time"

	"github.com/google/uuid"
	"github.com/ptamarov/go-cards/app/algorithm"
	"github.com/ptamarov/go-cards/app/card"
	"github.com/ptamarov/go-cards/app/history"
	"github.com/ptamarov/go-cards/app/judges"
)

// GetCardByID gets a card from its ID with a redacted prompt.
func (m *testDBRepo) GetCardByID(cardID uuid.UUID) (card.MemoryCard, error) {
	var newCard card.MemoryCard
	return newCard, nil
}

// GetCardToLearn gets the card for a given userID and given deckID that has the highest
// priority to be seen next. This means that the card's 'see next' is the smallest date
// which is larger or equal to current time in UTC.
func (m *testDBRepo) GetCardToLearn(userID, deckID uuid.UUID) (card.MemoryCard, error) {
	var newCard card.MemoryCard
	return newCard, nil
}

// GetCardStatus gets the status of a card for a userID and a deckID.
func (m *testDBRepo) GetCardStatus(userID, deckID, cardID uuid.UUID) (algorithm.CardStatus, error) {
	var cardStatus algorithm.CardStatus
	return cardStatus, nil
}

func (m *testDBRepo) GetCountCardsReady(userID, deckID uuid.UUID) (int, error) {
	var count int
	return count, nil
}

func (m *testDBRepo) GetAnsweredCorrectlyToday(userID, deckID uuid.UUID, j judges.Judge) (int, error) {
	var correct int
	return correct, nil
}

// UpdateCardStatus updates the status of a card for a userID and a deckID.
func (m *testDBRepo) UpdateCardStatus(userID, deckID, cardID uuid.UUID, status algorithm.CardStatus) error {
	return nil
}

// GetAllActionsForCard gets all the actions performed by a user, in the given deck, for the given card.
func (m *testDBRepo) GetAllActionsForCard(userID, deckID, cardID uuid.UUID) ([]history.UserAction, error) {
	var actions []history.UserAction
	return actions, nil
}

// GetCardPromptFromCardID gets the redacted prompt of a card from its card ID.
func (m *testDBRepo) GetRedactedPromptFromCardID(cardID uuid.UUID) (string, error) {
	var redactedPrompt string
	return redactedPrompt, nil
}

// GetNumberOfCardsAnsweredCorrectlyForDate gets the number of cards the user has
// answered correctly within the given time interval.
func (m *testDBRepo) GetAnsweredCorrectlyFromTo(userID, deckID uuid.UUID, start, end time.Time, judge judges.Judge) (int, error) {
	var count int
	return count, nil
}

// GetTimesCardAnsweredInDeck gets the number of times the user has answered a given
// card within a given deck.
func (m *testDBRepo) GetTimeSeen(userID, deckID, cardID uuid.UUID) (int, error) {
	var count int
	return count, nil
}

// GetAllActionsForDeck gets all actions with the input userID and deckID.
func (m *testDBRepo) GetAllActionsForDeck(userid, deckid uuid.UUID) ([]history.UserAction, error) {
	return []history.UserAction{}, nil
}

func (m *testDBRepo) GetCountCardsLearned(userID, deckID uuid.UUID) (int, error) {
	var count int
	return count, nil
}

func (m *testDBRepo) GetCountCardsNotSeen(userID, deckID uuid.UUID) (int, error) {
	var count int
	return count, nil
}

func (m *testDBRepo) GetCountCardsInProgress(userID, deckID uuid.UUID) (int, error) {
	var count int
	return count, nil
}

// RecordACtion records a user action in the database.
func (m *testDBRepo) RecordAction(action history.UserAction) error {
	return nil
}
