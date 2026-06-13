// Package hub provides a WebSocket hub for broadcasting real-time market data
// to connected clients.
package hub

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// upgrader configures WebSocket connection upgrades with permissive origin checks
// for development. In production, CheckOrigin should validate allowed origins.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Hub manages WebSocket client connections and broadcasts messages to subscribers.
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// Client represents a single WebSocket connection subscribed to market data.
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	symbol string
}

// WSMessage is the envelope for all WebSocket messages.
type WSMessage struct {
	Type   string          `json:"type"`
	Symbol string          `json:"symbol,omitempty"`
	Data   json.RawMessage `json:"data,omitempty"`
}

// NewHub creates a new Hub with initialized channels and client map.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub's main event loop, processing client registrations,
// disconnections, and message broadcasts.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("WebSocket client connected, symbol: %s, total: %d", client.symbol, len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("WebSocket client disconnected, symbol: %s, total: %d", client.symbol, len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					// Client buffer full; disconnect it to avoid blocking.
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// HandleWebSocket upgrades an HTTP request to a WebSocket connection and
// registers the new client with the hub.
func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	client := &Client{
		hub:  h,
		conn: conn,
		send: make(chan []byte, 256),
	}

	h.register <- client

	go client.writePump()
	go client.readPump()
}

// BroadcastTicker sends a ticker update to all clients subscribed to the
// given symbol (or all clients if symbol is empty).
func (h *Hub) BroadcastTicker(symbol string, data []byte) {
	msg := WSMessage{
		Type:   "ticker",
		Symbol: symbol,
		Data:   data,
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal ticker broadcast: %v", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.symbol == "" || client.symbol == symbol {
			select {
			case client.send <- payload:
			default:
				// Client buffer full; disconnect it.
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
}

// BroadcastOrderBook sends an order book update to all clients subscribed to
// the given symbol (or all clients if symbol is empty).
func (h *Hub) BroadcastOrderBook(symbol string, data []byte) {
	msg := WSMessage{
		Type:   "orderbook",
		Symbol: symbol,
		Data:   data,
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal orderbook broadcast: %v", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.symbol == "" || client.symbol == symbol {
			select {
			case client.send <- payload:
			default:
				// Client buffer full; disconnect it.
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
}

// readPump reads incoming WebSocket messages and handles subscribe/unsubscribe
// commands from the client.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		var msg WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		switch msg.Type {
		case "subscribe":
			c.symbol = msg.Symbol
			log.Printf("Client subscribed to: %s", msg.Symbol)
		case "unsubscribe":
			c.symbol = ""
			log.Printf("Client unsubscribed")
		}
	}
}

// writePump sends messages from the send channel to the WebSocket connection.
func (c *Client) writePump() {
	defer c.conn.Close()

	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			break
		}
	}
}
