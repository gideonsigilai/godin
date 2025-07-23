package widgets

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gideonsigilai/godin/pkg/state"
)

// Color represents a color value
type Color string

// Common color constants
const (
	// Basic colors
	ColorTransparent Color = "transparent"
	ColorBlack       Color = "#000000"
	ColorWhite       Color = "#ffffff"

	// Primary colors
	ColorRed        Color = "#f44336"
	ColorGreen      Color = "#4caf50"
	ColorBlue       Color = "#2196f3"
	ColorYellow     Color = "#ffeb3b"
	ColorOrange     Color = "#ff9800"
	ColorPurple     Color = "#9c27b0"
	ColorGrey       Color = "#9e9e9e"
	ColorCyan       Color = "#00bcd4"
	ColorPink       Color = "#e91e63"
	ColorIndigo     Color = "#3f51b5"
	ColorTeal       Color = "#009688"
	ColorLime       Color = "#cddc39"
	ColorAmber      Color = "#ffc107"
	ColorDeepOrange Color = "#ff5722"
	ColorDeepPurple Color = "#673ab7"
	ColorLightBlue  Color = "#03a9f4"
	ColorLightGreen Color = "#8bc34a"
	ColorBrown      Color = "#795548"
	ColorBlueGrey   Color = "#607d8b"

	// Red shades
	ColorRed50   Color = "#ffebee"
	ColorRed100  Color = "#ffcdd2"
	ColorRed200  Color = "#ef9a9a"
	ColorRed300  Color = "#e57373"
	ColorRed400  Color = "#ef5350"
	ColorRed500  Color = "#f44336"
	ColorRed600  Color = "#e53935"
	ColorRed700  Color = "#d32f2f"
	ColorRed800  Color = "#c62828"
	ColorRed900  Color = "#b71c1c"
	ColorRedA100 Color = "#ff8a80"
	ColorRedA200 Color = "#ff5252"
	ColorRedA400 Color = "#ff1744"
	ColorRedA700 Color = "#d50000"

	// Pink shades
	ColorPink50   Color = "#fce4ec"
	ColorPink100  Color = "#f8bbd9"
	ColorPink200  Color = "#f48fb1"
	ColorPink300  Color = "#f06292"
	ColorPink400  Color = "#ec407a"
	ColorPink500  Color = "#e91e63"
	ColorPink600  Color = "#d81b60"
	ColorPink700  Color = "#c2185b"
	ColorPink800  Color = "#ad1457"
	ColorPink900  Color = "#880e4f"
	ColorPinkA100 Color = "#ff80ab"
	ColorPinkA200 Color = "#ff4081"
	ColorPinkA400 Color = "#f50057"
	ColorPinkA700 Color = "#c51162"

	// Purple shades
	ColorPurple50   Color = "#f3e5f5"
	ColorPurple100  Color = "#e1bee7"
	ColorPurple200  Color = "#ce93d8"
	ColorPurple300  Color = "#ba68c8"
	ColorPurple400  Color = "#ab47bc"
	ColorPurple500  Color = "#9c27b0"
	ColorPurple600  Color = "#8e24aa"
	ColorPurple700  Color = "#7b1fa2"
	ColorPurple800  Color = "#6a1b9a"
	ColorPurple900  Color = "#4a148c"
	ColorPurpleA100 Color = "#ea80fc"
	ColorPurpleA200 Color = "#e040fb"
	ColorPurpleA400 Color = "#d500f9"
	ColorPurpleA700 Color = "#aa00ff"

	// Deep Purple shades
	ColorDeepPurple50   Color = "#ede7f6"
	ColorDeepPurple100  Color = "#d1c4e9"
	ColorDeepPurple200  Color = "#b39ddb"
	ColorDeepPurple300  Color = "#9575cd"
	ColorDeepPurple400  Color = "#7e57c2"
	ColorDeepPurple500  Color = "#673ab7"
	ColorDeepPurple600  Color = "#5e35b1"
	ColorDeepPurple700  Color = "#512da8"
	ColorDeepPurple800  Color = "#4527a0"
	ColorDeepPurple900  Color = "#311b92"
	ColorDeepPurpleA100 Color = "#b388ff"
	ColorDeepPurpleA200 Color = "#7c4dff"
	ColorDeepPurpleA400 Color = "#651fff"
	ColorDeepPurpleA700 Color = "#6200ea"

	// Indigo shades
	ColorIndigo50   Color = "#e8eaf6"
	ColorIndigo100  Color = "#c5cae9"
	ColorIndigo200  Color = "#9fa8da"
	ColorIndigo300  Color = "#7986cb"
	ColorIndigo400  Color = "#5c6bc0"
	ColorIndigo500  Color = "#3f51b5"
	ColorIndigo600  Color = "#3949ab"
	ColorIndigo700  Color = "#303f9f"
	ColorIndigo800  Color = "#283593"
	ColorIndigo900  Color = "#1a237e"
	ColorIndigoA100 Color = "#8c9eff"
	ColorIndigoA200 Color = "#536dfe"
	ColorIndigoA400 Color = "#3d5afe"
	ColorIndigoA700 Color = "#304ffe"

	// Blue shades
	ColorBlue50   Color = "#e3f2fd"
	ColorBlue100  Color = "#bbdefb"
	ColorBlue200  Color = "#90caf9"
	ColorBlue300  Color = "#64b5f6"
	ColorBlue400  Color = "#42a5f5"
	ColorBlue500  Color = "#2196f3"
	ColorBlue600  Color = "#1e88e5"
	ColorBlue700  Color = "#1976d2"
	ColorBlue800  Color = "#1565c0"
	ColorBlue900  Color = "#0d47a1"
	ColorBlueA100 Color = "#82b1ff"
	ColorBlueA200 Color = "#448aff"
	ColorBlueA400 Color = "#2979ff"
	ColorBlueA700 Color = "#2962ff"

	// Light Blue shades
	ColorLightBlue50   Color = "#e1f5fe"
	ColorLightBlue100  Color = "#b3e5fc"
	ColorLightBlue200  Color = "#81d4fa"
	ColorLightBlue300  Color = "#4fc3f7"
	ColorLightBlue400  Color = "#29b6f6"
	ColorLightBlue500  Color = "#03a9f4"
	ColorLightBlue600  Color = "#039be5"
	ColorLightBlue700  Color = "#0288d1"
	ColorLightBlue800  Color = "#0277bd"
	ColorLightBlue900  Color = "#01579b"
	ColorLightBlueA100 Color = "#80d8ff"
	ColorLightBlueA200 Color = "#40c4ff"
	ColorLightBlueA400 Color = "#00b0ff"
	ColorLightBlueA700 Color = "#0091ea"

	// Cyan shades
	ColorCyan50   Color = "#e0f6ff"
	ColorCyan100  Color = "#b2ebf2"
	ColorCyan200  Color = "#80deea"
	ColorCyan300  Color = "#4dd0e1"
	ColorCyan400  Color = "#26c6da"
	ColorCyan500  Color = "#00bcd4"
	ColorCyan600  Color = "#00acc1"
	ColorCyan700  Color = "#0097a7"
	ColorCyan800  Color = "#00838f"
	ColorCyan900  Color = "#006064"
	ColorCyanA100 Color = "#84ffff"
	ColorCyanA200 Color = "#18ffff"
	ColorCyanA400 Color = "#00e5ff"
	ColorCyanA700 Color = "#00b8d4"

	// Teal shades
	ColorTeal50   Color = "#e0f2f1"
	ColorTeal100  Color = "#b2dfdb"
	ColorTeal200  Color = "#80cbc4"
	ColorTeal300  Color = "#4db6ac"
	ColorTeal400  Color = "#26a69a"
	ColorTeal500  Color = "#009688"
	ColorTeal600  Color = "#00897b"
	ColorTeal700  Color = "#00796b"
	ColorTeal800  Color = "#00695c"
	ColorTeal900  Color = "#004d40"
	ColorTealA100 Color = "#a7ffeb"
	ColorTealA200 Color = "#64ffda"
	ColorTealA400 Color = "#1de9b6"
	ColorTealA700 Color = "#00bfa5"

	// Green shades
	ColorGreen50   Color = "#e8f5e8"
	ColorGreen100  Color = "#c8e6c9"
	ColorGreen200  Color = "#a5d6a7"
	ColorGreen300  Color = "#81c784"
	ColorGreen400  Color = "#66bb6a"
	ColorGreen500  Color = "#4caf50"
	ColorGreen600  Color = "#43a047"
	ColorGreen700  Color = "#388e3c"
	ColorGreen800  Color = "#2e7d32"
	ColorGreen900  Color = "#1b5e20"
	ColorGreenA100 Color = "#b9f6ca"
	ColorGreenA200 Color = "#69f0ae"
	ColorGreenA400 Color = "#00e676"
	ColorGreenA700 Color = "#00c853"

	// Light Green shades
	ColorLightGreen50   Color = "#f1f8e9"
	ColorLightGreen100  Color = "#dcedc8"
	ColorLightGreen200  Color = "#c5e1a5"
	ColorLightGreen300  Color = "#aed581"
	ColorLightGreen400  Color = "#9ccc65"
	ColorLightGreen500  Color = "#8bc34a"
	ColorLightGreen600  Color = "#7cb342"
	ColorLightGreen700  Color = "#689f38"
	ColorLightGreen800  Color = "#558b2f"
	ColorLightGreen900  Color = "#33691e"
	ColorLightGreenA100 Color = "#ccff90"
	ColorLightGreenA200 Color = "#b2ff59"
	ColorLightGreenA400 Color = "#76ff03"
	ColorLightGreenA700 Color = "#64dd17"

	// Lime shades
	ColorLime50   Color = "#f9fbe7"
	ColorLime100  Color = "#f0f4c3"
	ColorLime200  Color = "#e6ee9c"
	ColorLime300  Color = "#dce775"
	ColorLime400  Color = "#d4e157"
	ColorLime500  Color = "#cddc39"
	ColorLime600  Color = "#c0ca33"
	ColorLime700  Color = "#afb42b"
	ColorLime800  Color = "#9e9d24"
	ColorLime900  Color = "#827717"
	ColorLimeA100 Color = "#f4ff81"
	ColorLimeA200 Color = "#eeff41"
	ColorLimeA400 Color = "#c6ff00"
	ColorLimeA700 Color = "#aeea00"

	// Yellow shades
	ColorYellow50   Color = "#fffde7"
	ColorYellow100  Color = "#fff9c4"
	ColorYellow200  Color = "#fff59d"
	ColorYellow300  Color = "#fff176"
	ColorYellow400  Color = "#ffee58"
	ColorYellow500  Color = "#ffeb3b"
	ColorYellow600  Color = "#fdd835"
	ColorYellow700  Color = "#fbc02d"
	ColorYellow800  Color = "#f9a825"
	ColorYellow900  Color = "#f57f17"
	ColorYellowA100 Color = "#ffff8d"
	ColorYellowA200 Color = "#ffff00"
	ColorYellowA400 Color = "#ffea00"
	ColorYellowA700 Color = "#ffd600"

	// Amber shades
	ColorAmber50   Color = "#fff8e1"
	ColorAmber100  Color = "#ffecb3"
	ColorAmber200  Color = "#ffe082"
	ColorAmber300  Color = "#ffd54f"
	ColorAmber400  Color = "#ffca28"
	ColorAmber500  Color = "#ffc107"
	ColorAmber600  Color = "#ffb300"
	ColorAmber700  Color = "#ffa000"
	ColorAmber800  Color = "#ff8f00"
	ColorAmber900  Color = "#ff6f00"
	ColorAmberA100 Color = "#ffe57f"
	ColorAmberA200 Color = "#ffd740"
	ColorAmberA400 Color = "#ffc400"
	ColorAmberA700 Color = "#ffab00"

	// Orange shades
	ColorOrange50   Color = "#fff3e0"
	ColorOrange100  Color = "#ffe0b2"
	ColorOrange200  Color = "#ffcc80"
	ColorOrange300  Color = "#ffb74d"
	ColorOrange400  Color = "#ffa726"
	ColorOrange500  Color = "#ff9800"
	ColorOrange600  Color = "#fb8c00"
	ColorOrange700  Color = "#f57c00"
	ColorOrange800  Color = "#ef6c00"
	ColorOrange900  Color = "#e65100"
	ColorOrangeA100 Color = "#ffd180"
	ColorOrangeA200 Color = "#ffab40"
	ColorOrangeA400 Color = "#ff9100"
	ColorOrangeA700 Color = "#ff6d00"

	// Deep Orange shades
	ColorDeepOrange50   Color = "#fbe9e7"
	ColorDeepOrange100  Color = "#ffccbc"
	ColorDeepOrange200  Color = "#ffab91"
	ColorDeepOrange300  Color = "#ff8a65"
	ColorDeepOrange400  Color = "#ff7043"
	ColorDeepOrange500  Color = "#ff5722"
	ColorDeepOrange600  Color = "#f4511e"
	ColorDeepOrange700  Color = "#e64a19"
	ColorDeepOrange800  Color = "#d84315"
	ColorDeepOrange900  Color = "#bf360c"
	ColorDeepOrangeA100 Color = "#ff9e80"
	ColorDeepOrangeA200 Color = "#ff6e40"
	ColorDeepOrangeA400 Color = "#ff3d00"
	ColorDeepOrangeA700 Color = "#dd2c00"

	// Brown shades
	ColorBrown50  Color = "#efebe9"
	ColorBrown100 Color = "#d7ccc8"
	ColorBrown200 Color = "#bcaaa4"
	ColorBrown300 Color = "#a1887f"
	ColorBrown400 Color = "#8d6e63"
	ColorBrown500 Color = "#795548"
	ColorBrown600 Color = "#6d4c41"
	ColorBrown700 Color = "#5d4037"
	ColorBrown800 Color = "#4e342e"
	ColorBrown900 Color = "#3e2723"

	// Grey shades
	ColorGrey50  Color = "#fafafa"
	ColorGrey100 Color = "#f5f5f5"
	ColorGrey200 Color = "#eeeeee"
	ColorGrey300 Color = "#e0e0e0"
	ColorGrey400 Color = "#bdbdbd"
	ColorGrey500 Color = "#9e9e9e"
	ColorGrey600 Color = "#757575"
	ColorGrey700 Color = "#616161"
	ColorGrey800 Color = "#424242"
	ColorGrey900 Color = "#212121"

	// Blue Grey shades
	ColorBlueGrey50  Color = "#eceff1"
	ColorBlueGrey100 Color = "#cfd8dc"
	ColorBlueGrey200 Color = "#b0bec5"
	ColorBlueGrey300 Color = "#90a4ae"
	ColorBlueGrey400 Color = "#78909c"
	ColorBlueGrey500 Color = "#607d8b"
	ColorBlueGrey600 Color = "#546e7a"
	ColorBlueGrey700 Color = "#455a64"
	ColorBlueGrey800 Color = "#37474f"
	ColorBlueGrey900 Color = "#263238"

	// CSS Named Colors
	ColorAliceBlue            Color = "#f0f8ff"
	ColorAntiqueWhite         Color = "#faebd7"
	ColorAqua                 Color = "#00ffff"
	ColorAquamarine           Color = "#7fffd4"
	ColorAzure                Color = "#f0ffff"
	ColorBeige                Color = "#f5f5dc"
	ColorBisque               Color = "#ffe4c4"
	ColorBlanchedAlmond       Color = "#ffebcd"
	ColorBlueViolet           Color = "#8a2be2"
	ColorBurlywood            Color = "#deb887"
	ColorCadetBlue            Color = "#5f9ea0"
	ColorChartreuse           Color = "#7fff00"
	ColorChocolate            Color = "#d2691e"
	ColorCoral                Color = "#ff7f50"
	ColorCornflowerBlue       Color = "#6495ed"
	ColorCornsilk             Color = "#fff8dc"
	ColorCrimson              Color = "#dc143c"
	ColorDarkBlue             Color = "#00008b"
	ColorDarkCyan             Color = "#008b8b"
	ColorDarkGoldenrod        Color = "#b8860b"
	ColorDarkGray             Color = "#a9a9a9"
	ColorDarkGreen            Color = "#006400"
	ColorDarkKhaki            Color = "#bdb76b"
	ColorDarkMagenta          Color = "#8b008b"
	ColorDarkOliveGreen       Color = "#556b2f"
	ColorDarkOrange           Color = "#ff8c00"
	ColorDarkOrchid           Color = "#9932cc"
	ColorDarkRed              Color = "#8b0000"
	ColorDarkSalmon           Color = "#e9967a"
	ColorDarkSeaGreen         Color = "#8fbc8f"
	ColorDarkSlateBlue        Color = "#483d8b"
	ColorDarkSlateGray        Color = "#2f4f4f"
	ColorDarkTurquoise        Color = "#00ced1"
	ColorDarkViolet           Color = "#9400d3"
	ColorDeepSkyBlue          Color = "#00bfff"
	ColorDimGray              Color = "#696969"
	ColorDodgerBlue           Color = "#1e90ff"
	ColorFirebrick            Color = "#b22222"
	ColorFloralWhite          Color = "#fffaf0"
	ColorForestGreen          Color = "#228b22"
	ColorFuchsia              Color = "#ff00ff"
	ColorGainsboro            Color = "#dcdcdc"
	ColorGhostWhite           Color = "#f8f8ff"
	ColorGold                 Color = "#ffd700"
	ColorGoldenrod            Color = "#daa520"
	ColorGreenYellow          Color = "#adff2f"
	ColorHoneydew             Color = "#f0fff0"
	ColorHotPink              Color = "#ff69b4"
	ColorIndianRed            Color = "#cd5c5c"
	ColorIvory                Color = "#fffff0"
	ColorKhaki                Color = "#f0e68c"
	ColorLavender             Color = "#e6e6fa"
	ColorLavenderBlush        Color = "#fff0f5"
	ColorLawnGreen            Color = "#7cfc00"
	ColorLemonChiffon         Color = "#fffacd"
	ColorLightCoral           Color = "#f08080"
	ColorLightCyan            Color = "#e0ffff"
	ColorLightGoldenrodYellow Color = "#fafad2"
	ColorLightGray            Color = "#d3d3d3"
	ColorLightPink            Color = "#ffb6c1"
	ColorLightSalmon          Color = "#ffa07a"
	ColorLightSeaGreen        Color = "#20b2aa"
	ColorLightSkyBlue         Color = "#87cefa"
	ColorLightSlateGray       Color = "#778899"
	ColorLightSteelBlue       Color = "#b0c4de"
	ColorLightYellow          Color = "#ffffe0"
	ColorLimeGreen            Color = "#32cd32"
	ColorLinen                Color = "#faf0e6"
	ColorMagenta              Color = "#ff00ff"
	ColorMaroon               Color = "#800000"
	ColorMediumAquamarine     Color = "#66cdaa"
	ColorMediumBlue           Color = "#0000cd"
	ColorMediumOrchid         Color = "#ba55d3"
	ColorMediumPurple         Color = "#9370db"
	ColorMediumSeaGreen       Color = "#3cb371"
	ColorMediumSlateBlue      Color = "#7b68ee"
	ColorMediumSpringGreen    Color = "#00fa9a"
	ColorMediumTurquoise      Color = "#48d1cc"
	ColorMediumVioletRed      Color = "#c71585"
	ColorMidnightBlue         Color = "#191970"
	ColorMintCream            Color = "#f5fffa"
	ColorMistyRose            Color = "#ffe4e1"
	ColorMoccasin             Color = "#ffe4b5"
	ColorNavajoWhite          Color = "#ffdead"
	ColorNavy                 Color = "#000080"
	ColorOldLace              Color = "#fdf5e6"
	ColorOlive                Color = "#808000"
	ColorOliveDrab            Color = "#6b8e23"
	ColorOrangeRed            Color = "#ff4500"
	ColorOrchid               Color = "#da70d6"
	ColorPaleGoldenrod        Color = "#eee8aa"
	ColorPaleGreen            Color = "#98fb98"
	ColorPaleTurquoise        Color = "#afeeee"
	ColorPaleVioletRed        Color = "#db7093"
	ColorPapayaWhip           Color = "#ffefd5"
	ColorPeachPuff            Color = "#ffdab9"
	ColorPeru                 Color = "#cd853f"
	ColorPlum                 Color = "#dda0dd"
	ColorPowderBlue           Color = "#b0e0e6"
	ColorRebeccaPurple        Color = "#663399"
	ColorRosyBrown            Color = "#bc8f8f"
	ColorRoyalBlue            Color = "#4169e1"
	ColorSaddleBrown          Color = "#8b4513"
	ColorSalmon               Color = "#fa8072"
	ColorSandyBrown           Color = "#f4a460"
	ColorSeaGreen             Color = "#2e8b57"
	ColorSeashell             Color = "#fff5ee"
	ColorSienna               Color = "#a0522d"
	ColorSilver               Color = "#c0c0c0"
	ColorSkyBlue              Color = "#87ceeb"
	ColorSlateBlue            Color = "#6a5acd"
	ColorSlateGray            Color = "#708090"
	ColorSnow                 Color = "#fffafa"
	ColorSpringGreen          Color = "#00ff7f"
	ColorSteelBlue            Color = "#4682b4"
	ColorTan                  Color = "#d2b48c"
	ColorThistle              Color = "#d8bfd8"
	ColorTomato               Color = "#ff6347"
	ColorTurquoise            Color = "#40e0d0"
	ColorViolet               Color = "#ee82ee"
	ColorWheat                Color = "#f5deb3"
	ColorWhiteSmoke           Color = "#f5f5f5"
	ColorYellowGreen          Color = "#9acd32"

	// Additional Utility Colors
	ColorSuccess   Color = "#28a745"
	ColorInfo      Color = "#17a2b8"
	ColorWarning   Color = "#ffc107"
	ColorDanger    Color = "#dc3545"
	ColorLight     Color = "#f8f9fa"
	ColorDark      Color = "#343a40"
	ColorMuted     Color = "#6c757d"
	ColorPrimary   Color = "#007bff"
	ColorSecondary Color = "#6c757d"

	// Semantic Colors
	ColorError        Color = "#f44336"
	ColorAlert        Color = "#ff9800"
	ColorNotification Color = "#2196f3"
	ColorHighlight    Color = "#ffeb3b"
	ColorDisabled     Color = "#bdbdbd"
	ColorPlaceholder  Color = "#9e9e9e"
	ColorDivider      Color = "#e0e0e0"
	ColorSurface      Color = "#ffffff"
	ColorBackground   Color = "#fafafa"
	ColorOnSurface    Color = "#000000"
	ColorOnBackground Color = "#000000"
	ColorOnPrimary    Color = "#ffffff"
	ColorOnSecondary  Color = "#ffffff"
	ColorOnError      Color = "#ffffff"

	// Pastel Colors
	ColorPastelRed     Color = "#ffb3ba"
	ColorPastelOrange  Color = "#ffdfba"
	ColorPastelYellow  Color = "#ffffba"
	ColorPastelGreen   Color = "#baffc9"
	ColorPastelBlue    Color = "#bae1ff"
	ColorPastelPurple  Color = "#d4baff"
	ColorPastelPink    Color = "#ffb3ff"
	ColorPastelCyan    Color = "#baffff"
	ColorPastelLime    Color = "#e6ffba"
	ColorPastelMagenta Color = "#ffbaff"

	// Neon Colors
	ColorNeonRed     Color = "#ff073a"
	ColorNeonOrange  Color = "#ff8c00"
	ColorNeonYellow  Color = "#ffff00"
	ColorNeonGreen   Color = "#39ff14"
	ColorNeonBlue    Color = "#1b03a3"
	ColorNeonPurple  Color = "#bc13fe"
	ColorNeonPink    Color = "#ff10f0"
	ColorNeonCyan    Color = "#00ffff"
	ColorNeonLime    Color = "#ccff00"
	ColorNeonMagenta Color = "#ff00ff"

	// Earth Tones
	ColorEarthBrown      Color = "#8b4513"
	ColorEarthTan        Color = "#d2b48c"
	ColorEarthBeige      Color = "#f5f5dc"
	ColorEarthOlive      Color = "#808000"
	ColorEarthForest     Color = "#228b22"
	ColorEarthSage       Color = "#9caf88"
	ColorEarthTerracotta Color = "#e2725b"
	ColorEarthSienna     Color = "#a0522d"
	ColorEarthUmber      Color = "#635147"
	ColorEarthOchre      Color = "#cc7722"

	// Jewel Tones
	ColorJewelRuby      Color = "#e0115f"
	ColorJewelEmerald   Color = "#50c878"
	ColorJewelSapphire  Color = "#0f52ba"
	ColorJewelAmethyst  Color = "#9966cc"
	ColorJewelTopaz     Color = "#ffc87c"
	ColorJewelGarnet    Color = "#733635"
	ColorJewelTurquoise Color = "#40e0d0"
	ColorJewelOpal      Color = "#a8c3bc"
	ColorJewelPearl     Color = "#f8f6f0"
	ColorJewelOnyx      Color = "#353839"

	// Metallic Colors
	ColorMetallicGold     Color = "#d4af37"
	ColorMetallicSilver   Color = "#c0c0c0"
	ColorMetallicCopper   Color = "#b87333"
	ColorMetallicBronze   Color = "#cd7f32"
	ColorMetallicPlatinum Color = "#e5e4e2"
	ColorMetallicChrome   Color = "#c4c4c4"
	ColorMetallicTitanium Color = "#878681"
	ColorMetallicSteel    Color = "#71797e"
	ColorMetallicIron     Color = "#a19d94"
	ColorMetallicLead     Color = "#212121"

	// Vintage Colors
	ColorVintageRose     Color = "#fdeaa7"
	ColorVintageBlush    Color = "#fab1a0"
	ColorVintagePeach    Color = "#e17055"
	ColorVintageApricot  Color = "#fdcb6e"
	ColorVintageMint     Color = "#00b894"
	ColorVintageSage     Color = "#6c5ce7"
	ColorVintageNavy     Color = "#2d3436"
	ColorVintageBurgundy Color = "#a29bfe"
	ColorVintageMaroon   Color = "#fd79a8"
	ColorVintageOlive    Color = "#e84393"

	// Modern UI Colors
	ColorModernBlue   Color = "#4285f4"
	ColorModernGreen  Color = "#34a853"
	ColorModernYellow Color = "#fbbc05"
	ColorModernRed    Color = "#ea4335"
	ColorModernPurple Color = "#9c27b0"
	ColorModernTeal   Color = "#00acc1"
	ColorModernOrange Color = "#ff7043"
	ColorModernPink   Color = "#e91e63"
	ColorModernIndigo Color = "#3f51b5"
	ColorModernCyan   Color = "#00bcd4"

	// Gradient Start/End Colors
	ColorGradientSunrise1 Color = "#ff9a9e"
	ColorGradientSunrise2 Color = "#fecfef"
	ColorGradientSunset1  Color = "#fa709a"
	ColorGradientSunset2  Color = "#fee140"
	ColorGradientOcean1   Color = "#667eea"
	ColorGradientOcean2   Color = "#764ba2"
	ColorGradientForest1  Color = "#11998e"
	ColorGradientForest2  Color = "#38ef7d"
	ColorGradientFire1    Color = "#f093fb"
	ColorGradientFire2    Color = "#f5576c"
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

// TextEditingController manages text input state similar to Flutter
type TextEditingController struct {
	text      string
	selection TextSelection
	notifier  *state.StringNotifier
	listeners []func(string)
	mutex     sync.RWMutex
	id        string
}

// TextSelection represents text selection in a text field
type TextSelection struct {
	Start int // Start position of selection
	End   int // End position of selection
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

// TextEditingController methods

// NewTextEditingController creates a new TextEditingController with optional initial text
func NewTextEditingController(initialText string) *TextEditingController {
	controller := &TextEditingController{
		text:      initialText,
		selection: TextSelection{Start: len(initialText), End: len(initialText)},
		listeners: make([]func(string), 0),
		id:        generateControllerID(),
	}

	// Create a ValueNotifier for state management integration
	controller.notifier = state.NewStringNotifierWithID(controller.id, initialText)

	return controller
}

// Text returns the current text content
func (tec *TextEditingController) Text() string {
	tec.mutex.RLock()
	defer tec.mutex.RUnlock()
	return tec.text
}

// SetText updates the text content and notifies listeners
func (tec *TextEditingController) SetText(text string) {
	tec.mutex.Lock()
	oldText := tec.text
	tec.text = text

	// Update selection to end of text
	tec.selection = TextSelection{Start: len(text), End: len(text)}

	// Copy listeners to avoid holding lock during notification
	listeners := make([]func(string), len(tec.listeners))
	copy(listeners, tec.listeners)
	tec.mutex.Unlock()

	// Update the ValueNotifier
	tec.notifier.SetValue(text)

	// Notify listeners if text actually changed
	if oldText != text {
		for _, listener := range listeners {
			go listener(text)
		}
	}
}

// Clear clears the text content
func (tec *TextEditingController) Clear() {
	tec.SetText("")
}

// Selection returns the current text selection
func (tec *TextEditingController) Selection() TextSelection {
	tec.mutex.RLock()
	defer tec.mutex.RUnlock()
	return tec.selection
}

// SetSelection updates the text selection
func (tec *TextEditingController) SetSelection(selection TextSelection) {
	tec.mutex.Lock()
	defer tec.mutex.Unlock()

	// Validate selection bounds
	textLen := len(tec.text)
	if selection.Start < 0 {
		selection.Start = 0
	}
	if selection.Start > textLen {
		selection.Start = textLen
	}
	if selection.End < 0 {
		selection.End = 0
	}
	if selection.End > textLen {
		selection.End = textLen
	}

	tec.selection = selection
}

// AddListener adds a listener function that will be called when text changes
func (tec *TextEditingController) AddListener(listener func(string)) {
	tec.mutex.Lock()
	defer tec.mutex.Unlock()
	tec.listeners = append(tec.listeners, listener)
}

// RemoveListener removes a specific listener (simplified implementation)
func (tec *TextEditingController) RemoveListener(listener func(string)) {
	// Note: This is a simplified implementation
	// In practice, you'd need listener IDs for proper removal
	tec.mutex.Lock()
	defer tec.mutex.Unlock()

	// For now, we'll just clear all listeners as a workaround
	// A proper implementation would require listener registration with IDs
}

// ClearListeners removes all listeners
func (tec *TextEditingController) ClearListeners() {
	tec.mutex.Lock()
	defer tec.mutex.Unlock()
	tec.listeners = make([]func(string), 0)
}

// ValueNotifier returns the underlying ValueNotifier for state management integration
func (tec *TextEditingController) ValueNotifier() *state.StringNotifier {
	return tec.notifier
}

// ID returns the unique identifier for this controller
func (tec *TextEditingController) ID() string {
	return tec.id
}

// IsEmpty returns true if the text is empty
func (tec *TextEditingController) IsEmpty() bool {
	tec.mutex.RLock()
	defer tec.mutex.RUnlock()
	return len(tec.text) == 0
}

// Length returns the length of the text
func (tec *TextEditingController) Length() int {
	tec.mutex.RLock()
	defer tec.mutex.RUnlock()
	return len(tec.text)
}

// HasSelection returns true if there is a text selection (not just cursor position)
func (tec *TextEditingController) HasSelection() bool {
	tec.mutex.RLock()
	defer tec.mutex.RUnlock()
	return tec.selection.Start != tec.selection.End
}

// SelectedText returns the currently selected text
func (tec *TextEditingController) SelectedText() string {
	tec.mutex.RLock()
	defer tec.mutex.RUnlock()

	if tec.selection.Start == tec.selection.End {
		return ""
	}

	start := tec.selection.Start
	end := tec.selection.End

	// Ensure start <= end
	if start > end {
		start, end = end, start
	}

	// Validate bounds
	if start < 0 || end > len(tec.text) {
		return ""
	}

	return tec.text[start:end]
}

// SelectAll selects all text
func (tec *TextEditingController) SelectAll() {
	tec.mutex.Lock()
	defer tec.mutex.Unlock()
	tec.selection = TextSelection{Start: 0, End: len(tec.text)}
}

// MoveCursorToEnd moves the cursor to the end of the text
func (tec *TextEditingController) MoveCursorToEnd() {
	tec.mutex.Lock()
	defer tec.mutex.Unlock()
	textLen := len(tec.text)
	tec.selection = TextSelection{Start: textLen, End: textLen}
}

// MoveCursorToStart moves the cursor to the beginning of the text
func (tec *TextEditingController) MoveCursorToStart() {
	tec.mutex.Lock()
	defer tec.mutex.Unlock()
	tec.selection = TextSelection{Start: 0, End: 0}
}

// InsertText inserts text at the current cursor position
func (tec *TextEditingController) InsertText(insertText string) {
	tec.mutex.Lock()

	// Get current state
	currentText := tec.text
	selection := tec.selection

	// Ensure start <= end
	start := selection.Start
	end := selection.End
	if start > end {
		start, end = end, start
	}

	// Validate bounds
	if start < 0 {
		start = 0
	}
	if end > len(currentText) {
		end = len(currentText)
	}

	// Build new text
	newText := currentText[:start] + insertText + currentText[end:]

	// Update text and selection
	tec.text = newText
	newCursorPos := start + len(insertText)
	tec.selection = TextSelection{Start: newCursorPos, End: newCursorPos}

	// Copy listeners
	listeners := make([]func(string), len(tec.listeners))
	copy(listeners, tec.listeners)

	tec.mutex.Unlock()

	// Update the ValueNotifier
	tec.notifier.SetValue(newText)

	// Notify listeners
	for _, listener := range listeners {
		go listener(newText)
	}
}

// DeleteSelection deletes the currently selected text
func (tec *TextEditingController) DeleteSelection() {
	tec.mutex.Lock()

	selection := tec.selection
	if selection.Start == selection.End {
		// No selection, nothing to delete
		tec.mutex.Unlock()
		return
	}

	// Ensure start <= end
	start := selection.Start
	end := selection.End
	if start > end {
		start, end = end, start
	}

	// Build new text without selected portion
	currentText := tec.text
	newText := currentText[:start] + currentText[end:]

	// Update text and selection
	tec.text = newText
	tec.selection = TextSelection{Start: start, End: start}

	// Copy listeners
	listeners := make([]func(string), len(tec.listeners))
	copy(listeners, tec.listeners)

	tec.mutex.Unlock()

	// Update the ValueNotifier
	tec.notifier.SetValue(newText)

	// Notify listeners
	for _, listener := range listeners {
		go listener(newText)
	}
}

// Dispose cleans up the controller resources
func (tec *TextEditingController) Dispose() {
	tec.mutex.Lock()
	defer tec.mutex.Unlock()

	// Clear listeners
	tec.listeners = make([]func(string), 0)

	// Clear the ValueNotifier
	if tec.notifier != nil {
		tec.notifier.ClearListeners()
	}
}

// generateControllerID generates a unique ID for the controller
func generateControllerID() string {
	// Simple ID generation - in practice you might want something more sophisticated
	return fmt.Sprintf("tec_%d", time.Now().UnixNano())
}
