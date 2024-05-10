package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"github.com/pmwals09/yobs/internal/controllers"
	"github.com/pmwals09/yobs/internal/db"
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

	sqlDb, err := db.InitDb()
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
		r.Get("/logout", controllers.HandleLogout(&sessionRepo))
	})
	r.Mount("/", authenticatedRouter(&opptyRepo, &docRepo, &sessionRepo, &userRepo))

	log.Fatal(http.ListenAndServe(":8080", r))
}

func authenticatedRouter(opptyRepo opportunity.Repository, docRepo document.Repository, sessionRepo session.Repository, userRepo user.Repository) http.Handler {
	r := chi.NewRouter()
	r.Use(authOnly(sessionRepo, userRepo))
	r.Get("/home", controllers.HandleGetHomepage(opptyRepo))
	r.Route("/profile", func(r chi.Router) {
		r.Get("/", controllers.HandleGetProfilePage())
	})
	r.Route("/opportunities", func(r chi.Router) {
		r.Post("/", controllers.HandlePostOppty(opptyRepo))
		r.Route("/{opportunityId}", func(r chi.Router) {
			r.Get("/", controllers.HandleGetOpptyPage(opptyRepo, docRepo))
			r.Get("/edit", controllers.HandleEditOpptyPage(opptyRepo, docRepo))
			r.Post("/upload", controllers.HandleUploadToOppty(opptyRepo, docRepo))
			r.Post("/attach-existing", controllers.HandleAddExistingToOppty(opptyRepo, docRepo))
			r.Get("/contact-modal", controllers.HandleContactModal(opptyRepo))
		})
	})
	return r
}

func authOnly(sessionRepo session.Repository, userRepo user.Repository) func(http.Handler) http.Handler {
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
				sessionRepo.DeleteSessionByUUID(session.UUID)
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}

			session.Expiration = time.Now().Add(time.Minute * 30)
			err = sessionRepo.UpdateSession(session)
			if err != nil {
			}
			cookie.Expires = session.Expiration
			http.SetCookie(w, cookie)

			u, err := userRepo.GetUserById(session.UserID)
			if err != nil {
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
			ctx := context.WithValue(r.Context(), "user", u)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
