package creature

import (
	"math/rand"
)

// Breed creates a new creature from two parents
func Breed(parent1, parent2 *Creature) *Creature {
	// Create baby at same location as parents
	baby := NewCreature(parent1.X, parent2.Y, parent1.Type)

	// Set as baby
	baby.Age = 0
	baby.AgeStage = AgeBaby
	baby.Size = 0.7

	// Inherit genetics
	baby.Genetics = Combine(parent1.Genetics, parent2.Genetics)
	baby.applyGenetics()

	// Inherit some neural network weights from parents
	inheritBrain(baby.Brain, parent1.Brain, parent2.Brain)

	// Reset metabolism for baby
	baby.Metabolism = NewMetabolism()
	baby.Metabolism.Health = 80 // Babies start with less than full health
	baby.Metabolism.Energy = 60 // And less energy
	baby.Metabolism.Hunger = 50 // But moderately hungry

	// Inherit some personality traits
	baby.Emotions = NewEmotions()
	baby.Emotions.BaseHappiness = (parent1.Emotions.BaseHappiness + parent2.Emotions.BaseHappiness) / 2
	baby.Emotions.FearThreshold = (parent1.Emotions.FearThreshold + parent2.Emotions.FearThreshold) / 2
	baby.Emotions.AngerThreshold = (parent1.Emotions.AngerThreshold + parent2.Emotions.AngerThreshold) / 2

	// Start with basic skills inherited from parents
	baby.Learning = NewLearning()
	for skill := range parent1.Learning.Skills {
		parent1Skill := parent1.Learning.GetSkillLevel(skill)
		parent2Skill := parent2.Learning.GetSkillLevel(skill)
		inheritedSkill := (parent1Skill + parent2Skill) / 4 // Quarter of parent average
		baby.Learning.Skills[skill] = inheritedSkill
	}

	// Update breeding timers
	parent1.LastBreedTime = parent1.Age
	parent2.LastBreedTime = parent2.Age

	return baby
}

// inheritBrain combines neural networks from parents
func inheritBrain(childBrain, parent1Brain, parent2Brain *Brain) {
	parent1Weights := parent1Brain.GetWeights()
	parent2Weights := parent2Brain.GetWeights()

	if len(parent1Weights) != len(parent2Weights) {
		return // Incompatible brain structures
	}

	// Combine weights
	childWeights := make([][]float64, len(parent1Weights))
	for i := range parent1Weights {
		if len(parent1Weights[i]) != len(parent2Weights[i]) {
			continue
		}

		childWeights[i] = make([]float64, len(parent1Weights[i]))
		for j := range parent1Weights[i] {
			// Randomly choose from parents or average
			choice := rand.Float64()
			if choice < 0.45 {
				childWeights[i][j] = parent1Weights[i][j]
			} else if choice < 0.9 {
				childWeights[i][j] = parent2Weights[i][j]
			} else {
				// Average with small mutation
				childWeights[i][j] = (parent1Weights[i][j] + parent2Weights[i][j]) / 2
			}
		}
	}

	childBrain.SetWeights(childWeights)

	// Apply mutations
	childBrain.Mutate(0.1) // 10% mutation rate
}
