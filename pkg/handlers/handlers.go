package handlers

import (
	"net/http"
	"time"

	"github.com/ptamarov/go-cards/app/history"
	"github.com/ptamarov/go-cards/pkg/config"
	"github.com/ptamarov/go-cards/pkg/driver"
	"github.com/ptamarov/go-cards/pkg/helpers"
	"github.com/ptamarov/go-cards/pkg/models"
	"github.com/ptamarov/go-cards/pkg/renders"
	"github.com/ptamarov/go-cards/pkg/repository"
	"github.com/ptamarov/go-cards/pkg/repository/dbrepo"
)

// The repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepository
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{App: a, DB: dbrepo.NewPostgresRepo(db.SQL, a)}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// GetGuess gets the user guess and processes it, and redirects to show a card
func (m *Repository) GetGuess(w http.ResponseWriter, r *http.Request) {
	// once a guess is posted, session has started
	// must be a better way to do this
	m.App.NotFresh = true

	// measure delta
	newTime := time.Now()
	delta := newTime.Sub(m.App.Time).Seconds()

	m.App.InfoLog.Println("DURATION:", delta)
	guess := r.Form.Get("user_guess")
	lastCardData := Repo.App.UserData.CardData // get last card shown to user

	m.App.InfoLog.Println("USER_GUESS:", guess)
	m.App.InfoLog.Println("EXPECTED_ANSWER:", lastCardData.Answer)

	// prepare action payload
	newAction := history.UserAction{
		UserID:   m.App.UserID,
		CardID:   lastCardData.ID,
		DeckID:   m.App.UserID,
		Guess:    guess,
		Duration: delta,
		Date:     time.Now(),
	}

	// record in database
	err := m.DB.RecordActionForUserAndDeck(newAction)
	if err != nil {
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println(err)
	}

	// check if user guess was right
	if guess == lastCardData.Answer {
		// if right, get new card and populate template data
		newCardData, err := m.DB.GetRandomCardInDatabase()
		if err != nil {
			helpers.ServerError(w, err)
			panic(err)
		}
		Repo.App.UserData.CardData = newCardData
		Repo.App.UserData.LastAnswer = ""
		Repo.UpdateUserProgress()
		m.App.InfoLog.Println("NEW_PROMPT:", newCardData.Prompt)
		http.Redirect(w, r, "/show-card", http.StatusTemporaryRedirect)

	} else {
		Repo.App.UserData.CardData = lastCardData
		Repo.App.UserData.LastAnswer = lastCardData.Answer
		http.Redirect(w, r, "/show-card", http.StatusTemporaryRedirect)
	}

}

// ShowCard shows the user a card and handles a post request from the user
func (m *Repository) ShowCard(w http.ResponseWriter, r *http.Request) {
	m.App.Time = time.Now()

	td := models.TemplateData{} // generate new template data to pass on
	if !m.App.NotFresh {
		newCard, err := m.DB.GetRandomCardInDatabase()
		if err != nil {
			m.App.ErrorLog.Println(err)
		} else {
			td.CardData = newCard
			m.App.UserData.CardData = newCard
		}
	} else {
		// get template information from memory
		td.CardData = m.App.UserData.CardData
		td.Answer = m.App.UserData.LastAnswer
	}

	// progress and count always fetched from memory
	td.Progress = m.App.UserData.Progress
	td.Count = m.App.UserData.Count

	renders.RenderTemplate(w, r, "show-card.page.tmpl", &td)
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	renders.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) UpdateUserGuess(w http.ResponseWriter, r *http.Request) {
	td := models.TemplateData{}
	renders.RenderTemplate(w, r, "home.page.tmpl", &td)
}

func (m *Repository) UpdateUserProgress() {
	m.App.UserData.Progress++
	m.App.UserData.Progress++
	m.App.UserData.Count++
}
