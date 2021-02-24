package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/andrii-minchekov/lets-go/pkg/forms"
	"github.com/andrii-minchekov/lets-go/pkg/models"
)

// Home Display a "Hello from Snippetbox" message
func (app *App) Home(w http.ResponseWriter, r *http.Request) {

	// if r.URL.Path != "/" {
	// 	http.NotFound(w, r)
	// 	return
	// }

	snippets, err := app.Database.LatestSnippets()

	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.RenderHTML(w, r, "home.page.html", &HTMLData{
		Snippets: snippets,
	})

}

func (app *App) GetSnippets(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.Database.FindSnippets()

	if err != nil {
		app.ServerError(w, err)
		return
	}

	if snippets == nil {
		app.NotFound(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(snippets)
}

// ShowSnippet ...
func (app *App) ShowSnippet(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get(":id"))

	if err != nil || id < 1 {
		app.NotFound(w)
		return
	}

	snippet, err := app.Database.GetSnippet(id)

	if err != nil {
		app.ServerError(w, err)
		return
	}

	if snippet == nil {
		app.NotFound(w)
		return
	}

	session := app.Sessions.Load(r)
	flash, err := session.PopString(w, "flash")

	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.RenderHTML(w, r, "show.page.html", &HTMLData{
		Snippet: snippet,
		Flash:   flash,
	})
}

// NewSnippet ...
func (app *App) NewSnippet(w http.ResponseWriter, r *http.Request) {
	app.RenderHTML(w, r, "new.page.html", &HTMLData{
		Form: &forms.NewSnippet{},
	})
}

// CreateSnippet ...
func (app *App) CreateSnippet(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	err := r.ParseForm()

	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	form := &forms.NewSnippet{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: r.PostForm.Get("expires"),
	}

	if !form.Valid() {
		app.RenderHTML(w, r, "new.page.html", &HTMLData{Form: form})
		return
	}

	id, err := app.Database.InsertSnippet(form.Title, form.Content)

	if err != nil {
		app.ServerError(w, err)
		return
	}

	session := app.Sessions.Load(r)

	err = session.PutString(w, "flash", "Your snippet was saved!")

	if err != nil {
		app.ServerError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)

}

// SignupUser ...
func (app *App) SignupUser(w http.ResponseWriter, r *http.Request) {
	app.RenderHTML(w, r, "signup.page.html", &HTMLData{
		Form: &forms.SignupUser{},
	})
}

// CreateUser ...
func (app *App) CreateUser(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	form := &forms.SignupUser{
		Name:     r.PostForm.Get("name"),
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	if !form.Valid() {
		app.RenderHTML(w, r, "signup.page.html", &HTMLData{Form: form})
		return
	}

	err = app.Database.InsertUser(form.Name, form.Email, form.Password)

	if err != nil {
		app.ServerError(w, err)
		return
	}

	msg := "Your signup was successful. Please log in using your credentials."
	session := app.Sessions.Load(r)

	err = session.PutString(w, "flash", msg)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)

}

// LoginUser ...
func (app *App) LoginUser(w http.ResponseWriter, r *http.Request) {

	session := app.Sessions.Load(r)

	flash, err := session.PopString(w, "flash")
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.RenderHTML(w, r, "login.page.html", &HTMLData{
		Flash: flash,
		Form:  &forms.LoginUser{},
	})
}

// VerifyUser ...
func (app *App) VerifyUser(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	form := &forms.LoginUser{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	if !form.Valid() {
		app.RenderHTML(w, r, "login.page.html", &HTMLData{Form: form})
		return
	}

	currentUserID, err := app.Database.VerifyUser(form.Email, form.Password)

	if err == models.ErrInvalidCredentials {
		form.Failures["Generic"] = "Email or Password is incorrect"
		app.RenderHTML(w, r, "login.page.html", &HTMLData{Form: form})
		return
	} else if err != nil {
		app.ServerError(w, err)
		return
	}

	// Add the ID of the current user to the session
	session := app.Sessions.Load(r)
	err = session.PutInt(w, "currentUserID", currentUserID)

	if err != nil {
		app.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/snippet/new", http.StatusSeeOther)

}

// LogoutUser ...
func (app *App) LogoutUser(w http.ResponseWriter, r *http.Request) {

	// Remove the currentUserID from the session data.
	session := app.Sessions.Load(r)
	err := session.Remove(w, "currentUserID")

	if err != nil {
		app.ServerError(w, err)
		return

	}

	// Redirect the user to the homepage.
	http.Redirect(w, r, "/", 303)
}
