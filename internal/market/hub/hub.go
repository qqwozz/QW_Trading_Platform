package hub

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	symbol string
}

type WSMessage struct {
	Type   string          `json:"type"`
	Symbol string          `json:"symbol,omitempty"`
	Data   json.RawMessage `json:"data,omitempty"`
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			total := len(h.clients)
			h.mu.Unlock()
			log.Printf("WebSocket client connected, symbol: %s, total: %d", client.symbol, total)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			total := len(h.clients)
			h.mu.Unlock()
			log.Printf("WebSocket client disconnected, symbol: %s, total: %d", client.symbol, total)

		case message := <-h.broadcast:
			h.broadcastMessage(message, "")
		}
	}
}

func (h *Hub) broadcastMessage(payload []byte, symbol string) {
	h.mu.RLock()
	var stale []*Client
	for client := range h.clients {
		if symbol != "" && client.symbol != "" && client.symbol != symbol {
			continue
		}
		select {
		case client.send <- payload:
		default:
			stale = append(stale, client)
		}
	}
	h.mu.RUnlock()

	if len(stale) > 0 {
		h.mu.Lock()
		for _, client := range stale {
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		}
		h.mu.Unlock()
	}
}

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

func (h *Hub) BroadcastTicker(symbol string, data []byte) {
	h.marshalAndBroadcast("ticker", symbol, data)
}

func (h *Hub) BroadcastOrderBook(symbol string, data []byte) {
	h.marshalAndBroadcast("orderbook", symbol, data)
}

func (h *Hub) marshalAndBroadcast(msgType, symbol string, data []byte) {
	msg := WSMessage{Type: msgType, Symbol: symbol, Data: data}
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal %s broadcast: %v", msgType, err)
		return
	}
	h.broadcastMessage(payload, symbol)
}

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
		case "unsubscribe":
			c.symbol = ""
		}
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()
	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			break
		}
	}
}
