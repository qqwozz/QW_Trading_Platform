package main

import (
	"log"
	"net/http"
	"os"

	"github.com/qw-trading/platform/internal/db"
	"github.com/qw-trading/platform/internal/order/handler"
	"github.com/qw-trading/platform/internal/order/repository"
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
	tradeRepo := repository.NewTradeRepository(database)
	h := handler.New(repo, tradeRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/orders", h.CreateOrder)
	mux.HandleFunc("GET /v1/orders", h.ListOrders)
	mux.HandleFunc("GET /v1/orders/{id}", h.GetOrder)
	mux.HandleFunc("DELETE /v1/orders/{id}", h.CancelOrder)
	mux.HandleFunc("GET /health", db.HealthHandler(database, "order-service"))

	logger := applog.New("order-service")
	rl := middleware.NewRateLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst)
	authed := middleware.Auth(cfg.JWTSecret)(mux)
	wrapped := middleware.RequestID(middleware.Logger(logger)(rl.Middleware(middleware.CORS(cfg.AllowedOrigins)(authed))))

	server.Run("order-service", wrapped, server.DefaultConfig(), logger)
}
