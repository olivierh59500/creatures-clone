package objects

import (
	"github.com/olivierh59500/creatures-clone/utils"
)

// FoodType represents different types of food
type FoodType int

const (
	FoodApple FoodType = iota
	FoodCarrot
	FoodHoney
	FoodSeed
	FoodBerry
)

// Food represents an edible object
type Food struct {
	BaseObject

	// Food properties
	FoodType   FoodType
	Nutrition  float64
	Freshness  float64
	IsConsumed bool

	// Visual properties
	BounceOffset float64
	BounceSpeed  float64
}

// NewFood creates a new food item
func NewFood(x, y float64, foodType FoodType) *Food {
	f := &Food{
		BaseObject:   NewBaseObject(x, y),
		FoodType:     foodType,
		Nutrition:    getFoodNutrition(foodType),
		Freshness:    100,
		IsConsumed:   false,
		BounceOffset: 0,
		BounceSpeed:  0.1,
	}

	// Set food-specific properties
	f.Color = getFoodColor(foodType)
	f.Size = getFoodSize(foodType)

	return f
}

// Update updates the food's state
func (f *Food) Update() {
	// Food decays over time
	f.Freshness -= 0.01
	if f.Freshness <= 0 {
		f.Freshness = 0
		f.Nutrition *= 0.5 // Rotten food is less nutritious
	}

	// Animate bounce
	f.BounceOffset += f.BounceSpeed

	// Remove if consumed or completely rotten
	if f.IsConsumed || (f.Freshness <= 0 && f.Nutrition < 1) {
		f.Remove = true
	}
}

// GetType returns the object type
func (f *Food) GetType() string {
	return "food"
}

// Interact handles creature interaction
func (f *Food) Interact(creature interface{}) {
	// Food doesn't actively interact back
	// The creature will consume it if hungry
}

// CanInteract checks if the food can be interacted with
func (f *Food) CanInteract() bool {
	return !f.IsConsumed && f.Nutrition > 0
}

// Consume marks the food as eaten
func (f *Food) Consume() {
	f.IsConsumed = true
	f.Remove = true
}

// GetNutrition returns the nutritional value
func (f *Food) GetNutrition() float64 {
	// Adjust nutrition based on freshness
	return f.Nutrition * (0.5 + f.Freshness/200)
}

// GetSprite returns the sprite identifier
func (f *Food) GetSprite() string {
	switch f.FoodType {
	case FoodApple:
		return "apple"
	case FoodCarrot:
		return "carrot"
	case FoodHoney:
		return "honey"
	case FoodSeed:
		return "seed"
	case FoodBerry:
		return "berry"
	default:
		return "food"
	}
}

// GetBounceY returns the vertical offset for animation
func (f *Food) GetBounceY() float64 {
	return utils.Sin(f.BounceOffset) * 2
}

// Helper functions

func getFoodNutrition(foodType FoodType) float64 {
	switch foodType {
	case FoodApple:
		return 25
	case FoodCarrot:
		return 20
	case FoodHoney:
		return 40
	case FoodSeed:
		return 10
	case FoodBerry:
		return 15
	default:
		return 20
	}
}

func getFoodColor(foodType FoodType) utils.Color {
	switch foodType {
	case FoodApple:
		return utils.Color{R: 255, G: 0, B: 0, A: 255} // Red
	case FoodCarrot:
		return utils.Color{R: 255, G: 165, B: 0, A: 255} // Orange
	case FoodHoney:
		return utils.Color{R: 255, G: 215, B: 0, A: 255} // Gold
	case FoodSeed:
		return utils.Color{R: 139, G: 69, B: 19, A: 255} // Brown
	case FoodBerry:
		return utils.Color{R: 128, G: 0, B: 128, A: 255} // Purple
	default:
		return utils.Color{R: 200, G: 200, B: 200, A: 255}
	}
}

func getFoodSize(foodType FoodType) float64 {
	switch foodType {
	case FoodApple:
		return 1.0
	case FoodCarrot:
		return 0.8
	case FoodHoney:
		return 0.7
	case FoodSeed:
		return 0.4
	case FoodBerry:
		return 0.6
	default:
		return 1.0
	}
}
