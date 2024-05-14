package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	helpers "github.com/pmwals09/yobs/internal"
	"github.com/pmwals09/yobs/internal/models/contact"
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
		var od helpers.OpptyDetails
		var fd helpers.FormData

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
			if fd.Errors == nil {
				fd.Errors = map[string]string{}
			}

			fd.Errors["document-table"] = "Unable to retrieve opportunity documents."
		} else {
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
		}

		userDocuments, err := docRepo.GetAllUserDocuments(user)
		if err != nil {
			if fd.Errors == nil {
				fd.Errors = map[string]string{}
			}
			fd.Errors["existing-attachment"] = "Unable to retrieve user documents."
		}

		contacts, err := opptyRepo.GetAllContacts(opp)
		if err != nil {
			if fd.Errors == nil {
				fd.Errors = map[string]string{}
			}
			fd.Errors["contacts"] = fmt.Sprintf("Unable to retrieve opportunity contacts: %s", err.Error())

		} else {
			od.Contacts = contacts
		}

		templates.OpportunityDetailsPage(
			user,
			od,
			userDocuments,
			fd).Render(r.Context(), w)
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
			// fd := helpers.FormData{
			// 	Errors: map[string]string{
			// 		"overall": "Error retrieving opportunity",
			// 	},
			// }
			// TODO: What's the appropriate response?
			// should it use HX-Retarget: attachment-modal?
			returnAttachmentsSection(w, r, user, oppty, docRepo, opptyRepo)
			return
		}

		r.ParseMultipartForm(10 << 20)
		file, handler, err := r.FormFile("attachment-file")
		if file != nil {
			defer file.Close()
		}
		if err != nil {
			// fd := helpers.FormData{
			// 	Errors: map[string]string{
			// 		"attachment-file": "Problem parsing file - did you attach one?",
			// 	},
			// }
			// TODO: What's the appropriate response?
			// should it use HX-Retarget: attachment-modal?
			returnAttachmentsSection(w, r, user, oppty, docRepo, opptyRepo)
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
			// fd := helpers.FormData{
			// 	Errors: map[string]string{
			// 		"overall": "Unable to add document to the opportunity",
			// 	},
			// }
			// TODO: What's the appropriate response?
			// should it use HX-Retarget: attachment-modal?
			returnAttachmentsSection(w, r, user, oppty, docRepo, opptyRepo)
			return
		}

		docs, err := opptyRepo.GetAllDocuments(oppty, user)
		if err != nil {
			// fd := helpers.FormData{
			// 	Errors: map[string]string{
			// 		"overall": "Unable to retrieve associated documents after submission.",
			// 	},
			// }
			// TODO: What's the appropriate response?
			// should it use HX-Retarget: attachment-modal?
			returnAttachmentsSection(w, r, user, oppty, docRepo, opptyRepo)
			return
		}

		for i := range docs {
			_, err := docs[i].GetPresignedDownloadUrl()
			if err != nil {
				// fd := helpers.FormData{
				// 	Errors: map[string]string{
				// 		"document-table": "Unable to retrieve document URL for download.",
				// 	},
				// }
				// TODO: What's the appropriate response?
				// should it use HX-Retarget: attachment-modal?
				returnAttachmentsSection(w, r, user, oppty, docRepo, opptyRepo)
				return
			}
		}

		// fd := helpers.FormData{}
		// TODO: What's the appropriate response?
		// should it use HX-Retarget: attachment-modal?
		returnAttachmentsSection(w, r, user, oppty, docRepo, opptyRepo)
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
			// fd := helpers.FormData{
			// 	Errors: map[string]string{
			// 		"overall": "Unable to retrieve opportunity.",
			// 	},
			// }
			returnAttachmentsSection(w, r, user, oppty, docRepo, opptyRepo)
			return
		}

		// Get the selected document from the formdata
		err = r.ParseForm()
		if err != nil {
			// fd := helpers.FormData{
			// 	Errors: map[string]string{
			// 		"overall": "Unable to parse form.",
			// 	},
			// }
			returnAttachmentsSection(w, r, user, oppty, docRepo, opptyRepo)
			return
		}

		docIdStr := r.PostForm.Get("existing-attachment")
		if docIdStr == "" {
			// fd := helpers.FormData{
			// 	Errors: map[string]string{
			// 		"existing-attachment": "Must select a document.",
			// 	},
			// }
			returnAttachmentsSection(w, r, user, oppty, docRepo, opptyRepo)
			return
		}
		docId, err := strconv.ParseUint(docIdStr, 10, 64)
		if err != nil {
			// fd := helpers.FormData{
			// 	Errors: map[string]string{
			// 		"existing-attachment": "Unable to parse document ID.",
			// 	},
			// }
			returnAttachmentsSection(w, r, user, oppty, docRepo, opptyRepo)
			return
		}

		doc, err := docRepo.GetDocumentById(uint(docId), user)

		// Associate the existing document with this opportunity
		err = opptyRepo.AddDocument(oppty, &doc)
		if err != nil {
			// fd := helpers.FormData{
			// 	Errors: map[string]string{
			// 		"overall": "Unable to add document to opportunity.",
			// 	},
			// }
			returnAttachmentsSection(w, r, user, oppty, docRepo, opptyRepo)
			return
		}

		// fd := helpers.FormData{}
		returnAttachmentsSection(w, r, user, oppty, docRepo, opptyRepo)
		return
	}
}

func HandleContactModal(opptyRepo opportunity.Repository) http.HandlerFunc {
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

		templates.ContactModal(oppty).Render(r.Context(), w)
	}
}

func HandleAttachmentModal(opptyRepo opportunity.Repository, docRepo document.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "opportunityId")
		if idParam == "" {
			return
		}
		id, err := strconv.ParseUint(idParam, 10, 64) // Sqlite id's are 64-bit int
		if err != nil {
			return
		}

		u, ok := r.Context().Value("user").(*user.User)
		if !ok {
			helpers.WriteError(w, errors.New("Invalid user"))
			return
		}
		var fd helpers.FormData

		oppty, err := opptyRepo.GetOpportuntyById(uint(id), u)
		if err != nil {
			helpers.WriteError(w, err)
			return
		}

		userDocs, err := docRepo.GetAllUserDocuments(u)
		if err != nil {
			helpers.WriteError(w, err)
			return
		}
		templates.AttachmentModal(*oppty, userDocs, fd).Render(r.Context(), w)
	}
}

func HandleAddNewContactToOppty(opptyRepo opportunity.Repository, contactRepo contact.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newContact, err := newContactFromRequest(r)
		if err != nil {
			helpers.WriteError(w, err)
			return
		}

		user, ok := r.Context().Value("user").(*user.User)
		if !ok || user == nil {
			helpers.WriteError(w, errors.New("No user available"))
			return
		}

		err = contactRepo.CreateContact(newContact, *user)
		if err != nil {
			helpers.WriteError(w, err)
			return
		}

		idParam := chi.URLParam(r, "opportunityId")
		id, err := strconv.ParseUint(idParam, 10, 64)

		oppty, err := opptyRepo.GetOpportuntyById(uint(id), user)
		if err != nil {
			helpers.WriteError(w, err)
			return
		}

		err = opptyRepo.AddContact(oppty, *newContact)
		if err != nil {
			helpers.WriteError(w, err)
			return
		}

		// db stuff done, it's rendering time
		contacts, err := opptyRepo.GetAllContacts(oppty)
		if err != nil {
			helpers.WriteError(w, err)
			return
		}

		templates.ContactsTable(contacts).Render(r.Context(), w)
	}
}

func newContactFromRequest(r *http.Request) (*contact.Contact, error) {
	var c contact.Contact

	err := r.ParseForm()
	if err != nil {
		return &c, err
	}

	name := r.PostForm.Get("contact-name")
	c.Name = name
	company := r.PostForm.Get("company-name")
	c.CompanyName = company
	title := r.PostForm.Get("contact-title")
	c.Title = title
	phone := r.PostForm.Get("contact-phone")
	c.Phone = phone
	email := r.PostForm.Get("contact-email")
	c.Email = email

	if c.IsEmpty() {
		return &c, errors.New("Cannot create a valid contact from the provided information")
	}

	return &c, nil

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
) {
	docs, err := opptyRepo.GetAllDocuments(oppty, user)
	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	templates.AttachmentsTable(docs).Render(r.Context(), w)
	return
}
