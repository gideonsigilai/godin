package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gideonsigilai/godin/pkg/state"
	"github.com/gorilla/websocket"
)

// StateWebSocketHandler manages WebSocket connections for state synchronization
type StateWebSocketHandler struct {
	upgrader     websocket.Upgrader
	clients      map[string]*StateClient
	subscribers  map[string][]string // notifier_id -> client_ids
	mutex        sync.RWMutex
	stateManager *state.StateManager
}

// StateClient represents a WebSocket client connection for state management
type StateClient struct {
	ID         string
	Connection *websocket.Conn
	Send       chan []byte
	Handler    *StateWebSocketHandler
	mutex      sync.Mutex
}

// StateChangeMessage represents a state change message
type StateChangeMessage struct {
	Type       string      `json:"type"`
	NotifierID string      `json:"id"`
	Value      interface{} `json:"value"`
	Timestamp  int64       `json:"timestamp"`
	HTML       string      `json:"html,omitempty"`
}

// SubscriptionMessage represents a subscription message
type SubscriptionMessage struct {
	Type       string `json:"type"`
	NotifierID string `json:"notifier_id"`
	Action     string `json:"action"` // "subscribe" or "unsubscribe"
}

// NewStateWebSocketHandler creates a new StateWebSocketHandler
func NewStateWebSocketHandler(stateManager *state.StateManager) *StateWebSocketHandler {
	return &StateWebSocketHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow all origins for development - in production, implement proper origin checking
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		clients:      make(map[string]*StateClient),
		subscribers:  make(map[string][]string),
		stateManager: stateManager,
	}
}

// HandleWebSocket handles WebSocket connections
func (h *StateWebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	clientID := generateClientID()
	client := &StateClient{
		ID:         clientID,
		Connection: conn,
		Send:       make(chan []byte, 256),
		Handler:    h,
	}

	h.mutex.Lock()
	h.clients[clientID] = client
	h.mutex.Unlock()

	log.Printf("WebSocket client connected: %s", clientID)

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()
}

// HandleValueChange handles value changes and broadcasts to subscribers
func (h *StateWebSocketHandler) HandleValueChange(notifierID string, value interface{}) {
	h.mutex.RLock()
	subscriberIDs := h.subscribers[notifierID]
	h.mutex.RUnlock()

	if len(subscriberIDs) == 0 {
		return
	}

	message := StateChangeMessage{
		Type:       "value_change",
		NotifierID: notifierID,
		Value:      value,
		Timestamp:  time.Now().Unix(),
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling state change message: %v", err)
		return
	}

	// Broadcast to all subscribers
	h.mutex.RLock()
	for _, clientID := range subscriberIDs {
		if client, exists := h.clients[clientID]; exists {
			select {
			case client.Send <- messageBytes:
			default:
				// Client's send channel is full, close it
				close(client.Send)
				delete(h.clients, clientID)
			}
		}
	}
	h.mutex.RUnlock()
}

// BroadcastToSubscribers broadcasts a message to all subscribers of a notifier
func (h *StateWebSocketHandler) BroadcastToSubscribers(notifierID string, message interface{}) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling broadcast message: %v", err)
		return
	}

	h.mutex.RLock()
	subscriberIDs := h.subscribers[notifierID]
	for _, clientID := range subscriberIDs {
		if client, exists := h.clients[clientID]; exists {
			select {
			case client.Send <- messageBytes:
			default:
				// Client's send channel is full, close it
				close(client.Send)
				delete(h.clients, clientID)
			}
		}
	}
	h.mutex.RUnlock()
}

// Subscribe adds a client to a notifier's subscriber list
func (h *StateWebSocketHandler) Subscribe(clientID, notifierID string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.subscribers[notifierID] == nil {
		h.subscribers[notifierID] = make([]string, 0)
	}

	// Check if already subscribed
	for _, id := range h.subscribers[notifierID] {
		if id == clientID {
			return
		}
	}

	h.subscribers[notifierID] = append(h.subscribers[notifierID], clientID)
	log.Printf("Client %s subscribed to notifier %s", clientID, notifierID)
}

// Unsubscribe removes a client from a notifier's subscriber list
func (h *StateWebSocketHandler) Unsubscribe(clientID, notifierID string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	subscribers := h.subscribers[notifierID]
	for i, id := range subscribers {
		if id == clientID {
			h.subscribers[notifierID] = append(subscribers[:i], subscribers[i+1:]...)
			log.Printf("Client %s unsubscribed from notifier %s", clientID, notifierID)
			break
		}
	}
}

// RemoveClient removes a client and all its subscriptions
func (h *StateWebSocketHandler) RemoveClient(clientID string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Remove from clients map
	delete(h.clients, clientID)

	// Remove from all subscriptions
	for notifierID, subscribers := range h.subscribers {
		for i, id := range subscribers {
			if id == clientID {
				h.subscribers[notifierID] = append(subscribers[:i], subscribers[i+1:]...)
				break
			}
		}
	}

	log.Printf("Client %s removed", clientID)
}

// GetSubscriberCount returns the number of subscribers for a notifier
func (h *StateWebSocketHandler) GetSubscriberCount(notifierID string) int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.subscribers[notifierID])
}

// GetClientCount returns the total number of connected clients
func (h *StateWebSocketHandler) GetClientCount() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.clients)
}

// readPump handles reading messages from the WebSocket connection
func (c *StateClient) readPump() {
	defer func() {
		c.Handler.RemoveClient(c.ID)
		c.Connection.Close()
	}()

	c.Connection.SetReadLimit(512)
	c.Connection.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Connection.SetPongHandler(func(string) error {
		c.Connection.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Handle subscription messages
		var subMsg SubscriptionMessage
		if err := json.Unmarshal(message, &subMsg); err == nil {
			if subMsg.Type == "subscription" {
				if subMsg.Action == "subscribe" {
					c.Handler.Subscribe(c.ID, subMsg.NotifierID)
				} else if subMsg.Action == "unsubscribe" {
					c.Handler.Unsubscribe(c.ID, subMsg.NotifierID)
				}
			}
		}
	}
}

// writePump handles writing messages to the WebSocket connection
func (c *StateClient) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Connection.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Connection.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Connection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current message
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Connection.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// generateClientID generates a unique client ID
func generateClientID() string {
	return fmt.Sprintf("client_%d", time.Now().UnixNano())
}
