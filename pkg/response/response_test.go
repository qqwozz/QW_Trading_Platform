package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/qw-trading/platform/internal/models"
)

func TestJSON(t *testing.T) {
	w := httptest.NewRecorder()
	JSON(w, http.StatusOK, map[string]string{"key": "value"})

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}

	var body map[string]string
	json.NewDecoder(w.Body).Decode(&body)
	if body["key"] != "value" {
		t.Errorf("body key = %q, want %q", body["key"], "value")
	}
}

func TestSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	Success(w, "ok")
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestCreated(t *testing.T) {
	w := httptest.NewRecorder()
	Created(w, "created")
	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
	}
}

func TestError(t *testing.T) {
	w := httptest.NewRecorder()
	Error(w, http.StatusBadGateway, "bad gateway")
	if w.Code != http.StatusBadGateway {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadGateway)
	}
	var resp Response
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Error != "bad gateway" {
		t.Errorf("error = %q, want %q", resp.Error, "bad gateway")
	}
}

func TestBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	BadRequest(w, "bad")
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestUnauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	Unauthorized(w, "no auth")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	NotFound(w, "missing")
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestForbidden(t *testing.T) {
	w := httptest.NewRecorder()
	Forbidden(w, "denied")
	if w.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestConflict(t *testing.T) {
	w := httptest.NewRecorder()
	Conflict(w, "conflict")
	if w.Code != http.StatusConflict {
		t.Errorf("status = %d, want %d", w.Code, http.StatusConflict)
	}
}

func TestInternalError(t *testing.T) {
	w := httptest.NewRecorder()
	InternalError(w, "oops")
	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestPaginated(t *testing.T) {
	w := httptest.NewRecorder()
	Paginated(w, []string{"a", "b"}, 10, 5, 0)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp Response
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Meta == nil {
		t.Fatal("meta is nil")
	}
	if resp.Meta.Total != 10 || resp.Meta.Limit != 5 || resp.Meta.Offset != 0 {
		t.Errorf("meta = %+v, want total=10 limit=5 offset=0", resp.Meta)
	}
}

func TestParsePagination_Defaults(t *testing.T) {
	r := httptest.NewRequest("GET", "/test", nil)
	limit, offset := ParsePagination(r, 50)
	if limit != 50 {
		t.Errorf("limit = %d, want 50", limit)
	}
	if offset != 0 {
		t.Errorf("offset = %d, want 0", offset)
	}
}

func TestParsePagination_Custom(t *testing.T) {
	r := httptest.NewRequest("GET", "/test?limit=10&offset=20", nil)
	limit, offset := ParsePagination(r, 50)
	if limit != 10 {
		t.Errorf("limit = %d, want 10", limit)
	}
	if offset != 20 {
		t.Errorf("offset = %d, want 20", offset)
	}
}

func TestParsePagination_OverLimit(t *testing.T) {
	r := httptest.NewRequest("GET", "/test?limit=200", nil)
	limit, _ := ParsePagination(r, 50)
	if limit != 50 {
		t.Errorf("limit = %d, want 50 (clamped)", limit)
	}
}

func TestParsePagination_NegativeOffset(t *testing.T) {
	r := httptest.NewRequest("GET", "/test?offset=-5", nil)
	_, offset := ParsePagination(r, 50)
	if offset != 0 {
		t.Errorf("offset = %d, want 0 (clamped)", offset)
	}
}

func TestOrderFromModel(t *testing.T) {
	now := time.Now()
	order := &models.Order{
		ID: uuid.New(), Symbol: "BTC/USDT", Side: models.OrderSideBuy, Type: models.OrderTypeLimit,
		Quantity: 1.5, Status: models.OrderStatusOpen, TimeInForce: models.TimeInForceGTC,
		CreatedAt: now, UpdatedAt: now,
	}
	resp := OrderFromModel(order)
	if resp.Symbol != "BTC/USDT" {
		t.Errorf("Symbol = %q, want BTC/USDT", resp.Symbol)
	}
	if resp.Side != "BUY" {
		t.Errorf("Side = %q, want BUY", resp.Side)
	}
}

func TestTradeFromModel(t *testing.T) {
	now := time.Now()
	trade := &models.Trade{
		ID: uuid.New(), Symbol: "ETH/USDT", Price: 3000, Quantity: 2,
		ExecutedAt: now,
	}
	resp := TradeFromModel(trade)
	if resp.Symbol != "ETH/USDT" {
		t.Errorf("Symbol = %q, want ETH/USDT", resp.Symbol)
	}
	if resp.Price != 3000 {
		t.Errorf("Price = %f, want 3000", resp.Price)
	}
}
