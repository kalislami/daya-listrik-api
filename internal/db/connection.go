package db

import (
	"daya-listrik-api/internal/models"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func LoadConnectionGorm() (*gorm.DB, error) {
	// load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	driver := os.Getenv("DB_DRIVER")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	var dsn string
	var dialector gorm.Dialector

	switch driver {
	case "postgres":
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, pass, name,
		)
		dialector = postgres.Open(dsn)

	case "mysql":
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?parseTime=true",
			user, pass, host, port, name,
		)
		dialector = mysql.Open(dsn)

	case "sqlite":
		dsn = name // biasanya nama file, ex: mydb.sqlite
		dialector = sqlite.Open(dsn)

	default:
		return nil, fmt.Errorf("unsupported driver: %s", driver)
	}

	// open GORM connection
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// ambil *sql.DB kalau mau pakai untuk ping atau migration manual
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// jalankan migration otomatis hanya di parent process
	if !fiber.IsChild() {
		// daftar model yang mau di-migrate
		err = db.AutoMigrate(
			&models.EnergyRecord{},
			// tambahkan model lain di sini
		)
		if err != nil {
			log.Fatalf("Migration failed: %v", err)
		}

		log.Println("Migration success")
	}

	return db, nil
}
