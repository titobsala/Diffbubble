package intro

// AnimationType represents the type of intro animation to display
type AnimationType int

const (
	// Glitch represents the "Merge Conflict" chromatic aberration animation
	Glitch AnimationType = iota
	// Scan represents the "Diff Scan" animation (red to green transition)
	Scan
	// Stream represents the "Commit Stream" Matrix-style rain animation
	Stream
)

// SymbolStream represents a falling symbol stream for the Stream animation
type SymbolStream struct {
	X      int    // Column position
	Y      int    // Current Y position
	Speed  int    // How many frames per move
	Symbol string // The symbol to display
	Color  string // Color hex code
}

// GetMaxFrames returns the total number of frames for a given animation type
func GetMaxFrames(animType AnimationType) int {
	switch animType {
	case Glitch:
		return 30 // ~1.5 seconds at 20 FPS
	case Scan:
		return 40 // ~2 seconds at 20 FPS
	case Stream:
		return 50 // ~2.5 seconds at 20 FPS
	default:
		return 30
	}
}
