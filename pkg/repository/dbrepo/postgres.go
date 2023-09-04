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

func (m *postgresDBRepo) GetRandomCardInDatabase() (card.MemoryCard, error) {
	var newCard card.MemoryCard

	var db *sql.DB = m.DB

	query := `SELECT * 
	FROM cards 
	ORDER BY RANDOM() 
	LIMIT 1;`

	row := db.QueryRow(query)

	var raw_prompt string

	switch err := row.Scan(
		&newCard.ID,
		&newCard.Grammar,
		&newCard.Hint,
		&newCard.LanguageToLearn,
		&newCard.UserLanguage,
		&raw_prompt,
		&newCard.PromptTranslation); err {
	case sql.ErrNoRows:
		return newCard, err
	case nil:
		log.Println(raw_prompt)
		newCard.Prompt, err = card.GetRedactedPrompt(raw_prompt)
		if err != nil {
			log.Printf("while fetching prompt for %s: %v\n", newCard.ID, err)
		}
		newCard.Answer, err = card.GetWordToLearnFromPrompt(raw_prompt)
		if err != nil {
			log.Printf("while fetching answer for %s: %v\n", newCard.ID, err)
		}
		return newCard, nil
	default:
		panic(err)
	}
}

func (m *postgresDBRepo) GetCardPromptFromCardID(id uuid.UUID) string {
	var s string
	return s
}

func (m *postgresDBRepo) GetCardsWithGuessForDate(t time.Time) []string {
	var ss []string
	return ss
}

func (m *postgresDBRepo) GetNumberOfCardsAnsweredCorrectlyForDate(t time.Time) int {
	var z int
	return z
}

func (m *postgresDBRepo) GetNumberOfCardsAnsweredForDate(t time.Time) int {
	return 0
}

func (m *postgresDBRepo) GetTimesCardAnsweredInDeck(userid, deckid, cardid uuid.UUID) int {
	return 0
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
