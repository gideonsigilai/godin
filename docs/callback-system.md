# Godin Framework - Interactive Widget Callback System

## Overview

The Godin Framework now features a comprehensive callback system that enables Flutter-style interactive widgets with automatic HTMX integration. This system eliminates the need for manual endpoint creation and provides a seamless development experience.

## Key Features

- **Automatic HTMX Integration**: Callbacks are automatically converted to HTMX endpoints
- **Flutter-Style API**: Familiar callback patterns like `OnPressed`, `OnChanged`, etc.
- **State Management**: Built-in `setState` functionality for reactive UI updates
- **TextEditingController**: Advanced text input management with listeners
- **Real-time Updates**: WebSocket integration for live state synchronization
- **Error Handling**: Comprehensive error recovery and logging
- **Type Safety**: Strongly typed callback parameters

## Core Components

### 1. InteractiveWidget

The `InteractiveWidget` is the foundation of the callback system. All interactive widgets embed this struct to gain callback functionality.

```go
type InteractiveWidget struct {
    // Internal fields for callback management
}

// Key methods
func (iw *InteractiveWidget) Initialize(ctx *core.Context)
func (iw *InteractiveWidget) RegisterCallback(eventType string, callback interface{}) string
func (iw *InteractiveWidget) MergeAttributes(attrs map[string]string) map[string]string
```

### 2. CallbackRegistry

Manages the registration, execution, and cleanup of callbacks.

```go
type CallbackRegistry struct {
    // Thread-safe callback storage
}

// Key methods
func (cr *CallbackRegistry) RegisterCallback(widgetID, eventType string, callback interface{}) string
func (cr *CallbackRegistry) ExecuteCallback(callbackID string, params interface{}) error
func (cr *CallbackRegistry) UnregisterCallback(callbackID string)
```

### 3. setState Function

Provides reactive state updates with automatic UI rebuilds.

```go
// Global setState function
func setState(updateFunc func())

// Usage example
setState(func() {
    counter++
    log.Printf("Counter updated to: %d", counter)
})
```

### 4. TextEditingController

Advanced text input management with change notifications.

```go
type TextEditingController struct {
    Text      string
    Selection TextSelection
    // Internal listener management
}

// Key methods
func NewTextEditingController() *TextEditingController
func (tec *TextEditingController) SetText(text string)
func (tec *TextEditingController) AddListener(listener func())
```

## Widget Examples

### Button Widgets

All button widgets support the new callback system:

```go
// Basic Button
widgets.Button{
    ID:   "my-button",
    Text: "Click Me",
    Type: "primary",
    OnPressed: func() {
        log.Println("Button clicked!")
        setState(func() {
            // Update state here
        })
    },
}

// ElevatedButton with multiple callbacks
widgets.ElevatedButton{
    ID: "elevated-btn",
    Child: widgets.Text{Data: "Elevated Button"},
    OnPressed: func() {
        log.Println("Pressed!")
    },
    OnLongPress: func() {
        log.Println("Long pressed!")
    },
    OnHover: func(isHovered bool) {
        log.Printf("Hover state: %t", isHovered)
    },
}

// Other button types
widgets.TextButton{...}
widgets.OutlinedButton{...}
widgets.FilledButton{...}
widgets.IconButton{...}
widgets.FloatingActionButton{...}
```

### Form Widgets

Form widgets integrate seamlessly with the callback system:

```go
// TextField with TextEditingController
controller := widgets.NewTextEditingController()

widgets.TextField{
    ID:         "text-input",
    Controller: controller,
    OnChanged: func(value string) {
        log.Printf("Text changed: %s", value)
    },
    OnSubmitted: func(value string) {
        log.Printf("Text submitted: %s", value)
    },
    OnEditingComplete: func() {
        log.Println("Editing completed")
    },
}

// TextFormField with validation
widgets.TextFormField{
    ID: "form-field",
    OnChanged: func(value string) {
        // Validate input
    },
    OnFieldSubmitted: func(value string) {
        // Handle form submission
    },
    OnSaved: func(value string) {
        // Save form data
    },
}

// Switch widget
widgets.Switch{
    ID:    "toggle-switch",
    Value: isEnabled,
    OnChanged: func(value bool) {
        setState(func() {
            isEnabled = value
            log.Printf("Switch toggled: %t", value)
        })
    },
}
```

### Advanced Widgets

```go
// InkWell with comprehensive gesture support
widgets.InkWell{
    ID: "ink-well",
    Child: widgets.Text{Data: "Tap me"},
    OnTap: func() {
        log.Println("Tapped")
    },
    OnDoubleTap: func() {
        log.Println("Double tapped")
    },
    OnLongPress: func() {
        log.Println("Long pressed")
    },
    OnHover: func(isHovered bool) {
        log.Printf("Hover: %t", isHovered)
    },
}

// GestureDetector for complex gestures
widgets.GestureDetector{
    ID: "gesture-detector",
    Child: widgets.Container{...},
    OnTap: func() { /* Handle tap */ },
    OnPanStart: func() { /* Handle pan start */ },
    OnPanUpdate: func() { /* Handle pan update */ },
    OnScaleStart: func() { /* Handle scale start */ },
}
```

## State Management

### Using setState

The `setState` function triggers UI updates automatically:

```go
var counter int

// In your widget callback
widgets.Button{
    Text: "Increment",
    OnPressed: func() {
        setState(func() {
            counter++
            // UI will automatically update
        })
    },
}
```

### TextEditingController Integration

```go
// Create controller
controller := widgets.NewTextEditingController()

// Add listener for changes
controller.AddListener(func() {
    log.Printf("Text changed to: %s", controller.Text)
})

// Use with ValueListenableBuilder for reactive UI
widgets.ValueListenableBuilder{
    ValueListenable: controller,
    Builder: func(value interface{}) widgets.Widget {
        return widgets.Text{
            Data: fmt.Sprintf("Current text: %s", controller.Text),
        }
    },
}
```

### Real-time Updates with WebSocket

Enable WebSocket for real-time state synchronization:

```go
func main() {
    app := core.New()
    
    // Enable WebSocket
    app.WebSocket().Enable("/ws")
    
    // Your routes and handlers
    app.GET("/", HomeHandler)
    
    app.Serve(":8080")
}
```

## Migration Guide

### From Manual HTMX to Automatic Callbacks

**Before (Manual HTMX):**
```go
// Manual endpoint registration
app.POST("/increment", func(ctx *core.Context) widgets.Widget {
    counter++
    return widgets.Text{Data: fmt.Sprintf("Count: %d", counter)}
})

// Manual HTMX attributes
widgets.Button{
    Text: "Increment",
    // Manual HTMX setup required
}
```

**After (Automatic Callbacks):**
```go
// No manual endpoints needed!
widgets.Button{
    Text: "Increment",
    OnPressed: func() {
        setState(func() {
            counter++
        })
    },
}
```

### Updating Existing Widgets

1. **Add InteractiveWidget**: Embed `InteractiveWidget` in your custom widgets
2. **Register Callbacks**: Use `RegisterCallback` in your `Render` method
3. **Merge Attributes**: Call `MergeAttributes` to add HTMX attributes
4. **Use setState**: Replace manual state updates with `setState`

```go
type MyCustomWidget struct {
    InteractiveWidget  // Add this
    ID       string
    OnAction func()
}

func (mcw MyCustomWidget) Render(ctx *core.Context) string {
    // Initialize InteractiveWidget
    if !mcw.InteractiveWidget.IsInitialized() {
        mcw.InteractiveWidget.Initialize(ctx)
        mcw.InteractiveWidget.SetWidgetType("MyCustomWidget")
    }
    
    // Register callback
    if mcw.OnAction != nil {
        mcw.InteractiveWidget.RegisterCallback("OnAction", mcw.OnAction)
    }
    
    // Build attributes and merge
    attrs := map[string]string{"id": mcw.ID}
    attrs = mcw.InteractiveWidget.MergeAttributes(attrs)
    
    // Render element
    return htmlRenderer.RenderElement("div", attrs, content, false)
}
```

## Performance Best Practices

### 1. Callback Registration

- Callbacks are registered once during widget rendering
- Automatic cleanup prevents memory leaks
- Thread-safe operations ensure concurrent safety

### 2. State Updates

- Use `setState` for batched updates
- Avoid frequent state changes in tight loops
- Leverage debouncing for text input callbacks

### 3. WebSocket Optimization

- WebSocket connections are automatically managed
- Fallback to polling when WebSocket unavailable
- Connection pooling and reconnection logic built-in

### 4. Memory Management

- Callbacks are automatically cleaned up
- Expired callbacks are periodically removed
- No manual cleanup required

## Error Handling

The system includes comprehensive error handling:

```go
// Automatic error recovery
widgets.Button{
    OnPressed: func() {
        // If this panics, it's automatically caught and logged
        panic("Something went wrong")
    },
}

// Custom error handling
app.SetErrorHandler(core.NewDefaultErrorHandler(logger))
```

## Testing

The callback system is fully testable:

```go
func TestButtonCallback(t *testing.T) {
    ctx := &core.Context{
        CallbackRegistry: core.NewCallbackRegistry(),
    }
    
    called := false
    button := widgets.Button{
        OnPressed: func() {
            called = true
        },
    }
    
    // Render button (registers callback)
    button.Render(ctx)
    
    // Find and execute callback
    if callbackID, exists := button.InteractiveWidget.callbacks["OnPressed"]; exists {
        ctx.CallbackRegistry.ExecuteCallback(callbackID, nil)
        assert.True(t, called)
    }
}
```

## Advanced Features

### Custom Widgets

Create your own interactive widgets:

```go
type CustomSlider struct {
    InteractiveWidget
    ID        string
    Value     float64
    OnChanged func(float64)
}

func (cs CustomSlider) Render(ctx *core.Context) string {
    if !cs.InteractiveWidget.IsInitialized() {
        cs.InteractiveWidget.Initialize(ctx)
        cs.InteractiveWidget.SetWidgetType("CustomSlider")
    }
    
    if cs.OnChanged != nil {
        cs.InteractiveWidget.RegisterCallback("OnChanged", cs.OnChanged)
    }
    
    // Implementation details...
}
```

### Event Parameters

Callbacks can receive typed parameters:

```go
widgets.TextField{
    OnChanged: func(value string) {
        // value is automatically passed from the client
    },
}

widgets.Switch{
    OnChanged: func(checked bool) {
        // checked state is automatically passed
    },
}
```

## Troubleshooting

### Common Issues

1. **Callback not firing**: Ensure widget has an ID and is properly rendered
2. **State not updating**: Use `setState` for state changes
3. **WebSocket not connecting**: Check WebSocket endpoint configuration
4. **Memory leaks**: Callbacks are auto-cleaned, but check for circular references

### Debug Mode

Enable debug logging:

```go
app := core.New()
app.SetDebugMode(true)  // Enables detailed callback logging
```

## Examples

See the complete button demo at `examples/button-demo/main.go` for a comprehensive example of the new callback system in action.

The demo showcases:
- Counter with TextEditingController
- All button types with callbacks
- Switch widgets with state management
- Form widgets with validation
- Real-time state synchronization

## Conclusion

The new callback system provides a powerful, Flutter-inspired approach to building interactive web applications with Go. It eliminates boilerplate code, provides automatic HTMX integration, and offers a familiar development experience for Flutter developers.

For more examples and advanced usage, explore the examples directory and test files in the repository.