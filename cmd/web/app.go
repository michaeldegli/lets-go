package main

import (
	"github.com/alexedwards/scs"
	"github.com/michaeldegli/snippetbox.org/pkg/models"
)

// App defines a struct to hold the application-wide dependencies and configuration
// settings for our web application. For now we'll only include a HTMLDir
// field for the path to the HTML templates directory, but we'll add more
// to it as our build progresses.
type App struct {
	Addr      string
	Database  *models.Database
	HTMLDir   string
	StaticDir string
	Sessions  *scs.Manager
	TLSCert   string
	TLSKey    string
}
