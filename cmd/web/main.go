package main

import (
	"database/sql"
	"flag"
	"log"
	"time"

	"github.com/alexedwards/scs"
	_ "github.com/lib/pq"
	"github.com/michaeldegli/snippetbox.org/pkg/models"
)

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "user=local dbname=snippetbox password=local sslmode=disable", "Postgres Datasource Name")
	secret := flag.String("secret", "s6Nd%+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	htmlDir := flag.String("html-dir", "./ui/html", "Path to HTML templates")
	staticDir := flag.String("static-dir", "./ui/static", "Path to static assets")
	tlsCert := flag.String("tls-cert", "./tls/cert.pem", "Path to TLS certificate")
	tlsKey := flag.String("tls-key", "./tls/key.pem", "Path to TLS key")

	flag.Parse()

	db := connect(*dsn)
	defer db.Close()

	// Initialize the Session Store
	sessionManager := scs.NewCookieManager(*secret)
	sessionManager.Lifetime(12 * time.Hour)
	sessionManager.Persist(true)

	// Initialize a new instance of App containing the dependencies.
	app := &App{
		Database:  &models.Database{db},
		HTMLDir:   *htmlDir,
		StaticDir: *staticDir,
		Sessions:  sessionManager,
		Addr:      *addr,
		TLSCert:   *tlsCert,
		TLSKey:    *tlsKey,
	}

	app.RunServer()
}

func connect(dsn string) *sql.DB {

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db

}
