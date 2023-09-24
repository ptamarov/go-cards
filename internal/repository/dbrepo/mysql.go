package dbrepo

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/ptamarov/go-cards/app/card"
	"github.com/ptamarov/go-cards/translate"
)

type SQLRepo struct {
	DB *sql.DB
}

func MySQLDatabaseFromConfig(cfg mysql.Config) *SQLRepo {
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

	return &SQLRepo{DB: db}
}

func (m *SQLRepo) GetWordIDFromWord(word string) (int, error) {
	row := m.DB.QueryRow("SELECT word_id FROM words WHERE word = ?", word)

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

func (m *SQLRepo) GetWordFromWordID(id int) (string, error) {
	row := m.DB.QueryRow("SELECT word FROM words WHERE words_id = ?", id)

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

// GetSentenceFromID
func (m *SQLRepo) GetSentenceFromID(id int) (string, error) {
	row := m.DB.QueryRow(`
	SELECT sentence 
	FROM sentences 
	WHERE sentence_id = ?`, id)

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

func (m *SQLRepo) GetRandomSentenceIDsContainingWord(word string) (int, int, error) {
	word_id, err := m.GetWordIDFromWord(word)

	if err != nil {
		log.Fatal(err)
	}

	query_text := fmt.Sprintf(`
	SELECT 		sentence_id, position 
	FROM 		string_sentence_position 
	WHERE 		word_id = %d 
	ORDER BY 	RAND () 
	LIMIT 		1;`, word_id)

	row := m.DB.QueryRow(query_text)

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

func (m *SQLRepo) GetRandomWordID() (int, error) {
	query := `
	SELECT 		word_id 
	FROM 		words 
	ORDER BY 	RANDOM () 
	LIMIT 		1;
	`

	var id int

	row := m.DB.QueryRow(query)

	switch err := row.Scan(&id); err {
	case sql.ErrNoRows:
		return -1, err
	case nil:
		return id, nil
	default:
		panic(err)
	}
}

func (m *SQLRepo) GetRandomWord() (string, error) {
	query := `SELECT word 
	FROM words 
	ORDER BY RAND () 
	LIMIT 1;`

	var word string

	row := m.DB.QueryRow(query)

	switch err := row.Scan(&word); err {
	case sql.ErrNoRows:
		return "", err
	case nil:
		return word, nil
	default:
		panic(err)
	}
}

func (m *SQLRepo) GetRandomSentenceContainingWord(word string) (string, int, error) {
	randomSentenceID, position, err := m.GetRandomSentenceIDsContainingWord(word)
	if err != nil {
		return "", -1, err
	}

	sentence, err := m.GetSentenceFromID(randomSentenceID)

	if err != nil {
		return "", -1, err
	}

	return sentence, position, nil
}

// GetAllSentencesIDsContainingWord takes a word and returns all sentence ids
// for which the corresponding sentence contains that word.
func (m *SQLRepo) GetAllSentencesIDsContainingWord(word string) ([]int, error) {
	word_id, err := m.GetWordIDFromWord(word)

	if err != nil {
		log.Fatal(err)
	}

	var sentence_ids []int

	rows, err := m.DB.Query("SELECT sentence_id FROM word_sentence_position WHERE word_id = ?", word_id)
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

func (m *SQLRepo) GenerateRandomCardData(db *sql.DB, getTranslation bool, authKey string, targetLang string) (map[string]string, error) {
	cardData := make(map[string]string)

	// Pick a random sentence in DE from the database that contains this random string
	word, err := m.GetRandomWord()
	if err != nil {
		fmt.Println("error when getting random word", err)
		return cardData, err
	}

	log.Println("CHOSEN_WORD: ", word)

	// Get random sentence containing this word, and the position in which this word appears
	sentence, _, err := m.GetRandomSentenceContainingWord(word)
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
		wordToLearn, err = card.GetWordToLearnFromPrompt(translation)
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
