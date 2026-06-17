package main

import (
	"log"
	"net/http"
	"os"

	"github.com/qw-trading/platform/internal/db"
	"github.com/qw-trading/platform/internal/history/handler"
	"github.com/qw-trading/platform/internal/history/repository"
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
	mux.HandleFunc("GET /v1/history/orders", h.GetOrderHistory)
	mux.HandleFunc("GET /v1/history/trades", h.GetTradeHistory)
	mux.HandleFunc("GET /v1/history/balance", h.GetBalanceHistory)
	mux.HandleFunc("GET /v1/history/positions", h.GetPositionHistory)
	mux.HandleFunc("GET /health", db.HealthHandler(database, "history-service"))

	logger := applog.New("history-service")
	rl := middleware.NewRateLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst)
	wrapped := middleware.RequestID(middleware.Logger(logger)(rl.Middleware(middleware.CORS(cfg.AllowedOrigins)(mux))))

	server.Run("history-service", wrapped, server.DefaultConfig(), logger)
}
