package main

import (
	"log"
	"net/http"
	"os"

	"github.com/qw-trading/platform/internal/account/handler"
	"github.com/qw-trading/platform/internal/account/repository"
	"github.com/qw-trading/platform/internal/db"
	"github.com/qw-trading/platform/pkg/config"
	applog "github.com/qw-trading/platform/pkg/logger"
	"github.com/qw-trading/platform/pkg/middleware"
	"github.com/qw-trading/platform/pkg/server"
)

func main() {
	cfg := config.Load()
	os.Setenv("PORT", cfg.Port)

	database, err := db.Connect(cfg.DatabaseDSN(), cfg.DBMaxOpen, cfg.DBMaxIdle)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	repo := repository.New(database)
	h := handler.New(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /accounts", h.ListAccounts)
	mux.HandleFunc("POST /accounts/deposit", h.Deposit)
	mux.HandleFunc("GET /health", db.HealthHandler(database, "account-service"))

	logger := applog.New("account-service")
	rl := middleware.NewRateLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst)
	wrapped := middleware.RequestID(middleware.Logger(logger)(rl.Middleware(middleware.CORS(cfg.AllowedOrigins)(mux))))

	server.Run("account-service", wrapped, server.DefaultConfig(), logger)
}
