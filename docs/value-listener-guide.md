# ValueListener Widget Guide

## Overview

ValueListener is a powerful widget system in Godin that automatically rebuilds UI components when underlying data changes. It provides reactive programming capabilities with type safety, WebSocket integration, and comprehensive error handling.

## Table of Contents

1. [Basic Usage](#basic-usage)
2. [Type-Specific Implementations](#type-specific-implementations)
3. [Advanced Features](#advanced-features)
4. [Error Handling](#error-handling)
5. [Performance Considerations](#performance-considerations)
6. [Best Practices](#best-practices)
7. [API Reference](#api-reference)

## Basic Usage

### Creating a ValueListener

```go
package main

import (
    "fmt"
    "github.com/gideonsigilai/godin/pkg/state"
    "github.com/gideonsigilai/godin/pkg/widgets"
)

func main() {
    // Create a ValueNotifier
    counter := state.NewIntNotifier(0)
    
    // Create a ValueListener that rebuilds when counter changes
    listener := widgets.NewValueListenerInt(counter, func(value int) widgets.Widget {
        return &widgets.Text{
            Data: fmt.Sprintf("Count: %d", value),
        }
    })
    
    // Use the listener in your widget tree
    // The listener will automatically update when counter.SetValue() is called
}
```

### Updating Values

```go
// Update the value - all listeners will automatically rebuild
counter.SetValue(5)
counter.SetValue(10)

// Or use the Update method for functional updates
counter.Update(func(current int) int {
    return current + 1
})
```

## Type-Specific Implementations

Godin provides type-specific ValueListener implementations for common data types:

### IntNotifier and ValueListenerInt

```go
// Create an integer notifier
counter := state.NewIntNotifier(0)

// Create a listener for integer values
listener := widgets.NewValueListenerInt(counter, func(value int) widgets.Widget {
    return &widgets.Text{
        Data: fmt.Sprintf("Counter: %d", value),
        Style: "font-size: 18px; color: blue;",
    }
})
```

### StringNotifier and ValueListenerString

```go
// Create a string notifier
message := state.NewStringNotifier("Hello")

// Create a listener for string values
listener := widgets.NewValueListenerString(message, func(value string) widgets.Widget {
    return &widgets.Text{
        Data: fmt.Sprintf("Message: %s", value),
        Class: "message-display",
    }
})
```

### BoolNotifier and ValueListenerBool

```go
// Create a boolean notifier
isVisible := state.NewBoolNotifier(true)

// Create a listener for boolean values
listener := widgets.NewValueListenerBool(isVisible, func(value bool) widgets.Widget {
    if value {
        return &widgets.Text{Data: "Visible"}
    }
    return &widgets.Text{Data: "Hidden"}
})
```

### Float64Notifier and ValueListenerFloat64

```go
// Create a float64 notifier
temperature := state.NewFloat64Notifier(23.5)

// Create a listener for float64 values
listener := widgets.NewValueListenerFloat64(temperature, func(value float64) widgets.Widget {
    return &widgets.Text{
        Data: fmt.Sprintf("Temperature: %.1f°C", value),
    }
})
```

## Advanced Features

### Custom IDs and Styling

```go
listener := widgets.NewValueListenerInt(counter, func(value int) widgets.Widget {
    return &widgets.Text{Data: fmt.Sprintf("Count: %d", value)}
})

// Set custom properties
listener.ID = "my-counter"
listener.Style = "padding: 10px; border: 1px solid #ccc;"
listener.Class = "counter-widget primary"
```

### Value Change Callbacks

```go
listener := widgets.NewValueListenerInt(counter, func(value int) widgets.Widget {
    return &widgets.Text{Data: fmt.Sprintf("Count: %d", value)}
})

// Add a callback for value changes
listener.OnValueChanged = func(oldValue, newValue int) {
    fmt.Printf("Counter changed from %d to %d\n", oldValue, newValue)
}
```

### Generic ValueListener

For custom types, use the generic ValueListener:

```go
type User struct {
    Name  string
    Email string
}

// Create a generic notifier
userNotifier := state.NewValueNotifier(User{Name: "John", Email: "john@example.com"})

// Create a generic listener
listener := widgets.NewValueListener(userNotifier, func(user User) widgets.Widget {
    return &widgets.Container{
        Children: []widgets.Widget{
            &widgets.Text{Data: fmt.Sprintf("Name: %s", user.Name)},
            &widgets.Text{Data: fmt.Sprintf("Email: %s", user.Email)},
        },
    }
})
```

## Error Handling

### Basic Error Handling

```go
listener := widgets.NewValueListenerInt(counter, func(value int) widgets.Widget {
    return &widgets.Text{Data: fmt.Sprintf("Count: %d", value)}
})

// Add custom error handling
listener.ErrorBuilder = func(err error) widgets.Widget {
    return &widgets.Text{
        Data:  fmt.Sprintf("Error: %s", err.Error()),
        Style: "color: red; padding: 10px;",
    }
}
```

### Safe Rendering

```go
// Use SafeRender for automatic panic recovery
html := listener.SafeRender(ctx)
```

### Error Recovery Strategies

```go
// Configure error recovery
strategy := widgets.ErrorRecoveryStrategy{
    MaxRetries:    3,
    RetryInterval: time.Second * 2,
    EnableLogging: true,
    UserFriendly:  true,
}

// Use with error boundary
boundary := widgets.ErrorBoundary{
    ID:       "error-boundary",
    Child:    listener,
    Strategy: strategy,
}
```

## Performance Considerations

### Memory Management

```go
// Always cleanup listeners when done
defer listener.Cleanup()

// Or use lifecycle management
listener.OnDispose = func() {
    listener.Cleanup()
}
```

### Optimizing Renders

```go
// Use change detection to prevent unnecessary renders
listener := widgets.NewValueListenerInt(counter, func(value int) widgets.Widget {
    // Only rebuild if value actually changed
    if value == listener.lastValue {
        return nil // Skip rebuild
    }
    
    return &widgets.Text{Data: fmt.Sprintf("Count: %d", value)}
})
```

### Batching Updates

```go
// Batch multiple updates
counter.BeginUpdate()
counter.SetValue(1)
counter.SetValue(2)
counter.SetValue(3)
counter.EndUpdate() // Only triggers one rebuild
```

## Best Practices

### 1. Use Type-Specific Listeners

```go
// Preferred: Type-specific
listener := widgets.NewValueListenerInt(intNotifier, builder)

// Avoid: Generic when type-specific is available
listener := widgets.NewValueListener(intNotifier, builder)
```

### 2. Keep Builder Functions Simple

```go
// Good: Simple, focused builder
listener := widgets.NewValueListenerInt(counter, func(value int) widgets.Widget {
    return &widgets.Text{Data: fmt.Sprintf("Count: %d", value)}
})

// Avoid: Complex logic in builder
listener := widgets.NewValueListenerInt(counter, func(value int) widgets.Widget {
    // Complex calculations, API calls, etc.
    // This should be done elsewhere
})
```

### 3. Use Meaningful IDs

```go
listener.ID = "user-profile-counter"  // Good
listener.ID = "widget1"               // Avoid
```

### 4. Handle Errors Gracefully

```go
listener.ErrorBuilder = func(err error) widgets.Widget {
    return &widgets.ErrorWidget{
        Message: "Unable to display content",
        Details: err.Error(),
        Retry:   true,
    }
}
```

### 5. Cleanup Resources

```go
func (c *Component) Dispose() {
    c.listener.Cleanup()
    c.notifier.ClearListeners()
}
```

## API Reference

### ValueListener[T]

#### Constructor Functions

```go
// Generic constructor
func NewValueListener[T any](valueNotifier *state.ValueNotifier[T], builder func(value T) Widget) *ValueListener[T]

// With ID
func NewValueListenerWithID[T any](id string, valueNotifier *state.ValueNotifier[T], builder func(value T) Widget) *ValueListener[T]

// With options
func NewValueListenerWithOptions[T any](options ValueListenerOptions[T]) *ValueListener[T]
```

#### Type-Specific Constructors

```go
func NewValueListenerInt(valueNotifier *state.IntNotifier, builder func(value int) Widget) *ValueListenerInt
func NewValueListenerString(valueNotifier *state.StringNotifier, builder func(value string) Widget) *ValueListenerString
func NewValueListenerBool(valueNotifier *state.BoolNotifier, builder func(value bool) Widget) *ValueListenerBool
func NewValueListenerFloat64(valueNotifier *state.Float64Notifier, builder func(value float64) Widget) *ValueListenerFloat64
```

#### Properties

```go
type ValueListener[T] struct {
    ID             string                    // Widget ID
    Style          string                    // CSS styles
    Class          string                    // CSS classes
    ValueNotifier  *state.ValueNotifier[T]   // Data source
    Builder        func(value T) Widget      // Widget builder function
    OnValueChanged func(oldValue, newValue T) // Change callback
    ErrorBuilder   func(error) Widget        // Error widget builder
}
```

#### Methods

```go
// Render the widget
func (vl *ValueListener[T]) Render(ctx *core.Context) string

// Safe render with panic recovery
func (vl *ValueListener[T]) SafeRender(ctx *core.Context) string

// Cleanup resources
func (vl *ValueListener[T]) Cleanup()
```

### ValueNotifier[T]

#### Constructor Functions

```go
func NewValueNotifier[T any](initialValue T) *ValueNotifier[T]
func NewValueNotifierWithID[T any](id string, initialValue T) *ValueNotifier[T]
```

#### Type-Specific Constructors

```go
func NewIntNotifier(value int) *IntNotifier
func NewStringNotifier(value string) *StringNotifier
func NewBoolNotifier(value bool) *BoolNotifier
func NewFloat64Notifier(value float64) *Float64Notifier
```

#### Methods

```go
// Get current value
func (vn *ValueNotifier[T]) Value() T

// Set new value
func (vn *ValueNotifier[T]) SetValue(newValue T)

// Functional update
func (vn *ValueNotifier[T]) Update(updater func(T) T)

// Add change listener
func (vn *ValueNotifier[T]) AddListener(listener func(T))

// Clear all listeners
func (vn *ValueNotifier[T]) ClearListeners()

// Get unique ID
func (vn *ValueNotifier[T]) ID() string
```

## WebSocket Integration

ValueListener widgets automatically generate JavaScript code for real-time updates via WebSocket:

```html
<!-- Generated HTML includes WebSocket subscription -->
<div id="vl_counter" 
     class="value-listener value-listener-int"
     data-value-notifier-id="counter-123"
     data-current-value="5"
     data-value-type="int">
    <span>Count: 5</span>
</div>
<script>
    // Automatic WebSocket subscription
    window.godin.subscribe('state:counter-123', function(data) {
        // Update DOM when value changes
        element.innerHTML = data.html;
        element.setAttribute('data-current-value', data.value);
        
        // Trigger custom event
        element.dispatchEvent(new CustomEvent('valueChanged', {
            detail: { value: data.value, notifierId: 'counter-123' }
        }));
    });
</script>
```

## Migration Guide

### From Manual State Management

```go
// Before: Manual state management
type Component struct {
    counter int
}

func (c *Component) Render() string {
    return fmt.Sprintf("<span>Count: %d</span>", c.counter)
}

func (c *Component) Increment() {
    c.counter++
    // Manual re-render required
}

// After: ValueListener
counter := state.NewIntNotifier(0)
listener := widgets.NewValueListenerInt(counter, func(value int) widgets.Widget {
    return &widgets.Text{Data: fmt.Sprintf("Count: %d", value)}
})

// Automatic re-render
counter.SetValue(counter.Value() + 1)
```

### From Basic ValueNotifier

```go
// Before: Basic ValueNotifier usage
notifier := state.NewIntNotifier(0)
// Manual subscription and rendering

// After: ValueListener
listener := widgets.NewValueListenerInt(notifier, func(value int) widgets.Widget {
    return &widgets.Text{Data: fmt.Sprintf("Count: %d", value)}
})
// Automatic subscription and rendering
```

## Troubleshooting

### Common Issues

1. **Widget not updating**: Ensure ValueNotifier.SetValue() is called
2. **Memory leaks**: Always call Cleanup() when disposing widgets
3. **Performance issues**: Use type-specific listeners and avoid complex builders
4. **WebSocket not working**: Check WebSocket server configuration

### Debug Mode

```go
// Enable debug mode for detailed logging
listener.Debug = true

// Check registration status
if !listener.IsRegistered() {
    fmt.Println("Listener not registered with StateManager")
}
```

## Examples

See the [examples directory](../examples/) for complete working examples:

- [Counter App](../examples/counter-app/) - Basic ValueListener usage
- [Form Example](../examples/form-example/) - Form widgets with state management
- [Real-time Chat](../examples/chat-app/) - WebSocket integration
- [Dashboard](../examples/dashboard/) - Complex state management

## Contributing

To contribute to ValueListener development:

1. Read the [contributing guidelines](../CONTRIBUTING.md)
2. Check the [architecture documentation](./architecture.md)
3. Run tests: `go test ./pkg/widgets -v`
4. Submit pull requests with comprehensive tests

## License

This documentation is part of the Godin framework and is licensed under the same terms as the main project.