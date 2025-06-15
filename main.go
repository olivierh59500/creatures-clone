package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/olivierh59500/creatures-clone/game"
)

func main() {
	// Create new game instance
	g := game.NewGame()

	// Set window properties
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Creatures Clone - Artificial Life Simulation")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// Run the game
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
