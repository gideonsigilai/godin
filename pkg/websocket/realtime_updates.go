package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// RealtimeUpdateManager manages WebSocket connections for real-time updates
type RealtimeUpdateManager struct {
	clients    map[*websocket.Conn]*Client
	channels   map[string]map[*websocket.Conn]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *BroadcastMessage
	mutex      sync.RWMutex
	upgrader   websocket.Upgrader
}

// Client represents a WebSocket client
type Client struct {
	conn     *websocket.Conn
	send     chan []byte
	manager  *RealtimeUpdateManager
	channels map[string]bool
	lastPing time.Time
	isAlive  bool
}

// BroadcastMessage represents a message to broadcast to clients
type BroadcastMessage struct {
	Channel string
	Data    interface{}
}

// Message represents a WebSocket message
type Message struct {
	Type    string      `json:"type"`
	Channel string      `json:"channel,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// NewRealtimeUpdateManager creates a new realtime update manager
func NewRealtimeUpdateManager() *RealtimeUpdateManager {
	return &RealtimeUpdateManager{
		clients:    make(map[*websocket.Conn]*Client),
		channels:   make(map[string]map[*websocket.Conn]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *BroadcastMessage),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow connections from any origin in development
				// In production, you should implement proper origin checking
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

// Start starts the realtime update manager
func (rum *RealtimeUpdateManager) Start() {
	go rum.run()
	go rum.pingClients()
}

// run handles the main event loop
func (rum *RealtimeUpdateManager) run() {
	for {
		select {
		case client := <-rum.register:
			rum.registerClient(client)

		case client := <-rum.unregister:
			rum.unregisterClient(client)

		case message := <-rum.broadcast:
			rum.broadcastToChannel(message.Channel, message.Data)
		}
	}
}

// registerClient registers a new client
func (rum *RealtimeUpdateManager) registerClient(client *Client) {
	rum.mutex.Lock()
	defer rum.mutex.Unlock()

	rum.clients[client.conn] = client
	client.isAlive = true
	client.lastPing = time.Now()

	log.Printf("Client connected. Total clients: %d", len(rum.clients))

	// Send welcome message
	welcomeMsg := Message{
		Type: "connected",
		Data: map[string]interface{}{
			"message":   "Connected to Godin WebSocket",
			"timestamp": time.Now().Unix(),
		},
	}

	if data, err := json.Marshal(welcomeMsg); err == nil {
		select {
		case client.send <- data:
		default:
			close(client.send)
			delete(rum.clients, client.conn)
		}
	}
}

// unregisterClient unregisters a client
func (rum *RealtimeUpdateManager) unregisterClient(client *Client) {
	rum.mutex.Lock()
	defer rum.mutex.Unlock()

	if _, ok := rum.clients[client.conn]; ok {
		// Remove client from all channels
		for channel := range client.channels {
			if clients, exists := rum.channels[channel]; exists {
				delete(clients, client.conn)
				if len(clients) == 0 {
					delete(rum.channels, channel)
				}
			}
		}

		delete(rum.clients, client.conn)
		close(client.send)

		log.Printf("Client disconnected. Total clients: %d", len(rum.clients))
	}
}

// subscribeToChannel subscribes a client to a channel
func (rum *RealtimeUpdateManager) subscribeToChannel(client *Client, channel string) {
	rum.mutex.Lock()
	defer rum.mutex.Unlock()

	// Add client to channel
	if rum.channels[channel] == nil {
		rum.channels[channel] = make(map[*websocket.Conn]bool)
	}
	rum.channels[channel][client.conn] = true

	// Add channel to client
	if client.channels == nil {
		client.channels = make(map[string]bool)
	}
	client.channels[channel] = true

	log.Printf("Client subscribed to channel: %s", channel)
}

// unsubscribeFromChannel unsubscribes a client from a channel
func (rum *RealtimeUpdateManager) unsubscribeFromChannel(client *Client, channel string) {
	rum.mutex.Lock()
	defer rum.mutex.Unlock()

	// Remove client from channel
	if clients, exists := rum.channels[channel]; exists {
		delete(clients, client.conn)
		if len(clients) == 0 {
			delete(rum.channels, channel)
		}
	}

	// Remove channel from client
	if client.channels != nil {
		delete(client.channels, channel)
	}

	log.Printf("Client unsubscribed from channel: %s", channel)
}

// broadcastToChannel broadcasts a message to all clients in a channel
func (rum *RealtimeUpdateManager) broadcastToChannel(channel string, data interface{}) {
	rum.mutex.RLock()
	clients, exists := rum.channels[channel]
	rum.mutex.RUnlock()

	if !exists {
		return
	}

	message := Message{
		Type:    "broadcast",
		Channel: channel,
		Data:    data,
	}

	messageData, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling broadcast message: %v", err)
		return
	}

	for conn := range clients {
		if client, exists := rum.clients[conn]; exists && client.isAlive {
			select {
			case client.send <- messageData:
			default:
				// Client's send channel is full, remove client
				rum.unregister <- client
			}
		}
	}
}

// Broadcast broadcasts a message to a channel (implements WebSocketBroadcaster interface)
func (rum *RealtimeUpdateManager) Broadcast(channel string, data interface{}) {
	select {
	case rum.broadcast <- &BroadcastMessage{Channel: channel, Data: data}:
	default:
		log.Printf("Broadcast channel full, dropping message for channel: %s", channel)
	}
}

// HandleWebSocket handles WebSocket connections
func (rum *RealtimeUpdateManager) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := rum.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		conn:     conn,
		send:     make(chan []byte, 256),
		manager:  rum,
		channels: make(map[string]bool),
	}

	rum.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()
}

// readPump handles reading messages from the WebSocket connection
func (c *Client) readPump() {
	defer func() {
		c.manager.unregister <- c
		c.conn.Close()
	}()

	// Set read deadline and pong handler
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.lastPing = time.Now()
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, messageData, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(messageData, &msg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		c.handleMessage(&msg)
	}
}

// writePump handles writing messages to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage handles incoming messages from clients
func (c *Client) handleMessage(msg *Message) {
	switch msg.Type {
	case "subscribe":
		if msg.Channel != "" {
			c.manager.subscribeToChannel(c, msg.Channel)
		}

	case "unsubscribe":
		if msg.Channel != "" {
			c.manager.unsubscribeFromChannel(c, msg.Channel)
		}

	case "ping":
		c.lastPing = time.Now()
		pongMsg := Message{
			Type: "pong",
			Data: map[string]interface{}{
				"timestamp": time.Now().Unix(),
			},
		}
		if data, err := json.Marshal(pongMsg); err == nil {
			select {
			case c.send <- data:
			default:
				// Send channel is full
			}
		}

	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

// pingClients periodically checks client health
func (rum *RealtimeUpdateManager) pingClients() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rum.mutex.RLock()
			for _, client := range rum.clients {
				if time.Since(client.lastPing) > 90*time.Second {
					client.isAlive = false
					rum.mutex.RUnlock()
					rum.unregister <- client
					rum.mutex.RLock()
				}
			}
			rum.mutex.RUnlock()
		}
	}
}

// GetStats returns statistics about the WebSocket manager
func (rum *RealtimeUpdateManager) GetStats() map[string]interface{} {
	rum.mutex.RLock()
	defer rum.mutex.RUnlock()

	channelStats := make(map[string]int)
	for channel, clients := range rum.channels {
		channelStats[channel] = len(clients)
	}

	return map[string]interface{}{
		"total_clients":  len(rum.clients),
		"total_channels": len(rum.channels),
		"channel_stats":  channelStats,
		"timestamp":      time.Now().Unix(),
	}
}

// BroadcastToAll broadcasts a message to all connected clients
func (rum *RealtimeUpdateManager) BroadcastToAll(data interface{}) {
	rum.mutex.RLock()
	clients := make([]*Client, 0, len(rum.clients))
	for _, client := range rum.clients {
		if client.isAlive {
			clients = append(clients, client)
		}
	}
	rum.mutex.RUnlock()

	message := Message{
		Type: "broadcast_all",
		Data: data,
	}

	messageData, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling broadcast message: %v", err)
		return
	}

	for _, client := range clients {
		select {
		case client.send <- messageData:
		default:
			// Client's send channel is full, remove client
			rum.unregister <- client
		}
	}
}
