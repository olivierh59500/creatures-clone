package objects

import (
	"math"

	"github.com/olivierh59500/creatures-clone/utils"
)

// ToyType represents different types of toys
type ToyType int

const (
	ToyBall ToyType = iota
	ToyMusicBox
	ToyPuzzle
	ToyMirror
	ToyComputer
	ToyBed
)

// Toy represents an interactive plaything
type Toy struct {
	BaseObject

	// Toy properties
	ToyType      ToyType
	Durability   float64
	IsActivated  bool
	LastUsedTime float64

	// Animation properties
	Rotation      float64
	BounceHeight  float64
	AnimationTime float64

	// Interaction tracking
	TimesUsed  int
	LastUserID string
}

// NewToy creates a new toy
func NewToy(x, y float64, toyType ToyType) *Toy {
	t := &Toy{
		BaseObject:    NewBaseObject(x, y),
		ToyType:       toyType,
		Durability:    100,
		IsActivated:   false,
		LastUsedTime:  0,
		Rotation:      0,
		BounceHeight:  0,
		AnimationTime: 0,
		TimesUsed:     0,
	}

	// Set toy-specific properties
	t.Color = getToyColor(toyType)
	t.Size = getToySize(toyType)

	return t
}

// Update updates the toy's state
func (t *Toy) Update() {
	// Update animation
	t.AnimationTime += 0.016 // 60 FPS

	switch t.ToyType {
	case ToyBall:
		// Ball rolls and bounces
		if t.IsActivated {
			t.Rotation += 0.1
			t.BounceHeight = math.Abs(math.Sin(t.AnimationTime*3)) * 10

			// Deactivate after a while
			if t.AnimationTime > 3 {
				t.IsActivated = false
				t.AnimationTime = 0
			}
		}

	case ToyMusicBox:
		// Music box plays and rotates handle
		if t.IsActivated {
			t.Rotation += 0.05

			// Stop after one "song"
			if t.AnimationTime > 5 {
				t.IsActivated = false
				t.AnimationTime = 0
			}
		}

	case ToyPuzzle:
		// Puzzle pieces move
		if t.IsActivated {
			// Simple animation
			if t.AnimationTime > 2 {
				t.IsActivated = false
				t.AnimationTime = 0
			}
		}

	case ToyMirror:
		// Mirror doesn't animate much
		if t.IsActivated {
			if t.AnimationTime > 1 {
				t.IsActivated = false
				t.AnimationTime = 0
			}
		}

	case ToyComputer:
		// Computer shows learning animation
		if t.IsActivated {
			// Blinking screen effect
			if t.AnimationTime > 3 {
				t.IsActivated = false
				t.AnimationTime = 0
			}
		}

	case ToyBed:
		// Bed doesn't animate but provides comfort
		if t.IsActivated {
			if t.AnimationTime > 5 {
				t.IsActivated = false
				t.AnimationTime = 0
			}
		}
	}

	// Wear and tear
	if t.IsActivated {
		t.Durability -= 0.01
		if t.Durability <= 0 {
			t.Remove = true
		}
	}

	// Cool down time
	t.LastUsedTime += 0.016
}

// GetType returns the object type
func (t *Toy) GetType() string {
	return "toy"
}

// Interact handles creature interaction
func (t *Toy) Interact(creature interface{}) {
	// Can only interact if not already activated and cooled down
	if !t.IsActivated && t.LastUsedTime > 1 {
		t.IsActivated = true
		t.AnimationTime = 0
		t.LastUsedTime = 0
		t.TimesUsed++

		// Track who used it (would need creature ID)
		// t.LastUserID = creature.GetID()
	}
}

// CanInteract checks if the toy can be interacted with
func (t *Toy) CanInteract() bool {
	return !t.IsActivated && t.Durability > 0 && t.LastUsedTime > 1
}

// GetSprite returns the sprite identifier
func (t *Toy) GetSprite() string {
	switch t.ToyType {
	case ToyBall:
		return "ball"
	case ToyMusicBox:
		return "musicbox"
	case ToyPuzzle:
		return "puzzle"
	case ToyMirror:
		return "mirror"
	case ToyComputer:
		return "computer"
	case ToyBed:
		return "bed"
	default:
		return "toy"
	}
}

// GetRotation returns the current rotation angle
func (t *Toy) GetRotation() float64 {
	return t.Rotation
}

// GetBounceOffset returns the vertical bounce offset
func (t *Toy) GetBounceOffset() float64 {
	return t.BounceHeight
}

// IsPlaying checks if the toy is currently active
func (t *Toy) IsPlaying() bool {
	return t.IsActivated
}

// GetDurabilityPercent returns durability as a percentage
func (t *Toy) GetDurabilityPercent() float64 {
	return t.Durability
}

// Helper functions

func getToyColor(toyType ToyType) utils.Color {
	switch toyType {
	case ToyBall:
		// Multicolor ball - return primary color
		return utils.Color{R: 255, G: 0, B: 0, A: 255} // Red
	case ToyMusicBox:
		return utils.Color{R: 139, G: 69, B: 19, A: 255} // Wood brown
	case ToyPuzzle:
		return utils.Color{R: 100, G: 100, B: 255, A: 255} // Light blue
	case ToyMirror:
		return utils.Color{R: 192, G: 192, B: 192, A: 255} // Silver
	case ToyComputer:
		return utils.Color{R: 128, G: 128, B: 128, A: 255} // Gray
	case ToyBed:
		return utils.Color{R: 65, G: 105, B: 225, A: 255} // Royal blue
	default:
		return utils.Color{R: 200, G: 200, B: 200, A: 255}
	}
}

func getToySize(toyType ToyType) float64 {
	switch toyType {
	case ToyBall:
		return 1.2
	case ToyMusicBox:
		return 1.5
	case ToyPuzzle:
		return 1.3
	case ToyMirror:
		return 1.0
	case ToyComputer:
		return 1.8
	case ToyBed:
		return 2.0
	default:
		return 1.0
	}
}
