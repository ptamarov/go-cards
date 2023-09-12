package dbrepo

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/ptamarov/go-cards/app/card"
	"github.com/ptamarov/go-cards/app/history"
	"github.com/ptamarov/go-cards/app/judges"
)

const GO_DATE_FORMAT = "2006-01-02"
const GO_TIMESTAMP_FORMAT = "2006-01-02T15:04:05Z07:00"

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

	err := row.Scan(
		&newCard.ID,
		&newCard.Grammar,
		&newCard.Hint,
		&newCard.LanguageToLearn,
		&newCard.UserLanguage,
		&rawPrompt,
		&newCard.PromptTranslation)

	if err != nil {
		return newCard, err
	}

	newCard.Prompt, err = card.GetRedactedPrompt(rawPrompt)
	if err != nil {
		log.Printf("while fetching prompt for %s: %v\n", newCard.ID, err)
		return newCard, err
	}

	newCard.Answer, err = card.GetWordToLearnFromPrompt(rawPrompt)
	if err != nil {
		log.Printf("while fetching answer for %s: %v\n", newCard.ID, err)
		return newCard, err
	}

	return newCard, nil
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
	newCard.Prompt, err = card.GetRedactedPrompt(newCard.Prompt)
	if err != nil {
		log.Printf("while redacting prompt for %s: %v\n", newCard.ID, err)
	}
	return newCard, nil
}

// GetCardToLearn gets the card for a given userID and given deckID that has the highest
// priority to be seen next. This means that the card's 'see next' is the smallest date
// which is larger or equal to time.Now(). Returns an error if no such card exists.
func (m *postgresDBRepo) GetCardToLearn(userID, deckID uuid.UUID) (card.MemoryCard, error) {
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

func (m *postgresDBRepo) GetCardStatus(userID, deckID, cardID uuid.UUID) (card.CardStatus, error) {
	var cardStatus card.CardStatus

	query := `
	SELECT 	date_ready, card_learned, card_progress
	FROM 	decks 
	WHERE 	user_id = $1
	AND		card_id = $2
	AND 	deck_id = $3	
	`

	row := m.DB.QueryRow(query, userID, cardID, deckID)

	var learned int
	var timestamp string

	err := row.Scan(&timestamp, &learned, &cardStatus.CardProgress)
	if err != nil {
		m.App.ErrorLog.Println("while scannig row", err)
		return cardStatus, err
	}

	dateTime, err := time.Parse(GO_TIMESTAMP_FORMAT, timestamp)
	if err != nil {
		m.App.ErrorLog.Println("while parsing date", err)
		return cardStatus, err
	}
	cardStatus.NextAvailableDate = dateTime
	cardStatus.CardLearned = (learned != 0)

	return cardStatus, nil
}

func (m *postgresDBRepo) UpdateCardStatus(userID, deckID, cardID uuid.UUID, status card.CardStatus) error {

	dateString := status.NextAvailableDate.Format(GO_TIMESTAMP_FORMAT)

	query := `
		UPDATE 	decks 
		SET 	date_ready = $1, card_learned = $2, card_progress = $3
		WHERE 	user_id = $4
		AND		card_id = $5
		AND 	deck_id = $6`

	var cardLearned int
	if status.CardLearned {
		cardLearned = 1
	}
	result, err := m.DB.Exec(query, dateString, cardLearned, status.CardProgress, userID, cardID, deckID)
	if err != nil {
		return err
	}
	rowsAff, err := result.RowsAffected()
	if err != nil {
		return err
	}
	m.App.InfoLog.Println("rows affected:", rowsAff)
	return err
}

func (m *postgresDBRepo) GetAllActionsForCard(userID, deckID, cardID uuid.UUID) ([]history.UserAction, error) {
	var actions []history.UserAction

	query := `
	SELECT 	guess, duration, created_at
	FROM 	history
	WHERE 	user_id = $1
	AND 	deck_id = $2
	AND 	card_id = $3
	`

	rows, err := m.DB.Query(query, userID, deckID, cardID)
	if err != nil {
		m.App.ErrorLog.Println("while executing query:", err)
		return actions, err
	}

	for rows.Next() {
		newAction := history.UserAction{UserID: userID, CardID: cardID, DeckID: deckID}
		var date string
		err := rows.Scan(
			&newAction.Guess,
			&newAction.Duration,
			&date,
		)
		if err != nil {
			m.App.ErrorLog.Println("while reading new row:", err)
			return actions, err
		}
		parsedDate, err := time.Parse(GO_TIMESTAMP_FORMAT, date)
		if err != nil {
			m.App.ErrorLog.Println("while parsing date:", err)
			return actions, err
		}
		newAction.Date = parsedDate
		actions = append(actions, newAction)
	}
	return actions, nil
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
func (m *postgresDBRepo) GetAnsweredCorrectlyFromTo(userID, deckID uuid.UUID, start, end time.Time, judge judges.Judge) (int, error) {
	var actions []history.UserAction
	var count int

	query := `
	SELECT 	user_id, deck_id, card_id, guess, duration
	FROM 	history
	WHERE 	user_id = $1
	AND		deck_id = $2
	AND	    created_at >= $3 
	AND 	created_at < $4
	`

	s, e := start.Format(GO_DATE_FORMAT), end.Format(GO_DATE_FORMAT)

	m.App.InfoLog.Printf(`query:
SELECT	user_id, deck_id, card_id, guess, duration
FROM	history
WHERE	user_id = %s
AND	deck_id = %s
AND	created_at >= %s
AND	created_at < %s
	`, userID, deckID, s, e)

	rows, err := m.DB.Query(query, userID, deckID, s, e)

	if err != nil {
		m.App.ErrorLog.Println("while executing query", err)
		return count, err
	}

	for rows.Next() {
		var action history.UserAction
		err := rows.Scan(&action.UserID,
			&action.DeckID,
			&action.CardID,
			&action.Guess,
			&action.Duration,
		)
		if err != nil {
			m.App.ErrorLog.Println("while scanning next row", err)
			return count, err
		}
		actions = append(actions, action)
		m.App.InfoLog.Println("action card ID", action.CardID)

	}

	if err := rows.Err(); err != nil {
		m.App.ErrorLog.Println("while checking for error after iteration", err)
		return count, err
	}

	m.App.InfoLog.Println("ACTIONS RETRIEVED:", len(actions))
	for _, action := range actions {
		card, err := m.GetCardByID(action.CardID)
		if err != nil {
			m.App.ErrorLog.Println("while fetching card to get guess", err)
		}
		if judge.EvaluateUserAction(card, action) == 1 {
			count++
		}
	}
	return count, nil
}

// GetTimesCardAnsweredInDeck gets the number of times the user has answered a given
// card within a given deck
func (m *postgresDBRepo) GetTimeSeen(userID, deckID, cardID uuid.UUID) (int, error) {
	var count int
	query := `
	SELECT 	COUNT(*)
	FROM 	history
	WHERE 	user_id = $1
	AND		deck_id = $2
	AND 	card_id = $3
	`
	row := m.DB.QueryRow(query, userID, deckID, cardID)
	err := row.Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}

func (m *postgresDBRepo) GetAllActionsForDeck(userid, deckid uuid.UUID) ([]history.UserAction, error) {
	return []history.UserAction{}, nil
}

func (m *postgresDBRepo) RecordAction(action history.UserAction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	statement := `INSERT INTO history (user_id, deck_id, card_id, guess, duration, created_at) 
	VALUES ($1, $2, $3, $4, $5, $6)`

	dateString := action.Date.Format(GO_TIMESTAMP_FORMAT)

	_, err := m.DB.ExecContext(ctx, statement,
		action.UserID,
		action.DeckID,
		action.CardID,
		action.Guess,
		action.Duration,
		dateString,
	)
	if err != nil {
		return err
	}
	return nil

}
