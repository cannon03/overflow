package entity

import (
	"math"
	"math/rand"
	"overflow/internal/render"
)

// Particle is a visual effect entity.
type Particle struct {
	BaseEntity
	VX, VY  float64
	Life    float64
	MaxLife float64
	Rune    rune
	StartFG render.Color
	EndFG   render.Color
	BG      render.Color
}

// NewParticle creates a new particle.
func NewParticle(x, y, vx, vy float64, life float64, r rune, startFG, endFG, bg render.Color) *Particle {
	return &Particle{
		BaseEntity: BaseEntity{X: x, Y: y, Alive: true},
		VX:         vx,
		VY:         vy,
		Life:       life,
		MaxLife:    life,
		Rune:       r,
		StartFG:    startFG,
		EndFG:      endFG,
		BG:         bg,
	}
}

// NewExplosion creates explosion particles at the given position.
func NewExplosion(x, y float64) []*Particle {
	count := 8 + rand.Intn(5)
	particles := make([]*Particle, count)
	for i := 0; i < count; i++ {
		angle := 2.0 * math.Pi * float64(i) / float64(count)
		speed := 20 + rand.Float64()*40
		life := 0.3 + rand.Float64()*0.3

		r := '·'
		switch rand.Intn(3) {
		case 0:
			r = '·'
		case 1:
			r = '•'
		case 2:
			r = '✦'
		}

		particles[i] = NewParticle(
			x, y,
			math.Cos(angle)*speed,
			math.Sin(angle)*speed,
			life, r,
			render.ColorOrange,
			render.ColorRed,
			render.ColorBlack,
		)
	}
	return particles
}

// NewFloatingText creates particles for a score popup.
func NewFloatingText(x, y float64, text string) []*Particle {
	particles := make([]*Particle, len(text))
	for i, ch := range text {
		particles[i] = NewParticle(
			x+float64(i), y,
			0, -15,
			0.8,
			ch,
			render.ColorYellow,
			render.ColorOrange,
			render.ColorBlack,
		)
	}
	return particles
}

// Bounds returns the particle's bounds (1x1).
func (p *Particle) Bounds() (x, y, w, h int) {
	return int(p.X), int(p.Y), 1, 1
}

// Update updates the particle position and life.
func (p *Particle) Update(dt float64) {
	if !p.Alive {
		return
	}

	p.X += p.VX * dt
	p.Y += p.VY * dt
	p.Life -= dt
	p.VX *= 0.95
	p.VY *= 0.95

	if p.Life <= 0 {
		p.Alive = false
	}
}

// Draw renders the particle to the framebuffer.
func (p *Particle) Draw(fb *render.Framebuffer) {
	if !p.Alive {
		return
	}

	t := p.Life / p.MaxLife
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}

	fg := render.Color{
		R: uint8(float64(p.StartFG.R)*t + float64(p.EndFG.R)*(1-t)),
		G: uint8(float64(p.StartFG.G)*t + float64(p.EndFG.G)*(1-t)),
		B: uint8(float64(p.StartFG.B)*t + float64(p.EndFG.B)*(1-t)),
	}

	if t < 0.1 && rand.Float64() > t*10 {
		return
	}

	fb.Set(int(math.Round(p.Y)), int(math.Round(p.X)),
		render.Cell{Rune: p.Rune, FG: fg, BG: p.BG})
}
