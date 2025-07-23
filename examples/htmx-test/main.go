package main

import (
	"fmt"
	"log"

	"github.com/gideonsigilai/godin/pkg/core"
	"github.com/gideonsigilai/godin/pkg/widgets"
)

func main() {
	fmt.Println("=== STARTING HTMX TEST ===")

	// Create a new Godin application
	app := core.New()
	fmt.Println("=== APP CREATED ===")

	// Define the home page
	app.GET("/", func(ctx *core.Context) core.Widget {
		fmt.Println("=== HOME PAGE REQUESTED ===")
		return homePage()
	})

	// Simple test endpoint
	app.GET("/test", func(ctx *core.Context) core.Widget {
		fmt.Println("=== TEST ENDPOINT HIT ===")
		return widgets.Text{Data: "Test endpoint works!"}
	})

	// Start the server
	log.Println("Starting HTMX test server on :8087")
	log.Println("Visit http://localhost:8087 to test HTMX")
	if err := app.Serve(":8087"); err != nil {
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
					Data: "HTMX Test",
					TextStyle: &widgets.TextStyle{
						FontSize: &[]float64{24}[0],
						Color:    widgets.Color("#333"),
					},
				},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Manual HTMX button test
				widgets.Text{
					Data: "Manual HTMX test: Check browser console for HTMX activity",
					TextStyle: &widgets.TextStyle{
						FontSize: &[]float64{14}[0],
						Color:    widgets.Color("#666"),
					},
				},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Native button for comparison
				widgets.Button{
					Text: "Native Button",
					Type: "primary",
					OnPressed: func() {
						fmt.Println("=== NATIVE BUTTON CLICKED ===")
					},
				},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				widgets.Text{
					Data: "Check server logs for button clicks",
					TextStyle: &widgets.TextStyle{
						FontSize: &[]float64{14}[0],
						Color:    widgets.Color("#666"),
					},
				},
			},
		},
	}
}
