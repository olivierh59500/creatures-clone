package creature

import (
	"math/rand"
	"strings"
)

// Language manages the creature's vocabulary and communication
type Language struct {
	// Vocabulary mapping words to concepts
	Vocabulary map[string]Concept

	// Recent words heard
	RecentWords []string

	// Language parameters
	SpeechClarity   float64 // How clearly the creature speaks
	Comprehension   float64 // How well it understands speech
	VocabularyLimit int     // Maximum words it can remember

	// Current speech
	CurrentWord string
	SpeechTimer float64
}

// Concept represents what a word means to the creature
type Concept struct {
	Word         string
	ObjectType   string   // What type of object this refers to
	Associations []string // Related concepts
	Confidence   float64  // How sure the creature is about this word
	TimesUsed    int
	LastUsed     float64
}

// NewLanguage creates a new language system
func NewLanguage() *Language {
	return &Language{
		Vocabulary:      make(map[string]Concept),
		RecentWords:     make([]string, 0, 5),
		SpeechClarity:   0.5,
		Comprehension:   0.5,
		VocabularyLimit: 50,
	}
}

// HearWord processes a heard word and tries to learn it
func (l *Language) HearWord(word string, context interface{}) {
	// Normalize word
	word = strings.ToLower(strings.TrimSpace(word))
	if word == "" {
		return
	}

	// Add to recent words
	l.addRecentWord(word)

	// Try to understand based on context
	if context != nil {
		l.learnWordFromContext(word, context)
	}
}

// learnWordFromContext associates a word with an object or situation
func (l *Language) learnWordFromContext(word string, context interface{}) {
	// Determine object type from context
	// In a full implementation, this would use type assertion
	objectType := "unknown"

	// Check if we already know this word
	if concept, exists := l.Vocabulary[word]; exists {
		// Reinforce existing knowledge
		concept.Confidence = min(1.0, concept.Confidence+0.1)
		concept.TimesUsed++
		concept.LastUsed = 0
		l.Vocabulary[word] = concept
	} else {
		// Learn new word if under vocabulary limit
		if len(l.Vocabulary) < l.VocabularyLimit {
			l.Vocabulary[word] = Concept{
				Word:         word,
				ObjectType:   objectType,
				Associations: []string{},
				Confidence:   0.3 * l.Comprehension,
				TimesUsed:    1,
				LastUsed:     0,
			}
		} else {
			// Replace least used word
			l.replaceLeastUsedWord(word, objectType)
		}
	}
}

// replaceLeastUsedWord removes the least used word to make room
func (l *Language) replaceLeastUsedWord(newWord, objectType string) {
	var leastUsed string
	minUsage := int(^uint(0) >> 1) // Max int

	for word, concept := range l.Vocabulary {
		if concept.TimesUsed < minUsage {
			minUsage = concept.TimesUsed
			leastUsed = word
		}
	}

	if leastUsed != "" {
		delete(l.Vocabulary, leastUsed)
		l.Vocabulary[newWord] = Concept{
			Word:       newWord,
			ObjectType: objectType,
			Confidence: 0.3 * l.Comprehension,
			TimesUsed:  1,
			LastUsed:   0,
		}
	}
}

// Speak attempts to say a word based on current thoughts
func (l *Language) Speak(thought string) string {
	// Check if we know a word for this thought
	for word, concept := range l.Vocabulary {
		if concept.ObjectType == thought && concept.Confidence > 0.5 {
			// Add some speech imperfection based on clarity
			if rand.Float64() > l.SpeechClarity {
				return l.garbleWord(word)
			}

			l.CurrentWord = word
			l.SpeechTimer = 1.0 // Speech duration

			// Update usage
			concept.TimesUsed++
			concept.LastUsed = 0
			l.Vocabulary[word] = concept

			return word
		}
	}

	// Babble if we don't know the word
	return l.babble()
}

// garbleWord introduces speech errors
func (l *Language) garbleWord(word string) string {
	if len(word) < 2 {
		return word
	}

	// Random speech errors
	errorType := rand.Intn(3)
	switch errorType {
	case 0: // Repeat syllable
		mid := len(word) / 2
		return word[:mid] + word[mid-1:mid+1] + word[mid:]
	case 1: // Drop letter
		pos := rand.Intn(len(word))
		return word[:pos] + word[pos+1:]
	case 2: // Add random sound
		sounds := []string{"um", "ah", "er"}
		return sounds[rand.Intn(len(sounds))] + word
	}

	return word
}

// babble produces random baby talk
func (l *Language) babble() string {
	consonants := []string{"b", "d", "g", "m", "n", "p", "t", "w", "y"}
	vowels := []string{"a", "e", "i", "o", "u"}

	// Create simple CV or CVCV pattern
	pattern := rand.Intn(2)

	word := ""
	switch pattern {
	case 0: // CV
		word = consonants[rand.Intn(len(consonants))] +
			vowels[rand.Intn(len(vowels))]
	case 1: // CVCV (repetition common in baby talk)
		cv := consonants[rand.Intn(len(consonants))] +
			vowels[rand.Intn(len(vowels))]
		word = cv + cv
	}

	return word
}

// addRecentWord adds a word to recent memory
func (l *Language) addRecentWord(word string) {
	l.RecentWords = append(l.RecentWords, word)
	if len(l.RecentWords) > 5 {
		l.RecentWords = l.RecentWords[1:]
	}
}

// GetKnownWords returns all words the creature knows
func (l *Language) GetKnownWords() []string {
	words := make([]string, 0, len(l.Vocabulary))
	for word := range l.Vocabulary {
		words = append(words, word)
	}
	return words
}

// GetVocabularySize returns the number of known words
func (l *Language) GetVocabularySize() int {
	return len(l.Vocabulary)
}

// KnowsWord checks if the creature knows a specific word
func (l *Language) KnowsWord(word string) bool {
	word = strings.ToLower(strings.TrimSpace(word))
	concept, exists := l.Vocabulary[word]
	return exists && concept.Confidence > 0.3
}

// GetWordConfidence returns how well the creature knows a word
func (l *Language) GetWordConfidence(word string) float64 {
	word = strings.ToLower(strings.TrimSpace(word))
	if concept, exists := l.Vocabulary[word]; exists {
		return concept.Confidence
	}
	return 0
}

// TeachWord explicitly teaches a word with high confidence
func (l *Language) TeachWord(word, objectType string) {
	word = strings.ToLower(strings.TrimSpace(word))

	l.Vocabulary[word] = Concept{
		Word:         word,
		ObjectType:   objectType,
		Confidence:   0.8, // High confidence for taught words
		TimesUsed:    0,
		LastUsed:     0,
		Associations: []string{},
	}

	// Ensure vocabulary limit
	if len(l.Vocabulary) > l.VocabularyLimit {
		l.replaceLeastUsedWord("", "")
	}
}

// Update processes language updates
func (l *Language) Update() {
	// Update speech timer
	if l.SpeechTimer > 0 {
		l.SpeechTimer -= 0.016 // Assuming 60 FPS
		if l.SpeechTimer <= 0 {
			l.CurrentWord = ""
		}
	}

	// Age vocabulary usage
	for word, concept := range l.Vocabulary {
		concept.LastUsed += 0.016

		// Slowly forget unused words
		if concept.LastUsed > 600 { // 10 minutes
			concept.Confidence *= 0.999
			if concept.Confidence < 0.1 {
				delete(l.Vocabulary, word)
			} else {
				l.Vocabulary[word] = concept
			}
		}
	}
}

// IsSpeaking checks if the creature is currently speaking
func (l *Language) IsSpeaking() bool {
	return l.SpeechTimer > 0
}

// GetCurrentWord returns the word being spoken
func (l *Language) GetCurrentWord() string {
	return l.CurrentWord
}

// min returns the minimum of two floats
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
