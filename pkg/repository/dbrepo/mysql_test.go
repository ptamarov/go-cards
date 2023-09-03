package dbrepo

import (
	"fmt"
	"log"
	"testing"
)

func Test_WordThatExists(t *testing.T) {
	// TRY TO GET A WORD THAT EXISTS
	word := "foolish"
	id, err := SQLTest.GetWordIDFromWord(word)

	if err != nil {
		t.Errorf("%s not found !-> %s.\n", word, err)
	}
	log.Printf("The id of %s is %d.\n", word, id)
}

func Test_SentenceThatExists(t *testing.T) {
	// TRY TO GET A SENTENCE FROM A GOOD ID
	sentence_id := 2
	sentence, err := SQLTest.GetSentenceFromID(sentence_id)
	if err != nil {
		fmt.Printf("No sentence for id %d found !-> %s.\n", sentence_id, err)
		return
	}
	fmt.Printf("The id %d corresponds to: %s\n", sentence_id, sentence)
}

func Test_GetSentenceIDsFromWord(t *testing.T) {
	// TRY TO GET SENTENCE IDS FROM A WORD
	word1 := "is"
	sentence_ids, err := SQLTest.GetAllSentencesIDsContainingWord(word1)
	if err != nil {
		fmt.Printf("No sentences found for word %s!-> %s.\n", word1, err)
		return
	}
	fmt.Printf("The word \"%s\" corresponds to: %v\n", word1, sentence_ids)
}
