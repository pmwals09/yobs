package opportunity

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pmwals09/yobs/internal/models/contact"
	"github.com/pmwals09/yobs/internal/models/document"
	"github.com/pmwals09/yobs/internal/models/status"
	"github.com/pmwals09/yobs/internal/models/user"
	// "github.com/pmwals09/yobs/apps/backend/task"
)

type Opportunity struct {
	ID          uint            `json:"id"`
	CompanyName string          `json:"companyName"`
	Role        string          `json:"role"`
	Description string          `json:"description"`
	URL         string          `json:"url"`
	Statuses    []status.Status `json:"statuses"`
	User        *user.User      `json:"user"`
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
	AddContact(oppty *Opportunity, contact contact.Contact) error
	GetAllContacts(oppty *Opportunity) ([]contact.Contact, error)
	UpdateStatus(oppty *Opportunity, status status.Status) error
	RemoveDocument(oppty *Opportunity, doc document.Document) error
}

type OpportunityModel struct {
	DB *sql.DB
}

type NullableStatus struct {
	ID   sql.NullInt64
	Name sql.NullString
	Note sql.NullString
	Date sql.NullTime
}

func (n NullableStatus) ToStatus() (status.Status, bool) {
	var out status.Status
	if !n.ID.Valid {
		return out, false
	}
	out.ID = uint(n.ID.Int64)
	out.Name = n.Name.String
	out.Note = n.Note.String
	out.Date = n.Date.Time
	return out, true
}

// Create an opportunity and it's associated initial status entry
func (g *OpportunityModel) CreateOpportunity(opp *Opportunity) error {
	if opp.IsEmpty() {
		return errors.New("Empty opportunity - must have at least one of Role, Company Name, or URL")
	}

	ctx := context.Background()
	tx, err := g.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	res, err := tx.Exec(
		`
		INSERT INTO opportunities (
			company_name,
			role,
			description,
			url,
			user_id
		) VALUES (?, ?, ?, ?, ?, ?);
	`,
		opp.CompanyName,
		opp.Role,
		opp.Description,
		opp.URL,
		opp.User.ID)
	if err != nil {
		return txError(tx, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return txError(tx, err)
	}
	_, err = tx.Exec(`
		INSERT INTO statuses (
			name,
			note,
			date,
			opportunity_id
		) VALUES (?, ?, ?, ?);
`,
		opp.Statuses[0].Name,
		opp.Statuses[0].Note,
		opp.Statuses[0].Date,
		id)
	err = tx.Commit()
	if err != nil {
		return txError(tx, err)
	}

	return err
}

func txError(tx *sql.Tx, err error) error {
	rbErr := tx.Rollback()
	if rbErr != nil {
		err = errors.Join(err, rbErr)
		return err
	}
	return err
}

func (g *OpportunityModel) GetOpportuntyById(opptyId uint, user *user.User) (*Opportunity, error) {
	var oppty Opportunity
	res, err := g.DB.Query(`
		SELECT
			o.id,
			company_name,
			role,
			description,
			url,
			s.id,
			name,
			note,
			date
		FROM opportunities o FULL OUTER JOIN statuses s ON o.id = s.opportunity_id
		WHERE o.id = ? AND user_id = ?
		ORDER BY date DESC;
	`, opptyId, user.ID)
	if err != nil {
		return &oppty, err
	}
	for res.Next() {
		var nStatus NullableStatus
		err := res.Scan(
			&oppty.ID,
			&oppty.CompanyName,
			&oppty.Role,
			&oppty.Description,
			&oppty.URL,
			&nStatus.ID,
			&nStatus.Name,
			&nStatus.Note,
			&nStatus.Date)
		if err != nil {
			return &oppty, err
		}
		if status, ok := nStatus.ToStatus(); ok {
			oppty.Statuses = append(oppty.Statuses, status)
		}
	}
	oppty.User = user

	return &oppty, nil
}

func (g *OpportunityModel) GetAllOpportunities(user *user.User) ([]Opportunity, error) {
	var opptys []Opportunity
	rows, err := g.DB.Query(`
		SELECT
			o.id,
			company_name,
			role,
			description,
			url,
			s.id,
			name,
			note,
			date
		FROM opportunities o FULL OUTER JOIN statuses s ON o.id = s.opportunity_id
		WHERE user_id = ?
		ORDER BY date DESC;
	`, user.ID)

	if err != nil {
		return opptys, err
	}
	opptyMap := make(map[uint]Opportunity)
	for rows.Next() {
		var oppty Opportunity
		var nStatus NullableStatus
		err := rows.Scan(
			&oppty.ID,
			&oppty.CompanyName,
			&oppty.Role,
			&oppty.Description,
			&oppty.URL,
			&nStatus.ID,
			&nStatus.Name,
			&nStatus.Note,
			&nStatus.Date)
		if err != nil {
			return opptys, err
		}
		if val, ok := opptyMap[oppty.ID]; ok {
			if status, sOk := nStatus.ToStatus(); sOk {
				val.Statuses = append(val.Statuses, status)
			}
			opptyMap[oppty.ID] = val
		} else {
			if status, sOk := nStatus.ToStatus(); sOk {
				oppty.Statuses = append(oppty.Statuses, status)
			}
			oppty.User = user
			opptyMap[oppty.ID] = oppty
		}
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
			url = ?
		WHERE id = ? AND user_id = ?;
	`,
		opp.CompanyName,
		opp.Role,
		opp.Description,
		opp.URL,
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
			d.id,
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
		return docs, err
	}

	for rows.Next() {
		var d document.Document
		err := rows.Scan(
			&d.ID,
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

func (g *OpportunityModel) AddContact(oppty *Opportunity, contact contact.Contact) error {
	_, err := g.DB.Exec(`
		INSERT INTO opportunity_contacts (
			opportunity_id,
			contact_id
		) VALUES (?, ?);
	`, oppty.ID, contact.ID)
	return err
}

func (g *OpportunityModel) GetAllContacts(oppty *Opportunity) ([]contact.Contact, error) {
	var contacts []contact.Contact
	rows, err := g.DB.Query(`
		SELECT
			c.name,
			c.company_name,
			c.title,
			c.phone,
			c.email
		FROM contacts c
		JOIN opportunity_contacts oc ON c.id = oc.contact_id
		JOIN opportunities o ON o.id = oc.opportunity_id
		WHERE o.id = ?;
	`, oppty.ID)
	if err != nil {
		return contacts, err
	}

	for rows.Next() {
		var c contact.Contact
		err := rows.Scan(
			&c.Name,
			&c.CompanyName,
			&c.Title,
			&c.Phone,
			&c.Email)
		if err != nil {
			return contacts, err
		}
		contacts = append(contacts, c)
	}
	return contacts, nil
}

func (g *OpportunityModel) UpdateStatus(oppty *Opportunity, status status.Status) error {
	_, err := g.DB.Exec(`
		INSERT INTO statuses (
			name,
			note,
			date,
			opportunity_id
		) VALUES ( ?, ?, ?, ? );
`,
		status.Name,
		status.Note,
		status.Date,
		oppty.ID)
	return err
}

func (g *OpportunityModel) RemoveDocument(oppty *Opportunity, doc document.Document) error {
	_, err := g.DB.Exec(`
		DELETE FROM opportunity_documents
		WHERE opportunity_id = ? AND document_id = ?;
	`, oppty.ID, doc.ID)
	return err
}
