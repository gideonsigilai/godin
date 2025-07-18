package main

import (
	"log"

	"godin-framework/pkg/core"
	"godin-framework/pkg/widgets"
)

var clickCount = 0

func main() {
	app := core.New()

	// Register the home route
	app.GET("/", HomeHandler)

	log.Println("Starting server on :8080")
	if err := app.Serve(":8080"); err != nil {
		log.Fatal(err)
	}
}

// HomeHandler demonstrates the new Flutter-like widget pattern
func HomeHandler(ctx *core.Context) widgets.Widget {
	controller := widgets.TextEditingController{}
	return widgets.Container{
		Style: "max-width: 600px; margin: 50px auto; padding: 20px; font-family: Arial, sans-serif;",
		Child: widgets.Column{
			Children: []widgets.Widget{
				// Header
				widgets.Text{
					Data: "Godin Framework - Flutter-like Widgets",
					TextStyle: &widgets.TextStyle{
						FontSize:   &[]float64{24}[0],
						FontWeight: widgets.FontWeightBold,
						Color:      widgets.Color("#333"),
					},
				},

				// Spacer
				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Description
				widgets.Text{
					Data: "This demonstrates the new widget system where you can create widgets like Flutter and pass Go functions directly for event handlers.",
					TextStyle: &widgets.TextStyle{
						Color: widgets.Color("#666"),
					},
				},

				// Spacer
				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Interactive example
				widgets.Card{
					Style: "padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);",
					Child: widgets.Column{
						Children: []widgets.Widget{
							widgets.Text{
								Data: "Interactive Example",
								TextStyle: &widgets.TextStyle{
									FontSize:   &[]float64{18}[0],
									FontWeight: widgets.FontWeightBold,
								},
							},

							// Spacer
							widgets.SizedBox{Height: &[]float64{15}[0]},

							widgets.Row{
								Children: []widgets.Widget{
									widgets.Button{
										Text: "Click Me!",
										Type: "primary",
										OnClick: func() {
											clickCount++
											log.Printf("Button clicked! Count: %d", clickCount)
										},
									},

									// Spacer
									widgets.SizedBox{Width: &[]float64{10}[0]},

									widgets.Button{
										Text: "Reset",
										Type: "secondary",
										OnClick: func() {
											clickCount = 0
											log.Println("Counter reset!")
										},
									},
								},
							},

							// Spacer
							widgets.SizedBox{Height: &[]float64{15}[0]},

							widgets.Text{
								Data: "Check the server logs to see the click events!",
								TextStyle: &widgets.TextStyle{
									Color: widgets.Color("#888"),
								},
								Style: "font-style: italic;",
							},
						},
					},
				},

				// Spacer
				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Form example
				widgets.Card{
					Style: "padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);",
					Child: widgets.Column{
						Children: []widgets.Widget{
							widgets.Text{
								Data: "Form Example",
								TextStyle: &widgets.TextStyle{
									FontWeight: widgets.FontWeightBold,
									Color:      widgets.ColorBlue,
								},
							},

							// Spacer
							widgets.SizedBox{Height: &[]float64{15}[0]},

							widgets.TextField{
								Controller: &controller,
								ID:         "name-input",
								Decoration: &widgets.InputDecoration{
									HintText: "Enter your name",
								},
							},

							// Spacer
							widgets.SizedBox{Height: &[]float64{15}[0]},

							widgets.Row{
								Children: []widgets.Widget{
									widgets.Checkbox{
										OnChanged: func(checked bool) {
											log.Printf("Checkbox changed to: %t", checked)
										},
									},
									// Spacer
									widgets.SizedBox{Width: &[]float64{8}[0]},
									widgets.Text{
										Data: "I agree to the terms",
									},
								},
							},
						},
					},
				},

				// Spacer
				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Footer
				widgets.Text{
					Data: "All event handlers are Go functions - no more string-based HTMX endpoints!",
					TextStyle: &widgets.TextStyle{
						Color: widgets.Color("#4CAF50"),
					},
					TextAlign: widgets.TextAlignCenter,
					Style:     "font-weight: bold; margin-top: 20px;",
				},
			},
		},
	}
}
