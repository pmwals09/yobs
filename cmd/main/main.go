package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	helpers "github.com/pmwals09/yobs/internal"
	"github.com/pmwals09/yobs/internal/models/document"
	"github.com/pmwals09/yobs/internal/models/opportunity"
	"github.com/pmwals09/yobs/internal/controllers"
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

	r.Get("/", controllers.HandleGetHomepage())
	r.Get("/ping", controllers.HandlePing)
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

	log.Fatal(http.ListenAndServe(":8081", r))
}
