package main

import (
	"fmt"
	"log"

	"github.com/gideonsigilai/godin/pkg/core"
	"github.com/gideonsigilai/godin/pkg/widgets"
)

// App state for button demo using TextEditingController and setState
var (
	clickCountController = widgets.NewTextEditingController("0")
	buttonEnabled        = true
	selectedType         = "primary"
)

func main() {
	app := core.New()

	// Enable WebSocket for real-time state updates
	app.WebSocket().Enable("/ws")

	// Initialize state
	clickCountController.SetText("0")

	// Routes
	app.GET("/", HomeHandler)

	log.Println("Starting Button Demo on :8081")
	log.Println("Visit http://localhost:8081 to see the comprehensive button demo")
	log.Println("This demo now uses the new callback system with automatic HTMX integration")
	if err := app.Serve(":8081"); err != nil {
		log.Fatal(err)
	}
}

// HomeHandler renders the main button demo page
func HomeHandler(ctx *core.Context) widgets.Widget {
	x := "heelo"
	return widgets.Container{
		Style: "min-height: 100vh; font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background: #f5f5f5;",
		Child: widgets.Column{
			Children: []widgets.Widget{
				// Header
				widgets.Container{
					Style: "background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 40px 20px; text-align: center;",
					Child: widgets.Column{
						Children: []widgets.Widget{
							widgets.Text{
								Data: "Godin Framework - Button Demo " + x,
								TextStyle: &widgets.TextStyle{
									FontSize:   &[]float64{32}[0],
									FontWeight: widgets.FontWeightBold,
									Color:      widgets.ColorWhite,
								},
							},
							widgets.SizedBox{Height: &[]float64{10}[0]},
							widgets.Text{
								Data: "Comprehensive demonstration of all button types and state management",
								TextStyle: &widgets.TextStyle{
									FontSize: &[]float64{16}[0],
									Color:    widgets.ColorWhite,
								},
							},
						},
					},
				},

				// Main content
				widgets.Container{
					Style: "max-width: 1200px; margin: 0 auto; padding: 40px 20px;",
					Child: widgets.Column{
						Children: []widgets.Widget{
							// Counter Section
							CounterSection(),

							widgets.SizedBox{Height: &[]float64{40}[0]},

							// Basic Button Section
							BasicButtonSection(),

							widgets.SizedBox{Height: &[]float64{40}[0]},

							// Material Design Buttons Section
							MaterialButtonSection(),

							widgets.SizedBox{Height: &[]float64{40}[0]},

							// Interactive Buttons Section
							InteractiveButtonSection(),

							widgets.SizedBox{Height: &[]float64{40}[0]},

							// State Management Demo
							StateManagementSection(),
						},
					},
				},
			},
		},
	}
}

// CounterSection demonstrates basic button functionality with state
func CounterSection() widgets.Widget {
	return widgets.Card{
		Style: "padding: 30px; border-radius: 12px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); background: white;",
		Child: widgets.Column{
			Children: []widgets.Widget{
				widgets.Text{
					Data: "Counter Demo - New Callback System",
					TextStyle: &widgets.TextStyle{
						FontSize:   &[]float64{24}[0],
						FontWeight: widgets.FontWeightBold,
						Color:      widgets.Color("#333"),
					},
				},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Counter buttons with new callback system
				widgets.Row{
					MainAxisAlignment: widgets.MainAxisAlignmentCenter,
					Children: []widgets.Widget{
						widgets.Button{
							ID:   "decrement-btn",
							Text: "Decrement",
							Type: "secondary",
							OnPressed: func() {
								currentValue := parseIntFromController(clickCountController)
								newValue := currentValue - 1
								clickCountController.SetText(fmt.Sprintf("%d", newValue))
								log.Printf("Counter decremented to: %d", newValue)
							},
							Style: "margin-right: 10px;",
						},

						widgets.Button{
							ID:   "reset-btn",
							Text: "Reset",
							Type: "danger",
							OnPressed: func() {
								clickCountController.SetText("0")
								log.Println("Counter reset to: 0")
							},
							Style: "margin: 0 10px;",
						},

						widgets.Button{
							ID:   "increment-btn",
							Text: "Increment",
							Type: "primary",
							OnPressed: func() {
								currentValue := parseIntFromController(clickCountController)
								newValue := currentValue + 1
								clickCountController.SetText(fmt.Sprintf("%d", newValue))
								log.Printf("Counter incremented to: %d", newValue)
							},
							Style: "margin-left: 10px;",
						},
					},
				},

				widgets.SizedBox{Height: &[]float64{10}[0]},

				// Text field to demonstrate TextEditingController
				widgets.Text{
					Data: "Direct Input (demonstrates TextEditingController):",
					TextStyle: &widgets.TextStyle{
						FontSize:   &[]float64{16}[0],
						FontWeight: widgets.FontWeightBold,
						Color:      widgets.Color("#666"),
					},
				},

				widgets.SizedBox{Height: &[]float64{10}[0]},

				widgets.TextField{
					ID:         "counter-input",
					Controller: clickCountController,
					OnChanged: func(value string) {
						log.Printf("Counter value changed via text input: %s", value)
					},
					Style: "max-width: 200px; margin: 0 auto;",
				},
			},
		},
	}
}

// Helper function to parse integer from TextEditingController
func parseIntFromController(controller *widgets.TextEditingController) int {
	value := 0
	if controller.Text() != "" {
		if parsed, err := fmt.Sscanf(controller.Text(), "%d", &value); err != nil || parsed != 1 {
			value = 0
		}
	}
	return value
}

// BasicButtonSection demonstrates basic button types
func BasicButtonSection() widgets.Widget {
	return widgets.Card{
		Style: "padding: 30px; border-radius: 12px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); background: white;",
		Child: widgets.Column{
			Children: []widgets.Widget{
				widgets.Text{
					Data: "Basic Button Types",
					TextStyle: &widgets.TextStyle{
						FontSize:   &[]float64{24}[0],
						FontWeight: widgets.FontWeightBold,
						Color:      widgets.Color("#333"),
					},
				},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				widgets.Row{
					MainAxisAlignment: widgets.MainAxisAlignmentSpaceAround,
					Children: []widgets.Widget{
						widgets.Button{
							Text: "Default",
							OnPressed: func() {
								log.Println("Default button clicked")
							},
						},

						widgets.Button{
							Text: "Primary",
							Type: "primary",
							OnPressed: func() {
								log.Println("Primary button clicked")
							},
						},

						widgets.Button{
							Text: "Secondary",
							Type: "secondary",
							OnPressed: func() {
								log.Println("Secondary button clicked")
							},
						},

						widgets.Button{
							Text: "Danger",
							Type: "danger",
							OnPressed: func() {
								log.Println("Danger button clicked")
							},
						},
					},
				},
			},
		},
	}
}

// MaterialButtonSection demonstrates Material Design buttons
func MaterialButtonSection() widgets.Widget {
	return widgets.Card{
		Style: "padding: 30px; border-radius: 12px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); background: white;",
		Child: widgets.Column{
			Children: []widgets.Widget{
				widgets.Text{
					Data: "Material Design Buttons",
					TextStyle: &widgets.TextStyle{
						FontSize:   &[]float64{24}[0],
						FontWeight: widgets.FontWeightBold,
						Color:      widgets.Color("#333"),
					},
				},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Row 1: Elevated and Filled buttons
				widgets.Row{
					MainAxisAlignment: widgets.MainAxisAlignmentSpaceAround,
					Children: []widgets.Widget{
						widgets.ElevatedButton{
							Child: widgets.Text{Data: "Elevated Button"},
							OnPressed: func() {
								log.Println("Elevated button pressed")
							},
						},

						widgets.FilledButton{
							Child: widgets.Text{Data: "Filled Button"},
							OnPressed: func() {
								log.Println("Filled button pressed")
							},
						},
					},
				},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Row 2: Text and Outlined buttons
				widgets.Row{
					MainAxisAlignment: widgets.MainAxisAlignmentSpaceAround,
					Children: []widgets.Widget{
						widgets.TextButton{
							Child: widgets.Text{Data: "Text Button"},
							OnPressed: func() {
								log.Println("Text button pressed")
							},
						},

						widgets.OutlinedButton{
							Child: widgets.Text{Data: "Outlined Button"},
							OnPressed: func() {
								log.Println("Outlined button pressed")
							},
						},
					},
				},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Row 3: Icon and FAB buttons
				widgets.Row{
					MainAxisAlignment: widgets.MainAxisAlignmentSpaceAround,
					Children: []widgets.Widget{
						widgets.IconButton{
							Icon: widgets.Text{Data: "★"},
							OnPressed: func() {
								log.Println("Icon button pressed")
							},
						},

						widgets.FloatingActionButton{
							Child: widgets.Text{Data: "+"},
							OnPressed: func() {
								log.Println("FAB pressed")
							},
						},
					},
				},
			},
		},
	}
}

// InteractiveButtonSection demonstrates interactive button features
func InteractiveButtonSection() widgets.Widget {
	return widgets.Card{
		Style: "padding: 30px; border-radius: 12px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); background: white;",
		Child: widgets.Column{
			Children: []widgets.Widget{
				widgets.Text{
					Data: "Interactive Features - New Callback System",
					TextStyle: &widgets.TextStyle{
						FontSize:   &[]float64{24}[0],
						FontWeight: widgets.FontWeightBold,
						Color:      widgets.Color("#333"),
					},
				},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				widgets.Row{
					MainAxisAlignment: widgets.MainAxisAlignmentSpaceAround,
					Children: []widgets.Widget{
						// Dynamic button that changes based on state
						func() widgets.Widget {
							var button widgets.Widget
							if buttonEnabled {
								button = widgets.ElevatedButton{
									ID:    "dynamic-button",
									Child: widgets.Text{Data: "Enabled Button"},
									OnPressed: func() {
										log.Println("Enabled button clicked")
									},
								}
							} else {
								button = widgets.ElevatedButton{
									ID:    "dynamic-button",
									Child: widgets.Text{Data: "Disabled Button"},
									// OnPressed is nil, making it disabled
								}
							}
							return button
						}(),

						widgets.TextButton{
							ID:    "toggle-button",
							Child: widgets.Text{Data: "Toggle Enable/Disable"},
							OnPressed: func() {
								buttonEnabled = !buttonEnabled
								log.Printf("Button enabled: %t", buttonEnabled)
							},
						},
					},
				},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Switch widget demonstration
				widgets.Text{
					Data: "Switch Widget Demo:",
					TextStyle: &widgets.TextStyle{
						FontSize:   &[]float64{16}[0],
						FontWeight: widgets.FontWeightBold,
						Color:      widgets.Color("#666"),
					},
				},

				widgets.SizedBox{Height: &[]float64{10}[0]},

				widgets.Row{
					MainAxisAlignment: widgets.MainAxisAlignmentCenter,
					Children: []widgets.Widget{
						widgets.Text{
							Data: "Enable Features: ",
							TextStyle: &widgets.TextStyle{
								FontSize: &[]float64{16}[0],
								Color:    widgets.Color("#666"),
							},
						},

						widgets.Switch{
							ID:    "feature-switch",
							Value: buttonEnabled,
							OnChanged: func(value bool) {
								buttonEnabled = value
								log.Printf("Switch toggled, button enabled: %t", buttonEnabled)
							},
						},
					},
				},
			},
		},
	}
}

// StateManagementSection demonstrates state management integration
func StateManagementSection() widgets.Widget {
	return widgets.Card{
		Style: "padding: 30px; border-radius: 12px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); background: white;",
		Child: widgets.Column{
			Children: []widgets.Widget{
				widgets.Text{
					Data: "State Management Demo - New Callback System",
					TextStyle: &widgets.TextStyle{
						FontSize:   &[]float64{24}[0],
						FontWeight: widgets.FontWeightBold,
						Color:      widgets.Color("#333"),
					},
				},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				widgets.Text{
					Data: "This section demonstrates setState functionality with automatic UI updates. No manual HTMX endpoints needed!",
					TextStyle: &widgets.TextStyle{
						FontSize: &[]float64{14}[0],
						Color:    widgets.Color("#666"),
					},
				},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Global state display - now using simple state variable
				widgets.Container{
					ID:    "selected-type-display",
					Style: "padding: 15px; background: #e3f2fd; border-radius: 8px; text-align: center;",
					Child: widgets.Text{
						Data: fmt.Sprintf("Current Button Type: %s", selectedType),
						TextStyle: &widgets.TextStyle{
							FontWeight: widgets.FontWeightBold,
							Color:      widgets.Color("#1976d2"),
						},
					},
				},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Button type selectors with setState
				widgets.Row{
					MainAxisAlignment: widgets.MainAxisAlignmentSpaceAround,
					Children: []widgets.Widget{
						widgets.OutlinedButton{
							ID:    "primary-btn",
							Child: widgets.Text{Data: "Primary"},
							OnPressed: func() {
								selectedType = "primary"
								log.Printf("Selected type: %s", selectedType)
							},
						},

						widgets.OutlinedButton{
							ID:    "secondary-btn",
							Child: widgets.Text{Data: "Secondary"},
							OnPressed: func() {
								selectedType = "secondary"
								log.Printf("Selected type: %s", selectedType)
							},
						},

						widgets.OutlinedButton{
							ID:    "danger-btn",
							Child: widgets.Text{Data: "Danger"},
							OnPressed: func() {
								selectedType = "danger"
								log.Printf("Selected type: %s", selectedType)
							},
						},
					},
				},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Demonstration of form widgets with callbacks
				widgets.Text{
					Data: "Form Widgets Demo:",
					TextStyle: &widgets.TextStyle{
						FontSize:   &[]float64{18}[0],
						FontWeight: widgets.FontWeightBold,
						Color:      widgets.Color("#666"),
					},
				},

				widgets.SizedBox{Height: &[]float64{10}[0]},

				// Text form field with callback
				widgets.TextFormField{
					ID: "demo-form-field",
					Decoration: &widgets.InputDecoration{
						HintText: "Enter some text...",
					},
					OnChanged: func(value string) {
						log.Printf("Form field changed: %s", value)
					},
					OnFieldSubmitted: func(value string) {
						log.Printf("Form field submitted: %s", value)
					},
					Style: "margin-bottom: 10px;",
				},

				// Additional interactive widgets
				widgets.Row{
					MainAxisAlignment: widgets.MainAxisAlignmentSpaceAround,
					Children: []widgets.Widget{
						widgets.Text{
							Data: "Demo Switch: ",
							TextStyle: &widgets.TextStyle{
								FontSize: &[]float64{16}[0],
								Color:    widgets.Color("#666"),
							},
						},

						widgets.Switch{
							ID:    "demo-switch",
							Value: false,
							OnChanged: func(value bool) {
								log.Printf("Demo switch changed: %t", value)
							},
						},
					},
				},
			},
		},
	}
}

// Note: The old manual HTMX handler functions have been removed.
// All button interactions now use the new callback system with automatic
// HTMX integration through InteractiveWidget and setState functionality.
