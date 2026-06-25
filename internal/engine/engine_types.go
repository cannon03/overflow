package engine

import (
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
