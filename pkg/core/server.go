package core

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
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

	// Find the correct path to web/static directory
	webStaticPath := s.findWebStaticPath()
	webPath := s.findWebPath()

	log.Printf("Serving static files from: %s", webStaticPath)
	log.Printf("Serving web assets from: %s", webPath)

	// Serve static files from web/static
	s.router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir(webStaticPath))),
	)

	// Serve web assets
	s.router.PathPrefix("/web/").Handler(
		http.StripPrefix("/web/", http.FileServer(http.Dir(webPath))),
	)
}

// findWebStaticPath finds the correct path to the web/static directory
func (s *Server) findWebStaticPath() string {
	// Try current directory first
	if _, err := os.Stat("web/static"); err == nil {
		return "web/static"
	}

	// Try parent directory
	if _, err := os.Stat("../web/static"); err == nil {
		return "../web/static"
	}

	// Try grandparent directory (for examples in subdirectories)
	if _, err := os.Stat("../../web/static"); err == nil {
		return "../../web/static"
	}

	// Fallback to original path
	return "web/static"
}

// findWebPath finds the correct path to the web directory
func (s *Server) findWebPath() string {
	// Try current directory first
	if _, err := os.Stat("web"); err == nil {
		return "web"
	}

	// Try parent directory
	if _, err := os.Stat("../web"); err == nil {
		return "../web"
	}

	// Try grandparent directory (for examples in subdirectories)
	if _, err := os.Stat("../../web"); err == nil {
		return "../../web"
	}

	// Fallback to original path
	return "web"
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
	app     *App
	watcher *fsnotify.Watcher
	done    chan bool
}

// NewFileWatcher creates a new file watcher
func NewFileWatcher(app *App) *FileWatcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Error creating file watcher: %v", err)
		return nil
	}

	return &FileWatcher{
		app:     app,
		watcher: watcher,
		done:    make(chan bool),
	}
}

// Watch starts watching the specified directories for changes
func (fw *FileWatcher) Watch(paths []string) {
	if fw.watcher == nil {
		log.Println("File watcher not initialized")
		return
	}

	// Add paths to watcher
	for _, path := range paths {
		err := fw.addPathRecursively(path)
		if err != nil {
			log.Printf("Error watching path %s: %v", path, err)
		}
	}

	log.Println("File watcher started for paths:", paths)

	// Start watching for events
	go fw.watchEvents()
}

// addPathRecursively adds a path and all its subdirectories to the watcher
func (fw *FileWatcher) addPathRecursively(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden directories and common ignore patterns
		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") ||
				name == "node_modules" ||
				name == "dist" ||
				name == "bin" ||
				name == "vendor" {
				return filepath.SkipDir
			}
			return fw.watcher.Add(path)
		}
		return nil
	})
}

// watchEvents processes file system events
func (fw *FileWatcher) watchEvents() {
	debounceTimer := time.NewTimer(0)
	debounceTimer.Stop()

	for {
		select {
		case event, ok := <-fw.watcher.Events:
			if !ok {
				return
			}

			// Only process relevant file changes
			if fw.shouldProcessEvent(event) {
				log.Printf("File changed: %s", event.Name)

				// Debounce rapid file changes
				debounceTimer.Reset(500 * time.Millisecond)
				go func() {
					<-debounceTimer.C
					fw.handleFileChange(event)
				}()
			}

		case err, ok := <-fw.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("File watcher error: %v", err)

		case <-fw.done:
			return
		}
	}
}

// shouldProcessEvent determines if a file change event should trigger a reload
func (fw *FileWatcher) shouldProcessEvent(event fsnotify.Event) bool {
	// Only process write and create events
	if event.Op&fsnotify.Write == 0 && event.Op&fsnotify.Create == 0 {
		return false
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(event.Name))
	watchedExtensions := []string{".go", ".html", ".css", ".js", ".yaml", ".yml"}

	for _, watchedExt := range watchedExtensions {
		if ext == watchedExt {
			return true
		}
	}

	// Also watch package.yaml specifically
	if strings.HasSuffix(event.Name, "package.yaml") {
		return true
	}

	return false
}

// handleFileChange processes a file change and triggers appropriate actions
func (fw *FileWatcher) handleFileChange(event fsnotify.Event) {
	ext := strings.ToLower(filepath.Ext(event.Name))

	switch ext {
	case ".go", ".yaml", ".yml":
		// Go files or config changes require hot reload (restart)
		fw.triggerHotReload()
	case ".html", ".css", ".js":
		// Static files can use hot refresh (no restart)
		fw.triggerHotRefresh()
	default:
		// Default to hot reload for unknown files
		fw.triggerHotReload()
	}
}

// triggerHotReload sends a hot reload signal via WebSocket
func (fw *FileWatcher) triggerHotReload() {
	if fw.app.websocket.IsEnabled() {
		message := map[string]interface{}{
			"type":      "hot-reload",
			"timestamp": time.Now().Unix(),
		}
		fw.app.websocket.Broadcast("hot-reload", message)
		log.Println("🔥 Hot reload triggered")
	}
}

// triggerHotRefresh sends a hot refresh signal via WebSocket
func (fw *FileWatcher) triggerHotRefresh() {
	if fw.app.websocket.IsEnabled() {
		message := map[string]interface{}{
			"type":      "hot-refresh",
			"timestamp": time.Now().Unix(),
		}
		fw.app.websocket.Broadcast("hot-refresh", message)
		log.Println("🔄 Hot refresh triggered")
	}
}

// Stop stops the file watcher
func (fw *FileWatcher) Stop() {
	if fw.watcher != nil {
		fw.watcher.Close()
	}
	close(fw.done)
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
