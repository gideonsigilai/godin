package core

import (
	"fmt"
	"strings"
	"sync"
)

// ThemeData contains all theme configuration
type ThemeData struct {
	ColorScheme     *ColorScheme
	Typography      *Typography
	ComponentThemes map[string]interface{}
	Extensions      map[string]interface{}
	Brightness      Brightness
	CSS             map[string]string // CSS custom properties
	UseMaterial3    bool
	VisualDensity   VisualDensity
}

// ColorScheme defines the color palette for the theme
type ColorScheme struct {
	// Primary colors
	Primary            Color
	OnPrimary          Color
	PrimaryContainer   Color
	OnPrimaryContainer Color

	// Secondary colors
	Secondary            Color
	OnSecondary          Color
	SecondaryContainer   Color
	OnSecondaryContainer Color

	// Tertiary colors
	Tertiary            Color
	OnTertiary          Color
	TertiaryContainer   Color
	OnTertiaryContainer Color

	// Error colors
	Error            Color
	OnError          Color
	ErrorContainer   Color
	OnErrorContainer Color

	// Surface colors
	Surface          Color
	OnSurface        Color
	SurfaceVariant   Color
	OnSurfaceVariant Color
	SurfaceTint      Color

	// Background colors
	Background   Color
	OnBackground Color

	// Outline colors
	Outline        Color
	OutlineVariant Color

	// Other colors
	Shadow           Color
	Scrim            Color
	InverseSurface   Color
	InverseOnSurface Color
	InversePrimary   Color

	// Brightness
	Brightness Brightness
}

// Typography defines text styles for the theme
type Typography struct {
	// Display styles
	DisplayLarge  *TextStyle
	DisplayMedium *TextStyle
	DisplaySmall  *TextStyle

	// Headline styles
	HeadlineLarge  *TextStyle
	HeadlineMedium *TextStyle
	HeadlineSmall  *TextStyle

	// Title styles
	TitleLarge  *TextStyle
	TitleMedium *TextStyle
	TitleSmall  *TextStyle

	// Body styles
	BodyLarge  *TextStyle
	BodyMedium *TextStyle
	BodySmall  *TextStyle

	// Label styles
	LabelLarge  *TextStyle
	LabelMedium *TextStyle
	LabelSmall  *TextStyle
}

// VisualDensity represents the visual density of UI elements
type VisualDensity struct {
	Horizontal float64
	Vertical   float64
}

// ThemeProvider manages theme state and updates
type ThemeProvider struct {
	currentTheme *ThemeData
	lightTheme   *ThemeData
	darkTheme    *ThemeData
	themeMode    ThemeMode
	listeners    []func(*ThemeData)
	mutex        sync.RWMutex
	cssGenerator *CSSGenerator
}

// CSSGenerator generates CSS from theme data
type CSSGenerator struct {
	prefix string
}

// Standard visual densities
var (
	VisualDensityStandard    = VisualDensity{Horizontal: 0, Vertical: 0}
	VisualDensityComfortable = VisualDensity{Horizontal: -1, Vertical: -1}
	VisualDensityCompact     = VisualDensity{Horizontal: -2, Vertical: -2}
	VisualDensityAdaptive    = VisualDensity{Horizontal: 0, Vertical: 0}
)

// NewThemeData creates a new ThemeData with default values
func NewThemeData() *ThemeData {
	return &ThemeData{
		ColorScheme:     NewLightColorScheme(),
		Typography:      NewDefaultTypography(),
		ComponentThemes: make(map[string]interface{}),
		Extensions:      make(map[string]interface{}),
		Brightness:      BrightnessLight,
		CSS:             make(map[string]string),
		UseMaterial3:    true,
		VisualDensity:   VisualDensityStandard,
	}
}

// NewLightColorScheme creates a light color scheme
func NewLightColorScheme() *ColorScheme {
	return &ColorScheme{
		Primary:            ColorPrimary,
		OnPrimary:          ColorOnPrimary,
		PrimaryContainer:   ColorPrimary.Lighten(0.8),
		OnPrimaryContainer: ColorPrimary.Darken(0.2),

		Secondary:            ColorSecondary,
		OnSecondary:          ColorOnSecondary,
		SecondaryContainer:   ColorSecondary.Lighten(0.8),
		OnSecondaryContainer: ColorSecondary.Darken(0.2),

		Tertiary:            NewColor(125, 82, 96, 255),
		OnTertiary:          ColorWhite,
		TertiaryContainer:   NewColor(255, 216, 228, 255),
		OnTertiaryContainer: NewColor(55, 11, 30, 255),

		Error:            ColorError,
		OnError:          ColorWhite,
		ErrorContainer:   ColorError.Lighten(0.8),
		OnErrorContainer: ColorError.Darken(0.2),

		Surface:          ColorSurface,
		OnSurface:        ColorOnSurface,
		SurfaceVariant:   ColorGreyLight,
		OnSurfaceVariant: ColorGreyDark,
		SurfaceTint:      ColorPrimary,

		Background:   ColorBackground,
		OnBackground: ColorOnBackground,

		Outline:        ColorGrey,
		OutlineVariant: ColorGreyLight,

		Shadow:           ColorBlack,
		Scrim:            ColorBlack.WithOpacity(0.5),
		InverseSurface:   ColorGreyDark,
		InverseOnSurface: ColorWhite,
		InversePrimary:   ColorPrimaryLight,

		Brightness: BrightnessLight,
	}
}

// NewDarkColorScheme creates a dark color scheme
func NewDarkColorScheme() *ColorScheme {
	return &ColorScheme{
		Primary:            ColorPrimaryLight,
		OnPrimary:          ColorOnPrimaryDark,
		PrimaryContainer:   ColorPrimaryDark,
		OnPrimaryContainer: ColorPrimaryLight,

		Secondary:            ColorSecondaryLight,
		OnSecondary:          ColorOnSecondaryDark,
		SecondaryContainer:   ColorSecondaryDark,
		OnSecondaryContainer: ColorSecondaryLight,

		Tertiary:            NewColor(204, 187, 199, 255),
		OnTertiary:          NewColor(76, 32, 51, 255),
		TertiaryContainer:   NewColor(100, 58, 72, 255),
		OnTertiaryContainer: NewColor(255, 216, 228, 255),

		Error:            NewColor(244, 67, 54, 255),
		OnError:          ColorBlack,
		ErrorContainer:   NewColor(147, 0, 10, 255),
		OnErrorContainer: NewColor(255, 180, 171, 255),

		Surface:          ColorSurfaceDark,
		OnSurface:        ColorOnSurfaceDark,
		SurfaceVariant:   NewColor(68, 71, 78, 255),
		OnSurfaceVariant: NewColor(196, 199, 221, 255),
		SurfaceTint:      ColorPrimaryLight,

		Background:   ColorBackgroundDark,
		OnBackground: ColorOnBackgroundDark,

		Outline:        NewColor(142, 145, 153, 255),
		OutlineVariant: NewColor(68, 71, 78, 255),

		Shadow:           ColorBlack,
		Scrim:            ColorBlack.WithOpacity(0.5),
		InverseSurface:   ColorGreyLight,
		InverseOnSurface: ColorBlack,
		InversePrimary:   ColorPrimary,

		Brightness: BrightnessDark,
	}
}

// NewDefaultTypography creates default typography
func NewDefaultTypography() *Typography {
	return &Typography{
		DisplayLarge:  NewTextStyle().WithFontSize(57).WithFontWeight(FontWeightNormal),
		DisplayMedium: NewTextStyle().WithFontSize(45).WithFontWeight(FontWeightNormal),
		DisplaySmall:  NewTextStyle().WithFontSize(36).WithFontWeight(FontWeightNormal),

		HeadlineLarge:  NewTextStyle().WithFontSize(32).WithFontWeight(FontWeightNormal),
		HeadlineMedium: NewTextStyle().WithFontSize(28).WithFontWeight(FontWeightNormal),
		HeadlineSmall:  NewTextStyle().WithFontSize(24).WithFontWeight(FontWeightNormal),

		TitleLarge:  NewTextStyle().WithFontSize(22).WithFontWeight(FontWeightNormal),
		TitleMedium: NewTextStyle().WithFontSize(16).WithFontWeight(FontWeightMedium),
		TitleSmall:  NewTextStyle().WithFontSize(14).WithFontWeight(FontWeightMedium),

		BodyLarge:  NewTextStyle().WithFontSize(16).WithFontWeight(FontWeightNormal),
		BodyMedium: NewTextStyle().WithFontSize(14).WithFontWeight(FontWeightNormal),
		BodySmall:  NewTextStyle().WithFontSize(12).WithFontWeight(FontWeightNormal),

		LabelLarge:  NewTextStyle().WithFontSize(14).WithFontWeight(FontWeightMedium),
		LabelMedium: NewTextStyle().WithFontSize(12).WithFontWeight(FontWeightMedium),
		LabelSmall:  NewTextStyle().WithFontSize(11).WithFontWeight(FontWeightMedium),
	}
}

// NewThemeProvider creates a new ThemeProvider
func NewThemeProvider() *ThemeProvider {
	lightTheme := NewThemeData()
	darkTheme := NewThemeData()
	darkTheme.ColorScheme = NewDarkColorScheme()
	darkTheme.Brightness = BrightnessDark

	return &ThemeProvider{
		currentTheme: lightTheme,
		lightTheme:   lightTheme,
		darkTheme:    darkTheme,
		themeMode:    ThemeModeLight,
		listeners:    make([]func(*ThemeData), 0),
		cssGenerator: NewCSSGenerator("godin"),
	}
}

// SetTheme sets the current theme
func (tp *ThemeProvider) SetTheme(theme *ThemeData) {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()

	tp.currentTheme = theme
	tp.notifyListeners()
}

// SetThemeMode sets the theme mode
func (tp *ThemeProvider) SetThemeMode(mode ThemeMode) {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()

	tp.themeMode = mode

	switch mode {
	case ThemeModeLight:
		tp.currentTheme = tp.lightTheme
	case ThemeModeDark:
		tp.currentTheme = tp.darkTheme
	case ThemeModeSystem:
		// For now, default to light theme
		// In a real implementation, this would check system preference
		tp.currentTheme = tp.lightTheme
	}

	tp.notifyListeners()
}

// GetTheme returns the current theme
func (tp *ThemeProvider) GetTheme() *ThemeData {
	tp.mutex.RLock()
	defer tp.mutex.RUnlock()

	// Return a copy to prevent external modification
	return tp.copyTheme(tp.currentTheme)
}

// GetLightTheme returns the light theme
func (tp *ThemeProvider) GetLightTheme() *ThemeData {
	tp.mutex.RLock()
	defer tp.mutex.RUnlock()

	return tp.copyTheme(tp.lightTheme)
}

// GetDarkTheme returns the dark theme
func (tp *ThemeProvider) GetDarkTheme() *ThemeData {
	tp.mutex.RLock()
	defer tp.mutex.RUnlock()

	return tp.copyTheme(tp.darkTheme)
}

// SetLightTheme sets the light theme
func (tp *ThemeProvider) SetLightTheme(theme *ThemeData) {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()

	tp.lightTheme = theme
	if tp.themeMode == ThemeModeLight || (tp.themeMode == ThemeModeSystem && tp.currentTheme.Brightness == BrightnessLight) {
		tp.currentTheme = theme
		tp.notifyListeners()
	}
}

// SetDarkTheme sets the dark theme
func (tp *ThemeProvider) SetDarkTheme(theme *ThemeData) {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()

	tp.darkTheme = theme
	if tp.themeMode == ThemeModeDark || (tp.themeMode == ThemeModeSystem && tp.currentTheme.Brightness == BrightnessDark) {
		tp.currentTheme = theme
		tp.notifyListeners()
	}
}

// AddListener adds a theme change listener
func (tp *ThemeProvider) AddListener(listener func(*ThemeData)) {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()

	tp.listeners = append(tp.listeners, listener)
}

// RemoveListener removes a theme change listener
func (tp *ThemeProvider) RemoveListener(listener func(*ThemeData)) {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()

	// Note: This is a simplified implementation
	// In practice, you'd need a way to identify listeners uniquely
	for i, l := range tp.listeners {
		if &l == &listener {
			tp.listeners = append(tp.listeners[:i], tp.listeners[i+1:]...)
			break
		}
	}
}

// GenerateCSS generates CSS custom properties from the current theme
func (tp *ThemeProvider) GenerateCSS() string {
	tp.mutex.RLock()
	defer tp.mutex.RUnlock()

	return tp.cssGenerator.GenerateCSS(tp.currentTheme)
}

// notifyListeners notifies all listeners of theme changes
func (tp *ThemeProvider) notifyListeners() {
	for _, listener := range tp.listeners {
		go listener(tp.copyTheme(tp.currentTheme)) // Run in goroutine and pass copy
	}
}

// copyTheme creates a deep copy of a theme
func (tp *ThemeProvider) copyTheme(theme *ThemeData) *ThemeData {
	if theme == nil {
		return nil
	}

	copy := &ThemeData{
		Brightness:    theme.Brightness,
		UseMaterial3:  theme.UseMaterial3,
		VisualDensity: theme.VisualDensity,
	}

	// Copy ColorScheme
	if theme.ColorScheme != nil {
		colorScheme := *theme.ColorScheme
		copy.ColorScheme = &colorScheme
	}

	// Copy Typography
	if theme.Typography != nil {
		typography := *theme.Typography
		copy.Typography = &typography
	}

	// Copy maps
	copy.ComponentThemes = make(map[string]interface{})
	for k, v := range theme.ComponentThemes {
		copy.ComponentThemes[k] = v
	}

	copy.Extensions = make(map[string]interface{})
	for k, v := range theme.Extensions {
		copy.Extensions[k] = v
	}

	copy.CSS = make(map[string]string)
	for k, v := range theme.CSS {
		copy.CSS[k] = v
	}

	return copy
}

// NewCSSGenerator creates a new CSS generator
func NewCSSGenerator(prefix string) *CSSGenerator {
	return &CSSGenerator{
		prefix: prefix,
	}
}

// GenerateCSS generates CSS custom properties from theme data
func (cg *CSSGenerator) GenerateCSS(theme *ThemeData) string {
	if theme == nil {
		return ""
	}

	var css strings.Builder
	css.WriteString(":root {\n")

	// Generate color scheme CSS variables
	if theme.ColorScheme != nil {
		cg.writeColorSchemeCSS(&css, theme.ColorScheme)
	}

	// Generate typography CSS variables
	if theme.Typography != nil {
		cg.writeTypographyCSS(&css, theme.Typography)
	}

	// Add custom CSS properties
	for key, value := range theme.CSS {
		css.WriteString(fmt.Sprintf("  --%s-%s: %s;\n", cg.prefix, key, value))
	}

	css.WriteString("}\n")

	// Add component-specific CSS
	cg.writeComponentCSS(&css, theme)

	return css.String()
}

// writeColorSchemeCSS writes color scheme CSS variables
func (cg *CSSGenerator) writeColorSchemeCSS(css *strings.Builder, colorScheme *ColorScheme) {
	colors := map[string]Color{
		"primary":                colorScheme.Primary,
		"on-primary":             colorScheme.OnPrimary,
		"primary-container":      colorScheme.PrimaryContainer,
		"on-primary-container":   colorScheme.OnPrimaryContainer,
		"secondary":              colorScheme.Secondary,
		"on-secondary":           colorScheme.OnSecondary,
		"secondary-container":    colorScheme.SecondaryContainer,
		"on-secondary-container": colorScheme.OnSecondaryContainer,
		"tertiary":               colorScheme.Tertiary,
		"on-tertiary":            colorScheme.OnTertiary,
		"tertiary-container":     colorScheme.TertiaryContainer,
		"on-tertiary-container":  colorScheme.OnTertiaryContainer,
		"error":                  colorScheme.Error,
		"on-error":               colorScheme.OnError,
		"error-container":        colorScheme.ErrorContainer,
		"on-error-container":     colorScheme.OnErrorContainer,
		"surface":                colorScheme.Surface,
		"on-surface":             colorScheme.OnSurface,
		"surface-variant":        colorScheme.SurfaceVariant,
		"on-surface-variant":     colorScheme.OnSurfaceVariant,
		"surface-tint":           colorScheme.SurfaceTint,
		"background":             colorScheme.Background,
		"on-background":          colorScheme.OnBackground,
		"outline":                colorScheme.Outline,
		"outline-variant":        colorScheme.OutlineVariant,
		"shadow":                 colorScheme.Shadow,
		"scrim":                  colorScheme.Scrim,
		"inverse-surface":        colorScheme.InverseSurface,
		"inverse-on-surface":     colorScheme.InverseOnSurface,
		"inverse-primary":        colorScheme.InversePrimary,
	}

	for name, color := range colors {
		css.WriteString(fmt.Sprintf("  --%s-color-%s: %s;\n", cg.prefix, name, color.ToCSS()))
	}
}

// writeTypographyCSS writes typography CSS variables
func (cg *CSSGenerator) writeTypographyCSS(css *strings.Builder, typography *Typography) {
	styles := map[string]*TextStyle{
		"display-large":   typography.DisplayLarge,
		"display-medium":  typography.DisplayMedium,
		"display-small":   typography.DisplaySmall,
		"headline-large":  typography.HeadlineLarge,
		"headline-medium": typography.HeadlineMedium,
		"headline-small":  typography.HeadlineSmall,
		"title-large":     typography.TitleLarge,
		"title-medium":    typography.TitleMedium,
		"title-small":     typography.TitleSmall,
		"body-large":      typography.BodyLarge,
		"body-medium":     typography.BodyMedium,
		"body-small":      typography.BodySmall,
		"label-large":     typography.LabelLarge,
		"label-medium":    typography.LabelMedium,
		"label-small":     typography.LabelSmall,
	}

	for name, style := range styles {
		if style != nil {
			if style.FontSize != nil {
				css.WriteString(fmt.Sprintf("  --%s-typography-%s-size: %.1fpx;\n", cg.prefix, name, *style.FontSize))
			}
			if style.FontWeight != nil {
				css.WriteString(fmt.Sprintf("  --%s-typography-%s-weight: %d;\n", cg.prefix, name, *style.FontWeight))
			}
			if style.FontFamily != nil {
				css.WriteString(fmt.Sprintf("  --%s-typography-%s-family: %s;\n", cg.prefix, name, *style.FontFamily))
			}
			if style.LineHeight != nil {
				css.WriteString(fmt.Sprintf("  --%s-typography-%s-line-height: %.2f;\n", cg.prefix, name, *style.LineHeight))
			}
		}
	}
}

// writeComponentCSS writes component-specific CSS
func (cg *CSSGenerator) writeComponentCSS(css *strings.Builder, theme *ThemeData) {
	// Add common component styles
	css.WriteString(fmt.Sprintf(`
/* %s Component Styles */
.%s-button {
  background-color: var(--%s-color-primary);
  color: var(--%s-color-on-primary);
  border: none;
  border-radius: 4px;
  padding: 8px 16px;
  font-size: var(--%s-typography-label-large-size);
  font-weight: var(--%s-typography-label-large-weight);
  cursor: pointer;
  transition: all 0.2s ease;
}

.%s-button:hover {
  background-color: var(--%s-color-primary);
  opacity: 0.8;
}

.%s-card {
  background-color: var(--%s-color-surface);
  color: var(--%s-color-on-surface);
  border-radius: 12px;
  box-shadow: 0px 1px 3px rgba(0,0,0,0.12), 0px 1px 2px rgba(0,0,0,0.24);
  padding: 16px;
}

.%s-text-field {
  background-color: var(--%s-color-surface-variant);
  color: var(--%s-color-on-surface);
  border: 1px solid var(--%s-color-outline);
  border-radius: 4px;
  padding: 12px 16px;
  font-size: var(--%s-typography-body-large-size);
}

.%s-text-field:focus {
  border-color: var(--%s-color-primary);
  outline: none;
}
`,
		cg.prefix, cg.prefix, cg.prefix, cg.prefix, cg.prefix, cg.prefix,
		cg.prefix, cg.prefix,
		cg.prefix, cg.prefix, cg.prefix,
		cg.prefix, cg.prefix, cg.prefix, cg.prefix, cg.prefix,
		cg.prefix, cg.prefix,
	))
}

// Predefined themes
var (
	DefaultLightTheme = NewThemeData()
	DefaultDarkTheme  = func() *ThemeData {
		theme := NewThemeData()
		theme.ColorScheme = NewDarkColorScheme()
		theme.Brightness = BrightnessDark
		return theme
	}()
)
