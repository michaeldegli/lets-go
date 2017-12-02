package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"
)

// RunServer ...
func (app *App) RunServer() {

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:               tls.VersionTLS12,
		MaxVersion:               tls.VersionTLS12,
	}

	srv := &http.Server{
		Addr:           app.Addr,
		Handler:        app.Routes(),
		TLSConfig:      tlsConfig,
		IdleTimeout:    time.Minute,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   8 * time.Second,
		MaxHeaderBytes: 524288,
	}

	log.Printf("Starting server on %s", app.Addr)
	err := srv.ListenAndServeTLS(app.TLSCert, app.TLSKey)
	log.Fatal(err)

}
