package render

import (
	"bytes"
	"fmt"
	"os"
)

// ANSI escape sequences.
const (
	ansiHome       = "\033[H"
	ansiHideCursor = "\033[?25l"
	ansiShowCursor = "\033[?25h"
	ansiAltScreen  = "\033[?1049h"
	ansiExitAlt    = "\033[?1049l"
	ansiReset      = "\033[0m"
	ansiClearLine  = "\033[2K"
	ansiClearScreen = "\033[2J"
)

// ANSIRenderer renders a Framebuffer to the terminal using ANSI escape codes.
type ANSIRenderer struct {
	buf bytes.Buffer
}

// NewANSIRenderer creates a new ANSI renderer.
func NewANSIRenderer() *ANSIRenderer {
	return &ANSIRenderer{}
}

// EnterAltScreen switches to the alternate screen and hides the cursor.
func (r *ANSIRenderer) EnterAltScreen() {
	r.buf.Reset()
	r.buf.WriteString(ansiAltScreen)
	r.buf.WriteString(ansiHideCursor)
	r.buf.WriteString(ansiClearScreen)
	r.buf.WriteString(ansiHome)
	os.Stdout.Write(r.buf.Bytes())
}

// ExitAltScreen exits the alternate screen and shows the cursor.
func (r *ANSIRenderer) ExitAltScreen() {
	r.buf.Reset()
	r.buf.WriteString(ansiShowCursor)
	r.buf.WriteString(ansiExitAlt)
	os.Stdout.Write(r.buf.Bytes())
}

// Render draws only dirty rows from the framebuffer to the terminal.
// For each dirty row, it positions the cursor and outputs all cells sequentially.
// This produces a consistent, reliable output without per-cell positioning issues.
func (r *ANSIRenderer) Render(fb *Framebuffer) {
	r.buf.Reset()
	r.buf.WriteString(ansiHome)

	width := fb.Width
	height := fb.Height

	for row := 0; row < height; row++ {
		if !fb.dirtyRows[row] {
			continue
		}

		// Position cursor at start of this row
		fmt.Fprintf(&r.buf, "\033[%d;1H", row+1)

		// Output all cells in this row
		currentFG := Color{}
		currentBG := Color{}
		baseIdx := row * width

		for col := 0; col < width; col++ {
			cell := fb.cells[baseIdx+col]

			if cell.FG != currentFG {
				fmt.Fprintf(&r.buf, "\033[38;2;%d;%d;%dm", cell.FG.R, cell.FG.G, cell.FG.B)
				currentFG = cell.FG
			}
			if cell.BG != currentBG {
				fmt.Fprintf(&r.buf, "\033[48;2;%d;%d;%dm", cell.BG.R, cell.BG.G, cell.BG.B)
				currentBG = cell.BG
			}

			r.buf.WriteRune(cell.Rune)
		}

		// Clear the dirty flag for this row
		fb.dirtyRows[row] = false
	}

	// Reset colors and move cursor below the display area
	r.buf.WriteString(ansiReset)
	fmt.Fprintf(&r.buf, "\033[%d;1H", height+1)

	os.Stdout.Write(r.buf.Bytes())
}

// FullRender renders the entire framebuffer (used for initial draw and resize).
func (r *ANSIRenderer) FullRender(fb *Framebuffer) {
	// Mark all rows dirty and call Render
	for i := range fb.dirtyRows {
		fb.dirtyRows[i] = true
	}
	r.Render(fb)
}
