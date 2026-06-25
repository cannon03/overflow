package render

// Sprite represents a small grid of cells that can be drawn to a framebuffer.
type Sprite struct {
	Width  int
	Height int
	Cells  []Cell // row-major: cells[y*Width+x]
}

// NewSprite creates a sprite from a grid of runes with a single FG/BG color.
func NewSprite(width, height int, runes []rune, fg, bg Color) *Sprite {
	cells := make([]Cell, width*height)
	for i, r := range runes {
		cells[i] = Cell{Rune: r, FG: fg, BG: bg}
	}
	return &Sprite{
		Width:  width,
		Height: height,
		Cells:  cells,
	}
}

// NewSpriteFromString creates a sprite from a multiline string.
// Each line is a row. All rows must be the same length.
func NewSpriteFromString(s string, fg, bg Color) *Sprite {
	lines := splitLines(s)
	if len(lines) == 0 {
		return &Sprite{Width: 0, Height: 0}
	}
	width := 0
	for _, line := range lines {
		if len(line) > width {
			width = len(line)
		}
	}
	height := len(lines)

	cells := make([]Cell, width*height)
	for y, line := range lines {
		for x, ch := range line {
			if ch == ' ' {
				cells[y*width+x] = Cell{Rune: ' ', FG: fg, BG: bg, Transparent: true}
			} else {
				cells[y*width+x] = Cell{Rune: ch, FG: fg, BG: bg}
			}
		}
	}

	return &Sprite{
		Width:  width,
		Height: height,
		Cells:  cells,
	}
}

// NewSpriteWithColors creates a sprite where each character has its own color.
// runes and colors are the same layout. If colors is nil, fg/bg defaults are used.
func NewSpriteWithColors(width, height int, runes []rune, fg, bg Color, colors []Color) *Sprite {
	cells := make([]Cell, width*height)
	for i, r := range runes {
		c := Cell{Rune: r, FG: fg, BG: bg}
		if colors != nil && i < len(colors) {
			c.FG = colors[i]
		}
		cells[i] = c
	}
	return &Sprite{
		Width:  width,
		Height: height,
		Cells:  cells,
	}
}

// DrawAt draws the sprite onto a framebuffer centered at (x, y).
// x, y are in framebuffer column/row coordinates.
// Only non-transparent cells are drawn.
func (s *Sprite) DrawAt(fb *Framebuffer, x, y int) {
	ox := x - s.Width/2
	oy := y - s.Height/2

	for sy := 0; sy < s.Height; sy++ {
		for sx := 0; sx < s.Width; sx++ {
			cell := s.Cells[sy*s.Width+sx]
			if cell.Transparent {
				continue
			}
			fb.Set(oy+sy, ox+sx, cell)
		}
	}
}

// DrawAtTopLeft draws the sprite with top-left corner at (x, y).
func (s *Sprite) DrawAtTopLeft(fb *Framebuffer, x, y int) {
	for sy := 0; sy < s.Height; sy++ {
		for sx := 0; sx < s.Width; sx++ {
			cell := s.Cells[sy*s.Width+sx]
			if cell.Transparent {
				continue
			}
			fb.Set(y+sy, x+sx, cell)
		}
	}
}

func splitLines(s string) []string {
	var lines []string
	current := ""
	for _, ch := range s {
		if ch == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(ch)
		}
	}
	if current != "" || len(s) == 0 || s[len(s)-1] == '\n' {
		lines = append(lines, current)
	}
	return lines
}
