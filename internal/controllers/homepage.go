package controllers

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	helpers "github.com/pmwals09/yobs/internal"
	"github.com/pmwals09/yobs/internal/models/opportunity"
	"github.com/pmwals09/yobs/internal/models/session"
	"github.com/pmwals09/yobs/internal/models/user"
	homepage "github.com/pmwals09/yobs/web/home"
	indexpage "github.com/pmwals09/yobs/web/index"
	loginpage "github.com/pmwals09/yobs/web/login"
	signuppage "github.com/pmwals09/yobs/web/signup"
)

func HandleGetHomepage(opptyRepo opportunity.Repository, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if user, ok := r.Context().Value(user.UserCtxKey).(*user.User); ok && user != nil {
			opptys, err := opptyRepo.GetAllOpportunities(user)
			if err != nil {
				logger.Error("Can't get all user opptys", "error", err)
				helpers.WriteError(w, err)
				return
			}

			homepage.HomePage(user, opptys, helpers.FormData{}).Render(r.Context(), w)
		} else {
			msg := "no user available"
			logger.Error(msg)
			helpers.WriteError(w, errors.New(msg))
			return
		}
	}
}

func HandleGetLandingPage(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if user, ok := r.Context().Value(user.UserCtxKey).(*user.User); ok {
			indexpage.IndexPage(user).Render(r.Context(), w)
		} else {
			logger.Error("No user available")
			indexpage.IndexPage(nil).Render(r.Context(), w)
		}
	}
}

func HandleGetLoginPage(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if user, ok := r.Context().Value(user.UserCtxKey).(*user.User); ok {
			loginpage.LoginPage(user, helpers.FormData{}).Render(r.Context(), w)
		} else {
			logger.Error("No user available")
			loginpage.LoginPage(nil, helpers.FormData{}).Render(r.Context(), w)
		}
	}
}

func HandleGetSignUpPage(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if user, ok := r.Context().Value(user.UserCtxKey).(*user.User); ok {
			signuppage.SignupPage(user, helpers.FormData{}).Render(r.Context(), w)
		} else {
			logger.Error("No user available")
			signuppage.SignupPage(nil, helpers.FormData{}).Render(r.Context(), w)
		}
	}
}

func HandleLogout(sessionRepo session.Repository, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("yobs")
		if err != nil {
			logger.Error("Problem getting cookie", "error", err)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		if cookie == nil {
			logger.Error("No cookie")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		uuid, err := uuid.Parse(cookie.Value)
		if err != nil {
			logger.Error("Problem parsing cookie UUID", "error", err)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		cookie.Expires = time.Time{}
		http.SetCookie(w, cookie)

		sessionRepo.DeleteSessionByUUID(uuid)

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
