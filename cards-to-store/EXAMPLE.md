# Quick Start Examples

## Example 1: Add More German Cards

**Input file** (`sentences-de.txt`):
```
-- Everyday German
Ich *kaufe* ein Buch.
Sie *arbeitet* in Berlin.
Das *Kind* spielt im Park.
Wir *trinken* Kaffee.
```

**Command**:
```bash
go run main.go -input sentences-de.txt -output new-cards-de.sql
```

**Output** (`new-cards-de.sql`):
```sql
INSERT INTO cards (card_id, grammar, hint, lang_learn, lang_user, prompt, word_translation) VALUES ('...uuid1...', '', 'kaufe', 'DE', 'EN-GB', 'Ich *kaufe* ein Buch.', 'I *buy* a book.');
INSERT INTO cards (card_id, grammar, hint, lang_learn, lang_user, prompt, word_translation) VALUES ('...uuid2...', '', 'arbeitet', 'DE', 'EN-GB', 'Sie *arbeitet* in Berlin.', 'She *works* in Berlin.');
INSERT INTO cards (card_id, grammar, hint, lang_learn, lang_user, prompt, word_translation) VALUES ('...uuid3...', '', 'spielt', 'DE', 'EN-GB', 'Das *Kind* spielt im Park.', 'The *child* plays in the park.');
INSERT INTO cards (card_id, grammar, hint, lang_learn, lang_user, prompt, word_translation) VALUES ('...uuid4...', '', 'trinken', 'DE', 'EN-GB', 'Wir *trinken* Kaffee.', 'We *drink* coffee.');
```

**Import**:
```bash
psql -U postgres -d go_cards < new-cards-de.sql
```

---

## Example 2: Spanish to Ukrainian

**Input file** (`sentences-spanish.txt`):
```
-- Spanish nouns
El *gato* es blanco.
La *casa* es grande.
El *libro* está en la mesa.

-- Spanish verbs
Yo *como* manzanas.
Ellos *hablan* español.
```

**Command**:
```bash
go run main.go \
  -lang-learn ES \
  -lang-user UK \
  -target-lang UK \
  -input sentences-spanish.txt \
  -output new-cards-es-uk.sql
```

**Output**:
```sql
INSERT INTO cards (card_id, grammar, hint, lang_learn, lang_user, prompt, word_translation) VALUES ('...uuid1...', '', 'gato', 'ES', 'UK', 'El *gato* es blanco.', '*Кіт* білий.');
INSERT INTO cards (card_id, grammar, hint, lang_learn, lang_user, prompt, word_translation) VALUES ('...uuid2...', '', 'casa', 'ES', 'UK', 'La *casa* es grande.', '*Будинок* великий.');
INSERT INTO cards (card_id, grammar, hint, lang_learn, lang_user, prompt, word_translation) VALUES ('...uuid3...', '', 'libro', 'ES', 'UK', 'El *libro* está en la mesa.', '*Книга* на столі.');
INSERT INTO cards (card_id, grammar, hint, lang_learn, lang_user, prompt, word_translation) VALUES ('...uuid4...', '', 'como', 'ES', 'UK', 'Yo *como* manzanas.', 'Я їм *яблука*.');
INSERT INTO cards (card_id, grammar, hint, lang_learn, lang_user, prompt, word_translation) VALUES ('...uuid5...', '', 'hablan', 'ES', 'UK', 'Ellos *hablan* español.', 'Вони говорять *испанский*.');
```

**Import**:
```bash
psql -U postgres -d go_cards < new-cards-es-uk.sql
```

---

## Example 3: English to German (Reverse Direction)

**Input file** (`sentences-english.txt`):
```
-- Common English phrases
I *drink* coffee every morning.
She *works* at the hospital.
The *book* is very interesting.
```

**Command**:
```bash
go run main.go \
  -lang-learn EN \
  -lang-user DE \
  -target-lang DE \
  -input sentences-english.txt \
  -output new-cards-en-de.sql
```

**Output**:
```sql
INSERT INTO cards (card_id, grammar, hint, lang_learn, lang_user, prompt, word_translation) VALUES ('...uuid1...', '', 'drink', 'EN', 'DE', 'I *drink* coffee every morning.', 'Ich *trinke* jeden Morgen Kaffee.');
INSERT INTO cards (card_id, grammar, hint, lang_learn, lang_user, prompt, word_translation) VALUES ('...uuid2...', '', 'works', 'EN', 'DE', 'She *works* at the hospital.', 'Sie *arbeitet* im Krankenhaus.');
INSERT INTO cards (card_id, grammar, hint, lang_learn, lang_user, prompt, word_translation) VALUES ('...uuid3...', '', 'book', 'EN', 'DE', 'The *book* is very interesting.', 'Das *Buch* ist sehr interessant.');
```

---

## Example 4: Ukrainian to Spanish

**Input file** (`sentences-ukrainian.txt`):
```
-- Ukrainian everyday phrases
Я люблю *кілька* мов.
Вони *живуть* у Львові.
Це *моя* книга.
```

**Command**:
```bash
go run main.go \
  -lang-learn UK \
  -lang-user ES \
  -target-lang ES \
  -input sentences-ukrainian.txt \
  -output new-cards-uk-es.sql
```

**Output**:
```sql
INSERT INTO cards (card_id, grammar, hint, lang_learn, lang_user, prompt, word_translation) VALUES ('...uuid1...', '', 'кілька', 'UK', 'ES', 'Я люблю *кілька* мов.', 'Amo *varios* idiomas.');
INSERT INTO cards (card_id, grammar, hint, lang_learn, lang_user, prompt, word_translation) VALUES ('...uuid2...', '', 'живуть', 'UK', 'ES', 'Вони *живуть* у Львові.', 'Viven en *Lviv*.');
INSERT INTO cards (card_id, grammar, hint, lang_learn, lang_user, prompt, word_translation) VALUES ('...uuid3...', '', 'моя', 'UK', 'ES', 'Це *моя* книга.', 'Este es *mi* libro.');
```

---

## Tips for Quality Data

1. **Mark only one word per sentence**: The tool extracts the marked word as the "hint" and "answer"
   ```
   ✅ Good: "Ich *kaufe* ein Buch."
   ❌ Bad: "Ich *kaufe ein* Buch."
   ```

2. **Use complete sentences**: Flashcards work best with context
   ```
   ✅ Good: "Ich kaufe *täglich* Zeitung."
   ❌ Bad: "*täglich* newspaper"
   ```

3. **Check DeepL output**: Always review a few cards before bulk importing
   ```bash
   head -5 new-cards-es-uk.sql
   ```

4. **Use comment lines to organize**: Lines starting with `-` are skipped
   ```
   -- Verbs
   Yo *como* una manzana.
   
   -- Nouns
   El *gato* es negro.
   ```

5. **Set DEEPL_AUTH_KEY environment variable** to avoid typing the flag each time:
   ```bash
   export DEEPL_AUTH_KEY="your-api-key"
   # Now you can just run:
   go run main.go -input sentences.txt
   ```

---

## Troubleshooting

**Error: "could not read auth key file"**
- Ensure the auth key file exists at the path specified (default: `../local-auth.txt`)
- Or set the `DEEPL_AUTH_KEY` environment variable
- Or use `-auth-key-file` flag to specify a different path

**Error: "could not translate sentence"**
- Check your internet connection
- Verify your DeepL API key is valid
- Check if you've exceeded your monthly character limit

**Hint field is empty**
- Ensure each sentence in the input file has exactly one word marked with `*word*`
- Check the tool output for warnings about sentences that couldn't be parsed

**SQL import fails**
- Check that your database is running and accessible
- Ensure you're using the correct database name (default: `go_cards`)
- Verify the table schema matches what the INSERT expects
