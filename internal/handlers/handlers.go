package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/ptamarov/go-cards/app/algorithm"
	"github.com/ptamarov/go-cards/app/history"
	"github.com/ptamarov/go-cards/internal/config"
	"github.com/ptamarov/go-cards/internal/driver"
	"github.com/ptamarov/go-cards/internal/helpers"
	"github.com/ptamarov/go-cards/internal/models"
	"github.com/ptamarov/go-cards/internal/renders"
	"github.com/ptamarov/go-cards/internal/repository"
	"github.com/ptamarov/go-cards/internal/repository/dbrepo"
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

	lastCardData := Repo.App.User.CurrentCard // get last card shown to user

	// check card status before new user action
	userID := m.App.User.UserID
	cardID := lastCardData.ID
	deckID := m.App.User.DeckID
	cardStatus, err := m.DB.GetCardStatus(userID, deckID, cardID)

	m.App.InfoLog.Println("[GetGuess] CARD STATUS:", cardStatus)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	oldProgress := cardStatus.CardProgress

	// measure delta
	newTime := time.Now().UTC()
	delta := newTime.Sub(m.App.Time).Seconds()

	guess := r.Form.Get("user_guess")

	m.App.InfoLog.Println("[GetGuess] OLD PROGRESS:", oldProgress)
	m.App.InfoLog.Println("[GetGuess] DURATION:", delta)
	m.App.InfoLog.Println("[GetGuess] USER GUESS:", guess)
	m.App.InfoLog.Println("[GetGuess] EXPECTED ANSWER:", lastCardData.Answer)

	// prepare action payload
	newAction := history.UserAction{
		UserID:   m.App.User.UserID,
		CardID:   lastCardData.ID,
		DeckID:   m.App.User.DeckID,
		Guess:    guess,
		Duration: delta,
		Date:     time.Now().UTC(),
	}

	metric := m.Algorithm.GetJudge().EvaluateUserAction(lastCardData, newAction)

	if metric != 1.0 {
		// incorrect answer so dump (drop set to 0)
		err = m.DB.RecordAction(newAction)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		actions, err := m.DB.GetAllActionsForCard(userID, deckID, cardID)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		newStatus := m.Algorithm.ComputeNewCardStatus(lastCardData, actions)
		m.App.InfoLog.Println("[GetGuess] NEW STATUS:", newStatus)

		err = m.DB.UpdateCardStatus(userID, deckID, cardID, newStatus)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		// show card again
		m.App.User.CurrentCard = lastCardData
		m.App.User.LastAnswer = lastCardData.Answer
		m.App.AlreadyAnswered = true // mark card as answered
		http.Redirect(w, r, "/learn", http.StatusTemporaryRedirect)
	} else {
		if m.App.AlreadyAnswered {
			// drop correct guess from stats since user saw answer
			newAction.Drop = 1
		}
		// dump if card is fresh and correctly answered
		err = m.DB.RecordAction(newAction)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		actions, err := m.DB.GetAllActionsForCard(userID, deckID, cardID)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		newStatus := m.Algorithm.ComputeNewCardStatus(lastCardData, actions)
		m.App.InfoLog.Println("[GetGuess] NEW STATUS:", newStatus)

		err = m.DB.UpdateCardStatus(userID, deckID, cardID, newStatus)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		// get new card and populate template data
		newCard, err := m.DB.GetCardToLearn(m.App.User.UserID, m.App.User.DeckID)
		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/come-back-later", http.StatusSeeOther)
			return
		} else if err != nil {
			helpers.ServerError(w, err)
			return
		} else {
			m.App.AlreadyAnswered = false
			m.App.User.CurrentCard = newCard
			m.App.User.LastAnswer = ""
			m.UpdateUserProgress()
			http.Redirect(w, r, "/learn", http.StatusTemporaryRedirect)
			return
		}
	}

}

// ShowCard shows the user a card and handles a post request from the user.
func (m *Repository) ShowCard(w http.ResponseWriter, r *http.Request) {
	m.App.Time = time.Now().UTC()

	// 1. Check if daily goal is reached. If reached, redirect.
	if m.App.User.IsDailyGoalReached() {
		m.App.User.DailyGoalReached = true
		http.Redirect(w, r, "/come-back-later", http.StatusSeeOther)
		return
	}

	var td models.TemplateData // generate new template data to pass on

	// 2. Get the card to learn.
	if !m.App.NotFresh {
		newCard, err := m.DB.GetCardToLearn(m.App.User.UserID, m.App.User.DeckID)
		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/come-back-later", http.StatusSeeOther)
		} else if err != nil {
			helpers.ServerError(w, err)
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
	if td.FloatMap == nil {
		td.FloatMap = make(map[string]float64)
	}
	td.IntMap["daily_goal"] = m.App.User.DailyGoal
	td.IntMap["correct_today"] = m.App.User.CorrectToday
	td.FloatMap["bar_progress_perc"] = toPercentage(m.App.User.CorrectToday, m.App.User.DailyGoal)

	renders.RenderTemplate(w, r, "show-card.page.tmpl", &td)
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	userID, deckID := m.App.User.UserID, m.App.User.DeckID

	var correct int
	var err error
	var count int

	// get number of cards answered correctly today
	correct, err = m.DB.GetAnsweredCorrectlyToday(userID, deckID, m.Algorithm.GetJudge())
	if err != nil {
		m.App.ErrorLog.Println("while counting card answered correctly", err)
		helpers.ServerError(w, err)
		return
	}
	m.App.User.CorrectToday = correct

	count, err = m.DB.GetCountCardsReady(userID, deckID)
	if err != nil {
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("while counting cards that are ready", err)
		return
	}
	td := models.TemplateData{}
	err = m.PopulateTemplateWithCurrentStatistics(userID, deckID, &td)
	if err != nil {
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("while populating the template with stats", err)
		return
	}
	td.IntMap["cards_ready"] = count

	renders.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

// ComeBackLater alerts the user that there no more cards to learn for the day.
func (m *Repository) ComeBackLater(w http.ResponseWriter, r *http.Request) {
	td := models.TemplateData{}
	td.BoolMap = make(map[string]bool)
	td.BoolMap["goal_reached"] = m.App.User.DailyGoalReached
	renders.RenderTemplate(w, r, "come-back-later.page.tmpl", &td)
}

// UpdateUserProgress updates the user progress.
func (m *Repository) UpdateUserProgress() {
	m.App.User.CorrectToday++
}

// PopulateTemplateWithCurrentStatistics populates a template with the following information
//
// For IntMap:
//
//	daily_goal	<- Users daily goal
//	correct_today	<- Cards answered correctly today
//	not_seen	<-  Cards not seen to date
//	learned		<- Cards learned to date
//	in_progress	<- Cards active but not learned
//	total		<- Total cards in active deck
//
// For FloatMap:
//
//	not_seen_percentage	 <- Percentage of cards not seen.
//	learned_perc 		 <- Percentage of cards learned.
//	in_progress_perc	 <- Percentage of cards in progress.
//	bar_progress_perc	 <- Progress bar daily goal percentage.
func (m *Repository) PopulateTemplateWithCurrentStatistics(userID, deckID uuid.UUID, td *models.TemplateData) error {

	if td.IntMap == nil {
		td.IntMap = make(map[string]int)
	}
	td.IntMap["daily_goal"] = m.App.User.DailyGoal
	td.IntMap["correct_today"] = m.App.User.CorrectToday

	// get stats
	inProgress, err := m.DB.GetCountCardsInProgress(userID, deckID)
	if err != nil {
		return err
	}
	cardsLearned, err := m.DB.GetCountCardsLearned(userID, deckID)
	if err != nil {
		return err
	}
	notSeen, err := m.DB.GetCountCardsNotSeen(userID, deckID)
	if err != nil {
		return err
	}
	td.IntMap["not_seen"] = notSeen
	td.IntMap["learned"] = cardsLearned
	td.IntMap["in_progress"] = inProgress
	td.IntMap["total"] = notSeen + cardsLearned + inProgress

	notSeenPerC, progressPerC, learnedPerC := getStats(notSeen, inProgress, cardsLearned)
	if td.FloatMap == nil {
		td.FloatMap = make(map[string]float64)
	}
	td.FloatMap["not_seen_perc"] = notSeenPerC
	td.FloatMap["learned_perc"] = learnedPerC
	td.FloatMap["in_progress_perc"] = progressPerC
	td.FloatMap["bar_progress_perc"] = toPercentage(m.App.User.CorrectToday, m.App.User.DailyGoal)

	return nil
}

// Summary shows a summary of the learning session.
func (m *Repository) Summary(w http.ResponseWriter, r *http.Request) {
	// Display:
	// Cards answered
	// Time spent
	// New words
	// Correct rate
}
