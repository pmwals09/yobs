package controllers

import (
	"html/template"
	"net/http"
	"os"

	helpers "github.com/pmwals09/yobs/internal"
	"github.com/pmwals09/yobs/internal/models/document"
	"github.com/pmwals09/yobs/internal/models/user"
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

    u := r.Context().Value("user").(*user.User)
		pa := ProfileArgs {
			Username: u.Username,
			Email: u.Email,
      // TODO: Make preferred resume a field on the user
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
