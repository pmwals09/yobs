package controllers

import (
	"html/template"
	"net/http"
	"os"

	helpers "github.com/pmwals09/yobs/internal"
	"github.com/pmwals09/yobs/internal/models/document"
)

type ProfileArgs struct {
	Username string
	Email string
	Resume document.Document
}

func HandleGetProfilePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wd, err := os.Getwd()
		if err != nil {
			helpers.WriteError(w, err)
			return
		}

		t, err := template.ParseFiles(
			wd+"/web/template/profile-page.html",
			wd+"/web/template/base.html",
		)
		if err != nil {
			helpers.WriteError(w, err)
			return
		}

		// TODO: get the logged-in user
		pa := ProfileArgs {
			Username: "pmwals09",
			Email: "pmwals09@gmail.com",
			Resume: document.Document{
				FileName: "Anonymous Resume",
				Title: "My fancy resume",
				Type: "Resume",
				URL: "http://www.google.com",
			},
		}
		t.ExecuteTemplate(w, "base", pa)
	}
}
