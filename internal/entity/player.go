package entity

import (
	"math"
	"overflow/internal/render"
)

// Player is the player entity.
type Player struct {
	BaseEntity
	Width, Height       int
	Speed               float64
	HP                  int
	MaxHP               int
	ShootCooldown       float64
	shootTimer          float64
	InvincibleTimer     float64
	InvincibleDuration  float64
	HitFlash            bool
	Score               int
	Firing              bool

	// Sprites (injected)
	SpriteNormal *render.Sprite
	SpriteHit    *render.Sprite
}

// NewPlayer creates a new player at the given position.
func NewPlayer(x, y float64) *Player {
	return &Player{
		BaseEntity:         BaseEntity{X: x, Y: y, Alive: true},
		Width:              3,
		Height:             3,
		Speed:              180,
		HP:                 3,
		MaxHP:              3,
		ShootCooldown:      0.12,
		shootTimer:         0,
		InvincibleDuration: 1.5,
	}
}

// Bounds returns the player's bounding box (center-based).
func (p *Player) Bounds() (x, y, w, h int) {
	return int(p.X), int(p.Y), p.Width, p.Height
}

// Update updates the player state.
func (p *Player) Update(dt float64) {
	if p.shootTimer > 0 {
		p.shootTimer -= dt
	}
	if p.InvincibleTimer > 0 {
		p.InvincibleTimer -= dt
		p.HitFlash = math.Mod(float64(int(p.InvincibleTimer*10)), 2) < 1
	} else {
		p.HitFlash = false
	}
}

// Move moves the player by (dx, dy), clamped to world bounds.
// dx, dy are direction (-1, 0, or 1).
func (p *Player) Move(dx, dy float64, dt float64, worldWidth, worldHeight int) {
	p.X += dx * p.Speed * dt
	p.Y += dy * p.Speed * dt

	halfW := p.Width / 2
	halfH := p.Height / 2
	if p.X < float64(halfW) {
		p.X = float64(halfW)
	}
	if p.X > float64(worldWidth-halfW-1) {
		p.X = float64(worldWidth - halfW - 1)
	}
	if p.Y < float64(4+halfH) {
		p.Y = float64(4 + halfH)
	}
	if p.Y > float64(worldHeight-halfH-1) {
		p.Y = float64(worldHeight - halfH - 1)
	}
}

// CanShoot returns true if the player can fire.
func (p *Player) CanShoot() bool {
	return p.shootTimer <= 0 && p.Alive
}

// Shoot resets timer and returns the bullet spawn position.
func (p *Player) Shoot() (x, y float64) {
	p.shootTimer = p.ShootCooldown
	return p.X, p.Y - float64(p.Height/2) - 1
}

// TakeDamage reduces HP. Returns false if dead.
func (p *Player) TakeDamage() bool {
	if p.InvincibleTimer > 0 || !p.Alive {
		return p.Alive
	}
	p.HP--
	p.InvincibleTimer = p.InvincibleDuration
	p.HitFlash = true
	if p.HP <= 0 {
		p.Alive = false
	}
	return p.Alive
}

// Draw renders the player to the framebuffer.
func (p *Player) Draw(fb *render.Framebuffer) {
	if !p.Alive {
		return
	}
	if p.InvincibleTimer > 0 && p.HitFlash {
		return // blink off
	}

	sprite := p.SpriteNormal
	if p.InvincibleTimer > 0 && p.SpriteHit != nil {
		sprite = p.SpriteHit
	}
	if sprite != nil {
		sprite.DrawAt(fb, int(p.X), int(p.Y))
	}
}
