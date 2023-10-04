package controllers

import (
	"errors"
	"net/http"

	helpers "github.com/pmwals09/yobs/internal"
	"github.com/pmwals09/yobs/internal/models/opportunity"
	"github.com/pmwals09/yobs/internal/models/user"
	templates "github.com/pmwals09/yobs/web/template"
)

func HandleGetHomepage(opptyRepo opportunity.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
    user := r.Context().Value("user").(*user.User)
    if user == nil {
			helpers.WriteError(w, errors.New("No user available"))
			return
    }

    opptys, err := opptyRepo.GetAllOpportunities(user)
    if err != nil {
			helpers.WriteError(w, err)
			return
    }

    templates.HomePage(opptys, helpers.FormData{}).Render(r.Context(), w)
	}
}

func HandleGetLandingPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
    templates.IndexPage().Render(r.Context(), w)
	}
}

func HandleGetLoginPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
    templates.LoginPage(helpers.FormData{}).Render(r.Context(), w)
	}
}

func HandleGetSignUpPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
    templates.SignupPage(helpers.FormData{}).Render(r.Context(), w)
	}
}

func HandleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: set new cookie with a past expiration
	}
}
