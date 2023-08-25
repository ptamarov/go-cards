package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-sql-driver/mysql"
	"github.com/ptamarov/go-cards/database"
	"github.com/ptamarov/go-cards/pkg/config"
	"github.com/ptamarov/go-cards/pkg/handlers"
	"github.com/ptamarov/go-cards/pkg/renders"
	"github.com/ptamarov/go-cards/translate"
)

const portNumber = ":8880"

var app config.AppConfig // make it available to middleware too
var session *scs.SessionManager

func main() {
	// change this to true when in production
	app.InProduction = false

	// setting up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour // session lasts for 24 hours, can modify for more secure session handling
	session.Cookie.Persist = true     // session will persist after session ends
	session.Cookie.Secure = app.InProduction
	session.Cookie.SameSite = http.SameSiteLaxMode

	// make the session variable available to other packages
	app.Session = session

	tc, err := renders.CreateTemplateCache() // create the template cache
	if err != nil {
		log.Fatal("main; cannot create template cache", err)
	}

	app.TemplateCache = tc // assign the tc to the app configuration variable
	app.UseCache = false

	repo := handlers.NewRepo(&app) // create a (pointer to a) repository variable
	handlers.NewHandlers(repo)     // set the app repo to this variable
	renders.NewTemplates(&app)     // renders gets its own app variable so it can access cached templates

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	log.Printf("Working directory: %s\n", wd)
	log.Printf("Starting application on port %s\n", portNumber)

	// open database from config
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "germancorpora",
	}

	db := database.OpenDatabaseFromConfig(cfg)
	repo.App.DataBase = db

	// fill repo with information for DeepL API
	repo.App.UserData.TargetLanguage = "EN-GB"
	repo.App.UserData.AuthKey = translate.AUTHKEY
	repo.App.GetTranslations = false

	// serving with multiplexer
	srv := &http.Server{
		Addr:    portNumber,
		Handler: Routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(fmt.Sprintln("fatal error while listening and serving", err))
}
