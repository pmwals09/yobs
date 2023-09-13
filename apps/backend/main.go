package backend

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	db "github.com/pmwals09/yobs/apps/backend/db"
	"github.com/pmwals09/yobs/apps/backend/opportunity"
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

	opptyRepo := opportunity.GormRepository{DB: sqlDb}
	r.Get("/", handleGetHomepage())
	r.Get("/ping", handlePing)
	r.Route("/opportunities", func(r chi.Router) {
		r.Post("/", handlePostOppty(opptyRepo))
		// r.Get("/", handleGetAllOppty(opptyRepo))
		r.Get("/active", handleGetActiveOpptys(opptyRepo))
		r.Route("/{opportunityId}", func(r chi.Router) {
			r.Get("/", handleGetOppty(opptyRepo))
			r.Post("/upload", handleUploadToOppty(opptyRepo))
			// r.Get("/edit", handleEditOppty(opptyRepo))
			// r.Put("/", handleUpdateOppty(opptyRepo))
			// r.Delete("/", handleDeleteOppty(opptyRepo))
		})
	})

	http.ListenAndServe(":8080", r)
}

func handleGetHomepage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wd, err := os.Getwd()
		if err != nil {
			writeError(w, err)
			return
		}
		t, err := template.ParseFiles(
			wd+"/apps/clients/web/templates/opportunity-form-partial.html",
			wd+"/apps/clients/web/templates/index-page.html",
			wd+"/apps/clients/web/templates/base.html",
		)
		if err != nil {
			writeError(w, err)
			return
		}
		t.ExecuteTemplate(w, "base", nil)
	}
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("pong").Parse("<p>Pong</p>")
	if err != nil {
		writeError(w, err)
		return
	}
	t.Execute(w, nil)
}

func handlePostOppty(repo opportunity.GormRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		newOpportunity := newOpportunityFromRequest(r)

		wd, wdErr := os.Getwd()
		if wdErr != nil {
			writeError(w, wdErr)
			return
		}

		if _, err := repo.CreateOpportunity(newOpportunity); err != nil {
			handleCreateOpptyError(w, wd, err)
			return
		}

		r.Header.Add("HX-Retarget", "#main-content")
		tmpl, templateErr := template.New("opportunity-list").Funcs(getListFuncMap()).ParseFiles(
			wd + "/apps/clients/web/templates/opportunity-list-partial.html",
		)

		if templateErr != nil {
			writeError(w, templateErr)
			return
		}

		opportunities, opptyErr := repo.GetAllOpportunities()
		if opptyErr != nil {
			writeError(w, opptyErr)
			return
		}
		tmpl.ExecuteTemplate(w, "opportunity-list", opportunities)
	}
}

func handleCreateOpptyError(w http.ResponseWriter, wd string, err error) {
	t, templateError := template.ParseFiles(
		wd + "/apps/clients/web/templates/opportunity-form-partial.html",
	)

	if templateError != nil {
		writeError(w, templateError)
		return
	}

	data := map[string]string{
		"Errors": fmt.Sprintf(
			"<ul><li>An error occurred creating the opportunity: %s</li></ul>",
			err.Error(),
		),
	}
	t.ExecuteTemplate(
		w,
		"opportunity-form",
		data,
	)
}

// TODO: Would love to get rid of the overly verbose task objects when
// retrieving these
// func handleGetAllOppty(repo opportunity.GormRepository) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if opptys, err := repo.GetAllOpportunities(); err != nil {
// 			w.Write([]byte(fmt.Sprintf("<p>An error has occurred: %s", err.Error())))
// 			return
// 		} else {
// 			res := []render.Renderer{}
// 			for _, o := range opptys {
// 				res = append(res, o)
// 			}
// 			render.RenderList(w, r, res)
// 		}
// 	}
// }

func handleGetActiveOpptys(repo opportunity.GormRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opptys, err := repo.GetAllOpportunities()
		if err != nil {
			writeError(w, err)
			return
		}
		wd, err := os.Getwd()
		if err != nil {
			writeError(w, err)
			return
		}
		tmpl, err := template.New("opportunity-list").Funcs(getListFuncMap()).ParseFiles(
			wd + "/apps/clients/web/templates/opportunity-list-partial.html",
		)
		if err != nil {
			writeError(w, err)
			return
		}
		tmpl.ExecuteTemplate(w, "opportunity-list", opptys)
	}
}

func handleGetOppty(repo opportunity.GormRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wd, err := os.Getwd()
		if err != nil {
			writeError(w, err)
			return
		}
		t, err := template.New("base").Funcs(getListFuncMap()).ParseFiles(
			wd+"/apps/clients/web/templates/opportunity-details-page.html",
			wd+"/apps/clients/web/templates/base.html",
		)
		if err != nil {
			writeError(w, err)
			return
		}

		idParam := chi.URLParam(r, "opportunityId")
		id, err := strconv.ParseUint(idParam, 10, 64) // Sqlite id's are 64-bit int
		if err != nil {
			writeError(w, err)
			return
		}

		opp, err := repo.GetOpportuntyById(uint(id))
		if err != nil {
			writeError(w, err)
			return
		}

		t.ExecuteTemplate(w, "base", opp)
	}
}

func handleUploadToOppty(repo opportunity.GormRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(10 << 20)
		file, handler, err := r.FormFile("attachment-file")
		if err != nil {
			writeError(w, err)
			return
		}
		defer file.Close()

		fmt.Println(handler.Header)
	}
}

// TODO: How to get a partial update? Would be nice if this becomes a
// large model
//   - https://gorm.io/docs/update.html#Update-Selected-Fields
//   - Would need additional logic around what keys exist from user
//     provided JSON to handle in the Bind/marshalling/DTO model
// func handleUpdateOppty(repo opportunity.GormRepository) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if data, err := bindRequest(w, r); err != nil {
// 			return
// 		} else {
// 			idParam := chi.URLParam(r, "opportunityId")
// 			if id, err := strconv.ParseUint(idParam, 10, 64); err != nil {
// 				render.Render(w, r, ErrRender(r, err, 401))
// 				return
// 			} else {
// 				if oppty, err := handleGetOpptyById(repo, id, w, r); err != nil {
// 					return
// 				} else {
// 					updatedOppty := opportunity.NewOpportunity(data.Description, data.URL)
// 					updatedOppty.ID = oppty.ID
//
// 					if oppty, err := repo.UpdateOpporunity(updatedOppty); err != nil {
// 						render.Render(w, r, ErrRender(r, err, 500))
// 						return
// 					} else {
// 						res := OpportunityResponse{
// 							oppty,
// 						}
// 						render.Render(w, r, &res)
// 					}
// 				}
// 			}
// 		}
//
// 	}
// }

// func handleDeleteOppty(repo opportunity.GormRepository) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		idParam := chi.URLParam(r, "opportunityId")
// 		if id, err := strconv.ParseUint(idParam, 10, 64); err != nil {
// 			render.Render(w, r, ErrRender(r, err, 401))
// 			return
// 		} else {
// 			if err := repo.DeleteOpportunity(uint(id)); err != nil {
// 				render.Render(w, r, ErrRender(r, err, 500))
// 			} else {
// 				render.NoContent(w, r)
// 			}
// 		}
// 	}
// }

// func handleGetOpptyById(
// 	repo opportunity.GormRepository,
// 	id uint64,
// 	w http.ResponseWriter,
// 	r *http.Request,
// ) (opportunity.Opportunity, error) {
// 	if oppty, err := repo.GetOpportuntyById(uint(id)); err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			render.Render(w, r, ErrRender(r, err, 404))
// 			return oppty, err
// 		}
// 		render.Render(w, r, ErrRender(r, err, 500))
// 		return oppty, err
// 	} else {
// 		return oppty, nil
// 	}
// }

func newOpportunityFromRequest(r *http.Request) *opportunity.Opportunity {
	name := r.PostForm.Get("opportunity-name")
	description := r.PostForm.Get("opportunity-description")
	url := r.PostForm.Get("opportunity-url")
	date := r.PostForm.Get("opportunity-date")
	role := r.PostForm.Get("opportunity-role")
	o := opportunity.NewOpportunity().WithCompanyName(name).WithRole(role).WithDescription(description).WithURL(url).WithApplicationDateString(date)
	if o.ApplicationDate.IsZero() {
		o.Status = opportunity.None
	} else {
		o.Status = opportunity.Applied
	}
	return o
}

func writeError(w http.ResponseWriter, err error) {
	w.Write([]byte(fmt.Sprintf(
		"<p>An error has occurred: %s</p>",
		err.Error(),
	)))
}

func getListFuncMap() template.FuncMap {
	return template.FuncMap{
		"FormatApplicationDate": func(t time.Time) string {
			if t.IsZero() {
				return ""
			}
			return t.Format("2006-01-02")
		},
	}
}
