package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	Name     string
}

func LoadDatabaseConfig() *DatabaseConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return &DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
}

func Connect() (*sql.DB, error) {
	config := LoadDatabaseConfig()
    connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",
        config.User, config.Password, config.Name, config.Host)
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, fmt.Errorf("error connecting to database: %v", err)
    }
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("error pinging database: %v", err)
    }

	RunMigrations(db, "./migrations")
	
    return db, nil
}

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
