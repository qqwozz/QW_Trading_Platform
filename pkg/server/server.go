package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/qw-trading/platform/pkg/logger"
)

type Config struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	ShutdownWait time.Duration
}

func DefaultConfig() Config {
	return Config{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		ShutdownWait: 30 * time.Second,
	}
}

func Run(name string, handler http.Handler, cfg Config, log *logger.Logger) {
	srv := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		log.Info(name + " starting", map[string]interface{}{
			"addr": srv.Addr,
		})
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(name+" failed to start", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	log.Info(name+" received signal, shutting down", map[string]interface{}{
		"signal": sig.String(),
	})

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownWait)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error(name+" shutdown error", map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		log.Info(name+" stopped gracefully")
	}
}
