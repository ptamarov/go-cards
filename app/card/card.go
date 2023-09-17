package card

import (
	"errors"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

type MemoryCard struct {
	ID                uuid.UUID `json:"card_id"`            // unique card id
	Grammar           string    `json:"grammar"`            // explains the grammar of the missing word (in context)
	Hint              string    `json:"hint"`               // a hint indicates what the prompt is missing (marked word)
	LanguageToLearn   string    `json:"lang_learn"`         // the language to learn
	UserLanguage      string    `json:"lang_user"`          // the language the user uses to learn
	Prompt            string    `json:"prompt"`             // a prompt must be a sentence with a unique *marked* word to learn
	PromptTranslation string    `json:"prompt_translation"` // translation of the prompt in the user's language
	Answer            string    `json:"answer"`             // the word to learn
}

func GetWordToLearnFromPrompt(prompt string) (string, error) {
	count := strings.Count(prompt, "*")
	switch count {
	case 0:
		return "", errors.New("bad input: no opening * found")
	case 1:
		return "", errors.New("bad input: no closing * found")
	case 2:
		re := regexp.MustCompile(`\*(.*)\*`)
		return re.FindStringSubmatch(prompt)[1], nil
	default:
		return "", errors.New("bad input: too many delimiters found")
	}
}

func GetRedactedPrompt(prompt string) (string, error) {
	count := strings.Count(prompt, "*")
	switch count {
	case 0:
		return "", errors.New("bad input: no opening * found")
	case 1:
		return "", errors.New("bad input: no closing * found")
	case 2:
		var newPrompt string
		re := regexp.MustCompile(`\*.*\*`)
		newPrompt = re.ReplaceAllString(prompt, "_____")
		return newPrompt, nil

	default:
		return "", errors.New("bad input: too many delimiters found")
	}
}
