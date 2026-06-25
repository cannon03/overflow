package engine

import (
	"os"
	"os/signal"
	"syscall"
	"time"
	"unsafe"

	"overflow/internal/assets"
	"overflow/internal/entity"
	"overflow/internal/input"
	"overflow/internal/render"
	"overflow/internal/ui"
	"overflow/internal/world"
)

const (
	FixedDT   = 1.0 / 60.0
	MaxFrames = 5
)

// Engine manages the game loop and terminal.
type Engine struct {
	renderer *render.ANSIRenderer
	inputDev *input.Input
	world    *world.World
	ui       *ui.UI
	fb       *render.Framebuffer

	running bool

	// FPS tracking
	fpsAccum   float64
	fpsDisplay float64
	fpsTimer   float64

	// Terminal dimensions
	termWidth  int
	termHeight int
}

// New creates a new game engine.
func New() *Engine {
	return &Engine{
		renderer:   render.NewANSIRenderer(),
		inputDev:   input.New(),
		fpsDisplay: 60,
	}
}

// Run starts the game loop.
func (e *Engine) Run() {
	if err := e.setup(); err != nil {
		cleanupTerminal(e.renderer, e.inputDev)
		os.Exit(1)
	}
	defer cleanupTerminal(e.renderer, e.inputDev)

	e.running = true

	previous := time.Now()
	var accumulator float64

	sigwinch := make(chan os.Signal, 1)
	signal.Notify(sigwinch, syscall.SIGWINCH)

	for e.running {
		frameStart := time.Now()
		elapsed := frameStart.Sub(previous).Seconds()
		previous = frameStart

		accumulator += elapsed
		if accumulator > FixedDT*MaxFrames {
			accumulator = FixedDT * MaxFrames
		}

		select {
		case <-sigwinch:
			e.handleResize()
		default:
		}

		for accumulator >= FixedDT {
			e.processInput()
			e.world.Update(FixedDT)
			accumulator -= FixedDT
		}

		e.render()

		// Frame rate limiter: target 60 FPS
		frameElapsed := time.Since(frameStart)
		targetFrameTime := time.Second / 60
		if frameElapsed < targetFrameTime {
			time.Sleep(targetFrameTime - frameElapsed)
		}

		// FPS is the measured frame rate (before sleep)
		e.fpsAccum++
		e.fpsTimer += time.Since(frameStart).Seconds()
		if e.fpsTimer >= 0.5 {
			// Actual FPS = frames / time spent processing
			actualFPS := float64(e.fpsAccum) / e.fpsTimer
			// Clamp to max 60 FPS display
			if actualFPS > 61 {
				actualFPS = 60
			}
			e.fpsDisplay = actualFPS
			e.fpsAccum = 0
			e.fpsTimer = 0
		}
	}
}

func (e *Engine) setup() error {
	e.termWidth = 80
	e.termHeight = 24
	if w, h, err := getTermSize(); err == nil && w >= 40 && h >= 12 {
		e.termWidth = w
		e.termHeight = h
	}

	e.renderer.EnterAltScreen()

	if err := e.inputDev.EnableRawMode(); err != nil {
		return err
	}

	e.inputDev.Start()

	e.fb = render.NewFramebuffer(e.termWidth, e.termHeight)

	e.world = world.NewWorld(e.termWidth, e.termHeight)

	// Inject sprites into world (MUST be done before Reset)
	e.world.PlayerSpriteNormal = assets.Sprites.PlayerNormal
	e.world.PlayerSpriteHit = assets.Sprites.PlayerHit
	e.world.EnemyBasicSprite = assets.Sprites.EnemyBasic
	e.world.EnemyFastSprite = assets.Sprites.EnemyFast
	e.world.EnemyTankSprite = assets.Sprites.EnemyTank
	e.world.EnemyBossSprite = assets.Sprites.EnemyBoss
	e.world.EnemyBulletSprite = assets.Sprites.EnemyBullet
	e.world.PlayerBulletSprite = assets.Sprites.PlayerBullet
	e.world.Explosion1Sprite = assets.Sprites.Explosion1
	e.world.Explosion2Sprite = assets.Sprites.Explosion2
	e.world.Explosion3Sprite = assets.Sprites.Explosion3
	e.world.ParticleDotSprite = assets.Sprites.ParticleDot
	e.world.ParticleStarSprite = assets.Sprites.ParticleStar
	e.world.ParticleDiamondSprite = assets.Sprites.ParticleDiamond

	// Re-initialize player with sprites now injected
	e.world.Reset()

	e.ui = ui.New(e.world, e.termWidth)

	// Initial full render
	e.fb.Clear()
	e.world.Draw(e.fb)
	e.ui.Draw(e.fb)
	e.renderer.FullRender(e.fb)

	return nil
}

// processInput drains all buffered keys and processes them.
func (e *Engine) processInput() {
	for {
		key := e.inputDev.GetKey()
		if key == entity.KeyNone {
			break
		}

		if e.world.HandleInput(key) {
			e.running = false
			return
		}

		if e.world.State != world.StatePlaying {
			continue
		}

		switch key {
		case entity.KeyW, entity.KeyUp:
			e.world.PlayerMove(0, -1, FixedDT)
		case entity.KeyS, entity.KeyDown:
			e.world.PlayerMove(0, 1, FixedDT)
		case entity.KeyA, entity.KeyLeft:
			e.world.PlayerMove(-1, 0, FixedDT)
		case entity.KeyD, entity.KeyRight:
			e.world.PlayerMove(1, 0, FixedDT)
		case entity.KeySpace:
			e.world.PlayerShoot()
		}
	}
}

func (e *Engine) render() {
	e.fb.Clear()

	switch e.world.State {
	case world.StateTitle:
		// Title screen — draw fullscreen retro title card
		e.world.DrawTitle(e.fb)

	default:
		// Normal game rendering
		e.ui.Draw(e.fb)
		e.world.Draw(e.fb)

		// Draw overlays
		switch e.world.State {
		case world.StateGameOver:
			e.world.DrawGameOver(e.fb)
		case world.StatePaused:
			e.world.DrawPause(e.fb)
		}

		e.ui.SetFPS(e.fpsDisplay, FixedDT)
	}

	e.renderer.Render(e.fb)
}

func (e *Engine) handleResize() {
	w, h, err := getTermSize()
	if err != nil || w < 40 || h < 12 {
		return
	}
	e.termWidth = w
	e.termHeight = h

	e.fb.Resize(w, h)
	e.world.Resize(w, h)
	e.fb.Clear()
	e.world.Draw(e.fb)
	e.ui.Draw(e.fb)
	e.renderer.FullRender(e.fb)
}

func cleanupTerminal(r *render.ANSIRenderer, in *input.Input) {
	in.Stop()
	in.Restore()
	r.ExitAltScreen()
}

// getTermSize detects terminal size using ioctl.
func getTermSize() (width, height int, err error) {
	var winsize struct {
		Row    uint16
		Col    uint16
		Xpixel uint16
		Ypixel uint16
	}
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		syscall.TIOCGWINSZ,
		uintptr(unsafe.Pointer(&winsize)),
	); errno != 0 {
		return 80, 24, errno
	}
	return int(winsize.Col), int(winsize.Row), nil
}
