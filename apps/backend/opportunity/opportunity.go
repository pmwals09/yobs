package opportunity

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pmwals09/yobs/apps/backend/task"
	"gorm.io/gorm"
)

type Opportunity struct {
	ID              uint        `gorm:"primary_key" json:"id"`
	Name            string      `json:"name"`
	Description     string      `json:"description"`
	URL             string      `json:"url"`
	Tasks           []task.Task `json:"tasks"`
	ApplicationDate time.Time   `json:"applicationDate"`
	// Status
	// Tasks
	// Materials
	// Contacts
}

func (o Opportunity) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewOpportunity() *Opportunity {
	return &Opportunity{}
}

func (o *Opportunity) WithName(name string) *Opportunity {
	o.Name = name
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
	fmt.Printf("\n%s\n", applicationDate)
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
	CreateOpportunity(opp *Opportunity) (Opportunity, error)
	GetOpportuntyById(opptyId uint) (Opportunity, error)
	GetAllOpportunities() ([]Opportunity, error)
	UpdateOpporunity(opptyId uint, newOpportunity Opportunity) error
	DeleteOpportunity(opptyId uint) error
}

type GormRepository struct {
	DB *gorm.DB
}

func (g *GormRepository) CreateOpportunity(opp *Opportunity) (Opportunity, error) {
	if result := g.DB.Create(&opp); result.Error != nil {
		return *opp, fmt.Errorf("Error creating opportunity: %w", result.Error)
	}

	return *opp, nil
}

func (g *GormRepository) GetOpportuntyById(opptyId uint) (Opportunity, error) {
	var oppty Opportunity
	err := g.DB.Model(&Opportunity{}).Preload("Tasks").First(&oppty, opptyId).Error
	return oppty, err
}

func (g *GormRepository) GetAllOpportunities() ([]Opportunity, error) {
	var opptys []Opportunity
	err := g.DB.Model(&Opportunity{}).Preload("Tasks").Find(&opptys).Error
	return opptys, err
}

func (g *GormRepository) UpdateOpporunity(opp *Opportunity) (*Opportunity, error) {
	res := g.DB.Save(&opp)
	return opp, res.Error
}

func (g *GormRepository) DeleteOpportunity(opptyId uint) error {
	res := g.DB.Select("Tasks").Delete(&Opportunity{ID: opptyId})
	return res.Error
}

func (g *GormRepository) AddTask(opptyId uint, t []*task.Task) (*Opportunity, error) {
	if opp, err := g.GetOpportuntyById(opptyId); err != nil {
		return &opp, err
	} else {
		if appendErr := g.DB.Model(&opp).Association("Tasks").Append(t); appendErr != nil {
			return &opp, err
		} else {
			return &opp, nil
		}
	}
}
