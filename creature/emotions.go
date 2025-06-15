package creature

import (
	"math"

	"github.com/olivierh59500/creatures-clone/utils"
)

// Emotions represents the creature's emotional state
type Emotions struct {
	// Primary emotions (-100 to 100)
	Happiness float64
	Fear      float64
	Anger     float64
	Curiosity float64

	// Secondary emotions
	Loneliness float64
	Boredom    float64
	Love       float64
	Jealousy   float64

	// Emotional parameters
	BaseHappiness    float64 // Genetic happiness baseline
	EmotionalInertia float64 // How quickly emotions change

	// Thresholds
	FearThreshold  float64
	AngerThreshold float64

	// Social bonds (creature ID -> bond strength)
	SocialBonds map[string]float64

	// Recent emotional events
	RecentEvents []EmotionalEvent
}

// EmotionalEvent represents something that affected emotions
type EmotionalEvent struct {
	Type      string
	Intensity float64
	Timestamp float64
}

// NewEmotions creates a new emotion system
func NewEmotions() *Emotions {
	return &Emotions{
		Happiness: 30, // Start neutral-positive
		Fear:      0,
		Anger:     0,
		Curiosity: 50, // Start curious

		Loneliness: 20,
		Boredom:    30,
		Love:       0,
		Jealousy:   0,

		BaseHappiness:    0,
		EmotionalInertia: 0.9, // Emotions change gradually

		FearThreshold:  50,
		AngerThreshold: 60,

		SocialBonds:  make(map[string]float64),
		RecentEvents: make([]EmotionalEvent, 0, 10),
	}
}

// Update processes emotional changes
func (e *Emotions) Update(metabolism *Metabolism, brainOutput []float64) {
	// Apply emotional inertia (emotions don't change instantly)
	e.applyInertia()

	// Update based on physical state
	e.updateFromMetabolism(metabolism)

	// Update based on brain activity
	e.updateFromBrainActivity(brainOutput)

	// Process secondary emotions
	e.updateSecondaryEmotions()

	// Apply baseline drift
	e.applyBaselineDrift()

	// Clamp all emotions to valid ranges
	e.clampEmotions()

	// Age out old events
	e.cleanOldEvents()
}

// applyInertia makes emotions change gradually
func (e *Emotions) applyInertia() {
	inertia := e.EmotionalInertia

	// Primary emotions decay towards baseline
	e.Happiness = e.Happiness * inertia
	e.Fear = e.Fear * inertia
	e.Anger = e.Anger * inertia
	e.Curiosity = e.Curiosity * inertia

	// Secondary emotions decay faster
	fasterInertia := inertia * 0.9
	e.Loneliness = e.Loneliness * fasterInertia
	e.Boredom = e.Boredom * fasterInertia
	e.Love = e.Love * fasterInertia
	e.Jealousy = e.Jealousy * fasterInertia
}

// updateFromMetabolism adjusts emotions based on physical state
func (e *Emotions) updateFromMetabolism(m *Metabolism) {
	// Hunger affects happiness
	if m.Hunger > 70 {
		e.AdjustHappiness(-5)
		e.AdjustAnger(3)
	}

	// Low energy increases boredom and reduces curiosity
	if m.Energy < 30 {
		e.Boredom += 2
		e.Curiosity -= 3
	}

	// Poor health causes fear and unhappiness
	if m.Health < 50 {
		e.AdjustFear((50 - m.Health) * 0.2)
		e.AdjustHappiness(-3)
	}

	// Endorphins boost happiness
	e.AdjustHappiness(m.Endorphins * 0.1)

	// Adrenaline affects fear and anger
	if m.Adrenaline > 50 {
		e.AdjustFear(m.Adrenaline * 0.05)
		e.AdjustAnger(m.Adrenaline * 0.03)
	}
}

// updateFromBrainActivity adjusts emotions based on actions
func (e *Emotions) updateFromBrainActivity(brainOutput []float64) {
	// Playing increases happiness
	if brainOutput[OutputPlay] > 0.5 {
		e.AdjustHappiness(5)
		e.Boredom = utils.Clamp(e.Boredom-10, -100, 100)
	}

	// Social interaction reduces loneliness
	if brainOutput[OutputSpeak] > 0.5 {
		e.Loneliness = utils.Clamp(e.Loneliness-5, -100, 100)
		e.AdjustHappiness(2)
	}

	// Movement and exploration increase curiosity
	movement := math.Max(brainOutput[OutputMoveLeft], brainOutput[OutputMoveRight])
	if movement > 0.5 {
		e.Curiosity = utils.Clamp(e.Curiosity+1, -100, 100)
		e.Boredom = utils.Clamp(e.Boredom-2, -100, 100)
	}
}

// updateSecondaryEmotions processes complex emotional states
func (e *Emotions) updateSecondaryEmotions() {
	// Loneliness increases over time without social contact
	e.Loneliness = utils.Clamp(e.Loneliness+0.1, -100, 100)

	// Boredom increases without stimulation
	if math.Abs(e.Curiosity) < 20 {
		e.Boredom = utils.Clamp(e.Boredom+0.2, -100, 100)
	}

	// Love develops from positive social bonds
	maxBond := 0.0
	for _, bond := range e.SocialBonds {
		if bond > maxBond {
			maxBond = bond
		}
	}
	e.Love = utils.Clamp(maxBond*100, -100, 100)

	// High happiness reduces negative emotions
	if e.Happiness > 50 {
		e.Fear *= 0.95
		e.Anger *= 0.95
		e.Loneliness *= 0.95
	}

	// High fear suppresses other emotions
	if e.Fear > 70 {
		e.Curiosity *= 0.8
		e.Happiness *= 0.9
	}
}

// applyBaselineDrift slowly returns emotions to genetic baseline
func (e *Emotions) applyBaselineDrift() {
	driftRate := 0.01

	// Happiness drifts towards genetic baseline
	targetHappiness := e.BaseHappiness
	e.Happiness += (targetHappiness - e.Happiness) * driftRate

	// Other emotions drift towards zero
	e.Fear += (0 - e.Fear) * driftRate
	e.Anger += (0 - e.Anger) * driftRate
}

// clampEmotions ensures all emotions are within valid ranges
func (e *Emotions) clampEmotions() {
	e.Happiness = utils.Clamp(e.Happiness, -100, 100)
	e.Fear = utils.Clamp(e.Fear, -100, 100)
	e.Anger = utils.Clamp(e.Anger, -100, 100)
	e.Curiosity = utils.Clamp(e.Curiosity, -100, 100)

	e.Loneliness = utils.Clamp(e.Loneliness, -100, 100)
	e.Boredom = utils.Clamp(e.Boredom, -100, 100)
	e.Love = utils.Clamp(e.Love, -100, 100)
	e.Jealousy = utils.Clamp(e.Jealousy, -100, 100)
}

// AdjustHappiness modifies happiness with event tracking
func (e *Emotions) AdjustHappiness(amount float64) {
	e.Happiness = utils.Clamp(e.Happiness+amount, -100, 100)

	if math.Abs(amount) > 5 {
		e.addEvent("happiness", amount)
	}
}

// AdjustFear modifies fear with threshold checking
func (e *Emotions) AdjustFear(amount float64) {
	e.Fear = utils.Clamp(e.Fear+amount, -100, 100)

	if e.Fear > e.FearThreshold {
		// Trigger fear response
		e.addEvent("fear_response", e.Fear)
	}
}

// AdjustAnger modifies anger with threshold checking
func (e *Emotions) AdjustAnger(amount float64) {
	e.Anger = utils.Clamp(e.Anger+amount, -100, 100)

	if e.Anger > e.AngerThreshold {
		// Trigger anger response
		e.addEvent("anger_response", e.Anger)
	}
}

// AdjustCuriosity modifies curiosity
func (e *Emotions) AdjustCuriosity(amount float64) {
	e.Curiosity = utils.Clamp(e.Curiosity+amount, -100, 100)
}

// UpdateSocialBond updates relationship with another creature
func (e *Emotions) UpdateSocialBond(creatureID string, interaction float64) {
	current := e.SocialBonds[creatureID]
	e.SocialBonds[creatureID] = utils.Clamp(current+interaction, -1, 1)

	// Positive interactions reduce loneliness
	if interaction > 0 {
		e.Loneliness = utils.Clamp(e.Loneliness-interaction*10, -100, 100)
	}
}

// GetDominantEmotion returns the strongest current emotion
func (e *Emotions) GetDominantEmotion() string {
	emotions := map[string]float64{
		"happy":   math.Abs(e.Happiness),
		"afraid":  math.Abs(e.Fear),
		"angry":   math.Abs(e.Anger),
		"curious": math.Abs(e.Curiosity),
		"lonely":  math.Abs(e.Loneliness),
		"bored":   math.Abs(e.Boredom),
		"loving":  math.Abs(e.Love),
		"jealous": math.Abs(e.Jealousy),
	}

	maxEmotion := "neutral"
	maxValue := 20.0 // Threshold for "neutral"

	for emotion, value := range emotions {
		if value > maxValue {
			maxValue = value
			maxEmotion = emotion
		}
	}

	return maxEmotion
}

// GetMood returns overall mood (-1 to 1)
func (e *Emotions) GetMood() float64 {
	// Positive emotions
	positive := e.Happiness + e.Curiosity + e.Love

	// Negative emotions
	negative := e.Fear + e.Anger + e.Loneliness + e.Boredom + e.Jealousy

	mood := (positive - negative) / 500.0 // Normalize to -1 to 1
	return utils.Clamp(mood, -1, 1)
}

// addEvent records an emotional event
func (e *Emotions) addEvent(eventType string, intensity float64) {
	event := EmotionalEvent{
		Type:      eventType,
		Intensity: intensity,
		Timestamp: 0, // Would use game time
	}

	e.RecentEvents = append(e.RecentEvents, event)

	// Keep only recent events
	if len(e.RecentEvents) > 10 {
		e.RecentEvents = e.RecentEvents[1:]
	}
}

// cleanOldEvents removes old emotional events
func (e *Emotions) cleanOldEvents() {
	// In a full implementation, would check timestamps
	// For now, just limit the array size
	if len(e.RecentEvents) > 10 {
		e.RecentEvents = e.RecentEvents[len(e.RecentEvents)-10:]
	}
}
