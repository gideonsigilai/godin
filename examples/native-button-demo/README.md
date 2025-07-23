# Native Button Demo

This demo showcases the new native Go code execution system for buttons in the Godin Framework.

## Features Demonstrated

1. **Native Go Code in Button Callbacks**: All button logic is written in pure Go
2. **State Management**: Using `core.SetState()` and `core.GetState()` functions
3. **Reactive UI**: Consumer widgets that automatically update when state changes
4. **Real-time Updates**: WebSocket-based state synchronization
5. **Complex Logic**: Advanced calculations and conditional logic in button callbacks

## Running the Demo

```bash
cd examples/native-button-demo
go run main.go
```

Then visit http://localhost:8082 in your browser.

## Button Examples

### Simple State Update
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

### Complex Logic
```go
widgets.Button{
    Text: "Complex Logic",
    OnPressed: func() {
        counter := core.GetStateInt("counter")
        
        if counter%2 == 0 {
            newValue := counter * 2
            core.SetState("counter", newValue)
            core.SetState("message", fmt.Sprintf("Even! Doubled to %d", newValue))
        } else {
            newValue := counter + 10
            core.SetState("counter", newValue)
            core.SetState("message", fmt.Sprintf("Odd! Added 10 to %d", newValue))
        }
    },
}
```

## Key Differences from HTMX Approach

- **No HTTP endpoints needed**: Button logic runs directly in Go
- **No HTMX attributes**: Uses `OnPressed` callback instead
- **Type safety**: Full Go type checking and IDE support
- **Simpler debugging**: Standard Go debugging tools work
- **Better performance**: No HTTP round trips for simple state updates

## State Management

The demo uses the global state functions:

- `core.SetState(key, value)` - Update state and trigger UI updates
- `core.GetStateInt(key)` - Get integer state value
- `core.GetStateString(key)` - Get string state value
- `core.GetStateBool(key)` - Get boolean state value
- `core.GetState(key)` - Get any state value as interface{}

## UI Updates

The `Consumer` widget automatically re-renders when its associated state key changes:

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
        }
    },
}
```

This provides a reactive UI experience similar to Flutter's state management.
