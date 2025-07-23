package widgets

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/gideonsigilai/godin/pkg/core"
)

func TestDefaultErrorRecoveryStrategy(t *testing.T) {
	strategy := DefaultErrorRecoveryStrategy()

	if strategy.MaxRetries != 3 {
		t.Errorf("Expected MaxRetries to be 3, got %d", strategy.MaxRetries)
	}

	if strategy.RetryInterval != time.Second*2 {
		t.Errorf("Expected RetryInterval to be 2s, got %v", strategy.RetryInterval)
	}

	if !strategy.EnableLogging {
		t.Error("Expected EnableLogging to be true")
	}

	if !strategy.UserFriendly {
		t.Error("Expected UserFriendly to be true")
	}
}

func TestDevelopmentErrorRecoveryStrategy(t *testing.T) {
	strategy := DevelopmentErrorRecoveryStrategy()

	if strategy.MaxRetries != 1 {
		t.Errorf("Expected MaxRetries to be 1, got %d", strategy.MaxRetries)
	}

	if strategy.ShowStackTrace != true {
		t.Error("Expected ShowStackTrace to be true in development")
	}

	if strategy.UserFriendly != false {
		t.Error("Expected UserFriendly to be false in development")
	}
}

func TestProductionErrorRecoveryStrategy(t *testing.T) {
	strategy := ProductionErrorRecoveryStrategy()

	if strategy.MaxRetries != 5 {
		t.Errorf("Expected MaxRetries to be 5, got %d", strategy.MaxRetries)
	}

	if strategy.ShowStackTrace != false {
		t.Error("Expected ShowStackTrace to be false in production")
	}

	if strategy.UserFriendly != true {
		t.Error("Expected UserFriendly to be true in production")
	}
}

func TestWidgetError_Error(t *testing.T) {
	originalErr := errors.New("test error")
	widgetError := &WidgetError{
		Err: originalErr,
		Context: ErrorContext{
			WidgetType: "TestWidget",
			WidgetID:   "test-123",
			Operation:  "render",
		},
	}

	errorStr := widgetError.Error()
	expectedSubstrings := []string{
		"TestWidget",
		"test-123",
		"render",
		"test error",
	}

	for _, substr := range expectedSubstrings {
		if !strings.Contains(errorStr, substr) {
			t.Errorf("Expected error string to contain '%s', got: %s", substr, errorStr)
		}
	}
}

func TestErrorWidget_Render_UserFriendly(t *testing.T) {
	originalErr := errors.New("nil pointer dereference")
	widgetError := &WidgetError{
		Err: originalErr,
		Context: ErrorContext{
			WidgetType: "TestWidget",
			WidgetID:   "test-123",
			Operation:  "render",
			Timestamp:  time.Now(),
		},
		Recoverable: true,
	}

	strategy := DefaultErrorRecoveryStrategy() // User-friendly by default
	errorWidget := ErrorWidget{
		Error:    widgetError,
		Strategy: strategy,
	}

	ctx := &core.Context{}
	html := errorWidget.Render(ctx)

	// Should contain user-friendly message
	if !strings.Contains(html, "Content not available") {
		t.Error("Expected user-friendly error message for nil error")
	}

	// Should contain error widget class
	if !strings.Contains(html, "godin-error-widget") {
		t.Error("Expected HTML to contain error widget class")
	}

	// Should contain accessibility attributes
	if !strings.Contains(html, "role=\"alert\"") {
		t.Error("Expected HTML to contain alert role")
	}

	// Should contain retry button for recoverable errors
	if !strings.Contains(html, "Retry") {
		t.Error("Expected HTML to contain retry button for recoverable error")
	}

	// Should contain dismiss button
	if !strings.Contains(html, "Dismiss") {
		t.Error("Expected HTML to contain dismiss button")
	}

	// Should NOT contain technical details in user-friendly mode
	if strings.Contains(html, "Technical Details") {
		t.Error("Expected HTML to not contain technical details in user-friendly mode")
	}
}

func TestErrorWidget_Render_Development(t *testing.T) {
	originalErr := errors.New("detailed technical error")
	widgetError := &WidgetError{
		Err: originalErr,
		Context: ErrorContext{
			WidgetType: "TestWidget",
			WidgetID:   "test-123",
			Operation:  "render",
			Timestamp:  time.Now(),
			StackTrace: "stack trace here",
		},
		Recoverable: true,
	}

	strategy := DevelopmentErrorRecoveryStrategy()
	errorWidget := ErrorWidget{
		Error:    widgetError,
		Strategy: strategy,
	}

	ctx := &core.Context{}
	html := errorWidget.Render(ctx)

	// Should contain technical error message
	if !strings.Contains(html, "Error in TestWidget") {
		t.Error("Expected technical error message in development mode")
	}

	// Should contain technical details
	if !strings.Contains(html, "Technical Details") {
		t.Error("Expected HTML to contain technical details in development mode")
	}

	// Should contain stack trace
	if !strings.Contains(html, "stack trace here") {
		t.Error("Expected HTML to contain stack trace in development mode")
	}

	// Should contain widget context information
	if !strings.Contains(html, "TestWidget") {
		t.Error("Expected HTML to contain widget type")
	}

	if !strings.Contains(html, "test-123") {
		t.Error("Expected HTML to contain widget ID")
	}
}

func TestErrorWidget_Render_MaxRetriesReached(t *testing.T) {
	originalErr := errors.New("persistent error")
	widgetError := &WidgetError{
		Err: originalErr,
		Context: ErrorContext{
			WidgetType: "TestWidget",
			WidgetID:   "test-123",
			Operation:  "render",
		},
		Recoverable: true,
		RetryCount:  5, // Exceeds max retries
	}

	strategy := DefaultErrorRecoveryStrategy() // MaxRetries = 3
	errorWidget := ErrorWidget{
		Error:    widgetError,
		Strategy: strategy,
	}

	ctx := &core.Context{}
	html := errorWidget.Render(ctx)

	// Should NOT contain retry button when max retries reached
	if strings.Contains(html, "Retry") {
		t.Error("Expected HTML to not contain retry button when max retries reached")
	}

	// Should still contain dismiss button
	if !strings.Contains(html, "Dismiss") {
		t.Error("Expected HTML to contain dismiss button")
	}
}

func TestErrorWidget_Render_NonRecoverable(t *testing.T) {
	originalErr := errors.New("fatal error")
	widgetError := &WidgetError{
		Err: originalErr,
		Context: ErrorContext{
			WidgetType: "TestWidget",
			WidgetID:   "test-123",
			Operation:  "render",
		},
		Recoverable: false,
	}

	strategy := DefaultErrorRecoveryStrategy()
	errorWidget := ErrorWidget{
		Error:    widgetError,
		Strategy: strategy,
	}

	ctx := &core.Context{}
	html := errorWidget.Render(ctx)

	// Should NOT contain retry button for non-recoverable errors
	if strings.Contains(html, "Retry") {
		t.Error("Expected HTML to not contain retry button for non-recoverable error")
	}

	// Should still contain dismiss button
	if !strings.Contains(html, "Dismiss") {
		t.Error("Expected HTML to contain dismiss button")
	}
}

func TestErrorWidget_Render_NilError(t *testing.T) {
	errorWidget := ErrorWidget{
		Error:    nil,
		Strategy: DefaultErrorRecoveryStrategy(),
	}

	ctx := &core.Context{}
	html := errorWidget.Render(ctx)

	// Should return empty string for nil error
	if html != "" {
		t.Errorf("Expected empty string for nil error, got: %s", html)
	}
}

// Mock widget for testing SafeRenderWidget
type MockWidget struct {
	ShouldPanic bool
	ShouldError bool
	Content     string
}

func (mw MockWidget) Render(ctx *core.Context) string {
	if mw.ShouldPanic {
		panic("mock panic")
	}
	if mw.ShouldError {
		// Simulate an error by returning empty content
		return ""
	}
	return mw.Content
}

func TestSafeRenderWidget_Success(t *testing.T) {
	widget := MockWidget{
		Content: "<div>Success</div>",
	}

	strategy := DefaultErrorRecoveryStrategy()
	ctx := &core.Context{}

	result := SafeRenderWidget(widget, ctx, strategy)

	if result != "<div>Success</div>" {
		t.Errorf("Expected successful render, got: %s", result)
	}
}

func TestSafeRenderWidget_Panic(t *testing.T) {
	widget := MockWidget{
		ShouldPanic: true,
	}

	strategy := DefaultErrorRecoveryStrategy()
	ctx := &core.Context{}

	result := SafeRenderWidget(widget, ctx, strategy)

	// Should return error widget HTML instead of panicking
	if !strings.Contains(result, "godin-error-widget") {
		t.Error("Expected error widget HTML when widget panics")
	}

	if !strings.Contains(result, "Something went wrong") {
		t.Error("Expected user-friendly error message when widget panics")
	}
}

func TestSafeRenderWidget_NilWidget(t *testing.T) {
	strategy := DefaultErrorRecoveryStrategy()
	ctx := &core.Context{}

	result := SafeRenderWidget(nil, ctx, strategy)

	// Should return empty string for nil widget
	if result != "" {
		t.Errorf("Expected empty string for nil widget, got: %s", result)
	}
}

func TestErrorBoundary_Render_Success(t *testing.T) {
	child := MockWidget{
		Content: "<div>Child content</div>",
	}

	boundary := ErrorBoundary{
		ID:    "test-boundary",
		Child: child,
	}

	ctx := &core.Context{}
	result := boundary.Render(ctx)

	if result != "<div>Child content</div>" {
		t.Errorf("Expected child content, got: %s", result)
	}
}

func TestErrorBoundary_Render_ChildPanic(t *testing.T) {
	child := MockWidget{
		ShouldPanic: true,
	}

	boundary := ErrorBoundary{
		ID:    "test-boundary",
		Child: child,
	}

	ctx := &core.Context{}
	result := boundary.Render(ctx)

	// Should catch child panic and render error widget
	if !strings.Contains(result, "godin-error-widget") {
		t.Error("Expected error boundary to catch child panic and render error widget")
	}
}

func TestErrorBoundary_Render_NilChild(t *testing.T) {
	boundary := ErrorBoundary{
		ID:    "test-boundary",
		Child: nil,
	}

	ctx := &core.Context{}
	result := boundary.Render(ctx)

	// Should return empty string for nil child
	if result != "" {
		t.Errorf("Expected empty string for nil child, got: %s", result)
	}
}

func TestErrorLogger_LogError(t *testing.T) {
	var loggedError error
	strategy := ErrorRecoveryStrategy{
		EnableLogging: true,
		ErrorCallback: func(err error) {
			loggedError = err
		},
	}

	logger := NewErrorLogger(strategy)
	widgetError := &WidgetError{
		Err: errors.New("test error"),
		Context: ErrorContext{
			WidgetType: "TestWidget",
		},
	}

	logger.LogError(widgetError)

	if loggedError != widgetError {
		t.Error("Expected error callback to be called with the widget error")
	}
}

func TestErrorLogger_LogError_Disabled(t *testing.T) {
	var callbackCalled bool
	strategy := ErrorRecoveryStrategy{
		EnableLogging: false,
		ErrorCallback: func(err error) {
			callbackCalled = true
		},
	}

	logger := NewErrorLogger(strategy)
	widgetError := &WidgetError{
		Err: errors.New("test error"),
	}

	logger.LogError(widgetError)

	if callbackCalled {
		t.Error("Expected error callback to not be called when logging is disabled")
	}
}
