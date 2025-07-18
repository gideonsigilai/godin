package main

import (
	"fmt"
	"log"

	"github.com/gideonsigilai/godin/pkg/core"
	"github.com/gideonsigilai/godin/pkg/widgets"
)

// App state
var counter = 0

func main() {
	app := core.New()

	// Routes
	app.GET("/", HomeHandler)
	app.GET("/about", AboutHandler)
	app.GET("/settings", SettingsHandler)
	app.POST("/increment", IncrementHandler)
	app.POST("/decrement", DecrementHandler)
	app.POST("/reset", ResetHandler)

	log.Println("Starting Counter App on :8080")
	log.Println("Visit http://localhost:8080 to see the app")
	if err := app.Serve(":8080"); err != nil {
		log.Fatal(err)
	}
}

// HomeHandler renders the main counter page
func HomeHandler(ctx *core.Context) widgets.Widget {
	return widgets.Container{
		Style: "min-height: 100vh; font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;",
		Child: widgets.Column{
			Children: []widgets.Widget{
				// App Bar
				AppBarWidget("Flutter Demo Home Page"),

				// Main content
				widgets.Expanded{
					Child: widgets.Container{
						Style: "padding: 20px;",
						Child: widgets.Column{
							MainAxisAlignment: widgets.MainAxisAlignmentCenter,
							Children: []widgets.Widget{
								// Counter display section
								widgets.Text{
									Data:      "You have pushed the button this many times:",
									TextAlign: widgets.TextAlignCenter,
									TextStyle: &widgets.TextStyle{
										FontSize: &[]float64{16}[0],
									},
								},

								// Spacer
								widgets.SizedBox{Height: &[]float64{20}[0]},

								// Counter value with ID for HTMX updates
								widgets.Container{
									ID: "counter-display",
									Child: widgets.Text{
										Data:      fmt.Sprintf("%d", counter),
										TextAlign: widgets.TextAlignCenter,
										TextStyle: &widgets.TextStyle{
											FontSize:   &[]float64{48}[0],
											FontWeight: widgets.FontWeightBold,
											Color:      widgets.Color("#2196F3"),
										},
									},
								},

								// Spacer
								widgets.SizedBox{Height: &[]float64{40}[0]},

								// Action buttons row
								widgets.Row{
									MainAxisAlignment: widgets.MainAxisAlignmentCenter,
									Children: []widgets.Widget{
										// Decrement button
										widgets.ElevatedButton{
											Child: widgets.Text{Data: "−"},
											OnPressed: func() {
												counter--
												log.Printf("Counter decremented to: %d", counter)
											},
											Style: "min-width: 60px; height: 60px; font-size: 24px; border-radius: 30px; margin-right: 10px;",
										},

										// Spacer
										widgets.SizedBox{Width: &[]float64{20}[0]},

										// Reset button
										widgets.FilledButton{
											Child: widgets.Text{Data: "Reset"},
											OnPressed: func() {
												counter = 0
												log.Printf("Counter reset to: %d", counter)
											},
											Style: "min-width: 80px; height: 60px; font-size: 16px; border-radius: 30px; background-color: #f44336;",
										},

										// Spacer
										widgets.SizedBox{Width: &[]float64{20}[0]},

										// Increment button
										widgets.ElevatedButton{
											Child: widgets.Text{Data: "+"},
											OnPressed: func() {
												counter++
												log.Printf("Counter incremented to: %d", counter)
											},
											Style: "min-width: 60px; height: 60px; font-size: 24px; border-radius: 30px; margin-left: 10px;",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// AboutHandler renders the about page
func AboutHandler(ctx *core.Context) widgets.Widget {
	return widgets.Container{
		Style: "min-height: 100vh; font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;",
		Child: widgets.Column{
			Children: []widgets.Widget{
				// App Bar
				AppBarWidget("About"),

				// Main content
				widgets.Expanded{
					Child: widgets.Container{
						Style: "padding: 40px; max-width: 800px; margin: 0 auto;",
						Child: widgets.Column{
							CrossAxisAlignment: widgets.CrossAxisAlignmentStart,
							Children: []widgets.Widget{
								// Title
								widgets.Text{
									Data: "About Counter App",
									TextStyle: &widgets.TextStyle{
										FontSize:   &[]float64{32}[0],
										FontWeight: widgets.FontWeightBold,
										Color:      widgets.Color("#333"),
									},
								},

								// Spacer
								widgets.SizedBox{Height: &[]float64{30}[0]},

								// Description
								widgets.Text{
									Data: "This is a simple counter application built with the Godin Framework, inspired by the default Flutter template app.",
									TextStyle: &widgets.TextStyle{
										FontSize: &[]float64{16}[0],
										Color:    widgets.Color("#666"),
									},
								},

								// Spacer
								widgets.SizedBox{Height: &[]float64{20}[0]},

								// Features
								widgets.Text{
									Data: "Features:",
									TextStyle: &widgets.TextStyle{
										FontSize:   &[]float64{20}[0],
										FontWeight: widgets.FontWeightBold,
										Color:      widgets.Color("#333"),
									},
								},

								// Spacer
								widgets.SizedBox{Height: &[]float64{15}[0]},

								// Feature list
								FeatureItem("• Increment and decrement counter"),
								FeatureItem("• Reset counter to zero"),
								FeatureItem("• Navigation between pages"),
								FeatureItem("• Responsive design"),
								FeatureItem("• Flutter-like widget system"),

								// Spacer
								widgets.SizedBox{Height: &[]float64{40}[0]},

								// Current counter info
								widgets.Card{
									Style: "padding: 20px; border-radius: 8px; background-color: #f5f5f5;",
									Child: widgets.Column{
										Children: []widgets.Widget{
											widgets.Text{
												Data: "Current Counter Value:",
												TextStyle: &widgets.TextStyle{
													FontWeight: widgets.FontWeightBold,
													Color:      widgets.Color("#333"),
												},
											},
											widgets.SizedBox{Height: &[]float64{10}[0]},
											widgets.Text{
												Data: fmt.Sprintf("%d", counter),
												TextStyle: &widgets.TextStyle{
													FontSize:   &[]float64{24}[0],
													FontWeight: widgets.FontWeightBold,
													Color:      widgets.Color("#2196F3"),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// SettingsHandler renders the settings page
func SettingsHandler(ctx *core.Context) widgets.Widget {
	return widgets.Container{
		Style: "min-height: 100vh; font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;",
		Child: widgets.Column{
			Children: []widgets.Widget{
				// App Bar
				AppBarWidget("Settings"),

				// Main content
				widgets.Expanded{
					Child: widgets.Container{
						Style: "padding: 40px; max-width: 600px; margin: 0 auto;",
						Child: widgets.Column{
							CrossAxisAlignment: widgets.CrossAxisAlignmentStart,
							Children: []widgets.Widget{
								// Title
								widgets.Text{
									Data: "App Settings",
									TextStyle: &widgets.TextStyle{
										FontSize:   &[]float64{32}[0],
										FontWeight: widgets.FontWeightBold,
										Color:      widgets.Color("#333"),
									},
								},

								// Spacer
								widgets.SizedBox{Height: &[]float64{30}[0]},

								// Settings options
								SettingItem("Theme", "Light"),
								SettingItem("Language", "English"),
								SettingItem("Notifications", "Enabled"),
								SettingItem("Auto-save", "On"),

								// Spacer
								widgets.SizedBox{Height: &[]float64{40}[0]},

								// App info
								widgets.Card{
									Style: "padding: 20px; border-radius: 8px; background-color: #f0f8ff;",
									Child: widgets.Column{
										Children: []widgets.Widget{
											widgets.Text{
												Data: "App Information",
												TextStyle: &widgets.TextStyle{
													FontWeight: widgets.FontWeightBold,
													Color:      widgets.Color("#333"),
												},
											},
											widgets.SizedBox{Height: &[]float64{15}[0]},
											InfoRow("Version", "1.0.0"),
											InfoRow("Framework", "Godin"),
											InfoRow("Language", "Go"),
											InfoRow("Counter Value", fmt.Sprintf("%d", counter)),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Helper Components

// AppBarWidget creates a reusable app bar with navigation
func AppBarWidget(title string) widgets.Widget {
	return widgets.AppBar{
		Title: widgets.Text{
			Data: title,
			TextStyle: &widgets.TextStyle{
				Color:      widgets.ColorWhite,
				FontSize:   &[]float64{20}[0],
				FontWeight: widgets.FontWeightBold,
			},
		},
		Style: "background: linear-gradient(135deg, #2196F3, #1976D2); color: white; padding: 16px 20px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);",
		Actions: []widgets.Widget{
			// Navigation menu
			widgets.Row{
				Children: []widgets.Widget{
					NavButton("Home", "/"),
					widgets.SizedBox{Width: &[]float64{10}[0]},
					NavButton("About", "/about"),
					widgets.SizedBox{Width: &[]float64{10}[0]},
					NavButton("Settings", "/settings"),
				},
			},
		},
	}
}

// NavButton creates a navigation button with proper link
func NavButton(text, href string) widgets.Widget {
	return widgets.Container{
		Style: "padding: 8px 16px; border-radius: 4px; background: rgba(255,255,255,0.1); cursor: pointer; transition: background 0.2s;",
		Child: widgets.Container{
			Style: "text-decoration: none;",
			Child: widgets.Text{
				Data: text,
				TextStyle: &widgets.TextStyle{
					Color:      widgets.ColorWhite,
					FontWeight: widgets.FontWeightW500,
				},
			},
		},
	}
}

// CreateActionButton creates a button that triggers an action
func CreateActionButton(text, buttonType, endpoint, style string) widgets.Widget {
	// For now, we'll use a simple button that could be enhanced with HTMX later
	buttonStyle := style + "; cursor: pointer; border: none; color: white; font-weight: bold;"

	// Add button type specific styling
	switch buttonType {
	case "primary":
		buttonStyle += " background: #2196F3;"
	case "secondary":
		buttonStyle += " background: #757575;"
	case "danger":
		buttonStyle += " background: #f44336;"
	default:
		buttonStyle += " background: #2196F3;"
	}

	return widgets.Container{
		Style: buttonStyle,
		Child: widgets.Text{
			Data: text,
			TextStyle: &widgets.TextStyle{
				Color:      widgets.ColorWhite,
				FontWeight: widgets.FontWeightBold,
			},
		},
	}
}

// FeatureItem creates a feature list item
func FeatureItem(text string) widgets.Widget {
	return widgets.Container{
		Style: "margin-bottom: 8px;",
		Child: widgets.Text{
			Data: text,
			TextStyle: &widgets.TextStyle{
				FontSize: &[]float64{16}[0],
				Color:    widgets.Color("#555"),
			},
		},
	}
}

// SettingItem creates a setting row
func SettingItem(label, value string) widgets.Widget {
	return widgets.Container{
		Style: "padding: 16px 0; border-bottom: 1px solid #eee;",
		Child: widgets.Row{
			MainAxisAlignment: widgets.MainAxisAlignmentSpaceBetween,
			Children: []widgets.Widget{
				widgets.Text{
					Data: label,
					TextStyle: &widgets.TextStyle{
						FontSize: &[]float64{16}[0],
						Color:    widgets.Color("#333"),
					},
				},
				widgets.Text{
					Data: value,
					TextStyle: &widgets.TextStyle{
						FontSize: &[]float64{16}[0],
						Color:    widgets.Color("#666"),
					},
				},
			},
		},
	}
}

// InfoRow creates an information row
func InfoRow(label, value string) widgets.Widget {
	return widgets.Container{
		Style: "margin-bottom: 8px;",
		Child: widgets.Row{
			MainAxisAlignment: widgets.MainAxisAlignmentSpaceBetween,
			Children: []widgets.Widget{
				widgets.Text{
					Data: label + ":",
					TextStyle: &widgets.TextStyle{
						Color: widgets.Color("#666"),
					},
				},
				widgets.Text{
					Data: value,
					TextStyle: &widgets.TextStyle{
						FontWeight: widgets.FontWeightW500,
						Color:      widgets.Color("#333"),
					},
				},
			},
		},
	}
}

// Action Handlers

// IncrementHandler increments the counter
func IncrementHandler(ctx *core.Context) widgets.Widget {
	counter++
	log.Printf("Counter incremented to: %d", counter)

	// Return the updated counter display
	return widgets.Text{
		Data:      fmt.Sprintf("%d", counter),
		TextAlign: widgets.TextAlignCenter,
		TextStyle: &widgets.TextStyle{
			FontSize:   &[]float64{48}[0],
			FontWeight: widgets.FontWeightBold,
			Color:      widgets.Color("#2196F3"),
		},
	}
}

// DecrementHandler decrements the counter
func DecrementHandler(ctx *core.Context) widgets.Widget {
	counter--
	log.Printf("Counter decremented to: %d", counter)

	// Return the updated counter display
	return widgets.Text{
		Data:      fmt.Sprintf("%d", counter),
		TextAlign: widgets.TextAlignCenter,
		TextStyle: &widgets.TextStyle{
			FontSize:   &[]float64{48}[0],
			FontWeight: widgets.FontWeightBold,
			Color:      widgets.Color("#2196F3"),
		},
	}
}

// ResetHandler resets the counter to zero
func ResetHandler(ctx *core.Context) widgets.Widget {
	counter = 0
	log.Printf("Counter reset to: %d", counter)

	// Return the updated counter display
	return widgets.Text{
		Data:      fmt.Sprintf("%d", counter),
		TextAlign: widgets.TextAlignCenter,
		TextStyle: &widgets.TextStyle{
			FontSize:   &[]float64{48}[0],
			FontWeight: widgets.FontWeightBold,
			Color:      widgets.Color("#2196F3"),
		},
	}
}
