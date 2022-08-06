package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Use(app.logRequest, app.secureHeaders)

	//Snippet routes & index
	router.Group(func(router chi.Router) {
		router.Use(app.session.Enable, app.noSurf)
		router.Get("/", app.home)

		//  Creates a pattern route for "/snippet"
		router.Route("/snippet", func(router chi.Router) {

			// Creates a group of  snippet routes that requires authorization
			router.Group(func(router chi.Router) {
				router.Use(app.requireAuthenticatedUser)
				router.Get("/create", app.createSnippetForm)
				router.Post("/create", app.createSnippet)
			})

			router.Get("/{id}", app.showSnippet)
		})

	})

	// User authentication routes
	router.Route("/user", func(router chi.Router) {
		router.Use(app.session.Enable, app.noSurf)
		router.Get("/signup", app.signupUserForm)
		router.Post("/signup", app.signupUser)
		router.Get("/login", app.loginUserForm)
		router.Post("/login", app.loginUser)

		// Creates a group of user authentication routes that requires authorization
		router.Group(func(router chi.Router) {
			router.Post("/logout", app.logoutUser)
		})

	})

	router.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return router

}
