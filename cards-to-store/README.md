# Local Cards to Store Tool

This tool translates sentences using the DeepL free API and generates SQL INSERT statements for populating the `cards` table with flashcards.

## What Changed

### Features
- **CLI flags** for configurable language pairs:
  - `-lang-learn`: Language being learned (default: `DE`)
  - `-lang-user`: User's native language (default: `EN-GB`)
  - `-target-lang`: DeepL translation target (default: `EN`)
  - `-input`: Input sentence file (default: `sentences.txt`)
  - `-output`: Output SQL file (default: `new-cards.sql`)
  - `-auth-key-file`: Path to auth key file (default: `../local-auth.txt`)

- **Environment variable support**:
  - Set `DEEPL_AUTH_KEY` env var to use instead of the key file
  - Supports `.env` file via godotenv (already a project dependency)

- **Error handling**: Non-fatal errors (translation failures) are logged and skipped, allowing the tool to continue processing other sentences

## Usage

### Basic Usage (German to English)
```bash
cd local-cards-to-store
go run main.go
```
Reads from `sentences.txt`, translates to English, outputs to `new-cards.sql`.

### Spanish to Ukrainian
```bash
go run main.go \
  -lang-learn ES \
  -lang-user UK \
  -target-lang UK \
  -input sentences-spanish.txt \
  -output new-cards-es-uk.sql
```

### English to German (reverse)
```bash
go run main.go \
  -lang-learn EN \
  -lang-user DE \
  -target-lang DE \
  -input sentences-english.txt \
  -output new-cards-en-de.sql
```

### Using Environment Variable for Auth Key
```bash
export DEEPL_AUTH_KEY="your-deepl-api-key"
go run main.go
```

## Input File Format

Create a `.txt` file with one sentence per line. Mark the word to learn with asterisks:

```
-- German verbs (comment lines starting with - are skipped)
Aber das *haben* wir doch sogar gegessen.
Ich *ging* zum Markt.
Das *Buch* ist sehr interessant.

-- German nouns
Das ist eine *Katze*.
```

## Output

Generates SQL `INSERT` statements ready to be imported:

```sql
INSERT INTO cards (card_id, grammar, hint, lang_learn, lang_user, prompt, word_translation) VALUES ('...', '', 'haben', 'DE', 'EN-GB', 'Aber das *haben* wir doch sogar gegessen.', 'But we *did* eat that, didn't we?');
```

## Supported DeepL Language Codes

| Code | Language |
|------|----------|
| `DE` | German |
| `EN` | English |
| `ES` | Spanish |
| `FR` | French |
| `IT` | Italian |
| `UK` | Ukrainian |
| `RU` | Russian |
| `PL` | Polish |
| `PT` | Portuguese |
| `ZH` | Chinese |
| `JA` | Japanese |

See [DeepL API docs](https://www.deepl.com/docs-api/translate) for full list.

## Database Language Codes

Use ISO 639-1 codes with regional variants where applicable:

| Code | Language |
|------|----------|
| `DE` | German |
| `EN` | English (generic) |
| `EN-GB` | English (British) |
| `EN-BR` | English (Brazilian variant - custom) |
| `ES` | Spanish |
| `UK` | Ukrainian |
| `RU` | Russian |
| `FR` | French |

## Example Workflow: Add Spanish-Ukrainian Cards

1. **Create input file** (`sentences-es.txt`):
   ```
   -- Common Spanish phrases
   Buenos *días*, ¿cómo estás?
   Me *gusta* mucho este lugar.
   El *precio* es muy alto.
   ```

2. **Run the tool**:
   ```bash
   go run main.go \
     -lang-learn ES \
     -lang-user UK \
     -target-lang UK \
     -input sentences-es.txt \
     -output new-cards-es-uk.sql
   ```

3. **Review output** (check for correct translations):
   ```bash
   cat new-cards-es-uk.sql
   ```

4. **Import into database**:
   ```bash
   psql -U postgres -d go_cards < new-cards-es-uk.sql
   ```

5. **Repeat in reverse** for Ukrainian → Spanish by swapping `-lang-learn` and `-lang-user`:
   ```bash
   go run main.go \
     -lang-learn UK \
     -lang-user ES \
     -target-lang ES \
     -input sentences-uk.txt \
     -output new-cards-uk-es.sql
   ```

## Tips

- **Grammar field**: Currently empty (set to `''`). You can manually edit the SQL or enhance the tool to read grammar from a separate file.
- **Hint field**: Automatically set to the marked word. This helps learners remember what they're studying.
- **Validation**: Always review the generated SQL for translation quality before importing.
- **DeepL Free Tier**: Allows ~500,000 characters per month. Monitor usage to stay within limits.
