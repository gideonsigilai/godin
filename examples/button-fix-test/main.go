package main

import (
	"fmt"
	"log"

	"github.com/gideonsigilai/godin/pkg/core"
	"github.com/gideonsigilai/godin/pkg/widgets"
)

func main() {
	fmt.Println("=== BUTTON FIX TEST STARTING ===")
	
	// Create a new Godin application
	app := core.New()
	fmt.Println("=== APP CREATED ===")

	// Initialize state
	app.State().Set("counter", 0)
	fmt.Println("=== STATE INITIALIZED ===")

	// Define the home page
	app.GET("/", func(ctx *core.Context) core.Widget {
		fmt.Println("=== HOME PAGE REQUESTED ===")
		return homePage()
	})
	fmt.Println("=== ROUTES REGISTERED ===")

	// Start the server
	fmt.Println("=== STARTING SERVER ===")
	log.Println("Starting button fix test on :8088")
	log.Println("Visit http://localhost:8088 to test the fix")
	if err := app.Serve(":8088"); err != nil {
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
					Data: "Button Fix Test",
					TextStyle: &widgets.TextStyle{
						FontSize: &[]float64{24}[0],
						Color:    widgets.Color("#333"),
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
						return widgets.Text{
							Data: fmt.Sprintf("Counter: %d", counter),
							TextStyle: &widgets.TextStyle{
								FontSize: &[]float64{18}[0],
								Color:    widgets.Color("#666"),
							},
						}
					},
				},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Test button with native Go code
				widgets.Button{
					Text: "Test Button (Click Me!)",
					Type: "primary",
					OnPressed: func() {
						fmt.Println("🎉 BUTTON CLICKED! NATIVE GO CODE EXECUTING!")
						current := core.GetStateInt("counter")
						newValue := current + 1
						fmt.Printf("Incrementing counter from %d to %d\n", current, newValue)
						core.SetState("counter", newValue)
						fmt.Println("✅ BUTTON CALLBACK COMPLETE!")
					},
				},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				widgets.Text{
					Data: "If the button works, you'll see logs in the server console and the counter will update!",
					TextStyle: &widgets.TextStyle{
						FontSize: &[]float64{14}[0],
						Color:    widgets.Color("#888"),
					},
				},
			},
		},
	}
}
