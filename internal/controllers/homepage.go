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
			wd+"/web/template/home-page.html",
			wd+"/web/template/base.html",
		)
		if err != nil {
			helpers.WriteError(w, err)
			return
		}
		t.ExecuteTemplate(w, "base", nil)
	}
}

func HandleGetLandingPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wd, err := os.Getwd()
		if err != nil {
			helpers.WriteError(w, err)
			return
		}
		t, err := template.ParseFiles(
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
		// TODO: Get the user row from the database
		// TODO: Compare hash and password
		// TODO: set jwt token cookie for additional requests
	}
}

func HandleGetSignUpPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// https://www.youtube.com/watch?v=d4Y2DkKbxM0&ab_channel=freeCodeCamp.org
		// TODO: Create a user object
		// TODO: hash the password
		// TODO: save the user object in the db with the hashed password
    wd, err := os.Getwd()
    if err != nil {
      helpers.WriteError(w, err)
      return
    }
    t, err := template.ParseFiles(
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
