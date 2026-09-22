package config

import (
	"os"
	"testing"
)

func TestLoadWithoutDotEnv(t *testing.T) {
	tmp := t.TempDir()
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(old) })
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "6543")
	t.Setenv("DB_USER", "tester")
	t.Setenv("DB_PASSWORD", "test-only")
	t.Setenv("DB_NAME", "records")
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if c.DBPort != "6543" || c.DBHost != "localhost" || c.ServerPort == "" {
		t.Fatalf("unexpected config: %+v", c)
	}
}

func TestLoadRequiresDatabaseSettings(t *testing.T) {
	t.Setenv("DB_HOST", "")
	if _, err := Load(); err == nil {
		t.Fatal("expected missing DB_HOST error")
	}
}
