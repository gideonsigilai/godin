# Native Button System in Godin Framework

The Godin Framework now supports native Go code execution in button callbacks, similar to Flutter's `onPressed` functionality. This eliminates the need for HTMX dependencies and allows you to write pure Go code that executes when buttons are clicked.

## Key Features

- **Native Go Code Execution**: Write full Go code directly in button callbacks
- **Flutter-style API**: Uses `OnPressed` instead of HTMX attributes
- **State Management**: Built-in `setState()` function for reactive UI updates
- **Real-time Updates**: WebSocket-based state synchronization
- **No HTMX Dependencies**: Pure Go backend logic

## Basic Usage

### Simple Button with Native Go Code

```go
widgets.Button{
    Text: "Click Me",
    Type: "primary",
    OnPressed: func() {
        // Native Go code execution
        x := 1
        core.SetState("myValue", x)
    },
}
```

### Button with Complex Logic

```go
widgets.Button{
    Text: "Complex Logic",
    Type: "primary",
    OnPressed: func() {
        // Full native Go code
        counter := core.GetStateInt("counter")
        
        var message string
        var newCounter int
        
        if counter%2 == 0 {
            newCounter = counter * 2
            message = fmt.Sprintf("Even! Doubled to %d", newCounter)
        } else {
            newCounter = counter + 10
            message = fmt.Sprintf("Odd! Added 10 to %d", newCounter)
        }
        
        // Update multiple state values
        core.SetState("counter", newCounter)
        core.SetState("message", message)
    },
}
```

## State Management

### Setting State

Use the global `core.SetState()` function to update state values:

```go
core.SetState("counter", 42)
core.SetState("message", "Hello World")
core.SetState("isEnabled", true)
```

### Getting State

Retrieve state values using the global getter functions:

```go
counter := core.GetStateInt("counter")
message := core.GetStateString("message")
enabled := core.GetStateBool("isEnabled")
value := core.GetState("anyKey") // Returns interface{}
```

### Reactive UI Updates

Use the `Consumer` widget to automatically update UI when state changes:

```go
&widgets.Consumer{
    StateKey: "counter",
    Builder: func(value interface{}) widgets.Widget {
        counter := 0
        if v, ok := value.(int); ok {
            counter = v
        }
        return widgets.Text{
            Data: fmt.Sprintf("Counter: %d", counter),
            TextStyle: &widgets.TextStyle{
                FontSize: &[]float64{24}[0],
                Color:    widgets.Color("#2196F3"),
            },
        }
    },
}
```

## Complete Example

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/gideonsigilai/godin/pkg/core"
    "github.com/gideonsigilai/godin/pkg/widgets"
)

func main() {
    app := core.New()
    
    // Enable WebSocket for real-time updates
    app.WebSocket().Enable("/ws")
    
    // Initialize state
    app.State().Set("counter", 0)
    app.State().Set("message", "Welcome!")
    
    app.GET("/", func(ctx *core.Context) core.Widget {
        return homePage()
    })
    
    log.Println("Starting server on :8080")
    app.Serve(":8080")
}

func homePage() core.Widget {
    return widgets.Container{
        Style: "padding: 20px;",
        Child: widgets.Column{
            Children: []widgets.Widget{
                // State-reactive counter display
                &widgets.Consumer{
                    StateKey: "counter",
                    Builder: func(value interface{}) widgets.Widget {
                        counter := 0
                        if v, ok := value.(int); ok {
                            counter = v
                        }
                        return widgets.Text{
                            Data: fmt.Sprintf("Counter: %d", counter),
                        }
                    },
                },
                
                // Buttons with native Go code
                widgets.Row{
                    Children: []widgets.Widget{
                        widgets.Button{
                            Text: "Increment",
                            Type: "primary",
                            OnPressed: func() {
                                current := core.GetStateInt("counter")
                                core.SetState("counter", current + 1)
                                core.SetState("message", "Incremented!")
                            },
                        },
                        
                        widgets.Button{
                            Text: "Reset",
                            Type: "danger",
                            OnPressed: func() {
                                core.SetState("counter", 0)
                                core.SetState("message", "Reset!")
                            },
                        },
                    },
                },
            },
        },
    }
}
```

## Migration from HTMX Buttons

### Old HTMX Style (Deprecated)
```go
widgets.Button{
    Text:     "Increment",
    HxPost:   "/increment",
    HxTarget: "#counter-display",
    HxSwap:   "innerHTML",
}
```

### New Native Style (Recommended)
```go
widgets.Button{
    Text: "Increment",
    Type: "primary",
    OnPressed: func() {
        current := core.GetStateInt("counter")
        core.SetState("counter", current + 1)
    },
}
```

## Benefits

1. **Simplicity**: No need to create separate HTTP endpoints
2. **Type Safety**: Full Go type checking in button callbacks
3. **Performance**: Direct state updates without HTTP round trips
4. **Maintainability**: All logic in one place
5. **Debugging**: Standard Go debugging tools work
6. **Real-time**: Automatic UI updates via WebSocket

## Technical Implementation

- Button callbacks are registered as HTTP handlers automatically
- Global state manager tracks current context for `setState()` calls
- WebSocket broadcasts state changes to all connected clients
- Consumer widgets automatically re-render when their state keys change
- No HTMX attributes needed - pure Go backend logic

The native button system provides a Flutter-like development experience while maintaining the power and simplicity of Go on the backend.
