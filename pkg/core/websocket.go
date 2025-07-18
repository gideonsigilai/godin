package core

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocketManager manages WebSocket connections and channels
type WebSocketManager struct {
	connections map[string]*websocket.Conn
	channels    map[string][]chan interface{}
	upgrader    websocket.Upgrader
	mutex       sync.RWMutex
	enabled     bool
	path        string
}

// NewWebSocketManager creates a new WebSocket manager
func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		connections: make(map[string]*websocket.Conn),
		channels:    make(map[string][]chan interface{}),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins in development
			},
		},
		enabled: false,
		path:    "/ws",
	}
}

// Enable enables WebSocket support with the specified path
func (wsm *WebSocketManager) Enable(path string) {
	wsm.enabled = true
	if path != "" {
		wsm.path = path
	}
}

// IsEnabled returns whether WebSocket is enabled
func (wsm *WebSocketManager) IsEnabled() bool {
	return wsm.enabled
}

// GetPath returns the WebSocket path
func (wsm *WebSocketManager) GetPath() string {
	return wsm.path
}

// HandleConnection handles new WebSocket connections
func (wsm *WebSocketManager) HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := wsm.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// Generate connection ID
	connID := generateConnectionID()

	wsm.mutex.Lock()
	wsm.connections[connID] = conn
	wsm.mutex.Unlock()

	// Clean up on disconnect
	defer func() {
		wsm.mutex.Lock()
		delete(wsm.connections, connID)
		wsm.mutex.Unlock()
	}()

	// Handle incoming messages
	for {
		var message WebSocketMessage
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		wsm.handleMessage(connID, message)
	}
}

// WebSocketMessage represents a WebSocket message
type WebSocketMessage struct {
	Type    string      `json:"type"`
	Channel string      `json:"channel"`
	Data    interface{} `json:"data"`
}

// handleMessage processes incoming WebSocket messages
func (wsm *WebSocketManager) handleMessage(connID string, message WebSocketMessage) {
	switch message.Type {
	case "subscribe":
		wsm.subscribe(connID, message.Channel)
	case "unsubscribe":
		wsm.unsubscribe(connID, message.Channel)
	case "ping":
		wsm.sendToConnection(connID, WebSocketMessage{
			Type: "pong",
			Data: "pong",
		})
	}
}

// Subscribe subscribes a connection to a channel
func (wsm *WebSocketManager) subscribe(connID, channel string) {
	// Implementation for subscribing to channels
	log.Printf("Connection %s subscribed to channel %s", connID, channel)
}

// Unsubscribe unsubscribes a connection from a channel
func (wsm *WebSocketManager) unsubscribe(connID, channel string) {
	// Implementation for unsubscribing from channels
	log.Printf("Connection %s unsubscribed from channel %s", connID, channel)
}

// Broadcast sends data to all connections on a channel
func (wsm *WebSocketManager) Broadcast(channel string, data interface{}) {
	message := WebSocketMessage{
		Type:    "broadcast",
		Channel: channel,
		Data:    data,
	}

	wsm.mutex.RLock()
	defer wsm.mutex.RUnlock()

	for connID, conn := range wsm.connections {
		err := conn.WriteJSON(message)
		if err != nil {
			log.Printf("Error broadcasting to connection %s: %v", connID, err)
		}
	}
}

// SendToConnection sends a message to a specific connection
func (wsm *WebSocketManager) sendToConnection(connID string, message WebSocketMessage) {
	wsm.mutex.RLock()
	conn, exists := wsm.connections[connID]
	wsm.mutex.RUnlock()

	if !exists {
		log.Printf("Connection %s not found", connID)
		return
	}

	err := conn.WriteJSON(message)
	if err != nil {
		log.Printf("Error sending to connection %s: %v", connID, err)
	}
}

// Subscribe creates a channel for receiving data
func (wsm *WebSocketManager) Subscribe(channel string) chan interface{} {
	wsm.mutex.Lock()
	defer wsm.mutex.Unlock()

	ch := make(chan interface{}, 100)
	wsm.channels[channel] = append(wsm.channels[channel], ch)
	return ch
}

// Publish sends data to all subscribers of a channel
func (wsm *WebSocketManager) Publish(channel string, data interface{}) {
	wsm.mutex.RLock()
	channels, exists := wsm.channels[channel]
	wsm.mutex.RUnlock()

	if !exists {
		return
	}

	// Send to local channels
	for _, ch := range channels {
		select {
		case ch <- data:
		default:
			// Channel is full, skip
		}
	}

	// Broadcast to WebSocket connections
	wsm.Broadcast(channel, data)
}

// GetConnectionCount returns the number of active connections
func (wsm *WebSocketManager) GetConnectionCount() int {
	wsm.mutex.RLock()
	defer wsm.mutex.RUnlock()
	return len(wsm.connections)
}

// generateConnectionID generates a unique connection ID
func generateConnectionID() string {
	// Simple implementation - in production, use UUID or similar
	return "conn_" + randomString(8)
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[len(charset)/2] // Simplified for demo
	}
	return string(b)
}
