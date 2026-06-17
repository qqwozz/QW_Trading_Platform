// Command portfolio-service provides the portfolio and position management
// service for the QW Trading Platform.
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
	"github.com/qw-trading/platform/internal/portfolio/handler"
	"github.com/qw-trading/platform/internal/portfolio/repository"
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

	posRepo := repository.New(database)
	accRepo := repository.NewAccountRepository(database)
	h := handler.New(posRepo, accRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /portfolio", h.GetPortfolio)
	mux.HandleFunc("GET /positions", h.ListPositions)
	mux.HandleFunc("POST /positions", h.UpdatePosition)
	mux.HandleFunc("GET /balances", h.GetBalances)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	logger := logger.New("portfolio-service")
	rl := middleware.NewRateLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst)
	wrapped := middleware.RequestID(middleware.Logger(logger)(rl.Middleware(middleware.CORS(cfg.AllowedOrigins)(mux))))

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      wrapped,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Portfolio service starting on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down portfolio service...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("Portfolio service stopped")
}
