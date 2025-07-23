package core

import (
	"fmt"
	"runtime"
	"time"
)

// CallbackError represents an error that occurred during callback execution
type CallbackError struct {
	CallbackID  string      // ID of the callback that failed
	WidgetType  string      // Type of widget that triggered the callback
	WidgetID    string      // ID of the widget that triggered the callback
	EventType   string      // Type of event (OnPressed, OnChanged, etc.)
	OriginalErr error       // The original error that occurred
	Timestamp   time.Time   // When the error occurred
	StackTrace  string      // Stack trace at the time of error
	Context     interface{} // Additional context information
	Recoverable bool        // Whether the error is recoverable
}

// Error implements the error interface
func (ce *CallbackError) Error() string {
	return fmt.Sprintf("callback error [%s:%s] in %s widget '%s': %v",
		ce.EventType, ce.CallbackID, ce.WidgetType, ce.WidgetID, ce.OriginalErr)
}

// String returns a detailed string representation of the error
func (ce *CallbackError) String() string {
	return fmt.Sprintf(`CallbackError{
	CallbackID: %s
	WidgetType: %s
	WidgetID: %s
	EventType: %s
	Error: %v
	Timestamp: %s
	Recoverable: %t
	StackTrace: %s
}`, ce.CallbackID, ce.WidgetType, ce.WidgetID, ce.EventType,
		ce.OriginalErr, ce.Timestamp.Format(time.RFC3339), ce.Recoverable, ce.StackTrace)
}

// StateError represents an error that occurred during state management
type StateError struct {
	Operation   string      // The operation that failed (Set, Get, Watch, etc.)
	Key         string      // The state key involved
	Value       interface{} // The value involved (if applicable)
	OriginalErr error       // The original error that occurred
	Timestamp   time.Time   // When the error occurred
	StackTrace  string      // Stack trace at the time of error
	Context     interface{} // Additional context information
	Recoverable bool        // Whether the error is recoverable
}

// Error implements the error interface
func (se *StateError) Error() string {
	return fmt.Sprintf("state error [%s] on key '%s': %v",
		se.Operation, se.Key, se.OriginalErr)
}

// String returns a detailed string representation of the error
func (se *StateError) String() string {
	return fmt.Sprintf(`StateError{
	Operation: %s
	Key: %s
	Value: %v
	Error: %v
	Timestamp: %s
	Recoverable: %t
	StackTrace: %s
}`, se.Operation, se.Key, se.Value,
		se.OriginalErr, se.Timestamp.Format(time.RFC3339), se.Recoverable, se.StackTrace)
}

// ErrorHandler defines the interface for handling errors
type ErrorHandler interface {
	HandleCallbackError(error *CallbackError)
	HandleStateError(error *StateError)
	HandleGenericError(error error, context interface{})
}

// DefaultErrorHandler provides default error handling behavior
type DefaultErrorHandler struct {
	Logger ErrorLogger
}

// ErrorLogger defines the interface for logging errors
type ErrorLogger interface {
	LogError(level string, message string, error error, context interface{})
}

// NewDefaultErrorHandler creates a new default error handler
func NewDefaultErrorHandler(logger ErrorLogger) *DefaultErrorHandler {
	return &DefaultErrorHandler{
		Logger: logger,
	}
}

// HandleCallbackError handles callback errors
func (deh *DefaultErrorHandler) HandleCallbackError(err *CallbackError) {
	if deh.Logger != nil {
		deh.Logger.LogError("ERROR", "Callback execution failed", err, map[string]interface{}{
			"callback_id": err.CallbackID,
			"widget_type": err.WidgetType,
			"widget_id":   err.WidgetID,
			"event_type":  err.EventType,
			"recoverable": err.Recoverable,
		})
	}

	// If the error is recoverable, attempt recovery
	if err.Recoverable {
		deh.attemptCallbackRecovery(err)
	}
}

// HandleStateError handles state management errors
func (deh *DefaultErrorHandler) HandleStateError(err *StateError) {
	if deh.Logger != nil {
		deh.Logger.LogError("ERROR", "State operation failed", err, map[string]interface{}{
			"operation":   err.Operation,
			"key":         err.Key,
			"value":       err.Value,
			"recoverable": err.Recoverable,
		})
	}

	// If the error is recoverable, attempt recovery
	if err.Recoverable {
		deh.attemptStateRecovery(err)
	}
}

// HandleGenericError handles generic errors
func (deh *DefaultErrorHandler) HandleGenericError(err error, context interface{}) {
	if deh.Logger != nil {
		deh.Logger.LogError("ERROR", "Generic error occurred", err, context)
	}
}

// attemptCallbackRecovery attempts to recover from callback errors
func (deh *DefaultErrorHandler) attemptCallbackRecovery(err *CallbackError) {
	// Log recovery attempt
	if deh.Logger != nil {
		deh.Logger.LogError("INFO", "Attempting callback recovery", nil, map[string]interface{}{
			"callback_id": err.CallbackID,
			"widget_type": err.WidgetType,
		})
	}

	// Recovery strategies could include:
	// 1. Retry the callback with exponential backoff
	// 2. Reset the widget to a safe state
	// 3. Disable the callback temporarily
	// 4. Show user-friendly error message
}

// attemptStateRecovery attempts to recover from state errors
func (deh *DefaultErrorHandler) attemptStateRecovery(err *StateError) {
	// Log recovery attempt
	if deh.Logger != nil {
		deh.Logger.LogError("INFO", "Attempting state recovery", nil, map[string]interface{}{
			"operation": err.Operation,
			"key":       err.Key,
		})
	}

	// Recovery strategies could include:
	// 1. Reset the state key to a default value
	// 2. Clear corrupted state data
	// 3. Reinitialize the state manager
	// 4. Fallback to local storage
}

// ErrorRecoveryManager manages error recovery strategies
type ErrorRecoveryManager struct {
	strategies map[string]RecoveryStrategy
	maxRetries int
	retryDelay time.Duration
}

// RecoveryStrategy defines a strategy for recovering from errors
type RecoveryStrategy interface {
	CanRecover(err error) bool
	Recover(err error, context interface{}) error
}

// NewErrorRecoveryManager creates a new error recovery manager
func NewErrorRecoveryManager() *ErrorRecoveryManager {
	return &ErrorRecoveryManager{
		strategies: make(map[string]RecoveryStrategy),
		maxRetries: 3,
		retryDelay: time.Second,
	}
}

// RegisterStrategy registers a recovery strategy
func (erm *ErrorRecoveryManager) RegisterStrategy(name string, strategy RecoveryStrategy) {
	erm.strategies[name] = strategy
}

// Recover attempts to recover from an error using registered strategies
func (erm *ErrorRecoveryManager) Recover(err error, context interface{}) error {
	for _, strategy := range erm.strategies {
		if strategy.CanRecover(err) {
			for i := 0; i < erm.maxRetries; i++ {
				if recoveryErr := strategy.Recover(err, context); recoveryErr == nil {
					return nil // Recovery successful
				}

				// Wait before retrying
				if i < erm.maxRetries-1 {
					time.Sleep(erm.retryDelay * time.Duration(i+1))
				}
			}
		}
	}

	return fmt.Errorf("failed to recover from error: %v", err)
}

// CallbackRecoveryStrategy implements recovery for callback errors
type CallbackRecoveryStrategy struct{}

// CanRecover checks if this strategy can recover from the error
func (crs *CallbackRecoveryStrategy) CanRecover(err error) bool {
	_, ok := err.(*CallbackError)
	return ok
}

// Recover attempts to recover from a callback error
func (crs *CallbackRecoveryStrategy) Recover(err error, context interface{}) error {
	callbackErr, ok := err.(*CallbackError)
	if !ok {
		return fmt.Errorf("invalid error type for callback recovery")
	}

	// Implement specific recovery logic based on error type
	switch callbackErr.EventType {
	case "OnPressed":
		// For button press errors, we might disable the button temporarily
		return crs.recoverButtonPress(callbackErr, context)
	case "OnChanged":
		// For change events, we might reset to previous value
		return crs.recoverValueChange(callbackErr, context)
	default:
		return fmt.Errorf("no recovery strategy for event type: %s", callbackErr.EventType)
	}
}

// recoverButtonPress recovers from button press errors
func (crs *CallbackRecoveryStrategy) recoverButtonPress(err *CallbackError, context interface{}) error {
	// Implementation would depend on the specific context
	// For now, we'll just log the recovery attempt
	return nil
}

// recoverValueChange recovers from value change errors
func (crs *CallbackRecoveryStrategy) recoverValueChange(err *CallbackError, context interface{}) error {
	// Implementation would depend on the specific context
	// For now, we'll just log the recovery attempt
	return nil
}

// StateRecoveryStrategy implements recovery for state errors
type StateRecoveryStrategy struct{}

// CanRecover checks if this strategy can recover from the error
func (srs *StateRecoveryStrategy) CanRecover(err error) bool {
	_, ok := err.(*StateError)
	return ok
}

// Recover attempts to recover from a state error
func (srs *StateRecoveryStrategy) Recover(err error, context interface{}) error {
	stateErr, ok := err.(*StateError)
	if !ok {
		return fmt.Errorf("invalid error type for state recovery")
	}

	// Implement specific recovery logic based on operation
	switch stateErr.Operation {
	case "Set":
		return srs.recoverSetOperation(stateErr, context)
	case "Get":
		return srs.recoverGetOperation(stateErr, context)
	case "Watch":
		return srs.recoverWatchOperation(stateErr, context)
	default:
		return fmt.Errorf("no recovery strategy for operation: %s", stateErr.Operation)
	}
}

// recoverSetOperation recovers from state set errors
func (srs *StateRecoveryStrategy) recoverSetOperation(err *StateError, context interface{}) error {
	// Implementation would depend on the specific context
	return nil
}

// recoverGetOperation recovers from state get errors
func (srs *StateRecoveryStrategy) recoverGetOperation(err *StateError, context interface{}) error {
	// Implementation would depend on the specific context
	return nil
}

// recoverWatchOperation recovers from state watch errors
func (srs *StateRecoveryStrategy) recoverWatchOperation(err *StateError, context interface{}) error {
	// Implementation would depend on the specific context
	return nil
}

// CreateCallbackError creates a new callback error with stack trace
func CreateCallbackError(callbackID, widgetType, widgetID, eventType string, originalErr error, recoverable bool) *CallbackError {
	// Capture stack trace
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	stackTrace := string(buf[:n])

	return &CallbackError{
		CallbackID:  callbackID,
		WidgetType:  widgetType,
		WidgetID:    widgetID,
		EventType:   eventType,
		OriginalErr: originalErr,
		Timestamp:   time.Now(),
		StackTrace:  stackTrace,
		Recoverable: recoverable,
	}
}

// CreateStateError creates a new state error with stack trace
func CreateStateError(operation, key string, value interface{}, originalErr error, recoverable bool) *StateError {
	// Capture stack trace
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	stackTrace := string(buf[:n])

	return &StateError{
		Operation:   operation,
		Key:         key,
		Value:       value,
		OriginalErr: originalErr,
		Timestamp:   time.Now(),
		StackTrace:  stackTrace,
		Recoverable: recoverable,
	}
}

// SafeExecuteCallback executes a callback with error handling and recovery
func SafeExecuteCallback(callbackID, widgetType, widgetID, eventType string, callback func(), errorHandler ErrorHandler) {
	defer func() {
		if r := recover(); r != nil {
			var err error
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("panic: %v", r)
			}

			callbackErr := CreateCallbackError(callbackID, widgetType, widgetID, eventType, err, true)
			if errorHandler != nil {
				errorHandler.HandleCallbackError(callbackErr)
			}
		}
	}()

	// Execute the callback
	if callback != nil {
		callback()
	}
}

// SafeExecuteStateOperation executes a state operation with error handling
func SafeExecuteStateOperation(operation, key string, value interface{}, operation_func func() error, errorHandler ErrorHandler) error {
	defer func() {
		if r := recover(); r != nil {
			var err error
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("panic: %v", r)
			}

			stateErr := CreateStateError(operation, key, value, err, true)
			if errorHandler != nil {
				errorHandler.HandleStateError(stateErr)
			}
		}
	}()

	// Execute the operation
	if operation_func != nil {
		return operation_func()
	}
	return nil
}
