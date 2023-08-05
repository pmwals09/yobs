package opportunity

import (
	"fmt"
	"gorm.io/gorm"
	"net/http"
)

type OpportunityModel struct {
	Description string `json:"description"`
	URL         string `json:"url"`
}

// DTO that excludes the GORM Model but includes an ID
type OpportunityDTO struct {
	ID uint `json:"id"`
	OpportunityModel
}

func (o OpportunityDTO) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type Opportunity struct {
	gorm.Model
	*OpportunityModel
}

func NewOpportunity(description string, url string) *Opportunity {
	return &Opportunity{
		OpportunityModel: &OpportunityModel{
			Description: description, URL: url,
		},
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
		OpportunityModel: &OpportunityModel{
			Description: opp.Description,
			URL:         opp.URL,
		},
	}
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

func (g *GormRepository) UpdateOpporunity(opp *Opportunity) (*Opportunity, error) {
	res := g.DB.Save(&opp)
	return opp, res.Error
}

func (g *GormRepository) DeleteOpportunity(opptyId uint) error {
	res := g.DB.Delete(&Opportunity{}, opptyId)
	return res.Error
}
