package controllers

import (
	"errors"
	"html/template"
	"net/http"
	"os"

	helpers "github.com/pmwals09/yobs/internal"
	"github.com/pmwals09/yobs/internal/models/opportunity"
	"github.com/pmwals09/yobs/internal/models/user"
)

func HandleGetHomepage(opptyRepo opportunity.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wd, err := os.Getwd()
		if err != nil {
			helpers.WriteError(w, err)
			return
		}
		t, err := template.
      New("base").
      Funcs(helpers.GetListFuncMap()).
      ParseFiles(
        wd+"/web/template/opportunity-form-partial.html",
        wd+"/web/template/opportunity-list-partial.html",
        wd+"/web/template/home-page.html",
        wd+"/web/template/base.html",
      )
		if err != nil {
			helpers.WriteError(w, err)
			return
		}

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

    t.ExecuteTemplate(w, "base", map[string][]opportunity.Opportunity { "Opportunities": opptys })
	}
}

func HandleGetLandingPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wd, err := os.Getwd()
		if err != nil {
			helpers.WriteError(w, err)
			return
		}
		t, err := template.
      ParseFiles(
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

func HandleGetLoginPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
    wd, err := os.Getwd()
    if err != nil {
      helpers.WriteError(w, err)
      return
    }
    t, err := template.
      ParseFiles(
        wd+"/web/template/login-user-form-partial.html",
        wd+"/web/template/login-page.html",
        wd+"/web/template/base.html",
      )
    if err != nil {
      helpers.WriteError(w, err)
      return
    }
    t.ExecuteTemplate(w, "base", nil)
	}
}

func HandleGetSignUpPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
    wd, err := os.Getwd()
    if err != nil {
      helpers.WriteError(w, err)
      return
    }
    t, err := template.ParseFiles(
      wd+"/web/template/register-user-form-partial.html",
      wd+"/web/template/signup-page.html",
      wd+"/web/template/base.html",
    )
    if err != nil {
      helpers.WriteError(w, err)
      return
    }
    t.ExecuteTemplate(w, "base", nil)
	}
}

func HandleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: set new cookie with a past expiration
	}
}
