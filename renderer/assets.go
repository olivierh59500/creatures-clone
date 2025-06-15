package renderer

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// AssetManager manages all game assets
type AssetManager struct {
	// Creature assets
	creatureSprites map[string]*ebiten.Image

	// Object assets
	foodSprites  map[string]*ebiten.Image
	toySprites   map[string]*ebiten.Image
	plantSprites map[string]*ebiten.Image

	// UI assets
	uiSprites map[string]*ebiten.Image

	// Effect assets
	particleSprites map[string]*ebiten.Image
}

// NewAssetManager creates and initializes all game assets
func NewAssetManager() *AssetManager {
	am := &AssetManager{
		creatureSprites: make(map[string]*ebiten.Image),
		foodSprites:     make(map[string]*ebiten.Image),
		toySprites:      make(map[string]*ebiten.Image),
		plantSprites:    make(map[string]*ebiten.Image),
		uiSprites:       make(map[string]*ebiten.Image),
		particleSprites: make(map[string]*ebiten.Image),
	}

	am.generateAllAssets()
	return am
}

// generateAllAssets creates all programmatic assets
func (am *AssetManager) generateAllAssets() {
	am.generateCreatureAssets()
	am.generateFoodAssets()
	am.generateToyAssets()
	am.generatePlantAssets()
	am.generateUIAssets()
	am.generateParticleAssets()
}

// generateCreatureAssets creates all creature-related sprites
func (am *AssetManager) generateCreatureAssets() {
	// Creature body (oval shape)
	am.creatureSprites["body"] = am.createOval(40, 50, color.RGBA{100, 200, 100, 255})

	// Creature head (circle)
	am.creatureSprites["head"] = am.createCircle(30, color.RGBA{100, 200, 100, 255})

	// Eyes
	am.creatureSprites["eye_white"] = am.createCircle(8, color.White)
	am.creatureSprites["pupil"] = am.createCircle(4, color.Black)

	// Arms and legs
	am.creatureSprites["arm"] = am.createOval(15, 8, color.RGBA{100, 200, 100, 255})
	am.creatureSprites["leg"] = am.createOval(10, 15, color.RGBA{100, 200, 100, 255})

	// Tail
	am.creatureSprites["tail"] = am.createTriangle(10, 15, color.RGBA{100, 200, 100, 255})

	// Expression elements
	am.creatureSprites["smile"] = am.createArc(20, math.Pi*0.2, math.Pi*0.8, color.Black)
	am.creatureSprites["frown"] = am.createArc(20, math.Pi*1.2, math.Pi*1.8, color.Black)
	am.creatureSprites["tear"] = am.createTearDrop(6, color.RGBA{100, 200, 255, 200})

	// Generate color variations for different Norn types
	colors := map[string]color.RGBA{
		"forest":   {34, 139, 34, 255},
		"desert":   {210, 105, 30, 255},
		"ocean":    {70, 130, 180, 255},
		"mountain": {147, 112, 219, 255},
	}

	for name, col := range colors {
		am.creatureSprites["body_"+name] = am.createOval(40, 50, col)
		am.creatureSprites["head_"+name] = am.createCircle(30, col)
		am.creatureSprites["arm_"+name] = am.createOval(15, 8, col)
		am.creatureSprites["leg_"+name] = am.createOval(10, 15, col)
	}
}

// generateFoodAssets creates all food sprites
func (am *AssetManager) generateFoodAssets() {
	// Apple
	apple := am.createCircle(20, color.RGBA{255, 0, 0, 255})
	am.addStemToApple(apple)
	am.foodSprites["apple"] = apple

	// Carrot
	carrot := am.createCarrot(10, 25)
	am.foodSprites["carrot"] = carrot

	// Honey
	am.foodSprites["honey"] = am.createHexagon(15, color.RGBA{255, 215, 0, 255})

	// Seed
	am.foodSprites["seed"] = am.createOval(8, 6, color.RGBA{139, 69, 19, 255})

	// Berry cluster
	am.foodSprites["berry"] = am.createBerryCluster(12, color.RGBA{128, 0, 128, 255})
}

// generateToyAssets creates all toy sprites
func (am *AssetManager) generateToyAssets() {
	// Ball with stripes
	ball := am.createStripedBall(25)
	am.toySprites["ball"] = ball

	// Music box
	musicBox := am.createMusicBox(40, 30)
	am.toySprites["musicbox"] = musicBox

	// Puzzle
	am.toySprites["puzzle"] = am.createPuzzle(30)

	// Mirror
	am.toySprites["mirror"] = am.createMirror(25)

	// Learning computer
	am.toySprites["computer"] = am.createComputer(40, 30)

	// Bed
	am.toySprites["bed"] = am.createBed(60, 30)
}

// generatePlantAssets creates all plant sprites
func (am *AssetManager) generatePlantAssets() {
	// Tree components
	am.plantSprites["tree_trunk"] = am.createRectangle(20, 40, color.RGBA{139, 69, 19, 255})
	am.plantSprites["tree_leaves"] = am.createTreeLeaves(60, color.RGBA{34, 139, 34, 255})

	// Flower components
	am.plantSprites["flower_stem"] = am.createRectangle(2, 20, color.RGBA{0, 128, 0, 255})
	am.plantSprites["flower_petals"] = am.createFlower(20)

	// Grass
	am.plantSprites["grass"] = am.createGrass(20, 10)

	// Bush
	am.plantSprites["bush"] = am.createBush(40, 30, color.RGBA{0, 100, 0, 255})

	// Different growth stages
	for i := 0; i < 5; i++ {
		scale := 0.2 + float64(i)*0.2
		am.plantSprites[fmt.Sprintf("tree_stage_%d", i)] = am.createScaledTree(scale)
	}
}

// generateUIAssets creates all UI elements
func (am *AssetManager) generateUIAssets() {
	// Health bar components
	am.uiSprites["bar_bg"] = am.createRectangle(60, 8, color.RGBA{50, 50, 50, 255})
	am.uiSprites["bar_health"] = am.createGradientBar(60, 8, color.RGBA{255, 0, 0, 255}, color.RGBA{0, 255, 0, 255})
	am.uiSprites["bar_hunger"] = am.createGradientBar(60, 8, color.RGBA{255, 165, 0, 255}, color.RGBA{255, 255, 0, 255})
	am.uiSprites["bar_energy"] = am.createGradientBar(60, 8, color.RGBA{100, 100, 255, 255}, color.RGBA{200, 200, 255, 255})

	// Emotion icons
	am.uiSprites["emotion_happy"] = am.createEmotionIcon("happy", color.RGBA{255, 255, 0, 255})
	am.uiSprites["emotion_sad"] = am.createEmotionIcon("sad", color.RGBA{100, 100, 255, 255})
	am.uiSprites["emotion_angry"] = am.createEmotionIcon("angry", color.RGBA{255, 0, 0, 255})
	am.uiSprites["emotion_scared"] = am.createEmotionIcon("scared", color.RGBA{200, 0, 200, 255})
	am.uiSprites["emotion_curious"] = am.createEmotionIcon("curious", color.RGBA{0, 255, 0, 255})

	// Speech bubble
	am.uiSprites["speech_bubble"] = am.createSpeechBubble(80, 30)

	// Selection indicator
	am.uiSprites["selection"] = am.createSelectionRing(40)

	// Menu elements
	am.uiSprites["menu_bg"] = am.createRoundedRect(300, 200, 10, color.RGBA{0, 0, 0, 200})
	am.uiSprites["button"] = am.createButton(150, 40)
	am.uiSprites["button_hover"] = am.createButton(150, 40, color.RGBA{255, 255, 100, 255})
}

// generateParticleAssets creates all particle effects
func (am *AssetManager) generateParticleAssets() {
	// Star particle
	am.particleSprites["star"] = am.createStar(10, color.RGBA{255, 255, 0, 255})

	// Heart particle
	am.particleSprites["heart"] = am.createHeart(12, color.RGBA{255, 192, 203, 255})

	// Musical note
	am.particleSprites["note"] = am.createMusicNote(10, color.Black)

	// Z for sleeping
	am.particleSprites["z"] = am.createZ(15, color.RGBA{173, 216, 230, 255})

	// Exclamation mark
	am.particleSprites["exclamation"] = am.createExclamation(15, color.RGBA{255, 0, 0, 255})

	// Light bulb for learning
	am.particleSprites["lightbulb"] = am.createLightBulb(15, color.RGBA{255, 255, 0, 255})

	// Food particles
	am.particleSprites["food_particle"] = am.createSquare(5, color.RGBA{200, 100, 50, 255})
}

// Asset creation helper functions

func (am *AssetManager) createCircle(radius int, c color.Color) *ebiten.Image {
	img := ebiten.NewImage(radius*2, radius*2)
	vector.DrawFilledCircle(img, float32(radius), float32(radius), float32(radius), c, true)
	return img
}

func (am *AssetManager) createOval(width, height int, c color.Color) *ebiten.Image {
	img := ebiten.NewImage(width, height)
	// Draw oval using bezier curves or multiple circles
	cx := float32(width) / 2
	cy := float32(height) / 2

	// Approximate oval with ellipse
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			dx := float32(x) - cx
			dy := float32(y) - cy
			if (dx*dx)/(cx*cx)+(dy*dy)/(cy*cy) <= 1 {
				img.Set(x, y, c)
			}
		}
	}
	return img
}

func (am *AssetManager) createRectangle(width, height int, c color.Color) *ebiten.Image {
	img := ebiten.NewImage(width, height)
	vector.DrawFilledRect(img, 0, 0, float32(width), float32(height), c, true)
	return img
}

func (am *AssetManager) createTriangle(width, height int, c color.Color) *ebiten.Image {
	img := ebiten.NewImage(width, height)

	// Draw filled triangle
	for y := 0; y < height; y++ {
		lineWidth := int(float64(width) * float64(height-y) / float64(height))
		startX := (width - lineWidth) / 2
		for x := startX; x < startX+lineWidth; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}

func (am *AssetManager) createHexagon(size int, c color.Color) *ebiten.Image {
	img := ebiten.NewImage(size*2, size*2)
	center := float32(size)

	// Draw hexagon using 6 triangles from center
	for i := 0; i < 6; i++ {
		angle1 := float32(i) * math.Pi / 3
		angle2 := float32(i+1) * math.Pi / 3

		// Calculate vertices
		_ = center + float32(size)*float32(math.Cos(float64(angle1)))
		_ = center + float32(size)*float32(math.Sin(float64(angle1)))
		_ = center + float32(size)*float32(math.Cos(float64(angle2)))
		_ = center + float32(size)*float32(math.Sin(float64(angle2)))

		// Fill triangle from center to edge
		vector.DrawFilledRect(img, center, center, 1, 1, c, true) // Simplified
	}

	return img
}

func (am *AssetManager) createArc(radius int, startAngle, endAngle float64, c color.Color) *ebiten.Image {
	img := ebiten.NewImage(radius*2, radius*2)
	center := float32(radius)

	steps := 20
	angleStep := (endAngle - startAngle) / float64(steps)

	for i := 0; i < steps; i++ {
		angle := startAngle + float64(i)*angleStep
		x := center + float32(float64(radius)*0.8*math.Cos(angle))
		y := center + float32(float64(radius)*0.8*math.Sin(angle))
		vector.DrawFilledCircle(img, x, y, 2, c, true)
	}

	return img
}

// createGradientBar creates a health/energy bar with gradient
func (am *AssetManager) createGradientBar(width, height int, startColor, endColor color.RGBA) *ebiten.Image {
	img := ebiten.NewImage(width, height)

	// Create gradient effect
	for x := 0; x < width; x++ {
		t := float64(x) / float64(width)
		r := uint8(float64(startColor.R)*(1-t) + float64(endColor.R)*t)
		g := uint8(float64(startColor.G)*(1-t) + float64(endColor.G)*t)
		b := uint8(float64(startColor.B)*(1-t) + float64(endColor.B)*t)

		gradColor := color.RGBA{r, g, b, 255}
		vector.DrawFilledRect(img, float32(x), 0, 1, float32(height), gradColor, true)
	}

	return img
}

func (am *AssetManager) createTearDrop(size int, c color.Color) *ebiten.Image {
	img := ebiten.NewImage(size, size*2)

	// Top circle part
	vector.DrawFilledCircle(img, float32(size/2), float32(size/2), float32(size/2), c, true)

	// Bottom triangle part
	for y := size / 2; y < size*2; y++ {
		width := int(float64(size) * float64(size*2-y) / float64(size*3/2))
		startX := (size - width) / 2
		for x := startX; x < startX+width; x++ {
			img.Set(x, y, c)
		}
	}

	return img
}

func (am *AssetManager) addStemToApple(apple *ebiten.Image) {
	// Add green stem on top
	bounds := apple.Bounds()
	cx := bounds.Dx() / 2

	for y := 0; y < 5; y++ {
		for x := cx - 2; x < cx+2; x++ {
			apple.Set(x, y, color.RGBA{0, 128, 0, 255})
		}
	}
}

func (am *AssetManager) createCarrot(width, height int) *ebiten.Image {
	img := ebiten.NewImage(width, height+10)

	// Orange carrot body (triangle)
	orange := color.RGBA{255, 165, 0, 255}
	for y := 5; y < height+5; y++ {
		lineWidth := int(float64(width) * float64(height+5-y) / float64(height))
		startX := (width - lineWidth) / 2
		for x := startX; x < startX+lineWidth; x++ {
			img.Set(x, y, orange)
		}
	}

	// Green top
	green := color.RGBA{0, 255, 0, 255}
	vector.DrawFilledCircle(img, float32(width/2), 5, 5, green, true)

	return img
}

func (am *AssetManager) createBerryCluster(size int, c color.Color) *ebiten.Image {
	img := ebiten.NewImage(size, size)

	// Draw multiple small circles
	positions := []struct{ x, y float32 }{
		{float32(size) / 2, float32(size) / 2},
		{float32(size) / 3, float32(size) / 3},
		{float32(size) * 2 / 3, float32(size) / 3},
		{float32(size) / 2, float32(size) * 2 / 3},
	}

	for _, pos := range positions {
		vector.DrawFilledCircle(img, pos.x, pos.y, float32(size)/4, c, true)
	}

	return img
}

func (am *AssetManager) createStripedBall(radius int) *ebiten.Image {
	img := ebiten.NewImage(radius*2, radius*2)

	// Base ball
	vector.DrawFilledCircle(img, float32(radius), float32(radius), float32(radius), color.RGBA{255, 0, 0, 255}, true)

	// White stripes
	stripeWidth := radius / 3
	for i := 0; i < 3; i++ {
		x := float32(i*stripeWidth*2 - radius/2)
		vector.DrawFilledRect(img, float32(radius)+x, 0, float32(stripeWidth), float32(radius*2), color.White, true)
	}

	return img
}

func (am *AssetManager) createMusicBox(width, height int) *ebiten.Image {
	img := ebiten.NewImage(width, height)

	// Wood box
	wood := color.RGBA{139, 69, 19, 255}
	vector.DrawFilledRect(img, 0, 0, float32(width), float32(height), wood, true)

	// Gold details
	gold := color.RGBA{255, 215, 0, 255}
	vector.DrawFilledRect(img, float32(width-10), 0, 5, 10, gold, true)              // Handle
	vector.StrokeRect(img, 2, 2, float32(width-4), float32(height-4), 2, gold, true) // Border

	return img
}

func (am *AssetManager) createPuzzle(size int) *ebiten.Image {
	img := ebiten.NewImage(size, size)

	// Create 4 puzzle pieces in different colors
	colors := []color.Color{
		color.RGBA{255, 100, 100, 255},
		color.RGBA{100, 255, 100, 255},
		color.RGBA{100, 100, 255, 255},
		color.RGBA{255, 255, 100, 255},
	}

	halfSize := size / 2
	for i := 0; i < 4; i++ {
		x := float32((i % 2) * halfSize)
		y := float32((i / 2) * halfSize)
		vector.DrawFilledRect(img, x, y, float32(halfSize), float32(halfSize), colors[i], true)
	}

	// Draw puzzle lines
	vector.StrokeLine(img, float32(halfSize), 0, float32(halfSize), float32(size), 2, color.Black, true)
	vector.StrokeLine(img, 0, float32(halfSize), float32(size), float32(halfSize), 2, color.Black, true)

	return img
}

func (am *AssetManager) createMirror(size int) *ebiten.Image {
	img := ebiten.NewImage(size, size)

	// Mirror frame
	frame := color.RGBA{139, 69, 19, 255}
	vector.StrokeCircle(img, float32(size/2), float32(size/2), float32(size/2), 3, frame, true)

	// Reflective surface
	silver := color.RGBA{192, 192, 192, 255}
	vector.DrawFilledCircle(img, float32(size/2), float32(size/2), float32(size/2-3), silver, true)

	// Shine effect
	shine := color.RGBA{255, 255, 255, 100}
	vector.DrawFilledCircle(img, float32(size/3), float32(size/3), float32(size/6), shine, true)

	return img
}

func (am *AssetManager) createComputer(width, height int) *ebiten.Image {
	img := ebiten.NewImage(width, height)

	// Computer body
	gray := color.RGBA{128, 128, 128, 255}
	vector.DrawFilledRect(img, 0, 0, float32(width), float32(height), gray, true)

	// Screen
	blue := color.RGBA{0, 0, 255, 255}
	screenMargin := float32(4)
	vector.DrawFilledRect(img, screenMargin, screenMargin, float32(width)-screenMargin*2, float32(height)-screenMargin*2, blue, true)

	// Blinking cursor
	vector.DrawFilledRect(img, 10, 10, 2, 10, color.White, true)

	return img
}

func (am *AssetManager) createBed(width, height int) *ebiten.Image {
	img := ebiten.NewImage(width, height)

	// Bed frame
	brown := color.RGBA{139, 69, 19, 255}
	vector.DrawFilledRect(img, 0, float32(height/2), float32(width), float32(height/2), brown, true)

	// Mattress
	blue := color.RGBA{65, 105, 225, 255}
	vector.DrawFilledRect(img, 2, 2, float32(width-4), float32(height/2), blue, true)

	// Pillow
	white := color.RGBA{255, 255, 255, 255}
	vector.DrawFilledRect(img, 5, 5, 20, 10, white, true)

	return img
}

func (am *AssetManager) createTreeLeaves(size int, c color.Color) *ebiten.Image {
	img := ebiten.NewImage(size, size)

	// Create bushy leaves with multiple circles
	positions := []struct {
		x, y, r float32
	}{
		{float32(size) / 2, float32(size) / 2, float32(size) / 3},
		{float32(size) / 3, float32(size) / 2, float32(size) / 4},
		{float32(size) * 2 / 3, float32(size) / 2, float32(size) / 4},
		{float32(size) / 2, float32(size) / 3, float32(size) / 4},
		{float32(size) / 2, float32(size) * 2 / 3, float32(size) / 4},
	}

	for _, pos := range positions {
		vector.DrawFilledCircle(img, pos.x, pos.y, pos.r, c, true)
	}

	return img
}

func (am *AssetManager) createFlower(size int) *ebiten.Image {
	img := ebiten.NewImage(size, size)
	center := float32(size) / 2

	// Draw petals
	petalColor := color.RGBA{255, 192, 203, 255}
	petalRadius := float32(size) / 4

	for i := 0; i < 5; i++ {
		angle := float32(i) * 2 * math.Pi / 5
		x := center + petalRadius*float32(math.Cos(float64(angle)))
		y := center + petalRadius*float32(math.Sin(float64(angle)))
		vector.DrawFilledCircle(img, x, y, petalRadius, petalColor, true)
	}

	// Draw center
	yellow := color.RGBA{255, 255, 0, 255}
	vector.DrawFilledCircle(img, center, center, float32(size)/6, yellow, true)

	return img
}

func (am *AssetManager) createGrass(width, height int) *ebiten.Image {
	img := ebiten.NewImage(width, height)

	green := color.RGBA{144, 238, 144, 255}
	bladeWidth := 2

	// Draw multiple grass blades
	for i := 0; i < width; i += 4 {
		bladeHeight := height - rand.Intn(height/3)
		vector.DrawFilledRect(img, float32(i), float32(height-bladeHeight), float32(bladeWidth), float32(bladeHeight), green, true)
	}

	return img
}

func (am *AssetManager) createBush(width, height int, c color.Color) *ebiten.Image {
	img := ebiten.NewImage(width, height)

	// Create bushy shape with overlapping circles
	numCircles := 5
	for i := 0; i < numCircles; i++ {
		x := float32(rand.Intn(width))
		y := float32(rand.Intn(height))
		r := float32(width) / 4
		vector.DrawFilledCircle(img, x, y, r, c, true)
	}

	return img
}

func (am *AssetManager) createScaledTree(scale float64) *ebiten.Image {
	baseSize := 60
	size := int(float64(baseSize) * scale)
	img := ebiten.NewImage(size, size)

	// Trunk
	trunkWidth := int(20 * scale)
	trunkHeight := int(40 * scale)
	trunkX := (size - trunkWidth) / 2
	trunkY := size - trunkHeight

	brown := color.RGBA{139, 69, 19, 255}
	vector.DrawFilledRect(img, float32(trunkX), float32(trunkY), float32(trunkWidth), float32(trunkHeight), brown, true)

	// Leaves
	green := color.RGBA{34, 139, 34, 255}
	leafRadius := float32(size) / 3
	vector.DrawFilledCircle(img, float32(size/2), float32(size/2), leafRadius, green, true)

	return img
}

func (am *AssetManager) createEmotionIcon(emotion string, c color.Color) *ebiten.Image {
	size := 20
	img := ebiten.NewImage(size, size)

	// Base circle
	vector.DrawFilledCircle(img, float32(size/2), float32(size/2), float32(size/2), c, true)

	// Add expression
	switch emotion {
	case "happy":
		// Smiley face
		vector.DrawFilledCircle(img, 7, 7, 2, color.Black, true)
		vector.DrawFilledCircle(img, 13, 7, 2, color.Black, true)
		// Smile would be an arc
	case "sad":
		// Sad face
		vector.DrawFilledCircle(img, 7, 7, 2, color.Black, true)
		vector.DrawFilledCircle(img, 13, 7, 2, color.Black, true)
		// Frown would be an inverted arc
	case "angry":
		// Angry eyebrows
		vector.StrokeLine(img, 5, 5, 8, 7, 2, color.Black, true)
		vector.StrokeLine(img, 15, 5, 12, 7, 2, color.Black, true)
	case "scared":
		// Wide eyes
		vector.DrawFilledCircle(img, 7, 8, 3, color.White, true)
		vector.DrawFilledCircle(img, 13, 8, 3, color.White, true)
		vector.DrawFilledCircle(img, 7, 8, 1, color.Black, true)
		vector.DrawFilledCircle(img, 13, 8, 1, color.Black, true)
	case "curious":
		// Question mark overlay
		vector.DrawFilledCircle(img, 10, 8, 1, color.White, true)
		vector.DrawFilledCircle(img, 10, 14, 1, color.White, true)
	}

	return img
}

func (am *AssetManager) createSpeechBubble(width, height int) *ebiten.Image {
	img := ebiten.NewImage(width, height+10)

	// Main bubble
	vector.DrawFilledRect(img, 0, 0, float32(width), float32(height), color.White, true)
	vector.StrokeRect(img, 0, 0, float32(width), float32(height), 2, color.Black, true)

	// Tail
	tailSize := 10
	for y := 0; y < tailSize; y++ {
		lineWidth := tailSize - y
		startX := width/2 - lineWidth/2
		for x := 0; x < lineWidth; x++ {
			img.Set(startX+x, height+y, color.White)
		}
	}

	return img
}

func (am *AssetManager) createSelectionRing(radius int) *ebiten.Image {
	img := ebiten.NewImage(radius*2, radius*2)

	// Create dashed circle effect
	yellow := color.RGBA{255, 255, 0, 200}
	segments := 16
	segmentArc := 2 * math.Pi / float64(segments)

	for i := 0; i < segments; i += 2 {
		startAngle := float64(i) * segmentArc
		endAngle := float64(i+1) * segmentArc

		for a := startAngle; a < endAngle; a += 0.1 {
			x := float32(radius) + float32(radius)*float32(math.Cos(a))
			y := float32(radius) + float32(radius)*float32(math.Sin(a))
			vector.DrawFilledCircle(img, x, y, 2, yellow, true)
		}
	}

	return img
}

func (am *AssetManager) createRoundedRect(width, height, radius int, c color.Color) *ebiten.Image {
	img := ebiten.NewImage(width, height)

	// Main rectangle
	vector.DrawFilledRect(img, float32(radius), 0, float32(width-radius*2), float32(height), c, true)
	vector.DrawFilledRect(img, 0, float32(radius), float32(width), float32(height-radius*2), c, true)

	// Corners
	vector.DrawFilledCircle(img, float32(radius), float32(radius), float32(radius), c, true)
	vector.DrawFilledCircle(img, float32(width-radius), float32(radius), float32(radius), c, true)
	vector.DrawFilledCircle(img, float32(radius), float32(height-radius), float32(radius), c, true)
	vector.DrawFilledCircle(img, float32(width-radius), float32(height-radius), float32(radius), c, true)

	return img
}

func (am *AssetManager) createButton(width, height int, colors ...color.Color) *ebiten.Image {
	img := ebiten.NewImage(width, height)

	bgColor := color.RGBA{100, 100, 100, 255}
	if len(colors) > 0 {
		bgColor = colors[0].(color.RGBA)
	}

	// Button background
	am.drawRoundedRectToImage(img, 0, 0, width, height, 5, bgColor)

	// Button border
	borderColor := color.RGBA{200, 200, 200, 255}
	am.drawRoundedRectBorderToImage(img, 0, 0, width, height, 5, 2, borderColor)

	return img
}

func (am *AssetManager) createStar(size int, c color.Color) *ebiten.Image {
	img := ebiten.NewImage(size, size)
	center := float32(size) / 2

	// Draw 5-pointed star
	outerRadius := float32(size) / 2
	innerRadius := outerRadius / 2

	for i := 0; i < 10; i++ {
		angle := float32(i) * math.Pi / 5
		radius := outerRadius
		if i%2 == 1 {
			radius = innerRadius
		}

		x := center + radius*float32(math.Cos(float64(angle-math.Pi/2)))
		y := center + radius*float32(math.Sin(float64(angle-math.Pi/2)))

		if i == 0 {
			// Start of star
		} else {
			// Connect points (simplified - just dots)
			vector.DrawFilledCircle(img, x, y, 1, c, true)
		}
	}

	// Fill center
	vector.DrawFilledCircle(img, center, center, innerRadius, c, true)

	return img
}

func (am *AssetManager) createHeart(size int, c color.Color) *ebiten.Image {
	img := ebiten.NewImage(size, size)

	// Two circles for top
	quarter := float32(size) / 4
	vector.DrawFilledCircle(img, quarter, quarter, quarter, c, true)
	vector.DrawFilledCircle(img, 3*quarter, quarter, quarter, c, true)

	// Triangle for bottom
	for y := size / 4; y < size; y++ {
		width := int(float64(size) * float64(size-y) / float64(size*3/4))
		startX := (size - width) / 2
		for x := startX; x < startX+width; x++ {
			img.Set(x, y, c)
		}
	}

	return img
}

func (am *AssetManager) createMusicNote(size int, c color.Color) *ebiten.Image {
	img := ebiten.NewImage(size, size*2)

	// Note head
	vector.DrawFilledCircle(img, float32(size/2), float32(size*3/2), float32(size/3), c, true)

	// Stem
	vector.DrawFilledRect(img, float32(size/2+size/3-1), float32(size/2), 2, float32(size), c, true)

	// Flag
	vector.DrawFilledRect(img, float32(size/2+size/3-1), float32(size/2), float32(size/3), float32(size/4), c, true)

	return img
}

func (am *AssetManager) createZ(size int, c color.Color) *ebiten.Image {
	img := ebiten.NewImage(size, size)

	// Top line
	vector.DrawFilledRect(img, 0, 0, float32(size), 3, c, true)

	// Diagonal
	for i := 0; i < size; i++ {
		x := size - i - 1
		y := i
		vector.DrawFilledRect(img, float32(x), float32(y), 3, 3, c, true)
	}

	// Bottom line
	vector.DrawFilledRect(img, 0, float32(size-3), float32(size), 3, c, true)

	return img
}

func (am *AssetManager) createExclamation(size int, c color.Color) *ebiten.Image {
	img := ebiten.NewImage(size/3, size)

	// Main line
	vector.DrawFilledRect(img, float32(size/6-2), 0, 4, float32(size*2/3), c, true)

	// Dot
	vector.DrawFilledCircle(img, float32(size/6), float32(size*5/6), 3, c, true)

	return img
}

func (am *AssetManager) createLightBulb(size int, c color.Color) *ebiten.Image {
	img := ebiten.NewImage(size, size)

	// Bulb shape
	vector.DrawFilledCircle(img, float32(size/2), float32(size/3), float32(size/3), c, true)

	// Base
	gray := color.RGBA{128, 128, 128, 255}
	vector.DrawFilledRect(img, float32(size/3), float32(size*2/3), float32(size/3), float32(size/4), gray, true)

	// Light rays
	for i := 0; i < 8; i++ {
		angle := float32(i) * math.Pi / 4
		x1 := float32(size/2) + float32(size/3)*float32(math.Cos(float64(angle)))
		y1 := float32(size/3) + float32(size/3)*float32(math.Sin(float64(angle)))
		x2 := float32(size/2) + float32(size/2)*float32(math.Cos(float64(angle)))
		y2 := float32(size/3) + float32(size/2)*float32(math.Sin(float64(angle)))
		vector.StrokeLine(img, x1, y1, x2, y2, 1, c, true)
	}

	return img
}

func (am *AssetManager) createSquare(size int, c color.Color) *ebiten.Image {
	img := ebiten.NewImage(size, size)
	vector.DrawFilledRect(img, 0, 0, float32(size), float32(size), c, true)
	return img
}

// Helper methods for drawing

func (am *AssetManager) drawRoundedRectToImage(img *ebiten.Image, x, y, width, height, radius int, c color.Color) {
	// Main rectangle areas
	vector.DrawFilledRect(img, float32(x+radius), float32(y), float32(width-radius*2), float32(height), c, true)
	vector.DrawFilledRect(img, float32(x), float32(y+radius), float32(width), float32(height-radius*2), c, true)

	// Corners
	vector.DrawFilledCircle(img, float32(x+radius), float32(y+radius), float32(radius), c, true)
	vector.DrawFilledCircle(img, float32(x+width-radius), float32(y+radius), float32(radius), c, true)
	vector.DrawFilledCircle(img, float32(x+radius), float32(y+height-radius), float32(radius), c, true)
	vector.DrawFilledCircle(img, float32(x+width-radius), float32(y+height-radius), float32(radius), c, true)
}

func (am *AssetManager) drawRoundedRectBorderToImage(img *ebiten.Image, x, y, width, height, radius, thickness int, c color.Color) {
	// Top and bottom
	vector.StrokeLine(img, float32(x+radius), float32(y), float32(x+width-radius), float32(y), float32(thickness), c, true)
	vector.StrokeLine(img, float32(x+radius), float32(y+height), float32(x+width-radius), float32(y+height), float32(thickness), c, true)

	// Left and right
	vector.StrokeLine(img, float32(x), float32(y+radius), float32(x), float32(y+height-radius), float32(thickness), c, true)
	vector.StrokeLine(img, float32(x+width), float32(y+radius), float32(x+width), float32(y+height-radius), float32(thickness), c, true)

	// Corners (simplified)
	vector.StrokeCircle(img, float32(x+radius), float32(y+radius), float32(radius), float32(thickness), c, true)
	vector.StrokeCircle(img, float32(x+width-radius), float32(y+radius), float32(radius), float32(thickness), c, true)
	vector.StrokeCircle(img, float32(x+radius), float32(y+height-radius), float32(radius), float32(thickness), c, true)
	vector.StrokeCircle(img, float32(x+width-radius), float32(y+height-radius), float32(radius), float32(thickness), c, true)
}

// Getter methods

func (am *AssetManager) GetCreatureSprite(name string) *ebiten.Image {
	return am.creatureSprites[name]
}

func (am *AssetManager) GetFoodSprite(name string) *ebiten.Image {
	return am.foodSprites[name]
}

func (am *AssetManager) GetToySprite(name string) *ebiten.Image {
	return am.toySprites[name]
}

func (am *AssetManager) GetPlantSprite(name string) *ebiten.Image {
	return am.plantSprites[name]
}

func (am *AssetManager) GetUISprite(name string) *ebiten.Image {
	return am.uiSprites[name]
}

func (am *AssetManager) GetParticleSprite(name string) *ebiten.Image {
	return am.particleSprites[name]
}
