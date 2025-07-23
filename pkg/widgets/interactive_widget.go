package widgets

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/gideonsigilai/godin/pkg/core"
)

// InteractiveWidget provides base functionality for widgets with callbacks
type InteractiveWidget struct {
	HTMXWidget
	callbackRegistry    *core.CallbackRegistry
	widgetID            string
	widgetType          string
	registeredCallbacks map[string]string // callback type -> callback ID
	mutex               sync.RWMutex
	context             *core.Context
	isInitialized       bool
}

// NewInteractiveWidget creates a new InteractiveWidget instance
func NewInteractiveWidget(widgetType string, ctx *core.Context) *InteractiveWidget {
	widget := &InteractiveWidget{
		widgetType:          widgetType,
		widgetID:            generateWidgetID(),
		registeredCallbacks: make(map[string]string),
		context:             ctx,
	}

	if ctx != nil && ctx.App != nil {
		widget.callbackRegistry = ctx.App.CallbackRegistry()
	}

	return widget
}

// Initialize initializes the interactive widget with context
func (iw *InteractiveWidget) Initialize(ctx *core.Context) {
	iw.mutex.Lock()
	defer iw.mutex.Unlock()

	if iw.isInitialized {
		return
	}

	iw.context = ctx
	if ctx != nil && ctx.App != nil {
		iw.callbackRegistry = ctx.App.CallbackRegistry()
	}

	if iw.widgetID == "" {
		iw.widgetID = generateWidgetID()
	}

	iw.isInitialized = true
}

// RegisterCallback registers a callback function and returns the callback ID
func (iw *InteractiveWidget) RegisterCallback(callbackType string, fn interface{}) string {
	if fn == nil {
		return ""
	}

	iw.mutex.Lock()
	defer iw.mutex.Unlock()

	// Ensure widget is initialized
	if !iw.isInitialized && iw.context != nil {
		iw.Initialize(iw.context)
	}

	if iw.callbackRegistry == nil {
		return ""
	}

	// Register the callback
	callbackID := iw.callbackRegistry.RegisterCallback(
		iw.widgetID,
		iw.widgetType,
		callbackType,
		fn,
		iw.context,
	)

	// Store the callback ID for cleanup
	if callbackID != "" {
		// Initialize the map if it's nil
		if iw.registeredCallbacks == nil {
			iw.registeredCallbacks = make(map[string]string)
		}
		iw.registeredCallbacks[callbackType] = callbackID
	}

	return callbackID
}

// GetCallbackID returns the callback ID for a specific callback type
func (iw *InteractiveWidget) GetCallbackID(callbackType string) string {
	iw.mutex.RLock()
	defer iw.mutex.RUnlock()
	return iw.registeredCallbacks[callbackType]
}

// HasCallback returns true if the widget has a callback of the specified type
func (iw *InteractiveWidget) HasCallback(callbackType string) bool {
	iw.mutex.RLock()
	defer iw.mutex.RUnlock()
	_, exists := iw.registeredCallbacks[callbackType]
	return exists
}

// GenerateHTMXAttributes generates HTMX attributes for interactive elements
func (iw *InteractiveWidget) GenerateHTMXAttributes() map[string]string {
	iw.mutex.RLock()
	defer iw.mutex.RUnlock()

	attrs := make(map[string]string)

	// Add widget identification attributes
	attrs["data-widget-id"] = iw.widgetID
	attrs["data-widget-type"] = iw.widgetType

	// Generate HTMX attributes for each registered callback
	for callbackType, callbackID := range iw.registeredCallbacks {
		htmxAttrs := iw.generateHTMXForCallback(callbackType, callbackID)
		for key, value := range htmxAttrs {
			attrs[key] = value
		}
	}

	return attrs
}

// generateHTMXForCallback generates HTMX attributes for a specific callback
func (iw *InteractiveWidget) generateHTMXForCallback(callbackType, callbackID string) map[string]string {
	attrs := make(map[string]string)
	endpointPath := fmt.Sprintf("/api/callbacks/%s", callbackID)

	switch callbackType {
	case "OnPressed", "OnTap":
		attrs["hx-post"] = endpointPath
		attrs["hx-trigger"] = "click"
		attrs["hx-swap"] = "none" // Don't swap content by default

	case "OnChanged":
		attrs["hx-post"] = endpointPath
		attrs["hx-trigger"] = "change"
		attrs["hx-include"] = "this"
		attrs["hx-swap"] = "none"

	case "OnSubmitted", "OnFieldSubmitted":
		attrs["hx-post"] = endpointPath
		attrs["hx-trigger"] = "keyup[keyCode==13]" // Enter key
		attrs["hx-include"] = "this"
		attrs["hx-swap"] = "none"

	case "OnEditingComplete":
		attrs["hx-post"] = endpointPath
		attrs["hx-trigger"] = "blur"
		attrs["hx-include"] = "this"
		attrs["hx-swap"] = "none"

	case "OnDoubleTap":
		attrs["hx-post"] = endpointPath
		attrs["hx-trigger"] = "dblclick"
		attrs["hx-swap"] = "none"

	case "OnLongPress":
		attrs["hx-post"] = endpointPath
		attrs["hx-trigger"] = "contextmenu"
		attrs["hx-swap"] = "none"

	case "OnHover":
		attrs["hx-post"] = endpointPath
		attrs["hx-trigger"] = "mouseenter, mouseleave"
		attrs["hx-swap"] = "none"

	case "OnFocus":
		attrs["hx-post"] = endpointPath
		attrs["hx-trigger"] = "focus"
		attrs["hx-swap"] = "none"

	case "OnBlur":
		attrs["hx-post"] = endpointPath
		attrs["hx-trigger"] = "blur"
		attrs["hx-swap"] = "none"

	default:
		// Generic callback - use click as default trigger
		attrs["hx-post"] = endpointPath
		attrs["hx-trigger"] = "click"
		attrs["hx-swap"] = "none"
	}

	return attrs
}

// BuildEventHandlers builds JavaScript event handlers for fallback scenarios
func (iw *InteractiveWidget) BuildEventHandlers() map[string]string {
	iw.mutex.RLock()
	defer iw.mutex.RUnlock()

	handlers := make(map[string]string)

	for callbackType, callbackID := range iw.registeredCallbacks {
		handler := iw.generateEventHandler(callbackType, callbackID)
		if handler != "" {
			eventName := iw.getEventName(callbackType)
			if eventName != "" {
				handlers[eventName] = handler
			}
		}
	}

	return handlers
}

// generateEventHandler generates a JavaScript event handler for a callback
func (iw *InteractiveWidget) generateEventHandler(callbackType, callbackID string) string {
	endpointPath := fmt.Sprintf("/api/callbacks/%s", callbackID)

	switch callbackType {
	case "OnPressed", "OnTap":
		return fmt.Sprintf("handleWidgetCallback('%s', event)", endpointPath)

	case "OnChanged":
		return fmt.Sprintf("handleWidgetCallback('%s', event, this.value)", endpointPath)

	case "OnSubmitted", "OnFieldSubmitted":
		return fmt.Sprintf("if(event.key === 'Enter') handleWidgetCallback('%s', event, this.value)", endpointPath)

	case "OnEditingComplete":
		return fmt.Sprintf("handleWidgetCallback('%s', event, this.value)", endpointPath)

	case "OnDoubleTap":
		return fmt.Sprintf("handleWidgetCallback('%s', event)", endpointPath)

	case "OnLongPress":
		return fmt.Sprintf("handleWidgetCallback('%s', event); return false;", endpointPath)

	default:
		return fmt.Sprintf("handleWidgetCallback('%s', event)", endpointPath)
	}
}

// getEventName returns the JavaScript event name for a callback type
func (iw *InteractiveWidget) getEventName(callbackType string) string {
	switch callbackType {
	case "OnPressed", "OnTap":
		return "onclick"
	case "OnChanged":
		return "onchange"
	case "OnSubmitted", "OnFieldSubmitted":
		return "onkeypress"
	case "OnEditingComplete":
		return "onblur"
	case "OnDoubleTap":
		return "ondblclick"
	case "OnLongPress":
		return "oncontextmenu"
	case "OnHover":
		return "onmouseenter"
	case "OnFocus":
		return "onfocus"
	case "OnBlur":
		return "onblur"
	default:
		return "onclick"
	}
}

// GetWidgetID returns the unique widget identifier
func (iw *InteractiveWidget) GetWidgetID() string {
	iw.mutex.RLock()
	defer iw.mutex.RUnlock()
	return iw.widgetID
}

// GetWidgetType returns the widget type
func (iw *InteractiveWidget) GetWidgetType() string {
	return iw.widgetType
}

// SetWidgetType sets the widget type
func (iw *InteractiveWidget) SetWidgetType(widgetType string) {
	iw.mutex.Lock()
	defer iw.mutex.Unlock()
	iw.widgetType = widgetType
}

// GetRegisteredCallbacks returns a copy of registered callbacks
func (iw *InteractiveWidget) GetRegisteredCallbacks() map[string]string {
	iw.mutex.RLock()
	defer iw.mutex.RUnlock()

	result := make(map[string]string)
	for k, v := range iw.registeredCallbacks {
		result[k] = v
	}
	return result
}

// Cleanup removes all registered callbacks and cleans up resources
func (iw *InteractiveWidget) Cleanup() {
	iw.mutex.Lock()
	defer iw.mutex.Unlock()

	if iw.callbackRegistry != nil {
		// Clean up all registered callbacks
		for _, callbackID := range iw.registeredCallbacks {
			iw.callbackRegistry.CleanupCallback(callbackID)
		}
	}

	// Clear the callback map
	iw.registeredCallbacks = make(map[string]string)
	iw.isInitialized = false
}

// MergeAttributes merges HTMX attributes with existing attributes
func (iw *InteractiveWidget) MergeAttributes(existing map[string]string) map[string]string {
	htmxAttrs := iw.GenerateHTMXAttributes()
	eventHandlers := iw.BuildEventHandlers()

	// Start with existing attributes
	result := make(map[string]string)
	for k, v := range existing {
		result[k] = v
	}

	// Add HTMX attributes
	for k, v := range htmxAttrs {
		result[k] = v
	}

	// Add event handlers (only if not already present)
	for k, v := range eventHandlers {
		if _, exists := result[k]; !exists {
			result[k] = v
		}
	}

	return result
}

// IsInitialized returns true if the widget has been initialized
func (iw *InteractiveWidget) IsInitialized() bool {
	iw.mutex.RLock()
	defer iw.mutex.RUnlock()
	return iw.isInitialized
}

// SetContext sets the context for the widget
func (iw *InteractiveWidget) SetContext(ctx *core.Context) {
	iw.mutex.Lock()
	defer iw.mutex.Unlock()
	iw.context = ctx

	if ctx != nil && ctx.App != nil {
		iw.callbackRegistry = ctx.App.CallbackRegistry()
	}
}

// GetContext returns the current context
func (iw *InteractiveWidget) GetContext() *core.Context {
	iw.mutex.RLock()
	defer iw.mutex.RUnlock()
	return iw.context
}

// generateWidgetID generates a unique widget identifier
func generateWidgetID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return "widget_" + hex.EncodeToString(bytes)
}

// Helper function to generate JavaScript for widget callback handling
func GenerateCallbackHandlerScript() string {
	return `
<script>
function handleWidgetCallback(endpoint, event, value) {
	// Prevent default behavior
	if (event) {
		event.preventDefault();
	}

	// Prepare form data
	const formData = new FormData();
	if (value !== undefined) {
		formData.append('value', value);
	}

	// Add event information
	if (event) {
		formData.append('eventType', event.type);
		if (event.target) {
			formData.append('targetId', event.target.id || '');
			formData.append('targetValue', event.target.value || '');
		}
	}

	// Send request
	fetch(endpoint, {
		method: 'POST',
		body: formData
	})
	.then(response => {
		if (!response.ok) {
			console.error('Callback request failed:', response.statusText);
		}
		return response.json();
	})
	.then(data => {
		// Handle response if needed
		if (data && data.status === 'success') {
			// Callback executed successfully
			console.log('Callback executed successfully');
		}
	})
	.catch(error => {
		console.error('Callback error:', error);
	});
}
</script>`
}
