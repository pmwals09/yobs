package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func InitDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		return nil, err
	}

	err = migrate(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

type MigrationPair struct {
	Up   string
	Down string
}

func migrate(db *sql.DB) error {
	migrationDirectory, err := getMigrationDirectory()
	if err != nil {
		return err
	}

	migrations, err := getMigrationPairs(migrationDirectory)
	if err != nil {
		return err
	}

	if err = hasMigrationTable(db); err != nil {
		return err
	}

	for k, v := range migrations {
		ok, err := hasRunMigration(db, k, v)
		if err != nil {
			return fmt.Errorf("Error querying for migration")
		}

		if ok {
			continue
		}

		upQBytes, err := os.ReadFile(migrationDirectory + "/" + v.Up)
		if err != nil {
			return fmt.Errorf("Unable to read migration file")
		}
		downQBytes, err := os.ReadFile(migrationDirectory + "/" + v.Down)
		if err != nil {
			return fmt.Errorf("Unable to read reversion migration file. Must have a valid reversion available")
		}

		_, err = db.Exec(string(upQBytes))
		if err != nil {
			_, err = db.Exec(string(downQBytes))
			if err != nil {
				return fmt.Errorf("Error migrating and reverting. Please check db integrity.")
			}
			return fmt.Errorf("Error migrating %s - reverted using %s", v.Up, v.Down)
		}

		db.Exec(`
      INSERT INTO migrations (
        migration_number,
        filename
      ) VALUES (?, ?);
    `, k, v.Up)
	}
	return nil
}

func assignMigrationFile(entry *MigrationPair, upOrDown string, val string) error {
	if upOrDown == "up" {
		entry.Up = val
		return nil
	} else if upOrDown == "down" {
		entry.Down = val
		return nil
	} else {
		return errors.New("Invalid migration type. Valid options are 'up' or 'down', preceding '.sql' in the file name")
	}
}

func getMigrationDirectory() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return wd + "/internal/db/migrations", nil
}

func getMigrationPairs(migrationDir string) (map[string]*MigrationPair, error) {
	migrations := map[string]*MigrationPair{}
	migrationFiles, err := os.ReadDir(migrationDir)
	if err != nil {
		return migrations, fmt.Errorf("No migration files available.")
	}
	if len(migrationFiles) == 0 {
		return migrations, fmt.Errorf("No migration files available.")
	}

	for _, migrationFile := range migrationFiles {
		fileInfo, err := migrationFile.Info()
		if err != nil {
			return migrations, fmt.Errorf("Unable to fetch migration file information.")
		}
		n, _, found := strings.Cut(fileInfo.Name(), "-")
		if !found {
			return migrations, fmt.Errorf("Malformed migration file name. Must be ###-description-name.up/down.sql")
		}

		upOrDown := strings.Split(fileInfo.Name(), ".")[1]
		if entry, ok := migrations[n]; !ok {
			mPair := MigrationPair{}
			migrations[n] = &mPair
			assignMigrationFile(&mPair, upOrDown, fileInfo.Name())
		} else {
			assignMigrationFile(entry, upOrDown, fileInfo.Name())
		}
	}

	return migrations, nil
}

func hasMigrationTable(db *sql.DB) error {
	row := db.QueryRow(`
    SELECT name FROM sqlite_master WHERE type = 'table' AND name = 'migrations';
  `)

	var name string
	err := row.Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			db.Exec(`
        CREATE TABLE IF NOT EXISTS migrations (
          id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
          migration_number TEXT,
          filename TEXT
        );
      `)
		} else {
			return fmt.Errorf("Error creating migration table")
		}
	}

	return nil
}

func hasRunMigration(db *sql.DB, k string, v *MigrationPair) (bool, error) {
	row := db.QueryRow(`
      SELECT
        migration_number,
        filename
      FROM migrations WHERE migration_number = ? AND filename = ?;
    `, k, v.Up)
	var migrationNumber string
	var filename string
	err := row.Scan(&migrationNumber, &filename)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
