package handlers

import (
	"github.com/colingood026/bookings/pkg/config"
	"github.com/colingood026/bookings/pkg/models"
	"github.com/colingood026/bookings/pkg/render"
	"net/http"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (rp *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIp := r.RemoteAddr
	rp.App.Session.Put(r.Context(), "remote_ip", remoteIp)
	render.RenderTemplateFromCacheV2(w, "home.page.gohtml", &models.TemplateData{})
}

func (rp *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	stringMap["remote_ip"] = rp.App.Session.GetString(r.Context(), "remote_ip")
	tp := models.TemplateData{
		StringMap: stringMap,
	}

	render.RenderTemplateFromCacheV2(w, "about.page.gohtml", &tp)
}
