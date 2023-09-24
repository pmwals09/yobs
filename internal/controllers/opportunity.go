package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	helpers "github.com/pmwals09/yobs/internal"
	"github.com/pmwals09/yobs/internal/models/document"
	"github.com/pmwals09/yobs/internal/models/opportunity"
)

func HandlePostOppty(repo opportunity.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newOpportunity := newOpportunityFromRequest(r)

		wd, wdErr := os.Getwd()
		if wdErr != nil {
			helpers.WriteError(w, wdErr)
			return
		}

		if err := repo.CreateOpportunity(newOpportunity); err != nil {
			handleCreateOpptyError(w, wd, err)
			return
		}

		r.Header.Add("HX-Retarget", "#main-content")
		tmpl, templateErr := template.New("opportunity-list").Funcs(getListFuncMap()).ParseFiles(
			wd + "/web/template/opportunity-list-partial.html",
		)

		if templateErr != nil {
			helpers.WriteError(w, templateErr)
			return
		}

		opportunities, opptyErr := repo.GetAllOpportunities()
		if opptyErr != nil {
			helpers.WriteError(w, opptyErr)
			return
		}
		tmpl.ExecuteTemplate(w, "opportunity-list", opportunities)
	}
}

func handleCreateOpptyError(w http.ResponseWriter, wd string, err error) {
	t, templateError := template.ParseFiles(
		wd + "/web/template/opportunity-form-partial.html",
	)

	if templateError != nil {
		helpers.WriteError(w, templateError)
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

func HandleGetActiveOpptys(repo opportunity.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opptys, err := repo.GetAllOpportunities()
		if err != nil {
			helpers.WriteError(w, err)
			return
		}
		wd, err := os.Getwd()
		if err != nil {
			helpers.WriteError(w, err)
			return
		}
		tmpl, err := template.New("opportunity-list").Funcs(getListFuncMap()).ParseFiles(
			wd + "/web/template/opportunity-list-partial.html",
		)
		if err != nil {
			helpers.WriteError(w, err)
			return
		}
		tmpl.ExecuteTemplate(w, "opportunity-list", opptys)
	}
}

type OpptyDetails struct {
	Oppty     opportunity.Opportunity
	Documents []document.Document
}

func HandleGetOppty(repo opportunity.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wd, err := os.Getwd()
		if err != nil {
			helpers.WriteError(w, err)
			return
		}
		t, err := template.New("base").Funcs(getListFuncMap()).ParseFiles(
			wd+"/web/template/attachments-section-partial.html",
			wd+"/web/template/opportunity-details-page.html",
			wd+"/web/template/base.html",
		)
		if err != nil {
			helpers.WriteError(w, err)
			return
		}

		od := OpptyDetails{}

		idParam := chi.URLParam(r, "opportunityId")
		id, err := strconv.ParseUint(idParam, 10, 64) // Sqlite id's are 64-bit int
		if err != nil {
			helpers.WriteError(w, err)
			return
		}

		opp, err := repo.GetOpportuntyById(uint(id))
		if err != nil {
			helpers.WriteError(w, err)
			return
		}

		od.Oppty = *opp
		docs, err := repo.GetAllDocuments(opp.ID)
		if err != nil {
			helpers.WriteError(w, err)
			return
		}

		for i := range docs {
			_, err := docs[i].GetPresignedDownloadUrl()
			if err != nil {
				helpers.WriteError(w, err)
				return
			}
		}

		od.Documents = docs

		t.ExecuteTemplate(w, "base", od)
	}
}

func HandleUploadToOppty(
	oppRepo opportunity.Repository,
	docRepo document.Repository,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(10 << 20)
		file, handler, err := r.FormFile("attachment-file")
		defer file.Close()
		if err != nil {
			helpers.WriteError(w, err)
			return
		}

		// 1. Upload the file to its destination
		var docType document.DocumentType
		selectedType := r.FormValue("attachment-type")
		switch selectedType {
		case "Resume":
			docType = document.Resume
		case "CoverLetter":
			docType = document.CoverLetter
		case "Communication":
			docType = document.Communication
		}
		d := document.New(handler, docType)

		docTitle := r.FormValue("attachment-name")
		if docTitle != "" {
			d.WithTitle(docTitle)
		}

		err = d.Upload(file)
		if err != nil {
			helpers.WriteError(w, err)
			return
		}

		// 2. Insert a document entry into the db
		if err := docRepo.CreateDocument(d); err != nil {
			helpers.WriteError(w, err)
			return
		}

		// 3. Associate the document with this oppty
		idParam := chi.URLParam(r, "opportunityId")
		id, err := strconv.ParseUint(idParam, 10, 64) // Sqlite id's are 64-bit int
		if err != nil {
			helpers.WriteError(w, err)
			return
		}
		if err = oppRepo.AddDocument(uint(id), d.ID); err != nil {
			helpers.WriteError(w, err)
			return
		}

		// 4. What to return? And where?
		wd, err := os.Getwd()
		if err != nil {
			helpers.WriteError(w, err)
			return
		}
		t, err := template.New("attachments-section").Funcs(getListFuncMap()).ParseFiles(
			wd + "/web/template/attachments-section-partial.html",
		)
		o, err := oppRepo.GetOpportuntyById(uint(id))
		if err != nil {
			helpers.WriteError(w, err)
			return
		}

		docs, err := oppRepo.GetAllDocuments(uint(id))
		if err != nil {
			helpers.WriteError(w, err)
			return
		}

		for i := range docs {
			_, err := docs[i].GetPresignedDownloadUrl()
			if err != nil {
				fmt.Println(err.Error())
				helpers.WriteError(w, err)
				return
			}
		}

		od := OpptyDetails{
			Oppty:     *o,
			Documents: docs,
		}
		t.ExecuteTemplate(w, "attachments-section", od)
	}
}

func newOpportunityFromRequest(r *http.Request) *opportunity.Opportunity {
	r.ParseForm()
	name := r.PostForm.Get("opportunity-name")
	description := r.PostForm.Get("opportunity-description")
	url := r.PostForm.Get("opportunity-url")
	date := r.PostForm.Get("opportunity-date")
	role := r.PostForm.Get("opportunity-role")
	o := opportunity.New().WithCompanyName(name).WithRole(role).WithDescription(description).WithURL(url).WithApplicationDateString(date)
	if o.ApplicationDate.IsZero() {
		o.Status = opportunity.None
	} else {
		o.Status = opportunity.Applied
	}
	return o
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
