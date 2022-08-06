package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serveError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) authenticatedUser(r *http.Request) int {

	// Retrieves the userID from the session; if userID doesn't exist it will return 0
	return app.session.GetInt(r, "userID")
}
