//go:build js

package engine

import (
	"os"
	"time"

	"overflow/internal/assets"
	"overflow/internal/entity"
	"overflow/internal/render"
	"overflow/internal/ui"
	"overflow/internal/world"
)

// Run starts the WASM game loop. No signal handling, no raw terminal mode.
func (e *Engine) Run() {
	if err := e.setup(); err != nil {
		os.Exit(1)
	}

	e.running = true

	previous := time.Now()
	var accumulator float64

	for e.running {
		frameStart := time.Now()
		elapsed := frameStart.Sub(previous).Seconds()
		previous = frameStart

		accumulator += elapsed
		if accumulator > FixedDT*MaxFrames {
			accumulator = FixedDT * MaxFrames
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

		// FPS tracking
		e.fpsAccum++
		e.fpsTimer += time.Since(frameStart).Seconds()
		if e.fpsTimer >= 0.5 {
			actualFPS := float64(e.fpsAccum) / e.fpsTimer
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

	e.inputDev.Start()

	e.fb = render.NewFramebuffer(e.termWidth, e.termHeight)
	e.world = world.NewWorld(e.termWidth, e.termHeight)

	// Inject sprites
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

	e.world.Reset()
	e.ui = ui.New(e.world, e.termWidth)

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
		e.world.DrawTitle(e.fb)
	default:
		e.ui.Draw(e.fb)
		e.world.Draw(e.fb)
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
