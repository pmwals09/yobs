package backend

import (
	"fmt"
	"github.com/pmwals09/yobs/apps/backend/opportunity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDb() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return db, fmt.Errorf("Unable to open database: %w", err)
	}
	err = db.AutoMigrate(&opportunity.Opportunity{})
	if err != nil {
		return db, fmt.Errorf("unable to migrate database: %w", err)
	}
	return db, nil
}
