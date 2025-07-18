package widgets

import (
	"fmt"
	"godin-framework/pkg/core"
	"godin-framework/pkg/renderer"
)

// Dialog represents a dialog widget
type Dialog struct {
	HTMXWidget
	Title   string
	Content Widget
	Actions []Widget
	Open    bool
}

// Render renders the dialog as HTML
func (d Dialog) Render(ctx *core.Context) string {
	if !d.Open {
		return ""
	}

	htmlRenderer := renderer.NewHTMLRenderer()

	// Dialog overlay
	overlayAttrs := map[string]string{
		"class": "godin-dialog-overlay",
	}

	// Dialog container
	dialogAttrs := d.buildHTMXAttributes()
	dialogAttrs["class"] += " godin-dialog"

	if d.Style != "" {
		dialogAttrs["style"] = d.Style
	}

	var dialogContent string

	// Title
	if d.Title != "" {
		titleAttrs := map[string]string{"class": "godin-dialog-title"}
		dialogContent += htmlRenderer.RenderElement("h2", titleAttrs, d.Title, false)
	}

	// Content
	if d.Content != nil {
		contentAttrs := map[string]string{"class": "godin-dialog-content"}
		dialogContent += htmlRenderer.RenderElement("div", contentAttrs, d.Content.Render(ctx), false)
	}

	// Actions
	if len(d.Actions) > 0 {
		var actionElements []string
		for _, action := range d.Actions {
			if action != nil {
				actionElements = append(actionElements, action.Render(ctx))
			}
		}
		actionsAttrs := map[string]string{"class": "godin-dialog-actions"}
		dialogContent += htmlRenderer.RenderContainer("div", actionsAttrs, actionElements)
	}

	dialogHTML := htmlRenderer.RenderElement("div", dialogAttrs, dialogContent, false)
	overlayHTML := htmlRenderer.RenderElement("div", overlayAttrs, "", false)

	return overlayHTML + dialogHTML
}

// BottomSheet represents a bottom sheet widget
type BottomSheet struct {
	HTMXWidget
	Child  Widget
	Height string
	Open   bool
}

// Render renders the bottom sheet as HTML
func (bs BottomSheet) Render(ctx *core.Context) string {
	if !bs.Open {
		return ""
	}

	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := bs.buildHTMXAttributes()
	attrs["class"] += " godin-bottomsheet"

	// Build inline styles
	style := bs.Style
	if bs.Height != "" {
		style += "; height: " + bs.Height
	}

	if style != "" {
		attrs["style"] = style
	}

	content := ""
	if bs.Child != nil {
		content = bs.Child.Render(ctx)
	}

	return htmlRenderer.RenderElement("div", attrs, content, false)
}

// Snackbar represents a snackbar widget
type Snackbar struct {
	HTMXWidget
	Message  string
	Action   Widget
	Duration int
	Type     string // "info", "success", "warning", "error"
}

// Render renders the snackbar as HTML
func (s *Snackbar) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := s.buildHTMXAttributes()
	attrs["class"] += " godin-snackbar"

	if s.Type != "" {
		attrs["class"] += " godin-snackbar-" + s.Type
	}

	if s.Style != "" {
		attrs["style"] = s.Style
	}

	// Auto-hide after duration
	if s.Duration > 0 {
		attrs["data-duration"] = fmt.Sprintf("%d", s.Duration)
	}

	var content string

	// Message
	if s.Message != "" {
		messageAttrs := map[string]string{"class": "godin-snackbar-message"}
		content += htmlRenderer.RenderElement("span", messageAttrs, s.Message, false)
	}

	// Action
	if s.Action != nil {
		actionAttrs := map[string]string{"class": "godin-snackbar-action"}
		content += htmlRenderer.RenderElement("div", actionAttrs, s.Action.Render(ctx), false)
	}

	return htmlRenderer.RenderElement("div", attrs, content, false)
}

// Tooltip represents a tooltip widget
type Tooltip struct {
	HTMXWidget
	Child    Widget
	Message  string
	Position string // "top", "bottom", "left", "right"
}

// Render renders the tooltip as HTML
func (t *Tooltip) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := t.buildHTMXAttributes()
	attrs["class"] += " godin-tooltip"

	if t.Style != "" {
		attrs["style"] = t.Style
	}

	var content string

	// Child widget
	if t.Child != nil {
		content += t.Child.Render(ctx)
	}

	// Tooltip content
	if t.Message != "" {
		tooltipAttrs := map[string]string{
			"class": "godin-tooltip-content",
		}

		if t.Position != "" {
			tooltipAttrs["class"] += " godin-tooltip-" + t.Position
		}

		content += htmlRenderer.RenderElement("div", tooltipAttrs, t.Message, false)
	}

	return htmlRenderer.RenderElement("div", attrs, content, false)
}

// ProgressIndicator represents a progress indicator widget
type ProgressIndicator struct {
	HTMXWidget
	Value float64 // 0.0 to 1.0
	Type  string  // "linear", "circular"
}

// Render renders the progress indicator as HTML
func (pi *ProgressIndicator) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := pi.buildHTMXAttributes()

	if pi.Type == "circular" {
		attrs["class"] += " godin-progress-circular"

		if pi.Style != "" {
			attrs["style"] = pi.Style
		}

		return htmlRenderer.RenderElement("div", attrs, "", false)
	} else {
		// Linear progress
		attrs["class"] += " godin-progress-linear"

		if pi.Style != "" {
			attrs["style"] = pi.Style
		}

		// Progress bar
		barAttrs := map[string]string{
			"class": "godin-progress-linear-bar",
			"style": fmt.Sprintf("width: %.1f%%", pi.Value*100),
		}

		barHTML := htmlRenderer.RenderElement("div", barAttrs, "", false)
		return htmlRenderer.RenderElement("div", attrs, barHTML, false)
	}
}

// Modal represents a modal widget
type Modal struct {
	HTMXWidget
	Child       Widget
	Open        bool
	CloseOnTap  bool
	CloseButton bool
}

// Render renders the modal as HTML
func (m *Modal) Render(ctx *core.Context) string {
	if !m.Open {
		return ""
	}

	htmlRenderer := renderer.NewHTMLRenderer()

	// Modal overlay
	overlayAttrs := map[string]string{
		"class": "godin-modal-overlay",
	}

	if m.CloseOnTap {
		overlayAttrs["onclick"] = "this.parentElement.style.display='none'"
	}

	// Modal container
	modalAttrs := m.buildHTMXAttributes()
	modalAttrs["class"] += " godin-modal"

	if m.Style != "" {
		modalAttrs["style"] = m.Style
	}

	var modalContent string

	// Close button
	if m.CloseButton {
		closeAttrs := map[string]string{
			"class":   "godin-modal-close",
			"onclick": "this.closest('.godin-modal').style.display='none'",
		}
		modalContent += htmlRenderer.RenderElement("button", closeAttrs, "×", false)
	}

	// Child content
	if m.Child != nil {
		modalContent += m.Child.Render(ctx)
	}

	modalHTML := htmlRenderer.RenderElement("div", modalAttrs, modalContent, false)
	overlayHTML := htmlRenderer.RenderElement("div", overlayAttrs, "", false)

	// Container for both overlay and modal
	containerAttrs := map[string]string{
		"class": "godin-modal-container",
		"style": "position: fixed; top: 0; left: 0; right: 0; bottom: 0; z-index: 1000;",
	}

	return htmlRenderer.RenderElement("div", containerAttrs, overlayHTML+modalHTML, false)
}

// Popover represents a popover widget
type Popover struct {
	HTMXWidget
	Child    Widget
	Content  Widget
	Trigger  string // "click", "hover"
	Position string // "top", "bottom", "left", "right"
}

// Render renders the popover as HTML
func (p *Popover) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := p.buildHTMXAttributes()
	attrs["class"] += " godin-popover"

	if p.Style != "" {
		attrs["style"] = p.Style
	}

	var content string

	// Trigger element
	if p.Child != nil {
		triggerAttrs := map[string]string{
			"class": "godin-popover-trigger",
		}

		// Add trigger event
		if p.Trigger == "hover" {
			triggerAttrs["onmouseenter"] = "this.nextElementSibling.style.display='block'"
			triggerAttrs["onmouseleave"] = "this.nextElementSibling.style.display='none'"
		} else {
			triggerAttrs["onclick"] = "this.nextElementSibling.style.display=this.nextElementSibling.style.display==='block'?'none':'block'"
		}

		content += htmlRenderer.RenderElement("div", triggerAttrs, p.Child.Render(ctx), false)
	}

	// Popover content
	if p.Content != nil {
		popoverAttrs := map[string]string{
			"class": "godin-popover-content",
			"style": "display: none;",
		}

		if p.Position != "" {
			popoverAttrs["class"] += " godin-popover-" + p.Position
		}

		content += htmlRenderer.RenderElement("div", popoverAttrs, p.Content.Render(ctx), false)
	}

	return htmlRenderer.RenderElement("div", attrs, content, false)
}
