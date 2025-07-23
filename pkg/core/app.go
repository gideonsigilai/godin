package core

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gideonsigilai/godin/pkg/packages"
	"github.com/gideonsigilai/godin/pkg/state"

	"github.com/gorilla/mux"
)

// App represents the main Godin application
type App struct {
	router             *mux.Router
	server             *Server
	websocket          *WebSocketManager
	state              *state.StateManager
	packages           *packages.PackageManager
	config             *Config
	handlers           map[string]Handler // Global handler registry
	buttonCallbacks    map[string]func()  // Button callback registry for WebSocket (deprecated)
	callbackRegistry   *CallbackRegistry  // New comprehensive callback registry
	htmxIntegrator     *HTMXIntegrator    // HTMX integration system
	dialogManager      interface{}        // Dialog management system (will be properly typed later)
	navigator          interface{}        // Navigation system (will be properly typed later)
	mediaQueryProvider interface{}        // MediaQuery system (will be properly typed later)
	themeProvider      *ThemeProvider     // Theme management system
}

// Config holds application configuration
type Config struct {
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	WebSocket struct {
		Enabled bool   `yaml:"enabled"`
		Path    string `yaml:"path"`
	} `yaml:"websocket"`
	Static struct {
		Dir   string `yaml:"dir"`
		Cache bool   `yaml:"cache"`
	} `yaml:"static"`
}

// New creates a new Godin application
func New() *App {
	websocketManager := NewWebSocketManager()
	stateManager := state.NewStateManagerWithBroadcaster(websocketManager)

	app := &App{
		router:          mux.NewRouter(),
		websocket:       websocketManager,
		state:           stateManager,
		packages:        packages.NewPackageManager(),
		config:          &Config{},
		handlers:        make(map[string]Handler),
		buttonCallbacks: make(map[string]func()),
	}

	// Initialize callback registry
	app.callbackRegistry = NewCallbackRegistry(app)

	// Initialize HTMX integrator
	app.htmxIntegrator = NewHTMXIntegrator(&HTMXConfig{
		EndpointPrefix: "/api/callbacks",
		SwapStrategy:   "none",
	})

	// Initialize theme provider with default themes
	app.themeProvider = NewThemeProvider()

	// Initialize global state management for native Go code execution
	InitGlobalState()

	// Setup state API endpoints for Consumer widgets
	app.setupStateAPI()

	// Setup WebSocket button click handling
	app.setupButtonClickHandling()

	// Setup hot-reload endpoints for development
	app.setupHotReloadEndpoints()

	app.server = NewServer(app)
	return app
}

// Widget interface for rendering components
type Widget interface {
	Render(ctx *Context) string
}

// Handler represents a route handler function
type Handler func(ctx *Context) Widget

// GET registers a GET route handler
func (app *App) GET(path string, handler Handler) {
	app.router.HandleFunc(path, app.wrapHandler(handler)).Methods("GET")
}

// POST registers a POST route handler
func (app *App) POST(path string, handler Handler) {
	app.router.HandleFunc(path, app.wrapHandler(handler)).Methods("POST")
}

// PUT registers a PUT route handler
func (app *App) PUT(path string, handler Handler) {
	app.router.HandleFunc(path, app.wrapHandler(handler)).Methods("PUT")
}

// DELETE registers a DELETE route handler
func (app *App) DELETE(path string, handler Handler) {
	app.router.HandleFunc(path, app.wrapHandler(handler)).Methods("DELETE")
}

// wrapHandler wraps a Godin handler to work with HTTP
func (app *App) wrapHandler(handler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r, app)
		widget := handler(ctx)

		if widget != nil {
			// Use template rendering for full page responses
			ctx.RenderTemplate(widget, "Godin App")
		}
	}
}

// Serve starts the application server
func (app *App) Serve(addr string) error {
	return app.server.Start(addr)
}

// WebSocket returns the WebSocket manager
func (app *App) WebSocket() *WebSocketManager {
	return app.websocket
}

// GetEnvPort gets the port from environment variables
func GetEnvPort() string {
	if port := os.Getenv("PORT"); port != "" {
		if port[0] != ':' {
			return ":" + port
		}
		return port
	}
	if port := os.Getenv("GODIN_PORT"); port != "" {
		if port[0] != ':' {
			return ":" + port
		}
		return port
	}
	return ""
}

// setupHotReloadEndpoints sets up development hot-reload endpoints
func (app *App) setupHotReloadEndpoints() {
	// Only setup in development mode
	if os.Getenv("GODIN_DEV_MODE") != "true" {
		return
	}

	// Hot refresh endpoint - triggers browser refresh without server restart
	app.POST("/api/hot-refresh", func(ctx *Context) Widget {
		log.Println("🔄 Hot refresh triggered via API")

		// Broadcast hot refresh message to all connected clients
		if app.websocket.IsEnabled() {
			message := map[string]interface{}{
				"type":      "hot-refresh",
				"timestamp": time.Now().Unix(),
			}
			app.websocket.Broadcast("hot-refresh", message)
		}

		ctx.WriteJSON(map[string]string{"status": "success", "type": "hot-refresh"})
		return nil
	})

	// Hot reload endpoint - for manual triggers
	app.POST("/api/hot-reload", func(ctx *Context) Widget {
		log.Println("🔥 Hot reload triggered via API")

		// Broadcast hot reload message to all connected clients
		if app.websocket.IsEnabled() {
			message := map[string]interface{}{
				"type":      "hot-reload",
				"timestamp": time.Now().Unix(),
			}
			app.websocket.Broadcast("hot-reload", message)
		}

		ctx.WriteJSON(map[string]string{"status": "success", "type": "hot-reload"})
		return nil
	})
}

// State returns the state manager
func (app *App) State() *state.StateManager {
	return app.state
}

// RegisterHandler registers a handler globally and returns a unique ID
func (app *App) RegisterHandler(handler Handler) string {
	// Generate a unique ID for the handler
	handlerID := fmt.Sprintf("handler_%d", len(app.handlers))
	app.handlers[handlerID] = handler

	// Register the handler with the app's router
	app.router.HandleFunc("/handlers/"+handlerID, func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r, app)
		widget := handler(ctx)
		if widget != nil {
			html := widget.Render(ctx)
			ctx.WriteHTML(html)
		}
	}).Methods("GET", "POST", "PUT", "DELETE")

	return handlerID
}

// RegisterButtonCallback registers a button callback for WebSocket communication
func (app *App) RegisterButtonCallback(buttonID string, callback func()) {
	app.buttonCallbacks[buttonID] = callback
}

// GetHandlerCount returns the number of registered handlers (for generating unique IDs)
func (app *App) GetHandlerCount() int {
	return len(app.handlers)
}

// ExecuteButtonCallback executes a button callback by ID
func (app *App) ExecuteButtonCallback(buttonID string) bool {
	fmt.Printf("🔍 ExecuteButtonCallback called for buttonID: %s\n", buttonID)

	if callback, exists := app.buttonCallbacks[buttonID]; exists {
		fmt.Printf("✅ Button callback found for ID: %s\n", buttonID)

		// Create a proper context for state operations
		// We create a minimal context that has access to the app and state
		ctx := &Context{
			App:   app,
			state: make(map[string]interface{}),
		}

		fmt.Printf("🔧 Setting global context for button callback\n")
		// Set up global state context so core.SetState() and core.GetStateInt() work
		SetGlobalContext(ctx)

		fmt.Printf("🚀 Executing button callback for ID: %s\n", buttonID)
		// Execute the callback
		callback()

		fmt.Printf("🧹 Cleaning up global context\n")
		// Clean up global context
		SetGlobalContext(nil)

		fmt.Printf("✅ Button callback execution complete for ID: %s\n", buttonID)
		return true
	}

	fmt.Printf("❌ No button callback found for ID: %s\n", buttonID)
	return false
}

// setupButtonClickHandling sets up WebSocket-based button click handling
func (app *App) setupButtonClickHandling() {
	// Add HTTP endpoint for button clicks (fallback)
	app.router.HandleFunc("/api/button-click/{buttonId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		buttonID := vars["buttonId"]

		fmt.Printf("Button click received: %s\n", buttonID)

		if app.ExecuteButtonCallback(buttonID) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Button not found"))
		}
	}).Methods("POST", "GET")
}

// setupStateAPI sets up generic state API endpoints for Consumer widgets
func (app *App) setupStateAPI() {
	// Generic state endpoint that returns the current value of any state key
	app.router.HandleFunc("/api/state/{key}", func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r, app)
		key := ctx.Param("key")

		if key == "" {
			ctx.WriteText("State key is required")
			return
		}

		// Get the current state value
		value := app.state.Get(key)

		// Create a Consumer widget to render the state value
		consumer := &StateConsumer{
			StateKey: key,
			Value:    value,
		}

		html := consumer.Render(ctx)
		ctx.WriteHTML(html)
	}).Methods("GET")
}

// StateConsumer is a simple widget for rendering state values in API responses
type StateConsumer struct {
	StateKey string
	Value    interface{}
}

// Render renders the state consumer as HTML
func (sc *StateConsumer) Render(ctx *Context) string {
	if sc.Value == nil {
		return ""
	}

	// Simple rendering based on value type
	switch v := sc.Value.(type) {
	case string:
		return v
	case int:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%.2f", v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", v)
	}
}

// Router returns the underlying mux router for advanced routing
func (app *App) Router() *mux.Router {
	return app.router
}

// Packages returns the package manager
func (app *App) Packages() *packages.PackageManager {
	return app.packages
}

// CallbackRegistry returns the callback registry
func (app *App) CallbackRegistry() *CallbackRegistry {
	return app.callbackRegistry
}

// RegisterCallback registers a callback function and returns a unique callback ID
func (app *App) RegisterCallback(widgetID, widgetType, callbackType string, fn interface{}, ctx *Context) string {
	return app.callbackRegistry.RegisterCallback(widgetID, widgetType, callbackType, fn, ctx)
}

// ExecuteCallback executes a callback by ID with optional parameters
func (app *App) ExecuteCallback(callbackID string, params map[string]interface{}) error {
	return app.callbackRegistry.ExecuteCallback(callbackID, params)
}

// HTMXIntegrator returns the HTMX integrator
func (app *App) HTMXIntegrator() *HTMXIntegrator {
	return app.htmxIntegrator
}

// SetHTMXConfig updates the HTMX integrator configuration
func (app *App) SetHTMXConfig(config *HTMXConfig) {
	app.htmxIntegrator = NewHTMXIntegrator(config)
}

// GenerateCSRFTokenForApp generates and sets a CSRF token for the app
func (app *App) GenerateCSRFTokenForApp() string {
	token := GenerateCSRFToken()
	app.htmxIntegrator.SetCSRFToken(token)
	return token
}

// DialogManager returns the dialog manager
func (app *App) DialogManager() interface{} {
	return app.dialogManager
}

// SetDialogManager sets the dialog manager
func (app *App) SetDialogManager(manager interface{}) {
	app.dialogManager = manager
}

// Navigator returns the navigator
func (app *App) Navigator() interface{} {
	return app.navigator
}

// SetNavigator sets the navigator
func (app *App) SetNavigator(nav interface{}) {
	app.navigator = nav
}

// MediaQueryProvider returns the media query provider
func (app *App) MediaQueryProvider() *MediaQueryProvider {
	if provider, ok := app.mediaQueryProvider.(*MediaQueryProvider); ok {
		return provider
	}
	return nil
}

// SetMediaQueryProvider sets the media query provider
func (app *App) SetMediaQueryProvider(provider *MediaQueryProvider) {
	app.mediaQueryProvider = provider
}

// ThemeProvider returns the theme provider
func (app *App) ThemeProvider() *ThemeProvider {
	return app.themeProvider
}

// SetTheme sets the current theme
func (app *App) SetTheme(theme *ThemeData) {
	if app.themeProvider != nil {
		app.themeProvider.SetTheme(theme)
	}
}

// SetThemeMode sets the theme mode (light, dark, or system)
func (app *App) SetThemeMode(mode ThemeMode) {
	if app.themeProvider != nil {
		app.themeProvider.SetThemeMode(mode)
	}
}

// GetTheme returns the current theme
func (app *App) GetTheme() *ThemeData {
	if app.themeProvider != nil {
		return app.themeProvider.GetTheme()
	}
	return DefaultLightTheme
}

// SetLightTheme sets the light theme
func (app *App) SetLightTheme(theme *ThemeData) {
	if app.themeProvider != nil {
		app.themeProvider.SetLightTheme(theme)
	}
}

// SetDarkTheme sets the dark theme
func (app *App) SetDarkTheme(theme *ThemeData) {
	if app.themeProvider != nil {
		app.themeProvider.SetDarkTheme(theme)
	}
}

// GenerateThemeCSS generates CSS from the current theme
func (app *App) GenerateThemeCSS() string {
	if app.themeProvider != nil {
		return app.themeProvider.GenerateCSS()
	}
	return ""
}

// WithTheme creates a new app instance with a custom theme (builder pattern)
func (app *App) WithTheme(theme *ThemeData) *App {
	app.SetTheme(theme)
	return app
}

// WithLightTheme creates a new app instance with a custom light theme (builder pattern)
func (app *App) WithLightTheme(theme *ThemeData) *App {
	app.SetLightTheme(theme)
	return app
}

// WithDarkTheme creates a new app instance with a custom dark theme (builder pattern)
func (app *App) WithDarkTheme(theme *ThemeData) *App {
	app.SetDarkTheme(theme)
	return app
}

// WithThemeMode creates a new app instance with a specific theme mode (builder pattern)
func (app *App) WithThemeMode(mode ThemeMode) *App {
	app.SetThemeMode(mode)
	return app
}
