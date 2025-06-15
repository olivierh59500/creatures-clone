package creature

import (
	"math"
	"math/rand"

	"github.com/olivierh59500/creatures-clone/utils"
)

// Genetics represents a creature's genetic makeup
type Genetics struct {
	// Genes map with trait names and values (0.0 to 1.0)
	Genes map[string]float64

	// Appearance genes
	ColorR, ColorG, ColorB float64
	Pattern                string

	// Dominant/recessive tracking
	DominantGenes map[string]bool
}

// Gene trait names
const (
	GeneMetabolismRate = "metabolism_rate"
	GeneMovementSpeed  = "movement_speed"
	GeneLearningRate   = "learning_rate"
	GeneLifespan       = "lifespan"
	GeneFertility      = "fertility"
	GeneStrength       = "strength"
	GeneHappinessBias  = "happiness_bias"
	GeneFearThreshold  = "fear_threshold"
	GeneAngerThreshold = "anger_threshold"
	GeneCuriosity      = "curiosity"
	GeneSociability    = "sociability"
	GeneAggression     = "aggression"
)

// NewGenetics creates a new genetics instance
func NewGenetics() *Genetics {
	g := &Genetics{
		Genes:         make(map[string]float64),
		DominantGenes: make(map[string]bool),
	}

	// Initialize default genes
	g.initializeDefaultGenes()

	return g
}

// initializeDefaultGenes sets up default gene values
func (g *Genetics) initializeDefaultGenes() {
	// Set default values (0.5 is neutral)
	defaultGenes := map[string]float64{
		GeneMetabolismRate: 0.5,
		GeneMovementSpeed:  0.5,
		GeneLearningRate:   0.5,
		GeneLifespan:       0.5,
		GeneFertility:      0.5,
		GeneStrength:       0.5,
		GeneHappinessBias:  0.5,
		GeneFearThreshold:  0.5,
		GeneAngerThreshold: 0.5,
		GeneCuriosity:      0.5,
		GeneSociability:    0.5,
		GeneAggression:     0.5,
	}

	for gene, value := range defaultGenes {
		g.Genes[gene] = value
		g.DominantGenes[gene] = rand.Float64() > 0.5
	}

	// Default appearance
	g.ColorR = 0.5
	g.ColorG = 0.7
	g.ColorB = 0.3
	g.Pattern = "solid"
}

// Randomize creates random genetic values
func (g *Genetics) Randomize() {
	// Randomize trait genes
	for gene := range g.Genes {
		g.Genes[gene] = rand.Float64()
		g.DominantGenes[gene] = rand.Float64() > 0.5
	}

	// Randomize appearance
	g.randomizeAppearance()
}

// randomizeAppearance sets random appearance genes
func (g *Genetics) randomizeAppearance() {
	// Pick a random color scheme
	colorSchemes := []struct {
		r, g, b float64
		name    string
	}{
		{0.13, 0.55, 0.13, "forest"},   // Forest green
		{0.82, 0.41, 0.12, "desert"},   // Desert tan
		{0.27, 0.51, 0.71, "ocean"},    // Ocean blue
		{0.58, 0.44, 0.86, "mountain"}, // Mountain purple
		{0.86, 0.44, 0.58, "flower"},   // Flower pink
		{0.96, 0.87, 0.70, "sunshine"}, // Sunshine yellow
	}

	scheme := colorSchemes[rand.Intn(len(colorSchemes))]

	// Add some variation
	g.ColorR = utils.Clamp(scheme.r+rand.Float64()*0.2-0.1, 0, 1)
	g.ColorG = utils.Clamp(scheme.g+rand.Float64()*0.2-0.1, 0, 1)
	g.ColorB = utils.Clamp(scheme.b+rand.Float64()*0.2-0.1, 0, 1)

	// Random pattern
	patterns := []string{"solid", "spotted", "striped"}
	g.Pattern = patterns[rand.Intn(len(patterns))]
}

// Combine creates offspring genetics from two parents
func Combine(parent1, parent2 *Genetics) *Genetics {
	child := NewGenetics()

	// Combine trait genes
	for gene := range parent1.Genes {
		// Mendelian inheritance with dominance
		p1Value := parent1.Genes[gene]
		p2Value := parent2.Genes[gene]
		p1Dominant := parent1.DominantGenes[gene]
		p2Dominant := parent2.DominantGenes[gene]

		// Determine which genes are expressed
		if p1Dominant && p2Dominant {
			// Both dominant - average
			child.Genes[gene] = (p1Value + p2Value) / 2
			child.DominantGenes[gene] = true
		} else if p1Dominant && !p2Dominant {
			// P1 dominant
			child.Genes[gene] = p1Value
			child.DominantGenes[gene] = rand.Float64() > 0.25 // 75% chance dominant
		} else if !p1Dominant && p2Dominant {
			// P2 dominant
			child.Genes[gene] = p2Value
			child.DominantGenes[gene] = rand.Float64() > 0.25
		} else {
			// Both recessive - express recessive
			child.Genes[gene] = (p1Value + p2Value) / 2
			child.DominantGenes[gene] = false
		}
	}

	// Combine appearance genes
	child.ColorR = (parent1.ColorR + parent2.ColorR) / 2
	child.ColorG = (parent1.ColorG + parent2.ColorG) / 2
	child.ColorB = (parent1.ColorB + parent2.ColorB) / 2

	// Pattern inheritance (simplified)
	if rand.Float64() > 0.5 {
		child.Pattern = parent1.Pattern
	} else {
		child.Pattern = parent2.Pattern
	}

	// Apply mutations
	child.Mutate()

	return child
}

// Mutate applies random mutations to genes
func (g *Genetics) Mutate() {
	mutationRate := 0.1
	mutationStrength := 0.1

	// Mutate trait genes
	for gene := range g.Genes {
		if rand.Float64() < mutationRate {
			// Apply mutation
			change := (rand.Float64()*2 - 1) * mutationStrength
			g.Genes[gene] = utils.Clamp(g.Genes[gene]+change, 0, 1)

			// Small chance to flip dominance
			if rand.Float64() < 0.05 {
				g.DominantGenes[gene] = !g.DominantGenes[gene]
			}
		}
	}

	// Mutate appearance
	if rand.Float64() < mutationRate {
		g.ColorR = utils.Clamp(g.ColorR+(rand.Float64()*2-1)*mutationStrength, 0, 1)
		g.ColorG = utils.Clamp(g.ColorG+(rand.Float64()*2-1)*mutationStrength, 0, 1)
		g.ColorB = utils.Clamp(g.ColorB+(rand.Float64()*2-1)*mutationStrength, 0, 1)
	}

	// Rare pattern mutation
	if rand.Float64() < 0.02 {
		patterns := []string{"solid", "spotted", "striped"}
		g.Pattern = patterns[rand.Intn(len(patterns))]
	}
}

// GetColor returns the creature's color based on genetics
func (g *Genetics) GetColor() utils.Color {
	return utils.Color{
		R: uint8(g.ColorR * 255),
		G: uint8(g.ColorG * 255),
		B: uint8(g.ColorB * 255),
		A: 255,
	}
}

// GetTrait returns the value of a specific gene
func (g *Genetics) GetTrait(trait string) float64 {
	if value, exists := g.Genes[trait]; exists {
		return value
	}
	return 0.5 // Default neutral value
}

// SetTrait sets the value of a specific gene
func (g *Genetics) SetTrait(trait string, value float64) {
	g.Genes[trait] = utils.Clamp(value, 0, 1)
}

// GetDominance returns whether a gene is dominant
func (g *Genetics) GetDominance(trait string) bool {
	if dominant, exists := g.DominantGenes[trait]; exists {
		return dominant
	}
	return false
}

// Clone creates a copy of the genetics
func (g *Genetics) Clone() *Genetics {
	clone := NewGenetics()

	// Copy genes
	for gene, value := range g.Genes {
		clone.Genes[gene] = value
		clone.DominantGenes[gene] = g.DominantGenes[gene]
	}

	// Copy appearance
	clone.ColorR = g.ColorR
	clone.ColorG = g.ColorG
	clone.ColorB = g.ColorB
	clone.Pattern = g.Pattern

	return clone
}

// Similarity calculates genetic similarity between two creatures (0.0 to 1.0)
func (g *Genetics) Similarity(other *Genetics) float64 {
	if other == nil {
		return 0
	}

	totalDiff := 0.0
	count := 0

	// Compare trait genes
	for gene, value := range g.Genes {
		if otherValue, exists := other.Genes[gene]; exists {
			diff := math.Abs(value - otherValue)
			totalDiff += diff
			count++
		}
	}

	// Compare appearance
	totalDiff += math.Abs(g.ColorR - other.ColorR)
	totalDiff += math.Abs(g.ColorG - other.ColorG)
	totalDiff += math.Abs(g.ColorB - other.ColorB)
	count += 3

	if count == 0 {
		return 1.0
	}

	avgDiff := totalDiff / float64(count)
	return 1.0 - avgDiff
}
