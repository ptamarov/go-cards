package card

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type MemoryCard struct {
	ID                uuid.UUID `json:"card_id"`
	Grammar           string    `json:"grammar"` // explains the grammar of the missing word (in context)
	Hint              string    `json:"hint"`    // a hint indicates what the prompt is missing (marked word)
	LanguageToLearn   string    `json:"lang_learn"`
	UserLanguage      string    `json:"lang_user"`
	Prompt            string    `json:"prompt"`             // a prompt must be a sentence with a unique *marked* word to learn
	WordToLearn       string    `json:"word_to_learn"`      // a hint indicates what the prompt is missing (marked word)
	PromptTranslation string    `json:"prompt_translation"` // translation of the prompt in the user's language
	Answer            string    `json:"answer"`             // translation of the prompt in the user's language
}

func GetWordToLearnFromPrompt(prompt string) (string, error) {
	count := strings.Count(prompt, "*")
	switch count {
	case 0:
		return "", errors.New("bad input: no opening * found")
	case 1:
		return "", errors.New("bad input: no closing * found")
	case 2:
		start, end := 0, 0
	outer:
		for i := 0; i < len(prompt); i++ {
			if string(prompt[i]) == "*" {
				start = i + 1
				end = start
				for end < len(prompt) && string(prompt[end]) != "*" {
					end++

				}
				break outer
			}
		}
		return prompt[start:end], nil

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
		start, end := 0, 0
	outer:
		for i := 0; i < len(prompt); i++ {
			if string(prompt[i]) == "*" {
				start = i + 1
				end = start
				for end < len(prompt) && string(prompt[end]) != "*" {
					end++

				}
				break outer
			}
		}
		return prompt[:start-1] + "_____" + prompt[end+1:], nil

	default:
		return "", errors.New("bad input: too many delimiters found")
	}
}
