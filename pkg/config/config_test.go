package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	os.Clearenv()
	cfg := Load()

	if cfg.DBHost != "localhost" {
		t.Errorf("DBHost = %q, want %q", cfg.DBHost, "localhost")
	}
	if cfg.DBPort != 5432 {
		t.Errorf("DBPort = %d, want %d", cfg.DBPort, 5432)
	}
	if cfg.DBUser != "postgres" {
		t.Errorf("DBUser = %q, want %q", cfg.DBUser, "postgres")
	}
	if cfg.DBName != "qw_trading" {
		t.Errorf("DBName = %q, want %q", cfg.DBName, "qw_trading")
	}
	if cfg.Port != "8080" {
		t.Errorf("Port = %q, want %q", cfg.Port, "8080")
	}
	if cfg.Env != "development" {
		t.Errorf("Env = %q, want %q", cfg.Env, "development")
	}
}

func TestLoad_EnvOverrides(t *testing.T) {
	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("PORT", "9090")
	os.Setenv("APP_ENV", "production")
	defer os.Clearenv()

	cfg := Load()

	if cfg.DBHost != "db.example.com" {
		t.Errorf("DBHost = %q, want %q", cfg.DBHost, "db.example.com")
	}
	if cfg.DBPort != 5433 {
		t.Errorf("DBPort = %d, want %d", cfg.DBPort, 5433)
	}
	if cfg.Port != "9090" {
		t.Errorf("Port = %q, want %q", cfg.Port, "9090")
	}
	if cfg.Env != "production" {
		t.Errorf("Env = %q, want %q", cfg.Env, "production")
	}
}

func TestDatabaseDSN(t *testing.T) {
	cfg := &Config{
		DBHost:     "localhost",
		DBPort:     5432,
		DBUser:     "postgres",
		DBPassword: "secret",
		DBName:     "qw_trading",
		DBSSLMode:  "disable",
	}

	dsn := cfg.DatabaseDSN()
	expected := "host=localhost port=5432 user=postgres password=secret dbname=qw_trading sslmode=disable"
	if dsn != expected {
		t.Errorf("DatabaseDSN() = %q, want %q", dsn, expected)
	}
}

func TestGetEnvInt_Invalid(t *testing.T) {
	os.Setenv("TEST_INVALID_INT", "not-a-number")
	defer os.Unsetenv("TEST_INVALID_INT")

	result := getEnvInt("TEST_INVALID_INT", 42)
	if result != 42 {
		t.Errorf("getEnvInt() = %d, want 42 (fallback)", result)
	}
}

func TestGetEnvInt_Empty(t *testing.T) {
	result := getEnvInt("NONEXISTENT_KEY", 99)
	if result != 99 {
		t.Errorf("getEnvInt() = %d, want 99 (fallback)", result)
	}
}

func TestGetEnvFloat_Defaults(t *testing.T) {
	os.Clearenv()
	result := getEnvFloat("NONEXISTENT_FLOAT", 3.14)
	if result != 3.14 {
		t.Errorf("getEnvFloat() = %f, want 3.14", result)
	}
}

func TestGetEnvFloat_Parses(t *testing.T) {
	os.Setenv("TEST_FLOAT", "2.5")
	defer os.Unsetenv("TEST_FLOAT")

	result := getEnvFloat("TEST_FLOAT", 1.0)
	if result != 2.5 {
		t.Errorf("getEnvFloat() = %f, want 2.5", result)
	}
}
