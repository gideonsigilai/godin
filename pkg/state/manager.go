package state

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// WebSocketBroadcaster interface for broadcasting state changes
type WebSocketBroadcaster interface {
	Broadcast(channel string, data interface{})
}

// StateManager manages application state and notifications
type StateManager struct {
	data        map[string]interface{}
	watchers    map[string][]func(interface{})
	notifiers   map[string]interface{} // Store ValueNotifiers
	mutex       sync.RWMutex
	broadcaster WebSocketBroadcaster
	lastUpdated map[string]time.Time
}

// NewStateManager creates a new state manager
func NewStateManager() *StateManager {
	return &StateManager{
		data:        make(map[string]interface{}),
		watchers:    make(map[string][]func(interface{})),
		notifiers:   make(map[string]interface{}),
		lastUpdated: make(map[string]time.Time),
	}
}

// NewStateManagerWithBroadcaster creates a new state manager with WebSocket broadcaster
func NewStateManagerWithBroadcaster(broadcaster WebSocketBroadcaster) *StateManager {
	return &StateManager{
		data:        make(map[string]interface{}),
		watchers:    make(map[string][]func(interface{})),
		broadcaster: broadcaster,
	}
}

// SetBroadcaster sets the WebSocket broadcaster for real-time updates
func (sm *StateManager) SetBroadcaster(broadcaster WebSocketBroadcaster) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.broadcaster = broadcaster
}

// RegisterValueNotifier registers a ValueNotifier with the state manager
func (sm *StateManager) RegisterValueNotifier(id string, notifier interface{}) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.notifiers[id] = notifier
	sm.lastUpdated[id] = time.Now()
}

// GetValueNotifier retrieves a ValueNotifier by ID
func (sm *StateManager) GetValueNotifier(id string) (interface{}, bool) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	notifier, exists := sm.notifiers[id]
	return notifier, exists
}

// UnregisterValueNotifier removes a ValueNotifier from the state manager
func (sm *StateManager) UnregisterValueNotifier(id string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	delete(sm.notifiers, id)
	delete(sm.lastUpdated, id)
}

// notifyValueChange is called when a ValueNotifier's value changes
func (sm *StateManager) notifyValueChange(id string, value interface{}) {
	sm.mutex.Lock()
	sm.lastUpdated[id] = time.Now()
	sm.mutex.Unlock()

	// Broadcast the change via WebSocket if broadcaster is available
	if sm.broadcaster != nil {
		// Create the broadcast message
		message := map[string]interface{}{
			"type":      "value_change",
			"id":        id,
			"value":     value,
			"timestamp": time.Now().Unix(),
		}

		// Broadcast on the state channel for this specific notifier
		sm.broadcaster.Broadcast(fmt.Sprintf("state:%s", id), message)
	}

	// Notify local watchers
	sm.mutex.RLock()
	watchers, exists := sm.watchers[id]
	sm.mutex.RUnlock()

	if exists {
		for _, watcher := range watchers {
			go watcher(value)
		}
	}
}

// GetValueNotifierState returns the current state of a ValueNotifier
func (sm *StateManager) GetValueNotifierState(id string) (map[string]interface{}, error) {
	sm.mutex.RLock()
	notifier, exists := sm.notifiers[id]
	lastUpdated := sm.lastUpdated[id]
	sm.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("ValueNotifier with ID %s not found", id)
	}

	// Try to get the value using reflection or type assertion
	var value interface{}
	var jsonData []byte
	var err error

	// Try to call Value() method if it exists
	if valueGetter, ok := notifier.(interface{ Value() interface{} }); ok {
		value = valueGetter.Value()
	} else {
		value = notifier
	}

	// Try to serialize to JSON
	if jsonSerializer, ok := notifier.(interface{ ToJSON() ([]byte, error) }); ok {
		jsonData, err = jsonSerializer.ToJSON()
	} else {
		jsonData, err = json.Marshal(value)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to serialize value: %v", err)
	}

	return map[string]interface{}{
		"id":          id,
		"value":       value,
		"json":        string(jsonData),
		"lastUpdated": lastUpdated.Unix(),
		"type":        fmt.Sprintf("%T", notifier),
	}, nil
}

// ListValueNotifiers returns a list of all registered ValueNotifier IDs
func (sm *StateManager) ListValueNotifiers() []string {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	ids := make([]string, 0, len(sm.notifiers))
	for id := range sm.notifiers {
		ids = append(ids, id)
	}
	return ids
}

// Set sets a value in the state and notifies watchers
func (sm *StateManager) Set(key string, value interface{}) {
	sm.mutex.Lock()
	sm.data[key] = value
	watchers := sm.watchers[key]
	broadcaster := sm.broadcaster
	sm.mutex.Unlock()

	// Notify watchers
	for _, watcher := range watchers {
		go watcher(value)
	}

	// Broadcast state change via WebSocket for real-time UI updates
	if broadcaster != nil {
		go broadcaster.Broadcast("state:"+key, map[string]interface{}{
			"key":   key,
			"value": value,
		})
	}
}

// Get retrieves a value from the state
func (sm *StateManager) Get(key string) interface{} {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	return sm.data[key]
}

// GetString retrieves a string value from the state
func (sm *StateManager) GetString(key string) string {
	value := sm.Get(key)
	if str, ok := value.(string); ok {
		return str
	}
	return ""
}

// GetInt retrieves an integer value from the state
func (sm *StateManager) GetInt(key string) int {
	value := sm.Get(key)
	if i, ok := value.(int); ok {
		return i
	}
	return 0
}

// GetBool retrieves a boolean value from the state
func (sm *StateManager) GetBool(key string) bool {
	value := sm.Get(key)
	if b, ok := value.(bool); ok {
		return b
	}
	return false
}

// Watch creates a watcher for a state key
func (sm *StateManager) Watch(key string) *ValueListenable {
	return &ValueListenable{
		key:     key,
		manager: sm,
	}
}

// AddWatcher adds a watcher function for a state key
func (sm *StateManager) AddWatcher(key string, watcher func(interface{})) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.watchers[key] = append(sm.watchers[key], watcher)
}

// RemoveWatcher removes a watcher function for a state key
func (sm *StateManager) RemoveWatcher(key string, watcher func(interface{})) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	watchers := sm.watchers[key]
	for i, w := range watchers {
		// Note: This is a simplified comparison
		// In practice, you'd need a more sophisticated way to compare functions
		if &w == &watcher {
			sm.watchers[key] = append(watchers[:i], watchers[i+1:]...)
			break
		}
	}
}

// Delete removes a key from the state
func (sm *StateManager) Delete(key string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	delete(sm.data, key)
	delete(sm.watchers, key)
}

// Keys returns all keys in the state
func (sm *StateManager) Keys() []string {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	keys := make([]string, 0, len(sm.data))
	for key := range sm.data {
		keys = append(keys, key)
	}
	return keys
}

// Clear clears all state data
func (sm *StateManager) Clear() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.data = make(map[string]interface{})
	sm.watchers = make(map[string][]func(interface{}))
}

// ValueListenable represents a listenable value
type ValueListenable struct {
	key       string
	manager   *StateManager
	listeners []func()
}

// GetValue returns the current value
func (vl *ValueListenable) GetValue() interface{} {
	return vl.manager.Get(vl.key)
}

// AddListener adds a change listener
func (vl *ValueListenable) AddListener(listener func()) {
	vl.listeners = append(vl.listeners, listener)

	// Add watcher to state manager
	vl.manager.AddWatcher(vl.key, func(value interface{}) {
		listener()
	})
}

// RemoveListener removes a change listener
func (vl *ValueListenable) RemoveListener(listener func()) {
	for i, l := range vl.listeners {
		if &l == &listener {
			vl.listeners = append(vl.listeners[:i], vl.listeners[i+1:]...)
			break
		}
	}
}

// State represents a reactive state container
type State struct {
	manager *StateManager
	prefix  string
}

// NewState creates a new state container with optional prefix
func NewState(manager *StateManager, prefix string) *State {
	return &State{
		manager: manager,
		prefix:  prefix,
	}
}

// Set sets a value in the state
func (s *State) Set(key string, value interface{}) {
	fullKey := s.getFullKey(key)
	s.manager.Set(fullKey, value)
}

// Get retrieves a value from the state
func (s *State) Get(key string) interface{} {
	fullKey := s.getFullKey(key)
	return s.manager.Get(fullKey)
}

// Watch creates a watcher for a state key
func (s *State) Watch(key string) *ValueListenable {
	fullKey := s.getFullKey(key)
	return s.manager.Watch(fullKey)
}

// getFullKey returns the full key with prefix
func (s *State) getFullKey(key string) string {
	if s.prefix == "" {
		return key
	}
	return s.prefix + "." + key
}
