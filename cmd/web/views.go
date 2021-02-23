package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/justinas/nosurf"
	"github.com/andrii-minchekov/lets-go/pkg/models"
)

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 15:04")
}

// HTMLData ...
type HTMLData struct {
	CSRFToken string
	Flash     string
	Form      interface{}
	LoggedIn  bool
	Path      string
	Snippet   *models.Snippet
	Snippets  []*models.Snippet
}

// RenderHTML ...
func (app *App) RenderHTML(w http.ResponseWriter, r *http.Request, page string, data *HTMLData) {

	if data == nil {
		data = &HTMLData{}
	}

	data.Path = r.URL.Path

	// Add a token to every page
	data.CSRFToken = nosurf.Token(r)

	// Add the loggedIn status
	var err error
	data.LoggedIn, err = app.LoggedIn(r)

	if err != nil {
		app.ServerError(w, err)
		return
	}

	files := []string{
		filepath.Join(app.HTMLDir, "base.html"),
		filepath.Join(app.HTMLDir, page),
	}

	// Initialize a template function map object. This is a map
	// of our custom template functions
	fm := template.FuncMap{
		"humanDate": humanDate,
	}

	// Our template.FuncMap must be registered with the template set
	// before we call ParseFiles
	// 1. Create template set
	// 2. Register our Map
	// 3. Parse

	ts, err := template.New("").Funcs(fm).ParseFiles(files...)

	if err != nil {
		log.Println(err.Error())
		app.ServerError(w, err)
		return
	}

	// Initialize a buffer
	buf := new(bytes.Buffer)

	err = ts.ExecuteTemplate(buf, "base", data)

	if err != nil {
		log.Println(err.Error())
		app.ServerError(w, err)
		return
	}

	buf.WriteTo(w)

}
