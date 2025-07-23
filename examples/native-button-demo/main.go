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

	// Initialize counter state
	app.State().Set("counter", 0)
	app.State().Set("message", "Welcome to Native Button Demo!")

	// Define the home page
	app.GET("/", func(ctx *core.Context) core.Widget {
		return homePage()
	})

	// Start the server
	log.Println("Starting server on :8082")
	log.Println("Visit http://localhost:8082 to see the native button demo")
	if err := app.Serve(":8082"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func homePage() core.Widget {
	return widgets.Container{
		Style: "padding: 20px; font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto;",
		Child: widgets.Column{
			MainAxisAlignment: widgets.MainAxisAlignmentCenter,
			Children: []widgets.Widget{
				// Title
				widgets.Text{
					Data: "Native Button Demo",
					TextStyle: &widgets.TextStyle{
						FontSize:   &[]float64{28}[0],
						FontWeight: widgets.FontWeightBold,
						Color:      widgets.Color("#333333"),
					},
					TextAlign: widgets.TextAlignCenter,
				},

				widgets.SizedBox{Height: &[]float64{30}[0]},

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
								FontSize: &[]float64{18}[0],
								Color:    widgets.Color("#666666"),
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
							Style: "background-color: #f0f0f0; border-radius: 10px; text-align: center; padding: 20px; margin: 20px 0;",
							Child: widgets.Text{
								Data: fmt.Sprintf("Counter: %d", counter),
								TextStyle: &widgets.TextStyle{
									FontSize:   &[]float64{32}[0],
									FontWeight: widgets.FontWeightBold,
									Color:      widgets.Color("#2196F3"),
								},
							},
						}
					},
				},

				widgets.SizedBox{Height: &[]float64{30}[0]},

				// Button row
				widgets.Row{
					MainAxisAlignment: widgets.MainAxisAlignmentCenter,
					Children: []widgets.Widget{
						// Decrement button with native Go code
						widgets.Button{
							Text:  "Decrement",
							Type:  "secondary",
							Style: "margin-right: 10px;",
							OnPressed: func() {
								// Native Go code execution
								currentValue := core.GetStateInt("counter")
								newValue := currentValue - 1
								core.SetState("counter", newValue)
								core.SetState("message", fmt.Sprintf("Decremented to %d", newValue))
							},
						},

						// Reset button
						widgets.Button{
							Text:  "Reset",
							Type:  "danger",
							Style: "margin: 0 10px;",
							OnPressed: func() {
								// Native Go code execution
								core.SetState("counter", 0)
								core.SetState("message", "Counter reset to 0!")
							},
						},

						// Increment button
						widgets.Button{
							Text:  "Increment",
							Type:  "primary",
							Style: "margin-left: 10px;",
							OnPressed: func() {
								// Native Go code execution
								currentValue := core.GetStateInt("counter")
								newValue := currentValue + 1
								core.SetState("counter", newValue)
								core.SetState("message", fmt.Sprintf("Incremented to %d", newValue))
							},
						},
					},
				},

				widgets.SizedBox{Height: &[]float64{30}[0]},

				// Advanced example with complex logic
				widgets.Button{
					Text:  "Complex Logic",
					Type:  "primary",
					Style: "background-color: #9C27B0; width: 200px;",
					OnPressed: func() {
						// Complex native Go code
						counter := core.GetStateInt("counter")

						// Perform some calculations
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

				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Info text
				widgets.Text{
					Data: "All buttons use native Go code execution without HTMX dependencies!",
					TextStyle: &widgets.TextStyle{
						FontSize:  &[]float64{14}[0],
						FontStyle: widgets.FontStyleItalic,
						Color:     widgets.Color("#888888"),
					},
					TextAlign: widgets.TextAlignCenter,
				},
			},
		},
	}
}
