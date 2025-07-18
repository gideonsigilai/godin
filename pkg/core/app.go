package core

import (
	"godin-framework/pkg/packages"
	"godin-framework/pkg/state"
	"net/http"

	"github.com/gorilla/mux"
)

// App represents the main Godin application
type App struct {
	router    *mux.Router
	server    *Server
	websocket *WebSocketManager
	state     *state.StateManager
	packages  *packages.PackageManager
	config    *Config
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
		router:    mux.NewRouter(),
		websocket: websocketManager,
		state:     stateManager,
		packages:  packages.NewPackageManager(),
		config:    &Config{},
	}

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

// State returns the state manager
func (app *App) State() *state.StateManager {
	return app.state
}

// Router returns the underlying mux router for advanced routing
func (app *App) Router() *mux.Router {
	return app.router
}

// Packages returns the package manager
func (app *App) Packages() *packages.PackageManager {
	return app.packages
}
