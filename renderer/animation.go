package renderer

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Animation represents an animated sequence
type Animation struct {
	Frames        []int   // Frame indices
	FrameDuration float64 // Duration of each frame in seconds
	Loop          bool    // Whether animation loops
	CurrentFrame  int     // Current frame index
	Timer         float64 // Time accumulator
}

// NewAnimation creates a new animation
func NewAnimation(frames []int, frameDuration float64, loop bool) *Animation {
	return &Animation{
		Frames:        frames,
		FrameDuration: frameDuration,
		Loop:          loop,
		CurrentFrame:  0,
		Timer:         0,
	}
}

// Update advances the animation
func (a *Animation) Update(deltaTime float64) {
	a.Timer += deltaTime

	if a.Timer >= a.FrameDuration {
		a.Timer -= a.FrameDuration
		a.CurrentFrame++

		if a.CurrentFrame >= len(a.Frames) {
			if a.Loop {
				a.CurrentFrame = 0
			} else {
				a.CurrentFrame = len(a.Frames) - 1
			}
		}
	}
}

// GetCurrentFrame returns the current frame index
func (a *Animation) GetCurrentFrame() int {
	if a.CurrentFrame < len(a.Frames) {
		return a.Frames[a.CurrentFrame]
	}
	return 0
}

// Reset resets the animation to the beginning
func (a *Animation) Reset() {
	a.CurrentFrame = 0
	a.Timer = 0
}

// IsFinished checks if a non-looping animation has finished
func (a *Animation) IsFinished() bool {
	return !a.Loop && a.CurrentFrame >= len(a.Frames)-1
}

// ParticleType represents different particle effects
type ParticleType int

const (
	ParticleStar ParticleType = iota
	ParticleHeart
	ParticleNote
	ParticleFood
	ParticleZ
	ParticleExclamation
)

// Particle represents a visual effect particle
type Particle struct {
	X, Y     float32
	VX, VY   float32 // Velocity
	Life     float32 // Remaining life in frames
	Type     ParticleType
	Color    color.Color
	Size     float32
	Rotation float32
	RotSpeed float32
}

// Update updates the particle
func (p *Particle) Update() {
	p.X += p.VX
	p.Y += p.VY
	p.Life--
	p.Rotation += p.RotSpeed

	// Apply gravity for some particle types
	if p.Type == ParticleFood {
		p.VY += 0.1
	}

	// Fade out
	if c, ok := p.Color.(color.RGBA); ok {
		alpha := float32(c.A) * (p.Life / 60.0)
		if alpha < 0 {
			alpha = 0
		}
		p.Color = color.RGBA{c.R, c.G, c.B, uint8(alpha)}
	}
}

// Draw renders the particle
func (p *Particle) Draw(screen *ebiten.Image) {
	if p.Life <= 0 {
		return
	}

	switch p.Type {
	case ParticleStar:
		p.drawStar(screen)
	case ParticleHeart:
		p.drawHeart(screen)
	case ParticleNote:
		p.drawMusicNote(screen)
	case ParticleFood:
		p.drawFoodParticle(screen)
	case ParticleZ:
		p.drawZ(screen)
	case ParticleExclamation:
		p.drawExclamation(screen)
	}
}

func (p *Particle) drawStar(screen *ebiten.Image) {
	// Simplified star - just use a circle
	vector.DrawFilledCircle(screen, p.X, p.Y, p.Size, p.Color, false)
}

func (p *Particle) drawHeart(screen *ebiten.Image) {
	// Simplified heart - use two circles and a triangle
	vector.DrawFilledCircle(screen, p.X-p.Size/2, p.Y, p.Size/2, p.Color, false)
	vector.DrawFilledCircle(screen, p.X+p.Size/2, p.Y, p.Size/2, p.Color, false)
	vector.DrawFilledRect(screen, p.X-p.Size/2, p.Y, p.Size, p.Size, p.Color, false)
}

func (p *Particle) drawMusicNote(screen *ebiten.Image) {
	// Simple music note shape
	vector.DrawFilledCircle(screen, p.X, p.Y, p.Size, p.Color, false)
	vector.DrawFilledRect(screen, p.X+p.Size-2, p.Y-p.Size*2, 2, p.Size*2, p.Color, false)
}

func (p *Particle) drawFoodParticle(screen *ebiten.Image) {
	// Small square for food particles
	vector.DrawFilledRect(screen, p.X-p.Size/2, p.Y-p.Size/2, p.Size, p.Size, p.Color, false)
}

func (p *Particle) drawZ(screen *ebiten.Image) {
	// Simplified Z - just use rectangles
	vector.DrawFilledRect(screen, p.X-p.Size, p.Y-p.Size, p.Size*2, 2, p.Color, false)
	vector.DrawFilledRect(screen, p.X-p.Size, p.Y+p.Size-2, p.Size*2, 2, p.Color, false)
	vector.StrokeLine(screen, p.X+p.Size, p.Y-p.Size, p.X-p.Size, p.Y+p.Size, 2, p.Color, false)
}

func (p *Particle) drawExclamation(screen *ebiten.Image) {
	// Exclamation mark
	vector.DrawFilledRect(screen, p.X-2, p.Y-p.Size, 4, p.Size*0.7, p.Color, false)
	vector.DrawFilledCircle(screen, p.X, p.Y+p.Size*0.2, 2, p.Color, false)
}

// AnimationSet manages multiple animations
type AnimationSet struct {
	animations map[string]*Animation
	current    string
}

// NewAnimationSet creates a new animation set
func NewAnimationSet() *AnimationSet {
	return &AnimationSet{
		animations: make(map[string]*Animation),
		current:    "",
	}
}

// AddAnimation adds an animation to the set
func (as *AnimationSet) AddAnimation(name string, animation *Animation) {
	as.animations[name] = animation
}

// SetCurrent sets the current animation
func (as *AnimationSet) SetCurrent(name string) {
	if as.current != name {
		as.current = name
		if anim, exists := as.animations[name]; exists {
			anim.Reset()
		}
	}
}

// Update updates the current animation
func (as *AnimationSet) Update(deltaTime float64) {
	if anim, exists := as.animations[as.current]; exists {
		anim.Update(deltaTime)
	}
}

// GetCurrentFrame returns the current frame of the active animation
func (as *AnimationSet) GetCurrentFrame() int {
	if anim, exists := as.animations[as.current]; exists {
		return anim.GetCurrentFrame()
	}
	return 0
}

// CreateCreatureAnimations creates standard creature animations
func CreateCreatureAnimations() *AnimationSet {
	set := NewAnimationSet()

	// Idle animation
	set.AddAnimation("idle", NewAnimation([]int{0, 1, 0, 2}, 0.5, true))

	// Walking animation
	set.AddAnimation("walk", NewAnimation([]int{3, 4, 5, 6}, 0.15, true))

	// Eating animation
	set.AddAnimation("eat", NewAnimation([]int{7, 8, 7, 8}, 0.2, true))

	// Sleeping animation
	set.AddAnimation("sleep", NewAnimation([]int{9, 10}, 1.0, true))

	// Happy animation
	set.AddAnimation("happy", NewAnimation([]int{11, 12, 11, 12}, 0.1, true))

	return set
}
