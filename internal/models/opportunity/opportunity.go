package opportunity

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/pmwals09/yobs/internal/models/document"
	"github.com/pmwals09/yobs/internal/models/user"
	// "github.com/pmwals09/yobs/apps/backend/task"
)

type Status string

const (
	None          Status = "None"
	Applied              = "Applied"
	Rejected             = "Rejected"
	FollowedUp           = "Followed Up"
	Pending              = "Pending"
	Offer                = "Offer"
	AcceptedOffer        = "Accepted Offer"
)

type Opportunity struct {
	ID              uint       `json:"id"`
	CompanyName     string     `json:"companyName"`
	Role            string     `json:"role"`
	Description     string     `json:"description"`
	URL             string     `json:"url"`
	ApplicationDate time.Time  `json:"applicationDate"`
	Status          Status     `json:"status"`
	User            *user.User `json:"user"`
	// Tasks           []task.Task `json:"tasks"`
	// Documents       []document.Document `json:"documents"`
	// Tasks
	// Contacts
}

func New() *Opportunity {
	return &Opportunity{}
}

func (o *Opportunity) WithCompanyName(name string) *Opportunity {
	o.CompanyName = name
	return o
}

func (o *Opportunity) WithRole(role string) *Opportunity {
	o.Role = role
	return o
}

func (o *Opportunity) WithDescription(description string) *Opportunity {
	o.Description = description
	return o
}

func (o *Opportunity) WithURL(url string) *Opportunity {
	o.URL = url
	return o
}

func (o *Opportunity) WithUser(user *user.User) *Opportunity {
	o.User = user
  return o
}

func (o *Opportunity) WithApplicationDateString(applicationDate string) *Opportunity {
  if applicationDate == "" {
    o.ApplicationDate = time.Time{}
    return o
  }
	t, err := time.Parse("2006-01-02", applicationDate)
	if err != nil {
		fmt.Printf("\nError parsing date: %s\n", err.Error())
		return o
	}
	o.ApplicationDate = t
	return o
}

func (o *Opportunity) WithApplicationDateTime(applicationDate time.Time) *Opportunity {
	o.ApplicationDate = applicationDate
	return o
}

func (o *Opportunity) IsEmpty() bool {
  return o.CompanyName == "" && o.URL == "" && o.Role == ""
}

type Repository interface {
	CreateOpportunity(opp *Opportunity) error
	GetOpportuntyById(opptyId uint, user *user.User) (*Opportunity, error)
	GetAllOpportunities(user *user.User) ([]Opportunity, error)
	UpdateOpportunity(opp *Opportunity) error
	DeleteOpportunity(oppty *Opportunity) error
	AddDocument(oppty *Opportunity, document *document.Document) error
	GetAllDocuments(oppty *Opportunity, user *user.User) ([]document.Document, error)
}

type OpportunityModel struct {
	DB *sql.DB
}

func (g *OpportunityModel) CreateOpportunity(opp *Opportunity) error {
  if opp.IsEmpty() {
    return errors.New("Empty opportunity - must have at least one of Role, Company Name, or URL")
  }

	_, err := g.DB.Exec(`
		INSERT INTO opportunities (
			company_name,
			role,
			description,
			url,
			application_date,
			status,
      user_id
		) VALUES (?, ?, ?, ?, ?, ?, ?);
	`,
		opp.CompanyName,
		opp.Role,
		opp.Description,
		opp.URL,
		opp.ApplicationDate,
		opp.Status,
    opp.User.ID,
	)

	return err
}

func (g *OpportunityModel) GetOpportuntyById(opptyId uint, user *user.User) (*Opportunity, error) {
	var oppty Opportunity
	res := g.DB.QueryRow(`
		SELECT
			id,
			company_name,
			role,
			description,
			url,
			application_date,
			status
		FROM opportunities WHERE id = ? AND user_id = ?;
	`, opptyId, user.ID)
	err := res.Scan(
		&oppty.ID,
		&oppty.CompanyName,
		&oppty.Role,
		&oppty.Description,
		&oppty.URL,
		&oppty.ApplicationDate,
		&oppty.Status,
	)
  oppty.User = user

	return &oppty, err
}

func (g *OpportunityModel) GetAllOpportunities(user *user.User) ([]Opportunity, error) {
	var opptys []Opportunity
	rows, err := g.DB.Query(`
		SELECT
			id,
			company_name,
			role,
			description,
			url,
			application_date,
			status
		FROM opportunities WHERE user_id = ?;
	`, user.ID)

	if err != nil {
		return opptys, err
	}
	for rows.Next() {
		oppty := Opportunity{}
		err := rows.Scan(
			&oppty.ID,
			&oppty.CompanyName,
			&oppty.Role,
			&oppty.Description,
			&oppty.URL,
			&oppty.ApplicationDate,
			&oppty.Status,
		)
		if err != nil {
			return opptys, err
		}
    oppty.User = user
		opptys = append(opptys, oppty)
	}
	return opptys, nil
}

func (g *OpportunityModel) UpdateOpportunity(opp *Opportunity) error {
	_, err := g.DB.Exec(`
		UPDATE opportunities
		SET
			company_name = ?,
			role = ?,
			description = ?,
			url = ?,
			application_date = ?,
			status = ?
		WHERE id = ? AND user_id = ?;
	`,
		opp.CompanyName,
		opp.Role,
		opp.Description,
		opp.URL,
		opp.ApplicationDate,
		opp.Status,
		opp.ID,
    opp.User.ID,
	)
	return err
}

func (g *OpportunityModel) DeleteOpportunity(oppty *Opportunity) error {
	_, err := g.DB.Exec(`
		DELETE FROM opportunities WHERE id = ? AND user_id = ?
	`, oppty.ID, oppty.User.ID)
	return err
}

func (g *OpportunityModel) AddDocument(oppty *Opportunity, document *document.Document) error {
  // TODO: Only add it if you haven't already
	_, err := g.DB.Exec(`
		INSERT INTO opportunity_documents (
			opportunity_id,
			document_id
		) VALUES (?, ?);
	`, oppty.ID, document.ID)
	return err
}

func (o *OpportunityModel) GetAllDocuments(oppty *Opportunity, user *user.User) ([]document.Document, error) {
	var docs []document.Document
	rows, err := o.DB.Query(`
		SELECT
			file_name,
			title,
			type,
			content_type
		FROM documents d
		JOIN opportunity_documents od ON d.id = od.document_id
		JOIN opportunities o ON o.id = od.opportunity_id
		WHERE o.id = ?;
	`, oppty.ID)
	if err != nil {
		fmt.Println("QUERY ERROR", err.Error())
		return docs, err
	}

	for rows.Next() {
		var d document.Document
		err := rows.Scan(
			&d.FileName,
			&d.Title,
			&d.Type,
			&d.ContentType,
		)
		if err != nil {
			return docs, err
		}
    d.User = user
		docs = append(docs, d)
	}

	return docs, nil
}
