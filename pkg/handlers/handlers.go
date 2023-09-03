package handlers

import (
	"net/http"

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

// Home is the handler for the home page
func (m *Repository) Guess(w http.ResponseWriter, r *http.Request) {
	guess := r.Form.Get("user_guess")

	lastCardData := Repo.App.UserData.CardData // get last card shown to user

	td := models.TemplateData{} // generate new template data to pass on

	answer := lastCardData.Answer

	m.App.InfoLog.Println("USER_GUESS:", guess)
	m.App.InfoLog.Println("EXPECTED_ANSWER:", guess)

	// check if user guess was right
	if guess == answer {
		// if right, get new card and populate template data
		newCardData, err := m.DB.GetRandomCardInDatabase()
		if err != nil {
			helpers.ServerError(w, err)
			panic(err)
		}
		Repo.App.UserData.CardData = newCardData
		td.Answer = ""
		td.CardData = newCardData

		if len(answer) != 0 {
			Repo.UpdateUserProgress()
		}

		m.App.InfoLog.Println("NEW_PROMPT:", newCardData.Prompt)

	} else {
		td.CardData = lastCardData
		td.Answer = answer // shows user what the answer was
	}
	td.Progress = Repo.App.UserData.Progress
	td.Count = Repo.App.UserData.Count

	renders.RenderTemplate(w, r, "guess.page.tmpl", &td)
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
