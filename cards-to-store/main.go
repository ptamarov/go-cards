package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/ptamarov/go-cards/app/card"
	"github.com/ptamarov/go-cards/translate"
)

func main() {
	// Define CLI flags
	langLearn := flag.String("lang-learn", "DE", "Language to learn (e.g., DE, ES, UK)")
	langUser := flag.String("lang-user", "EN-GB", "User's native language (e.g., EN-GB, EN-BR, UK)")
	targetLang := flag.String("target-lang", "EN", "DeepL target language (e.g., EN, UK, ES, DE)")
	inputFile := flag.String("input", "sentences.txt", "Path to input sentence file")
	outputFile := flag.String("output", "new-cards.sql", "Path to output SQL file")
	authKeyFile := flag.String("auth-key-file", "../local-auth.txt", "Path to file containing DeepL auth key")
	flag.Parse()

	// Load auth key: try env var first, then fall back to file
	_ = godotenv.Load() // Load from .env if it exists
	authKey := os.Getenv("DEEPL_AUTH_KEY")

	if authKey == "" {
		// Fall back to reading from auth key file
		authKeyBytes, err := os.ReadFile(*authKeyFile)
		if err != nil {
			log.Fatalf("could not read auth key file (%s): %v", *authKeyFile, err)
		}
		authKey = strings.TrimSpace(string(authKeyBytes))
	}

	if authKey == "" {
		log.Fatal("DEEPL_AUTH_KEY environment variable not set and no auth key file provided")
	}

	// Read input sentence file
	sentenceBytes, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("could not read input file (%s): %v", *inputFile, err)
	}

	// Parse sentences (skip comment lines starting with '-' and empty lines)
	var toProcess []string
	for _, line := range strings.Split(string(sentenceBytes), "\n") {
		if len(line) > 1 && line[0] != '-' {
			toProcess = append(toProcess, line)
		}
	}

	// Translate sentences and build cards
	var cards []card.MemoryCard
	for _, sentence := range toProcess {
		newTranslation, err := translate.GetTranslation(authKey, sentence, *targetLang)
		if err != nil {
			log.Printf("warning: could not translate sentence '%s': %v\n", sentence, err)
			continue
		}

		// Extract the marked word from the source sentence (this becomes the answer)
		answer, err := card.GetWordToLearnFromPrompt(sentence)
		if err != nil {
			log.Printf("warning: could not extract word from prompt '%s': %v\n", sentence, err)
			continue
		}

		newCard := card.MemoryCard{
			ID:                uuid.New(),
			Grammar:           "", // Not provided by this tool; can be added manually later
			Hint:              answer,
			LanguageToLearn:   *langLearn,
			UserLanguage:      *langUser,
			Prompt:            sentence,
			PromptTranslation: newTranslation,
			Answer:            answer,
		}
		cards = append(cards, newCard)
	}

	// Write SQL to output file
	if err := cardsToSQL(cards, *outputFile); err != nil {
		log.Fatalf("could not write output file (%s): %v", *outputFile, err)
	}

	log.Printf("successfully wrote %d cards to %s\n", len(cards), *outputFile)
}

// cardsToSQL formats cards as SQL INSERT statements and writes them to a file.
func cardsToSQL(cards []card.MemoryCard, outputPath string) error {
	var sqlLines []string

	for _, c := range cards {
		line := fmt.Sprintf(
			"INSERT INTO cards (card_id, grammar, hint, lang_learn, lang_user, prompt, word_translation) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s');\n",
			c.ID, c.Grammar, c.Hint, c.LanguageToLearn, c.UserLanguage, c.Prompt, c.PromptTranslation,
		)
		sqlLines = append(sqlLines, line)
	}

	output := strings.Join(sqlLines, "")
	return os.WriteFile(outputPath, []byte(output), 0644)
}
