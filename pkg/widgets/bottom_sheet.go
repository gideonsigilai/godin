package widgets

import (
	"fmt"
	"strings"

	"github.com/gideonsigilai/godin/pkg/core"
	"github.com/gideonsigilai/godin/pkg/renderer"
)

// BottomSheet represents a bottom sheet widget
type BottomSheet struct {
	InteractiveWidget
	Child           core.Widget
	BackgroundColor *core.Color
	Elevation       float64
	Shape           *BorderRadius
	ClipBehavior    ClipBehavior
	EnableDrag      bool
	ShowDragHandle  bool
	OnDragStart     func()
	OnDragUpdate    func(float64)
	OnDragEnd       func()
}

// ModalBottomSheet represents a modal bottom sheet
type ModalBottomSheet struct {
	InteractiveWidget
	Child              core.Widget
	BackgroundColor    *core.Color
	Elevation          float64
	Shape              *BorderRadius
	IsScrollControlled bool
	UseRootNavigator   bool
	IsDismissible      bool
	EnableDrag         bool
	ShowDragHandle     bool
}

// DraggableScrollableSheet represents a draggable scrollable bottom sheet
type DraggableScrollableSheet struct {
	InteractiveWidget
	Child            core.Widget
	InitialChildSize float64
	MinChildSize     float64
	MaxChildSize     float64
	Expand           bool
	Snap             bool
	SnapSizes        []float64
}

// NewBottomSheet creates a new BottomSheet widget
func NewBottomSheet() *BottomSheet {
	return &BottomSheet{
		Elevation:      8.0,
		EnableDrag:     true,
		ShowDragHandle: true,
		ClipBehavior:   ClipBehaviorAntiAlias,
	}
}

// NewModalBottomSheet creates a new ModalBottomSheet widget
func NewModalBottomSheet() *ModalBottomSheet {
	return &ModalBottomSheet{
		Elevation:          16.0,
		IsScrollControlled: false,
		UseRootNavigator:   true,
		IsDismissible:      true,
		EnableDrag:         true,
		ShowDragHandle:     true,
	}
}

// NewDraggableScrollableSheet creates a new DraggableScrollableSheet widget
func NewDraggableScrollableSheet() *DraggableScrollableSheet {
	return &DraggableScrollableSheet{
		InitialChildSize: 0.5,
		MinChildSize:     0.25,
		MaxChildSize:     1.0,
		Expand:           true,
		Snap:             false,
	}
}

// Render renders the BottomSheet widget
func (bs *BottomSheet) Render(ctx *core.Context) string {
	if !bs.InteractiveWidget.IsInitialized() {
		bs.InteractiveWidget.Initialize(ctx)
		bs.InteractiveWidget.SetWidgetType("BottomSheet")
	}

	htmlRenderer := renderer.NewHTMLRenderer()

	// Build container attributes
	containerAttrs := map[string]string{
		"class": "godin-bottom-sheet-container",
		"style": bs.buildContainerStyle(),
	}

	// Build bottom sheet attributes
	sheetAttrs := map[string]string{
		"class":      "godin-bottom-sheet",
		"style":      bs.buildSheetStyle(),
		"role":       "dialog",
		"aria-modal": "false",
	}

	// Add drag attributes if enabled
	if bs.EnableDrag {
		sheetAttrs["data-draggable"] = "true"
		// Register drag callbacks
		if bs.OnDragStart != nil {
			bs.InteractiveWidget.RegisterCallback("OnDragStart", bs.OnDragStart)
		}
		if bs.OnDragUpdate != nil {
			bs.InteractiveWidget.RegisterCallback("OnDragUpdate", bs.OnDragUpdate)
		}
		if bs.OnDragEnd != nil {
			bs.InteractiveWidget.RegisterCallback("OnDragEnd", bs.OnDragEnd)
		}
	}

	// Merge with interactive widget attributes
	sheetAttrs = bs.InteractiveWidget.MergeAttributes(sheetAttrs)

	// Build sheet content
	var contentParts []string

	// Add drag handle if enabled
	if bs.ShowDragHandle {
		dragHandle := htmlRenderer.RenderElement("div", map[string]string{
			"class": "godin-bottom-sheet-drag-handle",
			"style": bs.buildDragHandleStyle(),
		}, "", true)
		contentParts = append(contentParts, dragHandle)
	}

	// Add child content
	if bs.Child != nil {
		childHTML := htmlRenderer.RenderElement("div", map[string]string{
			"class": "godin-bottom-sheet-content",
			"style": "flex: 1; overflow-y: auto;",
		}, bs.Child.Render(ctx), false)
		contentParts = append(contentParts, childHTML)
	}

	// Build complete bottom sheet
	sheetContent := strings.Join(contentParts, "")
	sheetHTML := htmlRenderer.RenderElement("div", sheetAttrs, sheetContent, false)

	// Wrap in container
	return htmlRenderer.RenderElement("div", containerAttrs, sheetHTML, false)
}

// buildContainerStyle builds the CSS style for the bottom sheet container
func (bs *BottomSheet) buildContainerStyle() string {
	styles := []string{
		"position: fixed",
		"bottom: 0",
		"left: 0",
		"right: 0",
		"z-index: 1000",
		"pointer-events: none", // Allow clicks to pass through to background
	}
	return strings.Join(styles, "; ")
}

// buildSheetStyle builds the CSS style for the bottom sheet
func (bs *BottomSheet) buildSheetStyle() string {
	styles := []string{
		"background-color: white",
		"border-top-left-radius: 16px",
		"border-top-right-radius: 16px",
		"max-height: 90vh",
		"display: flex",
		"flex-direction: column",
		"pointer-events: auto",     // Re-enable pointer events for the sheet
		"transform: translateY(0)", // Initial position
		"transition: transform 0.3s cubic-bezier(0.4, 0.0, 0.2, 1)",
	}

	if bs.BackgroundColor != nil {
		styles = append(styles, fmt.Sprintf("background-color: %s", bs.BackgroundColor.ToCSS()))
	}

	if bs.Elevation > 0 {
		shadowValue := bs.calculateBoxShadow(bs.Elevation)
		styles = append(styles, fmt.Sprintf("box-shadow: %s", shadowValue))
	}

	if bs.Shape != nil {
		styles = append(styles, fmt.Sprintf("border-radius: %s", bs.Shape.ToCSS()))
	}

	return strings.Join(styles, "; ")
}

// buildDragHandleStyle builds the CSS style for the drag handle
func (bs *BottomSheet) buildDragHandleStyle() string {
	styles := []string{
		"width: 32px",
		"height: 4px",
		"background-color: rgba(0, 0, 0, 0.2)",
		"border-radius: 2px",
		"margin: 12px auto 8px auto",
		"cursor: grab",
	}
	return strings.Join(styles, "; ")
}

// calculateBoxShadow calculates box-shadow based on elevation
func (bs *BottomSheet) calculateBoxShadow(elevation float64) string {
	// Material Design elevation shadows (upward shadows for bottom sheets)
	switch {
	case elevation <= 1:
		return "0px -2px 1px -1px rgba(0,0,0,0.2), 0px -1px 1px 0px rgba(0,0,0,0.14), 0px -1px 3px 0px rgba(0,0,0,0.12)"
	case elevation <= 2:
		return "0px -3px 1px -2px rgba(0,0,0,0.2), 0px -2px 2px 0px rgba(0,0,0,0.14), 0px -1px 5px 0px rgba(0,0,0,0.12)"
	case elevation <= 4:
		return "0px -2px 4px -1px rgba(0,0,0,0.2), 0px -4px 5px 0px rgba(0,0,0,0.14), 0px -1px 10px 0px rgba(0,0,0,0.12)"
	case elevation <= 8:
		return "0px -5px 5px -3px rgba(0,0,0,0.2), 0px -8px 10px 1px rgba(0,0,0,0.14), 0px -3px 14px 2px rgba(0,0,0,0.12)"
	case elevation <= 16:
		return "0px -8px 10px -5px rgba(0,0,0,0.2), 0px -16px 24px 2px rgba(0,0,0,0.14), 0px -6px 30px 5px rgba(0,0,0,0.12)"
	default:
		return "0px -11px 15px -7px rgba(0,0,0,0.2), 0px -24px 38px 3px rgba(0,0,0,0.14), 0px -9px 46px 8px rgba(0,0,0,0.12)"
	}
}

// Render renders the ModalBottomSheet widget
func (mbs *ModalBottomSheet) Render(ctx *core.Context) string {
	if !mbs.InteractiveWidget.IsInitialized() {
		mbs.InteractiveWidget.Initialize(ctx)
		mbs.InteractiveWidget.SetWidgetType("ModalBottomSheet")
	}

	htmlRenderer := renderer.NewHTMLRenderer()

	// Build backdrop attributes
	backdropAttrs := map[string]string{
		"class": "godin-modal-bottom-sheet-backdrop",
		"style": mbs.buildBackdropStyle(),
	}

	// Add dismiss callback if dismissible
	if mbs.IsDismissible {
		backdropAttrs["data-dismissible"] = "true"
	}

	// Build container attributes
	containerAttrs := map[string]string{
		"class": "godin-modal-bottom-sheet-container",
		"style": mbs.buildContainerStyle(),
	}

	// Build bottom sheet attributes
	sheetAttrs := map[string]string{
		"class":      "godin-modal-bottom-sheet",
		"style":      mbs.buildSheetStyle(),
		"role":       "dialog",
		"aria-modal": "true",
	}

	// Add drag attributes if enabled
	if mbs.EnableDrag {
		sheetAttrs["data-draggable"] = "true"
	}

	// Merge with interactive widget attributes
	sheetAttrs = mbs.InteractiveWidget.MergeAttributes(sheetAttrs)

	// Build sheet content
	var contentParts []string

	// Add drag handle if enabled
	if mbs.ShowDragHandle {
		dragHandle := htmlRenderer.RenderElement("div", map[string]string{
			"class": "godin-modal-bottom-sheet-drag-handle",
			"style": mbs.buildDragHandleStyle(),
		}, "", true)
		contentParts = append(contentParts, dragHandle)
	}

	// Add child content
	if mbs.Child != nil {
		childHTML := htmlRenderer.RenderElement("div", map[string]string{
			"class": "godin-modal-bottom-sheet-content",
			"style": mbs.buildContentStyle(),
		}, mbs.Child.Render(ctx), false)
		contentParts = append(contentParts, childHTML)
	}

	// Build complete bottom sheet
	sheetContent := strings.Join(contentParts, "")
	sheetHTML := htmlRenderer.RenderElement("div", sheetAttrs, sheetContent, false)
	containerHTML := htmlRenderer.RenderElement("div", containerAttrs, sheetHTML, false)

	// Wrap with backdrop
	return htmlRenderer.RenderElement("div", backdropAttrs, containerHTML, false)
}

// buildBackdropStyle builds the CSS style for the modal backdrop
func (mbs *ModalBottomSheet) buildBackdropStyle() string {
	styles := []string{
		"position: fixed",
		"top: 0",
		"left: 0",
		"right: 0",
		"bottom: 0",
		"background-color: rgba(0, 0, 0, 0.5)",
		"z-index: 1000",
		"display: flex",
		"align-items: flex-end",
	}

	if mbs.IsDismissible {
		styles = append(styles, "cursor: pointer")
	}

	return strings.Join(styles, "; ")
}

// buildContainerStyle builds the CSS style for the modal container
func (mbs *ModalBottomSheet) buildContainerStyle() string {
	styles := []string{
		"width: 100%",
		"pointer-events: none", // Allow clicks to pass through to backdrop
	}
	return strings.Join(styles, "; ")
}

// buildSheetStyle builds the CSS style for the modal bottom sheet
func (mbs *ModalBottomSheet) buildSheetStyle() string {
	styles := []string{
		"background-color: white",
		"border-top-left-radius: 16px",
		"border-top-right-radius: 16px",
		"display: flex",
		"flex-direction: column",
		"pointer-events: auto",     // Re-enable pointer events for the sheet
		"transform: translateY(0)", // Initial position
		"transition: transform 0.3s cubic-bezier(0.4, 0.0, 0.2, 1)",
		"width: 100%",
	}

	if mbs.IsScrollControlled {
		styles = append(styles, "max-height: 95vh")
	} else {
		styles = append(styles, "max-height: 50vh")
	}

	if mbs.BackgroundColor != nil {
		styles = append(styles, fmt.Sprintf("background-color: %s", mbs.BackgroundColor.ToCSS()))
	}

	if mbs.Elevation > 0 {
		shadowValue := mbs.calculateBoxShadow(mbs.Elevation)
		styles = append(styles, fmt.Sprintf("box-shadow: %s", shadowValue))
	}

	if mbs.Shape != nil {
		styles = append(styles, fmt.Sprintf("border-radius: %s", mbs.Shape.ToCSS()))
	}

	return strings.Join(styles, "; ")
}

// buildContentStyle builds the CSS style for the modal content
func (mbs *ModalBottomSheet) buildContentStyle() string {
	styles := []string{
		"flex: 1",
		"overflow-y: auto",
		"padding: 16px",
	}
	return strings.Join(styles, "; ")
}

// buildDragHandleStyle builds the CSS style for the modal drag handle
func (mbs *ModalBottomSheet) buildDragHandleStyle() string {
	styles := []string{
		"width: 32px",
		"height: 4px",
		"background-color: rgba(0, 0, 0, 0.2)",
		"border-radius: 2px",
		"margin: 12px auto 8px auto",
		"cursor: grab",
	}
	return strings.Join(styles, "; ")
}

// calculateBoxShadow calculates box-shadow based on elevation for modal bottom sheet
func (mbs *ModalBottomSheet) calculateBoxShadow(elevation float64) string {
	// Same as regular bottom sheet but with upward shadows
	switch {
	case elevation <= 1:
		return "0px -2px 1px -1px rgba(0,0,0,0.2), 0px -1px 1px 0px rgba(0,0,0,0.14), 0px -1px 3px 0px rgba(0,0,0,0.12)"
	case elevation <= 2:
		return "0px -3px 1px -2px rgba(0,0,0,0.2), 0px -2px 2px 0px rgba(0,0,0,0.14), 0px -1px 5px 0px rgba(0,0,0,0.12)"
	case elevation <= 4:
		return "0px -2px 4px -1px rgba(0,0,0,0.2), 0px -4px 5px 0px rgba(0,0,0,0.14), 0px -1px 10px 0px rgba(0,0,0,0.12)"
	case elevation <= 8:
		return "0px -5px 5px -3px rgba(0,0,0,0.2), 0px -8px 10px 1px rgba(0,0,0,0.14), 0px -3px 14px 2px rgba(0,0,0,0.12)"
	case elevation <= 16:
		return "0px -8px 10px -5px rgba(0,0,0,0.2), 0px -16px 24px 2px rgba(0,0,0,0.14), 0px -6px 30px 5px rgba(0,0,0,0.12)"
	default:
		return "0px -11px 15px -7px rgba(0,0,0,0.2), 0px -24px 38px 3px rgba(0,0,0,0.14), 0px -9px 46px 8px rgba(0,0,0,0.12)"
	}
}

// ShowBottomSheet displays a bottom sheet and returns a sheet ID
func ShowBottomSheet(ctx *core.Context, bottomSheet core.Widget, options ...BottomSheetOptions) string {
	if ctx == nil || ctx.App == nil {
		return ""
	}

	// Get or create dialog manager
	dialogManager := ctx.App.DialogManager()
	if dialogManager == nil {
		return ""
	}

	// Use default options if none provided
	opts := BottomSheetOptions{
		IsModal:        false,
		IsDraggable:    true,
		EnableDrag:     true,
		ShowDragHandle: true,
	}
	if len(options) > 0 {
		opts = options[0]
	}

	return ShowBottomSheet(ctx, bottomSheet, opts)
}

// ShowModalBottomSheet displays a modal bottom sheet and returns a sheet ID
func ShowModalBottomSheet(ctx *core.Context, bottomSheet core.Widget, options ...BottomSheetOptions) string {
	if ctx == nil || ctx.App == nil {
		return ""
	}

	// Get or create dialog manager
	dialogManager := ctx.App.DialogManager()
	if dialogManager == nil {
		return ""
	}

	// Use modal options
	opts := BottomSheetOptions{
		IsModal:        true,
		IsDraggable:    true,
		EnableDrag:     true,
		ShowDragHandle: true,
	}
	if len(options) > 0 {
		opts = options[0]
		opts.IsModal = true // Force modal
	}

	return ShowBottomSheet(ctx, bottomSheet, opts)
}

// DismissBottomSheet dismisses a bottom sheet by ID
func DismissBottomSheet(ctx *core.Context, sheetID string) bool {
	if ctx == nil || ctx.App == nil {
		return false
	}

	dialogManager := ctx.App.DialogManager()
	if dialogManager == nil {
		return false
	}

	return DismissBottomSheet(ctx, sheetID)
}
