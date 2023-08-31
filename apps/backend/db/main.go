package backend

import (
	"fmt"
	"log"
	"os"

	"github.com/pmwals09/yobs/apps/backend/opportunity"
	"github.com/pmwals09/yobs/apps/backend/task"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDb() (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel: logger.Info,
		},
	)
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return db, fmt.Errorf("Unable to open database: %w", err)
	}

	err = db.AutoMigrate(
		&opportunity.Opportunity{},
		&task.Task{},
	)
	if err != nil {
		return db, fmt.Errorf("unable to migrate database: %w", err)
	}
	return db, nil
}
