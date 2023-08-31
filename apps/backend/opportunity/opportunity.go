package opportunity

import (
	"fmt"
	"net/http"

	"github.com/pmwals09/yobs/apps/backend/task"
	"gorm.io/gorm"
)

type Opportunity struct {
	ID          uint        `gorm:"primary_key" json:"id"`
	Description string      `json:"description"`
	URL         string      `json:"url"`
	Tasks       []task.Task `json:"tasks"`
	// Status
	// Tasks
	// Materials
	// Contacts
}

func (o Opportunity) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewOpportunity(description string, url string, tasks []task.Task) *Opportunity {
	return &Opportunity{
		Description: description, URL: url, Tasks: tasks,
	}
}

type Repository interface {
	CreateOpportunity(description string, url string) (Opportunity, error)
	GetOpportuntyById(opptyId uint) (Opportunity, error)
	GetAllOpportunities() ([]Opportunity, error)
	UpdateOpporunity(opptyId uint, newOpportunity Opportunity) error
	DeleteOpportunity(opptyId uint) error
}

type GormRepository struct {
	DB *gorm.DB
}

func (g *GormRepository) CreateOpportunity(opp *Opportunity) (Opportunity, error) {
	o := Opportunity{
		Description: opp.Description,
		URL:         opp.URL,
		Tasks:       opp.Tasks,
	}
	if result := g.DB.Create(&o); result.Error != nil {
		return o, fmt.Errorf("Error creating opportunity: %w", result.Error)
	}

	return o, nil
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
