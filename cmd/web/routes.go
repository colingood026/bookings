package main

import (
	"github.com/colingood026/bookings/pkg/config"
	"github.com/colingood026/bookings/pkg/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()
	//mux.Use(WriteToConsole) // just demo how to implement middleware
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/generals-quarters", handlers.Repo.Generals)
	mux.Get("/majors-suite", handlers.Repo.Majors)
	mux.Get("/search-availability", handlers.Repo.Availability)
	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/make-reservation", handlers.Repo.Reservation)
	/**
	讓 client 可以透過 http request 取得專案內的檔案，
	例如專案內的圖檔(/static/images/house.jpg) 就可以使用這個 url 取得: http://localhost:8080/static/images/house.jpg
	*/
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
