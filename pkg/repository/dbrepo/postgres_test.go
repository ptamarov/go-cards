package dbrepo

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/ptamarov/go-cards/app/card"
)

type TestCaseGetCardByID struct {
	cardID       uuid.UUID
	expectedCard card.MemoryCard
}

// TODO: maintain a test database!
func TestGetCardByID(t *testing.T) {
	tests := []TestCaseGetCardByID{
		{cardID: uuid.MustParse("023daefa-b586-426f-b63e-2b57be380042"),
			expectedCard: card.MemoryCard{
				ID:                uuid.MustParse("023daefa-b586-426f-b63e-2b57be380042"),
				Grammar:           "verb, third person, sing.",
				Hint:              "warns of; phrasal verb",
				LanguageToLearn:   "DE",
				UserLanguage:      "EN-GB",
				Prompt:            "Der Gemeinde *warnt vor* Hochwasser.",
				PromptTranslation: "The community warns of floods.",
				Answer:            "warnt vor",
			}},
	}
	for _, test := range tests {
		gotCard, err := TestDB.GetCardByID(test.cardID)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(gotCard, test.expectedCard) {
			t.Errorf("cards not equal: \n expected: %v\n got: %v", test.expectedCard, gotCard)
		}
	}

}
