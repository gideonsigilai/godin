package main

import (
	"fmt"
	"log"

	"github.com/gideonsigilai/godin/pkg/core"
	"github.com/gideonsigilai/godin/pkg/widgets"
)

func main() {
	fmt.Println("Starting simple button test...")

	// Create a new Godin application
	app := core.New()
	fmt.Println("App created")

	// Enable WebSocket for real-time state updates
	app.WebSocket().Enable("/ws")
	fmt.Println("WebSocket enabled")

	// Initialize counter state
	app.State().Set("counter", 0)
	fmt.Println("State initialized")

	// Define the home page
	app.GET("/", func(ctx *core.Context) core.Widget {
		return homePage()
	})

	// Debug endpoint to see raw HTML
	app.GET("/debug", func(ctx *core.Context) core.Widget {
		button := widgets.Button{
			Text: "Debug Button",
			Type: "primary",
			OnPressed: func() {
				fmt.Println("Debug button clicked!")
			},
		}
		html := button.Render(ctx)
		fmt.Printf("Button HTML: %s\n", html)
		return widgets.Text{Data: html}
	})

	// Start the server
	log.Println("Starting simple button test server on :8085")
	log.Println("Visit http://localhost:8085 to test buttons")
	if err := app.Serve(":8085"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func homePage() core.Widget {
	return widgets.Container{
		Style: "padding: 20px; font-family: Arial, sans-serif;",
		Child: widgets.Column{
			Children: []widgets.Widget{
				// Title
				widgets.Text{
					Data: "Simple Button Test",
					TextStyle: &widgets.TextStyle{
						FontSize: &[]float64{24}[0],
						Color:    widgets.Color("#333"),
					},
				},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Counter display with Consumer widget
				&widgets.Consumer{
					StateKey: "counter",
					Builder: func(value interface{}) widgets.Widget {
						counter := 0
						if v, ok := value.(int); ok {
							counter = v
						}
						return widgets.Text{
							Data: fmt.Sprintf("Counter: %d (Click button to increment)", counter),
							TextStyle: &widgets.TextStyle{
								FontSize: &[]float64{18}[0],
								Color:    widgets.Color("#666"),
							},
						}
					},
				},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Simple button
				widgets.Button{
					Text: "Test Button",
					Type: "primary",
					OnPressed: func() {
						fmt.Println("=== BUTTON CLICKED ===")
						current := core.GetStateInt("counter")
						newValue := current + 1
						fmt.Printf("Current: %d, New: %d\n", current, newValue)
						core.SetState("counter", newValue)
						fmt.Println("=== BUTTON CALLBACK COMPLETE ===")
					},
				},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Manual state display (without Consumer)
				widgets.Text{
					Data: "Check server logs for button click events",
					TextStyle: &widgets.TextStyle{
						FontSize: &[]float64{14}[0],
						Color:    widgets.Color("#888"),
					},
				},
			},
		},
	}

}
