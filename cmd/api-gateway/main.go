package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/qw-trading/platform/pkg/config"
	"github.com/qw-trading/platform/pkg/middleware"
)

type Gateway struct {
	services map[string]string
	client   *http.Client
}

func NewGateway(cfg *config.Config) *Gateway {
	return &Gateway{
		services: map[string]string{
			"/auth":     "http://localhost:8081",
			"/users":    "http://localhost:8081",
			"/accounts": "http://localhost:8082",
			"/orders":   "http://localhost:8083",
			"/positions":"http://localhost:8084",
			"/portfolio":"http://localhost:8084",
			"/balances": "http://localhost:8084",
			"/market":   "http://localhost:8085",
			"/history":  "http://localhost:8086",
		},
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (g *Gateway) proxy(w http.ResponseWriter, r *http.Request, target string) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		http.Error(w, `{"error":"failed to read body"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	proxyReq, err := http.NewRequestWithContext(r.Context(), r.Method, target+r.URL.Path, strings.NewReader(string(body)))
	if err != nil {
		http.Error(w, `{"error":"failed to create request"}`, http.StatusInternalServerError)
		return
	}

	for key, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}
	proxyReq.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(proxyReq)
	if err != nil {
		http.Error(w, `{"error":"service unavailable"}`, http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (g *Gateway) route(r *http.Request) string {
	path := r.URL.Path
	for prefix, target := range g.services {
		if strings.HasPrefix(path, prefix) {
			return target
		}
	}
	return ""
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/health" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"healthy","service":"api-gateway"}`))
		return
	}

	target := g.route(r)
	if target == "" {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}

	g.proxy(w, r, target)
}

func main() {
	cfg := config.Load()
	gateway := NewGateway(cfg)

	handler := middleware.Logger(middleware.CORS(cfg.AllowedOrigins)(gateway))

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("API Gateway starting on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down API Gateway...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("API Gateway stopped")
}
