package widgets

import (
	"fmt"
	"html"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/gideonsigilai/godin/pkg/core"
)

// ErrorRecoveryStrategy defines how widgets should handle errors
type ErrorRecoveryStrategy struct {
	MaxRetries     int
	RetryInterval  time.Duration
	FallbackValue  interface{}
	ErrorCallback  func(error)
	EnableLogging  bool
	ShowStackTrace bool
	UserFriendly   bool
}

// DefaultErrorRecoveryStrategy returns a default error recovery strategy
func DefaultErrorRecoveryStrategy() ErrorRecoveryStrategy {
	return ErrorRecoveryStrategy{
		MaxRetries:     3,
		RetryInterval:  time.Second * 2,
		EnableLogging:  true,
		ShowStackTrace: false,
		UserFriendly:   true,
	}
}

// DevelopmentErrorRecoveryStrategy returns an error recovery strategy for development
func DevelopmentErrorRecoveryStrategy() ErrorRecoveryStrategy {
	return ErrorRecoveryStrategy{
		MaxRetries:     1,
		RetryInterval:  time.Second,
		EnableLogging:  true,
		ShowStackTrace: true,
		UserFriendly:   false,
	}
}

// ProductionErrorRecoveryStrategy returns an error recovery strategy for production
func ProductionErrorRecoveryStrategy() ErrorRecoveryStrategy {
	return ErrorRecoveryStrategy{
		MaxRetries:     5,
		RetryInterval:  time.Second * 5,
		EnableLogging:  true,
		ShowStackTrace: false,
		UserFriendly:   true,
	}
}

// ErrorContext provides context information about where an error occurred
type ErrorContext struct {
	WidgetType string
	WidgetID   string
	Operation  string
	Timestamp  time.Time
	StackTrace string
	RequestID  string
	UserAgent  string
	ClientIP   string
}

// WidgetError represents an error that occurred in a widget
type WidgetError struct {
	Err           error
	Context       ErrorContext
	Recoverable   bool
	RetryCount    int
	LastRetryTime time.Time
}

// Error implements the error interface
func (we *WidgetError) Error() string {
	return fmt.Sprintf("[%s:%s] %s: %v", we.Context.WidgetType, we.Context.WidgetID, we.Context.Operation, we.Err)
}

// ErrorLogger handles logging of widget errors
type ErrorLogger struct {
	strategy ErrorRecoveryStrategy
}

// NewErrorLogger creates a new error logger
func NewErrorLogger(strategy ErrorRecoveryStrategy) *ErrorLogger {
	return &ErrorLogger{
		strategy: strategy,
	}
}

// LogError logs a widget error
func (el *ErrorLogger) LogError(widgetError *WidgetError) {
	if !el.strategy.EnableLogging {
		return
	}

	logMessage := fmt.Sprintf("[WIDGET ERROR] %s", widgetError.Error())

	if el.strategy.ShowStackTrace && widgetError.Context.StackTrace != "" {
		logMessage += fmt.Sprintf("\nStack Trace:\n%s", widgetError.Context.StackTrace)
	}

	log.Printf("%s", logMessage)

	// Call custom error callback if provided
	if el.strategy.ErrorCallback != nil {
		el.strategy.ErrorCallback(widgetError)
	}
}

// ErrorWidget represents a widget that displays error information
type ErrorWidget struct {
	ID       string
	Style    string
	Class    string
	Error    *WidgetError
	Strategy ErrorRecoveryStrategy
}

// Render renders the error widget as HTML
func (ew ErrorWidget) Render(ctx *core.Context) string {
	if ew.Error == nil {
		return ""
	}

	// Generate error ID for tracking
	errorID := fmt.Sprintf("error_%d", time.Now().UnixNano())

	// Build CSS classes
	className := "godin-error-widget"
	if ew.Class != "" {
		className += " " + ew.Class
	}

	// Determine error message based on strategy
	var errorMessage string
	var errorDetails string

	if ew.Strategy.UserFriendly {
		errorMessage = "Something went wrong"
		errorDetails = "We're sorry, but something unexpected happened. Please try again."

		// Provide more specific user-friendly messages for common errors
		errorStr := ew.Error.Error()
		if strings.Contains(errorStr, "nil") {
			errorMessage = "Content not available"
			errorDetails = "The requested content is currently not available. Please try refreshing the page."
		} else if strings.Contains(errorStr, "network") || strings.Contains(errorStr, "connection") {
			errorMessage = "Connection problem"
			errorDetails = "There seems to be a connection issue. Please check your internet connection and try again."
		} else if strings.Contains(errorStr, "timeout") {
			errorMessage = "Request timed out"
			errorDetails = "The request took too long to complete. Please try again."
		}
	} else {
		errorMessage = fmt.Sprintf("Error in %s", ew.Error.Context.WidgetType)
		errorDetails = ew.Error.Error()
	}

	// Build inline styles
	var styles []string
	if ew.Style != "" {
		styles = append(styles, ew.Style)
	}

	// Default error widget styles
	styles = append(styles, "color: #d32f2f")
	styles = append(styles, "padding: 16px")
	styles = append(styles, "border: 1px solid #f44336")
	styles = append(styles, "border-radius: 4px")
	styles = append(styles, "background-color: #ffebee")
	styles = append(styles, "font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif")
	styles = append(styles, "margin: 8px 0")

	styleAttr := strings.Join(styles, "; ")

	// Build error content
	var content strings.Builder

	// Error icon and title
	content.WriteString(`<div style="display: flex; align-items: center; margin-bottom: 8px;">`)
	content.WriteString(`<span style="font-size: 18px; margin-right: 8px;">⚠️</span>`)
	content.WriteString(fmt.Sprintf(`<strong>%s</strong>`, html.EscapeString(errorMessage)))
	content.WriteString(`</div>`)

	// Error details
	content.WriteString(fmt.Sprintf(`<div style="margin-bottom: 12px; font-size: 14px;">%s</div>`, html.EscapeString(errorDetails)))

	// Show technical details in development mode
	if ew.Strategy.ShowStackTrace {
		content.WriteString(`<details style="margin-top: 12px;">`)
		content.WriteString(`<summary style="cursor: pointer; font-weight: bold; margin-bottom: 8px;">Technical Details</summary>`)
		content.WriteString(`<div style="background-color: #f5f5f5; padding: 8px; border-radius: 4px; font-family: monospace; font-size: 12px; white-space: pre-wrap;">`)

		// Widget context
		content.WriteString(fmt.Sprintf("Widget: %s (ID: %s)\n", ew.Error.Context.WidgetType, ew.Error.Context.WidgetID))
		content.WriteString(fmt.Sprintf("Operation: %s\n", ew.Error.Context.Operation))
		content.WriteString(fmt.Sprintf("Timestamp: %s\n", ew.Error.Context.Timestamp.Format(time.RFC3339)))
		content.WriteString(fmt.Sprintf("Error: %s\n", html.EscapeString(ew.Error.Err.Error())))

		if ew.Error.Context.StackTrace != "" {
			content.WriteString(fmt.Sprintf("\nStack Trace:\n%s", html.EscapeString(ew.Error.Context.StackTrace)))
		}

		content.WriteString(`</div>`)
		content.WriteString(`</details>`)
	}

	// Action buttons
	content.WriteString(`<div style="margin-top: 12px;">`)

	// Retry button if error is recoverable
	if ew.Error.Recoverable && ew.Error.RetryCount < ew.Strategy.MaxRetries {
		content.WriteString(`<button onclick="retryWidget(this)" style="background-color: #2196F3; color: white; border: none; padding: 8px 16px; border-radius: 4px; cursor: pointer; margin-right: 8px;">`)
		content.WriteString(`Retry`)
		content.WriteString(`</button>`)
	}

	// Dismiss button
	content.WriteString(`<button onclick="dismissError(this)" style="background-color: #757575; color: white; border: none; padding: 8px 16px; border-radius: 4px; cursor: pointer;">`)
	content.WriteString(`Dismiss`)
	content.WriteString(`</button>`)

	content.WriteString(`</div>`)

	// Build attributes
	attrs := map[string]string{
		"id":        errorID,
		"class":     className,
		"style":     styleAttr,
		"role":      "alert",
		"aria-live": "polite",
	}

	if ew.ID != "" {
		attrs["id"] = ew.ID
	}

	// Add data attributes for JavaScript handling
	attrs["data-widget-type"] = ew.Error.Context.WidgetType
	attrs["data-widget-id"] = ew.Error.Context.WidgetID
	attrs["data-retry-count"] = fmt.Sprintf("%d", ew.Error.RetryCount)
	attrs["data-max-retries"] = fmt.Sprintf("%d", ew.Strategy.MaxRetries)

	// Build the final HTML
	var attrStrings []string
	for key, value := range attrs {
		attrStrings = append(attrStrings, fmt.Sprintf(`%s="%s"`, key, html.EscapeString(value)))
	}

	return fmt.Sprintf(`<div %s>%s</div>`, strings.Join(attrStrings, " "), content.String())
}

// SafeRenderWidget safely renders a widget with error recovery
func SafeRenderWidget(widget Widget, ctx *core.Context, strategy ErrorRecoveryStrategy) string {
	if widget == nil {
		return ""
	}

	// Create error context
	errorContext := ErrorContext{
		WidgetType: fmt.Sprintf("%T", widget),
		Operation:  "render",
		Timestamp:  time.Now(),
	}

	// Extract widget ID if available
	if idWidget, ok := widget.(interface{ GetID() string }); ok {
		errorContext.WidgetID = idWidget.GetID()
	}

	// Extract request context information
	if ctx != nil && ctx.Request != nil {
		errorContext.UserAgent = ctx.Request.UserAgent()
		errorContext.ClientIP = ctx.Request.RemoteAddr
		if requestID := ctx.Request.Header.Get("X-Request-ID"); requestID != "" {
			errorContext.RequestID = requestID
		}
	}

	// Attempt to render with error recovery
	var result string
	var lastError error

	for attempt := 0; attempt <= strategy.MaxRetries; attempt++ {
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Capture stack trace
					buf := make([]byte, 4096)
					n := runtime.Stack(buf, false)
					errorContext.StackTrace = string(buf[:n])

					// Convert panic to error
					if err, ok := r.(error); ok {
						lastError = err
					} else {
						lastError = fmt.Errorf("panic: %v", r)
					}
				}
			}()

			// Attempt to render
			result = widget.Render(ctx)
			lastError = nil // Success
		}()

		// If successful, return result
		if lastError == nil {
			return result
		}

		// If this was the last attempt, break
		if attempt >= strategy.MaxRetries {
			break
		}

		// Wait before retrying
		if strategy.RetryInterval > 0 {
			time.Sleep(strategy.RetryInterval)
		}
	}

	// Create widget error
	widgetError := &WidgetError{
		Err:           lastError,
		Context:       errorContext,
		Recoverable:   true,
		RetryCount:    strategy.MaxRetries,
		LastRetryTime: time.Now(),
	}

	// Log the error
	logger := NewErrorLogger(strategy)
	logger.LogError(widgetError)

	// Render error widget
	errorWidget := ErrorWidget{
		Error:    widgetError,
		Strategy: strategy,
	}

	return errorWidget.Render(ctx)
}

// ErrorBoundary represents a widget that catches and handles errors from its children
type ErrorBoundary struct {
	ID             string
	Style          string
	Class          string
	Child          Widget
	ErrorBuilder   func(*WidgetError) Widget
	OnError        func(*WidgetError)
	Strategy       ErrorRecoveryStrategy
	FallbackWidget Widget
}

// Render renders the error boundary widget
func (eb ErrorBoundary) Render(ctx *core.Context) string {
	if eb.Child == nil {
		return ""
	}

	// Use custom strategy or default
	strategy := eb.Strategy
	if strategy.MaxRetries == 0 && strategy.RetryInterval == 0 {
		strategy = DefaultErrorRecoveryStrategy()
	}

	// Attempt to render child with error recovery
	result := SafeRenderWidget(eb.Child, ctx, strategy)

	// If result contains error widget HTML, we can detect it and potentially use custom error builder
	if strings.Contains(result, "godin-error-widget") && eb.ErrorBuilder != nil {
		// In a more sophisticated implementation, we would parse the error and use the custom builder
		// For now, we'll use the result as-is
	}

	return result
}

// JavaScript functions for client-side error handling
const ErrorHandlingJavaScript = `
<script>
function retryWidget(button) {
    const errorWidget = button.closest('.godin-error-widget');
    const widgetType = errorWidget.getAttribute('data-widget-type');
    const widgetId = errorWidget.getAttribute('data-widget-id');
    const retryCount = parseInt(errorWidget.getAttribute('data-retry-count')) || 0;
    const maxRetries = parseInt(errorWidget.getAttribute('data-max-retries')) || 3;
    
    if (retryCount >= maxRetries) {
        alert('Maximum retry attempts reached');
        return;
    }
    
    // Show loading state
    button.disabled = true;
    button.textContent = 'Retrying...';
    
    // Attempt to reload the widget (this would need to be implemented based on your architecture)
    setTimeout(() => {
        // In a real implementation, this would trigger a re-render of the specific widget
        window.location.reload();
    }, 1000);
}

function dismissError(button) {
    const errorWidget = button.closest('.godin-error-widget');
    errorWidget.style.display = 'none';
}

// Global error handler for unhandled JavaScript errors
window.addEventListener('error', function(event) {
    console.error('Unhandled error:', event.error);
    
    // You could send this error to your server for logging
    if (window.godin && window.godin.reportError) {
        window.godin.reportError({
            message: event.message,
            filename: event.filename,
            lineno: event.lineno,
            colno: event.colno,
            error: event.error ? event.error.stack : null
        });
    }
});

// Global handler for unhandled promise rejections
window.addEventListener('unhandledrejection', function(event) {
    console.error('Unhandled promise rejection:', event.reason);
    
    if (window.godin && window.godin.reportError) {
        window.godin.reportError({
            type: 'unhandledrejection',
            reason: event.reason
        });
    }
});
</script>
`
