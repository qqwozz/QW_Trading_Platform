package feeder

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/qw-trading/platform/internal/models"
	"github.com/qw-trading/platform/internal/market/repository"
)

type binanceTicker struct {
	Symbol             string `json:"symbol"`
	LastPrice          string `json:"lastPrice"`
	BidPrice           string `json:"bidPrice"`
	AskPrice           string `json:"askPrice"`
	BidQty             string `json:"bidQty"`
	AskQty             string `json:"askQty"`
	Volume             string `json:"volume"`
	HighPrice          string `json:"highPrice"`
	LowPrice           string `json:"lowPrice"`
	PriceChange        string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
}

type Feeder struct {
	repo   *repository.MarketRepository
	client *http.Client
	pairs  []pair
}

type pair struct {
	binanceSymbol string
	dbSymbol      string
}

var defaultPairs = []pair{
	{binanceSymbol: "BTCUSDT", dbSymbol: "BTC/USDT"},
	{binanceSymbol: "ETHUSDT", dbSymbol: "ETH/USDT"},
	{binanceSymbol: "SOLUSDT", dbSymbol: "SOL/USDT"},
}

func NewFeeder(repo *repository.MarketRepository) *Feeder {
	return &Feeder{
		repo:   repo,
		client: &http.Client{Timeout: 10 * time.Second},
		pairs:  defaultPairs,
	}
}

func (f *Feeder) Start(ctx context.Context) {
	log.Println("[feeder] starting Binance price feeder")
	go f.loop(ctx)
}

func (f *Feeder) loop(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	f.fetch(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Println("[feeder] shutting down")
			return
		case <-ticker.C:
			f.fetch(ctx)
		}
	}
}

func (f *Feeder) fetch(ctx context.Context) {
	for _, p := range f.pairs {
		if err := f.fetchPair(ctx, p); err != nil {
			log.Printf("[feeder] error fetching %s: %v", p.binanceSymbol, err)
		}
	}
}

func (f *Feeder) fetchPair(ctx context.Context, p pair) error {
	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/24hr?symbol=%s", p.binanceSymbol)
	resp, err := f.client.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Binance API returned status %d: %s", resp.StatusCode, string(body))
	}

	var tickerResp binanceTicker
	if err := json.Unmarshal(body, &tickerResp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	ticker := &models.MarketTicker{
		ID:           uuid.New(),
		Symbol:       p.dbSymbol,
		LastPrice:    parseFloat(tickerResp.LastPrice),
		BestBid:      parseFloat(tickerResp.BidPrice),
		BestAsk:      parseFloat(tickerResp.AskPrice),
		BidVolume:    parseFloat(tickerResp.BidQty),
		AskVolume:    parseFloat(tickerResp.AskQty),
		Volume24h:    parseFloat(tickerResp.Volume),
		High24h:      parseFloat(tickerResp.HighPrice),
		Low24h:       parseFloat(tickerResp.LowPrice),
		Change24h:    parseFloat(tickerResp.PriceChange),
		ChangePct24h: parseFloat(tickerResp.PriceChangePercent),
		UpdatedAt:    time.Now().UTC(),
	}

	if err := f.repo.UpsertTicker(ctx, ticker); err != nil {
		return fmt.Errorf("upsert ticker failed: %w", err)
	}

	log.Printf("[feeder] updated %s: price=%.2f", p.dbSymbol, ticker.LastPrice)
	return nil
}

func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
