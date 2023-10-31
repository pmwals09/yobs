package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	helpers "github.com/pmwals09/yobs/internal"
	"github.com/pmwals09/yobs/internal/models/document"
	"github.com/pmwals09/yobs/internal/models/opportunity"
	"github.com/pmwals09/yobs/internal/models/user"
	"github.com/pmwals09/yobs/web/templates"
)

func HandlePostOppty(repo opportunity.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newOpportunity, err := newOpportunityFromRequest(r)
		if err != nil {
			helpers.WriteError(w, err)
		}

		user := r.Context().Value("user").(*user.User)
		if user == nil {
			helpers.WriteError(w, errors.New("No user available"))
		}
		f := helpers.FormData{}
		if err := repo.CreateOpportunity(newOpportunity); err != nil {
			f.Errors["overall"] = fmt.Sprintf(
				"An error occurred creating the opportunity: %s",
				err.Error(),
			)
			templates.HomePage(user, []opportunity.Opportunity{}, f).Render(r.Context(), w)
			return
		}

		opportunities, opptyErr := repo.GetAllOpportunities(user)
		if opptyErr != nil {
			helpers.WriteError(w, opptyErr)
			return
		}

		templates.HomePage(user, opportunities, f).Render(r.Context(), w)
    return
	}
}

func HandleGetOpptyPage(
	opptyRepo opportunity.Repository,
	docRepo document.Repository,
) http.HandlerFunc {
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
		opp, err := opptyRepo.GetOpportuntyById(uint(id), user)
		if err != nil {
			helpers.WriteError(w, err)
			return
		}

		od.Oppty = *opp
		docs, err := opptyRepo.GetAllDocuments(opp, user)
		if err != nil {
			fd := helpers.FormData{
				Errors: map[string]string{
					"document-table": "Unable to retrieve opportunity documents.",
				},
			}
			templates.OpportunityDetailsPage(
				user,
				od,
				docs,
				fd,
			).Render(r.Context(), w)
			return
		}

		for i := range docs {
			_, err := docs[i].GetPresignedDownloadUrl()
			if err != nil {
				fd := helpers.FormData{
					Errors: map[string]string{
						"document-table": "Unable to retrieve document URL for download.",
					},
				}
				templates.OpportunityDetailsPage(
					user,
					od,
					docs,
					fd,
				).Render(r.Context(), w)
				return
			}
		}

		od.Documents = docs

		userDocuments, err := docRepo.GetAllUserDocuments(user)
		if err != nil {
			fd := helpers.FormData{
				Errors: map[string]string{
					"existing-attachment": "Unable to retrieve user documents.",
				},
			}
			templates.OpportunityDetailsPage(
				user,
				od,
				userDocuments,
				fd,
			).Render(r.Context(), w)
			return
		}

		templates.OpportunityDetailsPage(
			user,
			od,
			userDocuments,
			helpers.FormData{},
		).Render(r.Context(), w)
	}
}


func HandleUploadToOppty(
	opptyRepo opportunity.Repository,
	docRepo document.Repository,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(*user.User)
		if user == nil {
			helpers.WriteError(w, errors.New("No user available"))
			return
		}

		idParam := chi.URLParam(r, "opportunityId")
		id, err := strconv.ParseUint(idParam, 10, 64) // Sqlite id's are 64-bit int
		if err != nil {
			helpers.WriteError(w, err)
			return
		}
		oppty, err := opptyRepo.GetOpportuntyById(uint(id), user)
		if err != nil {
			fd := helpers.FormData{
				Errors: map[string]string{
					"overall": "Error retrieving opportunity",
				},
			}
			returnAttachmentsSection(w, r, user, oppty, docRepo, opptyRepo, fd)
			return
		}

		r.ParseMultipartForm(10 << 20)
		file, handler, err := r.FormFile("attachment-file")
		if file != nil {
			defer file.Close()
		}
		if err != nil {
			fd := helpers.FormData{
				Errors: map[string]string{
					"attachment-file": "Problem parsing file - did you attach one?",
				},
			}
			returnAttachmentsSection(w, r, user, oppty, docRepo, opptyRepo, fd)
			return
		}

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

		d.WithUser(user)

		err = d.Upload(file)
		if err != nil {
			helpers.WriteError(w, err)
			return
		}

		if err := docRepo.CreateDocument(d); err != nil {
			helpers.WriteError(w, err)
			return
		}

		if err = opptyRepo.AddDocument(oppty, d); err != nil {
			fd := helpers.FormData{
				Errors: map[string]string{
					"overall": "Unable to add document to the opportunity",
				},
			}
			returnAttachmentsSection(w, r, user, oppty, docRepo, opptyRepo, fd)
			return
		}

		docs, err := opptyRepo.GetAllDocuments(oppty, user)
		if err != nil {
			fd := helpers.FormData{
				Errors: map[string]string{
					"overall": "Unable to retrieve associated documents after submission.",
				},
			}
			returnAttachmentsSection(w, r, user, oppty, docRepo, opptyRepo, fd)
			return
		}

		for i := range docs {
			_, err := docs[i].GetPresignedDownloadUrl()
			if err != nil {
				fd := helpers.FormData{
					Errors: map[string]string{
						"document-table": "Unable to retrieve document URL for download.",
					},
				}
				returnAttachmentsSection(w, r, user, oppty, docRepo, opptyRepo, fd)
				return
			}
		}

		fd := helpers.FormData{}
		returnAttachmentsSection(w, r, user, oppty, docRepo, opptyRepo, fd)
    return
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

func HandleAddExistingToOppty(opptyRepo opportunity.Repository, docRepo document.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "opportunityId")
		if idParam == "" {
			return
		}

		user := r.Context().Value("user").(*user.User)
		if user == nil {
			helpers.WriteError(w, errors.New("No user available"))
			return
		}

		id, err := strconv.ParseUint(idParam, 10, 64) // Sqlite id's are 64-bit int
		if err != nil {
			helpers.WriteError(w, err)
			return
		}
		oppty, err := opptyRepo.GetOpportuntyById(uint(id), user)
		if err != nil {
			fd := helpers.FormData{
				Errors: map[string]string{
					"overall": "Unable to retrieve opportunity.",
				},
			}
			returnAttachmentsSection(w, r, user, oppty, docRepo, opptyRepo, fd)
			return
		}

		// Get the selected document from the formdata
		err = r.ParseForm()
		if err != nil {
			fd := helpers.FormData{
				Errors: map[string]string{
					"overall": "Unable to parse form.",
				},
			}
			returnAttachmentsSection(w, r, user, oppty, docRepo, opptyRepo, fd)
			return
		}

		docIdStr := r.PostForm.Get("existing-attachment")
		if docIdStr == "" {
			fd := helpers.FormData{
				Errors: map[string]string{
					"existing-attachment": "Must select a document.",
				},
			}
			returnAttachmentsSection(w, r, user, oppty, docRepo, opptyRepo, fd)
			return
		}
		docId, err := strconv.ParseUint(docIdStr, 10, 64)
		if err != nil {
			fd := helpers.FormData{
				Errors: map[string]string{
					"existing-attachment": "Unable to parse document ID.",
				},
			}
			returnAttachmentsSection(w, r, user, oppty, docRepo, opptyRepo, fd)
			return
		}

		doc, err := docRepo.GetDocumentById(uint(docId), user)

		// Associate the existing document with this opportunity
		err = opptyRepo.AddDocument(oppty, &doc)
		if err != nil {
			fd := helpers.FormData{
				Errors: map[string]string{
					"overall": "Unable to add document to opportunity.",
				},
			}
			returnAttachmentsSection(w, r, user, oppty, docRepo, opptyRepo, fd)
			return
		}

		fd := helpers.FormData{}
		returnAttachmentsSection(w, r, user, oppty, docRepo, opptyRepo, fd)
    return
	}
}

// TODO: How to update an existing opportunity?

func HandleEditOpptyPage(opptyRepo opportunity.Repository, docRepo document.Repository) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    templates.OpportunityEditPage(nil).Render(r.Context(), w)
    // get the id out of the url
    // get the opportunity
    // get the associated documents
    // populate a form for editing the opportunity
    // return the form
  }
}

// TODO: Update an existing opportunity

func returnAttachmentsSection(
	w http.ResponseWriter,
	r *http.Request,
	user *user.User,
	oppty *opportunity.Opportunity,
	docRepo document.Repository,
	opptyRepo opportunity.Repository,
	fd helpers.FormData,
) {
	docs, err := opptyRepo.GetAllDocuments(oppty, user)
	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	od := helpers.OpptyDetails{
		Oppty:     *oppty,
		Documents: docs,
	}

	userDocs, err := docRepo.GetAllUserDocuments(user)
	if err != nil {
		helpers.WriteError(w, err)
		return
	}

  templates.OpportunityDetailsPage(user, od, userDocs, fd).Render(r.Context(), w)
  return
}
