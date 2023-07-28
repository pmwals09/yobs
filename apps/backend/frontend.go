package backend

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func frontEndRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", serveHome)
	return r
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	// NOTE: this path needs to be from the root of the project, like a webroot
	t, _ := template.ParseFiles("apps/clients/web/templates/index.html")
	t.Execute(w, struct{}{})
}
