package entity

import (
	"math"
	"overflow/internal/render"
)

// BulletSpawn describes a bullet to be spawned.
type BulletSpawn struct {
	X, Y              float64
	VX, VY            float64
	IsEnemy           bool
	EnemyBulletSprite *render.Sprite
	PlayerBulletSprite *render.Sprite
}

// Bullet is a projectile entity.
type Bullet struct {
	BaseEntity
	VX, VY      float64
	IsEnemy     bool
	Damage      int
	Lifetime    float64
	maxLifetime float64
	WorldWidth  int
	WorldHeight int
	MinY        int

	// Sprites (injected)
	Sprite *render.Sprite
}

// NewPlayerBullet creates a player's bullet.
func NewPlayerBullet(x, y, vy float64, worldWidth, worldHeight int) *Bullet {
	return &Bullet{
		BaseEntity:  BaseEntity{X: x, Y: y, Alive: true},
		VX:          0,
		VY:          vy,
		IsEnemy:     false,
		Damage:      1,
		maxLifetime: 2.0,
		Lifetime:    2.0,
		WorldWidth:  worldWidth,
		WorldHeight: worldHeight,
		MinY:        4,
	}
}

// NewEnemyBullet creates an enemy's bullet.
func NewEnemyBullet(x, y, vx, vy float64, worldWidth, worldHeight int) *Bullet {
	return &Bullet{
		BaseEntity:  BaseEntity{X: x, Y: y, Alive: true},
		VX:          vx,
		VY:          vy,
		IsEnemy:     true,
		Damage:      1,
		maxLifetime: 5.0,
		Lifetime:    5.0,
		WorldWidth:  worldWidth,
		WorldHeight: worldHeight,
		MinY:        4,
	}
}

// Bounds returns the bullet's bounding box (1x1).
func (b *Bullet) Bounds() (x, y, w, h int) {
	return int(b.X), int(b.Y), 1, 1
}

// Update updates the bullet position and lifetime.
func (b *Bullet) Update(dt float64) {
	if !b.Alive {
		return
	}

	b.X += b.VX * dt
	b.Y += b.VY * dt
	b.Lifetime -= dt

	if b.Lifetime <= 0 ||
		b.X < 0 || b.X >= float64(b.WorldWidth) ||
		b.Y < float64(b.MinY) || b.Y >= float64(b.WorldHeight) {
		b.Alive = false
	}
}

// Draw renders the bullet to the framebuffer.
func (b *Bullet) Draw(fb *render.Framebuffer) {
	if !b.Alive {
		return
	}

	if b.Sprite != nil {
		b.Sprite.DrawAt(fb, int(math.Round(b.X)), int(math.Round(b.Y)))
		return
	}

	// Fallback: render as a colored dot
	x := int(math.Round(b.X))
	y := int(math.Round(b.Y))
	var cell render.Cell
	if b.IsEnemy {
		cell = render.Cell{Rune: '•', FG: render.ColorRed, BG: render.ColorBlack}
	} else {
		cell = render.Cell{Rune: '∙', FG: render.ColorBrightYellow, BG: render.ColorBlack}
	}
	fb.Set(y, x, cell)
}

// SpawnBulletFromDef creates a bullet from a BulletSpawn definition.
func SpawnBulletFromDef(spawn BulletSpawn, worldWidth, worldHeight int) *Bullet {
	var b *Bullet
	if spawn.IsEnemy {
		b = NewEnemyBullet(spawn.X, spawn.Y, spawn.VX, spawn.VY, worldWidth, worldHeight)
		if spawn.EnemyBulletSprite != nil {
			b.Sprite = spawn.EnemyBulletSprite
		}
	} else {
		b = NewPlayerBullet(spawn.X, spawn.Y, spawn.VY, worldWidth, worldHeight)
		if spawn.PlayerBulletSprite != nil {
			b.Sprite = spawn.PlayerBulletSprite
		}
	}
	return b
}
