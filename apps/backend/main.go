package backend

import (
	// "html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	db "github.com/pmwals09/yobs/apps/backend/db"
)

func Run() {
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

	r.Mount("/api", ApiRouter(sqlDb))
	r.Mount("/", frontEndRouter())

	http.ListenAndServe(":8080", r)
}
