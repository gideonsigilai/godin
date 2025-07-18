# Counter App - Godin Framework Example

A Flutter-inspired counter application built with the Godin Framework, demonstrating routing, state management, and widget composition.

## Features

- **Counter Operations**: Increment, decrement, and reset counter
- **Multi-page Navigation**: Home, About, and Settings pages
- **Flutter-like Widgets**: Uses the Godin Framework's widget system
- **Responsive Design**: Clean, modern UI with proper spacing
- **State Management**: Server-side state management
- **Routing**: Multiple routes with navigation

## Pages

### Home Page (`/`)
- Main counter display
- Increment (+), Decrement (-), and Reset buttons
- Large counter value display
- Navigation to other pages

### About Page (`/about`)
- Information about the application
- Feature list
- Current counter value display
- App description

### Settings Page (`/settings`)
- App configuration display
- System information
- Current counter value
- App metadata

## Routes

- `GET /` - Home page with counter
- `GET /about` - About page
- `GET /settings` - Settings page
- `POST /increment` - Increment counter action
- `POST /decrement` - Decrement counter action
- `POST /reset` - Reset counter action

## Running the App

1. Navigate to the counter-app directory:
   ```bash
   cd examples/counter-app
   ```

2. Run the application:
   ```bash
   go run main.go
   ```

3. Open your browser and visit:
   ```
   http://localhost:8080
   ```

## Architecture

### State Management
- Global counter variable maintained on the server
- Actions modify server state and return updated widgets
- State persists across page navigation

### Widget Structure
- **AppBarWidget**: Reusable navigation header
- **NavButton**: Navigation menu items
- **CreateActionButton**: Action buttons for counter operations
- **FeatureItem**: List items for features
- **SettingItem**: Settings display rows
- **InfoRow**: Information display rows

### Routing
The app demonstrates client-side navigation between different pages while maintaining server-side state. Each page is a separate route handler that returns a complete widget tree.

## Code Structure

```
main.go
├── main()                 # App setup and routing
├── HomeHandler()          # Main counter page
├── AboutHandler()         # About page
├── SettingsHandler()      # Settings page
├── Helper Components
│   ├── AppBarWidget()     # Navigation header
│   ├── NavButton()        # Navigation buttons
│   ├── CreateActionButton() # Action buttons
│   ├── FeatureItem()      # Feature list items
│   ├── SettingItem()      # Settings rows
│   └── InfoRow()          # Info display rows
└── Action Handlers
    ├── IncrementHandler() # Increment counter
    ├── DecrementHandler() # Decrement counter
    └── ResetHandler()     # Reset counter
```

## Widget Patterns

### Layout Composition
```go
widgets.Column{
    Children: []widgets.Widget{
        // Header
        AppBarWidget("Page Title"),
        
        // Content
        widgets.Expanded{
            Child: widgets.Container{
                // Page content
            },
        },
    },
}
```

### Spacing
Uses `widgets.SizedBox` for consistent spacing:
```go
widgets.SizedBox{Height: &[]float64{20}[0]}  // Vertical spacing
widgets.SizedBox{Width: &[]float64{10}[0]}   // Horizontal spacing
```

### Text Styling
Consistent text styling with `TextStyle`:
```go
widgets.Text{
    Data: "Counter Value",
    TextStyle: &widgets.TextStyle{
        FontSize:   &[]float64{24}[0],
        FontWeight: widgets.FontWeightBold,
        Color:      widgets.Color("#2196F3"),
    },
}
```

## Learning Points

This example demonstrates:

1. **Multi-page Applications**: How to structure apps with multiple routes
2. **State Management**: Server-side state that persists across requests
3. **Widget Composition**: Building reusable components
4. **Navigation**: Moving between different pages
5. **Action Handling**: Processing user interactions
6. **Responsive Design**: Creating layouts that work well on different screen sizes

## Next Steps

To extend this example, you could add:

- Persistent storage (database)
- User sessions
- More complex state management
- Client-side interactivity with HTMX
- Styling themes
- Animation effects
- Form validation
- Error handling
