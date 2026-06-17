package response

import (
	"github.com/google/uuid"
	"github.com/qw-trading/platform/internal/models"
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

func OrderFromModel(order *models.Order) OrderResponse {
	return OrderResponse{
		ID:             order.ID.String(),
		Symbol:         order.Symbol,
		Side:           string(order.Side),
		Type:           string(order.Type),
		Price:          order.Price,
		Quantity:       order.Quantity,
		FilledQuantity: order.FilledQuantity,
		Status:         string(order.Status),
		TimeInForce:    string(order.TimeInForce),
		CreatedAt:      order.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:      order.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func TradeFromModel(trade *models.Trade) TradeResponse {
	return TradeResponse{
		ID:            trade.ID.String(),
		Symbol:        trade.Symbol,
		BuyerOrderID:  trade.BuyerOrderID.String(),
		SellerOrderID: trade.SellerOrderID.String(),
		BuyerID:       trade.BuyerID.String(),
		SellerID:      trade.SellerID.String(),
		Price:         trade.Price,
		Quantity:      trade.Quantity,
		BuyerFee:      trade.BuyerFee,
		SellerFee:     trade.SellerFee,
		ExecutedAt:    trade.ExecutedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func MustParseUUID(s string) uuid.UUID {
	id, _ := uuid.Parse(s)
	return id
}
