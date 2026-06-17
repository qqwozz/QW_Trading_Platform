// Package config provides application configuration management by loading
// settings from environment variables with sensible defaults.
package config

import (
	"os"
	"strconv"
)

// Config holds all application configuration values loaded from environment variables.
type Config struct {
	// Database connection settings
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	DBMaxOpen  int
	DBMaxIdle  int

	// JWT authentication settings
	JWTSecret     string
	JWTExpiry     int
	RefreshExpiry int

	// Server settings
	Port           string
	GRPCPort       string
	Env            string
	AllowedOrigins string

	// Rate limiting
	RateLimitRPS   float64
	RateLimitBurst int

	// Trading fees (basis points, 1 bp = 0.01%)
	MakerFeeBPS int
	TakerFeeBPS int
}

// Load reads configuration from environment variables and returns a Config
// with fallback defaults for missing values.
func Load() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvInt("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "qw_trading"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		DBMaxOpen:  getEnvInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdle:  getEnvInt("DB_MAX_IDLE_CONNS", 5),

		JWTSecret:     getEnv("JWT_SECRET", ""),
		JWTExpiry:     getEnvInt("JWT_EXPIRY_HOURS", 1),
		RefreshExpiry: getEnvInt("REFRESH_EXPIRY_HOURS", 168),

		Port:           getEnv("PORT", "8080"),
		GRPCPort:       getEnv("GRPC_PORT", "9090"),
		Env:            getEnv("APP_ENV", "development"),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "*"),
		RateLimitRPS:   getEnvFloat("RATE_LIMIT_RPS", 100),
		RateLimitBurst: getEnvInt("RATE_LIMIT_BURST", 200),
		MakerFeeBPS:    getEnvInt("MAKER_FEE_BPS", 10),
		TakerFeeBPS:    getEnvInt("TAKER_FEE_BPS", 20),
	}
}

// DatabaseDSN constructs a PostgreSQL connection string from the config fields.
func (c *Config) DatabaseDSN() string {
	return "host=" + c.DBHost +
		" port=" + strconv.Itoa(c.DBPort) +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" sslmode=" + c.DBSSLMode
}

// getEnv returns the value of an environment variable, or fallback if unset/empty.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// getEnvInt parses an environment variable as an integer, returning fallback on
// any error (missing, empty, or non-numeric value).
func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvFloat(key string, fallback float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return fallback
}
