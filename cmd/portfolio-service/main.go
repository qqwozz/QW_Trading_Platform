package main

import (
	"log"
	"net/http"
	"os"

	"github.com/qw-trading/platform/internal/db"
	"github.com/qw-trading/platform/internal/portfolio/handler"
	"github.com/qw-trading/platform/internal/portfolio/repository"
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

	posRepo := repository.New(database)
	accRepo := repository.NewAccountRepository(database)
	h := handler.New(posRepo, accRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/portfolio", h.GetPortfolio)
	mux.HandleFunc("GET /v1/positions", h.ListPositions)
	mux.HandleFunc("POST /v1/positions", h.UpdatePosition)
	mux.HandleFunc("GET /v1/balances", h.GetBalances)
	mux.HandleFunc("GET /health", db.HealthHandler(database, "portfolio-service"))

	logger := applog.New("portfolio-service")
	rl := middleware.NewRateLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst)
	authed := middleware.Auth(cfg.JWTSecret)(mux)
	wrapped := middleware.RequestID(middleware.Logger(logger)(rl.Middleware(middleware.CORS(cfg.AllowedOrigins)(authed))))

	server.Run("portfolio-service", wrapped, server.DefaultConfig(), logger)
}
