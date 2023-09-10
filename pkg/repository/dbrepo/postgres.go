package dbrepo

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/ptamarov/go-cards/app/card"
	"github.com/ptamarov/go-cards/app/history"
)

// GetRandomCardInDatabase gets a random card in the cards table.
func (m *postgresDBRepo) GetRandomCardInDatabase() (card.MemoryCard, error) {
	var newCard card.MemoryCard
	var db *sql.DB = m.DB

	query := `SELECT * 
	FROM cards 
	ORDER BY RANDOM() 
	LIMIT 1;`

	row := db.QueryRow(query)

	var rawPrompt string

	switch err := row.Scan(
		&newCard.ID,
		&newCard.Grammar,
		&newCard.Hint,
		&newCard.LanguageToLearn,
		&newCard.UserLanguage,
		&rawPrompt,
		&newCard.PromptTranslation); err {
	case sql.ErrNoRows:
		return newCard, err
	case nil:
		newCard.Prompt, err = card.GetRedactedPrompt(rawPrompt)
		if err != nil {
			log.Printf("while fetching prompt for %s: %v\n", newCard.ID, err)
		}
		newCard.Answer, err = card.GetWordToLearnFromPrompt(rawPrompt)
		if err != nil {
			log.Printf("while fetching answer for %s: %v\n", newCard.ID, err)
		}
		return newCard, nil
	default:
		panic(err)
	}
}

// GetCardByID gets a card from its ID with a raw prompt
func (m *postgresDBRepo) GetCardByID(cardID uuid.UUID) (card.MemoryCard, error) {
	var newCard card.MemoryCard

	query := `
		SELECT 	*
		FROM	cards
		WHERE card_id = $1;`

	row := m.DB.QueryRow(query, cardID)

	err := row.Scan(
		&newCard.ID,
		&newCard.Grammar,
		&newCard.Hint,
		&newCard.LanguageToLearn,
		&newCard.UserLanguage,
		&newCard.Prompt,
		&newCard.PromptTranslation,
	)
	if err != nil {
		return newCard, err
	}
	newCard.Answer, err = card.GetWordToLearnFromPrompt(newCard.Prompt)
	if err != nil {
		log.Printf("while fetching answer for %s: %v\n", newCard.ID, err)
	}
	return newCard, nil
}

// GetTopCardFromDeck gets the card for a given userID and given deckID that has the highest
// priority to be seen next. This means that the card's 'see next' is the smallest date
// which is larger or equal to time.Now(). Returns an error if no such card exists.
func (m *postgresDBRepo) GetTopCardFromDeck(userID, deckID uuid.UUID) (card.MemoryCard, error) {
	var newCard card.MemoryCard
	var cardID uuid.UUID

	query := `
	SELECT 		card_id
	FROM 		decks
	WHERE 		user_id = $1
	AND 		deck_id = $2
	AND			date_ready < CURRENT_TIMESTAMP -- Card is ready to be learned.
	AND 		card_progress != 5			   -- Card is not learned.
	AND 		card_progress != -1 		   -- Card is not inactive.
	ORDER BY 	date_ready ASC
	LIMIT(1)
	`

	row := m.DB.QueryRow(query, userID, deckID)

	err := row.Scan(&cardID)
	if err != nil {
		return newCard, err
	}

	newCard, err = m.GetCardByID(cardID)
	if err != nil {
		return newCard, err
	}

	return newCard, nil
}

// PutCardBackWithNewDate updates the time in the future the card will be seen next for a
// given userID within the input deckID.
func (m *postgresDBRepo) UpdateCardSeeNextDate(userID, deckID, cardID uuid.UUID, newDate time.Time) error {
	return nil
}

// GetCardPromptFromCardID gets the redacted prompt of a card from its ID.
func (m *postgresDBRepo) GetRedactedPromptFromCardID(cardID uuid.UUID) (string, error) {
	var rawPrompt string

	query := `
		SELECT 	prompt
		FROM	cards
		WHERE card_id = $1;`

	row := m.DB.QueryRow(query, cardID)

	err := row.Scan(&rawPrompt)
	if err != nil {
		return "", err
	}

	redactedPrompt, err := card.GetRedactedPrompt(rawPrompt)
	if err != nil {
		return "", err
	}

	return redactedPrompt, nil
}

// GetNumberOfCardsAnsweredCorrectlyForDate gets the number of cards the user has
// answered correctly within the given time interval.
func (m *postgresDBRepo) GetNumberOfCardsAnsweredCorrectlyForEpoch(start, end time.Time) (int, error) {
	// query history table
	var z int
	return z, nil
}

// GetTimesCardAnsweredInDeck gets the number of times the user has answered a given
// card within a given deck
func (m *postgresDBRepo) GetTimesCardAnsweredInDeck(userID, deckID, cardID uuid.UUID) (int, error) {
	return 0, nil
}

func (m *postgresDBRepo) GetUserHistoryForDeck(userid, deckid uuid.UUID) history.UserHistoryForDeck {
	return history.UserHistoryForDeck{}
}

func (m *postgresDBRepo) RecordActionForUserAndDeck(action history.UserAction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	statement := `INSERT INTO history (user_id, deck_id, card_id, guess, duration, created_at) 
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := m.DB.ExecContext(ctx, statement,
		action.UserID,
		action.DeckID,
		action.CardID,
		action.Guess,
		action.Duration,
		action.Date,
	)
	if err != nil {
		return err
	}
	return nil

}
