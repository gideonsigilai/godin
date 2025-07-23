package widgets

import (
	"fmt"
	"strings"

	"github.com/gideonsigilai/godin/pkg/core"
	"github.com/gideonsigilai/godin/pkg/renderer"
)

// Text represents a text widget with full Flutter properties
type Text struct {
	ID                 string
	Style              string
	Class              string
	Data               string              // The text content
	TextStyle          *TextStyle          // Text styling
	StrutStyle         *StrutStyle         // Strut styling
	TextAlign          TextAlign           // Text alignment
	TextDirection      TextDirection       // Text direction
	Locale             *Locale             // Locale for text
	SoftWrap           *bool               // Whether text should break at soft line breaks
	Overflow           TextOverflow        // How to handle text overflow
	TextScaleFactor    *float64            // Text scale factor
	MaxLines           *int                // Maximum number of lines
	SemanticsLabel     string              // Semantic label for accessibility
	TextWidthBasis     TextWidthBasis      // How to measure text width
	TextHeightBehavior *TextHeightBehavior // Text height behavior
}

// Render renders the text as HTML
func (t Text) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(t.ID, t.Style, t.Class+" godin-text")

	// Build inline styles from various sources
	var styles []string

	// Add custom style if provided
	if t.Style != "" {
		styles = append(styles, t.Style)
	}

	// Add TextStyle CSS if provided
	if t.TextStyle != nil {
		if textStyleCSS := t.TextStyle.ToCSSString(); textStyleCSS != "" {
			styles = append(styles, textStyleCSS)
		}
	}

	// Add text alignment
	if t.TextAlign != "" {
		styles = append(styles, fmt.Sprintf("text-align: %s", t.TextAlign))
	}

	// Add text direction
	if t.TextDirection != "" {
		styles = append(styles, fmt.Sprintf("direction: %s", t.TextDirection))
		attrs["dir"] = string(t.TextDirection)
	}

	// Handle text overflow
	if t.Overflow != "" {
		switch t.Overflow {
		case TextOverflowEllipsis:
			styles = append(styles, "text-overflow: ellipsis; overflow: hidden; white-space: nowrap")
		case TextOverflowClip:
			styles = append(styles, "overflow: hidden")
		case TextOverflowFade:
			styles = append(styles, "overflow: hidden")
		}
	}

	// Handle max lines
	if t.MaxLines != nil && *t.MaxLines > 0 {
		if *t.MaxLines == 1 {
			styles = append(styles, "white-space: nowrap; overflow: hidden; text-overflow: ellipsis")
		} else {
			styles = append(styles, fmt.Sprintf("display: -webkit-box; -webkit-line-clamp: %d; -webkit-box-orient: vertical; overflow: hidden", *t.MaxLines))
		}
	}

	// Handle soft wrap
	if t.SoftWrap != nil && !*t.SoftWrap {
		styles = append(styles, "white-space: nowrap")
	}

	// Handle text scale factor
	if t.TextScaleFactor != nil && *t.TextScaleFactor != 1.0 {
		styles = append(styles, fmt.Sprintf("transform: scale(%.2f)", *t.TextScaleFactor))
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Add semantic label for accessibility
	if t.SemanticsLabel != "" {
		attrs["aria-label"] = t.SemanticsLabel
	}

	// Add locale if specified
	if t.Locale != nil {
		attrs["lang"] = t.Locale.LanguageCode
	}

	// Use Data as the text content, fallback to empty string
	content := t.Data
	if content == "" {
		content = ""
	}

	return htmlRenderer.RenderElement("span", attrs, content, false)
}

// RichText represents a rich text widget with HTML content
type RichText struct {
	ID    string
	Style string
	Class string
	HTML  string
}

// Render renders the rich text as HTML
func (rt RichText) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(rt.ID, rt.Style, rt.Class+" godin-rich-text")

	return htmlRenderer.RenderElement("div", attrs, rt.HTML, false)
}

// Image represents an image widget with full Flutter properties
type Image struct {
	ID                   string
	Style                string
	Class                string
	Image                ImageProvider     // Image source
	FrameBuilder         func() string     // Frame builder function
	LoadingBuilder       func() string     // Loading builder function
	ErrorBuilder         func() string     // Error builder function
	SemanticsLabel       string            // Semantic label for accessibility
	ExcludeFromSemantics bool              // Whether to exclude from semantics
	Width                *float64          // Image width
	Height               *float64          // Image height
	Color                Color             // Color to blend with image
	ColorBlendMode       BlendMode         // Blend mode for color
	Fit                  BoxFit            // How image should fit
	Alignment            AlignmentGeometry // Image alignment
	Repeat               ImageRepeat       // Image repeat behavior
	CenterSlice          *Rect             // Center slice for 9-patch
	MatchTextDirection   bool              // Whether to match text direction
	GaplessPlayback      bool              // Whether to use gapless playback
	IsAntiAlias          bool              // Whether to use anti-aliasing
	FilterQuality        FilterQuality     // Filter quality
}

// Render renders the image as HTML
func (i Image) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(i.ID, i.Style, i.Class+" godin-image")

	// Set image source from ImageProvider
	if i.Image != nil {
		attrs["src"] = i.Image.GetImageURL()
	}

	// Set alt text from semantic label or default
	if i.SemanticsLabel != "" {
		attrs["alt"] = i.SemanticsLabel
	} else {
		attrs["alt"] = ""
	}

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if i.Style != "" {
		styles = append(styles, i.Style)
	}

	// Add dimensions
	if i.Width != nil {
		styles = append(styles, fmt.Sprintf("width: %.1fpx", *i.Width))
	}
	if i.Height != nil {
		styles = append(styles, fmt.Sprintf("height: %.1fpx", *i.Height))
	}

	// Add object-fit based on BoxFit
	if i.Fit != "" {
		var objectFit string
		switch i.Fit {
		case BoxFitFill:
			objectFit = "fill"
		case BoxFitContain:
			objectFit = "contain"
		case BoxFitCover:
			objectFit = "cover"
		case BoxFitFitWidth:
			objectFit = "scale-down"
		case BoxFitFitHeight:
			objectFit = "scale-down"
		case BoxFitNone:
			objectFit = "none"
		case BoxFitScaleDown:
			objectFit = "scale-down"
		}
		if objectFit != "" {
			styles = append(styles, fmt.Sprintf("object-fit: %s", objectFit))
		}
	}

	// Add color blending if specified
	if i.Color != "" {
		styles = append(styles, "filter: hue-rotate(0deg) saturate(1) brightness(1)")
		// Note: CSS doesn't have direct color blending like Flutter, this is a simplified approach
	}

	// Add image repeat
	if i.Repeat != "" {
		styles = append(styles, fmt.Sprintf("background-repeat: %s", i.Repeat))
	}

	// Handle anti-aliasing
	if i.IsAntiAlias {
		styles = append(styles, "image-rendering: auto")
	} else {
		styles = append(styles, "image-rendering: pixelated")
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Add accessibility attributes
	if i.ExcludeFromSemantics {
		attrs["aria-hidden"] = "true"
	}

	// Add loading attribute for performance
	attrs["loading"] = "lazy"

	return htmlRenderer.RenderElement("img", attrs, "", true)
}

// Icon represents an icon widget with full Flutter properties
type Icon struct {
	ID             string
	Style          string
	Class          string
	Icon           IconData      // Icon data
	Size           *float64      // Icon size
	Color          Color         // Icon color
	SemanticsLabel string        // Semantic label for accessibility
	TextDirection  TextDirection // Text direction
	Shadows        []Shadow      // Icon shadows
}

// IconData represents icon data
type IconData struct {
	CodePoint          int
	FontFamily         string
	FontPackage        string
	MatchTextDirection bool
}

// Render renders the icon as HTML
func (i Icon) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	// Build class name from icon data
	className := i.Class + " godin-icon"
	if i.Icon.FontFamily != "" {
		className += " " + i.Icon.FontFamily
	}

	attrs := buildAttributes(i.ID, i.Style, className)

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if i.Style != "" {
		styles = append(styles, i.Style)
	}

	// Add icon size
	if i.Size != nil {
		styles = append(styles, fmt.Sprintf("font-size: %.1fpx", *i.Size))
	}

	// Add icon color
	if i.Color != "" {
		styles = append(styles, fmt.Sprintf("color: %s", i.Color))
	}

	// Add text direction
	if i.TextDirection != "" {
		styles = append(styles, fmt.Sprintf("direction: %s", i.TextDirection))
		attrs["dir"] = string(i.TextDirection)
	}

	// Add shadows if specified
	if len(i.Shadows) > 0 {
		var shadowStrings []string
		for _, shadow := range i.Shadows {
			shadowStrings = append(shadowStrings, fmt.Sprintf("%.1fpx %.1fpx %.1fpx %s",
				shadow.Offset.DX, shadow.Offset.DY, shadow.BlurRadius, shadow.Color))
		}
		styles = append(styles, fmt.Sprintf("text-shadow: %s", strings.Join(shadowStrings, ", ")))
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Add semantic label for accessibility
	if i.SemanticsLabel != "" {
		attrs["aria-label"] = i.SemanticsLabel
	}

	// Use Unicode character if available, otherwise empty content for CSS-based icons
	content := ""
	if i.Icon.CodePoint > 0 {
		content = fmt.Sprintf("&#%d;", i.Icon.CodePoint)
	}

	return htmlRenderer.RenderElement("i", attrs, content, false)
}

// Divider represents a divider widget with full Flutter properties
type Divider struct {
	ID        string
	Style     string
	Class     string
	Height    *float64 // Height of the divider
	Thickness *float64 // Thickness of the line
	Indent    *float64 // Left indent
	EndIndent *float64 // Right indent
	Color     Color    // Color of the divider
}

// Render renders the divider as HTML
func (d Divider) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(d.ID, d.Style, d.Class+" godin-divider")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if d.Style != "" {
		styles = append(styles, d.Style)
	}

	// Base divider styles
	styles = append(styles, "border: none")
	styles = append(styles, "margin: 0")

	// Set height (default to 16px for spacing)
	if d.Height != nil {
		styles = append(styles, fmt.Sprintf("height: %.1fpx", *d.Height))
	} else {
		styles = append(styles, "height: 16px")
	}

	// Set thickness and color
	thickness := 1.0
	if d.Thickness != nil {
		thickness = *d.Thickness
	}

	color := "#E0E0E0"
	if d.Color != "" {
		color = string(d.Color)
	}

	styles = append(styles, fmt.Sprintf("border-top: %.1fpx solid %s", thickness, color))

	// Add indents
	if d.Indent != nil {
		styles = append(styles, fmt.Sprintf("margin-left: %.1fpx", *d.Indent))
	}
	if d.EndIndent != nil {
		styles = append(styles, fmt.Sprintf("margin-right: %.1fpx", *d.EndIndent))
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	return htmlRenderer.RenderElement("hr", attrs, "", true)
}

// VerticalDivider represents a vertical divider widget with full Flutter properties
type VerticalDivider struct {
	ID        string
	Style     string
	Class     string
	Width     *float64 // Width of the divider
	Thickness *float64 // Thickness of the line
	Indent    *float64 // Top indent
	EndIndent *float64 // Bottom indent
	Color     Color    // Color of the divider
}

// Render renders the vertical divider as HTML
func (vd VerticalDivider) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(vd.ID, vd.Style, vd.Class+" godin-vertical-divider")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if vd.Style != "" {
		styles = append(styles, vd.Style)
	}

	// Base vertical divider styles
	styles = append(styles, "border: none")
	styles = append(styles, "margin: 0")
	styles = append(styles, "display: inline-block")

	// Set width (default to 16px for spacing)
	if vd.Width != nil {
		styles = append(styles, fmt.Sprintf("width: %.1fpx", *vd.Width))
	} else {
		styles = append(styles, "width: 16px")
	}

	// Set height to fill container
	styles = append(styles, "height: 100%")

	// Set thickness and color
	thickness := 1.0
	if vd.Thickness != nil {
		thickness = *vd.Thickness
	}

	color := "#E0E0E0"
	if vd.Color != "" {
		color = string(vd.Color)
	}

	styles = append(styles, fmt.Sprintf("border-left: %.1fpx solid %s", thickness, color))

	// Add indents
	if vd.Indent != nil {
		styles = append(styles, fmt.Sprintf("margin-top: %.1fpx", *vd.Indent))
	}
	if vd.EndIndent != nil {
		styles = append(styles, fmt.Sprintf("margin-bottom: %.1fpx", *vd.EndIndent))
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	return htmlRenderer.RenderElement("div", attrs, "", false)
}

// Spacer represents a spacer widget with full Flutter properties
type Spacer struct {
	ID    string
	Style string
	Class string
	Flex  int // Flex factor
}

// Render renders the spacer as HTML
func (s Spacer) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(s.ID, s.Style, s.Class+" godin-spacer")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if s.Style != "" {
		styles = append(styles, s.Style)
	}

	// Add flex property
	if s.Flex > 0 {
		styles = append(styles, fmt.Sprintf("flex: %d", s.Flex))
	} else {
		styles = append(styles, "flex: 1")
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	return htmlRenderer.RenderElement("div", attrs, "", false)
}

// Opacity represents an opacity widget with full Flutter properties
type Opacity struct {
	ID                     string
	Style                  string
	Class                  string
	Opacity                float64 // Opacity value (0.0 to 1.0)
	Child                  Widget  // Child widget
	AlwaysIncludeSemantics bool    // Always include semantics
}

// Render renders the opacity widget as HTML
func (o Opacity) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(o.ID, o.Style, o.Class+" godin-opacity")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if o.Style != "" {
		styles = append(styles, o.Style)
	}

	// Add opacity
	styles = append(styles, fmt.Sprintf("opacity: %.2f", o.Opacity))

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Render child content
	content := ""
	if o.Child != nil {
		content = o.Child.Render(ctx)
	}

	return htmlRenderer.RenderElement("div", attrs, content, false)
}

// Visibility represents a visibility widget with full Flutter properties
type Visibility struct {
	ID                    string
	Style                 string
	Class                 string
	Child                 Widget // Child widget
	Replacement           Widget // Replacement widget when not visible
	Visible               bool   // Visibility state
	MaintainState         bool   // Maintain state when hidden
	MaintainAnimation     bool   // Maintain animation when hidden
	MaintainSize          bool   // Maintain size when hidden
	MaintainSemantics     bool   // Maintain semantics when hidden
	MaintainInteractivity bool   // Maintain interactivity when hidden
}

// Render renders the visibility widget as HTML
func (v Visibility) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	// If not visible, render replacement or nothing
	if !v.Visible {
		if v.Replacement != nil {
			return v.Replacement.Render(ctx)
		}

		// If maintaining size, render invisible placeholder
		if v.MaintainSize {
			attrs := buildAttributes(v.ID, v.Style, v.Class+" godin-visibility-hidden")

			var styles []string
			if v.Style != "" {
				styles = append(styles, v.Style)
			}
			styles = append(styles, "visibility: hidden")

			if len(styles) > 0 {
				attrs["style"] = strings.Join(styles, "; ")
			}

			content := ""
			if v.Child != nil {
				content = v.Child.Render(ctx)
			}

			return htmlRenderer.RenderElement("div", attrs, content, false)
		}

		// Return empty string if not maintaining size
		return ""
	}

	// Render visible child
	attrs := buildAttributes(v.ID, v.Style, v.Class+" godin-visibility-visible")

	var styles []string
	if v.Style != "" {
		styles = append(styles, v.Style)
	}

	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	content := ""
	if v.Child != nil {
		content = v.Child.Render(ctx)
	}

	return htmlRenderer.RenderElement("div", attrs, content, false)
}

// ClipRRect represents a clip rounded rectangle widget with full Flutter properties
type ClipRRect struct {
	ID           string
	Style        string
	Class        string
	Child        Widget        // Child widget
	BorderRadius *BorderRadius // Border radius
	ClipBehavior Clip          // Clip behavior
	Clipper      CustomClipper // Custom clipper
}

// Render renders the clip rounded rectangle widget as HTML
func (crr ClipRRect) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(crr.ID, crr.Style, crr.Class+" godin-clip-rrect")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if crr.Style != "" {
		styles = append(styles, crr.Style)
	}

	// Add overflow hidden for clipping
	styles = append(styles, "overflow: hidden")

	// Add border radius
	if crr.BorderRadius != nil {
		styles = append(styles, crr.BorderRadius.ToCSSString())
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Render child content
	content := ""
	if crr.Child != nil {
		content = crr.Child.Render(ctx)
	}

	return htmlRenderer.RenderElement("div", attrs, content, false)
}

// ClipOval represents a clip oval widget with full Flutter properties
type ClipOval struct {
	ID           string
	Style        string
	Class        string
	Child        Widget        // Child widget
	ClipBehavior Clip          // Clip behavior
	Clipper      CustomClipper // Custom clipper
}

// Render renders the clip oval widget as HTML
func (co ClipOval) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(co.ID, co.Style, co.Class+" godin-clip-oval")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if co.Style != "" {
		styles = append(styles, co.Style)
	}

	// Add overflow hidden for clipping
	styles = append(styles, "overflow: hidden")

	// Make it circular/oval
	styles = append(styles, "border-radius: 50%")

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Render child content
	content := ""
	if co.Child != nil {
		content = co.Child.Render(ctx)
	}

	return htmlRenderer.RenderElement("div", attrs, content, false)
}

// ClipPath represents a clip path widget with full Flutter properties
type ClipPath struct {
	ID           string
	Style        string
	Class        string
	Child        Widget        // Child widget
	Clipper      CustomClipper // Custom clipper
	ClipBehavior Clip          // Clip behavior
}

// CustomClipper interface for custom clipping
type CustomClipper interface {
	GetClipPath() string // Returns CSS clip-path value
}

// Render renders the clip path widget as HTML
func (cp ClipPath) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(cp.ID, cp.Style, cp.Class+" godin-clip-path")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if cp.Style != "" {
		styles = append(styles, cp.Style)
	}

	// Add overflow hidden for clipping
	styles = append(styles, "overflow: hidden")

	// Add custom clip path if clipper is provided
	if cp.Clipper != nil {
		clipPath := cp.Clipper.GetClipPath()
		if clipPath != "" {
			styles = append(styles, fmt.Sprintf("clip-path: %s", clipPath))
		}
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Render child content
	content := ""
	if cp.Child != nil {
		content = cp.Child.Render(ctx)
	}

	return htmlRenderer.RenderElement("div", attrs, content, false)
}

// CircleAvatar represents a circular avatar widget with full Flutter properties
type CircleAvatar struct {
	ID                     string
	Style                  string
	Class                  string
	Child                  Widget        // Child widget
	BackgroundColor        Color         // Background color
	BackgroundImage        ImageProvider // Background image
	OnBackgroundImageError func()        // Background image error handler
	ForegroundColor        Color         // Foreground color
	Radius                 *float64      // Avatar radius
	MinRadius              *float64      // Minimum radius
	MaxRadius              *float64      // Maximum radius
	ForegroundImage        ImageProvider // Foreground image
	OnForegroundImageError func()        // Foreground image error handler
}

// Render renders the circle avatar as HTML
func (ca CircleAvatar) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(ca.ID, ca.Style, ca.Class+" godin-circle-avatar")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if ca.Style != "" {
		styles = append(styles, ca.Style)
	}

	// Make it circular
	styles = append(styles, "border-radius: 50%")
	styles = append(styles, "display: inline-flex")
	styles = append(styles, "align-items: center")
	styles = append(styles, "justify-content: center")
	styles = append(styles, "overflow: hidden")

	// Set radius (default to 20px if not specified)
	radius := 20.0
	if ca.Radius != nil {
		radius = *ca.Radius
	}
	if ca.MinRadius != nil && radius < *ca.MinRadius {
		radius = *ca.MinRadius
	}
	if ca.MaxRadius != nil && radius > *ca.MaxRadius {
		radius = *ca.MaxRadius
	}

	styles = append(styles, fmt.Sprintf("width: %.1fpx", radius*2))
	styles = append(styles, fmt.Sprintf("height: %.1fpx", radius*2))

	// Add background color
	if ca.BackgroundColor != "" {
		styles = append(styles, fmt.Sprintf("background-color: %s", ca.BackgroundColor))
	}

	// Add background image
	if ca.BackgroundImage != nil {
		styles = append(styles, fmt.Sprintf("background-image: url(%s)", ca.BackgroundImage.GetImageURL()))
		styles = append(styles, "background-size: cover")
		styles = append(styles, "background-position: center")
	}

	// Add foreground color for child content
	if ca.ForegroundColor != "" {
		styles = append(styles, fmt.Sprintf("color: %s", ca.ForegroundColor))
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Render child content or foreground image
	var content string
	if ca.ForegroundImage != nil {
		// Create an img element for foreground image
		imgAttrs := map[string]string{
			"src":   ca.ForegroundImage.GetImageURL(),
			"style": "width: 100%; height: 100%; object-fit: cover;",
		}
		content = htmlRenderer.RenderElement("img", imgAttrs, "", true)
	} else if ca.Child != nil {
		content = ca.Child.Render(ctx)
	}

	return htmlRenderer.RenderElement("div", attrs, content, false)
}

// OverflowBarAlignment enum
type OverflowBarAlignment string

const (
	OverflowBarAlignmentStart  OverflowBarAlignment = "start"
	OverflowBarAlignmentEnd    OverflowBarAlignment = "end"
	OverflowBarAlignmentCenter OverflowBarAlignment = "center"
)

// AlertDialog represents an alert dialog widget with full Flutter properties
type AlertDialog struct {
	ID                           string
	Style                        string
	Class                        string
	Title                        Widget               // Title widget
	TitlePadding                 *EdgeInsetsGeometry  // Title padding
	TitleTextStyle               *TextStyle           // Title text style
	Content                      Widget               // Content widget
	ContentPadding               *EdgeInsetsGeometry  // Content padding
	ContentTextStyle             *TextStyle           // Content text style
	Actions                      []Widget             // Action widgets
	ActionsPadding               *EdgeInsetsGeometry  // Actions padding
	ActionsAlignment             MainAxisAlignment    // Actions alignment
	ActionsOverflowAlignment     OverflowBarAlignment // Actions overflow alignment
	ActionsOverflowDirection     VerticalDirection    // Actions overflow direction
	ActionsOverflowButtonSpacing float64              // Actions overflow button spacing
	ButtonPadding                *EdgeInsetsGeometry  // Button padding
	BackgroundColor              Color                // Background color
	Elevation                    *float64             // Elevation
	ShadowColor                  Color                // Shadow color
	SurfaceTintColor             Color                // Surface tint color
	SemanticLabel                string               // Semantic label
	InsetPadding                 *EdgeInsetsGeometry  // Inset padding
	ClipBehavior                 Clip                 // Clip behavior
	Shape                        ShapeBorder          // Shape
	Alignment                    AlignmentGeometry    // Alignment
	Scrollable                   bool                 // Scrollable
}

// Render renders the alert dialog as HTML
func (ad AlertDialog) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	// Create overlay backdrop
	overlayAttrs := map[string]string{
		"class": "godin-dialog-overlay",
		"style": "position: fixed; top: 0; left: 0; width: 100%; height: 100%; background-color: rgba(0, 0, 0, 0.5); display: flex; align-items: center; justify-content: center; z-index: 1000",
	}

	attrs := buildAttributes(ad.ID, ad.Style, ad.Class+" godin-alert-dialog")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if ad.Style != "" {
		styles = append(styles, ad.Style)
	}

	// Base dialog styles
	styles = append(styles, "background-color: white")
	styles = append(styles, "border-radius: 4px")
	styles = append(styles, "min-width: 280px")
	styles = append(styles, "max-width: 560px")
	styles = append(styles, "max-height: 80vh")
	styles = append(styles, "overflow: hidden")

	// Add background color
	if ad.BackgroundColor != "" {
		styles = append(styles, fmt.Sprintf("background-color: %s", ad.BackgroundColor))
	}

	// Add elevation (box shadow)
	if ad.Elevation != nil {
		shadowBlur := *ad.Elevation * 2
		shadowColor := "rgba(0, 0, 0, 0.2)"
		if ad.ShadowColor != "" {
			shadowColor = string(ad.ShadowColor)
		}
		styles = append(styles, fmt.Sprintf("box-shadow: 0 %.1fpx %.1fpx %s", *ad.Elevation, shadowBlur, shadowColor))
	} else {
		styles = append(styles, "box-shadow: 0 11px 15px -7px rgba(0, 0, 0, 0.2), 0 24px 38px 3px rgba(0, 0, 0, 0.14), 0 9px 46px 8px rgba(0, 0, 0, 0.12)")
	}

	// Add surface tint color
	if ad.SurfaceTintColor != "" {
		styles = append(styles, fmt.Sprintf("background-image: linear-gradient(%s, %s)", ad.SurfaceTintColor, ad.SurfaceTintColor))
		styles = append(styles, "background-blend-mode: overlay")
	}

	// Add shape
	if ad.Shape != nil {
		if shapeCSS := ad.Shape.ToCSSString(); shapeCSS != "" {
			styles = append(styles, shapeCSS)
		}
	}

	// Add clip behavior
	if ad.ClipBehavior != "" && ad.ClipBehavior != ClipNone {
		styles = append(styles, "overflow: hidden")
	}

	// Add inset padding
	if ad.InsetPadding != nil {
		styles = append(styles, fmt.Sprintf("margin: %s", ad.InsetPadding.ToCSSString()))
	} else {
		styles = append(styles, "margin: 40px")
	}

	// Add scrollable behavior
	if ad.Scrollable {
		styles = append(styles, "overflow-y: auto")
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Add semantic attributes
	if ad.SemanticLabel != "" {
		attrs["aria-label"] = ad.SemanticLabel
	}
	attrs["role"] = "dialog"
	attrs["aria-modal"] = "true"

	var content string

	// Title
	if ad.Title != nil {
		titleAttrs := map[string]string{"class": "godin-dialog-title"}

		var titleStyles []string
		if ad.TitlePadding != nil {
			titleStyles = append(titleStyles, fmt.Sprintf("padding: %s", ad.TitlePadding.ToCSSString()))
		} else {
			titleStyles = append(titleStyles, "padding: 24px 24px 20px 24px")
		}

		if ad.TitleTextStyle != nil {
			if ad.TitleTextStyle.Color != "" {
				titleStyles = append(titleStyles, fmt.Sprintf("color: %s", ad.TitleTextStyle.Color))
			}
			if ad.TitleTextStyle.FontSize != nil {
				titleStyles = append(titleStyles, fmt.Sprintf("font-size: %.1fpx", *ad.TitleTextStyle.FontSize))
			}
			if ad.TitleTextStyle.FontWeight != "" {
				titleStyles = append(titleStyles, fmt.Sprintf("font-weight: %s", ad.TitleTextStyle.FontWeight))
			}
		} else {
			titleStyles = append(titleStyles, "font-size: 20px")
			titleStyles = append(titleStyles, "font-weight: 500")
		}

		if len(titleStyles) > 0 {
			titleAttrs["style"] = strings.Join(titleStyles, "; ")
		}

		content += htmlRenderer.RenderElement("div", titleAttrs, ad.Title.Render(ctx), false)
	}

	// Content
	if ad.Content != nil {
		contentAttrs := map[string]string{"class": "godin-dialog-content"}

		var contentStyles []string
		if ad.ContentPadding != nil {
			contentStyles = append(contentStyles, fmt.Sprintf("padding: %s", ad.ContentPadding.ToCSSString()))
		} else {
			contentStyles = append(contentStyles, "padding: 0 24px 24px 24px")
		}

		if ad.ContentTextStyle != nil {
			if ad.ContentTextStyle.Color != "" {
				contentStyles = append(contentStyles, fmt.Sprintf("color: %s", ad.ContentTextStyle.Color))
			}
			if ad.ContentTextStyle.FontSize != nil {
				contentStyles = append(contentStyles, fmt.Sprintf("font-size: %.1fpx", *ad.ContentTextStyle.FontSize))
			}
		} else {
			contentStyles = append(contentStyles, "font-size: 14px")
			contentStyles = append(contentStyles, "color: rgba(0, 0, 0, 0.6)")
		}

		if len(contentStyles) > 0 {
			contentAttrs["style"] = strings.Join(contentStyles, "; ")
		}

		content += htmlRenderer.RenderElement("div", contentAttrs, ad.Content.Render(ctx), false)
	}

	// Actions
	if len(ad.Actions) > 0 {
		actionsAttrs := map[string]string{"class": "godin-dialog-actions"}

		var actionsStyles []string
		actionsStyles = append(actionsStyles, "display: flex")
		actionsStyles = append(actionsStyles, "gap: 8px")

		if ad.ActionsPadding != nil {
			actionsStyles = append(actionsStyles, fmt.Sprintf("padding: %s", ad.ActionsPadding.ToCSSString()))
		} else {
			actionsStyles = append(actionsStyles, "padding: 8px")
		}

		// Add actions alignment
		switch ad.ActionsAlignment {
		case MainAxisAlignmentStart:
			actionsStyles = append(actionsStyles, "justify-content: flex-start")
		case MainAxisAlignmentEnd:
			actionsStyles = append(actionsStyles, "justify-content: flex-end")
		case MainAxisAlignmentCenter:
			actionsStyles = append(actionsStyles, "justify-content: center")
		case MainAxisAlignmentSpaceBetween:
			actionsStyles = append(actionsStyles, "justify-content: space-between")
		case MainAxisAlignmentSpaceAround:
			actionsStyles = append(actionsStyles, "justify-content: space-around")
		case MainAxisAlignmentSpaceEvenly:
			actionsStyles = append(actionsStyles, "justify-content: space-evenly")
		default:
			actionsStyles = append(actionsStyles, "justify-content: flex-end")
		}

		if len(actionsStyles) > 0 {
			actionsAttrs["style"] = strings.Join(actionsStyles, "; ")
		}

		var actionElements []string
		for _, action := range ad.Actions {
			if action != nil {
				actionElements = append(actionElements, action.Render(ctx))
			}
		}

		content += htmlRenderer.RenderContainer("div", actionsAttrs, actionElements)
	}

	dialogElement := htmlRenderer.RenderElement("div", attrs, content, false)
	return htmlRenderer.RenderElement("div", overlayAttrs, dialogElement, false)
}

// SimpleDialog represents a simple dialog widget with full Flutter properties
type SimpleDialog struct {
	ID               string
	Style            string
	Class            string
	Title            Widget              // Title widget
	TitlePadding     *EdgeInsetsGeometry // Title padding
	TitleTextStyle   *TextStyle          // Title text style
	Children         []Widget            // Child widgets
	ContentPadding   *EdgeInsetsGeometry // Content padding
	BackgroundColor  Color               // Background color
	Elevation        *float64            // Elevation
	ShadowColor      Color               // Shadow color
	SurfaceTintColor Color               // Surface tint color
	SemanticLabel    string              // Semantic label
	InsetPadding     *EdgeInsetsGeometry // Inset padding
	ClipBehavior     Clip                // Clip behavior
	Shape            ShapeBorder         // Shape
	Alignment        AlignmentGeometry   // Alignment
}

// Render renders the simple dialog as HTML
func (sd SimpleDialog) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	// Create overlay backdrop
	overlayAttrs := map[string]string{
		"class": "godin-dialog-overlay",
		"style": "position: fixed; top: 0; left: 0; width: 100%; height: 100%; background-color: rgba(0, 0, 0, 0.5); display: flex; align-items: center; justify-content: center; z-index: 1000",
	}

	attrs := buildAttributes(sd.ID, sd.Style, sd.Class+" godin-simple-dialog")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if sd.Style != "" {
		styles = append(styles, sd.Style)
	}

	// Base dialog styles
	styles = append(styles, "background-color: white")
	styles = append(styles, "border-radius: 4px")
	styles = append(styles, "min-width: 280px")
	styles = append(styles, "max-width: 560px")
	styles = append(styles, "max-height: 80vh")
	styles = append(styles, "overflow: hidden")

	// Add background color
	if sd.BackgroundColor != "" {
		styles = append(styles, fmt.Sprintf("background-color: %s", sd.BackgroundColor))
	}

	// Add elevation (box shadow)
	if sd.Elevation != nil {
		shadowBlur := *sd.Elevation * 2
		shadowColor := "rgba(0, 0, 0, 0.2)"
		if sd.ShadowColor != "" {
			shadowColor = string(sd.ShadowColor)
		}
		styles = append(styles, fmt.Sprintf("box-shadow: 0 %.1fpx %.1fpx %s", *sd.Elevation, shadowBlur, shadowColor))
	} else {
		styles = append(styles, "box-shadow: 0 11px 15px -7px rgba(0, 0, 0, 0.2), 0 24px 38px 3px rgba(0, 0, 0, 0.14), 0 9px 46px 8px rgba(0, 0, 0, 0.12)")
	}

	// Add surface tint color
	if sd.SurfaceTintColor != "" {
		styles = append(styles, fmt.Sprintf("background-image: linear-gradient(%s, %s)", sd.SurfaceTintColor, sd.SurfaceTintColor))
		styles = append(styles, "background-blend-mode: overlay")
	}

	// Add shape
	if sd.Shape != nil {
		if shapeCSS := sd.Shape.ToCSSString(); shapeCSS != "" {
			styles = append(styles, shapeCSS)
		}
	}

	// Add clip behavior
	if sd.ClipBehavior != "" && sd.ClipBehavior != ClipNone {
		styles = append(styles, "overflow: hidden")
	}

	// Add inset padding
	if sd.InsetPadding != nil {
		styles = append(styles, fmt.Sprintf("margin: %s", sd.InsetPadding.ToCSSString()))
	} else {
		styles = append(styles, "margin: 40px")
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Add semantic attributes
	if sd.SemanticLabel != "" {
		attrs["aria-label"] = sd.SemanticLabel
	}
	attrs["role"] = "dialog"
	attrs["aria-modal"] = "true"

	var content string

	// Title
	if sd.Title != nil {
		titleAttrs := map[string]string{"class": "godin-dialog-title"}

		var titleStyles []string
		if sd.TitlePadding != nil {
			titleStyles = append(titleStyles, fmt.Sprintf("padding: %s", sd.TitlePadding.ToCSSString()))
		} else {
			titleStyles = append(titleStyles, "padding: 24px 24px 20px 24px")
		}

		if sd.TitleTextStyle != nil {
			if sd.TitleTextStyle.Color != "" {
				titleStyles = append(titleStyles, fmt.Sprintf("color: %s", sd.TitleTextStyle.Color))
			}
			if sd.TitleTextStyle.FontSize != nil {
				titleStyles = append(titleStyles, fmt.Sprintf("font-size: %.1fpx", *sd.TitleTextStyle.FontSize))
			}
			if sd.TitleTextStyle.FontWeight != "" {
				titleStyles = append(titleStyles, fmt.Sprintf("font-weight: %s", sd.TitleTextStyle.FontWeight))
			}
		} else {
			titleStyles = append(titleStyles, "font-size: 20px")
			titleStyles = append(titleStyles, "font-weight: 500")
		}

		if len(titleStyles) > 0 {
			titleAttrs["style"] = strings.Join(titleStyles, "; ")
		}

		content += htmlRenderer.RenderElement("div", titleAttrs, sd.Title.Render(ctx), false)
	}

	// Content (children)
	if len(sd.Children) > 0 {
		contentAttrs := map[string]string{"class": "godin-dialog-content"}

		var contentStyles []string
		if sd.ContentPadding != nil {
			contentStyles = append(contentStyles, fmt.Sprintf("padding: %s", sd.ContentPadding.ToCSSString()))
		} else {
			contentStyles = append(contentStyles, "padding: 0 24px 24px 24px")
		}

		if len(contentStyles) > 0 {
			contentAttrs["style"] = strings.Join(contentStyles, "; ")
		}

		var childElements []string
		for _, child := range sd.Children {
			if child != nil {
				childElements = append(childElements, child.Render(ctx))
			}
		}

		content += htmlRenderer.RenderContainer("div", contentAttrs, childElements)
	}

	dialogElement := htmlRenderer.RenderElement("div", attrs, content, false)
	return htmlRenderer.RenderElement("div", overlayAttrs, dialogElement, false)
}

// SnackBarBehavior enum
type SnackBarBehavior string

const (
	SnackBarBehaviorFixed    SnackBarBehavior = "fixed"
	SnackBarBehaviorFloating SnackBarBehavior = "floating"
)

// DismissDirection enum
type DismissDirection string

const (
	DismissDirectionVertical   DismissDirection = "vertical"
	DismissDirectionHorizontal DismissDirection = "horizontal"
	DismissDirectionEndToStart DismissDirection = "endToStart"
	DismissDirectionStartToEnd DismissDirection = "startToEnd"
	DismissDirectionUp         DismissDirection = "up"
	DismissDirectionDown       DismissDirection = "down"
	DismissDirectionNone       DismissDirection = "none"
)

// SnackBar represents a snack bar widget with full Flutter properties
type SnackBar struct {
	ID                      string
	Style                   string
	Class                   string
	Content                 Widget              // Content widget
	BackgroundColor         Color               // Background color
	Elevation               *float64            // Elevation
	Margin                  *EdgeInsetsGeometry // Margin
	Padding                 *EdgeInsetsGeometry // Padding
	Width                   *float64            // Width
	Shape                   ShapeBorder         // Shape
	HitTestBehavior         HitTestBehavior     // Hit test behavior
	Behavior                SnackBarBehavior    // Behavior
	Action                  *SnackBarAction     // Action
	ActionOverflowThreshold *float64            // Action overflow threshold
	ShowCloseIcon           bool                // Show close icon
	CloseIconColor          Color               // Close icon color
	Duration                Duration            // Duration
	Animation               Animation           // Animation
	OnVisible               VoidCallback        // On visible callback
	DismissDirection        DismissDirection    // Dismiss direction
	ClipBehavior            Clip                // Clip behavior
}

// SnackBarAction represents a snack bar action
type SnackBarAction struct {
	Label                   string       // Action label
	TextColor               Color        // Text color
	DisabledTextColor       Color        // Disabled text color
	BackgroundColor         Color        // Background color
	DisabledBackgroundColor Color        // Disabled background color
	OnPressed               VoidCallback // On pressed callback
}

// HitTestBehavior enum
type HitTestBehavior string

const (
	HitTestBehaviorDeferToChild HitTestBehavior = "deferToChild"
	HitTestBehaviorOpaque       HitTestBehavior = "opaque"
	HitTestBehaviorTranslucent  HitTestBehavior = "translucent"
)

// Render renders the snack bar as HTML
func (sb SnackBar) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(sb.ID, sb.Style, sb.Class+" godin-snack-bar")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if sb.Style != "" {
		styles = append(styles, sb.Style)
	}

	// Base snack bar styles
	styles = append(styles, "position: fixed")
	styles = append(styles, "bottom: 16px")
	styles = append(styles, "left: 16px")
	styles = append(styles, "right: 16px")
	styles = append(styles, "z-index: 1000")
	styles = append(styles, "display: flex")
	styles = append(styles, "align-items: center")
	styles = append(styles, "min-height: 48px")
	styles = append(styles, "border-radius: 4px")

	// Add behavior-specific styles
	if sb.Behavior == SnackBarBehaviorFloating {
		styles = append(styles, "max-width: 344px")
		styles = append(styles, "margin: 0 auto")
		styles = append(styles, "left: 50%")
		styles = append(styles, "transform: translateX(-50%)")
	}

	// Add background color
	if sb.BackgroundColor != "" {
		styles = append(styles, fmt.Sprintf("background-color: %s", sb.BackgroundColor))
	} else {
		styles = append(styles, "background-color: #323232")
	}

	// Add elevation (box shadow)
	if sb.Elevation != nil {
		shadowBlur := *sb.Elevation * 2
		styles = append(styles, fmt.Sprintf("box-shadow: 0 %.1fpx %.1fpx rgba(0, 0, 0, 0.2)", *sb.Elevation, shadowBlur))
	} else {
		styles = append(styles, "box-shadow: 0 3px 5px -1px rgba(0, 0, 0, 0.2), 0 6px 10px 0 rgba(0, 0, 0, 0.14), 0 1px 18px 0 rgba(0, 0, 0, 0.12)")
	}

	// Add margin
	if sb.Margin != nil {
		styles = append(styles, fmt.Sprintf("margin: %s", sb.Margin.ToCSSString()))
	}

	// Add padding
	if sb.Padding != nil {
		styles = append(styles, fmt.Sprintf("padding: %s", sb.Padding.ToCSSString()))
	} else {
		styles = append(styles, "padding: 14px 16px")
	}

	// Add width
	if sb.Width != nil {
		styles = append(styles, fmt.Sprintf("width: %.1fpx", *sb.Width))
	}

	// Add shape
	if sb.Shape != nil {
		if shapeCSS := sb.Shape.ToCSSString(); shapeCSS != "" {
			styles = append(styles, shapeCSS)
		}
	}

	// Add clip behavior
	if sb.ClipBehavior != "" && sb.ClipBehavior != ClipNone {
		styles = append(styles, "overflow: hidden")
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	var content string

	// Content
	if sb.Content != nil {
		contentAttrs := map[string]string{"class": "godin-snack-bar-content"}
		contentAttrs["style"] = "flex: 1; color: white"
		content += htmlRenderer.RenderElement("div", contentAttrs, sb.Content.Render(ctx), false)
	}

	// Action
	if sb.Action != nil {
		actionAttrs := map[string]string{
			"class": "godin-snack-bar-action",
		}

		var actionStyles []string
		actionStyles = append(actionStyles, "margin-left: 16px")
		actionStyles = append(actionStyles, "padding: 8px 16px")
		actionStyles = append(actionStyles, "border: none")
		actionStyles = append(actionStyles, "border-radius: 4px")
		actionStyles = append(actionStyles, "cursor: pointer")
		actionStyles = append(actionStyles, "background: transparent")

		if sb.Action.TextColor != "" {
			actionStyles = append(actionStyles, fmt.Sprintf("color: %s", sb.Action.TextColor))
		} else {
			actionStyles = append(actionStyles, "color: #BB86FC")
		}

		if sb.Action.BackgroundColor != "" {
			actionStyles = append(actionStyles, fmt.Sprintf("background-color: %s", sb.Action.BackgroundColor))
		}

		if len(actionStyles) > 0 {
			actionAttrs["style"] = strings.Join(actionStyles, "; ")
		}

		// Add action handler
		if sb.Action.OnPressed != nil {
			handlerID := ctx.RegisterHandler(func(ctx *core.Context) Widget {
				sb.Action.OnPressed()
				return nil
			})
			actionAttrs["hx-post"] = "/handlers/" + handlerID
			actionAttrs["hx-trigger"] = "click"
		}

		content += htmlRenderer.RenderElement("button", actionAttrs, sb.Action.Label, false)
	}

	// Close icon
	if sb.ShowCloseIcon {
		closeAttrs := map[string]string{
			"class": "godin-snack-bar-close",
		}

		var closeStyles []string
		closeStyles = append(closeStyles, "margin-left: 16px")
		closeStyles = append(closeStyles, "padding: 8px")
		closeStyles = append(closeStyles, "border: none")
		closeStyles = append(closeStyles, "background: transparent")
		closeStyles = append(closeStyles, "cursor: pointer")
		closeStyles = append(closeStyles, "border-radius: 50%")

		if sb.CloseIconColor != "" {
			closeStyles = append(closeStyles, fmt.Sprintf("color: %s", sb.CloseIconColor))
		} else {
			closeStyles = append(closeStyles, "color: white")
		}

		if len(closeStyles) > 0 {
			closeAttrs["style"] = strings.Join(closeStyles, "; ")
		}

		content += htmlRenderer.RenderElement("button", closeAttrs, "×", false)
	}

	return htmlRenderer.RenderElement("div", attrs, content, false)
}
