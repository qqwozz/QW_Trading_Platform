package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/qw-trading/platform/internal/db"
	"github.com/qw-trading/platform/internal/market/feeder"
	"github.com/qw-trading/platform/internal/market/handler"
	"github.com/qw-trading/platform/internal/market/hub"
	"github.com/qw-trading/platform/internal/market/repository"
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

	h := hub.NewHub()
	go h.Run()

	repo := repository.New(database)
	hdl := handler.New(repo, h)

	f := feeder.NewFeeder(repo)
	f.Start(context.Background())

	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/market/tickers", hdl.ListTickers)
	mux.HandleFunc("GET /v1/market/tickers/{symbol}", hdl.GetTicker)
	mux.HandleFunc("GET /v1/market/orderbook/{symbol}", hdl.GetOrderBook)
	mux.HandleFunc("GET /v1/market/ws", hdl.HandleWebSocket)
	mux.HandleFunc("GET /health", db.HealthHandler(database, "market-data-service"))

	logger := applog.New("market-data-service")
	rl := middleware.NewRateLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst)
	wrapped := middleware.RequestID(middleware.Logger(logger)(rl.Middleware(middleware.CORS(cfg.AllowedOrigins)(mux))))

	server.Run("market-data-service", wrapped, server.DefaultConfig(), logger)
}
