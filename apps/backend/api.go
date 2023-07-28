package backend

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
)

func ApiRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Get("/ping", handlePing)
	return r
}

type PongRes struct {
	Msg     string `json:"msg"`
	Elapsed int64  `json:"elapsed"`
}

func (pr PongRes) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	res := PongRes{"pong", 0}
	render.Render(w, r, res)
}

// create opportunity
// read opportunity
// read all opportunities
