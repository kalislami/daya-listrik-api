package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort  string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	CORSOrigins string
}

func Load() (Config, error) {
	_ = godotenv.Load()
	c := Config{
		ServerPort:  defaultValue("PORT", "8080"),
		DBHost:      strings.TrimSpace(os.Getenv("DB_HOST")),
		DBPort:      defaultValue("DB_PORT", "5432"),
		DBUser:      strings.TrimSpace(os.Getenv("DB_USER")),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBName:      strings.TrimSpace(os.Getenv("DB_NAME")),
		DBSSLMode:   defaultValue("DB_SSLMODE", "disable"),
		CORSOrigins: defaultValue("CORS_ALLOWED_ORIGINS", "http://localhost:5173"),
	}
	for name, value := range map[string]string{
		"DB_HOST": c.DBHost, "DB_PORT": c.DBPort, "DB_USER": c.DBUser,
		"DB_PASSWORD": c.DBPassword, "DB_NAME": c.DBName,
	} {
		if value == "" {
			return Config{}, fmt.Errorf("%s is required", name)
		}
	}
	return c, nil
}

func defaultValue(name, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value
	}
	return fallback
}
