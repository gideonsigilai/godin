package core

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
)

// HTMXIntegrator handles automatic HTMX attribute generation
type HTMXIntegrator struct {
	endpointPrefix string
	csrfToken      string
	baseURL        string
	headers        map[string]string
	swapStrategy   string
	targetStrategy string
}

// HTMXConfig represents configuration for HTMX integration
type HTMXConfig struct {
	EndpointPrefix string            // Prefix for generated endpoints (default: "/api/callbacks")
	CSRFToken      string            // CSRF token for security
	BaseURL        string            // Base URL for endpoints
	Headers        map[string]string // Additional headers to include
	SwapStrategy   string            // Default swap strategy (none, innerHTML, outerHTML, etc.)
	TargetStrategy string            // Default target strategy
}

// NewHTMXIntegrator creates a new HTMX integrator with configuration
func NewHTMXIntegrator(config *HTMXConfig) *HTMXIntegrator {
	integrator := &HTMXIntegrator{
		endpointPrefix: "/api/callbacks",
		swapStrategy:   "none",
		targetStrategy: "",
		headers:        make(map[string]string),
	}

	if config != nil {
		if config.EndpointPrefix != "" {
			integrator.endpointPrefix = config.EndpointPrefix
		}
		integrator.csrfToken = config.CSRFToken
		integrator.baseURL = config.BaseURL
		if config.SwapStrategy != "" {
			integrator.swapStrategy = config.SwapStrategy
		}
		integrator.targetStrategy = config.TargetStrategy

		// Copy headers
		for k, v := range config.Headers {
			integrator.headers[k] = v
		}
	}

	return integrator
}

// GenerateClickHandler generates HTMX attributes for click events
func (hi *HTMXIntegrator) GenerateClickHandler(callbackID string) map[string]string {
	attrs := make(map[string]string)
	endpoint := hi.buildEndpointURL(callbackID)

	attrs["hx-post"] = endpoint
	attrs["hx-trigger"] = "click"
	attrs["hx-swap"] = hi.swapStrategy

	// Add target if specified
	if hi.targetStrategy != "" {
		attrs["hx-target"] = hi.targetStrategy
	}

	// Add headers
	hi.addHeaders(attrs)

	// Add CSRF token if available
	hi.addCSRFToken(attrs)

	return attrs
}

// GenerateChangeHandler generates HTMX attributes for change events
func (hi *HTMXIntegrator) GenerateChangeHandler(callbackID string) map[string]string {
	attrs := make(map[string]string)
	endpoint := hi.buildEndpointURL(callbackID)

	attrs["hx-post"] = endpoint
	attrs["hx-trigger"] = "change"
	attrs["hx-include"] = "this"
	attrs["hx-swap"] = hi.swapStrategy

	// Add target if specified
	if hi.targetStrategy != "" {
		attrs["hx-target"] = hi.targetStrategy
	}

	// Add headers
	hi.addHeaders(attrs)

	// Add CSRF token if available
	hi.addCSRFToken(attrs)

	return attrs
}

// GenerateSubmitHandler generates HTMX attributes for submit events
func (hi *HTMXIntegrator) GenerateSubmitHandler(callbackID string) map[string]string {
	attrs := make(map[string]string)
	endpoint := hi.buildEndpointURL(callbackID)

	attrs["hx-post"] = endpoint
	attrs["hx-trigger"] = "submit"
	attrs["hx-include"] = "this"
	attrs["hx-swap"] = hi.swapStrategy

	// Add target if specified
	if hi.targetStrategy != "" {
		attrs["hx-target"] = hi.targetStrategy
	}

	// Add headers
	hi.addHeaders(attrs)

	// Add CSRF token if available
	hi.addCSRFToken(attrs)

	return attrs
}

// GenerateInputHandler generates HTMX attributes for input events (real-time text changes)
func (hi *HTMXIntegrator) GenerateInputHandler(callbackID string) map[string]string {
	attrs := make(map[string]string)
	endpoint := hi.buildEndpointURL(callbackID)

	attrs["hx-post"] = endpoint
	attrs["hx-trigger"] = "input changed delay:300ms" // Debounce input events
	attrs["hx-include"] = "this"
	attrs["hx-swap"] = hi.swapStrategy

	// Add target if specified
	if hi.targetStrategy != "" {
		attrs["hx-target"] = hi.targetStrategy
	}

	// Add headers
	hi.addHeaders(attrs)

	// Add CSRF token if available
	hi.addCSRFToken(attrs)

	return attrs
}

// GenerateKeyHandler generates HTMX attributes for keyboard events
func (hi *HTMXIntegrator) GenerateKeyHandler(callbackID string, keyCode int) map[string]string {
	attrs := make(map[string]string)
	endpoint := hi.buildEndpointURL(callbackID)

	attrs["hx-post"] = endpoint
	attrs["hx-trigger"] = fmt.Sprintf("keyup[keyCode==%d]", keyCode)
	attrs["hx-include"] = "this"
	attrs["hx-swap"] = hi.swapStrategy

	// Add target if specified
	if hi.targetStrategy != "" {
		attrs["hx-target"] = hi.targetStrategy
	}

	// Add headers
	hi.addHeaders(attrs)

	// Add CSRF token if available
	hi.addCSRFToken(attrs)

	return attrs
}

// GenerateEnterKeyHandler generates HTMX attributes for Enter key events
func (hi *HTMXIntegrator) GenerateEnterKeyHandler(callbackID string) map[string]string {
	return hi.GenerateKeyHandler(callbackID, 13) // Enter key code
}

// GenerateFocusHandler generates HTMX attributes for focus events
func (hi *HTMXIntegrator) GenerateFocusHandler(callbackID string) map[string]string {
	attrs := make(map[string]string)
	endpoint := hi.buildEndpointURL(callbackID)

	attrs["hx-post"] = endpoint
	attrs["hx-trigger"] = "focus"
	attrs["hx-include"] = "this"
	attrs["hx-swap"] = hi.swapStrategy

	// Add target if specified
	if hi.targetStrategy != "" {
		attrs["hx-target"] = hi.targetStrategy
	}

	// Add headers
	hi.addHeaders(attrs)

	// Add CSRF token if available
	hi.addCSRFToken(attrs)

	return attrs
}

// GenerateBlurHandler generates HTMX attributes for blur events
func (hi *HTMXIntegrator) GenerateBlurHandler(callbackID string) map[string]string {
	attrs := make(map[string]string)
	endpoint := hi.buildEndpointURL(callbackID)

	attrs["hx-post"] = endpoint
	attrs["hx-trigger"] = "blur"
	attrs["hx-include"] = "this"
	attrs["hx-swap"] = hi.swapStrategy

	// Add target if specified
	if hi.targetStrategy != "" {
		attrs["hx-target"] = hi.targetStrategy
	}

	// Add headers
	hi.addHeaders(attrs)

	// Add CSRF token if available
	hi.addCSRFToken(attrs)

	return attrs
}

// GenerateHoverHandler generates HTMX attributes for hover events
func (hi *HTMXIntegrator) GenerateHoverHandler(callbackID string) map[string]string {
	attrs := make(map[string]string)
	endpoint := hi.buildEndpointURL(callbackID)

	attrs["hx-post"] = endpoint
	attrs["hx-trigger"] = "mouseenter, mouseleave"
	attrs["hx-swap"] = hi.swapStrategy

	// Add target if specified
	if hi.targetStrategy != "" {
		attrs["hx-target"] = hi.targetStrategy
	}

	// Add headers
	hi.addHeaders(attrs)

	// Add CSRF token if available
	hi.addCSRFToken(attrs)

	return attrs
}

// GenerateDoubleClickHandler generates HTMX attributes for double-click events
func (hi *HTMXIntegrator) GenerateDoubleClickHandler(callbackID string) map[string]string {
	attrs := make(map[string]string)
	endpoint := hi.buildEndpointURL(callbackID)

	attrs["hx-post"] = endpoint
	attrs["hx-trigger"] = "dblclick"
	attrs["hx-swap"] = hi.swapStrategy

	// Add target if specified
	if hi.targetStrategy != "" {
		attrs["hx-target"] = hi.targetStrategy
	}

	// Add headers
	hi.addHeaders(attrs)

	// Add CSRF token if available
	hi.addCSRFToken(attrs)

	return attrs
}

// GenerateContextMenuHandler generates HTMX attributes for context menu (right-click) events
func (hi *HTMXIntegrator) GenerateContextMenuHandler(callbackID string) map[string]string {
	attrs := make(map[string]string)
	endpoint := hi.buildEndpointURL(callbackID)

	attrs["hx-post"] = endpoint
	attrs["hx-trigger"] = "contextmenu"
	attrs["hx-swap"] = hi.swapStrategy

	// Add target if specified
	if hi.targetStrategy != "" {
		attrs["hx-target"] = hi.targetStrategy
	}

	// Add headers
	hi.addHeaders(attrs)

	// Add CSRF token if available
	hi.addCSRFToken(attrs)

	return attrs
}

// GenerateCustomHandler generates HTMX attributes for custom events
func (hi *HTMXIntegrator) GenerateCustomHandler(event, callbackID string) map[string]string {
	attrs := make(map[string]string)
	endpoint := hi.buildEndpointURL(callbackID)

	attrs["hx-post"] = endpoint
	attrs["hx-trigger"] = event
	attrs["hx-swap"] = hi.swapStrategy

	// Add target if specified
	if hi.targetStrategy != "" {
		attrs["hx-target"] = hi.targetStrategy
	}

	// Add headers
	hi.addHeaders(attrs)

	// Add CSRF token if available
	hi.addCSRFToken(attrs)

	return attrs
}

// GeneratePollingHandler generates HTMX attributes for polling/periodic updates
func (hi *HTMXIntegrator) GeneratePollingHandler(callbackID string, intervalMs int) map[string]string {
	attrs := make(map[string]string)
	endpoint := hi.buildEndpointURL(callbackID)

	attrs["hx-get"] = endpoint // Use GET for polling
	attrs["hx-trigger"] = fmt.Sprintf("every %dms", intervalMs)
	attrs["hx-swap"] = "innerHTML" // Usually want to update content for polling

	// Add target if specified
	if hi.targetStrategy != "" {
		attrs["hx-target"] = hi.targetStrategy
	}

	// Add headers
	hi.addHeaders(attrs)

	// Add CSRF token if available
	hi.addCSRFToken(attrs)

	return attrs
}

// GenerateWebSocketHandler generates HTMX attributes for WebSocket integration
func (hi *HTMXIntegrator) GenerateWebSocketHandler(wsEndpoint string) map[string]string {
	attrs := make(map[string]string)

	// HTMX WebSocket extension attributes
	attrs["hx-ws"] = fmt.Sprintf("connect:%s", wsEndpoint)

	return attrs
}

// buildEndpointURL builds the complete endpoint URL
func (hi *HTMXIntegrator) buildEndpointURL(callbackID string) string {
	endpoint := fmt.Sprintf("%s/%s", hi.endpointPrefix, callbackID)

	if hi.baseURL != "" {
		baseURL, err := url.Parse(hi.baseURL)
		if err == nil {
			endpointURL, err := url.Parse(endpoint)
			if err == nil {
				return baseURL.ResolveReference(endpointURL).String()
			}
		}
	}

	return endpoint
}

// addHeaders adds custom headers to HTMX attributes
func (hi *HTMXIntegrator) addHeaders(attrs map[string]string) {
	if len(hi.headers) == 0 {
		return
	}

	var headerPairs []string
	for key, value := range hi.headers {
		headerPairs = append(headerPairs, fmt.Sprintf("%s:%s", key, value))
	}

	if len(headerPairs) > 0 {
		attrs["hx-headers"] = fmt.Sprintf("{%s}", strings.Join(headerPairs, ","))
	}
}

// addCSRFToken adds CSRF token to HTMX attributes
func (hi *HTMXIntegrator) addCSRFToken(attrs map[string]string) {
	if hi.csrfToken == "" {
		return
	}

	// Add CSRF token as a header
	csrfHeader := fmt.Sprintf("\"X-CSRF-Token\":\"%s\"", hi.csrfToken)

	if existingHeaders, exists := attrs["hx-headers"]; exists {
		// Merge with existing headers
		existingHeaders = strings.TrimSuffix(existingHeaders, "}")
		attrs["hx-headers"] = fmt.Sprintf("%s,%s}", existingHeaders, csrfHeader)
	} else {
		attrs["hx-headers"] = fmt.Sprintf("{%s}", csrfHeader)
	}
}

// SetCSRFToken sets the CSRF token for security
func (hi *HTMXIntegrator) SetCSRFToken(token string) {
	hi.csrfToken = token
}

// SetSwapStrategy sets the default swap strategy
func (hi *HTMXIntegrator) SetSwapStrategy(strategy string) {
	hi.swapStrategy = strategy
}

// SetTargetStrategy sets the default target strategy
func (hi *HTMXIntegrator) SetTargetStrategy(target string) {
	hi.targetStrategy = target
}

// AddHeader adds a custom header
func (hi *HTMXIntegrator) AddHeader(key, value string) {
	hi.headers[key] = value
}

// RemoveHeader removes a custom header
func (hi *HTMXIntegrator) RemoveHeader(key string) {
	delete(hi.headers, key)
}

// GetConfig returns the current configuration
func (hi *HTMXIntegrator) GetConfig() HTMXConfig {
	return HTMXConfig{
		EndpointPrefix: hi.endpointPrefix,
		CSRFToken:      hi.csrfToken,
		BaseURL:        hi.baseURL,
		Headers:        hi.headers,
		SwapStrategy:   hi.swapStrategy,
		TargetStrategy: hi.targetStrategy,
	}
}

// GenerateCSRFToken generates a new CSRF token
func GenerateCSRFToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// ValidateCSRFToken validates a CSRF token (simplified implementation)
func ValidateCSRFToken(provided, expected string) bool {
	return provided == expected && provided != ""
}

// HTMXErrorHandler represents an error handler for HTMX requests
type HTMXErrorHandler struct {
	OnClientError  func(statusCode int, message string) string // 4xx errors
	OnServerError  func(statusCode int, message string) string // 5xx errors
	OnNetworkError func(message string) string                 // Network errors
}

// GenerateErrorHandlingAttributes generates HTMX attributes for error handling
func (hi *HTMXIntegrator) GenerateErrorHandlingAttributes(errorHandler *HTMXErrorHandler) map[string]string {
	attrs := make(map[string]string)

	if errorHandler != nil {
		// Add error handling events
		attrs["hx-on::after-request"] = "handleHTMXResponse(event)"
		attrs["hx-on::response-error"] = "handleHTMXError(event)"
		attrs["hx-on::send-error"] = "handleHTMXNetworkError(event)"
	}

	return attrs
}

// GenerateLoadingIndicatorAttributes generates HTMX attributes for loading indicators
func (hi *HTMXIntegrator) GenerateLoadingIndicatorAttributes(indicatorSelector string) map[string]string {
	attrs := make(map[string]string)

	if indicatorSelector != "" {
		attrs["hx-indicator"] = indicatorSelector
	}

	return attrs
}

// GenerateProgressAttributes generates HTMX attributes for progress tracking
func (hi *HTMXIntegrator) GenerateProgressAttributes() map[string]string {
	attrs := make(map[string]string)

	// Add progress tracking events
	attrs["hx-on::htmx:xhr:progress"] = "updateProgress(event)"
	attrs["hx-on::htmx:xhr:loadstart"] = "startProgress(event)"
	attrs["hx-on::htmx:xhr:loadend"] = "endProgress(event)"

	return attrs
}

// MergeAttributes merges multiple attribute maps
func MergeHTMXAttributes(attrMaps ...map[string]string) map[string]string {
	result := make(map[string]string)

	for _, attrs := range attrMaps {
		for key, value := range attrs {
			// Special handling for certain attributes that can be combined
			if existingValue, exists := result[key]; exists {
				switch key {
				case "hx-trigger":
					// Combine triggers with comma
					result[key] = fmt.Sprintf("%s, %s", existingValue, value)
				case "hx-headers":
					// Merge JSON headers
					result[key] = mergeJSONHeaders(existingValue, value)
				default:
					// For other attributes, the last one wins
					result[key] = value
				}
			} else {
				result[key] = value
			}
		}
	}

	return result
}

// mergeJSONHeaders merges two JSON header strings
func mergeJSONHeaders(existing, new string) string {
	// Simple implementation - in practice you might want more sophisticated JSON merging
	existing = strings.TrimPrefix(strings.TrimSuffix(existing, "}"), "{")
	new = strings.TrimPrefix(strings.TrimSuffix(new, "}"), "{")

	if existing == "" {
		return fmt.Sprintf("{%s}", new)
	}
	if new == "" {
		return fmt.Sprintf("{%s}", existing)
	}

	return fmt.Sprintf("{%s,%s}", existing, new)
}
