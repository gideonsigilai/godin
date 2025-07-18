package renderer

import (
	"fmt"
	"strings"
)

// HTMXAttributes represents HTMX attributes for widgets
type HTMXAttributes struct {
	Get     string // hx-get
	Post    string // hx-post
	Put     string // hx-put
	Delete  string // hx-delete
	Patch   string // hx-patch
	Target  string // hx-target
	Swap    string // hx-swap
	Trigger string // hx-trigger
	Vals    string // hx-vals
	Headers string // hx-headers
	Include string // hx-include
	Confirm string // hx-confirm
	Boost   bool   // hx-boost
	PushURL bool   // hx-push-url
}

// HTMXRenderer handles HTMX attribute generation
type HTMXRenderer struct{}

// NewHTMXRenderer creates a new HTMX renderer
func NewHTMXRenderer() *HTMXRenderer {
	return &HTMXRenderer{}
}

// RenderAttributes converts HTMXAttributes to HTML attributes map
func (hr *HTMXRenderer) RenderAttributes(htmx HTMXAttributes) map[string]string {
	attrs := make(map[string]string)
	
	if htmx.Get != "" {
		attrs["hx-get"] = htmx.Get
	}
	if htmx.Post != "" {
		attrs["hx-post"] = htmx.Post
	}
	if htmx.Put != "" {
		attrs["hx-put"] = htmx.Put
	}
	if htmx.Delete != "" {
		attrs["hx-delete"] = htmx.Delete
	}
	if htmx.Patch != "" {
		attrs["hx-patch"] = htmx.Patch
	}
	if htmx.Target != "" {
		attrs["hx-target"] = htmx.Target
	}
	if htmx.Swap != "" {
		attrs["hx-swap"] = htmx.Swap
	}
	if htmx.Trigger != "" {
		attrs["hx-trigger"] = htmx.Trigger
	}
	if htmx.Vals != "" {
		attrs["hx-vals"] = htmx.Vals
	}
	if htmx.Headers != "" {
		attrs["hx-headers"] = htmx.Headers
	}
	if htmx.Include != "" {
		attrs["hx-include"] = htmx.Include
	}
	if htmx.Confirm != "" {
		attrs["hx-confirm"] = htmx.Confirm
	}
	if htmx.Boost {
		attrs["hx-boost"] = "true"
	}
	if htmx.PushURL {
		attrs["hx-push-url"] = "true"
	}
	
	return attrs
}

// Common HTMX patterns and helpers

// LoadMore creates HTMX attributes for load more functionality
func (hr *HTMXRenderer) LoadMore(endpoint, target string) HTMXAttributes {
	return HTMXAttributes{
		Get:     endpoint,
		Target:  target,
		Swap:    "beforeend",
		Trigger: "click",
	}
}

// InfiniteScroll creates HTMX attributes for infinite scroll
func (hr *HTMXRenderer) InfiniteScroll(endpoint, target string) HTMXAttributes {
	return HTMXAttributes{
		Get:     endpoint,
		Target:  target,
		Swap:    "beforeend",
		Trigger: "revealed",
	}
}

// FormSubmit creates HTMX attributes for form submission
func (hr *HTMXRenderer) FormSubmit(endpoint, target, method string) HTMXAttributes {
	htmx := HTMXAttributes{
		Target:  target,
		Trigger: "submit",
	}
	
	switch strings.ToUpper(method) {
	case "POST":
		htmx.Post = endpoint
	case "PUT":
		htmx.Put = endpoint
	case "DELETE":
		htmx.Delete = endpoint
	case "PATCH":
		htmx.Patch = endpoint
	default:
		htmx.Post = endpoint
	}
	
	return htmx
}

// LiveSearch creates HTMX attributes for live search
func (hr *HTMXRenderer) LiveSearch(endpoint, target string) HTMXAttributes {
	return HTMXAttributes{
		Get:     endpoint,
		Target:  target,
		Trigger: "keyup changed delay:300ms",
		Include: "input[name='search']",
	}
}

// ToggleContent creates HTMX attributes for content toggling
func (hr *HTMXRenderer) ToggleContent(endpoint, target string) HTMXAttributes {
	return HTMXAttributes{
		Get:    endpoint,
		Target: target,
		Swap:   "outerHTML",
	}
}

// DeleteWithConfirm creates HTMX attributes for delete with confirmation
func (hr *HTMXRenderer) DeleteWithConfirm(endpoint, target, message string) HTMXAttributes {
	return HTMXAttributes{
		Delete:  endpoint,
		Target:  target,
		Swap:    "outerHTML",
		Confirm: message,
	}
}

// UpdateOnChange creates HTMX attributes for updating on input change
func (hr *HTMXRenderer) UpdateOnChange(endpoint, target string) HTMXAttributes {
	return HTMXAttributes{
		Post:    endpoint,
		Target:  target,
		Trigger: "change",
	}
}

// PollingUpdate creates HTMX attributes for polling updates
func (hr *HTMXRenderer) PollingUpdate(endpoint, target, interval string) HTMXAttributes {
	return HTMXAttributes{
		Get:     endpoint,
		Target:  target,
		Trigger: fmt.Sprintf("every %s", interval),
	}
}

// SwapOptions provides common swap options
type SwapOptions struct {
	Type     string // innerHTML, outerHTML, beforebegin, afterbegin, beforeend, afterend
	Settle   string // settle time
	Scroll   string // scroll behavior
	Show     string // show behavior
	Focus    bool   // focus-scroll
	Transition bool // view-transition
}

// BuildSwapString builds a swap string from options
func (hr *HTMXRenderer) BuildSwapString(opts SwapOptions) string {
	parts := []string{opts.Type}
	
	if opts.Settle != "" {
		parts = append(parts, fmt.Sprintf("settle:%s", opts.Settle))
	}
	if opts.Scroll != "" {
		parts = append(parts, fmt.Sprintf("scroll:%s", opts.Scroll))
	}
	if opts.Show != "" {
		parts = append(parts, fmt.Sprintf("show:%s", opts.Show))
	}
	if opts.Focus {
		parts = append(parts, "focus-scroll:true")
	}
	if opts.Transition {
		parts = append(parts, "transition:true")
	}
	
	return strings.Join(parts, " ")
}

// TriggerOptions provides common trigger options
type TriggerOptions struct {
	Event    string // click, change, keyup, etc.
	Modifier string // once, changed, delay, throttle, from, target, consume, queue
	Filter   string // CSS selector filter
	Delay    string // delay time
	Throttle string // throttle time
}

// BuildTriggerString builds a trigger string from options
func (hr *HTMXRenderer) BuildTriggerString(opts TriggerOptions) string {
	trigger := opts.Event
	
	if opts.Modifier != "" {
		trigger += " " + opts.Modifier
	}
	if opts.Filter != "" {
		trigger += fmt.Sprintf(" from:%s", opts.Filter)
	}
	if opts.Delay != "" {
		trigger += fmt.Sprintf(" delay:%s", opts.Delay)
	}
	if opts.Throttle != "" {
		trigger += fmt.Sprintf(" throttle:%s", opts.Throttle)
	}
	
	return trigger
}

// Common HTMX response headers
const (
	HXTrigger         = "HX-Trigger"
	HXTriggerAfterSettle = "HX-Trigger-After-Settle"
	HXTriggerAfterSwap   = "HX-Trigger-After-Swap"
	HXLocation        = "HX-Location"
	HXPushURL         = "HX-Push-Url"
	HXRedirect        = "HX-Redirect"
	HXRefresh         = "HX-Refresh"
	HXReplaceURL      = "HX-Replace-Url"
	HXReswap          = "HX-Reswap"
	HXRetarget        = "HX-Retarget"
	HXReselect        = "HX-Reselect"
)

// HTMXResponse helps build HTMX response headers
type HTMXResponse struct {
	headers map[string]string
}

// NewHTMXResponse creates a new HTMX response helper
func NewHTMXResponse() *HTMXResponse {
	return &HTMXResponse{
		headers: make(map[string]string),
	}
}

// Trigger sets the HX-Trigger header
func (hr *HTMXResponse) Trigger(event string) *HTMXResponse {
	hr.headers[HXTrigger] = event
	return hr
}

// Location sets the HX-Location header
func (hr *HTMXResponse) Location(url string) *HTMXResponse {
	hr.headers[HXLocation] = url
	return hr
}

// Redirect sets the HX-Redirect header
func (hr *HTMXResponse) Redirect(url string) *HTMXResponse {
	hr.headers[HXRedirect] = url
	return hr
}

// Refresh sets the HX-Refresh header
func (hr *HTMXResponse) Refresh() *HTMXResponse {
	hr.headers[HXRefresh] = "true"
	return hr
}

// GetHeaders returns all HTMX headers
func (hr *HTMXResponse) GetHeaders() map[string]string {
	return hr.headers
}
