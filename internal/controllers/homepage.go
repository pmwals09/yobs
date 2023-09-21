package controllers

import (
	"html/template"
	"net/http"
	"os"

	helpers "github.com/pmwals09/yobs/internal"
)

func HandleGetHomepage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wd, err := os.Getwd()
		if err != nil {
			helpers.WriteError(w, err)
			return
		}
		t, err := template.ParseFiles(
			wd+"/web/template/opportunity-form-partial.html",
			wd+"/web/template/index-page.html",
			wd+"/web/template/base.html",
		)
		if err != nil {
			helpers.WriteError(w, err)
			return
		}
		t.ExecuteTemplate(w, "base", nil)
	}
}
