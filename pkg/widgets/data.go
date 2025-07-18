package widgets

import (
	"fmt"
	"strings"

	"github.com/gideonsigilai/godin/pkg/core"
	"github.com/gideonsigilai/godin/pkg/renderer"
)

// ListView represents a list view widget with full Flutter properties
type ListView struct {
	ID                      string
	Style                   string
	Class                   string
	Children                []Widget                          // Child widgets
	ScrollDirection         Axis                              // Scroll direction
	Reverse                 bool                              // Reverse scroll direction
	Controller              *ScrollController                 // Scroll controller
	Primary                 *bool                             // Primary scroll view
	Physics                 ScrollPhysicsType                 // Scroll physics
	ShrinkWrap              bool                              // Shrink wrap
	Padding                 *EdgeInsetsGeometry               // Padding
	ItemExtent              *float64                          // Item extent
	PrototypeItem           Widget                            // Prototype item
	AddAutomaticKeepAlives  bool                              // Add automatic keep alives
	AddRepaintBoundaries    bool                              // Add repaint boundaries
	AddSemanticIndexes      bool                              // Add semantic indexes
	CacheExtent             *float64                          // Cache extent
	SemanticChildCount      *int                              // Semantic child count
	DragStartBehavior       DragStartBehavior                 // Drag start behavior
	KeyboardDismissBehavior ScrollViewKeyboardDismissBehavior // Keyboard dismiss behavior
	RestorationId           string                            // Restoration ID
	ClipBehavior            Clip                              // Clip behavior
}

// Render renders the list view as HTML
func (lv ListView) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(lv.ID, lv.Style, lv.Class+" godin-listview")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if lv.Style != "" {
		styles = append(styles, lv.Style)
	}

	// Base list view styles
	styles = append(styles, "display: flex")

	// Add scroll direction
	if lv.ScrollDirection == AxisHorizontal {
		styles = append(styles, "flex-direction: row")
		styles = append(styles, "overflow-x: auto")
		styles = append(styles, "overflow-y: hidden")
	} else {
		styles = append(styles, "flex-direction: column")
		styles = append(styles, "overflow-y: auto")
		styles = append(styles, "overflow-x: hidden")
	}

	// Add reverse direction
	if lv.Reverse {
		if lv.ScrollDirection == AxisHorizontal {
			styles = append(styles, "flex-direction: row-reverse")
		} else {
			styles = append(styles, "flex-direction: column-reverse")
		}
	}

	// Add padding
	if lv.Padding != nil {
		styles = append(styles, fmt.Sprintf("padding: %s", lv.Padding.ToCSSString()))
	}

	// Add shrink wrap
	if lv.ShrinkWrap {
		if lv.ScrollDirection == AxisHorizontal {
			styles = append(styles, "width: min-content")
		} else {
			styles = append(styles, "height: min-content")
		}
	} else {
		// Default to full size
		if lv.ScrollDirection == AxisHorizontal {
			styles = append(styles, "width: 100%")
		} else {
			styles = append(styles, "height: 100%")
		}
	}

	// Add clip behavior
	if lv.ClipBehavior != "" && lv.ClipBehavior != ClipNone {
		styles = append(styles, "overflow: hidden")
	}

	// Add scroll physics (simplified)
	if lv.Physics != "" {
		switch lv.Physics {
		case ScrollPhysicsNeverScrollable:
			styles = append(styles, "overflow: hidden")
		case ScrollPhysicsAlwaysScrollable:
			if lv.ScrollDirection == AxisHorizontal {
				styles = append(styles, "overflow-x: scroll")
			} else {
				styles = append(styles, "overflow-y: scroll")
			}
		}
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Render children
	var children []string
	for _, child := range lv.Children {
		if child != nil {
			// Wrap each child in a list item container if item extent is specified
			if lv.ItemExtent != nil {
				itemAttrs := map[string]string{"class": "godin-listview-item"}
				if lv.ScrollDirection == AxisHorizontal {
					itemAttrs["style"] = fmt.Sprintf("min-width: %.1fpx; max-width: %.1fpx", *lv.ItemExtent, *lv.ItemExtent)
				} else {
					itemAttrs["style"] = fmt.Sprintf("min-height: %.1fpx; max-height: %.1fpx", *lv.ItemExtent, *lv.ItemExtent)
				}
				itemHTML := htmlRenderer.RenderElement("div", itemAttrs, child.Render(ctx), false)
				children = append(children, itemHTML)
			} else {
				children = append(children, child.Render(ctx))
			}
		}
	}

	return htmlRenderer.RenderContainer("div", attrs, children)
}

// ListTile represents a list tile widget with full Flutter properties
type ListTile struct {
	ID                 string
	Style              string
	Class              string
	Leading            Widget                   // Leading widget
	Title              Widget                   // Title widget
	Subtitle           Widget                   // Subtitle widget
	Trailing           Widget                   // Trailing widget
	IsThreeLine        bool                     // Is three line
	Dense              *bool                    // Dense
	VisualDensity      *VisualDensity           // Visual density
	Shape              ShapeBorder              // Shape
	Style_             ListTileStyle            // List tile style (renamed to avoid conflict)
	SelectedColor      Color                    // Selected color
	IconColor          Color                    // Icon color
	TextColor          Color                    // Text color
	ContentPadding     *EdgeInsetsGeometry      // Content padding
	Enabled            bool                     // Enabled
	OnTap              GestureTapCallback       // On tap callback
	OnLongPress        GestureLongPressCallback // On long press callback
	MouseCursor        MouseCursor              // Mouse cursor
	Selected           bool                     // Selected
	FocusColor         Color                    // Focus color
	HoverColor         Color                    // Hover color
	SplashColor        Color                    // Splash color
	FocusNode          *FocusNode               // Focus node
	AutoFocus          bool                     // Auto focus
	TileColor          Color                    // Tile color
	SelectedTileColor  Color                    // Selected tile color
	EnableFeedback     *bool                    // Enable feedback
	HorizontalTitleGap *float64                 // Horizontal title gap
	MinVerticalPadding *float64                 // Min vertical padding
	MinLeadingWidth    *float64                 // Min leading width
}

// Render renders the list tile as HTML
func (lt ListTile) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(lt.ID, lt.Style, lt.Class+" godin-listtile")

	if lt.Selected {
		attrs["class"] += " selected"
	}

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if lt.Style != "" {
		styles = append(styles, lt.Style)
	}

	// Base list tile styles
	styles = append(styles, "display: flex")
	styles = append(styles, "align-items: center")
	styles = append(styles, "padding: 16px")

	// Add content padding
	if lt.ContentPadding != nil {
		styles = append(styles, fmt.Sprintf("padding: %s", lt.ContentPadding.ToCSSString()))
	}

	// Add tile colors
	if lt.Selected {
		if lt.SelectedTileColor != "" {
			styles = append(styles, fmt.Sprintf("background-color: %s", lt.SelectedTileColor))
		} else {
			styles = append(styles, "background-color: rgba(0, 0, 0, 0.08)")
		}
	} else if lt.TileColor != "" {
		styles = append(styles, fmt.Sprintf("background-color: %s", lt.TileColor))
	}

	// Add text color
	if lt.TextColor != "" {
		styles = append(styles, fmt.Sprintf("color: %s", lt.TextColor))
	}

	// Add dense styling
	if lt.Dense != nil && *lt.Dense {
		styles = append(styles, "padding: 8px 16px")
	}

	// Add minimum vertical padding
	if lt.MinVerticalPadding != nil {
		styles = append(styles, fmt.Sprintf("padding-top: %.1fpx", *lt.MinVerticalPadding))
		styles = append(styles, fmt.Sprintf("padding-bottom: %.1fpx", *lt.MinVerticalPadding))
	}

	// Add enabled/disabled styling
	if !lt.Enabled {
		styles = append(styles, "opacity: 0.6")
		styles = append(styles, "pointer-events: none")
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Add tap handler
	if lt.OnTap != nil && lt.Enabled {
		handlerID := ctx.RegisterHandler(func(ctx *core.Context) Widget {
			lt.OnTap()
			return nil
		})

		attrs["hx-post"] = "/handlers/" + handlerID
		attrs["hx-trigger"] = "click"
		styles = append(styles, "cursor: pointer")
	}

	// Add long press handler
	if lt.OnLongPress != nil && lt.Enabled {
		attrs["oncontextmenu"] = "handleListTileLongPress(event, this)"
	}

	// Add autofocus
	if lt.AutoFocus {
		attrs["autofocus"] = "true"
	}

	// Build content
	var content string

	// Leading widget
	if lt.Leading != nil {
		leadingAttrs := map[string]string{"class": "godin-listtile-leading"}
		if lt.MinLeadingWidth != nil {
			leadingAttrs["style"] = fmt.Sprintf("min-width: %.1fpx", *lt.MinLeadingWidth)
		}
		content += htmlRenderer.RenderElement("div", leadingAttrs, lt.Leading.Render(ctx), false)
	}

	// Main content
	mainContent := ""
	if lt.Title != nil {
		titleAttrs := map[string]string{"class": "godin-listtile-title"}
		mainContent += htmlRenderer.RenderElement("div", titleAttrs, lt.Title.Render(ctx), false)
	}
	if lt.Subtitle != nil {
		subtitleAttrs := map[string]string{"class": "godin-listtile-subtitle"}
		mainContent += htmlRenderer.RenderElement("div", subtitleAttrs, lt.Subtitle.Render(ctx), false)
	}

	if mainContent != "" {
		mainAttrs := map[string]string{"class": "godin-listtile-content"}
		if lt.HorizontalTitleGap != nil {
			mainAttrs["style"] = fmt.Sprintf("margin-left: %.1fpx; margin-right: %.1fpx", *lt.HorizontalTitleGap, *lt.HorizontalTitleGap)
		}
		content += htmlRenderer.RenderElement("div", mainAttrs, mainContent, false)
	}

	// Trailing widget
	if lt.Trailing != nil {
		trailingAttrs := map[string]string{"class": "godin-listtile-trailing"}
		content += htmlRenderer.RenderElement("div", trailingAttrs, lt.Trailing.Render(ctx), false)
	}

	return htmlRenderer.RenderElement("div", attrs, content, false)
}

// GridView represents a grid view widget with full Flutter properties
type GridView struct {
	ID                      string
	Style                   string
	Class                   string
	Children                []Widget                          // Child widgets
	ScrollDirection         Axis                              // Scroll direction
	Reverse                 bool                              // Reverse scroll direction
	Controller              *ScrollController                 // Scroll controller
	Primary                 *bool                             // Primary scroll view
	Physics                 ScrollPhysicsType                 // Scroll physics
	ShrinkWrap              bool                              // Shrink wrap
	Padding                 *EdgeInsetsGeometry               // Padding
	GridDelegate            SliverGridDelegate                // Grid delegate
	AddAutomaticKeepAlives  bool                              // Add automatic keep alives
	AddRepaintBoundaries    bool                              // Add repaint boundaries
	AddSemanticIndexes      bool                              // Add semantic indexes
	CacheExtent             *float64                          // Cache extent
	SemanticChildCount      *int                              // Semantic child count
	DragStartBehavior       DragStartBehavior                 // Drag start behavior
	ClipBehavior            Clip                              // Clip behavior
	KeyboardDismissBehavior ScrollViewKeyboardDismissBehavior // Keyboard dismiss behavior
	RestorationId           string                            // Restoration ID
}

// Render renders the grid view as HTML
func (gv GridView) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(gv.ID, gv.Style, gv.Class+" godin-gridview")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if gv.Style != "" {
		styles = append(styles, gv.Style)
	}

	// Base grid styles
	styles = append(styles, "display: grid")

	// Configure grid based on delegate
	if gv.GridDelegate != nil {
		crossAxisCount := gv.GridDelegate.GetCrossAxisCount()
		mainAxisSpacing := gv.GridDelegate.GetMainAxisSpacing()
		crossAxisSpacing := gv.GridDelegate.GetCrossAxisSpacing()
		aspectRatio := gv.GridDelegate.GetChildAspectRatio()

		// Set grid template columns
		if crossAxisCount > 0 {
			styles = append(styles, fmt.Sprintf("grid-template-columns: repeat(%d, 1fr)", crossAxisCount))
		}

		// Set grid gaps
		if mainAxisSpacing > 0 || crossAxisSpacing > 0 {
			styles = append(styles, fmt.Sprintf("gap: %.1fpx %.1fpx", mainAxisSpacing, crossAxisSpacing))
		}

		// Handle aspect ratio for grid items (will be applied to children)
		if aspectRatio > 0 {
			// Store aspect ratio for child rendering
			attrs["data-aspect-ratio"] = fmt.Sprintf("%.2f", aspectRatio)
		}
	} else {
		// Default grid configuration
		styles = append(styles, "grid-template-columns: repeat(auto-fit, minmax(200px, 1fr))")
		styles = append(styles, "gap: 16px")
	}

	// Add scroll direction
	if gv.ScrollDirection == AxisHorizontal {
		styles = append(styles, "overflow-x: auto")
		styles = append(styles, "overflow-y: hidden")
		styles = append(styles, "grid-auto-flow: column")
	} else {
		styles = append(styles, "overflow-y: auto")
		styles = append(styles, "overflow-x: hidden")
		styles = append(styles, "grid-auto-flow: row")
	}

	// Add reverse direction
	if gv.Reverse {
		if gv.ScrollDirection == AxisHorizontal {
			styles = append(styles, "direction: rtl")
		} else {
			styles = append(styles, "transform: scaleY(-1)")
		}
	}

	// Add padding
	if gv.Padding != nil {
		styles = append(styles, fmt.Sprintf("padding: %s", gv.Padding.ToCSSString()))
	}

	// Add shrink wrap
	if gv.ShrinkWrap {
		if gv.ScrollDirection == AxisHorizontal {
			styles = append(styles, "width: min-content")
		} else {
			styles = append(styles, "height: min-content")
		}
	} else {
		// Default to full size
		styles = append(styles, "width: 100%")
		styles = append(styles, "height: 100%")
	}

	// Add clip behavior
	if gv.ClipBehavior != "" && gv.ClipBehavior != ClipNone {
		styles = append(styles, "overflow: hidden")
	}

	// Add scroll physics (simplified)
	if gv.Physics != "" {
		switch gv.Physics {
		case ScrollPhysicsNeverScrollable:
			styles = append(styles, "overflow: hidden")
		case ScrollPhysicsAlwaysScrollable:
			if gv.ScrollDirection == AxisHorizontal {
				styles = append(styles, "overflow-x: scroll")
			} else {
				styles = append(styles, "overflow-y: scroll")
			}
		}
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Render children
	var children []string
	aspectRatio := 1.0
	if gv.GridDelegate != nil {
		aspectRatio = gv.GridDelegate.GetChildAspectRatio()
	}

	for _, child := range gv.Children {
		if child != nil {
			// Wrap each child in a grid item container
			itemAttrs := map[string]string{"class": "godin-gridview-item"}

			var itemStyles []string

			// Apply aspect ratio if specified
			if aspectRatio > 0 && aspectRatio != 1.0 {
				itemStyles = append(itemStyles, fmt.Sprintf("aspect-ratio: %.2f", aspectRatio))
			}

			// Handle reverse direction for individual items
			if gv.Reverse && gv.ScrollDirection != AxisHorizontal {
				itemStyles = append(itemStyles, "transform: scaleY(-1)")
			}

			if len(itemStyles) > 0 {
				itemAttrs["style"] = strings.Join(itemStyles, "; ")
			}

			itemHTML := htmlRenderer.RenderElement("div", itemAttrs, child.Render(ctx), false)
			children = append(children, itemHTML)
		}
	}

	return htmlRenderer.RenderContainer("div", attrs, children)
}

// SingleChildScrollView represents a single child scroll view widget with full Flutter properties
type SingleChildScrollView struct {
	ID                      string
	Style                   string
	Class                   string
	Child                   Widget                            // Child widget
	ScrollDirection         Axis                              // Scroll direction
	Reverse                 bool                              // Reverse scroll direction
	Padding                 *EdgeInsetsGeometry               // Padding
	Primary                 *bool                             // Primary scroll view
	Physics                 ScrollPhysicsType                 // Scroll physics
	Controller              *ScrollController                 // Scroll controller
	DragStartBehavior       DragStartBehavior                 // Drag start behavior
	ClipBehavior            Clip                              // Clip behavior
	RestorationId           string                            // Restoration ID
	KeyboardDismissBehavior ScrollViewKeyboardDismissBehavior // Keyboard dismiss behavior
}

// Render renders the single child scroll view as HTML
func (scsv SingleChildScrollView) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(scsv.ID, scsv.Style, scsv.Class+" godin-single-child-scroll-view")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if scsv.Style != "" {
		styles = append(styles, scsv.Style)
	}

	// Base scroll view styles
	styles = append(styles, "display: block")

	// Add scroll direction
	if scsv.ScrollDirection == AxisHorizontal {
		styles = append(styles, "overflow-x: auto")
		styles = append(styles, "overflow-y: hidden")
		styles = append(styles, "white-space: nowrap")
	} else {
		styles = append(styles, "overflow-y: auto")
		styles = append(styles, "overflow-x: hidden")
	}

	// Add reverse direction
	if scsv.Reverse {
		if scsv.ScrollDirection == AxisHorizontal {
			styles = append(styles, "direction: rtl")
		} else {
			styles = append(styles, "transform: scaleY(-1)")
		}
	}

	// Add padding
	if scsv.Padding != nil {
		styles = append(styles, fmt.Sprintf("padding: %s", scsv.Padding.ToCSSString()))
	}

	// Add clip behavior
	if scsv.ClipBehavior != "" && scsv.ClipBehavior != ClipNone {
		styles = append(styles, "overflow: hidden")
	}

	// Add scroll physics (simplified)
	if scsv.Physics != "" {
		switch scsv.Physics {
		case ScrollPhysicsNeverScrollable:
			styles = append(styles, "overflow: hidden")
		case ScrollPhysicsAlwaysScrollable:
			if scsv.ScrollDirection == AxisHorizontal {
				styles = append(styles, "overflow-x: scroll")
			} else {
				styles = append(styles, "overflow-y: scroll")
			}
		case ScrollPhysicsBouncingScrollable:
			// Add smooth scrolling for bouncing effect
			styles = append(styles, "scroll-behavior: smooth")
		}
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Render child content
	content := ""
	if scsv.Child != nil {
		childContent := scsv.Child.Render(ctx)

		// Handle reverse direction for child content
		if scsv.Reverse && scsv.ScrollDirection != AxisHorizontal {
			childAttrs := map[string]string{"style": "transform: scaleY(-1)"}
			content = htmlRenderer.RenderElement("div", childAttrs, childContent, false)
		} else {
			content = childContent
		}
	}

	return htmlRenderer.RenderElement("div", attrs, content, false)
}

// PageController represents a page controller
type PageController struct {
	InitialPage      int
	KeepPage         bool
	ViewportFraction float64
}

// PageView represents a page view widget with full Flutter properties
type PageView struct {
	ID                     string
	Style                  string
	Class                  string
	Children               []Widget          // Child widgets
	Controller             *PageController   // Page controller
	ScrollDirection        Axis              // Scroll direction
	Reverse                bool              // Reverse scroll direction
	Physics                ScrollPhysicsType // Scroll physics
	PageSnapping           bool              // Page snapping
	OnPageChanged          ValueChanged[int] // On page changed callback
	AllowImplicitScrolling bool              // Allow implicit scrolling
	RestorationId          string            // Restoration ID
	ClipBehavior           Clip              // Clip behavior
	DragStartBehavior      DragStartBehavior // Drag start behavior
	PadEnds                bool              // Pad ends
}

// Render renders the page view as HTML
func (pv PageView) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(pv.ID, pv.Style, pv.Class+" godin-page-view")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if pv.Style != "" {
		styles = append(styles, pv.Style)
	}

	// Base page view styles
	styles = append(styles, "position: relative")
	styles = append(styles, "overflow: hidden")
	styles = append(styles, "width: 100%")
	styles = append(styles, "height: 100%")

	// Add scroll direction
	if pv.ScrollDirection == AxisHorizontal {
		styles = append(styles, "display: flex")
		styles = append(styles, "flex-direction: row")
		styles = append(styles, "overflow-x: auto")
		styles = append(styles, "overflow-y: hidden")
		styles = append(styles, "scroll-snap-type: x mandatory")
	} else {
		styles = append(styles, "display: flex")
		styles = append(styles, "flex-direction: column")
		styles = append(styles, "overflow-y: auto")
		styles = append(styles, "overflow-x: hidden")
		styles = append(styles, "scroll-snap-type: y mandatory")
	}

	// Add page snapping
	if pv.PageSnapping {
		styles = append(styles, "scroll-behavior: smooth")
	}

	// Add reverse direction
	if pv.Reverse {
		if pv.ScrollDirection == AxisHorizontal {
			styles = append(styles, "flex-direction: row-reverse")
		} else {
			styles = append(styles, "flex-direction: column-reverse")
		}
	}

	// Add clip behavior
	if pv.ClipBehavior != "" && pv.ClipBehavior != ClipNone {
		styles = append(styles, "overflow: hidden")
	}

	// Add scroll physics (simplified)
	if pv.Physics != "" {
		switch pv.Physics {
		case ScrollPhysicsNeverScrollable:
			styles = append(styles, "overflow: hidden")
		case ScrollPhysicsAlwaysScrollable:
			// Keep default overflow behavior
		case ScrollPhysicsBouncingScrollable:
			styles = append(styles, "scroll-behavior: smooth")
		}
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Render children as pages
	var children []string
	for i, child := range pv.Children {
		pageAttrs := map[string]string{
			"class": "godin-page-view-page",
		}

		// Build page styles
		var pageStyles []string
		pageStyles = append(pageStyles, "flex: none")
		pageStyles = append(pageStyles, "scroll-snap-align: start")

		// Set viewport fraction
		viewportFraction := 1.0
		if pv.Controller != nil && pv.Controller.ViewportFraction > 0 {
			viewportFraction = pv.Controller.ViewportFraction
		}

		if pv.ScrollDirection == AxisHorizontal {
			pageStyles = append(pageStyles, fmt.Sprintf("width: %.1f%%", viewportFraction*100))
			pageStyles = append(pageStyles, "height: 100%")
		} else {
			pageStyles = append(pageStyles, "width: 100%")
			pageStyles = append(pageStyles, fmt.Sprintf("height: %.1f%%", viewportFraction*100))
		}

		if len(pageStyles) > 0 {
			pageAttrs["style"] = strings.Join(pageStyles, "; ")
		}

		// Add page change handler
		if pv.OnPageChanged != nil {
			handlerID := ctx.RegisterHandler(func(ctx *core.Context) Widget {
				pv.OnPageChanged(i)
				return nil
			})
			pageAttrs["hx-post"] = "/handlers/" + handlerID
			pageAttrs["hx-trigger"] = "intersect"
		}

		// Render child content
		childContent := ""
		if child != nil {
			childContent = child.Render(ctx)
		}

		children = append(children, htmlRenderer.RenderElement("div", pageAttrs, childContent, false))
	}

	return htmlRenderer.RenderContainer("div", attrs, children)
}

// CustomScrollView represents a custom scroll view widget with full Flutter properties
type CustomScrollView struct {
	ID                      string
	Style                   string
	Class                   string
	Slivers                 []Widget                          // Sliver widgets
	ScrollDirection         Axis                              // Scroll direction
	Reverse                 bool                              // Reverse scroll direction
	Controller              *ScrollController                 // Scroll controller
	Primary                 *bool                             // Primary scroll view
	Physics                 ScrollPhysicsType                 // Scroll physics
	ShrinkWrap              bool                              // Shrink wrap
	Center                  Key                               // Center key
	Anchor                  float64                           // Anchor
	CacheExtent             *float64                          // Cache extent
	SemanticChildCount      *int                              // Semantic child count
	DragStartBehavior       DragStartBehavior                 // Drag start behavior
	KeyboardDismissBehavior ScrollViewKeyboardDismissBehavior // Keyboard dismiss behavior
	RestorationId           string                            // Restoration ID
	ClipBehavior            Clip                              // Clip behavior
}

// Key represents a widget key
type Key interface {
	ToString() string
}

// Render renders the custom scroll view as HTML
func (csv CustomScrollView) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(csv.ID, csv.Style, csv.Class+" godin-custom-scroll-view")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if csv.Style != "" {
		styles = append(styles, csv.Style)
	}

	// Base custom scroll view styles
	styles = append(styles, "display: block")

	// Add scroll direction
	if csv.ScrollDirection == AxisHorizontal {
		styles = append(styles, "overflow-x: auto")
		styles = append(styles, "overflow-y: hidden")
		styles = append(styles, "white-space: nowrap")
	} else {
		styles = append(styles, "overflow-y: auto")
		styles = append(styles, "overflow-x: hidden")
	}

	// Add reverse direction
	if csv.Reverse {
		if csv.ScrollDirection == AxisHorizontal {
			styles = append(styles, "direction: rtl")
		} else {
			styles = append(styles, "transform: scaleY(-1)")
		}
	}

	// Add shrink wrap
	if csv.ShrinkWrap {
		if csv.ScrollDirection == AxisHorizontal {
			styles = append(styles, "width: fit-content")
		} else {
			styles = append(styles, "height: fit-content")
		}
	}

	// Add clip behavior
	if csv.ClipBehavior != "" && csv.ClipBehavior != ClipNone {
		styles = append(styles, "overflow: hidden")
	}

	// Add scroll physics (simplified)
	if csv.Physics != "" {
		switch csv.Physics {
		case ScrollPhysicsNeverScrollable:
			styles = append(styles, "overflow: hidden")
		case ScrollPhysicsAlwaysScrollable:
			if csv.ScrollDirection == AxisHorizontal {
				styles = append(styles, "overflow-x: scroll")
			} else {
				styles = append(styles, "overflow-y: scroll")
			}
		case ScrollPhysicsBouncingScrollable:
			styles = append(styles, "scroll-behavior: smooth")
		}
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Render slivers
	var children []string
	for _, sliver := range csv.Slivers {
		if sliver != nil {
			sliverContent := sliver.Render(ctx)

			// Handle reverse direction for sliver content
			if csv.Reverse && csv.ScrollDirection != AxisHorizontal {
				sliverAttrs := map[string]string{"style": "transform: scaleY(-1)"}
				sliverContent = htmlRenderer.RenderElement("div", sliverAttrs, sliverContent, false)
			}

			children = append(children, sliverContent)
		}
	}

	return htmlRenderer.RenderContainer("div", attrs, children)
}

// DataTable represents a data table widget
type DataTable struct {
	HTMXWidget
	Headers    []string
	Rows       [][]string
	Sortable   bool
	OnSort     string
	Pagination bool
	PageSize   int
}

// Render renders the data table as HTML
func (dt DataTable) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := dt.buildHTMXAttributes()
	attrs["class"] += " godin-datatable"

	if dt.Style != "" {
		attrs["style"] = dt.Style
	}

	return htmlRenderer.RenderTable(attrs, dt.Headers, dt.Rows)
}

// Card represents a card widget with full Flutter properties
type Card struct {
	ID                 string
	Style              string
	Class              string
	Child              Widget              // Child widget
	Color              Color               // Background color
	ShadowColor        Color               // Shadow color
	SurfaceTintColor   Color               // Surface tint color
	Elevation          *float64            // Elevation
	Shape              ShapeBorder         // Shape
	BorderOnForeground bool                // Border on foreground
	Margin             *EdgeInsetsGeometry // Margin
	ClipBehavior       Clip                // Clip behavior
	SemanticContainer  bool                // Semantic container
}

// Render renders the card as HTML
func (c Card) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(c.ID, c.Style, c.Class+" godin-card")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if c.Style != "" {
		styles = append(styles, c.Style)
	}

	// Base card styles
	styles = append(styles, "border-radius: 4px")
	styles = append(styles, "background-color: white")

	// Add background color
	if c.Color != "" {
		styles = append(styles, fmt.Sprintf("background-color: %s", c.Color))
	}

	// Add elevation (box shadow)
	if c.Elevation != nil {
		shadowBlur := *c.Elevation * 2
		shadowSpread := *c.Elevation * 0.5
		shadowColor := "rgba(0, 0, 0, 0.2)"
		if c.ShadowColor != "" {
			shadowColor = string(c.ShadowColor)
		}
		styles = append(styles, fmt.Sprintf("box-shadow: 0 %.1fpx %.1fpx %.1fpx %s", *c.Elevation, shadowBlur, shadowSpread, shadowColor))
	} else {
		// Default elevation
		styles = append(styles, "box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1)")
	}

	// Add margin
	if c.Margin != nil {
		styles = append(styles, fmt.Sprintf("margin: %s", c.Margin.ToCSSString()))
	}

	// Add surface tint color (simplified as overlay)
	if c.SurfaceTintColor != "" {
		styles = append(styles, fmt.Sprintf("background-image: linear-gradient(%s, %s)", c.SurfaceTintColor, c.SurfaceTintColor))
		styles = append(styles, "background-blend-mode: overlay")
	}

	// Add shape (simplified as border-radius)
	if c.Shape != nil {
		if shapeCSS := c.Shape.ToCSSString(); shapeCSS != "" {
			styles = append(styles, shapeCSS)
		}
	}

	// Add clip behavior
	if c.ClipBehavior != "" && c.ClipBehavior != ClipNone {
		styles = append(styles, "overflow: hidden")
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Add semantic attributes
	if c.SemanticContainer {
		attrs["role"] = "region"
	}

	content := ""
	if c.Child != nil {
		content = c.Child.Render(ctx)
	}

	return htmlRenderer.RenderElement("div", attrs, content, false)
}
