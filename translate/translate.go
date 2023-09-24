package translate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

func GetTranslation(authKey string, inputText string, targetLanguage string) (string, error) {
	// Define the request URL
	url := "https://api-free.deepl.com/v2/translate"
	var out string

	payload := map[string]interface{}{
		"text":        []string{inputText},
		"target_lang": targetLanguage,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return out, err
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return out, err
	}

	// Set the request headers
	req.Header.Set("Authorization", "DeepL-Auth-Key "+authKey)
	req.Header.Set("User-Agent", "myapp/0.0.0")
	req.Header.Set("Content-Type", "application/json")

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	// Read and print the response body (JSON output)
	var response map[string][]map[string]string

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return out, err
	}

	return response["translations"][0]["text"], nil
}

// Replaces all ocurrences of a word in a sentence by a string of six underscores: ______.
func BlankOutWordInSentence(sentence string, word string) string {
	re := regexp.MustCompile(fmt.Sprintf(`\b(?i)%s\b`, word))
	return re.ReplaceAllString(sentence, "______")
}

// Replaces all ocurrences of a word in a sentence by *word*.
func ProcessSentence(sentence string, word string) string {
	re := regexp.MustCompile(fmt.Sprintf(`\b(?i)%s\b`, word))
	return re.ReplaceAllString(sentence, "*"+word+"*")
}
