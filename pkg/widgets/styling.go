package widgets

import (
	"fmt"
	"strings"
	"time"
)

// Color represents a color value
type Color string

// Common color constants
const (
	ColorTransparent Color = "transparent"
	ColorBlack       Color = "#000000"
	ColorWhite       Color = "#ffffff"
	ColorRed         Color = "#f44336"
	ColorGreen       Color = "#4caf50"
	ColorBlue        Color = "#2196f3"
	ColorYellow      Color = "#ffeb3b"
	ColorOrange      Color = "#ff9800"
	ColorPurple      Color = "#9c27b0"
	ColorGrey        Color = "#9e9e9e"
)

// BoxConstraints represents layout constraints
// type BoxConstraints struct {
// 	MinWidth  *float64
// 	MaxWidth  *float64
// 	MinHeight *float64
// 	MaxHeight *float64
// }

// Matrix4 represents a 4x4 transformation matrix (simplified)
type Matrix4 struct {
	Values [16]float64
}

// BoxDecoration represents container decoration
type BoxDecoration struct {
	Color               Color
	Image               *DecorationImage
	Border              *BoxBorder
	BorderRadius        *BorderRadius
	BoxShadow           []BoxShadow
	Gradient            *Gradient
	BackgroundBlendMode BlendMode
	Shape               BoxShape
}

// ToCSSString converts BoxDecoration to CSS styles
func (bd BoxDecoration) ToCSSString() string {
	var styles []string

	if bd.Color != "" {
		styles = append(styles, fmt.Sprintf("background-color: %s", bd.Color))
	}

	if bd.BorderRadius != nil {
		styles = append(styles, bd.BorderRadius.ToCSSString())
	}

	if bd.Border != nil {
		styles = append(styles, bd.Border.ToCSSString())
	}

	if len(bd.BoxShadow) > 0 {
		var shadows []string
		for _, shadow := range bd.BoxShadow {
			shadows = append(shadows, shadow.ToCSSString())
		}
		styles = append(styles, fmt.Sprintf("box-shadow: %s", strings.Join(shadows, ", ")))
	}

	if bd.Shape == BoxShapeCircle {
		styles = append(styles, "border-radius: 50%")
	}

	return strings.Join(styles, "; ")
}

// BoxShadow represents a box shadow
type BoxShadow struct {
	Color        Color
	Offset       Offset
	BlurRadius   float64
	SpreadRadius float64
	BlurStyle    BlurStyle
}

// ToCSSString converts BoxShadow to CSS box-shadow value
func (bs BoxShadow) ToCSSString() string {
	return fmt.Sprintf("%.1fpx %.1fpx %.1fpx %.1fpx %s",
		bs.Offset.DX, bs.Offset.DY, bs.BlurRadius, bs.SpreadRadius, bs.Color)
}

// Offset represents a 2D offset
type Offset struct {
	DX float64
	DY float64
}

// BlurStyle enum
type BlurStyle string

const (
	BlurStyleNormal BlurStyle = "normal"
	BlurStyleSolid  BlurStyle = "solid"
	BlurStyleOuter  BlurStyle = "outer"
	BlurStyleInner  BlurStyle = "inner"
)

// BoxShape enum
type BoxShape string

const (
	BoxShapeRectangle BoxShape = "rectangle"
	BoxShapeCircle    BoxShape = "circle"
)

// BorderRadius represents border radius values
type BorderRadius struct {
	TopLeft     Radius
	TopRight    Radius
	BottomLeft  Radius
	BottomRight Radius
}

// BorderRadiusAll creates BorderRadius with all corners equal
func BorderRadiusAll(radius Radius) *BorderRadius {
	return &BorderRadius{
		TopLeft:     radius,
		TopRight:    radius,
		BottomLeft:  radius,
		BottomRight: radius,
	}
}

// BorderRadiusCircular creates circular BorderRadius
func BorderRadiusCircular(radius float64) *BorderRadius {
	r := Radius{X: radius, Y: radius}
	return BorderRadiusAll(r)
}

// ToCSSString converts BorderRadius to CSS border-radius
func (br BorderRadius) ToCSSString() string {
	return fmt.Sprintf("border-radius: %.1fpx %.1fpx %.1fpx %.1fpx",
		br.TopLeft.X, br.TopRight.X, br.BottomRight.X, br.BottomLeft.X)
}

// Radius represents a radius value
type Radius struct {
	X float64
	Y float64
}

// BoxBorder represents border properties
type BoxBorder struct {
	Top    BorderSide
	Right  BorderSide
	Bottom BorderSide
	Left   BorderSide
}

// BorderAll creates BoxBorder with all sides equal
func BorderAll(side BorderSide) *BoxBorder {
	return &BoxBorder{
		Top:    side,
		Right:  side,
		Bottom: side,
		Left:   side,
	}
}

// ToCSSString converts BoxBorder to CSS border
func (bb BoxBorder) ToCSSString() string {
	var styles []string

	if bb.Top.Width > 0 {
		styles = append(styles, fmt.Sprintf("border-top: %.1fpx %s %s", bb.Top.Width, bb.Top.Style, bb.Top.Color))
	}
	if bb.Right.Width > 0 {
		styles = append(styles, fmt.Sprintf("border-right: %.1fpx %s %s", bb.Right.Width, bb.Right.Style, bb.Right.Color))
	}
	if bb.Bottom.Width > 0 {
		styles = append(styles, fmt.Sprintf("border-bottom: %.1fpx %s %s", bb.Bottom.Width, bb.Bottom.Style, bb.Bottom.Color))
	}
	if bb.Left.Width > 0 {
		styles = append(styles, fmt.Sprintf("border-left: %.1fpx %s %s", bb.Left.Width, bb.Left.Style, bb.Left.Color))
	}

	return strings.Join(styles, "; ")
}

// BorderSide represents a single border side
type BorderSide struct {
	Color Color
	Width float64
	Style BorderStyle
}

// BorderStyle enum
type BorderStyle string

const (
	BorderStyleNone   BorderStyle = "none"
	BorderStyleSolid  BorderStyle = "solid"
	BorderStyleDashed BorderStyle = "dashed"
	BorderStyleDotted BorderStyle = "dotted"
)

// DecorationImage represents a background image
type DecorationImage struct {
	Image              ImageProvider
	ColorFilter        *ColorFilter
	Fit                BoxFit
	Alignment          AlignmentGeometry
	CenterSlice        *Rect
	Repeat             ImageRepeat
	MatchTextDirection bool
	Scale              float64
	Opacity            float64
	FilterQuality      FilterQuality
	InvertColors       bool
	IsAntiAlias        bool
}

// ImageProvider interface for image sources
type ImageProvider interface {
	GetImageURL() string
}

// NetworkImage implements ImageProvider for network images
type NetworkImage struct {
	URL string
}

func (ni NetworkImage) GetImageURL() string {
	return ni.URL
}

// AssetImage implements ImageProvider for asset images
type AssetImage struct {
	AssetPath string
}

func (ai AssetImage) GetImageURL() string {
	return ai.AssetPath
}

// ToCSSString converts DecorationImage to CSS background styles
func (di DecorationImage) ToCSSString() string {
	var styles []string

	if di.Image != nil {
		imageURL := di.Image.GetImageURL()
		styles = append(styles, fmt.Sprintf("background-image: url('%s')", imageURL))
	}

	// Add background fit
	switch di.Fit {
	case BoxFitFill:
		styles = append(styles, "background-size: 100% 100%")
	case BoxFitContain:
		styles = append(styles, "background-size: contain")
	case BoxFitCover:
		styles = append(styles, "background-size: cover")
	case BoxFitFitWidth:
		styles = append(styles, "background-size: 100% auto")
	case BoxFitFitHeight:
		styles = append(styles, "background-size: auto 100%")
	case BoxFitNone:
		styles = append(styles, "background-size: auto")
	case BoxFitScaleDown:
		styles = append(styles, "background-size: contain")
	}

	// Add background repeat
	switch di.Repeat {
	case ImageRepeatRepeat:
		styles = append(styles, "background-repeat: repeat")
	case ImageRepeatRepeatX:
		styles = append(styles, "background-repeat: repeat-x")
	case ImageRepeatRepeatY:
		styles = append(styles, "background-repeat: repeat-y")
	case ImageRepeatNoRepeat:
		styles = append(styles, "background-repeat: no-repeat")
	}

	// Add background position (simplified alignment)
	if di.Alignment != "" {
		alignmentStr := string(di.Alignment)
		if alignmentStr == "center" {
			styles = append(styles, "background-position: center")
		} else if alignmentStr == "topLeft" {
			styles = append(styles, "background-position: top left")
		} else if alignmentStr == "topRight" {
			styles = append(styles, "background-position: top right")
		} else if alignmentStr == "bottomLeft" {
			styles = append(styles, "background-position: bottom left")
		} else if alignmentStr == "bottomRight" {
			styles = append(styles, "background-position: bottom right")
		}
	}

	// Add opacity
	if di.Opacity > 0 && di.Opacity < 1 {
		styles = append(styles, fmt.Sprintf("opacity: %.2f", di.Opacity))
	}

	return strings.Join(styles, "; ")
}

// BoxFit enum for image fitting
type BoxFit string

const (
	BoxFitFill      BoxFit = "fill"
	BoxFitContain   BoxFit = "contain"
	BoxFitCover     BoxFit = "cover"
	BoxFitFitWidth  BoxFit = "fitWidth"
	BoxFitFitHeight BoxFit = "fitHeight"
	BoxFitNone      BoxFit = "none"
	BoxFitScaleDown BoxFit = "scaleDown"
)

// ImageRepeat enum
type ImageRepeat string

const (
	ImageRepeatRepeat   ImageRepeat = "repeat"
	ImageRepeatRepeatX  ImageRepeat = "repeat-x"
	ImageRepeatRepeatY  ImageRepeat = "repeat-y"
	ImageRepeatNoRepeat ImageRepeat = "no-repeat"
)

// FilterQuality enum
type FilterQuality string

const (
	FilterQualityNone   FilterQuality = "none"
	FilterQualityLow    FilterQuality = "low"
	FilterQualityMedium FilterQuality = "medium"
	FilterQualityHigh   FilterQuality = "high"
)

// ColorFilter represents color filtering
type ColorFilter struct {
	Color     Color
	BlendMode BlendMode
}

// BlendMode enum
type BlendMode string

const (
	BlendModeNormal   BlendMode = "normal"
	BlendModeMultiply BlendMode = "multiply"
	BlendModeScreen   BlendMode = "screen"
	BlendModeOverlay  BlendMode = "overlay"
)

// Gradient interface for gradients
type Gradient interface {
	ToCSSString() string
}

// LinearGradient implements Gradient
type LinearGradient struct {
	Begin  AlignmentGeometry
	End    AlignmentGeometry
	Colors []Color
	Stops  []float64
}

func (lg LinearGradient) ToCSSString() string {
	// Simplified linear gradient implementation
	if len(lg.Colors) >= 2 {
		return fmt.Sprintf("linear-gradient(to right, %s, %s)", lg.Colors[0], lg.Colors[1])
	}
	return ""
}

// Rect represents a rectangle
type Rect struct {
	Left   float64
	Top    float64
	Right  float64
	Bottom float64
}

// Clip enum
type Clip string

const (
	ClipNone                   Clip = "none"
	ClipHardEdge               Clip = "hardEdge"
	ClipAntiAlias              Clip = "antiAlias"
	ClipAntiAliasWithSaveLayer Clip = "antiAliasWithSaveLayer"
)

// Duration represents a time duration
type Duration time.Duration

// Curve represents animation curves
type Curve string

const (
	CurveLinear      Curve = "linear"
	CurveEase        Curve = "ease"
	CurveEaseIn      Curve = "ease-in"
	CurveEaseOut     Curve = "ease-out"
	CurveEaseInOut   Curve = "ease-in-out"
	CurveBounceIn    Curve = "bounce-in"
	CurveBounceOut   Curve = "bounce-out"
	CurveBounceInOut Curve = "bounce-in-out"
)

// TextStyle represents text styling properties
type TextStyle struct {
	Color               Color
	FontSize            *float64
	FontWeight          FontWeight
	FontStyle           FontStyle
	LetterSpacing       *float64
	WordSpacing         *float64
	TextBaseline        TextBaseline
	Height              *float64
	Locale              *Locale
	Foreground          *Paint
	Background          *Paint
	Shadows             []Shadow
	FontFeatures        []FontFeature
	Decoration          TextDecoration
	DecorationColor     Color
	DecorationStyle     TextDecorationStyle
	DecorationThickness *float64
	FontFamily          string
	FontFamilyFallback  []string
	Package             string
}

// ToCSSString converts TextStyle to CSS styles
func (ts TextStyle) ToCSSString() string {
	var styles []string

	if ts.Color != "" {
		styles = append(styles, fmt.Sprintf("color: %s", ts.Color))
	}

	if ts.FontSize != nil {
		styles = append(styles, fmt.Sprintf("font-size: %.1fpx", *ts.FontSize))
	}

	if ts.FontWeight != "" {
		styles = append(styles, fmt.Sprintf("font-weight: %s", ts.FontWeight))
	}

	if ts.FontStyle != "" {
		styles = append(styles, fmt.Sprintf("font-style: %s", ts.FontStyle))
	}

	if ts.FontFamily != "" {
		styles = append(styles, fmt.Sprintf("font-family: %s", ts.FontFamily))
	}

	if ts.LetterSpacing != nil {
		styles = append(styles, fmt.Sprintf("letter-spacing: %.1fpx", *ts.LetterSpacing))
	}

	if ts.WordSpacing != nil {
		styles = append(styles, fmt.Sprintf("word-spacing: %.1fpx", *ts.WordSpacing))
	}

	if ts.Height != nil {
		styles = append(styles, fmt.Sprintf("line-height: %.2f", *ts.Height))
	}

	if ts.Decoration != TextDecorationNone {
		styles = append(styles, fmt.Sprintf("text-decoration: %s", ts.Decoration))
	}

	if ts.DecorationColor != "" {
		styles = append(styles, fmt.Sprintf("text-decoration-color: %s", ts.DecorationColor))
	}

	if ts.DecorationStyle != "" {
		styles = append(styles, fmt.Sprintf("text-decoration-style: %s", ts.DecorationStyle))
	}

	return strings.Join(styles, "; ")
}

// FontWeight enum
type FontWeight string

const (
	FontWeightW100   FontWeight = "100"
	FontWeightW200   FontWeight = "200"
	FontWeightW300   FontWeight = "300"
	FontWeightW400   FontWeight = "400"
	FontWeightW500   FontWeight = "500"
	FontWeightW600   FontWeight = "600"
	FontWeightW700   FontWeight = "700"
	FontWeightW800   FontWeight = "800"
	FontWeightW900   FontWeight = "900"
	FontWeightNormal FontWeight = "normal"
	FontWeightBold   FontWeight = "bold"
)

// FontStyle enum
type FontStyle string

const (
	FontStyleNormal FontStyle = "normal"
	FontStyleItalic FontStyle = "italic"
)

// TextBaseline enum
type TextBaseline string

const (
	TextBaselineAlphabetic  TextBaseline = "alphabetic"
	TextBaselineIdeographic TextBaseline = "ideographic"
)

// TextDecoration enum
type TextDecoration string

const (
	TextDecorationNone        TextDecoration = "none"
	TextDecorationUnderline   TextDecoration = "underline"
	TextDecorationOverline    TextDecoration = "overline"
	TextDecorationLineThrough TextDecoration = "line-through"
)

// TextDecorationStyle enum
type TextDecorationStyle string

const (
	TextDecorationStyleSolid  TextDecorationStyle = "solid"
	TextDecorationStyleDouble TextDecorationStyle = "double"
	TextDecorationStyleDotted TextDecorationStyle = "dotted"
	TextDecorationStyleDashed TextDecorationStyle = "dashed"
	TextDecorationStyleWavy   TextDecorationStyle = "wavy"
)

// Paint represents paint properties
type Paint struct {
	Color Color
}

// Shadow represents text shadow
type Shadow struct {
	Color      Color
	Offset     Offset
	BlurRadius float64
}

// FontFeature represents font features
type FontFeature struct {
	Feature string
	Value   int
}

// Locale represents locale information
type Locale struct {
	LanguageCode string
	CountryCode  string
}

// StrutStyle represents strut styling
type StrutStyle struct {
	FontFamily       string
	FontSize         *float64
	Height           *float64
	Leading          *float64
	FontWeight       FontWeight
	FontStyle        FontStyle
	ForceStrutHeight bool
}

// TextWidthBasis enum
type TextWidthBasis string

const (
	TextWidthBasisParent      TextWidthBasis = "parent"
	TextWidthBasisLongestLine TextWidthBasis = "longestLine"
)

// TextHeightBehavior represents text height behavior
type TextHeightBehavior struct {
	ApplyHeightToFirstAscent bool
	ApplyHeightToLastDescent bool
}

// ButtonStyle represents button styling properties
type ButtonStyle struct {
	TextStyle         *TextStyle                                 // Text style
	BackgroundColor   *MaterialStateProperty[Color]              // Background color
	ForegroundColor   *MaterialStateProperty[Color]              // Foreground color
	OverlayColor      *MaterialStateProperty[Color]              // Overlay color
	ShadowColor       *MaterialStateProperty[Color]              // Shadow color
	SurfaceTintColor  *MaterialStateProperty[Color]              // Surface tint color
	Elevation         *MaterialStateProperty[float64]            // Elevation
	Padding           *MaterialStateProperty[EdgeInsetsGeometry] // Padding
	MinimumSize       *MaterialStateProperty[Size]               // Minimum size
	FixedSize         *MaterialStateProperty[Size]               // Fixed size
	MaximumSize       *MaterialStateProperty[Size]               // Maximum size
	Side              *MaterialStateProperty[BorderSide]         // Border side
	Shape             *MaterialStateProperty[OutlinedBorder]     // Shape
	MouseCursor       *MaterialStateProperty[MouseCursor]        // Mouse cursor
	VisualDensity     *VisualDensity                             // Visual density
	TapTargetSize     MaterialTapTargetSize                      // Tap target size
	AnimationDuration *Duration                                  // Animation duration
	EnableFeedback    *bool                                      // Enable feedback
	Alignment         AlignmentGeometry                          // Alignment
	SplashFactory     InteractiveInkFeatureFactory               // Splash factory
}

// MaterialStateProperty represents a property that can have different values for different states
type MaterialStateProperty[T any] struct {
	Default  T
	Hovered  *T
	Focused  *T
	Pressed  *T
	Dragged  *T
	Selected *T
	Scrolled *T
	Disabled *T
	Error    *T
}

// Size represents a size with width and height
type Size struct {
	Width  float64
	Height float64
}

// OutlinedBorder interface for outlined borders
type OutlinedBorder interface {
	ToCSSString() string
}

// RoundedRectangleBorder implements OutlinedBorder
type RoundedRectangleBorder struct {
	BorderRadius *BorderRadius
	Side         BorderSide
}

func (rrb RoundedRectangleBorder) ToCSSString() string {
	var styles []string

	if rrb.BorderRadius != nil {
		styles = append(styles, rrb.BorderRadius.ToCSSString())
	}

	if rrb.Side.Width > 0 {
		styles = append(styles, fmt.Sprintf("border: %.1fpx %s %s", rrb.Side.Width, rrb.Side.Style, rrb.Side.Color))
	}

	return strings.Join(styles, "; ")
}

// CircleBorder implements OutlinedBorder
type CircleBorder struct {
	Side BorderSide
}

func (cb CircleBorder) ToCSSString() string {
	var styles []string

	styles = append(styles, "border-radius: 50%")

	if cb.Side.Width > 0 {
		styles = append(styles, fmt.Sprintf("border: %.1fpx %s %s", cb.Side.Width, cb.Side.Style, cb.Side.Color))
	}

	return strings.Join(styles, "; ")
}

// MouseCursor enum
type MouseCursor string

const (
	MouseCursorBasic     MouseCursor = "default"
	MouseCursorClick     MouseCursor = "pointer"
	MouseCursorForbidden MouseCursor = "not-allowed"
	MouseCursorWait      MouseCursor = "wait"
	MouseCursorProgress  MouseCursor = "progress"
	MouseCursorPrecise   MouseCursor = "crosshair"
	MouseCursorText      MouseCursor = "text"
	MouseCursorHelp      MouseCursor = "help"
	MouseCursorMove      MouseCursor = "move"
	MouseCursorNone      MouseCursor = "none"
	MouseCursorGrab      MouseCursor = "grab"
	MouseCursorGrabbing  MouseCursor = "grabbing"
)

// VisualDensity represents visual density
type VisualDensity struct {
	Horizontal float64
	Vertical   float64
}

// MaterialTapTargetSize enum
type MaterialTapTargetSize string

const (
	MaterialTapTargetSizePadded     MaterialTapTargetSize = "padded"
	MaterialTapTargetSizeShrinkWrap MaterialTapTargetSize = "shrinkWrap"
)

// InteractiveInkFeatureFactory interface
type InteractiveInkFeatureFactory interface {
	Create() string
}

// VoidCallback represents a callback function with no parameters
type VoidCallback func()

// ValueChanged represents a callback function with a value parameter
type ValueChanged[T any] func(T)

// FormFieldSetter represents a callback function for saving form field values
type FormFieldSetter[T any] func(T)

// FormFieldValidator represents a validation function for form fields
type FormFieldValidator[T any] func(T) *string

// AutovalidateMode enum for form field validation
type AutovalidateMode string

const (
	AutovalidateModeDisabled          AutovalidateMode = "disabled"
	AutovalidateModeAlways            AutovalidateMode = "always"
	AutovalidateModeOnUserInteraction AutovalidateMode = "onUserInteraction"
)

// GestureTapCallback represents a tap gesture callback
type GestureTapCallback func()

// ImageErrorListener represents an image error callback
type ImageErrorListener func(error)

// GestureLongPressCallback represents a long press gesture callback
type GestureLongPressCallback func()

// ShapeBorder interface for shape borders
type ShapeBorder interface {
	ToCSSString() string
}

// ListTileStyle enum
type ListTileStyle string

const (
	ListTileStyleList   ListTileStyle = "list"
	ListTileStyleDrawer ListTileStyle = "drawer"
)

// Axis enum for scroll direction
type Axis string

const (
	AxisHorizontal Axis = "horizontal"
	AxisVertical   Axis = "vertical"
)

// ScrollViewKeyboardDismissBehavior enum
type ScrollViewKeyboardDismissBehavior string

const (
	ScrollViewKeyboardDismissBehaviorManual     ScrollViewKeyboardDismissBehavior = "manual"
	ScrollViewKeyboardDismissBehaviorOnDrag     ScrollViewKeyboardDismissBehavior = "onDrag"
	ScrollViewKeyboardDismissBehaviorOnDragDown ScrollViewKeyboardDismissBehavior = "onDragDown"
)

// ScrollPhysicsType enum for scroll physics types
type ScrollPhysicsType string

const (
	ScrollPhysicsAlwaysScrollable   ScrollPhysicsType = "alwaysScrollable"
	ScrollPhysicsNeverScrollable    ScrollPhysicsType = "neverScrollable"
	ScrollPhysicsBouncingScrollable ScrollPhysicsType = "bouncingScrollable"
	ScrollPhysicsClampingScrollable ScrollPhysicsType = "clampingScrollable"
)

// SliverGridDelegate interface for grid layout delegates
type SliverGridDelegate interface {
	GetCrossAxisCount() int
	GetMainAxisSpacing() float64
	GetCrossAxisSpacing() float64
	GetChildAspectRatio() float64
}

// SliverGridDelegateWithFixedCrossAxisCount implements SliverGridDelegate
type SliverGridDelegateWithFixedCrossAxisCount struct {
	CrossAxisCount   int      // Number of columns
	MainAxisSpacing  float64  // Spacing between rows
	CrossAxisSpacing float64  // Spacing between columns
	ChildAspectRatio float64  // Aspect ratio of each child
	MainAxisExtent   *float64 // Fixed main axis extent
}

func (d SliverGridDelegateWithFixedCrossAxisCount) GetCrossAxisCount() int {
	return d.CrossAxisCount
}

func (d SliverGridDelegateWithFixedCrossAxisCount) GetMainAxisSpacing() float64 {
	return d.MainAxisSpacing
}

func (d SliverGridDelegateWithFixedCrossAxisCount) GetCrossAxisSpacing() float64 {
	return d.CrossAxisSpacing
}

func (d SliverGridDelegateWithFixedCrossAxisCount) GetChildAspectRatio() float64 {
	if d.ChildAspectRatio <= 0 {
		return 1.0 // Default aspect ratio
	}
	return d.ChildAspectRatio
}

// SliverGridDelegateWithMaxCrossAxisExtent implements SliverGridDelegate
type SliverGridDelegateWithMaxCrossAxisExtent struct {
	MaxCrossAxisExtent float64  // Maximum cross axis extent
	MainAxisSpacing    float64  // Spacing between rows
	CrossAxisSpacing   float64  // Spacing between columns
	ChildAspectRatio   float64  // Aspect ratio of each child
	MainAxisExtent     *float64 // Fixed main axis extent
}

func (d SliverGridDelegateWithMaxCrossAxisExtent) GetCrossAxisCount() int {
	// This would need to be calculated based on available width
	// For now, return a default value
	return 2
}

func (d SliverGridDelegateWithMaxCrossAxisExtent) GetMainAxisSpacing() float64 {
	return d.MainAxisSpacing
}

func (d SliverGridDelegateWithMaxCrossAxisExtent) GetCrossAxisSpacing() float64 {
	return d.CrossAxisSpacing
}

func (d SliverGridDelegateWithMaxCrossAxisExtent) GetChildAspectRatio() float64 {
	if d.ChildAspectRatio <= 0 {
		return 1.0 // Default aspect ratio
	}
	return d.ChildAspectRatio
}

// TextField-related types and enums

// TextEditingController represents a text editing controller
type TextEditingController struct {
	Text      string
	Selection TextSelection
}

// TextSelection represents text selection
type TextSelection struct {
	BaseOffset   int
	ExtentOffset int
}

// InputDecoration represents input field decoration
type InputDecoration struct {
	Icon                   Widget
	IconColor              Color
	Label                  Widget
	LabelText              string
	LabelStyle             *TextStyle
	FloatingLabelStyle     *TextStyle
	HelperText             string
	HelperStyle            *TextStyle
	HelperMaxLines         *int
	HintText               string
	HintStyle              *TextStyle
	HintTextDirection      TextDirection
	HintMaxLines           *int
	ErrorText              string
	ErrorStyle             *TextStyle
	ErrorMaxLines          *int
	FloatingLabelBehavior  FloatingLabelBehavior
	FloatingLabelAlignment FloatingLabelAlignment
	IsCollapsed            bool
	IsDense                *bool
	ContentPadding         *EdgeInsetsGeometry
	PrefixIcon             Widget
	PrefixIconConstraints  *BoxConstraints
	Prefix                 Widget
	PrefixText             string
	PrefixStyle            *TextStyle
	PrefixIconColor        Color
	SuffixIcon             Widget
	Suffix                 Widget
	SuffixText             string
	SuffixStyle            *TextStyle
	SuffixIconColor        Color
	SuffixIconConstraints  *BoxConstraints
	Counter                Widget
	CounterText            string
	CounterStyle           *TextStyle
	Filled                 *bool
	FillColor              Color
	FocusColor             Color
	HoverColor             Color
	ErrorBorder            InputBorder
	FocusedBorder          InputBorder
	FocusedErrorBorder     InputBorder
	DisabledBorder         InputBorder
	EnabledBorder          InputBorder
	Border                 InputBorder
	Enabled                bool
	Semantics              string
	AlignLabelWithHint     bool
	Constraints            *BoxConstraints
}

// TextInputType enum
type TextInputType string

const (
	TextInputTypeText            TextInputType = "text"
	TextInputTypeMultiline       TextInputType = "multiline"
	TextInputTypeNumber          TextInputType = "number"
	TextInputTypePhone           TextInputType = "tel"
	TextInputTypeDatetime        TextInputType = "datetime-local"
	TextInputTypeEmailAddress    TextInputType = "email"
	TextInputTypeURL             TextInputType = "url"
	TextInputTypeVisiblePassword TextInputType = "password"
	TextInputTypeName            TextInputType = "text"
	TextInputTypeStreetAddress   TextInputType = "text"
	TextInputTypeNone            TextInputType = "text"
)

// TextInputAction enum
type TextInputAction string

const (
	TextInputActionNone           TextInputAction = "none"
	TextInputActionUnspecified    TextInputAction = "unspecified"
	TextInputActionDone           TextInputAction = "done"
	TextInputActionGo             TextInputAction = "go"
	TextInputActionSearch         TextInputAction = "search"
	TextInputActionSend           TextInputAction = "send"
	TextInputActionNext           TextInputAction = "next"
	TextInputActionPrevious       TextInputAction = "previous"
	TextInputActionContinueAction TextInputAction = "continue"
	TextInputActionJoin           TextInputAction = "join"
	TextInputActionRoute          TextInputAction = "route"
	TextInputActionEmergencyCall  TextInputAction = "emergencyCall"
	TextInputActionNewline        TextInputAction = "newline"
)

// TextCapitalization enum
type TextCapitalization string

const (
	TextCapitalizationNone       TextCapitalization = "none"
	TextCapitalizationWords      TextCapitalization = "words"
	TextCapitalizationSentences  TextCapitalization = "sentences"
	TextCapitalizationCharacters TextCapitalization = "characters"
)

// TextAlignVertical enum
type TextAlignVertical string

const (
	TextAlignVerticalTop    TextAlignVertical = "top"
	TextAlignVerticalCenter TextAlignVertical = "center"
	TextAlignVerticalBottom TextAlignVertical = "bottom"
)

// FloatingLabelBehavior enum
type FloatingLabelBehavior string

const (
	FloatingLabelBehaviorNever  FloatingLabelBehavior = "never"
	FloatingLabelBehaviorAuto   FloatingLabelBehavior = "auto"
	FloatingLabelBehaviorAlways FloatingLabelBehavior = "always"
)

// FloatingLabelAlignment enum
type FloatingLabelAlignment string

const (
	FloatingLabelAlignmentStart  FloatingLabelAlignment = "start"
	FloatingLabelAlignmentCenter FloatingLabelAlignment = "center"
)

// InputBorder interface
type InputBorder interface {
	ToCSSString() string
}

// OutlineInputBorder implements InputBorder
type OutlineInputBorder struct {
	BorderSide   BorderSide
	BorderRadius *BorderRadius
	GapPadding   float64
}

func (oib OutlineInputBorder) ToCSSString() string {
	var styles []string

	if oib.BorderSide.Width > 0 {
		styles = append(styles, fmt.Sprintf("border: %.1fpx %s %s", oib.BorderSide.Width, oib.BorderSide.Style, oib.BorderSide.Color))
	}

	if oib.BorderRadius != nil {
		styles = append(styles, oib.BorderRadius.ToCSSString())
	}

	return strings.Join(styles, "; ")
}

// UnderlineInputBorder implements InputBorder
type UnderlineInputBorder struct {
	BorderSide   BorderSide
	BorderRadius *BorderRadius
}

func (uib UnderlineInputBorder) ToCSSString() string {
	var styles []string

	if uib.BorderSide.Width > 0 {
		styles = append(styles, fmt.Sprintf("border-bottom: %.1fpx %s %s", uib.BorderSide.Width, uib.BorderSide.Style, uib.BorderSide.Color))
	}

	if uib.BorderRadius != nil {
		styles = append(styles, uib.BorderRadius.ToCSSString())
	}

	return strings.Join(styles, "; ")
}

// Additional TextField-related types

// ToolbarOptions represents toolbar options
type ToolbarOptions struct {
	Copy      bool
	Cut       bool
	Paste     bool
	SelectAll bool
}

// SmartDashesType enum
type SmartDashesType string

const (
	SmartDashesTypeDisabled SmartDashesType = "disabled"
	SmartDashesTypeEnabled  SmartDashesType = "enabled"
)

// SmartQuotesType enum
type SmartQuotesType string

const (
	SmartQuotesTypeDisabled SmartQuotesType = "disabled"
	SmartQuotesTypeEnabled  SmartQuotesType = "enabled"
)

// MaxLengthEnforcement enum
type MaxLengthEnforcement string

const (
	MaxLengthEnforcementNone                         MaxLengthEnforcement = "none"
	MaxLengthEnforcementEnforced                     MaxLengthEnforcement = "enforced"
	MaxLengthEnforcementTruncateAfterCompositionEnds MaxLengthEnforcement = "truncateAfterCompositionEnds"
)

// TextInputFormatter interface
type TextInputFormatter interface {
	FormatEditUpdate(oldValue, newValue TextEditingValue) TextEditingValue
}

// TextEditingValue represents text editing value
type TextEditingValue struct {
	Text      string
	Selection TextSelection
	Composing TextRange
}

// TextRange represents text range
type TextRange struct {
	Start int
	End   int
}

// BoxHeightStyle enum
type BoxHeightStyle string

const (
	BoxHeightStyleTight                    BoxHeightStyle = "tight"
	BoxHeightStyleMax                      BoxHeightStyle = "max"
	BoxHeightStyleIncludeLineSpacingMiddle BoxHeightStyle = "includeLineSpacingMiddle"
	BoxHeightStyleIncludeLineSpacingTop    BoxHeightStyle = "includeLineSpacingTop"
	BoxHeightStyleIncludeLineSpacingBottom BoxHeightStyle = "includeLineSpacingBottom"
	BoxHeightStyleStrut                    BoxHeightStyle = "strut"
)

// BoxWidthStyle enum
type BoxWidthStyle string

const (
	BoxWidthStyleTight BoxWidthStyle = "tight"
	BoxWidthStyleMax   BoxWidthStyle = "max"
)

// Brightness enum
type Brightness string

const (
	BrightnessLight Brightness = "light"
	BrightnessDark  Brightness = "dark"
)

// DragStartBehavior enum
type DragStartBehavior string

const (
	DragStartBehaviorStart DragStartBehavior = "start"
	DragStartBehaviorDown  DragStartBehavior = "down"
)

// TextSelectionControls interface
type TextSelectionControls interface {
	BuildToolbar() Widget
}

// ScrollController represents scroll controller
type ScrollController struct {
	InitialScrollOffset float64
	KeepScrollOffset    bool
}

// ScrollPhysics interface
type ScrollPhysics interface {
	ApplyPhysicsToUserOffset(offset float64) float64
}

// AutoFillHint enum
type AutoFillHint string

const (
	AutoFillHintEmail              AutoFillHint = "email"
	AutoFillHintName               AutoFillHint = "name"
	AutoFillHintNamePrefix         AutoFillHint = "namePrefix"
	AutoFillHintNameSuffix         AutoFillHint = "nameSuffix"
	AutoFillHintGivenName          AutoFillHint = "givenName"
	AutoFillHintMiddleName         AutoFillHint = "middleName"
	AutoFillHintFamilyName         AutoFillHint = "familyName"
	AutoFillHintUsername           AutoFillHint = "username"
	AutoFillHintPassword           AutoFillHint = "password"
	AutoFillHintNewPassword        AutoFillHint = "newPassword"
	AutoFillHintOneTimeCode        AutoFillHint = "oneTimeCode"
	AutoFillHintTelephoneNumber    AutoFillHint = "telephoneNumber"
	AutoFillHintStreetAddressLine1 AutoFillHint = "streetAddressLine1"
	AutoFillHintStreetAddressLine2 AutoFillHint = "streetAddressLine2"
	AutoFillHintAddressCity        AutoFillHint = "addressCity"
	AutoFillHintAddressState       AutoFillHint = "addressState"
	AutoFillHintPostalCode         AutoFillHint = "postalCode"
	AutoFillHintCountryName        AutoFillHint = "countryName"
	AutoFillHintCreditCardNumber   AutoFillHint = "creditCardNumber"
)
