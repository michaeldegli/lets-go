package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

// Routes ...
func (app *App) Routes() http.Handler {

	mux := pat.New()

	// mux := http.NewServeMux()

	mux.Get("/", NoSurf(app.Home))
	mux.Get("/snippet/new", app.RequireLogin(NoSurf(app.NewSnippet)))
	mux.Post("/snippet/new", app.RequireLogin(NoSurf(app.CreateSnippet)))
	mux.Get("/snippet/:id", NoSurf(app.ShowSnippet))

	// Application Routes [User]
	mux.Get("/user/signup", NoSurf(app.SignupUser))
	mux.Post("/user/signup", NoSurf(app.CreateUser))
	mux.Get("/user/login", NoSurf(app.LoginUser))
	mux.Post("/user/login", NoSurf(app.VerifyUser))
	mux.Post("/user/logout", app.RequireLogin(NoSurf(app.LogoutUser)))

	fileServer := http.FileServer(http.Dir(app.StaticDir))

	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return LogRequest(SecureHeaders(mux))
}
