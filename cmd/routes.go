package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ptamarov/go-cards/internal/handlers"
)

func routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	// Handle GET requests
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/get-guess", handlers.Repo.GetGuess)
	mux.Get("/come-back-later", handlers.Repo.ComeBackLater)
	mux.Get("/learn", handlers.Repo.ShowCard)
	mux.Get("/summary", handlers.Repo.Summary)
	mux.Get("/correct", handlers.Repo.ShowCorrectCard)

	// Handle POST requests
	mux.Post("/learn", handlers.Repo.ShowCard)
	mux.Post("/get-guess", handlers.Repo.GetGuess)

	// create a file server
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
