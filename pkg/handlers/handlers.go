package handlers

import (
	"log"
	"net/http"

	"github.com/ptamarov/go-cards/database"
	"github.com/ptamarov/go-cards/pkg/config"
	"github.com/ptamarov/go-cards/pkg/models"
	"github.com/ptamarov/go-cards/pkg/renders"
)

// The repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{App: a}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Guess(w http.ResponseWriter, r *http.Request) {
	guess := r.Form.Get("user_guess")
	lastCardData := Repo.App.UserData.CardData

	log.Println("USER_GUESS:", guess)

	td := models.TemplateData{}

	answer := lastCardData["answer"]
	log.Println("EXPECTED_ANSWER:", answer)

	if guess == answer {
		newCardData, err := database.GetRandomCardInDatabase(Repo.App.DataBase)
		log.Println("NEW_PROMPT:", newCardData["prompt"])
		if err != nil {
			panic(err)
		}
		Repo.App.UserData.CardData = newCardData
		td.Answer = ""
		td.CardData = newCardData
		if len(answer) != 0 {
			Repo.UpdateUserProgress()
		}
	} else {
		td.CardData = lastCardData
		td.Answer = answer
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
