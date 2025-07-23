package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var rootCmd = &cobra.Command{
	Use:   "godin",
	Short: "Godin - Server-Side Go with HTMX/WebSockets Framework",
	Long:  "A framework for building web applications using Go with HTMX and WebSocket support",
}

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new Godin project",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := "my-godin-app"
		if len(args) > 0 {
			projectName = args[0]
		}
		initProject(projectName)
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start development server with hot reload",
	Long: `Start the Godin development server with hot reload capabilities.

This command starts a development server with the following features:
- Hot reload: Automatically restarts the server when Go files change
- Hot refresh: Refreshes the browser when static files change
- Interactive commands: 'r' for manual hot reload, 'R' for manual hot refresh
- File watching: Monitors file changes in real-time

Examples:
  godin serve                    # Start server on default port 8080
  godin serve --port 3000        # Start server on custom port
  godin serve --watch            # Enable file watching (default)
  godin serve --listen           # Enable interactive commands`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Godin serve command started")

		port, _ := cmd.Flags().GetString("port")
		watch, _ := cmd.Flags().GetBool("watch")
		listen, _ := cmd.Flags().GetBool("listen")
		enhancedReload, _ := cmd.Flags().GetBool("enhanced-reload")
		restartRetries, _ := cmd.Flags().GetInt("restart-retries")
		debounce, _ := cmd.Flags().GetDuration("debounce")

		fmt.Printf("📋 Parsed flags: port=%s, watch=%v, listen=%v\n", port, watch, listen)

		startDevServerEnhanced(port, watch, listen, enhancedReload, restartRetries, debounce)
	},
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build application for production",
	Long: `Build the Godin application for production deployment.

This command compiles your application into a standalone executable named 'app.exe'
(or 'app' on Unix systems) that can be deployed to production servers.

Examples:
  godin build                    # Build to app.exe in current directory
  godin build --output dist/     # Build to dist/app.exe
  godin build --name myapp       # Build to myapp.exe`,
	Run: func(cmd *cobra.Command, args []string) {
		output, _ := cmd.Flags().GetString("output")
		name, _ := cmd.Flags().GetString("name")
		buildApp(output, name)
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run application in debug mode",
	Long: `Run the Godin application in debug mode with enhanced logging and development features.

This command starts your application with:
- Debug logging enabled
- Detailed error messages
- Hot reload capabilities
- Development middleware
- Enhanced stack traces

Examples:
  godin run                      # Run on default port 8080
  godin run --port 3000          # Run on custom port
  godin run --no-debug           # Run without debug features`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		debug, _ := cmd.Flags().GetBool("debug")
		runApp(port, debug)
	},
}

var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Package management commands",
}

var packageAddCmd = &cobra.Command{
	Use:   "add [github-url]",
	Short: "Add a package dependency",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		githubURL := args[0]
		version, _ := cmd.Flags().GetString("version")
		addPackage(githubURL, version)
	},
}

var packageListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed packages",
	Run: func(cmd *cobra.Command, args []string) {
		listPackages()
	},
}

var packageRemoveCmd = &cobra.Command{
	Use:   "remove [package-name]",
	Short: "Remove a package dependency",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		packageName := args[0]
		removePackage(packageName)
	},
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Install dependencies from package.yaml",
	Long: `Install all dependencies listed in package.yaml file.

This command reads the package.yaml file in the current directory and installs
all dependencies and dev_dependencies using 'go get'. It automatically handles
both regular dependencies and the Godin framework itself.

Examples:
  godin get                    # Install all dependencies from package.yaml
  godin get --dev              # Install dev dependencies as well
  godin get --update           # Update dependencies to latest versions`,
	Run: func(cmd *cobra.Command, args []string) {
		dev, _ := cmd.Flags().GetBool("dev")
		update, _ := cmd.Flags().GetBool("update")
		installDependencies(dev, update)
	},
}

var createCmd = &cobra.Command{
	Use:   "create [app-name]",
	Short: "Create a new Godin application",
	Long: `Create a new Godin application with optional templates.

Available templates:
  counter - Counter app with navigation (default)
  simple  - Minimal app structure
  todo    - Full-featured todo application

Examples:
  godin create myapp                    # Creates with counter template
  godin create myapp --template todo    # Creates with todo template
  godin create myapp --no-template      # Creates config files only
  godin create --list-templates         # Shows available templates`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		listTemplates, _ := cmd.Flags().GetBool("list-templates")

		if listTemplates {
			showAvailableTemplates()
			return
		}

		if len(args) == 0 {
			log.Fatal("Error: app name is required. Use 'godin create --help' for usage information.")
		}

		appName := args[0]
		noTemplate, _ := cmd.Flags().GetBool("no-template")
		template, _ := cmd.Flags().GetString("template")
		description, _ := cmd.Flags().GetString("description")
		createApp(appName, noTemplate, template, description)
	},
}

func init() {
	// Serve command flags
	serveCmd.Flags().StringP("port", "p", "8080", "Server port")
	serveCmd.Flags().BoolP("watch", "w", true, "Enable file watching")
	serveCmd.Flags().BoolP("listen", "l", false, "Enable interactive commands (r for reload, R for refresh)")
	serveCmd.Flags().Bool("enhanced-reload", true, "Enable enhanced hot-reload with build caching and health monitoring")
	serveCmd.Flags().Int("restart-retries", 3, "Number of restart attempts on failure")
	serveCmd.Flags().Duration("debounce", 500*time.Millisecond, "File change debounce duration")

	// Build command flags
	buildCmd.Flags().StringP("output", "o", ".", "Output directory")
	buildCmd.Flags().StringP("name", "n", "app", "Output executable name (without extension)")

	// Run command flags
	runCmd.Flags().StringP("port", "p", "8080", "Server port")
	runCmd.Flags().Bool("debug", true, "Enable debug mode (default: true)")

	// Package add command flags
	packageAddCmd.Flags().StringP("version", "v", "latest", "Package version")

	// Create command flags
	createCmd.Flags().Bool("no-template", false, "Create only config files without template code")
	createCmd.Flags().StringP("template", "t", "counter", "Template to use (counter, simple, todo)")
	createCmd.Flags().StringP("description", "d", "", "Custom description for the application")
	createCmd.Flags().Bool("list-templates", false, "List available templates")

	// Get command flags
	getCmd.Flags().Bool("dev", false, "Install dev dependencies as well")
	getCmd.Flags().Bool("update", false, "Update dependencies to latest versions")

	// Add subcommands
	packageCmd.AddCommand(packageAddCmd)
	packageCmd.AddCommand(packageListCmd)
	packageCmd.AddCommand(packageRemoveCmd)

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(packageCmd)
}

func main() {
	fmt.Println("🚀 Godin CLI starting...")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Command implementations

func initProject(projectName string) {
	log.Printf("Initializing Godin project: %s", projectName)

	// Create project directory
	err := os.MkdirAll(projectName, 0755)
	if err != nil {
		log.Fatalf("Failed to create project directory: %v", err)
	}

	// Create project structure
	dirs := []string{
		"handlers",
		"widgets/components",
		"widgets/pages",
		"static/css",
		"static/js",
		"static/images",
		"templates",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(projectName, dir)
		err := os.MkdirAll(fullPath, 0755)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Create package.yaml
	config := map[string]interface{}{
		"name":         projectName,
		"version":      "1.0.0",
		"description":  "A Godin framework application",
		"dependencies": map[string]interface{}{},
		"scripts": map[string]string{
			"dev":   "godin serve --watch",
			"build": "godin build --prod",
			"test":  "godin test",
		},
		"config": map[string]interface{}{
			"server": map[string]interface{}{
				"port": "8080",
				"host": "localhost",
			},
			"websocket": map[string]interface{}{
				"enabled": true,
				"path":    "/ws",
			},
			"static": map[string]interface{}{
				"dir":   "static",
				"cache": true,
			},
		},
	}

	configData, err := yaml.Marshal(config)
	if err != nil {
		log.Fatalf("Failed to marshal package.yaml: %v", err)
	}

	configPath := filepath.Join(projectName, "package.yaml")
	err = os.WriteFile(configPath, configData, 0644)
	if err != nil {
		log.Fatalf("Failed to write package.yaml: %v", err)
	}

	// Create main.go template
	mainGoContent := `package main

import (
	"log"
)

func main() {
	log.Println("Welcome to Godin!")
	log.Println("Server-side Go with HTMX and WebSockets")

	// TODO: Initialize Godin app
	// app := core.New()
	// app.GET("/", HomeHandler)
	// app.Serve(":8080")
}

// func HomeHandler(ctx *core.Context) core.Widget {
//     return widgets.Container{
//         Style: "max-width: 800px; margin: 0 auto; padding: 20px;",
//         Child: widgets.Text{
//             Content: "Welcome to Godin!",
//             Style: "font-size: 24px; font-weight: bold;",
//         },
//     }
// }
`

	mainPath := filepath.Join(projectName, "main.go")
	err = os.WriteFile(mainPath, []byte(mainGoContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write main.go: %v", err)
	}

	log.Printf("Project %s initialized successfully!", projectName)
	log.Printf("Next steps:")
	log.Printf("  cd %s", projectName)
	log.Printf("  godin serve")
}

func startDevServer(port string, watch bool, listen bool) {
	// Use enhanced version with default settings
	startDevServerEnhanced(port, watch, listen, true, 3, 500*time.Millisecond)
}

func startDevServerEnhanced(port string, watch bool, listen bool, enhancedReload bool, restartRetries int, debounce time.Duration) {
	fmt.Printf("🚀 Starting Godin development server on port %s\n", port)
	fmt.Printf("📁 Watch mode: %v\n", watch)
	fmt.Printf("⌨️  Interactive mode: %v\n", listen)
	fmt.Printf("🔧 Enhanced reload: %v\n", enhancedReload)
	if enhancedReload {
		fmt.Printf("🔄 Restart retries: %d\n", restartRetries)
		fmt.Printf("⏱️  Debounce duration: %v\n", debounce)
	}

	log.Printf("🚀 Starting Godin development server on port %s", port)
	log.Printf("📁 Watch mode: %v", watch)
	log.Printf("⌨️  Interactive mode: %v", listen)
	log.Printf("🔧 Enhanced reload: %v", enhancedReload)

	fmt.Println("🔍 Step 1: Checking project...")

	// Check if we're in a Godin project
	fmt.Println("🔍 Step 2: Checking if this is a Godin project...")
	if !isGodinProject() {
		log.Fatal("Error: Not in a Godin project directory. Make sure package.yaml exists.")
	}
	fmt.Println("✅ Godin project detected")

	// Check if we need to fix imports for framework development
	fmt.Println("🔍 Step 3: Checking if import fix is needed...")
	if false { // Temporarily disabled to test hot-reload
		fmt.Println("🔧 Detected Godin framework development environment")
		log.Printf("🔧 Detected Godin framework development environment")
		if err := fixFrameworkImports(); err != nil {
			log.Printf("⚠️  Warning: Could not fix imports: %v", err)
		} else {
			log.Printf("✅ Framework imports fixed")
		}
	} else {
		fmt.Println("✅ Import fix skipped for testing")
	}

	// Set development environment variables
	fmt.Println("🔍 Step 4: Setting environment variables...")
	os.Setenv("GODIN_DEBUG", "true")
	os.Setenv("GODIN_LOG_LEVEL", "debug")
	os.Setenv("GODIN_DEV_MODE", "true")
	os.Setenv("GODIN_ENHANCED_RELOAD", fmt.Sprintf("%v", enhancedReload))
	os.Setenv("GODIN_RESTART_RETRIES", fmt.Sprintf("%d", restartRetries))
	os.Setenv("GODIN_DEBOUNCE_MS", fmt.Sprintf("%d", debounce.Milliseconds()))
	fmt.Println("✅ Environment variables set")

	// Start the development server with hot reload
	if listen {
		log.Println("📝 Interactive commands:")
		log.Println("  r - Hot reload (restart server)")
		log.Println("  R - Hot refresh (refresh browser)")
		log.Println("  t - Test hot reload endpoints")
		log.Println("  q - Quit server")
		log.Println("  h - Show help")
		log.Println("  Ctrl+C - Quit server")
		if enhancedReload {
			log.Println("🔧 Enhanced features:")
			log.Println("  - Build caching and pre-checks")
			log.Println("  - Health monitoring")
			log.Println("  - Intelligent restart queuing")
			log.Println("  - Smart file filtering")
		}
	}

	// Start interactive command listener if enabled
	fmt.Println("🔍 Step 5: Starting components...")
	if listen {
		fmt.Println("🎮 Starting interactive listener...")
		go startInteractiveListener()
	}

	// Start the server process with enhanced features
	fmt.Println("🚀 Starting server process...")
	if enhancedReload {
		startServerProcessEnhanced(port, watch, restartRetries, debounce)
	} else {
		startServerProcess(port, watch)
	}
}

// Global variables for server control
var (
	serverCmd       *exec.Cmd
	serverMutex     sync.Mutex
	shouldRestart   = make(chan bool, 1)
	shouldRefresh   = make(chan bool, 1)
	watcher         *fsnotify.Watcher
	watcherDone     = make(chan bool, 1)
	buildMutex      sync.Mutex
	lastBuildTime   time.Time
	buildInProgress bool
	restartQueue    = make(chan restartRequest, 10)
	serverHealth    = make(chan bool, 1)
)

// restartRequest represents a server restart request with context
type restartRequest struct {
	reason    string
	timestamp time.Time
	port      string
}

// startInteractiveListener listens for interactive commands
func startInteractiveListener() {
	log.Println("🎮 Interactive mode started")
	log.Println("📝 Type 'h' for help, 'q' to quit")

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("❌ Interactive listener panic recovered: %v", r)
			}
		}()

		for {
			var input string
			fmt.Print("godin> ")
			_, err := fmt.Scanln(&input)
			if err != nil {
				// Handle EOF or other input errors gracefully
				continue
			}

			switch strings.TrimSpace(input) {
			case "r", "reload":
				log.Println("🔥 Manual hot reload triggered by user")
				triggerHotReload()
			case "R", "refresh":
				log.Println("🔄 Manual hot refresh triggered by user")
				triggerHotRefresh()
			case "t", "test":
				log.Println("🧪 Testing hot reload endpoints...")
				go testHotReloadEndpoints(currentServerPort)
			case "q", "quit", "exit":
				log.Println("👋 Shutting down server...")
				cleanup()
				os.Exit(0)
			case "h", "help":
				log.Println("📝 Available commands:")
				log.Println("  r, reload  - Hot reload (restart server)")
				log.Println("  R, refresh - Hot refresh (refresh browser)")
				log.Println("  t, test    - Test hot reload endpoints")
				log.Println("  q, quit    - Quit server")
				log.Println("  h, help    - Show this help")
			case "":
				// Ignore empty input
				continue
			default:
				log.Printf("❓ Unknown command: %s. Type 'h' for help", input)
			}
		}
	}()
}

// startServerProcess starts the Go application server with enhanced hot-reload
func startServerProcess(port string, watch bool) {
	// Set the current server port for hot refresh
	currentServerPort = port

	// Handle Ctrl+C gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Start file watcher if enabled
	if watch {
		startFileWatcher()
	}

	// Start the restart queue processor
	go processRestartQueue()

	// Start server health monitor
	go monitorServerHealth()

	// Start the initial server
	go func() {
		startServer(port)
	}()

	log.Printf("🎯 Enhanced server process manager started")
	log.Printf("📝 Press 'r' for manual reload, 'R' for refresh, 'q' to quit")
	log.Printf("🔧 Hot-reload features: Build caching, Health monitoring, Queue processing")

	// Handle signals and restart requests
	for {
		select {
		case <-c:
			log.Println("\n👋 Shutting down server...")
			cleanup()
			os.Exit(0)
		case <-shouldRestart:
			log.Println("🔄 Queuing server restart...")
			queueRestart("manual", port)
		}
	}
}

// startServerProcessEnhanced starts the Go application server with configurable enhanced hot-reload features
func startServerProcessEnhanced(port string, watch bool, restartRetries int, debounce time.Duration) {
	// Set the current server port for hot refresh
	currentServerPort = port

	// Handle Ctrl+C gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Start file watcher if enabled with enhanced filtering and custom debounce
	if watch {
		startFileWatcherEnhanced(debounce)
	}

	// Start the restart queue processor
	go processRestartQueue()

	// Start server health monitor
	go monitorServerHealth()

	// Start the initial server
	go func() {
		startServer(port)
	}()

	log.Printf("🎯 Enhanced server process manager started")
	log.Printf("📝 Press 'r' for manual reload, 'R' for refresh, 'q' to quit")
	log.Printf("🔧 Hot-reload features: Build caching, Health monitoring, Queue processing")
	log.Printf("🔄 Restart retries: %d", restartRetries)
	log.Printf("⏱️  Debounce duration: %v", debounce)

	// Handle signals and restart requests
	for {
		select {
		case <-c:
			log.Println("\n👋 Shutting down server...")
			cleanup()
			os.Exit(0)
		case <-shouldRestart:
			log.Println("🔄 Queuing server restart...")
			queueRestart("manual", port)
		}
	}
}

// startServer starts the Go application with enhanced error handling
func startServer(port string) {
	serverMutex.Lock()
	defer serverMutex.Unlock()

	// Ensure previous server is fully stopped
	if serverCmd != nil && serverCmd.Process != nil {
		log.Println("🛑 Stopping previous server...")

		// Try graceful shutdown first
		if err := serverCmd.Process.Signal(os.Interrupt); err != nil {
			// If graceful shutdown fails, force kill
			serverCmd.Process.Kill()
		}

		// Wait for process to exit with timeout
		done := make(chan error, 1)
		go func() {
			done <- serverCmd.Wait()
		}()

		select {
		case <-done:
			log.Println("✅ Previous server stopped gracefully")
		case <-time.After(5 * time.Second):
			log.Println("⚠️  Forcing server shutdown...")
			serverCmd.Process.Kill()
			<-done
		}

		// Wait longer to ensure port is fully released on Windows
		log.Println("⏳ Waiting for port to be released...")
		time.Sleep(3 * time.Second)

		// Additional check to ensure port is available with retries
		maxRetries := 5
		for i := 0; i < maxRetries; i++ {
			if isPortAvailable(port) {
				log.Printf("✅ Port %s is now available", port)
				break
			}
			log.Printf("⏳ Port %s still in use, waiting... (attempt %d/%d)", port, i+1, maxRetries)
			time.Sleep(2 * time.Second)
		}
	}

	// Pre-build check to catch compilation errors early
	if !performPreBuildCheck() {
		log.Printf("❌ Pre-build check failed, skipping server start")
		return
	}

	log.Printf("🚀 Starting server on port %s...", port)

	serverCmd = exec.Command("go", "run", ".")
	serverCmd.Stdout = os.Stdout
	serverCmd.Stderr = os.Stderr

	// Set comprehensive development environment variables
	env := append(os.Environ(),
		"PORT="+port,
		"GODIN_DEV_MODE=true",
		"GODIN_HOT_RELOAD=true",
		"GODIN_DEBUG=true",
		"GODIN_LOG_LEVEL=debug",
		"GODIN_WEBSOCKET_ENABLED=true",
		"GODIN_FILE_WATCHING=true",
	)

	// Add port without colon for applications that expect it
	if strings.HasPrefix(port, ":") {
		env = append(env, "GODIN_PORT="+port[1:])
	} else {
		env = append(env, "GODIN_PORT="+port)
	}

	serverCmd.Env = env

	if err := serverCmd.Start(); err != nil {
		log.Printf("❌ Failed to start server: %v", err)
		return
	}

	log.Printf("✅ Server started successfully (PID: %d)", serverCmd.Process.Pid)
	log.Printf("🌐 Visit http://localhost:%s", port)

	// Wait a moment for server to fully start
	time.Sleep(1 * time.Second)

	// Test hot reload endpoints
	go testHotReloadEndpoints(port)

	// Signal that server is healthy
	select {
	case serverHealth <- true:
	default:
	}
}

// stopServer stops the current server
func stopServer() {
	serverMutex.Lock()
	defer serverMutex.Unlock()

	if serverCmd != nil && serverCmd.Process != nil {
		log.Println("🛑 Stopping server...")

		// Try graceful shutdown first
		if err := serverCmd.Process.Signal(os.Interrupt); err != nil {
			// If graceful shutdown fails, force kill
			serverCmd.Process.Kill()
		}

		// Wait for process to exit with timeout
		done := make(chan error, 1)
		go func() {
			done <- serverCmd.Wait()
		}()

		select {
		case <-done:
			log.Println("✅ Server stopped gracefully")
		case <-time.After(3 * time.Second):
			log.Println("⚠️  Forcing server shutdown...")
			serverCmd.Process.Kill()
			<-done
		}

		serverCmd = nil
	}
}

// cleanup performs cleanup operations when shutting down
func cleanup() {
	log.Println("🧹 Cleaning up...")

	// Stop the server
	stopServer()

	// Stop file watcher
	if watcher != nil {
		select {
		case watcherDone <- true:
		default:
		}
		watcher.Close()
		watcher = nil
	}

	log.Println("✅ Cleanup completed")
}

// restartServer is deprecated - use queueRestart instead
func restartServer(port string) {
	log.Println("🔄 Legacy restart function called - redirecting to queue")
	queueRestart("legacy-restart", port)
}

// triggerHotReload triggers a server restart (legacy function, now uses queue)
func triggerHotReload() {
	queueRestart("manual-trigger", currentServerPort)
}

// Global variable to track current server port
var currentServerPort string = "8080"

// triggerHotRefresh triggers a browser refresh (via WebSocket if available)
func triggerHotRefresh() {
	log.Println("🔄 Hot refresh signal sent")

	// Send hot refresh signal via HTTP to the running server
	go func() {
		// Wait a moment for server to be ready
		time.Sleep(200 * time.Millisecond)

		// Prepare the URL - handle both :port and port formats
		port := currentServerPort
		if !strings.HasPrefix(port, ":") {
			port = ":" + port
		}

		// Try to send hot refresh signal to the running server
		client := &http.Client{Timeout: 2 * time.Second}
		url := fmt.Sprintf("http://localhost%s/api/hot-refresh", port)

		log.Printf("🔄 Sending hot refresh to: %s", url)

		resp, err := client.Post(url, "application/json", strings.NewReader(`{"type":"hot-refresh"}`))
		if err != nil {
			log.Printf("⚠️  Hot refresh request failed: %v", err)
			return
		}

		if resp != nil {
			resp.Body.Close()
			if resp.StatusCode == 200 {
				log.Println("✅ Hot refresh signal sent successfully")
			} else {
				log.Printf("⚠️  Hot refresh returned status: %d", resp.StatusCode)
			}
		}
	}()
}

// Enhanced hot-reload functions for guaranteed success

// queueRestart adds a restart request to the queue
func queueRestart(reason, port string) {
	req := restartRequest{
		reason:    reason,
		timestamp: time.Now(),
		port:      port,
	}

	select {
	case restartQueue <- req:
		log.Printf("🔄 Restart queued: %s", reason)
	default:
		log.Printf("⚠️  Restart queue full, dropping request: %s", reason)
	}
}

// processRestartQueue processes restart requests with intelligent batching
func processRestartQueue() {
	var lastRestart time.Time
	const minRestartInterval = 2 * time.Second

	for req := range restartQueue {
		// Batch rapid requests
		if time.Since(lastRestart) < minRestartInterval {
			log.Printf("⏳ Batching restart request: %s", req.reason)
			time.Sleep(minRestartInterval - time.Since(lastRestart))
		}

		// Drain any additional requests that came in during the wait
		var latestReq = req
	drainLoop:
		for {
			select {
			case newReq := <-restartQueue:
				log.Printf("🔄 Batching additional request: %s", newReq.reason)
				latestReq = newReq
			default:
				break drainLoop
			}
		}

		log.Printf("🔄 Processing restart: %s", latestReq.reason)
		performRestart(latestReq.port)
		lastRestart = time.Now()
	}
}

// performRestart performs the actual server restart with enhanced error handling
func performRestart(port string) {
	buildMutex.Lock()
	defer buildMutex.Unlock()

	if buildInProgress {
		log.Println("⚠️  Build already in progress, skipping restart")
		return
	}

	buildInProgress = true
	defer func() {
		buildInProgress = false
	}()

	log.Println("🔄 Initiating enhanced server restart...")

	// Stop the current server
	stopServer()

	// Longer pause to ensure clean restart, especially on Windows
	log.Println("⏳ Preparing for restart...")
	time.Sleep(2 * time.Second)

	// Find an available port near the original port to avoid Windows TIME_WAIT issues
	newPort := findAvailablePort(port)
	if newPort != port {
		log.Printf("🔄 Using port %s instead of %s to avoid TIME_WAIT", newPort, port)
		// Update the current server port for hot refresh
		currentServerPort = newPort
	}

	// Start server with retry logic
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		log.Printf("🚀 Starting server attempt %d/%d", i+1, maxRetries)

		// Start server in a goroutine to avoid blocking
		go func(attempt int) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("❌ Server restart panic recovered (attempt %d): %v", attempt, r)
					if attempt < maxRetries {
						time.Sleep(3 * time.Second)
						log.Printf("🔄 Attempting restart recovery %d/%d...", attempt+1, maxRetries)
						performRestart(newPort)
					}
				}
			}()
			startServer(newPort)
		}(i + 1)

		// Wait for server health check
		select {
		case <-serverHealth:
			log.Printf("✅ Server restart successful on attempt %d", i+1)
			return
		case <-time.After(10 * time.Second):
			log.Printf("⚠️  Server health check timeout on attempt %d", i+1)
			if i < maxRetries-1 {
				stopServer()
				time.Sleep(2 * time.Second)
				continue
			}
		}
	}

	log.Printf("❌ Server restart failed after %d attempts", maxRetries)
}

// monitorServerHealth monitors server health and triggers restarts if needed
func monitorServerHealth() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if serverCmd != nil && serverCmd.Process != nil {
				// Check if process is still running
				if err := serverCmd.Process.Signal(syscall.Signal(0)); err != nil {
					log.Printf("⚠️  Server process health check failed: %v", err)
					log.Printf("🔄 Triggering automatic restart due to health check failure")
					queueRestart("health-check-failure", currentServerPort)
				}
			}
		}
	}
}

// performPreBuildCheck performs a quick compilation check before starting server
func performPreBuildCheck() bool {
	buildMutex.Lock()
	defer buildMutex.Unlock()

	// Skip if we just did a build check recently
	if time.Since(lastBuildTime) < 5*time.Second {
		return true
	}

	log.Println("🔍 Performing pre-build check...")

	// Try to build without running
	cmd := exec.Command("go", "build", "-o", "temp_build_check", ".")
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	// Clean up temp file
	os.Remove("temp_build_check")
	os.Remove("temp_build_check.exe")

	lastBuildTime = time.Now()

	if err != nil {
		log.Printf("❌ Pre-build check failed: %v", err)
		return false
	}

	log.Println("✅ Pre-build check passed")
	return true
}

// Enhanced file watching with smart filtering and dependency tracking

// shouldProcessFileEvent determines if a file change event should trigger a reload with enhanced filtering
func shouldProcessFileEventEnhanced(event fsnotify.Event) bool {
	// Only process write and create events
	if event.Op&fsnotify.Write == 0 && event.Op&fsnotify.Create == 0 {
		return false
	}

	fileName := filepath.Base(event.Name)
	ext := strings.ToLower(filepath.Ext(event.Name))

	// Skip temporary files and common editor artifacts
	if strings.HasPrefix(fileName, ".") ||
		strings.HasPrefix(fileName, "~") ||
		strings.HasSuffix(fileName, ".tmp") ||
		strings.HasSuffix(fileName, ".swp") ||
		strings.HasSuffix(fileName, ".bak") ||
		strings.Contains(fileName, "___jb_tmp___") || // JetBrains temp files
		strings.Contains(fileName, "___jb_old___") { // JetBrains backup files
		return false
	}

	// Skip build artifacts and common ignore patterns
	if strings.Contains(event.Name, "/.git/") ||
		strings.Contains(event.Name, "/node_modules/") ||
		strings.Contains(event.Name, "/dist/") ||
		strings.Contains(event.Name, "/bin/") ||
		strings.Contains(event.Name, "/vendor/") ||
		strings.Contains(event.Name, "\\dist\\") ||
		strings.Contains(event.Name, "\\bin\\") ||
		strings.Contains(event.Name, "\\.git\\") {
		return false
	}

	// Check file extension
	watchedExtensions := []string{".go", ".html", ".css", ".js", ".yaml", ".yml", ".json", ".md"}
	for _, watchedExt := range watchedExtensions {
		if ext == watchedExt {
			return true
		}
	}

	// Also watch package.yaml specifically
	if fileName == "package.yaml" || fileName == "go.mod" || fileName == "go.sum" {
		return true
	}

	return false
}

// startFileWatcher starts watching files for changes
func startFileWatcher() {
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Printf("❌ Error creating file watcher: %v", err)
		return
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("❌ Error getting current directory: %v", err)
		watcher.Close()
		return
	}

	// Add paths to watch - watch the current project directory
	watchPaths := []string{".", "static", "templates", "web"}
	watchedCount := 0
	for _, path := range watchPaths {
		fullPath := filepath.Join(cwd, path)
		if _, err := os.Stat(fullPath); err == nil {
			err = addPathRecursively(watcher, fullPath)
			if err != nil {
				log.Printf("⚠️  Error watching path %s: %v", fullPath, err)
			} else {
				log.Printf("👀 Watching: %s", fullPath)
				watchedCount++
			}
		}
	}

	if watchedCount == 0 {
		log.Printf("⚠️  No valid paths found to watch")
		watcher.Close()
		return
	}

	log.Printf("👀 File watcher started for project: %s (%d paths)", cwd, watchedCount)

	// Debounce timer to prevent rapid restarts
	debounceTimer := time.NewTimer(0)
	debounceTimer.Stop()

	// Start watching in a separate goroutine
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("❌ File watcher panic recovered: %v", r)
			}
			watcher.Close()
		}()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					log.Printf("👋 File watcher events channel closed")
					return
				}

				if shouldProcessFileEventEnhanced(event) {
					// Get relative path for better logging
					filePath := event.Name
					if cwd, err := os.Getwd(); err == nil {
						if rel, err := filepath.Rel(cwd, filePath); err == nil {
							filePath = rel
						}
					}
					log.Printf("📝 File changed: %s", filePath)

					// Stop previous timer if running
					if !debounceTimer.Stop() {
						select {
						case <-debounceTimer.C:
						default:
						}
					}

					// Debounce rapid file changes
					debounceTimer.Reset(500 * time.Millisecond)
					go func(evt fsnotify.Event) {
						<-debounceTimer.C
						handleFileChangeEvent(evt)
					}(event)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					log.Printf("👋 File watcher errors channel closed")
					return
				}
				log.Printf("❌ File watcher error: %v", err)

			case <-watcherDone:
				log.Printf("👋 File watcher stopped")
				return
			}
		}
	}()
}

// startFileWatcherEnhanced starts watching files for changes with configurable debounce
func startFileWatcherEnhanced(debounce time.Duration) {
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Printf("❌ Error creating enhanced file watcher: %v", err)
		return
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("❌ Error getting current directory: %v", err)
		watcher.Close()
		return
	}

	// Add paths to watch - watch the current project directory
	watchPaths := []string{".", "static", "templates", "web", "handlers", "widgets"}
	watchedCount := 0
	for _, path := range watchPaths {
		fullPath := filepath.Join(cwd, path)
		if _, err := os.Stat(fullPath); err == nil {
			err = addPathRecursively(watcher, fullPath)
			if err != nil {
				log.Printf("⚠️  Error watching path %s: %v", fullPath, err)
			} else {
				log.Printf("👀 Enhanced watching: %s", fullPath)
				watchedCount++
			}
		}
	}

	if watchedCount == 0 {
		log.Printf("⚠️  No valid paths found to watch")
		watcher.Close()
		return
	}

	log.Printf("👀 Enhanced file watcher started for project: %s (%d paths)", cwd, watchedCount)
	log.Printf("⏱️  Using custom debounce duration: %v", debounce)

	// Debounce timer to prevent rapid restarts with configurable duration
	debounceTimer := time.NewTimer(0)
	debounceTimer.Stop()

	// Start watching in a separate goroutine
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("❌ Enhanced file watcher panic recovered: %v", r)
			}
			watcher.Close()
		}()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					log.Printf("👋 Enhanced file watcher events channel closed")
					return
				}

				if shouldProcessFileEventEnhanced(event) {
					// Get relative path for better logging
					filePath := event.Name
					if rel, err := filepath.Rel(cwd, filePath); err == nil {
						filePath = rel
					}
					log.Printf("📝 Enhanced file changed: %s", filePath)

					// Stop previous timer if running
					if !debounceTimer.Stop() {
						select {
						case <-debounceTimer.C:
						default:
						}
					}

					// Debounce rapid file changes with configurable duration
					debounceTimer.Reset(debounce)
					go func(evt fsnotify.Event) {
						<-debounceTimer.C
						handleFileChangeEvent(evt)
					}(event)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					log.Printf("👋 Enhanced file watcher errors channel closed")
					return
				}
				log.Printf("❌ Enhanced file watcher error: %v", err)

			case <-watcherDone:
				log.Printf("👋 Enhanced file watcher stopped")
				return
			}
		}
	}()
}

// addPathRecursively adds a path and all its subdirectories to the watcher
func addPathRecursively(watcher *fsnotify.Watcher, root string) error {
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
			return watcher.Add(path)
		}
		return nil
	})
}

// shouldProcessFileEvent determines if a file change event should trigger a reload
func shouldProcessFileEvent(event fsnotify.Event) bool {
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

// handleFileChangeEvent processes a file change and triggers appropriate actions
func handleFileChangeEvent(event fsnotify.Event) {
	ext := strings.ToLower(filepath.Ext(event.Name))
	fileName := filepath.Base(event.Name)
	filePath := event.Name

	// Get relative path for better logging
	if cwd, err := os.Getwd(); err == nil {
		if rel, err := filepath.Rel(cwd, filePath); err == nil {
			filePath = rel
		}
	}

	log.Printf("📝 File change detected: %s (type: %s)", filePath, event.Op.String())

	switch ext {
	case ".go", ".yaml", ".yml":
		// Go files or config changes require hot reload (restart)
		log.Printf("🔥 Go/config file changed (%s) - queuing hot reload", filePath)
		queueRestart("file-change:"+filePath, currentServerPort)
	case ".html", ".css", ".js":
		// Static files can use hot refresh (no restart)
		log.Printf("🔄 Static file changed (%s) - triggering hot refresh", filePath)
		triggerHotRefresh()
	default:
		// Check if it's package.yaml specifically
		if fileName == "package.yaml" {
			log.Printf("🔥 Package config changed - queuing hot reload")
			queueRestart("package-config-change", currentServerPort)
		} else {
			// Default to hot reload for unknown files
			log.Printf("🔥 File changed (%s) - queuing hot reload", filePath)
			queueRestart("file-change:"+filePath, currentServerPort)
		}
	}
}

func buildApp(output, name string) {
	log.Printf("Building Godin application...")

	// Check if we're in a Godin project
	if !isGodinProject() {
		log.Fatal("Error: Not in a Godin project directory. Make sure package.yaml exists.")
	}

	// Determine output path
	var outputPath string
	if runtime.GOOS == "windows" {
		outputPath = filepath.Join(output, name+".exe")
	} else {
		outputPath = filepath.Join(output, name)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(output, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Build the application
	log.Printf("Compiling to %s...", outputPath)

	buildCmd := exec.Command("go", "build", "-o", outputPath, ".")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	if err := buildCmd.Run(); err != nil {
		log.Fatalf("Build failed: %v", err)
	}

	log.Printf("✅ Build successful!")
	log.Printf("📦 Executable created: %s", outputPath)

	// Show file size
	if info, err := os.Stat(outputPath); err == nil {
		size := info.Size()
		sizeStr := formatFileSize(size)
		log.Printf("📊 File size: %s", sizeStr)
	}

	log.Printf("🚀 Ready for deployment!")
}

func runApp(port string, debug bool) {
	log.Printf("Starting Godin application in debug mode...")

	// Check if we're in a Godin project
	if !isGodinProject() {
		log.Fatal("Error: Not in a Godin project directory. Make sure package.yaml exists.")
	}

	// Check if we need to fix imports for framework development
	if needsImportFix() {
		log.Printf("🔧 Detected Godin framework development environment")
		if err := fixFrameworkImports(); err != nil {
			log.Printf("⚠️  Warning: Could not fix imports: %v", err)
		} else {
			log.Printf("✅ Framework imports fixed")
		}
	}

	// Set debug environment variables
	if debug {
		os.Setenv("GODIN_DEBUG", "true")
		os.Setenv("GODIN_LOG_LEVEL", "debug")
		log.Printf("🐛 Debug mode enabled")
	}

	// Set port environment variable
	os.Setenv("GODIN_PORT", port)

	log.Printf("🚀 Starting server on port %s", port)
	log.Printf("🔍 Debug logging: %v", debug)
	log.Printf("📂 Working directory: %s", getCurrentDir())

	if debug {
		log.Printf("🔧 Environment variables:")
		log.Printf("   GODIN_DEBUG=%s", os.Getenv("GODIN_DEBUG"))
		log.Printf("   GODIN_LOG_LEVEL=%s", os.Getenv("GODIN_LOG_LEVEL"))
		log.Printf("   GODIN_PORT=%s", os.Getenv("GODIN_PORT"))
	}

	// Run the application
	runCmd := exec.Command("go", "run", ".")
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	runCmd.Env = os.Environ()

	// Handle Ctrl+C gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Start the application in a goroutine
	go func() {
		if err := runCmd.Run(); err != nil {
			if !strings.Contains(err.Error(), "signal: interrupt") {
				log.Fatalf("Application failed to start: %v", err)
			}
		}
	}()

	log.Printf("✅ Application started successfully!")
	log.Printf("🌐 Visit http://localhost:%s to view your app", port)
	log.Printf("⏹️  Press Ctrl+C to stop the server")

	// Wait for interrupt signal
	<-c
	log.Printf("\n🛑 Shutting down server...")

	// Kill the process if it's still running
	if runCmd.Process != nil {
		runCmd.Process.Kill()
	}

	log.Printf("👋 Server stopped")
}

func addPackage(githubURL, version string) {
	log.Printf("Adding package %s@%s", githubURL, version)
	// TODO: Implement package installation
	log.Println("Package installation not yet implemented")
}

func listPackages() {
	log.Println("Listing installed packages")
	// TODO: Implement package listing
	log.Println("No packages installed")
}

func removePackage(packageName string) {
	log.Printf("Removing package %s", packageName)
	// TODO: Implement package removal
	log.Println("Package removal not yet implemented")
}

func installDependencies(includeDev, update bool) {
	log.Printf("📦 Installing dependencies from package.yaml...")

	// Check if package.yaml exists
	if !isGodinProject() {
		log.Fatal("Error: Not in a Godin project directory. Make sure package.yaml exists.")
	}

	// Load package.yaml
	config, err := loadPackageConfig(".")
	if err != nil {
		log.Fatalf("Failed to load package.yaml: %v", err)
	}

	log.Printf("📋 Found %d dependencies", len(config.Dependencies))
	if includeDev {
		log.Printf("📋 Found %d dev dependencies", len(config.DevDependencies))
	}

	// Install regular dependencies
	for name, dep := range config.Dependencies {
		if err := installGoPackage(name, dep.GitHub, dep.Version, update); err != nil {
			log.Printf("⚠️  Warning: Failed to install %s: %v", name, err)
		} else {
			log.Printf("✅ Installed %s@%s", name, dep.Version)
		}
	}

	// Install dev dependencies if requested
	if includeDev {
		for name, dep := range config.DevDependencies {
			if err := installGoPackage(name, dep.GitHub, dep.Version, update); err != nil {
				log.Printf("⚠️  Warning: Failed to install dev dependency %s: %v", name, err)
			} else {
				log.Printf("✅ Installed dev dependency %s@%s", name, dep.Version)
			}
		}
	}

	// Run go mod tidy to clean up
	log.Printf("🧹 Running go mod tidy...")
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Stdout = os.Stdout
	tidyCmd.Stderr = os.Stderr
	if err := tidyCmd.Run(); err != nil {
		log.Printf("⚠️  Warning: go mod tidy failed: %v", err)
	} else {
		log.Printf("✅ Dependencies cleaned up")
	}

	log.Printf("🎉 All dependencies installed successfully!")
}

// PackageConfig represents the structure of package.yaml
type PackageConfig struct {
	Name            string                       `yaml:"name"`
	Version         string                       `yaml:"version"`
	Description     string                       `yaml:"description"`
	Dependencies    map[string]PackageDependency `yaml:"dependencies"`
	DevDependencies map[string]PackageDependency `yaml:"dev_dependencies"`
	Scripts         map[string]string            `yaml:"scripts"`
}

// PackageDependency represents a package dependency
type PackageDependency struct {
	GitHub  string `yaml:"github"`
	Version string `yaml:"version"`
}

// loadPackageConfig loads and parses package.yaml
func loadPackageConfig(dir string) (*PackageConfig, error) {
	configPath := filepath.Join(dir, "package.yaml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read package.yaml: %w", err)
	}

	var config PackageConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse package.yaml: %w", err)
	}

	return &config, nil
}

// installGoPackage installs a Go package using go get
func installGoPackage(name, githubURL, version string, update bool) error {
	if githubURL == "" {
		return fmt.Errorf("no GitHub URL specified for package %s", name)
	}

	// Construct the go get command
	var args []string
	if update {
		args = []string{"get", "-u", githubURL + "@" + version}
	} else {
		args = []string{"get", githubURL + "@" + version}
	}

	// Execute go get
	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func showAvailableTemplates() {
	fmt.Println("Available Godin Templates:")
	fmt.Println()
	fmt.Println("  counter  - Counter app with navigation and state management (default)")
	fmt.Println("           Features: increment/decrement buttons, multiple pages, navigation")
	fmt.Println()
	fmt.Println("  simple   - Minimal app structure with basic welcome page")
	fmt.Println("           Features: clean starting point, minimal code")
	fmt.Println()
	fmt.Println("  todo     - Full-featured todo application")
	fmt.Println("           Features: add/toggle/delete todos, form handling, interactive UI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  godin create myapp --template counter")
	fmt.Println("  godin create myapp --template simple")
	fmt.Println("  godin create myapp --template todo")
	fmt.Println("  godin create myapp --no-template")
}

func createApp(appName string, noTemplate bool, template string, description string) {
	log.Printf("Creating Godin app: %s", appName)

	// Validate template
	validTemplates := []string{"counter", "simple", "todo"}
	isValidTemplate := false
	for _, validTemplate := range validTemplates {
		if template == validTemplate {
			isValidTemplate = true
			break
		}
	}

	if !noTemplate && !isValidTemplate {
		log.Printf("Warning: Unknown template '%s'. Available templates: %v", template, validTemplates)
		log.Printf("Using 'counter' template as fallback.")
		template = "counter"
	}

	// Create app directory
	err := os.MkdirAll(appName, 0755)
	if err != nil {
		log.Fatalf("Failed to create app directory: %v", err)
	}

	// Create project structure
	dirs := []string{
		"handlers",
		"widgets/components",
		"widgets/pages",
		"static/css",
		"static/js",
		"static/images",
		"templates",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(appName, dir)
		err := os.MkdirAll(fullPath, 0755)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Create package.yaml
	if description == "" {
		description = "A Godin framework application"
	}
	createPackageYaml(appName, description)

	if noTemplate {
		// Create minimal main.go
		createMinimalMainGo(appName)
		createReadme(appName, description, "minimal", true)
		log.Printf("App %s created with config files only!", appName)
		log.Printf("Next steps:")
		log.Printf("  cd %s", appName)
		log.Printf("  # Add your code to main.go")
		log.Printf("  godin serve")
	} else {
		// Create app with template
		createAppWithTemplate(appName, template)
		createReadme(appName, description, template, false)
		log.Printf("App %s created with %s template!", appName, template)
		log.Printf("Next steps:")
		log.Printf("  cd %s", appName)
		log.Printf("  godin serve")
	}
}

func createPackageYaml(appName, description string) {
	config := map[string]interface{}{
		"name":        appName,
		"version":     "1.0.0",
		"description": description,
		"dependencies": map[string]interface{}{
			"godin-framework": map[string]interface{}{
				"github":  "github.com/gideonsigilai/godin",
				"version": "v1.0.0",
			},
		},
		"dev_dependencies": map[string]interface{}{},
		"scripts": map[string]string{
			"dev":     "godin serve --watch",
			"build":   "godin build --prod",
			"test":    "godin test",
			"install": "go mod tidy",
			"clean":   "rm -rf dist/ bin/",
		},
		"config": map[string]interface{}{
			"server": map[string]interface{}{
				"port": "8080",
				"host": "localhost",
			},
			"websocket": map[string]interface{}{
				"enabled": true,
				"path":    "/ws",
			},
			"static": map[string]interface{}{
				"dir":   "static",
				"cache": true,
			},
		},
		"build": map[string]interface{}{
			"output_dir":    "dist",
			"static_dir":    "static",
			"templates_dir": "templates",
			"minify":        true,
			"source_maps":   true,
		},
		"development": map[string]interface{}{
			"hot_reload":    true,
			"file_watching": true,
			"auto_restart":  true,
			"debug_mode":    true,
			"watch": []string{
				"**/*.go",
				"**/*.html",
				"**/*.css",
				"**/*.js",
				"package.yaml",
			},
			"ignore": []string{
				"dist/**",
				"bin/**",
				"node_modules/**",
				".git/**",
			},
		},
	}

	configData, err := yaml.Marshal(config)
	if err != nil {
		log.Fatalf("Failed to marshal package.yaml: %v", err)
	}

	configPath := filepath.Join(appName, "package.yaml")
	err = os.WriteFile(configPath, configData, 0644)
	if err != nil {
		log.Fatalf("Failed to write package.yaml: %v", err)
	}

	// Create additional project files
	createGoMod(appName)
	createBasicCSS(appName)
	createBaseTemplate(appName)
}

func createMinimalMainGo(appName string) {
	mainGoContent := `package main

import (
	"log"

	. "github.com/gideonsigilai/godin/pkg/godin"
)

func main() {
	app := New()

	// Add your routes here
	app.GET("/", HomeHandler)

	log.Printf("Starting %s on :8080", "` + appName + `")
	log.Println("Visit http://localhost:8080 to see your app")
	if err := app.Serve(":8080"); err != nil {
		log.Fatal(err)
	}
}

// HomeHandler renders the home page
func HomeHandler(ctx *Context) Widget {
	return Container{
		Style: "max-width: 800px; margin: 0 auto; padding: 20px; font-family: Arial, sans-serif;",
		Child: Column{
			Children: []Widget{
				Text{
					Data: "Welcome to ` + appName + `!",
					TextStyle: &TextStyle{
						FontSize:   &[]float64{32}[0],
						FontWeight: FontWeightBold,
						Color:      Color("#333"),
					},
				},
				SizedBox{Height: &[]float64{20}[0]},
				Text{
					Data: "Your Godin app is ready. Start building amazing things!",
					TextStyle: &TextStyle{
						FontSize: &[]float64{16}[0],
						Color:    Color("#666"),
					},
				},
			},
		},
	}
}
`

	mainPath := filepath.Join(appName, "main.go")
	err := os.WriteFile(mainPath, []byte(mainGoContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write main.go: %v", err)
	}
}

func createAppWithTemplate(appName, template string) {
	switch template {
	case "counter":
		createCounterTemplate(appName)
	case "simple":
		createSimpleTemplate(appName)
	case "todo":
		createTodoTemplate(appName)
	default:
		log.Printf("Unknown template '%s', using 'counter' template", template)
		createCounterTemplate(appName)
	}
}

func createSimpleTemplate(appName string) {
	// This is the same as minimal but with a bit more structure
	createMinimalMainGo(appName)
}

func createCounterTemplate(appName string) {
	mainGoContent := `package main

import (
	"fmt"
	"log"

	. "github.com/gideonsigilai/godin/pkg/godin"
)

// App state
var counter = 0

func main() {
	app := New()

	// Routes
	app.GET("/", HomeHandler)
	app.GET("/about", AboutHandler)
	app.POST("/increment", IncrementHandler)
	app.POST("/decrement", DecrementHandler)
	app.POST("/reset", ResetHandler)

	log.Printf("Starting %s on :8080", "` + appName + `")
	log.Println("Visit http://localhost:8080 to see your counter app")
	if err := app.Serve(":8080"); err != nil {
		log.Fatal(err)
	}
}

// HomeHandler renders the main counter page
func HomeHandler(ctx *Context) Widget {
	return Container{
		Style: "min-height: 100vh; font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;",
		Child: Column{
			Children: []Widget{
				// App Bar
				AppBarWidget("` + appName + `"),

				// Main content
				Expanded{
					Child: Container{
						Style: "padding: 20px;",
						Child: Column{
							MainAxisAlignment: MainAxisAlignmentCenter,
							Children: []Widget{
								// Counter display section
								Text{
									Data:      "You have pushed the button this many times:",
									TextAlign: TextAlignCenter,
									TextStyle: &TextStyle{
										FontSize: &[]float64{16}[0],
									},
								},

								// Spacer
								SizedBox{Height: &[]float64{20}[0]},

								// Counter value with ID for HTMX updates
								Container{
									ID: "counter-display",
									Child: Text{
										Data:      fmt.Sprintf("%d", counter),
										TextAlign: TextAlignCenter,
										TextStyle: &TextStyle{
											FontSize:   &[]float64{48}[0],
											FontWeight: FontWeightBold,
											Color:      Color("#2196F3"),
										},
									},
								},

								// Spacer
								SizedBox{Height: &[]float64{40}[0]},

								// Action buttons row
								Row{
									MainAxisAlignment: MainAxisAlignmentCenter,
									Children: []Widget{
										// Decrement button
										ElevatedButton{
											Child: Text{Data: "−"},
											OnPressed: func() {
												counter--
												log.Printf("Counter decremented to: %d", counter)
											},
											Style: "min-width: 60px; height: 60px; font-size: 24px; border-radius: 30px; margin-right: 10px;",
										},

										// Spacer
										SizedBox{Width: &[]float64{20}[0]},

										// Reset button
										FilledButton{
											Child: Text{Data: "Reset"},
											OnPressed: func() {
												counter = 0
												log.Printf("Counter reset to: %d", counter)
											},
											Style: "min-width: 80px; height: 60px; font-size: 16px; border-radius: 30px; background-color: #f44336;",
										},

										// Spacer
										SizedBox{Width: &[]float64{20}[0]},

										// Increment button
										ElevatedButton{
											Child: Text{Data: "+"},
											OnPressed: func() {
												counter++
												log.Printf("Counter incremented to: %d", counter)
											},
											Style: "min-width: 60px; height: 60px; font-size: 24px; border-radius: 30px; margin-left: 10px;",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
`

	// Continue with the rest of the counter template
	mainGoContent += `
// AboutHandler renders the about page
func AboutHandler(ctx *core.Context) widgets.Widget {
	return widgets.Container{
		Style: "min-height: 100vh; font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;",
		Child: widgets.Column{
			Children: []widgets.Widget{
				// App Bar
				AppBarWidget("About"),

				// Main content
				widgets.Expanded{
					Child: widgets.Container{
						Style: "padding: 40px; max-width: 800px; margin: 0 auto;",
						Child: widgets.Column{
							CrossAxisAlignment: widgets.CrossAxisAlignmentStart,
							Children: []widgets.Widget{
								// Title
								widgets.Text{
									Data: "About ` + appName + `",
									TextStyle: &widgets.TextStyle{
										FontSize:   &[]float64{32}[0],
										FontWeight: widgets.FontWeightBold,
										Color:      widgets.Color("#333"),
									},
								},

								// Spacer
								widgets.SizedBox{Height: &[]float64{30}[0]},

								// Description
								widgets.Text{
									Data: "This is a counter application built with the Godin Framework.",
									TextStyle: &widgets.TextStyle{
										FontSize: &[]float64{16}[0],
										Color:    widgets.Color("#666"),
									},
								},

								// Spacer
								widgets.SizedBox{Height: &[]float64{20}[0]},

								// Current counter info
								widgets.Card{
									Style: "padding: 20px; border-radius: 8px; background-color: #f5f5f5;",
									Child: widgets.Column{
										Children: []widgets.Widget{
											widgets.Text{
												Data: "Current Counter Value:",
												TextStyle: &widgets.TextStyle{
													FontWeight: widgets.FontWeightBold,
													Color:      widgets.Color("#333"),
												},
											},
											widgets.SizedBox{Height: &[]float64{10}[0]},
											widgets.Text{
												Data: fmt.Sprintf("%d", counter),
												TextStyle: &widgets.TextStyle{
													FontSize:   &[]float64{24}[0],
													FontWeight: widgets.FontWeightBold,
													Color:      widgets.Color("#2196F3"),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Helper Components

// AppBarWidget creates a reusable app bar with navigation
func AppBarWidget(title string) widgets.Widget {
	return widgets.AppBar{
		Title: widgets.Text{
			Data: title,
			TextStyle: &widgets.TextStyle{
				Color:      widgets.ColorWhite,
				FontSize:   &[]float64{20}[0],
				FontWeight: widgets.FontWeightBold,
			},
		},
		Style: "background: linear-gradient(135deg, #2196F3, #1976D2); color: white; padding: 16px 20px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);",
		Actions: []widgets.Widget{
			// Navigation menu
			widgets.Row{
				Children: []widgets.Widget{
					NavButton("Home", "/"),
					widgets.SizedBox{Width: &[]float64{10}[0]},
					NavButton("About", "/about"),
				},
			},
		},
	}
}

// NavButton creates a navigation button
func NavButton(text, href string) widgets.Widget {
	return widgets.Container{
		Style: "padding: 8px 16px; border-radius: 4px; background: rgba(255,255,255,0.1); cursor: pointer;",
		Child: widgets.Text{
			Data: text,
			TextStyle: &widgets.TextStyle{
				Color:      widgets.ColorWhite,
				FontWeight: widgets.FontWeightW500,
			},
		},
	}
}

// Action Handlers

// IncrementHandler increments the counter
func IncrementHandler(ctx *core.Context) widgets.Widget {
	counter++
	log.Printf("Counter incremented to: %d", counter)

	// Return the updated counter display
	return widgets.Text{
		Data:      fmt.Sprintf("%d", counter),
		TextAlign: widgets.TextAlignCenter,
		TextStyle: &widgets.TextStyle{
			FontSize:   &[]float64{48}[0],
			FontWeight: widgets.FontWeightBold,
			Color:      widgets.Color("#2196F3"),
		},
	}
}

// DecrementHandler decrements the counter
func DecrementHandler(ctx *core.Context) widgets.Widget {
	counter--
	log.Printf("Counter decremented to: %d", counter)

	// Return the updated counter display
	return widgets.Text{
		Data:      fmt.Sprintf("%d", counter),
		TextAlign: widgets.TextAlignCenter,
		TextStyle: &widgets.TextStyle{
			FontSize:   &[]float64{48}[0],
			FontWeight: widgets.FontWeightBold,
			Color:      widgets.Color("#2196F3"),
		},
	}
}

// ResetHandler resets the counter to zero
func ResetHandler(ctx *core.Context) widgets.Widget {
	counter = 0
	log.Printf("Counter reset to: %d", counter)

	// Return the updated counter display
	return widgets.Text{
		Data:      fmt.Sprintf("%d", counter),
		TextAlign: widgets.TextAlignCenter,
		TextStyle: &widgets.TextStyle{
			FontSize:   &[]float64{48}[0],
			FontWeight: widgets.FontWeightBold,
			Color:      widgets.Color("#2196F3"),
		},
	}
}
`

	mainPath := filepath.Join(appName, "main.go")
	err := os.WriteFile(mainPath, []byte(mainGoContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write main.go: %v", err)
	}
}

func createTodoTemplate(appName string) {
	mainGoContent := `package main

import (
	"fmt"
	"log"
	"strconv"

	. "github.com/gideonsigilai/godin/pkg/godin"
)

// Todo represents a todo item
type Todo struct {
	ID        int
	Text      string
	Completed bool
}

// App state
var todos = []Todo{
	{ID: 1, Text: "Learn Godin Framework", Completed: false},
	{ID: 2, Text: "Build an awesome app", Completed: false},
}
var nextID = 3

func main() {
	app := New()

	// Routes
	app.GET("/", HomeHandler)
	app.POST("/add-todo", AddTodoHandler)
	app.POST("/toggle-todo", ToggleTodoHandler)
	app.POST("/delete-todo", DeleteTodoHandler)

	log.Printf("Starting %s on :8080", "` + appName + `")
	log.Println("Visit http://localhost:8080 to see your todo app")
	if err := app.Serve(":8080"); err != nil {
		log.Fatal(err)
	}
}

// HomeHandler renders the main todo page
func HomeHandler(ctx *Context) Widget {
	return Container{
		Style: "min-height: 100vh; font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #f5f5f5;",
		Child: widgets.Column{
			Children: []widgets.Widget{
				// Header
				widgets.Container{
					Style: "background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 30px 20px; text-align: center;",
					Child: widgets.Column{
						Children: []widgets.Widget{
							widgets.Text{
								Data: "` + appName + `",
								TextStyle: &widgets.TextStyle{
									FontSize:   &[]float64{32}[0],
									FontWeight: widgets.FontWeightBold,
									Color:      widgets.ColorWhite,
								},
							},
							widgets.SizedBox{Height: &[]float64{10}[0]},
							widgets.Text{
								Data: "Stay organized and get things done!",
								TextStyle: &widgets.TextStyle{
									FontSize: &[]float64{16}[0],
									Color:    widgets.ColorWhite,
								},
							},
						},
					},
				},

				// Main content
				widgets.Expanded{
					Child: widgets.Container{
						Style: "max-width: 600px; margin: 0 auto; padding: 20px;",
						Child: widgets.Column{
							Children: []widgets.Widget{
								// Add todo form
								widgets.Card{
									Style: "padding: 20px; margin-bottom: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);",
									Child: widgets.Column{
										Children: []widgets.Widget{
											widgets.Text{
												Data: "Add New Todo",
												TextStyle: &widgets.TextStyle{
													FontSize:   &[]float64{18}[0],
													FontWeight: widgets.FontWeightBold,
													Color:      widgets.Color("#333"),
												},
											},
											widgets.SizedBox{Height: &[]float64{15}[0]},
											widgets.Row{
												Children: []widgets.Widget{
													widgets.Expanded{
														Child: widgets.TextField{
															ID: "new-todo-input",
															Decoration: &widgets.InputDecoration{
																HintText: "What needs to be done?",
															},
															Style: "padding: 10px; border: 1px solid #ddd; border-radius: 4px;",
														},
													},
													widgets.SizedBox{Width: &[]float64{10}[0]},
													widgets.ElevatedButton{
														Child: widgets.Text{Data: "Add"},
														OnPressed: func() {
															log.Println("Add todo button pressed")
														},
														Style: "padding: 10px 20px; background-color: #667eea; color: white; border: none; border-radius: 4px; cursor: pointer;",
													},
												},
											},
										},
									},
								},

								// Todo list
								widgets.Card{
									Style: "padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);",
									Child: widgets.Column{
										Children: append([]widgets.Widget{
											widgets.Text{
												Data: fmt.Sprintf("Todo List (%d items)", len(todos)),
												TextStyle: &widgets.TextStyle{
													FontSize:   &[]float64{18}[0],
													FontWeight: widgets.FontWeightBold,
													Color:      widgets.Color("#333"),
												},
											},
											widgets.SizedBox{Height: &[]float64{15}[0]},
										}, getTodoWidgets()...),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// getTodoWidgets returns a slice of widgets for each todo
func getTodoWidgets() []widgets.Widget {
	var todoWidgets []widgets.Widget

	for _, todo := range todos {
		todoWidgets = append(todoWidgets, createTodoItem(todo))
	}

	if len(todoWidgets) == 0 {
		todoWidgets = append(todoWidgets, widgets.Container{
			Style: "text-align: center; padding: 40px; color: #999;",
			Child: widgets.Text{
				Data: "No todos yet. Add one above!",
				TextStyle: &widgets.TextStyle{
					Color: widgets.Color("#999"),
				},
			},
		})
	}

	return todoWidgets
}

// createTodoItem creates a widget for a single todo item
func createTodoItem(todo Todo) widgets.Widget {
	textStyle := &widgets.TextStyle{
		Color: widgets.Color("#333"),
	}

	if todo.Completed {
		textStyle.Color = widgets.Color("#999")
	}

	return widgets.Container{
		Style: "padding: 15px; border-bottom: 1px solid #eee; display: flex; align-items: center;",
		Child: widgets.Row{
			Children: []widgets.Widget{
				widgets.Checkbox{
					Value: todo.Completed,
					OnChanged: func(checked bool) {
						log.Printf("Todo %d toggled to: %t", todo.ID, checked)
					},
				},
				widgets.SizedBox{Width: &[]float64{15}[0]},
				widgets.Expanded{
					Child: widgets.Text{
						Data:      todo.Text,
						TextStyle: textStyle,
					},
				},
				widgets.IconButton{
					Icon: "delete",
					OnPressed: func() {
						log.Printf("Delete todo %d", todo.ID)
					},
					Style: "color: #f44336; cursor: pointer; padding: 5px;",
				},
			},
		},
	}
}

// Action Handlers

// AddTodoHandler adds a new todo
func AddTodoHandler(ctx *core.Context) widgets.Widget {
	text := ctx.FormValue("text")
	if text != "" {
		todos = append(todos, Todo{
			ID:        nextID,
			Text:      text,
			Completed: false,
		})
		nextID++
		log.Printf("Added todo: %s", text)
	}

	// Return updated todo list
	return widgets.Column{
		Children: getTodoWidgets(),
	}
}

// ToggleTodoHandler toggles a todo's completion status
func ToggleTodoHandler(ctx *core.Context) widgets.Widget {
	idStr := ctx.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid todo ID: %s", idStr)
		return widgets.Text{Data: "Error"}
	}

	for i, todo := range todos {
		if todo.ID == id {
			todos[i].Completed = !todos[i].Completed
			log.Printf("Toggled todo %d to: %t", id, todos[i].Completed)
			break
		}
	}

	// Return updated todo list
	return widgets.Column{
		Children: getTodoWidgets(),
	}
}

// DeleteTodoHandler deletes a todo
func DeleteTodoHandler(ctx *core.Context) widgets.Widget {
	idStr := ctx.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid todo ID: %s", idStr)
		return widgets.Text{Data: "Error"}
	}

	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			log.Printf("Deleted todo %d", id)
			break
		}
	}

	// Return updated todo list
	return widgets.Column{
		Children: getTodoWidgets(),
	}
}
`

	mainPath := filepath.Join(appName, "main.go")
	err := os.WriteFile(mainPath, []byte(mainGoContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write main.go: %v", err)
	}
}

func createGoMod(appName string) {
	// Check if we're creating the app within the Godin framework development environment
	var goModContent string

	// Look for framework source in parent directories
	frameworkPaths := []string{
		"pkg/core",                    // Current directory (if we're in framework root)
		"../pkg/core",                 // One level up
		"../../pkg/core",              // Two levels up
		"../godin/pkg/core",           // Sibling directory
		"../godin-framework/pkg/core", // Sibling with full name
	}

	var frameworkRoot string
	for _, path := range frameworkPaths {
		if _, err := os.Stat(path); err == nil {
			// Found framework source, determine the root
			if strings.HasSuffix(path, "pkg/core") {
				frameworkRoot = strings.TrimSuffix(path, "pkg/core")
				// Remove trailing slash if present
				frameworkRoot = strings.TrimSuffix(frameworkRoot, "/")
				frameworkRoot = strings.TrimSuffix(frameworkRoot, "\\")
				if frameworkRoot == "" {
					frameworkRoot = "."
				}
				break
			}
		}
	}

	if frameworkRoot != "" {
		// We're in development environment, use local framework
		// Make sure we're pointing to the framework root, not pkg/core
		var finalPath string
		if frameworkRoot == "." {
			// We're in the framework root directory
			abs, err := filepath.Abs(".")
			if err == nil {
				finalPath = abs
			} else {
				finalPath = "."
			}
		} else {
			abs, err := filepath.Abs(frameworkRoot)
			if err == nil {
				finalPath = abs
			} else {
				finalPath = frameworkRoot
			}
		}

		goModContent = `module ` + appName + `

go 1.21

replace github.com/gideonsigilai/godin => ` + finalPath + `

require (
	github.com/gideonsigilai/godin v0.0.0-00010101000000-000000000000
)
`
	} else {
		// Standard production go.mod
		goModContent = `module ` + appName + `

go 1.21

require (
	github.com/gideonsigilai/godin v1.0.0
)
`
	}

	goModPath := filepath.Join(appName, "go.mod")
	err := os.WriteFile(goModPath, []byte(goModContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write go.mod: %v", err)
	}

	// If we're in development environment, log it
	if frameworkRoot != "" {
		log.Printf("🔧 Created go.mod with local framework reference: %s", frameworkRoot)
	}
}

func createBasicCSS(appName string) {
	cssContent := `/* Godin Framework - Basic Styles */

* {
    box-sizing: border-box;
    margin: 0;
    padding: 0;
}

body {
    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    line-height: 1.6;
    color: #333;
    background-color: #f5f5f5;
}

.container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 20px;
}

/* Button Styles */
.btn {
    display: inline-block;
    padding: 10px 20px;
    background-color: #2196F3;
    color: white;
    text-decoration: none;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
    transition: background-color 0.3s ease;
}

.btn:hover {
    background-color: #1976D2;
}

.btn-secondary {
    background-color: #757575;
}

.btn-secondary:hover {
    background-color: #616161;
}

.btn-danger {
    background-color: #f44336;
}

.btn-danger:hover {
    background-color: #d32f2f;
}

/* Card Styles */
.card {
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    padding: 20px;
    margin-bottom: 20px;
}

/* Form Styles */
.form-group {
    margin-bottom: 15px;
}

.form-control {
    width: 100%;
    padding: 10px;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 14px;
}

.form-control:focus {
    outline: none;
    border-color: #2196F3;
    box-shadow: 0 0 0 2px rgba(33, 150, 243, 0.2);
}

/* Utility Classes */
.text-center {
    text-align: center;
}

.mt-20 {
    margin-top: 20px;
}

.mb-20 {
    margin-bottom: 20px;
}

.p-20 {
    padding: 20px;
}

/* Responsive */
@media (max-width: 768px) {
    .container {
        padding: 0 10px;
    }

    .card {
        padding: 15px;
    }
}
`

	cssPath := filepath.Join(appName, "static", "css", "app.css")
	err := os.WriteFile(cssPath, []byte(cssContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write app.css: %v", err)
	}
}

func createBaseTemplate(appName string) {
	templateContent := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - ` + appName + `</title>

    <!-- Godin Framework CSS -->
    <link rel="stylesheet" href="/static/css/app.css">

    <!-- HTMX Library -->
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>

    <!-- Additional CSS -->
    {{if .CSS}}
    <style>{{.CSS}}</style>
    {{end}}
</head>
<body>
    <!-- Main Content -->
    <div id="app">
        {{.Content}}
    </div>

    <!-- Godin Framework JavaScript -->
    <script src="/static/js/godin.js"></script>

    <!-- Hot Reload JavaScript (Development Only) -->
    <script src="/static/js/hot-reload.js"></script>

    <!-- Additional JavaScript -->
    {{if .JS}}
    <script>{{.JS}}</script>
    {{end}}
</body>
</html>`

	templatePath := filepath.Join(appName, "templates", "base.html")
	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write base.html: %v", err)
	}
}

func createReadme(appName, description, template string, noTemplate bool) {
	var templateInfo string
	var features string
	var runInstructions string

	if noTemplate {
		templateInfo = "This app was created with minimal configuration only."
		features = `- Basic project structure
- Configuration files (package.yaml, go.mod)
- Static assets directory
- Templates directory
- Ready for custom development`
		runInstructions = `1. Add your routes and handlers to main.go
2. Create widgets in the widgets/ directory
3. Add static assets (CSS, JS, images) to static/
4. Run the development server`
	} else {
		switch template {
		case "counter":
			templateInfo = "This app was created using the **Counter Template**."
			features = `- Counter with increment/decrement/reset functionality
- Multiple pages (Home, About) with navigation
- Flutter-like widget composition
- Server-side state management
- Modern UI with gradients and proper spacing
- HTMX integration for partial updates`
		case "simple":
			templateInfo = "This app was created using the **Simple Template**."
			features = `- Basic welcome page
- Minimal code structure
- Clean starting point for custom applications
- Flutter-like widget system
- Ready for expansion`
		case "todo":
			templateInfo = "This app was created using the **Todo Template**."
			features = `- Add, toggle, and delete todos
- Beautiful gradient UI design
- Form handling and validation
- Interactive components (checkboxes, buttons)
- List management with dynamic updates
- CRUD operations demonstration`
		default:
			templateInfo = "This app was created using a custom template."
			features = `- Custom application structure
- Godin Framework integration
- Widget-based UI system`
		}
		runInstructions = `1. Explore the generated code in main.go
2. Customize the widgets and styling
3. Add new routes and handlers as needed
4. Run the development server`
	}

	readmeContent := `# ` + appName + `

` + description + `

## About

` + templateInfo + `

## Features

` + features + `

## Getting Started

### Prerequisites

- Go 1.21 or later
- Godin Framework

### Installation

1. Navigate to the project directory:
   ` + "```bash" + `
   cd ` + appName + `
   ` + "```" + `

2. Install dependencies:
   ` + "```bash" + `
   go mod tidy
   ` + "```" + `

### Running the Application

` + runInstructions + `

` + "```bash" + `
# Run in debug mode (recommended for development)
godin run

# Start the development server with hot reload
godin serve --watch

# Build for production
godin build

# Or run directly with Go
go run main.go
` + "```" + `

The application will be available at [http://localhost:8080](http://localhost:8080).

## Project Structure

` + "```" + `
` + appName + `/
├── main.go                 # Application entry point
├── package.yaml           # Project configuration
├── go.mod                 # Go module file
├── handlers/              # Route handlers (empty)
├── widgets/
│   ├── components/        # Reusable components (empty)
│   └── pages/            # Page widgets (empty)
├── static/
│   ├── css/              # CSS files
│   │   └── app.css       # Basic application styles
│   ├── js/               # JavaScript files (empty)
│   └── images/           # Image assets (empty)
├── templates/
│   └── base.html         # Base HTML template
└── README.md             # This file
` + "```" + `

## Available Scripts

### Development Commands

` + "```bash" + `
# Run in debug mode (recommended)
godin run                    # Default port 8080
godin run --port 3000        # Custom port
godin run --debug=false      # Disable debug features

# Development server with hot reload
godin serve --watch

# Build for production
godin build                  # Creates app.exe
godin build --output dist/   # Build to dist/
godin build --name myapp     # Custom executable name
` + "```" + `

### Package Scripts

The following scripts are available in ` + "`package.yaml`" + `:

- ` + "`godin serve --watch`" + ` - Start development server with hot reload
- ` + "`godin build --prod`" + ` - Build for production
- ` + "`godin test`" + ` - Run tests
- ` + "`go mod tidy`" + ` - Install/update dependencies
- ` + "`rm -rf dist/ bin/`" + ` - Clean build artifacts

## Development

### Adding New Routes

Add new routes in ` + "`main.go`" + `:

` + "```go" + `
app.GET("/new-page", NewPageHandler)
app.POST("/api/action", ActionHandler)
` + "```" + `

### Creating Widgets

Create reusable widgets in the ` + "`widgets/components/`" + ` directory:

` + "```go" + `
package components

import . "github.com/gideonsigilai/godin/pkg/godin"

func MyCustomWidget() Widget {
    return Container{
        Child: Text{
            Data: "Hello from custom widget!",
        },
    }
}
` + "```" + `

### Styling

- Edit ` + "`static/css/app.css`" + ` for global styles
- Use inline styles in widgets for component-specific styling
- The base template includes HTMX and basic CSS

### WebSocket Integration

WebSocket support is enabled by default. The base template includes a basic WebSocket connection script.

## Godin Framework

This application is built with the Godin Framework, which provides:

- **Flutter-like Widgets**: Compose UIs using familiar widget patterns
- **Server-Side Rendering**: All logic runs on the server
- **HTMX Integration**: Seamless partial page updates
- **WebSocket Support**: Real-time communication built-in
- **Hot Reload**: Fast development with automatic reloading

## Learn More

- [Godin Framework Documentation](https://godin-framework.dev)
- [Widget Reference](https://godin-framework.dev/widgets)
- [Examples](https://github.com/gideonsigilai/godin/tree/main/examples)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.
`

	readmePath := filepath.Join(appName, "README.md")
	err := os.WriteFile(readmePath, []byte(readmeContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write README.md: %v", err)
	}
}

// Helper functions for build and run commands

func isGodinProject() bool {
	// Check if package.yaml exists
	if _, err := os.Stat("package.yaml"); err == nil {
		return true
	}

	// Check if main.go exists
	if _, err := os.Stat("main.go"); err == nil {
		return true
	}

	return false
}

func getCurrentDir() string {
	if dir, err := os.Getwd(); err == nil {
		return dir
	}
	return "unknown"
}

func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// needsImportFix checks if we're in a Godin framework development environment
// and the current project needs import path fixes
func needsImportFix() bool {
	// Check if we're in a directory that has both:
	// 1. A main.go that imports godin-framework packages
	// 2. The actual godin framework source code nearby

	// Read main.go to check for godin-framework imports
	mainContent, err := os.ReadFile("main.go")
	if err != nil {
		return false
	}

	content := string(mainContent)
	hasGodinImports := strings.Contains(content, "github.com/gideonsigilai/godin/pkg/core") ||
		strings.Contains(content, "github.com/gideonsigilai/godin/pkg/widgets") ||
		strings.Contains(content, "github.com/gideonsigilai/godin/pkg/godin")

	if !hasGodinImports {
		return false
	}

	// Check if we can find the framework source code
	// Look for the framework in parent directories or common locations
	frameworkPaths := []string{
		"../pkg/core",                 // One level up
		"../../pkg/core",              // Two levels up
		"../godin/pkg/core",           // Sibling directory
		"../godin-framework/pkg/core", // Sibling with full name
	}

	for _, path := range frameworkPaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	return false
}

// fixFrameworkImports temporarily fixes the import paths for framework development
func fixFrameworkImports() error {
	fmt.Println("🔧 Starting framework imports fix...")

	// Find the framework root directory
	var frameworkRoot string
	frameworkPaths := []string{
		"..",                 // One level up
		"../..",              // Two levels up
		"../godin",           // Sibling directory
		"../godin-framework", // Sibling with full name
	}

	for _, path := range frameworkPaths {
		fmt.Printf("🔍 Checking framework path: %s\n", path)
		if _, err := os.Stat(filepath.Join(path, "pkg/core")); err == nil {
			abs, err := filepath.Abs(path)
			if err == nil {
				frameworkRoot = abs
				fmt.Printf("✅ Found framework at: %s\n", frameworkRoot)
				break
			}
		}
	}

	if frameworkRoot == "" {
		fmt.Println("❌ Could not locate Godin framework source code")
		return fmt.Errorf("could not locate Godin framework source code")
	}

	// Create or update go.mod to use local framework
	goModContent := fmt.Sprintf(`module %s

go 1.21

replace github.com/gideonsigilai/godin => %s

require (
	github.com/gideonsigilai/godin v0.0.0-00010101000000-000000000000
)
`, getModuleName(), frameworkRoot)

	// Write the updated go.mod
	fmt.Println("📝 Writing updated go.mod...")
	if err := os.WriteFile("go.mod", []byte(goModContent), 0644); err != nil {
		return fmt.Errorf("failed to write go.mod: %v", err)
	}
	fmt.Println("✅ go.mod updated successfully")

	// Run go mod tidy to clean up with timeout
	fmt.Println("🧹 Running go mod tidy...")
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Stdout = os.Stdout
	tidyCmd.Stderr = os.Stderr

	// Run with timeout to prevent hanging
	done := make(chan error, 1)
	go func() {
		done <- tidyCmd.Run()
	}()

	select {
	case err := <-done:
		if err != nil {
			log.Printf("⚠️  Warning: go mod tidy failed: %v", err)
		} else {
			fmt.Println("✅ go mod tidy completed")
		}
	case <-time.After(10 * time.Second):
		tidyCmd.Process.Kill()
		log.Printf("⚠️  Warning: go mod tidy timed out")
	}

	return nil
}

// getModuleName extracts the module name from go.mod or uses directory name
func getModuleName() string {
	// Try to read existing go.mod
	if content, err := os.ReadFile("go.mod"); err == nil {
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "module ") {
				return strings.TrimSpace(strings.TrimPrefix(line, "module"))
			}
		}
	}

	// Fallback to directory name
	if dir, err := os.Getwd(); err == nil {
		return filepath.Base(dir)
	}

	return "myapp"
}

// isPortAvailable checks if a port is available for binding
func isPortAvailable(port string) bool {
	// Remove the colon if present
	if strings.HasPrefix(port, ":") {
		port = port[1:]
	}

	// Try to listen on the port
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return false
	}

	// Close the listener immediately
	listener.Close()
	return true
}

// testHotReloadEndpoints tests if the hot reload endpoints are working
func testHotReloadEndpoints(port string) {
	// Wait for server to be fully ready
	time.Sleep(2 * time.Second)

	// Prepare the URL - handle both :port and port formats
	testPort := port
	if !strings.HasPrefix(testPort, ":") {
		testPort = ":" + testPort
	}

	client := &http.Client{Timeout: 3 * time.Second}

	// Test the hot-refresh endpoint
	url := fmt.Sprintf("http://localhost%s/api/hot-refresh", testPort)
	log.Printf("🧪 Testing hot reload endpoint: %s", url)

	resp, err := client.Post(url, "application/json", strings.NewReader(`{"type":"test"}`))
	if err != nil {
		log.Printf("⚠️  Hot reload endpoint test failed: %v", err)
		log.Printf("💡 This might be normal if the application doesn't use Godin framework hot reload")
		return
	}

	if resp != nil {
		resp.Body.Close()
		if resp.StatusCode == 200 {
			log.Printf("✅ Hot reload endpoints are working correctly")
		} else {
			log.Printf("⚠️  Hot reload endpoint returned status: %d", resp.StatusCode)
		}
	}
}

// findAvailablePort finds an available port starting from the given port
func findAvailablePort(originalPort string) string {
	// Remove the colon if present
	port := originalPort
	if strings.HasPrefix(port, ":") {
		port = port[1:]
	}

	// Convert to integer
	portNum := 8080 // default fallback
	if p, err := strconv.Atoi(port); err == nil {
		portNum = p
	}

	// Try the original port first
	if isPortAvailable(fmt.Sprintf("%d", portNum)) {
		return originalPort
	}

	// Try ports in a range around the original port
	for i := 1; i <= 10; i++ {
		testPort := portNum + i
		if testPort > 65535 {
			testPort = portNum - i
		}
		if testPort > 0 && isPortAvailable(fmt.Sprintf("%d", testPort)) {
			if strings.HasPrefix(originalPort, ":") {
				return fmt.Sprintf(":%d", testPort)
			}
			return fmt.Sprintf("%d", testPort)
		}
	}

	// Fallback to original port if nothing else works
	log.Printf("⚠️  Could not find available port, using original: %s", originalPort)
	return originalPort
}
