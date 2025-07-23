# Performance Best Practices - Godin Interactive Widget System

This guide covers performance optimization techniques for the Godin Framework's interactive widget system.

## Overview

The new callback system is designed for performance, but following these best practices will ensure optimal application performance.

## Callback Performance

### 1. Efficient Callback Registration

**Good:**
```go
// Register callbacks once during render
func (w MyWidget) Render(ctx *core.Context) string {
    if !w.InteractiveWidget.IsInitialized() {
        w.InteractiveWidget.Initialize(ctx)
        w.InteractiveWidget.SetWidgetType("MyWidget")
        
        // Register all callbacks at once
        if w.OnClick != nil {
            w.InteractiveWidget.RegisterCallback("OnClick", w.OnClick)
        }
        if w.OnHover != nil {
            w.InteractiveWidget.RegisterCallback("OnHover", w.OnHover)
        }
    }
    
    // ... rest of render logic
}
```

**Avoid:**
```go
// Don't register callbacks multiple times
func (w MyWidget) Render(ctx *core.Context) string {
    // This will register callbacks on every render - inefficient!
    w.InteractiveWidget.RegisterCallback("OnClick", w.OnClick)
    
    // ... render logic
}
```

### 2. Lightweight Callback Functions

**Good:**
```go
widgets.Button{
    OnPressed: func() {
        // Keep callbacks lightweight
        setState(func() {
            counter++
        })
    },
}
```

**Avoid:**
```go
widgets.Button{
    OnPressed: func() {
        // Avoid heavy operations in callbacks
        for i := 0; i < 1000000; i++ {
            // Heavy computation
        }
        
        // Avoid blocking operations
        time.Sleep(5 * time.Second)
        
        setState(func() {
            counter++
        })
    },
}
```

### 3. Async Operations in Callbacks

**Good:**
```go
widgets.Button{
    OnPressed: func() {
        // Use goroutines for heavy operations
        go func() {
            result := heavyComputation()
            
            // Update state after completion
            setState(func() {
                data = result
            })
        }()
    },
}
```

## State Management Performance

### 1. Batch State Updates

**Good:**
```go
// Batch multiple state changes
setState(func() {
    counter++
    message = "Updated"
    isVisible = true
    // All changes in one batch
})
```

**Avoid:**
```go
// Multiple separate state updates
setState(func() { counter++ })
setState(func() { message = "Updated" })
setState(func() { isVisible = true })
```

### 2. Minimize State Update Frequency

**Good:**
```go
// Debounce frequent updates
var debounceTimer *time.Timer

widgets.TextField{
    OnChanged: func(value string) {
        if debounceTimer != nil {
            debounceTimer.Stop()
        }
        
        debounceTimer = time.AfterFunc(300*time.Millisecond, func() {
            setState(func() {
                searchQuery = value
                performSearch()
            })
        })
    },
}
```

**Avoid:**
```go
// Update on every keystroke
widgets.TextField{
    OnChanged: func(value string) {
        setState(func() {
            searchQuery = value
            performSearch() // Expensive operation on every keystroke
        })
    },
}
```

### 3. Efficient TextEditingController Usage

**Good:**
```go
// Reuse controllers
var globalController = widgets.NewTextEditingController()

// Add listeners once
func init() {
    globalController.AddListener(func() {
        log.Printf("Text changed: %s", globalController.Text)
    })
}

// Use in multiple places
widgets.TextField{
    Controller: globalController,
}
```

**Avoid:**
```go
// Creating new controllers frequently
func SomeWidget() widgets.Widget {
    // Don't create new controllers on every render
    controller := widgets.NewTextEditingController()
    
    return widgets.TextField{
        Controller: controller,
    }
}
```

## WebSocket Performance

### 1. Efficient Channel Subscriptions

**Good:**
```go
// Subscribe to specific channels
widgets.EnhancedValueListenableBuilder{
    ValueListenable: counterNotifier,
    UpdateMode:      widgets.UpdateModeWebSocket,
    DebounceMs:      100, // Debounce updates
}
```

**Avoid:**
```go
// Don't subscribe to unnecessary channels
// Avoid too frequent updates without debouncing
widgets.EnhancedValueListenableBuilder{
    ValueListenable: counterNotifier,
    UpdateMode:      widgets.UpdateModeWebSocket,
    DebounceMs:      0, // No debouncing - can cause performance issues
}
```

### 2. Connection Management

**Good:**
```go
func main() {
    app := core.New()
    
    // Enable WebSocket with reasonable limits
    wsManager := websocket.NewRealtimeUpdateManager()
    wsManager.SetMaxConnections(1000)
    wsManager.SetMessageBufferSize(256)
    
    app.WebSocket().SetManager(wsManager)
    app.WebSocket().Enable("/ws")
    
    app.Serve(":8080")
}
```

### 3. Message Optimization

**Good:**
```go
// Send only necessary data
type StateUpdate struct {
    Key   string `json:"key"`
    Value string `json:"value"`
}

// Broadcast minimal data
wsManager.Broadcast("state:counter", StateUpdate{
    Key:   "counter",
    Value: fmt.Sprintf("%d", counter),
})
```

**Avoid:**
```go
// Don't send large objects frequently
type LargeState struct {
    Counter     int                    `json:"counter"`
    LargeData   []ComplexStruct       `json:"large_data"`
    Metadata    map[string]interface{} `json:"metadata"`
    // ... many fields
}

// Avoid broadcasting large objects
wsManager.Broadcast("state:all", largeStateObject)
```

## Memory Management

### 1. Automatic Cleanup

The system handles most cleanup automatically:

```go
// Callbacks are automatically cleaned up
// No manual cleanup needed
widgets.Button{
    OnPressed: func() {
        // This callback will be automatically cleaned up
        // when the widget is no longer rendered
    },
}
```

### 2. Avoid Memory Leaks in Closures

**Good:**
```go
func CreateButton(id string) widgets.Widget {
    return widgets.Button{
        ID: id,
        OnPressed: func() {
            // Use only necessary variables in closure
            log.Printf("Button %s pressed", id)
        },
    }
}
```

**Avoid:**
```go
func CreateButton(id string, largeData []byte) widgets.Widget {
    return widgets.Button{
        ID: id,
        OnPressed: func() {
            // Don't capture large data unnecessarily
            log.Printf("Button %s pressed", id)
            // largeData is captured but not used - memory leak
        },
    }
}
```

### 3. Controller Lifecycle Management

**Good:**
```go
// Global controllers for persistent data
var (
    usernameController = widgets.NewTextEditingController()
    passwordController = widgets.NewTextEditingController()
)

// Local controllers for temporary data
func TemporaryForm() widgets.Widget {
    tempController := widgets.NewTextEditingController()
    
    return widgets.TextField{
        Controller: tempController,
        // Controller will be garbage collected when widget is removed
    }
}
```

## Rendering Performance

### 1. Minimize Widget Complexity

**Good:**
```go
// Simple, focused widgets
func CounterDisplay(count int) widgets.Widget {
    return widgets.Text{
        Data: fmt.Sprintf("Count: %d", count),
        TextStyle: &widgets.TextStyle{
            FontSize: &[]float64{24}[0],
        },
    }
}
```

**Avoid:**
```go
// Overly complex widgets
func ComplexWidget() widgets.Widget {
    return widgets.Container{
        Child: widgets.Column{
            Children: []widgets.Widget{
                // Many nested widgets with complex logic
                // Consider breaking into smaller components
            },
        },
    }
}
```

### 2. Conditional Rendering

**Good:**
```go
func ConditionalWidget(showDetails bool) widgets.Widget {
    children := []widgets.Widget{
        widgets.Text{Data: "Always visible"},
    }
    
    if showDetails {
        children = append(children, widgets.Text{Data: "Details"})
    }
    
    return widgets.Column{Children: children}
}
```

**Avoid:**
```go
func ConditionalWidget(showDetails bool) widgets.Widget {
    return widgets.Column{
        Children: []widgets.Widget{
            widgets.Text{Data: "Always visible"},
            // Don't render invisible widgets
            widgets.Visibility{
                Visible: showDetails,
                Child:   widgets.Text{Data: "Details"},
            },
        },
    }
}
```

## Error Handling Performance

### 1. Efficient Error Recovery

**Good:**
```go
// Use built-in error handling
widgets.Button{
    OnPressed: func() {
        // Errors are automatically caught and handled
        if err := riskyOperation(); err != nil {
            log.Printf("Operation failed: %v", err)
            return
        }
        
        setState(func() {
            // Update state
        })
    },
}
```

### 2. Avoid Panic in Callbacks

**Good:**
```go
widgets.Button{
    OnPressed: func() {
        if data == nil {
            log.Println("Data not available")
            return
        }
        
        // Safe operation
        processData(data)
    },
}
```

**Avoid:**
```go
widgets.Button{
    OnPressed: func() {
        // This could panic and impact performance
        result := data.SomeField.AnotherField.Value
        processResult(result)
    },
}
```

## Monitoring and Profiling

### 1. Enable Performance Monitoring

```go
func main() {
    app := core.New()
    
    // Enable performance monitoring
    app.SetDebugMode(true)
    app.EnableProfiling("/debug/pprof")
    
    // Monitor callback performance
    app.SetCallbackMetrics(true)
    
    app.Serve(":8080")
}
```

### 2. Callback Performance Metrics

```go
// Monitor callback execution time
widgets.Button{
    OnPressed: func() {
        start := time.Now()
        defer func() {
            duration := time.Since(start)
            if duration > 100*time.Millisecond {
                log.Printf("Slow callback: %v", duration)
            }
        }()
        
        // Your callback logic
        setState(func() {
            // State update
        })
    },
}
```

### 3. Memory Usage Monitoring

```go
import (
    "runtime"
    "time"
)

func monitorMemory() {
    ticker := time.NewTicker(30 * time.Second)
    go func() {
        for range ticker.C {
            var m runtime.MemStats
            runtime.ReadMemStats(&m)
            
            log.Printf("Memory: Alloc=%d KB, Sys=%d KB, NumGC=%d",
                m.Alloc/1024, m.Sys/1024, m.NumGC)
        }
    }()
}
```

## Performance Testing

### 1. Benchmark Callbacks

```go
func BenchmarkButtonCallback(b *testing.B) {
    ctx := &core.Context{
        CallbackRegistry: core.NewCallbackRegistry(),
    }
    
    button := widgets.Button{
        OnPressed: func() {
            // Callback logic
        },
    }
    
    button.Render(ctx)
    callbackID := button.InteractiveWidget.callbacks["OnPressed"]
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ctx.CallbackRegistry.ExecuteCallback(callbackID, nil)
    }
}
```

### 2. Load Testing

```go
func TestHighLoadCallbacks(t *testing.T) {
    const numConcurrentCallbacks = 1000
    
    ctx := &core.Context{
        CallbackRegistry: core.NewCallbackRegistry(),
    }
    
    var wg sync.WaitGroup
    
    for i := 0; i < numConcurrentCallbacks; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            
            button := widgets.Button{
                OnPressed: func() {
                    // Simulate work
                    time.Sleep(time.Millisecond)
                },
            }
            
            button.Render(ctx)
            // Execute callback
        }()
    }
    
    wg.Wait()
}
```

## Summary

Key performance principles:

1. **Register callbacks once** during widget initialization
2. **Batch state updates** using setState
3. **Debounce frequent updates** to avoid performance issues
4. **Use WebSocket efficiently** with proper debouncing
5. **Keep callbacks lightweight** and use goroutines for heavy work
6. **Monitor performance** with built-in metrics
7. **Test under load** to identify bottlenecks
8. **Leverage automatic cleanup** to prevent memory leaks

Following these practices will ensure your Godin applications perform optimally with the new interactive widget system.