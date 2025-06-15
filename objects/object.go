package objects

import (
	"github.com/olivierh59500/creatures-clone/utils"
)

// Object represents any interactive object in the world
type Object interface {
	// Core methods
	Update()
	GetPosition() utils.Vector2D
	GetType() string
	GetID() string

	// Interaction
	Interact(creature interface{})
	CanInteract() bool

	// State
	ShouldRemove() bool
	IsVisible() bool

	// Rendering info
	GetSprite() string
	GetColor() utils.Color
	GetSize() float64
	GetLayer() int // Rendering layer (0 = background, higher = foreground)
}

// BaseObject provides common object functionality
type BaseObject struct {
	ID       string
	Position utils.Vector2D
	Size     float64
	Color    utils.Color
	Visible  bool
	Remove   bool
	Layer    int
}

// NewBaseObject creates a new base object
func NewBaseObject(x, y float64) BaseObject {
	return BaseObject{
		ID:       utils.GenerateID(),
		Position: utils.Vector2D{X: x, Y: y},
		Size:     1.0,
		Color:    utils.Color{R: 255, G: 255, B: 255, A: 255},
		Visible:  true,
		Remove:   false,
		Layer:    1,
	}
}

// GetPosition returns the object's position
func (b *BaseObject) GetPosition() utils.Vector2D {
	return b.Position
}

// GetID returns the object's unique identifier
func (b *BaseObject) GetID() string {
	return b.ID
}

// ShouldRemove checks if the object should be removed
func (b *BaseObject) ShouldRemove() bool {
	return b.Remove
}

// IsVisible checks if the object is visible
func (b *BaseObject) IsVisible() bool {
	return b.Visible
}

// GetColor returns the object's color
func (b *BaseObject) GetColor() utils.Color {
	return b.Color
}

// GetSize returns the object's size multiplier
func (b *BaseObject) GetSize() float64 {
	return b.Size
}

// GetLayer returns the rendering layer
func (b *BaseObject) GetLayer() int {
	return b.Layer
}

// MarkForRemoval marks the object for removal
func (b *BaseObject) MarkForRemoval() {
	b.Remove = true
}

// SetVisible sets the object's visibility
func (b *BaseObject) SetVisible(visible bool) {
	b.Visible = visible
}

// SetPosition sets the object's position
func (b *BaseObject) SetPosition(x, y float64) {
	b.Position.X = x
	b.Position.Y = y
}

// Move moves the object by a delta
func (b *BaseObject) Move(dx, dy float64) {
	b.Position.X += dx
	b.Position.Y += dy
}
