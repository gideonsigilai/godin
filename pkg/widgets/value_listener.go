package widgets

import (
	"fmt"
	"html"
	"sync"
	"time"

	"github.com/gideonsigilai/godin/pkg/core"
	"github.com/gideonsigilai/godin/pkg/renderer"
	"github.com/gideonsigilai/godin/pkg/state"
)

// ValueListener is a generic widget that listens to ValueNotifier changes and rebuilds automatically
type ValueListener[T any] struct {
	ID             string
	Style          string
	Class          string
	ValueNotifier  *state.ValueNotifier[T]
	Builder        func(value T) Widget
	OnValueChanged func(oldValue, newValue T)
	ErrorBuilder   func(error) Widget

	// Internal state for lifecycle management
	listenerID     string
	isRegistered   bool
	mutex          sync.RWMutex
	lastValue      T
	lastRenderTime time.Time
}

// NewValueListener creates a new ValueListener widget
func NewValueListener[T any](valueNotifier *state.ValueNotifier[T], builder func(value T) Widget) *ValueListener[T] {
	return &ValueListener[T]{
		ValueNotifier: valueNotifier,
		Builder:       builder,
		listenerID:    generateListenerID(),
	}
}

// NewValueListenerWithID creates a new ValueListener widget with a specific ID
func NewValueListenerWithID[T any](id string, valueNotifier *state.ValueNotifier[T], builder func(value T) Widget) *ValueListener[T] {
	return &ValueListener[T]{
		ID:            id,
		ValueNotifier: valueNotifier,
		Builder:       builder,
		listenerID:    generateListenerID(),
	}
}

// NewValueListenerWithOptions creates a new ValueListener widget with all options
func NewValueListenerWithOptions[T any](options ValueListenerOptions[T]) *ValueListener[T] {
	return &ValueListener[T]{
		ID:             options.ID,
		Style:          options.Style,
		Class:          options.Class,
		ValueNotifier:  options.ValueNotifier,
		Builder:        options.Builder,
		OnValueChanged: options.OnValueChanged,
		ErrorBuilder:   options.ErrorBuilder,
		listenerID:     generateListenerID(),
	}
}

// ValueListenerOptions contains all configuration options for ValueListener
type ValueListenerOptions[T any] struct {
	ID             string
	Style          string
	Class          string
	ValueNotifier  *state.ValueNotifier[T]
	Builder        func(value T) Widget
	OnValueChanged func(oldValue, newValue T)
	ErrorBuilder   func(error) Widget
}

// Type-specific ValueListener implementations for common types

// ValueListenerInt is a ValueListener for int values
type ValueListenerInt struct {
	ID             string
	Style          string
	Class          string
	ValueNotifier  *state.IntNotifier
	Builder        func(value int) Widget
	OnValueChanged func(oldValue, newValue int)
	ErrorBuilder   func(error) Widget

	// Internal state
	listenerID     string
	isRegistered   bool
	mutex          sync.RWMutex
	lastValue      int
	lastRenderTime time.Time
}

// NewValueListenerInt creates a new ValueListener for int values
func NewValueListenerInt(valueNotifier *state.IntNotifier, builder func(value int) Widget) *ValueListenerInt {
	return &ValueListenerInt{
		ValueNotifier: valueNotifier,
		Builder:       builder,
		listenerID:    generateListenerID(),
	}
}

// NewValueListenerIntWithID creates a new ValueListener for int values with a specific ID
func NewValueListenerIntWithID(id string, valueNotifier *state.IntNotifier, builder func(value int) Widget) *ValueListenerInt {
	return &ValueListenerInt{
		ID:            id,
		ValueNotifier: valueNotifier,
		Builder:       builder,
		listenerID:    generateListenerID(),
	}
}

// ValueListenerString is a ValueListener for string values
type ValueListenerString struct {
	ID             string
	Style          string
	Class          string
	ValueNotifier  *state.StringNotifier
	Builder        func(value string) Widget
	OnValueChanged func(oldValue, newValue string)
	ErrorBuilder   func(error) Widget

	// Internal state
	listenerID     string
	isRegistered   bool
	mutex          sync.RWMutex
	lastValue      string
	lastRenderTime time.Time
}

// NewValueListenerString creates a new ValueListener for string values
func NewValueListenerString(valueNotifier *state.StringNotifier, builder func(value string) Widget) *ValueListenerString {
	return &ValueListenerString{
		ValueNotifier: valueNotifier,
		Builder:       builder,
		listenerID:    generateListenerID(),
	}
}

// NewValueListenerStringWithID creates a new ValueListener for string values with a specific ID
func NewValueListenerStringWithID(id string, valueNotifier *state.StringNotifier, builder func(value string) Widget) *ValueListenerString {
	return &ValueListenerString{
		ID:            id,
		ValueNotifier: valueNotifier,
		Builder:       builder,
		listenerID:    generateListenerID(),
	}
}

// ValueListenerBool is a ValueListener for bool values
type ValueListenerBool struct {
	ID             string
	Style          string
	Class          string
	ValueNotifier  *state.BoolNotifier
	Builder        func(value bool) Widget
	OnValueChanged func(oldValue, newValue bool)
	ErrorBuilder   func(error) Widget

	// Internal state
	listenerID     string
	isRegistered   bool
	mutex          sync.RWMutex
	lastValue      bool
	lastRenderTime time.Time
}

// NewValueListenerBool creates a new ValueListener for bool values
func NewValueListenerBool(valueNotifier *state.BoolNotifier, builder func(value bool) Widget) *ValueListenerBool {
	return &ValueListenerBool{
		ValueNotifier: valueNotifier,
		Builder:       builder,
		listenerID:    generateListenerID(),
	}
}

// NewValueListenerBoolWithID creates a new ValueListener for bool values with a specific ID
func NewValueListenerBoolWithID(id string, valueNotifier *state.BoolNotifier, builder func(value bool) Widget) *ValueListenerBool {
	return &ValueListenerBool{
		ID:            id,
		ValueNotifier: valueNotifier,
		Builder:       builder,
		listenerID:    generateListenerID(),
	}
}

// ValueListenerFloat64 is a ValueListener for float64 values
type ValueListenerFloat64 struct {
	ID             string
	Style          string
	Class          string
	ValueNotifier  *state.Float64Notifier
	Builder        func(value float64) Widget
	OnValueChanged func(oldValue, newValue float64)
	ErrorBuilder   func(error) Widget

	// Internal state
	listenerID     string
	isRegistered   bool
	mutex          sync.RWMutex
	lastValue      float64
	lastRenderTime time.Time
}

// NewValueListenerFloat64 creates a new ValueListener for float64 values
func NewValueListenerFloat64(valueNotifier *state.Float64Notifier, builder func(value float64) Widget) *ValueListenerFloat64 {
	return &ValueListenerFloat64{
		ValueNotifier: valueNotifier,
		Builder:       builder,
		listenerID:    generateListenerID(),
	}
}

// NewValueListenerFloat64WithID creates a new ValueListener for float64 values with a specific ID
func NewValueListenerFloat64WithID(id string, valueNotifier *state.Float64Notifier, builder func(value float64) Widget) *ValueListenerFloat64 {
	return &ValueListenerFloat64{
		ID:            id,
		ValueNotifier: valueNotifier,
		Builder:       builder,
		listenerID:    generateListenerID(),
	}
}

// Helper function to generate unique listener IDs
func generateListenerID() string {
	return fmt.Sprintf("vl_%d", time.Now().UnixNano())
}

// EnhancedValueListenableBuilder is an improved version of ValueListenableBuilder with automatic updates
type EnhancedValueListenableBuilder struct {
	InteractiveWidget // Embed InteractiveWidget for callback support
	ID                string
	Style             string
	Class             string
	ValueListenable   interface{}              // Value notifier to listen to (can be any ValueNotifier type)
	Builder           func(interface{}) Widget // Builder function
	Child             Widget                   // Optional child widget
	UpdateMode        UpdateMode               // How updates should be delivered
	DebounceMs        int                      // Debounce time in milliseconds
}

// UpdateMode defines how updates are delivered to the client
type UpdateMode string

// Constants for UpdateMode
const (
	UpdateModeWebSocket UpdateMode = "websocket" // Use WebSocket for real-time updates
	UpdateModePolling   UpdateMode = "polling"   // Use polling for updates
	UpdateModeHTMX      UpdateMode = "htmx"      // Use HTMX for updates
)

// Render renders the EnhancedValueListenableBuilder as HTML
func (vb EnhancedValueListenableBuilder) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	// Initialize the InteractiveWidget if needed
	if !vb.InteractiveWidget.IsInitialized() {
		vb.InteractiveWidget.Initialize(ctx)
		vb.InteractiveWidget.SetWidgetType("EnhancedValueListenableBuilder")
	}

	// Build attributes
	attrs := buildAttributes(vb.ID, vb.Style, vb.Class+" godin-value-listener")

	// Add data attributes if we have a value listenable
	if vb.ValueListenable != nil {
		// Try to get ID and type from the value listenable
		if idGetter, ok := vb.ValueListenable.(interface{ GetID() string }); ok {
			attrs["data-value-notifier-id"] = idGetter.GetID()
		}
		if typeGetter, ok := vb.ValueListenable.(interface{ GetValueType() string }); ok {
			attrs["data-value-type"] = typeGetter.GetValueType()
		}
		attrs["data-update-mode"] = string(vb.UpdateMode)
		if vb.DebounceMs > 0 {
			attrs["data-debounce"] = fmt.Sprintf("%d", vb.DebounceMs)
		}
	}

	// Set up WebSocket connection for real-time updates if needed
	if vb.UpdateMode == UpdateModeWebSocket && ctx != nil && ctx.App != nil {
		// Register this widget for WebSocket updates
		if vb.ValueListenable != nil {
			if idGetter, ok := vb.ValueListenable.(interface{ GetID() string }); ok {
				notifierID := idGetter.GetID()
				widgetID := vb.ID
				if widgetID == "" {
					widgetID = fmt.Sprintf("vlb_%p", &vb)
					attrs["id"] = widgetID
				}

				// Add WebSocket attributes
				attrs["data-ws-listen"] = "true"
				attrs["data-ws-channel"] = "value_notifier_" + notifierID
			}
		}
	} else if vb.UpdateMode == UpdateModePolling {
		// Set up polling for updates
		if vb.ValueListenable != nil {
			if idGetter, ok := vb.ValueListenable.(interface{ GetID() string }); ok {
				notifierID := idGetter.GetID()
				pollingEndpoint := "/api/value_notifier/" + notifierID + "/poll"
				attrs["data-polling-endpoint"] = pollingEndpoint
				attrs["data-polling-interval"] = "1000" // Default 1 second
			}
		}
	} else if vb.UpdateMode == UpdateModeHTMX {
		// Set up HTMX for updates
		if vb.ValueListenable != nil {
			if idGetter, ok := vb.ValueListenable.(interface{ GetID() string }); ok {
				notifierID := idGetter.GetID()
				attrs["hx-get"] = "/api/value_notifier/" + notifierID + "/value"
				attrs["hx-trigger"] = "valueChanged from:body"
				attrs["hx-swap"] = "innerHTML"
			}
		}
	}

	// Merge with interactive widget attributes
	attrs = vb.InteractiveWidget.MergeAttributes(attrs)

	// Build the current content using the builder function
	content := ""
	if vb.Builder != nil && vb.ValueListenable != nil {
		// Try to get current value from the value listenable
		if valueGetter, ok := vb.ValueListenable.(interface{ GetValue() interface{} }); ok {
			currentValue := valueGetter.GetValue()
			builtWidget := vb.Builder(currentValue)
			if builtWidget != nil {
				content = builtWidget.Render(ctx)
			}
		} else if valueGetter, ok := vb.ValueListenable.(interface{ Value() interface{} }); ok {
			currentValue := valueGetter.Value()
			builtWidget := vb.Builder(currentValue)
			if builtWidget != nil {
				content = builtWidget.Render(ctx)
			}
		}
	} else if vb.Child != nil {
		content = vb.Child.Render(ctx)
	}

	return htmlRenderer.RenderElement("div", attrs, content, false)
}

// Render renders the ValueListener widget with state integration and WebSocket support
func (vl *ValueListener[T]) Render(ctx *core.Context) string {
	vl.mutex.Lock()
	defer vl.mutex.Unlock()

	// Register with StateManager if not already registered
	if !vl.isRegistered && vl.ValueNotifier != nil {
		vl.registerWithStateManager(ctx)
	}

	// Get current value from ValueNotifier
	if vl.ValueNotifier == nil {
		return vl.renderError(fmt.Errorf("ValueNotifier is nil"))
	}

	currentValue := vl.ValueNotifier.Value()
	vl.lastValue = currentValue
	vl.lastRenderTime = time.Now()

	// Build child widget using the builder function
	if vl.Builder == nil {
		return vl.renderError(fmt.Errorf("Builder function is nil"))
	}

	childWidget := vl.Builder(currentValue)
	if childWidget == nil {
		return vl.renderError(fmt.Errorf("Builder function returned nil widget"))
	}

	// Generate unique ID if not provided
	id := vl.ID
	if id == "" {
		id = fmt.Sprintf("vl_%s", vl.ValueNotifier.ID())
	}

	// Build CSS classes
	containerClass := "value-listener"
	if vl.Class != "" {
		containerClass += " " + vl.Class
	}

	// Build style
	style := vl.Style
	if style == "" {
		style = "display: contents;" // Don't add extra layout by default
	}

	// Render child widget content
	childContent := childWidget.Render(ctx)

	// Create container with WebSocket update capabilities and data attributes
	return fmt.Sprintf(`
		<div id="%s"
			 class="%s"
			 style="%s"
			 data-value-notifier-id="%s"
			 data-listener-id="%s"
			 data-current-value="%s"
			 data-last-updated="%d">
			%s
		</div>
		<script>
			(function() {
				const element = document.getElementById('%s');
				const notifierId = '%s';
				const listenerId = '%s';

				// Initialize Godin state manager if not already done
				if (!window.godin) {
					window.godin = {
						subscriptions: new Map(),
						websocket: null,
						subscribe: function(channel, callback) {
							if (!this.subscriptions.has(channel)) {
								this.subscriptions.set(channel, []);
							}
							this.subscriptions.get(channel).push(callback);
							
							// Initialize WebSocket connection if not already done
							this.initWebSocket();
						},
						initWebSocket: function() {
							if (this.websocket && this.websocket.readyState === WebSocket.OPEN) {
								return;
							}
							
							const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
							const wsUrl = protocol + '//' + window.location.host + '/ws';
							
							this.websocket = new WebSocket(wsUrl);
							
							this.websocket.onmessage = (event) => {
								try {
									const data = JSON.parse(event.data);
									const channel = 'state:' + data.id;
									const callbacks = this.subscriptions.get(channel) || [];
									callbacks.forEach(callback => callback(data));
								} catch (e) {
									console.error('Error processing WebSocket message:', e);
								}
							};
							
							this.websocket.onerror = (error) => {
								console.error('WebSocket error:', error);
								// Fallback to polling
								this.startPolling();
							};
							
							this.websocket.onclose = () => {
								console.log('WebSocket connection closed, attempting to reconnect...');
								setTimeout(() => this.initWebSocket(), 1000);
							};
						},
						startPolling: function() {
							// Fallback polling mechanism
							setInterval(() => {
								this.subscriptions.forEach((callbacks, channel) => {
									if (channel.startsWith('state:')) {
										const notifierId = channel.replace('state:', '');
										fetch('/api/state/' + notifierId)
											.then(response => response.json())
											.then(data => {
												callbacks.forEach(callback => callback(data));
											})
											.catch(console.error);
									}
								});
							}, 1000);
						}
					};
				}

				// Register for WebSocket updates
				window.godin.subscribe('state:' + notifierId, function(data) {
					// Update the element with new content if provided
					if (data && data.html) {
						element.innerHTML = data.html;
					}

					// Update data attributes
					if (data && data.value !== undefined) {
						element.setAttribute('data-current-value', JSON.stringify(data.value));
						element.setAttribute('data-last-updated', Date.now().toString());
					}

					// Trigger custom event for other components to listen to
					element.dispatchEvent(new CustomEvent('valueChanged', {
						detail: { 
							value: data.value, 
							notifierId: notifierId,
							listenerId: listenerId,
							timestamp: data.timestamp
						}
					}));
				});

				// Mark element as initialized
				element.setAttribute('data-value-listener-initialized', 'true');
			})();
		</script>`,
		html.EscapeString(id),
		html.EscapeString(containerClass),
		html.EscapeString(style),
		html.EscapeString(vl.ValueNotifier.ID()),
		html.EscapeString(vl.listenerID),
		html.EscapeString(jsonStringify(currentValue)),
		vl.lastRenderTime.Unix(),
		childContent,
		html.EscapeString(id),
		html.EscapeString(vl.ValueNotifier.ID()),
		html.EscapeString(vl.listenerID),
	)
}

// registerWithStateManager registers this ValueListener with the StateManager for lifecycle management
func (vl *ValueListener[T]) registerWithStateManager(ctx *core.Context) {
	if vl.ValueNotifier == nil {
		return
	}

	// Register the ValueNotifier with the StateManager
	stateManager := ctx.App.State()
	stateManager.RegisterValueNotifier(vl.ValueNotifier.ID(), vl.ValueNotifier)

	// Set the StateManager on the ValueNotifier for WebSocket broadcasting
	vl.ValueNotifier.SetManager(stateManager)

	// Add listener to the ValueNotifier for change notifications
	vl.ValueNotifier.AddListener(func(newValue T) {
		// Call OnValueChanged callback if provided
		if vl.OnValueChanged != nil {
			vl.OnValueChanged(vl.lastValue, newValue)
		}
		vl.lastValue = newValue
	})

	vl.isRegistered = true
}

// renderError renders an error state for the ValueListener
func (vl *ValueListener[T]) renderError(err error) string {
	// Use custom error builder if provided
	if vl.ErrorBuilder != nil {
		errorWidget := vl.ErrorBuilder(err)
		if errorWidget != nil {
			// We can't call Render here without context, so return a simple error message
			// In a real implementation, you'd need to pass context through or handle this differently
			return fmt.Sprintf(`<div class="value-listener-error">%s</div>`, html.EscapeString(err.Error()))
		}
	}

	// Default error rendering
	return fmt.Sprintf(`
		<div class="value-listener-error" style="color: red; padding: 8px; border: 1px solid red; border-radius: 4px; background-color: #ffebee;">
			<strong>ValueListener Error:</strong> %s
		</div>`,
		html.EscapeString(err.Error()),
	)
}

// Cleanup removes the listener from the ValueNotifier (should be called when widget is disposed)
func (vl *ValueListener[T]) Cleanup() {
	vl.mutex.Lock()
	defer vl.mutex.Unlock()

	if vl.ValueNotifier != nil && vl.isRegistered {
		// Note: ValueNotifier doesn't have a RemoveListener method that takes a specific listener
		// In a production implementation, you'd need to enhance ValueNotifier to support this
		// For now, we'll clear all listeners (not ideal but functional)
		vl.ValueNotifier.ClearListeners()
		vl.isRegistered = false
	}
}

// Type-specific Render implementations

// Render renders the ValueListenerInt widget
func (vl *ValueListenerInt) Render(ctx *core.Context) string {
	vl.mutex.Lock()
	defer vl.mutex.Unlock()

	// Register with StateManager if not already registered
	if !vl.isRegistered && vl.ValueNotifier != nil {
		vl.registerWithStateManagerInt(ctx)
	}

	// Get current value from ValueNotifier
	if vl.ValueNotifier == nil {
		return vl.renderErrorInt(fmt.Errorf("ValueNotifier is nil"))
	}

	currentValue := vl.ValueNotifier.Value()
	vl.lastValue = currentValue
	vl.lastRenderTime = time.Now()

	// Build child widget using the builder function
	if vl.Builder == nil {
		return vl.renderErrorInt(fmt.Errorf("Builder function is nil"))
	}

	childWidget := vl.Builder(currentValue)
	if childWidget == nil {
		return vl.renderErrorInt(fmt.Errorf("Builder function returned nil widget"))
	}

	// Generate unique ID if not provided
	id := vl.ID
	if id == "" {
		id = fmt.Sprintf("vl_int_%s", vl.ValueNotifier.ID())
	}

	// Build CSS classes
	containerClass := "value-listener value-listener-int"
	if vl.Class != "" {
		containerClass += " " + vl.Class
	}

	// Build style
	style := vl.Style
	if style == "" {
		style = "display: contents;"
	}

	// Render child widget content
	childContent := childWidget.Render(ctx)

	// Create container with WebSocket update capabilities
	return fmt.Sprintf(`
		<div id="%s"
			 class="%s"
			 style="%s"
			 data-value-notifier-id="%s"
			 data-listener-id="%s"
			 data-current-value="%d"
			 data-value-type="int"
			 data-last-updated="%d">
			%s
		</div>`,
		html.EscapeString(id),
		html.EscapeString(containerClass),
		html.EscapeString(style),
		html.EscapeString(vl.ValueNotifier.ID()),
		html.EscapeString(vl.listenerID),
		currentValue,
		vl.lastRenderTime.Unix(),
		childContent,
	)
}

func (vl *ValueListenerInt) registerWithStateManagerInt(ctx *core.Context) {
	if vl.ValueNotifier == nil {
		return
	}

	stateManager := ctx.App.State()
	stateManager.RegisterValueNotifier(vl.ValueNotifier.ID(), vl.ValueNotifier)
	vl.ValueNotifier.SetManager(stateManager)

	vl.ValueNotifier.AddListener(func(newValue int) {
		if vl.OnValueChanged != nil {
			vl.OnValueChanged(vl.lastValue, newValue)
		}
		vl.lastValue = newValue
	})

	vl.isRegistered = true
}

func (vl *ValueListenerInt) renderErrorInt(err error) string {
	if vl.ErrorBuilder != nil {
		errorWidget := vl.ErrorBuilder(err)
		if errorWidget != nil {
			return fmt.Sprintf(`<div class="value-listener-error">%s</div>`, html.EscapeString(err.Error()))
		}
	}

	return fmt.Sprintf(`
		<div class="value-listener-error" style="color: red; padding: 8px; border: 1px solid red; border-radius: 4px; background-color: #ffebee;">
			<strong>ValueListenerInt Error:</strong> %s
		</div>`,
		html.EscapeString(err.Error()),
	)
}

// Render renders the ValueListenerString widget
func (vl *ValueListenerString) Render(ctx *core.Context) string {
	vl.mutex.Lock()
	defer vl.mutex.Unlock()

	if !vl.isRegistered && vl.ValueNotifier != nil {
		vl.registerWithStateManagerString(ctx)
	}

	if vl.ValueNotifier == nil {
		return vl.renderErrorString(fmt.Errorf("ValueNotifier is nil"))
	}

	currentValue := vl.ValueNotifier.Value()
	vl.lastValue = currentValue
	vl.lastRenderTime = time.Now()

	if vl.Builder == nil {
		return vl.renderErrorString(fmt.Errorf("Builder function is nil"))
	}

	childWidget := vl.Builder(currentValue)
	if childWidget == nil {
		return vl.renderErrorString(fmt.Errorf("Builder function returned nil widget"))
	}

	id := vl.ID
	if id == "" {
		id = fmt.Sprintf("vl_string_%s", vl.ValueNotifier.ID())
	}

	containerClass := "value-listener value-listener-string"
	if vl.Class != "" {
		containerClass += " " + vl.Class
	}

	style := vl.Style
	if style == "" {
		style = "display: contents;"
	}

	childContent := childWidget.Render(ctx)

	return fmt.Sprintf(`
		<div id="%s"
			 class="%s"
			 style="%s"
			 data-value-notifier-id="%s"
			 data-listener-id="%s"
			 data-current-value="%s"
			 data-value-type="string"
			 data-last-updated="%d">
			%s
		</div>`,
		html.EscapeString(id),
		html.EscapeString(containerClass),
		html.EscapeString(style),
		html.EscapeString(vl.ValueNotifier.ID()),
		html.EscapeString(vl.listenerID),
		html.EscapeString(currentValue),
		vl.lastRenderTime.Unix(),
		childContent,
	)
}

func (vl *ValueListenerString) registerWithStateManagerString(ctx *core.Context) {
	if vl.ValueNotifier == nil {
		return
	}

	stateManager := ctx.App.State()
	stateManager.RegisterValueNotifier(vl.ValueNotifier.ID(), vl.ValueNotifier)
	vl.ValueNotifier.SetManager(stateManager)

	vl.ValueNotifier.AddListener(func(newValue string) {
		if vl.OnValueChanged != nil {
			vl.OnValueChanged(vl.lastValue, newValue)
		}
		vl.lastValue = newValue
	})

	vl.isRegistered = true
}

func (vl *ValueListenerString) renderErrorString(err error) string {
	if vl.ErrorBuilder != nil {
		errorWidget := vl.ErrorBuilder(err)
		if errorWidget != nil {
			return fmt.Sprintf(`<div class="value-listener-error">%s</div>`, html.EscapeString(err.Error()))
		}
	}

	return fmt.Sprintf(`
		<div class="value-listener-error" style="color: red; padding: 8px; border: 1px solid red; border-radius: 4px; background-color: #ffebee;">
			<strong>ValueListenerString Error:</strong> %s
		</div>`,
		html.EscapeString(err.Error()),
	)
}

// Render renders the ValueListenerBool widget
func (vl *ValueListenerBool) Render(ctx *core.Context) string {
	vl.mutex.Lock()
	defer vl.mutex.Unlock()

	if !vl.isRegistered && vl.ValueNotifier != nil {
		vl.registerWithStateManagerBool(ctx)
	}

	if vl.ValueNotifier == nil {
		return vl.renderErrorBool(fmt.Errorf("ValueNotifier is nil"))
	}

	currentValue := vl.ValueNotifier.Value()
	vl.lastValue = currentValue
	vl.lastRenderTime = time.Now()

	if vl.Builder == nil {
		return vl.renderErrorBool(fmt.Errorf("Builder function is nil"))
	}

	childWidget := vl.Builder(currentValue)
	if childWidget == nil {
		return vl.renderErrorBool(fmt.Errorf("Builder function returned nil widget"))
	}

	id := vl.ID
	if id == "" {
		id = fmt.Sprintf("vl_bool_%s", vl.ValueNotifier.ID())
	}

	containerClass := "value-listener value-listener-bool"
	if vl.Class != "" {
		containerClass += " " + vl.Class
	}

	style := vl.Style
	if style == "" {
		style = "display: contents;"
	}

	childContent := childWidget.Render(ctx)

	return fmt.Sprintf(`
		<div id="%s"
			 class="%s"
			 style="%s"
			 data-value-notifier-id="%s"
			 data-listener-id="%s"
			 data-current-value="%t"
			 data-value-type="bool"
			 data-last-updated="%d">
			%s
		</div>`,
		html.EscapeString(id),
		html.EscapeString(containerClass),
		html.EscapeString(style),
		html.EscapeString(vl.ValueNotifier.ID()),
		html.EscapeString(vl.listenerID),
		currentValue,
		vl.lastRenderTime.Unix(),
		childContent,
	)
}

func (vl *ValueListenerBool) registerWithStateManagerBool(ctx *core.Context) {
	if vl.ValueNotifier == nil {
		return
	}

	stateManager := ctx.App.State()
	stateManager.RegisterValueNotifier(vl.ValueNotifier.ID(), vl.ValueNotifier)
	vl.ValueNotifier.SetManager(stateManager)

	vl.ValueNotifier.AddListener(func(newValue bool) {
		if vl.OnValueChanged != nil {
			vl.OnValueChanged(vl.lastValue, newValue)
		}
		vl.lastValue = newValue
	})

	vl.isRegistered = true
}

func (vl *ValueListenerBool) renderErrorBool(err error) string {
	if vl.ErrorBuilder != nil {
		errorWidget := vl.ErrorBuilder(err)
		if errorWidget != nil {
			return fmt.Sprintf(`<div class="value-listener-error">%s</div>`, html.EscapeString(err.Error()))
		}
	}

	return fmt.Sprintf(`
		<div class="value-listener-error" style="color: red; padding: 8px; border: 1px solid red; border-radius: 4px; background-color: #ffebee;">
			<strong>ValueListenerBool Error:</strong> %s
		</div>`,
		html.EscapeString(err.Error()),
	)
}

// Render renders the ValueListenerFloat64 widget
func (vl *ValueListenerFloat64) Render(ctx *core.Context) string {
	vl.mutex.Lock()
	defer vl.mutex.Unlock()

	if !vl.isRegistered && vl.ValueNotifier != nil {
		vl.registerWithStateManagerFloat64(ctx)
	}

	if vl.ValueNotifier == nil {
		return vl.renderErrorFloat64(fmt.Errorf("ValueNotifier is nil"))
	}

	currentValue := vl.ValueNotifier.Value()
	vl.lastValue = currentValue
	vl.lastRenderTime = time.Now()

	if vl.Builder == nil {
		return vl.renderErrorFloat64(fmt.Errorf("Builder function is nil"))
	}

	childWidget := vl.Builder(currentValue)
	if childWidget == nil {
		return vl.renderErrorFloat64(fmt.Errorf("Builder function returned nil widget"))
	}

	id := vl.ID
	if id == "" {
		id = fmt.Sprintf("vl_float64_%s", vl.ValueNotifier.ID())
	}

	containerClass := "value-listener value-listener-float64"
	if vl.Class != "" {
		containerClass += " " + vl.Class
	}

	style := vl.Style
	if style == "" {
		style = "display: contents;"
	}

	childContent := childWidget.Render(ctx)

	return fmt.Sprintf(`
		<div id="%s"
			 class="%s"
			 style="%s"
			 data-value-notifier-id="%s"
			 data-listener-id="%s"
			 data-current-value="%.6f"
			 data-value-type="float64"
			 data-last-updated="%d">
			%s
		</div>`,
		html.EscapeString(id),
		html.EscapeString(containerClass),
		html.EscapeString(style),
		html.EscapeString(vl.ValueNotifier.ID()),
		html.EscapeString(vl.listenerID),
		currentValue,
		vl.lastRenderTime.Unix(),
		childContent,
	)
}

func (vl *ValueListenerFloat64) registerWithStateManagerFloat64(ctx *core.Context) {
	if vl.ValueNotifier == nil {
		return
	}

	stateManager := ctx.App.State()
	stateManager.RegisterValueNotifier(vl.ValueNotifier.ID(), vl.ValueNotifier)
	vl.ValueNotifier.SetManager(stateManager)

	vl.ValueNotifier.AddListener(func(newValue float64) {
		if vl.OnValueChanged != nil {
			vl.OnValueChanged(vl.lastValue, newValue)
		}
		vl.lastValue = newValue
	})

	vl.isRegistered = true
}

func (vl *ValueListenerFloat64) renderErrorFloat64(err error) string {
	if vl.ErrorBuilder != nil {
		errorWidget := vl.ErrorBuilder(err)
		if errorWidget != nil {
			return fmt.Sprintf(`<div class="value-listener-error">%s</div>`, html.EscapeString(err.Error()))
		}
	}

	return fmt.Sprintf(`
		<div class="value-listener-error" style="color: red; padding: 8px; border: 1px solid red; border-radius: 4px; background-color: #ffebee;">
			<strong>ValueListenerFloat64 Error:</strong> %s
		</div>`,
		html.EscapeString(err.Error()),
	)
}

// Enhanced error handling and graceful degradation

// SafeRender wraps the Render method with panic recovery and error handling
func (vl *ValueListener[T]) SafeRender(ctx *core.Context) (result string) {
	// Recover from panics in the render process
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic in ValueListener.Render: %v", r)
			result = vl.renderErrorWithRecovery(err, ctx)
		}
	}()

	return vl.Render(ctx)
}

// renderErrorWithRecovery renders an error with enhanced recovery options
func (vl *ValueListener[T]) renderErrorWithRecovery(err error, ctx *core.Context) string {
	// Log the error if logging is enabled
	if ctx != nil {
		// In a real implementation, you'd use a proper logger
		fmt.Printf("[ValueListener Error] %s: %v\n", time.Now().Format(time.RFC3339), err)
	}

	// Use custom error builder if provided
	if vl.ErrorBuilder != nil {
		// Safely call the error builder with panic recovery
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Error builder itself panicked, fall back to default error rendering
					fmt.Printf("[ValueListener Error Builder Panic] %v\n", r)
				}
			}()

			errorWidget := vl.ErrorBuilder(err)
			if errorWidget != nil && ctx != nil {
				result := errorWidget.Render(ctx)
				if result != "" {
					return
				}
			}
		}()
	}

	// Default error rendering with graceful degradation
	return vl.renderDefaultError(err)
}

// renderDefaultError renders a default error widget with fallback mechanisms
func (vl *ValueListener[T]) renderDefaultError(err error) string {
	// Generate a safe ID for the error widget
	errorID := fmt.Sprintf("vl_error_%d", time.Now().UnixNano())

	// Create a user-friendly error message
	userMessage := "Unable to display content"
	if err != nil {
		// Only show technical details in development mode
		// In production, you'd want to log this and show a generic message
		userMessage = fmt.Sprintf("Error: %s", err.Error())
	}

	return fmt.Sprintf(`
		<div id="%s" 
			 class="value-listener-error" 
			 style="color: #d32f2f; padding: 12px; border: 1px solid #f44336; border-radius: 4px; background-color: #ffebee; font-family: monospace; margin: 4px 0;">
			<div style="font-weight: bold; margin-bottom: 4px;">⚠️ ValueListener Error</div>
			<div style="font-size: 0.9em;">%s</div>
			<div style="margin-top: 8px; font-size: 0.8em; color: #666;">
				<button onclick="this.parentElement.parentElement.style.display='none'" 
						style="background: none; border: 1px solid #ccc; padding: 2px 8px; cursor: pointer; border-radius: 2px;">
					Dismiss
				</button>
			</div>
		</div>`,
		html.EscapeString(errorID),
		html.EscapeString(userMessage),
	)
}

// Enhanced SafeRender methods for type-specific implementations

// SafeRender wraps the Render method with panic recovery for ValueListenerInt
func (vl *ValueListenerInt) SafeRender(ctx *core.Context) (result string) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic in ValueListenerInt.Render: %v", r)
			result = vl.renderErrorWithRecoveryInt(err, ctx)
		}
	}()

	return vl.Render(ctx)
}

func (vl *ValueListenerInt) renderErrorWithRecoveryInt(err error, ctx *core.Context) string {
	if ctx != nil {
		fmt.Printf("[ValueListenerInt Error] %s: %v\n", time.Now().Format(time.RFC3339), err)
	}

	if vl.ErrorBuilder != nil {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("[ValueListenerInt Error Builder Panic] %v\n", r)
				}
			}()

			errorWidget := vl.ErrorBuilder(err)
			if errorWidget != nil && ctx != nil {
				result := errorWidget.Render(ctx)
				if result != "" {
					return
				}
			}
		}()
	}

	return vl.renderDefaultErrorInt(err)
}

func (vl *ValueListenerInt) renderDefaultErrorInt(err error) string {
	errorID := fmt.Sprintf("vl_int_error_%d", time.Now().UnixNano())
	userMessage := "Unable to display integer value"
	if err != nil {
		userMessage = fmt.Sprintf("Error: %s", err.Error())
	}

	return fmt.Sprintf(`
		<div id="%s" 
			 class="value-listener-error value-listener-int-error" 
			 style="color: #d32f2f; padding: 12px; border: 1px solid #f44336; border-radius: 4px; background-color: #ffebee; font-family: monospace; margin: 4px 0;"
			 data-value-type="int"
			 data-error-time="%d">
			<div style="font-weight: bold; margin-bottom: 4px;">⚠️ ValueListenerInt Error</div>
			<div style="font-size: 0.9em;">%s</div>
			<div style="margin-top: 8px; font-size: 0.8em; color: #666;">
				<button onclick="this.parentElement.parentElement.style.display='none'" 
						style="background: none; border: 1px solid #ccc; padding: 2px 8px; cursor: pointer; border-radius: 2px;">
					Dismiss
				</button>
			</div>
		</div>`,
		html.EscapeString(errorID),
		time.Now().Unix(),
		html.EscapeString(userMessage),
	)
}

// SafeRender wraps the Render method with panic recovery for ValueListenerString
func (vl *ValueListenerString) SafeRender(ctx *core.Context) (result string) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic in ValueListenerString.Render: %v", r)
			result = vl.renderErrorWithRecoveryString(err, ctx)
		}
	}()

	return vl.Render(ctx)
}

func (vl *ValueListenerString) renderErrorWithRecoveryString(err error, ctx *core.Context) string {
	if ctx != nil {
		fmt.Printf("[ValueListenerString Error] %s: %v\n", time.Now().Format(time.RFC3339), err)
	}

	if vl.ErrorBuilder != nil {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("[ValueListenerString Error Builder Panic] %v\n", r)
				}
			}()

			errorWidget := vl.ErrorBuilder(err)
			if errorWidget != nil && ctx != nil {
				result := errorWidget.Render(ctx)
				if result != "" {
					return
				}
			}
		}()
	}

	return vl.renderDefaultErrorString(err)
}

func (vl *ValueListenerString) renderDefaultErrorString(err error) string {
	errorID := fmt.Sprintf("vl_string_error_%d", time.Now().UnixNano())
	userMessage := "Unable to display string value"
	if err != nil {
		userMessage = fmt.Sprintf("Error: %s", err.Error())
	}

	return fmt.Sprintf(`
		<div id="%s" 
			 class="value-listener-error value-listener-string-error" 
			 style="color: #d32f2f; padding: 12px; border: 1px solid #f44336; border-radius: 4px; background-color: #ffebee; font-family: monospace; margin: 4px 0;"
			 data-value-type="string"
			 data-error-time="%d">
			<div style="font-weight: bold; margin-bottom: 4px;">⚠️ ValueListenerString Error</div>
			<div style="font-size: 0.9em;">%s</div>
			<div style="margin-top: 8px; font-size: 0.8em; color: #666;">
				<button onclick="this.parentElement.parentElement.style.display='none'" 
						style="background: none; border: 1px solid #ccc; padding: 2px 8px; cursor: pointer; border-radius: 2px;">
					Dismiss
				</button>
			</div>
		</div>`,
		html.EscapeString(errorID),
		time.Now().Unix(),
		html.EscapeString(userMessage),
	)
}

// SafeRender wraps the Render method with panic recovery for ValueListenerBool
func (vl *ValueListenerBool) SafeRender(ctx *core.Context) (result string) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic in ValueListenerBool.Render: %v", r)
			result = vl.renderErrorWithRecoveryBool(err, ctx)
		}
	}()

	return vl.Render(ctx)
}

func (vl *ValueListenerBool) renderErrorWithRecoveryBool(err error, ctx *core.Context) string {
	if ctx != nil {
		fmt.Printf("[ValueListenerBool Error] %s: %v\n", time.Now().Format(time.RFC3339), err)
	}

	if vl.ErrorBuilder != nil {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("[ValueListenerBool Error Builder Panic] %v\n", r)
				}
			}()

			errorWidget := vl.ErrorBuilder(err)
			if errorWidget != nil && ctx != nil {
				result := errorWidget.Render(ctx)
				if result != "" {
					return
				}
			}
		}()
	}

	return vl.renderDefaultErrorBool(err)
}

func (vl *ValueListenerBool) renderDefaultErrorBool(err error) string {
	errorID := fmt.Sprintf("vl_bool_error_%d", time.Now().UnixNano())
	userMessage := "Unable to display boolean value"
	if err != nil {
		userMessage = fmt.Sprintf("Error: %s", err.Error())
	}

	return fmt.Sprintf(`
		<div id="%s" 
			 class="value-listener-error value-listener-bool-error" 
			 style="color: #d32f2f; padding: 12px; border: 1px solid #f44336; border-radius: 4px; background-color: #ffebee; font-family: monospace; margin: 4px 0;"
			 data-value-type="bool"
			 data-error-time="%d">
			<div style="font-weight: bold; margin-bottom: 4px;">⚠️ ValueListenerBool Error</div>
			<div style="font-size: 0.9em;">%s</div>
			<div style="margin-top: 8px; font-size: 0.8em; color: #666;">
				<button onclick="this.parentElement.parentElement.style.display='none'" 
						style="background: none; border: 1px solid #ccc; padding: 2px 8px; cursor: pointer; border-radius: 2px;">
					Dismiss
				</button>
			</div>
		</div>`,
		html.EscapeString(errorID),
		time.Now().Unix(),
		html.EscapeString(userMessage),
	)
}

// SafeRender wraps the Render method with panic recovery for ValueListenerFloat64
func (vl *ValueListenerFloat64) SafeRender(ctx *core.Context) (result string) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic in ValueListenerFloat64.Render: %v", r)
			result = vl.renderErrorWithRecoveryFloat64(err, ctx)
		}
	}()

	return vl.Render(ctx)
}

func (vl *ValueListenerFloat64) renderErrorWithRecoveryFloat64(err error, ctx *core.Context) string {
	if ctx != nil {
		fmt.Printf("[ValueListenerFloat64 Error] %s: %v\n", time.Now().Format(time.RFC3339), err)
	}

	if vl.ErrorBuilder != nil {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("[ValueListenerFloat64 Error Builder Panic] %v\n", r)
				}
			}()

			errorWidget := vl.ErrorBuilder(err)
			if errorWidget != nil && ctx != nil {
				result := errorWidget.Render(ctx)
				if result != "" {
					return
				}
			}
		}()
	}

	return vl.renderDefaultErrorFloat64(err)
}

func (vl *ValueListenerFloat64) renderDefaultErrorFloat64(err error) string {
	errorID := fmt.Sprintf("vl_float64_error_%d", time.Now().UnixNano())
	userMessage := "Unable to display float64 value"
	if err != nil {
		userMessage = fmt.Sprintf("Error: %s", err.Error())
	}

	return fmt.Sprintf(`
		<div id="%s" 
			 class="value-listener-error value-listener-float64-error" 
			 style="color: #d32f2f; padding: 12px; border: 1px solid #f44336; border-radius: 4px; background-color: #ffebee; font-family: monospace; margin: 4px 0;"
			 data-value-type="float64"
			 data-error-time="%d">
			<div style="font-weight: bold; margin-bottom: 4px;">⚠️ ValueListenerFloat64 Error</div>
			<div style="font-size: 0.9em;">%s</div>
			<div style="margin-top: 8px; font-size: 0.8em; color: #666;">
				<button onclick="this.parentElement.parentElement.style.display='none'" 
						style="background: none; border: 1px solid #ccc; padding: 2px 8px; cursor: pointer; border-radius: 2px;">
					Dismiss
				</button>
			</div>
		</div>`,
		html.EscapeString(errorID),
		time.Now().Unix(),
		html.EscapeString(userMessage),
	)
}

// WebSocket fallback mechanisms

// Enhanced Render method with WebSocket fallback for the generic ValueListener
func (vl *ValueListener[T]) RenderWithFallback(ctx *core.Context) string {
	// First try the normal render
	result := vl.SafeRender(ctx)

	// Add fallback polling script if WebSocket is not available
	fallbackScript := `
		<script>
			(function() {
				// Enhanced fallback mechanism for WebSocket failures
				const element = document.getElementById('%s');
				if (!element) return;
				
				const notifierId = element.getAttribute('data-value-notifier-id');
				let websocketFailed = false;
				let pollingInterval = null;
				
				// Check if WebSocket is available and working
				function checkWebSocket() {
					if (!window.godin || !window.godin.websocket || window.godin.websocket.readyState !== WebSocket.OPEN) {
						websocketFailed = true;
						startPolling();
					}
				}
				
				// Start polling as fallback
				function startPolling() {
					if (pollingInterval) return; // Already polling
					
					pollingInterval = setInterval(function() {
						fetch('/api/state/' + notifierId)
							.then(response => {
								if (!response.ok) throw new Error('Network response was not ok');
								return response.json();
							})
							.then(data => {
								const currentValue = element.getAttribute('data-current-value');
								const newValue = JSON.stringify(data.value);
								
								if (currentValue !== newValue) {
									// Update the data attribute
									element.setAttribute('data-current-value', newValue);
									element.setAttribute('data-last-updated', Date.now().toString());
									
									// Trigger a custom event
									element.dispatchEvent(new CustomEvent('valueChanged', {
										detail: { 
											value: data.value, 
											notifierId: notifierId,
											source: 'polling',
											timestamp: Date.now()
										}
									}));
									
									// If we have HTML content, update it
									if (data.html) {
										element.innerHTML = data.html;
									}
								}
							})
							.catch(error => {
								console.warn('Polling failed for ValueListener:', error);
								// Show a subtle indicator that updates are not working
								if (!element.querySelector('.connection-warning')) {
									const warning = document.createElement('div');
									warning.className = 'connection-warning';
									warning.style.cssText = 'position: absolute; top: 0; right: 0; background: orange; color: white; padding: 2px 6px; font-size: 10px; border-radius: 0 0 0 4px;';
									warning.textContent = 'Offline';
									element.style.position = 'relative';
									element.appendChild(warning);
								}
							});
					}, 2000); // Poll every 2 seconds
				}
				
				// Stop polling when WebSocket comes back online
				function stopPolling() {
					if (pollingInterval) {
						clearInterval(pollingInterval);
						pollingInterval = null;
						
						// Remove offline indicator
						const warning = element.querySelector('.connection-warning');
						if (warning) {
							warning.remove();
						}
					}
				}
				
				// Monitor WebSocket status
				setTimeout(checkWebSocket, 1000);
				
				// Listen for WebSocket reconnection
				document.addEventListener('websocket-reconnected', stopPolling);
			})();
		</script>`

	// Generate the element ID for the fallback script
	elementID := vl.ID
	if elementID == "" && vl.ValueNotifier != nil {
		elementID = fmt.Sprintf("vl_%s", vl.ValueNotifier.ID())
	}

	return result + fmt.Sprintf(fallbackScript, html.EscapeString(elementID))
}

// Additional methods and functionality for type-specific ValueListener implementations

// Cleanup methods for type-specific implementations

// Cleanup removes the listener from the IntNotifier
func (vl *ValueListenerInt) Cleanup() {
	vl.mutex.Lock()
	defer vl.mutex.Unlock()

	if vl.ValueNotifier != nil && vl.isRegistered {
		vl.ValueNotifier.ClearListeners()
		vl.isRegistered = false
	}
}

// Cleanup removes the listener from the StringNotifier
func (vl *ValueListenerString) Cleanup() {
	vl.mutex.Lock()
	defer vl.mutex.Unlock()

	if vl.ValueNotifier != nil && vl.isRegistered {
		vl.ValueNotifier.ClearListeners()
		vl.isRegistered = false
	}
}

// Cleanup removes the listener from the BoolNotifier
func (vl *ValueListenerBool) Cleanup() {
	vl.mutex.Lock()
	defer vl.mutex.Unlock()

	if vl.ValueNotifier != nil && vl.isRegistered {
		vl.ValueNotifier.ClearListeners()
		vl.isRegistered = false
	}
}

// Cleanup removes the listener from the Float64Notifier
func (vl *ValueListenerFloat64) Cleanup() {
	vl.mutex.Lock()
	defer vl.mutex.Unlock()

	if vl.ValueNotifier != nil && vl.isRegistered {
		vl.ValueNotifier.ClearListeners()
		vl.isRegistered = false
	}
}

// Enhanced constructor functions with options for type-specific implementations

// ValueListenerIntOptions contains all configuration options for ValueListenerInt
type ValueListenerIntOptions struct {
	ID             string
	Style          string
	Class          string
	ValueNotifier  *state.IntNotifier
	Builder        func(value int) Widget
	OnValueChanged func(oldValue, newValue int)
	ErrorBuilder   func(error) Widget
}

// NewValueListenerIntWithOptions creates a new ValueListenerInt with all options
func NewValueListenerIntWithOptions(options ValueListenerIntOptions) *ValueListenerInt {
	return &ValueListenerInt{
		ID:             options.ID,
		Style:          options.Style,
		Class:          options.Class,
		ValueNotifier:  options.ValueNotifier,
		Builder:        options.Builder,
		OnValueChanged: options.OnValueChanged,
		ErrorBuilder:   options.ErrorBuilder,
		listenerID:     generateListenerID(),
	}
}

// ValueListenerStringOptions contains all configuration options for ValueListenerString
type ValueListenerStringOptions struct {
	ID             string
	Style          string
	Class          string
	ValueNotifier  *state.StringNotifier
	Builder        func(value string) Widget
	OnValueChanged func(oldValue, newValue string)
	ErrorBuilder   func(error) Widget
}

// NewValueListenerStringWithOptions creates a new ValueListenerString with all options
func NewValueListenerStringWithOptions(options ValueListenerStringOptions) *ValueListenerString {
	return &ValueListenerString{
		ID:             options.ID,
		Style:          options.Style,
		Class:          options.Class,
		ValueNotifier:  options.ValueNotifier,
		Builder:        options.Builder,
		OnValueChanged: options.OnValueChanged,
		ErrorBuilder:   options.ErrorBuilder,
		listenerID:     generateListenerID(),
	}
}

// ValueListenerBoolOptions contains all configuration options for ValueListenerBool
type ValueListenerBoolOptions struct {
	ID             string
	Style          string
	Class          string
	ValueNotifier  *state.BoolNotifier
	Builder        func(value bool) Widget
	OnValueChanged func(oldValue, newValue bool)
	ErrorBuilder   func(error) Widget
}

// NewValueListenerBoolWithOptions creates a new ValueListenerBool with all options
func NewValueListenerBoolWithOptions(options ValueListenerBoolOptions) *ValueListenerBool {
	return &ValueListenerBool{
		ID:             options.ID,
		Style:          options.Style,
		Class:          options.Class,
		ValueNotifier:  options.ValueNotifier,
		Builder:        options.Builder,
		OnValueChanged: options.OnValueChanged,
		ErrorBuilder:   options.ErrorBuilder,
		listenerID:     generateListenerID(),
	}
}

// ValueListenerFloat64Options contains all configuration options for ValueListenerFloat64
type ValueListenerFloat64Options struct {
	ID             string
	Style          string
	Class          string
	ValueNotifier  *state.Float64Notifier
	Builder        func(value float64) Widget
	OnValueChanged func(oldValue, newValue float64)
	ErrorBuilder   func(error) Widget
}

// NewValueListenerFloat64WithOptions creates a new ValueListenerFloat64 with all options
func NewValueListenerFloat64WithOptions(options ValueListenerFloat64Options) *ValueListenerFloat64 {
	return &ValueListenerFloat64{
		ID:             options.ID,
		Style:          options.Style,
		Class:          options.Class,
		ValueNotifier:  options.ValueNotifier,
		Builder:        options.Builder,
		OnValueChanged: options.OnValueChanged,
		ErrorBuilder:   options.ErrorBuilder,
		listenerID:     generateListenerID(),
	}
}

// Utility methods for type-specific implementations

// GetCurrentValue returns the current value from the IntNotifier
func (vl *ValueListenerInt) GetCurrentValue() int {
	vl.mutex.RLock()
	defer vl.mutex.RUnlock()

	if vl.ValueNotifier == nil {
		return 0
	}
	return vl.ValueNotifier.Value()
}

// GetLastValue returns the last rendered value
func (vl *ValueListenerInt) GetLastValue() int {
	vl.mutex.RLock()
	defer vl.mutex.RUnlock()
	return vl.lastValue
}

// IsRegistered returns whether this listener is registered with the StateManager
func (vl *ValueListenerInt) IsRegistered() bool {
	vl.mutex.RLock()
	defer vl.mutex.RUnlock()
	return vl.isRegistered
}

// GetCurrentValue returns the current value from the StringNotifier
func (vl *ValueListenerString) GetCurrentValue() string {
	vl.mutex.RLock()
	defer vl.mutex.RUnlock()

	if vl.ValueNotifier == nil {
		return ""
	}
	return vl.ValueNotifier.Value()
}

// GetLastValue returns the last rendered value
func (vl *ValueListenerString) GetLastValue() string {
	vl.mutex.RLock()
	defer vl.mutex.RUnlock()
	return vl.lastValue
}

// IsRegistered returns whether this listener is registered with the StateManager
func (vl *ValueListenerString) IsRegistered() bool {
	vl.mutex.RLock()
	defer vl.mutex.RUnlock()
	return vl.isRegistered
}

// GetCurrentValue returns the current value from the BoolNotifier
func (vl *ValueListenerBool) GetCurrentValue() bool {
	vl.mutex.RLock()
	defer vl.mutex.RUnlock()

	if vl.ValueNotifier == nil {
		return false
	}
	return vl.ValueNotifier.Value()
}

// GetLastValue returns the last rendered value
func (vl *ValueListenerBool) GetLastValue() bool {
	vl.mutex.RLock()
	defer vl.mutex.RUnlock()
	return vl.lastValue
}

// IsRegistered returns whether this listener is registered with the StateManager
func (vl *ValueListenerBool) IsRegistered() bool {
	vl.mutex.RLock()
	defer vl.mutex.RUnlock()
	return vl.isRegistered
}

// GetCurrentValue returns the current value from the Float64Notifier
func (vl *ValueListenerFloat64) GetCurrentValue() float64 {
	vl.mutex.RLock()
	defer vl.mutex.RUnlock()

	if vl.ValueNotifier == nil {
		return 0.0
	}
	return vl.ValueNotifier.Value()
}

// GetLastValue returns the last rendered value
func (vl *ValueListenerFloat64) GetLastValue() float64 {
	vl.mutex.RLock()
	defer vl.mutex.RUnlock()
	return vl.lastValue
}

// IsRegistered returns whether this listener is registered with the StateManager
func (vl *ValueListenerFloat64) IsRegistered() bool {
	vl.mutex.RLock()
	defer vl.mutex.RUnlock()
	return vl.isRegistered
}

// Enhanced integration methods for type-specific implementations

// RenderWithFallback renders with WebSocket fallback for ValueListenerInt
func (vl *ValueListenerInt) RenderWithFallback(ctx *core.Context) string {
	result := vl.SafeRender(ctx)

	elementID := vl.ID
	if elementID == "" && vl.ValueNotifier != nil {
		elementID = fmt.Sprintf("vl_int_%s", vl.ValueNotifier.ID())
	}

	fallbackScript := fmt.Sprintf(`
		<script>
			(function() {
				const element = document.getElementById('%s');
				if (!element) return;
				
				const notifierId = element.getAttribute('data-value-notifier-id');
				let pollingInterval = null;
				
				function startPolling() {
					if (pollingInterval) return;
					
					pollingInterval = setInterval(function() {
						fetch('/api/state/' + notifierId)
							.then(response => response.json())
							.then(data => {
								const currentValue = parseInt(element.getAttribute('data-current-value'));
								const newValue = parseInt(data.value);
								
								if (currentValue !== newValue) {
									element.setAttribute('data-current-value', newValue.toString());
									element.setAttribute('data-last-updated', Date.now().toString());
									
									element.dispatchEvent(new CustomEvent('valueChanged', {
										detail: { 
											value: newValue, 
											notifierId: notifierId,
											source: 'polling',
											type: 'int'
										}
									}));
								}
							})
							.catch(console.error);
					}, 2000);
				}
				
				// Check WebSocket availability
				setTimeout(function() {
					if (!window.godin || !window.godin.websocket || window.godin.websocket.readyState !== WebSocket.OPEN) {
						startPolling();
					}
				}, 1000);
			})();
		</script>`, html.EscapeString(elementID))

	return result + fallbackScript
}

// RenderWithFallback renders with WebSocket fallback for ValueListenerString
func (vl *ValueListenerString) RenderWithFallback(ctx *core.Context) string {
	result := vl.SafeRender(ctx)

	elementID := vl.ID
	if elementID == "" && vl.ValueNotifier != nil {
		elementID = fmt.Sprintf("vl_string_%s", vl.ValueNotifier.ID())
	}

	fallbackScript := fmt.Sprintf(`
		<script>
			(function() {
				const element = document.getElementById('%s');
				if (!element) return;
				
				const notifierId = element.getAttribute('data-value-notifier-id');
				let pollingInterval = null;
				
				function startPolling() {
					if (pollingInterval) return;
					
					pollingInterval = setInterval(function() {
						fetch('/api/state/' + notifierId)
							.then(response => response.json())
							.then(data => {
								const currentValue = element.getAttribute('data-current-value');
								const newValue = data.value;
								
								if (currentValue !== newValue) {
									element.setAttribute('data-current-value', newValue);
									element.setAttribute('data-last-updated', Date.now().toString());
									
									element.dispatchEvent(new CustomEvent('valueChanged', {
										detail: { 
											value: newValue, 
											notifierId: notifierId,
											source: 'polling',
											type: 'string'
										}
									}));
								}
							})
							.catch(console.error);
					}, 2000);
				}
				
				setTimeout(function() {
					if (!window.godin || !window.godin.websocket || window.godin.websocket.readyState !== WebSocket.OPEN) {
						startPolling();
					}
				}, 1000);
			})();
		</script>`, html.EscapeString(elementID))

	return result + fallbackScript
}

// RenderWithFallback renders with WebSocket fallback for ValueListenerBool
func (vl *ValueListenerBool) RenderWithFallback(ctx *core.Context) string {
	result := vl.SafeRender(ctx)

	elementID := vl.ID
	if elementID == "" && vl.ValueNotifier != nil {
		elementID = fmt.Sprintf("vl_bool_%s", vl.ValueNotifier.ID())
	}

	fallbackScript := fmt.Sprintf(`
		<script>
			(function() {
				const element = document.getElementById('%s');
				if (!element) return;
				
				const notifierId = element.getAttribute('data-value-notifier-id');
				let pollingInterval = null;
				
				function startPolling() {
					if (pollingInterval) return;
					
					pollingInterval = setInterval(function() {
						fetch('/api/state/' + notifierId)
							.then(response => response.json())
							.then(data => {
								const currentValue = element.getAttribute('data-current-value') === 'true';
								const newValue = data.value === true;
								
								if (currentValue !== newValue) {
									element.setAttribute('data-current-value', newValue.toString());
									element.setAttribute('data-last-updated', Date.now().toString());
									
									element.dispatchEvent(new CustomEvent('valueChanged', {
										detail: { 
											value: newValue, 
											notifierId: notifierId,
											source: 'polling',
											type: 'bool'
										}
									}));
								}
							})
							.catch(console.error);
					}, 2000);
				}
				
				setTimeout(function() {
					if (!window.godin || !window.godin.websocket || window.godin.websocket.readyState !== WebSocket.OPEN) {
						startPolling();
					}
				}, 1000);
			})();
		</script>`, html.EscapeString(elementID))

	return result + fallbackScript
}

// RenderWithFallback renders with WebSocket fallback for ValueListenerFloat64
func (vl *ValueListenerFloat64) RenderWithFallback(ctx *core.Context) string {
	result := vl.SafeRender(ctx)

	elementID := vl.ID
	if elementID == "" && vl.ValueNotifier != nil {
		elementID = fmt.Sprintf("vl_float64_%s", vl.ValueNotifier.ID())
	}

	fallbackScript := fmt.Sprintf(`
		<script>
			(function() {
				const element = document.getElementById('%s');
				if (!element) return;
				
				const notifierId = element.getAttribute('data-value-notifier-id');
				let pollingInterval = null;
				
				function startPolling() {
					if (pollingInterval) return;
					
					pollingInterval = setInterval(function() {
						fetch('/api/state/' + notifierId)
							.then(response => response.json())
							.then(data => {
								const currentValue = parseFloat(element.getAttribute('data-current-value'));
								const newValue = parseFloat(data.value);
								
								if (Math.abs(currentValue - newValue) > 0.000001) { // Float comparison with epsilon
									element.setAttribute('data-current-value', newValue.toFixed(6));
									element.setAttribute('data-last-updated', Date.now().toString());
									
									element.dispatchEvent(new CustomEvent('valueChanged', {
										detail: { 
											value: newValue, 
											notifierId: notifierId,
											source: 'polling',
											type: 'float64'
										}
									}));
								}
							})
							.catch(console.error);
					}, 2000);
				}
				
				setTimeout(function() {
					if (!window.godin || !window.godin.websocket || window.godin.websocket.readyState !== WebSocket.OPEN) {
						startPolling();
					}
				}, 1000);
			})();
		</script>`, html.EscapeString(elementID))

	return result + fallbackScript
}
