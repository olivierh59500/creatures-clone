package creature

import (
	"math"
)

// Learning manages the creature's ability to learn from experience
type Learning struct {
	// Learning parameters
	LearningRate   float64
	MemoryCapacity int
	ForgetRate     float64

	// Experience memory
	Experiences []Experience

	// Learned associations
	Associations map[string]Association

	// Skill levels (0-100)
	Skills map[string]float64

	// Learning state
	AttentionSpan   float64
	Focus           float64
	LastLearnedTime float64
}

// Experience represents a memorable event
type Experience struct {
	Situation  []float64 // Sensory input at the time
	Action     int       // Action taken
	Outcome    float64   // Reward/punishment received
	Timestamp  float64
	Importance float64 // How significant this experience was
}

// Association represents a learned connection
type Association struct {
	Stimulus string
	Response string
	Strength float64
	LastUsed float64
}

// Skill names
const (
	SkillWalking  = "walking"
	SkillEating   = "eating"
	SkillSpeaking = "speaking"
	SkillPlaying  = "playing"
	SkillSurvival = "survival"
	SkillSocial   = "social"
)

// NewLearning creates a new learning system
func NewLearning() *Learning {
	l := &Learning{
		LearningRate:   0.1,
		MemoryCapacity: 100,
		ForgetRate:     0.001,

		Experiences:  make([]Experience, 0),
		Associations: make(map[string]Association),
		Skills:       make(map[string]float64),

		AttentionSpan: 50,
		Focus:         50,
	}

	// Initialize basic skills
	l.initializeSkills()

	return l
}

// initializeSkills sets up initial skill levels
func (l *Learning) initializeSkills() {
	// Babies start with low skills
	l.Skills[SkillWalking] = 10
	l.Skills[SkillEating] = 20
	l.Skills[SkillSpeaking] = 5
	l.Skills[SkillPlaying] = 15
	l.Skills[SkillSurvival] = 10
	l.Skills[SkillSocial] = 10
}

// Update processes learning over time
func (l *Learning) Update(brain *Brain, recentActions []int) {
	// Process forgetting
	l.forget()

	// Update attention and focus
	l.updateAttention()

	// Consolidate recent experiences
	l.consolidateMemories()
}

// LearnFromExperience records and learns from an experience
func (l *Learning) LearnFromExperience(situation []float64, action int, outcome float64) {
	// Calculate importance based on outcome magnitude
	importance := math.Abs(outcome)

	// Create experience
	exp := Experience{
		Situation:  situation,
		Action:     action,
		Outcome:    outcome,
		Timestamp:  0, // Would use game time
		Importance: importance,
	}

	// Add to memory
	l.addExperience(exp)

	// Learn from the experience if focused enough
	if l.Focus > 30 {
		l.updateBrainFromExperience(exp)

		// Update relevant skill
		l.updateSkillFromAction(action, outcome)
	}
}

// addExperience adds an experience to memory
func (l *Learning) addExperience(exp Experience) {
	l.Experiences = append(l.Experiences, exp)

	// Limit memory size
	if len(l.Experiences) > l.MemoryCapacity {
		// Remove least important old experience
		leastImportant := 0
		minImportance := l.Experiences[0].Importance

		for i, e := range l.Experiences[:len(l.Experiences)/2] { // Only check older half
			if e.Importance < minImportance {
				leastImportant = i
				minImportance = e.Importance
			}
		}

		// Remove the least important experience
		l.Experiences = append(l.Experiences[:leastImportant], l.Experiences[leastImportant+1:]...)
	}
}

// updateBrainFromExperience adjusts neural weights based on experience
func (l *Learning) updateBrainFromExperience(exp Experience) {
	// This would integrate with the brain's learning mechanism
	// For now, it's a placeholder for the learning algorithm
}

// updateSkillFromAction improves skills based on actions
func (l *Learning) updateSkillFromAction(action int, outcome float64) {
	skillUpdate := l.LearningRate * outcome

	switch action {
	case OutputMoveLeft, OutputMoveRight, OutputJump:
		l.improveSkill(SkillWalking, skillUpdate)
	case OutputEat:
		l.improveSkill(SkillEating, skillUpdate)
	case OutputPlay:
		l.improveSkill(SkillPlaying, skillUpdate)
	case OutputSpeak:
		l.improveSkill(SkillSpeaking, skillUpdate)
		l.improveSkill(SkillSocial, skillUpdate*0.5)
	case OutputBreed:
		l.improveSkill(SkillSocial, skillUpdate)
	}

	// Survival skill improves with positive outcomes
	if outcome > 0 {
		l.improveSkill(SkillSurvival, skillUpdate*0.3)
	}
}

// improveSkill increases a skill level
func (l *Learning) improveSkill(skill string, amount float64) {
	current := l.Skills[skill]

	// Skills improve faster at lower levels
	improvementRate := amount * (1.0 - current/100.0)

	l.Skills[skill] = math.Min(100, current+improvementRate)
}

// LearnAssociation creates or strengthens an association
func (l *Learning) LearnAssociation(stimulus, response string, success bool) {
	key := stimulus + "->" + response

	strength := 0.1
	if !success {
		strength = -0.1
	}

	if assoc, exists := l.Associations[key]; exists {
		// Strengthen existing association
		assoc.Strength += strength * l.LearningRate
		assoc.Strength = math.Max(-1, math.Min(1, assoc.Strength))
		assoc.LastUsed = 0 // Reset usage timer
		l.Associations[key] = assoc
	} else {
		// Create new association
		l.Associations[key] = Association{
			Stimulus: stimulus,
			Response: response,
			Strength: strength,
			LastUsed: 0,
		}
	}
}

// GetAssociationStrength returns how strongly two concepts are linked
func (l *Learning) GetAssociationStrength(stimulus, response string) float64 {
	key := stimulus + "->" + response
	if assoc, exists := l.Associations[key]; exists {
		return assoc.Strength
	}
	return 0
}

// forget processes memory decay
func (l *Learning) forget() {
	// Decay associations
	for key, assoc := range l.Associations {
		assoc.Strength *= (1.0 - l.ForgetRate)
		assoc.LastUsed += 1

		// Remove very weak or unused associations
		if math.Abs(assoc.Strength) < 0.01 || assoc.LastUsed > 1000 {
			delete(l.Associations, key)
		} else {
			l.Associations[key] = assoc
		}
	}

	// Decay experience importance
	for i := range l.Experiences {
		l.Experiences[i].Importance *= (1.0 - l.ForgetRate)
	}
}

// updateAttention manages attention and focus
func (l *Learning) updateAttention() {
	// Attention naturally wanders
	l.AttentionSpan -= 0.1

	// Focus decreases without stimulation
	l.Focus -= 0.05

	// Clamp values
	l.AttentionSpan = math.Max(0, math.Min(100, l.AttentionSpan))
	l.Focus = math.Max(0, math.Min(100, l.Focus))
}

// PayAttention increases focus on current activity
func (l *Learning) PayAttention(amount float64) {
	l.Focus = math.Min(100, l.Focus+amount)
	l.AttentionSpan = math.Min(100, l.AttentionSpan+amount*0.5)
}

// consolidateMemories strengthens important memories
func (l *Learning) consolidateMemories() {
	// During high focus, important experiences are strengthened
	if l.Focus > 70 {
		for i := range l.Experiences {
			if l.Experiences[i].Importance > 0.5 {
				l.Experiences[i].Importance *= 1.01 // Slight strengthening
			}
		}
	}
}

// GetSkillLevel returns the current level of a skill
func (l *Learning) GetSkillLevel(skill string) float64 {
	if level, exists := l.Skills[skill]; exists {
		return level
	}
	return 0
}

// CanLearn checks if the creature is in a good state to learn
func (l *Learning) CanLearn() bool {
	return l.Focus > 20 && l.AttentionSpan > 10
}

// GetMemoryUsage returns how full the memory is (0-1)
func (l *Learning) GetMemoryUsage() float64 {
	return float64(len(l.Experiences)) / float64(l.MemoryCapacity)
}

// RecallSimilarExperience finds a past experience similar to current situation
func (l *Learning) RecallSimilarExperience(currentSituation []float64) *Experience {
	if len(l.Experiences) == 0 {
		return nil
	}

	var bestMatch *Experience
	bestSimilarity := 0.0

	for i := range l.Experiences {
		exp := &l.Experiences[i]
		similarity := l.calculateSimilarity(currentSituation, exp.Situation)

		// Weight by importance and recency
		similarity *= exp.Importance

		if similarity > bestSimilarity {
			bestSimilarity = similarity
			bestMatch = exp
		}
	}

	// Only return if similarity is high enough
	if bestSimilarity > 0.7 {
		return bestMatch
	}

	return nil
}

// calculateSimilarity computes similarity between two situations
func (l *Learning) calculateSimilarity(situation1, situation2 []float64) float64 {
	if len(situation1) != len(situation2) {
		return 0
	}

	// Cosine similarity
	dotProduct := 0.0
	magnitude1 := 0.0
	magnitude2 := 0.0

	for i := range situation1 {
		dotProduct += situation1[i] * situation2[i]
		magnitude1 += situation1[i] * situation1[i]
		magnitude2 += situation2[i] * situation2[i]
	}

	if magnitude1 == 0 || magnitude2 == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(magnitude1) * math.Sqrt(magnitude2))
}
