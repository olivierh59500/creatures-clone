package renderer

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/olivierh59500/creatures-clone/creature"
	"github.com/olivierh59500/creatures-clone/objects"
	"github.com/olivierh59500/creatures-clone/utils"
)

// Renderer handles all game rendering
type Renderer struct {
	// Asset manager
	assets *AssetManager

	// Sprite system
	sprites map[string]*Sprite

	// Animation system
	animations map[string]*Animation

	// Particle system
	particles []Particle

	// Render settings
	enableShadows   bool
	enableParticles bool
}

// NewRenderer creates a new renderer
func NewRenderer() *Renderer {
	r := &Renderer{
		assets:          NewAssetManager(),
		sprites:         make(map[string]*Sprite),
		animations:      make(map[string]*Animation),
		particles:       make([]Particle, 0),
		enableShadows:   true,
		enableParticles: true,
	}

	// Initialize built-in sprites
	r.initializeSprites()

	return r
}

// initializeSprites creates programmatic sprites
func (r *Renderer) initializeSprites() {
	// Create basic shapes as sprites
	// In a full implementation, these would be loaded from files
	// or generated programmatically
}

// WorldInfo interface to avoid circular import
type WorldInfo interface {
	GetWidth() int
	GetHeight() int
}

// DrawWorldBackground draws the world background
func (r *Renderer) DrawWorldBackground(screen *ebiten.Image, world WorldInfo, transform *ebiten.GeoM) {
	bounds := screen.Bounds()

	// Get world dimensions
	worldWidth := float64(world.GetWidth())
	worldHeight := float64(world.GetHeight())

	// Calculate ground level in world coordinates (80% of world height)
	worldGroundY := worldHeight * 0.8

	// Draw sky - this fills the entire screen regardless of zoom
	for y := 0; y < bounds.Dy(); y++ {
		t := float64(y) / float64(bounds.Dy())

		skyColor := lerpColor(
			color.RGBA{255, 200, 150, 255}, // Light peach at horizon
			color.RGBA{135, 206, 235, 255}, // Sky blue at top
			t*t,                            // Non-linear gradient
		)

		vector.DrawFilledRect(screen, 0, float32(y), float32(bounds.Dx()), 1, skyColor, false)
	}

	// Draw clouds in world space
	r.drawCloudInWorld(screen, transform, worldWidth*0.2, worldHeight*0.2, 80)
	r.drawCloudInWorld(screen, transform, worldWidth*0.6, worldHeight*0.15, 100)
	r.drawCloudInWorld(screen, transform, worldWidth*0.8, worldHeight*0.25, 60)

	// Draw sun in world space
	sunX := worldWidth * 0.85
	sunY := worldHeight * 0.15
	r.drawSunInWorld(screen, transform, sunX, sunY, 40)

	// Draw ground in world coordinates
	// Create a temporary image for the ground that spans the entire world width
	groundImg := ebiten.NewImage(int(worldWidth), int(worldHeight-worldGroundY))

	// Fill with ground color
	groundColor := color.RGBA{139, 115, 85, 255}
	groundImg.Fill(groundColor)

	// Add grass texture
	grassHeight := 30
	for x := 0; x < int(worldWidth); x += 2 {
		variation := math.Sin(float64(x)*0.02) * 5
		grassColor := color.RGBA{
			R: 34 + uint8(variation),
			G: 139 + uint8(variation*2),
			B: 34,
			A: 255,
		}

		vector.DrawFilledRect(groundImg, float32(x), 0, 2, float32(grassHeight)+float32(variation), grassColor, false)
	}

	// Draw the ground with the camera transform
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, worldGroundY)
	op.GeoM.Concat(*transform)
	screen.DrawImage(groundImg, op)
}

// drawCloud draws a fluffy cloud
func (r *Renderer) drawCloud(screen *ebiten.Image, x, y, size float32) {
	cloudColor := color.RGBA{255, 255, 255, 200}

	// Create cloud with multiple circles
	positions := []struct{ x, y, r float32 }{
		{0, 0, size * 0.5},
		{size * 0.3, -size * 0.1, size * 0.4},
		{size * 0.6, 0, size * 0.45},
		{size * 0.9, -size * 0.05, size * 0.35},
		{size * 0.2, size * 0.15, size * 0.3},
		{size * 0.7, size * 0.1, size * 0.3},
	}

	for _, pos := range positions {
		vector.DrawFilledCircle(screen, x+pos.x, y+pos.y, pos.r, cloudColor, false)
	}
}

// drawCloudInWorld draws a cloud in world coordinates
func (r *Renderer) drawCloudInWorld(screen *ebiten.Image, transform *ebiten.GeoM, x, y, size float64) {
	// Transform world coordinates to screen coordinates
	screenX, screenY := transform.Apply(x, y)

	// Get approximate scale from transform
	// Since we only use uniform scaling, we can check the (0,0) element
	scale := transform.Element(0, 0)
	scaledSize := float32(size * scale)

	r.drawCloud(screen, float32(screenX), float32(screenY), scaledSize)
}

// drawSunInWorld draws the sun in world coordinates
func (r *Renderer) drawSunInWorld(screen *ebiten.Image, transform *ebiten.GeoM, x, y, radius float64) {
	// Transform world coordinates to screen coordinates
	screenX, screenY := transform.Apply(x, y)

	// Get scale from transform (assuming uniform scaling)
	scale := transform.Element(0, 0)
	scaledRadius := float32(radius * scale)

	r.drawSun(screen, float32(screenX), float32(screenY), scaledRadius)
}

// drawSun draws the sun with rays
func (r *Renderer) drawSun(screen *ebiten.Image, x, y, radius float32) {
	// Sun body
	sunColor := color.RGBA{255, 255, 100, 255}
	vector.DrawFilledCircle(screen, x, y, radius, sunColor, false)

	// Sun rays
	rayColor := color.RGBA{255, 255, 150, 100}
	numRays := 12
	for i := 0; i < numRays; i++ {
		angle := float32(i) * 2 * math.Pi / float32(numRays)
		x1 := x + radius*1.2*float32(math.Cos(float64(angle)))
		y1 := y + radius*1.2*float32(math.Sin(float64(angle)))
		x2 := x + radius*1.8*float32(math.Cos(float64(angle)))
		y2 := y + radius*1.8*float32(math.Sin(float64(angle)))

		vector.StrokeLine(screen, x1, y1, x2, y2, 3, rayColor, false)
	}
}

// DrawCreature renders a creature
func (r *Renderer) DrawCreature(screen *ebiten.Image, c *creature.Creature, transform *ebiten.GeoM, isSelected bool) {
	// Get screen position
	screenX, screenY := transform.Apply(c.X, c.Y)

	// Draw shadow if enabled
	if r.enableShadows {
		r.drawShadow(screen, screenX, screenY, 20*c.Size)
	}

	// Draw creature body
	r.drawCreatureBody(screen, c, screenX, screenY)

	// Draw selection indicator
	if isSelected {
		r.drawSelectionIndicator(screen, screenX, screenY, 30*c.Size)
	}

	// Draw speech bubble if speaking
	if c.Language.IsSpeaking() {
		r.drawSpeechBubble(screen, screenX, screenY-40, c.Language.GetCurrentWord())
	}

	// Draw emotion indicator
	r.drawEmotionIndicator(screen, c, screenX, screenY)
}

// drawCreatureBody draws the creature's body parts
func (r *Renderer) drawCreatureBody(screen *ebiten.Image, c *creature.Creature, x, y float64) {
	// Get creature color from genetics
	creatureColor := color.RGBA{
		R: c.Color.R,
		G: c.Color.G,
		B: c.Color.B,
		A: c.Color.A,
	}

	// Body (oval)
	bodyWidth := float32(40 * c.Size)
	bodyHeight := float32(50 * c.Size)
	r.drawOval(screen, float32(x), float32(y), bodyWidth, bodyHeight, creatureColor)

	// Head (circle)
	headSize := float32(30 * c.Size)
	headY := float32(y) - bodyHeight/2 - headSize/2
	r.drawCircle(screen, float32(x), headY, headSize/2, creatureColor)

	// Eyes
	eyeSize := float32(8 * c.Size)
	eyeY := headY - 5
	leftEyeX := float32(x) - 8*float32(c.Size)
	rightEyeX := float32(x) + 8*float32(c.Size)

	// Eye whites
	r.drawCircle(screen, leftEyeX, eyeY, eyeSize/2, color.White)
	r.drawCircle(screen, rightEyeX, eyeY, eyeSize/2, color.White)

	// Pupils (look in direction of movement)
	pupilOffset := float32(2)
	if c.VelocityX > 0 {
		pupilOffset = 2
	} else if c.VelocityX < 0 {
		pupilOffset = -2
	}

	pupilSize := float32(4 * c.Size)
	r.drawCircle(screen, leftEyeX+pupilOffset, eyeY, pupilSize/2, color.Black)
	r.drawCircle(screen, rightEyeX+pupilOffset, eyeY, pupilSize/2, color.Black)

	// Arms
	armWidth := float32(15 * c.Size)
	armHeight := float32(8 * c.Size)
	armY := float32(y) - 10
	r.drawOval(screen, float32(x)-bodyWidth/2-armWidth/2, armY, armWidth, armHeight, creatureColor)
	r.drawOval(screen, float32(x)+bodyWidth/2+armWidth/2, armY, armWidth, armHeight, creatureColor)

	// Legs with walking animation
	legWidth := float32(10 * c.Size)
	legHeight := float32(15 * c.Size)
	legY := float32(y) + bodyHeight/2

	// Get leg positions from movement system
	leftLegX, leftLegY := c.Movement.GetLegPosition(true)
	rightLegX, rightLegY := c.Movement.GetLegPosition(false)

	r.drawOval(screen, float32(x)-5+float32(leftLegX), legY+float32(leftLegY), legWidth, legHeight, creatureColor)
	r.drawOval(screen, float32(x)+5+float32(rightLegX), legY+float32(rightLegY), legWidth, legHeight, creatureColor)

	// Expression based on emotions
	if c.Emotions.Happiness > 50 {
		// Smile
		r.drawArc(screen, float32(x), headY+5, 10, math.Pi*0.2, math.Pi*0.8, color.Black)
	} else if c.Emotions.Fear > 50 {
		// Worried expression
		r.drawLine(screen, float32(x)-5, headY+5, float32(x)+5, headY+3, color.Black)
	}
}

// DrawObject renders a game object
func (r *Renderer) DrawObject(screen *ebiten.Image, obj objects.Object, transform *ebiten.GeoM) {
	pos := obj.GetPosition()
	screenX, screenY := transform.Apply(pos.X, pos.Y)

	// Draw shadow if enabled
	if r.enableShadows {
		r.drawShadow(screen, screenX, screenY, 15*obj.GetSize())
	}

	// Draw based on object type
	switch obj.GetType() {
	case "food":
		r.drawFood(screen, obj.(*objects.Food), screenX, screenY)
	case "toy":
		r.drawToy(screen, obj.(*objects.Toy), screenX, screenY)
	case "plant":
		r.drawPlant(screen, obj.(*objects.Plant), screenX, screenY)
	default:
		// Generic object rendering
		r.drawGenericObject(screen, obj, screenX, screenY)
	}
}

// drawFood renders food items
func (r *Renderer) drawFood(screen *ebiten.Image, food *objects.Food, x, y float64) {
	foodColor := color.RGBA{
		R: food.Color.R,
		G: food.Color.G,
		B: food.Color.B,
		A: food.Color.A,
	}

	// Apply bounce animation
	bounceY := food.GetBounceY()

	// Adjust y so food sits on ground
	y = y - bounceY

	switch food.GetSprite() {
	case "apple":
		// Draw apple shape sitting on ground
		radius := float32(20 * food.Size)
		r.drawCircle(screen, float32(x), float32(y)-radius, radius, foodColor)
		// Stem
		r.drawRect(screen, float32(x)-2, float32(y)-radius*2-5, 4, 8, color.RGBA{0, 128, 0, 255})

	case "carrot":
		// Draw carrot shape with point in ground
		r.drawTriangle(screen, float32(x), float32(y)-12, 10, 25, foodColor)
		// Green top
		r.drawCircle(screen, float32(x), float32(y)-32, 8, color.RGBA{0, 255, 0, 255})

	case "honey":
		// Draw hexagon sitting on ground
		size := float32(15 * food.Size)
		r.drawHexagon(screen, float32(x), float32(y)-size, size, foodColor)

	case "berry":
		// Draw berry cluster on ground
		offsets := []struct{ x, y float32 }{
			{0, -5}, {-5, -8}, {5, -8}, {0, -10},
		}
		for _, offset := range offsets {
			r.drawCircle(screen, float32(x)+offset.x, float32(y)+offset.y, 4, foodColor)
		}

	default:
		// Generic food on ground
		radius := float32(15 * food.Size)
		r.drawCircle(screen, float32(x), float32(y)-radius, radius, foodColor)
	}
}

// drawToy renders toy objects
func (r *Renderer) drawToy(screen *ebiten.Image, toy *objects.Toy, x, y float64) {
	toyColor := color.RGBA{
		R: toy.Color.R,
		G: toy.Color.G,
		B: toy.Color.B,
		A: toy.Color.A,
	}

	// Apply animations
	bounceOffset := float64(toy.GetBounceOffset())
	y = y - bounceOffset
	rotation := toy.GetRotation()

	switch toy.GetSprite() {
	case "ball":
		// Draw ball sitting on ground
		radius := float32(25 * toy.Size)
		r.drawCircleWithRotation(screen, float32(x), float32(y)-radius, radius, toyColor, rotation)
		// Add stripes
		r.drawArc(screen, float32(x), float32(y)-radius, radius, 0, math.Pi, color.White)

	case "musicbox":
		// Draw music box on ground
		r.drawRect(screen, float32(x)-15, float32(y)-20, 30, 20, toyColor)
		// Handle
		r.drawRectRotated(screen, float32(x)+10, float32(y)-25, 5, 10, color.RGBA{255, 215, 0, 255}, rotation)
		// Musical notes if playing
		if toy.IsPlaying() {
			r.addMusicNoteParticle(float32(x), float32(y)-40)
		}

	case "computer":
		// Draw computer on ground/desk level
		r.drawRect(screen, float32(x)-20, float32(y)-30, 40, 30, toyColor)
		// Screen
		r.drawRect(screen, float32(x)-15, float32(y)-25, 30, 20, color.RGBA{0, 0, 255, 255})
		// Show blinking cursor if active
		if toy.IsPlaying() {
			cursorBlink := int(toy.AnimationTime*2) % 2
			if cursorBlink == 0 {
				r.drawRect(screen, float32(x)-10, float32(y)-20, 2, 8, color.White)
			}
		}

	case "bed":
		// Draw bed on ground
		r.drawRect(screen, float32(x)-30, float32(y)-20, 60, 20, toyColor)
		// Pillow
		r.drawOval(screen, float32(x)-20, float32(y)-25, 20, 10, color.White)
		// Show Z's if creature is sleeping on it
		if toy.IsPlaying() {
			r.addSleepParticle(float32(x), float32(y)-45)
		}

	default:
		// Generic toy on ground
		size := float32(30 * toy.Size)
		r.drawRect(screen, float32(x)-size/2, float32(y)-size, size, size, toyColor)
	}
}

// drawPlant renders plant objects
func (r *Renderer) drawPlant(screen *ebiten.Image, plant *objects.Plant, x, y float64) {
	plantColor := color.RGBA{
		R: plant.Color.R,
		G: plant.Color.G,
		B: plant.Color.B,
		A: plant.Color.A,
	}

	// Apply sway animation
	swayX := plant.GetSwayX()

	switch plant.GetSprite() {
	case "tree":
		// Draw trunk - anchor at bottom
		trunkWidth := float32(20 * plant.Size)
		trunkHeight := float32(40 * plant.Size)
		r.drawRect(screen, float32(x)-trunkWidth/2, float32(y)-trunkHeight, trunkWidth, trunkHeight, color.RGBA{139, 69, 19, 255})

		// Draw leaves above trunk
		leavesSize := float32(60 * plant.Size)
		r.drawCircle(screen, float32(x+swayX), float32(y)-trunkHeight-leavesSize/3, leavesSize/2, plantColor)

	case "flower":
		// Draw stem from ground up
		stemHeight := float32(20 * plant.Size)
		r.drawRect(screen, float32(x)-1, float32(y)-stemHeight, 2, stemHeight, color.RGBA{0, 128, 0, 255})

		// Draw petals at top of stem
		petalSize := float32(20 * plant.Size)
		r.drawCircle(screen, float32(x+swayX), float32(y)-stemHeight-petalSize/2, petalSize/2, plantColor)

	default:
		// Generic plant anchored at ground
		plantHeight := float32(30 * plant.Size)
		r.drawRect(screen, float32(x)-5, float32(y)-plantHeight, 10, plantHeight, plantColor)
	}
}

// Helper drawing functions

func (r *Renderer) drawCircle(screen *ebiten.Image, x, y, radius float32, c color.Color) {
	vector.DrawFilledCircle(screen, x, y, radius, c, false)
}

func (r *Renderer) drawOval(screen *ebiten.Image, x, y, width, height float32, c color.Color) {
	// Approximate oval with scaled circle
	// In a full implementation, would draw proper oval
	vector.DrawFilledCircle(screen, x, y, width/2, c, false)
}

func (r *Renderer) drawRect(screen *ebiten.Image, x, y, width, height float32, c color.Color) {
	vector.DrawFilledRect(screen, x, y, width, height, c, false)
}

func (r *Renderer) drawTriangle(screen *ebiten.Image, x, y, width, height float32, c color.Color) {
	// Draw triangle using lines
	// In a full implementation, would fill the triangle
	vector.DrawFilledRect(screen, x-width/2, y, width, height, c, false)
}

func (r *Renderer) drawHexagon(screen *ebiten.Image, x, y, size float32, c color.Color) {
	// Approximate with circle for now
	vector.DrawFilledCircle(screen, x, y, size, c, false)
}

func (r *Renderer) drawLine(screen *ebiten.Image, x1, y1, x2, y2 float32, c color.Color) {
	vector.StrokeLine(screen, x1, y1, x2, y2, 2, c, false)
}

func (r *Renderer) drawArc(screen *ebiten.Image, x, y, radius, startAngle, endAngle float32, c color.Color) {
	// Approximate arc with lines
	steps := 10
	angleStep := (endAngle - startAngle) / float32(steps)

	for i := 0; i < steps; i++ {
		angle1 := startAngle + float32(i)*angleStep
		angle2 := startAngle + float32(i+1)*angleStep

		x1 := x + radius*float32(math.Cos(float64(angle1)))
		y1 := y + radius*float32(math.Sin(float64(angle1)))
		x2 := x + radius*float32(math.Cos(float64(angle2)))
		y2 := y + radius*float32(math.Sin(float64(angle2)))

		vector.StrokeLine(screen, x1, y1, x2, y2, 2, c, false)
	}
}

func (r *Renderer) drawCircleWithRotation(screen *ebiten.Image, x, y, radius float32, c color.Color, rotation float64) {
	// For now, just draw a regular circle
	// In a full implementation, would apply rotation to any patterns
	vector.DrawFilledCircle(screen, x, y, radius, c, false)
}

func (r *Renderer) drawRectRotated(screen *ebiten.Image, x, y, width, height float32, c color.Color, rotation float64) {
	// Simplified - just draw regular rect
	// In a full implementation, would apply rotation transform
	vector.DrawFilledRect(screen, x-width/2, y-height/2, width, height, c, false)
}

func (r *Renderer) drawShadow(screen *ebiten.Image, x, y, size float64) {
	shadowColor := color.RGBA{0, 0, 0, 64}
	vector.DrawFilledCircle(screen, float32(x), float32(y+size/2), float32(size*0.4), shadowColor, false)
}

func (r *Renderer) drawSelectionIndicator(screen *ebiten.Image, x, y, radius float64) {
	// Draw animated selection circle
	selectionColor := color.RGBA{255, 255, 0, 200}
	vector.StrokeCircle(screen, float32(x), float32(y), float32(radius), 3, selectionColor, false)
}

func (r *Renderer) drawSpeechBubble(screen *ebiten.Image, x, y float64, text string) {
	// Simple speech bubble
	bubbleWidth := float32(len(text)*8 + 20)
	bubbleHeight := float32(30)

	// Bubble body
	vector.DrawFilledRect(screen, float32(x)-bubbleWidth/2, float32(y)-bubbleHeight/2,
		bubbleWidth, bubbleHeight, color.White, false)

	// Tail
	vector.DrawFilledRect(screen, float32(x)-5, float32(y)+bubbleHeight/2, 10, 10, color.White, false)

	// Text (would use proper font rendering in full implementation)
	// For now, we'll skip the actual text rendering
}

func (r *Renderer) drawEmotionIndicator(screen *ebiten.Image, c *creature.Creature, x, y float64) {
	emotion := c.Emotions.GetDominantEmotion()

	// Position above head
	indicatorY := y - 60*c.Size

	switch emotion {
	case "happy":
		// Happy face
		r.drawCircle(screen, float32(x), float32(indicatorY), 8, color.RGBA{255, 255, 0, 200})
	case "afraid":
		// Exclamation mark
		r.drawRect(screen, float32(x)-2, float32(indicatorY)-10, 4, 8, color.RGBA{255, 0, 0, 200})
		r.drawCircle(screen, float32(x), float32(indicatorY)+5, 2, color.RGBA{255, 0, 0, 200})
	case "angry":
		// Angry symbol
		r.drawLine(screen, float32(x)-5, float32(indicatorY)-5, float32(x)+5, float32(indicatorY)+5, color.RGBA{255, 0, 0, 200})
		r.drawLine(screen, float32(x)-5, float32(indicatorY)+5, float32(x)+5, float32(indicatorY)-5, color.RGBA{255, 0, 0, 200})
	case "curious":
		// Question mark
		r.drawArc(screen, float32(x), float32(indicatorY)-5, 5, -math.Pi/2, math.Pi/2, color.RGBA{0, 200, 0, 200})
		r.drawCircle(screen, float32(x), float32(indicatorY)+8, 2, color.RGBA{0, 200, 0, 200})
	}
}

func (r *Renderer) drawGenericObject(screen *ebiten.Image, obj objects.Object, x, y float64) {
	objColor := color.RGBA{
		R: obj.GetColor().R,
		G: obj.GetColor().G,
		B: obj.GetColor().B,
		A: obj.GetColor().A,
	}

	// Draw as a simple rectangle
	size := float32(30 * obj.GetSize())
	r.drawRect(screen, float32(x)-size/2, float32(y)-size/2, size, size, objColor)
}

func (r *Renderer) addMusicNoteParticle(x, y float32) {
	if !r.enableParticles || len(r.particles) >= 100 {
		return
	}

	// Add musical note particle
	p := Particle{
		X:     x,
		Y:     y,
		VX:    float32(utils.RandomFloat(-0.5, 0.5)),
		VY:    -1,
		Life:  60,
		Type:  ParticleNote,
		Color: color.RGBA{255, 215, 0, 255},
		Size:  5,
	}

	r.particles = append(r.particles, p)
}

func (r *Renderer) addSleepParticle(x, y float32) {
	if !r.enableParticles || len(r.particles) >= 100 {
		return
	}

	// Add Z particle for sleeping
	p := Particle{
		X:        x + float32(utils.RandomFloat(-10, 10)),
		Y:        y,
		VX:       float32(utils.RandomFloat(-0.2, 0.2)),
		VY:       -0.5,
		Life:     90,
		Type:     ParticleZ,
		Color:    color.RGBA{173, 216, 230, 200},
		Size:     8,
		Rotation: float32(utils.RandomFloat(-0.2, 0.2)),
		RotSpeed: 0,
	}

	r.particles = append(r.particles, p)
}

// UpdateParticles updates all particles
func (r *Renderer) UpdateParticles() {
	for i := len(r.particles) - 1; i >= 0; i-- {
		p := &r.particles[i]
		p.Update()

		if p.Life <= 0 {
			r.particles = append(r.particles[:i], r.particles[i+1:]...)
		}
	}
}

// DrawParticles renders all particles
func (r *Renderer) DrawParticles(screen *ebiten.Image) {
	for _, p := range r.particles {
		p.Draw(screen)
	}
}

// Helper function to interpolate colors
func lerpColor(c1, c2 color.RGBA, t float64) color.RGBA {
	return color.RGBA{
		R: uint8(float64(c1.R) + (float64(c2.R)-float64(c1.R))*t),
		G: uint8(float64(c1.G) + (float64(c2.G)-float64(c1.G))*t),
		B: uint8(float64(c1.B) + (float64(c2.B)-float64(c1.B))*t),
		A: uint8(float64(c1.A) + (float64(c2.A)-float64(c1.A))*t),
	}
}
