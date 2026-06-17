package main

import (
	"log"
	"net/http"
	"os"

	"github.com/qw-trading/platform/internal/db"
	"github.com/qw-trading/platform/internal/user/handler"
	"github.com/qw-trading/platform/internal/user/repository"
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
	h := handler.New(repo, cfg.JWTSecret, cfg.JWTExpiry)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/register", h.Register)
	mux.HandleFunc("POST /auth/login", h.Login)
	mux.HandleFunc("GET /users/me", h.GetProfile)
	mux.HandleFunc("GET /health", db.HealthHandler(database, "user-service"))

	logger := applog.New("user-service")
	rl := middleware.NewRateLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst)
	wrapped := middleware.RequestID(middleware.Logger(logger)(rl.Middleware(middleware.CORS(cfg.AllowedOrigins)(mux))))

	server.Run("user-service", wrapped, server.DefaultConfig(), logger)
}
