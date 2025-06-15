package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Debug represents the debug overlay
type Debug struct {
	enabled bool

	// Debug info
	fps           float64
	creatureCount int
	objectCount   int
	worldTime     float64
	cameraPos     struct{ X, Y float64 }
	mouseWorldPos struct{ X, Y float64 }

	// Visual settings
	bgColor   color.RGBA
	textColor color.RGBA
}

// NewDebug creates a new debug overlay
func NewDebug() *Debug {
	return &Debug{
		enabled:   false,
		bgColor:   color.RGBA{0, 0, 0, 180},
		textColor: color.RGBA{0, 255, 0, 255},
	}
}

// Update updates debug information
func (d *Debug) Update(world, camera interface{}, mouseX, mouseY int) {
	if !d.enabled {
		return
	}

	// Update FPS
	d.fps = ebiten.ActualFPS()

	// Would extract actual values from world and camera
	// For now, using placeholder values
	d.creatureCount = 3
	d.objectCount = 15
	d.worldTime = 0
	d.cameraPos.X = 0
	d.cameraPos.Y = 0
	d.mouseWorldPos.X = float64(mouseX)
	d.mouseWorldPos.Y = float64(mouseY)
}

// Draw renders the debug overlay
func (d *Debug) Draw(screen *ebiten.Image) {
	if !d.enabled {
		return
	}

	// Draw background panel
	panelWidth := float32(250)
	panelHeight := float32(200)
	vector.DrawFilledRect(screen, 10, 40, panelWidth, panelHeight, d.bgColor, false)

	// Draw debug information
	x := 15
	y := 45
	lineHeight := 15

	debugInfo := []string{
		fmt.Sprintf("=== DEBUG INFO ==="),
		fmt.Sprintf("FPS: %.1f", d.fps),
		fmt.Sprintf("Creatures: %d", d.creatureCount),
		fmt.Sprintf("Objects: %d", d.objectCount),
		fmt.Sprintf("Time: %.1f", d.worldTime),
		fmt.Sprintf("Camera: (%.0f, %.0f)", d.cameraPos.X, d.cameraPos.Y),
		fmt.Sprintf("Mouse: (%.0f, %.0f)", d.mouseWorldPos.X, d.mouseWorldPos.Y),
		"",
		"Controls:",
		"Tab - Toggle Debug",
		"WASD - Move Camera",
		"Space - Pause",
		"ESC - Menu",
	}

	for i, line := range debugInfo {
		ebitenutil.DebugPrintAt(screen, line, x, y+i*lineHeight)
	}
}

// Toggle toggles the debug overlay
func (d *Debug) Toggle() {
	d.enabled = !d.enabled
}

// IsEnabled returns whether debug mode is enabled
func (d *Debug) IsEnabled() bool {
	return d.enabled
}

// SetEnabled sets debug mode state
func (d *Debug) SetEnabled(enabled bool) {
	d.enabled = enabled
}
