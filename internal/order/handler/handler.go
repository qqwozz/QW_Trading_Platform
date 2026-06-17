package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/qw-trading/platform/internal/models"
	"github.com/qw-trading/platform/internal/order/repository"
	"github.com/qw-trading/platform/pkg/errors"
	"github.com/qw-trading/platform/pkg/middleware"
	"github.com/qw-trading/platform/pkg/response"
)

type Handler struct {
	repo      repository.OrderRepositoryInterface
	tradeRepo repository.TradeRepositoryInterface
}

func New(repo repository.OrderRepositoryInterface, tradeRepo repository.TradeRepositoryInterface) *Handler {
	return &Handler{repo: repo, tradeRepo: tradeRepo}
}

type CreateOrderRequest struct {
	Symbol      string   `json:"symbol"`
	Side        string   `json:"side"`
	Type        string   `json:"type"`
	Price       *float64 `json:"price,omitempty"`
	Quantity    float64  `json:"quantity"`
	TimeInForce string   `json:"time_in_force"`
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Unauthorized(w, "unauthorized")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if err := validateOrder(&req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	order := &models.Order{
		ID:          uuid.New(),
		UserID:      userID,
		Symbol:      req.Symbol,
		Side:        models.OrderSide(req.Side),
		Type:        models.OrderType(req.Type),
		Price:       req.Price,
		Quantity:    req.Quantity,
		Status:      models.OrderStatusOpen,
		TimeInForce: models.TimeInForce(req.TimeInForce),
	}

	if order.TimeInForce == "" {
		order.TimeInForce = models.TimeInForceGTC
	}

	if err := h.repo.Create(r.Context(), order); err != nil {
		response.InternalError(w, "failed to create order")
		return
	}

	response.Created(w, toOrderResponse(order))
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Unauthorized(w, "unauthorized")
		return
	}

	orderID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.BadRequest(w, "invalid order id")
		return
	}

	order, err := h.repo.GetByID(r.Context(), orderID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok && appErr.Code == http.StatusNotFound {
			response.NotFound(w, "order not found")
			return
		}
		response.InternalError(w, "failed to get order")
		return
	}

	if order.UserID != userID {
		response.Forbidden(w, "access denied")
		return
	}

	response.Success(w, toOrderResponse(order))
}

func (h *Handler) ListOrders(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Unauthorized(w, "unauthorized")
		return
	}

	limit, offset := response.ParsePagination(r, 50)

	orders, total, err := h.repo.List(r.Context(), repository.ListFilter{
		UserID: userID,
		Symbol: r.URL.Query().Get("symbol"),
		Status: r.URL.Query().Get("status"),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		response.InternalError(w, "failed to list orders")
		return
	}

	result := make([]response.OrderResponse, len(orders))
	for i, order := range orders {
		result[i] = toOrderResponse(&order)
	}

	response.Paginated(w, result, total, limit, offset)
}

func (h *Handler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Unauthorized(w, "unauthorized")
		return
	}

	orderID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.BadRequest(w, "invalid order id")
		return
	}

	order, err := h.repo.GetByID(r.Context(), orderID)
	if err != nil {
		response.NotFound(w, "order not found")
		return
	}

	if order.UserID != userID {
		response.Forbidden(w, "access denied")
		return
	}

	if order.Status != models.OrderStatusOpen && order.Status != models.OrderStatusPartiallyFilled {
		response.Conflict(w, "order cannot be cancelled")
		return
	}

	if err := h.repo.UpdateStatus(r.Context(), orderID, models.OrderStatusCancelled, order.FilledQuantity); err != nil {
		response.InternalError(w, "failed to cancel order")
		return
	}

	order.Status = models.OrderStatusCancelled
	response.Success(w, toOrderResponse(order))
}

func validateOrder(req *CreateOrderRequest) error {
	if req.Symbol == "" {
		return errors.BadRequest("symbol is required")
	}
	if req.Side != "BUY" && req.Side != "SELL" {
		return errors.BadRequest("side must be BUY or SELL")
	}
	if req.Type != "LIMIT" && req.Type != "MARKET" {
		return errors.BadRequest("type must be LIMIT or MARKET")
	}
	if req.Quantity <= 0 {
		return errors.BadRequest("quantity must be positive")
	}
	if req.Type == "LIMIT" && (req.Price == nil || *req.Price <= 0) {
		return errors.BadRequest("price is required for limit orders")
	}
	tif := req.TimeInForce
	if tif == "" {
		tif = "GTC"
	}
	if tif != "GTC" && tif != "IOC" && tif != "FOK" {
		return errors.BadRequest("time_in_force must be GTC, IOC, or FOK")
	}
	return nil
}

func toOrderResponse(order *models.Order) response.OrderResponse {
	return response.OrderResponse{
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
