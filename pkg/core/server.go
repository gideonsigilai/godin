package core

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

// Server handles HTTP requests and WebSocket connections
type Server struct {
	app    *App
	router *mux.Router
}

// NewServer creates a new server instance
func NewServer(app *App) *Server {
	return &Server{
		app:    app,
		router: app.router,
	}
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	// Setup static file serving
	s.setupStaticFiles()

	// Setup WebSocket endpoint if enabled
	if s.app.websocket.IsEnabled() {
		s.setupWebSocket()
	}

	// Setup middleware
	s.setupMiddleware()

	log.Printf("Godin server starting on %s", addr)
	return http.ListenAndServe(addr, s.router)
}

// setupStaticFiles configures static file serving
func (s *Server) setupStaticFiles() {
	staticDir := s.app.config.Static.Dir
	if staticDir == "" {
		staticDir = "web/static"
	}

	// Serve static files from web/static
	s.router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))),
	)

	// Serve web assets
	s.router.PathPrefix("/web/").Handler(
		http.StripPrefix("/web/", http.FileServer(http.Dir("web"))),
	)
}

// setupWebSocket configures WebSocket endpoint
func (s *Server) setupWebSocket() {
	wsPath := s.app.websocket.GetPath()
	s.router.HandleFunc(wsPath, s.app.websocket.HandleConnection)
}

// setupMiddleware configures HTTP middleware
func (s *Server) setupMiddleware() {
	// CORS middleware
	s.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// Logging middleware
	s.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})
}

// DevServer extends Server with development features
type DevServer struct {
	*Server
	watcher *FileWatcher
}

// NewDevServer creates a development server with hot reload
func NewDevServer(app *App) *DevServer {
	server := NewServer(app)
	return &DevServer{
		Server:  server,
		watcher: NewFileWatcher(app),
	}
}

// Start starts the development server with file watching
func (ds *DevServer) Start(addr string) error {
	// Start file watcher
	go ds.watcher.Watch([]string{".", "pkg", "examples"})

	log.Printf("Godin development server starting on %s with hot reload", addr)
	return ds.Server.Start(addr)
}

// FileWatcher watches for file changes and triggers reloads
type FileWatcher struct {
	app *App
}

// NewFileWatcher creates a new file watcher
func NewFileWatcher(app *App) *FileWatcher {
	return &FileWatcher{app: app}
}

// Watch starts watching the specified directories for changes
func (fw *FileWatcher) Watch(paths []string) {
	// TODO: Implement file watching using fsnotify
	// This would watch for .go file changes and trigger:
	// 1. Recompilation
	// 2. WebSocket notification to browser for reload
	log.Println("File watcher started for paths:", paths)
}

// ServeFile serves a single file with proper content type
func ServeFile(w http.ResponseWriter, r *http.Request, filename string) error {
	ext := filepath.Ext(filename)

	switch ext {
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript")
	case ".html":
		w.Header().Set("Content-Type", "text/html")
	case ".json":
		w.Header().Set("Content-Type", "application/json")
	}

	http.ServeFile(w, r, filename)
	return nil
}
