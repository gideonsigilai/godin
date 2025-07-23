package main

import (
	"fmt"
	"log"

	"github.com/gideonsigilai/godin/pkg/core"
	"github.com/gideonsigilai/godin/pkg/widgets"
)

// App state
var (
	clickCountController = widgets.NewTextEditingController("0")
	currentThemeMode     = core.ThemeModeLight
	currentPage          = "home"
)

func main() {
	// Create app with theme injection
	app := core.New().
		WithThemeMode(core.ThemeModeLight).
		WithLightTheme(createCustomLightTheme()).
		WithDarkTheme(createCustomDarkTheme())

	// Enable WebSocket for real-time state updates
	app.WebSocket().Enable("/ws")

	// Initialize MediaQuery provider
	mediaQueryProvider := core.NewMediaQueryProvider()
	mediaQueryProvider.StartListening()
	app.SetMediaQueryProvider(mediaQueryProvider)

	// Initialize Dialog Manager
	dialogManager := widgets.NewDialogManager(nil)
	app.SetDialogManager(dialogManager)

	// Initialize Navigator
	navigator := widgets.NewNavigator(nil)
	app.SetNavigator(navigator)

	// Register routes
	app.GET("/", HomeHandler)
	app.GET("/dialogs", DialogDemoHandler)
	app.GET("/navigation", NavigationDemoHandler)
	app.GET("/responsive", ResponsiveDemoHandler)

	log.Println("🚀 Starting Godin UI Showcase on :8082")
	log.Println("📱 Visit http://localhost:8082 to see all new UI features")
	log.Println("✨ Features: Themes, Dialogs, Bottom Sheets, Navigation, MediaQuery, and more!")

	if err := app.Serve(":8082"); err != nil {
		log.Fatal(err)
	}
}

// createCustomLightTheme creates a custom light theme
func createCustomLightTheme() *core.ThemeData {
	theme := core.NewThemeData()
	theme.ColorScheme.Primary = core.NewColor(25, 118, 210, 255)     // Blue
	theme.ColorScheme.Secondary = core.NewColor(156, 39, 176, 255)   // Purple
	theme.ColorScheme.Surface = core.NewColor(255, 255, 255, 255)    // White
	theme.ColorScheme.Background = core.NewColor(248, 249, 250, 255) // Light gray
	return theme
}

// createCustomDarkTheme creates a custom dark theme
func createCustomDarkTheme() *core.ThemeData {
	theme := core.NewThemeData()
	theme.ColorScheme = core.NewDarkColorScheme()
	theme.ColorScheme.Primary = core.NewColor(144, 202, 249, 255)   // Light blue
	theme.ColorScheme.Secondary = core.NewColor(206, 147, 216, 255) // Light purple
	theme.Brightness = core.BrightnessDark
	return theme
}

// HomeHandler renders the main showcase page
func HomeHandler(ctx *core.Context) core.Widget {
	theme := ctx.Theme()
	mediaQuery := ctx.MediaQuery()

	return widgets.MediaQueryInheritedWidget{
		Data: mediaQuery,
		Child: widgets.Container{
			Color: &theme.ColorScheme.Background,
			Child: widgets.Column{
				Children: []core.Widget{
					// App Bar
					createAppBar(ctx),

					// Main Content
					widgets.Container{
						Padding: core.NewEdgeInsetsAll(20),
						Child: widgets.Column{
							Children: []core.Widget{
								// Welcome Section
								createWelcomeSection(ctx),

								widgets.SizedBox{Height: 30},

								// Theme Demo Section
								createThemeSection(ctx),

								widgets.SizedBox{Height: 30},

								// Dialog Demo Section
								createDialogSection(ctx),

								widgets.SizedBox{Height: 30},

								// Bottom Sheet Demo Section
								createBottomSheetSection(ctx),

								widgets.SizedBox{Height: 30},

								// Responsive Demo Section
								createResponsiveSection(ctx),

								widgets.SizedBox{Height: 30},

								// Navigation Demo Section
								createNavigationSection(ctx),
							},
						},
					},
				},
			},
		},
	}
}

// createAppBar creates a themed app bar
func createAppBar(ctx *core.Context) core.Widget {
	theme := ctx.Theme()

	return widgets.Container{
		Color:   &theme.ColorScheme.Primary,
		Height:  64,
		Padding: core.NewEdgeInsetsSymmetric(0, 20),
		Child: widgets.Row{
			MainAxisAlignment:  widgets.MainAxisAlignmentSpaceBetween,
			CrossAxisAlignment: widgets.CrossAxisAlignmentCenter,
			Children: []core.Widget{
				widgets.Text{
					Content: "🎨 Godin UI Showcase",
					Style: core.NewTextStyle().
						WithColor(theme.ColorScheme.OnPrimary).
						WithFontSize(24).
						WithFontWeight(core.FontWeightBold),
				},

				widgets.Row{
					Children: []core.Widget{
						// Theme toggle button
						widgets.IconButton{
							Icon: widgets.Text{
								Content: func() string {
									if currentThemeMode == core.ThemeModeLight {
										return "🌙"
									}
									return "☀️"
								}(),
								Style: core.NewTextStyle().WithFontSize(20),
							},
							OnPressed: func() {
								if currentThemeMode == core.ThemeModeLight {
									currentThemeMode = core.ThemeModeDark
									ctx.App.SetThemeMode(core.ThemeModeDark)
								} else {
									currentThemeMode = core.ThemeModeLight
									ctx.App.SetThemeMode(core.ThemeModeLight)
								}
								log.Printf("Theme switched to: %s", currentThemeMode)
							},
						},
					},
				},
			},
		},
	}
}

// createWelcomeSection creates the welcome section
func createWelcomeSection(ctx *core.Context) core.Widget {
	theme := ctx.Theme()
	mediaQuery := ctx.MediaQuery()

	return widgets.Container{
		Decoration: &widgets.BoxDecoration{
			Color:        &theme.ColorScheme.Surface,
			BorderRadius: widgets.NewBorderRadius(12),
			BoxShadow: []*widgets.BoxShadow{
				{
					Color:      theme.ColorScheme.Shadow.WithOpacity(0.1),
					Offset:     widgets.Offset{X: 0, Y: 2},
					BlurRadius: 8,
				},
			},
		},
		Padding: core.NewEdgeInsetsAll(24),
		Child: widgets.Column{
			Children: []core.Widget{
				widgets.Text{
					Content: "Welcome to Godin UI Framework!",
					Style: core.NewTextStyle().
						WithColor(theme.ColorScheme.OnSurface).
						WithFontSize(28).
						WithFontWeight(core.FontWeightBold),
				},

				widgets.SizedBox{Height: 16},

				widgets.Text{
					Content: fmt.Sprintf("Screen: %.0fx%.0f • %s • %s",
						mediaQuery.Size.Width,
						mediaQuery.Size.Height,
						mediaQuery.GetBreakpointName(),
						mediaQuery.Orientation),
					Style: core.NewTextStyle().
						WithColor(theme.ColorScheme.OnSurfaceVariant).
						WithFontSize(16),
				},

				widgets.SizedBox{Height: 16},

				widgets.Text{
					Content: "This showcase demonstrates all the new UI features including themes, dialogs, bottom sheets, navigation, responsive design, and more!",
					Style: core.NewTextStyle().
						WithColor(theme.ColorScheme.OnSurface).
						WithFontSize(16),
				},
			},
		},
	}
}

// createThemeSection creates the theme demonstration section
func createThemeSection(ctx *core.Context) core.Widget {
	theme := ctx.Theme()

	return widgets.Container{
		Decoration: &widgets.BoxDecoration{
			Color:        &theme.ColorScheme.Surface,
			BorderRadius: widgets.NewBorderRadius(12),
			BoxShadow: []*widgets.BoxShadow{
				{
					Color:      theme.ColorScheme.Shadow.WithOpacity(0.1),
					Offset:     widgets.Offset{X: 0, Y: 2},
					BlurRadius: 8,
				},
			},
		},
		Padding: core.NewEdgeInsetsAll(24),
		Child: widgets.Column{
			Children: []core.Widget{
				widgets.Text{
					Content: "🎨 Theme System",
					Style: core.NewTextStyle().
						WithColor(theme.ColorScheme.OnSurface).
						WithFontSize(20).
						WithFontWeight(core.FontWeightBold),
				},

				widgets.SizedBox{Height: 16},

				widgets.Text{
					Content: "Dynamic theme switching with Material Design colors",
					Style: core.NewTextStyle().
						WithColor(theme.ColorScheme.OnSurfaceVariant).
						WithFontSize(14),
				},

				widgets.SizedBox{Height: 20},

				// Color palette demonstration
				widgets.Row{
					MainAxisAlignment: widgets.MainAxisAlignmentSpaceAround,
					Children: []core.Widget{
						createColorSwatch("Primary", theme.ColorScheme.Primary),
						createColorSwatch("Secondary", theme.ColorScheme.Secondary),
						createColorSwatch("Surface", theme.ColorScheme.Surface),
						createColorSwatch("Error", theme.ColorScheme.Error),
					},
				},

				widgets.SizedBox{Height: 20},

				widgets.Row{
					MainAxisAlignment: widgets.MainAxisAlignmentCenter,
					Children: []core.Widget{
						widgets.ElevatedButton{
							Child: widgets.Text{
								Content: "Switch to " + func() string {
									if currentThemeMode == core.ThemeModeLight {
										return "Dark Theme"
									}
									return "Light Theme"
								}(),
							},
							OnPressed: func() {
								if currentThemeMode == core.ThemeModeLight {
									currentThemeMode = core.ThemeModeDark
									ctx.App.SetThemeMode(core.ThemeModeDark)
								} else {
									currentThemeMode = core.ThemeModeLight
									ctx.App.SetThemeMode(core.ThemeModeLight)
								}
								log.Printf("Theme switched to: %s", currentThemeMode)
							},
						},
					},
				},
			},
		},
	}
}

// createColorSwatch creates a color swatch display
func createColorSwatch(name string, color core.Color) core.Widget {
	return widgets.Column{
		Children: []core.Widget{
			widgets.Container{
				Width:  40,
				Height: 40,
				Decoration: &widgets.BoxDecoration{
					Color:        &color,
					BorderRadius: widgets.NewBorderRadius(8),
				},
			},
			widgets.SizedBox{Height: 8},
			widgets.Text{
				Content: name,
				Style:   core.NewTextStyle().WithFontSize(12),
			},
		},
	}
}

// createDialogSection creates the dialog demonstration section
func createDialogSection(ctx *core.Context) core.Widget {
	theme := ctx.Theme()

	return widgets.Container{
		Decoration: &widgets.BoxDecoration{
			Color:        &theme.ColorScheme.Surface,
			BorderRadius: widgets.NewBorderRadius(12),
			BoxShadow: []*widgets.BoxShadow{
				{
					Color:      theme.ColorScheme.Shadow.WithOpacity(0.1),
					Offset:     widgets.Offset{X: 0, Y: 2},
					BlurRadius: 8,
				},
			},
		},
		Padding: core.NewEdgeInsetsAll(24),
		Child: widgets.Column{
			Children: []core.Widget{
				widgets.Text{
					Content: "💬 Dialog System",
					Style: core.NewTextStyle().
						WithColor(theme.ColorScheme.OnSurface).
						WithFontSize(20).
						WithFontWeight(core.FontWeightBold),
				},

				widgets.SizedBox{Height: 16},

				widgets.Text{
					Content: "Modal dialogs with barrier dismissal and custom content",
					Style: core.NewTextStyle().
						WithColor(theme.ColorScheme.OnSurfaceVariant).
						WithFontSize(14),
				},

				widgets.SizedBox{Height: 20},

				widgets.Row{
					MainAxisAlignment: widgets.MainAxisAlignmentSpaceAround,
					Children: []core.Widget{
						widgets.ElevatedButton{
							Child: widgets.Text{Content: "Show Alert Dialog"},
							OnPressed: func() {
								showAlertDialog(ctx)
							},
						},

						widgets.OutlinedButton{
							Child: widgets.Text{Content: "Show Custom Dialog"},
							OnPressed: func() {
								showCustomDialog(ctx)
							},
						},
					},
				},
			},
		},
	}
}

// createBottomSheetSection creates the bottom sheet demonstration section
func createBottomSheetSection(ctx *core.Context) core.Widget {
	theme := ctx.Theme()

	return widgets.Container{
		Decoration: &widgets.BoxDecoration{
			Color:        &theme.ColorScheme.Surface,
			BorderRadius: widgets.NewBorderRadius(12),
			BoxShadow: []*widgets.BoxShadow{
				{
					Color:      theme.ColorScheme.Shadow.WithOpacity(0.1),
					Offset:     widgets.Offset{X: 0, Y: 2},
					BlurRadius: 8,
				},
			},
		},
		Padding: core.NewEdgeInsetsAll(24),
		Child: widgets.Column{
			Children: []core.Widget{
				widgets.Text{
					Content: "📱 Bottom Sheet System",
					Style: core.NewTextStyle().
						WithColor(theme.ColorScheme.OnSurface).
						WithFontSize(20).
						WithFontWeight(core.FontWeightBold),
				},

				widgets.SizedBox{Height: 16},

				widgets.Text{
					Content: "Sliding bottom panels with drag support",
					Style: core.NewTextStyle().
						WithColor(theme.ColorScheme.OnSurfaceVariant).
						WithFontSize(14),
				},

				widgets.SizedBox{Height: 20},

				widgets.Row{
					MainAxisAlignment: widgets.MainAxisAlignmentSpaceAround,
					Children: []core.Widget{
						widgets.ElevatedButton{
							Child: widgets.Text{Content: "Show Bottom Sheet"},
							OnPressed: func() {
								showBottomSheet(ctx)
							},
						},

						widgets.OutlinedButton{
							Child: widgets.Text{Content: "Show Modal Bottom Sheet"},
							OnPressed: func() {
								showModalBottomSheet(ctx)
							},
						},
					},
				},
			},
		},
	}
}

// createResponsiveSection creates the responsive design demonstration section
func createResponsiveSection(ctx *core.Context) core.Widget {
	theme := ctx.Theme()
	mediaQuery := ctx.MediaQuery()

	return widgets.Container{
		Decoration: &widgets.BoxDecoration{
			Color:        &theme.ColorScheme.Surface,
			BorderRadius: widgets.NewBorderRadius(12),
			BoxShadow: []*widgets.BoxShadow{
				{
					Color:      theme.ColorScheme.Shadow.WithOpacity(0.1),
					Offset:     widgets.Offset{X: 0, Y: 2},
					BlurRadius: 8,
				},
			},
		},
		Padding: core.NewEdgeInsetsAll(24),
		Child: widgets.Column{
			Children: []core.Widget{
				widgets.Text{
					Content: "📐 Responsive Design",
					Style: core.NewTextStyle().
						WithColor(theme.ColorScheme.OnSurface).
						WithFontSize(20).
						WithFontWeight(core.FontWeightBold),
				},

				widgets.SizedBox{Height: 16},

				widgets.Text{
					Content: "MediaQuery-based responsive layouts",
					Style: core.NewTextStyle().
						WithColor(theme.ColorScheme.OnSurfaceVariant).
						WithFontSize(14),
				},

				widgets.SizedBox{Height: 20},

				// Responsive content based on screen size
				widgets.ResponsiveBuilder{
					XS: widgets.Column{
						Children: []core.Widget{
							widgets.Text{Content: "📱 Extra Small Screen (XS)"},
							widgets.Text{Content: "Single column layout"},
						},
					},
					SM: widgets.Column{
						Children: []core.Widget{
							widgets.Text{Content: "📱 Small Screen (SM)"},
							widgets.Text{Content: "Compact layout"},
						},
					},
					MD: widgets.Row{
						Children: []core.Widget{
							widgets.Text{Content: "💻 Medium Screen (MD) - "},
							widgets.Text{Content: "Two column layout"},
						},
					},
					LG: widgets.Row{
						Children: []core.Widget{
							widgets.Text{Content: "🖥️ Large Screen (LG) - "},
							widgets.Text{Content: "Multi-column layout"},
						},
					},
					XL: widgets.Row{
						Children: []core.Widget{
							widgets.Text{Content: "🖥️ Extra Large Screen (XL) - "},
							widgets.Text{Content: "Full desktop layout"},
						},
					},
				},

				widgets.SizedBox{Height: 16},

				widgets.Container{
					Decoration: &widgets.BoxDecoration{
						Color:        &theme.ColorScheme.PrimaryContainer,
						BorderRadius: widgets.NewBorderRadius(8),
					},
					Padding: core.NewEdgeInsetsAll(12),
					Child: widgets.Text{
						Content: fmt.Sprintf("Current: %s • %.0fx%.0f • %s",
							mediaQuery.GetBreakpointName(),
							mediaQuery.Size.Width,
							mediaQuery.Size.Height,
							mediaQuery.Orientation),
						Style: core.NewTextStyle().
							WithColor(theme.ColorScheme.OnPrimaryContainer).
							WithFontWeight(core.FontWeightMedium),
					},
				},
			},
		},
	}
}

// createNavigationSection creates the navigation demonstration section
func createNavigationSection(ctx *core.Context) core.Widget {
	theme := ctx.Theme()

	return widgets.Container{
		Decoration: &widgets.BoxDecoration{
			Color:        &theme.ColorScheme.Surface,
			BorderRadius: widgets.NewBorderRadius(12),
			BoxShadow: []*widgets.BoxShadow{
				{
					Color:      theme.ColorScheme.Shadow.WithOpacity(0.1),
					Offset:     widgets.Offset{X: 0, Y: 2},
					BlurRadius: 8,
				},
			},
		},
		Padding: core.NewEdgeInsetsAll(24),
		Child: widgets.Column{
			Children: []core.Widget{
				widgets.Text{
					Content: "🧭 Navigation System",
					Style: core.NewTextStyle().
						WithColor(theme.ColorScheme.OnSurface).
						WithFontSize(20).
						WithFontWeight(core.FontWeightBold),
				},

				widgets.SizedBox{Height: 16},

				widgets.Text{
					Content: "Flutter-like navigation with push/pop operations",
					Style: core.NewTextStyle().
						WithColor(theme.ColorScheme.OnSurfaceVariant).
						WithFontSize(14),
				},

				widgets.SizedBox{Height: 20},

				widgets.Row{
					MainAxisAlignment: widgets.MainAxisAlignmentSpaceAround,
					Children: []core.Widget{
						widgets.ElevatedButton{
							Child: widgets.Text{Content: "Navigate to Page 2"},
							OnPressed: func() {
								// This would use the Navigator to push a new page
								log.Println("Navigation to Page 2 requested")
							},
						},

						widgets.OutlinedButton{
							Child: widgets.Text{Content: "Show Navigation Demo"},
							OnPressed: func() {
								// This would show a navigation demonstration
								log.Println("Navigation demo requested")
							},
						},
					},
				},
			},
		},
	}
}

// Dialog demonstration functions
func showAlertDialog(ctx *core.Context) {
	theme := ctx.Theme()

	alertDialog := widgets.NewAlertDialog()
	alertDialog.Title = widgets.Text{
		Content: "Alert Dialog",
		Style: core.NewTextStyle().
			WithColor(theme.ColorScheme.OnSurface).
			WithFontSize(20).
			WithFontWeight(core.FontWeightBold),
	}
	alertDialog.Content = widgets.Text{
		Content: "This is an alert dialog with custom styling and theme integration.",
		Style: core.NewTextStyle().
			WithColor(theme.ColorScheme.OnSurface).
			WithFontSize(16),
	}
	alertDialog.Actions = []core.Widget{
		widgets.TextButton{
			Child: widgets.Text{Content: "Cancel"},
			OnPressed: func() {
				widgets.DismissDialog(ctx, "current-dialog")
				log.Println("Dialog cancelled")
			},
		},
		widgets.ElevatedButton{
			Child: widgets.Text{Content: "OK"},
			OnPressed: func() {
				widgets.DismissDialog(ctx, "current-dialog")
				log.Println("Dialog confirmed")
			},
		},
	}

	widgets.ShowDialog(ctx, alertDialog, widgets.DialogOptions{
		BarrierDismissible: true,
	})
}

func showCustomDialog(ctx *core.Context) {
	theme := ctx.Theme()

	customDialog := widgets.NewDialog()
	customDialog.Title = widgets.Text{
		Content: "Custom Dialog",
		Style: core.NewTextStyle().
			WithColor(theme.ColorScheme.Primary).
			WithFontSize(24).
			WithFontWeight(core.FontWeightBold),
	}
	customDialog.Content = widgets.Container{
		Padding: core.NewEdgeInsetsAll(16),
		Child: widgets.Column{
			Children: []core.Widget{
				widgets.Text{
					Content: "This is a custom dialog with rich content:",
					Style:   core.NewTextStyle().WithFontSize(16),
				},
				widgets.SizedBox{Height: 16},
				widgets.TextField{
					Decoration: &widgets.InputDecoration{
						HintText: "Enter your name...",
					},
				},
				widgets.SizedBox{Height: 16},
				widgets.Switch{
					Value: false,
					OnChanged: func(value bool) {
						log.Printf("Dialog switch changed: %t", value)
					},
				},
			},
		},
	}

	widgets.ShowDialog(ctx, customDialog)
}

func showBottomSheet(ctx *core.Context) {
	theme := ctx.Theme()

	bottomSheet := widgets.NewBottomSheet()
	bottomSheet.Child = widgets.Container{
		Padding: core.NewEdgeInsetsAll(24),
		Child: widgets.Column{
			Children: []core.Widget{
				widgets.Text{
					Content: "Bottom Sheet",
					Style: core.NewTextStyle().
						WithColor(theme.ColorScheme.OnSurface).
						WithFontSize(20).
						WithFontWeight(core.FontWeightBold),
				},
				widgets.SizedBox{Height: 16},
				widgets.Text{
					Content: "This is a draggable bottom sheet. You can drag it up and down!",
					Style:   core.NewTextStyle().WithFontSize(16),
				},
				widgets.SizedBox{Height: 20},
				widgets.ElevatedButton{
					Child: widgets.Text{Content: "Close"},
					OnPressed: func() {
						widgets.DismissBottomSheet(ctx, "current-sheet")
					},
				},
			},
		},
	}

	widgets.ShowBottomSheet(ctx, bottomSheet)
}

func showModalBottomSheet(ctx *core.Context) {
	theme := ctx.Theme()

	modalBottomSheet := widgets.NewModalBottomSheet()
	modalBottomSheet.Child = widgets.Container{
		Padding: core.NewEdgeInsetsAll(24),
		Child: widgets.Column{
			Children: []core.Widget{
				widgets.Text{
					Content: "Modal Bottom Sheet",
					Style: core.NewTextStyle().
						WithColor(theme.ColorScheme.OnSurface).
						WithFontSize(20).
						WithFontWeight(core.FontWeightBold),
				},
				widgets.SizedBox{Height: 16},
				widgets.Text{
					Content: "This is a modal bottom sheet that blocks interaction with the background.",
					Style:   core.NewTextStyle().WithFontSize(16),
				},
				widgets.SizedBox{Height: 20},
				widgets.Row{
					MainAxisAlignment: widgets.MainAxisAlignmentSpaceAround,
					Children: []core.Widget{
						widgets.OutlinedButton{
							Child: widgets.Text{Content: "Cancel"},
							OnPressed: func() {
								widgets.DismissBottomSheet(ctx, "current-modal-sheet")
							},
						},
						widgets.ElevatedButton{
							Child: widgets.Text{Content: "Confirm"},
							OnPressed: func() {
								widgets.DismissBottomSheet(ctx, "current-modal-sheet")
								log.Println("Modal bottom sheet confirmed")
							},
						},
					},
				},
			},
		},
	}

	widgets.ShowModalBottomSheet(ctx, modalBottomSheet)
}

// Additional handlers for different demo pages
func DialogDemoHandler(ctx *core.Context) core.Widget {
	return widgets.Text{Content: "Dialog Demo Page - Coming Soon!"}
}

func NavigationDemoHandler(ctx *core.Context) core.Widget {
	return widgets.Text{Content: "Navigation Demo Page - Coming Soon!"}
}

func ResponsiveDemoHandler(ctx *core.Context) core.Widget {
	return widgets.Text{Content: "Responsive Demo Page - Coming Soon!"}
}
