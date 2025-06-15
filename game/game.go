package game

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/olivierh59500/creatures-clone/creature"
	"github.com/olivierh59500/creatures-clone/objects"
	"github.com/olivierh59500/creatures-clone/renderer"
	"github.com/olivierh59500/creatures-clone/ui"
	"github.com/olivierh59500/creatures-clone/utils"
)

// GameState represents the current state of the game
type GameState int

const (
	StateMenu GameState = iota
	StatePlaying
	StatePaused
)

// Game represents the main game structure
type Game struct {
	// Core systems
	world    *World
	camera   *Camera
	renderer *renderer.Renderer

	// UI systems
	hud   *ui.HUD
	menu  *ui.Menu
	debug *ui.Debug

	// Game state
	state          GameState
	selectedNorn   *creature.Creature
	mouseX, mouseY int
	currentWord    string // Word being typed
	message        string // Feedback message
	messageTimer   float64

	// Time tracking
	ticks uint64

	// Configuration
	config *utils.Config
}

// NewGame creates a new game instance
func NewGame() *Game {
	config := utils.LoadConfig()

	g := &Game{
		world:    NewWorld(config.WorldWidth, config.WorldHeight),
		camera:   NewCamera(config.ScreenWidth, config.ScreenHeight),
		renderer: renderer.NewRenderer(),
		hud:      ui.NewHUD(),
		menu:     ui.NewMenu(),
		debug:    ui.NewDebug(),
		state:    StateMenu,
		config:   config,
	}

	// Initialize the world with starting creatures and objects
	g.initializeWorld()

	return g
}

// initializeWorld sets up the initial game world
func (g *Game) initializeWorld() {
	// Calculate ground level
	groundY := float64(g.config.WorldHeight) * 0.8

	// Create starting Norns in a nice line on the ground
	startX := float64(g.config.WorldWidth) / 4
	for i := 0; i < g.config.StartingNorns; i++ {
		x := startX + float64(i*150)
		y := groundY - 50 // Just above ground

		norn := creature.NewCreature(x, y, creature.CreatureTypeNorn)
		norn.Genetics.Randomize() // Random genetics for variety

		// Give them slightly different starting stats
		norn.Metabolism.Hunger = 30 + float64(i*10)
		norn.Metabolism.Energy = 70 + float64(i*5)

		// Give each a unique name for easy identification
		names := []string{"Albie", "Bella", "Charlie", "Daisy", "Eddie"}
		if i < len(names) {
			norn.Name = names[i]
		}

		g.world.AddCreature(norn)
	}

	// Create organized food areas
	// Food garden on the left
	for i := 0; i < 6; i++ {
		x := 100.0 + float64(i%3)*80
		y := groundY - 30 - float64(i/3)*60

		foods := []objects.FoodType{objects.FoodApple, objects.FoodCarrot, objects.FoodBerry}
		food := objects.NewFood(x, y, foods[i%len(foods)])
		g.world.AddObject(food)
	}

	// Honey stash on the right
	for i := 0; i < 3; i++ {
		x := float64(g.config.WorldWidth) - 200 + float64(i*50)
		y := groundY - 30

		honey := objects.NewFood(x, y, objects.FoodHoney)
		g.world.AddObject(honey)
	}

	// Create a small forest area in the middle
	forestCenterX := float64(g.config.WorldWidth) / 2
	for i := 0; i < 4; i++ {
		x := forestCenterX + float64((i-2)*120)
		y := groundY

		tree := objects.NewPlant(x, y, objects.PlantTree)
		// Make some trees already grown
		if i%2 == 0 {
			tree.Age = 200
			tree.GrowthStage = objects.StageMature
			tree.Size = 1.0
		}
		g.world.AddObject(tree)
	}

	// Add some flowers around
	for i := 0; i < 8; i++ {
		x := utils.RandomFloat(100, float64(g.config.WorldWidth-100))
		y := groundY

		flower := objects.NewPlant(x, y, objects.PlantFlower)
		g.world.AddObject(flower)
	}

	// Place toys in accessible locations
	// Ball near the creatures
	ball := objects.NewToy(startX+100, groundY-30, objects.ToyBall)
	g.world.AddObject(ball)

	// Music box in the middle
	musicBox := objects.NewToy(forestCenterX, groundY-30, objects.ToyMusicBox)
	g.world.AddObject(musicBox)

	// Learning computer on a "table" (elevated position)
	computer := objects.NewToy(float64(g.config.WorldWidth)*0.75, groundY-60, objects.ToyComputer)
	g.world.AddObject(computer)

	// Create a cozy sleeping area with a bed
	bed := objects.NewToy(float64(g.config.WorldWidth)*0.85, groundY-20, objects.ToyBed)
	g.world.AddObject(bed)
}

// Update updates the game state
func (g *Game) Update() error {
	// Update mouse position
	g.mouseX, g.mouseY = ebiten.CursorPosition()

	// Handle state-specific updates
	switch g.state {
	case StateMenu:
		g.updateMenu()
	case StatePlaying:
		g.updatePlaying()
	case StatePaused:
		g.updatePaused()
	}

	return nil
}

// updateMenu handles menu state updates
func (g *Game) updateMenu() {
	action := g.menu.Update(g.mouseX, g.mouseY, inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft))

	switch action {
	case ui.MenuActionStart:
		g.state = StatePlaying
	case ui.MenuActionQuit:
		// In a real implementation, this would quit the game
		// For now, we'll just start the game
		g.state = StatePlaying
	}
}

// updatePlaying handles the main game state updates
func (g *Game) updatePlaying() {
	// Handle input
	g.handleInput()

	// Update camera
	g.camera.Update()

	// Update world
	g.world.Update()

	// Update HUD
	g.hud.Update(g.selectedNorn, g.world)

	// Update debug overlay if enabled
	if g.debug.IsEnabled() {
		g.debug.Update(g.world, g.camera, g.mouseX, g.mouseY)
	}

	// Increment tick counter
	g.ticks++
}

// updatePaused handles paused state updates
func (g *Game) updatePaused() {
	// Check for unpause
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.state = StatePlaying
	}
}

// handleInput processes user input
func (g *Game) handleInput() {
	// Camera movement
	moveSpeed := 5.0
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.camera.Move(-moveSpeed, 0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.camera.Move(moveSpeed, 0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.camera.Move(0, -moveSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.camera.Move(0, moveSpeed)
	}

	// Camera zoom
	_, scrollY := ebiten.Wheel()
	if scrollY != 0 {
		g.camera.Zoom(1 + scrollY*0.1)
	}

	// Pause/unpause
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if g.state == StatePlaying {
			g.state = StatePaused
		} else if g.state == StatePaused {
			g.state = StatePlaying
		}
	}

	// Toggle debug overlay
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		g.debug.Toggle()
	}

	// Escape to menu
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.state = StateMenu
	}

	// Mouse interactions
	worldX, worldY := g.camera.ScreenToWorld(float64(g.mouseX), float64(g.mouseY))

	// Left click - select creature or interact with object
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.selectedNorn = nil

		// Check creatures first
		for _, c := range g.world.GetCreatures() {
			if c.Contains(worldX, worldY) {
				g.selectedNorn = c
				break
			}
		}

		// If no creature selected, check objects
		if g.selectedNorn == nil {
			for _, obj := range g.world.GetObjects() {
				pos := obj.GetPosition()
				dist := utils.Distance(worldX, worldY, pos.X, pos.Y)
				if dist < 30 {
					// If we have a selected creature, make it interact with the object
					if g.selectedNorn != nil {
						// Guide creature to object
						g.selectedNorn.SetTarget(pos.X, pos.Y)
					}
				}
			}
		}
	}

	// Right click - place food or guide creature
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		if g.selectedNorn != nil {
			// Guide selected creature to location
			g.selectedNorn.SetTarget(worldX, worldY)
		} else {
			// Place food
			food := objects.NewFood(worldX, worldY, objects.FoodApple)
			g.world.AddObject(food)
		}
	}

	// Typing - teach words to selected creature
	if g.selectedNorn != nil {
		// Capture typed characters
		for _, r := range ebiten.AppendInputChars(nil) {
			if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' {
				g.currentWord += string(r)
			}
		}

		// On Enter, teach the word
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && g.currentWord != "" {
			// Find nearest object to associate with word
			nearestObj := g.findNearestObject(g.selectedNorn.X, g.selectedNorn.Y)
			if nearestObj != nil {
				g.selectedNorn.Language.TeachWord(g.currentWord, nearestObj.GetType())
				// Show feedback
				g.showMessage(fmt.Sprintf("Taught '%s' = %s", g.currentWord, nearestObj.GetType()))
			}
			g.currentWord = ""
		}
	}

	// B key - encourage breeding
	if inpututil.IsKeyJustPressed(ebiten.KeyB) && g.selectedNorn != nil {
		g.selectedNorn.EncourageBreeding()
	}
}

// Draw renders the game
func (g *Game) Draw(screen *ebiten.Image) {
	// Clear screen
	screen.Fill(color.RGBA{135, 206, 235, 255}) // Sky blue

	switch g.state {
	case StateMenu:
		g.menu.Draw(screen)
	case StatePlaying, StatePaused:
		g.drawGame(screen)

		if g.state == StatePaused {
			g.drawPausedOverlay(screen)
		}
	}

	// Always draw FPS in debug mode
	if g.debug.IsEnabled() {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()))
	}
}

// drawGame renders the main game view
func (g *Game) drawGame(screen *ebiten.Image) {
	// Create camera transform
	camTransform := g.camera.GetTransform()

	// Draw world background
	g.renderer.DrawWorldBackground(screen, g.world, camTransform)

	// Draw objects
	for _, obj := range g.world.GetObjects() {
		g.renderer.DrawObject(screen, obj, camTransform)
	}

	// Draw creatures
	for _, c := range g.world.GetCreatures() {
		isSelected := c == g.selectedNorn
		g.renderer.DrawCreature(screen, c, camTransform, isSelected)
	}

	// Update and draw particles
	g.renderer.UpdateParticles()
	g.renderer.DrawParticles(screen)

	// Draw UI elements
	g.hud.Draw(screen)

	// Draw creature info for selected creature
	if g.selectedNorn != nil {
		g.hud.DrawCreatureInfo(screen, g.selectedNorn)
	}

	if g.debug.IsEnabled() {
		g.debug.Draw(screen)
	}

	// Draw message if any
	if g.messageTimer > 0 && g.message != "" {
		msgX := screen.Bounds().Dx()/2 - len(g.message)*4
		msgY := screen.Bounds().Dy() - 100

		// Background for message
		bgWidth := float32(len(g.message)*8 + 20)
		bgHeight := float32(25)
		bgX := float32(msgX - 10)
		bgY := float32(msgY - 5)

		vector.DrawFilledRect(screen, bgX, bgY, bgWidth, bgHeight, color.RGBA{0, 0, 0, 200}, false)
		ebitenutil.DebugPrintAt(screen, g.message, msgX, msgY)
	}
}

// drawPausedOverlay draws the pause screen overlay
func (g *Game) drawPausedOverlay(screen *ebiten.Image) {
	// Semi-transparent overlay
	overlay := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())
	overlay.Fill(color.RGBA{0, 0, 0, 128})
	screen.DrawImage(overlay, nil)

	// Pause text
	text := "PAUSED"
	x := screen.Bounds().Dx()/2 - len(text)*4
	y := screen.Bounds().Dy() / 2
	ebitenutil.DebugPrintAt(screen, text, x, y)
	ebitenutil.DebugPrintAt(screen, "Press SPACE to continue", x-40, y+20)
}

// findNearestObject finds the nearest object to a position
func (g *Game) findNearestObject(x, y float64) objects.Object {
	var nearest objects.Object
	minDist := math.MaxFloat64

	for _, obj := range g.world.GetObjects() {
		pos := obj.GetPosition()
		dist := utils.Distance(x, y, pos.X, pos.Y)
		if dist < minDist {
			minDist = dist
			nearest = obj
		}
	}

	return nearest
}

// showMessage displays a temporary message
func (g *Game) showMessage(msg string) {
	g.message = msg
	g.messageTimer = 3.0 // Show for 3 seconds
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// Retourner la taille de l'écran définie dans la configuration
	return g.config.ScreenWidth, g.config.ScreenHeight
}
