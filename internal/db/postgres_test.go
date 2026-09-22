package db

import (
	"context"
	"daya-listrik-api/internal/config"
	"strings"
	"testing"
)

func TestConnectFailsWhenDatabaseUnavailable(t *testing.T) {
	_, err := Connect(context.Background(), config.Config{
		DBHost: "127.0.0.1", DBPort: "1", DBUser: "user", DBPassword: "secret", DBName: "records", DBSSLMode: "disable",
	})
	if err == nil || !strings.Contains(err.Error(), "ping database") {
		t.Fatalf("expected ping failure, got %v", err)
	}
	if strings.Contains(err.Error(), "secret") {
		t.Fatal("connection error exposed password")
	}
}
