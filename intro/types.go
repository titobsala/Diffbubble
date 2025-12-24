package intro

// AnimationType represents the type of intro animation to display
type AnimationType int

const (
	// Glitch represents the "Merge Conflict" chromatic aberration animation
	Glitch AnimationType = iota
	// Scan represents the "Diff Scan" animation (red to green transition)
	Scan
	// Lightning represents the "Cat Working at Computer" animation (red cat → loading → glitch → green cat)
	Lightning
)

// GetMaxFrames returns the total number of frames for a given animation type
func GetMaxFrames(animType AnimationType) int {
	switch animType {
	case Glitch:
		return 35 // ~1.75 seconds at 20 FPS (5 phases)
	case Scan:
		return 40 // ~2 seconds at 20 FPS
	case Lightning:
		return 50 // ~2.5 seconds at 20 FPS (5 phases: working, loading, complete, glitch, success)
	default:
		return 35
	}
}
