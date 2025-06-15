package objects

import (
	"math"

	"github.com/olivierh59500/creatures-clone/utils"
)

// PlantType represents different types of plants
type PlantType int

const (
	PlantTree PlantType = iota
	PlantFlower
	PlantGrass
	PlantBush
)

// GrowthStage represents plant growth stages
type GrowthStage int

const (
	StageSeed GrowthStage = iota
	StageSprout
	StageYoung
	StageMature
	StageFlowering
	StageDying
)

// Plant represents a growing plant
type Plant struct {
	BaseObject

	// Plant properties
	PlantType   PlantType
	GrowthStage GrowthStage
	Age         float64
	Health      float64
	GrowthRate  float64

	// Environmental factors
	WaterLevel  float64
	SunExposure float64

	// Visual properties
	SwayOffset float64
	SwaySpeed  float64

	// Production
	ProduceTimer float64
	FruitCount   int
}

// NewPlant creates a new plant
func NewPlant(x, y float64, plantType PlantType) *Plant {
	p := &Plant{
		BaseObject:   NewBaseObject(x, y),
		PlantType:    plantType,
		GrowthStage:  StageSeed,
		Age:          0,
		Health:       100,
		GrowthRate:   getGrowthRate(plantType),
		WaterLevel:   50,
		SunExposure:  50,
		SwayOffset:   utils.RandomFloat(0, math.Pi*2),
		SwaySpeed:    0.02 + utils.RandomFloat(-0.01, 0.01),
		ProduceTimer: 0,
		FruitCount:   0,
	}

	// Set plant-specific properties
	p.Color = getPlantColor(plantType)
	p.Size = 0.2 // Start small
	p.Layer = 0  // Background layer

	// Position plant at ground level
	p.Position.Y = y - 10 // Slightly embedded in ground

	return p
}

// Update updates the plant's state
func (p *Plant) Update() {
	// Age the plant
	p.Age += p.GrowthRate

	// Update growth stage
	p.updateGrowthStage()

	// Process environmental factors
	p.processEnvironment()

	// Update health
	p.updateHealth()

	// Animate swaying
	p.SwayOffset += p.SwaySpeed

	// Produce fruit/seeds if mature
	if p.GrowthStage == StageMature || p.GrowthStage == StageFlowering {
		p.ProduceTimer += 0.016
		if p.ProduceTimer > 30 { // Every 30 seconds
			p.produceFruit()
			p.ProduceTimer = 0
		}
	}

	// Remove if dead
	if p.Health <= 0 || p.GrowthStage == StageDying && p.Age > 1000 {
		p.Remove = true
	}
}

// updateGrowthStage advances the plant through growth stages
func (p *Plant) updateGrowthStage() {
	switch p.PlantType {
	case PlantTree:
		switch {
		case p.Age < 10:
			p.GrowthStage = StageSeed
			p.Size = 0.2
		case p.Age < 50:
			p.GrowthStage = StageSprout
			p.Size = 0.4
		case p.Age < 100:
			p.GrowthStage = StageYoung
			p.Size = 0.7
		case p.Age < 500:
			p.GrowthStage = StageMature
			p.Size = 1.0
		case p.Age < 800:
			p.GrowthStage = StageFlowering
			p.Size = 1.0
		default:
			p.GrowthStage = StageDying
			p.Size = 0.95
		}

	case PlantFlower:
		switch {
		case p.Age < 5:
			p.GrowthStage = StageSeed
			p.Size = 0.2
		case p.Age < 20:
			p.GrowthStage = StageSprout
			p.Size = 0.5
		case p.Age < 40:
			p.GrowthStage = StageYoung
			p.Size = 0.8
		case p.Age < 100:
			p.GrowthStage = StageFlowering
			p.Size = 1.0
		default:
			p.GrowthStage = StageDying
			p.Size = 0.8
		}

	default:
		// Simple growth for grass and bushes
		if p.Age < 20 {
			p.Size = p.Age / 20
		} else {
			p.Size = 1.0
		}
	}
}

// processEnvironment simulates environmental effects
func (p *Plant) processEnvironment() {
	// Water consumption
	p.WaterLevel -= 0.05
	if p.WaterLevel < 0 {
		p.WaterLevel = 0
	}

	// Sun exposure varies with time of day (simplified)
	// In full implementation, would get from world
	p.SunExposure = 50 + math.Sin(p.Age*0.01)*30
}

// updateHealth updates plant health based on conditions
func (p *Plant) updateHealth() {
	// Optimal conditions
	waterOptimal := p.WaterLevel > 30 && p.WaterLevel < 70
	sunOptimal := p.SunExposure > 40 && p.SunExposure < 80

	if waterOptimal && sunOptimal {
		// Heal in good conditions
		p.Health = utils.Clamp(p.Health+0.1, 0, 100)
	} else {
		// Damage from poor conditions
		if p.WaterLevel < 20 || p.WaterLevel > 80 {
			p.Health -= 0.2
		}
		if p.SunExposure < 20 || p.SunExposure > 90 {
			p.Health -= 0.1
		}
	}

	// Age-related health decline
	if p.GrowthStage == StageDying {
		p.Health -= 0.05
	}

	p.Health = utils.Clamp(p.Health, 0, 100)
}

// produceFruit creates fruit objects nearby
func (p *Plant) produceFruit() {
	if p.PlantType != PlantTree || p.FruitCount >= 3 {
		return
	}

	// Trees produce apples
	// In full implementation, this would add to world
	p.FruitCount++
}

// GetType returns the object type
func (p *Plant) GetType() string {
	return "plant"
}

// Interact handles creature interaction
func (p *Plant) Interact(creature interface{}) {
	// Plants don't actively interact
	// But creatures might eat them or water them
}

// CanInteract checks if the plant can be interacted with
func (p *Plant) CanInteract() bool {
	return p.Health > 0
}

// GetSprite returns the sprite identifier
func (p *Plant) GetSprite() string {
	switch p.PlantType {
	case PlantTree:
		return "tree"
	case PlantFlower:
		return "flower"
	case PlantGrass:
		return "grass"
	case PlantBush:
		return "bush"
	default:
		return "plant"
	}
}

// Water adds water to the plant
func (p *Plant) Water(amount float64) {
	p.WaterLevel = utils.Clamp(p.WaterLevel+amount, 0, 100)
}

// GetSwayX returns horizontal sway for animation
func (p *Plant) GetSwayX() float64 {
	return math.Sin(p.SwayOffset) * 2 * p.Size
}

// Helper functions

func getGrowthRate(plantType PlantType) float64 {
	switch plantType {
	case PlantTree:
		return 0.5
	case PlantFlower:
		return 1.0
	case PlantGrass:
		return 2.0
	case PlantBush:
		return 0.8
	default:
		return 1.0
	}
}

func getPlantColor(plantType PlantType) utils.Color {
	switch plantType {
	case PlantTree:
		// Green for leaves
		return utils.Color{R: 34, G: 139, B: 34, A: 255}
	case PlantFlower:
		// Various flower colors - default to pink
		return utils.Color{R: 255, G: 192, B: 203, A: 255}
	case PlantGrass:
		return utils.Color{R: 144, G: 238, B: 144, A: 255}
	case PlantBush:
		return utils.Color{R: 0, G: 100, B: 0, A: 255}
	default:
		return utils.Color{R: 0, G: 128, B: 0, A: 255}
	}
}
