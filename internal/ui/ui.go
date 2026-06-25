package ui

import (
	"overflow/internal/render"
	"overflow/internal/world"
)

// UI handles drawing the HUD and other UI elements.
type UI struct {
	world  *world.World
	width  int
	fps    float64
	sinceLastFPSCalc float64
	frameCount        int
	displayFPS        int
	fpsUpdateInterval float64
}

// New creates a new UI renderer.
func New(w *world.World, width int) *UI {
	return &UI{
		world:             w,
		width:             width,
		fpsUpdateInterval: 0.5,
	}
}

// SetFPS updates the FPS counter.
func (ui *UI) SetFPS(fps float64, dt float64) {
	ui.frameCount++
	ui.sinceLastFPSCalc += dt
	if ui.sinceLastFPSCalc >= ui.fpsUpdateInterval {
		ui.displayFPS = int(fps + 0.5)
		ui.sinceLastFPSCalc = 0
		ui.frameCount = 0
	}
}

// Draw renders the HUD to the framebuffer.
func (ui *UI) Draw(fb *render.Framebuffer) {
	ui.drawHUDBorder(fb)
	ui.drawHP(fb)
	ui.drawScore(fb)
	ui.drawWave(fb)
	ui.drawFPS(fb)
	ui.drawEnemyCount(fb)
}

func (ui *UI) drawHUDBorder(fb *render.Framebuffer) {
	w := ui.width
	if w > fb.Width {
		w = fb.Width
	}
	hudHeight := 4

	// Fill HUD background
	for r := 0; r < hudHeight; r++ {
		for c := 0; c < w; c++ {
			fb.Set(r, c, render.Cell{Rune: ' ', FG: render.ColorHudText, BG: render.ColorHudBG})
		}
	}

	// Draw border using box-drawing characters
	// Top border
	fb.Set(0, 0, render.Cell{Rune: '╔', FG: render.ColorHudBorder, BG: render.ColorHudBG})
	for c := 1; c < w-1; c++ {
		fb.Set(0, c, render.Cell{Rune: '═', FG: render.ColorHudBorder, BG: render.ColorHudBG})
	}
	fb.Set(0, w-1, render.Cell{Rune: '╗', FG: render.ColorHudBorder, BG: render.ColorHudBG})

	// Left/right borders
	for r := 1; r < hudHeight; r++ {
		fb.Set(r, 0, render.Cell{Rune: '║', FG: render.ColorHudBorder, BG: render.ColorHudBG})
		fb.Set(r, w-1, render.Cell{Rune: '║', FG: render.ColorHudBorder, BG: render.ColorHudBG})
	}

	// Bottom border
	fb.Set(hudHeight-1, 0, render.Cell{Rune: '╚', FG: render.ColorHudBorder, BG: render.ColorHudBG})
	for c := 1; c < w-1; c++ {
		fb.Set(hudHeight-1, c, render.Cell{Rune: '═', FG: render.ColorHudBorder, BG: render.ColorHudBG})
	}
	fb.Set(hudHeight-1, w-1, render.Cell{Rune: '╝', FG: render.ColorHudBorder, BG: render.ColorHudBG})
}

func (ui *UI) drawHP(fb *render.Framebuffer) {
	player := ui.world.Player
	if player == nil {
		return
	}

	// HP label
	fb.DrawString(1, 2, "HP", render.ColorHudText, render.ColorHudBG)

	// HP hearts
	x := 5
	for i := 0; i < player.MaxHP; i++ {
		if i < player.HP {
			fb.Set(1, x, render.Cell{Rune: '♥', FG: render.ColorHudHP, BG: render.ColorHudBG})
		} else {
			fb.Set(1, x, render.Cell{Rune: '♡', FG: render.ColorHudHPLost, BG: render.ColorHudBG})
		}
		x += 2
	}

	// Health bar (alternative to hearts)
	barWidth := 10
	barStart := x + 1
	hpFrac := float64(player.HP) / float64(player.MaxHP)
	hpFilled := int(hpFrac * float64(barWidth))
	if hpFilled > barWidth {
		hpFilled = barWidth
	}

	fb.Set(1, barStart-1, render.Cell{Rune: '▐', FG: render.ColorHudBorder, BG: render.ColorHudBG})

	for i := 0; i < barWidth; i++ {
		if i < hpFilled {
			fb.Set(1, barStart+i, render.Cell{Rune: '█', FG: render.ColorHudHP, BG: render.ColorHudBG})
		} else {
			fb.Set(1, barStart+i, render.Cell{Rune: '░', FG: render.ColorHudHPLost, BG: render.ColorHudBG})
		}
	}

	fb.Set(1, barStart+barWidth, render.Cell{Rune: '▌', FG: render.ColorHudBorder, BG: render.ColorHudBG})
}

func (ui *UI) drawScore(fb *render.Framebuffer) {
	scoreStr := itoa(ui.world.Score)
	// Pad to 8 digits
	for len(scoreStr) < 8 {
		scoreStr = "0" + scoreStr
	}

	label := "SCORE"
	text := label + " " + scoreStr
	x := ui.width - len(text) - 2
	fb.DrawString(1, x, text, render.ColorHudScore, render.ColorHudBG)
}

func (ui *UI) drawWave(fb *render.Framebuffer) {
	waveStr := ui.world.GetWaveText()
	text := "WAVE " + waveStr

	// Center in row 2
	x := (ui.width - len(text)) / 2
	fb.DrawString(2, x, text, render.ColorHudWave, render.ColorHudBG)

	// New wave announcement
	w := ui.world
	if w.Wave <= 1 {
		return
	}
	if w.AreEnemiesAlive() {
		return
	}

	// Show wave title briefly
	if w.GetEnemyCount() == 0 {
		title := "WAVE " + waveStr + " COMPLETE"
		tx := (ui.width - len(title)) / 2
		fb.DrawString(3, tx, title, render.ColorGreen, render.ColorHudBG)
	}
}

func (ui *UI) drawFPS(fb *render.Framebuffer) {
	fpsStr := itoa(ui.displayFPS)
	text := "FPS " + fpsStr
	fb.DrawString(2, 2, text, render.ColorMidGray, render.ColorHudBG)
}

func (ui *UI) drawEnemyCount(fb *render.Framebuffer) {
	count := ui.world.GetEnemyCount()
	enemyStr := "ENEMIES " + itoa(count)
	x := ui.width - len(enemyStr) - 2
	fb.DrawString(2, x, enemyStr, render.ColorRed, render.ColorHudBG)
}

// itoa converts int to string without strconv.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	negative := false
	if n < 0 {
		negative = true
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if negative {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
