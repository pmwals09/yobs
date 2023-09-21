package controllers

import (
	"html/template"
	"net/http"

	helpers "github.com/pmwals09/yobs/internal"
)

func HandlePing(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("pong").Parse("<p>Pong</p>")
	if err != nil {
		helpers.WriteError(w, err)
		return
	}
	t.Execute(w, nil)
}
