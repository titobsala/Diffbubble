package ui

// Theme defines colors for all UI elements
type Theme struct {
	Name string

	// Diff colors
	AdditionBg string
	AdditionFg string
	DeletionBg string
	DeletionFg string
	ContextFg  string
	HeaderFg   string

	// UI colors
	BorderColor        string
	FocusedBorderColor string
	TitleFg            string

	// File list colors
	ModifiedFg string
	AddedFg    string
	DeletedFg  string

	// General
	Background string
	Foreground string
}

// currentTheme stores the active theme
var currentTheme Theme

// themes holds all available themes
var themes = map[string]Theme{
	"dark":          DarkTheme(),
	"light":         LightTheme(),
	"high-contrast": HighContrastTheme(),
	"solarized":     SolarizedTheme(),
	"dracula":       DraculaTheme(),
	"github":        GitHubTheme(),
	"catppuccin":    CatppuccinTheme(),
	"tokyo-night":   TokyoNightTheme(),
	"one-dark":      OneDarkTheme(),
}

// SetTheme sets the active theme and updates all styles
func SetTheme(name string) {
	if t, ok := themes[name]; ok {
		currentTheme = t
	} else {
		currentTheme = DarkTheme() // fallback to dark
	}
	updateStyles()
}

// GetTheme returns the current theme
func GetTheme() Theme {
	return currentTheme
}

// ListThemes returns all available theme names
func ListThemes() []string {
	return []string{"dark", "light", "high-contrast", "solarized", "dracula", "github", "catppuccin", "tokyo-night", "one-dark"}
}

// ValidateTheme checks if a theme name is valid
func ValidateTheme(name string) bool {
	_, ok := themes[name]
	return ok
}

// DarkTheme - Current default dark theme
func DarkTheme() Theme {
	return Theme{
		Name: "dark",

		// Diff colors (keeping current colors)
		AdditionBg: "#1a3a1a", // dark green
		AdditionFg: "#43BF6D", // bright green
		DeletionBg: "#3a1a1a", // dark red
		DeletionFg: "#E05252", // bright red
		ContextFg:  "#8B8B8B", // gray
		HeaderFg:   "#666666", // darker gray

		// UI colors
		BorderColor:        "#5C5C5C",
		FocusedBorderColor: "#A855F7", // purple
		TitleFg:            "#FFFFFF",

		// File list colors
		ModifiedFg: "#F5C842", // yellow
		AddedFg:    "#43BF6D", // green
		DeletedFg:  "#E05252", // red

		// General
		Background: "#000000",
		Foreground: "#FFFFFF",
	}
}

// LightTheme - Light background theme
func LightTheme() Theme {
	return Theme{
		Name: "light",

		// Diff colors
		AdditionBg: "#D4F1D4", // light green background
		AdditionFg: "#0B6622", // dark green text
		DeletionBg: "#F1D4D4", // light red background
		DeletionFg: "#B62020", // dark red text
		ContextFg:  "#4A4A4A", // dark gray
		HeaderFg:   "#6A6A6A", // medium gray

		// UI colors
		BorderColor:        "#CCCCCC",
		FocusedBorderColor: "#8B5CF6", // purple
		TitleFg:            "#000000",

		// File list colors
		ModifiedFg: "#D97706", // orange
		AddedFg:    "#16A34A", // green
		DeletedFg:  "#DC2626", // red

		// General
		Background: "#FFFFFF",
		Foreground: "#000000",
	}
}

// HighContrastTheme - Maximum contrast for accessibility
func HighContrastTheme() Theme {
	return Theme{
		Name: "high-contrast",

		// Diff colors - very high contrast
		AdditionBg: "#003300", // very dark green
		AdditionFg: "#00FF00", // bright green
		DeletionBg: "#330000", // very dark red
		DeletionFg: "#FF0000", // bright red
		ContextFg:  "#FFFFFF", // white (high contrast)
		HeaderFg:   "#FFFF00", // yellow

		// UI colors
		BorderColor:        "#FFFFFF",
		FocusedBorderColor: "#FFFF00", // yellow for high visibility
		TitleFg:            "#FFFFFF",

		// File list colors
		ModifiedFg: "#FFFF00", // yellow
		AddedFg:    "#00FF00", // bright green
		DeletedFg:  "#FF0000", // bright red

		// General
		Background: "#000000",
		Foreground: "#FFFFFF",
	}
}

// SolarizedTheme - Solarized Dark color scheme
func SolarizedTheme() Theme {
	return Theme{
		Name: "solarized",

		// Diff colors - Solarized palette
		AdditionBg: "#0D3A2E", // base02 with green tint
		AdditionFg: "#859900", // green
		DeletionBg: "#3A0D0D", // base02 with red tint
		DeletionFg: "#DC322F", // red
		ContextFg:  "#657B83", // base00
		HeaderFg:   "#586E75", // base01

		// UI colors
		BorderColor:        "#073642", // base02
		FocusedBorderColor: "#6C71C4", // violet
		TitleFg:            "#839496", // base0

		// File list colors
		ModifiedFg: "#B58900", // yellow
		AddedFg:    "#859900", // green
		DeletedFg:  "#DC322F", // red

		// General
		Background: "#002B36", // base03
		Foreground: "#839496", // base0
	}
}

// DraculaTheme - Dracula color scheme
func DraculaTheme() Theme {
	return Theme{
		Name: "dracula",

		// Diff colors - Dracula palette
		AdditionBg: "#1A2A1A", // dark green
		AdditionFg: "#50FA7B", // green
		DeletionBg: "#2A1A1A", // dark red
		DeletionFg: "#FF5555", // red
		ContextFg:  "#F8F8F2", // foreground
		HeaderFg:   "#6272A4", // comment

		// UI colors
		BorderColor:        "#44475A", // current line
		FocusedBorderColor: "#BD93F9", // purple
		TitleFg:            "#F8F8F2", // foreground

		// File list colors
		ModifiedFg: "#F1FA8C", // yellow
		AddedFg:    "#50FA7B", // green
		DeletedFg:  "#FF5555", // red

		// General
		Background: "#282A36", // background
		Foreground: "#F8F8F2", // foreground
	}
}

// GitHubTheme - GitHub-style diff colors
func GitHubTheme() Theme {
	return Theme{
		Name: "github",

		// Diff colors - GitHub palette
		AdditionBg: "#E6FFED", // light green
		AdditionFg: "#24292F", // dark text
		DeletionBg: "#FFEBE9", // light red
		DeletionFg: "#24292F", // dark text
		ContextFg:  "#57606A", // gray
		HeaderFg:   "#6E7781", // muted gray

		// UI colors
		BorderColor:        "#D0D7DE",
		FocusedBorderColor: "#0969DA", // blue
		TitleFg:            "#24292F",

		// File list colors
		ModifiedFg: "#9A6700", // yellow-brown
		AddedFg:    "#1A7F37", // green
		DeletedFg:  "#CF222E", // red

		// General
		Background: "#FFFFFF",
		Foreground: "#24292F",
	}
}

// CatppuccinTheme - Catppuccin Mocha color scheme
func CatppuccinTheme() Theme {
	return Theme{
		Name: "catppuccin",

		// Diff colors - Catppuccin Mocha palette
		AdditionBg: "#1E2D2F", // base with green tint
		AdditionFg: "#A6E3A1", // green
		DeletionBg: "#2D1E1E", // base with red tint
		DeletionFg: "#F38BA8", // red
		ContextFg:  "#CDD6F4", // text
		HeaderFg:   "#6C7086", // overlay0

		// UI colors
		BorderColor:        "#45475A", // surface1
		FocusedBorderColor: "#CBA6F7", // mauve
		TitleFg:            "#CDD6F4", // text

		// File list colors
		ModifiedFg: "#F9E2AF", // yellow
		AddedFg:    "#A6E3A1", // green
		DeletedFg:  "#F38BA8", // red

		// General
		Background: "#1E1E2E", // base
		Foreground: "#CDD6F4", // text
	}
}

// TokyoNightTheme - Tokyo Night color scheme
func TokyoNightTheme() Theme {
	return Theme{
		Name: "tokyo-night",

		// Diff colors - Tokyo Night palette
		AdditionBg: "#1A2B32", // dark with green tint
		AdditionFg: "#9ECE6A", // green
		DeletionBg: "#2B1A1A", // dark with red tint
		DeletionFg: "#F7768E", // red
		ContextFg:  "#A9B1D6", // foreground
		HeaderFg:   "#565F89", // comment

		// UI colors
		BorderColor:        "#3B4261", // border
		FocusedBorderColor: "#BB9AF7", // purple
		TitleFg:            "#C0CAF5", // foreground bright

		// File list colors
		ModifiedFg: "#E0AF68", // yellow
		AddedFg:    "#9ECE6A", // green
		DeletedFg:  "#F7768E", // red

		// General
		Background: "#1A1B26", // background
		Foreground: "#A9B1D6", // foreground
	}
}

// OneDarkTheme - Atom's One Dark color scheme
func OneDarkTheme() Theme {
	return Theme{
		Name: "one-dark",

		// Diff colors - One Dark palette
		AdditionBg: "#1C2B1F", // dark with green tint
		AdditionFg: "#98C379", // green
		DeletionBg: "#2B1C1C", // dark with red tint
		DeletionFg: "#E06C75", // red
		ContextFg:  "#ABB2BF", // mono-1
		HeaderFg:   "#5C6370", // mono-3

		// UI colors
		BorderColor:        "#3E4451", // gutter
		FocusedBorderColor: "#C678DD", // purple
		TitleFg:            "#DCDFE4", // mono-0

		// File list colors
		ModifiedFg: "#E5C07B", // yellow
		AddedFg:    "#98C379", // green
		DeletedFg:  "#E06C75", // red

		// General
		Background: "#282C34", // background
		Foreground: "#ABB2BF", // mono-1
	}
}

func init() {
	// Initialize with dark theme by default
	SetTheme("dark")
}
