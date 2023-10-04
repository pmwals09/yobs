package controllers

import (
	"net/http"

	helpers "github.com/pmwals09/yobs/internal"
	"github.com/pmwals09/yobs/internal/models/document"
	"github.com/pmwals09/yobs/internal/models/user"
	templates "github.com/pmwals09/yobs/web/template"
)

func HandleGetProfilePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
    u := r.Context().Value("user").(*user.User)
		pa := helpers.ProfileArgs {
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
    templates.ProfilePage(pa).Render(r.Context(), w)
	}
}
