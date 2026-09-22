package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/url"
	"time"

	"daya-listrik-api/internal/config"
	"daya-listrik-api/migrations"

	"github.com/golang-migrate/migrate/v4"
	postgresmigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
)

func Connect(ctx context.Context, c config.Config) (*sql.DB, error) {
	dsn := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.DBUser, c.DBPassword),
		Host:   net.JoinHostPort(c.DBHost, c.DBPort),
		Path:   c.DBName,
	}
	query := dsn.Query()
	query.Set("sslmode", c.DBSSLMode)
	dsn.RawQuery = query.Encode()

	db, err := sql.Open("postgres", dsn.String())
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return db, nil
}

func RunMigrations(db *sql.DB) error {
	source, err := iofs.New(migrations.Files, ".")
	if err != nil {
		return fmt.Errorf("open migrations: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("get migration connection: %w", err)
	}
	driver, err := postgresmigrate.WithConnection(ctx, conn, &postgresmigrate.Config{StatementTimeout: 30 * time.Second})
	if err != nil {
		conn.Close()
		return fmt.Errorf("open migration database: %w", err)
	}
	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		driver.Close()
		return fmt.Errorf("initialize migrations: %w", err)
	}
	applyErr := m.Up()
	sourceErr, databaseErr := m.Close()
	if applyErr != nil && !errors.Is(applyErr, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations: %w", applyErr)
	}
	if sourceErr != nil || databaseErr != nil {
		return fmt.Errorf("close migrations: source: %v, database: %v", sourceErr, databaseErr)
	}
	return nil
}
