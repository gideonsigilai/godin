package widgets

import (
	"fmt"
	"strings"

	"github.com/gideonsigilai/godin/pkg/core"
	"github.com/gideonsigilai/godin/pkg/renderer"
)

// Dialog represents a modal dialog widget
type Dialog struct {
	InteractiveWidget
	Title           core.Widget
	Content         core.Widget
	Actions         []core.Widget
	BackgroundColor *core.Color
	Elevation       float64
	Shape           *BorderRadius
	InsetPadding    core.EdgeInsets
	ClipBehavior    ClipBehavior
}

// ClipBehavior represents how content should be clipped
type ClipBehavior string

const (
	ClipBehaviorNone                   ClipBehavior = "none"
	ClipBehaviorHardEdge               ClipBehavior = "hardEdge"
	ClipBehaviorAntiAlias              ClipBehavior = "antiAlias"
	ClipBehaviorAntiAliasWithSaveLayer ClipBehavior = "antiAliasWithSaveLayer"
)

// NewDialog creates a new Dialog widget
func NewDialog() *Dialog {
	return &Dialog{
		Elevation:    24.0,
		InsetPadding: core.NewEdgeInsetsSymmetric(40.0, 24.0),
		ClipBehavior: ClipBehaviorNone,
	}
}

// NewAlertDialog creates a new AlertDialog widget
func NewAlertDialog() *AlertDialog {
	return &AlertDialog{
		Elevation: &[]float64{24.0}[0],
		// TitlePadding:   core.EdgeInsets( &[]float64{8}[0]),
		// ContentPadding: core.NewEdgeInsetsAll(&[]float64{8}[0]),
		// ActionsPadding: core.NewEdgeInsetsAll(&[]float64{8}[0]),
		// ButtonPadding:  core.NewEdgeInsetsAll(&[]float64{8}[0]),
	}
}

// ToCSS converts BorderRadius to CSS border-radius property
func (br *BorderRadius) ToCSS() string {
	return fmt.Sprintf("%.1fpx %.1fpx %.1fpx %.1fpx", br.TopLeft, br.TopRight, br.BottomRight, br.BottomLeft)
}

// Render renders the Dialog widget
func (d *Dialog) Render(ctx *core.Context) string {
	if !d.InteractiveWidget.IsInitialized() {
		d.InteractiveWidget.Initialize(ctx)
		d.InteractiveWidget.SetWidgetType("Dialog")
	}

	htmlRenderer := renderer.NewHTMLRenderer()

	// Build dialog container attributes
	containerAttrs := map[string]string{
		"class": "godin-dialog-container",
		"style": d.buildContainerStyle(),
	}

	// Build dialog attributes
	dialogAttrs := map[string]string{
		"class":      "godin-dialog",
		"style":      d.buildDialogStyle(),
		"role":       "dialog",
		"aria-modal": "true",
	}

	// Merge with interactive widget attributes
	dialogAttrs = d.InteractiveWidget.MergeAttributes(dialogAttrs)

	// Build dialog content
	var contentParts []string

	// Add title if provided
	if d.Title != nil {
		titleHTML := htmlRenderer.RenderElement("div", map[string]string{
			"class": "godin-dialog-title",
			"style": "padding: 24px 24px 20px 24px; font-size: 20px; font-weight: 500;",
		}, d.Title.Render(ctx), false)
		contentParts = append(contentParts, titleHTML)
	}

	// Add content if provided
	if d.Content != nil {
		contentHTML := htmlRenderer.RenderElement("div", map[string]string{
			"class": "godin-dialog-content",
			"style": "padding: 20px 24px 24px 24px; flex: 1; overflow-y: auto;",
		}, d.Content.Render(ctx), false)
		contentParts = append(contentParts, contentHTML)
	}

	// Add actions if provided
	if len(d.Actions) > 0 {
		var actionHTMLs []string
		for _, action := range d.Actions {
			actionHTMLs = append(actionHTMLs, action.Render(ctx))
		}
		actionsHTML := htmlRenderer.RenderElement("div", map[string]string{
			"class": "godin-dialog-actions",
			"style": "padding: 8px; display: flex; justify-content: flex-end; gap: 8px;",
		}, strings.Join(actionHTMLs, ""), false)
		contentParts = append(contentParts, actionsHTML)
	}

	// Build complete dialog
	dialogContent := strings.Join(contentParts, "")
	dialogHTML := htmlRenderer.RenderElement("div", dialogAttrs, dialogContent, false)

	// Wrap in container
	return htmlRenderer.RenderElement("div", containerAttrs, dialogHTML, false)
}

// buildContainerStyle builds the CSS style for the dialog container
func (d *Dialog) buildContainerStyle() string {
	styles := []string{
		"position: fixed",
		"top: 0",
		"left: 0",
		"right: 0",
		"bottom: 0",
		"background-color: rgba(0, 0, 0, 0.5)",
		"display: flex",
		"align-items: center",
		"justify-content: center",
		"z-index: 1000",
		fmt.Sprintf("padding: %s", d.InsetPadding.ToCSS()),
	}
	return strings.Join(styles, "; ")
}

// buildDialogStyle builds the CSS style for the dialog
func (d *Dialog) buildDialogStyle() string {
	styles := []string{
		"background-color: white",
		"border-radius: 4px",
		"max-width: 560px",
		"max-height: calc(100vh - 80px)",
		"min-width: 280px",
		"display: flex",
		"flex-direction: column",
		"overflow: hidden",
	}

	if d.BackgroundColor != nil {
		styles = append(styles, fmt.Sprintf("background-color: %s", d.BackgroundColor.ToCSS()))
	}

	if d.Elevation > 0 {
		shadowValue := d.calculateBoxShadow(d.Elevation)
		styles = append(styles, fmt.Sprintf("box-shadow: %s", shadowValue))
	}

	if d.Shape != nil {
		styles = append(styles, fmt.Sprintf("border-radius: %s", d.Shape.ToCSS()))
	}

	return strings.Join(styles, "; ")
}

// calculateBoxShadow calculates box-shadow based on elevation
func (d *Dialog) calculateBoxShadow(elevation float64) string {
	// Material Design elevation shadows
	switch {
	case elevation <= 1:
		return "0px 2px 1px -1px rgba(0,0,0,0.2), 0px 1px 1px 0px rgba(0,0,0,0.14), 0px 1px 3px 0px rgba(0,0,0,0.12)"
	case elevation <= 2:
		return "0px 3px 1px -2px rgba(0,0,0,0.2), 0px 2px 2px 0px rgba(0,0,0,0.14), 0px 1px 5px 0px rgba(0,0,0,0.12)"
	case elevation <= 4:
		return "0px 2px 4px -1px rgba(0,0,0,0.2), 0px 4px 5px 0px rgba(0,0,0,0.14), 0px 1px 10px 0px rgba(0,0,0,0.12)"
	case elevation <= 8:
		return "0px 5px 5px -3px rgba(0,0,0,0.2), 0px 8px 10px 1px rgba(0,0,0,0.14), 0px 3px 14px 2px rgba(0,0,0,0.12)"
	case elevation <= 16:
		return "0px 8px 10px -5px rgba(0,0,0,0.2), 0px 16px 24px 2px rgba(0,0,0,0.14), 0px 6px 30px 5px rgba(0,0,0,0.12)"
	default:
		return "0px 11px 15px -7px rgba(0,0,0,0.2), 0px 24px 38px 3px rgba(0,0,0,0.14), 0px 9px 46px 8px rgba(0,0,0,0.12)"
	}
}

// showDialog displays a modal dialog and returns a dialog ID
func ShowDialog(ctx *core.Context, dialog core.Widget, options ...DialogOptions) string {
	if ctx == nil || ctx.App == nil {
		return ""
	}

	// Get or create dialog manager
	dialogManager := ctx.App.DialogManager()
	if dialogManager == nil {
		return ""
	}

	// Use default options if none provided
	opts := DialogOptions{
		BarrierDismissible: true,
	}
	if len(options) > 0 {
		opts = options[0]
	}

	return ShowDialog(ctx, dialog, opts)
}

// ShowAlertDialog is a convenience function for showing alert dialogs
func ShowAlertDialog(ctx *core.Context, title, content string, actions []core.Widget) string {
	alertDialog := NewAlertDialog()

	if title != "" {
		alertDialog.Title = Text{Data: title}
	}

	if content != "" {
		alertDialog.Content = Text{Data: content}
	}

	if len(actions) > 0 {
		alertDialog.Actions = actions
	}

	return ShowDialog(ctx, alertDialog)
}

// DismissDialog dismisses a dialog by ID
func DismissDialog(ctx *core.Context, dialogID string) bool {
	if ctx == nil || ctx.App == nil {
		return false
	}

	dialogManager := ctx.App.DialogManager()
	if dialogManager == nil {
		return false
	}

	return DismissDialog(ctx, dialogID)
}
