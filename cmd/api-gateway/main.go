package main

import (
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/qw-trading/platform/pkg/config"
	applog "github.com/qw-trading/platform/pkg/logger"
	"github.com/qw-trading/platform/pkg/middleware"
	"github.com/qw-trading/platform/pkg/server"
)

type route struct {
	prefix string
	target string
}

type Gateway struct {
	routes []route
	client *http.Client
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func NewGateway(_ *config.Config) *Gateway {
	userSvc := envOrDefault("USER_SERVICE_URL", "http://localhost:8081")
	accountSvc := envOrDefault("ACCOUNT_SERVICE_URL", "http://localhost:8082")
	orderSvc := envOrDefault("ORDER_SERVICE_URL", "http://localhost:8083")
	portfolioSvc := envOrDefault("PORTFOLIO_SERVICE_URL", "http://localhost:8084")
	marketSvc := envOrDefault("MARKET_SERVICE_URL", "http://localhost:8085")
	historySvc := envOrDefault("HISTORY_SERVICE_URL", "http://localhost:8086")

	routes := []route{
		{"/v1/auth", userSvc},
		{"/v1/users", userSvc},
		{"/v1/accounts", accountSvc},
		{"/v1/orders", orderSvc},
		{"/v1/positions", portfolioSvc},
		{"/v1/portfolio", portfolioSvc},
		{"/v1/balances", portfolioSvc},
		{"/v1/market", marketSvc},
		{"/v1/history", historySvc},
	}
	sort.Slice(routes, func(i, j int) bool {
		return len(routes[i].prefix) > len(routes[j].prefix)
	})

	return &Gateway{
		routes: routes,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (g *Gateway) route(r *http.Request) string {
	path := r.URL.Path
	for _, rt := range g.routes {
		if strings.HasPrefix(path, rt.prefix) {
			return rt.target
		}
	}
	return ""
}

func (g *Gateway) proxy(w http.ResponseWriter, r *http.Request, target string) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		http.Error(w, `{"error":"failed to read body"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	targetURL := target + r.URL.Path
	if r.URL.RawQuery != "" {
		targetURL += "?" + r.URL.RawQuery
	}

	proxyReq, err := http.NewRequestWithContext(r.Context(), r.Method, targetURL, strings.NewReader(string(body)))
	if err != nil {
		http.Error(w, `{"error":"failed to create request"}`, http.StatusInternalServerError)
		return
	}

	for key, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}

	resp, err := g.client.Do(proxyReq)
	if err != nil {
		http.Error(w, `{"error":"service unavailable"}`, http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
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
	os.Setenv("PORT", cfg.Port)
	gateway := NewGateway(cfg)

	logger := applog.New("api-gateway")
	rl := middleware.NewRateLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst)
	wrapped := middleware.RequestID(middleware.Logger(logger)(rl.Middleware(middleware.CORS(cfg.AllowedOrigins)(gateway))))

	server.Run("api-gateway", wrapped, server.DefaultConfig(), logger)
}
