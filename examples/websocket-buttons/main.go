package main

import (
	"fmt"
	"log"

	"github.com/gideonsigilai/godin/pkg/core"
	"github.com/gideonsigilai/godin/pkg/widgets"
)

func main() {
	fmt.Println("🚀 Starting WebSocket-based Native Button System")

	// Create a new Godin application
	app := core.New()

	// Enable WebSocket for real-time state updates
	app.WebSocket().Enable("/ws")

	// Initialize state
	app.State().Set("counter", 0)
	app.State().Set("message", "🎉 WebSocket Native Buttons - Ready to Test!")

	// Define the home page
	app.GET("/", func(ctx *core.Context) core.Widget {
		fmt.Println("📄 Home page requested")
		return homePage()
	})

	// Start the server
	fmt.Println("🚀 Starting WebSocket Native Button Demo on :8092")
	fmt.Println("📱 Visit http://localhost:8092 to test the new button system")
	fmt.Println("✅ WebSocket-based button clicks enabled!")
	fmt.Println("🔧 Server starting...")
	if err := app.Serve(":8092"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func homePage() core.Widget {
	return widgets.Container{
		Style: "padding: 20px; font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); min-height: 100vh;",
		Child: widgets.Column{
			Children: []widgets.Widget{
				// Title
				widgets.Container{
					Style: "background: white; padding: 20px; border-radius: 10px; margin-bottom: 20px; text-align: center;",
					Child: widgets.Text{
						Data: "🚀 WEBSOCKET NATIVE BUTTONS",
						TextStyle: &widgets.TextStyle{
							FontSize:   &[]float64{28}[0],
							FontWeight: widgets.FontWeightBold,
							Color:      widgets.Color("#2196F3"),
						},
					},
				},

				// Message display with state
				&widgets.Consumer{
					StateKey: "message",
					Builder: func(value interface{}) widgets.Widget {
						message := "Welcome!"
						if v, ok := value.(string); ok {
							message = v
						}
						return widgets.Container{
							Style: "background: rgba(255,255,255,0.9); padding: 15px; border-radius: 8px; margin: 10px 0;",
							Child: widgets.Text{
								Data: message,
								TextStyle: &widgets.TextStyle{
									FontSize: &[]float64{16}[0],
									Color:    widgets.Color("#495057"),
								},
								TextAlign: widgets.TextAlignCenter,
							},
						}
					},
				},

				// Counter display with state
				&widgets.Consumer{
					StateKey: "counter",
					Builder: func(value interface{}) widgets.Widget {
						counter := 0
						if v, ok := value.(int); ok {
							counter = v
						}
						return widgets.Container{
							Style: "background: rgba(255,255,255,0.95); padding: 25px; border-radius: 12px; text-align: center; margin: 20px 0; box-shadow: 0 4px 15px rgba(0,0,0,0.1);",
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

				widgets.SizedBox{Height: &[]float64{20}[0]},

				// YOUR EXACT SYNTAX - NOW WITH WEBSOCKET POWER!
				widgets.Button{
					Text:  "✨ Your Perfect Syntax (WebSocket)",
					Type:  "primary",
					Style: "margin: 10px; padding: 15px 30px; font-size: 16px; background: #28a745; border: none; border-radius: 8px;",
					OnPressed: func() {
						fmt.Println("🎉 WEBSOCKET BUTTON CLICKED! NATIVE GO CODE EXECUTING!")
						current := core.GetStateInt("counter")
						newValue := current + 1
						fmt.Printf("Current: %d, New: %d\n", current, newValue)
						core.SetState("counter", newValue)
						core.SetState("message", fmt.Sprintf("✅ WebSocket button worked! Counter: %d", newValue))
						fmt.Println("✅ WEBSOCKET BUTTON CALLBACK COMPLETE!")
					},
				},

				// Complex logic test
				widgets.Button{
					Text:  "🧮 Complex Logic (WebSocket)",
					Type:  "secondary",
					Style: "margin: 10px; padding: 15px 30px; font-size: 16px; background: #6f42c1; border: none; border-radius: 8px; color: white;",
					OnPressed: func() {
						fmt.Println("🧮 Complex logic button clicked!")
						counter := core.GetStateInt("counter")

						var message string
						var newCounter int

						if counter%2 == 0 {
							newCounter = counter * 2
							message = fmt.Sprintf("🔢 Even! WebSocket doubled %d → %d", counter, newCounter)
						} else {
							newCounter = counter + 10
							message = fmt.Sprintf("🔢 Odd! WebSocket added 10: %d → %d", counter, newCounter)
						}

						core.SetState("counter", newCounter)
						core.SetState("message", message)
						fmt.Printf("Complex logic executed: %s\n", message)
					},
				},

				// Reset button
				widgets.Button{
					Text:  "🔄 Reset (WebSocket)",
					Type:  "danger",
					Style: "margin: 10px; padding: 15px 30px; font-size: 16px; background: #dc3545; border: none; border-radius: 8px;",
					OnPressed: func() {
						fmt.Println("🔄 Reset button clicked!")
						core.SetState("counter", 0)
						core.SetState("message", "🔄 WebSocket reset complete! Ready for more testing.")
						fmt.Println("Reset executed successfully")
					},
				},

				widgets.SizedBox{Height: &[]float64{30}[0]},

				// Status
				widgets.Container{
					Style: "background: rgba(40, 167, 69, 0.1); border: 2px solid #28a745; border-radius: 8px; padding: 20px;",
					Child: widgets.Column{
						Children: []widgets.Widget{
							widgets.Text{
								Data: "✅ WEBSOCKET BUTTON SYSTEM: ACTIVE",
								TextStyle: &widgets.TextStyle{
									FontSize:   &[]float64{18}[0],
									FontWeight: widgets.FontWeightBold,
									Color:      widgets.Color("#155724"),
								},
								TextAlign: widgets.TextAlignCenter,
							},
							widgets.SizedBox{Height: &[]float64{10}[0]},
							widgets.Text{
								Data: "• No HTMX dependency: ✅\n• Direct HTTP requests: ✅\n• Native Go execution: ✅\n• Real-time UI updates: ✅\n• Your exact syntax: ✅",
								TextStyle: &widgets.TextStyle{
									FontSize: &[]float64{14}[0],
									Color:    widgets.Color("#155724"),
								},
								TextAlign: widgets.TextAlignCenter,
							},
						},
					},
				},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				widgets.Text{
					Data: "🔍 Check server console for native Go code execution logs!\n📡 Buttons use direct HTTP requests instead of HTMX",
					TextStyle: &widgets.TextStyle{
						FontSize:  &[]float64{12}[0],
						FontStyle: widgets.FontStyleItalic,
						Color:     widgets.ColorWhite,
					},
					TextAlign: widgets.TextAlignCenter,
				},
			},
		},
	}
}
