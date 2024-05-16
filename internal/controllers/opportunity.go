package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	helpers "github.com/pmwals09/yobs/internal"
	"github.com/pmwals09/yobs/internal/models/contact"
	"github.com/pmwals09/yobs/internal/models/document"
	"github.com/pmwals09/yobs/internal/models/opportunity"
	"github.com/pmwals09/yobs/internal/models/status"
	"github.com/pmwals09/yobs/internal/models/user"
	"github.com/pmwals09/yobs/web/templates"
)

func HandlePostOppty(repo opportunity.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newOpportunity, err := newOpportunityFromRequest(r)
		if err != nil {
			helpers.WriteError(w, err)
		}

		user, err := userFromRequest(r)
		if err != nil {
			helpers.WriteError(w, errors.New("No user available"))
		}
		f := helpers.FormData{}
		if err := repo.CreateOpportunity(newOpportunity); err != nil {
			f.AddError("overall", fmt.Sprintf(
				"An error occurred creating the opportunity: %s",
				err.Error()))
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

		user, err := userFromRequest(r)
		if err != nil {
			helpers.WriteError(w, err)
			return
		}
		opp, err := opptyFromRequest(r, opptyRepo, user)
		if err != nil {
			helpers.WriteError(w, err)
			return
		}
		od.Oppty = *opp

		docs, err := opptyRepo.GetAllDocuments(opp, user)
		if err != nil {
			fd.AddError("document-table", "Unable to retrieve opportunity documents.")
		} else {
			for i := range docs {
				_, err := docs[i].GetPresignedDownloadUrl()
				if err != nil {
					fd.AddError("document-table", "Unable to retrieve document URL for download.")
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
			fd.AddError("existing-attachment", "Unable to retrieve user documents.")
		}

		contacts, err := opptyRepo.GetAllContacts(opp)
		if err != nil {
			fd.AddError("contacts", fmt.Sprintf("Unable to retrieve opportunity contacts: %s", err.Error()))
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
		var fd helpers.FormData
		user, err := userFromRequest(r)
		if err != nil {
			fd.AddError("overall", "Error retrieving user")
			docs, err := docRepo.GetAllUserDocuments(user)
			if err != nil {
				fd.AddError("existing-attachment", "Error retrieving user documents")
			}
			var oppty opportunity.Opportunity
			retargetAttachmentModal(w, r, oppty, docs, fd)
			return
		}
		oppty, err := opptyFromRequest(r, opptyRepo, user)
		if err != nil {
			fd.AddError("overall", "Error retrieving opportunity")
			docs, err := docRepo.GetAllUserDocuments(user)
			if err != nil {
				fd.AddError("existing-attachment", "Error retrieving user documents")
			}
			retargetAttachmentModal(w, r, *oppty, docs, fd)
			return
		}

		r.ParseMultipartForm(10 << 20)
		file, handler, err := r.FormFile("attachment-file")
		if file != nil {
			defer file.Close()
		}
		if err != nil {
			fd.AddError("attachment-file", "Problem parsing file - did you attach one?")
			docs, err := docRepo.GetAllUserDocuments(user)
			if err != nil {
				fd.AddError("existing-attachment", "Error retrieving user documents")
			}
			retargetAttachmentModal(w, r, *oppty, docs, fd)
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
			fd.AddError("overall", "Unable to add document to the opportunity")
			docs, err := docRepo.GetAllUserDocuments(user)
			if err != nil {
				fd.AddError("existing-attachment", "Error retrieving user documents")
			}
			retargetAttachmentModal(w, r, *oppty, docs, fd)
			return
		}

		docs, err := opptyRepo.GetAllDocuments(oppty, user)
		if err != nil {
			fd.Errors["overall"] = "Unable to retrieve associated documents after submission."
			retargetAttachmentModal(w, r, *oppty, docs, fd)
			return
		}

		for i := range docs {
			_, err := docs[i].GetPresignedDownloadUrl()
			if err != nil {
				w.Header().Add("HX-Retarget", "attachment-modal")
				fd.AddError("document-table", "Unable to retrieve document URL for download.")
				retargetAttachmentModal(w, r, *oppty, docs, fd)
				return
			}
		}

		returnAttachmentsSection(w, r, user, oppty, opptyRepo)
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
		WithURL(url)
	var initialStatus status.Status
	if date == "" {
		initialStatus.Name = status.None
	} else {
		initialStatus.Name = status.Applied
		t, err := time.Parse(time.DateOnly, date)
		if err == nil {
			initialStatus.Date = t
		}
	}
	o.Statuses = []status.Status{ initialStatus }
	user, err := userFromRequest(r)
	if user == nil {
		return o, errors.New("No user available to associate with opportunity")
	}
	o.WithUser(user)
	return o, nil
}

func HandleAddExistingToOppty(opptyRepo opportunity.Repository, docRepo document.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := userFromRequest(r)
		if err == nil {
			helpers.WriteError(w, errors.New("No user available"))
			return
		}
		oppty, err := opptyFromRequest(r, opptyRepo, user)
		var fd helpers.FormData
		if err != nil {
			fd.AddError("overall", "Unable to retrieve opportunity.")
			docs, err := docRepo.GetAllUserDocuments(user)
			if err != nil {
				fd.AddError("existing-attachment", "Unable to retrieve user docs")
			}
			retargetAttachmentModal(w, r, *oppty, docs, fd)
			return
		}

		// Get the selected document from the formdata
		err = r.ParseForm()
		if err != nil {
			fd.AddError("overall", "Unable to parse form")
			docs, _ := docRepo.GetAllUserDocuments(user)
			retargetAttachmentModal(w, r, *oppty, docs, fd)
			return
		}

		docIdStr := r.PostForm.Get("existing-attachment")
		if docIdStr == "" {
			fd.AddError("existing-attachment", "Must select a document")
			docs, err := docRepo.GetAllUserDocuments(user)
			if err != nil {
				fd.AddError("overall", "Unable to retrieve user docs")
			}
			retargetAttachmentModal(w, r, *oppty, docs, fd)
			return
		}
		docId, err := strconv.ParseUint(docIdStr, 10, 64)
		if err != nil {
			fd.AddError("existing-attachment", "Unable to parse document ID.")
			docs, err := docRepo.GetAllUserDocuments(user)
			if err != nil {
				fd.AddError("overall", "Unable to retrieve user docs")
			}
			retargetAttachmentModal(w, r, *oppty, docs, fd)
			return
		}

		doc, err := docRepo.GetDocumentById(uint(docId), user)

		// Associate the existing document with this opportunity
		err = opptyRepo.AddDocument(oppty, &doc)
		if err != nil {
			fd.AddError("existing-attachment", "Unable to add document to opportunity.")
			docs, err := docRepo.GetAllUserDocuments(user)
			if err != nil {
				fd.AddError("overall", "Unable to retrieve user docs")
			}
			retargetAttachmentModal(w, r, *oppty, docs, fd)
			return
		}

		returnAttachmentsSection(w, r, user, oppty, opptyRepo)
		return
	}
}

func HandleContactModal(opptyRepo opportunity.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := userFromRequest(r)
		var fd helpers.FormData
		if err != nil {
			allErrors := err
			oppty, err := opptyFromRequest(r, opptyRepo, user)
			if err != nil {
				allErrors = errors.Join(allErrors, err)
			}
			fd.AddError("overall", allErrors.Error())
			retargetContactModal(w, r, *oppty, fd)
			return
		}

		oppty, err := opptyFromRequest(r, opptyRepo, user)
		if err != nil {
			fd.AddError("overall", err.Error())
			retargetContactModal(w, r, *oppty, fd)
			return
		}

		templates.ContactModal(oppty, fd).Render(r.Context(), w)
	}
}

func HandleAttachmentModal(opptyRepo opportunity.Repository, docRepo document.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := userFromRequest(r)
		if err != nil {
			helpers.WriteError(w, errors.New("Invalid user"))
			return
		}
		var fd helpers.FormData

		oppty, err := opptyFromRequest(r, opptyRepo, u)
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
		var fd helpers.FormData
		if err != nil {
			allErrors := err
			user, err := userFromRequest(r)
			if err != nil {
				allErrors = errors.Join(allErrors, err)
			}
			oppty, err := opptyFromRequest(r, opptyRepo, user)
			if err != nil {
				allErrors = errors.Join(allErrors, err)
			}
			fd.AddError("overall", allErrors.Error())
			retargetContactModal(w, r, *oppty, fd)
			return
		}

		user, err := userFromRequest(r)
		if err != nil {
			allErrors := err
			oppty, err := opptyFromRequest(r, opptyRepo, user)
			if err != nil {
				allErrors = errors.Join(allErrors, err)
			}
			fd.AddError("overall", allErrors.Error())
			retargetContactModal(w, r, *oppty, fd)
			return
		}

		err = contactRepo.CreateContact(newContact, *user)
		if err != nil {
			allErrors := err
			oppty, err := opptyFromRequest(r, opptyRepo, user)
			if err != nil {
				allErrors = errors.Join(allErrors, err)
			}
			fd.AddError("overall", allErrors.Error())
			retargetContactModal(w, r, *oppty, fd)
			return
		}

		oppty, err := opptyFromRequest(r, opptyRepo, user)
		if err != nil {
			fd.AddError("overall", err.Error())
			retargetContactModal(w, r, *oppty, fd)
			return
		}

		err = opptyRepo.AddContact(oppty, *newContact)
		if err != nil {
			fd.AddError("overall", err.Error())
			retargetContactModal(w, r, *oppty, fd)
			return
		}

		// db stuff done, it's rendering time
		contacts, err := opptyRepo.GetAllContacts(oppty)
		if err != nil {
			fd.AddError("overall", err.Error())
			retargetContactModal(w, r, *oppty, fd)
			return
		}

		templates.ContactsTable(contacts).Render(r.Context(), w)
	}
}

func userFromRequest(r *http.Request) (*user.User, error) {
	if u, ok := r.Context().Value("user").(*user.User); !ok {
		return u, errors.New("Unable to retrieve user from context.")
	} else if u == nil {
		return u, errors.New("No user available")
	} else {
		return u, nil
	}
}

func opptyFromRequest(r *http.Request, opptyRepo opportunity.Repository, u *user.User) (*opportunity.Opportunity, error) {
	idParam := chi.URLParam(r, "opportunityId")
	id, err := strconv.ParseUint(idParam, 10, 64) // Sqlite id's are 64-bit int
	if err != nil {
		var o opportunity.Opportunity
		return &o, err
	}
	o, err := opptyRepo.GetOpportuntyById(uint(id), u)
	return o, err
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

func retargetAttachmentModal(
	w http.ResponseWriter,
	r *http.Request,
	oppty opportunity.Opportunity,
	docs []document.Document,
	fd helpers.FormData) {
	w.Header().Add("HX-Retarget", "#attachment-modal")
	templates.AttachmentModal(oppty, docs, fd).Render(r.Context(), w)
}

func retargetContactModal(
	w http.ResponseWriter,
	r *http.Request,
	oppty opportunity.Opportunity,
	fd helpers.FormData) {
	w.Header().Add("HX-Retarget", "#contact-modal")
	templates.ContactModal(&oppty, fd).Render(r.Context(), w)
}
