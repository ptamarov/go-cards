package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/ptamarov/go-cards/app/card"
	"github.com/ptamarov/go-cards/app/history"
	"github.com/ptamarov/go-cards/translate"
)

type Word = string

type Database struct {
	Path string
}

func OpenDatabaseFromConfig(cfg mysql.Config) *sql.DB {
	// Get a database handle.
	var err error
	var db *sql.DB

	db, err = sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal("error while opening database", err)
	}

	// Ping the database.
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal("error while pinging", pingErr)
	}
	fmt.Println("*** Connected! ***")

	return db
}

// GetAllSentencesIDsContainingWord takes a word and returns all sentence ids
// for which the corresponding sentence contains that word.
func GetAllSentencesIDsContainingWord(db *sql.DB, word Word) ([]int, error) {
	word_id, err := GetWordIDFromWord(db, word)

	if err != nil {
		log.Fatal(err)
	}

	var sentence_ids []int

	rows, err := db.Query("SELECT sentence_id FROM word_sentence_position WHERE word_id = ?", word_id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var sentence_id int
		err := rows.Scan(&sentence_id)

		if err != nil {
			fmt.Println(err)
		}
		sentence_ids = append(sentence_ids, sentence_id)
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
	}

	return sentence_ids, nil
}

func GetRandomSentenceIDsContainingWord(db *sql.DB, word Word) (int, int, error) {
	word_id, err := GetWordIDFromWord(db, word)

	if err != nil {
		log.Fatal(err)
	}

	query_text := fmt.Sprintf("SELECT sentence_id, position FROM word_sentence_position WHERE word_id = %d ORDER BY RAND () LIMIT 1;", word_id)

	row := db.QueryRow(query_text)

	var id int
	var position int

	switch err := row.Scan(&id, &position); err {
	case sql.ErrNoRows:
		return -1, -1, err
	case nil:
		return id, position, nil
	default:
		panic(err)
	}

}

func (d *Database) GetCardPromptFromCardID(id uuid.UUID) string {
	return ""
}

func (d *Database) GetCardsWithGuessForDate(t time.Time) []string {
	return []string{}
}

func (d *Database) GetNMostFrequentWords(n int) ([]Word, error) {
	return []Word{}, nil
}

func (d *Database) GetNumberOfCardsAnsweredCorrectlyForDate(t time.Time) int {
	return 0
}

func (d *Database) GetNumberOfCardsAnsweredForDate(t time.Time) int {
	return 0
}

func GetRandomWordID(db *sql.DB) (int, error) {
	query := "SELECT word_id FROM words ORDER BY RANDOM () LIMIT 1;"

	var id int

	row := db.QueryRow(query)

	switch err := row.Scan(&id); err {
	case sql.ErrNoRows:
		return -1, err
	case nil:
		return id, nil
	default:
		panic(err)
	}
}

func GetRandomWord(db *sql.DB) (string, error) {
	query := `SELECT word 
	FROM words 
	ORDER BY RAND () 
	LIMIT 1;`

	var word string

	row := db.QueryRow(query)

	switch err := row.Scan(&word); err {
	case sql.ErrNoRows:
		return "", err
	case nil:
		return word, nil
	default:
		panic(err)
	}
}

func GetRandomSentenceContainingWord(db *sql.DB, word Word) (string, int, error) {
	randomSentenceID, position, err := GetRandomSentenceIDsContainingWord(db, word)
	if err != nil {
		return "", -1, err
	}

	sentence, err := GetSentenceFromID(db, randomSentenceID)

	if err != nil {
		return "", -1, err
	}

	return sentence, position, nil
}

// GetSentenceFromID
func GetSentenceFromID(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT sentence FROM sentences WHERE sentence_id = ?", id)

	var sentence string

	switch err := row.Scan(&sentence); err {
	case sql.ErrNoRows:
		return "", err
	case nil:
		return sentence, nil
	default:
		panic(err)
	}
}

func (d *Database) GetTimesCardAnsweredInDeck(userid, deckid, cardid uuid.UUID) int {
	return 0
}

func (d *Database) GetUserHistoryForDeck(userid, deckid uuid.UUID) history.UserHistoryForDeck {
	return history.UserHistoryForDeck{}
}

func GetWordIDFromWord(db *sql.DB, word Word) (int, error) {
	row := db.QueryRow("SELECT word_id FROM words WHERE word = ?", word)

	var id int

	switch err := row.Scan(&id); err {
	case sql.ErrNoRows:
		return id, err
	case nil:
		return id, nil
	default:
		panic(err)
	}

}

func GetWordFromWordID(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT word FROM words WHERE word_id = ?", id)

	var word string

	switch err := row.Scan(&word); err {
	case sql.ErrNoRows:
		return word, err
	case nil:
		return word, nil
	default:
		panic(err)
	}
}

func (d *Database) RecordActionForUserAndDeck(uuid.UUID, uuid.UUID, history.Action) {
}

func GetRandomCardInDatabase(db *sql.DB) (map[string]string, error) {
	var cardData = make(map[string]string)

	query := `SELECT * 
	FROM cards 
	ORDER BY RAND () 
	LIMIT 1;`

	row := db.QueryRow(query)

	var card_id string
	var grammar string
	var hint string
	var lang_learn string
	var lang_user string
	var prompt string
	var translation string

	switch err := row.Scan(
		&lang_learn,
		&lang_user,
		&card_id,
		&prompt,
		&hint,
		&grammar,
		&translation); err {
	case sql.ErrNoRows:
		return cardData, err
	case nil:
		log.Println(prompt)
		cardData["card_id"] = card_id
		cardData["prompt"], err = card.GetRedactedPrompt(prompt)
		if err != nil {
			log.Printf("while fetching prompt for %s: %v\n", card_id, err)
		}
		cardData["grammar"] = grammar
		cardData["translation"] = translation
		cardData["hint"] = hint
		cardData["answer"], err = card.GetWordToLearnFromPrompt(prompt)
		if err != nil {
			log.Printf("while fetching answer for %s: %v\n", card_id, err)
		}
		return cardData, nil
	default:
		panic(err)
	}
}

func GenerateRandomCardData(db *sql.DB, getTranslation bool, authKey string, targetLang string) (map[string]string, error) {
	cardData := make(map[string]string)

	// Pick a random sentence in DE from the database that contains this random word
	word, err := GetRandomWord(db)
	if err != nil {
		fmt.Println("error when getting random word", err)
		return cardData, err
	}

	log.Println("CHOSEN_WORD: ", word)

	// Get random sentence containing this word, and the position in which this word appears
	sentence, _, err := GetRandomSentenceContainingWord(db, word)
	if err != nil {
		fmt.Println("error when getting sentence and position", err)
		return cardData, err
	}

	log.Println("CHOSEN_SENTENCE: ", sentence)

	prompt := translate.BlankOutWordInSentence(sentence, word)
	cardData["prompt"] = prompt

	processedSentence := translate.ProcessSentence(sentence, word)
	log.Println("PROCESSED_SENTENCE: ", processedSentence)

	// Translate the sentence using the DeepL API
	var translation string
	var wordToLearn string

	if getTranslation {
		translation, err = translate.GetTranslation(authKey, processedSentence, targetLang)
		if err != nil {
			fmt.Println("error when translating:", err)
			return cardData, err
		}
		log.Print("TRANSLATION:", translation)

		// get the word to learn in target language (user language)
		wordToLearn, err = translate.GetMarkedWordFromPrompt(translation)
		log.Println("WORD_TO_LEARN:", wordToLearn)
		if err != nil {
			log.Println("error while fetching word to learn:", err)
		}
		translation = strings.Replace(translation, "*"+wordToLearn+"*", wordToLearn, -1)
		cardData["word_translation"] = wordToLearn
		cardData["full_translation"] = translation
	} else {
		cardData["word_translation"] = "Translations are off"
		cardData["full_translation"] = "Translations are off"
	}

	cardData["answer"] = word
	return cardData, nil
}
