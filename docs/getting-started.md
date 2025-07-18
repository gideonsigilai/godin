# Getting Started with Godin Framework

Welcome to the Godin Framework! This guide will help you get started building web applications with Go, HTMX, and WebSockets.

## Installation

First, make sure you have Go 1.21 or later installed on your system.

Clone the Godin Framework repository:
```bash
git clone https://github.com/godin-framework/godin
cd godin
```

Build the CLI tool:
```bash
go build -o bin/godin cmd/godin/main.go
```

## Creating Your First App

The Godin Framework provides a `create` command to quickly scaffold new applications with different templates.

### Basic Usage

```bash
godin create [app-name] [flags]
```

### Available Templates

#### 1. Counter Template (Default)
Creates a Flutter-inspired counter application with multiple pages and navigation.

```bash
godin create mycounterapp
# or explicitly specify the template
godin create mycounterapp --template counter
```

**Features:**
- Counter with increment, decrement, and reset functionality
- Multiple pages (Home, About)
- Navigation between pages
- Flutter-like widget composition
- Server-side state management

#### 2. Simple Template
Creates a minimal application with basic structure.

```bash
godin create mysimpleapp --template simple
```

**Features:**
- Basic welcome page
- Minimal code structure
- Good starting point for custom applications

#### 3. Todo Template
Creates a full-featured todo application.

```bash
godin create mytodoapp --template todo
```

**Features:**
- Add, toggle, and delete todos
- Beautiful gradient UI
- Form handling
- List management
- Interactive components

#### 4. No Template (Config Only)
Creates only the configuration files and directory structure without template code.

```bash
godin create myapp --no-template
```

**Features:**
- Directory structure only
- package.yaml configuration
- Minimal main.go with basic setup
- Perfect for starting from scratch

### Command Options

- `--template, -t`: Specify the template to use (counter, simple, todo)
- `--no-template`: Create only config files without template code
- `--help, -h`: Show help for the create command

### Examples

```bash
# Create a counter app (default template)
godin create my-awesome-app

# Create a todo app
godin create my-todo-app --template todo

# Create a minimal app structure
godin create my-custom-app --no-template

# Create a simple app
godin create my-simple-app --template simple
```

## Project Structure

When you create a new app, the following directory structure is generated:

```
my-app/
├── main.go                 # Application entry point
├── package.yaml           # Dependencies and configuration
├── handlers/               # Route handlers (empty)
├── widgets/
│   ├── components/         # Reusable components (empty)
│   └── pages/             # Page widgets (empty)
├── static/
│   ├── css/               # CSS files (empty)
│   ├── js/                # JavaScript files (empty)
│   └── images/            # Image assets (empty)
└── templates/             # HTML templates (empty)
```

## Configuration File (package.yaml)

Each app includes a `package.yaml` file with comprehensive configuration:

```yaml
name: my-app
version: 1.0.0
description: A Godin framework application

dependencies:
  godin-framework:
    github: github.com/godin-framework/godin
    version: v1.0.0

scripts:
  dev: godin serve --watch
  build: godin build --prod
  test: godin test
  install: go mod tidy
  clean: rm -rf dist/ bin/

config:
  server:
    port: "8080"
    host: localhost
  websocket:
    enabled: true
    path: /ws
  static:
    dir: static
    cache: true

development:
  hot_reload: true
  file_watching: true
  auto_restart: true
  debug_mode: true
  watch:
    - "**/*.go"
    - "**/*.html"
    - "**/*.css"
    - "**/*.js"
    - "package.yaml"
```

## Running Your App

After creating your app, you have several options for running it:

### Development Mode (Recommended)

Use `godin run` for development with enhanced debugging:

```bash
cd my-app
godin run
```

**Features:**
- 🐛 Debug mode enabled by default
- 📊 Environment variable logging
- 🔍 Detailed startup information
- 🚀 Graceful shutdown with Ctrl+C
- 🌐 Runs on http://localhost:8080

**Options:**
```bash
godin run --port 3000        # Custom port
godin run --debug=false      # Disable debug features
```

### Development Server (Hot Reload)

Use `godin serve` for hot reload during development:

```bash
cd my-app
godin serve --watch
```

### Production Build

Build your app into a standalone executable:

```bash
cd my-app
godin build
```

This creates `app.exe` (Windows) or `app` (Unix) in the current directory.

**Build Options:**
```bash
godin build --output dist/   # Build to dist/app.exe
godin build --name myapp     # Creates myapp.exe
```

Your app will be available at `http://localhost:8080` (or your custom port).

## Next Steps

1. **Explore the Templates**: Try creating apps with different templates to see various patterns
2. **Customize Your App**: Modify the generated code to fit your needs
3. **Add Routes**: Add new routes and handlers in your main.go
4. **Create Widgets**: Build reusable components in the widgets directory
5. **Add Static Assets**: Place CSS, JavaScript, and images in the static directory

## Template Details

### Counter Template
The counter template demonstrates:
- **State Management**: Global counter variable maintained on the server
- **Widget Composition**: Flutter-like widget hierarchy
- **Navigation**: Multiple pages with navigation between them
- **Event Handling**: Button clicks that modify server state
- **Styling**: Modern UI with gradients and proper spacing

### Todo Template
The todo template showcases:
- **Data Structures**: Todo struct with ID, text, and completion status
- **CRUD Operations**: Create, read, update, and delete todos
- **Form Handling**: Input fields and form submission
- **Dynamic UI**: List that updates based on data changes
- **Interactive Elements**: Checkboxes and delete buttons

### Simple Template
The simple template provides:
- **Basic Structure**: Minimal app with one page
- **Clean Code**: Easy to understand and extend
- **Starting Point**: Perfect foundation for custom applications

## Tips

1. **Start with Templates**: Use templates to understand patterns and best practices
2. **Incremental Development**: Start simple and add features gradually
3. **Widget Reuse**: Create reusable widgets in the components directory
4. **State Management**: Keep state on the server for simplicity
5. **Hot Reload**: Use `--watch` flag for automatic reloading during development

## Troubleshooting

### Common Issues

1. **Port Already in Use**: Change the port in package.yaml or use a different port
2. **Import Errors**: Make sure to run `go mod tidy` in your app directory
3. **Template Not Found**: Check the spelling of the template name (counter, simple, todo)

### Getting Help

- Use `godin create --help` for command help
- Check the examples directory for more complex applications
- Refer to the widgets documentation for available components
