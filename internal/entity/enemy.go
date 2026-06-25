package entity

import (
	"math"
	"math/rand"
	"overflow/internal/render"
)

// EnemyType defines the kind of enemy.
type EnemyType int

const (
	EnemyBasic EnemyType = iota
	EnemyFast
	EnemyTank
	EnemyBoss
)

// EnemyPattern defines the bullet pattern.
type EnemyPattern int

const (
	PatternNone EnemyPattern = iota
	PatternSingle
	PatternSpread
	PatternCross
	PatternSpiral
	PatternBurst
)

// Enemy is an enemy entity.
type Enemy struct {
	BaseEntity
	Width, Height int
	EnemyType     EnemyType
	Pattern       EnemyPattern
	HP            int
	MaxHP         int
	Speed         float64
	BulletSpeed   float64
	FireRate      float64
	fireTimer     float64
	ScoreValue    int

	spiralAngle float64
	spiralTimer float64
	spiralSpeed float64
	moveTimer   float64
	moveDirX    float64

	// Sprite (injected)
	Sprite *render.Sprite
}

// NewEnemy creates a new enemy.
func NewEnemy(x, y float64, etype EnemyType, pattern EnemyPattern, wave int) *Enemy {
	e := &Enemy{
		BaseEntity:  BaseEntity{X: x, Y: y, Alive: true},
		EnemyType:   etype,
		Pattern:     pattern,
		fireTimer:   rand.Float64() * 1.5,
		spiralSpeed: 3.0,
		moveDirX:    (rand.Float64() - 0.5) * 2,
	}

	switch etype {
	case EnemyBasic:
		e.Width = 4
		e.Height = 3
		e.HP = 1
		e.Speed = 12 + float64(wave)*3.5  // Wave 1: ~16, Wave 5: ~30, Wave 10: ~47
		e.BulletSpeed = 35 + float64(wave)*6  // Wave 1: 41, Wave 5: 65, Wave 10: 95
		e.FireRate = 3.5 - float64(wave)*0.12  // Wave 1: 3.4s, Wave 5: 2.9s, Wave 10: 2.3s
		e.ScoreValue = 100
		if e.FireRate < 0.8 {
			e.FireRate = 0.8
		}
	case EnemyFast:
		e.Width = 4
		e.Height = 3
		e.HP = 1
		e.Speed = 28 + float64(wave)*5  // Wave 1: 33, Wave 5: 53, Wave 10: 78
		e.BulletSpeed = 50 + float64(wave)*8  // Wave 1: 58, Wave 5: 90, Wave 10: 130
		e.FireRate = 3.0 - float64(wave)*0.15  // Wave 1: 2.9s, Wave 5: 2.3s, Wave 10: 1.5s
		e.ScoreValue = 150
		if e.FireRate < 0.5 {
			e.FireRate = 0.5
		}
	case EnemyTank:
		e.Width = 6
		e.Height = 3
		e.HP = 2 + wave/3
		e.Speed = 8 + float64(wave)*2  // Wave 1: 10, Wave 5: 18, Wave 10: 28
		e.BulletSpeed = 30 + float64(wave)*4  // Wave 1: 34, Wave 5: 50, Wave 10: 70
		e.FireRate = 4.0 - float64(wave)*0.15  // Wave 1: 3.9s, Wave 5: 3.3s, Wave 10: 2.5s
		e.ScoreValue = 250
		if e.FireRate < 1.0 {
			e.FireRate = 1.0
		}
	case EnemyBoss:
		e.Width = 8
		e.Height = 5
		e.HP = 5 + wave*3  // Wave 5: 20, Wave 10: 35
		e.Speed = 8 + float64(wave)*0.5  // Wave 5: ~11, Wave 10: 13
		e.BulletSpeed = 30 + float64(wave)*3  // Wave 5: 45, Wave 10: 60
		e.FireRate = 2.0 - float64(wave)*0.03  // Wave 5: 1.85s, Wave 10: 1.7s
		e.ScoreValue = 1000 + wave*200
		if e.FireRate < 0.5 {
			e.FireRate = 0.5
		}
		e.spiralSpeed = 2.0 + float64(wave)*0.2
	}

	e.MaxHP = e.HP
	e.moveTimer = rand.Float64() * 2
	return e
}

// Bounds returns the enemy's bounding box.
func (e *Enemy) Bounds() (x, y, w, h int) {
	return int(e.X), int(e.Y), e.Width, e.Height
}

// Update updates the enemy state.
func (e *Enemy) Update(dt float64) {
	if !e.Alive {
		return
	}

	e.moveTimer += dt

	switch e.EnemyType {
	case EnemyBasic:
		e.Y += e.Speed * dt
		e.X += math.Sin(e.moveTimer*1.5) * 30 * dt
	case EnemyFast:
		e.Y += e.Speed * dt
		e.X += math.Sin(e.moveTimer*3.0) * 60 * dt
	case EnemyTank:
		e.Y += e.Speed * dt
		e.X += math.Sin(e.moveTimer*0.8) * 20 * dt
	case EnemyBoss:
		e.Y += math.Sin(e.moveTimer*0.5) * 15 * dt
		e.X += math.Sin(e.moveTimer*0.7) * 40 * dt
	}

	// Kill enemies that go off-screen (past bottom)
	if e.Y > 100 {
		e.Alive = false
		return
	}

	e.fireTimer -= dt

	if e.Pattern == PatternSpiral {
		e.spiralTimer += dt
	}
}

// CanFire returns true if the enemy should fire this frame.
func (e *Enemy) CanFire() bool {
	return e.fireTimer <= 0 && e.Alive && e.Pattern != PatternNone
}

// GetBulletSpawns returns bullet spawns for the enemy's pattern.
func (e *Enemy) GetBulletSpawns() []BulletSpawn {
	if !e.CanFire() {
		return nil
	}
	// Reset fire timer AFTER confirming the enemy fires
	e.fireTimer = e.FireRate

	speed := e.BulletSpeed

	switch e.Pattern {
	case PatternSingle:
		return []BulletSpawn{
			{X: e.X, Y: e.Y + float64(e.Height/2), VX: 0, VY: speed},
		}
	case PatternSpread:
		count := 3 + rand.Intn(3)
		spawns := make([]BulletSpawn, count)
		spreadAngle := math.Pi / 3.0
		startAngle := math.Pi/2.0 - spreadAngle/2.0 // points downward
		for i := 0; i < count; i++ {
			angle := startAngle + spreadAngle*float64(i)/float64(count-1)
			spawns[i] = BulletSpawn{
				X: e.X, Y: e.Y + float64(e.Height/2),
				VX: math.Cos(angle) * speed,
				VY: math.Sin(angle) * speed,
			}
		}
		return spawns
	case PatternCross:
		return []BulletSpawn{
			{X: e.X, Y: e.Y + 1, VX: 0, VY: speed},
			{X: e.X, Y: e.Y - 1, VX: 0, VY: -speed},
			{X: e.X - 1, Y: e.Y, VX: -speed, VY: 0},
			{X: e.X + 1, Y: e.Y, VX: speed, VY: 0},
		}
	case PatternSpiral:
		e.spiralAngle += e.spiralSpeed * 0.1
		return []BulletSpawn{
			{X: e.X, Y: e.Y, VX: math.Cos(e.spiralAngle) * speed, VY: math.Sin(e.spiralAngle) * speed},
		}
	case PatternBurst:
		count := 8
		spawns := make([]BulletSpawn, count)
		for i := 0; i < count; i++ {
			angle := 2.0 * math.Pi * float64(i) / float64(count)
			spawns[i] = BulletSpawn{
				X: e.X, Y: e.Y,
				VX: math.Cos(angle) * speed * 1.2,
				VY: math.Sin(angle) * speed * 1.2,
			}
		}
		return spawns
	}

	return nil
}

// TakeDamage reduces HP. Returns true if enemy died.
func (e *Enemy) TakeDamage(amount int) bool {
	if !e.Alive {
		return false
	}
	e.HP -= amount
	if e.HP <= 0 {
		e.Alive = false
		return true
	}
	return false
}

// Draw renders the enemy to the framebuffer.
func (e *Enemy) Draw(fb *render.Framebuffer) {
	if !e.Alive || e.Sprite == nil {
		return
	}
	e.Sprite.DrawAt(fb, int(e.X), int(e.Y))
}
