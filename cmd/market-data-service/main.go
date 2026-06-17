// Command market-data-service provides real-time market data including tickers,
// order books, and WebSocket streaming for the QW Trading Platform.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/qw-trading/platform/internal/db"
	"github.com/qw-trading/platform/internal/market/handler"
	"github.com/qw-trading/platform/internal/market/hub"
	"github.com/qw-trading/platform/internal/market/repository"
	"github.com/qw-trading/platform/pkg/config"
	"github.com/qw-trading/platform/pkg/logger"
	"github.com/qw-trading/platform/pkg/middleware"
)

func main() {
	cfg := config.Load()

	database, err := db.Connect(cfg.DatabaseDSN(), cfg.DBMaxOpen, cfg.DBMaxIdle)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	h := hub.NewHub()
	go h.Run()

	repo := repository.New(database)
	hdl := handler.New(repo, h)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /market/tickers", hdl.ListTickers)
	mux.HandleFunc("GET /market/tickers/{symbol}", hdl.GetTicker)
	mux.HandleFunc("GET /market/orderbook/{symbol}", hdl.GetOrderBook)
	mux.HandleFunc("GET /market/ws", hdl.HandleWebSocket)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	logger := logger.New("market-data-service")
	rl := middleware.NewRateLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst)
	wrapped := middleware.RequestID(middleware.Logger(logger)(rl.Middleware(middleware.CORS(cfg.AllowedOrigins)(mux))))

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      wrapped,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Market data service starting on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down market data service...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("Market data service stopped")
}
