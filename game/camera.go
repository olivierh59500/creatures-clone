package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/creatures-clone/utils"
)

// Camera represents the game camera
type Camera struct {
	x, y          float64 // Position in world coordinates
	zoom          float64 // Zoom level (1.0 = normal)
	width, height int     // Screen dimensions

	// Smooth movement
	targetX, targetY float64
	smoothing        float64
}

// NewCamera creates a new camera
func NewCamera(screenWidth, screenHeight int) *Camera {
	return &Camera{
		x:         0,
		y:         0,
		zoom:      1.0,
		width:     screenWidth,
		height:    screenHeight,
		smoothing: 0.1,
	}
}

// Update updates the camera position with smoothing
func (c *Camera) Update() {
	// Smooth movement towards target
	c.x += (c.targetX - c.x) * c.smoothing
	c.y += (c.targetY - c.y) * c.smoothing
}

// Move moves the camera by the given offset
func (c *Camera) Move(dx, dy float64) {
	c.targetX += dx
	c.targetY += dy
}

// SetPosition sets the camera position directly
func (c *Camera) SetPosition(x, y float64) {
	c.x = x
	c.y = y
	c.targetX = x
	c.targetY = y
}

// Zoom adjusts the camera zoom level
func (c *Camera) Zoom(factor float64) {
	c.zoom = utils.Clamp(c.zoom*factor, 0.5, 2.0)
}

// SetZoom sets the zoom level directly
func (c *Camera) SetZoom(zoom float64) {
	c.zoom = utils.Clamp(zoom, 0.5, 2.0)
}

// FollowTarget makes the camera follow a target position
func (c *Camera) FollowTarget(targetX, targetY float64) {
	// Center the target on screen
	c.targetX = targetX - float64(c.width)/(2*c.zoom)
	c.targetY = targetY - float64(c.height)/(2*c.zoom)
}

// WorldToScreen converts world coordinates to screen coordinates
func (c *Camera) WorldToScreen(worldX, worldY float64) (float64, float64) {
	screenX := (worldX - c.x) * c.zoom
	screenY := (worldY - c.y) * c.zoom
	return screenX, screenY
}

// ScreenToWorld converts screen coordinates to world coordinates
func (c *Camera) ScreenToWorld(screenX, screenY float64) (float64, float64) {
	worldX := screenX/c.zoom + c.x
	worldY := screenY/c.zoom + c.y
	return worldX, worldY
}

// GetTransform returns the camera transformation matrix
func (c *Camera) GetTransform() *ebiten.GeoM {
	m := &ebiten.GeoM{}

	// Apply camera transformations
	m.Translate(-c.x, -c.y)
	m.Scale(c.zoom, c.zoom)

	return m
}

// GetBounds returns the visible world bounds
func (c *Camera) GetBounds() (minX, minY, maxX, maxY float64) {
	minX = c.x
	minY = c.y
	maxX = c.x + float64(c.width)/c.zoom
	maxY = c.y + float64(c.height)/c.zoom
	return
}

// IsVisible checks if a point is visible on screen
func (c *Camera) IsVisible(x, y, margin float64) bool {
	minX, minY, maxX, maxY := c.GetBounds()
	return x >= minX-margin && x <= maxX+margin && y >= minY-margin && y <= maxY+margin
}

// GetPosition returns the camera position
func (c *Camera) GetPosition() (float64, float64) {
	return c.x, c.y
}

// GetZoom returns the current zoom level
func (c *Camera) GetZoom() float64 {
	return c.zoom
}

// ConstrainToBounds keeps the camera within world bounds
func (c *Camera) ConstrainToBounds(worldWidth, worldHeight int) {
	maxX := float64(worldWidth) - float64(c.width)/c.zoom
	maxY := float64(worldHeight) - float64(c.height)/c.zoom

	c.x = utils.Clamp(c.x, 0, maxX)
	c.y = utils.Clamp(c.y, 0, maxY)
	c.targetX = utils.Clamp(c.targetX, 0, maxX)
	c.targetY = utils.Clamp(c.targetY, 0, maxY)
}
