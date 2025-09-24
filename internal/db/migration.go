package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// RunMigrations menjalankan semua migrasi dari folder yang ditentukan
func RunMigrations(db *sql.DB, migrationsDir string) error {
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("error reading migrations directory: %v", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			migrationFile := filepath.Join(migrationsDir, file.Name())
			err := executeMigration(db, migrationFile)
			if err != nil {
				return fmt.Errorf("error executing migration %s: %v", migrationFile, err)
			}
		}
	}

	return nil
}

// executeMigration membaca dan mengeksekusi query dari file migrasi
func executeMigration(db *sql.DB, migrationFile string) error {
	query, err := os.ReadFile(migrationFile)
	if err != nil {
		return fmt.Errorf("error reading migration file %s: %v", migrationFile, err)
	}

	_, err = db.Exec(string(query))
	if err != nil {
		return fmt.Errorf("error executing migration: %v", err)
	}

	log.Printf("Migration %s executed successfully\n", migrationFile)
	return nil
}
