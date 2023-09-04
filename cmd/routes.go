package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ptamarov/go-cards/pkg/config"
	"github.com/ptamarov/go-cards/pkg/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	// Handle GET requests
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/get-guess", handlers.Repo.GetGuess)
	mux.Post("/get-guess", handlers.Repo.GetGuess)

	// Handle POST requests
	mux.Post("/show-card", handlers.Repo.ShowCard)
	mux.Get("/show-card", handlers.Repo.ShowCard)

	// create a file server
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
