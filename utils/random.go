package utils

import (
	"math/rand"
	"time"
)

// Initialize random seed
func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt returns a random integer between min and max (exclusive)
func RandomInt(min, max int) int {
	if min >= max {
		return min
	}
	return min + rand.Intn(max-min)
}

// RandomFloat returns a random float64 between min and max
func RandomFloat(min, max float64) float64 {
	if min >= max {
		return min
	}
	return min + rand.Float64()*(max-min)
}

// RandomBool returns a random boolean
func RandomBool() bool {
	return rand.Float64() < 0.5
}

// RandomChoice returns a random element from a slice
func RandomChoice[T any](choices []T) T {
	if len(choices) == 0 {
		var zero T
		return zero
	}
	return choices[rand.Intn(len(choices))]
}

// RandomWeighted returns a random index based on weights
func RandomWeighted(weights []float64) int {
	if len(weights) == 0 {
		return -1
	}

	// Calculate total weight
	total := 0.0
	for _, w := range weights {
		total += w
	}

	if total == 0 {
		return rand.Intn(len(weights))
	}

	// Random value between 0 and total
	r := rand.Float64() * total

	// Find which weight range it falls into
	cumulative := 0.0
	for i, w := range weights {
		cumulative += w
		if r <= cumulative {
			return i
		}
	}

	return len(weights) - 1
}

// Chance returns true with the given probability (0-1)
func Chance(probability float64) bool {
	return rand.Float64() < probability
}

// RandomNormal returns a normally distributed random number
func RandomNormal(mean, stddev float64) float64 {
	return rand.NormFloat64()*stddev + mean
}

// RandomDirection returns a random unit vector
func RandomDirection() Vector2D {
	angle := rand.Float64() * 2 * 3.14159265359
	return Vector2D{
		X: Cos(angle),
		Y: Sin(angle),
	}
}

// RandomPointInCircle returns a random point within a circle
func RandomPointInCircle(centerX, centerY, radius float64) (float64, float64) {
	// Use square root for uniform distribution
	r := Sqrt(rand.Float64()) * radius
	theta := rand.Float64() * 2 * 3.14159265359

	x := centerX + r*Cos(theta)
	y := centerY + r*Sin(theta)

	return x, y
}

// RandomPointInRect returns a random point within a rectangle
func RandomPointInRect(x, y, width, height float64) (float64, float64) {
	px := x + rand.Float64()*width
	py := y + rand.Float64()*height
	return px, py
}

// Shuffle shuffles a slice in place
func Shuffle[T any](slice []T) {
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}
