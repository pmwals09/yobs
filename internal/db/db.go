package db

import (
	"database/sql"
	"fmt"

	"github.com/pmwals09/yobs/internal/config"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func InitDb(cfg config.Config) (*sql.DB, error) {
	var path string
	switch cfg.EnvironmentName {
	case config.LocalEnvironment:
		path = "file:./test.db"
	case config.ProductionEnvironment:
		path = fmt.Sprintf("%s?authToken=%s", cfg.DBConfig.URL, cfg.DBConfig.Token)
	}

	db, err := sql.Open("libsql", path)
	if err != nil {
		return nil, err
	}

	return db, nil
}
