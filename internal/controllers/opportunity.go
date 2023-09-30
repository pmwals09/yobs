package controllers

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	helpers "github.com/pmwals09/yobs/internal"
	"github.com/pmwals09/yobs/internal/models/document"
	"github.com/pmwals09/yobs/internal/models/opportunity"
	"github.com/pmwals09/yobs/internal/models/user"
)

func HandlePostOppty(repo opportunity.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newOpportunity, err := newOpportunityFromRequest(r)
    if err != nil {
      helpers.WriteError(w, err)
    }

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
		tmpl, templateErr := template.
      New("opportunity-list").
      Funcs(helpers.GetListFuncMap()).
      ParseFiles(
        wd + "/web/template/opportunity-list-partial.html",
      )

		if templateErr != nil {
			helpers.WriteError(w, templateErr)
			return
		}

    user := r.Context().Value("user").(*user.User)
    if user == nil {
      helpers.WriteError(w, errors.New("No user available"))
    }
		opportunities, opptyErr := repo.GetAllOpportunities(user)
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
    user := r.Context().Value("user").(*user.User)
    fmt.Println()
    fmt.Println("USER - GET", user)
    fmt.Println()
    if user == nil {
      helpers.WriteError(w, errors.New("No user available"))
    }
		opptys, err := repo.GetAllOpportunities(user)
		if err != nil {
			helpers.WriteError(w, err)
			return
		}
		wd, err := os.Getwd()
		if err != nil {
			helpers.WriteError(w, err)
			return
		}
		tmpl, err := template.
      New("opportunity-list").
      Funcs(helpers.GetListFuncMap()).
      ParseFiles(
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

// TODO: Should only do so for the current user
func HandleGetOppty(repo opportunity.Repository) http.HandlerFunc {
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

    user := r.Context().Value("user").(*user.User)
    if user == nil {
      helpers.WriteError(w, errors.New("No user available"))
    }
		opp, err := repo.GetOpportuntyById(uint(id), user)
		if err != nil {
			helpers.WriteError(w, err)
			return
		}

		od.Oppty = *opp
		docs, err := repo.GetAllDocuments(opp)
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

// TODO: ASsociate with user
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

    user := r.Context().Value("user").(*user.User)
    if user == nil {
      helpers.WriteError(w, errors.New("No user available"))
      return
    }
    d.WithUser(user)

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
    oppty, err := oppRepo.GetOpportuntyById(uint(id), user)
    if err != nil {
      helpers.WriteError(w, err)
      return
    }
		if err = oppRepo.AddDocument(oppty, d); err != nil {
			helpers.WriteError(w, err)
			return
		}

		// 4. What to return? And where?
		wd, err := os.Getwd()
		if err != nil {
			helpers.WriteError(w, err)
			return
		}
		t, err := template.
      New("attachments-section").
      Funcs(helpers.GetListFuncMap()).
      ParseFiles(
        wd + "/web/template/attachments-section-partial.html",
      )

		docs, err := oppRepo.GetAllDocuments(oppty)
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
			Oppty:     *oppty,
			Documents: docs,
		}
		t.ExecuteTemplate(w, "attachments-section", od)
	}
}

func newOpportunityFromRequest(r *http.Request) (*opportunity.Opportunity, error) {
  o := opportunity.New()

  err := r.ParseForm()
  if err != nil {
    return o, err
  }
	name := r.PostForm.Get("opportunity-name")
	description := r.PostForm.Get("opportunity-description")
	url := r.PostForm.Get("opportunity-url")
	date := r.PostForm.Get("opportunity-date")
	role := r.PostForm.Get("opportunity-role")
	o.WithCompanyName(name).
    WithRole(role).
    WithDescription(description).
    WithURL(url).
    WithApplicationDateString(date)

	if o.ApplicationDate.IsZero() {
		o.Status = opportunity.None
	} else {
		o.Status = opportunity.Applied
	}

  user := r.Context().Value("user").(*user.User)
  if user == nil {
    return o, errors.New("No user available to associate with opportunity")
  }
  o.WithUser(user)

	return o, nil
}

