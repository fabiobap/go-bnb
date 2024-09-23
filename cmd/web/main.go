package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fabiobap/go-bnb/internal/config"
	"github.com/fabiobap/go-bnb/internal/handlers"
	"github.com/fabiobap/go-bnb/internal/models"
	"github.com/fabiobap/go-bnb/internal/render"

	"github.com/alexedwards/scs/v2"
)

const HTTP_PORT = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	gob.Register(models.Reservation{})

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)

	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	fmt.Printf("Starting application on port: %s", HTTP_PORT)

	srv := &http.Server{
		Addr:    HTTP_PORT,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
