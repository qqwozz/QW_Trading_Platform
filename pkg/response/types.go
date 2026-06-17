package response

import (
	"encoding/json"
	"net/http"
)

type OrderResponse struct {
	ID             string   `json:"id"`
	Symbol         string   `json:"symbol"`
	Side           string   `json:"side"`
	Type           string   `json:"type"`
	Price          *float64 `json:"price,omitempty"`
	Quantity       float64  `json:"quantity"`
	FilledQuantity float64  `json:"filled_quantity"`
	Status         string   `json:"status"`
	TimeInForce    string   `json:"time_in_force"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
}

type TradeResponse struct {
	ID            string  `json:"id"`
	Symbol        string  `json:"symbol"`
	BuyerOrderID  string  `json:"buyer_order_id"`
	SellerOrderID string  `json:"seller_order_id"`
	BuyerID       string  `json:"buyer_id"`
	SellerID      string  `json:"seller_id"`
	Price         float64 `json:"price"`
	Quantity      float64 `json:"quantity"`
	BuyerFee      float64 `json:"buyer_fee"`
	SellerFee     float64 `json:"seller_fee"`
	ExecutedAt    string  `json:"executed_at"`
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
