package game

import (
	"github.com/olivierh59500/creatures-clone/creature"
	"github.com/olivierh59500/creatures-clone/objects"
	"github.com/olivierh59500/creatures-clone/utils"
)

// World represents the game world
type World struct {
	width, height int

	// Entities
	creatures []*creature.Creature
	objects   []objects.Object

	// World properties
	gravity   float64
	timeOfDay float64 // 0.0 to 1.0 (0=midnight, 0.5=noon)
	weather   WeatherType

	// Spatial partitioning for performance
	grid *SpatialGrid
}

// WeatherType represents different weather conditions
type WeatherType int

const (
	WeatherClear WeatherType = iota
	WeatherRain
	WeatherSnow
)

// NewWorld creates a new world instance
func NewWorld(width, height int) *World {
	return &World{
		width:     width,
		height:    height,
		creatures: make([]*creature.Creature, 0),
		objects:   make([]objects.Object, 0),
		gravity:   9.8,
		timeOfDay: 0.5, // Start at noon
		weather:   WeatherClear,
		grid:      NewSpatialGrid(width, height, 100), // 100x100 pixel cells
	}
}

// Update updates all entities in the world
func (w *World) Update() {
	// Update time of day (full cycle = 10 minutes)
	w.timeOfDay += 1.0 / (60.0 * 60.0 * 10.0) // 60 FPS * 60 seconds * 10 minutes
	if w.timeOfDay > 1.0 {
		w.timeOfDay -= 1.0
	}

	// Update spatial grid
	w.grid.Clear()
	for _, c := range w.creatures {
		w.grid.Add(c, c.X, c.Y)
	}
	for _, o := range w.objects {
		pos := o.GetPosition()
		w.grid.Add(o, pos.X, pos.Y)
	}

	// Update creatures
	for _, c := range w.creatures {
		// Find nearby entities for creature's sensors
		nearby := w.GetNearbyEntities(c.X, c.Y, 200) // 200 pixel vision range
		c.UpdateSensors(nearby, w)
		c.Update(w)

		// Apply gravity if creature is not on ground
		groundLevel := float64(w.height)*0.8 - 50 // 80% of world height minus creature height
		if c.Y < groundLevel {
			c.VelocityY += w.gravity * 0.016 // Assuming 60 FPS
		} else {
			c.Y = groundLevel
			c.VelocityY = 0
		}

		// Keep creatures in bounds
		c.X = utils.Clamp(c.X, 20, float64(w.width-20))
		c.Y = utils.Clamp(c.Y, 20, float64(w.height-20))
	}

	// Update objects
	for i := len(w.objects) - 1; i >= 0; i-- {
		obj := w.objects[i]
		obj.Update()

		// Remove consumed/destroyed objects
		if obj.ShouldRemove() {
			w.objects = append(w.objects[:i], w.objects[i+1:]...)
		}
	}

	// Handle creature interactions
	w.handleInteractions()

	// Handle breeding
	w.handleBreeding()

	// Remove dead creatures
	for i := len(w.creatures) - 1; i >= 0; i-- {
		if w.creatures[i].IsDead() {
			w.creatures = append(w.creatures[:i], w.creatures[i+1:]...)
		}
	}
}

// handleInteractions processes interactions between creatures and objects
func (w *World) handleInteractions() {
	for _, c := range w.creatures {
		// Check for food consumption
		for _, obj := range w.objects {
			if food, ok := obj.(*objects.Food); ok {
				pos := food.GetPosition()
				dist := utils.Distance(c.X, c.Y, pos.X, pos.Y)

				if dist < 30 && c.Brain.GetOutput()[creature.OutputEat] > 0.5 {
					nutritionValue := food.GetNutrition()
					c.Metabolism.Eat(nutritionValue)
					food.Consume()

					// Positive reinforcement for eating when hungry
					if c.Metabolism.Hunger > 50 {
						c.Brain.Reinforce(1.0)
					}
				}
			}

			// Check for toy interactions
			if toy, ok := obj.(*objects.Toy); ok {
				pos := toy.GetPosition()
				dist := utils.Distance(c.X, c.Y, pos.X, pos.Y)

				if dist < 40 && c.Brain.GetOutput()[creature.OutputPlay] > 0.5 {
					toy.Interact(c)
					c.Emotions.AdjustHappiness(10)

					// Positive reinforcement for playing
					c.Brain.Reinforce(0.5)
				}
			}
		}

		// Check for creature-to-creature interactions
		for _, other := range w.creatures {
			if c == other {
				continue
			}

			dist := utils.Distance(c.X, c.Y, other.X, other.Y)

			// Social interactions
			if dist < 50 {
				if c.Brain.GetOutput()[creature.OutputSpeak] > 0.5 {
					// Teaching/learning interactions
					if c.Language.GetVocabularySize() > other.Language.GetVocabularySize() {
						// Teach a word
						words := c.Language.GetKnownWords()
						if len(words) > 0 {
							word := words[utils.RandomInt(0, len(words))]
							// Convert objects slice to interface slice
							interfaceObjects := make([]interface{}, len(w.objects))
							for i, obj := range w.objects {
								interfaceObjects[i] = obj
							}
							other.Language.HearWord(word, c.GetNearestObject(interfaceObjects))
						}
					}
				}

				// Social bonding
				c.Emotions.AdjustHappiness(2)
				other.Emotions.AdjustHappiness(2)
			}
		}
	}
}

// handleBreeding checks for breeding conditions
func (w *World) handleBreeding() {
	// Limit population
	if len(w.creatures) >= w.GetMaxCreatures() {
		return
	}

	for i, c1 := range w.creatures {
		// Check if creature is ready to breed
		if !c1.CanBreed() {
			continue
		}

		for j := i + 1; j < len(w.creatures); j++ {
			c2 := w.creatures[j]

			if !c2.CanBreed() {
				continue
			}

			dist := utils.Distance(c1.X, c1.Y, c2.X, c2.Y)

			// Close enough and both willing to breed
			if dist < 60 && c1.Brain.GetOutput()[creature.OutputBreed] > 0.7 &&
				c2.Brain.GetOutput()[creature.OutputBreed] > 0.7 {
				// Create offspring
				baby := creature.Breed(c1, c2)
				baby.X = (c1.X + c2.X) / 2
				baby.Y = (c1.Y + c2.Y) / 2

				w.AddCreature(baby)

				// Parents can't breed again for a while
				c1.Metabolism.Energy -= 30
				c2.Metabolism.Energy -= 30

				// Only one breeding per update
				return
			}
		}
	}
}

// AddCreature adds a creature to the world
func (w *World) AddCreature(c *creature.Creature) {
	w.creatures = append(w.creatures, c)
}

// AddObject adds an object to the world
func (w *World) AddObject(obj objects.Object) {
	w.objects = append(w.objects, obj)
}

// GetCreatures returns all creatures in the world
func (w *World) GetCreatures() []*creature.Creature {
	return w.creatures
}

// GetObjects returns all objects in the world
func (w *World) GetObjects() []objects.Object {
	return w.objects
}

// GetNearbyEntities returns all entities within a radius of the given position
func (w *World) GetNearbyEntities(x, y, radius float64) []interface{} {
	return w.grid.GetNearby(x, y, radius)
}

// GetGravity returns the world's gravity
func (w *World) GetGravity() float64 {
	return w.gravity
}

// GetTimeOfDay returns the current time of day (0-1)
func (w *World) GetTimeOfDay() float64 {
	return w.timeOfDay
}

// GetWeather returns the current weather
func (w *World) GetWeather() WeatherType {
	return w.weather
}

// GetWidth returns the world width
func (w *World) GetWidth() int {
	return w.width
}

// GetHeight returns the world height
func (w *World) GetHeight() int {
	return w.height
}

// GetMaxCreatures returns the maximum number of creatures allowed
func (w *World) GetMaxCreatures() int {
	return 20 // Could be made configurable
}

// SpatialGrid provides efficient spatial queries
type SpatialGrid struct {
	width, height int
	cellSize      int
	cells         map[int][]interface{}
}

// NewSpatialGrid creates a new spatial grid
func NewSpatialGrid(width, height, cellSize int) *SpatialGrid {
	return &SpatialGrid{
		width:    width,
		height:   height,
		cellSize: cellSize,
		cells:    make(map[int][]interface{}),
	}
}

// Clear removes all entities from the grid
func (g *SpatialGrid) Clear() {
	g.cells = make(map[int][]interface{})
}

// Add adds an entity to the grid
func (g *SpatialGrid) Add(entity interface{}, x, y float64) {
	cellX := int(x) / g.cellSize
	cellY := int(y) / g.cellSize
	key := cellY*1000 + cellX // Simple hash

	g.cells[key] = append(g.cells[key], entity)
}

// GetNearby returns all entities within radius of the position
func (g *SpatialGrid) GetNearby(x, y, radius float64) []interface{} {
	result := make([]interface{}, 0)

	// Check cells that could contain entities within radius
	minCellX := int(x-radius) / g.cellSize
	maxCellX := int(x+radius) / g.cellSize
	minCellY := int(y-radius) / g.cellSize
	maxCellY := int(y+radius) / g.cellSize

	for cy := minCellY; cy <= maxCellY; cy++ {
		for cx := minCellX; cx <= maxCellX; cx++ {
			key := cy*1000 + cx
			if entities, ok := g.cells[key]; ok {
				result = append(result, entities...)
			}
		}
	}

	return result
}
