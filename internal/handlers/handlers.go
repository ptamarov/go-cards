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

	// if GetGuess is called, then DataCache should be cleared
	if m.App.SummaryDataCache != nil {
		m.App.SummaryDataCache = nil
	}
	if m.App.HomeDataCache != nil {
		m.App.HomeDataCache = nil
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
		// incorrect answer so dump (drop set to 0 by default)
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
		m.UpdateUserProgress() // user advances forward since guess was correct

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
			td := models.TemplateData{CardData: m.App.User.CurrentCard}
			td.Answer = Repo.App.User.CurrentCard.Answer
			m.App.CorrectCache = &td
			m.App.AlreadyAnswered = false
			m.App.User.CurrentCard = newCard
			m.App.User.LastAnswer = ""
			http.Redirect(w, r, "/correct", http.StatusSeeOther)
			return
		}
	}
}

func (m *Repository) ShowCorrectCard(w http.ResponseWriter, r *http.Request) {
	td := m.App.CorrectCache
	// progress and count always fetched from memory
	td.DailyGoal = m.App.User.DailyGoal
	td.AnsweredToday = m.App.User.AnsweredToday
	td.BarProgressPerc = toPercentage(m.App.User.AnsweredToday, m.App.User.DailyGoal)
	renders.RenderTemplate(w, r, "show-card-correct.page.tmpl", td)
}

// ShowCard shows the user a card and handles a post request from the user.
func (m *Repository) ShowCard(w http.ResponseWriter, r *http.Request) {
	m.App.Time = time.Now().UTC()

	// 1. Check if daily goal is reached. If reached, redirect.
	if m.App.User.IsDailyGoalReached() {
		m.App.User.DailyGoalReached = true
		http.Redirect(w, r, "/summary", http.StatusSeeOther)
		return
	}

	var td models.TemplateData // generate new template data to pass on

	// 2. Get the card to learn.
	if !m.App.NotFresh {
		newCard, err := m.DB.GetCardToLearn(m.App.User.UserID, m.App.User.DeckID)
		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/come-back-later", http.StatusSeeOther)
			return
		} else if err != nil {
			helpers.ServerError(w, err)
			return
		} else {
			td.CardData = newCard
			m.App.User.CurrentCard = newCard
		}
	} else {
		td.CardData = m.App.User.CurrentCard // get template information from memory
		td.Answer = Repo.App.User.LastAnswer
	}

	// progress and count always fetched from memory
	td.DailyGoal = m.App.User.DailyGoal
	td.AnsweredToday = m.App.User.AnsweredToday
	td.BarProgressPerc = toPercentage(m.App.User.AnsweredToday, m.App.User.DailyGoal)

	renders.RenderTemplate(w, r, "show-card.page.tmpl", &td)
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	tdPtr := &models.TemplateData{}

	if m.App.HomeDataCache != nil {
		tdPtr = m.App.HomeDataCache
	} else {
		err := m.GetAndPopulateTemplateWithCurrentStatisticsForHome(m.App.User.UserID, m.App.User.DeckID, tdPtr)
		if err != nil {
			helpers.ServerError(w, err)
			m.App.ErrorLog.Println("while populating the template with stats", err)
			return
		}
	}

	renders.RenderTemplate(w, r, "home.page.tmpl", tdPtr)
}

// ComeBackLater alerts the user that there no more cards to learn for the day.
func (m *Repository) ComeBackLater(w http.ResponseWriter, r *http.Request) {
	td := models.TemplateData{DailyGoalReached: m.App.User.DailyGoalReached}
	renders.RenderTemplate(w, r, "come-back-later.page.tmpl", &td)
}

// UpdateUserProgress updates the user progress.
func (m *Repository) UpdateUserProgress() {
	m.App.User.AnsweredToday++
}

// Retrieves and populates a template with statistics for the home page.
func (m *Repository) GetAndPopulateTemplateWithCurrentStatisticsForHome(userID, deckID uuid.UUID, td *models.TemplateData) error {

	// fetch all actions for the day
	actions, err := m.DB.GetAllActionsForToday(m.App.User.UserID, m.App.User.DeckID)
	if err != nil {
		m.App.ErrorLog.Println("while counting correct answers today", err)
		return err
	}
	m.App.InfoLog.Printf("fetched %d actions for today.\n", len(actions))

	var countCorrect int            // count how many actions were correct
	var countDrop int               // count how many actions were incorrect
	judge := m.Algorithm.GetJudge() // judge actions

	for _, action := range actions {
		card, err := m.DB.GetCardByID(action.CardID)
		if err != nil {
			m.App.ErrorLog.Println("while fetching card to get guess", err)
		} else if action.Drop == 1 {
			countDrop++
		} else {
			if judge.EvaluateUserAction(card, action) == 1 {
				countCorrect++
			}
		}
	}
	m.App.InfoLog.Printf("counted %d correct actions for today.\n", countCorrect)
	m.App.InfoLog.Printf("counted %d incorrect actions for today.\n", countDrop)
	m.App.User.AnsweredToday = countCorrect + countDrop
	countCardsReady, err := m.DB.GetCountCardsReady(m.App.User.UserID, m.App.User.DeckID)
	if err != nil {
		m.App.ErrorLog.Println("while counting cards that are ready", err)
		return err
	}
	m.App.InfoLog.Printf("counted %d cards ready.\n", countCardsReady)

	// stats for homepage
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
	notSeenPerC, progressPerC, learnedPerC := getStats(notSeen, inProgress, cardsLearned)

	// populate IntMap
	td.NotSeen = notSeen
	td.Learned = cardsLearned
	td.InProgress = inProgress
	td.Total = notSeen + cardsLearned + inProgress
	td.CardsReady = countCardsReady
	td.AnsweredToday = countCorrect + countDrop
	td.DailyGoal = m.App.User.DailyGoal

	// populate FloatMap
	td.BarProgressPerc = toPercentage(m.App.User.AnsweredToday, m.App.User.DailyGoal)
	td.NotSeenPerc = notSeenPerC
	td.LearnedPerc = learnedPerC
	td.InProgressPerc = progressPerC

	m.App.HomeDataCache = td
	return nil
}

// Retrieves and populates a template with statistics for the statistics page.
func (m *Repository) GetAndPopulateTemplateWithCurrentStatisticsForSummary(userID, deckID uuid.UUID, td *models.TemplateData) error {

	actions, err := m.DB.GetAllActionsForToday(m.App.User.UserID, m.App.User.DeckID)
	if err != nil {
		m.App.ErrorLog.Println("while getting all actions for today", err)
		return err
	}
	m.App.InfoLog.Printf("retrieved %d actions for today.\n", len(actions))

	var duration time.Duration // count time spent on actions
	var countCorrect int       // count how many actions were correct
	var countDropped int       // count how many actions were incorrect

	judge := m.Algorithm.GetJudge()
	for _, action := range actions {
		duration = duration + time.Duration(action.Duration*float64(time.Second))
		card, err := m.DB.GetCardByID(action.CardID)
		if err != nil {
			m.App.ErrorLog.Println("while fetching card to get guess", err)
			return err
		} else if action.Drop == 1 {
			countDropped++
		} else {
			if judge.EvaluateUserAction(card, action) == 1 {
				countCorrect++
			}
		}
	}

	m.App.InfoLog.Printf("%d correct guesses for today.\n", countCorrect)
	m.App.InfoLog.Printf("%d incorrect guesses for today.\n", countDropped)

	sec := int(duration.Seconds())
	min := sec / 60
	sec = sec - min*60

	td.Minutes = min
	td.Seconds = sec
	td.AnsweredToday = countDropped + countCorrect

	td.CorrectRate = toPercentage(countCorrect, countCorrect+countDropped)

	m.App.SummaryDataCache = td
	return nil
}

// Summary shows a summary of the learning session.
func (m *Repository) Summary(w http.ResponseWriter, r *http.Request) {

	tdPtr := &models.TemplateData{}

	if m.App.SummaryDataCache != nil {
		tdPtr = m.App.SummaryDataCache
	} else {
		err := m.GetAndPopulateTemplateWithCurrentStatisticsForSummary(m.App.User.UserID, m.App.User.DeckID, tdPtr)
		if err != nil {
			helpers.ServerError(w, err)
			m.App.ErrorLog.Println("while populating the template with stats", err)
			return
		}
	}

	renders.RenderTemplate(w, r, "summary.page.tmpl", tdPtr)

}
