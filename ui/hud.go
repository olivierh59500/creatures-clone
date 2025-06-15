package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/olivierh59500/creatures-clone/creature"
)

// HUD represents the heads-up display
type HUD struct {
	// Display settings
	visible bool

	// Colors
	bgColor     color.RGBA
	barBgColor  color.RGBA
	healthColor color.RGBA
	hungerColor color.RGBA
	energyColor color.RGBA
	textColor   color.RGBA

	// Layout
	padding      float32
	barWidth     float32
	barHeight    float32
	cornerRadius float32
}

// NewHUD creates a new HUD instance
func NewHUD() *HUD {
	return &HUD{
		visible:      true,
		bgColor:      color.RGBA{0, 0, 0, 180},
		barBgColor:   color.RGBA{50, 50, 50, 255},
		healthColor:  color.RGBA{0, 255, 0, 255},
		hungerColor:  color.RGBA{255, 165, 0, 255},
		energyColor:  color.RGBA{100, 100, 255, 255},
		textColor:    color.RGBA{255, 255, 255, 255},
		padding:      10,
		barWidth:     200,
		barHeight:    20,
		cornerRadius: 5,
	}
}

// Update updates the HUD state
func (h *HUD) Update(selectedCreature *creature.Creature, world interface{}) {
	// HUD doesn't need much updating
	// Could add animations here
}

// Draw renders the HUD
func (h *HUD) Draw(screen *ebiten.Image) {
	if !h.visible {
		return
	}

	// Draw basic world info at top
	h.drawWorldInfo(screen)

	// Draw help instructions
	h.drawHelpInstructions(screen)
}

// drawHelpInstructions shows basic controls
func (h *HUD) drawHelpInstructions(screen *ebiten.Image) {
	// Background panel for instructions
	panelX := float32(10)
	panelY := float32(50)
	panelWidth := float32(350)
	panelHeight := float32(200)

	// Semi-transparent background
	bgColor := color.RGBA{0, 0, 0, 160}
	vector.DrawFilledRect(screen, panelX, panelY, panelWidth, panelHeight, bgColor, false)

	// Title
	ebitenutil.DebugPrintAt(screen, "=== HOW TO PLAY ===", int(panelX+10), int(panelY+5))

	// Instructions
	instructions := []string{
		"Left Click: Select creature / Select object",
		"Right Click: Place food / Guide creature",
		"Type + Enter: Teach word to selected creature",
		"B: Encourage breeding (when adult selected)",
		"WASD/Arrows: Move camera",
		"Mouse Wheel: Zoom in/out",
		"Space: Pause/Resume",
		"Tab: Toggle debug info",
		"1-5: Place different food types",
		"",
		"Guide creatures to objects to interact!",
		"Teach them words to build vocabulary!",
		"Keep them fed, happy, and social!",
	}

	y := int(panelY + 25)
	for _, instruction := range instructions {
		ebitenutil.DebugPrintAt(screen, instruction, int(panelX+10), y)
		y += 12
	}
}

// DrawCreatureInfo renders information about a selected creature
func (h *HUD) DrawCreatureInfo(screen *ebiten.Image, c *creature.Creature) {
	if c == nil || !h.visible {
		return
	}

	// Position at bottom left
	x := h.padding
	y := float32(screen.Bounds().Dy()) - 150
	width := h.barWidth + h.padding*2
	height := float32(130)

	// Draw background panel
	h.drawPanel(screen, x, y, width, height)

	// Draw creature name and age
	textX := x + h.padding
	textY := y + h.padding

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s", c.Name), int(textX), int(textY))

	ageText := h.getAgeText(c.Age)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Age: %s", ageText), int(textX), int(textY+15))

	// Draw status bars
	barY := textY + 35

	// Health bar
	h.drawStatusBar(screen, textX, barY, "Health", c.Metabolism.Health, h.healthColor)

	// Hunger bar
	barY += 25
	h.drawStatusBar(screen, textX, barY, "Hunger", c.Metabolism.Hunger, h.hungerColor)

	// Energy bar
	barY += 25
	h.drawStatusBar(screen, textX, barY, "Energy", c.Metabolism.Energy, h.energyColor)

	// Draw emotion state
	emotion := c.Emotions.GetDominantEmotion()
	mood := c.Emotions.GetMood()
	moodText := h.getMoodText(mood)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Feeling: %s (%s)", emotion, moodText),
		int(textX), int(barY+25))
}

// drawWorldInfo renders general world information
func (h *HUD) drawWorldInfo(screen *ebiten.Image) {
	// Time of day indicator could go here
	// For now, just show FPS
	fps := fmt.Sprintf("FPS: %0.1f", ebiten.ActualFPS())
	ebitenutil.DebugPrintAt(screen, fps, screen.Bounds().Dx()-80, 10)
}

// drawPanel draws a rounded rectangle panel
func (h *HUD) drawPanel(screen *ebiten.Image, x, y, width, height float32) {
	// Simple filled rectangle for now
	// In a full implementation, would draw rounded corners
	vector.DrawFilledRect(screen, x, y, width, height, h.bgColor, false)
}

// drawStatusBar draws a labeled progress bar
func (h *HUD) drawStatusBar(screen *ebiten.Image, x, y float32, label string, value float64, barColor color.RGBA) {
	// Draw label
	ebitenutil.DebugPrintAt(screen, label, int(x), int(y))

	// Draw background bar
	barX := x + 60
	vector.DrawFilledRect(screen, barX, y+2, h.barWidth, h.barHeight, h.barBgColor, false)

	// Draw filled portion
	fillWidth := float32(value/100) * h.barWidth
	if fillWidth > 0 {
		// Adjust color based on value
		adjustedColor := h.adjustColorByValue(barColor, value)
		vector.DrawFilledRect(screen, barX, y+2, fillWidth, h.barHeight, adjustedColor, false)
	}

	// Draw value text
	valueText := fmt.Sprintf("%0.0f%%", value)
	ebitenutil.DebugPrintAt(screen, valueText, int(barX+h.barWidth+5), int(y))
}

// adjustColorByValue modifies color based on bar value
func (h *HUD) adjustColorByValue(baseColor color.RGBA, value float64) color.RGBA {
	if value < 30 {
		// Low values tend towards red
		return color.RGBA{255, uint8(value * 3), 0, 255}
	} else if value > 70 && baseColor == h.hungerColor {
		// High hunger is bad - tend towards red
		return color.RGBA{255, 0, 0, 255}
	}
	return baseColor
}

// getAgeText formats age into readable text
func (h *HUD) getAgeText(age float64) string {
	if age < 5 {
		return "Baby"
	} else if age < 15 {
		return fmt.Sprintf("Child (%0.0f min)", age)
	} else if age < 45 {
		return fmt.Sprintf("Adult (%0.0f min)", age)
	} else {
		return fmt.Sprintf("Elder (%0.0f min)", age)
	}
}

// getMoodText converts mood value to text
func (h *HUD) getMoodText(mood float64) string {
	switch {
	case mood > 0.5:
		return "😊"
	case mood > 0:
		return "🙂"
	case mood > -0.5:
		return "😐"
	default:
		return "😢"
	}
}

// Toggle toggles HUD visibility
func (h *HUD) Toggle() {
	h.visible = !h.visible
}

// IsVisible returns whether the HUD is visible
func (h *HUD) IsVisible() bool {
	return h.visible
}

// SetVisible sets HUD visibility
func (h *HUD) SetVisible(visible bool) {
	h.visible = visible
}
