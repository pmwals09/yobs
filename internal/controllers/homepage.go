package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	helpers "github.com/pmwals09/yobs/internal"
	"github.com/pmwals09/yobs/internal/models/opportunity"
	"github.com/pmwals09/yobs/internal/models/session"
	"github.com/pmwals09/yobs/internal/models/user"
	"github.com/pmwals09/yobs/web/templates"
)

func HandleGetHomepage(opptyRepo opportunity.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if user, ok := r.Context().Value("user").(*user.User); ok {
			if user == nil {
				helpers.WriteError(w, errors.New("No user available"))
				return
			}

			opptys, err := opptyRepo.GetAllOpportunities(user)
			if err != nil {
				helpers.WriteError(w, err)
				return
			}

			templates.HomePage(user, opptys, helpers.FormData{}).Render(r.Context(), w)
		} else {
			helpers.WriteError(w, errors.New("No user available"))
			return
		}
	}
}

func HandleGetLandingPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if user, ok := r.Context().Value("user").(*user.User); ok {
			templates.IndexPage(user).Render(r.Context(), w)
		} else {
			templates.IndexPage(nil).Render(r.Context(), w)
		}
	}
}

func HandleGetLoginPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if user, ok := r.Context().Value("user").(*user.User); ok {
			templates.LoginPage(user, helpers.FormData{}).Render(r.Context(), w)
		} else {
			templates.LoginPage(nil, helpers.FormData{}).Render(r.Context(), w)
		}
	}
}

func HandleGetSignUpPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if user, ok := r.Context().Value("user").(*user.User); ok {
			templates.SignupPage(user, helpers.FormData{}).Render(r.Context(), w)
		} else {
			templates.SignupPage(nil, helpers.FormData{}).Render(r.Context(), w)
		}
	}
}

func HandleLogout(sessionRepo session.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("yobs")
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		if cookie == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		uuid, err := uuid.Parse(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		cookie.Expires = time.Time{}
		http.SetCookie(w, cookie)

		err = sessionRepo.DeleteSessionByUUID(uuid)

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
