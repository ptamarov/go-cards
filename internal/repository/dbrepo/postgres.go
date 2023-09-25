package dbrepo

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/ptamarov/go-cards/app/algorithm"
	"github.com/ptamarov/go-cards/app/card"
	"github.com/ptamarov/go-cards/app/history"
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

// Counts all user actions on a given deck from start to end dates (UTC)
func (m *postgresDBRepo) GetCountAllActionsFromTo(userID, deckID uuid.UUID, start, end time.Time) (int, error) {
	var count int
	var err error

	query := `
	SELECT 	COUNT(*)
	FROM 	history
	WHERE 	user_id = $1
	AND		deck_id = $2
	AND	    created_at >= $3 
	AND 	created_at < $4
	`

	s, e := start.Format(GO_DATE_FORMAT), end.Format(GO_DATE_FORMAT)
	rows := m.DB.QueryRow(query, userID, deckID, s, e)

	if err != nil {
		m.App.ErrorLog.Println("while executing query", err)
		return count, err
	}

	if err := rows.Scan(&count); err != nil {
		m.App.ErrorLog.Println("while scanning count", err)
		return count, err
	}

	return count, nil
}

// Counts all user actions on a given deck today (UTC time)
func (m *postgresDBRepo) GetCountAllActionsForToday(userID, deckID uuid.UUID) (int, error) {
	dateNowUTC := time.Now().UTC().Format("2006-01-02")
	today, _ := time.Parse("2006-01-02", dateNowUTC)
	oneDay := (24 * 60) * time.Minute
	tomorrow := today.Add(oneDay)
	return m.GetCountAllActionsFromTo(userID, deckID, today, tomorrow)
}

// GetCardByID gets a card from its ID with a redacted prompt.
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
// which is larger or equal to current time in UTC.
func (m *postgresDBRepo) GetCardToLearn(userID, deckID uuid.UUID) (card.MemoryCard, error) {
	var newCard card.MemoryCard
	var cardID uuid.UUID
	timeNowUTC := time.Now().UTC().Format(GO_TIMESTAMP_FORMAT)
	query := `
	SELECT 		card_id
	FROM 		decks
	WHERE 		user_id = $1
	AND 		deck_id = $2
	AND			date_ready < $3 -- Card is ready to be learned.
	AND 		card_progress != 5			   -- Card is not learned.
	AND 		card_progress != -1 		   -- Card is not inactive.
	ORDER BY 	date_ready ASC
	LIMIT(1)
	`
	row := m.DB.QueryRow(query, userID, deckID, timeNowUTC)
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

// GetCardStatus gets the status of a card for a userID and a deckID.
func (m *postgresDBRepo) GetCardStatus(userID, deckID, cardID uuid.UUID) (algorithm.CardStatus, error) {
	var cardStatus algorithm.CardStatus

	query := `
	SELECT 	date_ready, card_learned, card_progress, times_seen
	FROM 	decks 
	WHERE 	user_id = $1
	AND		card_id = $2
	AND 	deck_id = $3	
	`

	row := m.DB.QueryRow(query, userID, cardID, deckID)

	var learned int
	var timestamp string
	var timesSeen int
	err := row.Scan(&timestamp, &learned, &cardStatus.CardProgress, &timesSeen)
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
	cardStatus.TimesSeen = timesSeen

	return cardStatus, nil
}

func (m *postgresDBRepo) GetCountCardsReady(userID, deckID uuid.UUID) (int, error) {
	var count int

	timeNowUTC := time.Now().UTC().Format(GO_TIMESTAMP_FORMAT)

	query := `
	SELECT 		COUNT(*)
	FROM 		decks
	WHERE 		user_id = $1
	AND 		deck_id = $2
	AND			date_ready < $3 -- Card is ready to be learned.
	AND 		card_progress != 5			   -- Card is not learned.
	AND 		card_progress != -1 		   -- Card is not inactive.
	`

	row := m.DB.QueryRow(query, userID, deckID, timeNowUTC)
	err := row.Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}

func (m *postgresDBRepo) GetAllActionsForToday(userID, deckID uuid.UUID) ([]history.UserAction, error) {
	dateNowUTC := time.Now().UTC().Format("2006-01-02")
	today, _ := time.Parse("2006-01-02", dateNowUTC)
	oneDay := (24 * 60) * time.Minute
	tomorrow := today.Add(oneDay)
	return m.GetAllActionsFromTo(userID, deckID, today, tomorrow)
}

// UpdateCardStatus updates the status of a card for a userID and a deckID.
func (m *postgresDBRepo) UpdateCardStatus(userID, deckID, cardID uuid.UUID, status algorithm.CardStatus) error {
	dateString := status.NextAvailableDate.Format(GO_TIMESTAMP_FORMAT)

	query := `
		UPDATE 	decks 
		SET 	date_ready = $1, card_learned = $2, card_progress = $3, times_seen = $7
		WHERE 	user_id = $4
		AND		card_id = $5
		AND 	deck_id = $6
		`

	var cardLearned int
	if status.CardLearned {
		cardLearned = 1
	}

	_, err := m.DB.Exec(query,
		dateString,
		cardLearned,
		status.CardProgress,
		userID,
		cardID,
		deckID,
		status.TimesSeen)
	if err != nil {
		return err
	}
	return err
}

// GetAllActionsForCard gets all the actions performed by a user, in the given deck, for the given card.
func (m *postgresDBRepo) GetAllActionsForCard(userID, deckID, cardID uuid.UUID) ([]history.UserAction, error) {
	var actions []history.UserAction
	query := `
	SELECT 	guess, duration, created_at
	FROM 	history
	WHERE 	user_id = $1
	AND 	deck_id = $2
	AND 	card_id = $3
	AND 	drop_action = 0
	ORDER BY created_at ASC
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

// GetCardPromptFromCardID gets the redacted prompt of a card from its card ID.
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

// Returns a slice with all the actions the user has performed within the datetime limits.
func (m *postgresDBRepo) GetAllActionsFromTo(userID, deckID uuid.UUID, start, end time.Time) ([]history.UserAction, error) {
	var actions []history.UserAction

	query := `
	SELECT 	user_id, deck_id, card_id, guess, duration, drop_action
	FROM 	history
	WHERE 	user_id = $1
	AND		deck_id = $2
	AND	    created_at >= $3 
	AND 	created_at < $4
	`

	s, e := start.Format(GO_DATE_FORMAT), end.Format(GO_DATE_FORMAT)
	rows, err := m.DB.Query(query, userID, deckID, s, e)

	if err != nil {
		m.App.ErrorLog.Println("while executing query", err)
		return actions, err
	}

	for rows.Next() {
		var action history.UserAction
		err := rows.Scan(&action.UserID,
			&action.DeckID,
			&action.CardID,
			&action.Guess,
			&action.Duration,
			&action.Drop,
		)
		if err != nil {
			m.App.ErrorLog.Println("while scanning next row", err)
			return actions, err
		}
		actions = append(actions, action)
	}

	if err := rows.Err(); err != nil {
		m.App.ErrorLog.Println("while checking for error after iteration", err)
		return actions, err
	}

	return actions, err
}

// GetTimesCardAnsweredInDeck gets the number of times the user has answered a given
// card within a given deck.
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
	return count, err
}

// GetAllActionsForDeck gets all actions with the input userID and deckID.
func (m *postgresDBRepo) GetAllActionsForDeck(userid, deckid uuid.UUID) ([]history.UserAction, error) {
	return []history.UserAction{}, nil
}

func (m *postgresDBRepo) GetCountCardsLearned(userID, deckID uuid.UUID) (int, error) {
	query := `
	SELECT 	COUNT(*)
	FROM 	decks
	WHERE	user_id = $1
	AND 	deck_id = $2
	AND 	card_learned = 1
	`
	row := m.DB.QueryRow(query, userID, deckID)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}

func (m *postgresDBRepo) GetCountCardsNotSeen(userID, deckID uuid.UUID) (int, error) {
	query := `
	SELECT 	COUNT(*)
	FROM 	decks
	WHERE	user_id = $1
	AND 	deck_id = $2
	AND 	times_seen = 0
	`
	row := m.DB.QueryRow(query, userID, deckID)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}

func (m *postgresDBRepo) GetCountCardsInProgress(userID, deckID uuid.UUID) (int, error) {
	query := `
	SELECT 	COUNT(*)
	FROM 	decks
	WHERE	user_id = $1
	AND 	deck_id = $2
	AND 	times_seen > 0
	AND 	card_learned = 0
	`
	row := m.DB.QueryRow(query, userID, deckID)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}

// RecordACtion records a user action in the database.
func (m *postgresDBRepo) RecordAction(action history.UserAction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	statement := `
	INSERT INTO history (user_id, deck_id, card_id, guess, duration, created_at, drop_action) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	dateString := action.Date.Format(GO_TIMESTAMP_FORMAT)

	_, err := m.DB.ExecContext(ctx, statement,
		action.UserID,
		action.DeckID,
		action.CardID,
		action.Guess,
		action.Duration,
		dateString,
		action.Drop,
	)
	if err != nil {
		return err
	}
	return nil

}
