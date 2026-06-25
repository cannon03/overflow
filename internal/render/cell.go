package render

// Color represents a 24-bit RGB color.
type Color struct {
	R, G, B uint8
}

// Cell represents a single character cell in the framebuffer.
type Cell struct {
	Rune        rune
	FG          Color
	BG          Color
	Transparent bool // if true, this cell is not drawn (for sprites with alpha)
}

// Predefined colors.
var (
	ColorBlack   = Color{0, 0, 0}
	ColorWhite   = Color{200, 200, 220}
	ColorRed     = Color{255, 50, 50}
	ColorGreen   = Color{50, 255, 50}
	ColorBlue    = Color{50, 50, 255}
	ColorYellow  = Color{255, 255, 80}
	ColorCyan    = Color{50, 255, 255}
	ColorMagenta = Color{255, 50, 255}
	ColorOrange  = Color{255, 160, 20}
	ColorPink    = Color{255, 100, 180}
	ColorDarkRed = Color{180, 0, 0}

	ColorDarkBlue  = Color{10, 10, 30}
	ColorDarkGray  = Color{60, 60, 70}
	ColorMidGray   = Color{120, 120, 130}
	ColorBrightCyan  = Color{100, 255, 255}
	ColorBrightYellow = Color{255, 255, 150}

	ColorHudBG    = Color{15, 15, 35}
	ColorHudBorder = Color{80, 80, 120}
	ColorHudText  = Color{180, 180, 200}
	ColorHudHP    = Color{255, 80, 80}
	ColorHudHPLost = Color{40, 20, 20}
	ColorHudScore = Color{255, 220, 80}
	ColorHudWave  = Color{80, 180, 255}
)

var DefaultCell = Cell{Rune: ' ', FG: ColorWhite, BG: ColorBlack}

// Cell comparisons with Transparent handling.
func (c Cell) Equal(other Cell) bool {
	if c.Transparent != other.Transparent {
		return false
	}
	if c.Rune != other.Rune {
		return false
	}
	if c.FG != other.FG {
		return false
	}
	if c.BG != other.BG {
		return false
	}
	return true
}
