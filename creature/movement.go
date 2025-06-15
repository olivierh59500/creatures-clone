package creature

import (
	"math"

	"github.com/olivierh59500/creatures-clone/utils"
)

// Movement handles creature locomotion and physics
type Movement struct {
	// Movement parameters
	Speed     float64
	JumpPower float64
	Agility   float64

	// Movement state
	IsMoving  bool
	IsJumping bool
	IsRunning bool

	// Gait parameters
	GaitCycle float64 // Current position in walk cycle
	GaitSpeed float64 // How fast the gait cycles

	// Physics modifiers
	Friction      float64
	AirResistance float64
}

// NewMovement creates a new movement system
func NewMovement() *Movement {
	return &Movement{
		Speed:     2.0,
		JumpPower: 8.0,
		Agility:   1.0,

		GaitSpeed: 0.1,

		Friction:      0.9,
		AirResistance: 0.98,
	}
}

// MoveLeft moves the creature left
func (m *Movement) MoveLeft(x *float64, velocityX *float64) {
	m.IsMoving = true

	// Apply acceleration
	acceleration := m.Speed * m.Agility
	if m.IsRunning {
		acceleration *= 1.5
	}

	*velocityX -= acceleration

	// Limit max speed
	maxSpeed := m.Speed * 3
	if m.IsRunning {
		maxSpeed *= 1.5
	}
	*velocityX = utils.Clamp(*velocityX, -maxSpeed, maxSpeed)

	// Update gait cycle
	m.updateGait()
}

// MoveRight moves the creature right
func (m *Movement) MoveRight(x *float64, velocityX *float64) {
	m.IsMoving = true

	// Apply acceleration
	acceleration := m.Speed * m.Agility
	if m.IsRunning {
		acceleration *= 1.5
	}

	*velocityX += acceleration

	// Limit max speed
	maxSpeed := m.Speed * 3
	if m.IsRunning {
		maxSpeed *= 1.5
	}
	*velocityX = utils.Clamp(*velocityX, -maxSpeed, maxSpeed)

	// Update gait cycle
	m.updateGait()
}

// Jump makes the creature jump
func (m *Movement) Jump(velocityY *float64, onGround bool) {
	if !m.IsJumping && onGround {
		*velocityY = -m.JumpPower
		m.IsJumping = true
	}
}

// Stop halts movement
func (m *Movement) Stop() {
	m.IsMoving = false
	m.IsRunning = false
}

// Run enables running mode
func (m *Movement) Run() {
	m.IsRunning = true
}

// Walk enables walking mode
func (m *Movement) Walk() {
	m.IsRunning = false
}

// updateGait advances the walking animation cycle
func (m *Movement) updateGait() {
	speed := m.GaitSpeed
	if m.IsRunning {
		speed *= 1.5
	}

	m.GaitCycle += speed
	if m.GaitCycle > 2*math.Pi {
		m.GaitCycle -= 2 * math.Pi
	}
}

// GetGaitOffset returns the vertical offset for walking animation
func (m *Movement) GetGaitOffset() float64 {
	if !m.IsMoving {
		return 0
	}

	// Create a bouncing motion
	return math.Sin(m.GaitCycle) * 2
}

// GetLegPosition returns leg positions for animation
func (m *Movement) GetLegPosition(isLeftLeg bool) (x, y float64) {
	if !m.IsMoving {
		return 0, 0
	}

	// Offset legs by half cycle
	cycle := m.GaitCycle
	if isLeftLeg {
		cycle += math.Pi
	}

	// Create walking motion
	x = math.Sin(cycle) * 5
	y = math.Max(0, math.Sin(cycle*2)) * 3

	return x, y
}

// ApplyPhysics applies physics constraints to movement
func (m *Movement) ApplyPhysics(velocityX, velocityY *float64, onGround bool) {
	// Apply friction when on ground
	if onGround {
		*velocityX *= m.Friction

		// Stop jumping when on ground
		if *velocityY >= 0 {
			m.IsJumping = false
		}
	} else {
		// Apply air resistance
		*velocityX *= m.AirResistance
		*velocityY *= m.AirResistance
	}

	// Stop movement if velocity is very small
	if math.Abs(*velocityX) < 0.01 {
		*velocityX = 0
		m.IsMoving = false
	}
}

// GetSpeed returns current movement speed
func (m *Movement) GetSpeed() float64 {
	if m.IsRunning {
		return m.Speed * 1.5
	}
	return m.Speed
}

// SetSpeed sets the base movement speed
func (m *Movement) SetSpeed(speed float64) {
	m.Speed = utils.Clamp(speed, 0.5, 5.0)
}

// SetJumpPower sets the jump strength
func (m *Movement) SetJumpPower(power float64) {
	m.JumpPower = utils.Clamp(power, 2.0, 15.0)
}

// SetAgility sets the movement agility
func (m *Movement) SetAgility(agility float64) {
	m.Agility = utils.Clamp(agility, 0.5, 2.0)
}

// GetEnergyUsage returns energy cost of current movement
func (m *Movement) GetEnergyUsage() float64 {
	if !m.IsMoving {
		return 0
	}

	energyUse := 0.1
	if m.IsRunning {
		energyUse *= 2
	}
	if m.IsJumping {
		energyUse *= 1.5
	}

	return energyUse
}

// CanMove checks if movement is possible
func (m *Movement) CanMove(energy float64) bool {
	requiredEnergy := 10.0
	if m.IsRunning {
		requiredEnergy = 20.0
	}

	return energy >= requiredEnergy
}

// UpdateFromGenetics applies genetic traits to movement
func (m *Movement) UpdateFromGenetics(speedGene, agilityGene, strengthGene float64) {
	// Map genes (0-1) to movement parameters
	m.Speed = 1.0 + speedGene*3.0        // 1.0 to 4.0
	m.Agility = 0.5 + agilityGene*1.5    // 0.5 to 2.0
	m.JumpPower = 4.0 + strengthGene*8.0 // 4.0 to 12.0
}
