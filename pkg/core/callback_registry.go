package core

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// CallbackRegistry manages widget callback registration and execution
type CallbackRegistry struct {
	callbacks    map[string]*CallbackInfo
	mutex        sync.RWMutex
	router       *mux.Router
	app          *App
	cleanupTimer *time.Timer
}

// CallbackInfo stores information about a registered callback
type CallbackInfo struct {
	ID           string                 // Unique callback identifier
	WidgetID     string                 // Widget instance identifier
	WidgetType   string                 // Type of widget (Button, TextField, etc.)
	CallbackType string                 // Type of callback (OnPressed, OnChanged, etc.)
	Function     interface{}            // The actual Go function to execute
	Context      *Context               // Context when callback was registered
	Parameters   map[string]interface{} // Additional parameters for the callback
	CreatedAt    time.Time              // When the callback was registered
	LastUsed     time.Time              // When the callback was last executed
}

// CallbackParameter represents a parameter passed to a callback
type CallbackParameter struct {
	Name  string
	Value interface{}
	Type  string
}

// NewCallbackRegistry creates a new callback registry
func NewCallbackRegistry(app *App) *CallbackRegistry {
	registry := &CallbackRegistry{
		callbacks: make(map[string]*CallbackInfo),
		router:    app.Router(),
		app:       app,
	}

	// Start cleanup timer to remove unused callbacks
	registry.startCleanupTimer()

	return registry
}

// RegisterCallback registers a callback function and returns a unique callback ID
func (cr *CallbackRegistry) RegisterCallback(widgetID, widgetType, callbackType string, fn interface{}, ctx *Context) string {
	if fn == nil {
		return ""
	}

	// Validate that fn is actually a function
	fnValue := reflect.ValueOf(fn)
	if fnValue.Kind() != reflect.Func {
		return ""
	}

	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	// Generate unique callback ID
	callbackID := cr.generateCallbackID()

	// Create callback info
	info := &CallbackInfo{
		ID:           callbackID,
		WidgetID:     widgetID,
		WidgetType:   widgetType,
		CallbackType: callbackType,
		Function:     fn,
		Context:      ctx,
		Parameters:   make(map[string]interface{}),
		CreatedAt:    time.Now(),
		LastUsed:     time.Now(),
	}

	// Store callback info
	cr.callbacks[callbackID] = info

	// Generate and register HTTP endpoint
	cr.generateEndpoint(callbackID)

	return callbackID
}

// ExecuteCallback executes a callback by ID with optional parameters
func (cr *CallbackRegistry) ExecuteCallback(callbackID string, params map[string]interface{}) error {
	cr.mutex.RLock()
	info, exists := cr.callbacks[callbackID]
	cr.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("callback with ID %s not found", callbackID)
	}

	// Update last used time
	cr.mutex.Lock()
	info.LastUsed = time.Now()
	cr.mutex.Unlock()

	// Set up global context for state operations
	if info.Context != nil {
		SetGlobalContext(info.Context)
		defer SetGlobalContext(nil)
	}

	// Execute the callback function
	return cr.executeFunction(info.Function, params)
}

// executeFunction executes a function with parameters using reflection
func (cr *CallbackRegistry) executeFunction(fn interface{}, params map[string]interface{}) error {
	fnValue := reflect.ValueOf(fn)
	fnType := fnValue.Type()

	// Prepare arguments based on function signature
	args := make([]reflect.Value, fnType.NumIn())

	for i := 0; i < fnType.NumIn(); i++ {
		paramType := fnType.In(i)

		// Try to find matching parameter by type or position
		var paramValue reflect.Value

		switch paramType.Kind() {
		case reflect.String:
			if val, ok := params["text"].(string); ok {
				paramValue = reflect.ValueOf(val)
			} else if val, ok := params["value"].(string); ok {
				paramValue = reflect.ValueOf(val)
			} else {
				paramValue = reflect.ValueOf("")
			}
		case reflect.Bool:
			if val, ok := params["checked"].(bool); ok {
				paramValue = reflect.ValueOf(val)
			} else if val, ok := params["value"].(bool); ok {
				paramValue = reflect.ValueOf(val)
			} else {
				paramValue = reflect.ValueOf(false)
			}
		case reflect.Int, reflect.Int32, reflect.Int64:
			if val, ok := params["value"].(int); ok {
				paramValue = reflect.ValueOf(val)
			} else if val, ok := params["index"].(int); ok {
				paramValue = reflect.ValueOf(val)
			} else {
				paramValue = reflect.ValueOf(0)
			}
		case reflect.Float32, reflect.Float64:
			if val, ok := params["value"].(float64); ok {
				paramValue = reflect.ValueOf(val)
			} else {
				paramValue = reflect.ValueOf(0.0)
			}
		default:
			// For complex types or no parameters, use zero value
			paramValue = reflect.Zero(paramType)
		}

		args[i] = paramValue
	}

	// Execute the function
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Callback execution panic: %v\n", r)
		}
	}()

	fnValue.Call(args)
	return nil
}

// GenerateEndpoint creates an HTTP endpoint for a callback
func (cr *CallbackRegistry) generateEndpoint(callbackID string) string {
	endpointPath := fmt.Sprintf("/api/callbacks/%s", callbackID)

	// Register the endpoint with the router
	cr.router.HandleFunc(endpointPath, func(w http.ResponseWriter, r *http.Request) {
		// Parse parameters from request
		params := make(map[string]interface{})

		// Parse form data
		if err := r.ParseForm(); err == nil {
			for key, values := range r.Form {
				if len(values) > 0 {
					// Try to parse as different types
					value := values[0]

					// Try boolean
					if value == "true" {
						params[key] = true
					} else if value == "false" {
						params[key] = false
					} else {
						// Keep as string
						params[key] = value
					}
				}
			}
		}

		// Parse JSON body if present
		if r.Header.Get("Content-Type") == "application/json" {
			var jsonParams map[string]interface{}
			if err := NewContext(w, r, cr.app).JSON(&jsonParams); err == nil {
				for key, value := range jsonParams {
					params[key] = value
				}
			}
		}

		// Execute the callback
		if err := cr.ExecuteCallback(callbackID, params); err != nil {
			http.Error(w, fmt.Sprintf("Callback execution failed: %v", err), http.StatusInternalServerError)
			return
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success"}`))

	}).Methods("GET", "POST", "PUT", "DELETE")

	return endpointPath
}

// CleanupCallback removes a callback from the registry
func (cr *CallbackRegistry) CleanupCallback(callbackID string) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()
	delete(cr.callbacks, callbackID)
}

// GetCallbackInfo returns information about a callback
func (cr *CallbackRegistry) GetCallbackInfo(callbackID string) (*CallbackInfo, bool) {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()
	info, exists := cr.callbacks[callbackID]
	return info, exists
}

// ListCallbacks returns all registered callbacks
func (cr *CallbackRegistry) ListCallbacks() map[string]*CallbackInfo {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()

	// Create a copy to avoid race conditions
	result := make(map[string]*CallbackInfo)
	for id, info := range cr.callbacks {
		result[id] = info
	}
	return result
}

// GetCallbacksByWidget returns all callbacks for a specific widget
func (cr *CallbackRegistry) GetCallbacksByWidget(widgetID string) []*CallbackInfo {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()

	var result []*CallbackInfo
	for _, info := range cr.callbacks {
		if info.WidgetID == widgetID {
			result = append(result, info)
		}
	}
	return result
}

// GetCallbacksByType returns all callbacks of a specific type
func (cr *CallbackRegistry) GetCallbacksByType(callbackType string) []*CallbackInfo {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()

	var result []*CallbackInfo
	for _, info := range cr.callbacks {
		if info.CallbackType == callbackType {
			result = append(result, info)
		}
	}
	return result
}

// generateCallbackID generates a unique callback identifier
func (cr *CallbackRegistry) generateCallbackID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// startCleanupTimer starts a timer to periodically clean up unused callbacks
func (cr *CallbackRegistry) startCleanupTimer() {
	cr.cleanupTimer = time.AfterFunc(time.Hour, func() {
		cr.cleanupUnusedCallbacks()
		cr.startCleanupTimer() // Restart timer
	})
}

// cleanupUnusedCallbacks removes callbacks that haven't been used in a while
func (cr *CallbackRegistry) cleanupUnusedCallbacks() {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	cutoff := time.Now().Add(-2 * time.Hour) // Remove callbacks unused for 2 hours

	for id, info := range cr.callbacks {
		if info.LastUsed.Before(cutoff) {
			delete(cr.callbacks, id)
		}
	}
}

// Stop stops the cleanup timer
func (cr *CallbackRegistry) Stop() {
	if cr.cleanupTimer != nil {
		cr.cleanupTimer.Stop()
	}
}

// GetStats returns statistics about the callback registry
func (cr *CallbackRegistry) GetStats() map[string]interface{} {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()

	stats := map[string]interface{}{
		"total_callbacks": len(cr.callbacks),
		"widget_types":    make(map[string]int),
		"callback_types":  make(map[string]int),
	}

	widgetTypes := stats["widget_types"].(map[string]int)
	callbackTypes := stats["callback_types"].(map[string]int)

	for _, info := range cr.callbacks {
		widgetTypes[info.WidgetType]++
		callbackTypes[info.CallbackType]++
	}

	return stats
}
