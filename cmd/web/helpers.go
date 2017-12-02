package main

import (
	"net/http"
)

// LoggedIn ...
func (app *App) LoggedIn(r *http.Request) (bool, error) {

	session := app.Sessions.Load(r)
	loggedIn, err := session.Exists("currentUserID")

	if err != nil {
		return false, err
	}

	return loggedIn, nil
}
