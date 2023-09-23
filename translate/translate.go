package translate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
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

// BlanOutWordInSentence replaces all ocurrences of a word in a sentence by
// a string of six underscores: ______.
func BlankOutWordInSentence(sentence string, word string) string {
	re := regexp.MustCompile(fmt.Sprintf(`(?i)%s`, word))
	return re.ReplaceAllString(sentence, "______")
}

func ProcessSentence(s string, word string) string {
	re := regexp.MustCompile(fmt.Sprintf(`(?i)%s`, word))
	return re.ReplaceAllString(s, "*"+word+"*")
}

// GetMarkedWordFromPrompt takes a sentence along with a unique
// marked word using the delimiter "*" and returns the marked word.
// "It is a *nice* day today" -> "nice".
// Returns an error if there are not exactly two ocurrences of * in the sentence.
func GetMarkedWordFromPrompt(sentence string) (string, error) {
	var out string

	count := strings.Count(sentence, "*")

	switch count {
	case 0:
		return out, fmt.Errorf("bad input: no delimiters found")
	case 1:
		return out, fmt.Errorf("bad input: only one delimiter found")
	case 2:
		start, end := 0, 0
		delimiter := "*"

		for string(sentence[start]) != delimiter {
			start++
		}
		end = start + 1
		for string(sentence[end]) != delimiter {
			end++
		}
		return sentence[start+1 : end], nil

	default:
		return out, fmt.Errorf("bad input: too many delimiters found")
	}
}
