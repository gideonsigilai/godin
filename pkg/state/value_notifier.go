package state

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"
)

// ValueNotifier is a generic value holder that notifies listeners when the value changes
type ValueNotifier[T any] struct {
	value     T
	listeners []func(T)
	mutex     sync.RWMutex
	id        string
	manager   *StateManager
}

// NewValueNotifier creates a new ValueNotifier with an initial value
func NewValueNotifier[T any](initialValue T) *ValueNotifier[T] {
	return &ValueNotifier[T]{
		value:     initialValue,
		listeners: make([]func(T), 0),
		id:        generateID(),
	}
}

// NewValueNotifierWithID creates a new ValueNotifier with a specific ID
func NewValueNotifierWithID[T any](id string, initialValue T) *ValueNotifier[T] {
	return &ValueNotifier[T]{
		value:     initialValue,
		listeners: make([]func(T), 0),
		id:        id,
	}
}

// Value returns the current value
func (vn *ValueNotifier[T]) Value() T {
	vn.mutex.RLock()
	defer vn.mutex.RUnlock()
	return vn.value
}

// SetValue updates the value and notifies all listeners
func (vn *ValueNotifier[T]) SetValue(newValue T) {
	vn.mutex.Lock()
	oldValue := vn.value
	vn.value = newValue
	listeners := make([]func(T), len(vn.listeners))
	copy(listeners, vn.listeners)
	vn.mutex.Unlock()

	// Check if value actually changed
	if !reflect.DeepEqual(oldValue, newValue) {
		// Notify all listeners
		for _, listener := range listeners {
			go listener(newValue)
		}

		// Notify state manager if attached
		if vn.manager != nil {
			vn.manager.notifyValueChange(vn.id, newValue)
		}
	}
}

// AddListener adds a listener function that will be called when the value changes
func (vn *ValueNotifier[T]) AddListener(listener func(T)) {
	vn.mutex.Lock()
	defer vn.mutex.Unlock()
	vn.listeners = append(vn.listeners, listener)
}

// RemoveListener removes a specific listener (not implemented for simplicity)
// In a real implementation, you'd need to store listener IDs
func (vn *ValueNotifier[T]) RemoveListener(listener func(T)) {
	// This is complex to implement without listener IDs
	// For now, we'll provide ClearListeners instead
}

// ClearListeners removes all listeners
func (vn *ValueNotifier[T]) ClearListeners() {
	vn.mutex.Lock()
	defer vn.mutex.Unlock()
	vn.listeners = make([]func(T), 0)
}

// ListenerCount returns the number of active listeners
func (vn *ValueNotifier[T]) ListenerCount() int {
	vn.mutex.RLock()
	defer vn.mutex.RUnlock()
	return len(vn.listeners)
}

// ID returns the unique identifier for this ValueNotifier
func (vn *ValueNotifier[T]) ID() string {
	return vn.id
}

// SetManager sets the state manager for this ValueNotifier
func (vn *ValueNotifier[T]) SetManager(manager *StateManager) {
	vn.manager = manager
}

// ToJSON converts the current value to JSON
func (vn *ValueNotifier[T]) ToJSON() ([]byte, error) {
	vn.mutex.RLock()
	defer vn.mutex.RUnlock()
	return json.Marshal(vn.value)
}

// FromJSON updates the value from JSON
func (vn *ValueNotifier[T]) FromJSON(data []byte) error {
	var newValue T
	if err := json.Unmarshal(data, &newValue); err != nil {
		return err
	}
	vn.SetValue(newValue)
	return nil
}

// Update applies a function to the current value
func (vn *ValueNotifier[T]) Update(updater func(T) T) {
	vn.mutex.Lock()
	oldValue := vn.value
	newValue := updater(oldValue)
	vn.value = newValue
	listeners := make([]func(T), len(vn.listeners))
	copy(listeners, vn.listeners)
	vn.mutex.Unlock()

	// Check if value actually changed
	if !reflect.DeepEqual(oldValue, newValue) {
		// Notify all listeners
		for _, listener := range listeners {
			go listener(newValue)
		}

		// Notify state manager if attached
		if vn.manager != nil {
			vn.manager.notifyValueChange(vn.id, newValue)
		}
	}
}

// String returns a string representation of the current value
func (vn *ValueNotifier[T]) String() string {
	vn.mutex.RLock()
	defer vn.mutex.RUnlock()
	return fmt.Sprintf("ValueNotifier[%s](%v)", vn.id, vn.value)
}

// Helper function to generate unique IDs
func generateID() string {
	return fmt.Sprintf("vn_%d", time.Now().UnixNano())
}

// Specialized ValueNotifiers for common types

// IntNotifier is a ValueNotifier for int values
type IntNotifier = ValueNotifier[int]

// Float64Notifier is a ValueNotifier for float64 values
type Float64Notifier = ValueNotifier[float64]

// StringNotifier is a ValueNotifier for string values
type StringNotifier = ValueNotifier[string]

// BoolNotifier is a ValueNotifier for bool values
type BoolNotifier = ValueNotifier[bool]

// Convenience constructors
func NewIntNotifier(value int) *IntNotifier {
	return NewValueNotifier(value)
}

func NewFloat64Notifier(value float64) *Float64Notifier {
	return NewValueNotifier(value)
}

func NewStringNotifier(value string) *StringNotifier {
	return NewValueNotifier(value)
}

func NewBoolNotifier(value bool) *BoolNotifier {
	return NewValueNotifier(value)
}

// Convenience constructors with ID
func NewIntNotifierWithID(id string, value int) *IntNotifier {
	return NewValueNotifierWithID(id, value)
}

func NewFloat64NotifierWithID(id string, value float64) *Float64Notifier {
	return NewValueNotifierWithID(id, value)
}

func NewStringNotifierWithID(id string, value string) *StringNotifier {
	return NewValueNotifierWithID(id, value)
}

func NewBoolNotifierWithID(id string, value bool) *BoolNotifier {
	return NewValueNotifierWithID(id, value)
}
