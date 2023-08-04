package opportunity

import (
	"fmt"

	"gorm.io/gorm"
)

type Opportunity struct {
	gorm.Model
	Description string
	URL         string
}

func NewOpportunity(description string, url string) *Opportunity {
	return &Opportunity{Description: description, URL: url}
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
	o := Opportunity{Description: opp.Description, URL: opp.URL}
	if result := g.DB.Create(&o); result.Error != nil {
		return o, fmt.Errorf("Error creating opportunity: %w", result.Error)
	}
	return o, nil
}

func (g *GormRepository) GetOpportuntyById(opptyId uint) (Opportunity, error) {
	oppty := Opportunity{}
	res := g.DB.First(&oppty, opptyId)
	return oppty, res.Error
}

func (g *GormRepository) GetAllOpportunities() ([]Opportunity, error) {
	opptys := []Opportunity{}
	res := g.DB.Find(&opptys)
	return opptys, res.Error
}

// func (g *GormRepository) UpdateOpporunity(opptyId uint, newOpportunity OpportunityDTO) error {
// }
//
// func (g *GormRepository) DeleteOpportunity(opptyId uint) error {
// }
