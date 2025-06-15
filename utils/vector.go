package utils

import (
	"math"
)

// Vector2D represents a 2D vector
type Vector2D struct {
	X, Y float64
}

// NewVector2D creates a new vector
func NewVector2D(x, y float64) Vector2D {
	return Vector2D{X: x, Y: y}
}

// Add adds two vectors
func (v Vector2D) Add(other Vector2D) Vector2D {
	return Vector2D{
		X: v.X + other.X,
		Y: v.Y + other.Y,
	}
}

// Sub subtracts two vectors
func (v Vector2D) Sub(other Vector2D) Vector2D {
	return Vector2D{
		X: v.X - other.X,
		Y: v.Y - other.Y,
	}
}

// Mul multiplies vector by scalar
func (v Vector2D) Mul(scalar float64) Vector2D {
	return Vector2D{
		X: v.X * scalar,
		Y: v.Y * scalar,
	}
}

// Div divides vector by scalar
func (v Vector2D) Div(scalar float64) Vector2D {
	if scalar == 0 {
		return v
	}
	return Vector2D{
		X: v.X / scalar,
		Y: v.Y / scalar,
	}
}

// Magnitude returns the length of the vector
func (v Vector2D) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Normalize returns a unit vector in the same direction
func (v Vector2D) Normalize() Vector2D {
	mag := v.Magnitude()
	if mag == 0 {
		return v
	}
	return v.Div(mag)
}

// Dot returns the dot product of two vectors
func (v Vector2D) Dot(other Vector2D) float64 {
	return v.X*other.X + v.Y*other.Y
}

// Distance returns the distance between two points
func (v Vector2D) Distance(other Vector2D) float64 {
	dx := v.X - other.X
	dy := v.Y - other.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// Angle returns the angle of the vector in radians
func (v Vector2D) Angle() float64 {
	return math.Atan2(v.Y, v.X)
}

// Rotate rotates the vector by the given angle in radians
func (v Vector2D) Rotate(angle float64) Vector2D {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	return Vector2D{
		X: v.X*cos - v.Y*sin,
		Y: v.X*sin + v.Y*cos,
	}
}

// Lerp performs linear interpolation between two vectors
func (v Vector2D) Lerp(other Vector2D, t float64) Vector2D {
	t = Clamp(t, 0, 1)
	return Vector2D{
		X: v.X + (other.X-v.X)*t,
		Y: v.Y + (other.Y-v.Y)*t,
	}
}

// Clamp clamps each component of the vector
func (v Vector2D) Clamp(min, max Vector2D) Vector2D {
	return Vector2D{
		X: Clamp(v.X, min.X, max.X),
		Y: Clamp(v.Y, min.Y, max.Y),
	}
}

// Zero returns a zero vector
func Zero() Vector2D {
	return Vector2D{X: 0, Y: 0}
}

// One returns a vector with both components set to 1
func One() Vector2D {
	return Vector2D{X: 1, Y: 1}
}

// Up returns an up vector (0, -1)
func Up() Vector2D {
	return Vector2D{X: 0, Y: -1}
}

// Down returns a down vector (0, 1)
func Down() Vector2D {
	return Vector2D{X: 0, Y: 1}
}

// Left returns a left vector (-1, 0)
func Left() Vector2D {
	return Vector2D{X: -1, Y: 0}
}

// Right returns a right vector (1, 0)
func Right() Vector2D {
	return Vector2D{X: 1, Y: 0}
}
