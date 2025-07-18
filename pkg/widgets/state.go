package widgets

import (
	"fmt"

	"github.com/gideonsigilai/godin/pkg/core"
	"github.com/gideonsigilai/godin/pkg/state"
)

// ValueListenableBuilder represents a widget that rebuilds when a value changes
type ValueListenableBuilder struct {
	HTMXWidget
	ValueListenable *state.ValueListenable
	Builder         func(value interface{}) Widget
}

// Render renders the value listenable builder as HTML
func (vlb *ValueListenableBuilder) Render(ctx *core.Context) string {
	if vlb.ValueListenable == nil || vlb.Builder == nil {
		return ""
	}

	// Get current value and build widget
	value := vlb.ValueListenable.GetValue()
	widget := vlb.Builder(value)

	if widget == nil {
		return ""
	}

	return widget.Render(ctx)
}

// StreamBuilder represents a widget that rebuilds when stream data changes
type StreamBuilder struct {
	HTMXWidget
	Stream  chan interface{}
	Builder func(data interface{}) Widget
}

// Render renders the stream builder as HTML
func (sb *StreamBuilder) Render(ctx *core.Context) string {
	if sb.Stream == nil || sb.Builder == nil {
		return ""
	}

	// For server-side rendering, we'll use WebSocket for real-time updates
	// This is a simplified implementation - in practice, you'd want to
	// register the stream with the WebSocket manager

	// For now, just render with nil data
	widget := sb.Builder(nil)
	if widget == nil {
		return ""
	}

	return widget.Render(ctx)
}

// FutureBuilder represents a widget that rebuilds when future completes
type FutureBuilder struct {
	HTMXWidget
	Future  func() interface{}
	Builder func(data interface{}, loading bool, err error) Widget
}

// Render renders the future builder as HTML
func (fb *FutureBuilder) Render(ctx *core.Context) string {
	if fb.Future == nil || fb.Builder == nil {
		return ""
	}

	// For server-side rendering, we'll execute the future immediately
	// In a real implementation, you might want to handle this asynchronously
	var data interface{}
	var err error

	// Execute future
	data = fb.Future()

	widget := fb.Builder(data, false, err)
	if widget == nil {
		return ""
	}

	return widget.Render(ctx)
}

// StateBuilder represents a widget that rebuilds when state changes
type StateBuilder struct {
	HTMXWidget
	State   *state.State
	Builder func(state *state.State) Widget
}

// Render renders the state builder as HTML
func (sb *StateBuilder) Render(ctx *core.Context) string {
	if sb.State == nil || sb.Builder == nil {
		return ""
	}

	widget := sb.Builder(sb.State)
	if widget == nil {
		return ""
	}

	return widget.Render(ctx)
}

// Consumer represents a widget that consumes state changes
type Consumer struct {
	HTMXWidget
	StateKey string
	Builder  func(value interface{}) Widget
}

// Render renders the consumer as HTML
func (c *Consumer) Render(ctx *core.Context) string {
	if c.StateKey == "" || c.Builder == nil {
		return ""
	}

	// Get state from context (assuming it's available)
	stateManager := ctx.App.State()
	value := stateManager.Get(c.StateKey)

	widget := c.Builder(value)
	if widget == nil {
		return ""
	}

	// Wrap the widget in a container with state tracking attributes
	containerHTML := fmt.Sprintf(`<div data-state-key="%s" data-state-endpoint="/api/state/%s">%s</div>`,
		c.StateKey, c.StateKey, widget.Render(ctx))

	return containerHTML
}

// Provider represents a widget that provides state to its children
type Provider struct {
	HTMXWidget
	StateKey string
	Value    interface{}
	Child    Widget
}

// Render renders the provider as HTML
func (p *Provider) Render(ctx *core.Context) string {
	if p.StateKey == "" || p.Child == nil {
		return ""
	}

	// Set state in context
	stateManager := ctx.App.State()
	stateManager.Set(p.StateKey, p.Value)

	return p.Child.Render(ctx)
}

// Selector represents a widget that selects specific parts of state
type Selector struct {
	HTMXWidget
	StateKey string
	Selector func(state interface{}) interface{}
	Builder  func(selected interface{}) Widget
}

// Render renders the selector as HTML
func (s *Selector) Render(ctx *core.Context) string {
	if s.StateKey == "" || s.Selector == nil || s.Builder == nil {
		return ""
	}

	// Get state from context
	stateManager := ctx.App.State()
	state := stateManager.Get(s.StateKey)

	// Select specific part of state
	selected := s.Selector(state)

	widget := s.Builder(selected)
	if widget == nil {
		return ""
	}

	return widget.Render(ctx)
}

// ChangeNotifierProvider represents a widget that provides a change notifier
type ChangeNotifierProvider struct {
	HTMXWidget
	Notifier ChangeNotifier
	Child    Widget
}

// ChangeNotifier interface for objects that can notify about changes
type ChangeNotifier interface {
	AddListener(func())
	RemoveListener(func())
	NotifyListeners()
}

// Render renders the change notifier provider as HTML
func (cnp *ChangeNotifierProvider) Render(ctx *core.Context) string {
	if cnp.Notifier == nil || cnp.Child == nil {
		return ""
	}

	// In a real implementation, you'd register the notifier with the context
	// For now, just render the child
	return cnp.Child.Render(ctx)
}

// AnimatedBuilder represents a widget that rebuilds on animation changes
type AnimatedBuilder struct {
	HTMXWidget
	Animation Animation
	Builder   func(animation Animation) Widget
}

// Animation interface for animation objects
type Animation interface {
	GetValue() float64
	AddListener(func())
	RemoveListener(func())
}

// Render renders the animated builder as HTML
func (ab *AnimatedBuilder) Render(ctx *core.Context) string {
	if ab.Animation == nil || ab.Builder == nil {
		return ""
	}

	widget := ab.Builder(ab.Animation)
	if widget == nil {
		return ""
	}

	return widget.Render(ctx)
}

// LayoutBuilderConstraints represents layout constraints for LayoutBuilder
type LayoutBuilderConstraints struct {
	MinWidth  float64
	MaxWidth  float64
	MinHeight float64
	MaxHeight float64
}

// LayoutBuilder represents a widget that builds based on layout constraints
type LayoutBuilder struct {
	HTMXWidget
	Builder func(constraints LayoutBuilderConstraints) Widget
}

// Render renders the layout builder as HTML
func (lb *LayoutBuilder) Render(ctx *core.Context) string {
	if lb.Builder == nil {
		return ""
	}

	// For server-side rendering, use default constraints
	constraints := LayoutBuilderConstraints{
		MinWidth:  0,
		MaxWidth:  1200, // Default max width
		MinHeight: 0,
		MaxHeight: 800, // Default max height
	}

	widget := lb.Builder(constraints)
	if widget == nil {
		return ""
	}

	return widget.Render(ctx)
}
