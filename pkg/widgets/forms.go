package widgets

import (
	"fmt"
	"godin-framework/pkg/core"
	"godin-framework/pkg/renderer"
	"strings"
)

// TextField represents a text input widget with full Flutter properties
type TextField struct {
	ID                            string
	Style                         string
	Class                         string
	Controller                    *TextEditingController                                                                // Text editing controller
	FocusNode                     *FocusNode                                                                            // Focus node
	Decoration                    *InputDecoration                                                                      // Input decoration
	KeyboardType                  TextInputType                                                                         // Keyboard type
	TextInputAction               TextInputAction                                                                       // Text input action
	TextCapitalization            TextCapitalization                                                                    // Text capitalization
	TextStyle                     *TextStyle                                                                            // Text style
	StrutStyle                    *StrutStyle                                                                           // Strut style
	TextAlign                     TextAlign                                                                             // Text alignment
	TextAlignVertical             TextAlignVertical                                                                     // Vertical text alignment
	TextDirection                 TextDirection                                                                         // Text direction
	ReadOnly                      bool                                                                                  // Read only
	ToolbarOptions                *ToolbarOptions                                                                       // Toolbar options
	ShowCursor                    *bool                                                                                 // Show cursor
	AutoFocus                     bool                                                                                  // Auto focus
	ObscuringCharacter            string                                                                                // Obscuring character
	ObscureText                   bool                                                                                  // Obscure text
	AutoCorrect                   bool                                                                                  // Auto correct
	SmartDashesType               SmartDashesType                                                                       // Smart dashes type
	SmartQuotesType               SmartQuotesType                                                                       // Smart quotes type
	EnableSuggestions             bool                                                                                  // Enable suggestions
	MaxLines                      *int                                                                                  // Maximum lines
	MinLines                      *int                                                                                  // Minimum lines
	Expands                       bool                                                                                  // Expands
	MaxLength                     *int                                                                                  // Maximum length
	MaxLengthEnforcement          MaxLengthEnforcement                                                                  // Max length enforcement
	OnChanged                     ValueChanged[string]                                                                  // On changed callback
	OnEditingComplete             VoidCallback                                                                          // On editing complete callback
	OnSubmitted                   ValueChanged[string]                                                                  // On submitted callback
	OnAppPrivateCommand           func(string, map[string]interface{})                                                  // On app private command
	InputFormatters               []TextInputFormatter                                                                  // Input formatters
	Enabled                       *bool                                                                                 // Enabled
	CursorWidth                   float64                                                                               // Cursor width
	CursorHeight                  *float64                                                                              // Cursor height
	CursorRadius                  *Radius                                                                               // Cursor radius
	CursorColor                   Color                                                                                 // Cursor color
	SelectionHeightStyle          BoxHeightStyle                                                                        // Selection height style
	SelectionWidthStyle           BoxWidthStyle                                                                         // Selection width style
	KeyboardAppearance            Brightness                                                                            // Keyboard appearance
	ScrollPadding                 EdgeInsetsGeometry                                                                    // Scroll padding
	DragStartBehavior             DragStartBehavior                                                                     // Drag start behavior
	EnableInteractiveSelection    *bool                                                                                 // Enable interactive selection
	SelectionControls             TextSelectionControls                                                                 // Selection controls
	OnTap                         VoidCallback                                                                          // On tap callback
	MouseCursor                   MouseCursor                                                                           // Mouse cursor
	BuildCounter                  func(context *core.Context, currentLength int, isFocused bool, maxLength *int) Widget // Build counter
	ScrollController              *ScrollController                                                                     // Scroll controller
	ScrollPhysics                 ScrollPhysics                                                                         // Scroll physics
	AutoFillHints                 []AutoFillHint                                                                        // Auto fill hints
	RestorationID                 string                                                                                // Restoration ID
	EnableIMEPersonalizedLearning bool                                                                                  // Enable IME personalized learning
}

// Render renders the text field as HTML
func (tf TextField) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	// Determine if this should be a textarea or input
	isTextarea := (tf.MaxLines != nil && *tf.MaxLines > 1) || tf.Expands || tf.KeyboardType == TextInputTypeMultiline

	attrs := buildAttributes(tf.ID, tf.Style, tf.Class+" godin-textfield")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if tf.Style != "" {
		styles = append(styles, tf.Style)
	}

	// Base input styles
	styles = append(styles, "box-sizing: border-box")
	styles = append(styles, "font-family: inherit")

	// Add text styling
	if tf.TextStyle != nil {
		if textStyleCSS := tf.TextStyle.ToCSSString(); textStyleCSS != "" {
			styles = append(styles, textStyleCSS)
		}
	}

	// Add text alignment
	if tf.TextAlign != "" {
		styles = append(styles, fmt.Sprintf("text-align: %s", tf.TextAlign))
	}

	// Add text direction
	if tf.TextDirection != "" {
		styles = append(styles, fmt.Sprintf("direction: %s", tf.TextDirection))
		attrs["dir"] = string(tf.TextDirection)
	}

	// Handle input type
	if !isTextarea {
		inputType := string(tf.KeyboardType)
		if inputType == "" || inputType == "text" {
			inputType = "text"
		}
		attrs["type"] = inputType
	}

	// Add value from controller or direct value
	if tf.Controller != nil && tf.Controller.Text != "" {
		if isTextarea {
			// For textarea, content goes inside the element
		} else {
			attrs["value"] = tf.Controller.Text
		}
	}

	// Add decoration properties
	if tf.Decoration != nil {
		if tf.Decoration.HintText != "" {
			attrs["placeholder"] = tf.Decoration.HintText
		}

		// Add padding from decoration
		if tf.Decoration.ContentPadding != nil {
			styles = append(styles, fmt.Sprintf("padding: %s", tf.Decoration.ContentPadding.ToCSSString()))
		}

		// Add border styling
		if tf.Decoration.Border != nil {
			styles = append(styles, tf.Decoration.Border.ToCSSString())
		}

		// Add fill color
		if tf.Decoration.FillColor != "" {
			styles = append(styles, fmt.Sprintf("background-color: %s", tf.Decoration.FillColor))
		}
	}

	// Handle enabled/disabled state
	enabled := true
	if tf.Enabled != nil {
		enabled = *tf.Enabled
	}
	if !enabled || tf.ReadOnly {
		attrs["disabled"] = "true"
		styles = append(styles, "opacity: 0.6")
	}

	if tf.ReadOnly {
		attrs["readonly"] = "true"
	}

	// Add max length
	if tf.MaxLength != nil {
		attrs["maxlength"] = fmt.Sprintf("%d", *tf.MaxLength)
	}

	// Handle obscure text (password)
	if tf.ObscureText && !isTextarea {
		attrs["type"] = "password"
	}

	// Add autocomplete attributes
	if len(tf.AutoFillHints) > 0 {
		attrs["autocomplete"] = string(tf.AutoFillHints[0])
	}

	// Add autofocus
	if tf.AutoFocus {
		attrs["autofocus"] = "true"
	}

	// Handle text capitalization
	if tf.TextCapitalization != "" {
		switch tf.TextCapitalization {
		case TextCapitalizationWords:
			styles = append(styles, "text-transform: capitalize")
		case TextCapitalizationCharacters:
			styles = append(styles, "text-transform: uppercase")
		}
	}

	// Add cursor styling
	if tf.CursorColor != "" {
		styles = append(styles, fmt.Sprintf("caret-color: %s", tf.CursorColor))
	}

	// Handle multiline properties
	if isTextarea {
		if tf.MinLines != nil {
			attrs["rows"] = fmt.Sprintf("%d", *tf.MinLines)
		}
		if !tf.Expands && tf.MaxLines != nil {
			styles = append(styles, "resize: none")
		}
	}

	// Add event handlers (simplified - would need proper HTMX integration)
	if tf.OnChanged != nil {
		attrs["oninput"] = "handleTextFieldChange(this)"
	}
	if tf.OnSubmitted != nil {
		attrs["onkeypress"] = "handleTextFieldSubmit(event, this)"
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Render the appropriate element
	if isTextarea {
		content := ""
		if tf.Controller != nil {
			content = tf.Controller.Text
		}
		return htmlRenderer.RenderElement("textarea", attrs, content, false)
	} else {
		return htmlRenderer.RenderElement("input", attrs, "", true)
	}
}

// TextFormField represents a text form field widget with full Flutter properties
type TextFormField struct {
	ID                            string
	Style                         string
	Class                         string
	Controller                    *TextEditingController                                                                // Text editing controller
	InitialValue                  string                                                                                // Initial value
	FocusNode                     *FocusNode                                                                            // Focus node
	Decoration                    *InputDecoration                                                                      // Input decoration
	KeyboardType                  TextInputType                                                                         // Keyboard type
	TextCapitalization            TextCapitalization                                                                    // Text capitalization
	TextInputAction               TextInputAction                                                                       // Text input action
	Style_                        *TextStyle                                                                            // Text style (renamed to avoid conflict)
	StrutStyle                    *StrutStyle                                                                           // Strut style
	TextDirection                 TextDirection                                                                         // Text direction
	TextAlign                     TextAlign                                                                             // Text alignment
	TextAlignVertical             TextAlignVertical                                                                     // Vertical text alignment
	AutoFocus                     bool                                                                                  // Auto focus
	ReadOnly                      bool                                                                                  // Read only
	ToolbarOptions                *ToolbarOptions                                                                       // Toolbar options
	ShowCursor                    *bool                                                                                 // Show cursor
	ObscuringCharacter            string                                                                                // Obscuring character
	ObscureText                   bool                                                                                  // Obscure text
	AutoCorrect                   bool                                                                                  // Auto correct
	SmartDashesType               SmartDashesType                                                                       // Smart dashes type
	SmartQuotesType               SmartQuotesType                                                                       // Smart quotes type
	EnableSuggestions             bool                                                                                  // Enable suggestions
	MaxLengthEnforcement          MaxLengthEnforcement                                                                  // Max length enforcement
	MaxLines                      *int                                                                                  // Maximum lines
	MinLines                      *int                                                                                  // Minimum lines
	Expands                       bool                                                                                  // Expands
	MaxLength                     *int                                                                                  // Maximum length
	OnChanged                     ValueChanged[string]                                                                  // On changed callback
	OnTap                         GestureTapCallback                                                                    // On tap callback
	OnEditingComplete             VoidCallback                                                                          // On editing complete callback
	OnFieldSubmitted              ValueChanged[string]                                                                  // On field submitted callback
	OnSaved                       FormFieldSetter[string]                                                               // On saved callback
	Validator                     FormFieldValidator[string]                                                            // Validator function
	InputFormatters               []TextInputFormatter                                                                  // Input formatters
	Enabled                       *bool                                                                                 // Enabled state
	CursorWidth                   *float64                                                                              // Cursor width
	CursorHeight                  *float64                                                                              // Cursor height
	CursorRadius                  *Radius                                                                               // Cursor radius
	CursorColor                   Color                                                                                 // Cursor color
	KeyboardAppearance            Brightness                                                                            // Keyboard appearance
	ScrollPadding                 EdgeInsetsGeometry                                                                    // Scroll padding
	EnableInteractiveSelection    *bool                                                                                 // Enable interactive selection
	SelectionControls             TextSelectionControls                                                                 // Selection controls
	BuildCounter                  func(context *core.Context, currentLength int, isFocused bool, maxLength *int) Widget // Build counter
	ScrollPhysics                 ScrollPhysics                                                                         // Scroll physics
	AutoFillHints                 []AutoFillHint                                                                        // Auto fill hints
	AutovalidateMode              AutovalidateMode                                                                      // Auto validate mode
	ScrollController              *ScrollController                                                                     // Scroll controller
	RestorationId                 string                                                                                // Restoration ID
	EnableIMEPersonalizedLearning *bool                                                                                 // Enable IME personalized learning
	MouseCursor                   MouseCursor                                                                           // Mouse cursor
}

// Render renders the text form field as HTML
func (tff TextFormField) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	// Determine if this should be a textarea or input
	isTextarea := (tff.MaxLines != nil && *tff.MaxLines > 1) || tff.Expands || tff.KeyboardType == TextInputTypeMultiline

	attrs := buildAttributes(tff.ID, tff.Style, tff.Class+" godin-textformfield")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if tff.Style != "" {
		styles = append(styles, tff.Style)
	}

	// Add text styling
	if tff.Style_ != nil {
		if tff.Style_.Color != "" {
			styles = append(styles, fmt.Sprintf("color: %s", tff.Style_.Color))
		}
		if tff.Style_.FontSize != nil {
			styles = append(styles, fmt.Sprintf("font-size: %.1fpx", *tff.Style_.FontSize))
		}
		if tff.Style_.FontWeight != "" {
			styles = append(styles, fmt.Sprintf("font-weight: %s", tff.Style_.FontWeight))
		}
		if tff.Style_.FontFamily != "" {
			styles = append(styles, fmt.Sprintf("font-family: %s", tff.Style_.FontFamily))
		}
	}

	// Add text alignment
	if tff.TextAlign != "" {
		styles = append(styles, fmt.Sprintf("text-align: %s", tff.TextAlign))
	}

	// Add keyboard type attributes
	if !isTextarea {
		switch tff.KeyboardType {
		case TextInputTypeEmailAddress:
			attrs["type"] = "email"
		case TextInputTypeNumber:
			attrs["type"] = "number"
		case TextInputTypePhone:
			attrs["type"] = "tel"
		case TextInputTypeURL:
			attrs["type"] = "url"
		case TextInputTypeDatetime:
			attrs["type"] = "datetime-local"
		default:
			attrs["type"] = "text"
		}
	}

	// Add decoration properties
	if tff.Decoration != nil {
		if tff.Decoration.HintText != "" {
			attrs["placeholder"] = tff.Decoration.HintText
		}

		// Add padding from decoration
		if tff.Decoration.ContentPadding != nil {
			styles = append(styles, fmt.Sprintf("padding: %s", tff.Decoration.ContentPadding.ToCSSString()))
		}

		// Add border styling
		if tff.Decoration.Border != nil {
			styles = append(styles, tff.Decoration.Border.ToCSSString())
		}

		// Add fill color
		if tff.Decoration.FillColor != "" {
			styles = append(styles, fmt.Sprintf("background-color: %s", tff.Decoration.FillColor))
		}
	}

	// Handle enabled/disabled state
	enabled := true
	if tff.Enabled != nil {
		enabled = *tff.Enabled
	}
	if !enabled || tff.ReadOnly {
		attrs["disabled"] = "true"
		styles = append(styles, "opacity: 0.6")
	}

	if tff.ReadOnly {
		attrs["readonly"] = "true"
	}

	// Add max length
	if tff.MaxLength != nil {
		attrs["maxlength"] = fmt.Sprintf("%d", *tff.MaxLength)
	}

	// Handle obscure text (password)
	if tff.ObscureText && !isTextarea {
		attrs["type"] = "password"
	}

	// Add autocomplete attributes
	if len(tff.AutoFillHints) > 0 {
		attrs["autocomplete"] = string(tff.AutoFillHints[0])
	}

	// Add autofocus
	if tff.AutoFocus {
		attrs["autofocus"] = "true"
	}

	// Handle text capitalization
	if tff.TextCapitalization != "" {
		switch tff.TextCapitalization {
		case TextCapitalizationWords:
			styles = append(styles, "text-transform: capitalize")
		case TextCapitalizationCharacters:
			styles = append(styles, "text-transform: uppercase")
		}
	}

	// Add cursor styling
	if tff.CursorColor != "" {
		styles = append(styles, fmt.Sprintf("caret-color: %s", tff.CursorColor))
	}

	// Handle multiline properties
	if isTextarea {
		if tff.MinLines != nil {
			attrs["rows"] = fmt.Sprintf("%d", *tff.MinLines)
		}
		if !tff.Expands && tff.MaxLines != nil {
			styles = append(styles, "resize: none")
		}
	}

	// Add form validation attributes
	if tff.Validator != nil {
		attrs["data-validator"] = "true"
	}

	// Add event handlers (simplified - would need proper HTMX integration)
	if tff.OnChanged != nil {
		attrs["oninput"] = "handleTextFormFieldChange(this)"
	}
	if tff.OnFieldSubmitted != nil {
		attrs["onkeypress"] = "handleTextFormFieldSubmit(event, this)"
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Determine initial value
	initialValue := tff.InitialValue
	if tff.Controller != nil && tff.Controller.Text != "" {
		initialValue = tff.Controller.Text
	}

	// Render the appropriate element
	if isTextarea {
		return htmlRenderer.RenderElement("textarea", attrs, initialValue, false)
	} else {
		if initialValue != "" {
			attrs["value"] = initialValue
		}
		return htmlRenderer.RenderElement("input", attrs, "", true)
	}
}

// Switch represents a switch widget with full Flutter properties
type Switch struct {
	ID                        string
	Style                     string
	Class                     string
	Value                     bool                          // Switch value
	OnChanged                 ValueChanged[bool]            // On changed callback
	ActiveColor               Color                         // Active color
	ActiveTrackColor          Color                         // Active track color
	InactiveThumbColor        Color                         // Inactive thumb color
	InactiveTrackColor        Color                         // Inactive track color
	ActiveThumbImage          ImageProvider                 // Active thumb image
	OnActiveThumbImageError   ImageErrorListener            // Active thumb image error callback
	InactiveThumbImage        ImageProvider                 // Inactive thumb image
	OnInactiveThumbImageError ImageErrorListener            // Inactive thumb image error callback
	ThumbColor                *MaterialStateProperty[Color] // Thumb color
	TrackColor                *MaterialStateProperty[Color] // Track color
	MaterialTapTargetSize     MaterialTapTargetSize         // Tap target size
	DragStartBehavior         DragStartBehavior             // Drag start behavior
	MouseCursor               MouseCursor                   // Mouse cursor
	FocusColor                Color                         // Focus color
	HoverColor                Color                         // Hover color
	OverlayColor              *MaterialStateProperty[Color] // Overlay color
	SplashRadius              *float64                      // Splash radius
	FocusNode                 *FocusNode                    // Focus node
	AutoFocus                 bool                          // Auto focus
}

// Render renders the switch as HTML
func (s Switch) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	// Create a container for the switch
	containerAttrs := buildAttributes(s.ID+"_container", s.Style, s.Class+" godin-switch-container")

	// Build inline styles for container
	var containerStyles []string

	// Add custom style if provided
	if s.Style != "" {
		containerStyles = append(containerStyles, s.Style)
	}

	// Base switch container styles
	containerStyles = append(containerStyles, "display: inline-flex")
	containerStyles = append(containerStyles, "align-items: center")
	containerStyles = append(containerStyles, "position: relative")

	// Handle tap target size
	if s.MaterialTapTargetSize == MaterialTapTargetSizePadded {
		containerStyles = append(containerStyles, "min-width: 48px")
		containerStyles = append(containerStyles, "min-height: 48px")
	}

	// Combine container styles
	if len(containerStyles) > 0 {
		containerAttrs["style"] = strings.Join(containerStyles, "; ")
	}

	// Create the actual input element
	inputAttrs := make(map[string]string)
	inputAttrs["type"] = "checkbox"
	inputAttrs["class"] = "godin-switch-input"

	if s.ID != "" {
		inputAttrs["id"] = s.ID
	}

	if s.Value {
		inputAttrs["checked"] = "checked"
	}

	// Add autofocus
	if s.AutoFocus {
		inputAttrs["autofocus"] = "true"
	}

	// Build input styles
	var inputStyles []string

	// Hide the default checkbox appearance
	inputStyles = append(inputStyles, "appearance: none")
	inputStyles = append(inputStyles, "-webkit-appearance: none")
	inputStyles = append(inputStyles, "-moz-appearance: none")
	inputStyles = append(inputStyles, "width: 52px")
	inputStyles = append(inputStyles, "height: 32px")
	inputStyles = append(inputStyles, "border-radius: 16px")
	inputStyles = append(inputStyles, "position: relative")
	inputStyles = append(inputStyles, "cursor: pointer")
	inputStyles = append(inputStyles, "transition: all 0.3s ease")

	// Add track colors
	if s.Value {
		if s.ActiveTrackColor != "" {
			inputStyles = append(inputStyles, fmt.Sprintf("background-color: %s", s.ActiveTrackColor))
		} else {
			inputStyles = append(inputStyles, "background-color: #2196F3")
		}
	} else {
		if s.InactiveTrackColor != "" {
			inputStyles = append(inputStyles, fmt.Sprintf("background-color: %s", s.InactiveTrackColor))
		} else {
			inputStyles = append(inputStyles, "background-color: #ccc")
		}
	}

	// Add border
	inputStyles = append(inputStyles, "border: none")
	inputStyles = append(inputStyles, "outline: none")

	// Add pseudo-element for thumb using CSS
	inputStyles = append(inputStyles, "position: relative")

	// Combine input styles
	if len(inputStyles) > 0 {
		inputAttrs["style"] = strings.Join(inputStyles, "; ")
	}

	// Add event handlers
	if s.OnChanged != nil {
		inputAttrs["onchange"] = "handleSwitchChange(this)"
	}

	// Create the thumb element (using a span positioned absolutely)
	thumbAttrs := make(map[string]string)
	thumbAttrs["class"] = "godin-switch-thumb"

	var thumbStyles []string
	thumbStyles = append(thumbStyles, "position: absolute")
	thumbStyles = append(thumbStyles, "top: 2px")
	thumbStyles = append(thumbStyles, "width: 28px")
	thumbStyles = append(thumbStyles, "height: 28px")
	thumbStyles = append(thumbStyles, "border-radius: 50%")
	thumbStyles = append(thumbStyles, "transition: all 0.3s ease")
	thumbStyles = append(thumbStyles, "pointer-events: none")

	// Position thumb based on switch state
	if s.Value {
		thumbStyles = append(thumbStyles, "left: 22px")
		if s.ActiveColor != "" {
			thumbStyles = append(thumbStyles, fmt.Sprintf("background-color: %s", s.ActiveColor))
		} else {
			thumbStyles = append(thumbStyles, "background-color: white")
		}
	} else {
		thumbStyles = append(thumbStyles, "left: 2px")
		if s.InactiveThumbColor != "" {
			thumbStyles = append(thumbStyles, fmt.Sprintf("background-color: %s", s.InactiveThumbColor))
		} else {
			thumbStyles = append(thumbStyles, "background-color: white")
		}
	}

	thumbAttrs["style"] = strings.Join(thumbStyles, "; ")

	// Render the input element
	inputHTML := htmlRenderer.RenderElement("input", inputAttrs, "", true)

	// Render the thumb element
	thumbHTML := htmlRenderer.RenderElement("span", thumbAttrs, "", false)

	// Combine input and thumb in container
	content := inputHTML + thumbHTML

	return htmlRenderer.RenderElement("div", containerAttrs, content, false)
}

// Button represents a button widget
type Button struct {
	ID       string
	Style    string
	Class    string
	Text     string
	OnClick  func() // Go function callback
	Type     string // "primary", "secondary", "danger"
	Disabled bool
	// HTMX attributes
	HxPost   string // hx-post
	HxGet    string // hx-get
	HxTarget string // hx-target
	HxSwap   string // hx-swap
}

// Render renders the button as HTML
func (b Button) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(b.ID, b.Style, b.Class+" godin-button")

	if b.Type != "" {
		attrs["class"] += " godin-button-" + b.Type
	}

	if b.Disabled {
		attrs["disabled"] = "disabled"
	}

	// Register Go function callback if provided
	if b.OnClick != nil {
		handlerID := ctx.RegisterHandler(func(ctx *core.Context) Widget {
			b.OnClick()
			return nil // Return nil for callbacks that don't return widgets
		})

		attrs["hx-post"] = "/handlers/" + handlerID
		attrs["hx-trigger"] = "click"
	}

	// Add HTMX attributes if provided (these override the OnClick handler)
	if b.HxPost != "" {
		attrs["hx-post"] = b.HxPost
		attrs["hx-trigger"] = "click"
	}
	if b.HxGet != "" {
		attrs["hx-get"] = b.HxGet
		attrs["hx-trigger"] = "click"
	}
	if b.HxTarget != "" {
		attrs["hx-target"] = b.HxTarget
	}
	if b.HxSwap != "" {
		attrs["hx-swap"] = b.HxSwap
	}

	return htmlRenderer.RenderElement("button", attrs, b.Text, false)
}

// Checkbox represents a checkbox widget with full Flutter properties
type Checkbox struct {
	ID                    string
	Style                 string
	Class                 string
	Value                 *bool                              // Checkbox value (null for indeterminate)
	Tristate              bool                               // Allow three states
	OnChanged             ValueChanged[bool]                 // On changed callback
	ActiveColor           Color                              // Active color
	FillColor             *MaterialStateProperty[Color]      // Fill color
	CheckColor            Color                              // Check color
	FocusColor            Color                              // Focus color
	HoverColor            Color                              // Hover color
	OverlayColor          *MaterialStateProperty[Color]      // Overlay color
	SplashRadius          *float64                           // Splash radius
	MaterialTapTargetSize MaterialTapTargetSize              // Tap target size
	VisualDensity         *VisualDensity                     // Visual density
	FocusNode             *FocusNode                         // Focus node
	AutoFocus             bool                               // Auto focus
	Shape                 OutlinedBorder                     // Shape
	Side                  *MaterialStateProperty[BorderSide] // Side
	IsError               bool                               // Is error state
	SemanticLabel         string                             // Semantic label
}

// Render renders the checkbox as HTML
func (c Checkbox) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(c.ID, c.Style, c.Class+" godin-checkbox")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if c.Style != "" {
		styles = append(styles, c.Style)
	}

	// Base checkbox styles
	styles = append(styles, "display: inline-flex")
	styles = append(styles, "align-items: center")
	styles = append(styles, "cursor: pointer")
	styles = append(styles, "user-select: none")

	// Checkbox input attributes
	attrs["type"] = "checkbox"

	// Handle checkbox value and states
	if c.Value != nil {
		if *c.Value {
			attrs["checked"] = "checked"
		}
	} else if c.Tristate {
		// Indeterminate state
		attrs["indeterminate"] = "true"
	}

	// Add colors
	if c.ActiveColor != "" {
		styles = append(styles, fmt.Sprintf("accent-color: %s", c.ActiveColor))
	}

	// Handle fill color from MaterialStateProperty
	if c.FillColor != nil {
		styles = append(styles, fmt.Sprintf("background-color: %s", c.FillColor.Default))
	}

	// Add focus and hover colors (simplified CSS approach)
	if c.FocusColor != "" {
		styles = append(styles, fmt.Sprintf("outline-color: %s", c.FocusColor))
	}

	// Handle disabled state (checkbox is disabled if OnChanged is nil)
	if c.OnChanged == nil {
		attrs["disabled"] = "true"
		styles = append(styles, "opacity: 0.6")
		styles = append(styles, "cursor: not-allowed")
	}

	// Add visual density adjustments
	if c.VisualDensity != nil {
		width := 18.0 + c.VisualDensity.Horizontal*4
		height := 18.0 + c.VisualDensity.Vertical*4
		styles = append(styles, fmt.Sprintf("width: %.1fpx", width))
		styles = append(styles, fmt.Sprintf("height: %.1fpx", height))
	} else {
		styles = append(styles, "width: 18px")
		styles = append(styles, "height: 18px")
	}

	// Handle tap target size
	if c.MaterialTapTargetSize == MaterialTapTargetSizePadded {
		styles = append(styles, "min-width: 48px")
		styles = append(styles, "min-height: 48px")
	}

	// Add shape styling
	if c.Shape != nil {
		styles = append(styles, c.Shape.ToCSSString())
	}

	// Handle error state
	if c.IsError {
		styles = append(styles, "border-color: #f44336")
	}

	// Add autofocus
	if c.AutoFocus {
		attrs["autofocus"] = "true"
	}

	// Add semantic label
	if c.SemanticLabel != "" {
		attrs["aria-label"] = c.SemanticLabel
	}

	// Add event handlers (simplified)
	if c.OnChanged != nil {
		attrs["onchange"] = "handleCheckboxChange(this)"
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	return htmlRenderer.RenderElement("input", attrs, "", true)
}

// Radio represents a radio button widget with full Flutter properties
type Radio[T comparable] struct {
	ID                         string
	Style                      string
	Class                      string
	Value                      T                             // Radio value
	GroupValue                 *T                            // Group value
	OnChanged                  ValueChanged[T]               // On changed callback
	MouseCursor                MouseCursor                   // Mouse cursor
	ToggleableActiveColor      Color                         // Toggleable active color
	FillColor                  *MaterialStateProperty[Color] // Fill color
	FocusColor                 Color                         // Focus color
	HoverColor                 Color                         // Hover color
	OverlayColor               *MaterialStateProperty[Color] // Overlay color
	SplashRadius               *float64                      // Splash radius
	MaterialTapTargetSize      MaterialTapTargetSize         // Tap target size
	VisualDensity              *VisualDensity                // Visual density
	FocusNode                  *FocusNode                    // Focus node
	AutoFocus                  bool                          // Auto focus
	UseCupertinoCheckmarkStyle bool                          // Use Cupertino checkmark style
}

// Render renders the radio button as HTML
func (r Radio[T]) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(r.ID, r.Style, r.Class+" godin-radio")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if r.Style != "" {
		styles = append(styles, r.Style)
	}

	// Base radio styles
	styles = append(styles, "display: inline-flex")
	styles = append(styles, "align-items: center")
	styles = append(styles, "cursor: pointer")
	styles = append(styles, "user-select: none")

	// Radio input attributes
	attrs["type"] = "radio"

	// Convert value to string for HTML
	attrs["value"] = fmt.Sprintf("%v", r.Value)

	// Check if this radio is selected
	if r.GroupValue != nil && *r.GroupValue == r.Value {
		attrs["checked"] = "checked"
	}

	// Add colors
	if r.ToggleableActiveColor != "" {
		styles = append(styles, fmt.Sprintf("accent-color: %s", r.ToggleableActiveColor))
	}

	// Handle fill color from MaterialStateProperty
	if r.FillColor != nil {
		styles = append(styles, fmt.Sprintf("background-color: %s", r.FillColor.Default))
	}

	// Add focus and hover colors (simplified CSS approach)
	if r.FocusColor != "" {
		styles = append(styles, fmt.Sprintf("outline-color: %s", r.FocusColor))
	}

	// Handle disabled state (radio is disabled if OnChanged is nil)
	if r.OnChanged == nil {
		attrs["disabled"] = "true"
		styles = append(styles, "opacity: 0.6")
		styles = append(styles, "cursor: not-allowed")
	}

	// Add visual density adjustments
	if r.VisualDensity != nil {
		width := 18.0 + r.VisualDensity.Horizontal*4
		height := 18.0 + r.VisualDensity.Vertical*4
		styles = append(styles, fmt.Sprintf("width: %.1fpx", width))
		styles = append(styles, fmt.Sprintf("height: %.1fpx", height))
	} else {
		styles = append(styles, "width: 18px")
		styles = append(styles, "height: 18px")
	}

	// Handle tap target size
	if r.MaterialTapTargetSize == MaterialTapTargetSizePadded {
		styles = append(styles, "min-width: 48px")
		styles = append(styles, "min-height: 48px")
	}

	// Add autofocus
	if r.AutoFocus {
		attrs["autofocus"] = "true"
	}

	// Add event handlers (simplified)
	if r.OnChanged != nil {
		attrs["onchange"] = "handleRadioChange(this)"
	}

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	return htmlRenderer.RenderElement("input", attrs, "", true)
}

// DropdownOption represents an option in a dropdown
type DropdownOption struct {
	Value string
	Label string
}

// Dropdown represents a dropdown widget
type Dropdown struct {
	HTMXWidget
	Value    string
	Options  []DropdownOption
	OnChange string
	Disabled bool
}

// Render renders the dropdown as HTML
func (d Dropdown) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := d.buildHTMXAttributes()
	attrs["class"] += " godin-dropdown"

	if d.Disabled {
		attrs["disabled"] = "disabled"
	}

	// Add HTMX for onChange
	if d.OnChange != "" && d.HTMX.Post == "" {
		attrs["hx-post"] = d.OnChange
		attrs["hx-trigger"] = "change"
	}

	if d.Style != "" {
		attrs["style"] = d.Style
	}

	// Build options
	var options []string
	for _, option := range d.Options {
		optionAttrs := map[string]string{
			"value": option.Value,
		}
		if option.Value == d.Value {
			optionAttrs["selected"] = "selected"
		}
		options = append(options, htmlRenderer.RenderElement("option", optionAttrs, option.Label, false))
	}

	return htmlRenderer.RenderContainer("select", attrs, options)
}

// Slider represents a slider widget with full Flutter properties
type Slider struct {
	ID                        string
	Style                     string
	Class                     string
	Value                     float64                   // Current value
	OnChanged                 ValueChanged[float64]     // On changed callback
	OnChangeStart             ValueChanged[float64]     // On change start callback
	OnChangeEnd               ValueChanged[float64]     // On change end callback
	Min                       float64                   // Minimum value
	Max                       float64                   // Maximum value
	Divisions                 *int                      // Number of divisions
	Label                     string                    // Label text
	ActiveColor               Color                     // Active color
	InactiveColor             Color                     // Inactive color
	ThumbColor                Color                     // Thumb color
	OverlayColor              Color                     // Overlay color
	MouseCursor               MouseCursor               // Mouse cursor
	SemanticFormatterCallback SemanticFormatterCallback // Semantic formatter callback
	FocusNode                 *FocusNode                // Focus node
	AutoFocus                 bool                      // Auto focus
}

// SemanticFormatterCallback represents a semantic formatter callback
type SemanticFormatterCallback func(float64) string

// Render renders the slider as HTML
func (s Slider) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	// Create a container for the slider
	containerAttrs := buildAttributes(s.ID+"_container", s.Style, s.Class+" godin-slider-container")

	// Build inline styles for container
	var containerStyles []string

	// Add custom style if provided
	if s.Style != "" {
		containerStyles = append(containerStyles, s.Style)
	}

	// Base slider container styles
	containerStyles = append(containerStyles, "display: flex")
	containerStyles = append(containerStyles, "flex-direction: column")
	containerStyles = append(containerStyles, "align-items: center")
	containerStyles = append(containerStyles, "width: 100%")

	// Combine container styles
	if len(containerStyles) > 0 {
		containerAttrs["style"] = strings.Join(containerStyles, "; ")
	}

	// Create the actual input element
	inputAttrs := make(map[string]string)
	inputAttrs["type"] = "range"
	inputAttrs["class"] = "godin-slider-input"

	if s.ID != "" {
		inputAttrs["id"] = s.ID
	}

	inputAttrs["value"] = fmt.Sprintf("%.2f", s.Value)
	inputAttrs["min"] = fmt.Sprintf("%.2f", s.Min)
	inputAttrs["max"] = fmt.Sprintf("%.2f", s.Max)

	// Add step based on divisions
	if s.Divisions != nil && *s.Divisions > 0 {
		step := (s.Max - s.Min) / float64(*s.Divisions)
		inputAttrs["step"] = fmt.Sprintf("%.6f", step)
	} else {
		inputAttrs["step"] = "any"
	}

	// Add autofocus
	if s.AutoFocus {
		inputAttrs["autofocus"] = "true"
	}

	// Build input styles
	var inputStyles []string

	// Base slider styles
	inputStyles = append(inputStyles, "width: 100%")
	inputStyles = append(inputStyles, "height: 4px")
	inputStyles = append(inputStyles, "border-radius: 2px")
	inputStyles = append(inputStyles, "outline: none")
	inputStyles = append(inputStyles, "cursor: pointer")

	// Add track colors
	if s.ActiveColor != "" {
		inputStyles = append(inputStyles, fmt.Sprintf("accent-color: %s", s.ActiveColor))
	}

	if s.InactiveColor != "" {
		inputStyles = append(inputStyles, fmt.Sprintf("background: %s", s.InactiveColor))
	} else {
		inputStyles = append(inputStyles, "background: #ddd")
	}

	// Combine input styles
	if len(inputStyles) > 0 {
		inputAttrs["style"] = strings.Join(inputStyles, "; ")
	}

	// Add event handlers
	if s.OnChanged != nil {
		inputAttrs["oninput"] = "handleSliderChange(this)"
	}
	if s.OnChangeStart != nil {
		inputAttrs["onmousedown"] = "handleSliderChangeStart(this)"
		inputAttrs["ontouchstart"] = "handleSliderChangeStart(this)"
	}
	if s.OnChangeEnd != nil {
		inputAttrs["onmouseup"] = "handleSliderChangeEnd(this)"
		inputAttrs["ontouchend"] = "handleSliderChangeEnd(this)"
	}

	// Render the input element
	inputHTML := htmlRenderer.RenderElement("input", inputAttrs, "", true)

	// Create label if provided
	var labelHTML string
	if s.Label != "" {
		labelAttrs := make(map[string]string)
		labelAttrs["class"] = "godin-slider-label"
		labelAttrs["style"] = "margin-top: 8px; font-size: 12px; color: #666;"
		labelHTML = htmlRenderer.RenderElement("div", labelAttrs, s.Label, false)
	}

	// Combine input and label in container
	content := inputHTML + labelHTML

	return htmlRenderer.RenderElement("div", containerAttrs, content, false)
}

// ElevatedButton represents an elevated button widget with full Flutter properties
type ElevatedButton struct {
	ID               string
	Style            string
	Class            string
	OnPressed        VoidCallback              // Callback when pressed
	OnLongPress      VoidCallback              // Callback when long pressed
	OnHover          ValueChanged[bool]        // Callback when hovered
	OnFocusChange    ValueChanged[bool]        // Callback when focus changes
	ButtonStyle      *ButtonStyle              // Button style
	FocusNode        *FocusNode                // Focus node
	AutoFocus        bool                      // Auto focus
	ClipBehavior     Clip                      // Clip behavior
	StatesController *MaterialStatesController // States controller
	Child            Widget                    // Child widget
}

// FocusNode represents a focus node (simplified)
type FocusNode struct {
	HasFocus bool
}

// MaterialStatesController represents material states controller (simplified)
type MaterialStatesController struct {
	States []MaterialState
}

// MaterialState enum
type MaterialState string

const (
	MaterialStateHovered  MaterialState = "hovered"
	MaterialStateFocused  MaterialState = "focused"
	MaterialStatePressed  MaterialState = "pressed"
	MaterialStateDragged  MaterialState = "dragged"
	MaterialStateSelected MaterialState = "selected"
	MaterialStateScrolled MaterialState = "scrolled"
	MaterialStateDisabled MaterialState = "disabled"
	MaterialStateError    MaterialState = "error"
)

// Render renders the elevated button as HTML
func (eb ElevatedButton) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(eb.ID, eb.Style, eb.Class+" godin-elevated-button")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if eb.Style != "" {
		styles = append(styles, eb.Style)
	}

	// Base button styles
	styles = append(styles, "display: inline-flex")
	styles = append(styles, "align-items: center")
	styles = append(styles, "justify-content: center")
	styles = append(styles, "border: none")
	styles = append(styles, "cursor: pointer")
	styles = append(styles, "text-decoration: none")
	styles = append(styles, "outline: none")
	styles = append(styles, "user-select: none")

	// Default elevated button styles
	styles = append(styles, "background-color: #1976d2")
	styles = append(styles, "color: white")
	styles = append(styles, "border-radius: 4px")
	styles = append(styles, "padding: 8px 16px")
	styles = append(styles, "min-height: 36px")
	styles = append(styles, "box-shadow: 0 2px 4px rgba(0,0,0,0.2)")

	// Apply ButtonStyle if provided
	if eb.ButtonStyle != nil {
		if eb.ButtonStyle.BackgroundColor != nil {
			styles = append(styles, fmt.Sprintf("background-color: %s", eb.ButtonStyle.BackgroundColor.Default))
		}
		if eb.ButtonStyle.ForegroundColor != nil {
			styles = append(styles, fmt.Sprintf("color: %s", eb.ButtonStyle.ForegroundColor.Default))
		}
		if eb.ButtonStyle.Padding != nil {
			styles = append(styles, fmt.Sprintf("padding: %s", eb.ButtonStyle.Padding.Default.ToCSSString()))
		}
		if eb.ButtonStyle.Shape != nil && eb.ButtonStyle.Shape.Default != nil {
			styles = append(styles, eb.ButtonStyle.Shape.Default.ToCSSString())
		}
		if eb.ButtonStyle.TextStyle != nil {
			if textStyleCSS := eb.ButtonStyle.TextStyle.ToCSSString(); textStyleCSS != "" {
				styles = append(styles, textStyleCSS)
			}
		}
	}

	// Handle disabled state
	if eb.OnPressed == nil {
		styles = append(styles, "opacity: 0.6")
		styles = append(styles, "cursor: not-allowed")
		attrs["disabled"] = "true"
	}

	// Add hover and focus styles via CSS
	styles = append(styles, "transition: all 0.2s ease")

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Add HTMX event handlers for OnPressed callback
	if eb.OnPressed != nil {
		handlerID := ctx.RegisterHandler(func(ctx *core.Context) Widget {
			eb.OnPressed()
			return nil // Return nil for callbacks that don't return widgets
		})

		attrs["hx-post"] = "/handlers/" + handlerID
		attrs["hx-trigger"] = "click"
	}

	// Add accessibility attributes
	attrs["role"] = "button"
	attrs["tabindex"] = "0"

	if eb.AutoFocus {
		attrs["autofocus"] = "true"
	}

	// Render child content
	content := ""
	if eb.Child != nil {
		content = eb.Child.Render(ctx)
	}

	return htmlRenderer.RenderElement("button", attrs, content, false)
}

// TextButton represents a text button widget with full Flutter properties
type TextButton struct {
	ID               string
	Style            string
	Class            string
	OnPressed        VoidCallback              // Callback when pressed
	OnLongPress      VoidCallback              // Callback when long pressed
	OnHover          ValueChanged[bool]        // Callback when hovered
	OnFocusChange    ValueChanged[bool]        // Callback when focus changes
	ButtonStyle      *ButtonStyle              // Button style
	FocusNode        *FocusNode                // Focus node
	AutoFocus        bool                      // Auto focus
	ClipBehavior     Clip                      // Clip behavior
	StatesController *MaterialStatesController // States controller
	Child            Widget                    // Child widget
}

// Render renders the text button as HTML
func (tb TextButton) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(tb.ID, tb.Style, tb.Class+" godin-text-button")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if tb.Style != "" {
		styles = append(styles, tb.Style)
	}

	// Base button styles
	styles = append(styles, "display: inline-flex")
	styles = append(styles, "align-items: center")
	styles = append(styles, "justify-content: center")
	styles = append(styles, "border: none")
	styles = append(styles, "cursor: pointer")
	styles = append(styles, "text-decoration: none")
	styles = append(styles, "outline: none")
	styles = append(styles, "user-select: none")

	// Default text button styles
	styles = append(styles, "background-color: transparent")
	styles = append(styles, "color: #1976d2")
	styles = append(styles, "border-radius: 4px")
	styles = append(styles, "padding: 8px 16px")
	styles = append(styles, "min-height: 36px")

	// Apply ButtonStyle if provided
	if tb.ButtonStyle != nil {
		if tb.ButtonStyle.BackgroundColor != nil {
			styles = append(styles, fmt.Sprintf("background-color: %s", tb.ButtonStyle.BackgroundColor.Default))
		}
		if tb.ButtonStyle.ForegroundColor != nil {
			styles = append(styles, fmt.Sprintf("color: %s", tb.ButtonStyle.ForegroundColor.Default))
		}
		if tb.ButtonStyle.Padding != nil {
			styles = append(styles, fmt.Sprintf("padding: %s", tb.ButtonStyle.Padding.Default.ToCSSString()))
		}
		if tb.ButtonStyle.Shape != nil && tb.ButtonStyle.Shape.Default != nil {
			styles = append(styles, tb.ButtonStyle.Shape.Default.ToCSSString())
		}
		if tb.ButtonStyle.TextStyle != nil {
			if textStyleCSS := tb.ButtonStyle.TextStyle.ToCSSString(); textStyleCSS != "" {
				styles = append(styles, textStyleCSS)
			}
		}
	}

	// Handle disabled state
	if tb.OnPressed == nil {
		styles = append(styles, "opacity: 0.6")
		styles = append(styles, "cursor: not-allowed")
		attrs["disabled"] = "true"
	}

	// Add hover and focus styles via CSS
	styles = append(styles, "transition: all 0.2s ease")

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Add HTMX event handlers for OnPressed callback
	if tb.OnPressed != nil {
		handlerID := ctx.RegisterHandler(func(ctx *core.Context) Widget {
			tb.OnPressed()
			return nil // Return nil for callbacks that don't return widgets
		})

		attrs["hx-post"] = "/handlers/" + handlerID
		attrs["hx-trigger"] = "click"
	}

	// Add accessibility attributes
	attrs["role"] = "button"
	attrs["tabindex"] = "0"

	if tb.AutoFocus {
		attrs["autofocus"] = "true"
	}

	// Render child content
	content := ""
	if tb.Child != nil {
		content = tb.Child.Render(ctx)
	}

	return htmlRenderer.RenderElement("button", attrs, content, false)
}

// OutlinedButton represents an outlined button widget with full Flutter properties
type OutlinedButton struct {
	ID               string
	Style            string
	Class            string
	OnPressed        VoidCallback              // Callback when pressed
	OnLongPress      VoidCallback              // Callback when long pressed
	OnHover          ValueChanged[bool]        // Callback when hovered
	OnFocusChange    ValueChanged[bool]        // Callback when focus changes
	ButtonStyle      *ButtonStyle              // Button style
	FocusNode        *FocusNode                // Focus node
	AutoFocus        bool                      // Auto focus
	ClipBehavior     Clip                      // Clip behavior
	StatesController *MaterialStatesController // States controller
	Child            Widget                    // Child widget
}

// Render renders the outlined button as HTML
func (ob OutlinedButton) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(ob.ID, ob.Style, ob.Class+" godin-outlined-button")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if ob.Style != "" {
		styles = append(styles, ob.Style)
	}

	// Base button styles
	styles = append(styles, "display: inline-flex")
	styles = append(styles, "align-items: center")
	styles = append(styles, "justify-content: center")
	styles = append(styles, "cursor: pointer")
	styles = append(styles, "text-decoration: none")
	styles = append(styles, "outline: none")
	styles = append(styles, "user-select: none")

	// Default outlined button styles
	styles = append(styles, "background-color: transparent")
	styles = append(styles, "color: #1976d2")
	styles = append(styles, "border: 1px solid #1976d2")
	styles = append(styles, "border-radius: 4px")
	styles = append(styles, "padding: 8px 16px")
	styles = append(styles, "min-height: 36px")

	// Apply ButtonStyle if provided
	if ob.ButtonStyle != nil {
		if ob.ButtonStyle.BackgroundColor != nil {
			styles = append(styles, fmt.Sprintf("background-color: %s", ob.ButtonStyle.BackgroundColor.Default))
		}
		if ob.ButtonStyle.ForegroundColor != nil {
			styles = append(styles, fmt.Sprintf("color: %s", ob.ButtonStyle.ForegroundColor.Default))
		}
		if ob.ButtonStyle.Side != nil {
			side := ob.ButtonStyle.Side.Default
			styles = append(styles, fmt.Sprintf("border: %.1fpx %s %s", side.Width, side.Style, side.Color))
		}
		if ob.ButtonStyle.Padding != nil {
			styles = append(styles, fmt.Sprintf("padding: %s", ob.ButtonStyle.Padding.Default.ToCSSString()))
		}
		if ob.ButtonStyle.Shape != nil && ob.ButtonStyle.Shape.Default != nil {
			styles = append(styles, ob.ButtonStyle.Shape.Default.ToCSSString())
		}
		if ob.ButtonStyle.TextStyle != nil {
			if textStyleCSS := ob.ButtonStyle.TextStyle.ToCSSString(); textStyleCSS != "" {
				styles = append(styles, textStyleCSS)
			}
		}
	}

	// Handle disabled state
	if ob.OnPressed == nil {
		styles = append(styles, "opacity: 0.6")
		styles = append(styles, "cursor: not-allowed")
		attrs["disabled"] = "true"
	}

	// Add hover and focus styles via CSS
	styles = append(styles, "transition: all 0.2s ease")

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Add HTMX event handlers for OnPressed callback
	if ob.OnPressed != nil {
		handlerID := ctx.RegisterHandler(func(ctx *core.Context) Widget {
			ob.OnPressed()
			return nil // Return nil for callbacks that don't return widgets
		})

		attrs["hx-post"] = "/handlers/" + handlerID
		attrs["hx-trigger"] = "click"
	}

	// Add accessibility attributes
	attrs["role"] = "button"
	attrs["tabindex"] = "0"

	if ob.AutoFocus {
		attrs["autofocus"] = "true"
	}

	// Render child content
	content := ""
	if ob.Child != nil {
		content = ob.Child.Render(ctx)
	}

	return htmlRenderer.RenderElement("button", attrs, content, false)
}

// FilledButton represents a filled button widget with full Flutter properties
type FilledButton struct {
	ID               string
	Style            string
	Class            string
	OnPressed        VoidCallback              // Callback when pressed
	OnLongPress      VoidCallback              // Callback when long pressed
	OnHover          ValueChanged[bool]        // Callback when hovered
	OnFocusChange    ValueChanged[bool]        // Callback when focus changes
	ButtonStyle      *ButtonStyle              // Button style
	FocusNode        *FocusNode                // Focus node
	AutoFocus        bool                      // Auto focus
	ClipBehavior     Clip                      // Clip behavior
	StatesController *MaterialStatesController // States controller
	Child            Widget                    // Child widget
}

// Render renders the filled button as HTML
func (fb FilledButton) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(fb.ID, fb.Style, fb.Class+" godin-filled-button")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if fb.Style != "" {
		styles = append(styles, fb.Style)
	}

	// Base button styles
	styles = append(styles, "display: inline-flex")
	styles = append(styles, "align-items: center")
	styles = append(styles, "justify-content: center")
	styles = append(styles, "border: none")
	styles = append(styles, "cursor: pointer")
	styles = append(styles, "text-decoration: none")
	styles = append(styles, "outline: none")
	styles = append(styles, "user-select: none")

	// Default filled button styles
	styles = append(styles, "background-color: #1976d2")
	styles = append(styles, "color: white")
	styles = append(styles, "border-radius: 20px")
	styles = append(styles, "padding: 10px 24px")
	styles = append(styles, "min-height: 40px")
	styles = append(styles, "font-weight: 500")
	styles = append(styles, "font-size: 14px")

	// Apply ButtonStyle if provided
	if fb.ButtonStyle != nil {
		if fb.ButtonStyle.BackgroundColor != nil {
			styles = append(styles, fmt.Sprintf("background-color: %s", fb.ButtonStyle.BackgroundColor.Default))
		}
		if fb.ButtonStyle.ForegroundColor != nil {
			styles = append(styles, fmt.Sprintf("color: %s", fb.ButtonStyle.ForegroundColor.Default))
		}
		if fb.ButtonStyle.Padding != nil {
			styles = append(styles, fmt.Sprintf("padding: %s", fb.ButtonStyle.Padding.Default.ToCSSString()))
		}
		if fb.ButtonStyle.Shape != nil {
			if fb.ButtonStyle.Shape.Default != nil {
				styles = append(styles, fb.ButtonStyle.Shape.Default.ToCSSString())
			}
		}
		if fb.ButtonStyle.TextStyle != nil {
			if textStyleCSS := fb.ButtonStyle.TextStyle.ToCSSString(); textStyleCSS != "" {
				styles = append(styles, textStyleCSS)
			}
		}
	}

	// Handle disabled state
	if fb.OnPressed == nil {
		styles = append(styles, "opacity: 0.6")
		styles = append(styles, "cursor: not-allowed")
		attrs["disabled"] = "true"
	}

	// Add hover and focus styles via CSS
	styles = append(styles, "transition: all 0.2s ease")

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Add HTMX event handlers for OnPressed callback
	if fb.OnPressed != nil {
		handlerID := ctx.RegisterHandler(func(ctx *core.Context) Widget {
			fb.OnPressed()
			return nil // Return nil for callbacks that don't return widgets
		})

		attrs["hx-post"] = "/handlers/" + handlerID
		attrs["hx-trigger"] = "click"
	}

	// Add accessibility attributes
	attrs["role"] = "button"
	attrs["tabindex"] = "0"

	if fb.AutoFocus {
		attrs["autofocus"] = "true"
	}

	// Render child content
	content := ""
	if fb.Child != nil {
		content = fb.Child.Render(ctx)
	}

	return htmlRenderer.RenderElement("button", attrs, content, false)
}

// IconButton represents an icon button widget with full Flutter properties
type IconButton struct {
	ID             string
	Style          string
	Class          string
	OnPressed      VoidCallback        // Callback when pressed
	Icon           Widget              // Icon widget
	IconSize       *float64            // Icon size
	VisualDensity  *VisualDensity      // Visual density
	Padding        *EdgeInsetsGeometry // Padding
	Alignment      AlignmentGeometry   // Alignment
	SplashRadius   *float64            // Splash radius
	Color          Color               // Icon color
	FocusColor     Color               // Focus color
	HoverColor     Color               // Hover color
	HighlightColor Color               // Highlight color
	SplashColor    Color               // Splash color
	DisabledColor  Color               // Disabled color
	MouseCursor    MouseCursor         // Mouse cursor
	FocusNode      *FocusNode          // Focus node
	AutoFocus      bool                // Auto focus
	Tooltip        string              // Tooltip text
	EnableFeedback *bool               // Enable feedback
	Constraints    *BoxConstraints     // Layout constraints
	ButtonStyle    *ButtonStyle        // Button style
	IsSelected     *bool               // Is selected
	SelectedIcon   Widget              // Selected icon
}

// Render renders the icon button as HTML
func (ib IconButton) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(ib.ID, ib.Style, ib.Class+" godin-icon-button")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if ib.Style != "" {
		styles = append(styles, ib.Style)
	}

	// Base button styles
	styles = append(styles, "display: inline-flex")
	styles = append(styles, "align-items: center")
	styles = append(styles, "justify-content: center")
	styles = append(styles, "border: none")
	styles = append(styles, "cursor: pointer")
	styles = append(styles, "text-decoration: none")
	styles = append(styles, "outline: none")
	styles = append(styles, "user-select: none")

	// Default icon button styles
	styles = append(styles, "background-color: transparent")
	styles = append(styles, "border-radius: 50%")
	styles = append(styles, "width: 48px")
	styles = append(styles, "height: 48px")

	// Add padding
	if ib.Padding != nil {
		styles = append(styles, fmt.Sprintf("padding: %s", ib.Padding.ToCSSString()))
	} else {
		styles = append(styles, "padding: 8px")
	}

	// Add icon size
	if ib.IconSize != nil {
		styles = append(styles, fmt.Sprintf("font-size: %.1fpx", *ib.IconSize))
	}

	// Add colors
	if ib.Color != "" {
		styles = append(styles, fmt.Sprintf("color: %s", ib.Color))
	}

	// Handle disabled state
	if ib.OnPressed == nil {
		if ib.DisabledColor != "" {
			styles = append(styles, fmt.Sprintf("color: %s", ib.DisabledColor))
		} else {
			styles = append(styles, "opacity: 0.6")
		}
		styles = append(styles, "cursor: not-allowed")
		attrs["disabled"] = "true"
	}

	// Add hover and focus styles via CSS
	styles = append(styles, "transition: all 0.2s ease")

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Add HTMX event handlers for OnPressed callback
	if ib.OnPressed != nil {
		handlerID := ctx.RegisterHandler(func(ctx *core.Context) Widget {
			ib.OnPressed()
			return nil // Return nil for callbacks that don't return widgets
		})

		attrs["hx-post"] = "/handlers/" + handlerID
		attrs["hx-trigger"] = "click"
	}

	// Add accessibility attributes
	attrs["role"] = "button"
	attrs["tabindex"] = "0"

	if ib.AutoFocus {
		attrs["autofocus"] = "true"
	}

	if ib.Tooltip != "" {
		attrs["title"] = ib.Tooltip
	}

	// Render icon content
	content := ""
	if ib.IsSelected != nil && *ib.IsSelected && ib.SelectedIcon != nil {
		content = ib.SelectedIcon.Render(ctx)
	} else if ib.Icon != nil {
		content = ib.Icon.Render(ctx)
	}

	return htmlRenderer.RenderElement("button", attrs, content, false)
}

// FloatingActionButton represents a floating action button widget with full Flutter properties
type FloatingActionButton struct {
	ID                    string
	Style                 string
	Class                 string
	Child                 Widget                // Child widget
	Tooltip               string                // Tooltip text
	ForegroundColor       Color                 // Foreground color
	BackgroundColor       Color                 // Background color
	FocusColor            Color                 // Focus color
	HoverColor            Color                 // Hover color
	SplashColor           Color                 // Splash color
	HeroTag               interface{}           // Hero tag
	Elevation             *float64              // Elevation
	FocusElevation        *float64              // Focus elevation
	HoverElevation        *float64              // Hover elevation
	HighlightElevation    *float64              // Highlight elevation
	DisabledElevation     *float64              // Disabled elevation
	OnPressed             VoidCallback          // Callback when pressed
	MouseCursor           MouseCursor           // Mouse cursor
	Mini                  bool                  // Is mini FAB
	Shape                 OutlinedBorder        // Shape
	ClipBehavior          Clip                  // Clip behavior
	FocusNode             *FocusNode            // Focus node
	AutoFocus             bool                  // Auto focus
	MaterialTapTargetSize MaterialTapTargetSize // Tap target size
	IsExtended            bool                  // Is extended FAB
	EnableFeedback        *bool                 // Enable feedback
}

// Render renders the floating action button as HTML
func (fab FloatingActionButton) Render(ctx *core.Context) string {
	htmlRenderer := renderer.NewHTMLRenderer()

	attrs := buildAttributes(fab.ID, fab.Style, fab.Class+" godin-floating-action-button")

	// Build inline styles
	var styles []string

	// Add custom style if provided
	if fab.Style != "" {
		styles = append(styles, fab.Style)
	}

	// Base FAB styles
	styles = append(styles, "display: inline-flex")
	styles = append(styles, "align-items: center")
	styles = append(styles, "justify-content: center")
	styles = append(styles, "border: none")
	styles = append(styles, "cursor: pointer")
	styles = append(styles, "text-decoration: none")
	styles = append(styles, "outline: none")
	styles = append(styles, "user-select: none")
	styles = append(styles, "position: fixed")
	styles = append(styles, "bottom: 16px")
	styles = append(styles, "right: 16px")

	// Default FAB styles
	if fab.Mini {
		styles = append(styles, "width: 40px")
		styles = append(styles, "height: 40px")
	} else if fab.IsExtended {
		styles = append(styles, "height: 56px")
		styles = append(styles, "padding: 0 16px")
		styles = append(styles, "border-radius: 28px")
	} else {
		styles = append(styles, "width: 56px")
		styles = append(styles, "height: 56px")
		styles = append(styles, "border-radius: 50%")
	}

	// Default colors
	if fab.BackgroundColor != "" {
		styles = append(styles, fmt.Sprintf("background-color: %s", fab.BackgroundColor))
	} else {
		styles = append(styles, "background-color: #1976d2")
	}

	if fab.ForegroundColor != "" {
		styles = append(styles, fmt.Sprintf("color: %s", fab.ForegroundColor))
	} else {
		styles = append(styles, "color: white")
	}

	// Add elevation (box-shadow)
	elevation := 6.0
	if fab.Elevation != nil {
		elevation = *fab.Elevation
	}
	if elevation > 0 {
		styles = append(styles, fmt.Sprintf("box-shadow: 0 %.0fpx %.0fpx rgba(0,0,0,0.2)", elevation, elevation*2))
	}

	// Handle disabled state
	if fab.OnPressed == nil {
		styles = append(styles, "opacity: 0.6")
		styles = append(styles, "cursor: not-allowed")
		attrs["disabled"] = "true"
	}

	// Add hover and focus styles via CSS
	styles = append(styles, "transition: all 0.2s ease")

	// Combine all styles
	if len(styles) > 0 {
		attrs["style"] = strings.Join(styles, "; ")
	}

	// Add HTMX event handlers for OnPressed callback
	if fab.OnPressed != nil {
		handlerID := ctx.RegisterHandler(func(ctx *core.Context) Widget {
			fab.OnPressed()
			return nil // Return nil for callbacks that don't return widgets
		})

		attrs["hx-post"] = "/handlers/" + handlerID
		attrs["hx-trigger"] = "click"
	}

	// Add accessibility attributes
	attrs["role"] = "button"
	attrs["tabindex"] = "0"

	if fab.AutoFocus {
		attrs["autofocus"] = "true"
	}

	if fab.Tooltip != "" {
		attrs["title"] = fab.Tooltip
	}

	// Render child content
	content := ""
	if fab.Child != nil {
		content = fab.Child.Render(ctx)
	}

	return htmlRenderer.RenderElement("button", attrs, content, false)
}
