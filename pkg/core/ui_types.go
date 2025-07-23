package core

import (
	"fmt"
	"math"
)

// Color represents an RGBA color
type Color struct {
	R, G, B, A uint8
}

// NewColor creates a new Color from RGBA values
func NewColor(r, g, b, a uint8) Color {
	return Color{R: r, G: g, B: b, A: a}
}

// NewColorFromHex creates a Color from a hex string (e.g., "#FF0000" or "FF0000")
func NewColorFromHex(hex string) (Color, error) {
	if len(hex) == 0 {
		return Color{}, fmt.Errorf("empty hex string")
	}

	// Remove # if present
	if hex[0] == '#' {
		hex = hex[1:]
	}

	if len(hex) != 6 && len(hex) != 8 {
		return Color{}, fmt.Errorf("invalid hex color length: %d", len(hex))
	}

	var r, g, b, a uint8 = 0, 0, 0, 255

	if _, err := fmt.Sscanf(hex[:2], "%02x", &r); err != nil {
		return Color{}, err
	}
	if _, err := fmt.Sscanf(hex[2:4], "%02x", &g); err != nil {
		return Color{}, err
	}
	if _, err := fmt.Sscanf(hex[4:6], "%02x", &b); err != nil {
		return Color{}, err
	}

	if len(hex) == 8 {
		if _, err := fmt.Sscanf(hex[6:8], "%02x", &a); err != nil {
			return Color{}, err
		}
	}

	return Color{R: r, G: g, B: b, A: a}, nil
}

// ToHex converts the color to a hex string
func (c Color) ToHex() string {
	if c.A == 255 {
		return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
	}
	return fmt.Sprintf("#%02X%02X%02X%02X", c.R, c.G, c.B, c.A)
}

// ToRGBA converts the color to CSS rgba() format
func (c Color) ToRGBA() string {
	alpha := float64(c.A) / 255.0
	return fmt.Sprintf("rgba(%d, %d, %d, %.2f)", c.R, c.G, c.B, alpha)
}

// ToCSS converts the color to the most appropriate CSS format
func (c Color) ToCSS() string {
	if c.A == 255 {
		return c.ToHex()
	}
	return c.ToRGBA()
}

// WithOpacity returns a new color with the specified opacity (0.0 to 1.0)
func (c Color) WithOpacity(opacity float64) Color {
	alpha := uint8(math.Max(0, math.Min(255, opacity*255)))
	return Color{R: c.R, G: c.G, B: c.B, A: alpha}
}

// Darken returns a darker version of the color
func (c Color) Darken(amount float64) Color {
	factor := 1.0 - math.Max(0, math.Min(1, amount))
	return Color{
		R: uint8(float64(c.R) * factor),
		G: uint8(float64(c.G) * factor),
		B: uint8(float64(c.B) * factor),
		A: c.A,
	}
}

// Lighten returns a lighter version of the color
func (c Color) Lighten(amount float64) Color {
	factor := math.Max(0, math.Min(1, amount))
	return Color{
		R: uint8(float64(c.R) + (255-float64(c.R))*factor),
		G: uint8(float64(c.G) + (255-float64(c.G))*factor),
		B: uint8(float64(c.B) + (255-float64(c.B))*factor),
		A: c.A,
	}
}

// Size represents width and height dimensions
type Size struct {
	Width  float64
	Height float64
}

// NewSize creates a new Size
func NewSize(width, height float64) Size {
	return Size{Width: width, Height: height}
}

// IsEmpty returns true if the size has zero or negative dimensions
func (s Size) IsEmpty() bool {
	return s.Width <= 0 || s.Height <= 0
}

// AspectRatio returns the width/height ratio
func (s Size) AspectRatio() float64 {
	if s.Height == 0 {
		return 0
	}
	return s.Width / s.Height
}

// EdgeInsets represents padding or margin values
type EdgeInsets struct {
	Top    float64
	Right  float64
	Bottom float64
	Left   float64
}

// NewEdgeInsets creates EdgeInsets with individual values
func NewEdgeInsets(top, right, bottom, left float64) EdgeInsets {
	return EdgeInsets{Top: top, Right: right, Bottom: bottom, Left: left}
}

// NewEdgeInsetsAll creates EdgeInsets with the same value for all sides
func NewEdgeInsetsAll(value float64) EdgeInsets {
	return EdgeInsets{Top: value, Right: value, Bottom: value, Left: value}
}

// NewEdgeInsetsSymmetric creates EdgeInsets with vertical and horizontal values
func NewEdgeInsetsSymmetric(vertical, horizontal float64) EdgeInsets {
	return EdgeInsets{Top: vertical, Right: horizontal, Bottom: vertical, Left: horizontal}
}

// NewEdgeInsetsOnly creates EdgeInsets with only specified sides
func NewEdgeInsetsOnly(top, right, bottom, left *float64) EdgeInsets {
	insets := EdgeInsets{}
	if top != nil {
		insets.Top = *top
	}
	if right != nil {
		insets.Right = *right
	}
	if bottom != nil {
		insets.Bottom = *bottom
	}
	if left != nil {
		insets.Left = *left
	}
	return insets
}

// ToCSS converts EdgeInsets to CSS padding/margin format
func (e EdgeInsets) ToCSS() string {
	return fmt.Sprintf("%.1fpx %.1fpx %.1fpx %.1fpx", e.Top, e.Right, e.Bottom, e.Left)
}

// Horizontal returns the total horizontal insets (left + right)
func (e EdgeInsets) Horizontal() float64 {
	return e.Left + e.Right
}

// Vertical returns the total vertical insets (top + bottom)
func (e EdgeInsets) Vertical() float64 {
	return e.Top + e.Bottom
}

// Brightness represents light or dark theme mode
type Brightness string

const (
	BrightnessLight Brightness = "light"
	BrightnessDark  Brightness = "dark"
)

// ThemeMode represents how theme should be determined
type ThemeMode string

const (
	ThemeModeSystem ThemeMode = "system" // Follow system preference
	ThemeModeLight  ThemeMode = "light"  // Always use light theme
	ThemeModeDark   ThemeMode = "dark"   // Always use dark theme
)

// Orientation represents screen orientation
type Orientation string

const (
	OrientationPortrait  Orientation = "portrait"
	OrientationLandscape Orientation = "landscape"
)

// Breakpoint represents responsive design breakpoints
type Breakpoint string

const (
	BreakpointXS Breakpoint = "xs" // < 576px
	BreakpointSM Breakpoint = "sm" // >= 576px
	BreakpointMD Breakpoint = "md" // >= 768px
	BreakpointLG Breakpoint = "lg" // >= 992px
	BreakpointXL Breakpoint = "xl" // >= 1200px
)

// BreakpointValues maps breakpoints to their pixel values
var BreakpointValues = map[Breakpoint]float64{
	BreakpointXS: 0,
	BreakpointSM: 576,
	BreakpointMD: 768,
	BreakpointLG: 992,
	BreakpointXL: 1200,
}

// GetBreakpoint returns the appropriate breakpoint for a given width
func GetBreakpoint(width float64) Breakpoint {
	if width >= BreakpointValues[BreakpointXL] {
		return BreakpointXL
	} else if width >= BreakpointValues[BreakpointLG] {
		return BreakpointLG
	} else if width >= BreakpointValues[BreakpointMD] {
		return BreakpointMD
	} else if width >= BreakpointValues[BreakpointSM] {
		return BreakpointSM
	}
	return BreakpointXS
}

// TextAlign represents text alignment options
type TextAlign string

const (
	TextAlignLeft    TextAlign = "left"
	TextAlignRight   TextAlign = "right"
	TextAlignCenter  TextAlign = "center"
	TextAlignJustify TextAlign = "justify"
	TextAlignStart   TextAlign = "start"
	TextAlignEnd     TextAlign = "end"
)

// FontWeight represents font weight values
type FontWeight int

const (
	FontWeightThin       FontWeight = 100
	FontWeightExtraLight FontWeight = 200
	FontWeightLight      FontWeight = 300
	FontWeightNormal     FontWeight = 400
	FontWeightMedium     FontWeight = 500
	FontWeightSemiBold   FontWeight = 600
	FontWeightBold       FontWeight = 700
	FontWeightExtraBold  FontWeight = 800
	FontWeightBlack      FontWeight = 900
)

// TextStyle represents text styling options
type TextStyle struct {
	Color          *Color
	FontSize       *float64
	FontWeight     *FontWeight
	FontFamily     *string
	LetterSpacing  *float64
	LineHeight     *float64
	TextAlign      *TextAlign
	TextDecoration *string
}

// NewTextStyle creates a new TextStyle
func NewTextStyle() *TextStyle {
	return &TextStyle{}
}

// WithColor sets the text color
func (ts *TextStyle) WithColor(color Color) *TextStyle {
	ts.Color = &color
	return ts
}

// WithFontSize sets the font size
func (ts *TextStyle) WithFontSize(size float64) *TextStyle {
	ts.FontSize = &size
	return ts
}

// WithFontWeight sets the font weight
func (ts *TextStyle) WithFontWeight(weight FontWeight) *TextStyle {
	ts.FontWeight = &weight
	return ts
}

// WithFontFamily sets the font family
func (ts *TextStyle) WithFontFamily(family string) *TextStyle {
	ts.FontFamily = &family
	return ts
}

// ToCSS converts TextStyle to CSS properties
func (ts *TextStyle) ToCSS() map[string]string {
	css := make(map[string]string)

	if ts.Color != nil {
		css["color"] = ts.Color.ToCSS()
	}
	if ts.FontSize != nil {
		css["font-size"] = fmt.Sprintf("%.1fpx", *ts.FontSize)
	}
	if ts.FontWeight != nil {
		css["font-weight"] = fmt.Sprintf("%d", *ts.FontWeight)
	}
	if ts.FontFamily != nil {
		css["font-family"] = *ts.FontFamily
	}
	if ts.LetterSpacing != nil {
		css["letter-spacing"] = fmt.Sprintf("%.2fem", *ts.LetterSpacing)
	}
	if ts.LineHeight != nil {
		css["line-height"] = fmt.Sprintf("%.2f", *ts.LineHeight)
	}
	if ts.TextAlign != nil {
		css["text-align"] = string(*ts.TextAlign)
	}
	if ts.TextDecoration != nil {
		css["text-decoration"] = *ts.TextDecoration
	}

	return css
}

// Predefined colors following Material Design
var (
	// Primary colors
	ColorPrimary      = NewColor(103, 80, 164, 255) // Material Purple
	ColorPrimaryLight = NewColor(187, 134, 252, 255)
	ColorPrimaryDark  = NewColor(77, 56, 138, 255)

	// Secondary colors
	ColorSecondary      = NewColor(3, 218, 198, 255) // Material Teal
	ColorSecondaryLight = NewColor(129, 199, 132, 255)
	ColorSecondaryDark  = NewColor(0, 150, 136, 255)

	// Surface colors
	ColorSurface        = NewColor(255, 255, 255, 255)
	ColorSurfaceDark    = NewColor(18, 18, 18, 255)
	ColorBackground     = NewColor(255, 255, 255, 255)
	ColorBackgroundDark = NewColor(0, 0, 0, 255)

	// Text colors
	ColorOnSurface        = NewColor(0, 0, 0, 255)
	ColorOnSurfaceDark    = NewColor(255, 255, 255, 255)
	ColorOnBackground     = NewColor(0, 0, 0, 255)
	ColorOnBackgroundDark = NewColor(255, 255, 255, 255)
	ColorOnPrimary        = NewColor(255, 255, 255, 255)
	ColorOnPrimaryDark    = NewColor(0, 0, 0, 255)
	ColorOnSecondary      = NewColor(0, 0, 0, 255)
	ColorOnSecondaryDark  = NewColor(0, 0, 0, 255)

	// Status colors
	ColorError   = NewColor(211, 47, 47, 255)
	ColorWarning = NewColor(255, 152, 0, 255)
	ColorSuccess = NewColor(76, 175, 80, 255)
	ColorInfo    = NewColor(33, 150, 243, 255)

	// Neutral colors
	ColorTransparent = NewColor(0, 0, 0, 0)
	ColorWhite       = NewColor(255, 255, 255, 255)
	ColorBlack       = NewColor(0, 0, 0, 255)
	ColorGrey        = NewColor(158, 158, 158, 255)
	ColorGreyLight   = NewColor(245, 245, 245, 255)
	ColorGreyDark    = NewColor(66, 66, 66, 255)
)

// Utility functions for common operations

// LerpColor interpolates between two colors
func LerpColor(a, b Color, t float64) Color {
	t = math.Max(0, math.Min(1, t))
	return Color{
		R: uint8(float64(a.R) + (float64(b.R)-float64(a.R))*t),
		G: uint8(float64(a.G) + (float64(b.G)-float64(a.G))*t),
		B: uint8(float64(a.B) + (float64(b.B)-float64(a.B))*t),
		A: uint8(float64(a.A) + (float64(b.A)-float64(a.A))*t),
	}
}

// LerpSize interpolates between two sizes
func LerpSize(a, b Size, t float64) Size {
	t = math.Max(0, math.Min(1, t))
	return Size{
		Width:  a.Width + (b.Width-a.Width)*t,
		Height: a.Height + (b.Height-a.Height)*t,
	}
}

// ClampFloat64 clamps a float64 value between min and max
func ClampFloat64(value, min, max float64) float64 {
	return math.Max(min, math.Min(max, value))
}

// ClampInt clamps an int value between min and max
func ClampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
