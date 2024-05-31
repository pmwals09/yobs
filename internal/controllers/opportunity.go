package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	helpers "github.com/pmwals09/yobs/internal"
	"github.com/pmwals09/yobs/internal/models/contact"
	"github.com/pmwals09/yobs/internal/models/document"
	"github.com/pmwals09/yobs/internal/models/opportunity"
	"github.com/pmwals09/yobs/internal/models/status"
	"github.com/pmwals09/yobs/internal/models/user"
	homepage "github.com/pmwals09/yobs/web/home"
	opptydetailspage "github.com/pmwals09/yobs/web/opportunity-details"
)

func HandlePostOppty(repo opportunity.Repository, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newOpportunity, err := newOpportunityFromRequest(r)
		if err != nil {
			logger.Error("Problem getting oppty from request", "error", err)
			helpers.WriteError(w, err)
		}

		user, err := userFromRequest(r)
		if err != nil {
			logger.Error("Problem getting user from request", "error", err)
			helpers.WriteError(w, errors.New("no user available"))
		}
		f := helpers.FormData{}
		if err := repo.CreateOpportunity(newOpportunity); err != nil {
			logger.Error("Problem creating opportunity", "error", err)
			f.AddError("overall", fmt.Sprintf(
				"An error occurred creating the opportunity: %s",
				err.Error()))
			homepage.HomePage(user, []opportunity.Opportunity{}, f).Render(r.Context(), w)
			return
		}

		opportunities, opptyErr := repo.GetAllOpportunities(user)
		if opptyErr != nil {
			logger.Error("Problem getting all opportunities", "error", err)
			helpers.WriteError(w, opptyErr)
			return
		}

		homepage.HomePage(user, opportunities, f).Render(r.Context(), w)
	}
}

func HandleGetOpptyPage(
	opptyRepo opportunity.Repository,
	docRepo document.Repository,
	logger *slog.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var od helpers.OpptyDetails
		var fd helpers.FormData

		user, err := userFromRequest(r)
		if err != nil {
			logger.Error("Problem getting user from request")
			helpers.WriteError(w, err)
			return
		}
		opp, err := opptyFromRequest(r, opptyRepo, user)
		if err != nil {
			logger.Error("Problem getting oppty from request", "error", err)
			helpers.WriteError(w, err)
			return
		}
		od.Oppty = *opp

		docs, err := opptyRepo.GetAllDocuments(opp, user)
		if err != nil {
			logger.Error("Problem getting all documents", "error", err)
			fd.AddError("document-table", "Unable to retrieve opportunity documents.")
		} else {
			cfg := docRepo.GetConfig()
			for i := range docs {
				_, err := docs[i].GetPresignedDownloadUrl(cfg)
				if err != nil {
					logger.Error("Problem getting document URL", "id", docs[i].ID, "error", err)
					fd.AddError("document-table", "Unable to retrieve document URL for download.")
					opptydetailspage.OpportunityDetailsPage(
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
			logger.Error("Problem getting all documents", "error", err)
			fd.AddError("existing-attachment", "Unable to retrieve user documents.")
		}

		contacts, err := opptyRepo.GetAllContacts(opp)
		if err != nil {
			logger.Error("Problem getting all contacts", "error", err)
			fd.AddError("contacts", fmt.Sprintf("Unable to retrieve opportunity contacts: %s", err.Error()))
		} else {
			od.Contacts = contacts
		}

		opptydetailspage.OpportunityDetailsPage(
			user,
			od,
			userDocuments,
			fd).Render(r.Context(), w)
	}
}

func HandleUploadToOppty(
	opptyRepo opportunity.Repository,
	docRepo document.Repository,
	logger *slog.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var fd helpers.FormData
		user, err := userFromRequest(r)
		if err != nil {
			logger.Error("Problem getting user from request", "error", err)
			fd.AddError("overall", "Error retrieving user")
			docs, err := docRepo.GetAllUserDocuments(user)
			if err != nil {
				logger.Error("Problem getting all user documents", "error", err)
				fd.AddError("existing-attachment", "Error retrieving user documents")
			}
			var oppty opportunity.Opportunity
			retargetAttachmentModal(w, r, oppty, docs, fd)
			return
		}
		oppty, err := opptyFromRequest(r, opptyRepo, user)
		if err != nil {
			logger.Error("Problem getting oppty from request", "error", err)
			fd.AddError("overall", "Error retrieving opportunity")
			docs, err := docRepo.GetAllUserDocuments(user)
			if err != nil {
				logger.Error("Problem getting all user documents", "error", err)
				fd.AddError("existing-attachment", "Error retrieving user documents")
			}
			retargetAttachmentModal(w, r, *oppty, docs, fd)
			return
		}

		err = r.ParseMultipartForm(10 << 20)
		if err != nil {
			msg := "Problem parsing file - did you attach one?"
			logger.Error(msg, "error", err)
			fd.AddError("attachment-file", msg)
			docs, err := docRepo.GetAllUserDocuments(user)
			if err != nil {
				logger.Error("Problem getting all user documents", "error", err)
				fd.AddError("existing-attachment", "Error retrieving user documents")
			}
			retargetAttachmentModal(w, r, *oppty, docs, fd)
			return
		}
		file, handler, err := r.FormFile("attachment-file")
		if file != nil {
			defer file.Close()
		}
		if err != nil {
			msg := "Problem parsing file - did you attach one?"
			logger.Error(msg, "error", err)
			fd.AddError("attachment-file", msg)
			docs, err := docRepo.GetAllUserDocuments(user)
			if err != nil {
				logger.Error("Problem getting all user documents", "error", err)
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

		cfg := docRepo.GetConfig()
		err = d.Upload(file, cfg)
		if err != nil {
			logger.Error("Problem uploading document", "error", err)
			helpers.WriteError(w, err)
			return
		}

		if err := docRepo.CreateDocument(d); err != nil {
			logger.Error("Problem creating document", "error", err)
			helpers.WriteError(w, err)
			return
		}

		if err = opptyRepo.AddDocument(oppty, d); err != nil {
			msg := "Unable to add document to the opportunity"
			logger.Error(msg, "error", err)
			fd.AddError("overall", msg)
			docs, err := docRepo.GetAllUserDocuments(user)
			if err != nil {
				msg := "Error retrieving user documents"
				logger.Error(msg, "error", err)
				fd.AddError("existing-attachment", msg)
			}
			retargetAttachmentModal(w, r, *oppty, docs, fd)
			return
		}

		docs, err := opptyRepo.GetAllDocuments(oppty, user)
		if err != nil {
			msg := "Unable to retrieve associated documents after submission."
			logger.Error(msg, "error", err)
			fd.AddError("overall", msg)
			retargetAttachmentModal(w, r, *oppty, docs, fd)
			return
		}

		for i := range docs {
			_, err := docs[i].GetPresignedDownloadUrl(cfg)
			if err != nil {
				w.Header().Add("HX-Retarget", "attachment-modal")
				msg := "Unable to retrieve document URL for download."
				logger.Error(msg, "id", docs[i].ID, "error", err)
				fd.AddError("document-table", msg)
				retargetAttachmentModal(w, r, *oppty, docs, fd)
				return
			}
		}

		returnAttachmentsSection(w, r, user, oppty, opptyRepo)
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
	o.Statuses = []status.Status{initialStatus}
	user, err := userFromRequest(r)
	if user == nil || err != nil {
		return o, errors.New("no user available to associate with opportunity")
	}
	o.WithUser(user)
	return o, nil
}

func HandleAddExistingToOppty(opptyRepo opportunity.Repository, docRepo document.Repository, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := userFromRequest(r)
		if err == nil {
			msg := "no user available"
			logger.Error(msg, "error", err)
			helpers.WriteError(w, errors.New(msg))
			return
		}
		oppty, err := opptyFromRequest(r, opptyRepo, user)
		var fd helpers.FormData
		if err != nil {
			msg := "Unable to retrieve opportunity."
			logger.Error(msg, "error", err)
			fd.AddError("overall", msg)
			docs, err := docRepo.GetAllUserDocuments(user)
			if err != nil {
				msg := "Unable to retrieve user docs"
				logger.Error(msg, "error", err)
				fd.AddError("existing-attachment", msg)
			}
			retargetAttachmentModal(w, r, *oppty, docs, fd)
			return
		}

		// Get the selected document from the formdata
		err = r.ParseForm()
		if err != nil {
			msg := "Unable to parse form"
			logger.Error(msg, "error", err)
			fd.AddError("overall", msg)
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
			msg := "Unable to parse document ID."
			logger.Error(msg, "error", err)
			fd.AddError("existing-attachment", msg)
			docs, err := docRepo.GetAllUserDocuments(user)
			if err != nil {
				msg := "Unable to retrieve user docs"
				logger.Error(msg, "error", err)
				fd.AddError("overall", msg)
			}
			retargetAttachmentModal(w, r, *oppty, docs, fd)
			return
		}

		doc, err := docRepo.GetDocumentById(uint(docId), user)
		if err != nil {
			msg := "Unable to retrieve document to add"
			logger.Error(msg, "error", err)
			fd.AddError("existing-attachment", msg)
			docs, err := docRepo.GetAllUserDocuments(user)
			if err != nil {
				msg := "Unable to retrieve user docs"
				logger.Error(msg, "error", err)
				fd.AddError("overall", msg)
			}
			retargetAttachmentModal(w, r, *oppty, docs, fd)
			return
		}

		// Associate the existing document with this opportunity
		err = opptyRepo.AddDocument(oppty, &doc)
		if err != nil {
			msg := "Unable to add document to opportunity."
			logger.Error(msg, "error", err)
			fd.AddError("existing-attachment", msg)
			docs, err := docRepo.GetAllUserDocuments(user)
			if err != nil {
				msg := "Unable to retrieve user docs"
				logger.Error(msg, "error", err)
				fd.AddError("overall", msg)
			}
			retargetAttachmentModal(w, r, *oppty, docs, fd)
			return
		}

		returnAttachmentsSection(w, r, user, oppty, opptyRepo)
	}
}

func HandleContactModal(opptyRepo opportunity.Repository, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := userFromRequest(r)
		var fd helpers.FormData
		if err != nil {
			logger.Error("Can't get user from request", "error", err)
			allErrors := err
			oppty, err := opptyFromRequest(r, opptyRepo, user)
			if err != nil {
				logger.Error("Can't get oppty from request", "error", err)
				allErrors = errors.Join(allErrors, err)
			}
			fd.AddError("overall", allErrors.Error())
			retargetContactModal(w, r, *oppty, fd)
			return
		}

		oppty, err := opptyFromRequest(r, opptyRepo, user)
		if err != nil {
			logger.Error("Can't get oppty from request", "error", err)
			fd.AddError("overall", err.Error())
			retargetContactModal(w, r, *oppty, fd)
			return
		}

		opptydetailspage.ContactModal(oppty, fd).Render(r.Context(), w)
	}
}

func HandleAttachmentModal(opptyRepo opportunity.Repository, docRepo document.Repository, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := userFromRequest(r)
		if err != nil {
			msg := "invalid user"
			logger.Error(msg, "error", err)
			helpers.WriteError(w, errors.New(msg))
			return
		}
		var fd helpers.FormData

		oppty, err := opptyFromRequest(r, opptyRepo, u)
		if err != nil {
			logger.Error("Can't get oppty from request", "error", err)
			helpers.WriteError(w, err)
			return
		}

		userDocs, err := docRepo.GetAllUserDocuments(u)
		if err != nil {
			logger.Error("Can't get all user documents", "error", err)
			helpers.WriteError(w, err)
			return
		}
		opptydetailspage.AttachmentModal(*oppty, userDocs, fd).Render(r.Context(), w)
	}
}

func HandleAddNewContactToOppty(opptyRepo opportunity.Repository, contactRepo contact.Repository, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newContact, err := newContactFromRequest(r)
		var fd helpers.FormData
		if err != nil {
			logger.Error("Can't get contact from request", "error", err)
			allErrors := err
			user, err := userFromRequest(r)
			if err != nil {
				logger.Error("Can't get user from request", "error", err)
				allErrors = errors.Join(allErrors, err)
			}
			oppty, err := opptyFromRequest(r, opptyRepo, user)
			if err != nil {
				logger.Error("Can't get oppty from request", "error", err)
				allErrors = errors.Join(allErrors, err)
			}
			fd.AddError("overall", allErrors.Error())
			retargetContactModal(w, r, *oppty, fd)
			return
		}

		user, err := userFromRequest(r)
		if err != nil {
			logger.Error("Can't get user from request", "error", err)
			allErrors := err
			oppty, err := opptyFromRequest(r, opptyRepo, user)
			if err != nil {
				logger.Error("Can't get oppty from request", "error", err)
				allErrors = errors.Join(allErrors, err)
			}
			fd.AddError("overall", allErrors.Error())
			retargetContactModal(w, r, *oppty, fd)
			return
		}

		err = contactRepo.CreateContact(newContact, *user)
		if err != nil {
			logger.Error("Error creating contact", "error", err)
			allErrors := err
			oppty, err := opptyFromRequest(r, opptyRepo, user)
			if err != nil {
				logger.Error("Can't get oppty from request", "error", err)
				allErrors = errors.Join(allErrors, err)
			}
			fd.AddError("overall", allErrors.Error())
			retargetContactModal(w, r, *oppty, fd)
			return
		}

		oppty, err := opptyFromRequest(r, opptyRepo, user)
		if err != nil {
			logger.Error("Can't get oppty from request", "error", err)
			fd.AddError("overall", err.Error())
			retargetContactModal(w, r, *oppty, fd)
			return
		}

		err = opptyRepo.AddContact(oppty, *newContact)
		if err != nil {
			logger.Error("Can't add contact to oppty", "error", err)
			fd.AddError("overall", err.Error())
			retargetContactModal(w, r, *oppty, fd)
			return
		}

		// db stuff done, it's rendering time
		contacts, err := opptyRepo.GetAllContacts(oppty)
		if err != nil {
			logger.Error("Can't get all oppty contacts", "error", err)
			fd.AddError("overall", err.Error())
			retargetContactModal(w, r, *oppty, fd)
			return
		}

		opptydetailspage.ContactsTable(oppty.ID, contacts, fd).Render(r.Context(), w)
	}
}

func HandleStatusModal(opptyRepo opportunity.Repository, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := userFromRequest(r)
		var fd helpers.FormData
		if err != nil {
			logger.Error("Can't get user from request", "error", err)
			fd.AddError("overall", err.Error())
			oppty, err := opptyFromRequest(r, opptyRepo, user)
			if err != nil {
				logger.Error("Can't get oppty from request", "error", err)
				fd.AddError("overall", err.Error())
			}
			retargetStatusModal(w, r, *oppty, fd)
			return
		}
		oppty, err := opptyFromRequest(r, opptyRepo, user)
		if err != nil {
			logger.Error("Can't get oppty from request", "error", err)
			fd.AddError("overall", err.Error())
			retargetStatusModal(w, r, *oppty, fd)
			return
		}
		opptydetailspage.StatusModal(*oppty, fd).Render(r.Context(), w)
	}
}

func HandleUpdateStatus(opptyRepo opportunity.Repository, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var fd helpers.FormData
		s, err := newStatusFromRequest(r)
		if s.IsEmpty() || err != nil {
			msg := "Unable to update status with provided data"
			logger.Error(msg, "error", err)
			fd.AddError("overall", msg)
			user, err := userFromRequest(r)
			if err != nil {
				logger.Error("Can't get user from request", "error", err)
				fd.AddError("overall", err.Error())
			}
			oppty, err := opptyFromRequest(r, opptyRepo, user)
			if err != nil {
				logger.Error("Can't get oppty from request", "error", err)
				fd.AddError("overall", err.Error())
			}
			retargetStatusModal(w, r, *oppty, fd)
			return
		}
		user, err := userFromRequest(r)
		if err != nil {
			logger.Error("Can't get user from request", "error", err)
			oppty, err := opptyFromRequest(r, opptyRepo, user)
			if err != nil {
				logger.Error("Can't get oppty from request", "error", err)
				fd.AddError("overall", err.Error())
			}
			retargetStatusModal(w, r, *oppty, fd)
			return
		}
		oppty, err := opptyFromRequest(r, opptyRepo, user)
		if err != nil {
			logger.Error("Can't get oppty from request", "error", err)
			fd.AddError("overall", err.Error())
			retargetStatusModal(w, r, *oppty, fd)
			return
		}
		err = opptyRepo.UpdateStatus(oppty, *s)
		if err != nil {
			logger.Error("Can't update oppty status", "error", err)
			fd.AddError("overall", err.Error())
			retargetStatusModal(w, r, *oppty, fd)
			return
		}

		oppty.Statuses = append(oppty.Statuses, *s)
		slices.SortFunc(oppty.Statuses, func(a, b status.Status) int {
			if a.Date.Equal(b.Date) {
				return 0
			}
			if a.Date.After(b.Date) {
				return -1
			}
			return 1
		})
		buf := new(bytes.Buffer)
		opptydetailspage.StatusTable(oppty.ID, oppty.Statuses).Render(r.Context(), buf)
		opptydetailspage.OpptyDetailGrid(*oppty, true).Render(r.Context(), buf)
		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
	}
}

func HandleGetEditDetailsForm(opptyRepo opportunity.Repository, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var fd helpers.FormData
		u, err := userFromRequest(r)
		if err != nil {
			logger.Error("Can't get user from request", "error", err)
			fd.AddError("overall", err.Error())
		}
		oppty, err := opptyFromRequest(r, opptyRepo, u)
		if err != nil {
			logger.Error("Can't get oppty from request", "error", err)
			fd.AddError("overall", err.Error())
		}
		opptydetailspage.OpportunityDetailForm(*oppty, fd).Render(r.Context(), w)
	}
}

func HandleUpdate(opptyRepo opportunity.Repository, docRepo document.Repository, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var fd helpers.FormData
		var od helpers.OpptyDetails
		u, err := userFromRequest(r)
		if err != nil {
			logger.Error("Can't get user from request", "error", err)
			fd.AddError("overall", err.Error())
		}
		oppty, err := opptyFromRequest(r, opptyRepo, u)
		if err != nil {
			logger.Error("Can't get oppty from request", "error", err)
			fd.AddError("overall", err.Error())
		}
		err = updateOpptyFromRequest(r, oppty)
		if err != nil {
			logger.Error("Can't update oppty from request", "error", err)
			fd.AddError("overall", err.Error())
			// NOTE: If we've got this far, we shouldn't error out here when we
			// run this query again
			oppty, _ := opptyFromRequest(r, opptyRepo, u)
			od.Oppty = *oppty
			docs, err := opptyRepo.GetAllDocuments(oppty, u)
			if err != nil {
				logger.Error("Can't get all documents", "error", err)
				fd.AddError("overall", err.Error())
			}
			opptydetailspage.OpportunityDetailsPage(u, od, docs, fd).Render(r.Context(), w)
			return
		}

		err = opptyRepo.UpdateOpportunity(oppty)
		if err != nil {
			logger.Error("Can't update oppty", "error", err)
			fd.AddError("overall", err.Error())
			// NOTE: If we've got this far, we shouldn't error out here when we
			// run this query again
			oppty, _ := opptyFromRequest(r, opptyRepo, u)
			od.Oppty = *oppty
			docs, err := opptyRepo.GetAllDocuments(oppty, u)
			if err != nil {
				logger.Error("Can't get all user documents", "error", err)
				fd.AddError("overall", err.Error())
			}
			opptydetailspage.OpportunityDetailsPage(u, od, docs, fd).Render(r.Context(), w)
			return
		}
		helpers.HTMXRedirect(w, fmt.Sprintf("/opportunities/%d", oppty.ID), http.StatusFound)
	}
}

func HandleDeleteStatus(statusRepo status.Repository, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statusId := chi.URLParam(r, "statusID")
		opptyId := chi.URLParam(r, "opportunityId")
		sId, _ := strconv.ParseUint(statusId, 10, 64)
		statusRepo.DeleteStatusByID(uint(sId))
		helpers.HTMXRedirect(w, fmt.Sprintf("/opportunities/%s", opptyId), http.StatusFound)
	}
}

func HandleRemoveDocFromOppty(opptyRepo opportunity.Repository, docRepo document.Repository, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := r.Context().Value(user.UserCtxKey).(*user.User)
		if !ok {
			msg := "Error getting user from context"
			logger.Error(msg)
			fmt.Println(msg)
			helpers.HTMXRedirect(w, "/home", http.StatusFound)
			return
		}
		opptyId, err := strconv.ParseUint(chi.URLParam(r, "opportunityId"), 10, 64)
		if err != nil {
			msg := "Error parsing opportunity ID"
			logger.Error(msg, "error", err)
			helpers.HTMXRedirect(w, fmt.Sprintf("/opportunities/%d", opptyId), http.StatusFound)
			return
		}
		oppty, err := opptyRepo.GetOpportuntyById(uint(opptyId), u)
		if err != nil {
			logger.Error("Error getting opportunity:", "error", err)
			helpers.HTMXRedirect(w, fmt.Sprintf("/opportunities/%d", opptyId), http.StatusFound)
			return
		}
		docId, err := strconv.ParseUint(chi.URLParam(r, "documentId"), 10, 64)
		if err != nil {
			logger.Error("Error parsing document ID:", "error", err)
			helpers.HTMXRedirect(w, fmt.Sprintf("/opportunities/%d", opptyId), http.StatusFound)
			return
		}
		doc, err := docRepo.GetDocumentById(uint(docId), u)
		if err != nil {
			logger.Error("Error getting document:", "error", err)
			helpers.HTMXRedirect(w, fmt.Sprintf("/opportunities/%d", opptyId), http.StatusFound)
			return
		}
		err = opptyRepo.RemoveDocument(oppty, doc)
		if err != nil {
			logger.Error("Error removing document from oppty:", "error", err)
			helpers.HTMXRedirect(w, fmt.Sprintf("/opportunities/%d", opptyId), http.StatusFound)
			return
		}
		helpers.HTMXRedirect(w, fmt.Sprintf("/opportunities/%d", opptyId), http.StatusFound)
	}
}

func HandleContactRowForm(contactRepo contact.Repository, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var fd helpers.FormData
		opptyID, err := strconv.ParseUint(chi.URLParam(r, "opportunityId"), 10, 64)
		if err != nil {
			msg := "Cannot parse opportunity ID"
			logger.Error(msg, "error", err)
			fd.AddError("overall", msg)
		}
		contactID, err := strconv.ParseUint(chi.URLParam(r, "contactId"), 10, 64)
		if err != nil {
			msg := "Cannot parse contact ID"
			logger.Error(msg, "error", err)
			fd.AddError("overall", msg)
		}
		if fd.Errors != nil && fd.Errors["overall"] != nil && len(fd.Errors["overall"]) > 0 {
			opptydetailspage.ContactTableRowForm(uint(opptyID), uint(contactID), fd).Render(r.Context(), w)
			return
		}
		u, ok := r.Context().Value(user.UserCtxKey).(*user.User)
		if !ok {
			msg := "No user available"
			logger.Error(msg, "error", err)
			fd.AddError("overall", msg)
			opptydetailspage.ContactTableRowForm(uint(opptyID), uint(contactID), fd).Render(r.Context(), w)
			return
		}
		contact, err := contactRepo.GetContactById(uint(contactID), *u)
		if err != nil {
			msg := "Cannot find contact"
			logger.Error(msg, "error", err)
			fd.AddError("overall", msg)
			opptydetailspage.ContactTableRowForm(uint(opptyID), uint(contactID), fd).Render(r.Context(), w)
			return
		}
		fd.Values = contact.ToFormDataValues()
		opptydetailspage.ContactTableRowForm(uint(opptyID), uint(contactID), fd).Render(r.Context(), w)
	}
}

func HandleUpdateContact(contactRepo contact.Repository, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var fd helpers.FormData
		opptyID, err := strconv.ParseUint(chi.URLParam(r, "opportunityId"), 10, 64)
		if err != nil {
			msg := "Cannot parse opportunity ID"
			logger.Error(msg, "error", err)
			fd.AddError("overall", msg)
		}
		contactID, err := strconv.ParseUint(chi.URLParam(r, "contactId"), 10, 64)
		if err != nil {
			msg := "Cannot parse contact ID"
			logger.Error(msg, "error", err)
			fd.AddError("overall", msg)
		}
		if fd.Errors != nil && fd.Errors["overall"] != nil && len(fd.Errors["overall"]) > 0 {
			opptydetailspage.ContactTableRowForm(uint(opptyID), uint(contactID), fd)
			return
		}
		err = r.ParseForm()
		if err != nil {
			msg := "Cannot parse form"
			logger.Error(msg, "error", err)
			fd.AddError("overall", msg)
		}
		contact := contact.Contact{
			ID:          uint(contactID),
			Name:        r.PostForm.Get("contact-name"),
			CompanyName: r.PostForm.Get("contact-company-name"),
			Title:       r.PostForm.Get("contact-title"),
			Phone:       r.PostForm.Get("contact-phone"),
			Email:       r.PostForm.Get("contact-email"),
		}
		err = contactRepo.UpdateContact(contact)
		if err != nil {
			msg := "Cannot update contact"
			logger.Error(msg, "error", err)
			fd.AddError("overall", msg)
			opptydetailspage.ContactTableRowForm(uint(opptyID), uint(contactID), fd)
			return
		}
		opptydetailspage.ContactTableRow(uint(opptyID), contact, fd).Render(r.Context(), w)
	}
}

func HandleDeleteContact(contactRepo contact.Repository, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var allErr error
		contactID, err := strconv.ParseUint(chi.URLParam(r, "contactId"), 10, 64)
		if err != nil {
			logger.Error("Can't get contactID from url", "error", err)
			allErr = errors.Join(allErr, err)
		}
		opptyID, err := strconv.ParseInt(chi.URLParam(r, "opportunityId"), 10, 64)
		if err != nil {
			logger.Error("Can't get opptyID from url", "error", err)
			allErr = errors.Join(allErr, err)
		}
		if allErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("<tr>Error deleting contact</tr>"))
			return
		}
		u, ok := r.Context().Value(user.UserCtxKey).(*user.User)
		if !ok {
			logger.Error("Invalid user")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("<tr>Invalid user</tr>"))
			return
		}
		contact, err := contactRepo.GetContactById(uint(contactID), *u)
		if err != nil {
			var fd helpers.FormData
			msg := "Error finding contact to delete."
			logger.Error(msg, "error", err)
			fd.AddError("actions", msg)
			opptydetailspage.ContactTableRow(uint(opptyID), contact, fd).Render(r.Context(), w)
		}
		err = contactRepo.DeleteContact(uint(opptyID), contact)
		if err != nil {
			var fd helpers.FormData
			msg := "Error deleting contact."
			logger.Error(msg, "error", err)
			fd.AddError("actions", msg)
			opptydetailspage.ContactTableRow(uint(opptyID), contact, fd).Render(r.Context(), w)
		}
	}
}

func updateOpptyFromRequest(r *http.Request, oppty *opportunity.Opportunity) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	company := r.PostForm.Get("company-name")
	oppty.CompanyName = company
	role := r.PostForm.Get("company-role")
	oppty.Role = role
	url := r.PostForm.Get("role-url")
	oppty.URL = url
	description := r.PostForm.Get("job-description")
	oppty.Description = description

	if oppty.IsEmpty() {
		return errors.New("insufficient data to update opportunity")
	}
	return nil
}

func userFromRequest(r *http.Request) (*user.User, error) {
	if u, ok := r.Context().Value(user.UserCtxKey).(*user.User); !ok {
		return u, errors.New("unable to retrieve user from context")
	} else if u == nil {
		return u, errors.New("no user available")
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
		return &c, errors.New("cannot create a valid contact from the provided information")
	}

	return &c, nil
}

func newStatusFromRequest(r *http.Request) (*status.Status, error) {
	var s status.Status

	err := r.ParseForm()
	if err != nil {
		return &s, err
	}

	s.Name = r.PostForm.Get("status-name")
	date := r.PostForm.Get("status-date")
	t, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return &s, err
	}
	s.Date = t
	s.Note = r.PostForm.Get("status-note")
	return &s, nil
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

	opptydetailspage.AttachmentsTable(oppty.ID, docs).Render(r.Context(), w)
}

func retargetAttachmentModal(
	w http.ResponseWriter,
	r *http.Request,
	oppty opportunity.Opportunity,
	docs []document.Document,
	fd helpers.FormData) {
	w.Header().Add("HX-Retarget", "#attachment-modal")
	opptydetailspage.AttachmentModal(oppty, docs, fd).Render(r.Context(), w)
}

func retargetContactModal(
	w http.ResponseWriter,
	r *http.Request,
	oppty opportunity.Opportunity,
	fd helpers.FormData) {
	w.Header().Add("HX-Retarget", "#contact-modal")
	opptydetailspage.ContactModal(&oppty, fd).Render(r.Context(), w)
}

func retargetStatusModal(
	w http.ResponseWriter,
	r *http.Request,
	oppty opportunity.Opportunity,
	fd helpers.FormData) {
	w.Header().Add("HX-Retarget", "#status-modal")
	opptydetailspage.StatusModal(oppty, fd).Render(r.Context(), w)
}
