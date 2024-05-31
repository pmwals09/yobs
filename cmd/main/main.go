package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"

	"github.com/pmwals09/yobs/internal/config"
	"github.com/pmwals09/yobs/internal/controllers"
	"github.com/pmwals09/yobs/internal/db"
	"github.com/pmwals09/yobs/internal/models/contact"
	"github.com/pmwals09/yobs/internal/models/document"
	"github.com/pmwals09/yobs/internal/models/opportunity"
	"github.com/pmwals09/yobs/internal/models/session"
	"github.com/pmwals09/yobs/internal/models/status"
	"github.com/pmwals09/yobs/internal/models/user"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	logger := slog.New(slog.Default().Handler())

	config, err := config.New()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	sqlDb, err := db.InitDb(*config)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	opptyRepo := opportunity.OpportunityModel{DB: sqlDb}
	docRepo := document.DocumentModel{DB: sqlDb, Config: *config}
	userRepo := user.UserModel{DB: sqlDb}
	sessionRepo := session.SessionModel{DB: sqlDb}
	contactRepo := contact.ContactModel{DB: sqlDb}
	statusRepo := status.StatusModel{DB: sqlDb}

	r.Get("/", controllers.HandleGetLandingPage(logger))
	r.Get("/ping", controllers.HandlePing)
	r.Get("/sign-up", controllers.HandleGetSignUpPage(logger))
	r.Get("/login", controllers.HandleGetLoginPage(logger))
	r.Route("/user", func(r chi.Router) {
		r.Post("/register", controllers.HandleRegisterUser(&userRepo, logger))
		r.Post("/login", controllers.HandleLogInUser(&userRepo, &sessionRepo, logger))
		r.Get("/logout", controllers.HandleLogout(&sessionRepo, logger))
	})
	r.Mount("/", authenticatedRouter(&opptyRepo, &docRepo, &sessionRepo, &userRepo, &contactRepo, &statusRepo, logger))

	log.Fatal(http.ListenAndServe(":8080", r))
}

func authenticatedRouter(opptyRepo opportunity.Repository, docRepo document.Repository, sessionRepo session.Repository, userRepo user.Repository, contactRepo contact.Repository, statusRepo status.Repository, logger *slog.Logger) http.Handler {
	r := chi.NewRouter()
	r.Use(authOnly(sessionRepo, userRepo, logger))
	r.Get("/home", controllers.HandleGetHomepage(opptyRepo, logger))
	r.Route("/profile", func(r chi.Router) {
		r.Get("/", controllers.HandleGetProfilePage(logger))
		r.Get("/basic-profile-form", controllers.HandleGetBasicProfileForm(logger))
		r.Put("/update", controllers.HandleUpdateProfile(userRepo, logger))
	})
	r.Route("/opportunities", func(r chi.Router) {
		r.Post("/", controllers.HandlePostOppty(opptyRepo, logger))
		r.Route("/{opportunityId}", func(r chi.Router) {
			r.Get("/", controllers.HandleGetOpptyPage(opptyRepo, docRepo, logger))
			r.Get("/contact-modal", controllers.HandleContactModal(opptyRepo, logger))
			r.Get("/attachment-modal", controllers.HandleAttachmentModal(opptyRepo, docRepo, logger))
			r.Get("/status-modal", controllers.HandleStatusModal(opptyRepo, logger))
			r.Get("/edit-details", controllers.HandleGetEditDetailsForm(opptyRepo, logger))
			r.Post("/upload", controllers.HandleUploadToOppty(opptyRepo, docRepo, logger))
			r.Post("/attach-existing", controllers.HandleAddExistingToOppty(opptyRepo, docRepo, logger))
			r.Post("/new-contact", controllers.HandleAddNewContactToOppty(opptyRepo, contactRepo, logger))
			r.Post("/update-status", controllers.HandleUpdateStatus(opptyRepo, logger))
			r.Put("/edit", controllers.HandleUpdate(opptyRepo, docRepo, logger))
			r.Delete("/documents/{documentId}", controllers.HandleRemoveDocFromOppty(opptyRepo, docRepo, logger))
			r.Route("/statuses", func(r chi.Router) {
				r.Delete("/{statusID}", controllers.HandleDeleteStatus(statusRepo, logger))
				r.Route("/{statusID}", func(r chi.Router) {
					r.Put("/", controllers.HandleUpdateStatusItem(statusRepo, opptyRepo, logger))
					r.Get("/status-row-form", controllers.HandleStatusRowForm(statusRepo, logger))
				})
			})
			r.Route("/contacts", func(r chi.Router) {
				r.Route("/{contactId}", func(r chi.Router) {
					r.Get("/contact-row-form", controllers.HandleContactRowForm(contactRepo, logger))
					r.Put("/", controllers.HandleUpdateContact(contactRepo, logger))
					r.Delete("/", controllers.HandleDeleteContact(contactRepo, logger))
				})
			})
		})
	})
	return r
}

func authOnly(sessionRepo session.Repository, userRepo user.Repository, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookies := r.Cookies()
			var cookie *http.Cookie

			if len(cookies) == 0 {
				logger.Warn("Cookies are len 0")
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}

			for _, c := range cookies {
				if c.Name == "yobs" {
					if cookie == nil || c.Expires.After(cookie.Expires) {
						cookie = c
					}
				}
			}

			if cookie == nil {
				logger.Warn("Cookie is nil")
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}

			uuid, err := uuid.Parse(cookie.Value)
			if err != nil {
				logger.Error("Problem parsing UUID", "error", err)
				http.Redirect(w, r, "/", http.StatusFound)
				err := sessionRepo.DeleteSessionByUUID(uuid)
				if err != nil {
					logger.Error("Problem deleting session", "error", err)
				}
				return
			}

			session, err := sessionRepo.GetSessionByUUID(uuid)
			if err != nil {
				logger.Error("Problem getting session by UUID", "error", err)
				http.Redirect(w, r, "/", http.StatusFound)
				err := sessionRepo.DeleteSessionByUUID(uuid)
				if err != nil {
					logger.Error("Problem deleting session", "error", err)
				}
				return
			}

			now := time.Now()
			if now.After(session.Expiration) {
				err := sessionRepo.DeleteSessionByUUID(session.UUID)
				if err != nil {
					logger.Error("Problem deleting session", "error", err)
				}
				logger.Info("Redirecting because session is expired")
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}

			session.Expiration = time.Now().Add(time.Minute * 30)
			err = sessionRepo.UpdateSession(session)
			if err != nil {
				logger.Error("Problem updating session", "error", err)
				err := sessionRepo.DeleteSessionByUUID(session.UUID)
				if err != nil {
					logger.Error("Problem deleting session", "error", err)
				}
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
			cookie.Expires = session.Expiration
			http.SetCookie(w, cookie)

			u, err := userRepo.GetUserById(session.UserID)
			if err != nil {
				logger.Error("Problem getting user", "error", err)
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
			logger.Info("Successfully navigated session mgmt")
			ctx := context.WithValue(r.Context(), user.UserCtxKey, u)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
