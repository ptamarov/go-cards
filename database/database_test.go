package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-sql-driver/mysql"
)

// set up database configuration
var cfg = mysql.Config{
	User:   os.Getenv("DBUSER"),
	Passwd: os.Getenv("DBPASS"),
	Net:    "tcp",
	Addr:   "127.0.0.1:3306",
	DBName: "cardappdb",
}

var db *sql.DB

func NewDatabaseForTesting(t *testing.T) (Database, error) {
	tempDir := t.TempDir()
	// TODO: create an SQL database  inside tempDir
	pathToDatabase := ""
	return Database{Path: filepath.Join(tempDir, pathToDatabase)}, nil

}

func Test_wordThatExists(t *testing.T) {

	db = OpenDatabaseFromConfig(cfg)
	// TRY TO GET A WORD THAT EXISTS
	word := "foolish"
	id, err := GetWordIDFromWord(db, word)

	if err != nil {
		t.Errorf("%s not found !-> %s.\n", word, err)
	}

	log.Printf("The id of %s is %d.\n", word, id)
}

func Test_sentenceThatExists(t *testing.T) {
	// TRY TO GET A SENTENCE FROM A GOOD ID
	db = OpenDatabaseFromConfig(cfg)
	sentence_id := 2
	sentence, err := GetSentenceFromID(db, sentence_id)

	if err != nil {
		fmt.Printf("No sentence for id %d found !-> %s.\n", sentence_id, err)
		return
	}
	fmt.Printf("The id %d corresponds to: %s\n", sentence_id, sentence)
}

func Test_getSentenceIDsFromWord(t *testing.T) {
	// TRY TO GET SENTENCE IDS FROM A WORD
	db = OpenDatabaseFromConfig(cfg)
	word1 := "is"

	sentence_ids, err := GetAllSentencesIDsContainingWord(db, word1)

	if err != nil {
		fmt.Printf("No sentences found for word %s!-> %s.\n", word1, err)
		return
	}
	fmt.Printf("The word \"%s\" corresponds to: %v\n", word1, sentence_ids)
}
