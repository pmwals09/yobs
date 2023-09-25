package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	helpers "github.com/pmwals09/yobs/internal"
	"github.com/pmwals09/yobs/internal/controllers"
	"github.com/pmwals09/yobs/internal/models/document"
	"github.com/pmwals09/yobs/internal/models/opportunity"
	"github.com/pmwals09/yobs/internal/models/session"
	"github.com/pmwals09/yobs/internal/models/user"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	sqlDb, err := helpers.InitDb()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	opptyRepo := opportunity.OpportunityModel{DB: sqlDb}
	docRepo := document.DocumentModel{DB: sqlDb}
	userRepo := user.UserModel{DB: sqlDb}
	sessionRepo := session.SessionModel{DB: sqlDb}

	r.Get("/", controllers.HandleGetLandingPage())
	r.Get("/ping", controllers.HandlePing)
	r.Get("/sign-up", controllers.HandleGetSignUpPage())
	r.Get("/login", controllers.HandleGetLoginPage())
	r.Route("/user", func(r chi.Router) {
		r.Post("/register", controllers.HandleRegisterUser(&userRepo))
		r.Post("/login", controllers.HandleLogInUser(&userRepo, &sessionRepo))
	})
	// TODO: All these routes should be behind Auth - only a valid user can see them
	// so we should add some middleware for these routes that confirms a user is:
	// - logged in
	//   - use previously set JWT
	// - allowed to see this content (i.e., admin vs. user)
	r.Mount("/", authenticatedRouter(&opptyRepo, &docRepo, &sessionRepo))

	log.Fatal(http.ListenAndServe(":8082", r))
}

func authenticatedRouter(opptyRepo opportunity.Repository, docRepo document.Repository, sessionRepo session.Repository) http.Handler {
	r := chi.NewRouter()
	r.Use(AuthOnly(sessionRepo))
	r.Get("/home", controllers.HandleGetHomepage())
	r.Route("/profile", func(r chi.Router) {
		r.Get("/", controllers.HandleGetProfilePage())
	})
	r.Route("/opportunities", func(r chi.Router) {
		r.Post("/", controllers.HandlePostOppty(opptyRepo))
		r.Get("/active", controllers.HandleGetActiveOpptys(opptyRepo))
		r.Route("/{opportunityId}", func(r chi.Router) {
			r.Get("/", controllers.HandleGetOppty(opptyRepo))
			r.Post("/upload", controllers.HandleUploadToOppty(opptyRepo, docRepo))
		})
	})
	return r
}

func AuthOnly(sessionRepo session.Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookies := r.Cookies()
			var cookie *http.Cookie

			if len(cookies) == 0 {
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}

			for _, c := range cookies {
				if c.Name == "yobs" {
					cookie = c
					break
				}
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

			session, err := sessionRepo.GetSessionByUUID(uuid)
			if err != nil {
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}

			now := time.Now()
			if now.After(session.Expiration) {
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}

			session.Expiration = time.Now().Add(time.Minute * 30)
			cookie.Expires = session.Expiration
			http.SetCookie(w, cookie)

			next.ServeHTTP(w, r)
		})
	}
}
