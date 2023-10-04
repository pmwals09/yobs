package controllers

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	helpers "github.com/pmwals09/yobs/internal"
	"github.com/pmwals09/yobs/internal/models/document"
	"github.com/pmwals09/yobs/internal/models/opportunity"
	"github.com/pmwals09/yobs/internal/models/user"
	templates "github.com/pmwals09/yobs/web/template"
)

func HandlePostOppty(repo opportunity.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newOpportunity, err := newOpportunityFromRequest(r)
    if err != nil {
      helpers.WriteError(w, err)
    }

		if err := repo.CreateOpportunity(newOpportunity); err != nil {
			handleCreateOpptyError(w, r, err)
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

		r.Header.Add("HX-Retarget", "#main-content")
    templates.OpportunityList(opportunities).Render(r.Context(), w)
	}
}

func handleCreateOpptyError(w http.ResponseWriter, r *http.Request, err error) {
  f := helpers.FormData {
    Errors: map[string]template.HTML{
      "overall": template.HTML(
        fmt.Sprintf(
          "<ul><li>An error occurred creating the opportunity: %s</li></ul>",
          err.Error(),
        ),
      ),
    },
  }
  templates.OpportunityForm(f).Render(r.Context(), w)
}

func HandleGetActiveOpptys(repo opportunity.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
    user := r.Context().Value("user").(*user.User)
    if user == nil {
      helpers.WriteError(w, errors.New("No user available"))
    }
		opptys, err := repo.GetAllOpportunities(user)
		if err != nil {
			helpers.WriteError(w, err)
			return
		}
    templates.OpportunityList(opptys).Render(r.Context(), w)
	}
}

// TODO: Should only do so for the current user
func HandleGetOppty(repo opportunity.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		od := helpers.OpptyDetails{}

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

    templates.OpportunityDetailsPage(od).Render(r.Context(), w)
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

    user := r.Context().Value("user").(*user.User)
    if user == nil {
			helpers.WriteError(w, errors.New("No user available"))
			return
    }
    d.WithUser(user)

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

		od := helpers.OpptyDetails{
			Oppty:     *oppty,
			Documents: docs,
		}
    templates.AttachmentsSection(od).Render(r.Context(), w)
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

