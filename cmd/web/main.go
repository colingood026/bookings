package main

import (
	"github.com/alexedwards/scs/v2"
	"github.com/colingood026/bookings/internal/config"
	"github.com/colingood026/bookings/internal/handlers"
	"github.com/colingood026/bookings/internal/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	app.InProduction = false // only dev set to false
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session
	tc, err := render.CreateTemplateCacheV2()
	if err != nil {
		log.Fatal("can not create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = false
	// handlers
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	// render
	render.NewTemplates(&app)
	// route
	log.Printf("Starting service on port %s ...", portNumber)

	srv := http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
