package controllers

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	helpers "github.com/pmwals09/yobs/internal"
	"github.com/pmwals09/yobs/internal/models/opportunity"
	"github.com/pmwals09/yobs/internal/models/status"
	"github.com/pmwals09/yobs/internal/models/user"
	"github.com/pmwals09/yobs/web/opportunity-details"
)

func HandleStatusRowForm(statusRepo status.Repository, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statusId := chi.URLParam(r, "statusID")
		id, err := strconv.ParseUint(statusId, 10, 64)
		opportunityId := chi.URLParam(r, "opportunityId")
		opptyId, oErr := strconv.ParseUint(opportunityId, 10, 64)
		err = errors.Join(err, oErr)
		var fd helpers.FormData
		if err != nil {
			msg := "Problem parsing status ID from path"
			logger.Error(msg, "error", err)
			fd.AddError("overall", msg)
			var status status.Status
			opportunitydetailspage.StatusTableRowForm(uint(opptyId), status, fd).Render(r.Context(), w)
			return

		}
		status, err := statusRepo.GetStatusByID(uint(id))
		if err != nil {
			msg := "Problem getting status by ID"
			logger.Error(msg, "error", err)
			fd.AddError("overall", msg)
			opportunitydetailspage.StatusTableRowForm(uint(opptyId), status, fd).Render(r.Context(), w)
			return
		}
		fd.Values = status.ToFormDataValues()

		opportunitydetailspage.StatusTableRowForm(uint(opptyId), status, fd).Render(r.Context(), w)
	}
}

func HandleUpdateStatusItem(statusRepo status.Repository, opptyRepo opportunity.Repository, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var fd helpers.FormData

		// get the status id from the request
		statusId := chi.URLParam(r, "statusID")
		id, err := strconv.ParseUint(statusId, 10, 64)
		opportunityId := chi.URLParam(r, "opportunityId")
		opptyId, oErr := strconv.ParseUint(opportunityId, 10, 64)
		err = errors.Join(err, oErr)
		if err != nil {
			msg := "Error parsing IDs from route"
			logger.Error(msg, "error", err)
			fd.AddError("overall", msg)
			err := r.ParseForm()
			if err != nil {
				msg := "Error parsing form"
				fd.AddError("overall", msg)
				logger.Error(msg, "error", err)
			}
			fd.Values = map[string]string{
				"status-name": r.PostForm.Get("status-name"),
				"status-date": r.PostForm.Get("status-date"),
				"status-note": r.PostForm.Get("status-note"),
			}
			// NOTE: We really just need the id for this particular template
			var s status.Status
			s.ID = uint(id)
			opportunitydetailspage.StatusTableRowForm(uint(opptyId), s, fd).Render(r.Context(), w)
			return
		}

		// get the updated status form data from the request
		s, err := statusFromRequest(r)
		if err != nil {
			msg := "Error parsing status data"
			logger.Error(msg, "error", err)
			fd.AddError("overall", msg)
			err := r.ParseForm()
			if err != nil {
				msg := "Error parsing form"
				fd.AddError("overall", msg)
				logger.Error(msg, "error", err)
			}
			fd.Values = map[string]string{
				"status-name": r.PostForm.Get("status-name"),
				"status-date": r.PostForm.Get("status-date"),
				"status-note": r.PostForm.Get("status-note"),
			}
			// NOTE: We really just need the id for this particular template
			var s status.Status
			s.ID = uint(id)
			opportunitydetailspage.StatusTableRowForm(uint(opptyId), s, fd).Render(r.Context(), w)
			return
		}
		s.ID = uint(id)

		// attempt to upate the resource in the database
		err = statusRepo.UpdateStatus(s)
		if err != nil {
			msg := "Error updating status in database"
			logger.Error(msg, "error", err)
			fd.AddError("overall", msg)
			fd.Values = map[string]string{
				"status-name": s.Name,
				"status-date": s.Date.Format(time.DateOnly),
				"status-note": s.Note,
			}
			// NOTE: We really just need the id for this particular template
			opportunitydetailspage.StatusTableRowForm(uint(opptyId), s, fd).Render(r.Context(), w)
			return
		}
		buf := new(bytes.Buffer)
		u, ok := r.Context().Value(user.UserCtxKey).(*user.User)
		if !ok {
			logger.Error("No user in ctx")
			w.Write(buf.Bytes())
			return
		}
		oppty, err := opptyRepo.GetOpportuntyById(uint(opptyId), u)
		if err != nil {
			logger.Error("Problem getting opportunity by ID", "error", err)
			w.Write(buf.Bytes())
			return
		}
		opportunitydetailspage.StatusTable(oppty.ID, oppty.Statuses).Render(r.Context(), buf)
		opportunitydetailspage.OpptyDetailGrid(*oppty, true).Render(r.Context(), buf)
		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
	}
}

func statusFromRequest(r *http.Request) (status.Status, error) {
	var s status.Status

	err := r.ParseForm()
	if err != nil {
		return s, nil
	}
	s.Name = r.PostForm.Get("status-name")
	date := r.PostForm.Get("status-date")
	t, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return s, nil
	}
	s.Date = t
	s.Note = r.PostForm.Get("status-note")
	return s, nil
}
