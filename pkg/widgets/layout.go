package widgets

import (
	"fmt"
	"godin-framework/pkg/core"
	"godin-framework/pkg/renderer"
	"math"
	"strings"
)

// BoxConstraints represents layout constraints for layout widgets
type BoxConstraints struct {
	MinWidth  *float64
	MaxWidth  *float64
	MinHeight *float64
	MaxHeight *float64
}

// Container represents a container widget with full Flutter properties
type Container struct {
	ID                   string
	Style                string
	Class                string
	Child                Widget              // Child widget
	Padding              *EdgeInsetsGeometry // Padding around child
	Margin               *EdgeInsetsGeometry // Margin around container
	Width                *float64            // Container width
	Height               *float64            // Container height
	Constraints          *BoxConstraints     // Layout constraints
	Decoration           *BoxDecoration      // Background decoration
	ForegroundDecoration *BoxDecoration      // Foreground decoration
	Transform            *Matrix4            // Transform matrix
	TransformAlignment   AlignmentGeometry   // Transform alignment
	Alignment            AlignmentGeometry   // Child alignment
	Color                Color               // Background color
	ClipBehavior         Clip                // Clip behavior
}

// Render renders the container as HTML
func (c Container) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(c.ID, c.Style, c.Class+" godin-container")

	// Build inline styles from various sources
	var styles []string

	// Add custom style if provided
	if c.Style != "" {
		styles = append(styles, c.Style)
	}

	// Add padding
	if c.Padding != nil {
		styles = append(styles, fmt.Sprintf("padding: %s", c.Padding.ToCSSString()))
	}

	// Add margin
	if c.Margin != nil {
		styles = append(styles, fmt.Sprintf("margin: %s", c.Margin.ToCSSString()))
	}

	// Add dimensions
	if c.Width != nil {
		styles = append(styles, fmt.Sprintf("width: %.1fpx", *c.Width))
	}
	if c.Height != nil {
		styles = append(styles, fmt.Sprintf("height: %.1fpx", *c.Height))
	}

	// Add background color
	if c.Color != "" {
		styles = append(styles, fmt.Sprintf("background-color: %s", c.Color))
	}

	// Add decoration styles
	if c.Decoration != nil {
		if decorationCSS := c.Decoration.ToCSSString(); decorationCSS != "" {
			styles = append(styles, decorationCSS)
		}
	}

	// Add alignment for child
	if c.Alignment != "" {
		alignParts := strings.Fields(string(c.Alignment))
		if len(alignParts) == 2 {
			styles = append(styles, "display: flex")
			styles = append(styles, fmt.Sprintf("align-items: %s", alignParts[0]))
			styles = append(styles, fmt.Sprintf("justify-content: %s", alignParts[1]))
		}
	}

	// Add constraints (simplified)
	if c.Constraints != nil {
		if c.Constraints.MinWidth != nil {
			styles = append(styles, fmt.Sprintf("min-width: %.1fpx", *c.Constraints.MinWidth))
		}
		if c.Constraints.MaxWidth != nil {
			styles = append(styles, fmt.Sprintf("max-width: %.1fpx", *c.Constraints.MaxWidth))
		}
		if c.Constraints.MinHeight != nil {
			styles = append(styles, fmt.Sprintf("min-height: %.1fpx", *c.Constraints.MinHeight))
		}
		if c.Constraints.MaxHeight != nil {
			styles = append(styles, fmt.Sprintf("max-height: %.1fpx", *c.Constraints.MaxHeight))
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

	// Render child content
	content := ""
	if c.Child != nil {
		content = c.Child.Render(ctx)
	}

	return htmlRenderer.RenderElement("div", attrs, content, false)
}

// Row represents a row layout widget with full Flutter properties
type Row struct {
	ID                 string
	Style              string
	Class              string
	Children           []Widget           // Child widgets
	MainAxisAlignment  MainAxisAlignment  // Main axis alignment
	CrossAxisAlignment CrossAxisAlignment // Cross axis alignment
	MainAxisSize       MainAxisSize       // Main axis size
	TextDirection      TextDirection      // Text direction
	VerticalDirection  VerticalDirection  // Vertical direction
	TextBaseline       TextBaseline       // Text baseline
}

// VerticalDirection enum
type VerticalDirection string

const (
	VerticalDirectionUp   VerticalDirection = "up"
	VerticalDirectionDown VerticalDirection = "down"
)

// Render renders the row as HTML
func (r Row) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(r.ID, r.Style, r.Class+" godin-row")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if r.Style != "" {
		styles = append(styles, r.Style)
	}

	// Base flex styles
	styles = append(styles, "display: flex")
	styles = append(styles, "flex-direction: row")

	// Add main axis alignment
	if r.MainAxisAlignment != "" {
		styles = append(styles, fmt.Sprintf("justify-content: %s", r.MainAxisAlignment))
	}

	// Add cross axis alignment
	if r.CrossAxisAlignment != "" {
		styles = append(styles, fmt.Sprintf("align-items: %s", r.CrossAxisAlignment))
	}

	// Add text direction
	if r.TextDirection != "" {
		styles = append(styles, fmt.Sprintf("direction: %s", r.TextDirection))
		attrs["dir"] = string(r.TextDirection)
	}

	// Handle vertical direction (reverse if up)
	if r.VerticalDirection == VerticalDirectionUp {
		styles = append(styles, "flex-direction: row-reverse")
	}

	// Handle main axis size
	if r.MainAxisSize == MainAxisSizeMin {
		styles = append(styles, "width: min-content")
	}

	// Add text baseline alignment if specified
	if r.TextBaseline != "" && r.CrossAxisAlignment == CrossAxisAlignmentBaseline {
		styles = append(styles, fmt.Sprintf("align-items: %s", r.CrossAxisAlignment))
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Render children
	var children []string
	for _, child := range r.Children {
		if child != nil {
			children = append(children, child.Render(ctx))
		}
	}

	return htmlRenderer.RenderContainer("div", attrs, children)
}

// Column represents a column layout widget with full Flutter properties
type Column struct {
	ID                 string
	Style              string
	Class              string
	Children           []Widget           // Child widgets
	MainAxisAlignment  MainAxisAlignment  // Main axis alignment
	CrossAxisAlignment CrossAxisAlignment // Cross axis alignment
	MainAxisSize       MainAxisSize       // Main axis size
	TextDirection      TextDirection      // Text direction
	VerticalDirection  VerticalDirection  // Vertical direction
	TextBaseline       TextBaseline       // Text baseline
}

// Render renders the column as HTML
func (c Column) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(c.ID, c.Style, c.Class+" godin-column")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if c.Style != "" {
		styles = append(styles, c.Style)
	}

	// Base flex styles
	styles = append(styles, "display: flex")
	styles = append(styles, "flex-direction: column")

	// Add main axis alignment
	if c.MainAxisAlignment != "" {
		styles = append(styles, fmt.Sprintf("justify-content: %s", c.MainAxisAlignment))
	}

	// Add cross axis alignment
	if c.CrossAxisAlignment != "" {
		styles = append(styles, fmt.Sprintf("align-items: %s", c.CrossAxisAlignment))
	}

	// Add text direction
	if c.TextDirection != "" {
		styles = append(styles, fmt.Sprintf("direction: %s", c.TextDirection))
		attrs["dir"] = string(c.TextDirection)
	}

	// Handle vertical direction (reverse if up)
	if c.VerticalDirection == VerticalDirectionUp {
		styles = append(styles, "flex-direction: column-reverse")
	}

	// Handle main axis size
	if c.MainAxisSize == MainAxisSizeMin {
		styles = append(styles, "height: min-content")
	}

	// Add text baseline alignment if specified
	if c.TextBaseline != "" && c.CrossAxisAlignment == CrossAxisAlignmentBaseline {
		styles = append(styles, fmt.Sprintf("align-items: %s", c.CrossAxisAlignment))
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Render children
	var children []string
	for _, child := range c.Children {
		if child != nil {
			children = append(children, child.Render(ctx))
		}
	}

	return htmlRenderer.RenderContainer("div", attrs, children)
}

// Stack represents a stack layout widget with full Flutter properties
type Stack struct {
	ID            string
	Style         string
	Class         string
	Children      []Widget          // Child widgets
	Alignment     AlignmentGeometry // Stack alignment
	TextDirection TextDirection     // Text direction
	Fit           StackFit          // How to size the stack
	ClipBehavior  Clip              // Clip behavior
}

// StackFit enum
type StackFit string

const (
	StackFitLoose       StackFit = "loose"
	StackFitExpand      StackFit = "expand"
	StackFitPassthrough StackFit = "passthrough"
)

// Render renders the stack as HTML
func (s Stack) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(s.ID, s.Style, s.Class+" godin-stack")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if s.Style != "" {
		styles = append(styles, s.Style)
	}

	// Base stack styles
	styles = append(styles, "position: relative")

	// Handle stack fit
	if s.Fit == StackFitExpand {
		styles = append(styles, "width: 100%")
		styles = append(styles, "height: 100%")
	}

	// Add alignment
	if s.Alignment != "" {
		alignParts := strings.Fields(string(s.Alignment))
		if len(alignParts) == 2 {
			styles = append(styles, "display: flex")
			styles = append(styles, fmt.Sprintf("align-items: %s", alignParts[0]))
			styles = append(styles, fmt.Sprintf("justify-content: %s", alignParts[1]))
		}
	}

	// Add text direction
	if s.TextDirection != "" {
		styles = append(styles, fmt.Sprintf("direction: %s", s.TextDirection))
		attrs["dir"] = string(s.TextDirection)
	}

	// Add clip behavior
	if s.ClipBehavior != "" && s.ClipBehavior != ClipNone {
		styles = append(styles, "overflow: hidden")
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Render children
	var children []string
	for _, child := range s.Children {
		if child != nil {
			children = append(children, child.Render(ctx))
		}
	}

	return htmlRenderer.RenderContainer("div", attrs, children)
}

// Positioned represents a positioned widget with full Flutter properties
type Positioned struct {
	ID     string
	Style  string
	Class  string
	Child  Widget   // Child widget
	Left   *float64 // Left position
	Top    *float64 // Top position
	Right  *float64 // Right position
	Bottom *float64 // Bottom position
	Width  *float64 // Width
	Height *float64 // Height
}

// Render renders the positioned widget as HTML
func (p Positioned) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(p.ID, p.Style, p.Class+" godin-positioned")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if p.Style != "" {
		styles = append(styles, p.Style)
	}

	// Base positioned styles
	styles = append(styles, "position: absolute")

	// Add positioning
	if p.Left != nil {
		styles = append(styles, fmt.Sprintf("left: %.1fpx", *p.Left))
	}
	if p.Top != nil {
		styles = append(styles, fmt.Sprintf("top: %.1fpx", *p.Top))
	}
	if p.Right != nil {
		styles = append(styles, fmt.Sprintf("right: %.1fpx", *p.Right))
	}
	if p.Bottom != nil {
		styles = append(styles, fmt.Sprintf("bottom: %.1fpx", *p.Bottom))
	}

	// Add dimensions
	if p.Width != nil {
		styles = append(styles, fmt.Sprintf("width: %.1fpx", *p.Width))
	}
	if p.Height != nil {
		styles = append(styles, fmt.Sprintf("height: %.1fpx", *p.Height))
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Render child content
	content := ""
	if p.Child != nil {
		content = p.Child.Render(ctx)
	}

	return htmlRenderer.RenderElement("div", attrs, content, false)
}

// Expanded represents an expanded layout widget with full Flutter properties
type Expanded struct {
	ID    string
	Style string
	Class string
	Child Widget // Child widget
	Flex  int    // Flex factor
}

// Render renders the expanded widget as HTML
func (e Expanded) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(e.ID, e.Style, e.Class+" godin-expanded")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if e.Style != "" {
		styles = append(styles, e.Style)
	}

	// Add flex property
	if e.Flex > 0 {
		styles = append(styles, fmt.Sprintf("flex: %d", e.Flex))
	} else {
		styles = append(styles, "flex: 1")
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Render child content
	content := ""
	if e.Child != nil {
		content = e.Child.Render(ctx)
	}

	return htmlRenderer.RenderElement("div", attrs, content, false)
}

// Flexible represents a flexible layout widget with full Flutter properties
type Flexible struct {
	ID    string
	Style string
	Class string
	Child Widget  // Child widget
	Flex  int     // Flex factor
	Fit   FlexFit // Flex fit
}

// FlexFit enum
type FlexFit string

const (
	FlexFitTight FlexFit = "tight"
	FlexFitLoose FlexFit = "loose"
)

// Render renders the flexible widget as HTML
func (f Flexible) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(f.ID, f.Style, f.Class+" godin-flexible")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if f.Style != "" {
		styles = append(styles, f.Style)
	}

	// Add flex property based on fit
	if f.Flex > 0 {
		if f.Fit == FlexFitTight {
			styles = append(styles, fmt.Sprintf("flex: %d 0 auto", f.Flex))
		} else {
			styles = append(styles, fmt.Sprintf("flex: %d 1 auto", f.Flex))
		}
	} else {
		styles = append(styles, "flex: 0 1 auto")
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Render child content
	content := ""
	if f.Child != nil {
		content = f.Child.Render(ctx)
	}

	return htmlRenderer.RenderElement("div", attrs, content, false)
}

// SizedBox represents a sized box widget with full Flutter properties
type SizedBox struct {
	ID     string
	Style  string
	Class  string
	Width  *float64 // Box width
	Height *float64 // Box height
	Child  Widget   // Child widget
}

// Render renders the sized box as HTML
func (sb SizedBox) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(sb.ID, sb.Style, sb.Class+" godin-sizedbox")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if sb.Style != "" {
		styles = append(styles, sb.Style)
	}

	// Add dimensions
	if sb.Width != nil {
		styles = append(styles, fmt.Sprintf("width: %.1fpx", *sb.Width))
	}
	if sb.Height != nil {
		styles = append(styles, fmt.Sprintf("height: %.1fpx", *sb.Height))
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Render child content
	content := ""
	if sb.Child != nil {
		content = sb.Child.Render(ctx)
	}

	return htmlRenderer.RenderElement("div", attrs, content, false)
}

// Padding represents a padding widget with full Flutter properties
type Padding struct {
	ID      string
	Style   string
	Class   string
	Padding EdgeInsetsGeometry // Padding values
	Child   Widget             // Child widget
}

// Render renders the padding widget as HTML
func (p Padding) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(p.ID, p.Style, p.Class+" godin-padding")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if p.Style != "" {
		styles = append(styles, p.Style)
	}

	// Add padding
	styles = append(styles, fmt.Sprintf("padding: %s", p.Padding.ToCSSString()))

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Render child content
	content := ""
	if p.Child != nil {
		content = p.Child.Render(ctx)
	}

	return htmlRenderer.RenderElement("div", attrs, content, false)
}

// Center represents a center widget with full Flutter properties
type Center struct {
	ID           string
	Style        string
	Class        string
	Child        Widget   // Child widget
	WidthFactor  *float64 // Width factor
	HeightFactor *float64 // Height factor
}

// Render renders the center widget as HTML
func (c Center) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(c.ID, c.Style, c.Class+" godin-center")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if c.Style != "" {
		styles = append(styles, c.Style)
	}

	// Base centering styles
	styles = append(styles, "display: flex")
	styles = append(styles, "align-items: center")
	styles = append(styles, "justify-content: center")

	// Add width and height factors
	if c.WidthFactor != nil {
		styles = append(styles, fmt.Sprintf("width: %.2f%%", *c.WidthFactor*100))
	}
	if c.HeightFactor != nil {
		styles = append(styles, fmt.Sprintf("height: %.2f%%", *c.HeightFactor*100))
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Render child content
	content := ""
	if c.Child != nil {
		content = c.Child.Render(ctx)
	}

	return htmlRenderer.RenderElement("div", attrs, content, false)
}

// Align represents an align widget with full Flutter properties
type Align struct {
	ID           string
	Style        string
	Class        string
	Child        Widget            // Child widget
	Alignment    AlignmentGeometry // Alignment
	WidthFactor  *float64          // Width factor
	HeightFactor *float64          // Height factor
}

// Render renders the align widget as HTML
func (a Align) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(a.ID, a.Style, a.Class+" godin-align")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if a.Style != "" {
		styles = append(styles, a.Style)
	}

	// Base alignment styles
	styles = append(styles, "display: flex")

	// Add alignment
	if a.Alignment != "" {
		alignParts := strings.Fields(string(a.Alignment))
		if len(alignParts) == 2 {
			styles = append(styles, fmt.Sprintf("align-items: %s", alignParts[0]))
			styles = append(styles, fmt.Sprintf("justify-content: %s", alignParts[1]))
		}
	}

	// Add width and height factors
	if a.WidthFactor != nil {
		styles = append(styles, fmt.Sprintf("width: %.2f%%", *a.WidthFactor*100))
	}
	if a.HeightFactor != nil {
		styles = append(styles, fmt.Sprintf("height: %.2f%%", *a.HeightFactor*100))
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Render child content
	content := ""
	if a.Child != nil {
		content = a.Child.Render(ctx)
	}

	return htmlRenderer.RenderElement("div", attrs, content, false)
}

// Transform represents a transform widget with full Flutter properties
type Transform struct {
	ID                string
	Style             string
	Class             string
	Child             Widget            // Child widget
	Transform         Matrix4           // Transform matrix
	Origin            Offset            // Transform origin
	Alignment         AlignmentGeometry // Alignment
	TransformHitTests bool              // Transform hit tests
	FilterQuality     FilterQuality     // Filter quality
}

// NewMatrix4Identity creates an identity matrix
func NewMatrix4Identity() Matrix4 {
	return Matrix4{
		Values: [16]float64{
			1, 0, 0, 0,
			0, 1, 0, 0,
			0, 0, 1, 0,
			0, 0, 0, 1,
		},
	}
}

// NewMatrix4RotationZ creates a rotation matrix around Z axis
func NewMatrix4RotationZ(radians float64) Matrix4 {
	cos := math.Cos(radians)
	sin := math.Sin(radians)
	return Matrix4{
		Values: [16]float64{
			cos, -sin, 0, 0,
			sin, cos, 0, 0,
			0, 0, 1, 0,
			0, 0, 0, 1,
		},
	}
}

// NewMatrix4Translation creates a translation matrix
func NewMatrix4Translation(x, y, z float64) Matrix4 {
	return Matrix4{
		Values: [16]float64{
			1, 0, 0, x,
			0, 1, 0, y,
			0, 0, 1, z,
			0, 0, 0, 1,
		},
	}
}

// NewMatrix4Scale creates a scale matrix
func NewMatrix4Scale(x, y, z float64) Matrix4 {
	return Matrix4{
		Values: [16]float64{
			x, 0, 0, 0,
			0, y, 0, 0,
			0, 0, z, 0,
			0, 0, 0, 1,
		},
	}
}

// ToCSSString converts the matrix to CSS transform matrix3d
func (m Matrix4) ToCSSString() string {
	return fmt.Sprintf("matrix3d(%.6f,%.6f,%.6f,%.6f,%.6f,%.6f,%.6f,%.6f,%.6f,%.6f,%.6f,%.6f,%.6f,%.6f,%.6f,%.6f)",
		m.Values[0], m.Values[1], m.Values[2], m.Values[3],
		m.Values[4], m.Values[5], m.Values[6], m.Values[7],
		m.Values[8], m.Values[9], m.Values[10], m.Values[11],
		m.Values[12], m.Values[13], m.Values[14], m.Values[15])
}

// Render renders the transform widget as HTML
func (t Transform) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(t.ID, t.Style, t.Class+" godin-transform")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if t.Style != "" {
		styles = append(styles, t.Style)
	}

	// Add transform matrix
	styles = append(styles, fmt.Sprintf("transform: %s", t.Transform.ToCSSString()))

	// Add transform origin
	if t.Origin.DX != 0 || t.Origin.DY != 0 {
		styles = append(styles, fmt.Sprintf("transform-origin: %.1fpx %.1fpx", t.Origin.DX, t.Origin.DY))
	}

	// Add filter quality (simplified as image-rendering)
	if t.FilterQuality != "" {
		switch t.FilterQuality {
		case FilterQualityHigh:
			styles = append(styles, "image-rendering: high-quality")
		case FilterQualityMedium:
			styles = append(styles, "image-rendering: auto")
		case FilterQualityLow:
			styles = append(styles, "image-rendering: pixelated")
		}
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Render child content
	content := ""
	if t.Child != nil {
		content = t.Child.Render(ctx)
	}

	return htmlRenderer.RenderElement("div", attrs, content, false)
}

// AnimatedContainer represents an animated container widget with full Flutter properties
type AnimatedContainer struct {
	ID                   string
	Style                string
	Class                string
	Child                Widget              // Child widget
	Alignment            AlignmentGeometry   // Alignment
	Padding              *EdgeInsetsGeometry // Padding
	Color                Color               // Background color
	Decoration           *BoxDecoration      // Box decoration
	ForegroundDecoration *BoxDecoration      // Foreground decoration
	Width                *float64            // Width
	Height               *float64            // Height
	Constraints          *BoxConstraints     // Box constraints
	Margin               *EdgeInsetsGeometry // Margin
	Transform            *Matrix4            // Transform
	TransformAlignment   AlignmentGeometry   // Transform alignment
	Curve                Curve               // Animation curve
	Duration             Duration            // Animation duration
	OnEnd                VoidCallback        // On animation end callback
	ClipBehavior         Clip                // Clip behavior
}

// Render renders the animated container as HTML
func (ac AnimatedContainer) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(ac.ID, ac.Style, ac.Class+" godin-animated-container")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if ac.Style != "" {
		styles = append(styles, ac.Style)
	}

	// Add animation transition
	transitionDuration := "300ms"
	if ac.Duration > 0 {
		transitionDuration = fmt.Sprintf("%dms", int(ac.Duration/1000000)) // Convert nanoseconds to milliseconds
	}

	transitionCurve := "ease"
	if ac.Curve != "" {
		transitionCurve = string(ac.Curve)
	}

	styles = append(styles, fmt.Sprintf("transition: all %s %s", transitionDuration, transitionCurve))

	// Add dimensions
	if ac.Width != nil {
		styles = append(styles, fmt.Sprintf("width: %.1fpx", *ac.Width))
	}
	if ac.Height != nil {
		styles = append(styles, fmt.Sprintf("height: %.1fpx", *ac.Height))
	}

	// Add background color
	if ac.Color != "" {
		styles = append(styles, fmt.Sprintf("background-color: %s", ac.Color))
	}

	// Add padding
	if ac.Padding != nil {
		styles = append(styles, fmt.Sprintf("padding: %s", ac.Padding.ToCSSString()))
	}

	// Add margin
	if ac.Margin != nil {
		styles = append(styles, fmt.Sprintf("margin: %s", ac.Margin.ToCSSString()))
	}

	// Add decoration
	if ac.Decoration != nil {
		decorationCSS := ac.Decoration.ToCSSString()
		if decorationCSS != "" {
			styles = append(styles, decorationCSS)
		}
	}

	// Add transform
	if ac.Transform != nil {
		styles = append(styles, fmt.Sprintf("transform: %s", ac.Transform.ToCSSString()))
	}

	// Add alignment
	if ac.Alignment != "" {
		alignParts := strings.Fields(string(ac.Alignment))
		if len(alignParts) == 2 {
			styles = append(styles, "display: flex")
			styles = append(styles, fmt.Sprintf("align-items: %s", alignParts[0]))
			styles = append(styles, fmt.Sprintf("justify-content: %s", alignParts[1]))
		}
	}

	// Add clip behavior
	if ac.ClipBehavior != "" && ac.ClipBehavior != ClipNone {
		styles = append(styles, "overflow: hidden")
	}

	// Add constraints (simplified as min/max width/height)
	if ac.Constraints != nil {
		if ac.Constraints.MinWidth != nil && *ac.Constraints.MinWidth > 0 {
			styles = append(styles, fmt.Sprintf("min-width: %.1fpx", *ac.Constraints.MinWidth))
		}
		if ac.Constraints.MaxWidth != nil && *ac.Constraints.MaxWidth < math.Inf(1) {
			styles = append(styles, fmt.Sprintf("max-width: %.1fpx", *ac.Constraints.MaxWidth))
		}
		if ac.Constraints.MinHeight != nil && *ac.Constraints.MinHeight > 0 {
			styles = append(styles, fmt.Sprintf("min-height: %.1fpx", *ac.Constraints.MinHeight))
		}
		if ac.Constraints.MaxHeight != nil && *ac.Constraints.MaxHeight < math.Inf(1) {
			styles = append(styles, fmt.Sprintf("max-height: %.1fpx", *ac.Constraints.MaxHeight))
		}
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Add animation end handler
	if ac.OnEnd != nil {
		attrs["ontransitionend"] = "handleAnimatedContainerEnd(this)"
	}

	// Render child content
	content := ""
	if ac.Child != nil {
		content = ac.Child.Render(ctx)
	}

	return htmlRenderer.RenderElement("div", attrs, content, false)
}
