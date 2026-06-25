package world

import (
	"math/rand"
	"overflow/internal/entity"
	"overflow/internal/render"
)

// GameState represents the overall game state.
type GameState int

const (
	StateTitle GameState = iota
	StatePlaying
	StatePaused
	StateGameOver
)

// World manages all game state and entities.
type World struct {
	Entities    *entity.EntityManager
	Player      *entity.Player
	State       GameState

	// World dimensions
	Width       int
	Height      int

	// Wave system
	Wave              int
	enemiesLeft       int
	enemiesMax        int
	spawnTimer        float64
	spawnInterval     float64
	waveDelay         float64
	waveCompleteTimer float64

	// Scoring
	Score       int
	HighScore   int

	// Sprites (injected)
	PlayerSpriteNormal  *render.Sprite
	PlayerSpriteHit     *render.Sprite
	EnemyBasicSprite    *render.Sprite
	EnemyFastSprite     *render.Sprite
	EnemyTankSprite     *render.Sprite
	EnemyBossSprite     *render.Sprite
	EnemyBulletSprite   *render.Sprite
	PlayerBulletSprite  *render.Sprite
	Explosion1Sprite    *render.Sprite
	Explosion2Sprite    *render.Sprite
	Explosion3Sprite    *render.Sprite
	ParticleDotSprite   *render.Sprite
	ParticleStarSprite  *render.Sprite
	ParticleDiamondSprite *render.Sprite
}

// NewWorld creates a new game world.
func NewWorld(width, height int) *World {
	w := &World{
		Entities:      entity.NewEntityManager(),
		Width:         width,
		Height:        height,
		State:         StateTitle,
		Wave:          0,
		spawnInterval: 4.0,
		HighScore:     0,
	}
	w.resetGame()
	w.State = StateTitle
	return w
}

// Resize handles terminal resize events.
func (w *World) Resize(width, height int) {
	w.Width = width
	w.Height = height
}

// resetGame resets game state for a new game.
func (w *World) resetGame() {
	w.Entities.RemoveAll()

	w.Player = entity.NewPlayer(float64(w.Width)/2, float64(w.Height)-3)
	w.Player.SpriteNormal = w.PlayerSpriteNormal
	w.Player.SpriteHit = w.PlayerSpriteHit

	w.Entities.Add(w.Player)

	w.Wave = 0
	w.Score = 0
	w.enemiesLeft = 0
	w.enemiesMax = 0
	w.spawnTimer = 0
	w.spawnInterval = 3.0
	w.waveDelay = 0
	w.waveCompleteTimer = 0
}

// Reset starts a new game from scratch.
func (w *World) Reset() {
	w.resetGame()
	w.State = StatePlaying
}

// StartGame transitions from title screen to gameplay.
func (w *World) StartGame() {
	w.resetGame()
	w.State = StatePlaying
}

// Update updates the entire world for one timestep.
func (w *World) Update(dt float64) {
	if w.State == StatePaused || w.State == StateGameOver || w.State == StateTitle {
		return
	}

	dtClamped := dt
	if dtClamped > 0.033 {
		dtClamped = 0.033
	}

	// Update all entities
	w.Entities.Update(dtClamped)
	w.Entities.GarbageCollect()

	// Handle wave system
	w.updateWaveSystem(dtClamped)

	// Handle enemy fire
	w.handleEnemyFire()

	// Handle collision
	w.handleCollisions()
}

// updateWaveSystem manages enemy waves and spawning.
func (w *World) updateWaveSystem(dt float64) {
	aliveEnemies := 0
	for _, e := range w.Entities.AllAlive() {
		if _, ok := e.(*entity.Enemy); ok {
			aliveEnemies++
		}
	}

	// Transition to next wave when all enemies are gone
	if aliveEnemies == 0 && w.enemiesLeft == 0 && w.State == StatePlaying {
		w.waveDelay += dt
		w.waveCompleteTimer += dt
		if w.waveDelay >= 2.0 {
			w.nextWave()
		}
		return
	}

	// Spawn enemies in bursts
	if w.enemiesLeft > 0 {
		w.spawnTimer -= dt
		if w.spawnTimer <= 0 {
			w.spawnTimer = w.spawnInterval
			// Spawn multiple enemies per tick for density
			burstSize := 1
			if w.Wave >= 2 {
				burstSize = 2
			}
			if w.Wave >= 4 {
				burstSize = 3
			}
			if w.Wave >= 7 {
				burstSize = 4
			}
			if w.Wave >= 10 {
				burstSize = 5
			}
			for i := 0; i < burstSize && w.enemiesLeft > 0; i++ {
				w.spawnEnemy()
				w.enemiesLeft--
			}
		}
	}
}

// nextWave advances to the next wave.
func (w *World) nextWave() {
	w.Wave++

	bossWave := w.Wave%5 == 0

	if bossWave {
		w.enemiesMax = 1
		w.enemiesLeft = 1
		w.spawnInterval = 2.0
	} else {
		w.enemiesMax = 4 + w.Wave*2
		if w.enemiesMax > 30 {
			w.enemiesMax = 30
		}
		w.enemiesLeft = w.enemiesMax
		// Slower spawns at early waves, faster as game progresses
		if w.Wave <= 2 {
			w.spawnInterval = 4.0
		} else {
			w.spawnInterval = 4.0 / float64(w.Wave-1)
		}
		if w.spawnInterval < 0.5 {
			w.spawnInterval = 0.5
		}
	}

	w.spawnTimer = 2.0
	w.waveDelay = 0
	w.waveCompleteTimer = 0
}

// spawnEnemy creates a new enemy based on current wave.
func (w *World) spawnEnemy() {
	if w.Width < 30 {
		return
	}
	x := float64(10 + rand.Intn(w.Width-20))
	y := float64(4 + 2 + rand.Intn(4))

	var etype entity.EnemyType
	var pattern entity.EnemyPattern
	var sprite *render.Sprite

	// Boss wave
	if w.Wave%5 == 0 {
		etype = entity.EnemyBoss
		pattern = entity.PatternBurst
		sprite = w.EnemyBossSprite
	} else {
		r := rand.Float64()
		switch {
		case r < 0.50:
			etype = entity.EnemyBasic
			pattern = entity.PatternSingle
			if w.Wave > 3 {
				pattern = entity.PatternSpread
			}
			sprite = w.EnemyBasicSprite
		case r < 0.80:
			etype = entity.EnemyFast
			pattern = entity.PatternSingle
			if w.Wave > 5 {
				pattern = entity.PatternSpread
			}
			sprite = w.EnemyFastSprite
		default:
			etype = entity.EnemyTank
			pattern = entity.PatternSpread
			if w.Wave > 7 {
				pattern = entity.PatternCross
			}
			sprite = w.EnemyTankSprite
		}
	}

	enemy := entity.NewEnemy(x, y, etype, pattern, w.Wave)
	enemy.Sprite = sprite
	w.Entities.Add(enemy)
}

// handleEnemyFire processes enemy bullet spawning.
func (w *World) handleEnemyFire() {
	for _, e := range w.Entities.AllAlive() {
		enemy, ok := e.(*entity.Enemy)
		if !ok {
			continue
		}

		spawns := enemy.GetBulletSpawns()
		for _, sp := range spawns {
			sp.IsEnemy = true
			sp.EnemyBulletSprite = w.EnemyBulletSprite
			bullet := entity.SpawnBulletFromDef(sp, w.Width, w.Height)
			w.Entities.Add(bullet)
		}
	}
}

// handleCollisions checks for bullet-enemy and bullet-player collisions.
func (w *World) handleCollisions() {
	for _, e := range w.Entities.AllAlive() {
		bullet, ok := e.(*entity.Bullet)
		if !ok || !bullet.IsAlive() {
			continue
		}

		bx, by, bw, bh := bullet.Bounds()

		if !bullet.IsEnemy {
			// Player bullet hitting enemies
			for _, target := range w.Entities.AllAlive() {
				enemy, ok := target.(*entity.Enemy)
				if !ok || !enemy.IsAlive() {
					continue
				}

				ex, ey, ew, eh := enemy.Bounds()
				if aabbCollide(bx, by, bw, bh, ex, ey, ew, eh) {
					// Hit!
					bullet.Alive = false
					if enemy.TakeDamage(bullet.Damage) {
						// Enemy killed
						w.Score += enemy.ScoreValue
						w.spawnExplosion(enemy.X, enemy.Y)
						w.spawnFloatingScore(enemy.X, enemy.Y-1, enemy.ScoreValue)
					}
					break
				}
			}
		} else {
			// Enemy bullet hitting player
			if w.Player.IsAlive() {
				px, py, pw, ph := w.Player.Bounds()
				if aabbCollide(bx, by, bw, bh, px, py, pw, ph) {
					bullet.Alive = false
					wasAlive := w.Player.TakeDamage()
					w.spawnExplosion(w.Player.X, w.Player.Y)
					if !wasAlive {
						w.State = StateGameOver
						if w.Score > w.HighScore {
							w.HighScore = w.Score
						}
						// Final explosion
						for i := 0; i < 5; i++ {
							w.spawnExplosion(
								w.Player.X + rand.Float64()*4 - 2,
								w.Player.Y + rand.Float64()*4 - 2,
							)
						}
					}
				}
			}
		}
	}
}

// spawnExplosion creates explosion particles at the given position.
func (w *World) spawnExplosion(x, y float64) {
	particles := entity.NewExplosion(x, y)
	for _, p := range particles {
		w.Entities.Add(p)
	}
}

// spawnFloatingScore creates floating score text.
func (w *World) spawnFloatingScore(x, y float64, score int) {
	text := itoa(score)
	particles := entity.NewFloatingText(x-float64(len(text))/2, y, text)
	for _, p := range particles {
		w.Entities.Add(p)
	}
}

// PlayerShoot creates a player bullet.
func (w *World) PlayerShoot() {
	if !w.Player.CanShoot() || w.State != StatePlaying {
		return
	}
	x, y := w.Player.Shoot()
	bullet := entity.NewPlayerBullet(x, y, -300, w.Width, w.Height)
	bullet.Sprite = w.PlayerBulletSprite
	w.Entities.Add(bullet)
}

// PlayerMove handles player movement.
func (w *World) PlayerMove(dx, dy float64, dt float64) {
	if w.State != StatePlaying || !w.Player.IsAlive() {
		return
	}
	w.Player.Move(dx, dy, dt, w.Width, w.Height)
}

// Draw renders the entire world to the framebuffer.
func (w *World) Draw(fb *render.Framebuffer) {
	if w.State == StateGameOver || w.State == StateTitle {
		return
	}

	// Draw game area background
	fb.FillRect(4, 0, fb.Width, fb.Height-4, render.Cell{Rune: ' ', FG: render.ColorWhite, BG: render.ColorDarkBlue})

	// Draw all entities
	w.Entities.Draw(fb)
}

// DrawTitle renders the retro title screen.
// Designed to fit in a standard 80×24 terminal (22 rows between borders).
func (w *World) DrawTitle(fb *render.Framebuffer) {
	cx := fb.Width / 2

	// Fill entire screen with dark blue
	for r := 0; r < fb.Height; r++ {
		for c := 0; c < fb.Width; c++ {
			fb.Set(r, c, render.Cell{Rune: ' ', FG: render.ColorWhite, BG: render.ColorDarkBlue})
		}
	}

	// Draw outer border
	for c := 0; c < fb.Width; c++ {
		fb.Set(0, c, render.Cell{Rune: '═', FG: render.ColorHudBorder, BG: render.ColorDarkBlue})
		fb.Set(fb.Height-1, c, render.Cell{Rune: '═', FG: render.ColorHudBorder, BG: render.ColorDarkBlue})
	}
	for r := 0; r < fb.Height; r++ {
		fb.Set(r, 0, render.Cell{Rune: '║', FG: render.ColorHudBorder, BG: render.ColorDarkBlue})
		fb.Set(r, fb.Width-1, render.Cell{Rune: '║', FG: render.ColorHudBorder, BG: render.ColorDarkBlue})
	}
	fb.Set(0, 0, render.Cell{Rune: '╔', FG: render.ColorHudBorder, BG: render.ColorDarkBlue})
	fb.Set(0, fb.Width-1, render.Cell{Rune: '╗', FG: render.ColorHudBorder, BG: render.ColorDarkBlue})
	fb.Set(fb.Height-1, 0, render.Cell{Rune: '╚', FG: render.ColorHudBorder, BG: render.ColorDarkBlue})
	fb.Set(fb.Height-1, fb.Width-1, render.Cell{Rune: '╝', FG: render.ColorHudBorder, BG: render.ColorDarkBlue})

	row := 3

	// — OVERFLOW banner —
	title := "OVERFLOW"
	bannerW := len(title) + 4
	bannerX := cx - bannerW/2
	fb.DrawString(row, bannerX, "▄"+repeat("▄", bannerW-2)+"▄", render.ColorCyan, render.ColorDarkBlue)
	row++
	fb.Set(row, bannerX, render.Cell{Rune: '█', FG: render.ColorCyan, BG: render.ColorDarkBlue})
	fb.DrawString(row, bannerX+2, title, render.ColorBrightCyan, render.ColorDarkBlue)
	fb.Set(row, bannerX+1, render.Cell{Rune: ' ', FG: render.ColorCyan, BG: render.ColorDarkBlue})
	fb.Set(row, bannerX+bannerW-2, render.Cell{Rune: ' ', FG: render.ColorCyan, BG: render.ColorDarkBlue})
	fb.Set(row, bannerX+bannerW-1, render.Cell{Rune: '█', FG: render.ColorCyan, BG: render.ColorDarkBlue})
	row++
	fb.DrawString(row, bannerX, "▀"+repeat("▀", bannerW-2)+"▀", render.ColorCyan, render.ColorDarkBlue)
	row++

	// Subtitle
	sub := "BULLET HELL"
	subX := cx - len(sub)/2
	fb.DrawString(row, subX, sub, render.ColorOrange, render.ColorDarkBlue)
	row++

	// Separator
	sep := repeat("─", 24)
	sepX := cx - len(sep)/2
	fb.DrawString(row, sepX, sep, render.ColorMidGray, render.ColorDarkBlue)
	row++

	// Controls
	controls := []struct {
		key  string
		action string
	}{
		{"WASD / ARROWS", "MOVE"},
		{"SPACE", "SHOOT"},
		{"P", "PAUSE"},
		{"ESC", "QUIT"},
	}
	for _, ctrl := range controls {
		line := ctrl.key + repeat(" ", 16-len(ctrl.key)) + ctrl.action
		lx := cx - len(line)/2
		fb.DrawString(row, lx, line, render.ColorWhite, render.ColorDarkBlue)
		row++
	}

	// Separator
	fb.DrawString(row, sepX, sep, render.ColorMidGray, render.ColorDarkBlue)
	row++

	// PRESS SPACE prompt (clean pill-shaped banner)
	row++
	prompt := "PRESS SPACE TO START"
	bw := len(prompt) + 8
	bx := cx - bw/2
	fb.DrawString(row, bx, "▄"+repeat("▄", bw-2)+"▄", render.ColorGreen, render.ColorDarkBlue)
	row++
	fb.Set(row, bx, render.Cell{Rune: '█', FG: render.ColorGreen, BG: render.ColorDarkBlue})
	fb.DrawString(row, bx+4, prompt, render.ColorBrightYellow, render.ColorDarkBlue)
	fb.Set(row, bx+bw-1, render.Cell{Rune: '█', FG: render.ColorGreen, BG: render.ColorDarkBlue})
	fb.Set(row, bx+1, render.Cell{Rune: ' ', FG: render.ColorGreen, BG: render.ColorDarkBlue})
	fb.Set(row, bx+2, render.Cell{Rune: '▌', FG: render.ColorGreen, BG: render.ColorDarkBlue})
	fb.Set(row, bx+3, render.Cell{Rune: ' ', FG: render.ColorGreen, BG: render.ColorDarkBlue})
	fb.Set(row, bx+bw-2, render.Cell{Rune: ' ', FG: render.ColorGreen, BG: render.ColorDarkBlue})
	fb.Set(row, bx+bw-3, render.Cell{Rune: '▐', FG: render.ColorGreen, BG: render.ColorDarkBlue})
	fb.Set(row, bx+bw-4, render.Cell{Rune: ' ', FG: render.ColorGreen, BG: render.ColorDarkBlue})
	row++
	fb.DrawString(row, bx, "▀"+repeat("▀", bw-2)+"▀", render.ColorGreen, render.ColorDarkBlue)
}

// repeat returns a string with ch repeated n times.
func repeat(ch string, n int) string {
	if n <= 0 {
		return ""
	}
	b := make([]byte, n*len([]byte(ch)))
	for i := 0; i < n; i++ {
		copy(b[i*len([]byte(ch)):], ch)
	}
	return string(b)
}

// aabbCollide checks axis-aligned bounding box collision.
// Boxes are centered at (cx, cy) with dimensions (w, h).
func aabbCollide(ax, ay, aw, ah int, bx, by, bw, bh int) bool {
	dx := ax - bx
	if dx < 0 {
		dx = -dx
	}
	dy := ay - by
	if dy < 0 {
		dy = -dy
	}
	return dx <= (aw+bw)/2 && dy <= (ah+bh)/2
}

// Simple integer to string conversion (avoids strconv import).
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

// DrawGameOver renders the game over overlay.
func (w *World) DrawGameOver(fb *render.Framebuffer) {
	// Semi-transparent overlay effect
	centerY := fb.Height / 2
	centerX := fb.Width / 2

	// Title
	title := "GAME OVER"
	titleX := centerX - len(title)/2
	fb.DrawString(centerY-2, titleX, title, render.ColorRed, render.ColorBlack)

	// Score
	scoreText := "SCORE: " + itoa(w.Score)
	scoreX := centerX - len(scoreText)/2
	fb.DrawString(centerY, scoreX, scoreText, render.ColorYellow, render.ColorBlack)

	// High score
	hiText := "HIGH SCORE: " + itoa(w.HighScore)
	hiX := centerX - len(hiText)/2
	fb.DrawString(centerY+1, hiX, hiText, render.ColorWhite, render.ColorBlack)

	// Instructions
	instrText := "PRESS 'R' TO RESTART  |  'Q' TO QUIT"
	instrX := centerX - len(instrText)/2
	fb.DrawString(centerY+3, instrX, instrText, render.ColorMidGray, render.ColorBlack)
}

// DrawPause renders the pause overlay.
func (w *World) DrawPause(fb *render.Framebuffer) {
	centerY := fb.Height / 2
	centerX := fb.Width / 2

	text := "PAUSED"
	textX := centerX - len(text)/2
	fb.DrawString(centerY, textX, text, render.ColorCyan, render.ColorBlack)

	instr := "PRESS 'P' TO CONTINUE"
	instrX := centerX - len(instr)/2
	fb.DrawString(centerY+1, instrX, instr, render.ColorMidGray, render.ColorBlack)
}

// GetWaveText returns the wave number as a formatted string.
func (w *World) GetWaveText() string {
	s := itoa(w.Wave)
	if w.Wave < 10 {
		s = "0" + s
	}
	if w.Wave < 100 {
		s = "0" + s
	}
	return s
}

// AreEnemiesAlive returns true if any enemies are alive (including those being spawned).
func (w *World) AreEnemiesAlive() bool {
	for _, e := range w.Entities.AllAlive() {
		if _, ok := e.(*entity.Enemy); ok {
			return true
		}
	}
	return w.enemiesLeft > 0
}

// GetEnemyCount returns the number of alive enemies.
func (w *World) GetEnemyCount() int {
	count := 0
	for _, e := range w.Entities.AllAlive() {
		if _, ok := e.(*entity.Enemy); ok {
			count++
		}
	}
	return count
}

// HandleInput processes a key press and returns whether to quit.
func (w *World) HandleInput(key entity.Key) bool {
	switch w.State {
	case StateTitle:
		switch key {
		case entity.KeyEsc:
			return true
		case entity.KeySpace:
			w.StartGame()
		}
	case StatePlaying:
		switch key {
		case entity.KeyEsc:
			return true // quit
		case entity.KeyP:
			w.State = StatePaused
		}
	case StatePaused:
		switch key {
		case entity.KeyEsc:
			return true
		case entity.KeyP:
			w.State = StatePlaying
		}
	case StateGameOver:
		switch key {
		case entity.KeyR:
			w.Reset()
		case entity.KeyQ, entity.KeyEsc:
			return true
		}
	}
	return false
}
