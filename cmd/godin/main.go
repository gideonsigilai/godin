package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

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
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		watch, _ := cmd.Flags().GetBool("watch")
		startDevServer(port, watch)
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

	// Add subcommands
	packageCmd.AddCommand(packageAddCmd)
	packageCmd.AddCommand(packageListCmd)
	packageCmd.AddCommand(packageRemoveCmd)

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(packageCmd)
}

func main() {
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

func startDevServer(port string, watch bool) {
	log.Printf("Starting Godin development server on port %s", port)
	log.Printf("Watch mode: %v", watch)

	// TODO: Implement server startup
	log.Println("Development server functionality not yet implemented")
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
				"github":  "github.com/godin-framework/godin",
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

	"godin-framework/pkg/core"
	"godin-framework/pkg/widgets"
)

func main() {
	app := core.New()

	// Add your routes here
	app.GET("/", HomeHandler)

	log.Printf("Starting %s on :8080", "` + appName + `")
	log.Println("Visit http://localhost:8080 to see your app")
	if err := app.Serve(":8080"); err != nil {
		log.Fatal(err)
	}
}

// HomeHandler renders the home page
func HomeHandler(ctx *core.Context) widgets.Widget {
	return widgets.Container{
		Style: "max-width: 800px; margin: 0 auto; padding: 20px; font-family: Arial, sans-serif;",
		Child: widgets.Column{
			Children: []widgets.Widget{
				widgets.Text{
					Data: "Welcome to ` + appName + `!",
					TextStyle: &widgets.TextStyle{
						FontSize:   &[]float64{32}[0],
						FontWeight: widgets.FontWeightBold,
						Color:      widgets.Color("#333"),
					},
				},
				widgets.SizedBox{Height: &[]float64{20}[0]},
				widgets.Text{
					Data: "Your Godin app is ready. Start building amazing things!",
					TextStyle: &widgets.TextStyle{
						FontSize: &[]float64{16}[0],
						Color:    widgets.Color("#666"),
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

	"godin-framework/pkg/core"
	"godin-framework/pkg/widgets"
)

// App state
var counter = 0

func main() {
	app := core.New()

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
func HomeHandler(ctx *core.Context) widgets.Widget {
	return widgets.Container{
		Style: "min-height: 100vh; font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;",
		Child: widgets.Column{
			Children: []widgets.Widget{
				// App Bar
				AppBarWidget("` + appName + `"),

				// Main content
				widgets.Expanded{
					Child: widgets.Container{
						Style: "padding: 20px;",
						Child: widgets.Column{
							MainAxisAlignment: widgets.MainAxisAlignmentCenter,
							Children: []widgets.Widget{
								// Counter display section
								widgets.Text{
									Data:      "You have pushed the button this many times:",
									TextAlign: widgets.TextAlignCenter,
									TextStyle: &widgets.TextStyle{
										FontSize: &[]float64{16}[0],
									},
								},

								// Spacer
								widgets.SizedBox{Height: &[]float64{20}[0]},

								// Counter value with ID for HTMX updates
								widgets.Container{
									ID: "counter-display",
									Child: widgets.Text{
										Data:      fmt.Sprintf("%d", counter),
										TextAlign: widgets.TextAlignCenter,
										TextStyle: &widgets.TextStyle{
											FontSize:   &[]float64{48}[0],
											FontWeight: widgets.FontWeightBold,
											Color:      widgets.Color("#2196F3"),
										},
									},
								},

								// Spacer
								widgets.SizedBox{Height: &[]float64{40}[0]},

								// Action buttons row
								widgets.Row{
									MainAxisAlignment: widgets.MainAxisAlignmentCenter,
									Children: []widgets.Widget{
										// Decrement button
										widgets.ElevatedButton{
											Child: widgets.Text{Data: "−"},
											OnPressed: func() {
												counter--
												log.Printf("Counter decremented to: %d", counter)
											},
											Style: "min-width: 60px; height: 60px; font-size: 24px; border-radius: 30px; margin-right: 10px;",
										},

										// Spacer
										widgets.SizedBox{Width: &[]float64{20}[0]},

										// Reset button
										widgets.FilledButton{
											Child: widgets.Text{Data: "Reset"},
											OnPressed: func() {
												counter = 0
												log.Printf("Counter reset to: %d", counter)
											},
											Style: "min-width: 80px; height: 60px; font-size: 16px; border-radius: 30px; background-color: #f44336;",
										},

										// Spacer
										widgets.SizedBox{Width: &[]float64{20}[0]},

										// Increment button
										widgets.ElevatedButton{
											Child: widgets.Text{Data: "+"},
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

	"godin-framework/pkg/core"
	"godin-framework/pkg/widgets"
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
	app := core.New()

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
func HomeHandler(ctx *core.Context) widgets.Widget {
	return widgets.Container{
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

replace godin-framework => ` + finalPath + `

require (
	godin-framework v0.0.0-00010101000000-000000000000
)
`
	} else {
		// Standard production go.mod
		goModContent = `module ` + appName + `

go 1.21

require (
	godin-framework v1.0.0
	github.com/gorilla/mux v1.8.1
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

    <!-- WebSocket Connection (if enabled) -->
    <script>
        // WebSocket connection for real-time updates
        if (window.location.protocol === 'https:') {
            var wsProtocol = 'wss:';
        } else {
            var wsProtocol = 'ws:';
        }

        var ws = new WebSocket(wsProtocol + '//' + window.location.host + '/ws');

        ws.onmessage = function(event) {
            // Handle WebSocket messages
            console.log('WebSocket message:', event.data);
        };

        ws.onopen = function() {
            console.log('WebSocket connected');
        };

        ws.onclose = function() {
            console.log('WebSocket disconnected');
        };
    </script>

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

import "godin-framework/pkg/widgets"

func MyCustomWidget() widgets.Widget {
    return widgets.Container{
        Child: widgets.Text{
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
- [Examples](https://github.com/godin-framework/godin/tree/main/examples)

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
	hasGodinImports := strings.Contains(content, "godin-framework/pkg/core") ||
		strings.Contains(content, "godin-framework/pkg/widgets")

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
	// Find the framework root directory
	var frameworkRoot string
	frameworkPaths := []string{
		"..",                 // One level up
		"../..",              // Two levels up
		"../godin",           // Sibling directory
		"../godin-framework", // Sibling with full name
	}

	for _, path := range frameworkPaths {
		if _, err := os.Stat(filepath.Join(path, "pkg/core")); err == nil {
			abs, err := filepath.Abs(path)
			if err == nil {
				frameworkRoot = abs
				break
			}
		}
	}

	if frameworkRoot == "" {
		return fmt.Errorf("could not locate Godin framework source code")
	}

	// Create or update go.mod to use local framework
	goModContent := fmt.Sprintf(`module %s

go 1.21

replace godin-framework => %s

require (
	godin-framework v0.0.0-00010101000000-000000000000
)
`, getModuleName(), frameworkRoot)

	// Write the updated go.mod
	if err := os.WriteFile("go.mod", []byte(goModContent), 0644); err != nil {
		return fmt.Errorf("failed to write go.mod: %v", err)
	}

	// Run go mod tidy to clean up
	tidyCmd := exec.Command("go", "mod", "tidy")
	if err := tidyCmd.Run(); err != nil {
		// Don't fail if tidy fails, just warn
		log.Printf("⚠️  Warning: go mod tidy failed: %v", err)
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
