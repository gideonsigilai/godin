package core

import (
	"sync"
)

// MediaQueryData contains screen and device information
type MediaQueryData struct {
	Size                  Size
	DevicePixelRatio      float64
	TextScaleFactor       float64
	Padding               EdgeInsets
	ViewInsets            EdgeInsets
	ViewPadding           EdgeInsets
	AlwaysUse24HourFormat bool
	AccessibleNavigation  bool
	InvertColors          bool
	HighContrast          bool
	DisableAnimations     bool
	BoldText              bool
	Orientation           Orientation
	PlatformBrightness    Brightness
	Breakpoint            Breakpoint
	SystemGestureInsets   EdgeInsets
	DisplayFeatures       []DisplayFeature
}

// DisplayFeature represents a display feature like notches or folds
type DisplayFeature struct {
	Bounds Rect
	Type   DisplayFeatureType
	State  DisplayFeatureState
}

// DisplayFeatureType represents the type of display feature
type DisplayFeatureType string

const (
	DisplayFeatureTypeFold   DisplayFeatureType = "fold"
	DisplayFeatureTypeHinge  DisplayFeatureType = "hinge"
	DisplayFeatureTypeCutout DisplayFeatureType = "cutout"
)

// DisplayFeatureState represents the state of a display feature
type DisplayFeatureState string

const (
	DisplayFeatureStateUnknown           DisplayFeatureState = "unknown"
	DisplayFeatureStatePostureFlat       DisplayFeatureState = "flat"
	DisplayFeatureStatePostureHalfOpened DisplayFeatureState = "halfOpened"
)

// Rect represents a rectangle
type Rect struct {
	Left   float64
	Top    float64
	Right  float64
	Bottom float64
}

// MediaQueryProvider manages screen information and updates
type MediaQueryProvider struct {
	currentData   *MediaQueryData
	listeners     []func(*MediaQueryData)
	mutex         sync.RWMutex
	updateChannel chan *MediaQueryData
	isListening   bool
}

// NewMediaQueryProvider creates a new MediaQueryProvider
func NewMediaQueryProvider() *MediaQueryProvider {
	return &MediaQueryProvider{
		currentData:   NewDefaultMediaQueryData(),
		listeners:     make([]func(*MediaQueryData), 0),
		updateChannel: make(chan *MediaQueryData, 10),
		isListening:   false,
	}
}

// NewDefaultMediaQueryData creates default MediaQueryData
func NewDefaultMediaQueryData() *MediaQueryData {
	return &MediaQueryData{
		Size:                  NewSize(1920, 1080), // Default desktop size
		DevicePixelRatio:      1.0,
		TextScaleFactor:       1.0,
		Padding:               EdgeInsets{},
		ViewInsets:            EdgeInsets{},
		ViewPadding:           EdgeInsets{},
		AlwaysUse24HourFormat: false,
		AccessibleNavigation:  false,
		InvertColors:          false,
		HighContrast:          false,
		DisableAnimations:     false,
		BoldText:              false,
		Orientation:           OrientationLandscape,
		PlatformBrightness:    BrightnessLight,
		Breakpoint:            GetBreakpoint(1920),
		SystemGestureInsets:   EdgeInsets{},
		DisplayFeatures:       make([]DisplayFeature, 0),
	}
}

// GetData returns the current MediaQueryData
func (mqp *MediaQueryProvider) GetData() *MediaQueryData {
	mqp.mutex.RLock()
	defer mqp.mutex.RUnlock()

	// Return a copy to prevent external modification
	return mqp.copyMediaQueryData(mqp.currentData)
}

// UpdateData updates the MediaQueryData and notifies listeners
func (mqp *MediaQueryProvider) UpdateData(data *MediaQueryData) {
	mqp.mutex.Lock()
	defer mqp.mutex.Unlock()

	// Update breakpoint based on size
	data.Breakpoint = GetBreakpoint(data.Size.Width)

	// Update orientation based on size
	if data.Size.Width > data.Size.Height {
		data.Orientation = OrientationLandscape
	} else {
		data.Orientation = OrientationPortrait
	}

	mqp.currentData = data
	mqp.notifyListeners()
}

// UpdateSize updates just the screen size
func (mqp *MediaQueryProvider) UpdateSize(width, height float64) {
	mqp.mutex.Lock()
	defer mqp.mutex.Unlock()

	mqp.currentData.Size = NewSize(width, height)
	mqp.currentData.Breakpoint = GetBreakpoint(width)

	if width > height {
		mqp.currentData.Orientation = OrientationLandscape
	} else {
		mqp.currentData.Orientation = OrientationPortrait
	}

	mqp.notifyListeners()
}

// UpdateDevicePixelRatio updates the device pixel ratio
func (mqp *MediaQueryProvider) UpdateDevicePixelRatio(ratio float64) {
	mqp.mutex.Lock()
	defer mqp.mutex.Unlock()

	mqp.currentData.DevicePixelRatio = ratio
	mqp.notifyListeners()
}

// UpdateTextScaleFactor updates the text scale factor
func (mqp *MediaQueryProvider) UpdateTextScaleFactor(factor float64) {
	mqp.mutex.Lock()
	defer mqp.mutex.Unlock()

	mqp.currentData.TextScaleFactor = factor
	mqp.notifyListeners()
}

// UpdatePlatformBrightness updates the platform brightness
func (mqp *MediaQueryProvider) UpdatePlatformBrightness(brightness Brightness) {
	mqp.mutex.Lock()
	defer mqp.mutex.Unlock()

	mqp.currentData.PlatformBrightness = brightness
	mqp.notifyListeners()
}

// UpdateAccessibilityFeatures updates accessibility features
func (mqp *MediaQueryProvider) UpdateAccessibilityFeatures(features AccessibilityFeatures) {
	mqp.mutex.Lock()
	defer mqp.mutex.Unlock()

	mqp.currentData.AccessibleNavigation = features.AccessibleNavigation
	mqp.currentData.InvertColors = features.InvertColors
	mqp.currentData.HighContrast = features.HighContrast
	mqp.currentData.DisableAnimations = features.DisableAnimations
	mqp.currentData.BoldText = features.BoldText

	mqp.notifyListeners()
}

// AccessibilityFeatures represents accessibility settings
type AccessibilityFeatures struct {
	AccessibleNavigation bool
	InvertColors         bool
	HighContrast         bool
	DisableAnimations    bool
	BoldText             bool
}

// AddListener adds a MediaQuery change listener
func (mqp *MediaQueryProvider) AddListener(listener func(*MediaQueryData)) {
	mqp.mutex.Lock()
	defer mqp.mutex.Unlock()

	mqp.listeners = append(mqp.listeners, listener)
}

// RemoveListener removes a MediaQuery change listener
func (mqp *MediaQueryProvider) RemoveListener(listener func(*MediaQueryData)) {
	mqp.mutex.Lock()
	defer mqp.mutex.Unlock()

	// Note: This is a simplified implementation
	// In practice, you'd need a way to identify listeners uniquely
	for i, l := range mqp.listeners {
		if &l == &listener {
			mqp.listeners = append(mqp.listeners[:i], mqp.listeners[i+1:]...)
			break
		}
	}
}

// StartListening starts listening for MediaQuery updates
func (mqp *MediaQueryProvider) StartListening() {
	mqp.mutex.Lock()
	defer mqp.mutex.Unlock()

	if mqp.isListening {
		return
	}

	mqp.isListening = true

	// Start background goroutine to process updates
	go mqp.processUpdates()
}

// StopListening stops listening for MediaQuery updates
func (mqp *MediaQueryProvider) StopListening() {
	mqp.mutex.Lock()
	defer mqp.mutex.Unlock()

	mqp.isListening = false
	close(mqp.updateChannel)
}

// processUpdates processes MediaQuery updates in the background
func (mqp *MediaQueryProvider) processUpdates() {
	for data := range mqp.updateChannel {
		mqp.UpdateData(data)
	}
}

// QueueUpdate queues a MediaQuery update
func (mqp *MediaQueryProvider) QueueUpdate(data *MediaQueryData) {
	select {
	case mqp.updateChannel <- data:
		// Update queued successfully
	default:
		// Channel is full, skip this update
	}
}

// notifyListeners notifies all listeners of MediaQuery changes
func (mqp *MediaQueryProvider) notifyListeners() {
	dataCopy := mqp.copyMediaQueryData(mqp.currentData)
	for _, listener := range mqp.listeners {
		go listener(dataCopy) // Run in goroutine and pass copy
	}
}

// copyMediaQueryData creates a deep copy of MediaQueryData
func (mqp *MediaQueryProvider) copyMediaQueryData(data *MediaQueryData) *MediaQueryData {
	if data == nil {
		return nil
	}

	copy := &MediaQueryData{
		Size:                  data.Size,
		DevicePixelRatio:      data.DevicePixelRatio,
		TextScaleFactor:       data.TextScaleFactor,
		Padding:               data.Padding,
		ViewInsets:            data.ViewInsets,
		ViewPadding:           data.ViewPadding,
		AlwaysUse24HourFormat: data.AlwaysUse24HourFormat,
		AccessibleNavigation:  data.AccessibleNavigation,
		InvertColors:          data.InvertColors,
		HighContrast:          data.HighContrast,
		DisableAnimations:     data.DisableAnimations,
		BoldText:              data.BoldText,
		Orientation:           data.Orientation,
		PlatformBrightness:    data.PlatformBrightness,
		Breakpoint:            data.Breakpoint,
		SystemGestureInsets:   data.SystemGestureInsets,
	}

	// Copy display features
	copy.DisplayFeatures = make([]DisplayFeature, len(data.DisplayFeatures))
	for i, feature := range data.DisplayFeatures {
		copy.DisplayFeatures[i] = feature
	}

	return copy
}

// GetBreakpointName returns the name of the current breakpoint
func (mqd *MediaQueryData) GetBreakpointName() string {
	return string(mqd.Breakpoint)
}

// IsPortrait returns true if the orientation is portrait
func (mqd *MediaQueryData) IsPortrait() bool {
	return mqd.Orientation == OrientationPortrait
}

// IsLandscape returns true if the orientation is landscape
func (mqd *MediaQueryData) IsLandscape() bool {
	return mqd.Orientation == OrientationLandscape
}

// IsSmallScreen returns true if the screen is considered small
func (mqd *MediaQueryData) IsSmallScreen() bool {
	return mqd.Breakpoint == BreakpointXS || mqd.Breakpoint == BreakpointSM
}

// IsMediumScreen returns true if the screen is considered medium
func (mqd *MediaQueryData) IsMediumScreen() bool {
	return mqd.Breakpoint == BreakpointMD
}

// IsLargeScreen returns true if the screen is considered large
func (mqd *MediaQueryData) IsLargeScreen() bool {
	return mqd.Breakpoint == BreakpointLG || mqd.Breakpoint == BreakpointXL
}

// GetSafeAreaInsets returns the safe area insets
func (mqd *MediaQueryData) GetSafeAreaInsets() EdgeInsets {
	return EdgeInsets{
		Top:    ClampFloat64(mqd.Padding.Top, 0, 100),
		Right:  ClampFloat64(mqd.Padding.Right, 0, 100),
		Bottom: ClampFloat64(mqd.Padding.Bottom, 0, 100),
		Left:   ClampFloat64(mqd.Padding.Left, 0, 100),
	}
}

// GetViewportSize returns the viewport size minus insets
func (mqd *MediaQueryData) GetViewportSize() Size {
	return Size{
		Width:  mqd.Size.Width - mqd.ViewInsets.Left - mqd.ViewInsets.Right,
		Height: mqd.Size.Height - mqd.ViewInsets.Top - mqd.ViewInsets.Bottom,
	}
}

// HasDisplayFeatures returns true if there are display features
func (mqd *MediaQueryData) HasDisplayFeatures() bool {
	return len(mqd.DisplayFeatures) > 0
}

// GetDisplayFeaturesByType returns display features of a specific type
func (mqd *MediaQueryData) GetDisplayFeaturesByType(featureType DisplayFeatureType) []DisplayFeature {
	var features []DisplayFeature
	for _, feature := range mqd.DisplayFeatures {
		if feature.Type == featureType {
			features = append(features, feature)
		}
	}
	return features
}

// MediaQueryInheritedWidget provides MediaQuery data to child widgets
type MediaQueryInheritedWidget struct {
	Data  *MediaQueryData
	Child Widget
}

// Render renders the MediaQueryInheritedWidget
func (mq MediaQueryInheritedWidget) Render(ctx *Context) string {
	// Store MediaQuery data in context for child widgets to access
	ctx.Set("mediaQuery", mq.Data)

	if mq.Child != nil {
		return mq.Child.Render(ctx)
	}

	return ""
}

// MediaQueryOf returns the MediaQueryData from the context
func MediaQueryOf(ctx *Context) *MediaQueryData {
	if data, ok := ctx.Get("mediaQuery").(*MediaQueryData); ok {
		return data
	}

	// Return default data if not found in context
	return NewDefaultMediaQueryData()
}

// ResponsiveBuilder builds different widgets based on screen size
type ResponsiveBuilder struct {
	XS    Widget // Extra small screens
	SM    Widget // Small screens
	MD    Widget // Medium screens
	LG    Widget // Large screens
	XL    Widget // Extra large screens
	Child Widget // Fallback widget
}

// Render renders the appropriate widget based on screen size
func (rb ResponsiveBuilder) Render(ctx *Context) string {
	mediaQuery := MediaQueryOf(ctx)

	var widget Widget

	switch mediaQuery.Breakpoint {
	case BreakpointXS:
		widget = rb.XS
	case BreakpointSM:
		widget = rb.SM
	case BreakpointMD:
		widget = rb.MD
	case BreakpointLG:
		widget = rb.LG
	case BreakpointXL:
		widget = rb.XL
	}

	// Fallback to child if no specific widget is defined
	if widget == nil {
		widget = rb.Child
	}

	if widget != nil {
		return widget.Render(ctx)
	}

	return ""
}

// OrientationBuilder builds different widgets based on orientation
type OrientationBuilder struct {
	Portrait  Widget
	Landscape Widget
	Child     Widget // Fallback widget
}

// Render renders the appropriate widget based on orientation
func (ob OrientationBuilder) Render(ctx *Context) string {
	mediaQuery := MediaQueryOf(ctx)

	var widget Widget

	if mediaQuery.IsPortrait() && ob.Portrait != nil {
		widget = ob.Portrait
	} else if mediaQuery.IsLandscape() && ob.Landscape != nil {
		widget = ob.Landscape
	} else {
		widget = ob.Child
	}

	if widget != nil {
		return widget.Render(ctx)
	}

	return ""
}

// LayoutBuilder provides constraints-based layout building
type LayoutBuilder struct {
	Builder func(ctx *Context, constraints BoxConstraints) Widget
}

// BoxConstraints represents layout constraints
type BoxConstraints struct {
	MinWidth  float64
	MaxWidth  float64
	MinHeight float64
	MaxHeight float64
}

// Render renders the LayoutBuilder
func (lb LayoutBuilder) Render(ctx *Context) string {
	if lb.Builder == nil {
		return ""
	}

	mediaQuery := MediaQueryOf(ctx)

	// Create constraints based on screen size
	constraints := BoxConstraints{
		MinWidth:  0,
		MaxWidth:  mediaQuery.Size.Width,
		MinHeight: 0,
		MaxHeight: mediaQuery.Size.Height,
	}

	widget := lb.Builder(ctx, constraints)
	if widget != nil {
		return widget.Render(ctx)
	}

	return ""
}
