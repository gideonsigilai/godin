package main

import (
	"fmt"
	"log"

	"github.com/gideonsigilai/godin/pkg/core"
	"github.com/gideonsigilai/godin/pkg/widgets"
)

func main() {
	// Create a new Godin application
	app := core.New()

	// Enable WebSocket for real-time state updates
	app.WebSocket().Enable("/ws")

	// Initialize state
	app.State().Set("counter", 0)
	app.State().Set("message", "Welcome to Native Buttons!")

	// Define the home page
	app.GET("/", func(ctx *core.Context) core.Widget {
		return homePage()
	})

	// Start the server
	log.Println("Starting server on :8084")
	log.Println("Visit http://localhost:8084 to see native buttons in action")
	if err := app.Serve(":8084"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func homePage() core.Widget {
	return widgets.Container{
		Style: "padding: 20px; font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;",
		Child: widgets.Column{
			Children: []widgets.Widget{
				// Title
				widgets.Text{
					Data: "🚀 Native Go Button System",
					TextStyle: &widgets.TextStyle{
						FontSize:   &[]float64{28}[0],
						FontWeight: widgets.FontWeightBold,
						Color:      widgets.Color("#2196F3"),
					},
					TextAlign: widgets.TextAlignCenter,
				},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Message display with state
				&widgets.Consumer{
					StateKey: "message",
					Builder: func(value interface{}) widgets.Widget {
						message := "Welcome!"
						if v, ok := value.(string); ok {
							message = v
						}
						return widgets.Text{
							Data: message,
							TextStyle: &widgets.TextStyle{
								FontSize: &[]float64{16}[0],
								Color:    widgets.Color("#666"),
							},
							TextAlign: widgets.TextAlignCenter,
						}
					},
				},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Counter display with state
				&widgets.Consumer{
					StateKey: "counter",
					Builder: func(value interface{}) widgets.Widget {
						counter := 0
						if v, ok := value.(int); ok {
							counter = v
						}
						return widgets.Container{
							Style: "background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 20px; border-radius: 10px; text-align: center; margin: 20px 0;",
							Child: widgets.Text{
								Data: fmt.Sprintf("Counter: %d", counter),
								TextStyle: &widgets.TextStyle{
									FontSize:   &[]float64{24}[0],
									FontWeight: widgets.FontWeightBold,
									Color:      widgets.ColorWhite,
								},
							},
						}
					},
				},

				widgets.SizedBox{Height: &[]float64{30}[0]},

				// Example 1: Simple button with native Go code
				widgets.Button{
					Text:  "Simple Increment",
					Type:  "primary",
					Style: "margin: 5px; padding: 10px 20px;",
					OnPressed: func() {
						// Native Go code execution - exactly as requested!
						x := core.GetStateInt("counter")
						newValue := x + 1
						core.SetState("counter", newValue)
						core.SetState("message", fmt.Sprintf("Incremented to %d", newValue))
					},
				},

				// Example 2: Button with complex logic
				widgets.Button{
					Text:  "Complex Logic",
					Type:  "secondary",
					Style: "margin: 5px; padding: 10px 20px;",
					OnPressed: func() {
						// Full native Go code with complex logic
						counter := core.GetStateInt("counter")
						
						var message string
						var newCounter int
						
						if counter%2 == 0 {
							newCounter = counter * 2
							message = fmt.Sprintf("Even number! Doubled from %d to %d", counter, newCounter)
						} else {
							newCounter = counter + 10
							message = fmt.Sprintf("Odd number! Added 10 from %d to %d", counter, newCounter)
						}
						
						// Update multiple state values
						core.SetState("counter", newCounter)
						core.SetState("message", message)
					},
				},

				// Example 3: Reset button
				widgets.Button{
					Text:  "Reset",
					Type:  "danger",
					Style: "margin: 5px; padding: 10px 20px;",
					OnPressed: func() {
						// Native Go code
						core.SetState("counter", 0)
						core.SetState("message", "Counter reset to 0!")
					},
				},

				widgets.SizedBox{Height: &[]float64{30}[0]},

				// Example 4: Advanced calculation
				widgets.Button{
					Text:  "Fibonacci Next",
					Type:  "primary",
					Style: "margin: 5px; padding: 10px 20px; background-color: #9C27B0;",
					OnPressed: func() {
						// Advanced native Go code
						current := core.GetStateInt("counter")
						
						// Calculate next Fibonacci number
						a, b := 0, 1
						for i := 0; i < current; i++ {
							a, b = b, a+b
						}
						
						nextFib := b
						core.SetState("counter", nextFib)
						core.SetState("message", fmt.Sprintf("Fibonacci sequence: F(%d) = %d", current+1, nextFib))
					},
				},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Info text
				widgets.Text{
					Data: "✅ All buttons use pure native Go code execution!\n🚀 No HTMX dependencies in button logic\n⚡ Real-time UI updates via WebSocket\n🎯 Flutter-style OnPressed callbacks",
					TextStyle: &widgets.TextStyle{
						FontSize:  &[]float64{12}[0],
						Color:     widgets.Color("#888"),
						FontStyle: widgets.FontStyleItalic,
					},
					TextAlign: widgets.TextAlignCenter,
				},
			},
		},
	}
}
