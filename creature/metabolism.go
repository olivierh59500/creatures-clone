package creature

import (
	"github.com/olivierh59500/creatures-clone/utils"
)

// Metabolism manages the creature's physical needs and health
type Metabolism struct {
	// Core stats (0-100)
	Health float64
	Hunger float64
	Energy float64

	// Rates
	HungerRate  float64 // How fast hunger increases
	EnergyRate  float64 // How fast energy depletes
	HealingRate float64 // How fast health recovers

	// Chemical levels (simplified)
	Glucose    float64 // From food
	Toxins     float64 // Harmful substances
	Endorphins float64 // Natural happiness chemicals
	Adrenaline float64 // Stress/excitement

	// Status tracking
	LastMealTime   float64
	LastSleepTime  float64
	TotalFoodEaten int
}

// NewMetabolism creates a new metabolism system
func NewMetabolism() *Metabolism {
	return &Metabolism{
		Health: 100,
		Hunger: 30, // Start slightly hungry
		Energy: 80,

		HungerRate:  0.05, // Hunger increases by 0.05 per update
		EnergyRate:  0.03, // Energy decreases by 0.03 per update
		HealingRate: 0.02, // Health recovers by 0.02 per update when fed

		Glucose:    50,
		Toxins:     0,
		Endorphins: 30,
		Adrenaline: 10,
	}
}

// Update processes metabolic changes
func (m *Metabolism) Update(activityLevel float64) {
	// Increase hunger over time
	m.Hunger = utils.Clamp(m.Hunger+m.HungerRate, 0, 100)

	// Energy depletion based on activity
	energyLoss := m.EnergyRate * (1 + activityLevel)
	m.Energy = utils.Clamp(m.Energy-energyLoss, 0, 100)

	// Process chemicals
	m.processChemicals()

	// Health effects from hunger and energy
	if m.Hunger > 80 {
		// Starvation damage
		m.Health -= 0.1
	} else if m.Hunger < 50 && m.Energy > 30 {
		// Natural healing when fed and rested
		m.Health = utils.Clamp(m.Health+m.HealingRate, 0, 100)
	}

	if m.Energy < 20 {
		// Exhaustion damage
		m.Health -= 0.05
	}

	// Toxin damage
	if m.Toxins > 50 {
		m.Health -= m.Toxins * 0.001
	}

	// Ensure health stays in bounds
	m.Health = utils.Clamp(m.Health, 0, 100)
}

// processChemicals updates chemical levels
func (m *Metabolism) processChemicals() {
	// Glucose consumption
	if m.Glucose > 0 {
		// Convert glucose to energy
		glucoseUsed := utils.Min(m.Glucose, 0.1)
		m.Glucose -= glucoseUsed
		m.Energy = utils.Clamp(m.Energy+glucoseUsed*2, 0, 100)

		// Reduce hunger when glucose is available
		m.Hunger = utils.Clamp(m.Hunger-glucoseUsed*3, 0, 100)
	}

	// Toxin processing
	if m.Toxins > 0 {
		// Liver processes toxins
		m.Toxins = utils.Clamp(m.Toxins-0.05, 0, 100)
	}

	// Endorphin decay
	m.Endorphins = utils.Clamp(m.Endorphins-0.02, 0, 100)

	// Adrenaline decay
	m.Adrenaline = utils.Clamp(m.Adrenaline-0.03, 0, 100)
}

// Eat processes food consumption
func (m *Metabolism) Eat(nutritionValue float64) {
	// Add glucose from food
	m.Glucose = utils.Clamp(m.Glucose+nutritionValue, 0, 100)

	// Immediate hunger reduction
	m.Hunger = utils.Clamp(m.Hunger-nutritionValue*0.5, 0, 100)

	// Small endorphin boost from eating
	m.Endorphins = utils.Clamp(m.Endorphins+5, 0, 100)

	// Track eating
	m.LastMealTime = 0 // Reset meal timer
	m.TotalFoodEaten++
}

// Sleep processes rest and recovery
func (m *Metabolism) Sleep() {
	// Energy recovery during sleep
	m.Energy = utils.Clamp(m.Energy+0.2, 0, 100)

	// Enhanced healing during sleep
	if m.Health < 100 && m.Hunger < 70 {
		m.Health = utils.Clamp(m.Health+m.HealingRate*2, 0, 100)
	}

	// Process toxins faster during sleep
	if m.Toxins > 0 {
		m.Toxins = utils.Clamp(m.Toxins-0.1, 0, 100)
	}

	m.LastSleepTime = 0 // Reset sleep timer
}

// IngestToxin adds harmful substances
func (m *Metabolism) IngestToxin(amount float64) {
	m.Toxins = utils.Clamp(m.Toxins+amount, 0, 100)

	// Immediate health impact for large doses
	if amount > 20 {
		m.Health -= amount * 0.2
	}
}

// Exercise increases activity and metabolism
func (m *Metabolism) Exercise(intensity float64) {
	// Burn energy
	m.Energy = utils.Clamp(m.Energy-intensity*0.1, 0, 100)

	// Increase hunger
	m.Hunger = utils.Clamp(m.Hunger+intensity*0.05, 0, 100)

	// Release endorphins
	m.Endorphins = utils.Clamp(m.Endorphins+intensity*0.2, 0, 100)

	// Increase adrenaline
	m.Adrenaline = utils.Clamp(m.Adrenaline+intensity*0.3, 0, 100)
}

// GetStress calculates current stress level
func (m *Metabolism) GetStress() float64 {
	stress := 0.0

	// Hunger stress
	if m.Hunger > 60 {
		stress += (m.Hunger - 60) * 0.02
	}

	// Energy stress
	if m.Energy < 30 {
		stress += (30 - m.Energy) * 0.02
	}

	// Health stress
	if m.Health < 50 {
		stress += (50 - m.Health) * 0.03
	}

	// Toxin stress
	stress += m.Toxins * 0.01

	// Adrenaline reduces perceived stress (fight or flight)
	stress -= m.Adrenaline * 0.005

	return utils.Clamp(stress, 0, 100)
}

// GetWellbeing calculates overall wellbeing
func (m *Metabolism) GetWellbeing() float64 {
	wellbeing := m.Health * 0.4
	wellbeing += (100 - m.Hunger) * 0.3
	wellbeing += m.Energy * 0.2
	wellbeing += m.Endorphins * 0.1

	// Penalty for toxins
	wellbeing -= m.Toxins * 0.2

	return utils.Clamp(wellbeing, 0, 100)
}

// NeedsFood checks if the creature needs to eat
func (m *Metabolism) NeedsFood() bool {
	return m.Hunger > 60 || m.Glucose < 20
}

// NeedsSleep checks if the creature needs rest
func (m *Metabolism) NeedsSleep() bool {
	return m.Energy < 30
}

// IsHealthy checks if the creature is in good health
func (m *Metabolism) IsHealthy() bool {
	return m.Health > 70 && m.Toxins < 20
}

// IsCritical checks if the creature is in critical condition
func (m *Metabolism) IsCritical() bool {
	return m.Health < 20 || m.Hunger > 90 || m.Toxins > 80
}
