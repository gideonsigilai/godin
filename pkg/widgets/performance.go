package widgets

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/gideonsigilai/godin/pkg/core"
	"github.com/gideonsigilai/godin/pkg/state"
)

// WidgetRegistry manages widget lifecycle and prevents memory leaks
type WidgetRegistry struct {
	widgets    map[string]*WidgetInfo
	listeners  map[string][]*ListenerInfo
	mutex      sync.RWMutex
	gcInterval time.Duration
	stopGC     chan bool
}

// WidgetInfo contains information about a registered widget
type WidgetInfo struct {
	ID          string
	Type        string
	NotifierIDs []string
	LastRender  time.Time
	RenderCount int
	CreatedAt   time.Time
	LastAccess  time.Time
	IsActive    bool
}

// ListenerInfo contains information about a registered listener
type ListenerInfo struct {
	ID         string
	WidgetID   string
	NotifierID string
	CreatedAt  time.Time
	LastUsed   time.Time
	IsActive   bool
}

// NewWidgetRegistry creates a new widget registry
func NewWidgetRegistry() *WidgetRegistry {
	registry := &WidgetRegistry{
		widgets:    make(map[string]*WidgetInfo),
		listeners:  make(map[string][]*ListenerInfo),
		gcInterval: time.Minute * 5, // Run garbage collection every 5 minutes
		stopGC:     make(chan bool),
	}

	// Start garbage collection goroutine
	go registry.startGarbageCollection()

	return registry
}

// RegisterWidget registers a widget in the registry
func (wr *WidgetRegistry) RegisterWidget(id, widgetType string, notifierIDs []string) {
	wr.mutex.Lock()
	defer wr.mutex.Unlock()

	now := time.Now()
	wr.widgets[id] = &WidgetInfo{
		ID:          id,
		Type:        widgetType,
		NotifierIDs: notifierIDs,
		LastRender:  now,
		RenderCount: 0,
		CreatedAt:   now,
		LastAccess:  now,
		IsActive:    true,
	}
}

// RegisterListener registers a listener in the registry
func (wr *WidgetRegistry) RegisterListener(listenerID, widgetID, notifierID string) {
	wr.mutex.Lock()
	defer wr.mutex.Unlock()

	now := time.Now()
	listener := &ListenerInfo{
		ID:         listenerID,
		WidgetID:   widgetID,
		NotifierID: notifierID,
		CreatedAt:  now,
		LastUsed:   now,
		IsActive:   true,
	}

	if wr.listeners[notifierID] == nil {
		wr.listeners[notifierID] = make([]*ListenerInfo, 0)
	}
	wr.listeners[notifierID] = append(wr.listeners[notifierID], listener)
}

// UnregisterWidget removes a widget from the registry
func (wr *WidgetRegistry) UnregisterWidget(id string) {
	wr.mutex.Lock()
	defer wr.mutex.Unlock()

	// Mark widget as inactive
	if widget, exists := wr.widgets[id]; exists {
		widget.IsActive = false
	}

	// Mark associated listeners as inactive
	for notifierID, listeners := range wr.listeners {
		for _, listener := range listeners {
			if listener.WidgetID == id {
				listener.IsActive = false
			}
		}
		// Clean up empty listener arrays
		activeListeners := make([]*ListenerInfo, 0)
		for _, listener := range listeners {
			if listener.IsActive {
				activeListeners = append(activeListeners, listener)
			}
		}
		wr.listeners[notifierID] = activeListeners
	}
}

// UnregisterListener removes a listener from the registry
func (wr *WidgetRegistry) UnregisterListener(listenerID string) {
	wr.mutex.Lock()
	defer wr.mutex.Unlock()

	for notifierID, listeners := range wr.listeners {
		for i, listener := range listeners {
			if listener.ID == listenerID {
				listener.IsActive = false
				// Remove from slice
				wr.listeners[notifierID] = append(listeners[:i], listeners[i+1:]...)
				break
			}
		}
	}
}

// UpdateWidgetAccess updates the last access time for a widget
func (wr *WidgetRegistry) UpdateWidgetAccess(id string) {
	wr.mutex.Lock()
	defer wr.mutex.Unlock()

	if widget, exists := wr.widgets[id]; exists {
		widget.LastAccess = time.Now()
		widget.RenderCount++
	}
}

// UpdateListenerAccess updates the last used time for a listener
func (wr *WidgetRegistry) UpdateListenerAccess(listenerID string) {
	wr.mutex.Lock()
	defer wr.mutex.Unlock()

	for _, listeners := range wr.listeners {
		for _, listener := range listeners {
			if listener.ID == listenerID {
				listener.LastUsed = time.Now()
				break
			}
		}
	}
}

// GetWidgetInfo returns information about a widget
func (wr *WidgetRegistry) GetWidgetInfo(id string) (*WidgetInfo, bool) {
	wr.mutex.RLock()
	defer wr.mutex.RUnlock()

	widget, exists := wr.widgets[id]
	return widget, exists
}

// GetActiveWidgetCount returns the number of active widgets
func (wr *WidgetRegistry) GetActiveWidgetCount() int {
	wr.mutex.RLock()
	defer wr.mutex.RUnlock()

	count := 0
	for _, widget := range wr.widgets {
		if widget.IsActive {
			count++
		}
	}
	return count
}

// GetActiveListenerCount returns the number of active listeners
func (wr *WidgetRegistry) GetActiveListenerCount() int {
	wr.mutex.RLock()
	defer wr.mutex.RUnlock()

	count := 0
	for _, listeners := range wr.listeners {
		for _, listener := range listeners {
			if listener.IsActive {
				count++
			}
		}
	}
	return count
}

// GetMemoryStats returns memory statistics
func (wr *WidgetRegistry) GetMemoryStats() runtime.MemStats {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	return stats
}

// startGarbageCollection starts the garbage collection process
func (wr *WidgetRegistry) startGarbageCollection() {
	ticker := time.NewTicker(wr.gcInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			wr.runGarbageCollection()
		case <-wr.stopGC:
			return
		}
	}
}

// runGarbageCollection performs garbage collection of inactive widgets and listeners
func (wr *WidgetRegistry) runGarbageCollection() {
	wr.mutex.Lock()
	defer wr.mutex.Unlock()

	now := time.Now()
	gcThreshold := time.Minute * 10 // Remove widgets/listeners inactive for 10 minutes

	// Clean up inactive widgets
	for id, widget := range wr.widgets {
		if !widget.IsActive && now.Sub(widget.LastAccess) > gcThreshold {
			delete(wr.widgets, id)
		}
	}

	// Clean up inactive listeners
	for notifierID, listeners := range wr.listeners {
		activeListeners := make([]*ListenerInfo, 0)
		for _, listener := range listeners {
			if listener.IsActive || now.Sub(listener.LastUsed) <= gcThreshold {
				activeListeners = append(activeListeners, listener)
			}
		}
		wr.listeners[notifierID] = activeListeners

		// Remove empty listener arrays
		if len(activeListeners) == 0 {
			delete(wr.listeners, notifierID)
		}
	}

	// Force Go garbage collection
	runtime.GC()
}

// Stop stops the garbage collection process
func (wr *WidgetRegistry) Stop() {
	close(wr.stopGC)
}

// ChangeDetector optimizes rendering by detecting actual changes
type ChangeDetector struct {
	lastValues map[string]interface{}
	mutex      sync.RWMutex
}

// NewChangeDetector creates a new change detector
func NewChangeDetector() *ChangeDetector {
	return &ChangeDetector{
		lastValues: make(map[string]interface{}),
	}
}

// HasChanged checks if a value has changed since the last check
func (cd *ChangeDetector) HasChanged(key string, newValue interface{}) bool {
	cd.mutex.Lock()
	defer cd.mutex.Unlock()

	lastValue, exists := cd.lastValues[key]
	if !exists {
		cd.lastValues[key] = newValue
		return true
	}

	// Simple comparison - in a more sophisticated implementation,
	// you might want to use deep equality or custom comparison functions
	changed := !isEqual(lastValue, newValue)
	if changed {
		cd.lastValues[key] = newValue
	}

	return changed
}

// Clear removes a value from the change detector
func (cd *ChangeDetector) Clear(key string) {
	cd.mutex.Lock()
	defer cd.mutex.Unlock()
	delete(cd.lastValues, key)
}

// ClearAll removes all values from the change detector
func (cd *ChangeDetector) ClearAll() {
	cd.mutex.Lock()
	defer cd.mutex.Unlock()
	cd.lastValues = make(map[string]interface{})
}

// isEqual performs a simple equality check
func isEqual(a, b interface{}) bool {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

// RenderBatcher batches multiple state changes to optimize rendering
type RenderBatcher struct {
	pendingUpdates map[string]interface{}
	batchInterval  time.Duration
	mutex          sync.RWMutex
	timer          *time.Timer
	callback       func(map[string]interface{})
}

// NewRenderBatcher creates a new render batcher
func NewRenderBatcher(batchInterval time.Duration, callback func(map[string]interface{})) *RenderBatcher {
	return &RenderBatcher{
		pendingUpdates: make(map[string]interface{}),
		batchInterval:  batchInterval,
		callback:       callback,
	}
}

// AddUpdate adds an update to the batch
func (rb *RenderBatcher) AddUpdate(key string, value interface{}) {
	rb.mutex.Lock()
	defer rb.mutex.Unlock()

	rb.pendingUpdates[key] = value

	// Reset or start the timer
	if rb.timer != nil {
		rb.timer.Stop()
	}

	rb.timer = time.AfterFunc(rb.batchInterval, func() {
		rb.flush()
	})
}

// flush processes all pending updates
func (rb *RenderBatcher) flush() {
	rb.mutex.Lock()
	defer rb.mutex.Unlock()

	if len(rb.pendingUpdates) == 0 {
		return
	}

	// Make a copy of pending updates
	updates := make(map[string]interface{})
	for k, v := range rb.pendingUpdates {
		updates[k] = v
	}

	// Clear pending updates
	rb.pendingUpdates = make(map[string]interface{})

	// Call the callback with the batched updates
	if rb.callback != nil {
		go rb.callback(updates)
	}
}

// Flush immediately processes all pending updates
func (rb *RenderBatcher) Flush() {
	if rb.timer != nil {
		rb.timer.Stop()
	}
	rb.flush()
}

// PerformanceMonitor monitors widget performance
type PerformanceMonitor struct {
	renderTimes map[string][]time.Duration
	mutex       sync.RWMutex
	maxSamples  int
}

// NewPerformanceMonitor creates a new performance monitor
func NewPerformanceMonitor(maxSamples int) *PerformanceMonitor {
	return &PerformanceMonitor{
		renderTimes: make(map[string][]time.Duration),
		maxSamples:  maxSamples,
	}
}

// RecordRenderTime records the render time for a widget
func (pm *PerformanceMonitor) RecordRenderTime(widgetID string, duration time.Duration) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	if pm.renderTimes[widgetID] == nil {
		pm.renderTimes[widgetID] = make([]time.Duration, 0, pm.maxSamples)
	}

	times := pm.renderTimes[widgetID]
	times = append(times, duration)

	// Keep only the most recent samples
	if len(times) > pm.maxSamples {
		times = times[len(times)-pm.maxSamples:]
	}

	pm.renderTimes[widgetID] = times
}

// GetAverageRenderTime returns the average render time for a widget
func (pm *PerformanceMonitor) GetAverageRenderTime(widgetID string) time.Duration {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	times := pm.renderTimes[widgetID]
	if len(times) == 0 {
		return 0
	}

	var total time.Duration
	for _, t := range times {
		total += t
	}

	return total / time.Duration(len(times))
}

// GetPerformanceStats returns performance statistics for all widgets
func (pm *PerformanceMonitor) GetPerformanceStats() map[string]PerformanceStats {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	stats := make(map[string]PerformanceStats)
	for widgetID, times := range pm.renderTimes {
		if len(times) == 0 {
			continue
		}

		var total, min, max time.Duration
		min = times[0]
		max = times[0]

		for _, t := range times {
			total += t
			if t < min {
				min = t
			}
			if t > max {
				max = t
			}
		}

		stats[widgetID] = PerformanceStats{
			WidgetID:    widgetID,
			SampleCount: len(times),
			Average:     total / time.Duration(len(times)),
			Min:         min,
			Max:         max,
		}
	}

	return stats
}

// PerformanceStats contains performance statistics for a widget
type PerformanceStats struct {
	WidgetID    string
	SampleCount int
	Average     time.Duration
	Min         time.Duration
	Max         time.Duration
}

// OptimizedValueListener is a ValueListener with performance optimizations
type OptimizedValueListener[T any] struct {
	*ValueListener[T]
	changeDetector *ChangeDetector
	registry       *WidgetRegistry
	monitor        *PerformanceMonitor
	batcher        *RenderBatcher
}

// NewOptimizedValueListener creates a new optimized ValueListener
func NewOptimizedValueListener[T any](
	valueNotifier *state.ValueNotifier[T],
	builder func(value T) Widget,
	registry *WidgetRegistry,
	monitor *PerformanceMonitor,
) *OptimizedValueListener[T] {
	base := NewValueListener(valueNotifier, builder)

	return &OptimizedValueListener[T]{
		ValueListener:  base,
		changeDetector: NewChangeDetector(),
		registry:       registry,
		monitor:        monitor,
	}
}

// Render renders the optimized ValueListener with performance monitoring
func (ovl *OptimizedValueListener[T]) Render(ctx *core.Context) string {
	startTime := time.Now()
	defer func() {
		if ovl.monitor != nil {
			duration := time.Since(startTime)
			ovl.monitor.RecordRenderTime(ovl.ID, duration)
		}
	}()

	// Update registry access
	if ovl.registry != nil {
		ovl.registry.UpdateWidgetAccess(ovl.ID)
	}

	// Check if value has actually changed
	if ovl.ValueNotifier != nil {
		currentValue := ovl.ValueNotifier.Value()
		changeKey := fmt.Sprintf("%s_%s", ovl.ID, ovl.ValueNotifier.ID())

		if !ovl.changeDetector.HasChanged(changeKey, currentValue) {
			// Value hasn't changed, return cached result or skip render
			// In a more sophisticated implementation, you might cache the HTML result
			return ""
		}
	}

	// Perform the actual render
	return ovl.ValueListener.Render(ctx)
}

// Cleanup performs cleanup with registry integration
func (ovl *OptimizedValueListener[T]) Cleanup() {
	// Call base cleanup
	ovl.ValueListener.Cleanup()

	// Unregister from registry
	if ovl.registry != nil {
		ovl.registry.UnregisterWidget(ovl.ID)
	}

	// Clear change detector
	if ovl.changeDetector != nil {
		changeKey := fmt.Sprintf("%s_%s", ovl.ID, ovl.ValueNotifier.ID())
		ovl.changeDetector.Clear(changeKey)
	}
}

// Global instances for performance optimization
var (
	GlobalWidgetRegistry     *WidgetRegistry
	GlobalPerformanceMonitor *PerformanceMonitor
	GlobalRenderBatcher      *RenderBatcher
)

// InitializePerformanceOptimizations initializes global performance optimization components
func InitializePerformanceOptimizations() {
	GlobalWidgetRegistry = NewWidgetRegistry()
	GlobalPerformanceMonitor = NewPerformanceMonitor(100) // Keep last 100 samples

	// Batch renders every 16ms (roughly 60fps)
	GlobalRenderBatcher = NewRenderBatcher(time.Millisecond*16, func(updates map[string]interface{}) {
		// Handle batched updates - this would integrate with your WebSocket system
		for key, value := range updates {
			fmt.Printf("Batched update: %s = %v\n", key, value)
		}
	})
}

// CleanupPerformanceOptimizations cleans up global performance optimization components
func CleanupPerformanceOptimizations() {
	if GlobalWidgetRegistry != nil {
		GlobalWidgetRegistry.Stop()
	}
	if GlobalRenderBatcher != nil {
		GlobalRenderBatcher.Flush()
	}
}
