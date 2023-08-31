package task

import (
	"net/http"
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID            uint      `json:"id" gorm:"primary_key"`
	Description   string    `json:"description"`
	DueDate       time.Time `json:"dueDate"`
	OpportunityID uint      `json:"opportunityId"`
	// Status
}

func (t Task) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewTask(description string, dueDate time.Time, opp uint) *Task {
	return &Task{
		Description: description, DueDate: dueDate, OpportunityID: opp,
	}
}

type Repository interface {
	// get all tasks associated with an opportunity
	// get a single task by id
	// update a task by id
	// delete a task
}

type GormRepository struct {
	DB *gorm.DB
}
