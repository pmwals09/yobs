package backend

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/pmwals09/yobs/apps/backend/opportunity"
	"gorm.io/gorm"
)

// DTO that excludes the GORM Model but includes an ID
type OpportunityDTO struct {
	ID          uint   `json:"id"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

func (o OpportunityDTO) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// ----- Request types -----

type OpportunityRequest struct {
	*OpportunityDTO
}

func (o *OpportunityRequest) Bind(r *http.Request) error {
	if o == nil {
		return errors.New("Missing required opportunity fields")
	}
	return nil
}

// Response types

type PongRes struct {
	Msg string `json:"msg"`
}

func (pr PongRes) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type ErrResponse struct {
	Err       error  `json:"error"`
	Code      int    `json:"code"`
	ErrorText string `json:"errorText"`
}

func (e ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Sets the render status to the provided code and returns and builds the
// ErrResponse
func ErrRender(r *http.Request, err error, code int) render.Renderer {
	render.Status(r, code)
	return &ErrResponse{
		Err:       err,
		Code:      code,
		ErrorText: err.Error(),
	}
}

type OpportunityResponse struct {
	*opportunity.Opportunity
}

func (o *OpportunityResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// API Router
func ApiRouter(db *gorm.DB) http.Handler {
	r := chi.NewRouter()

	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Get("/ping", handlePing)

	opptyRepo := opportunity.GormRepository{DB: db}
	r.Route("/opportunities", func(r chi.Router) {
		r.Post("/", handlePostOppty(opptyRepo))
		r.Get("/", handleGetAllOppty(opptyRepo))
		r.Route("/{opportunityId}", func(r chi.Router) {
			r.Get("/", handleGetOppty(opptyRepo))
			r.Put("/", handleUpdateOppty(opptyRepo))
			r.Delete("/", handleDeleteOppty(opptyRepo))
		})
	})
	return r
}

// Endpoint Handlers
func handlePing(w http.ResponseWriter, r *http.Request) {
	if err := render.Render(w, r, PongRes{"pong"}); err != nil {
		render.Render(w, r, ErrRender(r, err, 500))
		return
	}
}

func handlePostOppty(repo opportunity.GormRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &OpportunityRequest{}
		if err := render.Bind(r, data); err != nil {
			render.Render(w, r, ErrRender(r, err, 400))
			return
		}

		newOpportunity := opportunity.NewOpportunity(
			data.Description,
			data.URL,
		)

		if createdOppty, err := repo.CreateOpportunity(newOpportunity); err != nil {
			render.Render(w, r, ErrRender(r, err, 500))
			return
		} else {
			res := OpportunityDTO{
				ID:          createdOppty.ID,
				Description: createdOppty.Description,
				URL:         createdOppty.URL,
			}
			render.Render(w, r, res)
		}
	}
}

func handleGetAllOppty(repo opportunity.GormRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if opptys, err := repo.GetAllOpportunities(); err != nil {
			render.Render(w, r, ErrRender(r, err, 500))
			return
		} else {
			res := []render.Renderer{}
			for _, o := range opptys {
				res = append(res, OpportunityDTO{
					ID:          o.ID,
					Description: o.Description,
					URL:         o.URL,
				})
			}
			render.RenderList(w, r, res)
		}
	}
}
func handleGetOppty(repo opportunity.GormRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "opportunityId")
		if id, err := strconv.ParseUint(idParam, 10, 64); err != nil {
			render.Render(w, r, ErrRender(r, err, 401))
			return
		} else {
			if oppty, err := repo.GetOpportuntyById(uint(id)); err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					render.Render(w, r, ErrRender(r, err, 404))
					return
				}
				render.Render(w, r, ErrRender(r, err, 500))
				return
			} else {
				res := OpportunityResponse{
					&oppty,
				}
				render.Render(w, r, &res)
			}
		}
	}
}

func handleUpdateOppty(repo opportunity.GormRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the new details from the request body
		data := &OpportunityRequest{}
		if err := render.Bind(r, data); err != nil {
			render.Render(w, r, ErrRender(r, err, 400))
			return
		}

		updatedOppty := OpportunityDTO{
			ID:          data.ID,
			Description: data.Description,
			URL:         data.URL,
		}
		fmt.Println(updatedOppty)
		// pass them to the repository to db.Save per GORM
		// https://gorm.io/docs/update.html
		// ? How to get a partial update?
		//   - Would need additional logic around what keys exist from user
		//     provided JSON to handle in the Bind/marshalling/DTO model
	}
}
func handleDeleteOppty(repo opportunity.GormRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
