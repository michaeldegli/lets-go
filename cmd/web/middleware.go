package main

import (
	"log"
	"net/http"

	"github.com/justinas/nosurf"
)

// LogRequest ...
func LogRequest(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		pattern := `%s - "%s %s %s"`
		log.Printf(pattern, r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)

	})

}

// RequireLogin ...
func (app *App) RequireLogin(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		loggedIn, err := app.LoggedIn(r)

		if err != nil {
			app.ServerError(w, err)
			return
		}

		if !loggedIn {
			http.Redirect(w, r, "/user/login", 302)
			return
		}

		next.ServeHTTP(w, r)

	})

}

// SecureHeaders ...
func SecureHeaders(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header()["X-XSS-Protection"] = []string{"1; mode=block"}

		next.ServeHTTP(w, r)

	})
}

// NoSurf ...
func NoSurf(next http.HandlerFunc) http.Handler {

	csrfHandler := nosurf.New(next)
	//csrfHandler.SetBaseCookie(http.Cookie{
	//	HttpOnly: true,
	//	Path:     "/",
	//	Secure:   true,
	//})

	return csrfHandler
}
