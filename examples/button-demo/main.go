package main

import (
	"fmt"
	"log"
	"net/http"

	"godin-framework/pkg/core"
	"godin-framework/pkg/widgets"
)

// App state for button demo
var (
	clickCount    = 0
	buttonEnabled = true
	selectedType  = "primary"
)

func main() {
	app := core.New()

	// Enable WebSocket for real-time state updates
	app.WebSocket().Enable("/ws")

	// Routes
	app.GET("/", HomeHandler)

	// HTMX endpoints - these should return HTML fragments, not full pages
	app.Router().HandleFunc("/increment", func(w http.ResponseWriter, r *http.Request) {
		ctx := core.NewContext(w, r, app)
		widget := IncrementHandler(ctx)
		if widget != nil {
			html := widget.Render(ctx)
			ctx.WriteHTML(html)
		}
	}).Methods("POST")

	app.Router().HandleFunc("/decrement", func(w http.ResponseWriter, r *http.Request) {
		ctx := core.NewContext(w, r, app)
		widget := DecrementHandler(ctx)
		if widget != nil {
			html := widget.Render(ctx)
			ctx.WriteHTML(html)
		}
	}).Methods("POST")

	app.Router().HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		ctx := core.NewContext(w, r, app)
		widget := ResetHandler(ctx)
		if widget != nil {
			html := widget.Render(ctx)
			ctx.WriteHTML(html)
		}
	}).Methods("POST")

	app.POST("/toggle-enabled", ToggleEnabledHandler)
	app.POST("/change-type", ChangeTypeHandler)

	// State API endpoints for real-time updates
	app.GET("/api/state/clickCount", ClickCountStateHandler)
	app.GET("/api/state/buttonEnabled", ButtonEnabledStateHandler)
	app.GET("/api/state/selectedType", SelectedTypeStateHandler)

	log.Println("Starting Button Demo on :8081")
	log.Println("Visit http://localhost:8081 to see the comprehensive button demo")
	if err := app.Serve(":8081"); err != nil {
		log.Fatal(err)
	}
}

// HomeHandler renders the main button demo page
func HomeHandler(ctx *core.Context) widgets.Widget {
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
								Data: "Godin Framework - Button Demo",
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
					Data: "Counter Demo",
					TextStyle: &widgets.TextStyle{
						FontSize:   &[]float64{24}[0],
						FontWeight: widgets.FontWeightBold,
						Color:      widgets.Color("#333"),
					},
				},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Counter display with state management
				widgets.Container{
					ID:    "counter-display",
					Style: "text-align: center; padding: 20px; background: #f8f9fa; border-radius: 8px; margin: 20px 0;",
					Child: widgets.Text{
						Data: fmt.Sprintf("Click Count: %d", clickCount),
						TextStyle: &widgets.TextStyle{
							FontSize:   &[]float64{36}[0],
							FontWeight: widgets.FontWeightBold,
							Color:      widgets.Color("#2196F3"),
						},
					},
				},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Counter buttons
				widgets.Row{
					MainAxisAlignment: widgets.MainAxisAlignmentCenter,
					Children: []widgets.Widget{
						widgets.Button{
							Text:     "Decrement",
							Type:     "secondary",
							HxPost:   "/decrement",
							HxTarget: "#counter-display",
							HxSwap:   "innerHTML",
							Style:    "margin-right: 10px;",
						},

						widgets.Button{
							Text:     "Reset",
							Type:     "danger",
							HxPost:   "/reset",
							HxTarget: "#counter-display",
							HxSwap:   "innerHTML",
							Style:    "margin: 0 10px;",
						},

						widgets.Button{
							Text:     "Increment",
							Type:     "primary",
							HxPost:   "/increment",
							HxTarget: "#counter-display",
							HxSwap:   "innerHTML",
							Style:    "margin-left: 10px;",
						},
					},
				},
			},
		},
	}
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
							OnClick: func() {
								log.Println("Default button clicked")
							},
						},

						widgets.Button{
							Text: "Primary",
							Type: "primary",
							OnClick: func() {
								log.Println("Primary button clicked")
							},
						},

						widgets.Button{
							Text: "Secondary",
							Type: "secondary",
							OnClick: func() {
								log.Println("Secondary button clicked")
							},
						},

						widgets.Button{
							Text: "Danger",
							Type: "danger",
							OnClick: func() {
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
					Data: "Interactive Features",
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
						// Enabled/Disabled toggle
						&widgets.Consumer{
							StateKey: "buttonEnabled",
							Builder: func(value interface{}) widgets.Widget {
								enabled := true
								if value != nil {
									enabled = value.(bool)
								}

								var button widgets.Widget
								if enabled {
									button = widgets.ElevatedButton{
										Child: widgets.Text{Data: "Enabled Button"},
										OnPressed: func() {
											log.Println("Enabled button clicked")
										},
									}
								} else {
									button = widgets.ElevatedButton{
										Child: widgets.Text{Data: "Disabled Button"},
										// OnPressed is nil, making it disabled
									}
								}
								return button
							},
						},

						widgets.TextButton{
							Child: widgets.Text{Data: "Toggle Enable/Disable"},
							OnPressed: func() {
								buttonEnabled = !buttonEnabled
								log.Printf("Button enabled: %t", buttonEnabled)
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
					Data: "State Management Demo",
					TextStyle: &widgets.TextStyle{
						FontSize:   &[]float64{24}[0],
						FontWeight: widgets.FontWeightBold,
						Color:      widgets.Color("#333"),
					},
				},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				widgets.Text{
					Data: "This section demonstrates real-time state updates via WebSocket. Changes made here will be reflected across all connected clients.",
					TextStyle: &widgets.TextStyle{
						FontSize: &[]float64{14}[0],
						Color:    widgets.Color("#666"),
					},
				},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Global state display
				&widgets.Consumer{
					StateKey: "selectedType",
					Builder: func(value interface{}) widgets.Widget {
						selectedType := "primary"
						if value != nil {
							selectedType = value.(string)
						}
						return widgets.Container{
							Style: "padding: 15px; background: #e3f2fd; border-radius: 8px; text-align: center;",
							Child: widgets.Text{
								Data: fmt.Sprintf("Current Button Type: %s", selectedType),
								TextStyle: &widgets.TextStyle{
									FontWeight: widgets.FontWeightBold,
									Color:      widgets.Color("#1976d2"),
								},
							},
						}
					},
				},

				widgets.SizedBox{Height: &[]float64{20}[0]},

				// Button type selectors
				widgets.Row{
					MainAxisAlignment: widgets.MainAxisAlignmentSpaceAround,
					Children: []widgets.Widget{
						widgets.OutlinedButton{
							Child: widgets.Text{Data: "Primary"},
							OnPressed: func() {
								selectedType = "primary"
								log.Printf("Selected type: %s", selectedType)
							},
						},

						widgets.OutlinedButton{
							Child: widgets.Text{Data: "Secondary"},
							OnPressed: func() {
								selectedType = "secondary"
								log.Printf("Selected type: %s", selectedType)
							},
						},

						widgets.OutlinedButton{
							Child: widgets.Text{Data: "Danger"},
							OnPressed: func() {
								selectedType = "danger"
								log.Printf("Selected type: %s", selectedType)
							},
						},
					},
				},
			},
		},
	}
}

// Handler functions for button actions

func IncrementHandler(ctx *core.Context) widgets.Widget {
	clickCount++
	ctx.App.State().Set("clickCount", clickCount)
	log.Printf("Counter incremented to: %d", clickCount)

	// Return updated counter display
	return widgets.Text{
		Data: fmt.Sprintf("Click Count: %d", clickCount),
		TextStyle: &widgets.TextStyle{
			FontSize:   &[]float64{36}[0],
			FontWeight: widgets.FontWeightBold,
			Color:      widgets.Color("#2196F3"),
		},
	}
}

func DecrementHandler(ctx *core.Context) widgets.Widget {
	clickCount--
	ctx.App.State().Set("clickCount", clickCount)
	log.Printf("Counter decremented to: %d", clickCount)

	// Return updated counter display
	return widgets.Text{
		Data: fmt.Sprintf("Click Count: %d", clickCount),
		TextStyle: &widgets.TextStyle{
			FontSize:   &[]float64{36}[0],
			FontWeight: widgets.FontWeightBold,
			Color:      widgets.Color("#2196F3"),
		},
	}
}

func ResetHandler(ctx *core.Context) widgets.Widget {
	clickCount = 0
	ctx.App.State().Set("clickCount", clickCount)
	log.Printf("Counter reset to: %d", clickCount)

	// Return updated counter display
	return widgets.Text{
		Data: fmt.Sprintf("Click Count: %d", clickCount),
		TextStyle: &widgets.TextStyle{
			FontSize:   &[]float64{36}[0],
			FontWeight: widgets.FontWeightBold,
			Color:      widgets.Color("#2196F3"),
		},
	}
}

func ToggleEnabledHandler(ctx *core.Context) widgets.Widget {
	buttonEnabled = !buttonEnabled
	ctx.App.State().Set("buttonEnabled", buttonEnabled)
	log.Printf("Button enabled: %t", buttonEnabled)
	return nil
}

func ChangeTypeHandler(ctx *core.Context) widgets.Widget {
	newType := ctx.FormValue("type")
	if newType != "" {
		selectedType = newType
		ctx.App.State().Set("selectedType", selectedType)
		log.Printf("Selected type changed to: %s", selectedType)
	}
	return nil
}

// State API handlers for real-time updates

func ClickCountStateHandler(ctx *core.Context) widgets.Widget {
	return widgets.Text{
		Data: fmt.Sprintf("%d", clickCount),
		TextStyle: &widgets.TextStyle{
			FontSize:   &[]float64{36}[0],
			FontWeight: widgets.FontWeightBold,
			Color:      widgets.Color("#2196F3"),
		},
	}
}

func ButtonEnabledStateHandler(ctx *core.Context) widgets.Widget {
	status := "Enabled"
	if !buttonEnabled {
		status = "Disabled"
	}
	return widgets.Text{
		Data: status,
		TextStyle: &widgets.TextStyle{
			FontWeight: widgets.FontWeightBold,
			Color:      widgets.Color("#4CAF50"),
		},
	}
}

func SelectedTypeStateHandler(ctx *core.Context) widgets.Widget {
	return widgets.Text{
		Data: fmt.Sprintf("Current: %s", selectedType),
		TextStyle: &widgets.TextStyle{
			FontWeight: widgets.FontWeightBold,
			Color:      widgets.Color("#1976d2"),
		},
	}
}
