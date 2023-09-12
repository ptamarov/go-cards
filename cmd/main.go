package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/ptamarov/go-cards/app/algorithm"
	"github.com/ptamarov/go-cards/app/card"
	"github.com/ptamarov/go-cards/pkg/config"
	"github.com/ptamarov/go-cards/pkg/driver"
	"github.com/ptamarov/go-cards/pkg/handlers"
	"github.com/ptamarov/go-cards/pkg/helpers"
	"github.com/ptamarov/go-cards/pkg/renders"
)

const portNumber = ":8880"

var app config.AppConfig // make it available to middleware too
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Printf("Starting application on port %s\n", portNumber)

	// serving with multiplexer
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(fmt.Sprintln("fatal error while listening and serving", err))
}

func run() (*driver.DB, error) {
	// tell app about more complex types to be stored in session
	gob.Register(card.MemoryCard{})
	// change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// setting up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour // session lasts for 24 hours, can modify for more secure session handling
	session.Cookie.Persist = true     // session will persist after session ends
	session.Cookie.Secure = app.InProduction
	session.Cookie.SameSite = http.SameSiteLaxMode

	// make the session variable available to other packages
	app.Session = session

	fmt.Print("\n")
	// connect to database
	app.InfoLog.Println("connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=go_cards user=ptamarov password=")
	if err != nil {
		log.Fatal("while connecting to database:", err)
	}
	app.InfoLog.Println("...connected!")

	tc, err := renders.CreateTemplateCache() // create the template cache
	if err != nil {
		app.ErrorLog.Fatal("cannot create template cache")
		return nil, err
	}
	app.TemplateCache = tc // assign the tc to the app configuration variable
	app.UseCache = false

	algo := algorithm.SM2Algorithm{}
	repo := handlers.NewRepo(&app, db, &algo) // create a (pointer to a) repository variable
	handlers.NewHandlers(repo)                // set the app repo to this variable

	helpers.NewHelpers(&app)
	renders.NewRenders(&app) // link app to renders

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	fmt.Printf("Working directory: %s\n", wd)

	return db, nil
}
