package controllers

import (
	"bytes"
	"errors"
	"fmt"
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

func HandleStatusRowForm(statusRepo status.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statusId := chi.URLParam(r, "statusID")
		id, err := strconv.ParseUint(statusId, 10, 64)
		opportunityId := chi.URLParam(r, "opportunityId")
		opptyId, oErr := strconv.ParseUint(opportunityId, 10, 64)
		err = errors.Join(err, oErr)
		var fd helpers.FormData
		if err != nil {
			fd.AddError("overall", "Problem parsing status ID from path")
			var status status.Status
			opportunitydetailspage.StatusTableRowForm(uint(opptyId), status, fd).Render(r.Context(), w)
			return

		}
		status, err := statusRepo.GetStatusByID(uint(id))
		if err != nil {
			fd.AddError("overall", "Problem getting status by ID")
			opportunitydetailspage.StatusTableRowForm(uint(opptyId), status, fd).Render(r.Context(), w)
			return
		}
		fd.AddValue("status-name", status.Name)
		fd.AddValue("status-date", status.Date.Format(time.DateOnly))
		fd.AddValue("status-note", status.Note)

		opportunitydetailspage.StatusTableRowForm(uint(opptyId), status, fd).Render(r.Context(), w)
		return
	}
}

func HandleUpdateStatusItem(statusRepo status.Repository, opptyRepo opportunity.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var fd helpers.FormData

		// get the status id from the request
		statusId := chi.URLParam(r, "statusID")
		id, err := strconv.ParseUint(statusId, 10, 64)
		opportunityId := chi.URLParam(r, "opportunityId")
		opptyId, oErr := strconv.ParseUint(opportunityId, 10, 64)
		err = errors.Join(err, oErr)
		if err != nil {
			fd.AddError("overall", "Error parsing IDs from route")
			r.ParseForm()
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
			fd.AddError("overall", "Error parsing status data")
			r.ParseForm()
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
			fd.AddError("overall", "Error updating status in database")
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
		u, ok := r.Context().Value("user").(*user.User)
		if !ok {
			w.Write(buf.Bytes())
			return
		}
		oppty, err := opptyRepo.GetOpportuntyById(uint(opptyId), u)
		if err != nil {
			w.Write(buf.Bytes())
			return
		}
		opportunitydetailspage.StatusTable(oppty.ID, oppty.Statuses).Render(r.Context(), buf)
		opportunitydetailspage.OpptyDetailGrid(*oppty, true).Render(r.Context(), buf)
		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
		return
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
