# Godin Widget System API Reference

## ValueListener API

### Generic ValueListener[T]

The generic ValueListener provides type-safe reactive widgets for any data type.

#### Type Definition

```go
type ValueListener[T any] struct {
    ID             string
    Style          string
    Class          string
    ValueNotifier  *state.ValueNotifier[T]
    Builder        func(value T) Widget
    OnValueChanged func(oldValue, newValue T)
    ErrorBuilder   func(error) Widget
    
    // Internal fields (read-only)
    listenerID     string
    isRegistered   bool
    mutex          sync.RWMutex
    lastValue      T
    lastRenderTime time.Time
}
```

#### Constructor Functions

```go
// Create a new generic ValueListener
func NewValueListener[T any](valueNotifier *state.ValueNotifier[T], builder func(value T) Widget) *ValueListener[T]

// Create with custom ID
func NewValueListenerWithID[T any](id string, valueNotifier *state.ValueNotifier[T], builder func(value T) Widget) *ValueListener[T]

// Create with full options
func NewValueListenerWithOptions[T any](options ValueListenerOptions[T]) *ValueListener[T]
```

#### Methods

```go
// Render the widget to HTML
func (vl *ValueListener[T]) Render(ctx *core.Context) string

// Safe render with panic recovery
func (vl *ValueListener[T]) SafeRender(ctx *core.Context) string

// Clean up resources
func (vl *ValueListener[T]) Cleanup()

// Check if registered with StateManager
func (vl *ValueListener[T]) IsRegistered() bool
```

### Type-Specific ValueListeners

#### ValueListenerInt

```go
type ValueListenerInt struct {
    ID             string
    Style          string
    Class          string
    ValueNotifier  *state.IntNotifier
    Builder        func(value int) Widget
    OnValueChanged func(oldValue, newValue int)
    ErrorBuilder   func(error) Widget
    
    // Internal fields
    listenerID     string
    isRegistered   bool
    mutex          sync.RWMutex
    lastValue      int
    lastRenderTime time.Time
}

// Constructor
func NewValueListenerInt(valueNotifier *state.IntNotifier, builder func(value int) Widget) *ValueListenerInt
func NewValueListenerIntWithID(id string, valueNotifier *state.IntNotifier, builder func(value int) Widget) *ValueListenerInt

// Methods
func (vl *ValueListenerInt) Render(ctx *core.Context) string
func (vl *ValueListenerInt) SafeRender(ctx *core.Context) string
func (vl *ValueListenerInt) Cleanup()
```

#### ValueListenerString

```go
type ValueListenerString struct {
    ID             string
    Style          string
    Class          string
    ValueNotifier  *state.StringNotifier
    Builder        func(value string) Widget
    OnValueChanged func(oldValue, newValue string)
    ErrorBuilder   func(error) Widget
    
    // Internal fields
    listenerID     string
    isRegistered   bool
    mutex          sync.RWMutex
    lastValue      string
    lastRenderTime time.Time
}

// Constructor
func NewValueListenerString(valueNotifier *state.StringNotifier, builder func(value string) Widget) *ValueListenerString
func NewValueListenerStringWithID(id string, valueNotifier *state.StringNotifier, builder func(value string) Widget) *ValueListenerString

// Methods
func (vl *ValueListenerString) Render(ctx *core.Context) string
func (vl *ValueListenerString) SafeRender(ctx *core.Context) string
func (vl *ValueListenerString) Cleanup()
```

#### ValueListenerBool

```go
type ValueListenerBool struct {
    ID             string
    Style          string
    Class          string
    ValueNotifier  *state.BoolNotifier
    Builder        func(value bool) Widget
    OnValueChanged func(oldValue, newValue bool)
    ErrorBuilder   func(error) Widget
    
    // Internal fields
    listenerID     string
    isRegistered   bool
    mutex          sync.RWMutex
    lastValue      bool
    lastRenderTime time.Time
}

// Constructor
func NewValueListenerBool(valueNotifier *state.BoolNotifier, builder func(value bool) Widget) *ValueListenerBool
func NewValueListenerBoolWithID(id string, valueNotifier *state.BoolNotifier, builder func(value bool) Widget) *ValueListenerBool

// Methods
func (vl *ValueListenerBool) Render(ctx *core.Context) string
func (vl *ValueListenerBool) SafeRender(ctx *core.Context) string
func (vl *ValueListenerBool) Cleanup()
```

#### ValueListenerFloat64

```go
type ValueListenerFloat64 struct {
    ID             string
    Style          string
    Class          string
    ValueNotifier  *state.Float64Notifier
    Builder        func(value float64) Widget
    OnValueChanged func(oldValue, newValue float64)
    ErrorBuilder   func(error) Widget
    
    // Internal fields
    listenerID     string
    isRegistered   bool
    mutex          sync.RWMutex
    lastValue      float64
    lastRenderTime time.Time
}

// Constructor
func NewValueListenerFloat64(valueNotifier *state.Float64Notifier, builder func(value float64) Widget) *ValueListenerFloat64
func NewValueListenerFloat64WithID(id string, valueNotifier *state.Float64Notifier, builder func(value float64) Widget) *ValueListenerFloat64

// Methods
func (vl *ValueListenerFloat64) Render(ctx *core.Context) string
func (vl *ValueListenerFloat64) SafeRender(ctx *core.Context) string
func (vl *ValueListenerFloat64) Cleanup()
```

## ValueListenableBuilder API

### Generic ValueListenableBuilderGeneric[T]

```go
type ValueListenableBuilderGeneric[T any] struct {
    ValueNotifier *state.ValueNotifier[T]
    Builder       func(value T) Widget
    ErrorBuilder  func(err error) Widget
    InitialValue  *T
    ID            string
    Style         string
    Class         string
}

// Methods
func (vlb ValueListenableBuilderGeneric[T]) Render(ctx *core.Context) string
```

### Type-Specific ValueListenableBuilders

#### ValueListenableBuilderInt

```go
type ValueListenableBuilderInt struct {
    ValueListenable *state.IntNotifier
    Builder         func(value int) Widget
    ErrorBuilder    func(err error) Widget
    ID              string
    Style           string
    Class           string
    
    // Enhanced architecture fields
    listenerID     string
    isRegistered   bool
    mutex          sync.RWMutex
    lastValue      int
    lastRenderTime time.Time
}

// Constructor
func NewValueListenableBuilderInt(valueListenable *state.IntNotifier, builder func(value int) Widget) *ValueListenableBuilderInt

// Methods
func (vlb *ValueListenableBuilderInt) Render(ctx *core.Context) string
func (vlb *ValueListenableBuilderInt) Cleanup()
```

#### ValueListenableBuilderString

```go
type ValueListenableBuilderString struct {
    ValueListenable *state.StringNotifier
    Builder         func(value string) Widget
    ErrorBuilder    func(err error) Widget
    ID              string
    Style           string
    Class           string
    
    // Enhanced architecture fields
    listenerID     string
    isRegistered   bool
    mutex          sync.RWMutex
    lastValue      string
    lastRenderTime time.Time
}

// Constructor
func NewValueListenableBuilderString(valueListenable *state.StringNotifier, builder func(value string) Widget) *ValueListenableBuilderString

// Methods
func (vlb *ValueListenableBuilderString) Render(ctx *core.Context) string
func (vlb *ValueListenableBuilderString) Cleanup()
```

#### ValueListenableBuilderBool

```go
type ValueListenableBuilderBool struct {
    ValueListenable *state.BoolNotifier
    Builder         func(value bool) Widget
    ErrorBuilder    func(err error) Widget
    ID              string
    Style           string
    Class           string
    
    // Enhanced architecture fields
    listenerID     string
    isRegistered   bool
    mutex          sync.RWMutex
    lastValue      bool
    lastRenderTime time.Time
}

// Constructor
func NewValueListenableBuilderBool(valueListenable *state.BoolNotifier, builder func(value bool) Widget) *ValueListenableBuilderBool

// Methods
func (vlb *ValueListenableBuilderBool) Render(ctx *core.Context) string
func (vlb *ValueListenableBuilderBool) Cleanup()
```

#### ValueListenableBuilderFloat64

```go
type ValueListenableBuilderFloat64 struct {
    ValueListenable *state.Float64Notifier
    Builder         func(value float64) Widget
    ErrorBuilder    func(err error) Widget
    ID              string
    Style           string
    Class           string
    
    // Enhanced architecture fields
    listenerID     string
    isRegistered   bool
    mutex          sync.RWMutex
    lastValue      float64
    lastRenderTime time.Time
}

// Constructor
func NewValueListenableBuilderFloat64(valueListenable *state.Float64Notifier, builder func(value float64) Widget) *ValueListenableBuilderFloat64

// Methods
func (vlb *ValueListenableBuilderFloat64) Render(ctx *core.Context) string
func (vlb *ValueListenableBuilderFloat64) Cleanup()
```

## Form Widgets API

### TextField

```go
type TextField struct {
    ID          string
    Name        string
    Type        string // "text", "email", "password", etc.
    Value       string
    Placeholder string
    Label       string
    Required    bool
    Disabled    bool
    ReadOnly    bool
    MaxLength   int
    MinLength   int
    Pattern     string
    Style       string
    Class       string
    OnChange    func(value string)
    OnFocus     func()
    OnBlur      func()
}

func (tf TextField) Render(ctx *core.Context) string
```

### Checkbox

```go
type Checkbox struct {
    ID       string
    Name     string
    Value    string
    Label    string
    Checked  bool
    Disabled bool
    Required bool
    Style    string
    Class    string
    OnChange func(checked bool)
}

func (cb Checkbox) Render(ctx *core.Context) string
```

### Radio

```go
type Radio struct {
    ID       string
    Name     string
    Value    string
    Label    string
    Checked  bool
    Disabled bool
    Required bool
    Style    string
    Class    string
    OnChange func(value string)
}

func (r Radio) Render(ctx *core.Context) string
```

### Switch

```go
type Switch struct {
    ID       string
    Name     string
    Label    string
    Checked  bool
    Disabled bool
    Style    string
    Class    string
    OnChange func(checked bool)
}

func (s Switch) Render(ctx *core.Context) string
```

## Display Widgets API

### EnhancedImage

```go
type EnhancedImage struct {
    ID         string
    Src        string
    Alt        string
    Title      string
    Width      int
    Height     int
    Loading    string // "lazy", "eager"
    Responsive bool
    Style      string
    Class      string
    OnLoad     func()
    OnError    func()
}

func (img EnhancedImage) Render(ctx *core.Context) string
```

### EnhancedIcon

```go
type EnhancedIcon struct {
    ID    string
    Name  string
    Size  int
    Color string
    Style string
    Class string
}

func (icon EnhancedIcon) Render(ctx *core.Context) string
```

### CircularProgressIndicator

```go
type CircularProgressIndicator struct {
    ID              string
    Value           *float64 // nil for indeterminate
    Color           string
    BackgroundColor string
    StrokeWidth     float64
    Size            float64
    SemanticsLabel  string
    Style           string
    Class           string
}

func (cpi CircularProgressIndicator) Render(ctx *core.Context) string
```

### LinearProgressIndicator

```go
type LinearProgressIndicator struct {
    ID              string
    Value           *float64 // nil for indeterminate
    Color           string
    BackgroundColor string
    MinHeight       float64
    SemanticsLabel  string
    Style           string
    Class           string
}

func (lpi LinearProgressIndicator) Render(ctx *core.Context) string
```

## Interactive Widgets API

### GestureDetector

```go
type GestureDetector struct {
    ID            string
    Child         Widget
    OnTap         func()
    OnDoubleTap   func()
    OnLongPress   func()
    OnHover       func()
    OnFocusChange func(focused bool)
    Style         string
    Class         string
}

func (gd GestureDetector) Render(ctx *core.Context) string
```

### InkWell

```go
type InkWell struct {
    ID          string
    Child       Widget
    OnTap       func()
    SplashColor string
    BorderRadius float64
    Style       string
    Class       string
}

func (iw InkWell) Render(ctx *core.Context) string
```

### FloatingActionButton

```go
type FloatingActionButton struct {
    ID           string
    Icon         string
    Label        string
    OnTap        func()
    BackgroundColor string
    ForegroundColor string
    Size         string // "small", "regular", "large"
    Position     string // "bottomRight", "bottomLeft", etc.
    Tooltip      string
    Style        string
    Class        string
}

func (fab FloatingActionButton) Render(ctx *core.Context) string
```

## Error Handling API

### ErrorWidget

```go
type ErrorWidget struct {
    Error    *WidgetError
    Strategy ErrorRecoveryStrategy
}

func (ew ErrorWidget) Render(ctx *core.Context) string
```

### WidgetError

```go
type WidgetError struct {
    Err         error
    Context     ErrorContext
    Recoverable bool
    RetryCount  int
}

func (we *WidgetError) Error() string
```

### ErrorContext

```go
type ErrorContext struct {
    WidgetType string
    WidgetID   string
    Operation  string
    Timestamp  time.Time
    StackTrace string
}
```

### ErrorRecoveryStrategy

```go
type ErrorRecoveryStrategy struct {
    MaxRetries     int
    RetryInterval  time.Duration
    FallbackValue  interface{}
    ErrorCallback  func(error)
    EnableLogging  bool
    UserFriendly   bool
    ShowStackTrace bool
}

// Predefined strategies
func DefaultErrorRecoveryStrategy() ErrorRecoveryStrategy
func DevelopmentErrorRecoveryStrategy() ErrorRecoveryStrategy
func ProductionErrorRecoveryStrategy() ErrorRecoveryStrategy
```

### ErrorBoundary

```go
type ErrorBoundary struct {
    ID       string
    Child    Widget
    Strategy ErrorRecoveryStrategy
    OnError  func(error)
}

func (eb ErrorBoundary) Render(ctx *core.Context) string
```

### Utility Functions

```go
// Safe widget rendering with error recovery
func SafeRenderWidget(widget Widget, ctx *core.Context, strategy ErrorRecoveryStrategy) string

// Error logging
func NewErrorLogger(strategy ErrorRecoveryStrategy) *ErrorLogger
```

## State Management API

### ValueNotifier[T]

```go
type ValueNotifier[T any] struct {
    // Private fields
}

// Constructors
func NewValueNotifier[T any](initialValue T) *ValueNotifier[T]
func NewValueNotifierWithID[T any](id string, initialValue T) *ValueNotifier[T]

// Type-specific constructors
func NewIntNotifier(value int) *IntNotifier
func NewStringNotifier(value string) *StringNotifier
func NewBoolNotifier(value bool) *BoolNotifier
func NewFloat64Notifier(value float64) *Float64Notifier

// With ID variants
func NewIntNotifierWithID(id string, value int) *IntNotifier
func NewStringNotifierWithID(id string, value string) *StringNotifier
func NewBoolNotifierWithID(id string, value bool) *BoolNotifier
func NewFloat64NotifierWithID(id string, value float64) *Float64Notifier

// Methods
func (vn *ValueNotifier[T]) Value() T
func (vn *ValueNotifier[T]) SetValue(newValue T)
func (vn *ValueNotifier[T]) Update(updater func(T) T)
func (vn *ValueNotifier[T]) AddListener(listener func(T))
func (vn *ValueNotifier[T]) RemoveListener(listener func(T))
func (vn *ValueNotifier[T]) ClearListeners()
func (vn *ValueNotifier[T]) ListenerCount() int
func (vn *ValueNotifier[T]) ID() string
func (vn *ValueNotifier[T]) SetManager(manager *StateManager)
func (vn *ValueNotifier[T]) ToJSON() ([]byte, error)
func (vn *ValueNotifier[T]) FromJSON(data []byte) error
func (vn *ValueNotifier[T]) String() string
```

### StateManager

```go
type StateManager struct {
    // Private fields
}

func NewStateManager() *StateManager

// Methods
func (sm *StateManager) Get(key string) interface{}
func (sm *StateManager) Set(key string, value interface{})
func (sm *StateManager) RegisterValueNotifier(id string, notifier interface{})
func (sm *StateManager) GetValueNotifier(id string) interface{}
```

## WebSocket Integration

### Client-Side JavaScript API

The ValueListener widgets automatically generate JavaScript code that provides:

```javascript
// Global Godin state manager
window.godin = {
    subscriptions: Map,
    websocket: WebSocket,
    
    // Subscribe to state changes
    subscribe(channel, callback),
    
    // Initialize WebSocket connection
    initWebSocket(),
    
    // Start polling fallback
    startPolling()
}

// Custom events
element.addEventListener('valueChanged', function(event) {
    // event.detail contains:
    // - value: new value
    // - notifierId: ValueNotifier ID
    // - listenerId: ValueListener ID
    // - timestamp: change timestamp
    // - type: value type
});
```

### WebSocket Message Format

```json
{
    "id": "notifier-id",
    "value": "new-value",
    "timestamp": 1234567890,
    "html": "<updated-html-content>"
}
```

## Constants and Enums

### Colors

```go
const (
    ColorPrimary    = "#007bff"
    ColorSecondary  = "#6c757d"
    ColorSuccess    = "#28a745"
    ColorDanger     = "#dc3545"
    ColorWarning    = "#ffc107"
    ColorInfo       = "#17a2b8"
    ColorLight      = "#f8f9fa"
    ColorDark       = "#343a40"
    ColorWhite      = "#ffffff"
    ColorBlack      = "#000000"
    ColorRed        = "#ff0000"
    ColorGreen      = "#00ff00"
    ColorBlue       = "#0000ff"
    ColorYellow     = "#ffff00"
    ColorPurple     = "#800080"
    ColorOrange     = "#ffa500"
    ColorGray       = "#808080"
    ColorLightGray  = "#d3d3d3"
)
```

### Widget Types

```go
const (
    WidgetTypeValueListener         = "ValueListener"
    WidgetTypeValueListenableBuilder = "ValueListenableBuilder"
    WidgetTypeTextField            = "TextField"
    WidgetTypeCheckbox             = "Checkbox"
    WidgetTypeRadio                = "Radio"
    WidgetTypeSwitch               = "Switch"
    WidgetTypeGestureDetector      = "GestureDetector"
    WidgetTypeInkWell              = "InkWell"
    WidgetTypeFloatingActionButton = "FloatingActionButton"
)
```

This API reference provides comprehensive documentation for all the enhanced widget system components, including type definitions, constructor functions, methods, and usage examples.