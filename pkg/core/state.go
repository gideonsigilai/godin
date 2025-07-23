package core

import (
	"context"
	"fmt"
	"sync"
)

// Global state management for button callbacks and other native Go code

var (
	globalStateManager *GlobalStateManager
	globalStateMutex   sync.RWMutex
)

// GlobalStateManager manages global state for native Go code execution
type GlobalStateManager struct {
	currentContext *Context
	mutex          sync.RWMutex
}

// SetCurrentContext sets the current context for state operations
func (gsm *GlobalStateManager) SetCurrentContext(ctx *Context) {
	gsm.mutex.Lock()
	defer gsm.mutex.Unlock()
	gsm.currentContext = ctx
}

// GetCurrentContext gets the current context
func (gsm *GlobalStateManager) GetCurrentContext() *Context {
	gsm.mutex.RLock()
	defer gsm.mutex.RUnlock()
	return gsm.currentContext
}

// InitGlobalState initializes the global state manager
func InitGlobalState() {
	globalStateMutex.Lock()
	defer globalStateMutex.Unlock()

	if globalStateManager == nil {
		globalStateManager = &GlobalStateManager{}
	}
}

// SetGlobalContext sets the current context for global state operations
func SetGlobalContext(ctx *Context) {
	globalStateMutex.RLock()
	defer globalStateMutex.RUnlock()

	if globalStateManager != nil {
		globalStateManager.SetCurrentContext(ctx)
	}
}

// setState is the global function that can be called from button callbacks
func SetState(key string, value interface{}) {
	fmt.Printf("SetState called: %s = %v\n", key, value)

	globalStateMutex.RLock()
	defer globalStateMutex.RUnlock()

	if globalStateManager != nil {
		ctx := globalStateManager.GetCurrentContext()
		if ctx != nil {
			fmt.Printf("Setting state in context: %s = %v\n", key, value)
			ctx.SetState(key, value)
		} else {
			fmt.Printf("No current context available for SetState\n")
		}
	} else {
		fmt.Printf("Global state manager is nil\n")
	}
}

// GetState is the global function to get state values
func GetState(key string) interface{} {
	globalStateMutex.RLock()
	defer globalStateMutex.RUnlock()

	if globalStateManager != nil {
		ctx := globalStateManager.GetCurrentContext()
		if ctx != nil {
			return ctx.GetState(key)
		}
	}
	return nil
}

// GetStateString retrieves a string value from global state
func GetStateString(key string) string {
	if value, ok := GetState(key).(string); ok {
		return value
	}
	return ""
}

// GetStateInt retrieves an integer value from global state
func GetStateInt(key string) int {
	if value, ok := GetState(key).(int); ok {
		return value
	}
	return 0
}

// GetStateBool retrieves a boolean value from global state
func GetStateBool(key string) bool {
	if value, ok := GetState(key).(bool); ok {
		return value
	}
	return false
}

// ContextKey is used for context values
type ContextKey string

const (
	// StateContextKey is the key for storing state context
	StateContextKey ContextKey = "godin_state_context"
)

// WithStateContext creates a new context with state context
func WithStateContext(ctx context.Context, stateCtx *Context) context.Context {
	return context.WithValue(ctx, StateContextKey, stateCtx)
}

// StateContextFromContext retrieves the state context from a context
func StateContextFromContext(ctx context.Context) *Context {
	if stateCtx, ok := ctx.Value(StateContextKey).(*Context); ok {
		return stateCtx
	}
	return nil
}
