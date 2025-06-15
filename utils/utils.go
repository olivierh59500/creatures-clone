package utils

import (
	"crypto/rand"
	"encoding/hex"
	"math"
)

// Color represents an RGBA color
type Color struct {
	R, G, B, A uint8
}

// Math functions

// Clamp constrains a value between min and max
func Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// ClampInt constrains an integer between min and max
func ClampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Min returns the minimum of two values
func Min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// Max returns the maximum of two values
func Max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// Abs returns the absolute value
func Abs(x float64) float64 {
	return math.Abs(x)
}

// Sign returns -1 for negative, 1 for positive, 0 for zero
func Sign(x float64) float64 {
	if x < 0 {
		return -1
	} else if x > 0 {
		return 1
	}
	return 0
}

// Lerp performs linear interpolation
func Lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

// Distance calculates distance between two points
func Distance(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(dx*dx + dy*dy)
}

// AngleBetween calculates angle between two points
func AngleBetween(x1, y1, x2, y2 float64) float64 {
	return math.Atan2(y2-y1, x2-x1)
}

// NormalizeAngle normalizes an angle to [-π, π]
func NormalizeAngle(angle float64) float64 {
	for angle > math.Pi {
		angle -= 2 * math.Pi
	}
	for angle < -math.Pi {
		angle += 2 * math.Pi
	}
	return angle
}

// Trigonometric functions

// Sin returns sine of angle
func Sin(angle float64) float64 {
	return math.Sin(angle)
}

// Cos returns cosine of angle
func Cos(angle float64) float64 {
	return math.Cos(angle)
}

// Sqrt returns square root
func Sqrt(x float64) float64 {
	return math.Sqrt(x)
}

// ID generation

// GenerateID generates a unique identifier
func GenerateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// Color utilities

// NewColor creates a new color
func NewColor(r, g, b, a uint8) Color {
	return Color{R: r, G: g, B: b, A: a}
}

// Lerp interpolates between two colors
func (c Color) Lerp(other Color, t float64) Color {
	t = Clamp(t, 0, 1)
	return Color{
		R: uint8(float64(c.R) + (float64(other.R)-float64(c.R))*t),
		G: uint8(float64(c.G) + (float64(other.G)-float64(c.G))*t),
		B: uint8(float64(c.B) + (float64(other.B)-float64(c.B))*t),
		A: uint8(float64(c.A) + (float64(other.A)-float64(c.A))*t),
	}
}

// WithAlpha returns color with new alpha value
func (c Color) WithAlpha(alpha uint8) Color {
	return Color{R: c.R, G: c.G, B: c.B, A: alpha}
}

// Common colors
var (
	ColorWhite   = Color{255, 255, 255, 255}
	ColorBlack   = Color{0, 0, 0, 255}
	ColorRed     = Color{255, 0, 0, 255}
	ColorGreen   = Color{0, 255, 0, 255}
	ColorBlue    = Color{0, 0, 255, 255}
	ColorYellow  = Color{255, 255, 0, 255}
	ColorCyan    = Color{0, 255, 255, 255}
	ColorMagenta = Color{255, 0, 255, 255}
	ColorGray    = Color{128, 128, 128, 255}
)

// Time utilities

// FormatTime formats seconds into MM:SS format
func FormatTime(seconds float64) string {
	minutes := int(seconds) / 60
	secs := int(seconds) % 60
	return string(rune(minutes/10+'0')) + string(rune(minutes%10+'0')) + ":" +
		string(rune(secs/10+'0')) + string(rune(secs%10+'0'))
}

// Map remaps a value from one range to another
func Map(value, inMin, inMax, outMin, outMax float64) float64 {
	return outMin + (value-inMin)*(outMax-outMin)/(inMax-inMin)
}

// SmoothStep performs smooth Hermite interpolation
func SmoothStep(edge0, edge1, x float64) float64 {
	t := Clamp((x-edge0)/(edge1-edge0), 0.0, 1.0)
	return t * t * (3.0 - 2.0*t)
}
