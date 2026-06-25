package entity

import "overflow/internal/render"

// Entity is the base interface for all game entities.
type Entity interface {
	ID() int
	Position() (x, y float64)
	SetPosition(x, y float64)
	Update(dt float64)
	Draw(fb *render.Framebuffer)
	IsAlive() bool
	Destroy()
	Bounds() (x, y, w, h int)
}

// Key represents a keyboard key press.
type Key int

const (
	KeyNone Key = iota
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeyW
	KeyA
	KeyS
	KeyD
	KeySpace
	KeyEsc
	KeyP
	KeyR // restart after game over
	KeyQ // quit
)

// BaseEntity provides common fields for all entities.
type BaseEntity struct {
	IDVal int
	X, Y  float64
	Alive bool
}

func (e *BaseEntity) ID() int                { return e.IDVal }
func (e *BaseEntity) Position() (float64, float64) { return e.X, e.Y }
func (e *BaseEntity) SetPosition(x, y float64)     { e.X, e.Y = x, y }
func (e *BaseEntity) IsAlive() bool          { return e.Alive }
func (e *BaseEntity) Destroy()               { e.Alive = false }

// EntityManager manages all entities in the game world.
type EntityManager struct {
	entities []Entity
	nextID   int
}

// NewEntityManager creates a new entity manager.
func NewEntityManager() *EntityManager {
	return &EntityManager{
		entities: make([]Entity, 0, 256),
	}
}

// Add adds an entity to the manager.
func (em *EntityManager) Add(e Entity) {
	em.entities = append(em.entities, e)
}

// RemoveAll removes all entities.
func (em *EntityManager) RemoveAll() {
	em.entities = em.entities[:0]
}

// Update calls Update on all living entities.
func (em *EntityManager) Update(dt float64) {
	for _, e := range em.entities {
		if e.IsAlive() {
			e.Update(dt)
		}
	}
}

// Draw calls Draw on all living entities.
func (em *EntityManager) Draw(fb *render.Framebuffer) {
	for _, e := range em.entities {
		if e.IsAlive() {
			e.Draw(fb)
		}
	}
}

// GarbageCollect removes all dead entities.
func (em *EntityManager) GarbageCollect() {
	alive := em.entities[:0]
	for _, e := range em.entities {
		if e.IsAlive() {
			alive = append(alive, e)
		}
	}
	em.entities = alive
}

// All returns all entities.
func (em *EntityManager) All() []Entity {
	return em.entities
}

// AllAlive returns all living entities.
func (em *EntityManager) AllAlive() []Entity {
	var result []Entity
	for _, e := range em.entities {
		if e.IsAlive() {
			result = append(result, e)
		}
	}
	return result
}

// FilterByType returns all living entities matching the filter function.
func (em *EntityManager) FilterByType(filter func(Entity) bool) []Entity {
	var result []Entity
	for _, e := range em.entities {
		if e.IsAlive() && filter(e) {
			result = append(result, e)
		}
	}
	return result
}

// NextID returns and increments the entity ID counter.
func (em *EntityManager) NextID() int {
	id := em.nextID
	em.nextID++
	return id
}

// CountAlive returns the number of living entities.
func (em *EntityManager) CountAlive() int {
	count := 0
	for _, e := range em.entities {
		if e.IsAlive() {
			count++
		}
	}
	return count
}
