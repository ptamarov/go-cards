package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/ptamarov/go-cards/app/algorithm"
	"github.com/ptamarov/go-cards/app/history"
	"github.com/ptamarov/go-cards/app/judges"
	"github.com/ptamarov/go-cards/pkg/config"
	"github.com/ptamarov/go-cards/pkg/driver"
	"github.com/ptamarov/go-cards/pkg/helpers"
	"github.com/ptamarov/go-cards/pkg/models"
	"github.com/ptamarov/go-cards/pkg/renders"
	"github.com/ptamarov/go-cards/pkg/repository"
	"github.com/ptamarov/go-cards/pkg/repository/dbrepo"
)

// time conversion
const TIMESTAMP_FORMAT = "2006-01-02T15:04:05Z07:00"

// The repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App       *config.AppConfig
	DB        repository.DatabaseRepository
	Algorithm algorithm.NextCardAlgorithm
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB, algo algorithm.NextCardAlgorithm) *Repository {
	return &Repository{App: a, DB: dbrepo.NewPostgresRepo(db.SQL, a), Algorithm: algo}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// GetGuess gets the user guess and processes it, and redirects to show a card
func (m *Repository) GetGuess(w http.ResponseWriter, r *http.Request) {
	// once a guess is posted, session has started
	// must be a better way to do this
	if !m.App.NotFresh {
		m.App.NotFresh = true
	}

	// measure delta
	newTime := time.Now()
	delta := newTime.Sub(m.App.Time).Seconds()

	m.App.InfoLog.Println("DURATION:", delta)
	guess := r.Form.Get("user_guess")
	lastCardData := Repo.App.User.CurrentCard // get last card shown to user

	m.App.InfoLog.Println("USER_GUESS:", guess)
	m.App.InfoLog.Println("EXPECTED_ANSWER:", lastCardData.Answer)

	// prepare action payload
	newAction := history.UserAction{
		UserID:   m.App.User.UserID,
		CardID:   lastCardData.ID,
		DeckID:   m.App.User.DeckID,
		Guess:    guess,
		Duration: delta,
		Date:     time.Now().UTC(),
	}

	// record action to database
	err := m.DB.RecordAction(newAction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	userID := m.App.User.UserID
	cardID := lastCardData.ID
	deckID := m.App.User.DeckID

	actions, err := m.DB.GetAllActionsForCard(userID, deckID, cardID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	cardStatus, err := m.DB.GetCardStatus(userID, deckID, cardID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	oldProgress := cardStatus.CardProgress
	judge := &judges.LevenshsteinJudge{CaseInsensitive: true, UmlautInsensitive: true}
	newStatus := m.Algorithm.ComputeNewCardStatus(lastCardData, actions, judge, cardStatus)

	err = m.DB.UpdateCardStatus(userID, deckID, cardID, newStatus)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	madeProgress := (newStatus.CardProgress > oldProgress)

	if madeProgress {
		// if right, get new card and populate template data
		newCardData, err := m.DB.GetCardToLearn(m.App.User.UserID, m.App.User.DeckID)
		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/come-back-later", http.StatusSeeOther)
		} else if err != nil {
			helpers.ServerError(w, err)
			m.App.ErrorLog.Println("while getting top card:", err)
			return
		}
		Repo.App.User.CurrentCard = newCardData
		Repo.App.User.LastAnswer = ""
		Repo.UpdateUserProgress()
		m.App.InfoLog.Println("NEW_PROMPT:", newCardData.Prompt)
		http.Redirect(w, r, "/learn", http.StatusTemporaryRedirect)

	} else {
		Repo.App.User.CurrentCard = lastCardData
		Repo.App.User.LastAnswer = lastCardData.Answer
		http.Redirect(w, r, "/learn", http.StatusTemporaryRedirect)
	}

}

// ShowCard shows the user a card and handles a post request from the user
func (m *Repository) ShowCard(w http.ResponseWriter, r *http.Request) {
	m.App.Time = time.Now()

	if m.App.User.IsDailyGoalReached() {
		http.Redirect(w, r, "/come-back-later", http.StatusSeeOther)
	}

	var td models.TemplateData // generate new template data to pass on
	if !m.App.NotFresh {
		newCard, err := m.DB.GetCardToLearn(m.App.User.UserID, m.App.User.DeckID)
		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/come-back-later", http.StatusSeeOther)
		} else if err != nil {
			helpers.ServerError(w, err)
			m.App.ErrorLog.Println("while getting top card:", err)
			return
		} else {
			td.CardData = newCard
			m.App.User.CurrentCard = newCard
		}
	} else {
		td.CardData = m.App.User.CurrentCard // get template information from memory
		if td.StringMap == nil {
			td.StringMap = make(map[string]string)
		}
		td.StringMap["answer"] = Repo.App.User.LastAnswer
	}

	// progress and count always fetched from memory
	if td.IntMap == nil {
		td.IntMap = make(map[string]int)
	}
	td.IntMap["progress"] = 2 * m.App.User.CorrectToday
	td.IntMap["correct_today"] = m.App.User.CorrectToday

	m.App.InfoLog.Println("ANSWER IS:", td.StringMap["answer"])

	renders.RenderTemplate(w, r, "show-card.page.tmpl", &td)
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	dateNow := time.Now().UTC().Format("2006-01-02")
	today, _ := time.Parse("2006-01-02", dateNow)
	uID, dID := m.App.User.UserID, m.App.User.DeckID
	var correct int

	oneDay := (24 * 60) * time.Minute
	tomorrow := today.Add(oneDay)
	correct, err := m.DB.GetAnsweredCorrectlyFromTo(uID, dID, today, tomorrow, &judges.LevenshsteinJudge{})
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.App.User.CorrectToday = correct

	if m.App.User.IsDailyGoalReached() {
		http.Redirect(w, r, "/come-back-later", http.StatusSeeOther)
	}

	td := models.TemplateData{}
	td.IntMap = make(map[string]int)
	td.IntMap["progress"] = correct
	renders.RenderTemplate(w, r, "home.page.tmpl", &td)
}

func (m *Repository) ComeBackLater(w http.ResponseWriter, r *http.Request) {
	renders.RenderTemplate(w, r, "come-back-later.page.tmpl", &models.TemplateData{})
}

func (m *Repository) UpdateUserProgress() {
	m.App.User.CorrectToday++
}
