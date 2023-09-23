package dbrepo

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ptamarov/go-cards/internal/config"
	"github.com/ptamarov/go-cards/internal/driver"
	"github.com/ptamarov/go-cards/internal/helpers"
)

var TestDB *postgresDBRepo
var app config.AppConfig

func TestMain(m *testing.M) {
	var err error
	TestDB, err = run()
	if err != nil {
		log.Fatal(err)
	}
	defer TestDB.DB.Close()
	os.Exit(m.Run())
}

func run() (*postgresDBRepo, error) {
	// tell app about more complex types to be stored in session
	// change this to true when in production
	app.InProduction = false
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// connect to database
	app.InfoLog.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=go_cards user=ptamarov password=")
	if err != nil {
		log.Fatal("while connecting to database:", err)
	}

	helpers.NewHelpers(&app)

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	fmt.Printf("Working directory: %s\n", wd)

	pg := &postgresDBRepo{DB: db.SQL, App: &app}
	return pg, nil
}
