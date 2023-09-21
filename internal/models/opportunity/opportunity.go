package opportunity

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/pmwals09/yobs/internal/models/document"
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
	ID              uint      `json:"id"`
	CompanyName     string    `json:"companyName"`
	Role            string    `json:"role"`
	Description     string    `json:"description"`
	URL             string    `json:"url"`
	ApplicationDate time.Time `json:"applicationDate"`
	Status          Status    `json:"status"`
	// Tasks           []task.Task `json:"tasks"`
	// Documents       []document.Document `json:"documents"`
	// Tasks
	// Contacts
}

func CreateTable(db *sql.DB) error {
	createStr := `
		CREATE TABLE IF NOT EXISTS opportunities (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			company_name TEXT,
			role TEXT,
			description TEXT,
			url TEXT,
			application_date DATETIME,
			status TEXT
		);
	`
	_, err := db.Exec(createStr)
	return err
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

func (o *Opportunity) WithApplicationDateString(applicationDate string) *Opportunity {
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

type Repository interface {
	CreateOpportunity(opp *Opportunity) error
	GetOpportuntyById(opptyId uint) (*Opportunity, error)
	GetAllOpportunities() ([]Opportunity, error)
	UpdateOpporunity(opptyId uint, newOpportunity Opportunity) error
	DeleteOpportunity(opptyId uint) error
	AddDocument(opptyId uint, documentId uint) error
	GetAllDocuments() ([]document.Document, error)
}

type OpportunityModel struct {
	DB *sql.DB
}

func (g *OpportunityModel) CreateOpportunity(opp *Opportunity) error {
	_, err := g.DB.Exec(`
		INSERT INTO opportunities (
			company_name,
			role,
			description,
			url,
			application_date,
			status
		) VALUES (?, ?, ?, ?, ?, ?);
	`,
		opp.CompanyName,
		opp.Role,
		opp.Description,
		opp.URL,
		opp.ApplicationDate,
		opp.Status,
	)

	return err
}

func (g *OpportunityModel) GetOpportuntyById(opptyId uint) (*Opportunity, error) {
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
		FROM opportunities WHERE id = ?;
	`, opptyId)
	err := res.Scan(
		&oppty.ID,
		&oppty.CompanyName,
		&oppty.Role,
		&oppty.Description,
		&oppty.URL,
		&oppty.ApplicationDate,
		&oppty.Status,
	)

	return &oppty, err
}

func (g *OpportunityModel) GetAllOpportunities() ([]Opportunity, error) {
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
		FROM opportunities;
	`)

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
		opptys = append(opptys, oppty)
	}
	return opptys, nil
}

func (g *OpportunityModel) UpdateOpporunity(opp *Opportunity) error {
	_, err := g.DB.Exec(`
		UPDATE opportunities
		SET
			company_name = ?,
			role = ?,
			description = ?,
			url = ?,
			application_date = ?,
			status = ?
		WHERE id = ?;
	`,
		opp.CompanyName,
		opp.Role,
		opp.Description,
		opp.URL,
		opp.ApplicationDate,
		opp.Status,
		opp.ID,
	)
	return err
}

func (g *OpportunityModel) DeleteOpportunity(opptyId uint) error {
	_, err := g.DB.Exec(`
		DELETE FROM opportunities WHERE id = ?
	`, opptyId)
	return err
}

func (g *OpportunityModel) AddDocument(opptyId uint, documentId uint) error {
	_, err := g.DB.Exec(`
		INSERT INTO opportunity_documents (
			opportunity_id,
			document_id
		) VALUES (?, ?);
	`, opptyId, documentId)
	return err
}

func (o *OpportunityModel) GetAllDocuments(opptyId uint) ([]document.Document, error) {
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
	`, opptyId)
	if err != nil {
		fmt.Println("QUERY ERROR",err.Error())
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
		docs = append(docs, d)
	}
	
	return docs, nil
}

// func (g *GormRepository) AddTask(opptyId uint, t []*task.Task) (*Opportunity, error) {
// 	if opp, err := g.GetOpportuntyById(opptyId); err != nil {
// 		return &opp, err
// 	} else {
// 		if appendErr := g.DB.Model(&opp).Association("Tasks").Append(t); appendErr != nil {
// 			return &opp, err
// 		} else {
// 			return &opp, nil
// 		}
// 	}
// }
