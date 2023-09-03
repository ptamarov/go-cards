package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/ptamarov/go-cards/pkg/config"
)

var app *config.AppConfig // local variable to access app wide config

// NewHelpers sets up app config for helpers
func NewHelpers(a *config.AppConfig) {
	// get access to the app wide app configuration
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack()) // trace back the error
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
