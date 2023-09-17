package judges

import (
	"strings"

	"github.com/ptamarov/go-cards/app/card"
	"github.com/ptamarov/go-cards/app/history"
)

// NaiveJudge evaluates a guess by checking if the two strings are equal.
type NaiveJudge struct {
	CaseInsensitive   bool
	UmlautInsensitive bool
}

// EvaluateUserAction compares card.Answer is equal to g.Guess using strings.EqualFold. It
// returns 1 if EqualFold reports true, and 0 otherwise.
func (nj *NaiveJudge) EvaluateUserAction(card card.MemoryCard, g history.UserAction) float64 {
	if strings.EqualFold(card.Answer, g.Guess) {
		return 1.0
	} else {
		return 0.0
	}
}
