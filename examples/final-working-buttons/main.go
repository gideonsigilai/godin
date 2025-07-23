package main

import (
	"fmt"
	"log"

	"github.com/gideonsigilai/godin/pkg/core"
	"github.com/gideonsigilai/godin/pkg/widgets"
)

func main() {
	fmt.Println("=== STARTING FINAL WORKING BUTTONS DEMO ===")
	
	// Create a new Godin application
	app := core.New()
	fmt.Println("=== APP CREATED ===")

	// Enable WebSocket for real-time state updates
	app.WebSocket().Enable("/ws")
	fmt.Println("=== WEBSOCKET ENABLED ===")

	// Initialize state
	app.State().Set("counter", 0)
	app.State().Set("message", "Welcome! Click buttons to test native Go code execution.")
	fmt.Println("=== STATE INITIALIZED ===")

	// Define the home page
	app.GET("/", func(ctx *core.Context) core.Widget {
		fmt.Println("=== HOME PAGE REQUESTED ===")
		return homePage()
	})

	// Start the server
	log.Println("Starting final working buttons demo on :8086")
	log.Println("Visit http://localhost:8086 to see native buttons in action")
	log.Println("=== SERVER STARTING ===")
	if err := app.Serve(":8086"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func homePage() core.Widget {
	return widgets.Container{
		Style: "padding: 20px; font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; background: #f8f9fa;",
		Child: widgets.Column{
			Children: []widgets.Widget{
				// Title
				widgets.Text{
					Data: "🎉 Native Go Button System - WORKING!",
					TextStyle: &widgets.TextStyle{
						FontSize:   &[]float64{28}[0],
						FontWeight: widgets.FontWeightBold,
						Color:      widgets.Color("#28a745"),
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
						return widgets.Container{
							Style: "background: white; padding: 15px; border-radius: 8px; border-left: 4px solid #007bff; margin: 10px 0;",
							Child: widgets.Text{
								Data: message,
								TextStyle: &widgets.TextStyle{
									FontSize: &[]float64{16}[0],
									Color:    widgets.Color("#495057"),
								},
							},
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
							Style: "background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 25px; border-radius: 12px; text-align: center; margin: 20px 0; box-shadow: 0 4px 15px rgba(0,0,0,0.1);",
							Child: widgets.Text{
								Data: fmt.Sprintf("Counter: %d", counter),
								TextStyle: &widgets.TextStyle{
									FontSize:   &[]float64{32}[0],
									FontWeight: widgets.FontWeightBold,
									Color:      widgets.ColorWhite,
								},
							},
						}
					},
				},

				widgets.SizedBox{Height: &[]float64{30}[0]},

				// Perfect syntax example - exactly as you requested!
				widgets.Button{
					Text:  "✨ Perfect Syntax Test",
					Type:  "primary",
					Style: "margin: 5px; padding: 12px 24px; font-size: 16px;",
					OnPressed: func() {
						fmt.Println("=== BUTTON CLICKED ===")
						current := core.GetStateInt("counter")
						newValue := current + 1
						fmt.Printf("Current: %d, New: %d\n", current, newValue)
						core.SetState("counter", newValue)
						core.SetState("message", fmt.Sprintf("✅ Native Go code executed! Counter: %d", newValue))
						fmt.Println("=== BUTTON CALLBACK COMPLETE ===")
					},
				},

				// More examples
				widgets.Button{
					Text:  "🚀 Complex Logic",
					Type:  "secondary",
					Style: "margin: 5px; padding: 12px 24px;",
					OnPressed: func() {
						// Complex native Go code
						counter := core.GetStateInt("counter")
						
						var message string
						var newCounter int
						
						if counter%2 == 0 {
							newCounter = counter * 2
							message = fmt.Sprintf("🔢 Even number! Doubled from %d to %d", counter, newCounter)
						} else {
							newCounter = counter + 10
							message = fmt.Sprintf("🔢 Odd number! Added 10 from %d to %d", counter, newCounter)
						}
						
						// Update multiple state values
						core.SetState("counter", newCounter)
						core.SetState("message", message)
					},
				},

				widgets.Button{
					Text:  "🔄 Reset",
					Type:  "danger",
					Style: "margin: 5px; padding: 12px 24px;",
					OnPressed: func() {
						// Native Go code
						core.SetState("counter", 0)
						core.SetState("message", "🔄 Counter reset to 0! Ready for more testing.")
					},
				},

				widgets.SizedBox{Height: &[]float64{30}[0]},

				// Status info
				widgets.Container{
					Style: "background: #d4edda; border: 1px solid #c3e6cb; border-radius: 8px; padding: 15px;",
					Child: widgets.Column{
						Children: []widgets.Widget{
							widgets.Text{
								Data: "✅ System Status: FULLY WORKING",
								TextStyle: &widgets.TextStyle{
									FontSize:   &[]float64{16}[0],
									FontWeight: widgets.FontWeightBold,
									Color:      widgets.Color("#155724"),
								},
							},
							widgets.SizedBox{Height: &[]float64{10}[0]},
							widgets.Text{
								Data: "• Native Go code execution: ✅\n• State management: ✅\n• Real-time UI updates: ✅\n• Flutter-style OnPressed: ✅",
								TextStyle: &widgets.TextStyle{
									FontSize: &[]float64{14}[0],
									Color:    widgets.Color("#155724"),
								},
							},
						},
					},
				},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				widgets.Text{
					Data: "Check server logs to see native Go code execution in real-time!",
					TextStyle: &widgets.TextStyle{
						FontSize:  &[]float64{12}[0],
						FontStyle: widgets.FontStyleItalic,
						Color:     widgets.Color("#6c757d"),
					},
					TextAlign: widgets.TextAlignCenter,
				},
			},
		},
	}
}
