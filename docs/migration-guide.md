# Migration Guide: From Manual HTMX to Automatic Callbacks

This guide helps you migrate existing Godin applications from manual HTMX endpoint management to the new automatic callback system.

## Overview of Changes

The new callback system introduces:
- Automatic HTMX endpoint generation
- Flutter-style callback APIs
- Built-in state management with `setState`
- TextEditingController for advanced text input
- Real-time WebSocket updates
- Comprehensive error handling

## Step-by-Step Migration

### 1. Update Button Widgets

**Before:**
```go
// Manual endpoint registration
app.POST("/increment", func(ctx *core.Context) widgets.Widget {
    counter++
    return UpdateCounterDisplay(counter)
})

// Button without automatic callback
widgets.Button{
    Text: "Increment",
    // Manual HTMX attributes would be added separately
}
```

**After:**
```go
// No manual endpoints needed!
widgets.Button{
    ID:   "increment-btn",
    Text: "Increment",
    OnPressed: func() {
        setState(func() {
            counter++
            log.Printf("Counter: %d", counter)
        })
    },
}
```

### 2. Replace Manual State Management

**Before:**
```go
var counter int

// Manual state update endpoint
app.POST("/update-counter", func(ctx *core.Context) widgets.Widget {
    // Parse form data
    value := ctx.FormValue("value")
    if v, err := strconv.Atoi(value); err == nil {
        counter = v
    }
    
    // Return updated widget
    return widgets.Text{
        Data: fmt.Sprintf("Count: %d", counter),
    }
})
```

**After:**
```go
// Use TextEditingController for reactive state
counterController := widgets.NewTextEditingController()

// Initialize with current value
counterController.SetText(fmt.Sprintf("%d", counter))

// Use ValueListenableBuilder for automatic updates
widgets.ValueListenableBuilder{
    ValueListenable: counterController,
    Builder: func(value interface{}) widgets.Widget {
        return widgets.Text{
            Data: fmt.Sprintf("Count: %s", counterController.Text),
        }
    },
}
```

### 3. Update Form Handling

**Before:**
```go
// Manual form submission endpoint
app.POST("/submit-form", func(ctx *core.Context) widgets.Widget {
    name := ctx.FormValue("name")
    email := ctx.FormValue("email")
    
    // Process form data
    processForm(name, email)
    
    return widgets.Text{Data: "Form submitted!"}
})

// Form fields without callbacks
widgets.TextField{
    ID: "name-field",
    // Manual form handling required
}
```

**After:**
```go
// Form with automatic callbacks
nameController := widgets.NewTextEditingController()
emailController := widgets.NewTextEditingController()

widgets.Column{
    Children: []widgets.Widget{
        widgets.TextField{
            ID:         "name-field",
            Controller: nameController,
            OnChanged: func(value string) {
                log.Printf("Name changed: %s", value)
                // Real-time validation if needed
            },
        },
        
        widgets.TextField{
            ID:         "email-field",
            Controller: emailController,
            OnChanged: func(value string) {
                log.Printf("Email changed: %s", value)
            },
        },
        
        widgets.Button{
            Text: "Submit",
            OnPressed: func() {
                setState(func() {
                    name := nameController.Text
                    email := emailController.Text
                    processForm(name, email)
                })
            },
        },
    },
}
```

### 4. Migrate Switch and Toggle Widgets

**Before:**
```go
var isEnabled bool

app.POST("/toggle", func(ctx *core.Context) widgets.Widget {
    isEnabled = !isEnabled
    return UpdateToggleDisplay(isEnabled)
})

// Manual toggle handling
widgets.Container{
    Child: widgets.Text{
        Data: fmt.Sprintf("Status: %t", isEnabled),
    },
}
```

**After:**
```go
var isEnabled bool

widgets.Switch{
    ID:    "toggle-switch",
    Value: isEnabled,
    OnChanged: func(value bool) {
        setState(func() {
            isEnabled = value
            log.Printf("Toggle changed: %t", isEnabled)
        })
    },
}
```

### 5. Update Custom Widgets

**Before:**
```go
type CustomWidget struct {
    ID    string
    Value string
}

func (cw CustomWidget) Render(ctx *core.Context) string {
    // Manual HTMX attribute building
    attrs := map[string]string{
        "id":         cw.ID,
        "hx-post":    "/custom-action",
        "hx-trigger": "click",
    }
    
    return htmlRenderer.RenderElement("div", attrs, cw.Value, false)
}
```

**After:**
```go
type CustomWidget struct {
    InteractiveWidget  // Embed InteractiveWidget
    ID       string
    Value    string
    OnAction func()    // Add callback
}

func (cw CustomWidget) Render(ctx *core.Context) string {
    // Initialize InteractiveWidget
    if !cw.InteractiveWidget.IsInitialized() {
        cw.InteractiveWidget.Initialize(ctx)
        cw.InteractiveWidget.SetWidgetType("CustomWidget")
    }
    
    // Register callback
    if cw.OnAction != nil {
        cw.InteractiveWidget.RegisterCallback("OnAction", cw.OnAction)
    }
    
    // Build attributes
    attrs := map[string]string{"id": cw.ID}
    
    // Merge with interactive attributes (automatic HTMX)
    attrs = cw.InteractiveWidget.MergeAttributes(attrs)
    
    return htmlRenderer.RenderElement("div", attrs, cw.Value, false)
}
```

## Common Migration Patterns

### Pattern 1: Counter/Incrementer

**Before:**
```go
// Multiple endpoints for counter operations
app.POST("/increment", IncrementHandler)
app.POST("/decrement", DecrementHandler)
app.POST("/reset", ResetHandler)

func IncrementHandler(ctx *core.Context) widgets.Widget {
    counter++
    return widgets.Text{Data: fmt.Sprintf("%d", counter)}
}
```

**After:**
```go
// Single controller with multiple callbacks
counterController := widgets.NewTextEditingController()

widgets.Row{
    Children: []widgets.Widget{
        widgets.Button{
            Text: "Decrement",
            OnPressed: func() {
                setState(func() {
                    current := parseIntFromController(counterController)
                    counterController.SetText(fmt.Sprintf("%d", current-1))
                })
            },
        },
        widgets.Button{
            Text: "Reset",
            OnPressed: func() {
                setState(func() {
                    counterController.SetText("0")
                })
            },
        },
        widgets.Button{
            Text: "Increment",
            OnPressed: func() {
                setState(func() {
                    current := parseIntFromController(counterController)
                    counterController.SetText(fmt.Sprintf("%d", current+1))
                })
            },
        },
    },
}
```

### Pattern 2: Form Validation

**Before:**
```go
app.POST("/validate", func(ctx *core.Context) widgets.Widget {
    value := ctx.FormValue("input")
    if len(value) < 3 {
        return widgets.Text{
            Data: "Too short!",
            Style: "color: red;",
        }
    }
    return widgets.Text{
        Data: "Valid!",
        Style: "color: green;",
    }
})
```

**After:**
```go
var validationMessage string
var isValid bool

widgets.Column{
    Children: []widgets.Widget{
        widgets.TextField{
            ID: "input-field",
            OnChanged: func(value string) {
                setState(func() {
                    if len(value) < 3 {
                        validationMessage = "Too short!"
                        isValid = false
                    } else {
                        validationMessage = "Valid!"
                        isValid = true
                    }
                })
            },
        },
        widgets.Text{
            Data: validationMessage,
            Style: func() string {
                if isValid {
                    return "color: green;"
                }
                return "color: red;"
            }(),
        },
    },
}
```

### Pattern 3: Dynamic Content

**Before:**
```go
app.POST("/load-content", func(ctx *core.Context) widgets.Widget {
    contentType := ctx.FormValue("type")
    return LoadContentByType(contentType)
})
```

**After:**
```go
var currentContent widgets.Widget
var contentType string

widgets.Column{
    Children: []widgets.Widget{
        widgets.Row{
            Children: []widgets.Widget{
                widgets.Button{
                    Text: "Load Type A",
                    OnPressed: func() {
                        setState(func() {
                            contentType = "A"
                            currentContent = LoadContentByType("A")
                        })
                    },
                },
                widgets.Button{
                    Text: "Load Type B",
                    OnPressed: func() {
                        setState(func() {
                            contentType = "B"
                            currentContent = LoadContentByType("B")
                        })
                    },
                },
            },
        },
        currentContent,
    },
}
```

## Migration Checklist

### Phase 1: Preparation
- [ ] Update to latest Godin version
- [ ] Review existing HTMX endpoints
- [ ] Identify widgets that need callbacks
- [ ] Plan state management strategy

### Phase 2: Core Migration
- [ ] Replace manual endpoints with callbacks
- [ ] Update button widgets to use `OnPressed`
- [ ] Migrate form widgets to use controllers
- [ ] Implement `setState` for state updates
- [ ] Add TextEditingController where needed

### Phase 3: Advanced Features
- [ ] Enable WebSocket for real-time updates
- [ ] Implement error handling
- [ ] Add comprehensive testing
- [ ] Optimize performance

### Phase 4: Cleanup
- [ ] Remove old endpoint handlers
- [ ] Clean up manual HTMX attributes
- [ ] Update documentation
- [ ] Test thoroughly

## Testing Migration

### Before Migration Test
```go
func TestOldEndpoint(t *testing.T) {
    app := core.New()
    app.POST("/increment", IncrementHandler)
    
    req := httptest.NewRequest("POST", "/increment", nil)
    w := httptest.NewRecorder()
    
    app.ServeHTTP(w, req)
    
    assert.Equal(t, 200, w.Code)
}
```

### After Migration Test
```go
func TestNewCallback(t *testing.T) {
    ctx := &core.Context{
        CallbackRegistry: core.NewCallbackRegistry(),
    }
    
    called := false
    button := widgets.Button{
        OnPressed: func() {
            called = true
        },
    }
    
    button.Render(ctx)
    
    // Execute callback
    if callbackID, exists := button.InteractiveWidget.callbacks["OnPressed"]; exists {
        ctx.CallbackRegistry.ExecuteCallback(callbackID, nil)
        assert.True(t, called)
    }
}
```

## Performance Considerations

### Before Migration
- Multiple HTTP endpoints
- Manual state synchronization
- Complex HTMX attribute management
- Potential memory leaks from unmanaged handlers

### After Migration
- Automatic endpoint generation
- Built-in state management
- Automatic cleanup
- WebSocket optimization
- Better error handling

## Troubleshooting

### Common Issues During Migration

1. **Callbacks not firing**
   - Ensure widgets have unique IDs
   - Check that InteractiveWidget is properly initialized
   - Verify callback registration in Render method

2. **State not updating**
   - Use `setState` for all state changes
   - Ensure ValueListenableBuilder is used for reactive UI
   - Check TextEditingController initialization

3. **Memory leaks**
   - Callbacks are automatically cleaned up
   - Remove manual cleanup code
   - Check for circular references in closures

4. **Performance issues**
   - Use debouncing for frequent updates
   - Batch state updates with setState
   - Enable WebSocket for real-time features

### Debug Tips

```go
// Enable debug mode
app := core.New()
app.SetDebugMode(true)

// Add logging to callbacks
widgets.Button{
    OnPressed: func() {
        log.Println("Button pressed - callback working!")
        setState(func() {
            // Your state update
        })
    },
}
```

## Benefits After Migration

1. **Reduced Boilerplate**: No manual endpoint creation
2. **Better Developer Experience**: Flutter-style APIs
3. **Automatic Optimization**: Built-in HTMX integration
4. **Real-time Updates**: WebSocket support
5. **Error Handling**: Comprehensive error recovery
6. **Type Safety**: Strongly typed callbacks
7. **Testing**: Easier unit testing
8. **Performance**: Optimized state management

## Next Steps

After completing the migration:

1. Explore advanced features like custom widgets
2. Implement real-time features with WebSocket
3. Add comprehensive error handling
4. Optimize performance with best practices
5. Create additional interactive components

For more detailed examples, see the updated button demo at `examples/button-demo/main.go`.