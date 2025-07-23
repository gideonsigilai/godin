package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gideonsigilai/godin/pkg/core"
	"github.com/gideonsigilai/godin/pkg/state"
	"github.com/gideonsigilai/godin/pkg/widgets"
)

// CounterApp demonstrates ValueListener usage with a simple counter
type CounterApp struct {
	counter *state.IntNotifier
}

// NewCounterApp creates a new counter application
func NewCounterApp() *CounterApp {
	return &CounterApp{
		counter: state.NewIntNotifierWithID("main-counter", 0),
	}
}

// BuildUI builds the counter UI using ValueListener
func (app *CounterApp) BuildUI() widgets.Widget {
	// Create ValueListener that rebuilds when counter changes
	counterDisplay := widgets.NewValueListenerInt(app.counter, func(value int) widgets.Widget {
		return &widgets.Column{
			Style: "text-align: center; padding: 20px; font-size: 24px; color: #333;",
			Children: []widgets.Widget{
				&widgets.Text{
					Data:  fmt.Sprintf("Count: %d", value),
					Style: "font-weight: bold; margin-bottom: 20px;",
				},
			},
		}
	})

	// Set custom properties for the ValueListener
	counterDisplay.ID = "counter-display"
	counterDisplay.Class = "counter-widget"
	counterDisplay.OnValueChanged = func(oldValue, newValue int) {
		fmt.Printf("Counter changed from %d to %d\n", oldValue, newValue)
	}

	// Create buttons for incrementing/decrementing
	incrementButton := &widgets.Button{
		Text:  "Increment",
		Style: "padding: 10px 20px; margin: 5px; background-color: #007bff; color: white; border: none; border-radius: 4px; cursor: pointer;",
		OnPressed: func() {
			app.counter.SetValue(app.counter.Value() + 1)
		},
	}

	decrementButton := &widgets.Button{
		Text:  "Decrement",
		Style: "padding: 10px 20px; margin: 5px; background-color: #dc3545; color: white; border: none; border-radius: 4px; cursor: pointer;",
		OnPressed: func() {
			app.counter.SetValue(app.counter.Value() - 1)
		},
	}

	resetButton := &widgets.Button{
		Text:  "Reset",
		Style: "padding: 10px 20px; margin: 5px; background-color: #6c757d; color: white; border: none; border-radius: 4px; cursor: pointer;",
		OnPressed: func() {
			app.counter.SetValue(0)
		},
	}

	// Create the main layout
	return &widgets.Column{
		Style: "max-width: 600px; margin: 50px auto; padding: 30px; border: 1px solid #ddd; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);",
		Children: []widgets.Widget{
			&widgets.Text{
				Data:  "ValueListener Counter Example",
				Style: "font-size: 28px; font-weight: bold; text-align: center; margin-bottom: 30px; color: #333;",
			},
			counterDisplay,
			&widgets.Column{
				Style: "text-align: center; margin-top: 20px;",
				Children: []widgets.Widget{
					incrementButton,
					decrementButton,
					resetButton,
				},
			},
			&widgets.Text{
				Data:  "This counter automatically updates using ValueListener. The display rebuilds whenever the counter value changes.",
				Style: "margin-top: 30px; padding: 15px; background-color: #f8f9fa; border-radius: 4px; color: #666; font-size: 14px; text-align: center;",
			},
		},
	}
}

// Handler handles HTTP requests
func (app *CounterApp) Handler(w http.ResponseWriter, r *http.Request) {
	// Create Godin app context
	godinApp := core.New()
	ctx := core.NewContext(w, r, godinApp)

	// Build and render the UI
	ui := app.BuildUI()
	html := ui.Render(ctx)

	// Create complete HTML page
	page := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>ValueListener Counter Example</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            margin: 0;
            padding: 0;
            background-color: #f5f5f5;
        }
        .counter-widget {
            transition: all 0.3s ease;
        }
        .counter-widget:hover {
            transform: scale(1.02);
        }
        button:hover {
            opacity: 0.9;
            transform: translateY(-1px);
        }
        button:active {
            transform: translateY(0);
        }
    </style>
</head>
<body>
    %s
    <script>
        // Add some interactivity
        document.addEventListener('DOMContentLoaded', function() {
            // Listen for value changes
            document.addEventListener('valueChanged', function(event) {
                console.log('Counter value changed:', event.detail);
                
                // Add visual feedback
                const counterElement = document.getElementById('counter-display');
                if (counterElement) {
                    counterElement.style.transform = 'scale(1.1)';
                    setTimeout(() => {
                        counterElement.style.transform = 'scale(1)';
                    }, 200);
                }
            });
        });
    </script>
</body>
</html>`, html)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(page))
}

func main() {
	app := NewCounterApp()

	http.HandleFunc("/", app.Handler)

	fmt.Println("ValueListener Counter Example running on http://localhost:8080")
	fmt.Println("Open your browser and try clicking the buttons!")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
