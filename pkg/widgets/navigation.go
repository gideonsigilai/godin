package widgets

import (
	"fmt"
	"strings"

	"github.com/gideonsigilai/godin/pkg/core"
	"github.com/gideonsigilai/godin/pkg/renderer"
)

// AppBar represents an app bar widget with full Flutter properties
type AppBar struct {
	ID                        string
	Style                     string
	Class                     string
	Leading                   Widget               // Leading widget
	AutomaticallyImplyLeading bool                 // Automatically imply leading
	Title                     Widget               // Title widget
	Actions                   []Widget             // Action widgets
	FlexibleSpace             Widget               // Flexible space widget
	Bottom                    PreferredSizeWidget  // Bottom widget
	Elevation                 *float64             // Elevation
	ShadowColor               Color                // Shadow color
	SurfaceTintColor          Color                // Surface tint color
	Shape                     ShapeBorder          // Shape
	BackgroundColor           Color                // Background color
	ForegroundColor           Color                // Foreground color
	IconTheme                 *IconThemeData       // Icon theme
	ActionsIconTheme          *IconThemeData       // Actions icon theme
	Primary                   bool                 // Primary
	CenterTitle               *bool                // Center title
	ExcludeHeaderSemantics    bool                 // Exclude header semantics
	TitleSpacing              *float64             // Title spacing
	ToolbarOpacity            float64              // Toolbar opacity
	BottomOpacity             float64              // Bottom opacity
	ToolbarHeight             *float64             // Toolbar height
	LeadingWidth              *float64             // Leading width
	ToolbarTextStyle          *TextStyle           // Toolbar text style
	TitleTextStyle            *TextStyle           // Title text style
	SystemOverlayStyle        SystemUiOverlayStyle // System overlay style
	ForceMaterialTransparency bool                 // Force material transparency
	ClipBehavior              Clip                 // Clip behavior
}

// PreferredSizeWidget interface for widgets with preferred size
type PreferredSizeWidget interface {
	Widget
	GetPreferredSize() Size
}

// IconThemeData represents icon theme data
type IconThemeData struct {
	Color   Color    // Icon color
	Opacity *float64 // Icon opacity
	Size    *float64 // Icon size
}

// SystemUiOverlayStyle represents system UI overlay style
type SystemUiOverlayStyle string

const (
	SystemUiOverlayStyleLight SystemUiOverlayStyle = "light"
	SystemUiOverlayStyleDark  SystemUiOverlayStyle = "dark"
)

// Render renders the app bar as HTML
func (ab AppBar) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(ab.ID, ab.Style, ab.Class+" godin-appbar")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if ab.Style != "" {
		styles = append(styles, ab.Style)
	}

	// Base app bar styles
	styles = append(styles, "display: flex")
	styles = append(styles, "align-items: center")
	styles = append(styles, "padding: 0 16px")

	// Set default height
	toolbarHeight := 56.0
	if ab.ToolbarHeight != nil {
		toolbarHeight = *ab.ToolbarHeight
	}
	styles = append(styles, fmt.Sprintf("height: %.1fpx", toolbarHeight))

	// Add background color
	if ab.BackgroundColor != "" {
		styles = append(styles, fmt.Sprintf("background-color: %s", ab.BackgroundColor))
	} else {
		styles = append(styles, "background-color: #2196F3") // Default primary color
	}

	// Add foreground color
	if ab.ForegroundColor != "" {
		styles = append(styles, fmt.Sprintf("color: %s", ab.ForegroundColor))
	} else {
		styles = append(styles, "color: white") // Default white text
	}

	// Add elevation (box shadow)
	if ab.Elevation != nil {
		shadowBlur := *ab.Elevation * 2
		shadowColor := "rgba(0, 0, 0, 0.2)"
		if ab.ShadowColor != "" {
			shadowColor = string(ab.ShadowColor)
		}
		styles = append(styles, fmt.Sprintf("box-shadow: 0 %.1fpx %.1fpx %s", *ab.Elevation, shadowBlur, shadowColor))
	} else {
		// Default elevation
		styles = append(styles, "box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2)")
	}

	// Add surface tint color (simplified as overlay)
	if ab.SurfaceTintColor != "" {
		styles = append(styles, fmt.Sprintf("background-image: linear-gradient(%s, %s)", ab.SurfaceTintColor, ab.SurfaceTintColor))
		styles = append(styles, "background-blend-mode: overlay")
	}

	// Add shape (simplified as border-radius)
	if ab.Shape != nil {
		if shapeCSS := ab.Shape.ToCSSString(); shapeCSS != "" {
			styles = append(styles, shapeCSS)
		}
	}

	// Add toolbar opacity
	if ab.ToolbarOpacity > 0 && ab.ToolbarOpacity < 1 {
		styles = append(styles, fmt.Sprintf("opacity: %.2f", ab.ToolbarOpacity))
	}

	// Add clip behavior
	if ab.ClipBehavior != "" && ab.ClipBehavior != ClipNone {
		styles = append(styles, "overflow: hidden")
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	var content string

	// Leading widget
	if ab.Leading != nil {
		leadingAttrs := map[string]string{"class": "godin-appbar-leading"}
		if ab.LeadingWidth != nil {
			leadingAttrs["style"] = fmt.Sprintf("width: %.1fpx", *ab.LeadingWidth)
		}
		content += htmlRenderer.RenderElement("div", leadingAttrs, ab.Leading.Render(ctx), false)
	}

	// Title
	if ab.Title != nil {
		titleAttrs := map[string]string{"class": "godin-appbar-title"}

		var titleStyles []string
		titleStyles = append(titleStyles, "flex: 1")

		// Add title spacing
		if ab.TitleSpacing != nil {
			titleStyles = append(titleStyles, fmt.Sprintf("margin-left: %.1fpx", *ab.TitleSpacing))
		}

		// Center title if specified
		if ab.CenterTitle != nil && *ab.CenterTitle {
			titleStyles = append(titleStyles, "text-align: center")
		}

		// Add title text style
		if ab.TitleTextStyle != nil {
			if ab.TitleTextStyle.Color != "" {
				titleStyles = append(titleStyles, fmt.Sprintf("color: %s", ab.TitleTextStyle.Color))
			}
			if ab.TitleTextStyle.FontSize != nil {
				titleStyles = append(titleStyles, fmt.Sprintf("font-size: %.1fpx", *ab.TitleTextStyle.FontSize))
			}
			if ab.TitleTextStyle.FontWeight != "" {
				titleStyles = append(titleStyles, fmt.Sprintf("font-weight: %s", ab.TitleTextStyle.FontWeight))
			}
		}

		if len(titleStyles) > 0 {
			titleAttrs["style"] = strings.Join(titleStyles, "; ")
		}

		content += htmlRenderer.RenderElement("div", titleAttrs, ab.Title.Render(ctx), false)
	}

	// Actions
	if len(ab.Actions) > 0 {
		var actionElements []string
		for _, action := range ab.Actions {
			if action != nil {
				actionElements = append(actionElements, action.Render(ctx))
			}
		}
		actionsAttrs := map[string]string{"class": "godin-appbar-actions"}
		actionsAttrs["style"] = "display: flex; align-items: center; gap: 8px"
		content += htmlRenderer.RenderContainer("div", actionsAttrs, actionElements)
	}

	// Flexible space (if provided)
	if ab.FlexibleSpace != nil {
		flexibleAttrs := map[string]string{"class": "godin-appbar-flexible"}
		flexibleAttrs["style"] = "position: absolute; top: 0; left: 0; right: 0; bottom: 0; z-index: -1"
		content += htmlRenderer.RenderElement("div", flexibleAttrs, ab.FlexibleSpace.Render(ctx), false)
	}

	// Bottom widget (if provided)
	if ab.Bottom != nil {
		bottomAttrs := map[string]string{"class": "godin-appbar-bottom"}
		bottomAttrs["style"] = "position: absolute; bottom: 0; left: 0; right: 0"
		if ab.BottomOpacity > 0 && ab.BottomOpacity < 1 {
			bottomAttrs["style"] += fmt.Sprintf("; opacity: %.2f", ab.BottomOpacity)
		}
		content += htmlRenderer.RenderElement("div", bottomAttrs, ab.Bottom.Render(ctx), false)
	}

	return htmlRenderer.RenderElement("header", attrs, content, false)
}

// Drawer represents a drawer widget with full Flutter properties
type Drawer struct {
	ID               string
	Style            string
	Class            string
	Child            Widget      // Child widget
	BackgroundColor  Color       // Background color
	Elevation        *float64    // Elevation
	ShadowColor      Color       // Shadow color
	SurfaceTintColor Color       // Surface tint color
	Shape            ShapeBorder // Shape
	Width            *float64    // Width
	ClipBehavior     Clip        // Clip behavior
	SemanticLabel    string      // Semantic label
}

// Render renders the drawer as HTML
func (d Drawer) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(d.ID, d.Style, d.Class+" godin-drawer")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if d.Style != "" {
		styles = append(styles, d.Style)
	}

	// Base drawer styles
	styles = append(styles, "position: fixed")
	styles = append(styles, "top: 0")
	styles = append(styles, "left: 0")
	styles = append(styles, "height: 100vh")
	styles = append(styles, "z-index: 1000")
	styles = append(styles, "transform: translateX(-100%)")
	styles = append(styles, "transition: transform 0.3s ease")

	// Set width
	if d.Width != nil {
		styles = append(styles, fmt.Sprintf("width: %.1fpx", *d.Width))
	} else {
		styles = append(styles, "width: 304px") // Default Material Design drawer width
	}

	// Add background color
	if d.BackgroundColor != "" {
		styles = append(styles, fmt.Sprintf("background-color: %s", d.BackgroundColor))
	} else {
		styles = append(styles, "background-color: white")
	}

	// Add elevation (box shadow)
	if d.Elevation != nil {
		shadowBlur := *d.Elevation * 2
		shadowColor := "rgba(0, 0, 0, 0.16)"
		if d.ShadowColor != "" {
			shadowColor = string(d.ShadowColor)
		}
		styles = append(styles, fmt.Sprintf("box-shadow: %.1fpx 0 %.1fpx %s", *d.Elevation, shadowBlur, shadowColor))
	} else {
		// Default elevation
		styles = append(styles, "box-shadow: 2px 0 4px rgba(0, 0, 0, 0.16)")
	}

	// Add surface tint color (simplified as overlay)
	if d.SurfaceTintColor != "" {
		styles = append(styles, fmt.Sprintf("background-image: linear-gradient(%s, %s)", d.SurfaceTintColor, d.SurfaceTintColor))
		styles = append(styles, "background-blend-mode: overlay")
	}

	// Add shape (simplified as border-radius)
	if d.Shape != nil {
		if shapeCSS := d.Shape.ToCSSString(); shapeCSS != "" {
			styles = append(styles, shapeCSS)
		}
	}

	// Add clip behavior
	if d.ClipBehavior != "" && d.ClipBehavior != ClipNone {
		styles = append(styles, "overflow: hidden")
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Add semantic attributes
	if d.SemanticLabel != "" {
		attrs["aria-label"] = d.SemanticLabel
	}
	attrs["role"] = "navigation"

	// Render child content
	content := ""
	if d.Child != nil {
		content = d.Child.Render(ctx)
	}

	return htmlRenderer.RenderElement("aside", attrs, content, false)
}

// BottomNavigationBarItem represents an item in bottom navigation with full Flutter properties
type BottomNavigationBarItem struct {
	Icon            Widget // Icon widget
	Label           string // Label text
	ActiveIcon      Widget // Active icon widget
	BackgroundColor Color  // Background color
	Tooltip         string // Tooltip text
}

// BottomNavigationBarType enum
type BottomNavigationBarType string

const (
	BottomNavigationBarTypeFixed    BottomNavigationBarType = "fixed"
	BottomNavigationBarTypeShifting BottomNavigationBarType = "shifting"
)

// BottomNavigationBar represents a bottom navigation bar widget with full Flutter properties
type BottomNavigationBar struct {
	ID                   string
	Style                string
	Class                string
	Items                []BottomNavigationBarItem          // Navigation items
	OnTap                ValueChanged[int]                  // On tap callback
	CurrentIndex         int                                // Current selected index
	Elevation            *float64                           // Elevation
	Type                 BottomNavigationBarType            // Type
	FixedColor           Color                              // Fixed color (deprecated, use selectedItemColor)
	BackgroundColor      Color                              // Background color
	IconSize             *float64                           // Icon size
	SelectedItemColor    Color                              // Selected item color
	UnselectedItemColor  Color                              // Unselected item color
	SelectedIconTheme    *IconThemeData                     // Selected icon theme
	UnselectedIconTheme  *IconThemeData                     // Unselected icon theme
	SelectedFontSize     *float64                           // Selected font size
	UnselectedFontSize   *float64                           // Unselected font size
	SelectedLabelStyle   *TextStyle                         // Selected label style
	UnselectedLabelStyle *TextStyle                         // Unselected label style
	ShowSelectedLabels   *bool                              // Show selected labels
	ShowUnselectedLabels *bool                              // Show unselected labels
	MouseCursor          MouseCursor                        // Mouse cursor
	EnableFeedback       *bool                              // Enable feedback
	LandscapeLayout      BottomNavigationBarLandscapeLayout // Landscape layout
	UseLegacyColorScheme bool                               // Use legacy color scheme
}

// BottomNavigationBarLandscapeLayout enum
type BottomNavigationBarLandscapeLayout string

const (
	BottomNavigationBarLandscapeLayoutSpread   BottomNavigationBarLandscapeLayout = "spread"
	BottomNavigationBarLandscapeLayoutCentered BottomNavigationBarLandscapeLayout = "centered"
	BottomNavigationBarLandscapeLayoutLinear   BottomNavigationBarLandscapeLayout = "linear"
)

// Render renders the bottom navigation bar as HTML
func (bnb BottomNavigationBar) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(bnb.ID, bnb.Style, bnb.Class+" godin-bottom-nav")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if bnb.Style != "" {
		styles = append(styles, bnb.Style)
	}

	// Base bottom navigation styles
	styles = append(styles, "display: flex")
	styles = append(styles, "position: fixed")
	styles = append(styles, "bottom: 0")
	styles = append(styles, "left: 0")
	styles = append(styles, "right: 0")
	styles = append(styles, "height: 56px")

	// Add background color
	if bnb.BackgroundColor != "" {
		styles = append(styles, fmt.Sprintf("background-color: %s", bnb.BackgroundColor))
	} else {
		styles = append(styles, "background-color: white")
	}

	// Add elevation (box shadow)
	if bnb.Elevation != nil {
		shadowBlur := *bnb.Elevation * 2
		styles = append(styles, fmt.Sprintf("box-shadow: 0 %.1fpx %.1fpx rgba(0, 0, 0, 0.2)", -*bnb.Elevation, shadowBlur))
	} else {
		// Default elevation
		styles = append(styles, "box-shadow: 0 -2px 4px rgba(0, 0, 0, 0.1)")
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Render items
	var children []string
	for i, item := range bnb.Items {
		itemAttrs := map[string]string{
			"class": "godin-bottom-nav-item",
		}

		// Build item styles
		var itemStyles []string
		itemStyles = append(itemStyles, "flex: 1")
		itemStyles = append(itemStyles, "display: flex")
		itemStyles = append(itemStyles, "flex-direction: column")
		itemStyles = append(itemStyles, "align-items: center")
		itemStyles = append(itemStyles, "justify-content: center")
		itemStyles = append(itemStyles, "cursor: pointer")
		itemStyles = append(itemStyles, "padding: 8px")

		isSelected := i == bnb.CurrentIndex

		// Add background color for selected item
		if isSelected && item.BackgroundColor != "" {
			itemStyles = append(itemStyles, fmt.Sprintf("background-color: %s", item.BackgroundColor))
		}

		// Add selected/unselected colors
		if isSelected {
			if bnb.SelectedItemColor != "" {
				itemStyles = append(itemStyles, fmt.Sprintf("color: %s", bnb.SelectedItemColor))
			} else if bnb.FixedColor != "" {
				itemStyles = append(itemStyles, fmt.Sprintf("color: %s", bnb.FixedColor))
			} else {
				itemStyles = append(itemStyles, "color: #2196F3")
			}
		} else {
			if bnb.UnselectedItemColor != "" {
				itemStyles = append(itemStyles, fmt.Sprintf("color: %s", bnb.UnselectedItemColor))
			} else {
				itemStyles = append(itemStyles, "color: #757575")
			}
		}

		if len(itemStyles) > 0 {
			itemAttrs["style"] = strings.Join(itemStyles, "; ")
		}

		if isSelected {
			itemAttrs["class"] += " active"
		}

		// Add tap handler
		if bnb.OnTap != nil {
			handlerID := ctx.RegisterHandler(func(ctx *core.Context) Widget {
				bnb.OnTap(i)
				return nil
			})
			itemAttrs["hx-post"] = "/handlers/" + handlerID
			itemAttrs["hx-trigger"] = "click"
		}

		// Add tooltip
		if item.Tooltip != "" {
			itemAttrs["title"] = item.Tooltip
		}

		var itemContent string

		// Icon
		var iconWidget Widget
		if isSelected && item.ActiveIcon != nil {
			iconWidget = item.ActiveIcon
		} else if item.Icon != nil {
			iconWidget = item.Icon
		}

		if iconWidget != nil {
			iconAttrs := map[string]string{"class": "godin-bottom-nav-icon"}

			var iconStyles []string
			if bnb.IconSize != nil {
				iconStyles = append(iconStyles, fmt.Sprintf("font-size: %.1fpx", *bnb.IconSize))
			}

			if len(iconStyles) > 0 {
				iconAttrs["style"] = strings.Join(iconStyles, "; ")
			}

			itemContent += htmlRenderer.RenderElement("div", iconAttrs, iconWidget.Render(ctx), false)
		}

		// Label
		showLabel := true
		if isSelected {
			if bnb.ShowSelectedLabels != nil {
				showLabel = *bnb.ShowSelectedLabels
			}
		} else {
			if bnb.ShowUnselectedLabels != nil {
				showLabel = *bnb.ShowUnselectedLabels
			}
		}

		if item.Label != "" && showLabel {
			labelAttrs := map[string]string{"class": "godin-bottom-nav-label"}

			var labelStyles []string
			labelStyles = append(labelStyles, "margin-top: 4px")

			// Add font size
			if isSelected && bnb.SelectedFontSize != nil {
				labelStyles = append(labelStyles, fmt.Sprintf("font-size: %.1fpx", *bnb.SelectedFontSize))
			} else if !isSelected && bnb.UnselectedFontSize != nil {
				labelStyles = append(labelStyles, fmt.Sprintf("font-size: %.1fpx", *bnb.UnselectedFontSize))
			} else {
				labelStyles = append(labelStyles, "font-size: 12px")
			}

			// Add text style
			var textStyle *TextStyle
			if isSelected {
				textStyle = bnb.SelectedLabelStyle
			} else {
				textStyle = bnb.UnselectedLabelStyle
			}

			if textStyle != nil {
				if textStyle.Color != "" {
					labelStyles = append(labelStyles, fmt.Sprintf("color: %s", textStyle.Color))
				}
				if textStyle.FontWeight != "" {
					labelStyles = append(labelStyles, fmt.Sprintf("font-weight: %s", textStyle.FontWeight))
				}
			}

			if len(labelStyles) > 0 {
				labelAttrs["style"] = strings.Join(labelStyles, "; ")
			}

			itemContent += htmlRenderer.RenderElement("span", labelAttrs, item.Label, false)
		}

		children = append(children, htmlRenderer.RenderElement("div", itemAttrs, itemContent, false))
	}

	return htmlRenderer.RenderContainer("nav", attrs, children)
}

// TabController represents a tab controller
type TabController struct {
	Length       int
	InitialIndex int
	VSync        TickerProvider
}

// TickerProvider interface for animation
type TickerProvider interface {
	CreateTicker() Ticker
}

// Ticker interface for animation tickers
type Ticker interface {
	Start()
	Stop()
	IsActive() bool
}

// TabBarIndicatorSize enum
type TabBarIndicatorSize string

const (
	TabBarIndicatorSizeTab   TabBarIndicatorSize = "tab"
	TabBarIndicatorSizeLabel TabBarIndicatorSize = "label"
)

// TabBar represents a tab bar widget with full Flutter properties
type TabBar struct {
	ID                                string
	Style                             string
	Class                             string
	Tabs                              []Widget                     // Tab widgets
	Controller                        *TabController               // Tab controller
	IsScrollable                      bool                         // Is scrollable
	Padding                           *EdgeInsetsGeometry          // Padding
	IndicatorColor                    Color                        // Indicator color
	AutomaticIndicatorColorAdjustment bool                         // Automatic indicator color adjustment
	IndicatorWeight                   float64                      // Indicator weight
	IndicatorPadding                  *EdgeInsetsGeometry          // Indicator padding
	Indicator                         Decoration                   // Indicator decoration
	IndicatorSize                     TabBarIndicatorSize          // Indicator size
	LabelColor                        Color                        // Label color
	LabelStyle                        *TextStyle                   // Label style
	LabelPadding                      *EdgeInsetsGeometry          // Label padding
	UnselectedLabelColor              Color                        // Unselected label color
	UnselectedLabelStyle              *TextStyle                   // Unselected label style
	DragStartBehavior                 DragStartBehavior            // Drag start behavior
	OverlayColor                      MaterialStateProperty[Color] // Overlay color
	MouseCursor                       MouseCursor                  // Mouse cursor
	EnableFeedback                    *bool                        // Enable feedback
	OnTap                             ValueChanged[int]            // On tap callback
	Physics                           ScrollPhysicsType            // Scroll physics
	SplashFactory                     InteractiveInkFeatureFactory // Splash factory
	SplashBorderRadius                *BorderRadius                // Splash border radius
}

// Decoration interface for decorations
type Decoration interface {
	ToCSSString() string
}

// Render renders the tab bar as HTML
func (tb TabBar) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(tb.ID, tb.Style, tb.Class+" godin-tab-bar")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if tb.Style != "" {
		styles = append(styles, tb.Style)
	}

	// Base tab bar styles
	styles = append(styles, "display: flex")
	styles = append(styles, "position: relative")
	styles = append(styles, "border-bottom: 1px solid #e0e0e0")

	// Add scrollable behavior
	if tb.IsScrollable {
		styles = append(styles, "overflow-x: auto")
		styles = append(styles, "scrollbar-width: none")
		styles = append(styles, "-ms-overflow-style: none")
	}

	// Add padding
	if tb.Padding != nil {
		styles = append(styles, fmt.Sprintf("padding: %s", tb.Padding.ToCSSString()))
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Render tabs
	var children []string
	for i, tab := range tb.Tabs {
		tabAttrs := map[string]string{
			"class": "godin-tab-item",
		}

		// Build tab styles
		var tabStyles []string
		tabStyles = append(tabStyles, "flex: 1")
		tabStyles = append(tabStyles, "display: flex")
		tabStyles = append(tabStyles, "align-items: center")
		tabStyles = append(tabStyles, "justify-content: center")
		tabStyles = append(tabStyles, "padding: 12px 16px")
		tabStyles = append(tabStyles, "cursor: pointer")
		tabStyles = append(tabStyles, "position: relative")
		tabStyles = append(tabStyles, "transition: color 0.2s ease")

		// Add label color
		if tb.LabelColor != "" {
			tabStyles = append(tabStyles, fmt.Sprintf("color: %s", tb.LabelColor))
		}

		// Add label padding
		if tb.LabelPadding != nil {
			tabStyles = append(tabStyles, fmt.Sprintf("padding: %s", tb.LabelPadding.ToCSSString()))
		}

		if len(tabStyles) > 0 {
			tabAttrs["style"] = strings.Join(tabStyles, "; ")
		}

		// Add tap handler
		if tb.OnTap != nil {
			handlerID := ctx.RegisterHandler(func(ctx *core.Context) Widget {
				tb.OnTap(i)
				return nil
			})
			tabAttrs["hx-post"] = "/handlers/" + handlerID
			tabAttrs["hx-trigger"] = "click"
		}

		// Render tab content
		tabContent := ""
		if tab != nil {
			tabContent = tab.Render(ctx)
		}

		children = append(children, htmlRenderer.RenderElement("div", tabAttrs, tabContent, false))
	}

	// Add indicator
	indicatorAttrs := map[string]string{"class": "godin-tab-indicator"}
	var indicatorStyles []string
	indicatorStyles = append(indicatorStyles, "position: absolute")
	indicatorStyles = append(indicatorStyles, "bottom: 0")
	indicatorStyles = append(indicatorStyles, "left: 0")
	indicatorStyles = append(indicatorStyles, fmt.Sprintf("height: %.1fpx", tb.IndicatorWeight))
	indicatorStyles = append(indicatorStyles, "transition: all 0.3s ease")

	if tb.IndicatorColor != "" {
		indicatorStyles = append(indicatorStyles, fmt.Sprintf("background-color: %s", tb.IndicatorColor))
	} else {
		indicatorStyles = append(indicatorStyles, "background-color: #2196F3")
	}

	// Calculate indicator width based on tab count
	if len(tb.Tabs) > 0 {
		indicatorWidth := 100.0 / float64(len(tb.Tabs))
		indicatorStyles = append(indicatorStyles, fmt.Sprintf("width: %.1f%%", indicatorWidth))
	}

	indicatorAttrs["style"] = strings.Join(indicatorStyles, "; ")
	children = append(children, htmlRenderer.RenderElement("div", indicatorAttrs, "", false))

	return htmlRenderer.RenderContainer("div", attrs, children)
}

// TabBarView represents a tab bar view widget with full Flutter properties
type TabBarView struct {
	ID                string
	Style             string
	Class             string
	Children          []Widget          // Child widgets
	Controller        *TabController    // Tab controller
	Physics           ScrollPhysicsType // Scroll physics
	DragStartBehavior DragStartBehavior // Drag start behavior
	ViewportFraction  float64           // Viewport fraction
	ClipBehavior      Clip              // Clip behavior
}

// Render renders the tab bar view as HTML
func (tbv TabBarView) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(tbv.ID, tbv.Style, tbv.Class+" godin-tab-bar-view")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if tbv.Style != "" {
		styles = append(styles, tbv.Style)
	}

	// Base tab bar view styles
	styles = append(styles, "position: relative")
	styles = append(styles, "overflow: hidden")
	styles = append(styles, "flex: 1")

	// Add clip behavior
	if tbv.ClipBehavior != "" && tbv.ClipBehavior != ClipNone {
		styles = append(styles, "overflow: hidden")
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Render children as tab panels
	var children []string
	for i, child := range tbv.Children {
		panelAttrs := map[string]string{
			"class": "godin-tab-panel",
		}

		// Build panel styles
		var panelStyles []string
		panelStyles = append(panelStyles, "position: absolute")
		panelStyles = append(panelStyles, "top: 0")
		panelStyles = append(panelStyles, "left: 0")
		panelStyles = append(panelStyles, "width: 100%")
		panelStyles = append(panelStyles, "height: 100%")
		panelStyles = append(panelStyles, "transition: transform 0.3s ease")

		// Only show the first panel by default (in a real implementation, this would be controlled by the TabController)
		if i == 0 {
			panelStyles = append(panelStyles, "transform: translateX(0)")
			panelStyles = append(panelStyles, "opacity: 1")
		} else {
			panelStyles = append(panelStyles, "transform: translateX(100%)")
			panelStyles = append(panelStyles, "opacity: 0")
		}

		if len(panelStyles) > 0 {
			panelAttrs["style"] = strings.Join(panelStyles, "; ")
		}

		// Render child content
		childContent := ""
		if child != nil {
			childContent = child.Render(ctx)
		}

		children = append(children, htmlRenderer.RenderElement("div", panelAttrs, childContent, false))
	}

	return htmlRenderer.RenderContainer("div", attrs, children)
}
