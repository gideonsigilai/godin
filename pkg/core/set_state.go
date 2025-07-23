package core

import (
	"fmt"
	"sync"
	"time"
)

// StateUpdater provides Flutter-like setState functionality
type StateUpdater struct {
	context      *Context
	updateQueue  chan StateUpdate
	batchTimer   *time.Timer
	mutex        sync.RWMutex
	isProcessing bool
	batchDelay   time.Duration
}

// StateUpdate represents a state update operation
type StateUpdate struct {
	UpdateFunc func()                 // Function to execute for the state update
	WidgetIDs  []string               // Widget IDs that should be rebuilt
	Timestamp  time.Time              // When the update was queued
	Context    *Context               // Context for the update
	Metadata   map[string]interface{} // Additional metadata
}

// Global state updater instance
var globalStateUpdater *StateUpdater
var globalStateUpdaterMutex sync.RWMutex

// InitializeStateUpdater initializes the global state updater
func InitializeStateUpdater(ctx *Context) {
	globalStateUpdaterMutex.Lock()
	defer globalStateUpdaterMutex.Unlock()

	if globalStateUpdater == nil {
		globalStateUpdater = NewStateUpdater(ctx)
	}
}

// NewStateUpdater creates a new StateUpdater instance
func NewStateUpdater(ctx *Context) *StateUpdater {
	updater := &StateUpdater{
		context:     ctx,
		updateQueue: make(chan StateUpdate, 100), // Buffer for 100 updates
		batchDelay:  50 * time.Millisecond,       // 50ms batch delay
	}

	// Start the update processor goroutine
	go updater.processUpdates()

	return updater
}

// SetStateFunc executes a state update function and triggers UI rebuilds
// This is the main setState function similar to Flutter
func SetStateFunc(ctx *Context, updateFunc func()) error {
	if ctx == nil {
		// Try to use global context if available
		globalStateUpdaterMutex.RLock()
		if globalStateUpdater != nil && globalStateUpdater.context != nil {
			ctx = globalStateUpdater.context
		}
		globalStateUpdaterMutex.RUnlock()

		if ctx == nil {
			return fmt.Errorf("no context available for setState")
		}
	}

	// Initialize global state updater if needed
	InitializeStateUpdater(ctx)

	return globalStateUpdater.QueueUpdate(StateUpdate{
		UpdateFunc: updateFunc,
		Timestamp:  time.Now(),
		Context:    ctx,
	})
}

// SetStateWithWidgets executes a state update and rebuilds specific widgets
func SetStateWithWidgets(ctx *Context, updateFunc func(), widgetIDs []string) error {
	if ctx == nil {
		return fmt.Errorf("context is required for setState")
	}

	// Initialize global state updater if needed
	InitializeStateUpdater(ctx)

	return globalStateUpdater.QueueUpdate(StateUpdate{
		UpdateFunc: updateFunc,
		WidgetIDs:  widgetIDs,
		Timestamp:  time.Now(),
		Context:    ctx,
	})
}

// QueueUpdate queues a state update for batch processing
func (su *StateUpdater) QueueUpdate(update StateUpdate) error {
	select {
	case su.updateQueue <- update:
		return nil
	default:
		// Queue is full, process immediately
		return su.processUpdate(update)
	}
}

// processUpdates processes queued state updates in batches
func (su *StateUpdater) processUpdates() {
	var batch []StateUpdate
	batchTimer := time.NewTimer(su.batchDelay)
	batchTimer.Stop() // Stop initially

	for {
		select {
		case update := <-su.updateQueue:
			batch = append(batch, update)

			// Start or reset the batch timer
			if !batchTimer.Stop() {
				select {
				case <-batchTimer.C:
				default:
				}
			}
			batchTimer.Reset(su.batchDelay)

		case <-batchTimer.C:
			if len(batch) > 0 {
				su.processBatch(batch)
				batch = nil
			}
		}
	}
}

// processBatch processes a batch of state updates
func (su *StateUpdater) processBatch(batch []StateUpdate) {
	su.mutex.Lock()
	su.isProcessing = true
	su.mutex.Unlock()

	defer func() {
		su.mutex.Lock()
		su.isProcessing = false
		su.mutex.Unlock()
	}()

	// Collect all widget IDs that need rebuilding
	widgetIDSet := make(map[string]bool)
	var allContexts []*Context

	// Execute all update functions
	for _, update := range batch {
		// Set global context for the update
		if update.Context != nil {
			SetGlobalContext(update.Context)
			allContexts = append(allContexts, update.Context)
		}

		// Execute the update function
		if update.UpdateFunc != nil {
			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("setState update function panic: %v\n", r)
					}
				}()
				update.UpdateFunc()
			}()
		}

		// Collect widget IDs
		for _, widgetID := range update.WidgetIDs {
			widgetIDSet[widgetID] = true
		}
	}

	// Clean up global context
	SetGlobalContext(nil)

	// Convert widget ID set to slice
	var widgetIDs []string
	for widgetID := range widgetIDSet {
		widgetIDs = append(widgetIDs, widgetID)
	}

	// Trigger rebuilds for affected widgets
	if len(widgetIDs) > 0 {
		su.TriggerRebuild(widgetIDs)
	}

	// Broadcast state change event via WebSocket if available
	su.broadcastStateChange(batch, widgetIDs)
}

// processUpdate processes a single state update immediately
func (su *StateUpdater) processUpdate(update StateUpdate) error {
	// Set global context for the update
	if update.Context != nil {
		SetGlobalContext(update.Context)
		defer SetGlobalContext(nil)
	}

	// Execute the update function
	if update.UpdateFunc != nil {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("setState update function panic: %v\n", r)
				}
			}()
			update.UpdateFunc()
		}()
	}

	// Trigger rebuilds for affected widgets
	if len(update.WidgetIDs) > 0 {
		su.TriggerRebuild(update.WidgetIDs)
	}

	// Broadcast state change event
	su.broadcastStateChange([]StateUpdate{update}, update.WidgetIDs)

	return nil
}

// TriggerRebuild triggers rebuilds for specific widget IDs
func (su *StateUpdater) TriggerRebuild(widgetIDs []string) {
	if su.context == nil || su.context.App == nil {
		return
	}

	// For now, we'll broadcast a generic rebuild message
	// In a more sophisticated implementation, we would track which widgets
	// need rebuilding and send targeted updates

	if su.context.App.WebSocket().IsEnabled() {
		message := map[string]interface{}{
			"type":      "widget_rebuild",
			"widgetIDs": widgetIDs,
			"timestamp": time.Now().Unix(),
		}

		// Broadcast to all clients - in practice you might want to be more selective
		su.context.App.WebSocket().Broadcast("state_update", message)
	}
}

// broadcastStateChange broadcasts state change events via WebSocket
func (su *StateUpdater) broadcastStateChange(batch []StateUpdate, affectedWidgetIDs []string) {
	if su.context == nil || su.context.App == nil || !su.context.App.WebSocket().IsEnabled() {
		return
	}

	message := map[string]interface{}{
		"type":            "state_batch_update",
		"batchSize":       len(batch),
		"affectedWidgets": affectedWidgetIDs,
		"timestamp":       time.Now().Unix(),
	}

	su.context.App.WebSocket().Broadcast("state_update", message)
}

// IsProcessing returns true if the state updater is currently processing updates
func (su *StateUpdater) IsProcessing() bool {
	su.mutex.RLock()
	defer su.mutex.RUnlock()
	return su.isProcessing
}

// SetBatchDelay sets the delay for batching state updates
func (su *StateUpdater) SetBatchDelay(delay time.Duration) {
	su.mutex.Lock()
	defer su.mutex.Unlock()
	su.batchDelay = delay
}

// GetStats returns statistics about the state updater
func (su *StateUpdater) GetStats() map[string]interface{} {
	su.mutex.RLock()
	defer su.mutex.RUnlock()

	return map[string]interface{}{
		"queueLength":  len(su.updateQueue),
		"isProcessing": su.isProcessing,
		"batchDelay":   su.batchDelay.String(),
	}
}

// Stop stops the state updater (for cleanup)
func (su *StateUpdater) Stop() {
	// Close the update queue to stop the processor goroutine
	close(su.updateQueue)
}

// Global convenience functions

// GetGlobalStateUpdater returns the global state updater instance
func GetGlobalStateUpdater() *StateUpdater {
	globalStateUpdaterMutex.RLock()
	defer globalStateUpdaterMutex.RUnlock()
	return globalStateUpdater
}

// SetGlobalStateUpdater sets the global state updater instance
func SetGlobalStateUpdater(updater *StateUpdater) {
	globalStateUpdaterMutex.Lock()
	defer globalStateUpdaterMutex.Unlock()
	globalStateUpdater = updater
}

// Helper functions for common setState patterns

// SetStateValue is a convenience function for setting a single state value
func SetStateValue(ctx *Context, key string, value interface{}) error {
	return SetStateFunc(ctx, func() {
		if ctx != nil && ctx.App != nil {
			ctx.App.State().Set(key, value)
		}
	})
}

// SetStateValues is a convenience function for setting multiple state values
func SetStateValues(ctx *Context, values map[string]interface{}) error {
	return SetStateFunc(ctx, func() {
		if ctx != nil && ctx.App != nil {
			stateManager := ctx.App.State()
			for key, value := range values {
				stateManager.Set(key, value)
			}
		}
	})
}

// SetStateWithCallback executes a setState with a completion callback
func SetStateWithCallback(ctx *Context, updateFunc func(), callback func()) error {
	return SetStateFunc(ctx, func() {
		updateFunc()
		if callback != nil {
			// Execute callback after a short delay to ensure state is propagated
			go func() {
				time.Sleep(10 * time.Millisecond)
				callback()
			}()
		}
	})
}

// Flutter-style setState function that can be called from widget callbacks
// This provides the familiar setState(() => { ... }) pattern
func SetStateCallback(ctx *Context, callback func()) error {
	return SetStateFunc(ctx, callback)
}

// SetStateWithRebuild executes a state update and rebuilds specific widgets
func SetStateWithRebuild(ctx *Context, updateFunc func(), widgetIDs ...string) error {
	return SetStateWithWidgets(ctx, updateFunc, widgetIDs)
}
