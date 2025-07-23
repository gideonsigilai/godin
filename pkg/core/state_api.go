package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gideonsigilai/godin/pkg/state"
)

// StateAPIHandler handles state-related API endpoints
type StateAPIHandler struct {
	stateManager *state.StateManager
}

// NewStateAPIHandler creates a new state API handler
func NewStateAPIHandler(stateManager *state.StateManager) *StateAPIHandler {
	return &StateAPIHandler{
		stateManager: stateManager,
	}
}

// RegisterRoutes registers state API routes with the app
func (h *StateAPIHandler) RegisterRoutes(app *App) {
	app.router.HandleFunc("/api/state", h.handleStateList)
	app.router.HandleFunc("/api/state/", h.handleStateItem)
}

// handleStateList handles GET /api/state - lists all ValueNotifiers
func (h *StateAPIHandler) handleStateList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		h.listValueNotifiers(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleStateItem handles /api/state/{id} - get/update specific ValueNotifier
func (h *StateAPIHandler) handleStateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract ID from path
	path := strings.TrimPrefix(r.URL.Path, "/api/state/")
	if path == "" {
		http.Error(w, "ValueNotifier ID required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getValueNotifierState(w, r, path)
	case http.MethodPost:
		h.updateValueNotifier(w, r, path)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// listValueNotifiers returns a list of all registered ValueNotifiers
func (h *StateAPIHandler) listValueNotifiers(w http.ResponseWriter, r *http.Request) {
	ids := h.stateManager.ListValueNotifiers()

	response := map[string]interface{}{
		"success":   true,
		"count":     len(ids),
		"notifiers": ids,
	}

	json.NewEncoder(w).Encode(response)
}

// getValueNotifierState returns the current state of a specific ValueNotifier
func (h *StateAPIHandler) getValueNotifierState(w http.ResponseWriter, r *http.Request, id string) {
	stateData, err := h.stateManager.GetValueNotifierState(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting state: %v", err), http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"state":   stateData,
	}

	json.NewEncoder(w).Encode(response)
}

// updateValueNotifier updates a ValueNotifier's value
func (h *StateAPIHandler) updateValueNotifier(w http.ResponseWriter, r *http.Request, id string) {
	var requestData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	notifier, exists := h.stateManager.GetValueNotifier(id)
	if !exists {
		http.Error(w, "ValueNotifier not found", http.StatusNotFound)
		return
	}

	// Try to update the value
	newValue, hasValue := requestData["value"]
	if !hasValue {
		http.Error(w, "Missing 'value' field", http.StatusBadRequest)
		return
	}

	// Try to call SetValue method if it exists
	if valueSetter, ok := notifier.(interface{ SetValue(interface{}) }); ok {
		valueSetter.SetValue(newValue)
	} else if jsonUpdater, ok := notifier.(interface{ FromJSON([]byte) error }); ok {
		// Try to update via JSON
		jsonData, err := json.Marshal(newValue)
		if err != nil {
			http.Error(w, "Failed to serialize value", http.StatusBadRequest)
			return
		}
		if err := jsonUpdater.FromJSON(jsonData); err != nil {
			http.Error(w, fmt.Sprintf("Failed to update value: %v", err), http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "ValueNotifier does not support updates", http.StatusBadRequest)
		return
	}

	// Return updated state
	stateData, err := h.stateManager.GetValueNotifierState(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting updated state: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Value updated successfully",
		"state":   stateData,
	}

	json.NewEncoder(w).Encode(response)
}

// StateWebSocketHandler handles WebSocket connections for real-time state updates
type StateWebSocketHandler struct {
	stateManager *state.StateManager
}

// NewStateWebSocketHandler creates a new state WebSocket handler
func NewStateWebSocketHandler(stateManager *state.StateManager) *StateWebSocketHandler {
	return &StateWebSocketHandler{
		stateManager: stateManager,
	}
}

// HandleWebSocket handles WebSocket connections for state updates
func (h *StateWebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// This would integrate with the existing WebSocket system
	// For now, we'll just return a placeholder
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message":  "WebSocket state handler - integrate with existing WebSocket system",
		"endpoint": "/ws",
	}
	json.NewEncoder(w).Encode(response)
}

// StateMiddleware provides state management middleware
func StateMiddleware(stateManager *state.StateManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Add state manager to request context if needed
			// For now, just pass through
			next.ServeHTTP(w, r)
		})
	}
}

// Helper function to create a complete state API setup
func SetupStateAPI(app *App, stateManager *state.StateManager) {
	// Create and register state API handler
	stateAPI := NewStateAPIHandler(stateManager)
	stateAPI.RegisterRoutes(app)

	// Create and register WebSocket handler
	wsHandler := NewStateWebSocketHandler(stateManager)
	app.router.HandleFunc("/api/state/ws", wsHandler.HandleWebSocket)

	// Add state middleware if needed
	// app.router.Use(StateMiddleware(stateManager))
}
