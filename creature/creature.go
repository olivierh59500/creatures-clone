// generateName generates a random name for a creature
func generateName(creatureType CreatureType) string {
	prefixes := []string{"Ala", "Bel", "Cor", "Dex", "Eva", "Flo", "Gus", "Hex", "Ira", "Jax"}
	suffixes := []string{"bert", "mina", "dor", "thy", "ron", "lia", "max", "win", "zor", "bella"}
	
	prefix := prefixes[utils.RandomInt(0, len(prefixes))]
	suffix := suffixes[utils.RandomInt(0, len(suffixes))]
	
	return prefix + suffix
}package creature

import (
	"math"
	"math/rand"
	
	"github.com/olivierh59500/creatures-clone/utils"
)

// CreatureType represents different species of creatures
type CreatureType int

const (
	CreatureTypeNorn CreatureType = iota
	CreatureTypeGrendel
	CreatureTypeEttin
)

// AgeStage represents the life stage of a creature
type AgeStage int

const (
	AgeBaby AgeStage = iota
	AgeChild
	AgeAdult
	AgeElder
)

// Creature represents a living entity in the game
type Creature struct {
	// Identity
	ID   string
	Type CreatureType
	Name string
	
	// Position and physics
	X, Y           float64
	VelocityX      float64
	VelocityY      float64
	Direction      float64 // Facing direction in radians
	
	// Core systems
	Brain      *Brain
	Genetics   *Genetics
	Metabolism *Metabolism
	Emotions   *Emotions
	Movement   *Movement
	Learning   *Learning
	Language   *Language
	
	// Physical attributes
	Age       float64 // Age in game minutes
	AgeStage  AgeStage
	Size      float64
	Color     utils.Color
	
	// State
	IsAsleep  bool
	IsSick    bool
	
	// Goals
	TargetX   float64
	TargetY   float64
	HasTarget bool
	
	// Animation
	AnimationState string
	AnimationFrame int
	AnimationTimer float64
	
	// Sensory input
	Vision      []float64 // What the creature sees
	Hearing     []string  // Words heard recently
	Touch       []float64 // Physical sensations
	
	// Memory
	RecentActions []int     // Recent action history
	LastBreedTime float64   // Time since last breeding
}

// Neural network output indices
const (
	OutputMoveLeft = iota
	OutputMoveRight
	OutputJump
	OutputEat
	OutputSleep
	OutputPlay
	OutputSpeak
	OutputBreed
	OutputMax
)

// NewCreature creates a new creature instance
func NewCreature(x, y float64, creatureType CreatureType) *Creature {
	id := utils.GenerateID()
	
	c := &Creature{
		ID:        id,
		Type:      creatureType,
		Name:      generateName(creatureType),
		X:         x,
		Y:         y,
		Direction: 0,
		Size:      1.0,
		AgeStage:  AgeAdult,
		
		// Initialize systems
		Brain:      NewBrain(),
		Genetics:   NewGenetics(),
		Metabolism: NewMetabolism(),
		Emotions:   NewEmotions(),
		Movement:   NewMovement(),
		Learning:   NewLearning(),
		Language:   NewLanguage(),
		
		// Sensory arrays
		Vision:  make([]float64, 20), // 20 vision sensors
		Hearing: make([]string, 5),   // Remember last 5 words
		Touch:   make([]float64, 4),  // 4 touch sensors
		
		RecentActions: make([]int, 10),
		
		AnimationState: "idle",
	}
	
	// Apply genetic traits
	c.applyGenetics()
	
	return c
}

// Update updates the creature's state
func (c *Creature) Update(world interface{}) {
	// Update age
	c.Age += 1.0 / (60.0 * 60.0) // 1 game minute = 1 real second at 60 FPS
	c.updateAgeStage()
	
	// Update metabolism
	c.Metabolism.Update(c.Movement.GetSpeed())
	
	// Check health conditions
	c.updateHealthStatus()
	
	// Process sensory input through brain
	brainInput := c.prepareBrainInput()
	c.Brain.Process(brainInput)
	
	// Execute actions based on brain output
	c.executeActions()
	
	// Update emotions based on current state
	c.Emotions.Update(c.Metabolism, c.Brain.GetOutput())
	
	// Update animation
	c.updateAnimation()
	
	// Learning from experiences
	c.Learning.Update(c.Brain, c.RecentActions)
}

// UpdateSensors updates the creature's sensory input
func (c *Creature) UpdateSensors(nearbyEntities []interface{}, world interface{}) {
	// Clear vision
	for i := range c.Vision {
		c.Vision[i] = 0
	}
	
	// Process nearby entities for vision
	for _, entity := range nearbyEntities {
		// Different processing for different entity types
		// This is simplified - in a full implementation, you'd calculate
		// which vision "pixel" the entity appears in based on angle
		switch e := entity.(type) {
		case *Creature:
			if e != c {
				angle := math.Atan2(e.Y-c.Y, e.X-c.X) - c.Direction
				visionIndex := c.angleToVisionIndex(angle)
				if visionIndex >= 0 && visionIndex < len(c.Vision) {
					c.Vision[visionIndex] = 1.0 // Creature detected
				}
			}
		default:
			// Handle other object types
		}
	}
	
	// Update touch sensors based on collisions
	// Simplified - would check actual collisions
	c.Touch[0] = 0 // Front
	c.Touch[1] = 0 // Back
	c.Touch[2] = 0 // Left
	c.Touch[3] = 0 // Right
}

// prepareBrainInput prepares input vector for the neural network
func (c *Creature) prepareBrainInput() []float64 {
	input := make([]float64, 0)
	
	// Add vision sensors
	input = append(input, c.Vision...)
	
	// Add internal state sensors
	input = append(input, 
		c.Metabolism.Hunger/100.0,
		c.Metabolism.Energy/100.0,
		c.Metabolism.Health/100.0,
		c.Emotions.Happiness/100.0,
		c.Emotions.Fear/100.0,
		c.Emotions.Anger/100.0,
		c.Emotions.Curiosity/100.0,
	)
	
	// Add touch sensors
	input = append(input, c.Touch...)
	
	// Add time of day sensor (would get from world)
	input = append(input, 0.5) // Placeholder
	
	return input
}

// executeActions performs actions based on brain output
func (c *Creature) executeActions() {
	output := c.Brain.GetOutput()
	
	// Check if we have a target to move towards
	if c.HasTarget {
		c.MoveTowardsTarget()
	} else {
		// Normal AI-driven movement
		if output[OutputMoveLeft] > 0.5 {
			c.Movement.MoveLeft(&c.X, &c.VelocityX)
			c.Direction = math.Pi // Face left
			c.recordAction(OutputMoveLeft)
		}
		if output[OutputMoveRight] > 0.5 {
			c.Movement.MoveRight(&c.X, &c.VelocityX)
			c.Direction = 0 // Face right
			c.recordAction(OutputMoveRight)
		}
	}
	
	if output[OutputJump] > 0.5 {
		// Check if on ground (80% of world height)
		onGround := c.Y >= 400 // This will be updated by world physics
		c.Movement.Jump(&c.VelocityY, onGround)
		c.recordAction(OutputJump)
	}
	
	// Other actions are handled by world interaction system
	// but we record the intention
	if output[OutputEat] > 0.5 {
		c.recordAction(OutputEat)
	}
	if output[OutputSleep] > 0.5 {
		c.IsAsleep = true
		c.recordAction(OutputSleep)
	} else {
		c.IsAsleep = false
	}
	if output[OutputPlay] > 0.5 {
		c.recordAction(OutputPlay)
	}
	if output[OutputSpeak] > 0.5 {
		c.recordAction(OutputSpeak)
	}
	if output[OutputBreed] > 0.5 {
		c.recordAction(OutputBreed)
	}
	
	// Apply physics
	c.X += c.VelocityX
	c.Y += c.VelocityY
	
	// Friction
	c.VelocityX *= 0.9
}

// recordAction adds an action to recent history
func (c *Creature) recordAction(action int) {
	// Shift array and add new action
	copy(c.RecentActions[1:], c.RecentActions[:len(c.RecentActions)-1])
	c.RecentActions[0] = action
}

// updateAgeStage updates the creature's life stage
func (c *Creature) updateAgeStage() {
	switch {
	case c.Age < 5:
		c.AgeStage = AgeBaby
		c.Size = 0.7
	case c.Age < 15:
		c.AgeStage = AgeChild
		c.Size = 0.85
	case c.Age < 45:
		c.AgeStage = AgeAdult
		c.Size = 1.0
	default:
		c.AgeStage = AgeElder
		c.Size = 0.95
	}
}

// updateHealthStatus updates sickness and other health states
func (c *Creature) updateHealthStatus() {
	// Check for sickness conditions
	if c.Metabolism.Health < 30 || c.Metabolism.Hunger > 80 {
		c.IsSick = true
	} else if c.Metabolism.Health > 50 {
		c.IsSick = false
	}
}

// updateAnimation updates the creature's animation state
func (c *Creature) updateAnimation() {
	c.AnimationTimer += 1.0 / 60.0 // 60 FPS
	
	// Determine animation state
	newState := "idle"
	if c.IsAsleep {
		newState = "sleep"
	} else if c.IsSick {
		newState = "sick"
	} else if math.Abs(c.VelocityX) > 0.1 {
		newState = "walk"
	} else if c.Metabolism.Hunger > 70 {
		newState = "hungry"
	} else if c.Emotions.Happiness > 50 {
		newState = "happy"
	} else if c.Emotions.Fear > 50 {
		newState = "scared"
	}
	
	// Reset animation if state changed
	if newState != c.AnimationState {
		c.AnimationState = newState
		c.AnimationFrame = 0
		c.AnimationTimer = 0
	}
	
	// Advance animation frame
	animationSpeed := 0.1 // 10 FPS animation
	if c.AnimationTimer > animationSpeed {
		c.AnimationFrame++
		c.AnimationTimer = 0
		
		// Loop animation (assume 4 frames per animation)
		if c.AnimationFrame >= 4 {
			c.AnimationFrame = 0
		}
	}
}

// applyGenetics applies genetic traits to the creature
func (c *Creature) applyGenetics() {
	genes := c.Genetics.Genes
	
	// Apply color from genetics
	c.Color = c.Genetics.GetColor()
	
	// Apply genetic modifiers to systems
	c.Metabolism.HungerRate *= genes["metabolism_rate"]
	c.Movement.Speed *= genes["movement_speed"]
	c.Learning.LearningRate *= genes["learning_rate"]
	
	// Apply personality traits
	c.Emotions.BaseHappiness = (genes["happiness_bias"] - 0.5) * 40
	c.Emotions.FearThreshold = genes["fear_threshold"] * 100
	c.Emotions.AngerThreshold = genes["anger_threshold"] * 100
}

// angleToVisionIndex converts an angle to a vision array index
func (c *Creature) angleToVisionIndex(angle float64) int {
	// Normalize angle to -π to π
	for angle > math.Pi {
		angle -= 2 * math.Pi
	}
	for angle < -math.Pi {
		angle += 2 * math.Pi
	}
	
	// Convert to vision index (assuming 180 degree field of view)
	if math.Abs(angle) > math.Pi/2 {
		return -1 // Outside field of view
	}
	
	index := int((angle + math.Pi/2) / (math.Pi / float64(len(c.Vision))))
	return utils.ClampInt(index, 0, len(c.Vision)-1)
}

// Contains checks if a point is within the creature
func (c *Creature) Contains(x, y float64) bool {
	// Simple circular hit detection
	radius := 20.0 * c.Size
	dx := x - c.X
	dy := y - c.Y
	return dx*dx+dy*dy <= radius*radius
}

// GetNearestObject finds the nearest object from a list
func (c *Creature) GetNearestObject(objects []interface{}) interface{} {
	var nearest interface{}
	minDist := math.MaxFloat64
	
	for _, obj := range objects {
		// Would need type assertion and position extraction
		// This is simplified
		dist := 100.0 // Placeholder
		if dist < minDist {
			minDist = dist
			nearest = obj
		}
	}
	
	return nearest
}

// CanBreed checks if the creature can breed
func (c *Creature) CanBreed() bool {
	return c.AgeStage == AgeAdult &&
		c.Metabolism.Health > 70 &&
		c.Metabolism.Energy > 50 &&
		c.Age-c.LastBreedTime > 10 // 10 minute cooldown
}

// IsDead checks if the creature has died
func (c *Creature) IsDead() bool {
	return c.Metabolism.Health <= 0 || c.Age > 60 // 60 minute lifespan
}

// SetTarget sets a movement target for the creature
func (c *Creature) SetTarget(x, y float64) {
	c.TargetX = x
	c.TargetY = y
	c.HasTarget = true
	
	// Increase curiosity when given a target
	c.Emotions.AdjustCuriosity(10)
}

// ClearTarget removes the movement target
func (c *Creature) ClearTarget() {
	c.HasTarget = false
}

// EncourageBreeding increases breeding desire
func (c *Creature) EncourageBreeding() {
	// Temporarily boost breeding output in brain
	if c.CanBreed() {
		c.Brain.GetOutput()[OutputBreed] = 0.9
		c.Emotions.AdjustHappiness(20)
	}
}

// MoveTowardsTarget moves creature towards its target
func (c *Creature) MoveTowardsTarget() {
	if !c.HasTarget {
		return
	}
	
	// Calculate direction to target
	dx := c.TargetX - c.X
	dy := c.TargetY - c.Y
	dist := math.Sqrt(dx*dx + dy*dy)
	
	// If close enough, clear target
	if dist < 20 {
		c.ClearTarget()
		return
	}
	
	// Move towards target
	speed := c.Movement.GetSpeed()
	c.VelocityX = (dx / dist) * speed
	c.VelocityY = (dy / dist) * speed
	
	// Face direction of movement
	if dx > 0 {
		c.Direction = 0
	} else {
		c.Direction = math.Pi
	}
}