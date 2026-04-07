# go-cards

A spaced repetition flashcard learning application built with Go. Learn languages efficiently using the SuperMemo 2 (SM2) algorithm.

## Features

- **Spaced Repetition**: Uses the SM2 algorithm to optimize learning intervals based on recall difficulty
- **Language Learning**: Designed specifically for learning foreign languages with contextual sentences
- **Intelligent Answer Judging**: Uses Levenshtein distance to evaluate answers, tolerating minor spelling variations and case differences
- **Session Management**: Persistent user sessions for continuous learning progress
- **Web Interface**: Clean web UI for studying flashcards

## Architecture

### Core Components

- **Algorithm**: Implements the SuperMemo 2 (SM2) spaced repetition algorithm (`app/algorithm/`)
- **Card Management**: Handles flashcard data structures and operations (`app/card/`)
- **Judges**: Evaluates user answers using configurable matching strategies (`app/judges/`)
- **User Management**: Tracks user profiles and learning progress (`app/user/`)
- **Handlers**: HTTP request handlers for the web interface (`internal/handlers/`)
- **Database**: PostgreSQL integration for persistent storage (`internal/repository/`)

### Key Files

- `cmd/main.go` - Application entry point and configuration
- `cmd/routes.go` - HTTP route definitions
- `cmd/middleware.go` - Middleware for CSRF protection and session management
- `internal/config/` - Application configuration
- `internal/driver/` - Database connection management
- `internal/helpers/` - Template and utility helpers
- `internal/renders/` - Template rendering

## Getting Started

### Prerequisites

- Go 1.20 or later
- PostgreSQL database
- Port 8880 available (configurable)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/ptamarov/go-cards.git
cd go-cards
```

2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file for your local database configuration:
```bash
cp .env.example .env
```

4. Edit `.env` with your PostgreSQL credentials:
```
DATABASE_URL=host=localhost port=5432 dbname=go_cards user=<your_username> password=<your_password>
```

5. Create a PostgreSQL database:
```bash
createdb go_cards
```

6. Set up the database schema and seed data:
```bash
psql -U <your_username> -d go_cards < internal/repository/create-tables.sql
psql -U <your_username> -d go_cards < internal/repository/populate-cards.sql
```

7. Run the application:
```bash
go run ./cmd/main.go
```

8. Open your browser to `http://localhost:8880`

**Note**: The app will use the `DATABASE_URL` from your `.env` file. If no `.env` file exists, it defaults to:
```
host=localhost port=5432 dbname=go_cards user=postgres password=
```

## Usage

### Learning Flow

1. **Home**: View your daily learning goal and start a study session
2. **Learn**: Review a flashcard with the target word hidden
3. **Guess**: Submit your answer (spelling variations are tolerated)
4. **Feedback**: See the correct answer and continue
5. **Summary**: Review your progress after completing a session

### Flashcard Format

Flashcards include:
- **Prompt**: A sentence with the word to learn marked with `*asterisks*`
- **Prompt Translation**: Translation of the sentence in your native language
- **Grammar**: Grammatical notes about the target word
- **Hint**: A helpful hint for remembering the word
- **Answer**: The correct answer
- **Language Fields**: Target language and user's language

Example:
```
Prompt: "Ich *bin* Deutscher."
Answer: "bin"
Hint: "to be (present tense)"
Grammar: "1st person singular present of 'sein'"
```

## SM2 Algorithm

The SuperMemo 2 algorithm tracks:
- **Ease Factor**: Difficulty rating that affects future review intervals (default 2.5)
- **Interval**: Days until the card is shown again
- **Progress**: Learning progression score (0-5)
- **Times Seen**: How many times the card has been reviewed

Cards graduate from learning to reviewed status as they're successfully recalled.

## Answer Judgment

The application includes two judge implementations:
- **Levenshtein Judge**: Tolerates spelling variations and can be configured for case/umlaut sensitivity
- **Naive Judge**: Exact matching

Default configuration uses Levenshtein with case and umlaut insensitivity for German language learning.

## Project Structure

```
go-cards/
├── cmd/                          # Command line entry point
│   ├── main.go                   # App initialization
│   ├── routes.go                 # HTTP routes
│   └── middleware.go             # CSRF & session middleware
├── app/                          # Business logic
│   ├── algorithm/                # SM2 spaced repetition algorithm
│   ├── card/                     # Card data structures
│   ├── judges/                   # Answer evaluation
│   ├── user/                     # User management
│   ├── history/                  # Learning history
│   └── orders/                   # Card ordering
├── internal/                     # Internal packages
│   ├── config/                   # Application configuration
│   ├── driver/                   # Database driver
│   ├── handlers/                 # HTTP handlers
│   ├── helpers/                  # Template helpers
│   ├── renders/                  # Template rendering
│   └── repository/               # Database queries & schema
├── static/                       # CSS and client-side assets
├── template/                     # HTML templates
└── go.mod                        # Module definition
```

## Dependencies

- `github.com/go-chi/chi/v5` - HTTP router
- `github.com/alexedwards/scs/v2` - Session management
- `github.com/jackc/pgx/v4` - PostgreSQL driver
- `github.com/texttheater/golang-levenshtein` - String similarity
- `github.com/justinas/nosurf` - CSRF protection
- `github.com/google/uuid` - UUID generation

## Testing

Run the test suite:
```bash
go test ./...
```

Test coverage includes:
- SM2 algorithm calculations
- Levenshtein distance judging
- Card operations
- User management

## Development

### Hot Reload

By default, templates are not cached (`app.UseTemplateCache = false` in `main.go`). This allows editing templates without restarting the server.

### Production Mode

To run in production:
1. Set `app.InProduction = true` in `cmd/main.go`
2. Configure secure session cookies
3. Set appropriate database connection parameters
4. Enable template caching for performance

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

