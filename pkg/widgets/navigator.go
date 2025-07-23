package widgets

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gideonsigilai/godin/pkg/core"
)

// Navigator manages page navigation and routing
type Navigator struct {
	pageStack      []*PageInfo
	routeTable     map[string]RouteHandler
	currentIndex   int
	mutex          sync.RWMutex
	context        *core.Context
	onRouteChanged func(route string)
	observers      []NavigatorObserver
	canPopCallback func() bool
}

// PageInfo contains information about a page in the navigation stack
type PageInfo struct {
	ID         string
	Route      string
	Widget     core.Widget
	Parameters map[string]interface{}
	Title      string
	CanPop     bool
	CreatedAt  time.Time
	Arguments  interface{}
	Settings   *RouteSettings
}

// RouteHandler is a function that creates a widget for a route
type RouteHandler func(ctx *core.Context, params map[string]interface{}) core.Widget

// NavigatorObserver observes navigation events
type NavigatorObserver interface {
	DidPush(route string, previousRoute string)
	DidPop(route string, previousRoute string)
	DidRemove(route string, previousRoute string)
	DidReplace(oldRoute string, newRoute string)
}

// RouteSettings contains settings for routes
type RouteSettings struct {
	Name               string
	Arguments          interface{}
	BarrierDismissible bool
	BarrierColor       *core.Color
	BarrierLabel       string
	UseSafeArea        bool
	UseRootNavigator   bool
	Maintainstate      bool
	FullscreenDialog   bool
}

// NavigationResult represents the result of a navigation operation
type NavigationResult struct {
	Success bool
	Error   error
	Data    interface{}
}

// NewNavigator creates a new Navigator instance
func NewNavigator(ctx *core.Context) *Navigator {
	return &Navigator{
		pageStack:    make([]*PageInfo, 0),
		routeTable:   make(map[string]RouteHandler),
		currentIndex: -1,
		context:      ctx,
		observers:    make([]NavigatorObserver, 0),
	}
}

// RegisterRoute registers a route handler
func (n *Navigator) RegisterRoute(route string, handler RouteHandler) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.routeTable[route] = handler
}

// RegisterRoutes registers multiple routes at once
func (n *Navigator) RegisterRoutes(routes map[string]RouteHandler) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	for route, handler := range routes {
		n.routeTable[route] = handler
	}
}

// Push adds a new page to the navigation stack
func (n *Navigator) Push(route string, widget core.Widget, args ...interface{}) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	// Generate unique page ID
	pageID := fmt.Sprintf("page_%d_%d", time.Now().UnixNano(), len(n.pageStack))

	// Extract parameters from route if it contains query parameters
	params, cleanRoute := n.extractRouteParameters(route)

	// Create page info
	pageInfo := &PageInfo{
		ID:         pageID,
		Route:      cleanRoute,
		Widget:     widget,
		Parameters: params,
		Title:      cleanRoute,
		CanPop:     true,
		CreatedAt:  time.Now(),
	}

	// Add arguments if provided
	if len(args) > 0 {
		pageInfo.Arguments = args[0]
	}

	// Get previous route for observers
	var previousRoute string
	if len(n.pageStack) > 0 {
		previousRoute = n.pageStack[len(n.pageStack)-1].Route
	}

	// Add to stack
	n.pageStack = append(n.pageStack, pageInfo)
	n.currentIndex = len(n.pageStack) - 1

	// Update browser URL if context is available
	if n.context != nil {
		n.updateBrowserURL(route)
	}

	// Notify observers
	n.notifyObservers(func(observer NavigatorObserver) {
		observer.DidPush(cleanRoute, previousRoute)
	})

	// Call route changed callback
	if n.onRouteChanged != nil {
		go n.onRouteChanged(cleanRoute)
	}

	return nil
}

// Pop removes the current page from the navigation stack
func (n *Navigator) Pop(result ...interface{}) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if len(n.pageStack) <= 1 {
		return fmt.Errorf("cannot pop the last page")
	}

	// Check if we can pop
	if n.canPopCallback != nil && !n.canPopCallback() {
		return fmt.Errorf("navigation blocked by canPop callback")
	}

	// Get current and previous routes
	currentRoute := n.pageStack[n.currentIndex].Route
	previousRoute := ""
	if n.currentIndex > 0 {
		previousRoute = n.pageStack[n.currentIndex-1].Route
	}

	// Remove current page
	n.pageStack = n.pageStack[:len(n.pageStack)-1]
	n.currentIndex = len(n.pageStack) - 1

	// Update browser URL
	if n.context != nil && len(n.pageStack) > 0 {
		n.updateBrowserURL(n.pageStack[n.currentIndex].Route)
	}

	// Notify observers
	n.notifyObservers(func(observer NavigatorObserver) {
		observer.DidPop(currentRoute, previousRoute)
	})

	// Call route changed callback
	if n.onRouteChanged != nil && len(n.pageStack) > 0 {
		go n.onRouteChanged(n.pageStack[n.currentIndex].Route)
	}

	return nil
}

// PushReplacement replaces the current page with a new one
func (n *Navigator) PushReplacement(route string, widget core.Widget, args ...interface{}) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if len(n.pageStack) == 0 {
		// If stack is empty, just push
		n.mutex.Unlock()
		return n.Push(route, widget, args...)
	}

	// Get current route for observers
	oldRoute := n.pageStack[n.currentIndex].Route

	// Extract parameters from route
	params, cleanRoute := n.extractRouteParameters(route)

	// Generate unique page ID
	pageID := fmt.Sprintf("page_%d_%d", time.Now().UnixNano(), len(n.pageStack))

	// Create new page info
	pageInfo := &PageInfo{
		ID:         pageID,
		Route:      cleanRoute,
		Widget:     widget,
		Parameters: params,
		Title:      cleanRoute,
		CanPop:     len(n.pageStack) > 1, // Can pop if there are other pages
		CreatedAt:  time.Now(),
	}

	// Add arguments if provided
	if len(args) > 0 {
		pageInfo.Arguments = args[0]
	}

	// Replace current page
	n.pageStack[n.currentIndex] = pageInfo

	// Update browser URL
	if n.context != nil {
		n.updateBrowserURL(route)
	}

	// Notify observers
	n.notifyObservers(func(observer NavigatorObserver) {
		observer.DidReplace(oldRoute, cleanRoute)
	})

	// Call route changed callback
	if n.onRouteChanged != nil {
		go n.onRouteChanged(cleanRoute)
	}

	return nil
}

// PushAndRemoveUntil pushes a new page and removes pages until predicate returns true
func (n *Navigator) PushAndRemoveUntil(route string, widget core.Widget, predicate func(*PageInfo) bool, args ...interface{}) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	// First, push the new page
	n.mutex.Unlock()
	err := n.Push(route, widget, args...)
	if err != nil {
		return err
	}
	n.mutex.Lock()

	// Find the index where predicate returns true
	removeIndex := -1
	for i := len(n.pageStack) - 2; i >= 0; i-- { // Start from second-to-last (before the newly pushed page)
		if predicate(n.pageStack[i]) {
			removeIndex = i + 1 // Keep the page where predicate is true
			break
		}
	}

	if removeIndex > 0 && removeIndex < len(n.pageStack)-1 {
		// Remove pages between removeIndex and the newly pushed page
		newStack := make([]*PageInfo, 0, removeIndex+1)
		newStack = append(newStack, n.pageStack[:removeIndex]...)
		newStack = append(newStack, n.pageStack[len(n.pageStack)-1]) // Keep the newly pushed page

		n.pageStack = newStack
		n.currentIndex = len(n.pageStack) - 1
	}

	return nil
}

// PopUntil pops pages until predicate returns true
func (n *Navigator) PopUntil(predicate func(*PageInfo) bool) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if len(n.pageStack) <= 1 {
		return fmt.Errorf("cannot pop the last page")
	}

	// Find the target page
	targetIndex := -1
	for i := len(n.pageStack) - 1; i >= 0; i-- {
		if predicate(n.pageStack[i]) {
			targetIndex = i
			break
		}
	}

	if targetIndex == -1 {
		return fmt.Errorf("no page matches the predicate")
	}

	if targetIndex == n.currentIndex {
		return nil // Already at target page
	}

	// Get routes for observers
	currentRoute := n.pageStack[n.currentIndex].Route
	targetRoute := n.pageStack[targetIndex].Route

	// Remove pages after target
	n.pageStack = n.pageStack[:targetIndex+1]
	n.currentIndex = len(n.pageStack) - 1

	// Update browser URL
	if n.context != nil {
		n.updateBrowserURL(targetRoute)
	}

	// Notify observers
	n.notifyObservers(func(observer NavigatorObserver) {
		observer.DidPop(currentRoute, targetRoute)
	})

	// Call route changed callback
	if n.onRouteChanged != nil {
		go n.onRouteChanged(targetRoute)
	}

	return nil
}

// CanPop returns true if there are pages that can be popped
func (n *Navigator) CanPop() bool {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	if len(n.pageStack) <= 1 {
		return false
	}

	if n.canPopCallback != nil {
		return n.canPopCallback()
	}

	return n.pageStack[n.currentIndex].CanPop
}

// GetCurrentPage returns the current page info
func (n *Navigator) GetCurrentPage() (*PageInfo, bool) {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	if n.currentIndex < 0 || n.currentIndex >= len(n.pageStack) {
		return nil, false
	}

	// Return a copy to prevent external modification
	page := *n.pageStack[n.currentIndex]
	return &page, true
}

// GetPageStack returns a copy of the current page stack
func (n *Navigator) GetPageStack() []*PageInfo {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	stack := make([]*PageInfo, len(n.pageStack))
	for i, page := range n.pageStack {
		pageCopy := *page
		stack[i] = &pageCopy
	}
	return stack
}

// NavigateToRoute navigates to a named route
func (n *Navigator) NavigateToRoute(route string, args ...interface{}) error {
	n.mutex.RLock()
	handler, exists := n.routeTable[route]
	n.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("route '%s' not found", route)
	}

	// Extract parameters from route
	params, cleanRoute := n.extractRouteParameters(route)

	// Create widget using route handler
	widget := handler(n.context, params)
	if widget == nil {
		return fmt.Errorf("route handler for '%s' returned nil widget", route)
	}

	return n.Push(cleanRoute, widget, args...)
}

// SetCanPopCallback sets a callback to determine if navigation can be popped
func (n *Navigator) SetCanPopCallback(callback func() bool) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.canPopCallback = callback
}

// SetOnRouteChanged sets a callback for route changes
func (n *Navigator) SetOnRouteChanged(callback func(route string)) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.onRouteChanged = callback
}

// AddObserver adds a navigation observer
func (n *Navigator) AddObserver(observer NavigatorObserver) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.observers = append(n.observers, observer)
}

// RemoveObserver removes a navigation observer
func (n *Navigator) RemoveObserver(observer NavigatorObserver) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	for i, obs := range n.observers {
		if obs == observer {
			n.observers = append(n.observers[:i], n.observers[i+1:]...)
			break
		}
	}
}

// extractRouteParameters extracts query parameters from a route
func (n *Navigator) extractRouteParameters(route string) (map[string]interface{}, string) {
	params := make(map[string]interface{})

	// Parse URL to extract query parameters
	if strings.Contains(route, "?") {
		parts := strings.SplitN(route, "?", 2)
		cleanRoute := parts[0]

		if len(parts) > 1 {
			values, err := url.ParseQuery(parts[1])
			if err == nil {
				for key, vals := range values {
					if len(vals) > 0 {
						params[key] = vals[0] // Take first value
					}
				}
			}
		}

		return params, cleanRoute
	}

	return params, route
}

// updateBrowserURL updates the browser URL (placeholder for HTMX integration)
func (n *Navigator) updateBrowserURL(route string) {
	// This would integrate with HTMX to update browser URL
	// For now, we'll store it in context for potential use
	if n.context != nil {
		n.context.Set("current_route", route)
	}
}

// notifyObservers notifies all observers with the given function
func (n *Navigator) notifyObservers(notify func(NavigatorObserver)) {
	for _, observer := range n.observers {
		go notify(observer) // Run in goroutine to avoid blocking
	}
}

// Clear clears the navigation stack
func (n *Navigator) Clear() {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	n.pageStack = make([]*PageInfo, 0)
	n.currentIndex = -1
}

// GetRouteTable returns a copy of the route table
func (n *Navigator) GetRouteTable() map[string]RouteHandler {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	table := make(map[string]RouteHandler)
	for route, handler := range n.routeTable {
		table[route] = handler
	}
	return table
}

// HasRoute checks if a route is registered
func (n *Navigator) HasRoute(route string) bool {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	_, exists := n.routeTable[route]
	return exists
}
