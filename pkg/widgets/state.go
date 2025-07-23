package widgets

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"sync"
	"time"

	"github.com/gideonsigilai/godin/pkg/core"
	"github.com/gideonsigilai/godin/pkg/state"
)

// ValueListenableBuilder represents a widget that rebuilds when a value changes
type ValueListenableBuilder struct {
	HTMXWidget
	ValueListenable *state.ValueListenable
	Builder         func(value interface{}) Widget
}

// Render renders the value listenable builder as HTML
func (vlb *ValueListenableBuilder) Render(ctx *core.Context) string {
	if vlb.ValueListenable == nil || vlb.Builder == nil {
		return ""
	}

	// Get current value and build widget
	value := vlb.ValueListenable.GetValue()
	widget := vlb.Builder(value)

	if widget == nil {
		return ""
	}

	return widget.Render(ctx)
}

// StreamBuilder represents a widget that rebuilds when stream data changes
type StreamBuilder struct {
	HTMXWidget
	Stream  chan interface{}
	Builder func(data interface{}) Widget
}

// Render renders the stream builder as HTML
func (sb *StreamBuilder) Render(ctx *core.Context) string {
	if sb.Stream == nil || sb.Builder == nil {
		return ""
	}

	// For server-side rendering, we'll use WebSocket for real-time updates
	// This is a simplified implementation - in practice, you'd want to
	// register the stream with the WebSocket manager

	// For now, just render with nil data
	widget := sb.Builder(nil)
	if widget == nil {
		return ""
	}

	return widget.Render(ctx)
}

// FutureBuilder represents a widget that rebuilds when future completes
type FutureBuilder struct {
	HTMXWidget
	Future  func() interface{}
	Builder func(data interface{}, loading bool, err error) Widget
}

// Render renders the future builder as HTML
func (fb *FutureBuilder) Render(ctx *core.Context) string {
	if fb.Future == nil || fb.Builder == nil {
		return ""
	}

	// For server-side rendering, we'll execute the future immediately
	// In a real implementation, you might want to handle this asynchronously
	var data interface{}
	var err error

	// Execute future
	data = fb.Future()

	widget := fb.Builder(data, false, err)
	if widget == nil {
		return ""
	}

	return widget.Render(ctx)
}

// StateBuilder represents a widget that rebuilds when state changes
type StateBuilder struct {
	HTMXWidget
	State   *state.State
	Builder func(state *state.State) Widget
}

// Render renders the state builder as HTML
func (sb *StateBuilder) Render(ctx *core.Context) string {
	if sb.State == nil || sb.Builder == nil {
		return ""
	}

	widget := sb.Builder(sb.State)
	if widget == nil {
		return ""
	}

	return widget.Render(ctx)
}

// Consumer represents a widget that consumes state changes
type Consumer struct {
	HTMXWidget
	StateKey string
	Builder  func(value interface{}) Widget
}

// Render renders the consumer as HTML
func (c *Consumer) Render(ctx *core.Context) string {
	if c.StateKey == "" || c.Builder == nil {
		return ""
	}

	// Get state from context (assuming it's available)
	stateManager := ctx.App.State()
	value := stateManager.Get(c.StateKey)

	widget := c.Builder(value)
	if widget == nil {
		return ""
	}

	// Register a custom endpoint for this specific Consumer widget
	// This ensures the state updates use the same Builder function
	consumerID := fmt.Sprintf("consumer_%s_%p", c.StateKey, c.Builder)
	endpointPath := fmt.Sprintf("/api/consumer/%s", consumerID)

	// Register the endpoint that uses this Consumer's Builder function
	ctx.App.Router().HandleFunc(endpointPath, func(w http.ResponseWriter, r *http.Request) {
		consumerCtx := core.NewContext(w, r, ctx.App)
		currentValue := ctx.App.State().Get(c.StateKey)

		// Use the same Builder function to render the updated content
		updatedWidget := c.Builder(currentValue)
		if updatedWidget != nil {
			html := updatedWidget.Render(consumerCtx)
			consumerCtx.WriteHTML(html)
		}
	}).Methods("GET")

	// Wrap the widget in a container with state tracking attributes
	// Use the custom endpoint instead of the generic state endpoint
	containerHTML := fmt.Sprintf(`<div data-state-key="%s" data-state-endpoint="%s">%s</div>`,
		c.StateKey, endpointPath, widget.Render(ctx))

	return containerHTML
}

// Provider represents a widget that provides state to its children
type Provider struct {
	HTMXWidget
	StateKey string
	Value    interface{}
	Child    Widget
}

// Render renders the provider as HTML
func (p *Provider) Render(ctx *core.Context) string {
	if p.StateKey == "" || p.Child == nil {
		return ""
	}

	// Set state in context
	stateManager := ctx.App.State()
	stateManager.Set(p.StateKey, p.Value)

	return p.Child.Render(ctx)
}

// Selector represents a widget that selects specific parts of state
type Selector struct {
	HTMXWidget
	StateKey string
	Selector func(state interface{}) interface{}
	Builder  func(selected interface{}) Widget
}

// Render renders the selector as HTML
func (s *Selector) Render(ctx *core.Context) string {
	if s.StateKey == "" || s.Selector == nil || s.Builder == nil {
		return ""
	}

	// Get state from context
	stateManager := ctx.App.State()
	state := stateManager.Get(s.StateKey)

	// Select specific part of state
	selected := s.Selector(state)

	widget := s.Builder(selected)
	if widget == nil {
		return ""
	}

	return widget.Render(ctx)
}

// ChangeNotifierProvider represents a widget that provides a change notifier
type ChangeNotifierProvider struct {
	HTMXWidget
	Notifier ChangeNotifier
	Child    Widget
}

// ChangeNotifier interface for objects that can notify about changes
type ChangeNotifier interface {
	AddListener(func())
	RemoveListener(func())
	NotifyListeners()
}

// Render renders the change notifier provider as HTML
func (cnp *ChangeNotifierProvider) Render(ctx *core.Context) string {
	if cnp.Notifier == nil || cnp.Child == nil {
		return ""
	}

	// In a real implementation, you'd register the notifier with the context
	// For now, just render the child
	return cnp.Child.Render(ctx)
}

// AnimatedBuilder represents a widget that rebuilds on animation changes
type AnimatedBuilder struct {
	HTMXWidget
	Animation Animation
	Builder   func(animation Animation) Widget
}

// Animation interface for animation objects
type Animation interface {
	GetValue() float64
	AddListener(func())
	RemoveListener(func())
}

// Render renders the animated builder as HTML
func (ab *AnimatedBuilder) Render(ctx *core.Context) string {
	if ab.Animation == nil || ab.Builder == nil {
		return ""
	}

	widget := ab.Builder(ab.Animation)
	if widget == nil {
		return ""
	}

	return widget.Render(ctx)
}

// LayoutBuilderConstraints represents layout constraints for LayoutBuilder
type LayoutBuilderConstraints struct {
	MinWidth  float64
	MaxWidth  float64
	MinHeight float64
	MaxHeight float64
}

// LayoutBuilder represents a widget that builds based on layout constraints
type LayoutBuilder struct {
	HTMXWidget
	Builder func(constraints LayoutBuilderConstraints) Widget
}

// Render renders the layout builder as HTML
func (lb *LayoutBuilder) Render(ctx *core.Context) string {
	if lb.Builder == nil {
		return ""
	}

	// For server-side rendering, use default constraints
	constraints := LayoutBuilderConstraints{
		MinWidth:  0,
		MaxWidth:  1200, // Default max width
		MinHeight: 0,
		MaxHeight: 800, // Default max height
	}

	widget := lb.Builder(constraints)
	if widget == nil {
		return ""
	}

	return widget.Render(ctx)
}

// Generic ValueListenableBuilder types for type-safe ValueNotifiers

// ValueListenableBuilderGeneric is a generic ValueListenableBuilder that works with any ValueNotifier[T]
type ValueListenableBuilderGeneric[T any] struct {
	ValueNotifier *state.ValueNotifier[T]
	Builder       func(value T) Widget
	ErrorBuilder  func(err error) Widget
	InitialValue  *T
	ID            string
	Style         string
	Class         string
}

// Render renders the generic ValueListenableBuilder widget
func (vlb ValueListenableBuilderGeneric[T]) Render(ctx *core.Context) string {
	if vlb.ValueNotifier == nil {
		if vlb.ErrorBuilder != nil {
			errorWidget := vlb.ErrorBuilder(fmt.Errorf("ValueNotifier is nil"))
			if errorWidget != nil {
				return errorWidget.Render(ctx)
			}
		}
		return `<div class="error">ValueNotifier is nil</div>`
	}

	if vlb.Builder == nil {
		if vlb.ErrorBuilder != nil {
			errorWidget := vlb.ErrorBuilder(fmt.Errorf("Builder function is nil"))
			if errorWidget != nil {
				return errorWidget.Render(ctx)
			}
		}
		return `<div class="error">Builder function is nil</div>`
	}

	// Get current value
	var currentValue T
	if vlb.InitialValue != nil {
		currentValue = *vlb.InitialValue
	} else {
		currentValue = vlb.ValueNotifier.Value()
	}

	// Build the child widget with error recovery
	var childWidget Widget
	func() {
		defer func() {
			if r := recover(); r != nil {
				if vlb.ErrorBuilder != nil {
					childWidget = vlb.ErrorBuilder(fmt.Errorf("Builder panic: %v", r))
				}
			}
		}()
		childWidget = vlb.Builder(currentValue)
	}()

	if childWidget == nil {
		if vlb.ErrorBuilder != nil {
			errorWidget := vlb.ErrorBuilder(fmt.Errorf("Builder returned nil widget"))
			if errorWidget != nil {
				return errorWidget.Render(ctx)
			}
		}
		return `<div class="error">Builder returned nil widget</div>`
	}

	// Generate unique ID if not provided
	id := vlb.ID
	if id == "" {
		id = fmt.Sprintf("vlb_generic_%s", vlb.ValueNotifier.ID())
	}

	// Create container with real-time update capabilities
	containerClass := "value-listenable-builder-generic"
	if vlb.Class != "" {
		containerClass += " " + vlb.Class
	}

	style := vlb.Style
	if style == "" {
		style = "display: contents;" // Don't add extra layout by default
	}

	// Render the child widget content
	childContent := childWidget.Render(ctx)

	// Create the container with WebSocket update capabilities
	return fmt.Sprintf(`
		<div id="%s"
			 class="%s"
			 style="%s"
			 data-value-notifier-id="%s"
			 data-current-value="%s"
			 data-value-type="%T">
			%s
		</div>
		<script>
			(function() {
				const element = document.getElementById('%s');
				const notifierId = '%s';

				// Register for WebSocket updates
				if (window.godin && window.godin.subscribe) {
					window.godin.subscribe('state:' + notifierId, function(data) {
						// Update the element with new content
						if (data && data.html) {
							element.innerHTML = data.html;
						}

						// Update data attribute
						if (data && data.value !== undefined) {
							element.setAttribute('data-current-value', JSON.stringify(data.value));
						}

						// Trigger custom event
						element.dispatchEvent(new CustomEvent('valueChanged', {
							detail: { value: data.value, notifierId: notifierId }
						}));
					});
				}

				// Fallback polling if WebSocket not available
				if (!window.godin || !window.godin.subscribe) {
					setInterval(function() {
						fetch('/api/state/' + notifierId)
							.then(response => response.json())
							.then(data => {
								if (data.html && data.html !== element.innerHTML) {
									element.innerHTML = data.html;
									element.setAttribute('data-current-value', JSON.stringify(data.value));
									element.dispatchEvent(new CustomEvent('valueChanged', {
										detail: { value: data.value, notifierId: notifierId }
									}));
								}
							})
							.catch(console.error);
					}, 1000);
				}
			})();
		</script>`,
		html.EscapeString(id),
		html.EscapeString(containerClass),
		html.EscapeString(style),
		html.EscapeString(vlb.ValueNotifier.ID()),
		html.EscapeString(jsonStringify(currentValue)),
		currentValue,
		childContent,
		html.EscapeString(id),
		html.EscapeString(vlb.ValueNotifier.ID()),
	)
}

// ValueListenableBuilderInt is a type-safe ValueListenableBuilder for int values
type ValueListenableBuilderInt struct {
	ValueListenable *state.IntNotifier
	Builder         func(value int) Widget
	ErrorBuilder    func(err error) Widget
	ID              string
	Style           string
	Class           string

	// Enhanced architecture fields for lifecycle management
	listenerID     string
	isRegistered   bool
	mutex          sync.RWMutex
	lastValue      int
	lastRenderTime time.Time
}

// NewValueListenableBuilderInt creates a new ValueListenableBuilderInt with enhanced architecture
func NewValueListenableBuilderInt(valueListenable *state.IntNotifier, builder func(value int) Widget) *ValueListenableBuilderInt {
	return &ValueListenableBuilderInt{
		ValueListenable: valueListenable,
		Builder:         builder,
		listenerID:      generateListenerID(),
	}
}

// Render renders the ValueListenableBuilder widget with enhanced architecture
func (vlb *ValueListenableBuilderInt) Render(ctx *core.Context) string {
	vlb.mutex.Lock()
	defer vlb.mutex.Unlock()

	// Register with StateManager if not already registered
	if !vlb.isRegistered && vlb.ValueListenable != nil {
		vlb.registerWithStateManager(ctx)
	}

	if vlb.ValueListenable == nil {
		return vlb.renderError(fmt.Errorf("ValueListenable is nil"))
	}

	if vlb.Builder == nil {
		return vlb.renderError(fmt.Errorf("Builder function is nil"))
	}

	// Get current value
	currentValue := vlb.ValueListenable.Value()
	vlb.lastValue = currentValue
	vlb.lastRenderTime = time.Now()

	// Build the child widget with error recovery
	var childWidget Widget
	func() {
		defer func() {
			if r := recover(); r != nil {
				if vlb.ErrorBuilder != nil {
					childWidget = vlb.ErrorBuilder(fmt.Errorf("Builder panic: %v", r))
				}
			}
		}()
		childWidget = vlb.Builder(currentValue)
	}()

	if childWidget == nil {
		return vlb.renderError(fmt.Errorf("Builder returned nil widget"))
	}

	// Generate unique ID if not provided
	id := vlb.ID
	if id == "" {
		id = fmt.Sprintf("vlb_int_%s", vlb.ValueListenable.ID())
	}

	// Create container with real-time update capabilities
	containerClass := "value-listenable-builder value-listenable-builder-int"
	if vlb.Class != "" {
		containerClass += " " + vlb.Class
	}

	style := vlb.Style
	if style == "" {
		style = "display: contents;" // Don't add extra layout by default
	}

	// Render the child widget content
	childContent := childWidget.Render(ctx)

	// Create the container with enhanced WebSocket update capabilities
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
		</div>
		<script>
			(function() {
				const element = document.getElementById('%s');
				const notifierId = '%s';
				const listenerId = '%s';

				// Initialize enhanced Godin state manager if not already done
				if (!window.godin) {
					window.godin = {
						subscriptions: new Map(),
						websocket: null,
						subscribe: function(channel, callback) {
							if (!this.subscriptions.has(channel)) {
								this.subscriptions.set(channel, []);
							}
							this.subscriptions.get(channel).push(callback);
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
								this.startPolling();
							};
							
							this.websocket.onclose = () => {
								console.log('WebSocket connection closed, attempting to reconnect...');
								setTimeout(() => this.initWebSocket(), 1000);
							};
						},
						startPolling: function() {
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
						element.setAttribute('data-current-value', data.value.toString());
						element.setAttribute('data-last-updated', Date.now().toString());
					}

					// Trigger custom event for other components to listen to
					element.dispatchEvent(new CustomEvent('valueChanged', {
						detail: { 
							value: data.value, 
							notifierId: notifierId,
							listenerId: listenerId,
							timestamp: data.timestamp,
							type: 'int'
						}
					}));
				});

				// Mark element as initialized
				element.setAttribute('data-value-listenable-builder-initialized', 'true');
			})();
		</script>`,
		html.EscapeString(id),
		html.EscapeString(containerClass),
		html.EscapeString(style),
		html.EscapeString(vlb.ValueListenable.ID()),
		html.EscapeString(vlb.listenerID),
		currentValue,
		vlb.lastRenderTime.Unix(),
		childContent,
		html.EscapeString(id),
		html.EscapeString(vlb.ValueListenable.ID()),
		html.EscapeString(vlb.listenerID),
	)
}

// registerWithStateManager registers this ValueListenableBuilder with the StateManager for lifecycle management
func (vlb *ValueListenableBuilderInt) registerWithStateManager(ctx *core.Context) {
	if vlb.ValueListenable == nil {
		return
	}

	// Register the ValueNotifier with the StateManager
	stateManager := ctx.App.State()
	stateManager.RegisterValueNotifier(vlb.ValueListenable.ID(), vlb.ValueListenable)

	// Set the StateManager on the ValueNotifier for WebSocket broadcasting
	vlb.ValueListenable.SetManager(stateManager)

	// Add listener to the ValueNotifier for change notifications
	vlb.ValueListenable.AddListener(func(newValue int) {
		vlb.lastValue = newValue
	})

	vlb.isRegistered = true
}

// renderError renders an error state for the ValueListenableBuilderInt
func (vlb *ValueListenableBuilderInt) renderError(err error) string {
	// Use custom error builder if provided
	if vlb.ErrorBuilder != nil {
		errorWidget := vlb.ErrorBuilder(err)
		if errorWidget != nil {
			return fmt.Sprintf(`<div class="value-listenable-builder-error">%s</div>`, html.EscapeString(err.Error()))
		}
	}

	// Default error rendering with enhanced styling
	return fmt.Sprintf(`
		<div class="value-listenable-builder-error" style="color: #d32f2f; padding: 8px; border: 1px solid #f44336; border-radius: 4px; background-color: #ffebee;">
			<strong>ValueListenableBuilderInt Error:</strong> %s
		</div>`,
		html.EscapeString(err.Error()),
	)
}

// Cleanup removes the listener from the ValueNotifier (should be called when widget is disposed)
func (vlb *ValueListenableBuilderInt) Cleanup() {
	vlb.mutex.Lock()
	defer vlb.mutex.Unlock()

	if vlb.ValueListenable != nil && vlb.isRegistered {
		vlb.ValueListenable.ClearListeners()
		vlb.isRegistered = false
	}
}

// ValueListenableBuilderFloat64 is a type-safe ValueListenableBuilder for float64 values
type ValueListenableBuilderFloat64 struct {
	ValueListenable *state.Float64Notifier
	Builder         func(value float64) Widget
	ErrorBuilder    func(err error) Widget
	ID              string
	Style           string
	Class           string

	// Enhanced architecture fields for lifecycle management
	listenerID     string
	isRegistered   bool
	mutex          sync.RWMutex
	lastValue      float64
	lastRenderTime time.Time
}

// NewValueListenableBuilderFloat64 creates a new ValueListenableBuilderFloat64 with enhanced architecture
func NewValueListenableBuilderFloat64(valueListenable *state.Float64Notifier, builder func(value float64) Widget) *ValueListenableBuilderFloat64 {
	return &ValueListenableBuilderFloat64{
		ValueListenable: valueListenable,
		Builder:         builder,
		listenerID:      generateListenerID(),
	}
}

// Render renders the ValueListenableBuilder widget with enhanced architecture
func (vlb *ValueListenableBuilderFloat64) Render(ctx *core.Context) string {
	vlb.mutex.Lock()
	defer vlb.mutex.Unlock()

	// Register with StateManager if not already registered
	if !vlb.isRegistered && vlb.ValueListenable != nil {
		vlb.registerWithStateManager(ctx)
	}

	if vlb.ValueListenable == nil {
		return vlb.renderError(fmt.Errorf("ValueListenable is nil"))
	}

	if vlb.Builder == nil {
		return vlb.renderError(fmt.Errorf("Builder function is nil"))
	}

	// Get current value
	currentValue := vlb.ValueListenable.Value()
	vlb.lastValue = currentValue
	vlb.lastRenderTime = time.Now()

	// Build the child widget with error recovery
	var childWidget Widget
	func() {
		defer func() {
			if r := recover(); r != nil {
				if vlb.ErrorBuilder != nil {
					childWidget = vlb.ErrorBuilder(fmt.Errorf("Builder panic: %v", r))
				}
			}
		}()
		childWidget = vlb.Builder(currentValue)
	}()

	if childWidget == nil {
		return vlb.renderError(fmt.Errorf("Builder returned nil widget"))
	}

	// Generate unique ID if not provided
	id := vlb.ID
	if id == "" {
		id = fmt.Sprintf("vlb_float64_%s", vlb.ValueListenable.ID())
	}

	// Create container with real-time update capabilities
	containerClass := "value-listenable-builder value-listenable-builder-float64"
	if vlb.Class != "" {
		containerClass += " " + vlb.Class
	}

	style := vlb.Style
	if style == "" {
		style = "display: contents;" // Don't add extra layout by default
	}

	// Render the child widget content
	childContent := childWidget.Render(ctx)

	// Create the container with enhanced WebSocket update capabilities
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
		</div>
		<script>
			(function() {
				const element = document.getElementById('%s');
				const notifierId = '%s';
				const listenerId = '%s';

				// Initialize enhanced Godin state manager if not already done
				if (!window.godin) {
					window.godin = {
						subscriptions: new Map(),
						websocket: null,
						subscribe: function(channel, callback) {
							if (!this.subscriptions.has(channel)) {
								this.subscriptions.set(channel, []);
							}
							this.subscriptions.get(channel).push(callback);
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
								this.startPolling();
							};
							
							this.websocket.onclose = () => {
								console.log('WebSocket connection closed, attempting to reconnect...');
								setTimeout(() => this.initWebSocket(), 1000);
							};
						},
						startPolling: function() {
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
						element.setAttribute('data-current-value', data.value.toString());
						element.setAttribute('data-last-updated', Date.now().toString());
					}

					// Trigger custom event for other components to listen to
					element.dispatchEvent(new CustomEvent('valueChanged', {
						detail: { 
							value: data.value, 
							notifierId: notifierId,
							listenerId: listenerId,
							timestamp: data.timestamp,
							type: 'float64'
						}
					}));
				});

				// Mark element as initialized
				element.setAttribute('data-value-listenable-builder-initialized', 'true');
			})();
		</script>`,
		html.EscapeString(id),
		html.EscapeString(containerClass),
		html.EscapeString(style),
		html.EscapeString(vlb.ValueListenable.ID()),
		html.EscapeString(vlb.listenerID),
		currentValue,
		vlb.lastRenderTime.Unix(),
		childContent,
		html.EscapeString(id),
		html.EscapeString(vlb.ValueListenable.ID()),
		html.EscapeString(vlb.listenerID),
	)
}

// registerWithStateManager registers this ValueListenableBuilder with the StateManager for lifecycle management
func (vlb *ValueListenableBuilderFloat64) registerWithStateManager(ctx *core.Context) {
	if vlb.ValueListenable == nil {
		return
	}

	// Register the ValueNotifier with the StateManager
	stateManager := ctx.App.State()
	stateManager.RegisterValueNotifier(vlb.ValueListenable.ID(), vlb.ValueListenable)

	// Set the StateManager on the ValueNotifier for WebSocket broadcasting
	vlb.ValueListenable.SetManager(stateManager)

	// Add listener to the ValueNotifier for change notifications
	vlb.ValueListenable.AddListener(func(newValue float64) {
		vlb.lastValue = newValue
	})

	vlb.isRegistered = true
}

// renderError renders an error state for the ValueListenableBuilderFloat64
func (vlb *ValueListenableBuilderFloat64) renderError(err error) string {
	// Use custom error builder if provided
	if vlb.ErrorBuilder != nil {
		errorWidget := vlb.ErrorBuilder(err)
		if errorWidget != nil {
			return fmt.Sprintf(`<div class="value-listenable-builder-error">%s</div>`, html.EscapeString(err.Error()))
		}
	}

	// Default error rendering with enhanced styling
	return fmt.Sprintf(`
		<div class="value-listenable-builder-error" style="color: #d32f2f; padding: 8px; border: 1px solid #f44336; border-radius: 4px; background-color: #ffebee;">
			<strong>ValueListenableBuilderFloat64 Error:</strong> %s
		</div>`,
		html.EscapeString(err.Error()),
	)
}

// Cleanup removes the listener from the ValueNotifier (should be called when widget is disposed)
func (vlb *ValueListenableBuilderFloat64) Cleanup() {
	vlb.mutex.Lock()
	defer vlb.mutex.Unlock()

	if vlb.ValueListenable != nil && vlb.isRegistered {
		vlb.ValueListenable.ClearListeners()
		vlb.isRegistered = false
	}
}

// ValueListenableBuilderString is a type-safe ValueListenableBuilder for string values
type ValueListenableBuilderString struct {
	ValueListenable *state.StringNotifier
	Builder         func(value string) Widget
	ErrorBuilder    func(err error) Widget
	ID              string
	Style           string
	Class           string

	// Enhanced architecture fields for lifecycle management
	listenerID     string
	isRegistered   bool
	mutex          sync.RWMutex
	lastValue      string
	lastRenderTime time.Time
}

// NewValueListenableBuilderString creates a new ValueListenableBuilderString with enhanced architecture
func NewValueListenableBuilderString(valueListenable *state.StringNotifier, builder func(value string) Widget) *ValueListenableBuilderString {
	return &ValueListenableBuilderString{
		ValueListenable: valueListenable,
		Builder:         builder,
		listenerID:      generateListenerID(),
	}
}

// Render renders the ValueListenableBuilder widget with enhanced architecture
func (vlb *ValueListenableBuilderString) Render(ctx *core.Context) string {
	vlb.mutex.Lock()
	defer vlb.mutex.Unlock()

	// Register with StateManager if not already registered
	if !vlb.isRegistered && vlb.ValueListenable != nil {
		vlb.registerWithStateManager(ctx)
	}

	if vlb.ValueListenable == nil {
		return vlb.renderError(fmt.Errorf("ValueListenable is nil"))
	}

	if vlb.Builder == nil {
		return vlb.renderError(fmt.Errorf("Builder function is nil"))
	}

	// Get current value
	currentValue := vlb.ValueListenable.Value()
	vlb.lastValue = currentValue
	vlb.lastRenderTime = time.Now()

	// Build the child widget with error recovery
	var childWidget Widget
	func() {
		defer func() {
			if r := recover(); r != nil {
				if vlb.ErrorBuilder != nil {
					childWidget = vlb.ErrorBuilder(fmt.Errorf("Builder panic: %v", r))
				}
			}
		}()
		childWidget = vlb.Builder(currentValue)
	}()

	if childWidget == nil {
		return vlb.renderError(fmt.Errorf("Builder returned nil widget"))
	}

	// Generate unique ID if not provided
	id := vlb.ID
	if id == "" {
		id = fmt.Sprintf("vlb_string_%s", vlb.ValueListenable.ID())
	}

	// Create container with real-time update capabilities
	containerClass := "value-listenable-builder value-listenable-builder-string"
	if vlb.Class != "" {
		containerClass += " " + vlb.Class
	}

	style := vlb.Style
	if style == "" {
		style = "display: contents;" // Don't add extra layout by default
	}

	// Render the child widget content
	childContent := childWidget.Render(ctx)

	// Create the container with enhanced WebSocket update capabilities
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
		</div>
		<script>
			(function() {
				const element = document.getElementById('%s');
				const notifierId = '%s';
				const listenerId = '%s';

				// Initialize enhanced Godin state manager if not already done
				if (!window.godin) {
					window.godin = {
						subscriptions: new Map(),
						websocket: null,
						subscribe: function(channel, callback) {
							if (!this.subscriptions.has(channel)) {
								this.subscriptions.set(channel, []);
							}
							this.subscriptions.get(channel).push(callback);
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
								this.startPolling();
							};
							
							this.websocket.onclose = () => {
								console.log('WebSocket connection closed, attempting to reconnect...');
								setTimeout(() => this.initWebSocket(), 1000);
							};
						},
						startPolling: function() {
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
						element.setAttribute('data-current-value', data.value.toString());
						element.setAttribute('data-last-updated', Date.now().toString());
					}

					// Trigger custom event for other components to listen to
					element.dispatchEvent(new CustomEvent('valueChanged', {
						detail: { 
							value: data.value, 
							notifierId: notifierId,
							listenerId: listenerId,
							timestamp: data.timestamp,
							type: 'string'
						}
					}));
				});

				// Mark element as initialized
				element.setAttribute('data-value-listenable-builder-initialized', 'true');
			})();
		</script>`,
		html.EscapeString(id),
		html.EscapeString(containerClass),
		html.EscapeString(style),
		html.EscapeString(vlb.ValueListenable.ID()),
		html.EscapeString(vlb.listenerID),
		html.EscapeString(currentValue),
		vlb.lastRenderTime.Unix(),
		childContent,
		html.EscapeString(id),
		html.EscapeString(vlb.ValueListenable.ID()),
		html.EscapeString(vlb.listenerID),
	)
}

// registerWithStateManager registers this ValueListenableBuilder with the StateManager for lifecycle management
func (vlb *ValueListenableBuilderString) registerWithStateManager(ctx *core.Context) {
	if vlb.ValueListenable == nil {
		return
	}

	// Register the ValueNotifier with the StateManager
	stateManager := ctx.App.State()
	stateManager.RegisterValueNotifier(vlb.ValueListenable.ID(), vlb.ValueListenable)

	// Set the StateManager on the ValueNotifier for WebSocket broadcasting
	vlb.ValueListenable.SetManager(stateManager)

	// Add listener to the ValueNotifier for change notifications
	vlb.ValueListenable.AddListener(func(newValue string) {
		vlb.lastValue = newValue
	})

	vlb.isRegistered = true
}

// renderError renders an error state for the ValueListenableBuilderString
func (vlb *ValueListenableBuilderString) renderError(err error) string {
	// Use custom error builder if provided
	if vlb.ErrorBuilder != nil {
		errorWidget := vlb.ErrorBuilder(err)
		if errorWidget != nil {
			return fmt.Sprintf(`<div class="value-listenable-builder-error">%s</div>`, html.EscapeString(err.Error()))
		}
	}

	// Default error rendering with enhanced styling
	return fmt.Sprintf(`
		<div class="value-listenable-builder-error" style="color: #d32f2f; padding: 8px; border: 1px solid #f44336; border-radius: 4px; background-color: #ffebee;">
			<strong>ValueListenableBuilderString Error:</strong> %s
		</div>`,
		html.EscapeString(err.Error()),
	)
}

// Cleanup removes the listener from the ValueNotifier (should be called when widget is disposed)
func (vlb *ValueListenableBuilderString) Cleanup() {
	vlb.mutex.Lock()
	defer vlb.mutex.Unlock()

	if vlb.ValueListenable != nil && vlb.isRegistered {
		vlb.ValueListenable.ClearListeners()
		vlb.isRegistered = false
	}
}

// ValueListenableBuilderBool is a type-safe ValueListenableBuilder for bool values
type ValueListenableBuilderBool struct {
	ValueListenable *state.BoolNotifier
	Builder         func(value bool) Widget
	ErrorBuilder    func(err error) Widget
	ID              string
	Style           string
	Class           string

	// Enhanced architecture fields for lifecycle management
	listenerID     string
	isRegistered   bool
	mutex          sync.RWMutex
	lastValue      bool
	lastRenderTime time.Time
}

// NewValueListenableBuilderBool creates a new ValueListenableBuilderBool with enhanced architecture
func NewValueListenableBuilderBool(valueListenable *state.BoolNotifier, builder func(value bool) Widget) *ValueListenableBuilderBool {
	return &ValueListenableBuilderBool{
		ValueListenable: valueListenable,
		Builder:         builder,
		listenerID:      generateListenerID(),
	}
}

// Render renders the ValueListenableBuilder widget
func (vlb ValueListenableBuilderBool) Render(ctx *core.Context) string {
	if vlb.ValueListenable == nil {
		if vlb.ErrorBuilder != nil {
			errorWidget := vlb.ErrorBuilder(fmt.Errorf("ValueListenable is nil"))
			if errorWidget != nil {
				return errorWidget.Render(ctx)
			}
		}
		return `<div class="error">ValueListenable is nil</div>`
	}

	if vlb.Builder == nil {
		if vlb.ErrorBuilder != nil {
			errorWidget := vlb.ErrorBuilder(fmt.Errorf("Builder function is nil"))
			if errorWidget != nil {
				return errorWidget.Render(ctx)
			}
		}
		return `<div class="error">Builder function is nil</div>`
	}

	currentValue := vlb.ValueListenable.Value()

	// Build the child widget with error recovery
	var childWidget Widget
	func() {
		defer func() {
			if r := recover(); r != nil {
				if vlb.ErrorBuilder != nil {
					childWidget = vlb.ErrorBuilder(fmt.Errorf("Builder panic: %v", r))
				}
			}
		}()
		childWidget = vlb.Builder(currentValue)
	}()

	if childWidget == nil {
		if vlb.ErrorBuilder != nil {
			errorWidget := vlb.ErrorBuilder(fmt.Errorf("Builder returned nil widget"))
			if errorWidget != nil {
				return errorWidget.Render(ctx)
			}
		}
		return `<div class="error">Builder returned nil widget</div>`
	}

	id := vlb.ID
	if id == "" {
		id = fmt.Sprintf("vlb_%s", vlb.ValueListenable.ID())
	}

	containerClass := "value-listenable-builder"
	if vlb.Class != "" {
		containerClass += " " + vlb.Class
	}

	style := vlb.Style
	if style == "" {
		style = "display: contents;"
	}

	childContent := childWidget.Render(ctx)

	return fmt.Sprintf(`
		<div id="%s"
			 class="%s"
			 style="%s"
			 data-value-notifier-id="%s"
			 data-current-value="%s">
			%s
		</div>
		<script>
			(function() {
				const element = document.getElementById('%s');
				const notifierId = '%s';

				// Register for WebSocket updates
				if (window.godin && window.godin.subscribe) {
					window.godin.subscribe('state:' + notifierId, function(data) {
						// Update the element with new content
						if (data && data.html) {
							element.innerHTML = data.html;
						}

						// Update data attribute
						if (data && data.value !== undefined) {
							element.setAttribute('data-current-value', JSON.stringify(data.value));
						}

						// Trigger custom event
						element.dispatchEvent(new CustomEvent('valueChanged', {
							detail: { value: data.value, notifierId: notifierId }
						}));
					});
				}

				// Fallback polling if WebSocket not available
				if (!window.godin || !window.godin.subscribe) {
					setInterval(function() {
						fetch('/api/state/' + notifierId)
							.then(response => response.json())
							.then(data => {
								if (data.html && data.html !== element.innerHTML) {
									element.innerHTML = data.html;
									element.setAttribute('data-current-value', JSON.stringify(data.value));
									element.dispatchEvent(new CustomEvent('valueChanged', {
										detail: { value: data.value, notifierId: notifierId }
									}));
								}
							})
							.catch(console.error);
					}, 1000);
				}
			})();
		</script>`,
		html.EscapeString(id),
		html.EscapeString(containerClass),
		html.EscapeString(style),
		html.EscapeString(vlb.ValueListenable.ID()),
		html.EscapeString(jsonStringify(currentValue)),
		childContent,
		html.EscapeString(id),
		html.EscapeString(vlb.ValueListenable.ID()),
	)
}

// HTML is a simple widget that renders raw HTML
type HTML struct {
	Content string
}

// Render renders the HTML widget
func (h HTML) Render(ctx *core.Context) string {
	return h.Content
}

// Helper function to safely convert value to JSON string
func jsonStringify(value interface{}) string {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprintf("%v", value)
	}
	return string(data)
}
