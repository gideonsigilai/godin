# Godin Framework

> Server-Side Go with HTMX/WebSockets Framework

Godin is a modern web framework for Go that combines the power of server-side rendering with the interactivity of HTMX and WebSockets. Build Flutter-like applications with Go on the backend.

## ✨ Features

- **Flutter-like Widgets**: Compose UIs using familiar widget patterns
- **Server-Side Rendering**: All logic runs on the server for simplicity
- **HTMX Integration**: Seamless partial page updates
- **WebSocket Support**: Real-time communication built-in
- **Hot Reload**: Fast development with automatic reloading
- **Package Management**: Custom package system for reusable components
- **Template System**: Quick project scaffolding with multiple templates

## 🚀 Quick Start

### Installation

```bash
git clone https://github.com/gideonsigilai/godin
cd godin
go build -o bin/godin cmd/godin/main.go
```

### Create Your First App

```bash
# Create a counter app (default template)
./bin/godin create my-awesome-app

# Create a todo app
./bin/godin create my-todo-app --template todo

# Create minimal structure only
./bin/godin create my-custom-app --no-template
```

### Run Your App

```bash
cd my-awesome-app

# Run in debug mode (recommended for development)
godin run

# Or use the development server with hot reload
godin serve --watch

# Build for production
godin build
```

Visit `http://localhost:8080` to see your app!

## 📋 Available Templates

### Counter Template (Default)
A Flutter-inspired counter app with navigation and state management.

```bash
godin create my-counter-app
```

### Todo Template
A full-featured todo application with CRUD operations.

```bash
godin create my-todo-app --template todo
```

### Simple Template
A minimal starting point for custom applications.

```bash
godin create my-simple-app --template simple
```

### No Template
Just the configuration and directory structure.

```bash
godin create my-custom-app --no-template
```

## 🏗️ Architecture

```
Godin App
├── Server-Side Logic (Go)
├── Widget System (Flutter-like)
├── HTMX Integration (Partial Updates)
├── WebSocket Manager (Real-time)
└── Static Assets (CSS/JS/Images)
```

## 📖 Documentation

- [Getting Started](docs/getting-started.md) - Complete setup and usage guide
- [Widget System](docs/widgets.md) - Available widgets and patterns
- [Examples](docs/examples.md) - Sample applications and code

## 🛠️ CLI Commands

```bash
# Create new applications
godin create <app-name> [--template <template>] [--no-template]

# Run in debug mode
godin run [--port 8080] [--debug]

# Development server with hot reload
godin serve [--port 8080] [--watch]

# Build for production
godin build [--output .] [--name app]

# Package management
godin package add <github-url>
godin package list
godin package remove <package-name>

# Initialize new project (legacy)
godin init <project-name>
```

## 🎯 Example: Counter App

```go
package main

import (
    "fmt"
    "log"
    "godin-framework/pkg/core"
    "godin-framework/pkg/widgets"
)

var counter = 0

func main() {
    app := core.New()

    app.GET("/", HomeHandler)
    app.POST("/increment", IncrementHandler)

    log.Println("Starting on :8080")
    app.Serve(":8080")
}

func HomeHandler(ctx *core.Context) widgets.Widget {
    return widgets.Container{
        Child: widgets.Column{
            Children: []widgets.Widget{
                widgets.Text{
                    Data: fmt.Sprintf("Count: %d", counter),
                    TextStyle: &widgets.TextStyle{
                        FontSize: &[]float64{24}[0],
                    },
                },
                widgets.ElevatedButton{
                    Child: widgets.Text{Data: "+"},
                    OnPressed: func() {
                        counter++
                    },
                },
            },
        },
    }
}

func IncrementHandler(ctx *core.Context) widgets.Widget {
    counter++
    return widgets.Text{
        Data: fmt.Sprintf("%d", counter),
    }
}
```

## 🔧 Configuration (package.yaml)

```yaml
name: my-app
version: 1.0.0
description: A Godin framework application

dependencies:
  godin-framework:
    github: github.com/gideonsigilai/godin
    version: v1.0.0

scripts:
  dev: godin serve --watch
  build: godin build --prod
  test: godin test

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
```

## 🔧 Development Workflow

### Debug Mode (`godin run`)
Perfect for development with enhanced debugging features:

```bash
godin run                    # Run on port 8080 with debug enabled
godin run --port 3000        # Run on custom port
godin run --debug=false      # Run without debug features
```

**Debug Features:**
- 🐛 Enhanced error messages and stack traces
- 📊 Environment variable logging
- 🔍 Detailed startup information
- 🚀 Graceful shutdown with Ctrl+C

### Production Build (`godin build`)
Compile your app into a standalone executable:

```bash
godin build                  # Creates app.exe in current directory
godin build --output dist/   # Build to dist/app.exe
godin build --name myapp     # Creates myapp.exe
```

**Build Features:**
- 📦 Single executable file (app.exe on Windows)
- 📊 File size reporting
- 🚀 Ready for deployment
- ✅ Cross-platform support

## 🌟 Why Godin?

- **Familiar Patterns**: If you know Flutter, you'll feel at home
- **Server-Side Simplicity**: No complex client-side state management
- **Real-Time Ready**: WebSockets built-in for live updates
- **Fast Development**: Hot reload and template scaffolding
- **Go Performance**: Leverage Go's speed and reliability
- **HTMX Power**: Modern web interactions without heavy JavaScript

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- [Documentation](docs/)
- [Examples](examples/)
- [Package Registry](https://registry.godin-framework.dev)
- [Community Discord](https://discord.gg/godin-framework)

---

**Built with ❤️ by the Godin Framework Team**