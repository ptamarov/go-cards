package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
	"github.com/ptamarov/go-cards/app/algorithm"
	"github.com/ptamarov/go-cards/internal/config"
	"github.com/ptamarov/go-cards/internal/driver"
	"github.com/ptamarov/go-cards/internal/helpers"
	"github.com/ptamarov/go-cards/internal/renders"
	"github.com/ptamarov/go-cards/internal/repository/dbrepo"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../cmd/templates"
var functions = template.FuncMap{}

func TestMain(m *testing.M) {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.InProduction = false
	app.InfoLog = infoLog
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.Secure = app.InProduction
	session.Cookie.SameSite = http.SameSiteLaxMode
	app.Session = session

	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("main; cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseTemplateCache = true

	repo := NewTestRepo(&app, &driver.DB{})
	NewHandlers(repo)
	helpers.NewHelpers(&app)
	renders.NewRenders(&app)
	a := algorithm.NewSM2(true, true)
	repo.Algorithm = &a
	os.Exit(m.Run())
}

func getRoutes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	// Handle GET requests
	mux.Get("/", Repo.Home)
	mux.Get("/get-guess", Repo.GetGuess)
	mux.Get("/come-back-later", Repo.ComeBackLater)
	mux.Get("/learn", Repo.ShowCard)

	// Handle POST requests
	mux.Post("/learn", Repo.ShowCard)
	mux.Post("/get-guess", Repo.GetGuess)

	// create a file server
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads and saves the session in every request state
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	// myCache := make(map[string]*template.Template)
	myCache := map[string]*template.Template{}

	// get all of the files named *.page.tmpl from the ./templates folder
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates)) // Glob returns all files matching expression
	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page) // Base gets the last element of the path
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		// now get all layouts
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}

	return myCache, nil
}

// NewTestRepo creates a new repository
func NewTestRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{App: a, DB: dbrepo.NewTestingRepo(a)}
}
