package render

// Framebuffer is a 2D grid of Cells with double-buffering and dirty-row tracking.
type Framebuffer struct {
	Width     int
	Height    int
	cells     []Cell
	dirtyRows []bool
}

// NewFramebuffer creates a new framebuffer of the given size.
func NewFramebuffer(width, height int) *Framebuffer {
	cells := make([]Cell, width*height)
	for i := range cells {
		cells[i] = DefaultCell
	}
	return &Framebuffer{
		Width:     width,
		Height:    height,
		cells:     cells,
		dirtyRows: make([]bool, height),
	}
}

// Resize recreates the framebuffer with new dimensions (clears it).
func (fb *Framebuffer) Resize(width, height int) {
	fb.Width = width
	fb.Height = height
	fb.cells = make([]Cell, width*height)
	fb.dirtyRows = make([]bool, height)
	for i := range fb.cells {
		fb.cells[i] = DefaultCell
	}
	for i := range fb.dirtyRows {
		fb.dirtyRows[i] = true
	}
}

// Clear fills the entire framebuffer with the default cell.
func (fb *Framebuffer) Clear() {
	for i := range fb.cells {
		fb.cells[i] = DefaultCell
	}
	for i := range fb.dirtyRows {
		fb.dirtyRows[i] = true
	}
}

// Set places a cell at the given position. Clips to bounds.
func (fb *Framebuffer) Set(row, col int, cell Cell) {
	if row < 0 || row >= fb.Height || col < 0 || col >= fb.Width {
		return
	}
	idx := row*fb.Width + col
	if fb.cells[idx] != cell {
		fb.cells[idx] = cell
		fb.dirtyRows[row] = true
	}
}

// Get returns the cell at the given position.
func (fb *Framebuffer) Get(row, col int) Cell {
	if row < 0 || row >= fb.Height || col < 0 || col >= fb.Width {
		return DefaultCell
	}
	return fb.cells[row*fb.Width+col]
}

// FillRect fills a rectangle with the given cell.
func (fb *Framebuffer) FillRect(row, col, width, height int, cell Cell) {
	for r := row; r < row+height; r++ {
		for c := col; c < col+width; c++ {
			fb.Set(r, c, cell)
		}
	}
}

// DrawString draws a string at the given position with the given color.
// Supports multi-byte UTF-8 characters (uses rune index, not byte index).
func (fb *Framebuffer) DrawString(row, col int, s string, fg, bg Color) {
	c := col
	for _, ch := range s {
		fb.Set(row, c, Cell{Rune: ch, FG: fg, BG: bg})
		c++
	}
}

// DrawStringClipped draws a string, clipped to framebuffer bounds.
// Supports multi-byte UTF-8 characters (uses rune index, not byte index).
func (fb *Framebuffer) DrawStringClipped(row, col int, s string, fg, bg Color) {
	c := col
	for _, ch := range s {
		if c >= 0 && c < fb.Width && row >= 0 && row < fb.Height {
			fb.Set(row, c, Cell{Rune: ch, FG: fg, BG: bg})
		}
		c++
	}
}

// IsDirty returns whether any row is dirty.
func (fb *Framebuffer) IsDirty() bool {
	for _, d := range fb.dirtyRows {
		if d {
			return true
		}
	}
	return false
}

// MarkAllDirty marks every row as dirty.
func (fb *Framebuffer) MarkAllDirty() {
	for i := range fb.dirtyRows {
		fb.dirtyRows[i] = true
	}
}

// MarkDirty marks a specific row as dirty.
func (fb *Framebuffer) MarkDirty(row int) {
	if row >= 0 && row < fb.Height {
		fb.dirtyRows[row] = true
	}
}

// IsRowDirty returns whether the given row is dirty.
func (fb *Framebuffer) IsRowDirty(row int) bool {
	if row < 0 || row >= fb.Height {
		return false
	}
	return fb.dirtyRows[row]
}

// ClearDirty clears the dirty flag for a specific row.
func (fb *Framebuffer) ClearDirty(row int) {
	if row >= 0 && row < fb.Height {
		fb.dirtyRows[row] = false
	}
}

// ClearAllDirty clears all dirty flags.
func (fb *Framebuffer) ClearAllDirty() {
	for i := range fb.dirtyRows {
		fb.dirtyRows[i] = false
	}
}

// Diff returns the dirty rectangle bounds (minRow, maxRow, minCol, maxCol).
// Returns (-1, -1, -1, -1) if nothing is dirty.
func (fb *Framebuffer) Diff() (minRow, maxRow, minCol, maxCol int) {
	minRow = fb.Height
	maxRow = -1
	minCol = fb.Width
	maxCol = -1

	for row := 0; row < fb.Height; row++ {
		if !fb.dirtyRows[row] {
			continue
		}
		if row < minRow {
			minRow = row
		}
		if row > maxRow {
			maxRow = row
		}
		for col := 0; col < fb.Width; col++ {
			if col < minCol {
				minCol = col
			}
			if col > maxCol {
				maxCol = col
			}
		}
	}

	if maxRow < 0 {
		return -1, -1, -1, -1
	}
	return
}
