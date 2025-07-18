package widgets

import (
	"fmt"
	"godin-framework/pkg/core"
	"godin-framework/pkg/renderer"
)

// Widget is an alias for core.Widget for convenience
type Widget = core.Widget

// Common Flutter-style types and enums

// EdgeInsetsGeometry represents padding/margin values
type EdgeInsetsGeometry struct {
	Top    float64
	Right  float64
	Bottom float64
	Left   float64
}

// EdgeInsets creates EdgeInsetsGeometry with all sides equal
func EdgeInsets(value float64) EdgeInsetsGeometry {
	return EdgeInsetsGeometry{Top: value, Right: value, Bottom: value, Left: value}
}

// EdgeInsetsOnly creates EdgeInsetsGeometry with specific sides
func EdgeInsetsOnly(top, right, bottom, left float64) EdgeInsetsGeometry {
	return EdgeInsetsGeometry{Top: top, Right: right, Bottom: bottom, Left: left}
}

// EdgeInsetsSymmetric creates EdgeInsetsGeometry with symmetric values
func EdgeInsetsSymmetric(vertical, horizontal float64) EdgeInsetsGeometry {
	return EdgeInsetsGeometry{Top: vertical, Right: horizontal, Bottom: vertical, Left: horizontal}
}

// ToCSSString converts EdgeInsetsGeometry to CSS padding/margin string
func (e EdgeInsetsGeometry) ToCSSString() string {
	return fmt.Sprintf("%.1fpx %.1fpx %.1fpx %.1fpx", e.Top, e.Right, e.Bottom, e.Left)
}

// AlignmentGeometry represents alignment values
type AlignmentGeometry string

const (
	AlignmentTopLeft      AlignmentGeometry = "flex-start flex-start"
	AlignmentTopCenter    AlignmentGeometry = "flex-start center"
	AlignmentTopRight     AlignmentGeometry = "flex-start flex-end"
	AlignmentCenterLeft   AlignmentGeometry = "center flex-start"
	AlignmentCenter       AlignmentGeometry = "center center"
	AlignmentCenterRight  AlignmentGeometry = "center flex-end"
	AlignmentBottomLeft   AlignmentGeometry = "flex-end flex-start"
	AlignmentBottomCenter AlignmentGeometry = "flex-end center"
	AlignmentBottomRight  AlignmentGeometry = "flex-end flex-end"
)

// MainAxisAlignment for Row/Column widgets
type MainAxisAlignment string

const (
	MainAxisAlignmentStart        MainAxisAlignment = "flex-start"
	MainAxisAlignmentEnd          MainAxisAlignment = "flex-end"
	MainAxisAlignmentCenter       MainAxisAlignment = "center"
	MainAxisAlignmentSpaceBetween MainAxisAlignment = "space-between"
	MainAxisAlignmentSpaceAround  MainAxisAlignment = "space-around"
	MainAxisAlignmentSpaceEvenly  MainAxisAlignment = "space-evenly"
)

// CrossAxisAlignment for Row/Column widgets
type CrossAxisAlignment string

const (
	CrossAxisAlignmentStart    CrossAxisAlignment = "flex-start"
	CrossAxisAlignmentEnd      CrossAxisAlignment = "flex-end"
	CrossAxisAlignmentCenter   CrossAxisAlignment = "center"
	CrossAxisAlignmentStretch  CrossAxisAlignment = "stretch"
	CrossAxisAlignmentBaseline CrossAxisAlignment = "baseline"
)

// MainAxisSize for Row/Column widgets
type MainAxisSize string

const (
	MainAxisSizeMin MainAxisSize = "min"
	MainAxisSizeMax MainAxisSize = "max"
)

// TextAlign enum
type TextAlign string

const (
	TextAlignLeft    TextAlign = "left"
	TextAlignRight   TextAlign = "right"
	TextAlignCenter  TextAlign = "center"
	TextAlignJustify TextAlign = "justify"
	TextAlignStart   TextAlign = "start"
	TextAlignEnd     TextAlign = "end"
)

// TextDirection enum
type TextDirection string

const (
	TextDirectionLTR TextDirection = "ltr"
	TextDirectionRTL TextDirection = "rtl"
)

// TextOverflow enum
type TextOverflow string

const (
	TextOverflowClip     TextOverflow = "clip"
	TextOverflowEllipsis TextOverflow = "ellipsis"
	TextOverflowFade     TextOverflow = "fade"
	TextOverflowVisible  TextOverflow = "visible"
)

// buildAttributes builds HTML attributes for a widget
func buildAttributes(id, style, class string) map[string]string {
	attrs := make(map[string]string)

	if id != "" {
		attrs["id"] = id
	}

	// Combine custom style with class-based styles
	cssClass := "godin-widget"
	if class != "" {
		cssClass += " " + class
	}
	attrs["class"] = cssClass

	if style != "" {
		attrs["style"] = style
	}

	return attrs
}

// buildHTMXAttributes builds HTML attributes including HTMX attributes
func buildHTMXAttributes(id, style, class string, htmx renderer.HTMXAttributes) map[string]string {
	attrs := buildAttributes(id, style, class)

	// Add HTMX attributes
	htmxRenderer := renderer.NewHTMXRenderer()
	htmxAttrs := htmxRenderer.RenderAttributes(htmx)

	for key, value := range htmxAttrs {
		attrs[key] = value
	}

	return attrs
}

// HTMXWidget is a temporary stub for widgets that haven't been converted yet
type HTMXWidget struct {
	ID    string
	Style string
	Class string
	HTMX  renderer.HTMXAttributes
}

// buildHTMXAttributes method for backward compatibility
func (w HTMXWidget) buildHTMXAttributes() map[string]string {
	return buildHTMXAttributes(w.ID, w.Style, w.Class, w.HTMX)
}
