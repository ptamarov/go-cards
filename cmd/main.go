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
	"github.com/ptamarov/go-cards/app/user"
	"github.com/ptamarov/go-cards/internal/config"
	"github.com/ptamarov/go-cards/internal/driver"
	"github.com/ptamarov/go-cards/internal/handlers"
	"github.com/ptamarov/go-cards/internal/helpers"
	"github.com/ptamarov/go-cards/internal/renders"
)

const portNumber = ":8880"

var app config.AppConfig // make it available to middleware too
var session *scs.SessionManager

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Printf("*** Starting application on port %s ***\n", portNumber)

	// serving with multiplexer
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(),
	}

	err = srv.ListenAndServe()
	log.Fatal(fmt.Sprintln("fatal error while listening and serving", err))
}

func run() (*driver.DB, error) {
	gob.Register(card.MemoryCard{}) // tell app about more complex types to be stored in session
	app.InProduction = false        // change this to true when in production

	////// Set up logging /////////////////////////////////////////////////////
	app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	////// Set up the session /////////////////////////////////////////////////
	session = scs.New()
	session.Lifetime = 24 * time.Hour        // can modify for more secure session handling
	session.Cookie.Persist = true            // session will persist after session ends
	app.Session = session                    // make the session variable available to other modules
	session.Cookie.Secure = app.InProduction // make apps secure if in production
	session.Cookie.SameSite = http.SameSiteLaxMode
	fmt.Print("\n")

	////// Connect to database ////////////////////////////////////////////////
	app.InfoLog.Println("connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=go_cards user=ptamarov password=")
	if err != nil {
		log.Fatal("while connecting to database:", err)
	}
	app.InfoLog.Println("...connected!")

	////// Create template caches /////////////////////////////////////////////
	tc, err := renders.CreateTemplateCache()
	if err != nil {
		app.ErrorLog.Fatal("cannot create template cache", err)
		return nil, err
	}
	app.TemplateCache = tc       // assign cache to the app configuration variable
	app.UseTemplateCache = false // but do not use it to allow template edition and reloading

	////// Mock a user ////////////////////////////////////////////////////////
	///// This will be handled by a log-in page in the future
	app.User = user.User{UserName: "test_user", DailyGoal: 100}

	algo := algorithm.NewSM2(true, true)      // case and umlaut insensitive judge
	repo := handlers.NewRepo(&app, db, &algo) // create a (pointer to a) repository variable
	handlers.NewHandlers(repo)                // set the app repo to this variable
	helpers.NewHelpers(&app)                  // connect helpers
	renders.NewRenders(&app)                  // link app to renders

	return db, nil
}
